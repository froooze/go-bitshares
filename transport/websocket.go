package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

// WebsocketClient speaks to a BitShares node over websocket JSON-RPC.
type WebsocketClient struct {
	endpoint string

	mu            sync.Mutex
	conn          *websocket.Conn
	nextID        atomic.Uint64
	pending       map[uint64]chan RPCResponse
	closed        bool
	done          chan struct{}
	notifyHandler func(json.RawMessage)
	closeHandler  func(error)
}

func NewWebsocketClient(endpoint string) (*WebsocketClient, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("empty endpoint")
	}

	return &WebsocketClient{
		endpoint: endpoint,
		pending:  make(map[uint64]chan RPCResponse),
		done:     make(chan struct{}),
	}, nil
}

func (c *WebsocketClient) Connect(ctx context.Context) error {
	if err := requireContext(ctx); err != nil {
		return err
	}

	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return ErrShutdown
	}
	if c.conn != nil {
		c.mu.Unlock()
		return nil
	}
	c.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	cfg, err := websocket.NewConfig(c.endpoint, "http://localhost/")
	if err != nil {
		return err
	}
	cfg.Dialer = &net.Dialer{}

	conn, err := websocket.DialConfig(cfg)
	if err != nil {
		return err
	}

	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		_ = conn.Close()
		return ErrShutdown
	}
	if c.conn != nil {
		c.mu.Unlock()
		_ = conn.Close()
		return nil
	}
	c.conn = conn
	c.mu.Unlock()
	go c.readLoop()
	return nil
}

func (c *WebsocketClient) Close() error {
	return c.closeWithError(ErrShutdown)
}

// Connected reports whether the websocket transport currently has an open connection.
func (c *WebsocketClient) Connected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn != nil && !c.closed
}

// SetNotificationHandler registers a callback for unsolicited websocket messages.
func (c *WebsocketClient) SetNotificationHandler(handler func(json.RawMessage)) {
	c.mu.Lock()
	c.notifyHandler = handler
	c.mu.Unlock()
}

// SetCloseHandler registers a callback fired when the websocket closes.
func (c *WebsocketClient) SetCloseHandler(handler func(error)) {
	c.mu.Lock()
	c.closeHandler = handler
	c.mu.Unlock()
}

func (c *WebsocketClient) closeWithError(err error) error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil
	}
	c.closed = true
	conn := c.conn
	c.conn = nil
	pending := c.pending
	c.pending = make(map[uint64]chan RPCResponse)
	closeHandler := c.closeHandler
	close(c.done)
	c.mu.Unlock()

	c.failPending(pending, err)

	if closeHandler != nil {
		go closeHandler(err)
	}

	if conn == nil {
		return nil
	}

	return conn.Close()
}

func (c *WebsocketClient) disconnect(err error) error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil
	}
	conn := c.conn
	c.conn = nil
	pending := c.pending
	c.pending = make(map[uint64]chan RPCResponse)
	closeHandler := c.closeHandler
	c.mu.Unlock()

	c.failPending(pending, err)

	if closeHandler != nil {
		go closeHandler(err)
	}

	if conn == nil {
		return nil
	}

	return conn.Close()
}

func (c *WebsocketClient) failPending(pending map[uint64]chan RPCResponse, err error) {
	if len(pending) == 0 {
		return
	}

	rpcErr := &RPCError{Message: ErrNotConnected.Error()}
	if err != nil {
		rpcErr.Message = err.Error()
	}

	for _, ch := range pending {
		select {
		case ch <- RPCResponse{Error: rpcErr}:
		default:
		}
		close(ch)
	}
}

func (c *WebsocketClient) CallAPI(ctx context.Context, apiID int, method string, args []any, reply any) error {
	if err := requireContext(ctx); err != nil {
		return err
	}
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return ErrShutdown
	}
	if c.conn == nil {
		c.mu.Unlock()
		return ErrNotConnected
	}

	id := c.nextID.Add(1)
	ch := make(chan RPCResponse, 1)
	c.pending[id] = ch
	conn := c.conn

	req := RPCRequest{
		Method: "call",
		Params: []any{apiID, method, args},
		ID:     id,
	}

	c.mu.Unlock()

	if err := websocket.JSON.Send(conn, req); err != nil {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return err
	}

	select {
	case <-ctx.Done():
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return ctx.Err()
	case <-c.done:
		return ErrShutdown
	case resp := <-ch:
		if resp.Error != nil {
			return resp.Error
		}
		if reply == nil || resp.Result == nil {
			return nil
		}
		return json.Unmarshal(*resp.Result, reply)
	}
}

func (c *WebsocketClient) readLoop() {
	for {
		var raw string
		c.mu.Lock()
		conn := c.conn
		c.mu.Unlock()

		if conn == nil {
			return
		}

		if err := websocket.Message.Receive(conn, &raw); err != nil {
			_ = c.disconnect(err)
			return
		}

		var resp RPCResponse
		if err := json.Unmarshal([]byte(raw), &resp); err != nil {
			c.dispatchNotification(json.RawMessage(raw))
			continue
		}

		c.mu.Lock()
		ch := c.pending[resp.ID]
		if ch != nil {
			delete(c.pending, resp.ID)
		}
		notify := c.notifyHandler
		c.mu.Unlock()

		if ch != nil {
			ch <- resp
			continue
		}

		if notify != nil && resp.ID == 0 && resp.Result == nil && resp.Error == nil {
			notify(json.RawMessage(raw))
		}
	}
}

func (c *WebsocketClient) dispatchNotification(raw json.RawMessage) {
	c.mu.Lock()
	notify := c.notifyHandler
	c.mu.Unlock()
	if notify != nil {
		notify(raw)
	}
}

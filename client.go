package bitshares

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/froooze/go-bitshares/protocol"
	"github.com/froooze/go-bitshares/sign"
	"github.com/froooze/go-bitshares/transport"
)

// Client is a small facade over the available transports.
type Client struct {
	ws     *transport.WebsocketClient
	wsMgr  *transport.SubscriptionManager
	rpc    *transport.RPCClient
	Signer sign.Signer
}

// NewWebsocketClient creates a websocket-based BitShares client.
func NewWebsocketClient(endpoint string) (*Client, error) {
	ws, err := transport.NewWebsocketClient(endpoint)
	if err != nil {
		return nil, err
	}

	return &Client{ws: ws}, nil
}

// NewRPCClient creates a JSON-RPC client.
func NewRPCClient(endpoint string) *Client {
	return &Client{rpc: transport.NewRPCClient(endpoint)}
}

// NewManagedWebsocketClient creates a websocket client with reconnect and resubscribe support.
func NewManagedWebsocketClient(endpoint string, reconnectDelay time.Duration) (*Client, error) {
	ws, err := transport.NewWebsocketClient(endpoint)
	if err != nil {
		return nil, err
	}

	return &Client{
		ws:    ws,
		wsMgr: transport.NewSubscriptionManager(ws, reconnectDelay),
	}, nil
}

// Connect opens the websocket transport when one is configured.
func (c *Client) Connect(ctx context.Context) error {
	if err := requireContext(ctx); err != nil {
		return err
	}
	if c.wsMgr != nil {
		return c.wsMgr.Connect(ctx)
	}
	if c.ws == nil {
		return nil
	}

	return c.ws.Connect(ctx)
}

// Close releases transport resources.
func (c *Client) Close() error {
	var errs []error

	if c.wsMgr != nil {
		if err := c.wsMgr.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if c.ws != nil {
		if err := c.ws.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if c.rpc != nil {
		if err := c.rpc.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// SignTransaction delegates to the configured signer.
func (c *Client) SignTransaction(tx *protocol.Transaction) (*protocol.SignedTransaction, error) {
	if c.Signer == nil {
		return nil, fmt.Errorf("signer is not configured")
	}

	return c.Signer.Sign(tx)
}

// Call invokes a JSON-RPC method on the HTTP transport.
func (c *Client) Call(ctx context.Context, method string, args []any, reply any) error {
	if c.rpc == nil {
		return fmt.Errorf("rpc transport is not configured")
	}
	if err := requireContext(ctx); err != nil {
		return err
	}

	return c.rpc.Call(ctx, method, args, reply)
}

// CallAPI invokes a BitShares websocket API method.
func (c *Client) CallAPI(ctx context.Context, apiID int, method string, args []any, reply any) error {
	if c.wsMgr != nil {
		return c.wsMgr.CallAPI(ctx, apiID, method, args, reply)
	}
	if c.ws == nil {
		return fmt.Errorf("websocket transport is not configured")
	}
	if err := requireContext(ctx); err != nil {
		return err
	}

	return c.ws.CallAPI(ctx, apiID, method, args, reply)
}

// Subscribe registers a websocket resubscription callback on the managed websocket client.
func (c *Client) Subscribe(name string, fn transport.SubscriptionFunc) error {
	if c.wsMgr == nil {
		return fmt.Errorf("websocket subscription manager is not configured")
	}
	return c.wsMgr.Register(name, fn)
}

// Notifications exposes raw websocket notifications when the managed websocket client is used.
func (c *Client) Notifications() <-chan json.RawMessage {
	if c.wsMgr == nil {
		return nil
	}
	return c.wsMgr.Notifications()
}

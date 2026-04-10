package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// SubscriptionFunc re-registers a websocket subscription after connect or reconnect.
type SubscriptionFunc func(context.Context, *WebsocketClient) error

// SubscriptionManager keeps websocket subscriptions alive across reconnects.
type SubscriptionManager struct {
	client         *WebsocketClient
	reconnectDelay time.Duration

	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	started bool
	closed  bool

	subs map[string]SubscriptionFunc

	notifications chan json.RawMessage
	errors        chan error
	closeCh       chan error
}

// NewSubscriptionManager creates a manager for the given websocket client.
func NewSubscriptionManager(client *WebsocketClient, reconnectDelay time.Duration) *SubscriptionManager {
	if reconnectDelay <= 0 {
		reconnectDelay = 5 * time.Second
	}

	m := &SubscriptionManager{
		client:         client,
		reconnectDelay: reconnectDelay,
		subs:           make(map[string]SubscriptionFunc),
		notifications:  make(chan json.RawMessage, 128),
		errors:         make(chan error, 16),
		closeCh:        make(chan error, 1),
	}

	if client != nil {
		client.SetNotificationHandler(func(msg json.RawMessage) {
			if msg == nil {
				return
			}
			copyMsg := append(json.RawMessage(nil), msg...)
			select {
			case m.notifications <- copyMsg:
			default:
			}
		})

		client.SetCloseHandler(func(err error) {
			select {
			case m.closeCh <- err:
			default:
			}
		})
	}

	return m
}

// Notifications returns raw websocket notifications.
func (m *SubscriptionManager) Notifications() <-chan json.RawMessage {
	return m.notifications
}

// Errors returns reconnect and subscription errors.
func (m *SubscriptionManager) Errors() <-chan error {
	return m.errors
}

// Register adds a websocket resubscription callback.
func (m *SubscriptionManager) Register(name string, fn SubscriptionFunc) error {
	if name == "" {
		return fmt.Errorf("empty subscription name")
	}
	if fn == nil {
		return fmt.Errorf("nil subscription callback")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return ErrShutdown
	}
	if _, exists := m.subs[name]; exists {
		return fmt.Errorf("subscription %q already registered", name)
	}
	m.subs[name] = fn

	if m.started && m.client != nil && m.client.Connected() {
		ctx := m.ctx
		go func() {
			if err := m.invokeSubscription(ctx, name, fn); err != nil {
				m.sendError(err)
			}
		}()
	}

	return nil
}

// Start begins the reconnect loop and performs the initial connect.
func (m *SubscriptionManager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return ErrShutdown
	}
	if m.started {
		return nil
	}
	if m.client == nil {
		return fmt.Errorf("websocket client is not configured")
	}
	if err := requireContext(ctx); err != nil {
		return err
	}

	m.ctx, m.cancel = context.WithCancel(ctx)
	m.started = true
	go m.loop()
	return nil
}

// Connect is an alias for Start.
func (m *SubscriptionManager) Connect(ctx context.Context) error {
	return m.Start(ctx)
}

// CallAPI delegates to the websocket client.
func (m *SubscriptionManager) CallAPI(ctx context.Context, apiID int, method string, args []any, reply any) error {
	if m.client == nil {
		return fmt.Errorf("websocket client is not configured")
	}
	return m.client.CallAPI(ctx, apiID, method, args, reply)
}

// Close stops reconnecting and closes the websocket client.
func (m *SubscriptionManager) Close() error {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return nil
	}
	m.closed = true
	cancel := m.cancel
	m.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	if m.client != nil {
		return m.client.Close()
	}
	return nil
}

func (m *SubscriptionManager) loop() {
	for {
		if err := m.waitContextDone(); err != nil {
			return
		}

		if err := m.client.Connect(m.ctx); err != nil {
			m.sendError(err)
			if !m.waitReconnectDelay() {
				return
			}
			continue
		}

		if err := m.resubscribeAll(); err != nil {
			m.sendError(err)
			_ = m.client.Close()
			if !m.waitReconnectDelay() {
				return
			}
			continue
		}

		select {
		case err := <-m.closeCh:
			if m.isClosed() {
				return
			}
			if err != nil && err != ErrShutdown {
				m.sendError(err)
			}
			if !m.waitReconnectDelay() {
				return
			}
		case <-m.ctx.Done():
			return
		}
	}
}

func (m *SubscriptionManager) resubscribeAll() error {
	m.mu.Lock()
	subs := make([]struct {
		name string
		fn   SubscriptionFunc
	}, 0, len(m.subs))
	for name, fn := range m.subs {
		subs = append(subs, struct {
			name string
			fn   SubscriptionFunc
		}{name: name, fn: fn})
	}
	ctx := m.ctx
	m.mu.Unlock()

	for _, sub := range subs {
		if err := m.invokeSubscription(ctx, sub.name, sub.fn); err != nil {
			return err
		}
	}
	return nil
}

func (m *SubscriptionManager) invokeSubscription(ctx context.Context, name string, fn SubscriptionFunc) error {
	if err := fn(ctx, m.client); err != nil {
		return fmt.Errorf("subscription %q: %w", name, err)
	}
	return nil
}

func (m *SubscriptionManager) sendError(err error) {
	if err == nil {
		return
	}
	select {
	case m.errors <- err:
	default:
	}
}

func (m *SubscriptionManager) waitReconnectDelay() bool {
	timer := time.NewTimer(m.reconnectDelay)
	defer timer.Stop()

	select {
	case <-timer.C:
		return true
	case <-m.ctx.Done():
		return false
	}
}

func (m *SubscriptionManager) waitContextDone() error {
	select {
	case <-m.ctx.Done():
		return m.ctx.Err()
	default:
		return nil
	}
}

func (m *SubscriptionManager) isClosed() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.closed
}

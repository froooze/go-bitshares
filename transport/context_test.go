package transport

import (
	"errors"
	"net/url"
	"testing"
)

func TestRESTClientDoRequiresContext(t *testing.T) {
	t.Parallel()

	base, err := url.Parse("https://node.example")
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}

	client := &RESTClient{baseURL: base, client: normalizeHTTPClient(nil)}
	err = client.Do(nil, "GET", "/api/v1/test", nil, nil)
	if !errors.Is(err, ErrNilContext) {
		t.Fatalf("Do(nil, ...) error = %v, want %v", err, ErrNilContext)
	}
}

func TestRPCClientCallRequiresContext(t *testing.T) {
	t.Parallel()

	client := NewRPCClient("https://node.example")
	err := client.Call(nil, "call", nil, nil)
	if !errors.Is(err, ErrNilContext) {
		t.Fatalf("Call(nil, ...) error = %v, want %v", err, ErrNilContext)
	}
}

func TestSubscriptionManagerStartRequiresContext(t *testing.T) {
	t.Parallel()

	ws, err := NewWebsocketClient("ws://node.example")
	if err != nil {
		t.Fatalf("NewWebsocketClient() error = %v", err)
	}

	manager := NewSubscriptionManager(ws, 0)
	err = manager.Start(nil)
	if !errors.Is(err, ErrNilContext) {
		t.Fatalf("Start(nil) error = %v, want %v", err, ErrNilContext)
	}
}

func TestWebsocketClientConnectRequiresContext(t *testing.T) {
	t.Parallel()

	client, err := NewWebsocketClient("ws://node.example")
	if err != nil {
		t.Fatalf("NewWebsocketClient() error = %v", err)
	}

	err = client.Connect(nil)
	if !errors.Is(err, ErrNilContext) {
		t.Fatalf("Connect(nil) error = %v, want %v", err, ErrNilContext)
	}
}

func TestWebsocketClientCallAPIRequiresContext(t *testing.T) {
	t.Parallel()

	client, err := NewWebsocketClient("ws://node.example")
	if err != nil {
		t.Fatalf("NewWebsocketClient() error = %v", err)
	}

	err = client.CallAPI(nil, 0, "call", nil, nil)
	if !errors.Is(err, ErrNilContext) {
		t.Fatalf("CallAPI(nil, ...) error = %v, want %v", err, ErrNilContext)
	}
}

package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// RPCClient performs JSON-RPC requests over HTTP.
type RPCClient struct {
	endpoint string
	client   *http.Client
	nextID   atomic.Uint64
}

func NewRPCClient(endpoint string) *RPCClient {
	return &RPCClient{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *RPCClient) Close() error { return nil }

func (c *RPCClient) Call(ctx context.Context, method string, args []any, reply any) error {
	if err := requireContext(ctx); err != nil {
		return err
	}

	req := RPCRequest{
		Method: method,
		Params: args,
		ID:     c.nextID.Add(1),
	}

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var envelope RPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return err
	}

	if envelope.Error != nil {
		return envelope.Error
	}

	if reply == nil || envelope.Result == nil {
		return nil
	}

	if err := json.Unmarshal(*envelope.Result, reply); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

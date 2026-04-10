package bitshares

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestBitSharesRESTClientPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		call       func(*BitSharesRESTClient, any) error
		wantMethod string
		wantPath   string
	}{
		{
			name: "account balances",
			call: func(c *BitSharesRESTClient, reply any) error {
				return c.Account("alice").Balances(context.Background(), reply)
			},
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/accounts/alice/balances",
		},
		{
			name: "asset dynamic data",
			call: func(c *BitSharesRESTClient, reply any) error {
				return c.Asset("BTS").DynamicData(context.Background(), reply)
			},
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/assets/BTS/dynamic-data",
		},
		{
			name: "market ticker",
			call: func(c *BitSharesRESTClient, reply any) error {
				return c.Market("BTS", "USD").Ticker(context.Background(), reply)
			},
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/markets/BTS/USD/ticker",
		},
		{
			name: "block transactions",
			call: func(c *BitSharesRESTClient, reply any) error {
				return c.Block(42).Transactions(context.Background(), reply)
			},
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/blocks/42/transactions",
		},
		{
			name: "transaction operations",
			call: func(c *BitSharesRESTClient, reply any) error {
				return c.Transaction("tx-123").Operations(context.Background(), reply)
			},
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/transactions/tx-123/operations",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var gotMethod string
			var gotPath string

			httpClient := &http.Client{
				Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
					gotMethod = req.Method
					gotPath = req.URL.Path
					body := io.NopCloser(strings.NewReader(`{"ok":true}`))
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     make(http.Header),
						Body:       body,
						Request:    req,
					}, nil
				}),
			}

			client, err := NewBitSharesRESTClientWithHTTPClient("https://node.example", httpClient)
			if err != nil {
				t.Fatalf("NewBitSharesRESTClient() error = %v", err)
			}

			var reply map[string]bool
			err = tc.call(client, &reply)
			if err != nil {
				t.Fatalf("call() error = %v", err)
			}

			if gotMethod != tc.wantMethod {
				t.Fatalf("method = %q, want %q", gotMethod, tc.wantMethod)
			}
			if gotPath != tc.wantPath {
				t.Fatalf("path = %q, want %q", gotPath, tc.wantPath)
			}
			if reply["ok"] != true {
				t.Fatalf("reply = %#v, want ok=true", reply)
			}
		})
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

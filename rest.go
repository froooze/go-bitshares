package bitshares

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/froooze/go-bitshares/transport"
)

const bitsharesAPIPrefix = "/api/v1"

// BitSharesRESTClient exposes BitShares-specific REST resources.
type BitSharesRESTClient struct {
	rest *transport.RESTClient
}

// NewBitSharesRESTClient creates a BitShares-specific REST client rooted at the provided endpoint.
func NewBitSharesRESTClient(endpoint string) (*BitSharesRESTClient, error) {
	return NewBitSharesRESTClientWithHTTPClient(endpoint, nil)
}

// NewBitSharesRESTClientWithHTTPClient creates a BitShares-specific REST client with a custom HTTP client.
func NewBitSharesRESTClientWithHTTPClient(endpoint string, httpClient *http.Client) (*BitSharesRESTClient, error) {
	rest, err := transport.NewRESTClientWithHTTPClient(endpoint, httpClient)
	if err != nil {
		return nil, err
	}

	return &BitSharesRESTClient{rest: rest}, nil
}

func (c *BitSharesRESTClient) accountPath(parts ...string) apiPath {
	if c == nil {
		return apiPath{}
	}

	return apiPath{
		client: c.rest,
		path:   bitsharesPath(parts...),
	}
}

// Close releases REST resources. It is a no-op for the current HTTP transport.
func (c *BitSharesRESTClient) Close() error {
	if c == nil || c.rest == nil {
		return nil
	}
	return c.rest.Close()
}

// Account returns the account resource identified by name or account ID.
func (c *BitSharesRESTClient) Account(name string) *AccountResource {
	return &AccountResource{apiPath: c.accountPath("accounts", name)}
}

// Asset returns the asset resource identified by symbol or asset ID.
func (c *BitSharesRESTClient) Asset(symbol string) *AssetResource {
	return &AssetResource{apiPath: c.accountPath("assets", symbol)}
}

// Block returns the block resource identified by block number or block ID.
func (c *BitSharesRESTClient) Block(number uint32) *BlockResource {
	return &BlockResource{apiPath: c.accountPath("blocks", fmt.Sprint(number))}
}

// Transaction returns the transaction resource identified by transaction ID.
func (c *BitSharesRESTClient) Transaction(txID string) *TransactionResource {
	return &TransactionResource{apiPath: c.accountPath("transactions", txID)}
}

// Market returns the market resource for the given base and quote symbols.
func (c *BitSharesRESTClient) Market(baseSymbol, quoteSymbol string) *MarketResource {
	return &MarketResource{apiPath: c.accountPath("markets", baseSymbol, quoteSymbol)}
}

// Witness returns the witness resource identified by name or witness ID.
func (c *BitSharesRESTClient) Witness(name string) *WitnessResource {
	return &WitnessResource{apiPath: c.accountPath("witnesses", name)}
}

// CommitteeMember returns the committee member resource identified by name or ID.
func (c *BitSharesRESTClient) CommitteeMember(name string) *CommitteeMemberResource {
	return &CommitteeMemberResource{apiPath: c.accountPath("committee-members", name)}
}

// Proposal returns the proposal resource identified by proposal ID.
func (c *BitSharesRESTClient) Proposal(id string) *ProposalResource {
	return &ProposalResource{apiPath: c.accountPath("proposals", id)}
}

type apiPath struct {
	client *transport.RESTClient
	path   string
}

func (p apiPath) child(parts ...string) apiPath {
	return apiPath{
		client: p.client,
		path:   joinBitSharesPath(p.path, parts...),
	}
}

func (p apiPath) get(ctx context.Context, reply any) error {
	if p.client == nil {
		return fmt.Errorf("bitshares REST client is not configured")
	}
	return p.client.Get(ctx, p.path, reply)
}

func (p apiPath) post(ctx context.Context, body any, reply any) error {
	if p.client == nil {
		return fmt.Errorf("bitshares REST client is not configured")
	}
	return p.client.Post(ctx, p.path, body, reply)
}

func (p apiPath) put(ctx context.Context, body any, reply any) error {
	if p.client == nil {
		return fmt.Errorf("bitshares REST client is not configured")
	}
	return p.client.Put(ctx, p.path, body, reply)
}

func (p apiPath) patch(ctx context.Context, body any, reply any) error {
	if p.client == nil {
		return fmt.Errorf("bitshares REST client is not configured")
	}
	return p.client.Patch(ctx, p.path, body, reply)
}

func (p apiPath) delete(ctx context.Context, reply any) error {
	if p.client == nil {
		return fmt.Errorf("bitshares REST client is not configured")
	}
	return p.client.Delete(ctx, p.path, reply)
}

// AccountResource models a BitShares account endpoint.
type AccountResource struct {
	apiPath
}

// Get fetches the account object.
func (r *AccountResource) Get(ctx context.Context, reply any) error { return r.get(ctx, reply) }

// Balances fetches the account balances collection.
func (r *AccountResource) Balances(ctx context.Context, reply any) error {
	return r.child("balances").get(ctx, reply)
}

// History fetches the account history collection.
func (r *AccountResource) History(ctx context.Context, reply any) error {
	return r.child("history").get(ctx, reply)
}

// Orders fetches the account orders collection.
func (r *AccountResource) Orders(ctx context.Context, reply any) error {
	return r.child("orders").get(ctx, reply)
}

// OpenOrders fetches the account's open orders collection.
func (r *AccountResource) OpenOrders(ctx context.Context, reply any) error {
	return r.child("open-orders").get(ctx, reply)
}

// AssetResource models a BitShares asset endpoint.
type AssetResource struct {
	apiPath
}

// Get fetches the asset object.
func (r *AssetResource) Get(ctx context.Context, reply any) error { return r.get(ctx, reply) }

// DynamicData fetches the asset dynamic data object.
func (r *AssetResource) DynamicData(ctx context.Context, reply any) error {
	return r.child("dynamic-data").get(ctx, reply)
}

// Holders fetches the asset holders collection.
func (r *AssetResource) Holders(ctx context.Context, reply any) error {
	return r.child("holders").get(ctx, reply)
}

// BlockResource models a BitShares block endpoint.
type BlockResource struct {
	apiPath
}

// Get fetches the block object.
func (r *BlockResource) Get(ctx context.Context, reply any) error { return r.get(ctx, reply) }

// Transactions fetches the transactions contained in the block.
func (r *BlockResource) Transactions(ctx context.Context, reply any) error {
	return r.child("transactions").get(ctx, reply)
}

// TransactionResource models a BitShares transaction endpoint.
type TransactionResource struct {
	apiPath
}

// Get fetches the transaction object.
func (r *TransactionResource) Get(ctx context.Context, reply any) error { return r.get(ctx, reply) }

// Operations fetches the operations contained in the transaction.
func (r *TransactionResource) Operations(ctx context.Context, reply any) error {
	return r.child("operations").get(ctx, reply)
}

// MarketResource models a BitShares market endpoint.
type MarketResource struct {
	apiPath
}

// Get fetches the market object.
func (r *MarketResource) Get(ctx context.Context, reply any) error { return r.get(ctx, reply) }

// Ticker fetches the market ticker.
func (r *MarketResource) Ticker(ctx context.Context, reply any) error {
	return r.child("ticker").get(ctx, reply)
}

// OrderBook fetches the market order book.
func (r *MarketResource) OrderBook(ctx context.Context, reply any) error {
	return r.child("order-book").get(ctx, reply)
}

// Trades fetches the market trade history.
func (r *MarketResource) Trades(ctx context.Context, reply any) error {
	return r.child("trades").get(ctx, reply)
}

// WitnessResource models a BitShares witness endpoint.
type WitnessResource struct {
	apiPath
}

// Get fetches the witness object.
func (r *WitnessResource) Get(ctx context.Context, reply any) error { return r.get(ctx, reply) }

// Votes fetches the witness votes collection.
func (r *WitnessResource) Votes(ctx context.Context, reply any) error {
	return r.child("votes").get(ctx, reply)
}

// CommitteeMemberResource models a BitShares committee member endpoint.
type CommitteeMemberResource struct {
	apiPath
}

// Get fetches the committee member object.
func (r *CommitteeMemberResource) Get(ctx context.Context, reply any) error { return r.get(ctx, reply) }

// Votes fetches the committee member votes collection.
func (r *CommitteeMemberResource) Votes(ctx context.Context, reply any) error {
	return r.child("votes").get(ctx, reply)
}

// ProposalResource models a BitShares proposal endpoint.
type ProposalResource struct {
	apiPath
}

// Get fetches the proposal object.
func (r *ProposalResource) Get(ctx context.Context, reply any) error { return r.get(ctx, reply) }

// Operations fetches the proposal operations collection.
func (r *ProposalResource) Operations(ctx context.Context, reply any) error {
	return r.child("operations").get(ctx, reply)
}

func bitsharesPath(parts ...string) string {
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			filtered = append(filtered, strings.Trim(trimmed, "/"))
		}
	}
	if len(filtered) == 0 {
		return bitsharesAPIPrefix
	}

	return path.Join(bitsharesAPIPrefix, path.Join(filtered...))
}

func joinBitSharesPath(base string, parts ...string) string {
	segments := make([]string, 0, 1+len(parts))
	if trimmed := strings.TrimSpace(base); trimmed != "" {
		segments = append(segments, strings.Trim(trimmed, "/"))
	}
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			segments = append(segments, strings.Trim(trimmed, "/"))
		}
	}
	if len(segments) == 0 {
		return bitsharesAPIPrefix
	}

	return "/" + path.Join(segments...)
}

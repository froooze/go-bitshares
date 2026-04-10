package bitshares

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/froooze/go-bitshares/ecc"
	"github.com/froooze/go-bitshares/protocol"
	"github.com/froooze/go-bitshares/transport"
)

const bootstrapSubscriptionName = "__bootstrap__"

// ChainConfig describes the active BitShares network.
type ChainConfig struct {
	Name                  string
	CoreAsset             string
	AddressPrefix         string
	ExpireInSecs          int
	ExpireInSecsProposal  int
	ReviewInSecsCommittee int
	ChainID               string
}

var defaultChains = []ChainConfig{
	{
		Name:                  "BitShares",
		CoreAsset:             "BTS",
		AddressPrefix:         "BTS",
		ExpireInSecs:          15,
		ExpireInSecsProposal:  24 * 60 * 60,
		ReviewInSecsCommittee: 24 * 60 * 60,
		ChainID:               "4018d7844c78f6a6c41c6a552b898022310fc5dec06da467ee7905a8dad512c8",
	},
	{
		Name:                  "TestNet",
		CoreAsset:             "TEST",
		AddressPrefix:         "TEST",
		ExpireInSecs:          15,
		ExpireInSecsProposal:  24 * 60 * 60,
		ReviewInSecsCommittee: 24 * 60 * 60,
		ChainID:               "39f5e2ede1f8bc1a3a54a7914414e3779e33193f1f5693510e73cb7a87617447",
	},
}

// KeyWeightPair models a BitShares key/weight authority entry.
type KeyWeightPair struct {
	Key    string
	Weight uint16
}

// UnmarshalJSON parses the BitShares array form [pubkey, weight].
func (p *KeyWeightPair) UnmarshalJSON(data []byte) error {
	var raw []any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) != 2 {
		return fmt.Errorf("invalid key auth pair")
	}
	key, ok := raw[0].(string)
	if !ok {
		return fmt.Errorf("invalid key auth key")
	}
	weight, err := toUint16(raw[1])
	if err != nil {
		return err
	}
	p.Key = key
	p.Weight = weight
	return nil
}

// AccountWeightPair models a BitShares account/weight authority entry.
type AccountWeightPair struct {
	Account string
	Weight  uint16
}

// UnmarshalJSON parses the BitShares array form [account, weight].
func (p *AccountWeightPair) UnmarshalJSON(data []byte) error {
	var raw []any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) != 2 {
		return fmt.Errorf("invalid account auth pair")
	}
	account, ok := raw[0].(string)
	if !ok {
		return fmt.Errorf("invalid account auth account")
	}
	weight, err := toUint16(raw[1])
	if err != nil {
		return err
	}
	p.Account = account
	p.Weight = weight
	return nil
}

// Authority describes an account authority tree.
type Authority struct {
	WeightThreshold uint32              `json:"weight_threshold"`
	KeyAuths        []KeyWeightPair     `json:"key_auths"`
	AccountAuths    []AccountWeightPair `json:"account_auths"`
	Extensions      []json.RawMessage   `json:"extensions"`
}

// AccountOptions contains the account options field subset used by the library.
type AccountOptions struct {
	MemoKey    string            `json:"memo_key"`
	Extensions []json.RawMessage `json:"extensions"`
}

// AccountInfo is a cached BitShares account record.
type AccountInfo struct {
	ID      protocol.ObjectID `json:"id"`
	Name    string            `json:"name"`
	Owner   Authority         `json:"owner"`
	Active  Authority         `json:"active"`
	Options AccountOptions    `json:"options"`
}

// AssetInfo is a cached BitShares asset record.
type AssetInfo struct {
	ID               protocol.ObjectID `json:"id"`
	Symbol           string            `json:"symbol"`
	Precision        uint8             `json:"precision"`
	MarketFeePercent uint16            `json:"market_fee_percent"`
	Options          json.RawMessage   `json:"options"`
}

// ChainClient bootstraps and caches a BitShares websocket connection.
type ChainClient struct {
	ws    *transport.WebsocketClient
	wsMgr *transport.SubscriptionManager

	mu      sync.RWMutex
	apiIDs  map[string]int
	chain   ChainConfig
	ready   chan struct{}
	bootErr chan error
	started bool

	accountMu   sync.RWMutex
	accounts    map[string]*AccountInfo
	accountByID map[string]*AccountInfo

	assetMu   sync.RWMutex
	assets    map[string]*AssetInfo
	assetByID map[string]*AssetInfo
}

// NewChainClient creates a websocket-backed BitShares chain client.
func NewChainClient(endpoint string, reconnectDelay time.Duration) (*ChainClient, error) {
	ws, err := transport.NewWebsocketClient(endpoint)
	if err != nil {
		return nil, err
	}

	c := &ChainClient{
		ws:          ws,
		wsMgr:       transport.NewSubscriptionManager(ws, reconnectDelay),
		apiIDs:      make(map[string]int),
		ready:       make(chan struct{}),
		bootErr:     make(chan error, 1),
		accounts:    make(map[string]*AccountInfo),
		accountByID: make(map[string]*AccountInfo),
		assets:      make(map[string]*AssetInfo),
		assetByID:   make(map[string]*AssetInfo),
	}

	if err := c.wsMgr.Register(bootstrapSubscriptionName, c.bootstrap); err != nil {
		return nil, err
	}

	return c, nil
}

// Connect opens the websocket, discovers APIs, and waits for the initial bootstrap.
func (c *ChainClient) Connect(ctx context.Context) error {
	if err := requireContext(ctx); err != nil {
		return err
	}

	c.mu.Lock()
	if c.started {
		ready := c.ready
		c.mu.Unlock()
		select {
		case <-ready:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	c.started = true
	c.mu.Unlock()

	if err := c.wsMgr.Connect(ctx); err != nil {
		return err
	}

	select {
	case <-c.ready:
		return nil
	case err := <-c.bootErr:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close closes the websocket chain client.
func (c *ChainClient) Close() error {
	if c == nil || c.wsMgr == nil {
		return nil
	}
	return c.wsMgr.Close()
}

// Chain returns the current chain configuration.
func (c *ChainClient) Chain() ChainConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.chain
}

// APIID returns a discovered websocket API id.
func (c *ChainClient) APIID(name string) (int, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	id, ok := c.apiIDs[name]
	return id, ok
}

// Account returns a cached account by name.
func (c *ChainClient) Account(ctx context.Context, name string) (*AccountInfo, error) {
	if err := requireContext(ctx); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("empty account name")
	}
	if cached := c.cachedAccount(strings.ToLower(name)); cached != nil {
		return cached, nil
	}

	var reply AccountInfo
	if err := c.callDatabase(ctx, "get_account_by_name", []any{strings.ToLower(name)}, &reply); err != nil {
		return nil, err
	}
	c.storeAccount(&reply)
	return &reply, nil
}

// AccountByID returns a cached account by object id.
func (c *ChainClient) AccountByID(ctx context.Context, id string) (*AccountInfo, error) {
	if err := requireContext(ctx); err != nil {
		return nil, err
	}
	if id == "" {
		return nil, fmt.Errorf("empty account id")
	}
	if cached := c.cachedAccountByID(id); cached != nil {
		return cached, nil
	}

	var reply []AccountInfo
	if err := c.callDatabase(ctx, "get_accounts", []any{[]string{id}}, &reply); err != nil {
		return nil, err
	}
	if len(reply) == 0 {
		return nil, fmt.Errorf("account %s not found", id)
	}
	c.storeAccount(&reply[0])
	return &reply[0], nil
}

// Asset returns a cached asset by symbol.
func (c *ChainClient) Asset(ctx context.Context, symbol string) (*AssetInfo, error) {
	if err := requireContext(ctx); err != nil {
		return nil, err
	}
	if symbol == "" {
		return nil, fmt.Errorf("empty asset symbol")
	}
	if cached := c.cachedAsset(strings.ToUpper(symbol)); cached != nil {
		return cached, nil
	}

	var reply []AssetInfo
	if err := c.callDatabase(ctx, "lookup_asset_symbols", []any{[]string{strings.ToUpper(symbol)}}, &reply); err != nil {
		if err := c.callDatabase(ctx, "list_assets", []any{strings.ToUpper(symbol), 1}, &reply); err != nil {
			return nil, err
		}
	}
	if len(reply) == 0 {
		return nil, fmt.Errorf("asset %s not found", symbol)
	}
	c.storeAsset(&reply[0])
	return &reply[0], nil
}

// AssetByID returns a cached asset by object id.
func (c *ChainClient) AssetByID(ctx context.Context, id string) (*AssetInfo, error) {
	if err := requireContext(ctx); err != nil {
		return nil, err
	}
	if id == "" {
		return nil, fmt.Errorf("empty asset id")
	}
	if cached := c.cachedAssetByID(id); cached != nil {
		return cached, nil
	}

	var reply []AssetInfo
	if err := c.callDatabase(ctx, "get_assets", []any{[]string{id}}, &reply); err != nil {
		return nil, err
	}
	if len(reply) == 0 {
		return nil, fmt.Errorf("asset %s not found", id)
	}
	c.storeAsset(&reply[0])
	return &reply[0], nil
}

// Errors exposes reconnect and subscription errors from the websocket manager.
func (c *ChainClient) Errors() <-chan error {
	if c == nil || c.wsMgr == nil {
		return nil
	}
	return c.wsMgr.Errors()
}

// Notifications exposes raw websocket notifications.
func (c *ChainClient) Notifications() <-chan json.RawMessage {
	if c == nil || c.wsMgr == nil {
		return nil
	}
	return c.wsMgr.Notifications()
}

func (c *ChainClient) bootstrap(ctx context.Context, ws *transport.WebsocketClient) (err error) {
	defer func() {
		if err != nil {
			c.signalReady(err)
		}
	}()

	if ws == nil {
		return fmt.Errorf("websocket client is not configured")
	}
	if err := requireContext(ctx); err != nil {
		return err
	}

	var ok bool
	if err = ws.CallAPI(ctx, 1, "login", []any{"", ""}, &ok); err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("login rejected by node")
	}

	apiNames := []string{"database", "history", "network_broadcast", "block", "asset", "orders", "crypto"}
	nextAPIs := make(map[string]int, len(apiNames))
	for _, name := range apiNames {
		var id int
		if err = ws.CallAPI(ctx, 1, name, []any{}, &id); err != nil {
			if name == "database" {
				return err
			}
			continue
		}
		nextAPIs[name] = id
	}

	var chainID string
	if err = ws.CallAPI(ctx, nextAPIs["database"], "get_chain_id", []any{}, &chainID); err != nil {
		return err
	}

	cfg, ok := lookupChain(chainID)
	if !ok {
		return fmt.Errorf("unknown chain id %s", chainID)
	}

	c.mu.Lock()
	c.apiIDs = nextAPIs
	c.chain = cfg
	c.mu.Unlock()

	ecc.SetAddressPrefix(cfg.AddressPrefix)
	c.signalReady(nil)
	return nil
}

func (c *ChainClient) signalReady(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-c.ready:
		// already signaled
		return
	default:
	}

	if err != nil {
		select {
		case c.bootErr <- err:
		default:
		}
		return
	}

	close(c.ready)
}

func lookupChain(chainID string) (ChainConfig, bool) {
	for _, cfg := range defaultChains {
		if strings.EqualFold(cfg.ChainID, chainID) {
			return cfg, true
		}
	}
	return ChainConfig{}, false
}

func (c *ChainClient) callDatabase(ctx context.Context, method string, args []any, reply any) error {
	if err := requireContext(ctx); err != nil {
		return err
	}
	apiID, ok := c.APIID("database")
	if !ok {
		return fmt.Errorf("database API is not available")
	}
	if c.wsMgr != nil {
		return c.wsMgr.CallAPI(ctx, apiID, method, args, reply)
	}
	if c.ws == nil {
		return fmt.Errorf("websocket client is not configured")
	}
	return c.ws.CallAPI(ctx, apiID, method, args, reply)
}

func (c *ChainClient) cachedAccount(name string) *AccountInfo {
	c.accountMu.RLock()
	defer c.accountMu.RUnlock()
	return c.accounts[strings.ToLower(name)]
}

func (c *ChainClient) cachedAccountByID(id string) *AccountInfo {
	c.accountMu.RLock()
	defer c.accountMu.RUnlock()
	return c.accountByID[id]
}

func (c *ChainClient) storeAccount(acc *AccountInfo) {
	if acc == nil {
		return
	}
	c.accountMu.Lock()
	defer c.accountMu.Unlock()
	c.accounts[strings.ToLower(acc.Name)] = acc
	c.accountByID[acc.ID.String()] = acc
}

func (c *ChainClient) cachedAsset(symbol string) *AssetInfo {
	c.assetMu.RLock()
	defer c.assetMu.RUnlock()
	return c.assets[strings.ToUpper(symbol)]
}

func (c *ChainClient) cachedAssetByID(id string) *AssetInfo {
	c.assetMu.RLock()
	defer c.assetMu.RUnlock()
	return c.assetByID[id]
}

func (c *ChainClient) storeAsset(asset *AssetInfo) {
	if asset == nil {
		return
	}
	c.assetMu.Lock()
	defer c.assetMu.Unlock()
	c.assets[strings.ToUpper(asset.Symbol)] = asset
	c.assetByID[asset.ID.String()] = asset
}

func toUint16(v any) (uint16, error) {
	switch n := v.(type) {
	case float64:
		return uint16(n), nil
	case float32:
		return uint16(n), nil
	case int:
		return uint16(n), nil
	case int64:
		return uint16(n), nil
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			return 0, err
		}
		return uint16(i), nil
	default:
		return 0, fmt.Errorf("invalid weight value %T", v)
	}
}

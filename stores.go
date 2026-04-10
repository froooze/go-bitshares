package bitshares

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/froooze/go-bitshares/protocol"
)

// AccountStore wraps the chain cache with a btsdex-style account lookup API.
type AccountStore struct {
	chain *ChainClient
}

func (s *AccountStore) Get(ctx context.Context, name string) (*AccountInfo, error) {
	if s == nil || s.chain == nil {
		return nil, fmt.Errorf("account store is not configured")
	}
	return s.chain.Account(ctx, name)
}

func (s *AccountStore) ID(ctx context.Context, id string) (*AccountInfo, error) {
	if s == nil || s.chain == nil {
		return nil, fmt.Errorf("account store is not configured")
	}
	return s.chain.AccountByID(ctx, id)
}

func (s *AccountStore) Update(ctx context.Context) error {
	if s == nil || s.chain == nil {
		return fmt.Errorf("account store is not configured")
	}
	s.chain.accountMu.RLock()
	ids := make([]string, 0, len(s.chain.accountByID))
	for id := range s.chain.accountByID {
		ids = append(ids, id)
	}
	s.chain.accountMu.RUnlock()
	if len(ids) == 0 {
		return nil
	}

	var reply []AccountInfo
	if err := s.chain.callDatabase(ctx, "get_accounts", []any{ids}, &reply); err != nil {
		return err
	}
	for i := range reply {
		s.chain.storeAccount(&reply[i])
	}
	return nil
}

// AssetStore wraps the chain cache with a btsdex-style asset lookup API.
type AssetStore struct {
	chain *ChainClient
}

func (s *AssetStore) Get(ctx context.Context, symbol string) (*AssetInfo, error) {
	if s == nil || s.chain == nil {
		return nil, fmt.Errorf("asset store is not configured")
	}
	return s.chain.Asset(ctx, symbol)
}

func (s *AssetStore) ID(ctx context.Context, id string) (*AssetInfo, error) {
	if s == nil || s.chain == nil {
		return nil, fmt.Errorf("asset store is not configured")
	}
	return s.chain.AssetByID(ctx, id)
}

func (s *AssetStore) Update(ctx context.Context) error {
	if s == nil || s.chain == nil {
		return fmt.Errorf("asset store is not configured")
	}
	s.chain.assetMu.RLock()
	ids := make([]string, 0, len(s.chain.assetByID))
	for id := range s.chain.assetByID {
		ids = append(ids, id)
	}
	s.chain.assetMu.RUnlock()
	if len(ids) == 0 {
		return nil
	}

	var reply []AssetInfo
	if err := s.chain.callDatabase(ctx, "get_assets", []any{ids}, &reply); err != nil {
		return err
	}
	for i := range reply {
		s.chain.storeAsset(&reply[i])
	}
	return nil
}

// FeeTable keeps a lightweight snapshot of the current chain fees.
type FeeTable struct {
	chain *ChainClient

	mu        sync.RWMutex
	byName    map[string]float64
	byType    map[protocol.OperationType]float64
	operation []string
}

func newFeeTable(chain *ChainClient) *FeeTable {
	return &FeeTable{
		chain:  chain,
		byName: make(map[string]float64),
		byType: make(map[protocol.OperationType]float64),
	}
}

// Update refreshes the cached fee snapshot from the chain.
func (f *FeeTable) Update(ctx context.Context) error {
	if f == nil || f.chain == nil {
		return fmt.Errorf("fee table is not configured")
	}

	var props GlobalProperties
	if err := f.chain.GetGlobalProperties(ctx, &props); err != nil {
		return err
	}

	byName := make(map[string]float64, len(props.Parameters.CurrentFees.Parameters))
	byType := make(map[protocol.OperationType]float64, len(props.Parameters.CurrentFees.Parameters))
	names := make([]string, 0, len(props.Parameters.CurrentFees.Parameters))

	for _, param := range props.Parameters.CurrentFees.Parameters {
		feeValue, ok := param.FeeValue()
		if !ok {
			continue
		}
		normalized := feeValue / 100000
		byType[param.OperationType] = normalized
		byName[param.OperationType.String()] = normalized
		names = append(names, param.OperationType.String())
	}

	sort.Strings(names)

	f.mu.Lock()
	f.byName = byName
	f.byType = byType
	f.operation = names
	f.mu.Unlock()

	return nil
}

// Get returns a cached fee by operation name.
func (f *FeeTable) Get(name string) (float64, bool) {
	if f == nil {
		return 0, false
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	fee, ok := f.byName[strings.ToLower(name)]
	if ok {
		return fee, true
	}
	fee, ok = f.byName[name]
	return fee, ok
}

// ByType returns a cached fee by operation type.
func (f *FeeTable) ByType(kind protocol.OperationType) (float64, bool) {
	if f == nil {
		return 0, false
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	fee, ok := f.byType[kind]
	return fee, ok
}

// Operations lists the currently cached fee operation names.
func (f *FeeTable) Operations() []string {
	if f == nil {
		return nil
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	out := make([]string, len(f.operation))
	copy(out, f.operation)
	return out
}

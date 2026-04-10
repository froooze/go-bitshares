package bitshares

import (
	"context"
	"fmt"

	"github.com/froooze/go-bitshares/protocol"
)

// GlobalProperties mirrors the subset of get_global_properties used by the library.
type GlobalProperties struct {
	Parameters ChainParameters `json:"parameters"`
}

// ChainParameters mirrors the subset of chain_parameters used by the library.
type ChainParameters struct {
	CurrentFees protocol.FeeSchedule `json:"current_fees"`
}

type FeeSchedule = protocol.FeeSchedule
type FeeScheduleParameter = protocol.FeeScheduleParameter

func (c *ChainClient) callAPI(ctx context.Context, api string, method string, args []any, reply any) error {
	if err := requireContext(ctx); err != nil {
		return err
	}
	apiID, ok := c.APIID(api)
	if !ok {
		return fmt.Errorf("%s API is not available", api)
	}
	if c.wsMgr != nil {
		return c.wsMgr.CallAPI(ctx, apiID, method, args, reply)
	}
	if c.ws == nil {
		return fmt.Errorf("websocket client is not configured")
	}
	return c.ws.CallAPI(ctx, apiID, method, args, reply)
}

// GetGlobalProperties fetches the chain global properties into reply.
func (c *ChainClient) GetGlobalProperties(ctx context.Context, reply any) error {
	return c.callDatabase(ctx, "get_global_properties", []any{}, reply)
}

// GetDynamicGlobalProperties fetches the dynamic global properties into reply.
func (c *ChainClient) GetDynamicGlobalProperties(ctx context.Context, reply any) error {
	return c.callDatabase(ctx, "get_dynamic_global_properties", []any{}, reply)
}

// GetConfig fetches compile-time configuration into reply.
func (c *ChainClient) GetConfig(ctx context.Context, reply any) error {
	return c.callDatabase(ctx, "get_config", []any{}, reply)
}

// GetObjects fetches raw chain objects by id.
func (c *ChainClient) GetObjects(ctx context.Context, ids []string, reply any) error {
	return c.callDatabase(ctx, "get_objects", []any{ids}, reply)
}

// LookupAccounts fetches accounts by lower-bound name and limit.
func (c *ChainClient) LookupAccounts(ctx context.Context, lowerBoundName string, limit uint16, reply any) error {
	return c.callDatabase(ctx, "lookup_accounts", []any{lowerBoundName, limit}, reply)
}

// LookupAssetSymbols fetches assets by symbol or id.
func (c *ChainClient) LookupAssetSymbols(ctx context.Context, symbols []string, reply any) error {
	return c.callDatabase(ctx, "lookup_asset_symbols", []any{symbols}, reply)
}

// ListAssets fetches assets beginning at the lower bound symbol.
func (c *ChainClient) ListAssets(ctx context.Context, lowerBoundSymbol string, limit int, reply any) error {
	return c.callDatabase(ctx, "list_assets", []any{lowerBoundSymbol, limit}, reply)
}

// GetAccountBalances returns the balances for one account and a set of asset ids.
func (c *ChainClient) GetAccountBalances(ctx context.Context, accountID string, assetIDs []string, reply any) error {
	return c.callDatabase(ctx, "get_account_balances", []any{accountID, assetIDs}, reply)
}

// GetNamedAccountBalances returns balances for an account name.
func (c *ChainClient) GetNamedAccountBalances(ctx context.Context, account string, assetIDs []string, reply any) error {
	return c.callDatabase(ctx, "get_named_account_balances", []any{account, assetIDs}, reply)
}

// GetFullAccounts fetches full account records and optionally subscribes.
func (c *ChainClient) GetFullAccounts(ctx context.Context, ids []string, subscribe bool, reply any) error {
	return c.callDatabase(ctx, "get_full_accounts", []any{ids, subscribe}, reply)
}

// GetAccountHistory fetches account history entries.
func (c *ChainClient) GetAccountHistory(ctx context.Context, accountID, stop string, limit int, start string, reply any) error {
	return c.callAPI(ctx, "history", "get_account_history", []any{accountID, stop, limit, start}, reply)
}

// GetRelativeAccountHistory fetches relative account history entries.
func (c *ChainClient) GetRelativeAccountHistory(ctx context.Context, accountID, stop string, limit uint32, start string, reply any) error {
	return c.callDatabase(ctx, "get_relative_account_history", []any{accountID, stop, limit, start}, reply)
}

// GetFillOrderHistory fetches fill order history for a market.
func (c *ChainClient) GetFillOrderHistory(ctx context.Context, baseID, quoteID string, limit int, reply any) error {
	return c.callAPI(ctx, "history", "get_fill_order_history", []any{baseID, quoteID, limit}, reply)
}

// GetMarketHistoryBuckets fetches the market history buckets.
func (c *ChainClient) GetMarketHistoryBuckets(ctx context.Context, reply any) error {
	return c.callAPI(ctx, "history", "get_market_history_buckets", []any{}, reply)
}

// GetMarketHistory fetches bucketed trade history for a market.
func (c *ChainClient) GetMarketHistory(ctx context.Context, baseID, quoteID string, bucketSeconds int, start, stop string, reply any) error {
	return c.callAPI(ctx, "history", "get_market_history", []any{baseID, quoteID, bucketSeconds, start, stop}, reply)
}

// GetTicker fetches the market ticker for a pair of asset ids.
func (c *ChainClient) GetTicker(ctx context.Context, baseID, quoteID string, reply any) error {
	return c.callDatabase(ctx, "get_ticker", []any{baseID, quoteID}, reply)
}

// GetOrderBook fetches the order book for a pair of asset ids.
func (c *ChainClient) GetOrderBook(ctx context.Context, baseID, quoteID string, limit int, reply any) error {
	return c.callDatabase(ctx, "get_order_book", []any{baseID, quoteID, limit}, reply)
}

// GetLimitOrders fetches open limit orders for a market.
func (c *ChainClient) GetLimitOrders(ctx context.Context, baseID, quoteID string, limit int, reply any) error {
	return c.callDatabase(ctx, "get_limit_orders", []any{baseID, quoteID, limit}, reply)
}

// GetRequiredFees fetches required fees for the supplied operations.
func (c *ChainClient) GetRequiredFees(ctx context.Context, ops []protocol.OperationEnvelope, feeAssetID string, reply any) error {
	return c.callDatabase(ctx, "get_required_fees", []any{ops, feeAssetID}, reply)
}

// SetSubscribeCallback registers a database subscription callback.
func (c *ChainClient) SetSubscribeCallback(ctx context.Context, callback any, clearFilter bool) error {
	return c.callDatabase(ctx, "set_subscribe_callback", []any{callback, clearFilter}, nil)
}

// SetBlockAppliedCallback registers a block applied callback.
func (c *ChainClient) SetBlockAppliedCallback(ctx context.Context, callback any) error {
	return c.callDatabase(ctx, "set_block_applied_callback", []any{callback}, nil)
}

// SubscribeToMarket registers a market subscription callback.
func (c *ChainClient) SubscribeToMarket(ctx context.Context, callback any, baseID, quoteID string, clearFilter bool) error {
	return c.callDatabase(ctx, "subscribe_to_market", []any{callback, baseID, quoteID, clearFilter}, nil)
}

// CancelAllSubscriptions cancels all websocket subscriptions.
func (c *ChainClient) CancelAllSubscriptions(ctx context.Context) error {
	return c.callDatabase(ctx, "cancel_all_subscriptions", []any{}, nil)
}

// BroadcastTransaction sends a signed transaction to the network broadcast API.
func (c *ChainClient) BroadcastTransaction(ctx context.Context, tx *protocol.SignedTransaction, reply any) error {
	return c.callAPI(ctx, "network_broadcast", "broadcast_transaction", []any{tx}, reply)
}

// BroadcastTransactionSynchronous sends a signed transaction and waits for confirmation.
func (c *ChainClient) BroadcastTransactionSynchronous(ctx context.Context, tx *protocol.SignedTransaction, reply any) error {
	return c.callAPI(ctx, "network_broadcast", "broadcast_transaction_synchronous", []any{tx}, reply)
}

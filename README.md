# go-bitshares

Standalone BitShares Go library.

First-release disclaimer: this is the initial public release of the library.
The API surface, behavior, and compatibility guarantees may still evolve, so
pin a tagged version before relying on it in production.

The module is intentionally small and split into a few stable packages:

- `bitshares`: client facade and BitShares-specific REST resources
- `transport`: HTTP JSON-RPC, websocket RPC, and subscription management
- `protocol`: core chain primitives, transaction envelopes, and operation types
- `sign`: signer abstraction
- `ecc`: BitShares key, memo, and address helpers

## What It Exposes

- BitShares core 7.0.2 operation tags
- typed protocol objects plus wallet-side builders for the common trading, governance, and asset-admin flows
- binary marshal/unmarshal coverage for the non-virtual BitShares 7.0.2 operations used in signed transactions
- high-level `BitShares` facade with account and asset stores
- fee-table snapshots from `get_global_properties`
- wallet/session helpers for login, backup restore, memo encryption, balance lookup, and typed operation builders
- transaction serialization, digesting, and compact signing for broadcast-ready payloads
- raw operation fallback for unmodeled chain operations
- JSON-RPC and websocket transports for node access
- BitShares-specific REST resources for accounts, assets, blocks, markets, witnesses, committee members, proposals, and transactions
- websocket subscription management with automatic reconnect and resubscribe

## Wallet Builder Coverage

The wallet layer now exposes typed builders for the user-signable core operation
families used by trading, account management, governance, and asset
administration:

- transfers, asset issue, asset reserve, balance claim, override transfer
- account create, update, whitelist, upgrade, transfer
- asset create, update, update bitasset, update feed producers, fund fee pool, publish feed, settle, global settle, claim fees
- limit orders, buy/sell helpers, order cancel, call-order update, bid collateral
- proposal create, update, delete
- witness create, update
- committee member create, update, global-parameter update
- withdraw-permission create, update, claim, delete
- HTLC create

Virtual or chain-internal records such as `fill_order`, `asset_settle_cancel`,
`fba_distribute`, and `execute_bid` are intentionally not wrapped as wallet
builders because they are chain-emitted history records, not user-authored
transaction payloads.

## Notes

The old generated compatibility tree is no longer included. This module is
intended for direct use as a standalone dependency, not as a compatibility
shim.

Backup restore uses in-process LZMA decoding and no longer depends on host
`xz` or `lzma` executables.

Public network-facing APIs require an explicit `context.Context`; passing
`nil` returns `ErrNilContext`.

Secret-handling policy:

- secret inputs are accepted as byte slices where practical
- the library does not persist secrets or credentials
- callers should wipe their own password, WIF, and backup buffers after use
- wallet sessions expose public keys only, and `Wallet.Wipe()` should be called
  when signing and memo operations are finished

Binary set serialization now canonicalizes duplicate object IDs, public keys,
and addresses before writing the wire payload, which keeps encoded transactions
aligned with core `flat_set` semantics.

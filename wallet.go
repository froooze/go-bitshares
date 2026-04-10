package bitshares

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/froooze/go-bitshares/ecc"
	"github.com/froooze/go-bitshares/protocol"
	"github.com/froooze/go-bitshares/sign"
)

// MemoData captures a BitShares memo payload.
type MemoData = protocol.MemoData

const maxChainURLLength = 127
const graphene100Percent = 10000
const minAccountNameLength = 1
const maxAccountNameLength = 63
const minAssetSymbolLength = 3
const maxAssetSymbolLength = 16
const maxAssetPrecision = 12
const voteTypeWitness = 0
const voteTypeCommittee = 1
const assetIssuerPermissionWhiteList = 0x02
const assetIssuerPermissionGlobalSettle = 0x20
const assetIssuerPermissionWitnessFed = 0x80
const assetIssuerPermissionCommitteeFed = 0x100
const assetIssuerPermissionDisableBSRMUpdate = 0x4000
const assetIssuerPermissionMask = 0xFFFF

// Wallet represents an authenticated account session.
type Wallet struct {
	parent *BitShares

	account   *AccountInfo
	activeKey *ecc.PrivateKey
	memoKey   *ecc.PrivateKey
	feeAsset  *AssetInfo
}

// Wipe clears the in-memory private keys held by the wallet session.
// It is best-effort only, but should be called as soon as the session is no longer needed.
func (w *Wallet) Wipe() {
	if w == nil {
		return
	}
	if w.activeKey != nil {
		w.activeKey.Wipe()
	}
	if w.memoKey != nil && w.memoKey != w.activeKey {
		w.memoKey.Wipe()
	}
	w.activeKey = nil
	w.memoKey = nil
}

// Account returns the authenticated account record.
func (w *Wallet) Account() *AccountInfo { return w.account }

// FeeAsset returns the current fee asset.
func (w *Wallet) FeeAsset() *AssetInfo { return w.feeAsset }

// ActivePublicKey returns the active public key without exposing the private key.
func (w *Wallet) ActivePublicKey() *ecc.PublicKey {
	if w == nil || w.activeKey == nil {
		return nil
	}
	return w.activeKey.PublicKey()
}

// MemoPublicKey returns the memo public key without exposing the private key.
func (w *Wallet) MemoPublicKey() *ecc.PublicKey {
	if w == nil || w.memoKey == nil {
		return nil
	}
	return w.memoKey.PublicKey()
}

// SetFeeAsset switches the wallet fee asset.
func (w *Wallet) SetFeeAsset(ctx context.Context, symbol string) error {
	if w == nil || w.parent == nil {
		return fmt.Errorf("wallet is not configured")
	}
	asset, err := w.parent.Asset(ctx, symbol)
	if err != nil {
		return err
	}
	w.feeAsset = asset
	return nil
}

// SetMemoKey sets the wallet memo key from WIF bytes.
// The caller owns the input buffer and may wipe it after the call returns.
func (w *Wallet) SetMemoKey(wif []byte) error {
	if len(wif) == 0 {
		return fmt.Errorf("empty memo key")
	}
	key, err := ecc.PrivateKeyFromWIF(wif)
	if err != nil {
		return err
	}
	w.memoKey = key
	return nil
}

// Balances returns balances for the selected asset symbols.
func (w *Wallet) Balances(ctx context.Context, assetSymbols ...string) ([]protocol.AssetAmount, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}

	ids := make([]string, 0, len(assetSymbols))
	for _, symbol := range assetSymbols {
		asset, err := w.parent.Asset(ctx, symbol)
		if err != nil {
			return nil, err
		}
		ids = append(ids, asset.ID.String())
	}

	var reply []protocol.AssetAmount
	if err := w.parent.chain.GetAccountBalances(ctx, w.account.ID.String(), ids, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

// Memo builds an encrypted memo for the target account.
func (w *Wallet) Memo(ctx context.Context, toName, message string) (*MemoData, error) {
	if w == nil || w.parent == nil || w.memoKey == nil {
		return nil, fmt.Errorf("memo key is not configured")
	}
	toAccount, err := w.parent.Account(ctx, toName)
	if err != nil {
		return nil, err
	}
	if toAccount.Options.MemoKey == "" {
		return nil, fmt.Errorf("target account has no memo key")
	}
	toPub, err := ecc.PublicKeyFromString(toAccount.Options.MemoKey)
	if err != nil {
		return nil, err
	}

	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	enc, err := ecc.EncryptWithChecksum(w.memoKey, toPub, nonce, []byte(message))
	if err != nil {
		return nil, err
	}

	return &MemoData{
		From:    w.memoKey.PublicKey().String(),
		To:      toPub.String(),
		Nonce:   nonce,
		Message: hex.EncodeToString(enc),
	}, nil
}

// MemoDecode decrypts a memo payload for the wallet memo key.
func (w *Wallet) MemoDecode(memo *MemoData) (string, error) {
	if w == nil || w.memoKey == nil {
		return "", fmt.Errorf("memo key is not configured")
	}
	if memo == nil {
		return "", fmt.Errorf("nil memo")
	}

	fromPub, err := ecc.PublicKeyFromString(memo.From)
	if err != nil {
		return "", err
	}
	toPub, err := ecc.PublicKeyFromString(memo.To)
	if err != nil {
		return "", err
	}
	raw, err := hex.DecodeString(strings.TrimSpace(memo.Message))
	if err != nil {
		return "", err
	}

	// Try the recipient side first, then fall back to the sender-side variant.
	if out, err := ecc.DecryptWithChecksum(w.memoKey, fromPub, memo.Nonce, raw, false); err == nil {
		return string(out), nil
	}
	out, err := ecc.DecryptWithChecksum(w.memoKey, toPub, memo.Nonce, raw, false)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// BuildTransferOperation creates a typed transfer operation with required fee filled in.
func (w *Wallet) BuildTransferOperation(ctx context.Context, toName, assetSymbol string, amount float64, memo any) (*protocol.TransferOperation, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}

	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}
	toAccount, err := w.parent.Account(ctx, toName)
	if err != nil {
		return nil, err
	}

	op := &protocol.TransferOperation{
		From:       w.account.ID,
		To:         toAccount.ID,
		Amount:     protocol.AssetAmount{Amount: mustAmountFromFloat(amount, asset.Precision), AssetID: asset.ID},
		Extensions: []json.RawMessage{},
	}
	if op.Amount.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	if memo != nil {
		rawMemo, err := w.memoPayload(ctx, toName, memo)
		if err != nil {
			return nil, err
		}
		op.Memo = rawMemo
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// Transfer is a convenience alias for BuildTransferOperation.
func (w *Wallet) Transfer(ctx context.Context, toName, assetSymbol string, amount float64, memo any) (*protocol.TransferOperation, error) {
	return w.BuildTransferOperation(ctx, toName, assetSymbol, amount, memo)
}

// BuildAssetIssueOperation creates a typed asset issue operation.
func (w *Wallet) BuildAssetIssueOperation(ctx context.Context, toName, assetSymbol string, amount float64, memo any) (*protocol.AssetIssueOperation, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}

	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}
	toAccount, err := w.parent.Account(ctx, toName)
	if err != nil {
		return nil, err
	}

	op := &protocol.AssetIssueOperation{
		Issuer:         w.account.ID,
		AssetToIssue:   protocol.AssetAmount{Amount: mustAmountFromFloat(amount, asset.Precision), AssetID: asset.ID},
		IssueToAccount: toAccount.ID,
		Extensions:     []json.RawMessage{},
	}
	if op.AssetToIssue.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	if memo != nil {
		rawMemo, err := w.memoPayload(ctx, toName, memo)
		if err != nil {
			return nil, err
		}
		op.Memo = rawMemo
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// AssetIssue is a convenience alias for BuildAssetIssueOperation.
func (w *Wallet) AssetIssue(ctx context.Context, toName, assetSymbol string, amount float64, memo any) (*protocol.AssetIssueOperation, error) {
	return w.BuildAssetIssueOperation(ctx, toName, assetSymbol, amount, memo)
}

// BuildAssetReserveOperation creates a typed reserve operation.
func (w *Wallet) BuildAssetReserveOperation(ctx context.Context, assetSymbol string, amount float64) (*protocol.AssetReserveOperation, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}

	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}

	op := &protocol.AssetReserveOperation{
		Payer:           w.account.ID,
		AmountToReserve: protocol.AssetAmount{Amount: mustAmountFromFloat(amount, asset.Precision), AssetID: asset.ID},
		Extensions:      []json.RawMessage{},
	}
	if op.AmountToReserve.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// AssetReserve is a convenience alias for BuildAssetReserveOperation.
func (w *Wallet) AssetReserve(ctx context.Context, assetSymbol string, amount float64) (*protocol.AssetReserveOperation, error) {
	return w.BuildAssetReserveOperation(ctx, assetSymbol, amount)
}

// BuildAccountCreateOperation creates a typed account create operation using the
// current wallet account as registrar and, by default, as referrer.
func (w *Wallet) BuildAccountCreateOperation(ctx context.Context, name string, owner protocol.Authority, active protocol.Authority, options protocol.AccountOptions, referrerName string, referrerPercent uint16, extensions *protocol.AccountCreateExtensions) (*protocol.AccountCreateOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}

	referrerID := w.account.ID
	if strings.TrimSpace(referrerName) != "" {
		referrer, err := w.parent.Account(ctx, referrerName)
		if err != nil {
			return nil, err
		}
		referrerID = referrer.ID
	}

	op, err := buildAccountCreateOperation(w.account.ID, referrerID, referrerPercent, name, owner, active, options, extensions)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// AccountCreate is a convenience alias for BuildAccountCreateOperation.
func (w *Wallet) AccountCreate(ctx context.Context, name string, owner protocol.Authority, active protocol.Authority, options protocol.AccountOptions, referrerName string, referrerPercent uint16, extensions *protocol.AccountCreateExtensions) (*protocol.AccountCreateOperation, error) {
	return w.BuildAccountCreateOperation(ctx, name, owner, active, options, referrerName, referrerPercent, extensions)
}

// BuildAccountUpdateOperation creates a typed account update operation.
func (w *Wallet) BuildAccountUpdateOperation(ctx context.Context, owner *protocol.Authority, active *protocol.Authority, newOptions *protocol.AccountOptions, extensions *protocol.AccountUpdateExtensions) (*protocol.AccountUpdateOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}

	op, err := buildAccountUpdateOperation(w.account.ID, owner, active, newOptions, extensions)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// AccountUpdate is a convenience alias for BuildAccountUpdateOperation.
func (w *Wallet) AccountUpdate(ctx context.Context, owner *protocol.Authority, active *protocol.Authority, newOptions *protocol.AccountOptions, extensions *protocol.AccountUpdateExtensions) (*protocol.AccountUpdateOperation, error) {
	return w.BuildAccountUpdateOperation(ctx, owner, active, newOptions, extensions)
}

// BuildAccountUpgradeOperation creates a typed account upgrade operation.
func (w *Wallet) BuildAccountUpgradeOperation(ctx context.Context, upgradeToLifetimeMember bool) (*protocol.AccountUpgradeOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}

	op := buildAccountUpgradeOperation(w.account.ID, upgradeToLifetimeMember)
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// AccountUpgrade is a convenience alias for BuildAccountUpgradeOperation.
func (w *Wallet) AccountUpgrade(ctx context.Context, upgradeToLifetimeMember bool) (*protocol.AccountUpgradeOperation, error) {
	return w.BuildAccountUpgradeOperation(ctx, upgradeToLifetimeMember)
}

// BuildCallOrderUpdateOperation creates a typed call-order update operation.
func (w *Wallet) BuildCallOrderUpdateOperation(ctx context.Context, debtSymbol string, deltaDebt float64, collateralSymbol string, deltaCollateral float64, targetCollateralRatio *uint16) (*protocol.CallOrderUpdateOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	debtAsset, err := w.parent.Asset(ctx, debtSymbol)
	if err != nil {
		return nil, err
	}
	collateralAsset, err := w.parent.Asset(ctx, collateralSymbol)
	if err != nil {
		return nil, err
	}

	op, err := buildCallOrderUpdateOperation(
		w.account.ID,
		debtAsset.ID,
		debtAsset.Precision,
		deltaDebt,
		collateralAsset.ID,
		collateralAsset.Precision,
		deltaCollateral,
		targetCollateralRatio,
	)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// CallOrderUpdate is a convenience alias for BuildCallOrderUpdateOperation.
func (w *Wallet) CallOrderUpdate(ctx context.Context, debtSymbol string, deltaDebt float64, collateralSymbol string, deltaCollateral float64, targetCollateralRatio *uint16) (*protocol.CallOrderUpdateOperation, error) {
	return w.BuildCallOrderUpdateOperation(ctx, debtSymbol, deltaDebt, collateralSymbol, deltaCollateral, targetCollateralRatio)
}

// BuildAssetFundFeePoolOperation creates a typed fee-pool funding operation using the chain core asset.
func (w *Wallet) BuildAssetFundFeePoolOperation(ctx context.Context, assetSymbol string, amount float64) (*protocol.AssetFundFeePoolOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}
	coreAsset, err := w.parent.Asset(ctx, w.parent.Chain().CoreAsset)
	if err != nil {
		return nil, err
	}

	op, err := buildAssetFundFeePoolOperation(w.account.ID, asset.ID, coreAsset.Precision, amount)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// AssetFundFeePool is a convenience alias for BuildAssetFundFeePoolOperation.
func (w *Wallet) AssetFundFeePool(ctx context.Context, assetSymbol string, amount float64) (*protocol.AssetFundFeePoolOperation, error) {
	return w.BuildAssetFundFeePoolOperation(ctx, assetSymbol, amount)
}

// BuildAssetUpdateFeedProducersOperation creates a typed feed-producer update operation.
func (w *Wallet) BuildAssetUpdateFeedProducersOperation(ctx context.Context, assetSymbol string, producerNames []string) (*protocol.AssetUpdateFeedProducersOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}

	producerIDs := make([]protocol.ObjectID, 0, len(producerNames))
	for _, name := range producerNames {
		account, err := w.parent.Account(ctx, name)
		if err != nil {
			return nil, err
		}
		producerIDs = append(producerIDs, account.ID)
	}

	op, err := buildAssetUpdateFeedProducersOperation(w.account.ID, asset.ID, producerIDs)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// AssetUpdateFeedProducers is a convenience alias for BuildAssetUpdateFeedProducersOperation.
func (w *Wallet) AssetUpdateFeedProducers(ctx context.Context, assetSymbol string, producerNames []string) (*protocol.AssetUpdateFeedProducersOperation, error) {
	return w.BuildAssetUpdateFeedProducersOperation(ctx, assetSymbol, producerNames)
}

// BuildAssetUpdateOperation creates a typed asset update operation.
func (w *Wallet) BuildAssetUpdateOperation(ctx context.Context, assetSymbol string, newOptions protocol.AssetOptions, newIssuerName string, extensions *protocol.AssetUpdateExtensions) (*protocol.AssetUpdateOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}

	var newIssuerID *protocol.ObjectID
	if strings.TrimSpace(newIssuerName) != "" {
		account, err := w.parent.Account(ctx, newIssuerName)
		if err != nil {
			return nil, err
		}
		id := account.ID
		newIssuerID = &id
	}

	op := buildAssetUpdateOperation(w.account.ID, asset.ID, newOptions, newIssuerID, extensions)
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// AssetUpdate is a convenience alias for BuildAssetUpdateOperation.
func (w *Wallet) AssetUpdate(ctx context.Context, assetSymbol string, newOptions protocol.AssetOptions, newIssuerName string, extensions *protocol.AssetUpdateExtensions) (*protocol.AssetUpdateOperation, error) {
	return w.BuildAssetUpdateOperation(ctx, assetSymbol, newOptions, newIssuerName, extensions)
}

// BuildAssetUpdateBitassetOperation creates a typed bitasset update operation.
func (w *Wallet) BuildAssetUpdateBitassetOperation(ctx context.Context, assetSymbol string, newOptions protocol.BitAssetOptions) (*protocol.AssetUpdateBitassetOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}

	op := buildAssetUpdateBitassetOperation(w.account.ID, asset.ID, newOptions)
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// AssetUpdateBitasset is a convenience alias for BuildAssetUpdateBitassetOperation.
func (w *Wallet) AssetUpdateBitasset(ctx context.Context, assetSymbol string, newOptions protocol.BitAssetOptions) (*protocol.AssetUpdateBitassetOperation, error) {
	return w.BuildAssetUpdateBitassetOperation(ctx, assetSymbol, newOptions)
}

// BuildAssetCreateOperation creates a typed asset create operation after local
// validation of symbol, precision, and asset option invariants.
func (w *Wallet) BuildAssetCreateOperation(ctx context.Context, symbol string, precision uint8, commonOptions protocol.AssetOptions, bitassetOpts *protocol.BitAssetOptions, isPredictionMarket bool) (*protocol.AssetCreateOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}

	op, err := buildAssetCreateOperation(w.account.ID, symbol, precision, commonOptions, bitassetOpts, isPredictionMarket)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// AssetCreate is a convenience alias for BuildAssetCreateOperation.
func (w *Wallet) AssetCreate(ctx context.Context, symbol string, precision uint8, commonOptions protocol.AssetOptions, bitassetOpts *protocol.BitAssetOptions, isPredictionMarket bool) (*protocol.AssetCreateOperation, error) {
	return w.BuildAssetCreateOperation(ctx, symbol, precision, commonOptions, bitassetOpts, isPredictionMarket)
}

// BuildAssetGlobalSettleOperation creates a typed asset global settle operation.
func (w *Wallet) BuildAssetGlobalSettleOperation(ctx context.Context, assetSymbol string, settlePrice protocol.Price) (*protocol.AssetGlobalSettleOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}

	op, err := buildAssetGlobalSettleOperation(w.account.ID, asset.ID, settlePrice)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// AssetGlobalSettle is a convenience alias for BuildAssetGlobalSettleOperation.
func (w *Wallet) AssetGlobalSettle(ctx context.Context, assetSymbol string, settlePrice protocol.Price) (*protocol.AssetGlobalSettleOperation, error) {
	return w.BuildAssetGlobalSettleOperation(ctx, assetSymbol, settlePrice)
}

// BuildLimitOrderCreateOperation creates a typed limit order create operation.
func (w *Wallet) BuildLimitOrderCreateOperation(ctx context.Context, sellSymbol string, sellAmount float64, buySymbol string, buyAmount float64, fillOrKill bool, expiration time.Time) (*protocol.LimitOrderCreateOperation, error) {
	if sellAmount <= 0 || buyAmount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}

	sellAsset, err := w.parent.Asset(ctx, sellSymbol)
	if err != nil {
		return nil, err
	}
	buyAsset, err := w.parent.Asset(ctx, buySymbol)
	if err != nil {
		return nil, err
	}
	if expiration.IsZero() {
		expiration = time.Now().UTC().AddDate(5, 0, 0)
	}

	op := &protocol.LimitOrderCreateOperation{
		Seller:       w.account.ID,
		AmountToSell: protocol.AssetAmount{Amount: mustAmountFromFloat(sellAmount, sellAsset.Precision), AssetID: sellAsset.ID},
		MinToReceive: protocol.AssetAmount{Amount: mustAmountFromFloat(buyAmount, buyAsset.Precision), AssetID: buyAsset.ID},
		Expiration:   protocol.Time{Time: expiration},
		FillOrKill:   fillOrKill,
		Extensions:   protocol.LimitOrderCreateExtensions{},
	}
	if op.AmountToSell.Amount <= 0 || op.MinToReceive.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// LimitOrderCreate is a convenience alias for BuildLimitOrderCreateOperation.
func (w *Wallet) LimitOrderCreate(ctx context.Context, sellSymbol string, sellAmount float64, buySymbol string, buyAmount float64, fillOrKill bool, expiration time.Time) (*protocol.LimitOrderCreateOperation, error) {
	return w.BuildLimitOrderCreateOperation(ctx, sellSymbol, sellAmount, buySymbol, buyAmount, fillOrKill, expiration)
}

// BuildBuyOperation creates a limit order that buys the requested asset.
func (w *Wallet) BuildBuyOperation(ctx context.Context, buySymbol, baseSymbol string, amount, price float64, fillOrKill bool, expiration time.Time) (*protocol.LimitOrderCreateOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	buyAsset, err := w.parent.Asset(ctx, buySymbol)
	if err != nil {
		return nil, err
	}
	baseAsset, err := w.parent.Asset(ctx, baseSymbol)
	if err != nil {
		return nil, err
	}
	buyAmount := mustAmountFromFloat(amount, buyAsset.Precision)
	sellAmount, err := multiplyAmountPrice(amount, price, baseAsset.Precision)
	if err != nil {
		return nil, err
	}
	if buyAmount <= 0 || sellAmount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	return w.buildOrder(ctx, sellAmount, baseAsset, buyAmount, buyAsset, fillOrKill, expiration)
}

// Buy is a convenience alias for BuildBuyOperation.
func (w *Wallet) Buy(ctx context.Context, buySymbol, baseSymbol string, amount, price float64, fillOrKill bool, expiration time.Time) (*protocol.LimitOrderCreateOperation, error) {
	return w.BuildBuyOperation(ctx, buySymbol, baseSymbol, amount, price, fillOrKill, expiration)
}

// BuildSellOperation creates a limit order that sells the requested asset.
func (w *Wallet) BuildSellOperation(ctx context.Context, sellSymbol, baseSymbol string, amount, price float64, fillOrKill bool, expiration time.Time) (*protocol.LimitOrderCreateOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	sellAsset, err := w.parent.Asset(ctx, sellSymbol)
	if err != nil {
		return nil, err
	}
	baseAsset, err := w.parent.Asset(ctx, baseSymbol)
	if err != nil {
		return nil, err
	}
	sellAmount := mustAmountFromFloat(amount, sellAsset.Precision)
	buyAmount, err := multiplyAmountPrice(amount, price, baseAsset.Precision)
	if err != nil {
		return nil, err
	}
	if buyAmount <= 0 || sellAmount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	return w.buildOrder(ctx, sellAmount, sellAsset, buyAmount, baseAsset, fillOrKill, expiration)
}

// Sell is a convenience alias for BuildSellOperation.
func (w *Wallet) Sell(ctx context.Context, sellSymbol, baseSymbol string, amount, price float64, fillOrKill bool, expiration time.Time) (*protocol.LimitOrderCreateOperation, error) {
	return w.BuildSellOperation(ctx, sellSymbol, baseSymbol, amount, price, fillOrKill, expiration)
}

// BuildHTLCCreateOperation creates a typed HTLC create operation.
func (w *Wallet) BuildHTLCCreateOperation(ctx context.Context, toName, assetSymbol string, amount float64, preimageHash protocol.HTLCPreimageHash, preimageSize uint16, claimPeriodSeconds uint32, memo any) (*protocol.HTLCCreateOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}
	toAccount, err := w.parent.Account(ctx, toName)
	if err != nil {
		return nil, err
	}

	var rawMemo json.RawMessage
	if memo != nil {
		rawMemo, err = w.memoPayload(ctx, toName, memo)
		if err != nil {
			return nil, err
		}
	}

	op, err := buildHTLCCreateOperation(
		w.account.ID,
		toAccount.ID,
		asset.ID,
		asset.Precision,
		amount,
		preimageHash,
		preimageSize,
		claimPeriodSeconds,
		rawMemo,
	)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// HTLCCreate is a convenience alias for BuildHTLCCreateOperation.
func (w *Wallet) HTLCCreate(ctx context.Context, toName, assetSymbol string, amount float64, preimageHash protocol.HTLCPreimageHash, preimageSize uint16, claimPeriodSeconds uint32, memo any) (*protocol.HTLCCreateOperation, error) {
	return w.BuildHTLCCreateOperation(ctx, toName, assetSymbol, amount, preimageHash, preimageSize, claimPeriodSeconds, memo)
}

// BuildWitnessCreateOperation creates a typed witness create operation.
func (w *Wallet) BuildWitnessCreateOperation(ctx context.Context, url string, blockSigningKey string) (*protocol.WitnessCreateOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	signingKey, err := protocol.ParsePublicKey(blockSigningKey)
	if err != nil {
		return nil, err
	}

	op, err := buildWitnessCreateOperation(w.account.ID, url, signingKey)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// WitnessCreate is a convenience alias for BuildWitnessCreateOperation.
func (w *Wallet) WitnessCreate(ctx context.Context, url string, blockSigningKey string) (*protocol.WitnessCreateOperation, error) {
	return w.BuildWitnessCreateOperation(ctx, url, blockSigningKey)
}

// BuildWitnessUpdateOperation creates a typed witness update operation.
func (w *Wallet) BuildWitnessUpdateOperation(ctx context.Context, witnessID string, newURL string, newSigningKey string) (*protocol.WitnessUpdateOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	witness, err := protocol.ParseObjectID(witnessID)
	if err != nil {
		return nil, err
	}

	var signingKey *protocol.PublicKey
	if strings.TrimSpace(newSigningKey) != "" {
		key, err := protocol.ParsePublicKey(newSigningKey)
		if err != nil {
			return nil, err
		}
		signingKey = &key
	}

	op, err := buildWitnessUpdateOperation(witness, w.account.ID, newURL, signingKey)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// WitnessUpdate is a convenience alias for BuildWitnessUpdateOperation.
func (w *Wallet) WitnessUpdate(ctx context.Context, witnessID string, newURL string, newSigningKey string) (*protocol.WitnessUpdateOperation, error) {
	return w.BuildWitnessUpdateOperation(ctx, witnessID, newURL, newSigningKey)
}

// BuildAssetPublishFeedOperation creates a typed asset publish feed operation.
func (w *Wallet) BuildAssetPublishFeedOperation(ctx context.Context, assetSymbol string, feed protocol.PriceFeed, extensions *protocol.AssetPublishFeedExtensions) (*protocol.AssetPublishFeedOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}

	op := buildAssetPublishFeedOperation(w.account.ID, asset.ID, feed, extensions)
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// AssetPublishFeed is a convenience alias for BuildAssetPublishFeedOperation.
func (w *Wallet) AssetPublishFeed(ctx context.Context, assetSymbol string, feed protocol.PriceFeed, extensions *protocol.AssetPublishFeedExtensions) (*protocol.AssetPublishFeedOperation, error) {
	return w.BuildAssetPublishFeedOperation(ctx, assetSymbol, feed, extensions)
}

// BuildAssetSettleOperation creates a typed asset settle operation.
func (w *Wallet) BuildAssetSettleOperation(ctx context.Context, assetSymbol string, amount float64) (*protocol.AssetSettleOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}

	op, err := buildAssetSettleOperation(w.account.ID, asset.ID, asset.Precision, amount)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// AssetSettle is a convenience alias for BuildAssetSettleOperation.
func (w *Wallet) AssetSettle(ctx context.Context, assetSymbol string, amount float64) (*protocol.AssetSettleOperation, error) {
	return w.BuildAssetSettleOperation(ctx, assetSymbol, amount)
}

// BuildAccountWhitelistOperation creates a typed account whitelist operation.
func (w *Wallet) BuildAccountWhitelistOperation(ctx context.Context, accountName string, newListing uint8) (*protocol.AccountWhitelistOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	account, err := w.parent.Account(ctx, accountName)
	if err != nil {
		return nil, err
	}

	op, err := buildAccountWhitelistOperation(w.account.ID, account.ID, newListing)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// AccountWhitelist is a convenience alias for BuildAccountWhitelistOperation.
func (w *Wallet) AccountWhitelist(ctx context.Context, accountName string, newListing uint8) (*protocol.AccountWhitelistOperation, error) {
	return w.BuildAccountWhitelistOperation(ctx, accountName, newListing)
}

// BuildProposalCreateOperation creates a typed proposal create operation from
// existing typed operations. A zero expiration uses the chain proposal default.
func (w *Wallet) BuildProposalCreateOperation(ctx context.Context, expiration time.Time, reviewPeriodSeconds *uint32, ops ...protocol.Operation) (*protocol.ProposalCreateOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}

	if expiration.IsZero() {
		var props struct {
			Time string `json:"time"`
		}
		if err := w.parent.chain.GetDynamicGlobalProperties(ctx, &props); err != nil {
			return nil, err
		}
		expiration = buildTransactionExpiration(props.Time, w.parent.Chain().ExpireInSecsProposal).Time
	}

	op, err := buildProposalCreateOperation(w.account.ID, expiration, reviewPeriodSeconds, ops...)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// ProposalCreate is a convenience alias for BuildProposalCreateOperation.
func (w *Wallet) ProposalCreate(ctx context.Context, expiration time.Time, reviewPeriodSeconds *uint32, ops ...protocol.Operation) (*protocol.ProposalCreateOperation, error) {
	return w.BuildProposalCreateOperation(ctx, expiration, reviewPeriodSeconds, ops...)
}

// BuildProposalUpdateOperation creates a typed proposal update operation.
func (w *Wallet) BuildProposalUpdateOperation(ctx context.Context, proposalID string, activeApprovalsToAdd []string, activeApprovalsToRemove []string, ownerApprovalsToAdd []string, ownerApprovalsToRemove []string, keyApprovalsToAdd []string, keyApprovalsToRemove []string) (*protocol.ProposalUpdateOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	proposal, err := protocol.ParseObjectID(proposalID)
	if err != nil {
		return nil, err
	}

	activeAdd, err := w.resolveAccountIDs(ctx, activeApprovalsToAdd)
	if err != nil {
		return nil, err
	}
	activeRemove, err := w.resolveAccountIDs(ctx, activeApprovalsToRemove)
	if err != nil {
		return nil, err
	}
	ownerAdd, err := w.resolveAccountIDs(ctx, ownerApprovalsToAdd)
	if err != nil {
		return nil, err
	}
	ownerRemove, err := w.resolveAccountIDs(ctx, ownerApprovalsToRemove)
	if err != nil {
		return nil, err
	}
	keyAdd, err := parsePublicKeys(keyApprovalsToAdd)
	if err != nil {
		return nil, err
	}
	keyRemove, err := parsePublicKeys(keyApprovalsToRemove)
	if err != nil {
		return nil, err
	}

	op, err := buildProposalUpdateOperation(w.account.ID, proposal, activeAdd, activeRemove, ownerAdd, ownerRemove, keyAdd, keyRemove)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// ProposalUpdate is a convenience alias for BuildProposalUpdateOperation.
func (w *Wallet) ProposalUpdate(ctx context.Context, proposalID string, activeApprovalsToAdd []string, activeApprovalsToRemove []string, ownerApprovalsToAdd []string, ownerApprovalsToRemove []string, keyApprovalsToAdd []string, keyApprovalsToRemove []string) (*protocol.ProposalUpdateOperation, error) {
	return w.BuildProposalUpdateOperation(ctx, proposalID, activeApprovalsToAdd, activeApprovalsToRemove, ownerApprovalsToAdd, ownerApprovalsToRemove, keyApprovalsToAdd, keyApprovalsToRemove)
}

// BuildProposalDeleteOperation creates a typed proposal delete operation.
func (w *Wallet) BuildProposalDeleteOperation(ctx context.Context, proposalID string, usingOwnerAuthority bool) (*protocol.ProposalDeleteOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	proposal, err := protocol.ParseObjectID(proposalID)
	if err != nil {
		return nil, err
	}

	op := buildProposalDeleteOperation(w.account.ID, proposal, usingOwnerAuthority)
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// ProposalDelete is a convenience alias for BuildProposalDeleteOperation.
func (w *Wallet) ProposalDelete(ctx context.Context, proposalID string, usingOwnerAuthority bool) (*protocol.ProposalDeleteOperation, error) {
	return w.BuildProposalDeleteOperation(ctx, proposalID, usingOwnerAuthority)
}

// BuildWithdrawPermissionCreateOperation creates a typed withdraw-permission create operation.
func (w *Wallet) BuildWithdrawPermissionCreateOperation(ctx context.Context, authorizedAccountName string, assetSymbol string, amount float64, withdrawalPeriodSec uint32, periodsUntilExpiration uint32, periodStartTime time.Time) (*protocol.WithdrawPermissionCreateOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	authorizedAccount, err := w.parent.Account(ctx, authorizedAccountName)
	if err != nil {
		return nil, err
	}
	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}

	op, err := buildWithdrawPermissionCreateOperation(
		w.account.ID,
		authorizedAccount.ID,
		asset.ID,
		asset.Precision,
		amount,
		withdrawalPeriodSec,
		periodsUntilExpiration,
		periodStartTime,
	)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// WithdrawPermissionCreate is a convenience alias for BuildWithdrawPermissionCreateOperation.
func (w *Wallet) WithdrawPermissionCreate(ctx context.Context, authorizedAccountName string, assetSymbol string, amount float64, withdrawalPeriodSec uint32, periodsUntilExpiration uint32, periodStartTime time.Time) (*protocol.WithdrawPermissionCreateOperation, error) {
	return w.BuildWithdrawPermissionCreateOperation(ctx, authorizedAccountName, assetSymbol, amount, withdrawalPeriodSec, periodsUntilExpiration, periodStartTime)
}

// BuildWithdrawPermissionUpdateOperation creates a typed withdraw-permission update operation.
func (w *Wallet) BuildWithdrawPermissionUpdateOperation(ctx context.Context, permissionID string, authorizedAccountName string, assetSymbol string, amount float64, withdrawalPeriodSec uint32, periodsUntilExpiration uint32, periodStartTime time.Time) (*protocol.WithdrawPermissionUpdateOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	permission, err := protocol.ParseObjectID(permissionID)
	if err != nil {
		return nil, err
	}
	authorizedAccount, err := w.parent.Account(ctx, authorizedAccountName)
	if err != nil {
		return nil, err
	}
	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}

	op, err := buildWithdrawPermissionUpdateOperation(
		w.account.ID,
		authorizedAccount.ID,
		permission,
		asset.ID,
		asset.Precision,
		amount,
		withdrawalPeriodSec,
		periodsUntilExpiration,
		periodStartTime,
	)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// WithdrawPermissionUpdate is a convenience alias for BuildWithdrawPermissionUpdateOperation.
func (w *Wallet) WithdrawPermissionUpdate(ctx context.Context, permissionID string, authorizedAccountName string, assetSymbol string, amount float64, withdrawalPeriodSec uint32, periodsUntilExpiration uint32, periodStartTime time.Time) (*protocol.WithdrawPermissionUpdateOperation, error) {
	return w.BuildWithdrawPermissionUpdateOperation(ctx, permissionID, authorizedAccountName, assetSymbol, amount, withdrawalPeriodSec, periodsUntilExpiration, periodStartTime)
}

// BuildWithdrawPermissionClaimOperation creates a typed withdraw-permission claim operation.
func (w *Wallet) BuildWithdrawPermissionClaimOperation(ctx context.Context, permissionID string, withdrawFromAccountName string, assetSymbol string, amount float64, memo any) (*protocol.WithdrawPermissionClaimOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	permission, err := protocol.ParseObjectID(permissionID)
	if err != nil {
		return nil, err
	}
	withdrawFromAccount, err := w.parent.Account(ctx, withdrawFromAccountName)
	if err != nil {
		return nil, err
	}
	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}

	var rawMemo json.RawMessage
	if memo != nil {
		rawMemo, err = w.memoPayload(ctx, withdrawFromAccountName, memo)
		if err != nil {
			return nil, err
		}
	}

	op, err := buildWithdrawPermissionClaimOperation(
		permission,
		withdrawFromAccount.ID,
		w.account.ID,
		asset.ID,
		asset.Precision,
		amount,
		rawMemo,
	)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// WithdrawPermissionClaim is a convenience alias for BuildWithdrawPermissionClaimOperation.
func (w *Wallet) WithdrawPermissionClaim(ctx context.Context, permissionID string, withdrawFromAccountName string, assetSymbol string, amount float64, memo any) (*protocol.WithdrawPermissionClaimOperation, error) {
	return w.BuildWithdrawPermissionClaimOperation(ctx, permissionID, withdrawFromAccountName, assetSymbol, amount, memo)
}

// BuildWithdrawPermissionDeleteOperation creates a typed withdraw-permission delete operation.
func (w *Wallet) BuildWithdrawPermissionDeleteOperation(ctx context.Context, permissionID string, authorizedAccountName string) (*protocol.WithdrawPermissionDeleteOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	permission, err := protocol.ParseObjectID(permissionID)
	if err != nil {
		return nil, err
	}
	authorizedAccount, err := w.parent.Account(ctx, authorizedAccountName)
	if err != nil {
		return nil, err
	}

	op, err := buildWithdrawPermissionDeleteOperation(w.account.ID, authorizedAccount.ID, permission)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// WithdrawPermissionDelete is a convenience alias for BuildWithdrawPermissionDeleteOperation.
func (w *Wallet) WithdrawPermissionDelete(ctx context.Context, permissionID string, authorizedAccountName string) (*protocol.WithdrawPermissionDeleteOperation, error) {
	return w.BuildWithdrawPermissionDeleteOperation(ctx, permissionID, authorizedAccountName)
}

// BuildAccountTransferOperation creates a typed account transfer operation.
func (w *Wallet) BuildAccountTransferOperation(ctx context.Context, newOwnerName string) (*protocol.AccountTransferOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	newOwner, err := w.parent.Account(ctx, newOwnerName)
	if err != nil {
		return nil, err
	}

	op := buildAccountTransferOperation(w.account.ID, newOwner.ID)
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// AccountTransfer is a convenience alias for BuildAccountTransferOperation.
func (w *Wallet) AccountTransfer(ctx context.Context, newOwnerName string) (*protocol.AccountTransferOperation, error) {
	return w.BuildAccountTransferOperation(ctx, newOwnerName)
}

// BuildBalanceClaimOperation creates a typed balance claim operation. Core
// treats this operation as zero-fee, so no fee lookup is required.
func (w *Wallet) BuildBalanceClaimOperation(ctx context.Context, balanceID string, balanceOwnerKey string, assetSymbol string, amount float64) (*protocol.BalanceClaimOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	balance, err := protocol.ParseObjectID(balanceID)
	if err != nil {
		return nil, err
	}
	ownerKey, err := protocol.ParsePublicKey(balanceOwnerKey)
	if err != nil {
		return nil, err
	}
	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}

	return buildBalanceClaimOperation(w.account.ID, balance, ownerKey, asset.ID, asset.Precision, amount)
}

// BalanceClaim is a convenience alias for BuildBalanceClaimOperation.
func (w *Wallet) BalanceClaim(ctx context.Context, balanceID string, balanceOwnerKey string, assetSymbol string, amount float64) (*protocol.BalanceClaimOperation, error) {
	return w.BuildBalanceClaimOperation(ctx, balanceID, balanceOwnerKey, assetSymbol, amount)
}

// BuildOverrideTransferOperation creates a typed override transfer operation.
func (w *Wallet) BuildOverrideTransferOperation(ctx context.Context, fromName string, toName string, assetSymbol string, amount float64, memo any) (*protocol.OverrideTransferOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	fromAccount, err := w.parent.Account(ctx, fromName)
	if err != nil {
		return nil, err
	}
	toAccount, err := w.parent.Account(ctx, toName)
	if err != nil {
		return nil, err
	}
	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}

	var rawMemo json.RawMessage
	if memo != nil {
		rawMemo, err = w.memoPayload(ctx, toName, memo)
		if err != nil {
			return nil, err
		}
	}

	op, err := buildOverrideTransferOperation(w.account.ID, fromAccount.ID, toAccount.ID, asset.ID, asset.Precision, amount, rawMemo)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// OverrideTransfer is a convenience alias for BuildOverrideTransferOperation.
func (w *Wallet) OverrideTransfer(ctx context.Context, fromName string, toName string, assetSymbol string, amount float64, memo any) (*protocol.OverrideTransferOperation, error) {
	return w.BuildOverrideTransferOperation(ctx, fromName, toName, assetSymbol, amount, memo)
}

// BuildBidCollateralOperation creates a typed bid collateral operation.
func (w *Wallet) BuildBidCollateralOperation(ctx context.Context, collateralSymbol string, additionalCollateral float64, debtSymbol string, debtCovered float64) (*protocol.BidCollateralOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	collateralAsset, err := w.parent.Asset(ctx, collateralSymbol)
	if err != nil {
		return nil, err
	}
	debtAsset, err := w.parent.Asset(ctx, debtSymbol)
	if err != nil {
		return nil, err
	}

	op, err := buildBidCollateralOperation(w.account.ID, collateralAsset.ID, collateralAsset.Precision, additionalCollateral, debtAsset.ID, debtAsset.Precision, debtCovered)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// BidCollateral is a convenience alias for BuildBidCollateralOperation.
func (w *Wallet) BidCollateral(ctx context.Context, collateralSymbol string, additionalCollateral float64, debtSymbol string, debtCovered float64) (*protocol.BidCollateralOperation, error) {
	return w.BuildBidCollateralOperation(ctx, collateralSymbol, additionalCollateral, debtSymbol, debtCovered)
}

// BuildCommitteeMemberCreateOperation creates a typed committee member create operation.
func (w *Wallet) BuildCommitteeMemberCreateOperation(ctx context.Context, url string) (*protocol.CommitteeMemberCreateOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}

	op, err := buildCommitteeMemberCreateOperation(w.account.ID, url)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// CommitteeMemberCreate is a convenience alias for BuildCommitteeMemberCreateOperation.
func (w *Wallet) CommitteeMemberCreate(ctx context.Context, url string) (*protocol.CommitteeMemberCreateOperation, error) {
	return w.BuildCommitteeMemberCreateOperation(ctx, url)
}

// BuildCommitteeMemberUpdateOperation creates a typed committee member update operation.
func (w *Wallet) BuildCommitteeMemberUpdateOperation(ctx context.Context, committeeMemberID string, newURL string) (*protocol.CommitteeMemberUpdateOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	member, err := protocol.ParseObjectID(committeeMemberID)
	if err != nil {
		return nil, err
	}

	op, err := buildCommitteeMemberUpdateOperation(member, w.account.ID, newURL)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// CommitteeMemberUpdate is a convenience alias for BuildCommitteeMemberUpdateOperation.
func (w *Wallet) CommitteeMemberUpdate(ctx context.Context, committeeMemberID string, newURL string) (*protocol.CommitteeMemberUpdateOperation, error) {
	return w.BuildCommitteeMemberUpdateOperation(ctx, committeeMemberID, newURL)
}

// BuildCommitteeMemberUpdateGlobalParametersOperation creates a typed committee
// global-parameters update operation. In practice this operation is meant to be
// wrapped in a proposal before broadcast.
func (w *Wallet) BuildCommitteeMemberUpdateGlobalParametersOperation(ctx context.Context, newParameters protocol.ChainParameters) (*protocol.CommitteeMemberUpdateGlobalParametersOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}

	op := buildCommitteeMemberUpdateGlobalParametersOperation(newParameters)
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// CommitteeMemberUpdateGlobalParameters is a convenience alias for BuildCommitteeMemberUpdateGlobalParametersOperation.
func (w *Wallet) CommitteeMemberUpdateGlobalParameters(ctx context.Context, newParameters protocol.ChainParameters) (*protocol.CommitteeMemberUpdateGlobalParametersOperation, error) {
	return w.BuildCommitteeMemberUpdateGlobalParametersOperation(ctx, newParameters)
}

// BuildAssetClaimFeesOperation creates a typed asset claim fees operation.
func (w *Wallet) BuildAssetClaimFeesOperation(ctx context.Context, assetSymbol string, amount float64, claimFromAssetSymbol string) (*protocol.AssetClaimFeesOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	asset, err := w.parent.Asset(ctx, assetSymbol)
	if err != nil {
		return nil, err
	}

	var claimFromAssetID *protocol.ObjectID
	if strings.TrimSpace(claimFromAssetSymbol) != "" {
		claimAsset, err := w.parent.Asset(ctx, claimFromAssetSymbol)
		if err != nil {
			return nil, err
		}
		id := claimAsset.ID
		claimFromAssetID = &id
	}

	op, err := buildAssetClaimFeesOperation(w.account.ID, asset.ID, asset.Precision, amount, claimFromAssetID)
	if err != nil {
		return nil, err
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// AssetClaimFees is a convenience alias for BuildAssetClaimFeesOperation.
func (w *Wallet) AssetClaimFees(ctx context.Context, assetSymbol string, amount float64, claimFromAssetSymbol string) (*protocol.AssetClaimFeesOperation, error) {
	return w.BuildAssetClaimFeesOperation(ctx, assetSymbol, amount, claimFromAssetSymbol)
}

// BuildCancelOrderOperation creates a typed order cancel operation.
func (w *Wallet) BuildCancelOrderOperation(ctx context.Context, orderID string) (*protocol.LimitOrderCancelOperation, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	order, err := protocol.ParseObjectID(orderID)
	if err != nil {
		return nil, err
	}
	op := &protocol.LimitOrderCancelOperation{
		Order:            order,
		FeePayingAccount: w.account.ID,
		Extensions:       []json.RawMessage{},
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

// CancelOrder is a convenience alias for BuildCancelOrderOperation.
func (w *Wallet) CancelOrder(ctx context.Context, orderID string) (*protocol.LimitOrderCancelOperation, error) {
	return w.BuildCancelOrderOperation(ctx, orderID)
}

// BuildTransaction assembles a BitShares transaction envelope ready for signing.
func (w *Wallet) BuildTransaction(ctx context.Context, ops ...protocol.Operation) (*protocol.Transaction, error) {
	if w == nil || w.parent == nil || w.account == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	if len(ops) == 0 {
		return nil, fmt.Errorf("no operations provided")
	}

	var props struct {
		HeadBlockNumber uint32 `json:"head_block_number"`
		HeadBlockID     string `json:"head_block_id"`
		Time            string `json:"time"`
	}
	if err := w.parent.chain.GetDynamicGlobalProperties(ctx, &props); err != nil {
		return nil, err
	}

	headID := strings.TrimSpace(props.HeadBlockID)
	headBytes, err := hex.DecodeString(headID)
	if err != nil {
		return nil, err
	}
	if len(headBytes) < 8 {
		return nil, fmt.Errorf("invalid head block id")
	}

	tx := &protocol.Transaction{
		RefBlockNum:    uint16(props.HeadBlockNumber & 0xffff),
		RefBlockPrefix: binary.LittleEndian.Uint32(headBytes[4:8]),
		Expiration:     buildTransactionExpiration(props.Time, w.parent.Chain().ExpireInSecs),
		Extensions:     []json.RawMessage{},
	}
	for _, op := range ops {
		tx.Push(op)
	}
	return tx, nil
}

// SignTransaction signs a transaction using the wallet active key unless additional keys are supplied.
func (w *Wallet) SignTransaction(tx *protocol.Transaction, keys ...*ecc.PrivateKey) (*protocol.SignedTransaction, error) {
	if w == nil || w.parent == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}
	if tx == nil {
		return nil, fmt.Errorf("nil transaction")
	}
	signers := keys
	if len(signers) == 0 {
		if w.activeKey == nil {
			return nil, fmt.Errorf("active key is not configured")
		}
		signers = []*ecc.PrivateKey{w.activeKey}
	}
	signer := sign.TransactionSigner{
		ChainID: w.parent.Chain().ChainID,
		Keys:    signers,
	}
	return signer.Sign(tx)
}

// BuildSignedTransaction builds and signs a transaction from the supplied operations.
func (w *Wallet) BuildSignedTransaction(ctx context.Context, ops ...protocol.Operation) (*protocol.SignedTransaction, error) {
	tx, err := w.BuildTransaction(ctx, ops...)
	if err != nil {
		return nil, err
	}
	return w.SignTransaction(tx)
}

// BroadcastSignedTransaction broadcasts an already signed transaction.
func (w *Wallet) BroadcastSignedTransaction(ctx context.Context, tx *protocol.SignedTransaction, reply any) error {
	if w == nil || w.parent == nil {
		return fmt.Errorf("wallet is not configured")
	}
	if tx == nil {
		return fmt.Errorf("nil signed transaction")
	}
	return w.parent.chain.BroadcastTransactionSynchronous(ctx, tx, reply)
}

// BroadcastOperations builds, signs, and broadcasts a transaction from the supplied operations.
func (w *Wallet) BroadcastOperations(ctx context.Context, ops ...protocol.Operation) (*protocol.SignedTransaction, error) {
	return w.BroadcastOperationsWithReply(ctx, nil, ops...)
}

// BroadcastOperationsWithReply builds, signs, and broadcasts a transaction, returning the signed payload.
func (w *Wallet) BroadcastOperationsWithReply(ctx context.Context, reply any, ops ...protocol.Operation) (*protocol.SignedTransaction, error) {
	tx, err := w.BuildSignedTransaction(ctx, ops...)
	if err != nil {
		return nil, err
	}
	if err := w.parent.chain.BroadcastTransactionSynchronous(ctx, tx, reply); err != nil {
		return nil, err
	}
	return tx, nil
}

func (w *Wallet) buildOrder(ctx context.Context, sellAmount int64, sellAsset *AssetInfo, buyAmount int64, buyAsset *AssetInfo, fillOrKill bool, expiration time.Time) (*protocol.LimitOrderCreateOperation, error) {
	if expiration.IsZero() {
		expiration = time.Now().UTC().AddDate(5, 0, 0)
	}
	op := &protocol.LimitOrderCreateOperation{
		Seller:       w.account.ID,
		AmountToSell: protocol.AssetAmount{Amount: sellAmount, AssetID: sellAsset.ID},
		MinToReceive: protocol.AssetAmount{Amount: buyAmount, AssetID: buyAsset.ID},
		Expiration:   protocol.Time{Time: expiration},
		FillOrKill:   fillOrKill,
		Extensions:   protocol.LimitOrderCreateExtensions{},
	}
	if err := w.fillRequiredFee(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}

func (w *Wallet) fillRequiredFee(ctx context.Context, op protocol.Operation) error {
	if w == nil || w.parent == nil || w.feeAsset == nil {
		return fmt.Errorf("fee asset is not configured")
	}
	if op == nil {
		return fmt.Errorf("nil operation")
	}

	envelope := []protocol.OperationEnvelope{{Operation: op}}
	var fees []protocol.AssetAmount
	if err := w.parent.chain.GetRequiredFees(ctx, envelope, w.feeAsset.ID.String(), &fees); err != nil {
		return err
	}
	if len(fees) == 0 {
		return fmt.Errorf("no fee returned")
	}
	return setOperationFee(op, fees[0])
}

func (w *Wallet) memoPayload(ctx context.Context, toName string, memo any) (json.RawMessage, error) {
	if memo == nil {
		return nil, nil
	}

	switch v := memo.(type) {
	case string:
		data, err := w.Memo(ctx, toName, v)
		if err != nil {
			return nil, err
		}
		return json.Marshal(data)
	case MemoData:
		return json.Marshal(v)
	case *MemoData:
		return json.Marshal(v)
	default:
		return json.Marshal(memo)
	}
}

func mustAmountFromFloat(amount float64, precision uint8) int64 {
	return RoundAmount(amount, precision)
}

func signedAmountFromFloat(amount float64, precision uint8) (int64, error) {
	if amount == 0 {
		return 0, nil
	}
	return parseAmount(strconv.FormatFloat(amount, 'f', -1, 64), precision, true)
}

func buildAccountUpdateOperation(accountID protocol.ObjectID, owner *protocol.Authority, active *protocol.Authority, newOptions *protocol.AccountOptions, extensions *protocol.AccountUpdateExtensions) (*protocol.AccountUpdateOperation, error) {
	if owner == nil && active == nil && newOptions == nil && extensions == nil {
		return nil, fmt.Errorf("no account update fields provided")
	}
	op := &protocol.AccountUpdateOperation{
		Account: accountID,
		Owner:   owner,
		Active:  active,
	}
	if newOptions != nil {
		op.NewOptions = newOptions
	}
	if extensions != nil {
		op.Extensions = *extensions
	}
	return op, nil
}

func buildAccountCreateOperation(registrarID, referrerID protocol.ObjectID, referrerPercent uint16, name string, owner protocol.Authority, active protocol.Authority, options protocol.AccountOptions, extensions *protocol.AccountCreateExtensions) (*protocol.AccountCreateOperation, error) {
	if !isValidAccountName(name) {
		return nil, fmt.Errorf("invalid account name")
	}
	if referrerPercent > graphene100Percent {
		return nil, fmt.Errorf("referrer percent must be at most %d", graphene100Percent)
	}
	if err := validateCreateAuthority(owner, "owner"); err != nil {
		return nil, err
	}
	if err := validateCreateAuthority(active, "active"); err != nil {
		return nil, err
	}
	if err := validateAccountOptions(options); err != nil {
		return nil, err
	}
	if err := validateAccountCreateExtensions(owner, active, extensions); err != nil {
		return nil, err
	}

	op := &protocol.AccountCreateOperation{
		Registrar:       registrarID,
		Referrer:        referrerID,
		ReferrerPercent: referrerPercent,
		Name:            name,
		Owner:           owner,
		Active:          active,
		Options:         options,
	}
	if extensions != nil {
		op.Extensions = *extensions
	}
	return op, nil
}

func buildAccountUpgradeOperation(accountID protocol.ObjectID, upgradeToLifetimeMember bool) *protocol.AccountUpgradeOperation {
	return &protocol.AccountUpgradeOperation{
		AccountToUpgrade:        accountID,
		UpgradeToLifetimeMember: upgradeToLifetimeMember,
		Extensions:              []json.RawMessage{},
	}
}

func buildCallOrderUpdateOperation(accountID protocol.ObjectID, debtAssetID protocol.ObjectID, debtPrecision uint8, deltaDebt float64, collateralAssetID protocol.ObjectID, collateralPrecision uint8, deltaCollateral float64, targetCollateralRatio *uint16) (*protocol.CallOrderUpdateOperation, error) {
	if deltaDebt == 0 && deltaCollateral == 0 && targetCollateralRatio == nil {
		return nil, fmt.Errorf("no call order update fields provided")
	}

	debtAmount, err := signedAmountFromFloat(deltaDebt, debtPrecision)
	if err != nil {
		return nil, err
	}
	collateralAmount, err := signedAmountFromFloat(deltaCollateral, collateralPrecision)
	if err != nil {
		return nil, err
	}

	return &protocol.CallOrderUpdateOperation{
		FundingAccount:  accountID,
		DeltaDebt:       protocol.AssetAmount{Amount: debtAmount, AssetID: debtAssetID},
		DeltaCollateral: protocol.AssetAmount{Amount: collateralAmount, AssetID: collateralAssetID},
		Extensions:      protocol.CallOrderUpdateExtensions{TargetCollateralRatio: targetCollateralRatio},
	}, nil
}

func buildHTLCCreateOperation(fromID, toID, assetID protocol.ObjectID, precision uint8, amount float64, preimageHash protocol.HTLCPreimageHash, preimageSize uint16, claimPeriodSeconds uint32, memo json.RawMessage) (*protocol.HTLCCreateOperation, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	if preimageSize == 0 {
		return nil, fmt.Errorf("preimage size must be greater than zero")
	}
	if claimPeriodSeconds == 0 {
		return nil, fmt.Errorf("claim period seconds must be greater than zero")
	}

	op := &protocol.HTLCCreateOperation{
		From:               fromID,
		To:                 toID,
		Amount:             protocol.AssetAmount{Amount: mustAmountFromFloat(amount, precision), AssetID: assetID},
		PreimageHash:       preimageHash,
		PreimageSize:       preimageSize,
		ClaimPeriodSeconds: claimPeriodSeconds,
		Extensions:         protocol.HTLCCreateExtensions{Memo: memo},
	}
	if op.Amount.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	return op, nil
}

func buildAssetFundFeePoolOperation(accountID, assetID protocol.ObjectID, corePrecision uint8, amount float64) (*protocol.AssetFundFeePoolOperation, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	op := &protocol.AssetFundFeePoolOperation{
		FromAccount: accountID,
		AssetID:     assetID,
		Amount:      mustAmountFromFloat(amount, corePrecision),
		Extensions:  []json.RawMessage{},
	}
	if op.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	return op, nil
}

func buildAssetUpdateOperation(issuerID, assetID protocol.ObjectID, newOptions protocol.AssetOptions, newIssuerID *protocol.ObjectID, extensions *protocol.AssetUpdateExtensions) *protocol.AssetUpdateOperation {
	op := &protocol.AssetUpdateOperation{
		Issuer:        issuerID,
		AssetToUpdate: assetID,
		NewOptions:    newOptions,
		NewIssuer:     newIssuerID,
	}
	if extensions != nil {
		op.Extensions = *extensions
	}
	return op
}

func buildAssetUpdateBitassetOperation(issuerID, assetID protocol.ObjectID, newOptions protocol.BitAssetOptions) *protocol.AssetUpdateBitassetOperation {
	return &protocol.AssetUpdateBitassetOperation{
		Issuer:        issuerID,
		AssetToUpdate: assetID,
		NewOptions:    newOptions,
		Extensions:    []json.RawMessage{},
	}
}

func buildAssetCreateOperation(issuerID protocol.ObjectID, symbol string, precision uint8, commonOptions protocol.AssetOptions, bitassetOpts *protocol.BitAssetOptions, isPredictionMarket bool) (*protocol.AssetCreateOperation, error) {
	if !isValidAssetSymbol(symbol) {
		return nil, fmt.Errorf("invalid asset symbol")
	}
	if precision > maxAssetPrecision {
		return nil, fmt.Errorf("precision must be at most %d", maxAssetPrecision)
	}
	if err := validateAssetOptions(commonOptions); err != nil {
		return nil, err
	}
	if bitassetOpts != nil {
		if err := validateBitAssetOptions(*bitassetOpts); err != nil {
			return nil, err
		}
	}
	if isPredictionMarket {
		if bitassetOpts == nil {
			return nil, fmt.Errorf("prediction markets require bitasset options")
		}
		if commonOptions.IssuerPermissions&assetIssuerPermissionGlobalSettle == 0 {
			return nil, fmt.Errorf("prediction markets require global settle issuer permission")
		}
		if commonOptions.IssuerPermissions&assetIssuerPermissionDisableBSRMUpdate != 0 {
			return nil, fmt.Errorf("prediction markets cannot disable BSRM updates")
		}
	}

	return &protocol.AssetCreateOperation{
		Issuer:             issuerID,
		Symbol:             symbol,
		Precision:          precision,
		CommonOptions:      commonOptions,
		BitassetOpts:       bitassetOpts,
		IsPredictionMarket: isPredictionMarket,
		Extensions:         []json.RawMessage{},
	}, nil
}

func buildAssetGlobalSettleOperation(issuerID, assetID protocol.ObjectID, settlePrice protocol.Price) (*protocol.AssetGlobalSettleOperation, error) {
	if settlePrice.Base.AssetID != assetID {
		return nil, fmt.Errorf("settle price base asset must match asset to settle")
	}
	return &protocol.AssetGlobalSettleOperation{
		Issuer:        issuerID,
		AssetToSettle: assetID,
		SettlePrice:   settlePrice,
		Extensions:    []json.RawMessage{},
	}, nil
}

func buildAssetUpdateFeedProducersOperation(issuerID, assetID protocol.ObjectID, producerIDs []protocol.ObjectID) (*protocol.AssetUpdateFeedProducersOperation, error) {
	if len(producerIDs) == 0 {
		return nil, fmt.Errorf("at least one feed producer is required")
	}
	out := make([]protocol.ObjectID, len(producerIDs))
	copy(out, producerIDs)
	return &protocol.AssetUpdateFeedProducersOperation{
		Issuer:           issuerID,
		AssetToUpdate:    assetID,
		NewFeedProducers: out,
		Extensions:       []json.RawMessage{},
	}, nil
}

func buildAssetPublishFeedOperation(publisherID, assetID protocol.ObjectID, feed protocol.PriceFeed, extensions *protocol.AssetPublishFeedExtensions) *protocol.AssetPublishFeedOperation {
	op := &protocol.AssetPublishFeedOperation{
		Publisher: publisherID,
		AssetID:   assetID,
		Feed:      feed,
	}
	if extensions != nil {
		op.Extensions = *extensions
	}
	return op
}

func buildWitnessCreateOperation(accountID protocol.ObjectID, url string, blockSigningKey protocol.PublicKey) (*protocol.WitnessCreateOperation, error) {
	if len(url) >= maxChainURLLength {
		return nil, fmt.Errorf("url must be shorter than %d characters", maxChainURLLength)
	}
	return &protocol.WitnessCreateOperation{
		WitnessAccount:  accountID,
		URL:             url,
		BlockSigningKey: blockSigningKey,
	}, nil
}

func buildWitnessUpdateOperation(witnessID, accountID protocol.ObjectID, newURL string, newSigningKey *protocol.PublicKey) (*protocol.WitnessUpdateOperation, error) {
	trimmedURL := strings.TrimSpace(newURL)
	if trimmedURL == "" && newSigningKey == nil {
		return nil, fmt.Errorf("no witness update fields provided")
	}
	if trimmedURL != "" && len(trimmedURL) >= maxChainURLLength {
		return nil, fmt.Errorf("url must be shorter than %d characters", maxChainURLLength)
	}

	op := &protocol.WitnessUpdateOperation{
		Witness:        witnessID,
		WitnessAccount: accountID,
		NewSigningKey:  newSigningKey,
	}
	if trimmedURL != "" {
		op.NewURL = &trimmedURL
	}
	return op, nil
}

func buildAssetSettleOperation(accountID, assetID protocol.ObjectID, precision uint8, amount float64) (*protocol.AssetSettleOperation, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	op := &protocol.AssetSettleOperation{
		Account:    accountID,
		Amount:     protocol.AssetAmount{Amount: mustAmountFromFloat(amount, precision), AssetID: assetID},
		Extensions: []json.RawMessage{},
	}
	if op.Amount.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	return op, nil
}

func buildAccountWhitelistOperation(authorizingAccountID, accountToListID protocol.ObjectID, newListing uint8) (*protocol.AccountWhitelistOperation, error) {
	if newListing >= 0x4 {
		return nil, fmt.Errorf("new listing must be between 0 and 3")
	}
	return &protocol.AccountWhitelistOperation{
		AuthorizingAccount: authorizingAccountID,
		AccountToList:      accountToListID,
		NewListing:         newListing,
		Extensions:         []json.RawMessage{},
	}, nil
}

func buildAssetClaimFeesOperation(issuerID, assetID protocol.ObjectID, precision uint8, amount float64, claimFromAssetID *protocol.ObjectID) (*protocol.AssetClaimFeesOperation, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	op := &protocol.AssetClaimFeesOperation{
		Issuer:        issuerID,
		AmountToClaim: protocol.AssetAmount{Amount: mustAmountFromFloat(amount, precision), AssetID: assetID},
	}
	if claimFromAssetID != nil {
		op.Extensions = &protocol.AssetClaimFeesExtensions{ClaimFromAssetID: claimFromAssetID}
	}
	if op.AmountToClaim.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	return op, nil
}

func buildProposalCreateOperation(accountID protocol.ObjectID, expiration time.Time, reviewPeriodSeconds *uint32, ops ...protocol.Operation) (*protocol.ProposalCreateOperation, error) {
	if len(ops) == 0 {
		return nil, fmt.Errorf("no proposed operations provided")
	}
	wrapped := make([]protocol.OpWrapper, 0, len(ops))
	for _, op := range ops {
		if op == nil {
			return nil, fmt.Errorf("nil proposed operation")
		}
		wrapped = append(wrapped, protocol.OpWrapper{Op: protocol.OperationEnvelope{Operation: op}})
	}

	out := &protocol.ProposalCreateOperation{
		FeePayingAccount: accountID,
		ExpirationTime:   protocol.NewTime(expiration),
		ProposedOps:      wrapped,
		Extensions:       []json.RawMessage{},
	}
	if reviewPeriodSeconds != nil {
		out.ReviewPeriodSeconds = reviewPeriodSeconds
	}
	return out, nil
}

func buildProposalUpdateOperation(
	accountID protocol.ObjectID,
	proposalID protocol.ObjectID,
	activeApprovalsToAdd []protocol.ObjectID,
	activeApprovalsToRemove []protocol.ObjectID,
	ownerApprovalsToAdd []protocol.ObjectID,
	ownerApprovalsToRemove []protocol.ObjectID,
	keyApprovalsToAdd []protocol.PublicKey,
	keyApprovalsToRemove []protocol.PublicKey,
) (*protocol.ProposalUpdateOperation, error) {
	if len(activeApprovalsToAdd) == 0 && len(activeApprovalsToRemove) == 0 &&
		len(ownerApprovalsToAdd) == 0 && len(ownerApprovalsToRemove) == 0 &&
		len(keyApprovalsToAdd) == 0 && len(keyApprovalsToRemove) == 0 {
		return nil, fmt.Errorf("no proposal approval changes provided")
	}
	if hasObjectIDOverlap(activeApprovalsToAdd, activeApprovalsToRemove) {
		return nil, fmt.Errorf("cannot add and remove the same active approval")
	}
	if hasObjectIDOverlap(ownerApprovalsToAdd, ownerApprovalsToRemove) {
		return nil, fmt.Errorf("cannot add and remove the same owner approval")
	}
	if hasPublicKeyOverlap(keyApprovalsToAdd, keyApprovalsToRemove) {
		return nil, fmt.Errorf("cannot add and remove the same key approval")
	}

	return &protocol.ProposalUpdateOperation{
		FeePayingAccount:        accountID,
		Proposal:                proposalID,
		ActiveApprovalsToAdd:    cloneObjectIDs(activeApprovalsToAdd),
		ActiveApprovalsToRemove: cloneObjectIDs(activeApprovalsToRemove),
		OwnerApprovalsToAdd:     cloneObjectIDs(ownerApprovalsToAdd),
		OwnerApprovalsToRemove:  cloneObjectIDs(ownerApprovalsToRemove),
		KeyApprovalsToAdd:       clonePublicKeys(keyApprovalsToAdd),
		KeyApprovalsToRemove:    clonePublicKeys(keyApprovalsToRemove),
		Extensions:              []json.RawMessage{},
	}, nil
}

func buildProposalDeleteOperation(accountID, proposalID protocol.ObjectID, usingOwnerAuthority bool) *protocol.ProposalDeleteOperation {
	return &protocol.ProposalDeleteOperation{
		FeePayingAccount:    accountID,
		UsingOwnerAuthority: usingOwnerAuthority,
		Proposal:            proposalID,
		Extensions:          []json.RawMessage{},
	}
}

func buildWithdrawPermissionCreateOperation(withdrawFromID, authorizedID, assetID protocol.ObjectID, precision uint8, amount float64, withdrawalPeriodSec uint32, periodsUntilExpiration uint32, periodStartTime time.Time) (*protocol.WithdrawPermissionCreateOperation, error) {
	if err := validateWithdrawPermissionCommon(withdrawFromID, authorizedID, amount, withdrawalPeriodSec, periodsUntilExpiration, periodStartTime); err != nil {
		return nil, err
	}
	op := &protocol.WithdrawPermissionCreateOperation{
		WithdrawFromAccount:    withdrawFromID,
		AuthorizedAccount:      authorizedID,
		WithdrawalLimit:        protocol.AssetAmount{Amount: mustAmountFromFloat(amount, precision), AssetID: assetID},
		WithdrawalPeriodSec:    withdrawalPeriodSec,
		PeriodsUntilExpiration: periodsUntilExpiration,
		PeriodStartTime:        protocol.NewTime(periodStartTime.UTC()),
	}
	if op.WithdrawalLimit.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	return op, nil
}

func buildWithdrawPermissionUpdateOperation(withdrawFromID, authorizedID, permissionID, assetID protocol.ObjectID, precision uint8, amount float64, withdrawalPeriodSec uint32, periodsUntilExpiration uint32, periodStartTime time.Time) (*protocol.WithdrawPermissionUpdateOperation, error) {
	if err := validateWithdrawPermissionCommon(withdrawFromID, authorizedID, amount, withdrawalPeriodSec, periodsUntilExpiration, periodStartTime); err != nil {
		return nil, err
	}
	op := &protocol.WithdrawPermissionUpdateOperation{
		WithdrawFromAccount:    withdrawFromID,
		AuthorizedAccount:      authorizedID,
		PermissionToUpdate:     permissionID,
		WithdrawalLimit:        protocol.AssetAmount{Amount: mustAmountFromFloat(amount, precision), AssetID: assetID},
		WithdrawalPeriodSec:    withdrawalPeriodSec,
		PeriodStartTime:        protocol.NewTime(periodStartTime.UTC()),
		PeriodsUntilExpiration: periodsUntilExpiration,
	}
	if op.WithdrawalLimit.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	return op, nil
}

func buildWithdrawPermissionClaimOperation(permissionID, withdrawFromID, withdrawToID, assetID protocol.ObjectID, precision uint8, amount float64, memo json.RawMessage) (*protocol.WithdrawPermissionClaimOperation, error) {
	if withdrawFromID == withdrawToID {
		return nil, fmt.Errorf("withdraw from account and withdraw to account must differ")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	op := &protocol.WithdrawPermissionClaimOperation{
		WithdrawPermission:  permissionID,
		WithdrawFromAccount: withdrawFromID,
		WithdrawToAccount:   withdrawToID,
		AmountToWithdraw:    protocol.AssetAmount{Amount: mustAmountFromFloat(amount, precision), AssetID: assetID},
		Memo:                memo,
	}
	if op.AmountToWithdraw.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	return op, nil
}

func buildWithdrawPermissionDeleteOperation(withdrawFromID, authorizedID, permissionID protocol.ObjectID) (*protocol.WithdrawPermissionDeleteOperation, error) {
	if withdrawFromID == authorizedID {
		return nil, fmt.Errorf("withdraw from account and authorized account must differ")
	}
	return &protocol.WithdrawPermissionDeleteOperation{
		WithdrawFromAccount:  withdrawFromID,
		AuthorizedAccount:    authorizedID,
		WithdrawalPermission: permissionID,
	}, nil
}

func buildAccountTransferOperation(accountID, newOwnerID protocol.ObjectID) *protocol.AccountTransferOperation {
	return &protocol.AccountTransferOperation{
		AccountID:  accountID,
		NewOwner:   newOwnerID,
		Extensions: []json.RawMessage{},
	}
}

func buildBalanceClaimOperation(depositToID, balanceID protocol.ObjectID, balanceOwnerKey protocol.PublicKey, assetID protocol.ObjectID, precision uint8, amount float64) (*protocol.BalanceClaimOperation, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	op := &protocol.BalanceClaimOperation{
		Fee:              protocol.AssetAmount{AssetID: protocol.ObjectID{Space: 1, Type: 3, ID: 0}},
		DepositToAccount: depositToID,
		BalanceToClaim:   balanceID,
		BalanceOwnerKey:  balanceOwnerKey,
		TotalClaimed:     protocol.AssetAmount{Amount: mustAmountFromFloat(amount, precision), AssetID: assetID},
	}
	if op.TotalClaimed.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	return op, nil
}

func buildOverrideTransferOperation(issuerID, fromID, toID, assetID protocol.ObjectID, precision uint8, amount float64, memo json.RawMessage) (*protocol.OverrideTransferOperation, error) {
	if issuerID == fromID {
		return nil, fmt.Errorf("issuer and from account must differ")
	}
	if fromID == toID {
		return nil, fmt.Errorf("from and to account must differ")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	op := &protocol.OverrideTransferOperation{
		Issuer:     issuerID,
		From:       fromID,
		To:         toID,
		Amount:     protocol.AssetAmount{Amount: mustAmountFromFloat(amount, precision), AssetID: assetID},
		Memo:       memo,
		Extensions: []json.RawMessage{},
	}
	if op.Amount.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	return op, nil
}

func buildBidCollateralOperation(bidderID, collateralAssetID protocol.ObjectID, collateralPrecision uint8, additionalCollateral float64, debtAssetID protocol.ObjectID, debtPrecision uint8, debtCovered float64) (*protocol.BidCollateralOperation, error) {
	if additionalCollateral == 0 && debtCovered == 0 {
		return nil, fmt.Errorf("no bid collateral fields provided")
	}
	if debtCovered > 0 && additionalCollateral <= 0 {
		return nil, fmt.Errorf("additional collateral must be greater than zero when debt is covered")
	}
	if additionalCollateral < 0 || debtCovered < 0 {
		return nil, fmt.Errorf("amounts must not be negative")
	}

	op := &protocol.BidCollateralOperation{
		Bidder:               bidderID,
		AdditionalCollateral: protocol.AssetAmount{Amount: mustAmountFromFloat(additionalCollateral, collateralPrecision), AssetID: collateralAssetID},
		DebtCovered:          protocol.AssetAmount{Amount: mustAmountFromFloat(debtCovered, debtPrecision), AssetID: debtAssetID},
		Extensions:           []json.RawMessage{},
	}
	return op, nil
}

func buildCommitteeMemberCreateOperation(accountID protocol.ObjectID, url string) (*protocol.CommitteeMemberCreateOperation, error) {
	if len(url) >= maxChainURLLength {
		return nil, fmt.Errorf("url must be shorter than %d characters", maxChainURLLength)
	}
	return &protocol.CommitteeMemberCreateOperation{
		CommitteeMemberAccount: accountID,
		URL:                    url,
	}, nil
}

func buildCommitteeMemberUpdateOperation(memberID, accountID protocol.ObjectID, newURL string) (*protocol.CommitteeMemberUpdateOperation, error) {
	trimmedURL := strings.TrimSpace(newURL)
	if trimmedURL == "" {
		return nil, fmt.Errorf("no committee member update fields provided")
	}
	if len(trimmedURL) >= maxChainURLLength {
		return nil, fmt.Errorf("url must be shorter than %d characters", maxChainURLLength)
	}
	return &protocol.CommitteeMemberUpdateOperation{
		CommitteeMember:        memberID,
		CommitteeMemberAccount: accountID,
		NewURL:                 &trimmedURL,
	}, nil
}

func buildCommitteeMemberUpdateGlobalParametersOperation(newParameters protocol.ChainParameters) *protocol.CommitteeMemberUpdateGlobalParametersOperation {
	return &protocol.CommitteeMemberUpdateGlobalParametersOperation{
		NewParameters: newParameters,
	}
}

func setOperationFee(op protocol.Operation, fee protocol.AssetAmount) error {
	if op == nil {
		return fmt.Errorf("nil operation")
	}

	value := reflect.ValueOf(op)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return fmt.Errorf("operation %T must be a non-nil pointer", op)
	}

	elem := value.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("operation %T must point to a struct", op)
	}

	field := elem.FieldByName("Fee")
	if !field.IsValid() || !field.CanSet() || field.Type() != reflect.TypeOf(protocol.AssetAmount{}) {
		return fmt.Errorf("operation %T does not expose a writable Fee field", op)
	}

	field.Set(reflect.ValueOf(fee))
	return nil
}

func buildTransactionExpiration(chainTime string, expireInSecs int) protocol.Time {
	base := time.Now().UTC()
	if parsed, ok := parseChainTime(chainTime); ok {
		base = parsed
	}
	return protocol.NewTime(base.Add(time.Duration(expireInSecs) * time.Second))
}

func parseChainTime(value string) (time.Time, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return time.Time{}, false
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
	}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, trimmed)
		if err == nil {
			return parsed.UTC(), true
		}
	}
	return time.Time{}, false
}

func validateWithdrawPermissionCommon(withdrawFromID, authorizedID protocol.ObjectID, amount float64, withdrawalPeriodSec uint32, periodsUntilExpiration uint32, periodStartTime time.Time) error {
	if withdrawFromID == authorizedID {
		return fmt.Errorf("withdraw from account and authorized account must differ")
	}
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if withdrawalPeriodSec == 0 {
		return fmt.Errorf("withdrawal period seconds must be greater than zero")
	}
	if periodsUntilExpiration == 0 {
		return fmt.Errorf("periods until expiration must be greater than zero")
	}
	if periodStartTime.IsZero() {
		return fmt.Errorf("period start time must be provided")
	}
	return nil
}

func isValidAccountName(name string) bool {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) < minAccountNameLength || len(trimmed) > maxAccountNameLength {
		return false
	}

	begin := 0
	for {
		end := strings.IndexByte(trimmed[begin:], '.')
		if end >= 0 {
			end += begin
		} else {
			end = len(trimmed)
		}
		if end-begin < minAccountNameLength {
			return false
		}
		if !isLowerAlpha(trimmed[begin]) {
			return false
		}
		last := trimmed[end-1]
		if !isLowerAlpha(last) && !isDigit(last) {
			return false
		}
		for i := begin + 1; i < end-1; i++ {
			c := trimmed[i]
			if !isLowerAlpha(c) && !isDigit(c) && c != '-' {
				return false
			}
		}
		if end == len(trimmed) {
			return true
		}
		begin = end + 1
		if begin >= len(trimmed) {
			return false
		}
	}
}

func isValidAssetSymbol(symbol string) bool {
	trimmed := strings.TrimSpace(symbol)
	if len(trimmed) < minAssetSymbolLength || len(trimmed) > maxAssetSymbolLength {
		return false
	}
	if strings.HasPrefix(trimmed, "BIT") {
		return false
	}
	if !isUpperAlpha(trimmed[0]) {
		return false
	}
	last := trimmed[len(trimmed)-1]
	if !isUpperAlpha(last) && !isDigit(last) {
		return false
	}

	dotSeen := false
	for i := 0; i < len(trimmed); i++ {
		c := trimmed[i]
		if isUpperAlpha(c) || isDigit(c) {
			continue
		}
		if c == '.' && !dotSeen {
			dotSeen = true
			continue
		}
		return false
	}
	return true
}

func validateCreateAuthority(authority protocol.Authority, label string) error {
	if authorityNumAuths(authority) == 0 {
		return fmt.Errorf("%s authority must include at least one auth", label)
	}
	if len(authority.AddressAuths) != 0 {
		return fmt.Errorf("%s authority address auths are not supported", label)
	}
	if authorityIsImpossible(authority) {
		return fmt.Errorf("%s authority threshold is impossible", label)
	}
	return nil
}

func validateAccountOptions(options protocol.AccountOptions) error {
	neededWitnesses := int(options.NumWitness)
	neededCommittee := int(options.NumCommittee)
	for _, vote := range options.Votes {
		switch vote.Type {
		case voteTypeWitness:
			if neededWitnesses > 0 {
				neededWitnesses--
			}
		case voteTypeCommittee:
			if neededCommittee > 0 {
				neededCommittee--
			}
		}
	}
	if neededWitnesses != 0 || neededCommittee != 0 {
		return fmt.Errorf("account options require at least as many witness and committee votes as requested")
	}
	return nil
}

func validateAccountCreateExtensions(owner protocol.Authority, active protocol.Authority, extensions *protocol.AccountCreateExtensions) error {
	if extensions == nil || extensions.BuybackOptions == nil {
		return nil
	}
	if extensions.OwnerSpecialAuthority != nil || extensions.ActiveSpecialAuthority != nil {
		return fmt.Errorf("buyback options cannot be combined with special authorities")
	}
	if !isNullAuthority(owner) || !isNullAuthority(active) {
		return fmt.Errorf("buyback accounts require null owner and active authorities")
	}
	if len(extensions.BuybackOptions.Markets) == 0 {
		return fmt.Errorf("buyback options require at least one market")
	}
	for _, market := range extensions.BuybackOptions.Markets {
		if market == extensions.BuybackOptions.AssetToBuy {
			return fmt.Errorf("buyback market cannot match asset to buy")
		}
	}
	return nil
}

func validateAssetOptions(options protocol.AssetOptions) error {
	if options.MaxSupply <= 0 {
		return fmt.Errorf("max supply must be greater than zero")
	}
	if options.MarketFeePercent > graphene100Percent {
		return fmt.Errorf("market fee percent must be at most %d", graphene100Percent)
	}
	if options.Extensions != nil {
		if options.Extensions.RewardPercent != nil && *options.Extensions.RewardPercent > graphene100Percent {
			return fmt.Errorf("reward percent must be at most %d", graphene100Percent)
		}
		if options.Extensions.TakerFeePercent != nil && *options.Extensions.TakerFeePercent > graphene100Percent {
			return fmt.Errorf("taker fee percent must be at most %d", graphene100Percent)
		}
	}
	if options.Flags&assetIssuerPermissionGlobalSettle != 0 {
		return fmt.Errorf("global settle cannot be set in flags")
	}
	if options.Flags&(assetIssuerPermissionWitnessFed|assetIssuerPermissionCommitteeFed) == (assetIssuerPermissionWitnessFed | assetIssuerPermissionCommitteeFed) {
		return fmt.Errorf("witness-fed and committee-fed flags cannot both be set")
	}
	if options.CoreExchangeRate.Base.Amount == 0 || options.CoreExchangeRate.Quote.Amount == 0 {
		return fmt.Errorf("core exchange rate amounts must be non-zero")
	}
	coreAsset := protocol.MustParseObjectID("1.3.0")
	if options.CoreExchangeRate.Base.AssetID != coreAsset && options.CoreExchangeRate.Quote.AssetID != coreAsset {
		return fmt.Errorf("core exchange rate must involve the core asset")
	}
	if (len(options.WhitelistAuthorities) > 0 || len(options.BlacklistAuthorities) > 0) && options.Flags&assetIssuerPermissionWhiteList == 0 {
		return fmt.Errorf("whitelist authorities require the white_list flag")
	}
	if hasObjectIDOverlap(options.WhitelistMarkets, options.BlacklistMarkets) {
		return fmt.Errorf("whitelist and blacklist markets must not overlap")
	}
	return nil
}

func validateBitAssetOptions(options protocol.BitAssetOptions) error {
	if options.MinimumFeeds == 0 {
		return fmt.Errorf("minimum feeds must be greater than zero")
	}
	if options.ForceSettlementOffsetPercent > graphene100Percent {
		return fmt.Errorf("force settlement offset percent must be at most %d", graphene100Percent)
	}
	if options.MaximumForceSettlementVolume > graphene100Percent {
		return fmt.Errorf("maximum force settlement volume must be at most %d", graphene100Percent)
	}
	if options.Extensions != nil {
		if options.Extensions.ForceSettleFeePercent != nil && *options.Extensions.ForceSettleFeePercent > graphene100Percent {
			return fmt.Errorf("force settle fee percent must be at most %d", graphene100Percent)
		}
	}
	return nil
}

func authorityNumAuths(authority protocol.Authority) int {
	return len(authority.AccountAuths) + len(authority.KeyAuths) + len(authority.AddressAuths)
}

func authorityIsImpossible(authority protocol.Authority) bool {
	var total uint64
	for _, weight := range authority.AccountAuths {
		total += uint64(weight)
	}
	for _, weight := range authority.KeyAuths {
		total += uint64(weight)
	}
	for _, weight := range authority.AddressAuths {
		total += uint64(weight)
	}
	return total < uint64(authority.WeightThreshold)
}

func isNullAuthority(authority protocol.Authority) bool {
	if authority.WeightThreshold != 1 || len(authority.KeyAuths) != 0 || len(authority.AddressAuths) != 0 || len(authority.AccountAuths) != 1 {
		return false
	}
	nullAccount := protocol.MustParseObjectID("1.2.0")
	weight, ok := authority.AccountAuths[nullAccount]
	return ok && weight == 1
}

func isLowerAlpha(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isUpperAlpha(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (w *Wallet) resolveAccountIDs(ctx context.Context, names []string) ([]protocol.ObjectID, error) {
	if len(names) == 0 {
		return nil, nil
	}
	ids := make([]protocol.ObjectID, 0, len(names))
	for _, name := range names {
		account, err := w.parent.Account(ctx, name)
		if err != nil {
			return nil, err
		}
		ids = append(ids, account.ID)
	}
	return ids, nil
}

func parsePublicKeys(values []string) ([]protocol.PublicKey, error) {
	if len(values) == 0 {
		return nil, nil
	}
	keys := make([]protocol.PublicKey, 0, len(values))
	for _, value := range values {
		key, err := protocol.ParsePublicKey(value)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func hasObjectIDOverlap(additions []protocol.ObjectID, removals []protocol.ObjectID) bool {
	if len(additions) == 0 || len(removals) == 0 {
		return false
	}
	seen := make(map[string]struct{}, len(removals))
	for _, id := range removals {
		seen[id.String()] = struct{}{}
	}
	for _, id := range additions {
		if _, ok := seen[id.String()]; ok {
			return true
		}
	}
	return false
}

func hasPublicKeyOverlap(additions []protocol.PublicKey, removals []protocol.PublicKey) bool {
	if len(additions) == 0 || len(removals) == 0 {
		return false
	}
	seen := make(map[string]struct{}, len(removals))
	for _, key := range removals {
		seen[key.String()] = struct{}{}
	}
	for _, key := range additions {
		if _, ok := seen[key.String()]; ok {
			return true
		}
	}
	return false
}

func cloneObjectIDs(values []protocol.ObjectID) []protocol.ObjectID {
	if len(values) == 0 {
		return nil
	}
	out := make([]protocol.ObjectID, len(values))
	copy(out, values)
	return out
}

func clonePublicKeys(values []protocol.PublicKey) []protocol.PublicKey {
	if len(values) == 0 {
		return nil
	}
	out := make([]protocol.PublicKey, len(values))
	copy(out, values)
	return out
}

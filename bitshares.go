package bitshares

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/froooze/go-bitshares/ecc"
	"github.com/ulikunitz/xz/lzma"
)

// BitShares is the high-level client facade for the library.
type BitShares struct {
	chain *ChainClient

	Accounts *AccountStore
	Assets   *AssetStore
	Fees     *FeeTable
}

// NewBitShares creates a new BitShares client facade.
func NewBitShares(endpoint string, reconnectDelay time.Duration) (*BitShares, error) {
	chain, err := NewChainClient(endpoint, reconnectDelay)
	if err != nil {
		return nil, err
	}

	return &BitShares{
		chain:    chain,
		Accounts: &AccountStore{chain: chain},
		Assets:   &AssetStore{chain: chain},
		Fees:     newFeeTable(chain),
	}, nil
}

// Connect opens the websocket connection and refreshes the fee snapshot.
func (b *BitShares) Connect(ctx context.Context) error {
	if b == nil || b.chain == nil {
		return fmt.Errorf("bitshares client is not configured")
	}
	if err := b.chain.Connect(ctx); err != nil {
		return err
	}
	if b.Fees != nil {
		_ = b.Fees.Update(ctx)
	}
	return nil
}

// Close shuts down the underlying chain client.
func (b *BitShares) Close() error {
	if b == nil || b.chain == nil {
		return nil
	}
	return b.chain.Close()
}

// Chain returns the active chain configuration.
func (b *BitShares) Chain() ChainConfig {
	if b == nil || b.chain == nil {
		return ChainConfig{}
	}
	return b.chain.Chain()
}

// Account resolves an account by name.
func (b *BitShares) Account(ctx context.Context, name string) (*AccountInfo, error) {
	return b.Accounts.Get(ctx, name)
}

// AccountByID resolves an account by object id.
func (b *BitShares) AccountByID(ctx context.Context, id string) (*AccountInfo, error) {
	return b.Accounts.ID(ctx, id)
}

// Asset resolves an asset by symbol.
func (b *BitShares) Asset(ctx context.Context, symbol string) (*AssetInfo, error) {
	return b.Assets.Get(ctx, symbol)
}

// AssetByID resolves an asset by object id.
func (b *BitShares) AssetByID(ctx context.Context, id string) (*AssetInfo, error) {
	return b.Assets.ID(ctx, id)
}

// Ticker fetches the market ticker for a pair of symbols.
func (b *BitShares) Ticker(ctx context.Context, baseSymbol, quoteSymbol string, reply any) error {
	base, err := b.Asset(ctx, strings.ToUpper(baseSymbol))
	if err != nil {
		return err
	}
	quote, err := b.Asset(ctx, strings.ToUpper(quoteSymbol))
	if err != nil {
		return err
	}
	return b.chain.GetTicker(ctx, base.ID.String(), quote.ID.String(), reply)
}

// OrderBook fetches the market order book for a pair of symbols.
func (b *BitShares) OrderBook(ctx context.Context, quoteSymbol, baseSymbol string, limit int, reply any) error {
	quote, err := b.Asset(ctx, strings.ToUpper(quoteSymbol))
	if err != nil {
		return err
	}
	base, err := b.Asset(ctx, strings.ToUpper(baseSymbol))
	if err != nil {
		return err
	}
	if limit <= 0 {
		limit = 50
	}
	return b.chain.GetOrderBook(ctx, quote.ID.String(), base.ID.String(), limit, reply)
}

// LimitOrders fetches active limit orders for a market.
func (b *BitShares) LimitOrders(ctx context.Context, quoteSymbol, baseSymbol string, limit int, reply any) error {
	quote, err := b.Asset(ctx, strings.ToUpper(quoteSymbol))
	if err != nil {
		return err
	}
	base, err := b.Asset(ctx, strings.ToUpper(baseSymbol))
	if err != nil {
		return err
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	return b.chain.GetLimitOrders(ctx, quote.ID.String(), base.ID.String(), limit, reply)
}

// TradeHistory fetches historical trades for a market.
func (b *BitShares) TradeHistory(ctx context.Context, quoteSymbol, baseSymbol string, start, stop string, bucketSeconds int, reply any) error {
	quote, err := b.Asset(ctx, strings.ToUpper(quoteSymbol))
	if err != nil {
		return err
	}
	base, err := b.Asset(ctx, strings.ToUpper(baseSymbol))
	if err != nil {
		return err
	}
	return b.chain.GetMarketHistory(ctx, quote.ID.String(), base.ID.String(), bucketSeconds, start, stop, reply)
}

// GenerateKeys derives BitShares key pairs for the requested roles from password bytes.
// Callers own the input buffer and may wipe it after the call returns.
func GenerateKeys(accountName string, password []byte, roles []string, prefix string) (map[string]*ecc.PrivateKey, map[string]string, error) {
	return ecc.GenerateKeys(accountName, password, roles, prefix)
}

// Login derives the active and memo keys for an account from password bytes and returns a wallet session.
// The returned wallet keeps the keys in memory until Wallet.Wipe is called.
func (b *BitShares) Login(ctx context.Context, accountName string, password []byte, feeSymbol string) (*Wallet, error) {
	if b == nil || b.chain == nil {
		return nil, fmt.Errorf("bitshares client is not configured")
	}
	if strings.TrimSpace(accountName) == "" || len(password) == 0 {
		return nil, fmt.Errorf("account name or password required")
	}

	account, err := b.Account(ctx, accountName)
	if err != nil {
		return nil, err
	}

	activeSeed := appendAccountSecretSeed(accountName, "active", password)
	activeKey := ecc.PrivateKeyFromSeed(activeSeed)
	zeroBytes(activeSeed)
	if !authorityHasKey(account.Active, activeKey.PublicKey().String()) {
		return nil, fmt.Errorf("the pair of login and password do not match")
	}

	feeAssetSymbol := feeSymbol
	if strings.TrimSpace(feeAssetSymbol) == "" {
		feeAssetSymbol = b.Chain().CoreAsset
	}
	feeAsset, err := b.Asset(ctx, feeAssetSymbol)
	if err != nil {
		return nil, err
	}

	memoSeed := appendAccountSecretSeed(accountName, "memo", password)
	memoKey := ecc.PrivateKeyFromSeed(memoSeed)
	zeroBytes(memoSeed)
	if account.Options.MemoKey == activeKey.PublicKey().String() {
		memoKey = activeKey
	}

	return &Wallet{
		parent:    b,
		account:   account,
		activeKey: activeKey,
		memoKey:   memoKey,
		feeAsset:  feeAsset,
	}, nil
}

// LoginWithWIF creates a wallet session from an existing active key passed as WIF bytes.
// The returned wallet keeps the key in memory until Wallet.Wipe is called.
func (b *BitShares) LoginWithWIF(ctx context.Context, accountName string, activeWIF []byte, feeSymbol string) (*Wallet, error) {
	if b == nil || b.chain == nil {
		return nil, fmt.Errorf("bitshares client is not configured")
	}
	if strings.TrimSpace(accountName) == "" || len(activeWIF) == 0 {
		return nil, fmt.Errorf("account name and active wif are required")
	}
	account, err := b.Account(ctx, accountName)
	if err != nil {
		return nil, err
	}
	activeKey, err := ecc.PrivateKeyFromWIF(activeWIF)
	if err != nil {
		return nil, err
	}
	if !authorityHasKey(account.Active, activeKey.PublicKey().String()) {
		return nil, fmt.Errorf("active key is not authorized for account %s", accountName)
	}

	feeAssetSymbol := feeSymbol
	if strings.TrimSpace(feeAssetSymbol) == "" {
		feeAssetSymbol = b.Chain().CoreAsset
	}
	feeAsset, err := b.Asset(ctx, feeAssetSymbol)
	if err != nil {
		return nil, err
	}

	var memoKey *ecc.PrivateKey
	if account.Options.MemoKey == activeKey.PublicKey().String() {
		memoKey = activeKey
	}
	return &Wallet{
		parent:    b,
		account:   account,
		activeKey: activeKey,
		memoKey:   memoKey,
		feeAsset:  feeAsset,
	}, nil
}

// LoginFromFile restores a wallet session from an encrypted BitShares backup file.
// The recovered keys stay in memory until Wallet.Wipe is called.
func (b *BitShares) LoginFromFile(ctx context.Context, backup []byte, password []byte, accountName, feeSymbol string) (*Wallet, error) {
	if b == nil || b.chain == nil {
		return nil, fmt.Errorf("bitshares client is not configured")
	}
	if len(backup) < 33 {
		return nil, fmt.Errorf("backup file is too small")
	}
	if len(password) == 0 || strings.TrimSpace(accountName) == "" {
		return nil, fmt.Errorf("password and account name are required")
	}

	headerKey, err := ecc.PublicKeyFromBytes(backup[:33])
	if err != nil {
		return nil, err
	}

	plain, err := ecc.DecryptWithChecksum(ecc.PrivateKeyFromSeed(password), headerKey, "", backup[33:], false)
	if err != nil {
		return nil, err
	}
	defer zeroBytes(plain)

	jsonBlob, err := decompressBackupPayload(ctx, plain)
	if err != nil {
		return nil, err
	}

	var payload struct {
		Wallet []struct {
			EncryptionKey string `json:"encryption_key"`
		} `json:"wallet"`
		PrivateKeys []backupKeyRecord `json:"private_keys"`
	}
	if err := json.Unmarshal(jsonBlob, &payload); err != nil {
		return nil, err
	}
	if len(payload.Wallet) == 0 {
		return nil, fmt.Errorf("backup file does not contain a wallet record")
	}

	account, err := b.Account(ctx, accountName)
	if err != nil {
		return nil, err
	}

	passwordAES := ecc.FromSeed(password)
	defer passwordAES.Wipe()
	encryptionSeed, err := passwordAES.DecryptHexToBuffer(payload.Wallet[0].EncryptionKey)
	if err != nil {
		return nil, err
	}
	defer zeroBytes(encryptionSeed)
	secretAES := ecc.FromSeed(encryptionSeed)
	defer secretAES.Wipe()

	activeEntry := findBackupActiveKey(payload.PrivateKeys, account.Active)
	if activeEntry == nil {
		return nil, fmt.Errorf("not found active key for account %s", accountName)
	}

	activeHex, err := secretAES.DecryptHex(activeEntry.EncryptedKey)
	if err != nil {
		return nil, err
	}
	activeBytes, err := hex.DecodeString(strings.TrimSpace(activeHex))
	if err != nil {
		return nil, err
	}
	defer zeroBytes(activeBytes)
	activeKey := ecc.PrivateKeyFromBytes(activeBytes)

	feeAssetSymbol := feeSymbol
	if strings.TrimSpace(feeAssetSymbol) == "" {
		feeAssetSymbol = b.Chain().CoreAsset
	}
	feeAsset, err := b.Asset(ctx, feeAssetSymbol)
	if err != nil {
		return nil, err
	}

	memoKey := activeKey
	if account.Options.MemoKey != activeKey.PublicKey().String() {
		memoEntry := findBackupKey(payload.PrivateKeys, account.Options.MemoKey)
		if memoEntry == nil {
			return nil, fmt.Errorf("not found memo key for account %s", accountName)
		}
		memoHex, err := secretAES.DecryptHex(memoEntry.EncryptedKey)
		if err != nil {
			return nil, err
		}
		memoBytes, err := hex.DecodeString(strings.TrimSpace(memoHex))
		if err != nil {
			return nil, err
		}
		defer zeroBytes(memoBytes)
		memoKey = ecc.PrivateKeyFromBytes(memoBytes)
	}

	return &Wallet{
		parent:    b,
		account:   account,
		activeKey: activeKey,
		memoKey:   memoKey,
		feeAsset:  feeAsset,
	}, nil
}

// FeeAmount returns the cached fee value for a given operation name.
func (b *BitShares) FeeAmount(name string) (float64, bool) {
	if b == nil || b.Fees == nil {
		return 0, false
	}
	return b.Fees.Get(name)
}

// RoundAmount converts a human amount into the chain integer representation.
func RoundAmount(amount float64, precision uint8) int64 {
	if amount <= 0 {
		return 0
	}
	rounded, err := ParseAmount(strconv.FormatFloat(amount, 'f', -1, 64), precision)
	if err != nil {
		return 0
	}
	return rounded
}

// ParseAmount converts a decimal string into the chain integer representation.
func ParseAmount(amount string, precision uint8) (int64, error) {
	return parseAmount(amount, precision, false)
}

func parseAmount(amount string, precision uint8, allowNegative bool) (int64, error) {
	value := strings.TrimSpace(amount)
	if value == "" {
		return 0, fmt.Errorf("empty amount")
	}

	rat, ok := new(big.Rat).SetString(value)
	if !ok {
		return 0, fmt.Errorf("invalid amount %q", amount)
	}
	if rat.Sign() < 0 && !allowNegative {
		return 0, fmt.Errorf("amount must not be negative")
	}

	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision)), nil)
	scaled := new(big.Rat).Mul(rat, new(big.Rat).SetInt(scale))

	num := new(big.Int).Set(scaled.Num())
	den := new(big.Int).Set(scaled.Denom())
	if den.Sign() == 0 {
		return 0, fmt.Errorf("invalid amount %q", amount)
	}

	quotient := new(big.Int).Quo(num, den)
	if !quotient.IsInt64() {
		return 0, fmt.Errorf("amount %q exceeds int64 range", amount)
	}
	return quotient.Int64(), nil
}

func multiplyAmountPrice(amount, price float64, precision uint8) (int64, error) {
	if amount <= 0 || price <= 0 {
		return 0, nil
	}
	amountRat, ok := new(big.Rat).SetString(strconv.FormatFloat(amount, 'f', -1, 64))
	if !ok {
		return 0, fmt.Errorf("invalid amount %v", amount)
	}
	priceRat, ok := new(big.Rat).SetString(strconv.FormatFloat(price, 'f', -1, 64))
	if !ok {
		return 0, fmt.Errorf("invalid price %v", price)
	}
	product := new(big.Rat).Mul(amountRat, priceRat)
	return ParseAmount(product.FloatString(int(precision)+18), precision)
}

func decompressBackupPayload(ctx context.Context, data []byte) ([]byte, error) {
	if err := requireContext(ctx); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("empty backup payload")
	}

	reader, err := lzma.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	out, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("empty decompressed backup payload")
	}
	return out, nil
}

type backupKeyRecord struct {
	PubKey       string `json:"pubkey"`
	EncryptedKey string `json:"encrypted_key"`
}

func findBackupKey(keys []backupKeyRecord, pub string) *backupKeyRecord {
	for i := range keys {
		if keys[i].PubKey == pub {
			return &keys[i]
		}
	}
	return nil
}

func findBackupActiveKey(keys []backupKeyRecord, auth Authority) *backupKeyRecord {
	for _, keyAuth := range auth.KeyAuths {
		if match := findBackupKey(keys, keyAuth.Key); match != nil {
			return match
		}
	}
	return nil
}

func authorityHasKey(auth Authority, pub string) bool {
	for _, keyAuth := range auth.KeyAuths {
		if keyAuth.Key == pub {
			return true
		}
	}
	return false
}

func appendAccountSecretSeed(accountName, role string, password []byte) []byte {
	seed := make([]byte, 0, len(accountName)+len(role)+len(password))
	seed = append(seed, accountName...)
	seed = append(seed, role...)
	seed = append(seed, password...)
	return seed
}

func zeroBytes(buf []byte) {
	for i := range buf {
		buf[i] = 0
	}
}

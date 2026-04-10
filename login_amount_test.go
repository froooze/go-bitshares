package bitshares

import (
	"context"
	"errors"
	"testing"

	"github.com/froooze/go-bitshares/ecc"
)

func TestRoundAmountHandlesFloatArtifacts(t *testing.T) {
	t.Parallel()

	if got, want := RoundAmount(0.1+0.2, 2), int64(30); got != want {
		t.Fatalf("RoundAmount(0.1+0.2, 2) = %d, want %d", got, want)
	}
}

func TestParseAmountTruncatesTowardZero(t *testing.T) {
	t.Parallel()

	got, err := ParseAmount("1.239", 2)
	if err != nil {
		t.Fatalf("ParseAmount() error = %v", err)
	}
	if want := int64(123); got != want {
		t.Fatalf("ParseAmount(1.239, 2) = %d, want %d", got, want)
	}
}

func TestParseAmountRejectsNegativeValues(t *testing.T) {
	t.Parallel()

	_, err := ParseAmount("-1.0", 5)
	if err == nil {
		t.Fatal("ParseAmount(-1.0, 5) error = nil, want error")
	}
}

func TestLoginWithWIFRejectsUnauthorizedKey(t *testing.T) {
	t.Parallel()

	client := newCachedTestBitShares(t)
	ctx := context.Background()
	accountKey := ecc.PrivateKeyFromSeed([]byte("alice-active"))
	otherKey := ecc.PrivateKeyFromSeed([]byte("other-active"))

	client.chain.storeAccount(&AccountInfo{
		Name: "alice",
		Active: Authority{
			KeyAuths: []KeyWeightPair{{Key: accountKey.PublicKey().String(), Weight: 1}},
		},
		Options: AccountOptions{MemoKey: accountKey.PublicKey().String()},
	})
	client.chain.storeAsset(&AssetInfo{Symbol: "BTS"})

	_, err := client.LoginWithWIF(ctx, "alice", []byte(otherKey.WIF()), "BTS")
	if err == nil {
		t.Fatal("LoginWithWIF() error = nil, want unauthorized key error")
	}
}

func TestLoginWithWIFLeavesMemoUnsetWhenDistinct(t *testing.T) {
	t.Parallel()

	client := newCachedTestBitShares(t)
	ctx := context.Background()
	activeKey := ecc.PrivateKeyFromSeed([]byte("alice-active"))
	memoKey := ecc.PrivateKeyFromSeed([]byte("alice-memo"))

	client.chain.storeAccount(&AccountInfo{
		Name: "alice",
		Active: Authority{
			KeyAuths: []KeyWeightPair{{Key: activeKey.PublicKey().String(), Weight: 1}},
		},
		Options: AccountOptions{MemoKey: memoKey.PublicKey().String()},
	})
	client.chain.storeAsset(&AssetInfo{Symbol: "BTS"})

	wallet, err := client.LoginWithWIF(ctx, "alice", []byte(activeKey.WIF()), "BTS")
	if err != nil {
		t.Fatalf("LoginWithWIF() error = %v", err)
	}
	if wallet.MemoPublicKey() != nil {
		t.Fatal("LoginWithWIF() memo key = active key, want nil when account memo key differs")
	}
}

func TestLoginWithWIFUsesActiveKeyAsMemoWhenMatching(t *testing.T) {
	t.Parallel()

	client := newCachedTestBitShares(t)
	ctx := context.Background()
	activeKey := ecc.PrivateKeyFromSeed([]byte("alice-active"))

	client.chain.storeAccount(&AccountInfo{
		Name: "alice",
		Active: Authority{
			KeyAuths: []KeyWeightPair{{Key: activeKey.PublicKey().String(), Weight: 1}},
		},
		Options: AccountOptions{MemoKey: activeKey.PublicKey().String()},
	})
	client.chain.storeAsset(&AssetInfo{Symbol: "BTS"})

	wallet, err := client.LoginWithWIF(ctx, "alice", []byte(activeKey.WIF()), "BTS")
	if err != nil {
		t.Fatalf("LoginWithWIF() error = %v", err)
	}
	if wallet.MemoPublicKey() == nil {
		t.Fatal("LoginWithWIF() memo key = nil, want active key when memo key matches")
	}
	if got, want := wallet.MemoPublicKey().String(), activeKey.PublicKey().String(); got != want {
		t.Fatalf("LoginWithWIF() memo pubkey = %q, want %q", got, want)
	}
}

func TestLoginRequiresContext(t *testing.T) {
	t.Parallel()

	client := newCachedTestBitShares(t)
	_, err := client.Login(nil, "alice", []byte("password"), "BTS")
	if !errors.Is(err, ErrNilContext) {
		t.Fatalf("Login(nil, ...) error = %v, want %v", err, ErrNilContext)
	}
}

func TestWalletWipeClearsPrivateKeys(t *testing.T) {
	t.Parallel()

	activeKey := ecc.PrivateKeyFromSeed([]byte("wallet-wipe-active"))
	memoKey := ecc.PrivateKeyFromSeed([]byte("wallet-wipe-memo"))
	wallet := &Wallet{
		activeKey: activeKey,
		memoKey:   memoKey,
	}

	wallet.Wipe()

	if wallet.ActivePublicKey() != nil {
		t.Fatal("Wallet.Wipe() left active key available")
	}
	if wallet.MemoPublicKey() != nil {
		t.Fatal("Wallet.Wipe() left memo key available")
	}
	if wallet.activeKey != nil || wallet.memoKey != nil {
		t.Fatal("Wallet.Wipe() did not clear internal key pointers")
	}
}

func newCachedTestBitShares(t *testing.T) *BitShares {
	t.Helper()

	chain := &ChainClient{
		chain:       defaultChains[0],
		accounts:    make(map[string]*AccountInfo),
		accountByID: make(map[string]*AccountInfo),
		assets:      make(map[string]*AssetInfo),
		assetByID:   make(map[string]*AssetInfo),
	}

	return &BitShares{
		chain:    chain,
		Accounts: &AccountStore{chain: chain},
		Assets:   &AssetStore{chain: chain},
		Fees:     newFeeTable(chain),
	}
}

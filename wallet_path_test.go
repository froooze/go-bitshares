package bitshares

import (
	"strings"
	"testing"
	"time"

	"github.com/froooze/go-bitshares/protocol"
)

func TestBuildTransactionExpirationUsesChainTimeWithoutZone(t *testing.T) {
	t.Parallel()

	expiration := buildTransactionExpiration("2026-04-11T12:34:56", 15)
	want := time.Date(2026, 4, 11, 12, 35, 11, 0, time.UTC)
	if got := expiration.Time; !got.Equal(want) {
		t.Fatalf("buildTransactionExpiration() = %s, want %s", got.Format(time.RFC3339), want.Format(time.RFC3339))
	}
}

func TestBuildTransactionExpirationUsesChainTimeRFC3339(t *testing.T) {
	t.Parallel()

	expiration := buildTransactionExpiration("2026-04-11T12:34:56Z", 30)
	want := time.Date(2026, 4, 11, 12, 35, 26, 0, time.UTC)
	if got := expiration.Time; !got.Equal(want) {
		t.Fatalf("buildTransactionExpiration() = %s, want %s", got.Format(time.RFC3339), want.Format(time.RFC3339))
	}
}

func TestSignTransactionRequiresActiveKeyWhenNoOverridesProvided(t *testing.T) {
	t.Parallel()

	wallet := &Wallet{
		parent: &BitShares{
			chain: &ChainClient{chain: defaultChains[0]},
		},
	}

	_, err := wallet.SignTransaction(&protocol.Transaction{})
	if err == nil || !strings.Contains(err.Error(), "active key is not configured") {
		t.Fatalf("SignTransaction() error = %v, want active key error", err)
	}
}

func TestSetOperationFeeSupportsBroaderTypedOperations(t *testing.T) {
	t.Parallel()

	op := &protocol.AccountCreateOperation{}
	fee := protocol.AssetAmount{Amount: 77, AssetID: protocol.MustParseObjectID("1.3.0")}

	if err := setOperationFee(op, fee); err != nil {
		t.Fatalf("setOperationFee() error = %v", err)
	}
	if got := op.Fee; got != fee {
		t.Fatalf("setOperationFee() fee = %#v, want %#v", got, fee)
	}
}

func TestSetOperationFeeRejectsOperationsWithoutFeeField(t *testing.T) {
	t.Parallel()

	op := protocol.NewRawOperationBody(protocol.OperationTypeTransfer, []byte(`{}`))
	fee := protocol.AssetAmount{Amount: 77, AssetID: protocol.MustParseObjectID("1.3.0")}

	err := setOperationFee(op, fee)
	if err == nil || !strings.Contains(err.Error(), "Fee field") {
		t.Fatalf("setOperationFee() error = %v, want missing Fee field error", err)
	}
}

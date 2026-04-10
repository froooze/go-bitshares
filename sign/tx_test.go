package sign

import (
	"strings"
	"testing"
	"time"

	"github.com/froooze/go-bitshares/ecc"
	"github.com/froooze/go-bitshares/protocol"
)

func TestTransactionSignerSigns(t *testing.T) {
	t.Parallel()

	tx := &protocol.Transaction{
		RefBlockNum:    1,
		RefBlockPrefix: 2,
		Expiration:     protocol.NewTime(time.Unix(1700000000, 0)),
	}
	tx.Push(&protocol.TransferOperation{
		Fee:    protocol.AssetAmount{Amount: 1, AssetID: protocol.MustParseObjectID("1.3.0")},
		From:   protocol.MustParseObjectID("1.2.1"),
		To:     protocol.MustParseObjectID("1.2.2"),
		Amount: protocol.AssetAmount{Amount: 1, AssetID: protocol.MustParseObjectID("1.3.0")},
	})

	signer := TransactionSigner{
		ChainID: strings.Repeat("01", 32),
		Keys:    []*ecc.PrivateKey{ecc.PrivateKeyFromSeed([]byte("signature-seed"))},
	}
	signed, err := signer.Sign(tx)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	if len(signed.Signatures) != 1 {
		t.Fatalf("Signatures length = %d, want 1", len(signed.Signatures))
	}
	if len(signed.Signatures[0]) != 130 {
		t.Fatalf("signature hex length = %d, want 130", len(signed.Signatures[0]))
	}
}

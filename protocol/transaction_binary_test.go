package protocol

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/froooze/go-bitshares/ecc"
)

func TestTransactionBinaryRoundTrip(t *testing.T) {
	t.Parallel()

	sender := ecc.PrivateKeyFromSeed([]byte("binary-sender"))
	recipient := ecc.PrivateKeyFromSeed([]byte("binary-recipient"))
	memo := MemoData{
		From:    sender.PublicKey().String(),
		To:      recipient.PublicKey().String(),
		Nonce:   "42",
		Message: "68656c6c6f",
	}
	rawMemo, err := json.Marshal(memo)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	original := &Transaction{
		RefBlockNum:    7,
		RefBlockPrefix: 9,
		Expiration:     NewTime(time.Unix(1700000000, 0)),
		Operations: []OperationEnvelope{
			{
				Operation: &TransferOperation{
					Fee:    AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
					From:   MustParseObjectID("1.2.1"),
					To:     MustParseObjectID("1.2.2"),
					Amount: AssetAmount{Amount: 99, AssetID: MustParseObjectID("1.3.0")},
					Memo:   rawMemo,
				},
			},
		},
	}

	raw, err := original.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary() error = %v", err)
	}

	var decoded Transaction
	if err := decoded.UnmarshalBinary(raw); err != nil {
		t.Fatalf("UnmarshalBinary() error = %v", err)
	}

	if got, want := decoded.RefBlockNum, original.RefBlockNum; got != want {
		t.Fatalf("RefBlockNum = %d, want %d", got, want)
	}
	if got, want := decoded.RefBlockPrefix, original.RefBlockPrefix; got != want {
		t.Fatalf("RefBlockPrefix = %d, want %d", got, want)
	}
	if got, want := decoded.Expiration.Unix(), original.Expiration.Unix(); got != want {
		t.Fatalf("Expiration = %d, want %d", got, want)
	}
	if len(decoded.Operations) != 1 {
		t.Fatalf("Operations length = %d, want 1", len(decoded.Operations))
	}

	decodedTransfer, ok := decoded.Operations[0].Operation.(*TransferOperation)
	if !ok {
		t.Fatalf("decoded operation type = %T, want *TransferOperation", decoded.Operations[0].Operation)
	}
	if !reflect.DeepEqual(decodedTransfer.Fee, original.Operations[0].Operation.(*TransferOperation).Fee) {
		t.Fatalf("fee mismatch: got %#v want %#v", decodedTransfer.Fee, original.Operations[0].Operation.(*TransferOperation).Fee)
	}
	if got, want := decodedTransfer.Memo, rawMemo; string(got) != string(want) {
		t.Fatalf("memo mismatch: got %s want %s", got, want)
	}
}

func TestLimitOrderUpdateAutoActionBinaryRoundTrip(t *testing.T) {
	t.Parallel()

	original := &LimitOrderUpdateOperation{
		Fee:    AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
		Seller: MustParseObjectID("1.2.1"),
		Order:  MustParseObjectID("1.7.42"),
		NewPrice: &Price{
			Base:  AssetAmount{Amount: 100, AssetID: MustParseObjectID("1.3.0")},
			Quote: AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.1")},
		},
		DeltaAmountToSell: &AssetAmount{Amount: 50, AssetID: MustParseObjectID("1.3.0")},
		NewExpiration:     func() *Time { t := NewTime(time.Unix(1700000000, 0)); return &t }(),
		OnFill: []LimitOrderAutoAction{
			{
				Kind: 0,
				TakeProfit: &CreateTakeProfitOrderAction{
					FeeAssetID:        MustParseObjectID("1.3.0"),
					SpreadPercent:     250,
					SizePercent:       5000,
					ExpirationSeconds: 3600,
					Repeat:            true,
				},
			},
		},
	}

	raw, err := original.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary() error = %v", err)
	}

	var envelope OperationEnvelope
	if err := envelope.UnmarshalBinary(raw); err != nil {
		t.Fatalf("UnmarshalBinary() error = %v", err)
	}
	decoded, ok := envelope.Operation.(*LimitOrderUpdateOperation)
	if !ok {
		t.Fatalf("decoded operation type = %T, want *LimitOrderUpdateOperation", envelope.Operation)
	}
	if len(decoded.OnFill) != 1 {
		t.Fatalf("OnFill length = %d, want 1", len(decoded.OnFill))
	}
	if decoded.OnFill[0].Kind != 0 || decoded.OnFill[0].TakeProfit == nil {
		t.Fatalf("decoded on_fill = %#v", decoded.OnFill[0])
	}
	if got, want := decoded.OnFill[0].TakeProfit.SpreadPercent, original.OnFill[0].TakeProfit.SpreadPercent; got != want {
		t.Fatalf("SpreadPercent = %d, want %d", got, want)
	}
}

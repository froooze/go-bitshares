package protocol

import (
	"encoding/json"
	"testing"
)

func TestOperationEnvelopeRoundTrip(t *testing.T) {
	original := OperationEnvelope{
		Operation: &TransferOperation{
			Fee:    AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
			From:   MustParseObjectID("1.2.1"),
			To:     MustParseObjectID("1.2.2"),
			Amount: AssetAmount{Amount: 99, AssetID: MustParseObjectID("1.3.0")},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded OperationEnvelope
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if decoded.Type() != OperationTypeTransfer {
		t.Fatalf("unexpected type: %v", decoded.Type())
	}
	if decoded.Operation == nil || decoded.Operation.Type() != OperationTypeTransfer {
		t.Fatalf("unexpected operation: %#v", decoded.Operation)
	}
	if _, ok := decoded.Operation.(*TransferOperation); !ok {
		t.Fatalf("expected typed transfer operation, got %#v", decoded.Operation)
	}
}

func TestOperationEnvelopeFallbackRoundTrip(t *testing.T) {
	original := OperationEnvelope{
		Operation: &RawOperation{
			OperationBody: OperationBody{
				Kind:    OperationType(99),
				Payload: json.RawMessage(`{"opaque":true}`),
			},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded OperationEnvelope
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if _, ok := decoded.Operation.(*RawOperation); !ok {
		t.Fatalf("expected raw fallback operation, got %#v", decoded.Operation)
	}
}

func TestTypedAssetOperationsRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []Operation{
		&AssetIssueOperation{
			Fee:            AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
			Issuer:         MustParseObjectID("1.2.1"),
			AssetToIssue:   AssetAmount{Amount: 99, AssetID: MustParseObjectID("1.3.1")},
			IssueToAccount: MustParseObjectID("1.2.2"),
		},
		&AssetReserveOperation{
			Fee:             AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
			Payer:           MustParseObjectID("1.2.1"),
			AmountToReserve: AssetAmount{Amount: 99, AssetID: MustParseObjectID("1.3.1")},
		},
	}

	for _, op := range tests {
		data, err := json.Marshal(OperationEnvelope{Operation: op})
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}

		var decoded OperationEnvelope
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}

		if decoded.Type() != op.Type() {
			t.Fatalf("decoded type = %v, want %v", decoded.Type(), op.Type())
		}
	}
}

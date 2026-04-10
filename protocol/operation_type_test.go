package protocol

import "testing"

func TestOperationTypeNames(t *testing.T) {
	if got := OperationTypeTransfer.String(); got != "transfer" {
		t.Fatalf("unexpected name: %s", got)
	}

	if OperationTypeAssetClaimPool != 47 {
		t.Fatalf("unexpected asset_claim_pool tag: %d", OperationTypeAssetClaimPool)
	}

	if OperationTypeAssetUpdateIssuer != 48 {
		t.Fatalf("unexpected asset_update_issuer tag: %d", OperationTypeAssetUpdateIssuer)
	}

	if got := OperationTypeLimitOrderUpdate.String(); got != "limit_order_update" {
		t.Fatalf("unexpected new core op name: %s", got)
	}

	if OperationTypeLimitOrderUpdate != 77 {
		t.Fatalf("unexpected limit_order_update tag: %d", OperationTypeLimitOrderUpdate)
	}

	if !IsKnownOperationType(OperationTypeCreditDealUpdate) {
		t.Fatalf("expected credit_deal_update to be known")
	}

	if IsKnownOperationType(OperationType(999)) {
		t.Fatalf("unexpected known type")
	}

	parsed, err := ParseOperationType("limit_order_update")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if parsed != OperationTypeLimitOrderUpdate {
		t.Fatalf("unexpected parsed kind: %v", parsed)
	}
}

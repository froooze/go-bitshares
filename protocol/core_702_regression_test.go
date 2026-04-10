package protocol

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/froooze/go-bitshares/ecc"
)

func TestCallOrderUpdateJSONUsesExtensionsObject(t *testing.T) {
	t.Parallel()

	op := OperationEnvelope{
		Operation: &CallOrderUpdateOperation{
			Fee:             AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
			FundingAccount:  MustParseObjectID("1.2.1"),
			DeltaCollateral: AssetAmount{Amount: 10, AssetID: MustParseObjectID("1.3.0")},
			DeltaDebt:       AssetAmount{Amount: -5, AssetID: MustParseObjectID("1.3.1")},
			Extensions: CallOrderUpdateExtensions{
				TargetCollateralRatio: uint16Ptr(1750),
			},
		},
	}

	raw, err := json.Marshal(op)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	if !strings.Contains(string(raw), `"extensions":{"target_collateral_ratio":1750}`) {
		t.Fatalf("extensions payload missing from %s", raw)
	}
	if strings.Contains(string(raw), `"target_collateral_ratio":1750`) && !strings.Contains(string(raw), `"extensions":{"target_collateral_ratio":1750}`) {
		t.Fatalf("target_collateral_ratio was flattened in %s", raw)
	}

	var decoded OperationEnvelope
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}
	callOrder, ok := decoded.Operation.(*CallOrderUpdateOperation)
	if !ok {
		t.Fatalf("decoded operation type = %T, want *CallOrderUpdateOperation", decoded.Operation)
	}
	if callOrder.Extensions.TargetCollateralRatio == nil || *callOrder.Extensions.TargetCollateralRatio != 1750 {
		t.Fatalf("decoded extensions = %#v", callOrder.Extensions)
	}
}

func TestExtensionBearingOperationsBinaryRoundTrip(t *testing.T) {
	t.Parallel()

	sender := ecc.PrivateKeyFromSeed([]byte("htlc-create-extension"))
	recipient := ecc.PrivateKeyFromSeed([]byte("htlc-create-extension-recipient"))
	memo := MemoData{
		From:    sender.PublicKey().String(),
		To:      recipient.PublicKey().String(),
		Nonce:   "7",
		Message: "68656c6c6f",
	}
	rawMemo, err := json.Marshal(memo)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	tests := []struct {
		name string
		op   Operation
	}{
		{
			name: "call_order_update",
			op: &CallOrderUpdateOperation{
				Fee:             AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
				FundingAccount:  MustParseObjectID("1.2.1"),
				DeltaCollateral: AssetAmount{Amount: 10, AssetID: MustParseObjectID("1.3.0")},
				DeltaDebt:       AssetAmount{Amount: -5, AssetID: MustParseObjectID("1.3.1")},
				Extensions: CallOrderUpdateExtensions{
					TargetCollateralRatio: uint16Ptr(1600),
				},
			},
		},
		{
			name: "asset_publish_feed",
			op: &AssetPublishFeedOperation{
				Fee:       AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
				Publisher: MustParseObjectID("1.2.1"),
				AssetID:   MustParseObjectID("1.3.121"),
				Feed: PriceFeed{
					SettlementPrice: Price{
						Base:  AssetAmount{Amount: 1000, AssetID: MustParseObjectID("1.3.121")},
						Quote: AssetAmount{Amount: 100, AssetID: MustParseObjectID("1.3.0")},
					},
					MaintenanceCollateralRatio: 1750,
					MaximumShortSqueezeRatio:   1100,
					CoreExchangeRate: Price{
						Base:  AssetAmount{Amount: 1000, AssetID: MustParseObjectID("1.3.121")},
						Quote: AssetAmount{Amount: 100, AssetID: MustParseObjectID("1.3.0")},
					},
				},
				Extensions: AssetPublishFeedExtensions{
					InitialCollateralRatio: uint16Ptr(1900),
				},
			},
		},
		{
			name: "htlc_create",
			op: &HTLCCreateOperation{
				Fee:                AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
				From:               MustParseObjectID("1.2.1"),
				To:                 MustParseObjectID("1.2.2"),
				Amount:             AssetAmount{Amount: 50, AssetID: MustParseObjectID("1.3.0")},
				PreimageHash:       HTLCPreimageHash{Kind: 2, Value: "00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"},
				PreimageSize:       32,
				ClaimPeriodSeconds: 3600,
				Extensions: HTLCCreateExtensions{
					Memo: rawMemo,
				},
			},
		},
		{
			name: "credit_offer_accept",
			op: &CreditOfferAcceptOperation{
				Fee:                AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
				Borrower:           MustParseObjectID("1.2.1"),
				OfferID:            MustParseObjectID("1.21.7"),
				BorrowAmount:       AssetAmount{Amount: 100, AssetID: MustParseObjectID("1.3.1")},
				Collateral:         AssetAmount{Amount: 200, AssetID: MustParseObjectID("1.3.0")},
				MaxFeeRate:         25,
				MinDurationSeconds: 600,
				Extensions: CreditOfferAcceptExtensions{
					AutoRepay: uint8Ptr(2),
				},
			},
		},
	}

	for _, tt := range tests {
		raw, err := OperationEnvelope{Operation: tt.op}.MarshalBinary()
		if err != nil {
			t.Fatalf("%s MarshalBinary() error = %v", tt.name, err)
		}
		var decoded OperationEnvelope
		if err := decoded.UnmarshalBinary(raw); err != nil {
			t.Fatalf("%s UnmarshalBinary() error = %v", tt.name, err)
		}
		if decoded.Type() != tt.op.Type() {
			t.Fatalf("%s decoded type = %v, want %v", tt.name, decoded.Type(), tt.op.Type())
		}
	}
}

func TestCustomAuthorityRestrictionBinaryRoundTrip(t *testing.T) {
	t.Parallel()

	original := &CustomAuthorityCreateOperation{
		Fee:           AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
		Account:       MustParseObjectID("1.2.1"),
		Enabled:       true,
		ValidFrom:     NewTime(time.Unix(1700000000, 0)),
		ValidTo:       NewTime(time.Unix(1700003600, 0)),
		OperationType: uint64(OperationTypeTransfer),
		Auth: Authority{
			WeightThreshold: 1,
			AccountAuths: map[ObjectID]uint16{
				MustParseObjectID("1.2.9"): 1,
			},
		},
		Restrictions: []Restriction{
			{
				MemberIndex:     0,
				RestrictionType: 0,
				Argument: RestrictionArgument{
					Kind:  2,
					Int64: int64Ptr(42),
				},
			},
		},
	}

	raw, err := original.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary() error = %v", err)
	}

	var decodedEnvelope OperationEnvelope
	if err := decodedEnvelope.UnmarshalBinary(raw); err != nil {
		t.Fatalf("UnmarshalBinary() error = %v", err)
	}
	decoded, ok := decodedEnvelope.Operation.(*CustomAuthorityCreateOperation)
	if !ok {
		t.Fatalf("decoded operation type = %T, want *CustomAuthorityCreateOperation", decodedEnvelope.Operation)
	}
	if len(decoded.Restrictions) != 1 {
		t.Fatalf("Restrictions length = %d, want 1", len(decoded.Restrictions))
	}
	if decoded.Restrictions[0].Argument.Int64 == nil || *decoded.Restrictions[0].Argument.Int64 != 42 {
		t.Fatalf("decoded restriction = %#v", decoded.Restrictions[0])
	}
}

func TestAdditionalNonVirtualOperationsBinaryRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		op   Operation
	}{
		{
			name: "custom",
			op: &CustomOperation{
				Fee:           AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
				Payer:         MustParseObjectID("1.2.1"),
				RequiredAuths: []ObjectID{MustParseObjectID("1.2.9")},
				Id:            7,
				Data:          json.RawMessage(`"cafe"`),
			},
		},
		{
			name: "ticket_create",
			op: &TicketCreateOperation{
				Fee:        AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
				Account:    MustParseObjectID("1.2.1"),
				TargetType: 4,
				Amount:     AssetAmount{Amount: 1000, AssetID: MustParseObjectID("1.3.0")},
			},
		},
		{
			name: "liquidity_pool_create",
			op: &LiquidityPoolCreateOperation{
				Fee:                  AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
				Account:              MustParseObjectID("1.2.1"),
				AssetA:               MustParseObjectID("1.3.0"),
				AssetB:               MustParseObjectID("1.3.1"),
				ShareAsset:           MustParseObjectID("1.3.2"),
				TakerFeePercent:      30,
				WithdrawalFeePercent: 15,
			},
		},
		{
			name: "samet_fund_create",
			op: &SametFundCreateOperation{
				Fee:          AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
				OwnerAccount: MustParseObjectID("1.2.1"),
				AssetType:    MustParseObjectID("1.3.1"),
				Balance:      5000,
				FeeRate:      25,
			},
		},
		{
			name: "credit_offer_create",
			op: &CreditOfferCreateOperation{
				Fee:                AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
				OwnerAccount:       MustParseObjectID("1.2.1"),
				AssetType:          MustParseObjectID("1.3.1"),
				Balance:            10000,
				FeeRate:            30,
				MaxDurationSeconds: 3600,
				MinDealAmount:      100,
				Enabled:            true,
				AutoDisableTime:    NewTime(time.Unix(1700000000, 0)),
				AcceptableCollateral: []CreditOfferCollateral{
					{
						AssetID: MustParseObjectID("1.3.0"),
						Price: Price{
							Base:  AssetAmount{Amount: 2, AssetID: MustParseObjectID("1.3.0")},
							Quote: AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.1")},
						},
					},
				},
				AcceptableBorrowers: []CreditOfferBorrower{
					{AccountID: MustParseObjectID("1.2.9"), Amount: 123},
				},
			},
		},
		{
			name: "committee_member_update_global_parameters",
			op: &CommitteeMemberUpdateGlobalParametersOperation{
				Fee: AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
				NewParameters: ChainParameters{
					BlockInterval:       3,
					MaintenanceInterval: 3600,
					MaxAuthorityDepth:   5,
				},
			},
		},
	}

	for _, tt := range tests {
		raw, err := OperationEnvelope{Operation: tt.op}.MarshalBinary()
		if err != nil {
			t.Fatalf("%s MarshalBinary() error = %v", tt.name, err)
		}
		var decoded OperationEnvelope
		if err := decoded.UnmarshalBinary(raw); err != nil {
			t.Fatalf("%s UnmarshalBinary() error = %v", tt.name, err)
		}
		if decoded.Type() != tt.op.Type() {
			t.Fatalf("%s decoded type = %v, want %v", tt.name, decoded.Type(), tt.op.Type())
		}
	}
}

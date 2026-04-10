package protocol

import (
	"testing"
	"time"

	"github.com/froooze/go-bitshares/ecc"
)

func TestFillOrderBinaryRoundTrip(t *testing.T) {
	original := FillOrderOperation{
		Fee:       AssetAmount{Amount: 5, AssetID: MustParseObjectID("1.3.0")},
		OrderID:   MustParseObjectID("1.7.42"),
		AccountID: MustParseObjectID("1.2.9"),
		Pays:      AssetAmount{Amount: 11, AssetID: MustParseObjectID("1.3.0")},
		Receives:  AssetAmount{Amount: 7, AssetID: MustParseObjectID("1.3.1")},
		FillPrice: Price{
			Base:  AssetAmount{Amount: 11, AssetID: MustParseObjectID("1.3.0")},
			Quote: AssetAmount{Amount: 7, AssetID: MustParseObjectID("1.3.1")},
		},
		IsMaker: true,
	}

	raw, err := original.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary() error = %v", err)
	}

	var decodedEnvelope OperationEnvelope
	if err := decodedEnvelope.UnmarshalBinary(raw); err != nil {
		t.Fatalf("UnmarshalBinary() error = %v", err)
	}
	decoded, ok := decodedEnvelope.Operation.(*FillOrderOperation)
	if !ok {
		t.Fatalf("decoded operation type = %T, want *FillOrderOperation", decodedEnvelope.Operation)
	}
	if decoded.OrderID != original.OrderID || decoded.AccountID != original.AccountID {
		t.Fatalf("decoded fill order mismatch: %#v", decoded)
	}
	if decoded.FillPrice.Base.AssetID != original.FillPrice.Base.AssetID || decoded.FillPrice.Quote.AssetID != original.FillPrice.Quote.AssetID {
		t.Fatalf("decoded fill price mismatch: %#v", decoded.FillPrice)
	}
	if !decoded.IsMaker {
		t.Fatalf("decoded IsMaker = false, want true")
	}
}

func TestLimitOrderCreateExtensionsBinaryRoundTrip(t *testing.T) {
	original := &LimitOrderCreateOperation{
		Fee:          AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
		Seller:       MustParseObjectID("1.2.1"),
		AmountToSell: AssetAmount{Amount: 100, AssetID: MustParseObjectID("1.3.0")},
		MinToReceive: AssetAmount{Amount: 50, AssetID: MustParseObjectID("1.3.1")},
		Expiration:   NewTime(time.Unix(1700000000, 0)),
		FillOrKill:   true,
		Extensions: LimitOrderCreateExtensions{
			OnFill: []LimitOrderAutoAction{
				{
					Kind: 0,
					TakeProfit: &CreateTakeProfitOrderAction{
						FeeAssetID:        MustParseObjectID("1.3.0"),
						SpreadPercent:     200,
						SizePercent:       4000,
						ExpirationSeconds: 7200,
						Repeat:            false,
					},
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
	decoded, ok := envelope.Operation.(*LimitOrderCreateOperation)
	if !ok {
		t.Fatalf("decoded operation type = %T, want *LimitOrderCreateOperation", envelope.Operation)
	}
	if len(decoded.Extensions.OnFill) != 1 {
		t.Fatalf("OnFill length = %d, want 1", len(decoded.Extensions.OnFill))
	}
	if decoded.Extensions.OnFill[0].Kind != 0 || decoded.Extensions.OnFill[0].TakeProfit == nil {
		t.Fatalf("decoded on_fill = %#v", decoded.Extensions.OnFill[0])
	}
	if got, want := decoded.Extensions.OnFill[0].TakeProfit.SpreadPercent, original.Extensions.OnFill[0].TakeProfit.SpreadPercent; got != want {
		t.Fatalf("SpreadPercent = %d, want %d", got, want)
	}
}

func TestHTLCPreimageHashKind3BinaryRoundTrip(t *testing.T) {
	original := HTLCPreimageHash{Kind: 3, Value: "00112233445566778899aabbccddeeff00112233"}
	raw, err := original.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary() error = %v", err)
	}

	var decoded HTLCPreimageHash
	if err := decoded.UnmarshalBinaryFrom(newBinaryReader(raw)); err != nil {
		t.Fatalf("UnmarshalBinaryFrom() error = %v", err)
	}
	if decoded.Kind != original.Kind || decoded.Value != original.Value {
		t.Fatalf("decoded hash mismatch: %#v", decoded)
	}
}

func TestAccountCreateExtensionsBinaryRoundTrip(t *testing.T) {
	pub := MustPublicKey(ecc.PrivateKeyFromSeed([]byte("account-create-extension")).PublicKey().String())
	original := AccountCreateOperation{
		Fee:             AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
		Registrar:       MustParseObjectID("1.2.1"),
		Referrer:        MustParseObjectID("1.2.2"),
		ReferrerPercent: 1000,
		Name:            "alice",
		Owner: Authority{
			WeightThreshold: 1,
			KeyAuths: map[PublicKey]uint16{
				pub: 1,
			},
		},
		Active: Authority{
			WeightThreshold: 1,
			AccountAuths: map[ObjectID]uint16{
				MustParseObjectID("1.2.3"): 1,
			},
		},
		Options: AccountOptions{
			MemoKey:       pub,
			VotingAccount: MustParseObjectID("1.2.5"),
			NumWitness:    1,
			NumCommittee:  1,
			Votes:         []VoteID{{Type: 1, ID: 2}},
		},
		Extensions: AccountCreateExtensions{
			OwnerSpecialAuthority: &SpecialAuthority{
				Kind: 1,
				TopHolders: &TopHoldersSpecialAuthority{
					Asset:         MustParseObjectID("1.3.0"),
					NumTopHolders: 7,
				},
			},
			ActiveSpecialAuthority: &SpecialAuthority{
				Kind: 0,
				No:   &NoSpecialAuthority{},
			},
			BuybackOptions: &BuybackAccountOptions{
				AssetToBuy:       MustParseObjectID("1.3.0"),
				AssetToBuyIssuer: MustParseObjectID("1.2.9"),
				Markets: []ObjectID{
					MustParseObjectID("1.3.1"),
					MustParseObjectID("1.3.2"),
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
	decoded, ok := decodedEnvelope.Operation.(*AccountCreateOperation)
	if !ok {
		t.Fatalf("decoded operation type = %T, want *AccountCreateOperation", decodedEnvelope.Operation)
	}
	if decoded.Extensions.OwnerSpecialAuthority == nil || decoded.Extensions.OwnerSpecialAuthority.Kind != 1 {
		t.Fatalf("decoded owner special authority mismatch: %#v", decoded.Extensions)
	}
	if decoded.Extensions.BuybackOptions == nil || len(decoded.Extensions.BuybackOptions.Markets) != 2 {
		t.Fatalf("decoded buyback options mismatch: %#v", decoded.Extensions)
	}
}

func TestAssetAndChainExtensionsBinaryRoundTrip(t *testing.T) {
	assetUpdate := AssetUpdateOperation{
		Fee:           AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
		Issuer:        MustParseObjectID("1.2.1"),
		AssetToUpdate: MustParseObjectID("1.3.1"),
		NewOptions:    AssetOptions{MaxSupply: 1000, CoreExchangeRate: Price{Base: AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")}, Quote: AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.1")}}},
		Extensions: AssetUpdateExtensions{
			NewPrecision:         uint8Ptr(8),
			SkipCoreExchangeRate: boolPtr(true),
		},
	}
	raw, err := assetUpdate.MarshalBinary()
	if err != nil {
		t.Fatalf("asset update MarshalBinary() error = %v", err)
	}
	var assetEnvelope OperationEnvelope
	if err := assetEnvelope.UnmarshalBinary(raw); err != nil {
		t.Fatalf("asset update UnmarshalBinary() error = %v", err)
	}
	decodedAssetUpdate, ok := assetEnvelope.Operation.(*AssetUpdateOperation)
	if !ok {
		t.Fatalf("decoded operation type = %T, want *AssetUpdateOperation", assetEnvelope.Operation)
	}
	if decodedAssetUpdate.Extensions.NewPrecision == nil || *decodedAssetUpdate.Extensions.NewPrecision != 8 {
		t.Fatalf("decoded asset update precision mismatch: %#v", decodedAssetUpdate.Extensions)
	}
	if decodedAssetUpdate.Extensions.SkipCoreExchangeRate == nil || !*decodedAssetUpdate.Extensions.SkipCoreExchangeRate {
		t.Fatalf("decoded asset update skip_core_exchange_rate mismatch: %#v", decodedAssetUpdate.Extensions)
	}

	params := ChainParameters{
		BlockInterval:       3,
		MaintenanceInterval: 3600,
		MaxAuthorityDepth:   5,
		Extensions: ChainParametersExtensions{
			UpdatableHTLCOptions: &HTLCOptions{
				MaxTimeoutSecs:  86400,
				MaxPreimageSize: 64,
			},
			CustomAuthorityOptions: &CustomAuthorityOptions{
				MaxCustomAuthorityLifetimeSeconds: 100,
				MaxCustomAuthoritiesPerAccount:    10,
				MaxCustomAuthoritiesPerAccountOp:  3,
				MaxCustomAuthorityRestrictions:    8,
			},
			MarketFeeNetworkPercent: uint16Ptr(15),
			MakerFeeDiscountPercent: uint16Ptr(5),
		},
	}
	rawParams, err := params.MarshalBinary()
	if err != nil {
		t.Fatalf("chain params MarshalBinary() error = %v", err)
	}
	var decodedParams ChainParameters
	if err := decodedParams.UnmarshalBinaryFrom(newBinaryReader(rawParams)); err != nil {
		t.Fatalf("chain params UnmarshalBinaryFrom() error = %v", err)
	}
	if decodedParams.Extensions.UpdatableHTLCOptions == nil || decodedParams.Extensions.CustomAuthorityOptions == nil {
		t.Fatalf("decoded chain parameter extensions mismatch: %#v", decodedParams.Extensions)
	}
	if got, want := decodedParams.Extensions.MarketFeeNetworkPercent, uint16Ptr(15); got == nil || *got != *want {
		t.Fatalf("market_fee_network_percent mismatch: %#v", decodedParams.Extensions)
	}
}

func TestFeeScheduleBinaryRoundTrip(t *testing.T) {
	original := FeeSchedule{
		Scale: 100,
		Parameters: []FeeScheduleParameter{
			{
				OperationType: OperationTypeTransfer,
				Value: &FeeAndPricePerKbyteParameters{
					Fee:           123,
					PricePerKbyte: 7,
				},
			},
			{
				OperationType: OperationTypeFillOrder,
				Value:         &EmptyFeeParameters{},
			},
		},
	}

	w := newBinaryWriter()
	if err := original.MarshalBinaryInto(w); err != nil {
		t.Fatalf("MarshalBinaryInto() error = %v", err)
	}
	raw := w.Bytes()

	var decoded FeeSchedule
	if err := decoded.UnmarshalBinaryFrom(newBinaryReader(raw)); err != nil {
		t.Fatalf("UnmarshalBinaryFrom() error = %v", err)
	}
	if got, want := decoded.Scale, original.Scale; got != want {
		t.Fatalf("Scale = %d, want %d", got, want)
	}
	if len(decoded.Parameters) != 2 {
		t.Fatalf("Parameters length = %d, want 2", len(decoded.Parameters))
	}
	if got, ok := decoded.Parameters[0].Value.(*FeeAndPricePerKbyteParameters); !ok || got.Fee != 123 || got.PricePerKbyte != 7 {
		t.Fatalf("decoded transfer fee parameter mismatch: %#v", decoded.Parameters[0].Value)
	}
	if _, ok := decoded.Parameters[1].Value.(*EmptyFeeParameters); !ok {
		t.Fatalf("decoded fill_order fee parameter type = %T, want *EmptyFeeParameters", decoded.Parameters[1].Value)
	}
}

func boolPtr(v bool) *bool { return &v }

package bitshares

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/froooze/go-bitshares/ecc"
	"github.com/froooze/go-bitshares/protocol"
)

func TestSignedAmountFromFloatSupportsNegativeValues(t *testing.T) {
	t.Parallel()

	got, err := signedAmountFromFloat(-1.25, 2)
	if err != nil {
		t.Fatalf("signedAmountFromFloat() error = %v", err)
	}
	if want := int64(-125); got != want {
		t.Fatalf("signedAmountFromFloat(-1.25, 2) = %d, want %d", got, want)
	}
}

func TestBuildAccountCreateOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	ownerKey := protocol.MustPublicKey(ecc.PrivateKeyFromSeed([]byte("account-create-owner")).PublicKey().String())
	activeKey := protocol.MustPublicKey(ecc.PrivateKeyFromSeed([]byte("account-create-active")).PublicKey().String())
	op, err := buildAccountCreateOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.678"),
		2500,
		"alice",
		protocol.Authority{
			WeightThreshold: 1,
			KeyAuths:        map[protocol.PublicKey]uint16{ownerKey: 1},
		},
		protocol.Authority{
			WeightThreshold: 1,
			KeyAuths:        map[protocol.PublicKey]uint16{activeKey: 1},
		},
		protocol.AccountOptions{
			MemoKey:       activeKey,
			VotingAccount: protocol.MustParseObjectID("1.2.5"),
		},
		nil,
	)
	if err != nil {
		t.Fatalf("buildAccountCreateOperation() error = %v", err)
	}
	if got, want := op.Name, "alice"; got != want {
		t.Fatalf("name = %q, want %q", got, want)
	}
	if got, want := op.ReferrerPercent, uint16(2500); got != want {
		t.Fatalf("referrer percent = %d, want %d", got, want)
	}
}

func TestBuildAccountCreateOperationRejectsInvalidAccountName(t *testing.T) {
	t.Parallel()

	key := protocol.MustPublicKey(ecc.PrivateKeyFromSeed([]byte("account-create-invalid")).PublicKey().String())
	_, err := buildAccountCreateOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.678"),
		0,
		"Alice",
		protocol.Authority{WeightThreshold: 1, KeyAuths: map[protocol.PublicKey]uint16{key: 1}},
		protocol.Authority{WeightThreshold: 1, KeyAuths: map[protocol.PublicKey]uint16{key: 1}},
		protocol.AccountOptions{MemoKey: key, VotingAccount: protocol.MustParseObjectID("1.2.5")},
		nil,
	)
	if err == nil {
		t.Fatal("buildAccountCreateOperation() error = nil, want error")
	}
}

func TestBuildAccountCreateOperationRejectsImpossibleAuthority(t *testing.T) {
	t.Parallel()

	key := protocol.MustPublicKey(ecc.PrivateKeyFromSeed([]byte("account-create-impossible")).PublicKey().String())
	_, err := buildAccountCreateOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.678"),
		0,
		"alice",
		protocol.Authority{WeightThreshold: 2, KeyAuths: map[protocol.PublicKey]uint16{key: 1}},
		protocol.Authority{WeightThreshold: 1, KeyAuths: map[protocol.PublicKey]uint16{key: 1}},
		protocol.AccountOptions{MemoKey: key, VotingAccount: protocol.MustParseObjectID("1.2.5")},
		nil,
	)
	if err == nil {
		t.Fatal("buildAccountCreateOperation() error = nil, want error")
	}
}

func TestBuildAssetCreateOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	op, err := buildAssetCreateOperation(
		protocol.MustParseObjectID("1.2.345"),
		"HONEST.USD",
		5,
		protocol.AssetOptions{
			MaxSupply:        1000000,
			MarketFeePercent: 25,
			MaxMarketFee:     1000,
			CoreExchangeRate: protocol.Price{
				Base:  protocol.AssetAmount{Amount: 1000, AssetID: protocol.MustParseObjectID("1.3.0")},
				Quote: protocol.AssetAmount{Amount: 1000, AssetID: protocol.MustParseObjectID("1.3.121")},
			},
		},
		nil,
		false,
	)
	if err != nil {
		t.Fatalf("buildAssetCreateOperation() error = %v", err)
	}
	if got, want := op.Symbol, "HONEST.USD"; got != want {
		t.Fatalf("symbol = %q, want %q", got, want)
	}
	if got, want := op.Precision, uint8(5); got != want {
		t.Fatalf("precision = %d, want %d", got, want)
	}
}

func TestBuildAssetCreateOperationRejectsInvalidSymbol(t *testing.T) {
	t.Parallel()

	_, err := buildAssetCreateOperation(
		protocol.MustParseObjectID("1.2.345"),
		"bitusd",
		5,
		protocol.AssetOptions{
			MaxSupply:        1000000,
			CoreExchangeRate: protocol.Price{Base: protocol.AssetAmount{Amount: 1000, AssetID: protocol.MustParseObjectID("1.3.0")}, Quote: protocol.AssetAmount{Amount: 1000, AssetID: protocol.MustParseObjectID("1.3.121")}},
		},
		nil,
		false,
	)
	if err == nil {
		t.Fatal("buildAssetCreateOperation() error = nil, want error")
	}
}

func TestBuildAssetCreateOperationRejectsPredictionMarketWithoutBitassetOptions(t *testing.T) {
	t.Parallel()

	_, err := buildAssetCreateOperation(
		protocol.MustParseObjectID("1.2.345"),
		"PMARKET",
		5,
		protocol.AssetOptions{
			MaxSupply:         1000000,
			IssuerPermissions: assetIssuerPermissionGlobalSettle,
			CoreExchangeRate:  protocol.Price{Base: protocol.AssetAmount{Amount: 1000, AssetID: protocol.MustParseObjectID("1.3.0")}, Quote: protocol.AssetAmount{Amount: 1000, AssetID: protocol.MustParseObjectID("1.3.121")}},
			MarketFeePercent:  0,
			MaxMarketFee:      0,
		},
		nil,
		true,
	)
	if err == nil {
		t.Fatal("buildAssetCreateOperation() error = nil, want error")
	}
}

func TestBuildAccountUpdateOperationRequiresChanges(t *testing.T) {
	t.Parallel()

	_, err := buildAccountUpdateOperation(protocol.MustParseObjectID("1.2.345"), nil, nil, nil, nil)
	if err == nil {
		t.Fatal("buildAccountUpdateOperation() error = nil, want error")
	}
}

func TestBuildAccountUpdateOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	key := protocol.MustPublicKey(ecc.PrivateKeyFromSeed([]byte("account-update-memo")).PublicKey().String())
	options := &protocol.AccountOptions{
		MemoKey:       key,
		VotingAccount: protocol.MustParseObjectID("1.2.5"),
	}
	owner := &protocol.Authority{WeightThreshold: 1}

	op, err := buildAccountUpdateOperation(protocol.MustParseObjectID("1.2.345"), owner, nil, options, &protocol.AccountUpdateExtensions{})
	if err != nil {
		t.Fatalf("buildAccountUpdateOperation() error = %v", err)
	}
	if got, want := op.Account.String(), "1.2.345"; got != want {
		t.Fatalf("account = %q, want %q", got, want)
	}
	if op.Owner == nil || op.NewOptions == nil {
		t.Fatal("expected owner and new options to be preserved")
	}
}

func TestBuildCallOrderUpdateOperationSupportsSignedDeltas(t *testing.T) {
	t.Parallel()

	tcr := uint16(1750)
	op, err := buildCallOrderUpdateOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.3.0"),
		5,
		-1.25,
		protocol.MustParseObjectID("1.3.1"),
		5,
		2.5,
		&tcr,
	)
	if err != nil {
		t.Fatalf("buildCallOrderUpdateOperation() error = %v", err)
	}
	if got, want := op.DeltaDebt.Amount, int64(-125000); got != want {
		t.Fatalf("delta debt = %d, want %d", got, want)
	}
	if got, want := op.DeltaCollateral.Amount, int64(250000); got != want {
		t.Fatalf("delta collateral = %d, want %d", got, want)
	}
	if op.Extensions.TargetCollateralRatio == nil || *op.Extensions.TargetCollateralRatio != tcr {
		t.Fatal("target collateral ratio was not preserved")
	}
}

func TestBuildAccountUpgradeOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	op := buildAccountUpgradeOperation(protocol.MustParseObjectID("1.2.345"), true)
	if got, want := op.AccountToUpgrade.String(), "1.2.345"; got != want {
		t.Fatalf("account = %q, want %q", got, want)
	}
	if !op.UpgradeToLifetimeMember {
		t.Fatal("expected lifetime member flag to be true")
	}
}

func TestBuildAssetFundFeePoolOperationUsesCorePrecision(t *testing.T) {
	t.Parallel()

	op, err := buildAssetFundFeePoolOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.3.121"),
		5,
		1.25,
	)
	if err != nil {
		t.Fatalf("buildAssetFundFeePoolOperation() error = %v", err)
	}
	if got, want := op.Amount, int64(125000); got != want {
		t.Fatalf("amount = %d, want %d", got, want)
	}
}

func TestBuildAssetSettleOperationPopulatesAmount(t *testing.T) {
	t.Parallel()

	op, err := buildAssetSettleOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.3.121"),
		4,
		2.5,
	)
	if err != nil {
		t.Fatalf("buildAssetSettleOperation() error = %v", err)
	}
	if got, want := op.Amount.Amount, int64(25000); got != want {
		t.Fatalf("amount = %d, want %d", got, want)
	}
}

func TestBuildAssetSettleOperationRejectsZeroAmount(t *testing.T) {
	t.Parallel()

	_, err := buildAssetSettleOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.3.121"),
		4,
		0,
	)
	if err == nil {
		t.Fatal("buildAssetSettleOperation() error = nil, want error")
	}
}

func TestBuildAssetUpdateFeedProducersOperationRequiresProducers(t *testing.T) {
	t.Parallel()

	_, err := buildAssetUpdateFeedProducersOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.3.121"),
		nil,
	)
	if err == nil {
		t.Fatal("buildAssetUpdateFeedProducersOperation() error = nil, want error")
	}
}

func TestBuildAssetUpdateFeedProducersOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	op, err := buildAssetUpdateFeedProducersOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.3.121"),
		[]protocol.ObjectID{protocol.MustParseObjectID("1.2.10"), protocol.MustParseObjectID("1.2.11")},
	)
	if err != nil {
		t.Fatalf("buildAssetUpdateFeedProducersOperation() error = %v", err)
	}
	if got, want := len(op.NewFeedProducers), 2; got != want {
		t.Fatalf("producer count = %d, want %d", got, want)
	}
}

func TestBuildAssetUpdateOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	newIssuer := protocol.MustParseObjectID("1.2.999")
	op := buildAssetUpdateOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.3.121"),
		protocol.AssetOptions{
			MaxSupply:        1000000,
			MarketFeePercent: 25,
			MaxMarketFee:     1000,
			CoreExchangeRate: protocol.Price{
				Base:  protocol.AssetAmount{Amount: 1000, AssetID: protocol.MustParseObjectID("1.3.121")},
				Quote: protocol.AssetAmount{Amount: 1000, AssetID: protocol.MustParseObjectID("1.3.0")},
			},
		},
		&newIssuer,
		&protocol.AssetUpdateExtensions{},
	)
	if got, want := op.AssetToUpdate.String(), "1.3.121"; got != want {
		t.Fatalf("asset id = %q, want %q", got, want)
	}
	if op.NewIssuer == nil || op.NewIssuer.String() != "1.2.999" {
		t.Fatal("expected new issuer to be preserved")
	}
}

func TestBuildAssetUpdateBitassetOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	op := buildAssetUpdateBitassetOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.3.121"),
		protocol.BitAssetOptions{
			FeedLifetimeSec:              3600,
			MinimumFeeds:                 7,
			ForceSettlementDelaySec:      600,
			ForceSettlementOffsetPercent: 100,
			MaximumForceSettlementVolume: 20,
			ShortBackingAsset:            protocol.MustParseObjectID("1.3.0"),
		},
	)
	if got, want := op.NewOptions.MinimumFeeds, uint8(7); got != want {
		t.Fatalf("minimum feeds = %d, want %d", got, want)
	}
}

func TestBuildAssetGlobalSettleOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	op, err := buildAssetGlobalSettleOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.3.121"),
		protocol.Price{
			Base:  protocol.AssetAmount{Amount: 1000, AssetID: protocol.MustParseObjectID("1.3.121")},
			Quote: protocol.AssetAmount{Amount: 2000, AssetID: protocol.MustParseObjectID("1.3.0")},
		},
	)
	if err != nil {
		t.Fatalf("buildAssetGlobalSettleOperation() error = %v", err)
	}
	if got, want := op.AssetToSettle.String(), "1.3.121"; got != want {
		t.Fatalf("asset to settle = %q, want %q", got, want)
	}
}

func TestBuildAssetGlobalSettleOperationRejectsMismatchedBaseAsset(t *testing.T) {
	t.Parallel()

	_, err := buildAssetGlobalSettleOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.3.121"),
		protocol.Price{
			Base:  protocol.AssetAmount{Amount: 1000, AssetID: protocol.MustParseObjectID("1.3.999")},
			Quote: protocol.AssetAmount{Amount: 2000, AssetID: protocol.MustParseObjectID("1.3.0")},
		},
	)
	if err == nil {
		t.Fatal("buildAssetGlobalSettleOperation() error = nil, want error")
	}
}

func TestBuildAssetPublishFeedOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	icr := uint16(1800)
	op := buildAssetPublishFeedOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.3.121"),
		protocol.PriceFeed{
			SettlementPrice: protocol.Price{
				Base:  protocol.AssetAmount{Amount: 1000, AssetID: protocol.MustParseObjectID("1.3.121")},
				Quote: protocol.AssetAmount{Amount: 2000, AssetID: protocol.MustParseObjectID("1.3.0")},
			},
			MaintenanceCollateralRatio: 1750,
			MaximumShortSqueezeRatio:   1100,
			CoreExchangeRate: protocol.Price{
				Base:  protocol.AssetAmount{Amount: 1000, AssetID: protocol.MustParseObjectID("1.3.121")},
				Quote: protocol.AssetAmount{Amount: 1500, AssetID: protocol.MustParseObjectID("1.3.0")},
			},
		},
		&protocol.AssetPublishFeedExtensions{InitialCollateralRatio: &icr},
	)
	if got, want := op.AssetID.String(), "1.3.121"; got != want {
		t.Fatalf("asset id = %q, want %q", got, want)
	}
	if op.Extensions.InitialCollateralRatio == nil || *op.Extensions.InitialCollateralRatio != icr {
		t.Fatal("initial collateral ratio was not preserved")
	}
}

func TestBuildProposalCreateOperationWrapsOperations(t *testing.T) {
	t.Parallel()

	review := uint32(3600)
	expiration := time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC)
	proposedTransfer := &protocol.TransferOperation{
		From:   protocol.MustParseObjectID("1.2.345"),
		To:     protocol.MustParseObjectID("1.2.678"),
		Amount: protocol.AssetAmount{Amount: 1000, AssetID: protocol.MustParseObjectID("1.3.0")},
	}

	op, err := buildProposalCreateOperation(protocol.MustParseObjectID("1.2.345"), expiration, &review, proposedTransfer)
	if err != nil {
		t.Fatalf("buildProposalCreateOperation() error = %v", err)
	}
	if got, want := len(op.ProposedOps), 1; got != want {
		t.Fatalf("proposed op count = %d, want %d", got, want)
	}
	if got, want := op.ExpirationTime.Time, expiration.UTC(); !got.Equal(want) {
		t.Fatalf("expiration = %s, want %s", got.Format(time.RFC3339), want.Format(time.RFC3339))
	}
}

func TestBuildProposalCreateOperationRejectsEmptyOps(t *testing.T) {
	t.Parallel()

	_, err := buildProposalCreateOperation(protocol.MustParseObjectID("1.2.345"), time.Now().UTC(), nil)
	if err == nil {
		t.Fatal("buildProposalCreateOperation() error = nil, want error")
	}
}

func TestBuildAccountWhitelistOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	op, err := buildAccountWhitelistOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.678"),
		3,
	)
	if err != nil {
		t.Fatalf("buildAccountWhitelistOperation() error = %v", err)
	}
	if got, want := op.AuthorizingAccount.String(), "1.2.345"; got != want {
		t.Fatalf("authorizing account = %q, want %q", got, want)
	}
	if got, want := op.AccountToList.String(), "1.2.678"; got != want {
		t.Fatalf("account to list = %q, want %q", got, want)
	}
	if got, want := op.NewListing, uint8(3); got != want {
		t.Fatalf("new listing = %d, want %d", got, want)
	}
}

func TestBuildAccountWhitelistOperationRejectsInvalidListing(t *testing.T) {
	t.Parallel()

	_, err := buildAccountWhitelistOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.678"),
		4,
	)
	if err == nil {
		t.Fatal("buildAccountWhitelistOperation() error = nil, want error")
	}
}

func TestBuildProposalUpdateOperationRejectsEmptyChanges(t *testing.T) {
	t.Parallel()

	_, err := buildProposalUpdateOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.10.999"),
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	if err == nil {
		t.Fatal("buildProposalUpdateOperation() error = nil, want error")
	}
}

func TestBuildProposalUpdateOperationRejectsOverlappingApprovals(t *testing.T) {
	t.Parallel()

	key := protocol.MustPublicKey(ecc.PrivateKeyFromSeed([]byte("proposal-update-key")).PublicKey().String())
	_, err := buildProposalUpdateOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.10.999"),
		[]protocol.ObjectID{protocol.MustParseObjectID("1.2.101")},
		[]protocol.ObjectID{protocol.MustParseObjectID("1.2.101")},
		nil,
		nil,
		[]protocol.PublicKey{key},
		[]protocol.PublicKey{key},
	)
	if err == nil {
		t.Fatal("buildProposalUpdateOperation() error = nil, want error")
	}
}

func TestBuildProposalUpdateOperationPopulatesApprovals(t *testing.T) {
	t.Parallel()

	keyAdd := protocol.MustPublicKey(ecc.PrivateKeyFromSeed([]byte("proposal-update-add")).PublicKey().String())
	keyRemove := protocol.MustPublicKey(ecc.PrivateKeyFromSeed([]byte("proposal-update-remove")).PublicKey().String())
	op, err := buildProposalUpdateOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.10.999"),
		[]protocol.ObjectID{protocol.MustParseObjectID("1.2.101")},
		[]protocol.ObjectID{protocol.MustParseObjectID("1.2.102")},
		[]protocol.ObjectID{protocol.MustParseObjectID("1.2.201")},
		[]protocol.ObjectID{protocol.MustParseObjectID("1.2.202")},
		[]protocol.PublicKey{keyAdd},
		[]protocol.PublicKey{keyRemove},
	)
	if err != nil {
		t.Fatalf("buildProposalUpdateOperation() error = %v", err)
	}
	if got, want := op.Proposal.String(), "1.10.999"; got != want {
		t.Fatalf("proposal = %q, want %q", got, want)
	}
	if got, want := len(op.ActiveApprovalsToAdd), 1; got != want {
		t.Fatalf("active approvals to add = %d, want %d", got, want)
	}
	if got, want := len(op.KeyApprovalsToRemove), 1; got != want {
		t.Fatalf("key approvals to remove = %d, want %d", got, want)
	}
}

func TestBuildProposalDeleteOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	op := buildProposalDeleteOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.10.999"),
		true,
	)
	if got, want := op.FeePayingAccount.String(), "1.2.345"; got != want {
		t.Fatalf("fee paying account = %q, want %q", got, want)
	}
	if got, want := op.Proposal.String(), "1.10.999"; got != want {
		t.Fatalf("proposal = %q, want %q", got, want)
	}
	if !op.UsingOwnerAuthority {
		t.Fatal("expected using owner authority to be true")
	}
}

func TestBuildWithdrawPermissionCreateOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	start := time.Date(2026, 4, 12, 8, 0, 0, 0, time.UTC)
	op, err := buildWithdrawPermissionCreateOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.678"),
		protocol.MustParseObjectID("1.3.121"),
		5,
		1.25,
		3600,
		10,
		start,
	)
	if err != nil {
		t.Fatalf("buildWithdrawPermissionCreateOperation() error = %v", err)
	}
	if got, want := op.WithdrawalLimit.Amount, int64(125000); got != want {
		t.Fatalf("withdrawal limit = %d, want %d", got, want)
	}
	if got, want := op.PeriodStartTime.Time, start; !got.Equal(want) {
		t.Fatalf("period start = %s, want %s", got.Format(time.RFC3339), want.Format(time.RFC3339))
	}
}

func TestBuildWithdrawPermissionCreateOperationRejectsInvalidPayload(t *testing.T) {
	t.Parallel()

	_, err := buildWithdrawPermissionCreateOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.3.121"),
		5,
		1.25,
		3600,
		10,
		time.Date(2026, 4, 12, 8, 0, 0, 0, time.UTC),
	)
	if err == nil {
		t.Fatal("buildWithdrawPermissionCreateOperation() error = nil, want error")
	}
}

func TestBuildWithdrawPermissionUpdateOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	start := time.Date(2026, 4, 13, 8, 0, 0, 0, time.UTC)
	op, err := buildWithdrawPermissionUpdateOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.678"),
		protocol.MustParseObjectID("1.9.1"),
		protocol.MustParseObjectID("1.3.121"),
		5,
		2.5,
		7200,
		15,
		start,
	)
	if err != nil {
		t.Fatalf("buildWithdrawPermissionUpdateOperation() error = %v", err)
	}
	if got, want := op.PermissionToUpdate.String(), "1.9.1"; got != want {
		t.Fatalf("permission id = %q, want %q", got, want)
	}
	if got, want := op.WithdrawalLimit.Amount, int64(250000); got != want {
		t.Fatalf("withdrawal limit = %d, want %d", got, want)
	}
}

func TestBuildWithdrawPermissionClaimOperationCarriesMemo(t *testing.T) {
	t.Parallel()

	memo, err := json.Marshal(protocol.MemoData{
		From:    ecc.PrivateKeyFromSeed([]byte("withdraw-claim-from")).PublicKey().String(),
		To:      ecc.PrivateKeyFromSeed([]byte("withdraw-claim-to")).PublicKey().String(),
		Nonce:   "7",
		Message: "abcd",
	})
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	op, err := buildWithdrawPermissionClaimOperation(
		protocol.MustParseObjectID("1.9.1"),
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.678"),
		protocol.MustParseObjectID("1.3.121"),
		5,
		1.75,
		memo,
	)
	if err != nil {
		t.Fatalf("buildWithdrawPermissionClaimOperation() error = %v", err)
	}
	if got, want := op.AmountToWithdraw.Amount, int64(175000); got != want {
		t.Fatalf("amount to withdraw = %d, want %d", got, want)
	}
	if string(op.Memo) == "" {
		t.Fatal("expected memo to be preserved")
	}
}

func TestBuildWithdrawPermissionDeleteOperationRejectsSameAccounts(t *testing.T) {
	t.Parallel()

	_, err := buildWithdrawPermissionDeleteOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.9.1"),
	)
	if err == nil {
		t.Fatal("buildWithdrawPermissionDeleteOperation() error = nil, want error")
	}
}

func TestBuildWithdrawPermissionDeleteOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	op, err := buildWithdrawPermissionDeleteOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.678"),
		protocol.MustParseObjectID("1.9.1"),
	)
	if err != nil {
		t.Fatalf("buildWithdrawPermissionDeleteOperation() error = %v", err)
	}
	if got, want := op.WithdrawalPermission.String(), "1.9.1"; got != want {
		t.Fatalf("withdrawal permission = %q, want %q", got, want)
	}
}

func TestBuildAccountTransferOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	op := buildAccountTransferOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.678"),
	)
	if got, want := op.AccountID.String(), "1.2.345"; got != want {
		t.Fatalf("account id = %q, want %q", got, want)
	}
	if got, want := op.NewOwner.String(), "1.2.678"; got != want {
		t.Fatalf("new owner = %q, want %q", got, want)
	}
}

func TestBuildBalanceClaimOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	key := protocol.MustPublicKey(ecc.PrivateKeyFromSeed([]byte("balance-claim")).PublicKey().String())
	op, err := buildBalanceClaimOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.15.1"),
		key,
		protocol.MustParseObjectID("1.3.121"),
		5,
		1.5,
	)
	if err != nil {
		t.Fatalf("buildBalanceClaimOperation() error = %v", err)
	}
	if got, want := op.TotalClaimed.Amount, int64(150000); got != want {
		t.Fatalf("total claimed = %d, want %d", got, want)
	}
	if got, want := op.Fee.AssetID.String(), "1.3.0"; got != want {
		t.Fatalf("fee asset id = %q, want %q", got, want)
	}
}

func TestBuildOverrideTransferOperationRejectsIssuerFromMatch(t *testing.T) {
	t.Parallel()

	_, err := buildOverrideTransferOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.678"),
		protocol.MustParseObjectID("1.3.121"),
		5,
		1.5,
		nil,
	)
	if err == nil {
		t.Fatal("buildOverrideTransferOperation() error = nil, want error")
	}
}

func TestBuildOverrideTransferOperationCarriesMemo(t *testing.T) {
	t.Parallel()

	memo, err := json.Marshal(protocol.MemoData{
		From:    ecc.PrivateKeyFromSeed([]byte("override-from")).PublicKey().String(),
		To:      ecc.PrivateKeyFromSeed([]byte("override-to")).PublicKey().String(),
		Nonce:   "9",
		Message: "beef",
	})
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	op, err := buildOverrideTransferOperation(
		protocol.MustParseObjectID("1.2.999"),
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.678"),
		protocol.MustParseObjectID("1.3.121"),
		5,
		1.5,
		memo,
	)
	if err != nil {
		t.Fatalf("buildOverrideTransferOperation() error = %v", err)
	}
	if got, want := op.Amount.Amount, int64(150000); got != want {
		t.Fatalf("amount = %d, want %d", got, want)
	}
	if string(op.Memo) == "" {
		t.Fatal("expected memo to be preserved")
	}
}

func TestBuildBidCollateralOperationRejectsNoop(t *testing.T) {
	t.Parallel()

	_, err := buildBidCollateralOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.3.0"),
		5,
		0,
		protocol.MustParseObjectID("1.3.121"),
		5,
		0,
	)
	if err == nil {
		t.Fatal("buildBidCollateralOperation() error = nil, want error")
	}
}

func TestBuildBidCollateralOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	op, err := buildBidCollateralOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.3.0"),
		5,
		3.5,
		protocol.MustParseObjectID("1.3.121"),
		5,
		1.25,
	)
	if err != nil {
		t.Fatalf("buildBidCollateralOperation() error = %v", err)
	}
	if got, want := op.AdditionalCollateral.Amount, int64(350000); got != want {
		t.Fatalf("additional collateral = %d, want %d", got, want)
	}
	if got, want := op.DebtCovered.Amount, int64(125000); got != want {
		t.Fatalf("debt covered = %d, want %d", got, want)
	}
}

func TestBuildCommitteeMemberUpdateGlobalParametersOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	op := buildCommitteeMemberUpdateGlobalParametersOperation(protocol.ChainParameters{
		BlockInterval:                    3,
		MaintenanceInterval:              3600,
		MaintenanceSkipSlots:             3,
		CommitteeProposalReviewPeriod:    3600,
		MaximumTransactionSize:           2048,
		MaximumBlockSize:                 65536,
		MaximumTimeUntilExpiration:       86400,
		MaximumProposalLifetime:          86400,
		MaximumAssetWhitelistAuthorities: 10,
		MaximumAssetFeedPublishers:       25,
		MaximumWitnessCount:              21,
		MaximumCommitteeCount:            11,
	})
	if got, want := op.NewParameters.BlockInterval, uint8(3); got != want {
		t.Fatalf("block interval = %d, want %d", got, want)
	}
}

func TestBuildWitnessCreateOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	key := protocol.MustPublicKey(ecc.PrivateKeyFromSeed([]byte("witness-create")).PublicKey().String())
	op, err := buildWitnessCreateOperation(protocol.MustParseObjectID("1.2.345"), "https://example.com/witness", key)
	if err != nil {
		t.Fatalf("buildWitnessCreateOperation() error = %v", err)
	}
	if got, want := op.WitnessAccount.String(), "1.2.345"; got != want {
		t.Fatalf("witness account = %q, want %q", got, want)
	}
	if got, want := op.BlockSigningKey.String(), key.String(); got != want {
		t.Fatalf("block signing key = %q, want %q", got, want)
	}
}

func TestBuildWitnessCreateOperationRejectsLongURL(t *testing.T) {
	t.Parallel()

	key := protocol.MustPublicKey(ecc.PrivateKeyFromSeed([]byte("witness-create-long-url")).PublicKey().String())
	_, err := buildWitnessCreateOperation(
		protocol.MustParseObjectID("1.2.345"),
		strings.Repeat("x", maxChainURLLength),
		key,
	)
	if err == nil {
		t.Fatal("buildWitnessCreateOperation() error = nil, want error")
	}
}

func TestBuildWitnessUpdateOperationRequiresChanges(t *testing.T) {
	t.Parallel()

	_, err := buildWitnessUpdateOperation(
		protocol.MustParseObjectID("1.6.1"),
		protocol.MustParseObjectID("1.2.345"),
		"",
		nil,
	)
	if err == nil {
		t.Fatal("buildWitnessUpdateOperation() error = nil, want error")
	}
}

func TestBuildWitnessUpdateOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	key := protocol.MustPublicKey(ecc.PrivateKeyFromSeed([]byte("witness-update")).PublicKey().String())
	op, err := buildWitnessUpdateOperation(
		protocol.MustParseObjectID("1.6.1"),
		protocol.MustParseObjectID("1.2.345"),
		"https://example.com/updated",
		&key,
	)
	if err != nil {
		t.Fatalf("buildWitnessUpdateOperation() error = %v", err)
	}
	if op.NewURL == nil || *op.NewURL != "https://example.com/updated" {
		t.Fatal("expected new url to be preserved")
	}
	if op.NewSigningKey == nil || op.NewSigningKey.String() != key.String() {
		t.Fatal("expected new signing key to be preserved")
	}
}

func TestBuildCommitteeMemberCreateOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	op, err := buildCommitteeMemberCreateOperation(protocol.MustParseObjectID("1.2.345"), "https://example.com/committee")
	if err != nil {
		t.Fatalf("buildCommitteeMemberCreateOperation() error = %v", err)
	}
	if got, want := op.CommitteeMemberAccount.String(), "1.2.345"; got != want {
		t.Fatalf("committee account = %q, want %q", got, want)
	}
}

func TestBuildCommitteeMemberCreateOperationRejectsLongURL(t *testing.T) {
	t.Parallel()

	_, err := buildCommitteeMemberCreateOperation(
		protocol.MustParseObjectID("1.2.345"),
		strings.Repeat("x", maxChainURLLength),
	)
	if err == nil {
		t.Fatal("buildCommitteeMemberCreateOperation() error = nil, want error")
	}
}

func TestBuildCommitteeMemberUpdateOperationRequiresChanges(t *testing.T) {
	t.Parallel()

	_, err := buildCommitteeMemberUpdateOperation(
		protocol.MustParseObjectID("1.5.1"),
		protocol.MustParseObjectID("1.2.345"),
		"",
	)
	if err == nil {
		t.Fatal("buildCommitteeMemberUpdateOperation() error = nil, want error")
	}
}

func TestBuildCommitteeMemberUpdateOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	op, err := buildCommitteeMemberUpdateOperation(
		protocol.MustParseObjectID("1.5.1"),
		protocol.MustParseObjectID("1.2.345"),
		"https://example.com/new-committee",
	)
	if err != nil {
		t.Fatalf("buildCommitteeMemberUpdateOperation() error = %v", err)
	}
	if op.NewURL == nil || *op.NewURL != "https://example.com/new-committee" {
		t.Fatal("expected new url to be preserved")
	}
}

func TestBuildAssetClaimFeesOperationPopulatesPayload(t *testing.T) {
	t.Parallel()

	claimFrom := protocol.MustParseObjectID("1.3.555")
	op, err := buildAssetClaimFeesOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.3.121"),
		5,
		1.25,
		&claimFrom,
	)
	if err != nil {
		t.Fatalf("buildAssetClaimFeesOperation() error = %v", err)
	}
	if got, want := op.AmountToClaim.Amount, int64(125000); got != want {
		t.Fatalf("claim amount = %d, want %d", got, want)
	}
	if op.Extensions == nil || op.Extensions.ClaimFromAssetID == nil || op.Extensions.ClaimFromAssetID.String() != "1.3.555" {
		t.Fatal("expected claim_from_asset_id to be preserved")
	}
}

func TestBuildAssetClaimFeesOperationRejectsZeroAmount(t *testing.T) {
	t.Parallel()

	_, err := buildAssetClaimFeesOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.3.121"),
		5,
		0,
		nil,
	)
	if err == nil {
		t.Fatal("buildAssetClaimFeesOperation() error = nil, want error")
	}
}

func TestBuildHTLCCreateOperationCarriesMemoExtension(t *testing.T) {
	t.Parallel()

	fromKey := ecc.PrivateKeyFromSeed([]byte("htlc-memo-from")).PublicKey().String()
	toKey := ecc.PrivateKeyFromSeed([]byte("htlc-memo-to")).PublicKey().String()
	memo, err := json.Marshal(protocol.MemoData{
		From:    fromKey,
		To:      toKey,
		Nonce:   "1",
		Message: "00",
	})
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	op, err := buildHTLCCreateOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.678"),
		protocol.MustParseObjectID("1.3.0"),
		5,
		1.25,
		protocol.HTLCPreimageHash{Kind: 2, Value: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"},
		32,
		3600,
		memo,
	)
	if err != nil {
		t.Fatalf("buildHTLCCreateOperation() error = %v", err)
	}
	if got, want := op.Amount.Amount, int64(125000); got != want {
		t.Fatalf("amount = %d, want %d", got, want)
	}
	if string(op.Extensions.Memo) == "" {
		t.Fatal("expected memo extension to be present")
	}
}

func TestBuildHTLCCreateOperationRejectsZeroClaimPeriod(t *testing.T) {
	t.Parallel()

	_, err := buildHTLCCreateOperation(
		protocol.MustParseObjectID("1.2.345"),
		protocol.MustParseObjectID("1.2.678"),
		protocol.MustParseObjectID("1.3.0"),
		5,
		1.25,
		protocol.HTLCPreimageHash{Kind: 2, Value: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"},
		32,
		0,
		nil,
	)
	if err == nil {
		t.Fatal("buildHTLCCreateOperation() error = nil, want error")
	}
}

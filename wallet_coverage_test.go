package bitshares

import (
	"reflect"
	"testing"
)

func TestWalletSignedOperationBuilderCoverage(t *testing.T) {
	t.Parallel()

	walletType := reflect.TypeOf(&Wallet{})
	required := []string{
		"BuildTransferOperation",
		"BuildAssetIssueOperation",
		"BuildAssetReserveOperation",
		"BuildAccountCreateOperation",
		"BuildAccountUpdateOperation",
		"BuildAccountWhitelistOperation",
		"BuildAccountUpgradeOperation",
		"BuildAccountTransferOperation",
		"BuildAssetCreateOperation",
		"BuildAssetUpdateOperation",
		"BuildAssetUpdateBitassetOperation",
		"BuildAssetUpdateFeedProducersOperation",
		"BuildAssetFundFeePoolOperation",
		"BuildAssetSettleOperation",
		"BuildAssetGlobalSettleOperation",
		"BuildAssetPublishFeedOperation",
		"BuildWitnessCreateOperation",
		"BuildWitnessUpdateOperation",
		"BuildProposalCreateOperation",
		"BuildProposalUpdateOperation",
		"BuildProposalDeleteOperation",
		"BuildWithdrawPermissionCreateOperation",
		"BuildWithdrawPermissionUpdateOperation",
		"BuildWithdrawPermissionClaimOperation",
		"BuildWithdrawPermissionDeleteOperation",
		"BuildCommitteeMemberCreateOperation",
		"BuildCommitteeMemberUpdateOperation",
		"BuildCommitteeMemberUpdateGlobalParametersOperation",
		"BuildBalanceClaimOperation",
		"BuildOverrideTransferOperation",
		"BuildAssetClaimFeesOperation",
		"BuildBidCollateralOperation",
		"BuildCallOrderUpdateOperation",
		"BuildLimitOrderCreateOperation",
		"BuildCancelOrderOperation",
		"BuildHTLCCreateOperation",
	}

	for _, method := range required {
		if _, ok := walletType.MethodByName(method); !ok {
			t.Fatalf("Wallet is missing required signed-operation builder %s", method)
		}
	}
}

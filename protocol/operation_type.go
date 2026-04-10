package protocol

import (
	"fmt"
	"strings"
)

// OperationType matches the BitShares operation variant tag order from core 7.0.2.
type OperationType uint16

const (
	OperationTypeTransfer OperationType = iota
	OperationTypeLimitOrderCreate
	OperationTypeLimitOrderCancel
	OperationTypeCallOrderUpdate
	OperationTypeFillOrder
	OperationTypeAccountCreate
	OperationTypeAccountUpdate
	OperationTypeAccountWhitelist
	OperationTypeAccountUpgrade
	OperationTypeAccountTransfer
	OperationTypeAssetCreate
	OperationTypeAssetUpdate
	OperationTypeAssetUpdateBitasset
	OperationTypeAssetUpdateFeedProducers
	OperationTypeAssetIssue
	OperationTypeAssetReserve
	OperationTypeAssetFundFeePool
	OperationTypeAssetSettle
	OperationTypeAssetGlobalSettle
	OperationTypeAssetPublishFeed
	OperationTypeWitnessCreate
	OperationTypeWitnessUpdate
	OperationTypeProposalCreate
	OperationTypeProposalUpdate
	OperationTypeProposalDelete
	OperationTypeWithdrawPermissionCreate
	OperationTypeWithdrawPermissionUpdate
	OperationTypeWithdrawPermissionClaim
	OperationTypeWithdrawPermissionDelete
	OperationTypeCommitteeMemberCreate
	OperationTypeCommitteeMemberUpdate
	OperationTypeCommitteeMemberUpdateGlobalParameters
	OperationTypeVestingBalanceCreate
	OperationTypeVestingBalanceWithdraw
	OperationTypeWorkerCreate
	OperationTypeCustom
	OperationTypeAssert
	OperationTypeBalanceClaim
	OperationTypeOverrideTransfer
	OperationTypeTransferToBlind
	OperationTypeBlindTransfer
	OperationTypeTransferFromBlind
	OperationTypeAssetSettleCancel
	OperationTypeAssetClaimFees
	OperationTypeFBADistribute
	OperationTypeBidCollateral
	OperationTypeExecuteBid
	OperationTypeAssetClaimPool
	OperationTypeAssetUpdateIssuer
	OperationTypeHTLCCreate
	OperationTypeHTLCRedeem
	OperationTypeHTLCRedeemed
	OperationTypeHTLCExtend
	OperationTypeHTLCRefund
	OperationTypeCustomAuthorityCreate
	OperationTypeCustomAuthorityUpdate
	OperationTypeCustomAuthorityDelete
	OperationTypeTicketCreate
	OperationTypeTicketUpdate
	OperationTypeLiquidityPoolCreate
	OperationTypeLiquidityPoolDelete
	OperationTypeLiquidityPoolDeposit
	OperationTypeLiquidityPoolWithdraw
	OperationTypeLiquidityPoolExchange
	OperationTypeSametFundCreate
	OperationTypeSametFundDelete
	OperationTypeSametFundUpdate
	OperationTypeSametFundBorrow
	OperationTypeSametFundRepay
	OperationTypeCreditOfferCreate
	OperationTypeCreditOfferDelete
	OperationTypeCreditOfferUpdate
	OperationTypeCreditOfferAccept
	OperationTypeCreditDealRepay
	OperationTypeCreditDealExpired
	OperationTypeLiquidityPoolUpdate
	OperationTypeCreditDealUpdate
	OperationTypeLimitOrderUpdate
)

var operationTypeNames = [...]string{
	"transfer",
	"limit_order_create",
	"limit_order_cancel",
	"call_order_update",
	"fill_order",
	"account_create",
	"account_update",
	"account_whitelist",
	"account_upgrade",
	"account_transfer",
	"asset_create",
	"asset_update",
	"asset_update_bitasset",
	"asset_update_feed_producers",
	"asset_issue",
	"asset_reserve",
	"asset_fund_fee_pool",
	"asset_settle",
	"asset_global_settle",
	"asset_publish_feed",
	"witness_create",
	"witness_update",
	"proposal_create",
	"proposal_update",
	"proposal_delete",
	"withdraw_permission_create",
	"withdraw_permission_update",
	"withdraw_permission_claim",
	"withdraw_permission_delete",
	"committee_member_create",
	"committee_member_update",
	"committee_member_update_global_parameters",
	"vesting_balance_create",
	"vesting_balance_withdraw",
	"worker_create",
	"custom",
	"assert",
	"balance_claim",
	"override_transfer",
	"transfer_to_blind",
	"blind_transfer",
	"transfer_from_blind",
	"asset_settle_cancel",
	"asset_claim_fees",
	"fba_distribute",
	"bid_collateral",
	"execute_bid",
	"asset_claim_pool",
	"asset_update_issuer",
	"htlc_create",
	"htlc_redeem",
	"htlc_redeemed",
	"htlc_extend",
	"htlc_refund",
	"custom_authority_create",
	"custom_authority_update",
	"custom_authority_delete",
	"ticket_create",
	"ticket_update",
	"liquidity_pool_create",
	"liquidity_pool_delete",
	"liquidity_pool_deposit",
	"liquidity_pool_withdraw",
	"liquidity_pool_exchange",
	"samet_fund_create",
	"samet_fund_delete",
	"samet_fund_update",
	"samet_fund_borrow",
	"samet_fund_repay",
	"credit_offer_create",
	"credit_offer_delete",
	"credit_offer_update",
	"credit_offer_accept",
	"credit_deal_repay",
	"credit_deal_expired",
	"liquidity_pool_update",
	"credit_deal_update",
	"limit_order_update",
}

func (t OperationType) String() string {
	if int(t) < len(operationTypeNames) {
		return operationTypeNames[t]
	}
	return fmt.Sprintf("operation_type_%d", t)
}

func (t OperationType) OperationName() string {
	return t.String()
}

var operationTypeByName = func() map[string]OperationType {
	out := make(map[string]OperationType, len(operationTypeNames))
	for i, name := range operationTypeNames {
		out[name] = OperationType(i)
	}
	return out
}()

// ParseOperationType resolves a BitShares operation name to its numeric tag.
func ParseOperationType(name string) (OperationType, error) {
	kind, ok := operationTypeByName[strings.ToLower(name)]
	if !ok {
		return 0, fmt.Errorf("unknown operation type %q", name)
	}
	return kind, nil
}

// MustParseOperationType panics if the name is not a recognized operation type.
func MustParseOperationType(name string) OperationType {
	kind, err := ParseOperationType(name)
	if err != nil {
		panic(err)
	}
	return kind
}

package protocol

import "sync"

var KnownOperationTypes = map[OperationType]struct{}{
	OperationTypeTransfer:                              {},
	OperationTypeLimitOrderCreate:                      {},
	OperationTypeLimitOrderCancel:                      {},
	OperationTypeCallOrderUpdate:                       {},
	OperationTypeFillOrder:                             {},
	OperationTypeAccountCreate:                         {},
	OperationTypeAccountUpdate:                         {},
	OperationTypeAccountWhitelist:                      {},
	OperationTypeAccountUpgrade:                        {},
	OperationTypeAccountTransfer:                       {},
	OperationTypeAssetCreate:                           {},
	OperationTypeAssetUpdate:                           {},
	OperationTypeAssetUpdateBitasset:                   {},
	OperationTypeAssetUpdateFeedProducers:              {},
	OperationTypeAssetIssue:                            {},
	OperationTypeAssetReserve:                          {},
	OperationTypeAssetFundFeePool:                      {},
	OperationTypeAssetSettle:                           {},
	OperationTypeAssetGlobalSettle:                     {},
	OperationTypeAssetPublishFeed:                      {},
	OperationTypeWitnessCreate:                         {},
	OperationTypeWitnessUpdate:                         {},
	OperationTypeProposalCreate:                        {},
	OperationTypeProposalUpdate:                        {},
	OperationTypeProposalDelete:                        {},
	OperationTypeWithdrawPermissionCreate:              {},
	OperationTypeWithdrawPermissionUpdate:              {},
	OperationTypeWithdrawPermissionClaim:               {},
	OperationTypeWithdrawPermissionDelete:              {},
	OperationTypeCommitteeMemberCreate:                 {},
	OperationTypeCommitteeMemberUpdate:                 {},
	OperationTypeCommitteeMemberUpdateGlobalParameters: {},
	OperationTypeVestingBalanceCreate:                  {},
	OperationTypeVestingBalanceWithdraw:                {},
	OperationTypeWorkerCreate:                          {},
	OperationTypeCustom:                                {},
	OperationTypeAssert:                                {},
	OperationTypeBalanceClaim:                          {},
	OperationTypeOverrideTransfer:                      {},
	OperationTypeTransferToBlind:                       {},
	OperationTypeBlindTransfer:                         {},
	OperationTypeTransferFromBlind:                     {},
	OperationTypeAssetSettleCancel:                     {},
	OperationTypeAssetClaimFees:                        {},
	OperationTypeFBADistribute:                         {},
	OperationTypeBidCollateral:                         {},
	OperationTypeExecuteBid:                            {},
	OperationTypeAssetClaimPool:                        {},
	OperationTypeAssetUpdateIssuer:                     {},
	OperationTypeHTLCCreate:                            {},
	OperationTypeHTLCRedeem:                            {},
	OperationTypeHTLCRedeemed:                          {},
	OperationTypeHTLCExtend:                            {},
	OperationTypeHTLCRefund:                            {},
	OperationTypeCustomAuthorityCreate:                 {},
	OperationTypeCustomAuthorityUpdate:                 {},
	OperationTypeCustomAuthorityDelete:                 {},
	OperationTypeTicketCreate:                          {},
	OperationTypeTicketUpdate:                          {},
	OperationTypeLiquidityPoolCreate:                   {},
	OperationTypeLiquidityPoolDelete:                   {},
	OperationTypeLiquidityPoolDeposit:                  {},
	OperationTypeLiquidityPoolWithdraw:                 {},
	OperationTypeLiquidityPoolExchange:                 {},
	OperationTypeSametFundCreate:                       {},
	OperationTypeSametFundDelete:                       {},
	OperationTypeSametFundUpdate:                       {},
	OperationTypeSametFundBorrow:                       {},
	OperationTypeSametFundRepay:                        {},
	OperationTypeCreditOfferCreate:                     {},
	OperationTypeCreditOfferDelete:                     {},
	OperationTypeCreditOfferUpdate:                     {},
	OperationTypeCreditOfferAccept:                     {},
	OperationTypeCreditDealRepay:                       {},
	OperationTypeCreditDealExpired:                     {},
	OperationTypeLiquidityPoolUpdate:                   {},
	OperationTypeCreditDealUpdate:                      {},
	OperationTypeLimitOrderUpdate:                      {},
}

func IsKnownOperationType(kind OperationType) bool {
	_, ok := KnownOperationTypes[kind]
	return ok
}

var (
	factoryMu          sync.RWMutex
	operationFactories = map[OperationType]func() Operation{}
)

// RegisterOperationFactory registers a typed decoder for a known operation type.
func RegisterOperationFactory(kind OperationType, factory func() Operation) {
	factoryMu.Lock()
	defer factoryMu.Unlock()
	operationFactories[kind] = factory
}

func newOperation(kind OperationType) Operation {
	factoryMu.RLock()
	factory := operationFactories[kind]
	factoryMu.RUnlock()
	if factory == nil {
		return nil
	}
	return factory()
}

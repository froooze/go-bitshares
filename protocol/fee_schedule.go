package protocol

import (
	"encoding/json"
	"fmt"
	"sort"
)

// FeeParameterValue is a BitShares fee-parameter payload.
type FeeParameterValue interface {
	MarshalBinaryInto(*binaryWriter) error
	UnmarshalBinaryFrom(*binaryReader) error
	feeValue() (float64, bool)
}

// EmptyFeeParameters represents an empty fee parameter payload.
type EmptyFeeParameters struct{}

func (e *EmptyFeeParameters) MarshalBinaryInto(*binaryWriter) error   { return nil }
func (e *EmptyFeeParameters) UnmarshalBinaryFrom(*binaryReader) error { return nil }
func (e *EmptyFeeParameters) feeValue() (float64, bool)               { return 0, false }

// FeeOnlyParameters represents { fee }.
type FeeOnlyParameters struct {
	Fee uint64 `json:"fee"`
}

func (p *FeeOnlyParameters) MarshalBinaryInto(w *binaryWriter) error {
	w.writeUint64(p.Fee)
	return nil
}

func (p *FeeOnlyParameters) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := r.readUint64()
	if err != nil {
		return err
	}
	p.Fee = fee
	return nil
}

func (p *FeeOnlyParameters) feeValue() (float64, bool) { return float64(p.Fee), true }

// SignedFeeParameters represents { fee } with signed values.
type SignedFeeParameters struct {
	Fee int64 `json:"fee"`
}

func (p *SignedFeeParameters) MarshalBinaryInto(w *binaryWriter) error {
	w.writeInt64(p.Fee)
	return nil
}

func (p *SignedFeeParameters) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := r.readInt64()
	if err != nil {
		return err
	}
	p.Fee = fee
	return nil
}

func (p *SignedFeeParameters) feeValue() (float64, bool) { return float64(p.Fee), true }

// FeeAndPricePerKbyteParameters represents { fee, price_per_kbyte }.
type FeeAndPricePerKbyteParameters struct {
	Fee           uint64 `json:"fee"`
	PricePerKbyte uint32 `json:"price_per_kbyte"`
}

func (p *FeeAndPricePerKbyteParameters) MarshalBinaryInto(w *binaryWriter) error {
	w.writeUint64(p.Fee)
	w.writeUint32(p.PricePerKbyte)
	return nil
}

func (p *FeeAndPricePerKbyteParameters) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := r.readUint64()
	if err != nil {
		return err
	}
	price, err := r.readUint32()
	if err != nil {
		return err
	}
	p.Fee = fee
	p.PricePerKbyte = price
	return nil
}

func (p *FeeAndPricePerKbyteParameters) feeValue() (float64, bool) { return float64(p.Fee), true }

// SignedFeeAndPricePerKbyteParameters represents { fee, price_per_kbyte } with signed fee.
type SignedFeeAndPricePerKbyteParameters struct {
	Fee           int64  `json:"fee"`
	PricePerKbyte uint32 `json:"price_per_kbyte"`
}

func (p *SignedFeeAndPricePerKbyteParameters) MarshalBinaryInto(w *binaryWriter) error {
	w.writeInt64(p.Fee)
	w.writeUint32(p.PricePerKbyte)
	return nil
}

func (p *SignedFeeAndPricePerKbyteParameters) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := r.readInt64()
	if err != nil {
		return err
	}
	price, err := r.readUint32()
	if err != nil {
		return err
	}
	p.Fee = fee
	p.PricePerKbyte = price
	return nil
}

func (p *SignedFeeAndPricePerKbyteParameters) feeValue() (float64, bool) { return float64(p.Fee), true }

// BasicPremiumPricePerKbyteParameters represents { basic_fee, premium_fee, price_per_kbyte }.
type BasicPremiumPricePerKbyteParameters struct {
	BasicFee      uint64 `json:"basic_fee"`
	PremiumFee    uint64 `json:"premium_fee"`
	PricePerKbyte uint32 `json:"price_per_kbyte"`
}

func (p *BasicPremiumPricePerKbyteParameters) MarshalBinaryInto(w *binaryWriter) error {
	w.writeUint64(p.BasicFee)
	w.writeUint64(p.PremiumFee)
	w.writeUint32(p.PricePerKbyte)
	return nil
}

func (p *BasicPremiumPricePerKbyteParameters) UnmarshalBinaryFrom(r *binaryReader) error {
	basic, err := r.readUint64()
	if err != nil {
		return err
	}
	premium, err := r.readUint64()
	if err != nil {
		return err
	}
	price, err := r.readUint32()
	if err != nil {
		return err
	}
	p.BasicFee = basic
	p.PremiumFee = premium
	p.PricePerKbyte = price
	return nil
}

func (p *BasicPremiumPricePerKbyteParameters) feeValue() (float64, bool) {
	return float64(p.BasicFee), true
}

// MembershipFeeParameters represents { membership_annual_fee, membership_lifetime_fee }.
type MembershipFeeParameters struct {
	MembershipAnnualFee   uint64 `json:"membership_annual_fee"`
	MembershipLifetimeFee uint64 `json:"membership_lifetime_fee"`
}

func (p *MembershipFeeParameters) MarshalBinaryInto(w *binaryWriter) error {
	w.writeUint64(p.MembershipAnnualFee)
	w.writeUint64(p.MembershipLifetimeFee)
	return nil
}

func (p *MembershipFeeParameters) UnmarshalBinaryFrom(r *binaryReader) error {
	annual, err := r.readUint64()
	if err != nil {
		return err
	}
	lifetime, err := r.readUint64()
	if err != nil {
		return err
	}
	p.MembershipAnnualFee = annual
	p.MembershipLifetimeFee = lifetime
	return nil
}

func (p *MembershipFeeParameters) feeValue() (float64, bool) {
	return float64(p.MembershipAnnualFee), true
}

// SymbolFeeParameters represents { symbol3, symbol4, long_symbol, price_per_kbyte }.
type SymbolFeeParameters struct {
	Symbol3       uint64 `json:"symbol3"`
	Symbol4       uint64 `json:"symbol4"`
	LongSymbol    uint64 `json:"long_symbol"`
	PricePerKbyte uint32 `json:"price_per_kbyte"`
}

func (p *SymbolFeeParameters) MarshalBinaryInto(w *binaryWriter) error {
	w.writeUint64(p.Symbol3)
	w.writeUint64(p.Symbol4)
	w.writeUint64(p.LongSymbol)
	w.writeUint32(p.PricePerKbyte)
	return nil
}

func (p *SymbolFeeParameters) UnmarshalBinaryFrom(r *binaryReader) error {
	s3, err := r.readUint64()
	if err != nil {
		return err
	}
	s4, err := r.readUint64()
	if err != nil {
		return err
	}
	long, err := r.readUint64()
	if err != nil {
		return err
	}
	price, err := r.readUint32()
	if err != nil {
		return err
	}
	p.Symbol3 = s3
	p.Symbol4 = s4
	p.LongSymbol = long
	p.PricePerKbyte = price
	return nil
}

func (p *SymbolFeeParameters) feeValue() (float64, bool) { return float64(p.Symbol3), true }

// BasicFeePerByteParameters represents { basic_fee, price_per_byte }.
type BasicFeePerByteParameters struct {
	BasicFee     uint64 `json:"basic_fee"`
	PricePerByte uint32 `json:"price_per_byte"`
}

func (p *BasicFeePerByteParameters) MarshalBinaryInto(w *binaryWriter) error {
	w.writeUint64(p.BasicFee)
	w.writeUint32(p.PricePerByte)
	return nil
}

func (p *BasicFeePerByteParameters) UnmarshalBinaryFrom(r *binaryReader) error {
	basic, err := r.readUint64()
	if err != nil {
		return err
	}
	price, err := r.readUint32()
	if err != nil {
		return err
	}
	p.BasicFee = basic
	p.PricePerByte = price
	return nil
}

func (p *BasicFeePerByteParameters) feeValue() (float64, bool) { return float64(p.BasicFee), true }

// FeeAndPricePerDayParameters represents { fee, fee_per_day }.
type FeeAndPricePerDayParameters struct {
	Fee       uint64 `json:"fee"`
	FeePerDay uint64 `json:"fee_per_day"`
}

func (p *FeeAndPricePerDayParameters) MarshalBinaryInto(w *binaryWriter) error {
	w.writeUint64(p.Fee)
	w.writeUint64(p.FeePerDay)
	return nil
}

func (p *FeeAndPricePerDayParameters) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := r.readUint64()
	if err != nil {
		return err
	}
	perDay, err := r.readUint64()
	if err != nil {
		return err
	}
	p.Fee = fee
	p.FeePerDay = perDay
	return nil
}

func (p *FeeAndPricePerDayParameters) feeValue() (float64, bool) { return float64(p.Fee), true }

// FeeAndPricePerKbParameters represents { fee, fee_per_kb }.
type FeeAndPricePerKbParameters struct {
	Fee      uint64 `json:"fee"`
	FeePerKb uint64 `json:"fee_per_kb"`
}

func (p *FeeAndPricePerKbParameters) MarshalBinaryInto(w *binaryWriter) error {
	w.writeUint64(p.Fee)
	w.writeUint64(p.FeePerKb)
	return nil
}

func (p *FeeAndPricePerKbParameters) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := r.readUint64()
	if err != nil {
		return err
	}
	perKb, err := r.readUint64()
	if err != nil {
		return err
	}
	p.Fee = fee
	p.FeePerKb = perKb
	return nil
}

func (p *FeeAndPricePerKbParameters) feeValue() (float64, bool) { return float64(p.Fee), true }

// FeeAndPricePerOutputParameters represents { fee, price_per_output }.
type FeeAndPricePerOutputParameters struct {
	Fee            uint64 `json:"fee"`
	PricePerOutput uint32 `json:"price_per_output"`
}

func (p *FeeAndPricePerOutputParameters) MarshalBinaryInto(w *binaryWriter) error {
	w.writeUint64(p.Fee)
	w.writeUint32(p.PricePerOutput)
	return nil
}

func (p *FeeAndPricePerOutputParameters) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := r.readUint64()
	if err != nil {
		return err
	}
	price, err := r.readUint32()
	if err != nil {
		return err
	}
	p.Fee = fee
	p.PricePerOutput = price
	return nil
}

func (p *FeeAndPricePerOutputParameters) feeValue() (float64, bool) { return float64(p.Fee), true }

// FeeScheduleParameter represents one fee schedule entry.
type FeeScheduleParameter struct {
	OperationType OperationType     `json:"-"`
	Value         FeeParameterValue `json:"-"`
}

func newFeeParameterValue(kind OperationType) (FeeParameterValue, error) {
	switch kind {
	case OperationTypeTransfer, OperationTypeProposalCreate, OperationTypeProposalUpdate, OperationTypeAssetIssue, OperationTypeCustom, OperationTypeOverrideTransfer, OperationTypeCreditOfferCreate, OperationTypeCreditOfferUpdate:
		return &FeeAndPricePerKbyteParameters{}, nil
	case OperationTypeLimitOrderCreate, OperationTypeLimitOrderCancel, OperationTypeCallOrderUpdate, OperationTypeAccountTransfer, OperationTypeAssetReserve, OperationTypeAssetFundFeePool, OperationTypeAssetSettle, OperationTypeAssetGlobalSettle, OperationTypeAssetPublishFeed, OperationTypeWitnessCreate, OperationTypeProposalDelete, OperationTypeWithdrawPermissionCreate, OperationTypeWithdrawPermissionUpdate, OperationTypeWithdrawPermissionDelete, OperationTypeCommitteeMemberCreate, OperationTypeCommitteeMemberUpdate, OperationTypeCommitteeMemberUpdateGlobalParameters, OperationTypeVestingBalanceCreate, OperationTypeVestingBalanceWithdraw, OperationTypeWorkerCreate, OperationTypeAssert, OperationTypeTransferFromBlind, OperationTypeAssetClaimFees, OperationTypeBidCollateral, OperationTypeAssetClaimPool, OperationTypeAssetUpdateIssuer, OperationTypeCustomAuthorityDelete, OperationTypeTicketCreate, OperationTypeTicketUpdate, OperationTypeLiquidityPoolCreate, OperationTypeLiquidityPoolDelete, OperationTypeLiquidityPoolDeposit, OperationTypeLiquidityPoolWithdraw, OperationTypeLiquidityPoolExchange, OperationTypeSametFundCreate, OperationTypeSametFundDelete, OperationTypeSametFundUpdate, OperationTypeSametFundBorrow, OperationTypeSametFundRepay, OperationTypeCreditOfferDelete, OperationTypeCreditOfferAccept, OperationTypeCreditDealRepay, OperationTypeCreditDealUpdate, OperationTypeLiquidityPoolUpdate, OperationTypeLimitOrderUpdate:
		return &FeeOnlyParameters{}, nil
	case OperationTypeFillOrder, OperationTypeBalanceClaim, OperationTypeAssetSettleCancel, OperationTypeFBADistribute, OperationTypeExecuteBid, OperationTypeCreditDealExpired, OperationTypeHTLCRedeemed, OperationTypeHTLCRefund:
		return &EmptyFeeParameters{}, nil
	case OperationTypeAccountCreate:
		return &BasicPremiumPricePerKbyteParameters{}, nil
	case OperationTypeAccountUpdate:
		return &SignedFeeAndPricePerKbyteParameters{}, nil
	case OperationTypeAccountWhitelist, OperationTypeWitnessUpdate:
		return &SignedFeeParameters{}, nil
	case OperationTypeAccountUpgrade:
		return &MembershipFeeParameters{}, nil
	case OperationTypeAssetCreate:
		return &SymbolFeeParameters{}, nil
	case OperationTypeAssetUpdate:
		return &FeeAndPricePerKbyteParameters{}, nil
	case OperationTypeAssetUpdateBitasset, OperationTypeAssetUpdateFeedProducers:
		return &FeeOnlyParameters{}, nil
	case OperationTypeWithdrawPermissionClaim:
		return &FeeAndPricePerKbyteParameters{}, nil
	case OperationTypeBlindTransfer, OperationTypeTransferToBlind:
		return &FeeAndPricePerOutputParameters{}, nil
	case OperationTypeHTLCCreate, OperationTypeHTLCExtend:
		return &FeeAndPricePerDayParameters{}, nil
	case OperationTypeHTLCRedeem:
		return &FeeAndPricePerKbParameters{}, nil
	case OperationTypeCustomAuthorityCreate, OperationTypeCustomAuthorityUpdate:
		return &BasicFeePerByteParameters{}, nil
	default:
		return nil, fmt.Errorf("unsupported fee parameter type %d", kind)
	}
}

func (p FeeScheduleParameter) MarshalJSON() ([]byte, error) {
	if p.Value == nil {
		value, err := newFeeParameterValue(p.OperationType)
		if err != nil {
			return nil, err
		}
		p.Value = value
	}
	payload, err := json.Marshal(p.Value)
	if err != nil {
		return nil, err
	}
	return json.Marshal([]any{uint16(p.OperationType), json.RawMessage(payload)})
}

func (p *FeeScheduleParameter) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) != 2 {
		return fmt.Errorf("invalid fee schedule parameter")
	}
	var kind uint16
	if err := json.Unmarshal(raw[0], &kind); err != nil {
		return err
	}
	value, err := newFeeParameterValue(OperationType(kind))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw[1], value); err != nil {
		return err
	}
	p.OperationType = OperationType(kind)
	p.Value = value
	return nil
}

func (p FeeScheduleParameter) MarshalBinaryInto(w *binaryWriter) error {
	if p.Value == nil {
		value, err := newFeeParameterValue(p.OperationType)
		if err != nil {
			return err
		}
		p.Value = value
	}
	w.writeVarUint64(uint64(p.OperationType))
	return p.Value.MarshalBinaryInto(w)
}

func (p *FeeScheduleParameter) UnmarshalBinaryFrom(r *binaryReader) error {
	kind, err := r.readVarUint64()
	if err != nil {
		return err
	}
	value, err := newFeeParameterValue(OperationType(kind))
	if err != nil {
		return err
	}
	if err := value.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	p.OperationType = OperationType(kind)
	p.Value = value
	return nil
}

func (p FeeScheduleParameter) FeeValue() (float64, bool) {
	if p.Value == nil {
		return 0, false
	}
	return p.Value.feeValue()
}

// FeeSchedule mirrors the chain fee schedule JSON shape.
type FeeSchedule struct {
	Scale      uint32                 `json:"scale"`
	Parameters []FeeScheduleParameter `json:"parameters"`
}

// Operation-specific fee parameter aliases keep the public API readable.
type TransferOperationFeeParameters = FeeAndPricePerKbyteParameters
type LimitOrderCreateOperationFeeParameters = FeeOnlyParameters
type LimitOrderCancelOperationFeeParameters = FeeOnlyParameters
type CallOrderUpdateOperationFeeParameters = FeeOnlyParameters
type FillOrderOperationFeeParameters = EmptyFeeParameters
type AccountCreateOperationFeeParameters = BasicPremiumPricePerKbyteParameters
type AccountUpdateOperationFeeParameters = SignedFeeAndPricePerKbyteParameters
type AccountWhitelistOperationFeeParameters = SignedFeeParameters
type AccountUpgradeOperationFeeParameters = MembershipFeeParameters
type AccountTransferOperationFeeParameters = FeeOnlyParameters
type AssetCreateOperationFeeParameters = SymbolFeeParameters
type AssetUpdateOperationFeeParameters = FeeAndPricePerKbyteParameters
type AssetUpdateBitassetOperationFeeParameters = FeeOnlyParameters
type AssetUpdateFeedProducersOperationFeeParameters = FeeOnlyParameters
type AssetIssueOperationFeeParameters = FeeAndPricePerKbyteParameters
type AssetReserveOperationFeeParameters = FeeOnlyParameters
type AssetFundFeePoolOperationFeeParameters = FeeOnlyParameters
type AssetSettleOperationFeeParameters = FeeOnlyParameters
type AssetGlobalSettleOperationFeeParameters = FeeOnlyParameters
type AssetPublishFeedOperationFeeParameters = FeeOnlyParameters
type WitnessCreateOperationFeeParameters = FeeOnlyParameters
type WitnessUpdateOperationFeeParameters = SignedFeeParameters
type ProposalCreateOperationFeeParameters = FeeAndPricePerKbyteParameters
type ProposalUpdateOperationFeeParameters = FeeAndPricePerKbyteParameters
type ProposalDeleteOperationFeeParameters = FeeOnlyParameters
type WithdrawPermissionCreateOperationFeeParameters = FeeOnlyParameters
type WithdrawPermissionUpdateOperationFeeParameters = FeeOnlyParameters
type WithdrawPermissionClaimOperationFeeParameters = FeeAndPricePerKbyteParameters
type WithdrawPermissionDeleteOperationFeeParameters = FeeOnlyParameters
type CommitteeMemberCreateOperationFeeParameters = FeeOnlyParameters
type CommitteeMemberUpdateOperationFeeParameters = FeeOnlyParameters
type CommitteeMemberUpdateGlobalParametersOperationFeeParameters = FeeOnlyParameters
type VestingBalanceCreateOperationFeeParameters = FeeOnlyParameters
type VestingBalanceWithdrawOperationFeeParameters = FeeOnlyParameters
type WorkerCreateOperationFeeParameters = FeeOnlyParameters
type CustomOperationFeeParameters = FeeAndPricePerKbyteParameters
type AssertOperationFeeParameters = FeeOnlyParameters
type BalanceClaimOperationFeeParameters = EmptyFeeParameters
type OverrideTransferOperationFeeParameters = FeeAndPricePerKbyteParameters
type TransferToBlindOperationFeeParameters = FeeAndPricePerOutputParameters
type BlindTransferOperationFeeParameters = FeeAndPricePerOutputParameters
type TransferFromBlindOperationFeeParameters = FeeOnlyParameters
type AssetSettleCancelOperationFeeParameters = EmptyFeeParameters
type AssetClaimFeesOperationFeeParameters = FeeOnlyParameters
type FBADistributeOperationFeeParameters = EmptyFeeParameters
type BidCollateralOperationFeeParameters = FeeOnlyParameters
type ExecuteBidOperationFeeParameters = EmptyFeeParameters
type AssetClaimPoolOperationFeeParameters = FeeOnlyParameters
type AssetUpdateIssuerOperationFeeParameters = FeeOnlyParameters
type HTLCCreateOperationFeeParameters = FeeAndPricePerDayParameters
type HTLCRedeemOperationFeeParameters = FeeAndPricePerKbParameters
type HTLCRedeemedOperationFeeParameters = EmptyFeeParameters
type HTLCExtendOperationFeeParameters = FeeAndPricePerDayParameters
type HTLCRefundOperationFeeParameters = EmptyFeeParameters
type CustomAuthorityCreateOperationFeeParameters = BasicFeePerByteParameters
type CustomAuthorityUpdateOperationFeeParameters = BasicFeePerByteParameters
type CustomAuthorityDeleteOperationFeeParameters = FeeOnlyParameters
type TicketCreateOperationFeeParameters = FeeOnlyParameters
type TicketUpdateOperationFeeParameters = FeeOnlyParameters
type LiquidityPoolCreateOperationFeeParameters = FeeOnlyParameters
type LiquidityPoolDeleteOperationFeeParameters = FeeOnlyParameters
type LiquidityPoolDepositOperationFeeParameters = FeeOnlyParameters
type LiquidityPoolWithdrawOperationFeeParameters = FeeOnlyParameters
type LiquidityPoolExchangeOperationFeeParameters = FeeOnlyParameters
type SametFundCreateOperationFeeParameters = FeeOnlyParameters
type SametFundDeleteOperationFeeParameters = FeeOnlyParameters
type SametFundUpdateOperationFeeParameters = FeeOnlyParameters
type SametFundBorrowOperationFeeParameters = FeeOnlyParameters
type SametFundRepayOperationFeeParameters = FeeOnlyParameters
type CreditOfferCreateOperationFeeParameters = FeeAndPricePerKbyteParameters
type CreditOfferDeleteOperationFeeParameters = FeeOnlyParameters
type CreditOfferUpdateOperationFeeParameters = FeeAndPricePerKbyteParameters
type CreditOfferAcceptOperationFeeParameters = FeeOnlyParameters
type CreditDealRepayOperationFeeParameters = FeeOnlyParameters
type CreditDealExpiredOperationFeeParameters = EmptyFeeParameters
type LiquidityPoolUpdateOperationFeeParameters = FeeOnlyParameters
type CreditDealUpdateOperationFeeParameters = FeeOnlyParameters
type LimitOrderUpdateOperationFeeParameters = FeeOnlyParameters

func (f FeeSchedule) MarshalBinaryInto(w *binaryWriter) error {
	params := append([]FeeScheduleParameter(nil), f.Parameters...)
	sort.Slice(params, func(i, j int) bool {
		return params[i].OperationType < params[j].OperationType
	})
	w.writeVarUint64(uint64(len(params)))
	for i := range params {
		if err := params[i].MarshalBinaryInto(w); err != nil {
			return err
		}
	}
	w.writeUint32(f.Scale)
	return nil
}

func (f *FeeSchedule) UnmarshalBinaryFrom(r *binaryReader) error {
	count, err := r.readVarUint64()
	if err != nil {
		return err
	}
	params := make([]FeeScheduleParameter, 0, count)
	for i := uint64(0); i < count; i++ {
		var param FeeScheduleParameter
		if err := param.UnmarshalBinaryFrom(r); err != nil {
			return err
		}
		params = append(params, param)
	}
	scale, err := r.readUint32()
	if err != nil {
		return err
	}
	f.Parameters = params
	f.Scale = scale
	return nil
}

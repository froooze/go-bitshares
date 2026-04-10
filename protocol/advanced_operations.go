package protocol

import "encoding/json"

// LiquidityPoolCreateOperation creates a liquidity pool.
type LiquidityPoolCreateOperation struct {
	Fee                  AssetAmount       `json:"fee"`
	Account              ObjectID          `json:"account"`
	AssetA               ObjectID          `json:"asset_a"`
	AssetB               ObjectID          `json:"asset_b"`
	ShareAsset           ObjectID          `json:"share_asset"`
	TakerFeePercent      uint16            `json:"taker_fee_percent"`
	WithdrawalFeePercent uint16            `json:"withdrawal_fee_percent"`
	Extensions           []json.RawMessage `json:"extensions"`
}

func (o LiquidityPoolCreateOperation) Type() OperationType { return OperationTypeLiquidityPoolCreate }

func (o LiquidityPoolCreateOperation) MarshalJSON() ([]byte, error) {
	type alias LiquidityPoolCreateOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *LiquidityPoolCreateOperation) UnmarshalJSON(data []byte) error {
	type alias LiquidityPoolCreateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeLiquidityPoolCreate, &payload); err != nil {
		return err
	}
	*o = LiquidityPoolCreateOperation(payload)
	return nil
}

// LiquidityPoolDeleteOperation deletes a pool.
type LiquidityPoolDeleteOperation struct {
	Fee        AssetAmount       `json:"fee"`
	Account    ObjectID          `json:"account"`
	Pool       ObjectID          `json:"pool"`
	Extensions []json.RawMessage `json:"extensions"`
}

func (o LiquidityPoolDeleteOperation) Type() OperationType { return OperationTypeLiquidityPoolDelete }

func (o LiquidityPoolDeleteOperation) MarshalJSON() ([]byte, error) {
	type alias LiquidityPoolDeleteOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *LiquidityPoolDeleteOperation) UnmarshalJSON(data []byte) error {
	type alias LiquidityPoolDeleteOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeLiquidityPoolDelete, &payload); err != nil {
		return err
	}
	*o = LiquidityPoolDeleteOperation(payload)
	return nil
}

type LiquidityPoolDepositOperation struct {
	Fee        AssetAmount       `json:"fee"`
	Account    ObjectID          `json:"account"`
	Pool       ObjectID          `json:"pool"`
	AmountA    AssetAmount       `json:"amount_a"`
	AmountB    AssetAmount       `json:"amount_b"`
	Extensions []json.RawMessage `json:"extensions"`
}

func (o LiquidityPoolDepositOperation) Type() OperationType { return OperationTypeLiquidityPoolDeposit }

func (o LiquidityPoolDepositOperation) MarshalJSON() ([]byte, error) {
	type alias LiquidityPoolDepositOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *LiquidityPoolDepositOperation) UnmarshalJSON(data []byte) error {
	type alias LiquidityPoolDepositOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeLiquidityPoolDeposit, &payload); err != nil {
		return err
	}
	*o = LiquidityPoolDepositOperation(payload)
	return nil
}

type LiquidityPoolWithdrawOperation struct {
	Fee         AssetAmount       `json:"fee"`
	Account     ObjectID          `json:"account"`
	Pool        ObjectID          `json:"pool"`
	ShareAmount AssetAmount       `json:"share_amount"`
	Extensions  []json.RawMessage `json:"extensions"`
}

func (o LiquidityPoolWithdrawOperation) Type() OperationType {
	return OperationTypeLiquidityPoolWithdraw
}

func (o LiquidityPoolWithdrawOperation) MarshalJSON() ([]byte, error) {
	type alias LiquidityPoolWithdrawOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *LiquidityPoolWithdrawOperation) UnmarshalJSON(data []byte) error {
	type alias LiquidityPoolWithdrawOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeLiquidityPoolWithdraw, &payload); err != nil {
		return err
	}
	*o = LiquidityPoolWithdrawOperation(payload)
	return nil
}

type LiquidityPoolExchangeOperation struct {
	Fee          AssetAmount       `json:"fee"`
	Account      ObjectID          `json:"account"`
	Pool         ObjectID          `json:"pool"`
	AmountToSell AssetAmount       `json:"amount_to_sell"`
	MinToReceive AssetAmount       `json:"min_to_receive"`
	Extensions   []json.RawMessage `json:"extensions"`
}

func (o LiquidityPoolExchangeOperation) Type() OperationType {
	return OperationTypeLiquidityPoolExchange
}

func (o LiquidityPoolExchangeOperation) MarshalJSON() ([]byte, error) {
	type alias LiquidityPoolExchangeOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *LiquidityPoolExchangeOperation) UnmarshalJSON(data []byte) error {
	type alias LiquidityPoolExchangeOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeLiquidityPoolExchange, &payload); err != nil {
		return err
	}
	*o = LiquidityPoolExchangeOperation(payload)
	return nil
}

type LiquidityPoolUpdateOperation struct {
	Fee                  AssetAmount       `json:"fee"`
	Account              ObjectID          `json:"account"`
	Pool                 ObjectID          `json:"pool"`
	TakerFeePercent      *uint16           `json:"taker_fee_percent,omitempty"`
	WithdrawalFeePercent *uint16           `json:"withdrawal_fee_percent,omitempty"`
	Extensions           []json.RawMessage `json:"extensions"`
}

func (o LiquidityPoolUpdateOperation) Type() OperationType { return OperationTypeLiquidityPoolUpdate }

func (o LiquidityPoolUpdateOperation) MarshalJSON() ([]byte, error) {
	type alias LiquidityPoolUpdateOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *LiquidityPoolUpdateOperation) UnmarshalJSON(data []byte) error {
	type alias LiquidityPoolUpdateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeLiquidityPoolUpdate, &payload); err != nil {
		return err
	}
	*o = LiquidityPoolUpdateOperation(payload)
	return nil
}

type SametFundCreateOperation struct {
	Fee          AssetAmount       `json:"fee"`
	OwnerAccount ObjectID          `json:"owner_account"`
	AssetType    ObjectID          `json:"asset_type"`
	Balance      int64             `json:"balance"`
	FeeRate      uint32            `json:"fee_rate"`
	Extensions   []json.RawMessage `json:"extensions"`
}

func (o SametFundCreateOperation) Type() OperationType { return OperationTypeSametFundCreate }

func (o SametFundCreateOperation) MarshalJSON() ([]byte, error) {
	type alias SametFundCreateOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *SametFundCreateOperation) UnmarshalJSON(data []byte) error {
	type alias SametFundCreateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeSametFundCreate, &payload); err != nil {
		return err
	}
	*o = SametFundCreateOperation(payload)
	return nil
}

type SametFundDeleteOperation struct {
	Fee          AssetAmount       `json:"fee"`
	OwnerAccount ObjectID          `json:"owner_account"`
	FundID       ObjectID          `json:"fund_id"`
	Extensions   []json.RawMessage `json:"extensions"`
}

func (o SametFundDeleteOperation) Type() OperationType { return OperationTypeSametFundDelete }

func (o SametFundDeleteOperation) MarshalJSON() ([]byte, error) {
	type alias SametFundDeleteOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *SametFundDeleteOperation) UnmarshalJSON(data []byte) error {
	type alias SametFundDeleteOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeSametFundDelete, &payload); err != nil {
		return err
	}
	*o = SametFundDeleteOperation(payload)
	return nil
}

type SametFundUpdateOperation struct {
	Fee          AssetAmount       `json:"fee"`
	OwnerAccount ObjectID          `json:"owner_account"`
	FundID       ObjectID          `json:"fund_id"`
	DeltaAmount  *AssetAmount      `json:"delta_amount,omitempty"`
	NewFeeRate   *uint32           `json:"new_fee_rate,omitempty"`
	Extensions   []json.RawMessage `json:"extensions"`
}

func (o SametFundUpdateOperation) Type() OperationType { return OperationTypeSametFundUpdate }

func (o SametFundUpdateOperation) MarshalJSON() ([]byte, error) {
	type alias SametFundUpdateOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *SametFundUpdateOperation) UnmarshalJSON(data []byte) error {
	type alias SametFundUpdateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeSametFundUpdate, &payload); err != nil {
		return err
	}
	*o = SametFundUpdateOperation(payload)
	return nil
}

type SametFundBorrowOperation struct {
	Fee          AssetAmount       `json:"fee"`
	Borrower     ObjectID          `json:"borrower"`
	FundID       ObjectID          `json:"fund_id"`
	BorrowAmount AssetAmount       `json:"borrow_amount"`
	Extensions   []json.RawMessage `json:"extensions"`
}

func (o SametFundBorrowOperation) Type() OperationType { return OperationTypeSametFundBorrow }

func (o SametFundBorrowOperation) MarshalJSON() ([]byte, error) {
	type alias SametFundBorrowOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *SametFundBorrowOperation) UnmarshalJSON(data []byte) error {
	type alias SametFundBorrowOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeSametFundBorrow, &payload); err != nil {
		return err
	}
	*o = SametFundBorrowOperation(payload)
	return nil
}

type SametFundRepayOperation struct {
	Fee         AssetAmount       `json:"fee"`
	Account     ObjectID          `json:"account"`
	FundID      ObjectID          `json:"fund_id"`
	RepayAmount AssetAmount       `json:"repay_amount"`
	FundFee     AssetAmount       `json:"fund_fee"`
	Extensions  []json.RawMessage `json:"extensions"`
}

func (o SametFundRepayOperation) Type() OperationType { return OperationTypeSametFundRepay }

func (o SametFundRepayOperation) MarshalJSON() ([]byte, error) {
	type alias SametFundRepayOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *SametFundRepayOperation) UnmarshalJSON(data []byte) error {
	type alias SametFundRepayOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeSametFundRepay, &payload); err != nil {
		return err
	}
	*o = SametFundRepayOperation(payload)
	return nil
}

type CreditOfferCollateral struct {
	AssetID ObjectID `json:"asset_id"`
	Price   Price    `json:"price"`
}

type CreditOfferBorrower struct {
	AccountID ObjectID `json:"account_id"`
	Amount    int64    `json:"amount"`
}

type CreditOfferCreateOperation struct {
	Fee                  AssetAmount             `json:"fee"`
	OwnerAccount         ObjectID                `json:"owner_account"`
	AssetType            ObjectID                `json:"asset_type"`
	Balance              int64                   `json:"balance"`
	FeeRate              uint32                  `json:"fee_rate"`
	MaxDurationSeconds   uint32                  `json:"max_duration_seconds"`
	MinDealAmount        int64                   `json:"min_deal_amount"`
	Enabled              bool                    `json:"enabled"`
	AutoDisableTime      Time                    `json:"auto_disable_time"`
	AcceptableCollateral []CreditOfferCollateral `json:"acceptable_collateral,omitempty"`
	AcceptableBorrowers  []CreditOfferBorrower   `json:"acceptable_borrowers,omitempty"`
	Extensions           []json.RawMessage       `json:"extensions"`
}

func (o CreditOfferCreateOperation) Type() OperationType { return OperationTypeCreditOfferCreate }

func (o CreditOfferCreateOperation) MarshalJSON() ([]byte, error) {
	type alias CreditOfferCreateOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *CreditOfferCreateOperation) UnmarshalJSON(data []byte) error {
	type alias CreditOfferCreateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeCreditOfferCreate, &payload); err != nil {
		return err
	}
	*o = CreditOfferCreateOperation(payload)
	return nil
}

type CreditOfferDeleteOperation struct {
	Fee          AssetAmount       `json:"fee"`
	OwnerAccount ObjectID          `json:"owner_account"`
	OfferID      ObjectID          `json:"offer_id"`
	Extensions   []json.RawMessage `json:"extensions"`
}

func (o CreditOfferDeleteOperation) Type() OperationType { return OperationTypeCreditOfferDelete }

func (o CreditOfferDeleteOperation) MarshalJSON() ([]byte, error) {
	type alias CreditOfferDeleteOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *CreditOfferDeleteOperation) UnmarshalJSON(data []byte) error {
	type alias CreditOfferDeleteOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeCreditOfferDelete, &payload); err != nil {
		return err
	}
	*o = CreditOfferDeleteOperation(payload)
	return nil
}

type CreditOfferUpdateOperation struct {
	Fee                  AssetAmount             `json:"fee"`
	OwnerAccount         ObjectID                `json:"owner_account"`
	OfferID              ObjectID                `json:"offer_id"`
	DeltaAmount          *AssetAmount            `json:"delta_amount,omitempty"`
	FeeRate              *uint32                 `json:"fee_rate,omitempty"`
	MaxDurationSeconds   *uint32                 `json:"max_duration_seconds,omitempty"`
	MinDealAmount        *int64                  `json:"min_deal_amount,omitempty"`
	Enabled              *bool                   `json:"enabled,omitempty"`
	AutoDisableTime      *Time                   `json:"auto_disable_time,omitempty"`
	AcceptableCollateral []CreditOfferCollateral `json:"acceptable_collateral,omitempty"`
	AcceptableBorrowers  []CreditOfferBorrower   `json:"acceptable_borrowers,omitempty"`
	Extensions           []json.RawMessage       `json:"extensions"`
}

func (o CreditOfferUpdateOperation) Type() OperationType { return OperationTypeCreditOfferUpdate }

func (o CreditOfferUpdateOperation) MarshalJSON() ([]byte, error) {
	type alias CreditOfferUpdateOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *CreditOfferUpdateOperation) UnmarshalJSON(data []byte) error {
	type alias CreditOfferUpdateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeCreditOfferUpdate, &payload); err != nil {
		return err
	}
	*o = CreditOfferUpdateOperation(payload)
	return nil
}

type CreditOfferAcceptOperation struct {
	Fee                AssetAmount                 `json:"fee"`
	Borrower           ObjectID                    `json:"borrower"`
	OfferID            ObjectID                    `json:"offer_id"`
	BorrowAmount       AssetAmount                 `json:"borrow_amount"`
	Collateral         AssetAmount                 `json:"collateral"`
	MaxFeeRate         uint32                      `json:"max_fee_rate"`
	MinDurationSeconds uint32                      `json:"min_duration_seconds"`
	Extensions         CreditOfferAcceptExtensions `json:"extensions"`
}

func (o CreditOfferAcceptOperation) Type() OperationType { return OperationTypeCreditOfferAccept }

func (o CreditOfferAcceptOperation) MarshalJSON() ([]byte, error) {
	type alias CreditOfferAcceptOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *CreditOfferAcceptOperation) UnmarshalJSON(data []byte) error {
	type alias CreditOfferAcceptOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeCreditOfferAccept, &payload); err != nil {
		return err
	}
	*o = CreditOfferAcceptOperation(payload)
	return nil
}

type CreditDealRepayOperation struct {
	Fee         AssetAmount       `json:"fee"`
	Account     ObjectID          `json:"account"`
	DealID      ObjectID          `json:"deal_id"`
	RepayAmount AssetAmount       `json:"repay_amount"`
	CreditFee   AssetAmount       `json:"credit_fee"`
	Extensions  []json.RawMessage `json:"extensions"`
}

func (o CreditDealRepayOperation) Type() OperationType { return OperationTypeCreditDealRepay }

func (o CreditDealRepayOperation) MarshalJSON() ([]byte, error) {
	type alias CreditDealRepayOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *CreditDealRepayOperation) UnmarshalJSON(data []byte) error {
	type alias CreditDealRepayOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeCreditDealRepay, &payload); err != nil {
		return err
	}
	*o = CreditDealRepayOperation(payload)
	return nil
}

type CreditDealExpiredOperation struct {
	Fee          AssetAmount `json:"fee"`
	DealID       ObjectID    `json:"deal_id"`
	OfferID      ObjectID    `json:"offer_id"`
	OfferOwner   ObjectID    `json:"offer_owner"`
	Borrower     ObjectID    `json:"borrower"`
	UnpaidAmount AssetAmount `json:"unpaid_amount"`
	Collateral   AssetAmount `json:"collateral"`
	FeeRate      uint32      `json:"fee_rate"`
}

func (o CreditDealExpiredOperation) Type() OperationType { return OperationTypeCreditDealExpired }

func (o CreditDealExpiredOperation) MarshalJSON() ([]byte, error) {
	type alias CreditDealExpiredOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *CreditDealExpiredOperation) UnmarshalJSON(data []byte) error {
	type alias CreditDealExpiredOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeCreditDealExpired, &payload); err != nil {
		return err
	}
	*o = CreditDealExpiredOperation(payload)
	return nil
}

type CreditDealUpdateOperation struct {
	Fee        AssetAmount       `json:"fee"`
	Account    ObjectID          `json:"account"`
	DealID     ObjectID          `json:"deal_id"`
	AutoRepay  uint8             `json:"auto_repay"`
	Extensions []json.RawMessage `json:"extensions"`
}

func (o CreditDealUpdateOperation) Type() OperationType { return OperationTypeCreditDealUpdate }

func (o CreditDealUpdateOperation) MarshalJSON() ([]byte, error) {
	type alias CreditDealUpdateOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *CreditDealUpdateOperation) UnmarshalJSON(data []byte) error {
	type alias CreditDealUpdateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeCreditDealUpdate, &payload); err != nil {
		return err
	}
	*o = CreditDealUpdateOperation(payload)
	return nil
}

func init() {
	RegisterOperationFactory(OperationTypeLiquidityPoolCreate, func() Operation { return &LiquidityPoolCreateOperation{} })
	RegisterOperationFactory(OperationTypeLiquidityPoolDelete, func() Operation { return &LiquidityPoolDeleteOperation{} })
	RegisterOperationFactory(OperationTypeLiquidityPoolDeposit, func() Operation { return &LiquidityPoolDepositOperation{} })
	RegisterOperationFactory(OperationTypeLiquidityPoolWithdraw, func() Operation { return &LiquidityPoolWithdrawOperation{} })
	RegisterOperationFactory(OperationTypeLiquidityPoolExchange, func() Operation { return &LiquidityPoolExchangeOperation{} })
	RegisterOperationFactory(OperationTypeLiquidityPoolUpdate, func() Operation { return &LiquidityPoolUpdateOperation{} })
	RegisterOperationFactory(OperationTypeSametFundCreate, func() Operation { return &SametFundCreateOperation{} })
	RegisterOperationFactory(OperationTypeSametFundDelete, func() Operation { return &SametFundDeleteOperation{} })
	RegisterOperationFactory(OperationTypeSametFundUpdate, func() Operation { return &SametFundUpdateOperation{} })
	RegisterOperationFactory(OperationTypeSametFundBorrow, func() Operation { return &SametFundBorrowOperation{} })
	RegisterOperationFactory(OperationTypeSametFundRepay, func() Operation { return &SametFundRepayOperation{} })
	RegisterOperationFactory(OperationTypeCreditOfferCreate, func() Operation { return &CreditOfferCreateOperation{} })
	RegisterOperationFactory(OperationTypeCreditOfferDelete, func() Operation { return &CreditOfferDeleteOperation{} })
	RegisterOperationFactory(OperationTypeCreditOfferUpdate, func() Operation { return &CreditOfferUpdateOperation{} })
	RegisterOperationFactory(OperationTypeCreditOfferAccept, func() Operation { return &CreditOfferAcceptOperation{} })
	RegisterOperationFactory(OperationTypeCreditDealRepay, func() Operation { return &CreditDealRepayOperation{} })
	RegisterOperationFactory(OperationTypeCreditDealExpired, func() Operation { return &CreditDealExpiredOperation{} })
	RegisterOperationFactory(OperationTypeCreditDealUpdate, func() Operation { return &CreditDealUpdateOperation{} })
}

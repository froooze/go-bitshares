package protocol

import "encoding/json"

type CustomOperation struct {
	Fee           AssetAmount     `json:"fee"`
	Payer         ObjectID        `json:"payer"`
	RequiredAuths []ObjectID      `json:"required_auths,omitempty"`
	Id            uint16          `json:"id"`
	Data          json.RawMessage `json:"data"`
}

func (o CustomOperation) Type() OperationType { return OperationTypeCustom }

func (o CustomOperation) MarshalJSON() ([]byte, error) {
	type alias CustomOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *CustomOperation) UnmarshalJSON(data []byte) error {
	type alias CustomOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeCustom, &payload); err != nil {
		return err
	}
	*o = CustomOperation(payload)
	return nil
}

type AssertOperation struct {
	Fee              AssetAmount       `json:"fee"`
	FeePayingAccount ObjectID          `json:"fee_paying_account"`
	Predicates       []Predicate       `json:"predicates"`
	RequiredAuths    []ObjectID        `json:"required_auths,omitempty"`
	Extensions       []json.RawMessage `json:"extensions"`
}

func (o AssertOperation) Type() OperationType { return OperationTypeAssert }

func (o AssertOperation) MarshalJSON() ([]byte, error) {
	type alias AssertOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *AssertOperation) UnmarshalJSON(data []byte) error {
	type alias AssertOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeAssert, &payload); err != nil {
		return err
	}
	*o = AssertOperation(payload)
	return nil
}

type VestingBalanceCreateOperation struct {
	Fee     AssetAmount              `json:"fee"`
	Creator ObjectID                 `json:"creator"`
	Owner   ObjectID                 `json:"owner"`
	Amount  AssetAmount              `json:"amount"`
	Policy  VestingPolicyInitializer `json:"policy"`
}

func (o VestingBalanceCreateOperation) Type() OperationType { return OperationTypeVestingBalanceCreate }

func (o VestingBalanceCreateOperation) MarshalJSON() ([]byte, error) {
	type alias VestingBalanceCreateOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *VestingBalanceCreateOperation) UnmarshalJSON(data []byte) error {
	type alias VestingBalanceCreateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeVestingBalanceCreate, &payload); err != nil {
		return err
	}
	*o = VestingBalanceCreateOperation(payload)
	return nil
}

type VestingBalanceWithdrawOperation struct {
	Fee            AssetAmount `json:"fee"`
	VestingBalance ObjectID    `json:"vesting_balance"`
	Owner          ObjectID    `json:"owner"`
	Amount         AssetAmount `json:"amount"`
}

func (o VestingBalanceWithdrawOperation) Type() OperationType {
	return OperationTypeVestingBalanceWithdraw
}

func (o VestingBalanceWithdrawOperation) MarshalJSON() ([]byte, error) {
	type alias VestingBalanceWithdrawOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *VestingBalanceWithdrawOperation) UnmarshalJSON(data []byte) error {
	type alias VestingBalanceWithdrawOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeVestingBalanceWithdraw, &payload); err != nil {
		return err
	}
	*o = VestingBalanceWithdrawOperation(payload)
	return nil
}

type WorkerCreateOperation struct {
	Fee           AssetAmount       `json:"fee"`
	Owner         ObjectID          `json:"owner"`
	WorkBeginDate Time              `json:"work_begin_date"`
	WorkEndDate   Time              `json:"work_end_date"`
	DailyPay      int64             `json:"daily_pay"`
	Name          string            `json:"name"`
	URL           string            `json:"url"`
	Initializer   WorkerInitializer `json:"initializer"`
}

func (o WorkerCreateOperation) Type() OperationType { return OperationTypeWorkerCreate }

func (o WorkerCreateOperation) MarshalJSON() ([]byte, error) {
	type alias WorkerCreateOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *WorkerCreateOperation) UnmarshalJSON(data []byte) error {
	type alias WorkerCreateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeWorkerCreate, &payload); err != nil {
		return err
	}
	*o = WorkerCreateOperation(payload)
	return nil
}

type HTLCCreateOperation struct {
	Fee                AssetAmount          `json:"fee"`
	From               ObjectID             `json:"from"`
	To                 ObjectID             `json:"to"`
	Amount             AssetAmount          `json:"amount"`
	PreimageHash       HTLCPreimageHash     `json:"preimage_hash"`
	PreimageSize       uint16               `json:"preimage_size"`
	ClaimPeriodSeconds uint32               `json:"claim_period_seconds"`
	Extensions         HTLCCreateExtensions `json:"extensions"`
}

func (o HTLCCreateOperation) Type() OperationType { return OperationTypeHTLCCreate }

func (o HTLCCreateOperation) MarshalJSON() ([]byte, error) {
	type alias HTLCCreateOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *HTLCCreateOperation) UnmarshalJSON(data []byte) error {
	type alias HTLCCreateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeHTLCCreate, &payload); err != nil {
		return err
	}
	*o = HTLCCreateOperation(payload)
	return nil
}

type HTLCRedeemOperation struct {
	Fee        AssetAmount       `json:"fee"`
	HTLCID     ObjectID          `json:"htlc_id"`
	Redeemer   ObjectID          `json:"redeemer"`
	Preimage   HexBytes          `json:"preimage"`
	Extensions []json.RawMessage `json:"extensions"`
}

func (o HTLCRedeemOperation) Type() OperationType { return OperationTypeHTLCRedeem }

func (o HTLCRedeemOperation) MarshalJSON() ([]byte, error) {
	type alias HTLCRedeemOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *HTLCRedeemOperation) UnmarshalJSON(data []byte) error {
	type alias HTLCRedeemOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeHTLCRedeem, &payload); err != nil {
		return err
	}
	*o = HTLCRedeemOperation(payload)
	return nil
}

type HTLCRedeemedOperation struct {
	Fee              AssetAmount      `json:"fee"`
	HTLCID           ObjectID         `json:"htlc_id"`
	From             ObjectID         `json:"from"`
	To               ObjectID         `json:"to"`
	Redeemer         ObjectID         `json:"redeemer"`
	Amount           AssetAmount      `json:"amount"`
	HTLCPreimageHash HTLCPreimageHash `json:"htlc_preimage_hash"`
	HTLCPreimageSize uint16           `json:"htlc_preimage_size"`
	Preimage         HexBytes         `json:"preimage"`
}

func (o HTLCRedeemedOperation) Type() OperationType { return OperationTypeHTLCRedeemed }

func (o HTLCRedeemedOperation) MarshalJSON() ([]byte, error) {
	type alias HTLCRedeemedOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *HTLCRedeemedOperation) UnmarshalJSON(data []byte) error {
	type alias HTLCRedeemedOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeHTLCRedeemed, &payload); err != nil {
		return err
	}
	*o = HTLCRedeemedOperation(payload)
	return nil
}

type HTLCExtendOperation struct {
	Fee          AssetAmount       `json:"fee"`
	HTLCID       ObjectID          `json:"htlc_id"`
	UpdateIssuer ObjectID          `json:"update_issuer"`
	SecondsToAdd uint32            `json:"seconds_to_add"`
	Extensions   []json.RawMessage `json:"extensions"`
}

func (o HTLCExtendOperation) Type() OperationType { return OperationTypeHTLCExtend }

func (o HTLCExtendOperation) MarshalJSON() ([]byte, error) {
	type alias HTLCExtendOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *HTLCExtendOperation) UnmarshalJSON(data []byte) error {
	type alias HTLCExtendOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeHTLCExtend, &payload); err != nil {
		return err
	}
	*o = HTLCExtendOperation(payload)
	return nil
}

type HTLCRefundOperation struct {
	Fee                   AssetAmount      `json:"fee"`
	HTLCID                ObjectID         `json:"htlc_id"`
	To                    ObjectID         `json:"to"`
	OriginalHTLCRecipient ObjectID         `json:"original_htlc_recipient"`
	HTLCAmount            AssetAmount      `json:"htlc_amount"`
	HTLCPreimageHash      HTLCPreimageHash `json:"htlc_preimage_hash"`
	HTLCPreimageSize      uint16           `json:"htlc_preimage_size"`
}

func (o HTLCRefundOperation) Type() OperationType { return OperationTypeHTLCRefund }

func (o HTLCRefundOperation) MarshalJSON() ([]byte, error) {
	type alias HTLCRefundOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *HTLCRefundOperation) UnmarshalJSON(data []byte) error {
	type alias HTLCRefundOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeHTLCRefund, &payload); err != nil {
		return err
	}
	*o = HTLCRefundOperation(payload)
	return nil
}

type TicketCreateOperation struct {
	Fee        AssetAmount       `json:"fee"`
	Account    ObjectID          `json:"account"`
	TargetType uint64            `json:"target_type"`
	Amount     AssetAmount       `json:"amount"`
	Extensions []json.RawMessage `json:"extensions"`
}

func (o TicketCreateOperation) Type() OperationType { return OperationTypeTicketCreate }

func (o TicketCreateOperation) MarshalJSON() ([]byte, error) {
	type alias TicketCreateOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *TicketCreateOperation) UnmarshalJSON(data []byte) error {
	type alias TicketCreateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeTicketCreate, &payload); err != nil {
		return err
	}
	*o = TicketCreateOperation(payload)
	return nil
}

type TicketUpdateOperation struct {
	Fee                AssetAmount       `json:"fee"`
	Ticket             ObjectID          `json:"ticket"`
	Account            ObjectID          `json:"account"`
	TargetType         uint64            `json:"target_type"`
	AmountForNewTarget *AssetAmount      `json:"amount_for_new_target,omitempty"`
	Extensions         []json.RawMessage `json:"extensions"`
}

func (o TicketUpdateOperation) Type() OperationType { return OperationTypeTicketUpdate }

func (o TicketUpdateOperation) MarshalJSON() ([]byte, error) {
	type alias TicketUpdateOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *TicketUpdateOperation) UnmarshalJSON(data []byte) error {
	type alias TicketUpdateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeTicketUpdate, &payload); err != nil {
		return err
	}
	*o = TicketUpdateOperation(payload)
	return nil
}

type CustomAuthorityCreateOperation struct {
	Fee           AssetAmount       `json:"fee"`
	Account       ObjectID          `json:"account"`
	Enabled       bool              `json:"enabled"`
	ValidFrom     Time              `json:"valid_from"`
	ValidTo       Time              `json:"valid_to"`
	OperationType uint64            `json:"operation_type"`
	Auth          Authority         `json:"auth"`
	Restrictions  []Restriction     `json:"restrictions"`
	Extensions    []json.RawMessage `json:"extensions"`
}

func (o CustomAuthorityCreateOperation) Type() OperationType {
	return OperationTypeCustomAuthorityCreate
}

func (o CustomAuthorityCreateOperation) MarshalJSON() ([]byte, error) {
	type alias CustomAuthorityCreateOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *CustomAuthorityCreateOperation) UnmarshalJSON(data []byte) error {
	type alias CustomAuthorityCreateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeCustomAuthorityCreate, &payload); err != nil {
		return err
	}
	*o = CustomAuthorityCreateOperation(payload)
	return nil
}

type CustomAuthorityUpdateOperation struct {
	Fee                  AssetAmount       `json:"fee"`
	Account              ObjectID          `json:"account"`
	AuthorityToUpdate    ObjectID          `json:"authority_to_update"`
	NewEnabled           *bool             `json:"new_enabled,omitempty"`
	NewValidFrom         *Time             `json:"new_valid_from,omitempty"`
	NewValidTo           *Time             `json:"new_valid_to,omitempty"`
	NewAuth              *Authority        `json:"new_auth,omitempty"`
	RestrictionsToRemove []uint16          `json:"restrictions_to_remove,omitempty"`
	RestrictionsToAdd    []Restriction     `json:"restrictions_to_add,omitempty"`
	Extensions           []json.RawMessage `json:"extensions"`
}

func (o CustomAuthorityUpdateOperation) Type() OperationType {
	return OperationTypeCustomAuthorityUpdate
}

func (o CustomAuthorityUpdateOperation) MarshalJSON() ([]byte, error) {
	type alias CustomAuthorityUpdateOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *CustomAuthorityUpdateOperation) UnmarshalJSON(data []byte) error {
	type alias CustomAuthorityUpdateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeCustomAuthorityUpdate, &payload); err != nil {
		return err
	}
	*o = CustomAuthorityUpdateOperation(payload)
	return nil
}

type CustomAuthorityDeleteOperation struct {
	Fee               AssetAmount       `json:"fee"`
	Account           ObjectID          `json:"account"`
	AuthorityToDelete ObjectID          `json:"authority_to_delete"`
	Extensions        []json.RawMessage `json:"extensions"`
}

func (o CustomAuthorityDeleteOperation) Type() OperationType {
	return OperationTypeCustomAuthorityDelete
}

func (o CustomAuthorityDeleteOperation) MarshalJSON() ([]byte, error) {
	type alias CustomAuthorityDeleteOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *CustomAuthorityDeleteOperation) UnmarshalJSON(data []byte) error {
	type alias CustomAuthorityDeleteOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeCustomAuthorityDelete, &payload); err != nil {
		return err
	}
	*o = CustomAuthorityDeleteOperation(payload)
	return nil
}

func init() {
	RegisterOperationFactory(OperationTypeCustom, func() Operation { return &CustomOperation{} })
	RegisterOperationFactory(OperationTypeAssert, func() Operation { return &AssertOperation{} })
	RegisterOperationFactory(OperationTypeVestingBalanceCreate, func() Operation { return &VestingBalanceCreateOperation{} })
	RegisterOperationFactory(OperationTypeVestingBalanceWithdraw, func() Operation { return &VestingBalanceWithdrawOperation{} })
	RegisterOperationFactory(OperationTypeWorkerCreate, func() Operation { return &WorkerCreateOperation{} })
	RegisterOperationFactory(OperationTypeHTLCCreate, func() Operation { return &HTLCCreateOperation{} })
	RegisterOperationFactory(OperationTypeHTLCRedeem, func() Operation { return &HTLCRedeemOperation{} })
	RegisterOperationFactory(OperationTypeHTLCRedeemed, func() Operation { return &HTLCRedeemedOperation{} })
	RegisterOperationFactory(OperationTypeHTLCExtend, func() Operation { return &HTLCExtendOperation{} })
	RegisterOperationFactory(OperationTypeHTLCRefund, func() Operation { return &HTLCRefundOperation{} })
	RegisterOperationFactory(OperationTypeTicketCreate, func() Operation { return &TicketCreateOperation{} })
	RegisterOperationFactory(OperationTypeTicketUpdate, func() Operation { return &TicketUpdateOperation{} })
	RegisterOperationFactory(OperationTypeCustomAuthorityCreate, func() Operation { return &CustomAuthorityCreateOperation{} })
	RegisterOperationFactory(OperationTypeCustomAuthorityUpdate, func() Operation { return &CustomAuthorityUpdateOperation{} })
	RegisterOperationFactory(OperationTypeCustomAuthorityDelete, func() Operation { return &CustomAuthorityDeleteOperation{} })
}

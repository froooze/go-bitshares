package protocol

import (
	"encoding/json"
	"fmt"
)

func marshalOperation(kind OperationType, payload any) ([]byte, error) {
	return json.Marshal([]any{uint16(kind), payload})
}

func unmarshalOperationBody(data []byte, expected OperationType, payload any) error {
	if len(data) > 0 && data[0] != '[' {
		return json.Unmarshal(data, payload)
	}

	var body OperationBody
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	if body.Kind != expected {
		return fmt.Errorf("unexpected operation type %d", body.Kind)
	}
	return json.Unmarshal(body.Payload, payload)
}

type transferOperationJSON struct {
	Fee        AssetAmount       `json:"fee"`
	From       ObjectID          `json:"from"`
	To         ObjectID          `json:"to"`
	Amount     AssetAmount       `json:"amount"`
	Memo       json.RawMessage   `json:"memo,omitempty"`
	Extensions []json.RawMessage `json:"extensions"`
}

// TransferOperation moves funds between accounts.
type TransferOperation struct {
	Fee        AssetAmount
	From       ObjectID
	To         ObjectID
	Amount     AssetAmount
	Memo       json.RawMessage
	Extensions []json.RawMessage
}

func (o TransferOperation) Type() OperationType { return OperationTypeTransfer }

func (o TransferOperation) MarshalJSON() ([]byte, error) {
	payload := transferOperationJSON{
		Fee:        o.Fee,
		From:       o.From,
		To:         o.To,
		Amount:     o.Amount,
		Memo:       o.Memo,
		Extensions: o.Extensions,
	}
	if payload.Extensions == nil {
		payload.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), payload)
}

func (o *TransferOperation) UnmarshalJSON(data []byte) error {
	var payload transferOperationJSON
	if err := unmarshalOperationBody(data, OperationTypeTransfer, &payload); err != nil {
		return err
	}
	o.Fee = payload.Fee
	o.From = payload.From
	o.To = payload.To
	o.Amount = payload.Amount
	o.Memo = payload.Memo
	o.Extensions = payload.Extensions
	return nil
}

type limitOrderCreateOperationJSON struct {
	Fee          AssetAmount                `json:"fee"`
	Seller       ObjectID                   `json:"seller"`
	AmountToSell AssetAmount                `json:"amount_to_sell"`
	MinToReceive AssetAmount                `json:"min_to_receive"`
	Expiration   Time                       `json:"expiration"`
	FillOrKill   bool                       `json:"fill_or_kill"`
	Extensions   LimitOrderCreateExtensions `json:"extensions"`
}

// LimitOrderCreateOperation creates a market order.
type LimitOrderCreateOperation struct {
	Fee          AssetAmount
	Seller       ObjectID
	AmountToSell AssetAmount
	MinToReceive AssetAmount
	Expiration   Time
	FillOrKill   bool
	Extensions   LimitOrderCreateExtensions
}

func (o LimitOrderCreateOperation) Type() OperationType { return OperationTypeLimitOrderCreate }

func (o LimitOrderCreateOperation) MarshalJSON() ([]byte, error) {
	payload := limitOrderCreateOperationJSON{
		Fee:          o.Fee,
		Seller:       o.Seller,
		AmountToSell: o.AmountToSell,
		MinToReceive: o.MinToReceive,
		Expiration:   o.Expiration,
		FillOrKill:   o.FillOrKill,
		Extensions:   o.Extensions,
	}
	return marshalOperation(o.Type(), payload)
}

func (o *LimitOrderCreateOperation) UnmarshalJSON(data []byte) error {
	var payload limitOrderCreateOperationJSON
	if err := unmarshalOperationBody(data, OperationTypeLimitOrderCreate, &payload); err != nil {
		return err
	}
	o.Fee = payload.Fee
	o.Seller = payload.Seller
	o.AmountToSell = payload.AmountToSell
	o.MinToReceive = payload.MinToReceive
	o.Expiration = payload.Expiration
	o.FillOrKill = payload.FillOrKill
	o.Extensions = payload.Extensions
	return nil
}

type limitOrderCancelOperationJSON struct {
	Fee              AssetAmount       `json:"fee"`
	Order            ObjectID          `json:"order"`
	FeePayingAccount ObjectID          `json:"fee_paying_account"`
	Extensions       []json.RawMessage `json:"extensions"`
}

// LimitOrderCancelOperation cancels an existing market order.
type LimitOrderCancelOperation struct {
	Fee              AssetAmount
	Order            ObjectID
	FeePayingAccount ObjectID
	Extensions       []json.RawMessage
}

func (o LimitOrderCancelOperation) Type() OperationType { return OperationTypeLimitOrderCancel }

func (o LimitOrderCancelOperation) MarshalJSON() ([]byte, error) {
	payload := limitOrderCancelOperationJSON{
		Fee:              o.Fee,
		Order:            o.Order,
		FeePayingAccount: o.FeePayingAccount,
		Extensions:       o.Extensions,
	}
	if payload.Extensions == nil {
		payload.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), payload)
}

func (o *LimitOrderCancelOperation) UnmarshalJSON(data []byte) error {
	var payload limitOrderCancelOperationJSON
	if err := unmarshalOperationBody(data, OperationTypeLimitOrderCancel, &payload); err != nil {
		return err
	}
	o.Fee = payload.Fee
	o.Order = payload.Order
	o.FeePayingAccount = payload.FeePayingAccount
	o.Extensions = payload.Extensions
	return nil
}

type limitOrderUpdateOperationJSON struct {
	Fee               AssetAmount            `json:"fee"`
	Seller            ObjectID               `json:"seller"`
	Order             ObjectID               `json:"order"`
	NewPrice          *Price                 `json:"new_price,omitempty"`
	DeltaAmountToSell *AssetAmount           `json:"delta_amount_to_sell,omitempty"`
	NewExpiration     *Time                  `json:"new_expiration,omitempty"`
	OnFill            []LimitOrderAutoAction `json:"on_fill,omitempty"`
	Extensions        []json.RawMessage      `json:"extensions"`
}

// LimitOrderUpdateOperation updates an existing order.
type LimitOrderUpdateOperation struct {
	Fee               AssetAmount
	Seller            ObjectID
	Order             ObjectID
	NewPrice          *Price
	DeltaAmountToSell *AssetAmount
	NewExpiration     *Time
	OnFill            []LimitOrderAutoAction
	Extensions        []json.RawMessage
}

func (o LimitOrderUpdateOperation) Type() OperationType { return OperationTypeLimitOrderUpdate }

func (o LimitOrderUpdateOperation) MarshalJSON() ([]byte, error) {
	payload := limitOrderUpdateOperationJSON{
		Fee:               o.Fee,
		Seller:            o.Seller,
		Order:             o.Order,
		NewPrice:          o.NewPrice,
		DeltaAmountToSell: o.DeltaAmountToSell,
		NewExpiration:     o.NewExpiration,
		OnFill:            o.OnFill,
		Extensions:        o.Extensions,
	}
	if payload.Extensions == nil {
		payload.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), payload)
}

func (o *LimitOrderUpdateOperation) UnmarshalJSON(data []byte) error {
	var payload limitOrderUpdateOperationJSON
	if err := unmarshalOperationBody(data, OperationTypeLimitOrderUpdate, &payload); err != nil {
		return err
	}
	o.Fee = payload.Fee
	o.Seller = payload.Seller
	o.Order = payload.Order
	o.NewPrice = payload.NewPrice
	o.DeltaAmountToSell = payload.DeltaAmountToSell
	o.NewExpiration = payload.NewExpiration
	o.OnFill = payload.OnFill
	o.Extensions = payload.Extensions
	return nil
}

type createTakeProfitOrderActionJSON struct {
	FeeAssetID        ObjectID          `json:"fee_asset_id"`
	SpreadPercent     uint16            `json:"spread_percent"`
	SizePercent       uint16            `json:"size_percent"`
	ExpirationSeconds uint32            `json:"expiration_seconds"`
	Repeat            bool              `json:"repeat"`
	Extensions        []json.RawMessage `json:"extensions"`
}

type LimitOrderCreateExtensions struct {
	OnFill []LimitOrderAutoAction `json:"on_fill,omitempty"`
}

func (e LimitOrderCreateExtensions) MarshalJSON() ([]byte, error) {
	type alias LimitOrderCreateExtensions
	return json.Marshal(alias(e))
}

func (e *LimitOrderCreateExtensions) UnmarshalJSON(data []byte) error {
	type alias LimitOrderCreateExtensions
	var payload alias
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	*e = LimitOrderCreateExtensions(payload)
	return nil
}

func (e LimitOrderCreateExtensions) MarshalBinaryInto(w *binaryWriter) error {
	count := 0
	if len(e.OnFill) != 0 {
		count = 1
	}
	w.writeVarUint64(uint64(count))
	if count == 0 {
		return nil
	}
	w.writeVarUint64(0)
	return writeLimitOrderAutoActions(w, e.OnFill)
}

func (e *LimitOrderCreateExtensions) UnmarshalBinaryFrom(r *binaryReader) error {
	count, err := r.readVarUint64()
	if err != nil {
		return err
	}
	e.OnFill = nil
	for i := uint64(0); i < count; i++ {
		index, err := r.readVarUint64()
		if err != nil {
			return err
		}
		switch index {
		case 0:
			actions, err := readLimitOrderAutoActions(r)
			if err != nil {
				return err
			}
			e.OnFill = actions
		default:
			return fmt.Errorf("unknown limit order create extension %d", index)
		}
	}
	return nil
}

// CreateTakeProfitOrderAction is the only currently supported limit-order auto action.
type CreateTakeProfitOrderAction struct {
	FeeAssetID        ObjectID
	SpreadPercent     uint16
	SizePercent       uint16
	ExpirationSeconds uint32
	Repeat            bool
	Extensions        []json.RawMessage
}

func (o CreateTakeProfitOrderAction) MarshalJSON() ([]byte, error) {
	payload := createTakeProfitOrderActionJSON{
		FeeAssetID:        o.FeeAssetID,
		SpreadPercent:     o.SpreadPercent,
		SizePercent:       o.SizePercent,
		ExpirationSeconds: o.ExpirationSeconds,
		Repeat:            o.Repeat,
		Extensions:        o.Extensions,
	}
	if payload.Extensions == nil {
		payload.Extensions = []json.RawMessage{}
	}
	return json.Marshal(payload)
}

func (o *CreateTakeProfitOrderAction) UnmarshalJSON(data []byte) error {
	var payload createTakeProfitOrderActionJSON
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	o.FeeAssetID = payload.FeeAssetID
	o.SpreadPercent = payload.SpreadPercent
	o.SizePercent = payload.SizePercent
	o.ExpirationSeconds = payload.ExpirationSeconds
	o.Repeat = payload.Repeat
	o.Extensions = payload.Extensions
	return nil
}

func (o CreateTakeProfitOrderAction) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	if err := o.FeeAssetID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeUint16(o.SpreadPercent)
	w.writeUint16(o.SizePercent)
	w.writeUint32(o.ExpirationSeconds)
	w.writeUint8(boolByte(o.Repeat))
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("create_take_profit_order_action extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *CreateTakeProfitOrderAction) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *CreateTakeProfitOrderAction) UnmarshalBinaryFrom(r *binaryReader) error {
	feeAsset, err := readObjectID(r)
	if err != nil {
		return err
	}
	spread, err := r.readUint16()
	if err != nil {
		return err
	}
	size, err := r.readUint16()
	if err != nil {
		return err
	}
	expiration, err := r.readUint32()
	if err != nil {
		return err
	}
	repeat, err := r.readUint8()
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("create_take_profit_order_action extensions are not supported in binary serialization")
	}
	o.FeeAssetID = feeAsset
	o.SpreadPercent = spread
	o.SizePercent = size
	o.ExpirationSeconds = expiration
	o.Repeat = repeat != 0
	o.Extensions = nil
	return nil
}

type LimitOrderAutoAction struct {
	Kind       uint16
	TakeProfit *CreateTakeProfitOrderAction
}

func (a LimitOrderAutoAction) MarshalJSON() ([]byte, error) {
	switch a.Kind {
	case 0:
		if a.TakeProfit == nil {
			return nil, fmt.Errorf("missing take profit action")
		}
		return json.Marshal([]any{uint16(0), a.TakeProfit})
	default:
		return nil, fmt.Errorf("unsupported limit order auto action type %d", a.Kind)
	}
}

func (a *LimitOrderAutoAction) UnmarshalJSON(data []byte) error {
	var body []json.RawMessage
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	if len(body) != 2 {
		return fmt.Errorf("invalid limit order auto action")
	}
	var kind uint16
	if err := json.Unmarshal(body[0], &kind); err != nil {
		return err
	}
	a.Kind = kind
	switch kind {
	case 0:
		var action CreateTakeProfitOrderAction
		if err := json.Unmarshal(body[1], &action); err != nil {
			return err
		}
		a.TakeProfit = &action
		return nil
	default:
		return fmt.Errorf("unsupported limit order auto action type %d", kind)
	}
}

func (a LimitOrderAutoAction) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(a.Kind))
	switch a.Kind {
	case 0:
		if a.TakeProfit == nil {
			return nil, fmt.Errorf("missing take profit action")
		}
		raw, err := a.TakeProfit.MarshalBinary()
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(raw); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported limit order auto action type %d", a.Kind)
	}
	return w.Bytes(), nil
}

func (a *LimitOrderAutoAction) UnmarshalBinary(data []byte) error {
	return a.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (a *LimitOrderAutoAction) UnmarshalBinaryFrom(r *binaryReader) error {
	kind, err := r.readVarUint64()
	if err != nil {
		return err
	}
	a.Kind = uint16(kind)
	switch a.Kind {
	case 0:
		var action CreateTakeProfitOrderAction
		if err := action.UnmarshalBinaryFrom(r); err != nil {
			return err
		}
		a.TakeProfit = &action
		return nil
	default:
		return fmt.Errorf("unsupported limit order auto action type %d", a.Kind)
	}
}

type assetUpdateIssuerOperationJSON struct {
	Fee           AssetAmount       `json:"fee"`
	Issuer        ObjectID          `json:"issuer"`
	AssetToUpdate ObjectID          `json:"asset_to_update"`
	NewIssuer     ObjectID          `json:"new_issuer"`
	Extensions    []json.RawMessage `json:"extensions"`
}

// AssetUpdateIssuerOperation changes the controlling issuer of an asset.
type AssetUpdateIssuerOperation struct {
	Fee           AssetAmount
	Issuer        ObjectID
	AssetToUpdate ObjectID
	NewIssuer     ObjectID
	Extensions    []json.RawMessage
}

func (o AssetUpdateIssuerOperation) Type() OperationType { return OperationTypeAssetUpdateIssuer }

func (o AssetUpdateIssuerOperation) MarshalJSON() ([]byte, error) {
	payload := assetUpdateIssuerOperationJSON{
		Fee:           o.Fee,
		Issuer:        o.Issuer,
		AssetToUpdate: o.AssetToUpdate,
		NewIssuer:     o.NewIssuer,
		Extensions:    o.Extensions,
	}
	if payload.Extensions == nil {
		payload.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), payload)
}

func (o *AssetUpdateIssuerOperation) UnmarshalJSON(data []byte) error {
	var payload assetUpdateIssuerOperationJSON
	if err := unmarshalOperationBody(data, OperationTypeAssetUpdateIssuer, &payload); err != nil {
		return err
	}
	o.Fee = payload.Fee
	o.Issuer = payload.Issuer
	o.AssetToUpdate = payload.AssetToUpdate
	o.NewIssuer = payload.NewIssuer
	o.Extensions = payload.Extensions
	return nil
}

type assetClaimPoolOperationJSON struct {
	Fee           AssetAmount       `json:"fee"`
	Issuer        ObjectID          `json:"issuer"`
	AssetID       ObjectID          `json:"asset_id"`
	AmountToClaim AssetAmount       `json:"amount_to_claim"`
	Extensions    []json.RawMessage `json:"extensions"`
}

// AssetClaimPoolOperation withdraws BTS from an asset fee pool.
type AssetClaimPoolOperation struct {
	Fee           AssetAmount
	Issuer        ObjectID
	AssetID       ObjectID
	AmountToClaim AssetAmount
	Extensions    []json.RawMessage
}

func (o AssetClaimPoolOperation) Type() OperationType { return OperationTypeAssetClaimPool }

func (o AssetClaimPoolOperation) MarshalJSON() ([]byte, error) {
	payload := assetClaimPoolOperationJSON{
		Fee:           o.Fee,
		Issuer:        o.Issuer,
		AssetID:       o.AssetID,
		AmountToClaim: o.AmountToClaim,
		Extensions:    o.Extensions,
	}
	if payload.Extensions == nil {
		payload.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), payload)
}

func (o *AssetClaimPoolOperation) UnmarshalJSON(data []byte) error {
	var payload assetClaimPoolOperationJSON
	if err := unmarshalOperationBody(data, OperationTypeAssetClaimPool, &payload); err != nil {
		return err
	}
	o.Fee = payload.Fee
	o.Issuer = payload.Issuer
	o.AssetID = payload.AssetID
	o.AmountToClaim = payload.AmountToClaim
	o.Extensions = payload.Extensions
	return nil
}

type assetIssueOperationJSON struct {
	Fee            AssetAmount       `json:"fee"`
	Issuer         ObjectID          `json:"issuer"`
	AssetToIssue   AssetAmount       `json:"asset_to_issue"`
	IssueToAccount ObjectID          `json:"issue_to_account"`
	Memo           json.RawMessage   `json:"memo,omitempty"`
	Extensions     []json.RawMessage `json:"extensions"`
}

// AssetIssueOperation issues an asset to an account.
type AssetIssueOperation struct {
	Fee            AssetAmount
	Issuer         ObjectID
	AssetToIssue   AssetAmount
	IssueToAccount ObjectID
	Memo           json.RawMessage
	Extensions     []json.RawMessage
}

func (o AssetIssueOperation) Type() OperationType { return OperationTypeAssetIssue }

func (o AssetIssueOperation) MarshalJSON() ([]byte, error) {
	payload := assetIssueOperationJSON{
		Fee:            o.Fee,
		Issuer:         o.Issuer,
		AssetToIssue:   o.AssetToIssue,
		IssueToAccount: o.IssueToAccount,
		Memo:           o.Memo,
		Extensions:     o.Extensions,
	}
	if payload.Extensions == nil {
		payload.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), payload)
}

func (o *AssetIssueOperation) UnmarshalJSON(data []byte) error {
	var payload assetIssueOperationJSON
	if err := unmarshalOperationBody(data, OperationTypeAssetIssue, &payload); err != nil {
		return err
	}
	o.Fee = payload.Fee
	o.Issuer = payload.Issuer
	o.AssetToIssue = payload.AssetToIssue
	o.IssueToAccount = payload.IssueToAccount
	o.Memo = payload.Memo
	o.Extensions = payload.Extensions
	return nil
}

type assetReserveOperationJSON struct {
	Fee             AssetAmount       `json:"fee"`
	Payer           ObjectID          `json:"payer"`
	AmountToReserve AssetAmount       `json:"amount_to_reserve"`
	Extensions      []json.RawMessage `json:"extensions"`
}

// AssetReserveOperation reserves issued asset supply.
type AssetReserveOperation struct {
	Fee             AssetAmount
	Payer           ObjectID
	AmountToReserve AssetAmount
	Extensions      []json.RawMessage
}

func (o AssetReserveOperation) Type() OperationType { return OperationTypeAssetReserve }

func (o AssetReserveOperation) MarshalJSON() ([]byte, error) {
	payload := assetReserveOperationJSON{
		Fee:             o.Fee,
		Payer:           o.Payer,
		AmountToReserve: o.AmountToReserve,
		Extensions:      o.Extensions,
	}
	if payload.Extensions == nil {
		payload.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), payload)
}

func (o *AssetReserveOperation) UnmarshalJSON(data []byte) error {
	var payload assetReserveOperationJSON
	if err := unmarshalOperationBody(data, OperationTypeAssetReserve, &payload); err != nil {
		return err
	}
	o.Fee = payload.Fee
	o.Payer = payload.Payer
	o.AmountToReserve = payload.AmountToReserve
	o.Extensions = payload.Extensions
	return nil
}

func init() {
	RegisterOperationFactory(OperationTypeTransfer, func() Operation { return &TransferOperation{} })
	RegisterOperationFactory(OperationTypeLimitOrderCreate, func() Operation { return &LimitOrderCreateOperation{} })
	RegisterOperationFactory(OperationTypeLimitOrderCancel, func() Operation { return &LimitOrderCancelOperation{} })
	RegisterOperationFactory(OperationTypeLimitOrderUpdate, func() Operation { return &LimitOrderUpdateOperation{} })
	RegisterOperationFactory(OperationTypeAssetUpdateIssuer, func() Operation { return &AssetUpdateIssuerOperation{} })
	RegisterOperationFactory(OperationTypeAssetClaimPool, func() Operation { return &AssetClaimPoolOperation{} })
	RegisterOperationFactory(OperationTypeAssetIssue, func() Operation { return &AssetIssueOperation{} })
	RegisterOperationFactory(OperationTypeAssetReserve, func() Operation { return &AssetReserveOperation{} })
}

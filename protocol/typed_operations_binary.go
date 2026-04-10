package protocol

import (
	"encoding/json"
	"fmt"
)

func (o TransferOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.From.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.To.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Amount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := writeOptionalMemo(w, o.Memo); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("transfer extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *TransferOperation) unmarshalBinaryBody(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	from, err := readObjectID(r)
	if err != nil {
		return err
	}
	to, err := readObjectID(r)
	if err != nil {
		return err
	}
	amount, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	memo, err := readOptionalMemo(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("transfer extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.From = from
	o.To = to
	o.Amount = amount
	if memo != nil {
		raw, err := json.Marshal(memo)
		if err != nil {
			return err
		}
		o.Memo = raw
	} else {
		o.Memo = nil
	}
	o.Extensions = nil
	return nil
}

func (o LimitOrderCreateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Seller.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AmountToSell.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.MinToReceive.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Expiration.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeUint8(boolByte(o.FillOrKill))
	if err := o.Extensions.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *LimitOrderCreateOperation) unmarshalBinaryBody(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	seller, err := readObjectID(r)
	if err != nil {
		return err
	}
	sell, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	buy, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	exp, err := readTime(r)
	if err != nil {
		return err
	}
	fill, err := r.readUint8()
	if err != nil {
		return err
	}
	var ext LimitOrderCreateExtensions
	if err := ext.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	o.Fee = fee
	o.Seller = seller
	o.AmountToSell = sell
	o.MinToReceive = buy
	o.Expiration = exp
	o.FillOrKill = fill != 0
	o.Extensions = ext
	return nil
}

func (o LimitOrderCancelOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Order.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.FeePayingAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("limit order cancel extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *LimitOrderCancelOperation) unmarshalBinaryBody(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	order, err := readObjectID(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("limit order cancel extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Order = order
	o.FeePayingAccount = account
	o.Extensions = nil
	return nil
}

func (o LimitOrderUpdateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Seller.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Order.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := writeOptionalPrice(w, o.NewPrice); err != nil {
		return nil, err
	}
	if err := writeOptionalAssetAmount(w, o.DeltaAmountToSell); err != nil {
		return nil, err
	}
	if err := writeOptionalTime(w, o.NewExpiration); err != nil {
		return nil, err
	}
	if err := writeOptionalAutoActions(w, o.OnFill); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("limit order update extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *LimitOrderUpdateOperation) unmarshalBinaryBody(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	seller, err := readObjectID(r)
	if err != nil {
		return err
	}
	order, err := readObjectID(r)
	if err != nil {
		return err
	}
	price, err := readOptionalPrice(r)
	if err != nil {
		return err
	}
	delta, err := readOptionalAssetAmount(r)
	if err != nil {
		return err
	}
	exp, err := readOptionalTime(r)
	if err != nil {
		return err
	}
	onFill, err := readOptionalAutoActions(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("limit order update extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Seller = seller
	o.Order = order
	o.NewPrice = price
	o.DeltaAmountToSell = delta
	o.NewExpiration = exp
	o.OnFill = onFill
	o.Extensions = nil
	return nil
}

func (o AssetUpdateIssuerOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Issuer.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AssetToUpdate.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.NewIssuer.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("asset update issuer extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *AssetUpdateIssuerOperation) unmarshalBinaryBody(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	issuer, err := readObjectID(r)
	if err != nil {
		return err
	}
	asset, err := readObjectID(r)
	if err != nil {
		return err
	}
	newIssuer, err := readObjectID(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("asset update issuer extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Issuer = issuer
	o.AssetToUpdate = asset
	o.NewIssuer = newIssuer
	o.Extensions = nil
	return nil
}

func (o AssetClaimPoolOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Issuer.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AssetID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AmountToClaim.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("asset claim pool extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *AssetClaimPoolOperation) unmarshalBinaryBody(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	issuer, err := readObjectID(r)
	if err != nil {
		return err
	}
	asset, err := readObjectID(r)
	if err != nil {
		return err
	}
	claim, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("asset claim pool extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Issuer = issuer
	o.AssetID = asset
	o.AmountToClaim = claim
	o.Extensions = nil
	return nil
}

func (o AssetIssueOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Issuer.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AssetToIssue.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.IssueToAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := writeOptionalMemo(w, o.Memo); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("asset issue extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *AssetIssueOperation) unmarshalBinaryBody(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	issuer, err := readObjectID(r)
	if err != nil {
		return err
	}
	asset, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	memo, err := readOptionalMemo(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("asset issue extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Issuer = issuer
	o.AssetToIssue = asset
	o.IssueToAccount = account
	if memo != nil {
		raw, err := json.Marshal(memo)
		if err != nil {
			return err
		}
		o.Memo = raw
	} else {
		o.Memo = nil
	}
	o.Extensions = nil
	return nil
}

func (o AssetReserveOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Payer.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AmountToReserve.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("asset reserve extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *AssetReserveOperation) unmarshalBinaryBody(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	payer, err := readObjectID(r)
	if err != nil {
		return err
	}
	amount, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("asset reserve extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Payer = payer
	o.AmountToReserve = amount
	o.Extensions = nil
	return nil
}

func writeOptionalMemo(w *binaryWriter, raw json.RawMessage) error {
	memo, err := memoDataFromRaw(raw)
	if err != nil {
		return err
	}
	if memo == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	return memo.writeBinary(w)
}

func readOptionalMemo(r *binaryReader) (*MemoData, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	var memo MemoData
	if err := memo.readBinary(r); err != nil {
		return nil, err
	}
	return &memo, nil
}

func writeOptionalPrice(w *binaryWriter, value *Price) error {
	if value == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	return value.MarshalBinaryInto(w)
}

func readOptionalPrice(r *binaryReader) (*Price, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	value, err := readPrice(r)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func writeOptionalAssetAmount(w *binaryWriter, value *AssetAmount) error {
	if value == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	return value.MarshalBinaryInto(w)
}

func readOptionalAssetAmount(r *binaryReader) (*AssetAmount, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	value, err := readAssetAmount(r)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func writeOptionalTime(w *binaryWriter, value *Time) error {
	if value == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	return value.MarshalBinaryInto(w)
}

func readOptionalTime(r *binaryReader) (*Time, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	value, err := readTime(r)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func writeOptionalAutoActions(w *binaryWriter, actions []LimitOrderAutoAction) error {
	if len(actions) == 0 {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	w.writeVarUint64(uint64(len(actions)))
	for _, action := range actions {
		raw, err := action.MarshalBinary()
		if err != nil {
			return err
		}
		if _, err := w.Write(raw); err != nil {
			return err
		}
	}
	return nil
}

func readOptionalAutoActions(r *binaryReader) ([]LimitOrderAutoAction, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	actions := make([]LimitOrderAutoAction, 0, count)
	for i := uint64(0); i < count; i++ {
		var action LimitOrderAutoAction
		if err := action.UnmarshalBinaryFrom(r); err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	return actions, nil
}

func writeLimitOrderAutoActions(w *binaryWriter, actions []LimitOrderAutoAction) error {
	w.writeVarUint64(uint64(len(actions)))
	for _, action := range actions {
		raw, err := action.MarshalBinary()
		if err != nil {
			return err
		}
		if _, err := w.Write(raw); err != nil {
			return err
		}
	}
	return nil
}

func readLimitOrderAutoActions(r *binaryReader) ([]LimitOrderAutoAction, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	actions := make([]LimitOrderAutoAction, 0, count)
	for i := uint64(0); i < count; i++ {
		var action LimitOrderAutoAction
		if err := action.UnmarshalBinaryFrom(r); err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	return actions, nil
}

func boolByte(v bool) uint8 {
	if v {
		return 1
	}
	return 0
}

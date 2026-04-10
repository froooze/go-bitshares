package protocol

import "fmt"

func (o AssertOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.FeePayingAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := writePredicateArray(w, o.Predicates); err != nil {
		return nil, err
	}
	if err := writeObjectIDSet(w, o.RequiredAuths); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("assert extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *AssertOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *AssertOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	predicates, err := readPredicateArray(r)
	if err != nil {
		return err
	}
	required, err := readObjectIDSet(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("assert extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.FeePayingAccount = account
	o.Predicates = predicates
	o.RequiredAuths = required
	o.Extensions = nil
	return nil
}

func (o VestingBalanceCreateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Creator.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Owner.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Amount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	raw, err := o.Policy.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(raw); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *VestingBalanceCreateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *VestingBalanceCreateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	creator, err := readObjectID(r)
	if err != nil {
		return err
	}
	owner, err := readObjectID(r)
	if err != nil {
		return err
	}
	amount, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	var policy VestingPolicyInitializer
	if err := policy.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	o.Fee = fee
	o.Creator = creator
	o.Owner = owner
	o.Amount = amount
	o.Policy = policy
	return nil
}

func (o VestingBalanceWithdrawOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.VestingBalance.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Owner.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Amount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *VestingBalanceWithdrawOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *VestingBalanceWithdrawOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	vb, err := readObjectID(r)
	if err != nil {
		return err
	}
	owner, err := readObjectID(r)
	if err != nil {
		return err
	}
	amount, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	o.Fee = fee
	o.VestingBalance = vb
	o.Owner = owner
	o.Amount = amount
	return nil
}

func (o WorkerCreateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Owner.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.WorkBeginDate.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.WorkEndDate.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeInt64(o.DailyPay)
	w.writeString(o.Name)
	w.writeString(o.URL)
	raw, err := o.Initializer.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(raw); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *WorkerCreateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *WorkerCreateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	owner, err := readObjectID(r)
	if err != nil {
		return err
	}
	begin, err := readTime(r)
	if err != nil {
		return err
	}
	end, err := readTime(r)
	if err != nil {
		return err
	}
	daily, err := r.readInt64()
	if err != nil {
		return err
	}
	name, err := r.readString()
	if err != nil {
		return err
	}
	url, err := r.readString()
	if err != nil {
		return err
	}
	var init WorkerInitializer
	if err := init.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	o.Fee = fee
	o.Owner = owner
	o.WorkBeginDate = begin
	o.WorkEndDate = end
	o.DailyPay = daily
	o.Name = name
	o.URL = url
	o.Initializer = init
	return nil
}

func (o HTLCCreateOperation) MarshalBinary() ([]byte, error) {
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
	raw, err := o.PreimageHash.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(raw); err != nil {
		return nil, err
	}
	w.writeUint16(o.PreimageSize)
	w.writeUint32(o.ClaimPeriodSeconds)
	if err := o.Extensions.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *HTLCCreateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *HTLCCreateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
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
	var hash HTLCPreimageHash
	if err := hash.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	size, err := r.readUint16()
	if err != nil {
		return err
	}
	period, err := r.readUint32()
	if err != nil {
		return err
	}
	var ext HTLCCreateExtensions
	if err := ext.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	o.Fee = fee
	o.From = from
	o.To = to
	o.Amount = amount
	o.PreimageHash = hash
	o.PreimageSize = size
	o.ClaimPeriodSeconds = period
	o.Extensions = ext
	return nil
}

func (o HTLCRedeemOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.HTLCID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Redeemer.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Preimage.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *HTLCRedeemOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *HTLCRedeemOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	htlc, err := readObjectID(r)
	if err != nil {
		return err
	}
	redeemer, err := readObjectID(r)
	if err != nil {
		return err
	}
	var preimage HexBytes
	if err := preimage.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("htlc redeem extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.HTLCID = htlc
	o.Redeemer = redeemer
	o.Preimage = preimage
	o.Extensions = nil
	return nil
}

func (o HTLCRedeemedOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.HTLCID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.From.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.To.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Redeemer.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Amount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	raw, err := o.HTLCPreimageHash.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(raw); err != nil {
		return nil, err
	}
	w.writeUint16(o.HTLCPreimageSize)
	if err := o.Preimage.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *HTLCRedeemedOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *HTLCRedeemedOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	htlc, err := readObjectID(r)
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
	redeemer, err := readObjectID(r)
	if err != nil {
		return err
	}
	amount, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	var hash HTLCPreimageHash
	if err := hash.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	size, err := r.readUint16()
	if err != nil {
		return err
	}
	var preimage HexBytes
	if err := preimage.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	o.Fee = fee
	o.HTLCID = htlc
	o.From = from
	o.To = to
	o.Redeemer = redeemer
	o.Amount = amount
	o.HTLCPreimageHash = hash
	o.HTLCPreimageSize = size
	o.Preimage = preimage
	return nil
}

func (o HTLCExtendOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.HTLCID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.UpdateIssuer.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeUint32(o.SecondsToAdd)
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *HTLCExtendOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *HTLCExtendOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	htlc, err := readObjectID(r)
	if err != nil {
		return err
	}
	issuer, err := readObjectID(r)
	if err != nil {
		return err
	}
	seconds, err := r.readUint32()
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("htlc extend extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.HTLCID = htlc
	o.UpdateIssuer = issuer
	o.SecondsToAdd = seconds
	o.Extensions = nil
	return nil
}

func (o HTLCRefundOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.HTLCID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.To.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.OriginalHTLCRecipient.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.HTLCAmount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	raw, err := o.HTLCPreimageHash.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(raw); err != nil {
		return nil, err
	}
	w.writeUint16(o.HTLCPreimageSize)
	return w.Bytes(), nil
}

func (o *HTLCRefundOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *HTLCRefundOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	htlc, err := readObjectID(r)
	if err != nil {
		return err
	}
	to, err := readObjectID(r)
	if err != nil {
		return err
	}
	originalRecipient, err := readObjectID(r)
	if err != nil {
		return err
	}
	amount, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	var hash HTLCPreimageHash
	if err := hash.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	size, err := r.readUint16()
	if err != nil {
		return err
	}
	o.Fee = fee
	o.HTLCID = htlc
	o.To = to
	o.OriginalHTLCRecipient = originalRecipient
	o.HTLCAmount = amount
	o.HTLCPreimageHash = hash
	o.HTLCPreimageSize = size
	return nil
}

func (o CustomAuthorityCreateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeUint8(boolByte(o.Enabled))
	if err := o.ValidFrom.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.ValidTo.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeVarUint64(o.OperationType)
	if err := o.Auth.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := writeRestrictionArray(w, o.Restrictions); err != nil {
		return nil, err
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *CustomAuthorityCreateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *CustomAuthorityCreateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	enabled, err := readBool(r)
	if err != nil {
		return err
	}
	from, err := readTime(r)
	if err != nil {
		return err
	}
	to, err := readTime(r)
	if err != nil {
		return err
	}
	opType, err := r.readVarUint64()
	if err != nil {
		return err
	}
	var auth Authority
	if err := auth.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	restrictions, err := readRestrictionArray(r)
	if err != nil {
		return err
	}
	if _, err := r.readVarUint64(); err != nil {
		return err
	}
	o.Fee = fee
	o.Account = account
	o.Enabled = enabled
	o.ValidFrom = from
	o.ValidTo = to
	o.OperationType = opType
	o.Auth = auth
	o.Restrictions = restrictions
	o.Extensions = nil
	return nil
}

func (o CustomAuthorityUpdateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AuthorityToUpdate.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := writeOptionalBool(w, o.NewEnabled); err != nil {
		return nil, err
	}
	if err := writeOptionalTime(w, o.NewValidFrom); err != nil {
		return nil, err
	}
	if err := writeOptionalTime(w, o.NewValidTo); err != nil {
		return nil, err
	}
	if err := writeOptionalAuthority(w, o.NewAuth); err != nil {
		return nil, err
	}
	w.writeVarUint64(uint64(len(o.RestrictionsToRemove)))
	for _, value := range o.RestrictionsToRemove {
		w.writeUint16(value)
	}
	if err := writeRestrictionArray(w, o.RestrictionsToAdd); err != nil {
		return nil, err
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *CustomAuthorityUpdateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *CustomAuthorityUpdateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	authority, err := readObjectID(r)
	if err != nil {
		return err
	}
	enabled, err := readOptionalBool(r)
	if err != nil {
		return err
	}
	validFrom, err := readOptionalTime(r)
	if err != nil {
		return err
	}
	validTo, err := readOptionalTime(r)
	if err != nil {
		return err
	}
	auth, err := readOptionalAuthority(r)
	if err != nil {
		return err
	}
	removeCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	remove := make([]uint16, 0, removeCount)
	for i := uint64(0); i < removeCount; i++ {
		value, err := r.readUint16()
		if err != nil {
			return err
		}
		remove = append(remove, value)
	}
	add, err := readRestrictionArray(r)
	if err != nil {
		return err
	}
	if _, err := r.readVarUint64(); err != nil {
		return err
	}
	o.Fee = fee
	o.Account = account
	o.AuthorityToUpdate = authority
	o.NewEnabled = enabled
	o.NewValidFrom = validFrom
	o.NewValidTo = validTo
	o.NewAuth = auth
	o.RestrictionsToRemove = remove
	o.RestrictionsToAdd = add
	o.Extensions = nil
	return nil
}

func (o CustomAuthorityDeleteOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AuthorityToDelete.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *CustomAuthorityDeleteOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *CustomAuthorityDeleteOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	authority, err := readObjectID(r)
	if err != nil {
		return err
	}
	if _, err := r.readVarUint64(); err != nil {
		return err
	}
	o.Fee = fee
	o.Account = account
	o.AuthorityToDelete = authority
	o.Extensions = nil
	return nil
}

func writePredicateArray(w *binaryWriter, values []Predicate) error {
	w.writeVarUint64(uint64(len(values)))
	for _, value := range values {
		raw, err := value.MarshalBinary()
		if err != nil {
			return err
		}
		if _, err := w.Write(raw); err != nil {
			return err
		}
	}
	return nil
}

func readPredicateArray(r *binaryReader) ([]Predicate, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make([]Predicate, 0, count)
	for i := uint64(0); i < count; i++ {
		var value Predicate
		if err := value.UnmarshalBinaryFrom(r); err != nil {
			return nil, err
		}
		out = append(out, value)
	}
	return out, nil
}

func writeOptionalBool(w *binaryWriter, value *bool) error {
	if value == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	w.writeUint8(boolByte(*value))
	return nil
}

func readOptionalBool(r *binaryReader) (*bool, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	value, err := readBool(r)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

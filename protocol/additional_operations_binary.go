package protocol

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

func (o CustomOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Payer.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := writeObjectIDSet(w, o.RequiredAuths); err != nil {
		return nil, err
	}
	w.writeUint16(o.Id)
	data, err := customDataFromRaw(o.Data)
	if err != nil {
		return nil, err
	}
	w.writeBytes(data)
	return w.Bytes(), nil
}

func (o *CustomOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	payer, err := readObjectID(r)
	if err != nil {
		return err
	}
	requiredAuths, err := readObjectIDSet(r)
	if err != nil {
		return err
	}
	id, err := r.readUint16()
	if err != nil {
		return err
	}
	data, err := r.readBytes()
	if err != nil {
		return err
	}
	raw, err := json.Marshal(hex.EncodeToString(data))
	if err != nil {
		return err
	}
	o.Fee = fee
	o.Payer = payer
	o.RequiredAuths = requiredAuths
	o.Id = id
	o.Data = raw
	return nil
}

func (o *CustomOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o TicketCreateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeVarUint64(o.TargetType)
	if err := o.Amount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("ticket create extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *TicketCreateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	targetType, err := r.readVarUint64()
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
		return fmt.Errorf("ticket create extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Account = account
	o.TargetType = targetType
	o.Amount = amount
	o.Extensions = nil
	return nil
}

func (o *TicketCreateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o TicketUpdateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Ticket.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeVarUint64(o.TargetType)
	if err := writeOptionalAssetAmount(w, o.AmountForNewTarget); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("ticket update extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *TicketUpdateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	ticket, err := readObjectID(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	targetType, err := r.readVarUint64()
	if err != nil {
		return err
	}
	amountForNewTarget, err := readOptionalAssetAmount(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("ticket update extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Ticket = ticket
	o.Account = account
	o.TargetType = targetType
	o.AmountForNewTarget = amountForNewTarget
	o.Extensions = nil
	return nil
}

func (o *TicketUpdateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o LiquidityPoolCreateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AssetA.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AssetB.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.ShareAsset.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeUint16(o.TakerFeePercent)
	w.writeUint16(o.WithdrawalFeePercent)
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("liquidity pool create extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *LiquidityPoolCreateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	assetA, err := readObjectID(r)
	if err != nil {
		return err
	}
	assetB, err := readObjectID(r)
	if err != nil {
		return err
	}
	shareAsset, err := readObjectID(r)
	if err != nil {
		return err
	}
	takerFeePercent, err := r.readUint16()
	if err != nil {
		return err
	}
	withdrawalFeePercent, err := r.readUint16()
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("liquidity pool create extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Account = account
	o.AssetA = assetA
	o.AssetB = assetB
	o.ShareAsset = shareAsset
	o.TakerFeePercent = takerFeePercent
	o.WithdrawalFeePercent = withdrawalFeePercent
	o.Extensions = nil
	return nil
}

func (o *LiquidityPoolCreateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o LiquidityPoolDeleteOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Pool.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("liquidity pool delete extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *LiquidityPoolDeleteOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	pool, err := readObjectID(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("liquidity pool delete extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Account = account
	o.Pool = pool
	o.Extensions = nil
	return nil
}

func (o *LiquidityPoolDeleteOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o LiquidityPoolDepositOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Pool.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AmountA.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AmountB.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("liquidity pool deposit extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *LiquidityPoolDepositOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	pool, err := readObjectID(r)
	if err != nil {
		return err
	}
	amountA, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	amountB, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("liquidity pool deposit extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Account = account
	o.Pool = pool
	o.AmountA = amountA
	o.AmountB = amountB
	o.Extensions = nil
	return nil
}

func (o *LiquidityPoolDepositOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o LiquidityPoolWithdrawOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Pool.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.ShareAmount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("liquidity pool withdraw extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *LiquidityPoolWithdrawOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	pool, err := readObjectID(r)
	if err != nil {
		return err
	}
	shareAmount, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("liquidity pool withdraw extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Account = account
	o.Pool = pool
	o.ShareAmount = shareAmount
	o.Extensions = nil
	return nil
}

func (o *LiquidityPoolWithdrawOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o LiquidityPoolExchangeOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Pool.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AmountToSell.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.MinToReceive.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("liquidity pool exchange extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *LiquidityPoolExchangeOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	pool, err := readObjectID(r)
	if err != nil {
		return err
	}
	amountToSell, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	minToReceive, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("liquidity pool exchange extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Account = account
	o.Pool = pool
	o.AmountToSell = amountToSell
	o.MinToReceive = minToReceive
	o.Extensions = nil
	return nil
}

func (o *LiquidityPoolExchangeOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o LiquidityPoolUpdateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Pool.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := writeOptionalUint16(w, o.TakerFeePercent); err != nil {
		return nil, err
	}
	if err := writeOptionalUint16(w, o.WithdrawalFeePercent); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("liquidity pool update extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *LiquidityPoolUpdateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	pool, err := readObjectID(r)
	if err != nil {
		return err
	}
	takerFeePercent, err := readOptionalUint16(r)
	if err != nil {
		return err
	}
	withdrawalFeePercent, err := readOptionalUint16(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("liquidity pool update extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Account = account
	o.Pool = pool
	o.TakerFeePercent = takerFeePercent
	o.WithdrawalFeePercent = withdrawalFeePercent
	o.Extensions = nil
	return nil
}

func (o *LiquidityPoolUpdateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o SametFundCreateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.OwnerAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AssetType.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeInt64(o.Balance)
	w.writeUint32(o.FeeRate)
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("samet fund create extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *SametFundCreateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	ownerAccount, err := readObjectID(r)
	if err != nil {
		return err
	}
	assetType, err := readObjectID(r)
	if err != nil {
		return err
	}
	balance, err := r.readInt64()
	if err != nil {
		return err
	}
	feeRate, err := r.readUint32()
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("samet fund create extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.OwnerAccount = ownerAccount
	o.AssetType = assetType
	o.Balance = balance
	o.FeeRate = feeRate
	o.Extensions = nil
	return nil
}

func (o *SametFundCreateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o SametFundDeleteOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.OwnerAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.FundID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("samet fund delete extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *SametFundDeleteOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	ownerAccount, err := readObjectID(r)
	if err != nil {
		return err
	}
	fundID, err := readObjectID(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("samet fund delete extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.OwnerAccount = ownerAccount
	o.FundID = fundID
	o.Extensions = nil
	return nil
}

func (o *SametFundDeleteOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o SametFundUpdateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.OwnerAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.FundID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := writeOptionalAssetAmount(w, o.DeltaAmount); err != nil {
		return nil, err
	}
	if err := writeOptionalUint32(w, o.NewFeeRate); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("samet fund update extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *SametFundUpdateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	ownerAccount, err := readObjectID(r)
	if err != nil {
		return err
	}
	fundID, err := readObjectID(r)
	if err != nil {
		return err
	}
	deltaAmount, err := readOptionalAssetAmount(r)
	if err != nil {
		return err
	}
	newFeeRate, err := readOptionalUint32(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("samet fund update extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.OwnerAccount = ownerAccount
	o.FundID = fundID
	o.DeltaAmount = deltaAmount
	o.NewFeeRate = newFeeRate
	o.Extensions = nil
	return nil
}

func (o *SametFundUpdateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o SametFundBorrowOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Borrower.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.FundID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.BorrowAmount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("samet fund borrow extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *SametFundBorrowOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	borrower, err := readObjectID(r)
	if err != nil {
		return err
	}
	fundID, err := readObjectID(r)
	if err != nil {
		return err
	}
	borrowAmount, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("samet fund borrow extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Borrower = borrower
	o.FundID = fundID
	o.BorrowAmount = borrowAmount
	o.Extensions = nil
	return nil
}

func (o *SametFundBorrowOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o SametFundRepayOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.FundID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.RepayAmount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.FundFee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("samet fund repay extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *SametFundRepayOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	fundID, err := readObjectID(r)
	if err != nil {
		return err
	}
	repayAmount, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	fundFee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("samet fund repay extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Account = account
	o.FundID = fundID
	o.RepayAmount = repayAmount
	o.FundFee = fundFee
	o.Extensions = nil
	return nil
}

func (o *SametFundRepayOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o CreditOfferCreateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.OwnerAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AssetType.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeInt64(o.Balance)
	w.writeUint32(o.FeeRate)
	w.writeUint32(o.MaxDurationSeconds)
	w.writeInt64(o.MinDealAmount)
	writeBool(w, o.Enabled)
	if err := o.AutoDisableTime.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := writeCreditOfferCollateralMap(w, o.AcceptableCollateral); err != nil {
		return nil, err
	}
	if err := writeCreditOfferBorrowerMap(w, o.AcceptableBorrowers); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("credit offer create extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *CreditOfferCreateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	ownerAccount, err := readObjectID(r)
	if err != nil {
		return err
	}
	assetType, err := readObjectID(r)
	if err != nil {
		return err
	}
	balance, err := r.readInt64()
	if err != nil {
		return err
	}
	feeRate, err := r.readUint32()
	if err != nil {
		return err
	}
	maxDurationSeconds, err := r.readUint32()
	if err != nil {
		return err
	}
	minDealAmount, err := r.readInt64()
	if err != nil {
		return err
	}
	enabled, err := readBool(r)
	if err != nil {
		return err
	}
	autoDisableTime, err := readTime(r)
	if err != nil {
		return err
	}
	acceptableCollateral, err := readCreditOfferCollateralMap(r)
	if err != nil {
		return err
	}
	acceptableBorrowers, err := readCreditOfferBorrowerMap(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("credit offer create extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.OwnerAccount = ownerAccount
	o.AssetType = assetType
	o.Balance = balance
	o.FeeRate = feeRate
	o.MaxDurationSeconds = maxDurationSeconds
	o.MinDealAmount = minDealAmount
	o.Enabled = enabled
	o.AutoDisableTime = autoDisableTime
	o.AcceptableCollateral = acceptableCollateral
	o.AcceptableBorrowers = acceptableBorrowers
	o.Extensions = nil
	return nil
}

func (o *CreditOfferCreateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o CreditOfferDeleteOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.OwnerAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.OfferID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("credit offer delete extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *CreditOfferDeleteOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	ownerAccount, err := readObjectID(r)
	if err != nil {
		return err
	}
	offerID, err := readObjectID(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("credit offer delete extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.OwnerAccount = ownerAccount
	o.OfferID = offerID
	o.Extensions = nil
	return nil
}

func (o *CreditOfferDeleteOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o CreditOfferUpdateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.OwnerAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.OfferID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := writeOptionalAssetAmount(w, o.DeltaAmount); err != nil {
		return nil, err
	}
	if err := writeOptionalUint32(w, o.FeeRate); err != nil {
		return nil, err
	}
	if err := writeOptionalUint32(w, o.MaxDurationSeconds); err != nil {
		return nil, err
	}
	if err := writeOptionalInt64(w, o.MinDealAmount); err != nil {
		return nil, err
	}
	if err := writeOptionalBool(w, o.Enabled); err != nil {
		return nil, err
	}
	if err := writeOptionalTime(w, o.AutoDisableTime); err != nil {
		return nil, err
	}
	if err := writeOptionalCreditOfferCollateralMap(w, o.AcceptableCollateral); err != nil {
		return nil, err
	}
	if err := writeOptionalCreditOfferBorrowerMap(w, o.AcceptableBorrowers); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("credit offer update extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *CreditOfferUpdateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	ownerAccount, err := readObjectID(r)
	if err != nil {
		return err
	}
	offerID, err := readObjectID(r)
	if err != nil {
		return err
	}
	deltaAmount, err := readOptionalAssetAmount(r)
	if err != nil {
		return err
	}
	feeRate, err := readOptionalUint32(r)
	if err != nil {
		return err
	}
	maxDurationSeconds, err := readOptionalUint32(r)
	if err != nil {
		return err
	}
	minDealAmount, err := readOptionalInt64(r)
	if err != nil {
		return err
	}
	enabled, err := readOptionalBool(r)
	if err != nil {
		return err
	}
	autoDisableTime, err := readOptionalTime(r)
	if err != nil {
		return err
	}
	acceptableCollateral, err := readOptionalCreditOfferCollateralMap(r)
	if err != nil {
		return err
	}
	acceptableBorrowers, err := readOptionalCreditOfferBorrowerMap(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("credit offer update extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.OwnerAccount = ownerAccount
	o.OfferID = offerID
	o.DeltaAmount = deltaAmount
	o.FeeRate = feeRate
	o.MaxDurationSeconds = maxDurationSeconds
	o.MinDealAmount = minDealAmount
	o.Enabled = enabled
	o.AutoDisableTime = autoDisableTime
	o.AcceptableCollateral = acceptableCollateral
	o.AcceptableBorrowers = acceptableBorrowers
	o.Extensions = nil
	return nil
}

func (o *CreditOfferUpdateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o CreditOfferAcceptOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Borrower.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.OfferID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.BorrowAmount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Collateral.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeUint32(o.MaxFeeRate)
	w.writeUint32(o.MinDurationSeconds)
	if err := o.Extensions.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *CreditOfferAcceptOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	borrower, err := readObjectID(r)
	if err != nil {
		return err
	}
	offerID, err := readObjectID(r)
	if err != nil {
		return err
	}
	borrowAmount, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	collateral, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	maxFeeRate, err := r.readUint32()
	if err != nil {
		return err
	}
	minDurationSeconds, err := r.readUint32()
	if err != nil {
		return err
	}
	var ext CreditOfferAcceptExtensions
	if err := ext.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	o.Fee = fee
	o.Borrower = borrower
	o.OfferID = offerID
	o.BorrowAmount = borrowAmount
	o.Collateral = collateral
	o.MaxFeeRate = maxFeeRate
	o.MinDurationSeconds = minDurationSeconds
	o.Extensions = ext
	return nil
}

func (o *CreditOfferAcceptOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o CreditDealRepayOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.DealID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.RepayAmount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.CreditFee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("credit deal repay extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *CreditDealRepayOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	dealID, err := readObjectID(r)
	if err != nil {
		return err
	}
	repayAmount, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	creditFee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("credit deal repay extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Account = account
	o.DealID = dealID
	o.RepayAmount = repayAmount
	o.CreditFee = creditFee
	o.Extensions = nil
	return nil
}

func (o *CreditDealRepayOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o CreditDealExpiredOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.DealID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.OfferID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.OfferOwner.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Borrower.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.UnpaidAmount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Collateral.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeUint32(o.FeeRate)
	return w.Bytes(), nil
}

func (o *CreditDealExpiredOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	dealID, err := readObjectID(r)
	if err != nil {
		return err
	}
	offerID, err := readObjectID(r)
	if err != nil {
		return err
	}
	offerOwner, err := readObjectID(r)
	if err != nil {
		return err
	}
	borrower, err := readObjectID(r)
	if err != nil {
		return err
	}
	unpaidAmount, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	collateral, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	feeRate, err := r.readUint32()
	if err != nil {
		return err
	}
	o.Fee = fee
	o.DealID = dealID
	o.OfferID = offerID
	o.OfferOwner = offerOwner
	o.Borrower = borrower
	o.UnpaidAmount = unpaidAmount
	o.Collateral = collateral
	o.FeeRate = feeRate
	return nil
}

func (o *CreditDealExpiredOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o CreditDealUpdateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.DealID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeUint8(o.AutoRepay)
	if !extensionsEmpty(o.Extensions) {
		return nil, fmt.Errorf("credit deal update extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *CreditDealUpdateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	dealID, err := readObjectID(r)
	if err != nil {
		return err
	}
	autoRepay, err := r.readUint8()
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("credit deal update extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Account = account
	o.DealID = dealID
	o.AutoRepay = autoRepay
	o.Extensions = nil
	return nil
}

func (o *CreditDealUpdateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func writeOptionalInt64(w *binaryWriter, value *int64) error {
	if value == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	w.writeInt64(*value)
	return nil
}

func readOptionalInt64(r *binaryReader) (*int64, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	value, err := r.readInt64()
	if err != nil {
		return nil, err
	}
	return int64Ptr(value), nil
}

func writeCreditOfferCollateralMap(w *binaryWriter, values []CreditOfferCollateral) error {
	sorted := append([]CreditOfferCollateral(nil), values...)
	if len(sorted) > 1 {
		sortCreditOfferCollateral(sorted)
	}
	w.writeVarUint64(uint64(len(sorted)))
	for _, value := range sorted {
		if err := value.AssetID.MarshalBinaryInto(w); err != nil {
			return err
		}
		if err := value.Price.MarshalBinaryInto(w); err != nil {
			return err
		}
	}
	return nil
}

func readCreditOfferCollateralMap(r *binaryReader) ([]CreditOfferCollateral, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make([]CreditOfferCollateral, 0, count)
	for i := uint64(0); i < count; i++ {
		assetID, err := readObjectID(r)
		if err != nil {
			return nil, err
		}
		price, err := readPrice(r)
		if err != nil {
			return nil, err
		}
		out = append(out, CreditOfferCollateral{AssetID: assetID, Price: price})
	}
	return out, nil
}

func writeOptionalCreditOfferCollateralMap(w *binaryWriter, values []CreditOfferCollateral) error {
	if values == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	return writeCreditOfferCollateralMap(w, values)
}

func readOptionalCreditOfferCollateralMap(r *binaryReader) ([]CreditOfferCollateral, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	return readCreditOfferCollateralMap(r)
}

func writeCreditOfferBorrowerMap(w *binaryWriter, values []CreditOfferBorrower) error {
	sorted := append([]CreditOfferBorrower(nil), values...)
	if len(sorted) > 1 {
		sortCreditOfferBorrowers(sorted)
	}
	w.writeVarUint64(uint64(len(sorted)))
	for _, value := range sorted {
		if err := value.AccountID.MarshalBinaryInto(w); err != nil {
			return err
		}
		w.writeInt64(value.Amount)
	}
	return nil
}

func readCreditOfferBorrowerMap(r *binaryReader) ([]CreditOfferBorrower, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make([]CreditOfferBorrower, 0, count)
	for i := uint64(0); i < count; i++ {
		accountID, err := readObjectID(r)
		if err != nil {
			return nil, err
		}
		amount, err := r.readInt64()
		if err != nil {
			return nil, err
		}
		out = append(out, CreditOfferBorrower{AccountID: accountID, Amount: amount})
	}
	return out, nil
}

func writeOptionalCreditOfferBorrowerMap(w *binaryWriter, values []CreditOfferBorrower) error {
	if values == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	return writeCreditOfferBorrowerMap(w, values)
}

func readOptionalCreditOfferBorrowerMap(r *binaryReader) ([]CreditOfferBorrower, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	return readCreditOfferBorrowerMap(r)
}

func sortCreditOfferCollateral(values []CreditOfferCollateral) {
	for i := 1; i < len(values); i++ {
		j := i
		for j > 0 && compareObjectID(values[j-1].AssetID, values[j].AssetID) > 0 {
			values[j-1], values[j] = values[j], values[j-1]
			j--
		}
	}
}

func sortCreditOfferBorrowers(values []CreditOfferBorrower) {
	for i := 1; i < len(values); i++ {
		j := i
		for j > 0 && compareObjectID(values[j-1].AccountID, values[j].AccountID) > 0 {
			values[j-1], values[j] = values[j], values[j-1]
			j--
		}
	}
}

func customDataFromRaw(raw json.RawMessage) ([]byte, error) {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return nil, nil
	}
	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		if decoded, err := hex.DecodeString(strings.TrimSpace(asString)); err == nil {
			return decoded, nil
		}
		return []byte(asString), nil
	}
	var asBytes []byte
	if err := json.Unmarshal(raw, &asBytes); err == nil {
		return asBytes, nil
	}
	return nil, fmt.Errorf("custom operation data must be a JSON string or byte array")
}

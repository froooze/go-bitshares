package protocol

import (
	"encoding/json"
	"fmt"
)

func (o CallOrderUpdateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.FundingAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.DeltaCollateral.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.DeltaDebt.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Extensions.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *CallOrderUpdateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *CallOrderUpdateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	funding, err := readObjectID(r)
	if err != nil {
		return err
	}
	collateral, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	debt, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	var ext CallOrderUpdateExtensions
	if err := ext.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	o.Fee = fee
	o.FundingAccount = funding
	o.DeltaCollateral = collateral
	o.DeltaDebt = debt
	o.Extensions = ext
	return nil
}

func (o FillOrderOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.OrderID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AccountID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Pays.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Receives.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.FillPrice.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	writeBool(w, o.IsMaker)
	return w.Bytes(), nil
}

func (o *FillOrderOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *FillOrderOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	order, err := readObjectID(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	pays, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	receives, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	fillPrice, err := readPrice(r)
	if err != nil {
		return err
	}
	isMaker, err := readBool(r)
	if err != nil {
		return err
	}
	o.OrderID = order
	o.AccountID = account
	o.Pays = pays
	o.Receives = receives
	o.Fee = fee
	o.FillPrice = fillPrice
	o.IsMaker = isMaker
	return nil
}

func (o AccountCreateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Registrar.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Referrer.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeUint16(o.ReferrerPercent)
	w.writeString(o.Name)
	if err := o.Owner.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Active.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Options.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Extensions.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *AccountCreateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *AccountCreateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	registrar, err := readObjectID(r)
	if err != nil {
		return err
	}
	referrer, err := readObjectID(r)
	if err != nil {
		return err
	}
	referrerPercent, err := r.readUint16()
	if err != nil {
		return err
	}
	name, err := r.readString()
	if err != nil {
		return err
	}
	var owner Authority
	if err := owner.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	var active Authority
	if err := active.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	var options AccountOptions
	if err := options.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	var ext AccountCreateExtensions
	if err := ext.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	o.Fee = fee
	o.Registrar = registrar
	o.Referrer = referrer
	o.ReferrerPercent = referrerPercent
	o.Name = name
	o.Owner = owner
	o.Active = active
	o.Options = options
	o.Extensions = ext
	return nil
}

func (o AccountUpdateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := writeOptionalAuthority(w, o.Owner); err != nil {
		return nil, err
	}
	if err := writeOptionalAuthority(w, o.Active); err != nil {
		return nil, err
	}
	if err := writeOptionalAccountOptions(w, o.NewOptions); err != nil {
		return nil, err
	}
	if err := o.Extensions.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *AccountUpdateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *AccountUpdateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	owner, err := readOptionalAuthority(r)
	if err != nil {
		return err
	}
	active, err := readOptionalAuthority(r)
	if err != nil {
		return err
	}
	options, err := readOptionalAccountOptions(r)
	if err != nil {
		return err
	}
	var ext AccountUpdateExtensions
	if err := ext.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	o.Fee = fee
	o.Account = account
	o.Owner = owner
	o.Active = active
	o.NewOptions = options
	o.Extensions = ext
	return nil
}

func (o AccountWhitelistOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AuthorizingAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AccountToList.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeUint8(o.NewListing)
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *AccountWhitelistOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *AccountWhitelistOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	auth, err := readObjectID(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	listing, err := r.readUint8()
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("account whitelist extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.AuthorizingAccount = auth
	o.AccountToList = account
	o.NewListing = listing
	o.Extensions = nil
	return nil
}

func (o AccountUpgradeOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AccountToUpgrade.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	writeBool(w, o.UpgradeToLifetimeMember)
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *AccountUpgradeOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *AccountUpgradeOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	upgrade, err := readBool(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("account upgrade extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.AccountToUpgrade = account
	o.UpgradeToLifetimeMember = upgrade
	o.Extensions = nil
	return nil
}

func (o AccountTransferOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AccountID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.NewOwner.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *AccountTransferOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *AccountTransferOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	owner, err := readObjectID(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("account transfer extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.AccountID = account
	o.NewOwner = owner
	o.Extensions = nil
	return nil
}

func (o AssetCreateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Issuer.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeString(o.Symbol)
	w.writeUint8(o.Precision)
	if err := o.CommonOptions.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if o.BitassetOpts == nil {
		w.writeUint8(0)
	} else {
		w.writeUint8(1)
		if err := o.BitassetOpts.MarshalBinaryInto(w); err != nil {
			return nil, err
		}
	}
	writeBool(w, o.IsPredictionMarket)
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *AssetCreateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *AssetCreateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	issuer, err := readObjectID(r)
	if err != nil {
		return err
	}
	symbol, err := r.readString()
	if err != nil {
		return err
	}
	precision, err := r.readUint8()
	if err != nil {
		return err
	}
	var options AssetOptions
	if err := options.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	present, err := r.readUint8()
	if err != nil {
		return err
	}
	var bitasset *BitAssetOptions
	if present != 0 {
		var value BitAssetOptions
		if err := value.UnmarshalBinaryFrom(r); err != nil {
			return err
		}
		bitasset = &value
	}
	prediction, err := readBool(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("asset create extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Issuer = issuer
	o.Symbol = symbol
	o.Precision = precision
	o.CommonOptions = options
	o.BitassetOpts = bitasset
	o.IsPredictionMarket = prediction
	o.Extensions = nil
	return nil
}

func (o AssetUpdateOperation) MarshalBinary() ([]byte, error) {
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
	if err := writeOptionalObjectID(w, o.NewIssuer); err != nil {
		return nil, err
	}
	if err := o.NewOptions.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Extensions.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *AssetUpdateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *AssetUpdateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
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
	newIssuer, err := readOptionalObjectID(r)
	if err != nil {
		return err
	}
	var options AssetOptions
	if err := options.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	var ext AssetUpdateExtensions
	if err := ext.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	o.Fee = fee
	o.Issuer = issuer
	o.AssetToUpdate = asset
	o.NewIssuer = newIssuer
	o.NewOptions = options
	o.Extensions = ext
	return nil
}

func (o AssetUpdateBitassetOperation) MarshalBinary() ([]byte, error) {
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
	if err := o.NewOptions.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *AssetUpdateBitassetOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *AssetUpdateBitassetOperation) UnmarshalBinaryFrom(r *binaryReader) error {
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
	var options BitAssetOptions
	if err := options.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("asset update bitasset extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Issuer = issuer
	o.AssetToUpdate = asset
	o.NewOptions = options
	o.Extensions = nil
	return nil
}

func (o AssetUpdateFeedProducersOperation) MarshalBinary() ([]byte, error) {
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
	if err := writeObjectIDSet(w, o.NewFeedProducers); err != nil {
		return nil, err
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *AssetUpdateFeedProducersOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *AssetUpdateFeedProducersOperation) UnmarshalBinaryFrom(r *binaryReader) error {
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
	producers, err := readObjectIDSet(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("asset update feed producers extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Issuer = issuer
	o.AssetToUpdate = asset
	o.NewFeedProducers = producers
	o.Extensions = nil
	return nil
}

func (o AssetFundFeePoolOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.FromAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AssetID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeInt64(o.Amount)
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *AssetFundFeePoolOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *AssetFundFeePoolOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	from, err := readObjectID(r)
	if err != nil {
		return err
	}
	asset, err := readObjectID(r)
	if err != nil {
		return err
	}
	amount, err := r.readInt64()
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("asset fund fee pool extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.FromAccount = from
	o.AssetID = asset
	o.Amount = amount
	o.Extensions = nil
	return nil
}

func (o AssetSettleOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Amount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *AssetSettleOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *AssetSettleOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
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
		return fmt.Errorf("asset settle extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Account = account
	o.Amount = amount
	o.Extensions = nil
	return nil
}

func (o AssetGlobalSettleOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Issuer.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AssetToSettle.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.SettlePrice.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *AssetGlobalSettleOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *AssetGlobalSettleOperation) UnmarshalBinaryFrom(r *binaryReader) error {
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
	price, err := readPrice(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("asset global settle extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Issuer = issuer
	o.AssetToSettle = asset
	o.SettlePrice = price
	o.Extensions = nil
	return nil
}

func (o AssetPublishFeedOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Publisher.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AssetID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Feed.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Extensions.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *AssetPublishFeedOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *AssetPublishFeedOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	publisher, err := readObjectID(r)
	if err != nil {
		return err
	}
	asset, err := readObjectID(r)
	if err != nil {
		return err
	}
	var feed PriceFeed
	if err := feed.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	var ext AssetPublishFeedExtensions
	if err := ext.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	o.Fee = fee
	o.Publisher = publisher
	o.AssetID = asset
	o.Feed = feed
	o.Extensions = ext
	return nil
}

func (o WitnessCreateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.WitnessAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeString(o.URL)
	if err := o.BlockSigningKey.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *WitnessCreateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *WitnessCreateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	url, err := r.readString()
	if err != nil {
		return err
	}
	key, err := readPublicKey(r)
	if err != nil {
		return err
	}
	o.Fee = fee
	o.WitnessAccount = account
	o.URL = url
	o.BlockSigningKey = key
	return nil
}

func (o WitnessUpdateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Witness.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.WitnessAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := writeOptionalString(w, o.NewURL); err != nil {
		return nil, err
	}
	if err := writeOptionalPublicKey(w, o.NewSigningKey); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *WitnessUpdateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *WitnessUpdateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	witness, err := readObjectID(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	url, err := readOptionalString(r)
	if err != nil {
		return err
	}
	key, err := readOptionalPublicKey(r)
	if err != nil {
		return err
	}
	o.Fee = fee
	o.Witness = witness
	o.WitnessAccount = account
	o.NewURL = url
	o.NewSigningKey = key
	return nil
}

func (o OpWrapper) MarshalBinary() ([]byte, error) {
	return o.Op.MarshalBinary()
}

func (o *OpWrapper) UnmarshalBinary(data []byte) error {
	var env OperationEnvelope
	if err := env.UnmarshalBinary(data); err != nil {
		return err
	}
	o.Op = env
	return nil
}

func (o ProposalCreateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.FeePayingAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.ExpirationTime.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeVarUint64(uint64(len(o.ProposedOps)))
	for _, op := range o.ProposedOps {
		raw, err := op.MarshalBinary()
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(raw); err != nil {
			return nil, err
		}
	}
	if err := writeOptionalUint32(w, o.ReviewPeriodSeconds); err != nil {
		return nil, err
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *ProposalCreateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *ProposalCreateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	expiration, err := readTime(r)
	if err != nil {
		return err
	}
	count, err := r.readVarUint64()
	if err != nil {
		return err
	}
	ops := make([]OpWrapper, 0, count)
	for i := uint64(0); i < count; i++ {
		env, err := readOperationEnvelope(r)
		if err != nil {
			return err
		}
		ops = append(ops, OpWrapper{Op: env})
	}
	review, err := readOptionalUint32(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("proposal create extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.FeePayingAccount = account
	o.ExpirationTime = expiration
	o.ProposedOps = ops
	o.ReviewPeriodSeconds = review
	o.Extensions = nil
	return nil
}

func (o ProposalUpdateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.FeePayingAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Proposal.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := writeObjectIDSet(w, o.ActiveApprovalsToAdd); err != nil {
		return nil, err
	}
	if err := writeObjectIDSet(w, o.ActiveApprovalsToRemove); err != nil {
		return nil, err
	}
	if err := writeObjectIDSet(w, o.OwnerApprovalsToAdd); err != nil {
		return nil, err
	}
	if err := writeObjectIDSet(w, o.OwnerApprovalsToRemove); err != nil {
		return nil, err
	}
	if err := writePublicKeySet(w, o.KeyApprovalsToAdd); err != nil {
		return nil, err
	}
	if err := writePublicKeySet(w, o.KeyApprovalsToRemove); err != nil {
		return nil, err
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *ProposalUpdateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *ProposalUpdateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	proposal, err := readObjectID(r)
	if err != nil {
		return err
	}
	activeAdd, err := readObjectIDSet(r)
	if err != nil {
		return err
	}
	activeRemove, err := readObjectIDSet(r)
	if err != nil {
		return err
	}
	ownerAdd, err := readObjectIDSet(r)
	if err != nil {
		return err
	}
	ownerRemove, err := readObjectIDSet(r)
	if err != nil {
		return err
	}
	keyAdd, err := readPublicKeySet(r)
	if err != nil {
		return err
	}
	keyRemove, err := readPublicKeySet(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("proposal update extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.FeePayingAccount = account
	o.Proposal = proposal
	o.ActiveApprovalsToAdd = activeAdd
	o.ActiveApprovalsToRemove = activeRemove
	o.OwnerApprovalsToAdd = ownerAdd
	o.OwnerApprovalsToRemove = ownerRemove
	o.KeyApprovalsToAdd = keyAdd
	o.KeyApprovalsToRemove = keyRemove
	o.Extensions = nil
	return nil
}

func (o ProposalDeleteOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.FeePayingAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	writeBool(w, o.UsingOwnerAuthority)
	if err := o.Proposal.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *ProposalDeleteOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *ProposalDeleteOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	owner, err := readBool(r)
	if err != nil {
		return err
	}
	proposal, err := readObjectID(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("proposal delete extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.FeePayingAccount = account
	o.UsingOwnerAuthority = owner
	o.Proposal = proposal
	o.Extensions = nil
	return nil
}

func (o WithdrawPermissionCreateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.WithdrawFromAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AuthorizedAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.WithdrawalLimit.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeUint32(o.WithdrawalPeriodSec)
	w.writeUint32(o.PeriodsUntilExpiration)
	if err := o.PeriodStartTime.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *WithdrawPermissionCreateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *WithdrawPermissionCreateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	from, err := readObjectID(r)
	if err != nil {
		return err
	}
	auth, err := readObjectID(r)
	if err != nil {
		return err
	}
	limit, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	period, err := r.readUint32()
	if err != nil {
		return err
	}
	expire, err := r.readUint32()
	if err != nil {
		return err
	}
	start, err := readTime(r)
	if err != nil {
		return err
	}
	o.Fee = fee
	o.WithdrawFromAccount = from
	o.AuthorizedAccount = auth
	o.WithdrawalLimit = limit
	o.WithdrawalPeriodSec = period
	o.PeriodsUntilExpiration = expire
	o.PeriodStartTime = start
	return nil
}

func (o WithdrawPermissionUpdateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.WithdrawFromAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AuthorizedAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.PermissionToUpdate.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.WithdrawalLimit.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeUint32(o.WithdrawalPeriodSec)
	if err := o.PeriodStartTime.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeUint32(o.PeriodsUntilExpiration)
	return w.Bytes(), nil
}

func (o *WithdrawPermissionUpdateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *WithdrawPermissionUpdateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	from, err := readObjectID(r)
	if err != nil {
		return err
	}
	auth, err := readObjectID(r)
	if err != nil {
		return err
	}
	perm, err := readObjectID(r)
	if err != nil {
		return err
	}
	limit, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	period, err := r.readUint32()
	if err != nil {
		return err
	}
	start, err := readTime(r)
	if err != nil {
		return err
	}
	expire, err := r.readUint32()
	if err != nil {
		return err
	}
	o.Fee = fee
	o.WithdrawFromAccount = from
	o.AuthorizedAccount = auth
	o.PermissionToUpdate = perm
	o.WithdrawalLimit = limit
	o.WithdrawalPeriodSec = period
	o.PeriodStartTime = start
	o.PeriodsUntilExpiration = expire
	return nil
}

func (o WithdrawPermissionClaimOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.WithdrawPermission.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.WithdrawFromAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.WithdrawToAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AmountToWithdraw.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := writeOptionalMemo(w, o.Memo); err != nil {
		return nil, err
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *WithdrawPermissionClaimOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *WithdrawPermissionClaimOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	perm, err := readObjectID(r)
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
		return fmt.Errorf("withdraw permission claim extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.WithdrawPermission = perm
	o.WithdrawFromAccount = from
	o.WithdrawToAccount = to
	o.AmountToWithdraw = amount
	if memo != nil {
		raw, err := json.Marshal(memo)
		if err != nil {
			return err
		}
		o.Memo = raw
	}
	return nil
}

func (o WithdrawPermissionDeleteOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.WithdrawFromAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AuthorizedAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.WithdrawalPermission.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *WithdrawPermissionDeleteOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *WithdrawPermissionDeleteOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	from, err := readObjectID(r)
	if err != nil {
		return err
	}
	auth, err := readObjectID(r)
	if err != nil {
		return err
	}
	perm, err := readObjectID(r)
	if err != nil {
		return err
	}
	o.Fee = fee
	o.WithdrawFromAccount = from
	o.AuthorizedAccount = auth
	o.WithdrawalPermission = perm
	return nil
}

func (o CommitteeMemberCreateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.CommitteeMemberAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeString(o.URL)
	return w.Bytes(), nil
}

func (o *CommitteeMemberCreateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *CommitteeMemberCreateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	url, err := r.readString()
	if err != nil {
		return err
	}
	o.Fee = fee
	o.CommitteeMemberAccount = account
	o.URL = url
	return nil
}

func (o CommitteeMemberUpdateOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.CommitteeMember.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.CommitteeMemberAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := writeOptionalString(w, o.NewURL); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *CommitteeMemberUpdateOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *CommitteeMemberUpdateOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	member, err := readObjectID(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	url, err := readOptionalString(r)
	if err != nil {
		return err
	}
	o.Fee = fee
	o.CommitteeMember = member
	o.CommitteeMemberAccount = account
	o.NewURL = url
	return nil
}

func (o CommitteeMemberUpdateGlobalParametersOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	raw, err := o.NewParameters.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(raw); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *CommitteeMemberUpdateGlobalParametersOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *CommitteeMemberUpdateGlobalParametersOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	var params ChainParameters
	if err := params.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	o.Fee = fee
	o.NewParameters = params
	return nil
}

func (o BalanceClaimOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.DepositToAccount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.BalanceToClaim.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.BalanceOwnerKey.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.TotalClaimed.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *BalanceClaimOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *BalanceClaimOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	deposit, err := readObjectID(r)
	if err != nil {
		return err
	}
	claim, err := readObjectID(r)
	if err != nil {
		return err
	}
	key, err := readPublicKey(r)
	if err != nil {
		return err
	}
	total, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	o.Fee = fee
	o.DepositToAccount = deposit
	o.BalanceToClaim = claim
	o.BalanceOwnerKey = key
	o.TotalClaimed = total
	return nil
}

func (o OverrideTransferOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Issuer.MarshalBinaryInto(w); err != nil {
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
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *OverrideTransferOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *OverrideTransferOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	issuer, err := readObjectID(r)
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
		return fmt.Errorf("override transfer extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Issuer = issuer
	o.From = from
	o.To = to
	o.Amount = amount
	if memo != nil {
		raw, err := json.Marshal(memo)
		if err != nil {
			return err
		}
		o.Memo = raw
	}
	o.Extensions = nil
	return nil
}

func (o AssetSettleCancelOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Settlement.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Account.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Amount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *AssetSettleCancelOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *AssetSettleCancelOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	settlement, err := readObjectID(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
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
		return fmt.Errorf("asset settle cancel extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Settlement = settlement
	o.Account = account
	o.Amount = amount
	o.Extensions = nil
	return nil
}

func (o AssetClaimFeesOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Issuer.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AmountToClaim.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if o.Extensions == nil {
		w.writeVarUint64(0)
		return w.Bytes(), nil
	}
	if err := o.Extensions.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *AssetClaimFeesOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *AssetClaimFeesOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	issuer, err := readObjectID(r)
	if err != nil {
		return err
	}
	amount, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	var ext AssetClaimFeesExtensions
	if err := ext.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	o.Fee = fee
	o.Issuer = issuer
	o.AmountToClaim = amount
	if ext.ClaimFromAssetID != nil {
		o.Extensions = &ext
	} else {
		o.Extensions = nil
	}
	return nil
}

func (o FBADistributeOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AccountID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.FBAID.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeInt64(o.Amount)
	return w.Bytes(), nil
}

func (o *FBADistributeOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *FBADistributeOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	account, err := readObjectID(r)
	if err != nil {
		return err
	}
	fba, err := readObjectID(r)
	if err != nil {
		return err
	}
	amount, err := r.readInt64()
	if err != nil {
		return err
	}
	o.Fee = fee
	o.AccountID = account
	o.FBAID = fba
	o.Amount = amount
	return nil
}

func (o BidCollateralOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Bidder.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.AdditionalCollateral.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.DebtCovered.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

func (o *BidCollateralOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *BidCollateralOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	bidder, err := readObjectID(r)
	if err != nil {
		return err
	}
	collateral, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	debt, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("bid collateral extensions are not supported in binary serialization")
	}
	o.Fee = fee
	o.Bidder = bidder
	o.AdditionalCollateral = collateral
	o.DebtCovered = debt
	o.Extensions = nil
	return nil
}

func (o ExecuteBidOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Bidder.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Debt.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Collateral.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (o *ExecuteBidOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *ExecuteBidOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	bidder, err := readObjectID(r)
	if err != nil {
		return err
	}
	debt, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	collateral, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	o.Fee = fee
	o.Bidder = bidder
	o.Debt = debt
	o.Collateral = collateral
	return nil
}

func writeOptionalUint16(w *binaryWriter, value *uint16) error {
	if value == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	w.writeUint16(*value)
	return nil
}

func readOptionalUint16(r *binaryReader) (*uint16, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	value, err := r.readUint16()
	if err != nil {
		return nil, err
	}
	return uint16Ptr(value), nil
}

func writeOptionalUint32(w *binaryWriter, value *uint32) error {
	if value == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	w.writeUint32(*value)
	return nil
}

func readOptionalUint32(r *binaryReader) (*uint32, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	value, err := r.readUint32()
	if err != nil {
		return nil, err
	}
	return uint32Ptr(value), nil
}

func writeOptionalString(w *binaryWriter, value *string) error {
	if value == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	w.writeString(*value)
	return nil
}

func readOptionalString(r *binaryReader) (*string, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	value, err := r.readString()
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func writeOptionalAuthority(w *binaryWriter, value *Authority) error {
	if value == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	return value.MarshalBinaryInto(w)
}

func readOptionalAuthority(r *binaryReader) (*Authority, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	var value Authority
	if err := value.UnmarshalBinaryFrom(r); err != nil {
		return nil, err
	}
	return &value, nil
}

func writeOptionalAccountOptions(w *binaryWriter, value *AccountOptions) error {
	if value == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	return value.MarshalBinaryInto(w)
}

func readOptionalAccountOptions(r *binaryReader) (*AccountOptions, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	var value AccountOptions
	if err := value.UnmarshalBinaryFrom(r); err != nil {
		return nil, err
	}
	return &value, nil
}

func writeOptionalBitAssetOptions(w *binaryWriter, value *BitAssetOptions) error {
	if value == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	return value.MarshalBinaryInto(w)
}

func readOptionalBitAssetOptions(r *binaryReader) (*BitAssetOptions, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	var value BitAssetOptions
	if err := value.UnmarshalBinaryFrom(r); err != nil {
		return nil, err
	}
	return &value, nil
}

package protocol

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// HexBytes encodes byte slices as hex strings in JSON while preserving raw bytes on the wire.
type HexBytes []byte

func (b HexBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(b))
}

func (b *HexBytes) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || strings.EqualFold(strings.TrimSpace(string(data)), "null") {
		*b = nil
		return nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	raw, err := hex.DecodeString(strings.TrimSpace(value))
	if err != nil {
		return err
	}
	*b = append(HexBytes(nil), raw...)
	return nil
}

func (b HexBytes) MarshalBinaryInto(w *binaryWriter) error {
	writeVarBytes(w, b)
	return nil
}

func (b *HexBytes) UnmarshalBinaryFrom(r *binaryReader) error {
	raw, err := readVarBytes(r)
	if err != nil {
		return err
	}
	*b = append(HexBytes(nil), raw...)
	return nil
}

// NoSpecialAuthority mirrors the empty core special authority variant.
type NoSpecialAuthority struct{}

// TopHoldersSpecialAuthority selects the top holders of an asset.
type TopHoldersSpecialAuthority struct {
	Asset         ObjectID `json:"asset"`
	NumTopHolders uint8    `json:"num_top_holders"`
}

// SpecialAuthority is the static variant used by account extensions.
type SpecialAuthority struct {
	Kind       uint16
	No         *NoSpecialAuthority
	TopHolders *TopHoldersSpecialAuthority
}

func (a SpecialAuthority) MarshalJSON() ([]byte, error) {
	switch a.Kind {
	case 0:
		return json.Marshal([]any{uint16(0), struct{}{}})
	case 1:
		if a.TopHolders == nil {
			return nil, fmt.Errorf("missing top holders special authority")
		}
		return json.Marshal([]any{uint16(1), a.TopHolders})
	default:
		return nil, fmt.Errorf("unsupported special authority type %d", a.Kind)
	}
}

func (a *SpecialAuthority) UnmarshalJSON(data []byte) error {
	var body []json.RawMessage
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	if len(body) != 2 {
		return fmt.Errorf("invalid special authority")
	}
	var kind uint16
	if err := json.Unmarshal(body[0], &kind); err != nil {
		return err
	}
	a.Kind = kind
	switch kind {
	case 0:
		a.No = &NoSpecialAuthority{}
		a.TopHolders = nil
	case 1:
		var value TopHoldersSpecialAuthority
		if err := json.Unmarshal(body[1], &value); err != nil {
			return err
		}
		a.No = nil
		a.TopHolders = &value
	default:
		return fmt.Errorf("unsupported special authority type %d", kind)
	}
	return nil
}

func (a SpecialAuthority) MarshalBinaryInto(w *binaryWriter) error {
	w.writeVarUint64(uint64(a.Kind))
	switch a.Kind {
	case 0:
		return nil
	case 1:
		if a.TopHolders == nil {
			return fmt.Errorf("missing top holders special authority")
		}
		if err := a.TopHolders.Asset.MarshalBinaryInto(w); err != nil {
			return err
		}
		w.writeUint8(a.TopHolders.NumTopHolders)
		return nil
	default:
		return fmt.Errorf("unsupported special authority type %d", a.Kind)
	}
}

func (a *SpecialAuthority) UnmarshalBinaryFrom(r *binaryReader) error {
	kind, err := r.readVarUint64()
	if err != nil {
		return err
	}
	a.Kind = uint16(kind)
	switch a.Kind {
	case 0:
		a.No = &NoSpecialAuthority{}
		a.TopHolders = nil
	case 1:
		asset, err := readObjectID(r)
		if err != nil {
			return err
		}
		top, err := r.readUint8()
		if err != nil {
			return err
		}
		a.No = nil
		a.TopHolders = &TopHoldersSpecialAuthority{Asset: asset, NumTopHolders: top}
	default:
		return fmt.Errorf("unsupported special authority type %d", a.Kind)
	}
	return nil
}

// BuybackAccountOptions mirrors the account buyback extension.
type BuybackAccountOptions struct {
	AssetToBuy       ObjectID   `json:"asset_to_buy"`
	AssetToBuyIssuer ObjectID   `json:"asset_to_buy_issuer"`
	Markets          []ObjectID `json:"markets"`
}

func (o BuybackAccountOptions) MarshalBinaryInto(w *binaryWriter) error {
	if err := o.AssetToBuy.MarshalBinaryInto(w); err != nil {
		return err
	}
	if err := o.AssetToBuyIssuer.MarshalBinaryInto(w); err != nil {
		return err
	}
	return writeObjectIDSet(w, o.Markets)
}

func (o *BuybackAccountOptions) UnmarshalBinaryFrom(r *binaryReader) error {
	asset, err := readObjectID(r)
	if err != nil {
		return err
	}
	issuer, err := readObjectID(r)
	if err != nil {
		return err
	}
	markets, err := readObjectIDSet(r)
	if err != nil {
		return err
	}
	o.AssetToBuy = asset
	o.AssetToBuyIssuer = issuer
	o.Markets = markets
	return nil
}

// AccountCreateExtensions models the BitShares account_create extension block.
type AccountCreateExtensions struct {
	OwnerSpecialAuthority  *SpecialAuthority      `json:"owner_special_authority,omitempty"`
	ActiveSpecialAuthority *SpecialAuthority      `json:"active_special_authority,omitempty"`
	BuybackOptions         *BuybackAccountOptions `json:"buyback_options,omitempty"`
}

func (e AccountCreateExtensions) MarshalBinaryInto(w *binaryWriter) error {
	count := 0
	if e.OwnerSpecialAuthority != nil {
		count++
	}
	if e.ActiveSpecialAuthority != nil {
		count++
	}
	if e.BuybackOptions != nil {
		count++
	}
	w.writeVarUint64(uint64(count))
	if e.OwnerSpecialAuthority != nil {
		w.writeVarUint64(1)
		if err := e.OwnerSpecialAuthority.MarshalBinaryInto(w); err != nil {
			return err
		}
	}
	if e.ActiveSpecialAuthority != nil {
		w.writeVarUint64(2)
		if err := e.ActiveSpecialAuthority.MarshalBinaryInto(w); err != nil {
			return err
		}
	}
	if e.BuybackOptions != nil {
		w.writeVarUint64(3)
		if err := e.BuybackOptions.MarshalBinaryInto(w); err != nil {
			return err
		}
	}
	return nil
}

func (e *AccountCreateExtensions) UnmarshalBinaryFrom(r *binaryReader) error {
	count, err := r.readVarUint64()
	if err != nil {
		return err
	}
	for i := uint64(0); i < count; i++ {
		index, err := r.readVarUint64()
		if err != nil {
			return err
		}
		switch index {
		case 0:
			// null_ext placeholder, no payload
		case 1:
			var value SpecialAuthority
			if err := value.UnmarshalBinaryFrom(r); err != nil {
				return err
			}
			e.OwnerSpecialAuthority = &value
		case 2:
			var value SpecialAuthority
			if err := value.UnmarshalBinaryFrom(r); err != nil {
				return err
			}
			e.ActiveSpecialAuthority = &value
		case 3:
			var value BuybackAccountOptions
			if err := value.UnmarshalBinaryFrom(r); err != nil {
				return err
			}
			e.BuybackOptions = &value
		default:
			return fmt.Errorf("unknown account create extension %d", index)
		}
	}
	return nil
}

// AccountUpdateExtensions models the BitShares account_update extension block.
type AccountUpdateExtensions struct {
	OwnerSpecialAuthority  *SpecialAuthority `json:"owner_special_authority,omitempty"`
	ActiveSpecialAuthority *SpecialAuthority `json:"active_special_authority,omitempty"`
}

func (e AccountUpdateExtensions) MarshalBinaryInto(w *binaryWriter) error {
	count := 0
	if e.OwnerSpecialAuthority != nil {
		count++
	}
	if e.ActiveSpecialAuthority != nil {
		count++
	}
	w.writeVarUint64(uint64(count))
	if e.OwnerSpecialAuthority != nil {
		w.writeVarUint64(1)
		if err := e.OwnerSpecialAuthority.MarshalBinaryInto(w); err != nil {
			return err
		}
	}
	if e.ActiveSpecialAuthority != nil {
		w.writeVarUint64(2)
		if err := e.ActiveSpecialAuthority.MarshalBinaryInto(w); err != nil {
			return err
		}
	}
	return nil
}

func (e *AccountUpdateExtensions) UnmarshalBinaryFrom(r *binaryReader) error {
	count, err := r.readVarUint64()
	if err != nil {
		return err
	}
	for i := uint64(0); i < count; i++ {
		index, err := r.readVarUint64()
		if err != nil {
			return err
		}
		switch index {
		case 0:
			// null_ext placeholder, no payload
		case 1:
			var value SpecialAuthority
			if err := value.UnmarshalBinaryFrom(r); err != nil {
				return err
			}
			e.OwnerSpecialAuthority = &value
		case 2:
			var value SpecialAuthority
			if err := value.UnmarshalBinaryFrom(r); err != nil {
				return err
			}
			e.ActiveSpecialAuthority = &value
		default:
			return fmt.Errorf("unknown account update extension %d", index)
		}
	}
	return nil
}

// AssetUpdateExtensions models the BitShares asset_update extension block.
type AssetUpdateExtensions struct {
	NewPrecision         *uint8 `json:"new_precision,omitempty"`
	SkipCoreExchangeRate *bool  `json:"skip_core_exchange_rate,omitempty"`
}

func (e AssetUpdateExtensions) MarshalBinaryInto(w *binaryWriter) error {
	count := 0
	if e.NewPrecision != nil {
		count++
	}
	if e.SkipCoreExchangeRate != nil {
		count++
	}
	w.writeVarUint64(uint64(count))
	if e.NewPrecision != nil {
		w.writeVarUint64(0)
		w.writeUint8(*e.NewPrecision)
	}
	if e.SkipCoreExchangeRate != nil {
		w.writeVarUint64(1)
		w.writeUint8(boolByte(*e.SkipCoreExchangeRate))
	}
	return nil
}

func (e *AssetUpdateExtensions) UnmarshalBinaryFrom(r *binaryReader) error {
	count, err := r.readVarUint64()
	if err != nil {
		return err
	}
	for i := uint64(0); i < count; i++ {
		index, err := r.readVarUint64()
		if err != nil {
			return err
		}
		switch index {
		case 0:
			value, err := r.readUint8()
			if err != nil {
				return err
			}
			e.NewPrecision = uint8Ptr(value)
		case 1:
			value, err := r.readUint8()
			if err != nil {
				return err
			}
			flag := value != 0
			e.SkipCoreExchangeRate = &flag
		default:
			return fmt.Errorf("unknown asset update extension %d", index)
		}
	}
	return nil
}

// CallOrderUpdateExtensions mirrors call_order_update_operation::options_type.
type CallOrderUpdateExtensions struct {
	TargetCollateralRatio *uint16 `json:"target_collateral_ratio,omitempty"`
}

func (e CallOrderUpdateExtensions) MarshalBinaryInto(w *binaryWriter) error {
	if e.TargetCollateralRatio == nil {
		w.writeVarUint64(0)
		return nil
	}
	w.writeVarUint64(1)
	w.writeVarUint64(0)
	w.writeUint16(*e.TargetCollateralRatio)
	return nil
}

func (e *CallOrderUpdateExtensions) UnmarshalBinaryFrom(r *binaryReader) error {
	count, err := r.readVarUint64()
	if err != nil {
		return err
	}
	for i := uint64(0); i < count; i++ {
		index, err := r.readVarUint64()
		if err != nil {
			return err
		}
		switch index {
		case 0:
			value, err := r.readUint16()
			if err != nil {
				return err
			}
			e.TargetCollateralRatio = uint16Ptr(value)
		default:
			return fmt.Errorf("unknown call order update extension %d", index)
		}
	}
	return nil
}

// AssetPublishFeedExtensions mirrors asset_publish_feed_operation::ext.
type AssetPublishFeedExtensions struct {
	InitialCollateralRatio *uint16 `json:"initial_collateral_ratio,omitempty"`
}

func (e AssetPublishFeedExtensions) MarshalBinaryInto(w *binaryWriter) error {
	if e.InitialCollateralRatio == nil {
		w.writeVarUint64(0)
		return nil
	}
	w.writeVarUint64(1)
	w.writeVarUint64(0)
	w.writeUint16(*e.InitialCollateralRatio)
	return nil
}

func (e *AssetPublishFeedExtensions) UnmarshalBinaryFrom(r *binaryReader) error {
	count, err := r.readVarUint64()
	if err != nil {
		return err
	}
	for i := uint64(0); i < count; i++ {
		index, err := r.readVarUint64()
		if err != nil {
			return err
		}
		switch index {
		case 0:
			value, err := r.readUint16()
			if err != nil {
				return err
			}
			e.InitialCollateralRatio = uint16Ptr(value)
		default:
			return fmt.Errorf("unknown asset publish feed extension %d", index)
		}
	}
	return nil
}

// HTLCCreateExtensions mirrors htlc_create_operation::additional_options_type.
type HTLCCreateExtensions struct {
	Memo json.RawMessage `json:"memo,omitempty"`
}

func (e HTLCCreateExtensions) MarshalBinaryInto(w *binaryWriter) error {
	memo, err := memoDataFromRaw(e.Memo)
	if err != nil {
		return err
	}
	if memo == nil {
		w.writeVarUint64(0)
		return nil
	}
	w.writeVarUint64(1)
	w.writeVarUint64(0)
	return memo.writeBinary(w)
}

func (e *HTLCCreateExtensions) UnmarshalBinaryFrom(r *binaryReader) error {
	count, err := r.readVarUint64()
	if err != nil {
		return err
	}
	for i := uint64(0); i < count; i++ {
		index, err := r.readVarUint64()
		if err != nil {
			return err
		}
		switch index {
		case 0:
			var memo MemoData
			if err := memo.readBinary(r); err != nil {
				return err
			}
			raw, err := json.Marshal(memo)
			if err != nil {
				return err
			}
			e.Memo = raw
		default:
			return fmt.Errorf("unknown htlc create extension %d", index)
		}
	}
	return nil
}

// CreditOfferAcceptExtensions mirrors credit_offer_accept_operation::ext.
type CreditOfferAcceptExtensions struct {
	AutoRepay *uint8 `json:"auto_repay,omitempty"`
}

func (e CreditOfferAcceptExtensions) MarshalBinaryInto(w *binaryWriter) error {
	if e.AutoRepay == nil {
		w.writeVarUint64(0)
		return nil
	}
	w.writeVarUint64(1)
	w.writeVarUint64(0)
	w.writeUint8(*e.AutoRepay)
	return nil
}

func (e *CreditOfferAcceptExtensions) UnmarshalBinaryFrom(r *binaryReader) error {
	count, err := r.readVarUint64()
	if err != nil {
		return err
	}
	for i := uint64(0); i < count; i++ {
		index, err := r.readVarUint64()
		if err != nil {
			return err
		}
		switch index {
		case 0:
			value, err := r.readUint8()
			if err != nil {
				return err
			}
			e.AutoRepay = uint8Ptr(value)
		default:
			return fmt.Errorf("unknown credit offer accept extension %d", index)
		}
	}
	return nil
}

// HTLCOptions models the HTLC chain parameter extension.
type HTLCOptions struct {
	MaxTimeoutSecs  uint32 `json:"max_timeout_secs"`
	MaxPreimageSize uint32 `json:"max_preimage_size"`
}

func (o HTLCOptions) MarshalBinaryInto(w *binaryWriter) error {
	w.writeUint32(o.MaxTimeoutSecs)
	w.writeUint32(o.MaxPreimageSize)
	return nil
}

func (o *HTLCOptions) UnmarshalBinaryFrom(r *binaryReader) error {
	timeout, err := r.readUint32()
	if err != nil {
		return err
	}
	size, err := r.readUint32()
	if err != nil {
		return err
	}
	o.MaxTimeoutSecs = timeout
	o.MaxPreimageSize = size
	return nil
}

// CustomAuthorityOptions models the custom authority chain parameter extension.
type CustomAuthorityOptions struct {
	MaxCustomAuthorityLifetimeSeconds uint32 `json:"max_custom_authority_lifetime_seconds"`
	MaxCustomAuthoritiesPerAccount    uint32 `json:"max_custom_authorities_per_account"`
	MaxCustomAuthoritiesPerAccountOp  uint32 `json:"max_custom_authorities_per_account_op"`
	MaxCustomAuthorityRestrictions    uint32 `json:"max_custom_authority_restrictions"`
}

func (o CustomAuthorityOptions) MarshalBinaryInto(w *binaryWriter) error {
	w.writeUint32(o.MaxCustomAuthorityLifetimeSeconds)
	w.writeUint32(o.MaxCustomAuthoritiesPerAccount)
	w.writeUint32(o.MaxCustomAuthoritiesPerAccountOp)
	w.writeUint32(o.MaxCustomAuthorityRestrictions)
	return nil
}

func (o *CustomAuthorityOptions) UnmarshalBinaryFrom(r *binaryReader) error {
	lifetime, err := r.readUint32()
	if err != nil {
		return err
	}
	perAccount, err := r.readUint32()
	if err != nil {
		return err
	}
	perOp, err := r.readUint32()
	if err != nil {
		return err
	}
	restrictions, err := r.readUint32()
	if err != nil {
		return err
	}
	o.MaxCustomAuthorityLifetimeSeconds = lifetime
	o.MaxCustomAuthoritiesPerAccount = perAccount
	o.MaxCustomAuthoritiesPerAccountOp = perOp
	o.MaxCustomAuthorityRestrictions = restrictions
	return nil
}

// ChainParametersExtensions models the committee global parameter extension block.
type ChainParametersExtensions struct {
	UpdatableHTLCOptions    *HTLCOptions            `json:"updatable_htlc_options,omitempty"`
	CustomAuthorityOptions  *CustomAuthorityOptions `json:"custom_authority_options,omitempty"`
	MarketFeeNetworkPercent *uint16                 `json:"market_fee_network_percent,omitempty"`
	MakerFeeDiscountPercent *uint16                 `json:"maker_fee_discount_percent,omitempty"`
}

func (e ChainParametersExtensions) MarshalBinaryInto(w *binaryWriter) error {
	count := 0
	if e.UpdatableHTLCOptions != nil {
		count++
	}
	if e.CustomAuthorityOptions != nil {
		count++
	}
	if e.MarketFeeNetworkPercent != nil {
		count++
	}
	if e.MakerFeeDiscountPercent != nil {
		count++
	}
	w.writeVarUint64(uint64(count))
	if e.UpdatableHTLCOptions != nil {
		w.writeVarUint64(0)
		if err := e.UpdatableHTLCOptions.MarshalBinaryInto(w); err != nil {
			return err
		}
	}
	if e.CustomAuthorityOptions != nil {
		w.writeVarUint64(1)
		if err := e.CustomAuthorityOptions.MarshalBinaryInto(w); err != nil {
			return err
		}
	}
	if e.MarketFeeNetworkPercent != nil {
		w.writeVarUint64(2)
		w.writeUint16(*e.MarketFeeNetworkPercent)
	}
	if e.MakerFeeDiscountPercent != nil {
		w.writeVarUint64(3)
		w.writeUint16(*e.MakerFeeDiscountPercent)
	}
	return nil
}

func (e *ChainParametersExtensions) UnmarshalBinaryFrom(r *binaryReader) error {
	count, err := r.readVarUint64()
	if err != nil {
		return err
	}
	for i := uint64(0); i < count; i++ {
		index, err := r.readVarUint64()
		if err != nil {
			return err
		}
		switch index {
		case 0:
			var value HTLCOptions
			if err := value.UnmarshalBinaryFrom(r); err != nil {
				return err
			}
			e.UpdatableHTLCOptions = &value
		case 1:
			var value CustomAuthorityOptions
			if err := value.UnmarshalBinaryFrom(r); err != nil {
				return err
			}
			e.CustomAuthorityOptions = &value
		case 2:
			value, err := r.readUint16()
			if err != nil {
				return err
			}
			e.MarketFeeNetworkPercent = uint16Ptr(value)
		case 3:
			value, err := r.readUint16()
			if err != nil {
				return err
			}
			e.MakerFeeDiscountPercent = uint16Ptr(value)
		default:
			return fmt.Errorf("unknown chain parameters extension %d", index)
		}
	}
	return nil
}

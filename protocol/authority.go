package protocol

import (
	"encoding/json"
	"fmt"
	"sort"
)

// Authority models a BitShares authority structure.
type Authority struct {
	WeightThreshold uint32               `json:"weight_threshold"`
	AccountAuths    map[ObjectID]uint16  `json:"account_auths,omitempty"`
	KeyAuths        map[PublicKey]uint16 `json:"key_auths,omitempty"`
	AddressAuths    map[Address]uint16   `json:"address_auths,omitempty"`
}

func (a Authority) MarshalBinaryInto(w *binaryWriter) error {
	w.writeUint32(a.WeightThreshold)
	if err := writeAuthorityObjectAuths(w, a.AccountAuths); err != nil {
		return err
	}
	if err := writeAuthorityPublicKeyAuths(w, a.KeyAuths); err != nil {
		return err
	}
	if err := writeAuthorityAddressAuths(w, a.AddressAuths); err != nil {
		return err
	}
	return nil
}

func (a *Authority) UnmarshalBinaryFrom(r *binaryReader) error {
	threshold, err := r.readUint32()
	if err != nil {
		return err
	}
	accountAuths, err := readAuthorityObjectAuths(r)
	if err != nil {
		return err
	}
	keyAuths, err := readAuthorityPublicKeyAuths(r)
	if err != nil {
		return err
	}
	addressAuths, err := readAuthorityAddressAuths(r)
	if err != nil {
		return err
	}
	a.WeightThreshold = threshold
	a.AccountAuths = accountAuths
	a.KeyAuths = keyAuths
	a.AddressAuths = addressAuths
	return nil
}

func writeAuthorityObjectAuths(w *binaryWriter, values map[ObjectID]uint16) error {
	keys := make([]ObjectID, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return compareObjectID(keys[i], keys[j]) < 0 })
	w.writeVarUint64(uint64(len(keys)))
	for _, key := range keys {
		if err := key.MarshalBinaryInto(w); err != nil {
			return err
		}
		w.writeUint16(values[key])
	}
	return nil
}

func readAuthorityObjectAuths(r *binaryReader) (map[ObjectID]uint16, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make(map[ObjectID]uint16, count)
	for i := uint64(0); i < count; i++ {
		key, err := readObjectID(r)
		if err != nil {
			return nil, err
		}
		weight, err := r.readUint16()
		if err != nil {
			return nil, err
		}
		out[key] = weight
	}
	return out, nil
}

func writeAuthorityPublicKeyAuths(w *binaryWriter, values map[PublicKey]uint16) error {
	keys := make([]PublicKey, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return compareStrings(keys[i].String(), keys[j].String()) < 0 })
	w.writeVarUint64(uint64(len(keys)))
	for _, key := range keys {
		if err := key.MarshalBinaryInto(w); err != nil {
			return err
		}
		w.writeUint16(values[key])
	}
	return nil
}

func readAuthorityPublicKeyAuths(r *binaryReader) (map[PublicKey]uint16, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make(map[PublicKey]uint16, count)
	for i := uint64(0); i < count; i++ {
		key, err := readPublicKey(r)
		if err != nil {
			return nil, err
		}
		weight, err := r.readUint16()
		if err != nil {
			return nil, err
		}
		out[key] = weight
	}
	return out, nil
}

func writeAuthorityAddressAuths(w *binaryWriter, values map[Address]uint16) error {
	keys := make([]Address, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return compareStrings(keys[i].String(), keys[j].String()) < 0 })
	w.writeVarUint64(uint64(len(keys)))
	for _, key := range keys {
		if err := key.MarshalBinaryInto(w); err != nil {
			return err
		}
		w.writeUint16(values[key])
	}
	return nil
}

func readAuthorityAddressAuths(r *binaryReader) (map[Address]uint16, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make(map[Address]uint16, count)
	for i := uint64(0); i < count; i++ {
		key, err := readAddress(r)
		if err != nil {
			return nil, err
		}
		weight, err := r.readUint16()
		if err != nil {
			return nil, err
		}
		out[key] = weight
	}
	return out, nil
}

// AccountOptions holds the voting and memo settings for an account.
type AccountOptions struct {
	MemoKey       PublicKey         `json:"memo_key"`
	VotingAccount ObjectID          `json:"voting_account"`
	NumWitness    uint16            `json:"num_witness"`
	NumCommittee  uint16            `json:"num_committee"`
	Votes         []VoteID          `json:"votes,omitempty"`
	Extensions    []json.RawMessage `json:"extensions,omitempty"`
}

func (o AccountOptions) MarshalBinaryInto(w *binaryWriter) error {
	if err := o.MemoKey.MarshalBinaryInto(w); err != nil {
		return err
	}
	if err := o.VotingAccount.MarshalBinaryInto(w); err != nil {
		return err
	}
	w.writeUint16(o.NumWitness)
	w.writeUint16(o.NumCommittee)
	if err := writeVoteIDSet(w, o.Votes); err != nil {
		return err
	}
	if !extensionsEmpty(o.Extensions) {
		return fmt.Errorf("account options extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return nil
}

func (o *AccountOptions) UnmarshalBinaryFrom(r *binaryReader) error {
	memo, err := readPublicKey(r)
	if err != nil {
		return err
	}
	voting, err := readObjectID(r)
	if err != nil {
		return err
	}
	witness, err := r.readUint16()
	if err != nil {
		return err
	}
	committee, err := r.readUint16()
	if err != nil {
		return err
	}
	votes, err := readVoteIDSet(r)
	if err != nil {
		return err
	}
	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("account options extensions are not supported in binary serialization")
	}
	o.MemoKey = memo
	o.VotingAccount = voting
	o.NumWitness = witness
	o.NumCommittee = committee
	o.Votes = votes
	o.Extensions = nil
	return nil
}

// AssetOptionsExtensions models the extensions field of asset_options.
type AssetOptionsExtensions struct {
	RewardPercent             *uint16    `json:"reward_percent,omitempty"`
	WhitelistMarketFeeSharing []ObjectID `json:"whitelist_market_fee_sharing,omitempty"`
	TakerFeePercent           *uint16    `json:"taker_fee_percent,omitempty"`
}

func (e AssetOptionsExtensions) MarshalBinaryInto(w *binaryWriter) error {
	count := 0
	if e.RewardPercent != nil {
		count++
	}
	if len(e.WhitelistMarketFeeSharing) > 0 {
		count++
	}
	if e.TakerFeePercent != nil {
		count++
	}
	w.writeVarUint64(uint64(count))
	if e.RewardPercent != nil {
		w.writeVarUint64(0)
		w.writeUint16(*e.RewardPercent)
	}
	if len(e.WhitelistMarketFeeSharing) > 0 {
		w.writeVarUint64(1)
		if err := writeObjectIDSet(w, e.WhitelistMarketFeeSharing); err != nil {
			return err
		}
	}
	if e.TakerFeePercent != nil {
		w.writeVarUint64(2)
		w.writeUint16(*e.TakerFeePercent)
	}
	return nil
}

func (e *AssetOptionsExtensions) UnmarshalBinaryFrom(r *binaryReader) error {
	count, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
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
			e.RewardPercent = uint16Ptr(value)
		case 1:
			value, err := readObjectIDSet(r)
			if err != nil {
				return err
			}
			e.WhitelistMarketFeeSharing = value
		case 2:
			value, err := r.readUint16()
			if err != nil {
				return err
			}
			e.TakerFeePercent = uint16Ptr(value)
		default:
			return fmt.Errorf("unknown asset options extension %d", index)
		}
	}
	return nil
}

// AssetOptions models common asset settings.
type AssetOptions struct {
	MaxSupply            int64                   `json:"max_supply"`
	MarketFeePercent     uint16                  `json:"market_fee_percent"`
	MaxMarketFee         int64                   `json:"max_market_fee"`
	IssuerPermissions    uint16                  `json:"issuer_permissions"`
	Flags                uint16                  `json:"flags"`
	CoreExchangeRate     Price                   `json:"core_exchange_rate"`
	WhitelistAuthorities []ObjectID              `json:"whitelist_authorities,omitempty"`
	BlacklistAuthorities []ObjectID              `json:"blacklist_authorities,omitempty"`
	WhitelistMarkets     []ObjectID              `json:"whitelist_markets,omitempty"`
	BlacklistMarkets     []ObjectID              `json:"blacklist_markets,omitempty"`
	Description          string                  `json:"description"`
	Extensions           *AssetOptionsExtensions `json:"extensions,omitempty"`
}

func (o AssetOptions) MarshalBinaryInto(w *binaryWriter) error {
	w.writeInt64(o.MaxSupply)
	w.writeUint16(o.MarketFeePercent)
	w.writeInt64(o.MaxMarketFee)
	w.writeUint16(o.IssuerPermissions)
	w.writeUint16(o.Flags)
	if err := o.CoreExchangeRate.MarshalBinaryInto(w); err != nil {
		return err
	}
	if err := writeObjectIDSet(w, o.WhitelistAuthorities); err != nil {
		return err
	}
	if err := writeObjectIDSet(w, o.BlacklistAuthorities); err != nil {
		return err
	}
	if err := writeObjectIDSet(w, o.WhitelistMarkets); err != nil {
		return err
	}
	if err := writeObjectIDSet(w, o.BlacklistMarkets); err != nil {
		return err
	}
	binaryWriteString(w, o.Description)
	if o.Extensions == nil {
		w.writeVarUint64(0)
		return nil
	}
	return o.Extensions.MarshalBinaryInto(w)
}

func (o *AssetOptions) UnmarshalBinaryFrom(r *binaryReader) error {
	maxSupply, err := r.readInt64()
	if err != nil {
		return err
	}
	marketFeePercent, err := r.readUint16()
	if err != nil {
		return err
	}
	maxMarketFee, err := r.readInt64()
	if err != nil {
		return err
	}
	issuerPermissions, err := r.readUint16()
	if err != nil {
		return err
	}
	flags, err := r.readUint16()
	if err != nil {
		return err
	}
	coreRate, err := readPrice(r)
	if err != nil {
		return err
	}
	whitelistAuthorities, err := readObjectIDSet(r)
	if err != nil {
		return err
	}
	blacklistAuthorities, err := readObjectIDSet(r)
	if err != nil {
		return err
	}
	whitelistMarkets, err := readObjectIDSet(r)
	if err != nil {
		return err
	}
	blacklistMarkets, err := readObjectIDSet(r)
	if err != nil {
		return err
	}
	description, err := binaryReadString(r)
	if err != nil {
		return err
	}
	var ext AssetOptionsExtensions
	if err := ext.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	o.MaxSupply = maxSupply
	o.MarketFeePercent = marketFeePercent
	o.MaxMarketFee = maxMarketFee
	o.IssuerPermissions = issuerPermissions
	o.Flags = flags
	o.CoreExchangeRate = coreRate
	o.WhitelistAuthorities = whitelistAuthorities
	o.BlacklistAuthorities = blacklistAuthorities
	o.WhitelistMarkets = whitelistMarkets
	o.BlacklistMarkets = blacklistMarkets
	o.Description = description
	if ext.RewardPercent != nil || len(ext.WhitelistMarketFeeSharing) > 0 || ext.TakerFeePercent != nil {
		o.Extensions = &ext
	} else {
		o.Extensions = nil
	}
	return nil
}

// BitAssetOptionsExtensions models the extensions field of bitasset_options.
type BitAssetOptionsExtensions struct {
	InitialCollateralRatio     *uint16 `json:"initial_collateral_ratio,omitempty"`
	MaintenanceCollateralRatio *uint16 `json:"maintenance_collateral_ratio,omitempty"`
	MaximumShortSqueezeRatio   *uint16 `json:"maximum_short_squeeze_ratio,omitempty"`
	MarginCallFeeRatio         *uint16 `json:"margin_call_fee_ratio,omitempty"`
	ForceSettleFeePercent      *uint16 `json:"force_settle_fee_percent,omitempty"`
	BlackSwanResponseMethod    *uint8  `json:"black_swan_response_method,omitempty"`
}

func (e BitAssetOptionsExtensions) MarshalBinaryInto(w *binaryWriter) error {
	count := 0
	if e.InitialCollateralRatio != nil {
		count++
	}
	if e.MaintenanceCollateralRatio != nil {
		count++
	}
	if e.MaximumShortSqueezeRatio != nil {
		count++
	}
	if e.MarginCallFeeRatio != nil {
		count++
	}
	if e.ForceSettleFeePercent != nil {
		count++
	}
	if e.BlackSwanResponseMethod != nil {
		count++
	}
	w.writeVarUint64(uint64(count))
	if e.InitialCollateralRatio != nil {
		w.writeVarUint64(0)
		w.writeUint16(*e.InitialCollateralRatio)
	}
	if e.MaintenanceCollateralRatio != nil {
		w.writeVarUint64(1)
		w.writeUint16(*e.MaintenanceCollateralRatio)
	}
	if e.MaximumShortSqueezeRatio != nil {
		w.writeVarUint64(2)
		w.writeUint16(*e.MaximumShortSqueezeRatio)
	}
	if e.MarginCallFeeRatio != nil {
		w.writeVarUint64(3)
		w.writeUint16(*e.MarginCallFeeRatio)
	}
	if e.ForceSettleFeePercent != nil {
		w.writeVarUint64(4)
		w.writeUint16(*e.ForceSettleFeePercent)
	}
	if e.BlackSwanResponseMethod != nil {
		w.writeVarUint64(5)
		w.writeUint8(*e.BlackSwanResponseMethod)
	}
	return nil
}

func (e *BitAssetOptionsExtensions) UnmarshalBinaryFrom(r *binaryReader) error {
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
		case 1:
			value, err := r.readUint16()
			if err != nil {
				return err
			}
			e.MaintenanceCollateralRatio = uint16Ptr(value)
		case 2:
			value, err := r.readUint16()
			if err != nil {
				return err
			}
			e.MaximumShortSqueezeRatio = uint16Ptr(value)
		case 3:
			value, err := r.readUint16()
			if err != nil {
				return err
			}
			e.MarginCallFeeRatio = uint16Ptr(value)
		case 4:
			value, err := r.readUint16()
			if err != nil {
				return err
			}
			e.ForceSettleFeePercent = uint16Ptr(value)
		case 5:
			value, err := r.readUint8()
			if err != nil {
				return err
			}
			e.BlackSwanResponseMethod = uint8Ptr(value)
		default:
			return fmt.Errorf("unknown bitasset options extension %d", index)
		}
	}
	return nil
}

// BitAssetOptions models the BTS bitasset settings.
type BitAssetOptions struct {
	FeedLifetimeSec              uint32                     `json:"feed_lifetime_sec"`
	MinimumFeeds                 uint8                      `json:"minimum_feeds"`
	ForceSettlementDelaySec      uint32                     `json:"force_settlement_delay_sec"`
	ForceSettlementOffsetPercent uint16                     `json:"force_settlement_offset_percent"`
	MaximumForceSettlementVolume uint16                     `json:"maximum_force_settlement_volume"`
	ShortBackingAsset            ObjectID                   `json:"short_backing_asset"`
	Extensions                   *BitAssetOptionsExtensions `json:"extensions,omitempty"`
}

func (o BitAssetOptions) MarshalBinaryInto(w *binaryWriter) error {
	w.writeUint32(o.FeedLifetimeSec)
	w.writeUint8(o.MinimumFeeds)
	w.writeUint32(o.ForceSettlementDelaySec)
	w.writeUint16(o.ForceSettlementOffsetPercent)
	w.writeUint16(o.MaximumForceSettlementVolume)
	if err := o.ShortBackingAsset.MarshalBinaryInto(w); err != nil {
		return err
	}
	if o.Extensions == nil {
		w.writeVarUint64(0)
		return nil
	}
	return o.Extensions.MarshalBinaryInto(w)
}

func (o *BitAssetOptions) UnmarshalBinaryFrom(r *binaryReader) error {
	feedLifetimeSec, err := r.readUint32()
	if err != nil {
		return err
	}
	minimumFeeds, err := r.readUint8()
	if err != nil {
		return err
	}
	forceSettlementDelaySec, err := r.readUint32()
	if err != nil {
		return err
	}
	offsetPercent, err := r.readUint16()
	if err != nil {
		return err
	}
	maxVolume, err := r.readUint16()
	if err != nil {
		return err
	}
	shortBackingAsset, err := readObjectID(r)
	if err != nil {
		return err
	}
	var ext BitAssetOptionsExtensions
	if err := ext.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	o.FeedLifetimeSec = feedLifetimeSec
	o.MinimumFeeds = minimumFeeds
	o.ForceSettlementDelaySec = forceSettlementDelaySec
	o.ForceSettlementOffsetPercent = offsetPercent
	o.MaximumForceSettlementVolume = maxVolume
	o.ShortBackingAsset = shortBackingAsset
	if ext.InitialCollateralRatio != nil || ext.MaintenanceCollateralRatio != nil || ext.MaximumShortSqueezeRatio != nil || ext.MarginCallFeeRatio != nil || ext.ForceSettleFeePercent != nil || ext.BlackSwanResponseMethod != nil {
		o.Extensions = &ext
	} else {
		o.Extensions = nil
	}
	return nil
}

// PriceFeed is the structured feed payload used by asset_publish_feed.
type PriceFeed struct {
	SettlementPrice            Price  `json:"settlement_price"`
	MaintenanceCollateralRatio uint16 `json:"maintenance_collateral_ratio"`
	MaximumShortSqueezeRatio   uint16 `json:"maximum_short_squeeze_ratio"`
	CoreExchangeRate           Price  `json:"core_exchange_rate"`
}

func (p PriceFeed) MarshalBinaryInto(w *binaryWriter) error {
	if err := p.SettlementPrice.MarshalBinaryInto(w); err != nil {
		return err
	}
	w.writeUint16(p.MaintenanceCollateralRatio)
	w.writeUint16(p.MaximumShortSqueezeRatio)
	return p.CoreExchangeRate.MarshalBinaryInto(w)
}

func (p *PriceFeed) UnmarshalBinaryFrom(r *binaryReader) error {
	settlement, err := readPrice(r)
	if err != nil {
		return err
	}
	maintenance, err := r.readUint16()
	if err != nil {
		return err
	}
	maximum, err := r.readUint16()
	if err != nil {
		return err
	}
	core, err := readPrice(r)
	if err != nil {
		return err
	}
	p.SettlementPrice = settlement
	p.MaintenanceCollateralRatio = maintenance
	p.MaximumShortSqueezeRatio = maximum
	p.CoreExchangeRate = core
	return nil
}

package protocol

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// AccountNameEqLitPredicate matches an account by name.
type AccountNameEqLitPredicate struct {
	AccountID ObjectID `json:"account_id"`
	Name      string   `json:"name"`
}

// AssetSymbolEqLitPredicate matches an asset by symbol.
type AssetSymbolEqLitPredicate struct {
	AssetID ObjectID `json:"asset_id"`
	Symbol  string   `json:"symbol"`
}

// BlockIDPredicate matches a specific block id.
type BlockIDPredicate struct {
	ID string `json:"id"`
}

// Predicate is a static variant used by assert operations.
type Predicate struct {
	Kind             uint16
	AccountNameEqLit *AccountNameEqLitPredicate
	AssetSymbolEqLit *AssetSymbolEqLitPredicate
	BlockIDPredicate *BlockIDPredicate
}

func (p Predicate) MarshalJSON() ([]byte, error) {
	switch p.Kind {
	case 0:
		if p.AccountNameEqLit == nil {
			return nil, fmt.Errorf("missing account_name_eq_lit predicate")
		}
		return json.Marshal([]any{uint16(0), p.AccountNameEqLit})
	case 1:
		if p.AssetSymbolEqLit == nil {
			return nil, fmt.Errorf("missing asset_symbol_eq_lit predicate")
		}
		return json.Marshal([]any{uint16(1), p.AssetSymbolEqLit})
	case 2:
		if p.BlockIDPredicate == nil {
			return nil, fmt.Errorf("missing block_id predicate")
		}
		return json.Marshal([]any{uint16(2), p.BlockIDPredicate})
	default:
		return nil, fmt.Errorf("unsupported predicate type %d", p.Kind)
	}
}

func (p *Predicate) UnmarshalJSON(data []byte) error {
	var body []json.RawMessage
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	if len(body) != 2 {
		return fmt.Errorf("invalid predicate")
	}
	var kind uint16
	if err := json.Unmarshal(body[0], &kind); err != nil {
		return err
	}
	p.Kind = kind
	switch kind {
	case 0:
		var value AccountNameEqLitPredicate
		if err := json.Unmarshal(body[1], &value); err != nil {
			return err
		}
		p.AccountNameEqLit = &value
	case 1:
		var value AssetSymbolEqLitPredicate
		if err := json.Unmarshal(body[1], &value); err != nil {
			return err
		}
		p.AssetSymbolEqLit = &value
	case 2:
		var value BlockIDPredicate
		if err := json.Unmarshal(body[1], &value); err != nil {
			return err
		}
		p.BlockIDPredicate = &value
	default:
		return fmt.Errorf("unsupported predicate type %d", kind)
	}
	return nil
}

func (p Predicate) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(p.Kind))
	switch p.Kind {
	case 0:
		if p.AccountNameEqLit == nil {
			return nil, fmt.Errorf("missing account_name_eq_lit predicate")
		}
		if err := p.AccountNameEqLit.AccountID.MarshalBinaryInto(w); err != nil {
			return nil, err
		}
		w.writeString(p.AccountNameEqLit.Name)
	case 1:
		if p.AssetSymbolEqLit == nil {
			return nil, fmt.Errorf("missing asset_symbol_eq_lit predicate")
		}
		if err := p.AssetSymbolEqLit.AssetID.MarshalBinaryInto(w); err != nil {
			return nil, err
		}
		w.writeString(p.AssetSymbolEqLit.Symbol)
	case 2:
		if p.BlockIDPredicate == nil {
			return nil, fmt.Errorf("missing block_id predicate")
		}
		raw, err := hex.DecodeString(strings.TrimSpace(p.BlockIDPredicate.ID))
		if err != nil {
			return nil, err
		}
		if len(raw) != 20 {
			return nil, fmt.Errorf("block id must be 20 bytes")
		}
		if err := writeFixedBytes(w, raw, 20); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported predicate type %d", p.Kind)
	}
	return w.Bytes(), nil
}

func (p *Predicate) UnmarshalBinaryFrom(r *binaryReader) error {
	kind, err := r.readVarUint64()
	if err != nil {
		return err
	}
	p.Kind = uint16(kind)
	switch p.Kind {
	case 0:
		account, err := readObjectID(r)
		if err != nil {
			return err
		}
		name, err := r.readString()
		if err != nil {
			return err
		}
		p.AccountNameEqLit = &AccountNameEqLitPredicate{AccountID: account, Name: name}
	case 1:
		asset, err := readObjectID(r)
		if err != nil {
			return err
		}
		symbol, err := r.readString()
		if err != nil {
			return err
		}
		p.AssetSymbolEqLit = &AssetSymbolEqLitPredicate{AssetID: asset, Symbol: symbol}
	case 2:
		raw, err := readFixedBytes(r, 20)
		if err != nil {
			return err
		}
		p.BlockIDPredicate = &BlockIDPredicate{ID: hex.EncodeToString(raw)}
	default:
		return fmt.Errorf("unsupported predicate type %d", p.Kind)
	}
	return nil
}

// LinearVestingPolicyInitializer is the standard vesting policy.
type LinearVestingPolicyInitializer struct {
	BeginTimestamp         Time   `json:"begin_timestamp"`
	VestingCliffSeconds    uint32 `json:"vesting_cliff_seconds"`
	VestingDurationSeconds uint32 `json:"vesting_duration_seconds"`
}

// CDDVestingPolicyInitializer is the coin-days-destroyed policy.
type CDDVestingPolicyInitializer struct {
	StartClaim     Time   `json:"start_claim"`
	VestingSeconds uint32 `json:"vesting_seconds"`
}

// VestingPolicyInitializer is a static variant.
type VestingPolicyInitializer struct {
	Kind   uint16
	Linear *LinearVestingPolicyInitializer
	CDD    *CDDVestingPolicyInitializer
}

func (v VestingPolicyInitializer) MarshalJSON() ([]byte, error) {
	switch v.Kind {
	case 0:
		if v.Linear == nil {
			return nil, fmt.Errorf("missing linear vesting initializer")
		}
		return json.Marshal([]any{uint16(0), v.Linear})
	case 1:
		if v.CDD == nil {
			return nil, fmt.Errorf("missing cdd vesting initializer")
		}
		return json.Marshal([]any{uint16(1), v.CDD})
	default:
		return nil, fmt.Errorf("unsupported vesting policy type %d", v.Kind)
	}
}

func (v *VestingPolicyInitializer) UnmarshalJSON(data []byte) error {
	var body []json.RawMessage
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	if len(body) != 2 {
		return fmt.Errorf("invalid vesting policy initializer")
	}
	var kind uint16
	if err := json.Unmarshal(body[0], &kind); err != nil {
		return err
	}
	v.Kind = kind
	switch kind {
	case 0:
		var value LinearVestingPolicyInitializer
		if err := json.Unmarshal(body[1], &value); err != nil {
			return err
		}
		v.Linear = &value
	case 1:
		var value CDDVestingPolicyInitializer
		if err := json.Unmarshal(body[1], &value); err != nil {
			return err
		}
		v.CDD = &value
	default:
		return fmt.Errorf("unsupported vesting policy type %d", kind)
	}
	return nil
}

func (v VestingPolicyInitializer) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(v.Kind))
	switch v.Kind {
	case 0:
		if v.Linear == nil {
			return nil, fmt.Errorf("missing linear vesting initializer")
		}
		if err := v.Linear.BeginTimestamp.MarshalBinaryInto(w); err != nil {
			return nil, err
		}
		w.writeUint32(v.Linear.VestingCliffSeconds)
		w.writeUint32(v.Linear.VestingDurationSeconds)
	case 1:
		if v.CDD == nil {
			return nil, fmt.Errorf("missing cdd vesting initializer")
		}
		if err := v.CDD.StartClaim.MarshalBinaryInto(w); err != nil {
			return nil, err
		}
		w.writeUint32(v.CDD.VestingSeconds)
	default:
		return nil, fmt.Errorf("unsupported vesting policy type %d", v.Kind)
	}
	return w.Bytes(), nil
}

func (v *VestingPolicyInitializer) UnmarshalBinaryFrom(r *binaryReader) error {
	kind, err := r.readVarUint64()
	if err != nil {
		return err
	}
	v.Kind = uint16(kind)
	switch v.Kind {
	case 0:
		begin, err := readTime(r)
		if err != nil {
			return err
		}
		cliff, err := r.readUint32()
		if err != nil {
			return err
		}
		duration, err := r.readUint32()
		if err != nil {
			return err
		}
		v.Linear = &LinearVestingPolicyInitializer{BeginTimestamp: begin, VestingCliffSeconds: cliff, VestingDurationSeconds: duration}
	case 1:
		start, err := readTime(r)
		if err != nil {
			return err
		}
		seconds, err := r.readUint32()
		if err != nil {
			return err
		}
		v.CDD = &CDDVestingPolicyInitializer{StartClaim: start, VestingSeconds: seconds}
	default:
		return fmt.Errorf("unsupported vesting policy type %d", v.Kind)
	}
	return nil
}

// RefundWorkerInitializer is the default worker initializer.
type RefundWorkerInitializer struct{}

// VestingBalanceWorkerInitializer pays into a vesting balance.
type VestingBalanceWorkerInitializer struct {
	PayVestingPeriodDays uint16 `json:"pay_vesting_period_days"`
}

// BurnWorkerInitializer destroys the worker payout.
type BurnWorkerInitializer struct{}

// WorkerInitializer is a static variant.
type WorkerInitializer struct {
	Kind    uint16
	Refund  *RefundWorkerInitializer
	Vesting *VestingBalanceWorkerInitializer
	Burn    *BurnWorkerInitializer
}

func (w WorkerInitializer) MarshalJSON() ([]byte, error) {
	switch w.Kind {
	case 0:
		return json.Marshal([]any{uint16(0), struct{}{}})
	case 1:
		if w.Vesting == nil {
			return nil, fmt.Errorf("missing vesting worker initializer")
		}
		return json.Marshal([]any{uint16(1), w.Vesting})
	case 2:
		return json.Marshal([]any{uint16(2), struct{}{}})
	default:
		return nil, fmt.Errorf("unsupported worker initializer type %d", w.Kind)
	}
}

func (w *WorkerInitializer) UnmarshalJSON(data []byte) error {
	var body []json.RawMessage
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	if len(body) != 2 {
		return fmt.Errorf("invalid worker initializer")
	}
	var kind uint16
	if err := json.Unmarshal(body[0], &kind); err != nil {
		return err
	}
	w.Kind = kind
	switch kind {
	case 0:
		w.Refund = &RefundWorkerInitializer{}
	case 1:
		var value VestingBalanceWorkerInitializer
		if err := json.Unmarshal(body[1], &value); err != nil {
			return err
		}
		w.Vesting = &value
	case 2:
		w.Burn = &BurnWorkerInitializer{}
	default:
		return fmt.Errorf("unsupported worker initializer type %d", kind)
	}
	return nil
}

func (w WorkerInitializer) MarshalBinary() ([]byte, error) {
	wr := newBinaryWriter()
	wr.writeVarUint64(uint64(w.Kind))
	switch w.Kind {
	case 0:
	case 1:
		if w.Vesting == nil {
			return nil, fmt.Errorf("missing vesting worker initializer")
		}
		wr.writeUint16(w.Vesting.PayVestingPeriodDays)
	case 2:
	default:
		return nil, fmt.Errorf("unsupported worker initializer type %d", w.Kind)
	}
	return wr.Bytes(), nil
}

func (w *WorkerInitializer) UnmarshalBinaryFrom(r *binaryReader) error {
	kind, err := r.readVarUint64()
	if err != nil {
		return err
	}
	w.Kind = uint16(kind)
	switch w.Kind {
	case 0:
		w.Refund = &RefundWorkerInitializer{}
	case 1:
		days, err := r.readUint16()
		if err != nil {
			return err
		}
		w.Vesting = &VestingBalanceWorkerInitializer{PayVestingPeriodDays: days}
	case 2:
		w.Burn = &BurnWorkerInitializer{}
	default:
		return fmt.Errorf("unsupported worker initializer type %d", w.Kind)
	}
	return nil
}

// HTLCPreimageHash identifies the hash type and bytes used by an HTLC.
type HTLCPreimageHash struct {
	Kind  uint16
	Value string
}

func (h HTLCPreimageHash) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{h.Kind, h.Value})
}

func (h *HTLCPreimageHash) UnmarshalJSON(data []byte) error {
	var body []json.RawMessage
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	if len(body) != 2 {
		return fmt.Errorf("invalid htlc preimage hash")
	}
	if err := json.Unmarshal(body[0], &h.Kind); err != nil {
		return err
	}
	if err := json.Unmarshal(body[1], &h.Value); err != nil {
		return err
	}
	return nil
}

func (h HTLCPreimageHash) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(h.Kind))
	raw, err := hex.DecodeString(strings.TrimSpace(h.Value))
	if err != nil {
		return nil, err
	}
	switch h.Kind {
	case 0, 1, 3:
		if len(raw) != 20 {
			return nil, fmt.Errorf("htlc hash must be 20 bytes for kind %d", h.Kind)
		}
	case 2:
		if len(raw) != 32 {
			return nil, fmt.Errorf("htlc hash must be 32 bytes for kind %d", h.Kind)
		}
	default:
		return nil, fmt.Errorf("unsupported htlc hash kind %d", h.Kind)
	}
	if err := writeFixedBytes(w, raw, len(raw)); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (h *HTLCPreimageHash) UnmarshalBinaryFrom(r *binaryReader) error {
	kind, err := r.readVarUint64()
	if err != nil {
		return err
	}
	var size int
	switch uint16(kind) {
	case 0, 1, 3:
		size = 20
	case 2:
		size = 32
	default:
		return fmt.Errorf("unsupported htlc hash kind %d", kind)
	}
	raw, err := readFixedBytes(r, size)
	if err != nil {
		return err
	}
	h.Kind = uint16(kind)
	h.Value = hex.EncodeToString(raw)
	return nil
}

// AssetClaimFeesExtensions carries the optional claim asset id.
type AssetClaimFeesExtensions struct {
	ClaimFromAssetID *ObjectID `json:"claim_from_asset_id,omitempty"`
}

func (e AssetClaimFeesExtensions) MarshalBinaryInto(w *binaryWriter) error {
	if e.ClaimFromAssetID == nil {
		w.writeVarUint64(0)
		return nil
	}
	w.writeVarUint64(1)
	w.writeVarUint64(0)
	return e.ClaimFromAssetID.MarshalBinaryInto(w)
}

func (e *AssetClaimFeesExtensions) UnmarshalBinaryFrom(r *binaryReader) error {
	count, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	if _, err := r.readVarUint64(); err != nil {
		return err
	}
	value, err := readObjectID(r)
	if err != nil {
		return err
	}
	e.ClaimFromAssetID = &value
	return nil
}

// ChainParameters mirrors the committee global parameters structure.
type ChainParameters struct {
	CurrentFees                      FeeSchedule               `json:"current_fees,omitempty"`
	BlockInterval                    uint8                     `json:"block_interval"`
	MaintenanceInterval              uint32                    `json:"maintenance_interval"`
	MaintenanceSkipSlots             uint8                     `json:"maintenance_skip_slots"`
	CommitteeProposalReviewPeriod    uint32                    `json:"committee_proposal_review_period"`
	MaximumTransactionSize           uint32                    `json:"maximum_transaction_size"`
	MaximumBlockSize                 uint32                    `json:"maximum_block_size"`
	MaximumTimeUntilExpiration       uint32                    `json:"maximum_time_until_expiration"`
	MaximumProposalLifetime          uint32                    `json:"maximum_proposal_lifetime"`
	MaximumAssetWhitelistAuthorities uint8                     `json:"maximum_asset_whitelist_authorities"`
	MaximumAssetFeedPublishers       uint8                     `json:"maximum_asset_feed_publishers"`
	MaximumWitnessCount              uint16                    `json:"maximum_witness_count"`
	MaximumCommitteeCount            uint16                    `json:"maximum_committee_count"`
	MaximumAuthorityMembership       uint16                    `json:"maximum_authority_membership"`
	ReservePercentOfFee              uint16                    `json:"reserve_percent_of_fee"`
	NetworkPercentOfFee              uint16                    `json:"network_percent_of_fee"`
	LifetimeReferrerPercentOfFee     uint16                    `json:"lifetime_referrer_percent_of_fee"`
	CashbackVestingPeriodSeconds     uint32                    `json:"cashback_vesting_period_seconds"`
	CashbackVestingThreshold         int64                     `json:"cashback_vesting_threshold"`
	CountNonMemberVotes              bool                      `json:"count_non_member_votes"`
	AllowNonMemberWhitelists         bool                      `json:"allow_non_member_whitelists"`
	WitnessPayPerBlock               int64                     `json:"witness_pay_per_block"`
	WitnessPayVestingSeconds         uint32                    `json:"witness_pay_vesting_seconds"`
	WorkerBudgetPerDay               int64                     `json:"worker_budget_per_day"`
	MaxPredicateOpcode               uint16                    `json:"max_predicate_opcode"`
	FeeLiquidationThreshold          int64                     `json:"fee_liquidation_threshold"`
	AccountsPerFeeScale              uint16                    `json:"accounts_per_fee_scale"`
	AccountFeeScaleBitshifts         uint8                     `json:"account_fee_scale_bitshifts"`
	MaxAuthorityDepth                uint8                     `json:"max_authority_depth"`
	Extensions                       ChainParametersExtensions `json:"extensions"`
}

func (c ChainParameters) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	if err := c.CurrentFees.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeUint8(c.BlockInterval)
	w.writeUint32(c.MaintenanceInterval)
	w.writeUint8(c.MaintenanceSkipSlots)
	w.writeUint32(c.CommitteeProposalReviewPeriod)
	w.writeUint32(c.MaximumTransactionSize)
	w.writeUint32(c.MaximumBlockSize)
	w.writeUint32(c.MaximumTimeUntilExpiration)
	w.writeUint32(c.MaximumProposalLifetime)
	w.writeUint8(c.MaximumAssetWhitelistAuthorities)
	w.writeUint8(c.MaximumAssetFeedPublishers)
	w.writeUint16(c.MaximumWitnessCount)
	w.writeUint16(c.MaximumCommitteeCount)
	w.writeUint16(c.MaximumAuthorityMembership)
	w.writeUint16(c.ReservePercentOfFee)
	w.writeUint16(c.NetworkPercentOfFee)
	w.writeUint16(c.LifetimeReferrerPercentOfFee)
	w.writeUint32(c.CashbackVestingPeriodSeconds)
	w.writeInt64(c.CashbackVestingThreshold)
	w.writeUint8(boolByte(c.CountNonMemberVotes))
	w.writeUint8(boolByte(c.AllowNonMemberWhitelists))
	w.writeInt64(c.WitnessPayPerBlock)
	w.writeUint32(c.WitnessPayVestingSeconds)
	w.writeInt64(c.WorkerBudgetPerDay)
	w.writeUint16(c.MaxPredicateOpcode)
	w.writeInt64(c.FeeLiquidationThreshold)
	w.writeUint16(c.AccountsPerFeeScale)
	w.writeUint8(c.AccountFeeScaleBitshifts)
	w.writeUint8(c.MaxAuthorityDepth)
	if err := c.Extensions.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (c *ChainParameters) UnmarshalBinaryFrom(r *binaryReader) error {
	var currentFees FeeSchedule
	if err := currentFees.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	blockInterval, err := r.readUint8()
	if err != nil {
		return err
	}
	maintenanceInterval, err := r.readUint32()
	if err != nil {
		return err
	}
	maintenanceSkipSlots, err := r.readUint8()
	if err != nil {
		return err
	}
	review, err := r.readUint32()
	if err != nil {
		return err
	}
	maxTx, err := r.readUint32()
	if err != nil {
		return err
	}
	maxBlock, err := r.readUint32()
	if err != nil {
		return err
	}
	maxExpiration, err := r.readUint32()
	if err != nil {
		return err
	}
	maxProposal, err := r.readUint32()
	if err != nil {
		return err
	}
	maxWhitelist, err := r.readUint8()
	if err != nil {
		return err
	}
	maxFeed, err := r.readUint8()
	if err != nil {
		return err
	}
	maxWitness, err := r.readUint16()
	if err != nil {
		return err
	}
	maxCommittee, err := r.readUint16()
	if err != nil {
		return err
	}
	maxAuthority, err := r.readUint16()
	if err != nil {
		return err
	}
	reserve, err := r.readUint16()
	if err != nil {
		return err
	}
	network, err := r.readUint16()
	if err != nil {
		return err
	}
	lifetime, err := r.readUint16()
	if err != nil {
		return err
	}
	cashbackPeriod, err := r.readUint32()
	if err != nil {
		return err
	}
	cashbackThreshold, err := r.readInt64()
	if err != nil {
		return err
	}
	countVotes, err := readBool(r)
	if err != nil {
		return err
	}
	allowWhitelists, err := readBool(r)
	if err != nil {
		return err
	}
	witnessPay, err := r.readInt64()
	if err != nil {
		return err
	}
	witnessVestingSeconds, err := r.readUint32()
	if err != nil {
		return err
	}
	workerBudget, err := r.readInt64()
	if err != nil {
		return err
	}
	maxPredicate, err := r.readUint16()
	if err != nil {
		return err
	}
	feeLiquidation, err := r.readInt64()
	if err != nil {
		return err
	}
	accountsPerFeeScale, err := r.readUint16()
	if err != nil {
		return err
	}
	accountBitshifts, err := r.readUint8()
	if err != nil {
		return err
	}
	maxDepth, err := r.readUint8()
	if err != nil {
		return err
	}
	var ext ChainParametersExtensions
	if err := ext.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	c.CurrentFees = currentFees
	c.BlockInterval = blockInterval
	c.MaintenanceInterval = maintenanceInterval
	c.MaintenanceSkipSlots = maintenanceSkipSlots
	c.CommitteeProposalReviewPeriod = review
	c.MaximumTransactionSize = maxTx
	c.MaximumBlockSize = maxBlock
	c.MaximumTimeUntilExpiration = maxExpiration
	c.MaximumProposalLifetime = maxProposal
	c.MaximumAssetWhitelistAuthorities = maxWhitelist
	c.MaximumAssetFeedPublishers = maxFeed
	c.MaximumWitnessCount = maxWitness
	c.MaximumCommitteeCount = maxCommittee
	c.MaximumAuthorityMembership = maxAuthority
	c.ReservePercentOfFee = reserve
	c.NetworkPercentOfFee = network
	c.LifetimeReferrerPercentOfFee = lifetime
	c.CashbackVestingPeriodSeconds = cashbackPeriod
	c.CashbackVestingThreshold = cashbackThreshold
	c.CountNonMemberVotes = countVotes
	c.AllowNonMemberWhitelists = allowWhitelists
	c.WitnessPayPerBlock = witnessPay
	c.WitnessPayVestingSeconds = witnessVestingSeconds
	c.WorkerBudgetPerDay = workerBudget
	c.MaxPredicateOpcode = maxPredicate
	c.FeeLiquidationThreshold = feeLiquidation
	c.AccountsPerFeeScale = accountsPerFeeScale
	c.AccountFeeScaleBitshifts = accountBitshifts
	c.MaxAuthorityDepth = maxDepth
	c.Extensions = ext
	return nil
}

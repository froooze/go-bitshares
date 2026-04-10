package protocol

import "encoding/json"

// CallOrderUpdateOperation updates collateral and debt on a margin position.
type CallOrderUpdateOperation struct {
	Fee             AssetAmount               `json:"fee"`
	FundingAccount  ObjectID                  `json:"funding_account"`
	DeltaCollateral AssetAmount               `json:"delta_collateral"`
	DeltaDebt       AssetAmount               `json:"delta_debt"`
	Extensions      CallOrderUpdateExtensions `json:"extensions"`
}

func (o CallOrderUpdateOperation) Type() OperationType { return OperationTypeCallOrderUpdate }

func (o CallOrderUpdateOperation) MarshalJSON() ([]byte, error) {
	type alias CallOrderUpdateOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *CallOrderUpdateOperation) UnmarshalJSON(data []byte) error {
	type alias CallOrderUpdateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeCallOrderUpdate, &payload); err != nil {
		return err
	}
	*o = CallOrderUpdateOperation(payload)
	return nil
}

// FillOrderOperation records a matched order fill.
type FillOrderOperation struct {
	Fee       AssetAmount `json:"fee"`
	OrderID   ObjectID    `json:"order_id"`
	AccountID ObjectID    `json:"account_id"`
	Pays      AssetAmount `json:"pays"`
	Receives  AssetAmount `json:"receives"`
	FillPrice Price       `json:"fill_price"`
	IsMaker   bool        `json:"is_maker"`
}

func (o FillOrderOperation) Type() OperationType { return OperationTypeFillOrder }

func (o FillOrderOperation) MarshalJSON() ([]byte, error) {
	type alias FillOrderOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *FillOrderOperation) UnmarshalJSON(data []byte) error {
	type alias FillOrderOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeFillOrder, &payload); err != nil {
		return err
	}
	*o = FillOrderOperation(payload)
	return nil
}

// AccountCreateOperation creates a new account.
type AccountCreateOperation struct {
	Fee             AssetAmount             `json:"fee"`
	Registrar       ObjectID                `json:"registrar"`
	Referrer        ObjectID                `json:"referrer"`
	ReferrerPercent uint16                  `json:"referrer_percent"`
	Name            string                  `json:"name"`
	Owner           Authority               `json:"owner"`
	Active          Authority               `json:"active"`
	Options         AccountOptions          `json:"options"`
	Extensions      AccountCreateExtensions `json:"extensions"`
}

func (o AccountCreateOperation) Type() OperationType { return OperationTypeAccountCreate }

func (o AccountCreateOperation) MarshalJSON() ([]byte, error) {
	type alias AccountCreateOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *AccountCreateOperation) UnmarshalJSON(data []byte) error {
	type alias AccountCreateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeAccountCreate, &payload); err != nil {
		return err
	}
	*o = AccountCreateOperation(payload)
	return nil
}

// AccountUpdateOperation updates account authorities and options.
type AccountUpdateOperation struct {
	Fee        AssetAmount             `json:"fee"`
	Account    ObjectID                `json:"account"`
	Owner      *Authority              `json:"owner,omitempty"`
	Active     *Authority              `json:"active,omitempty"`
	NewOptions *AccountOptions         `json:"new_options,omitempty"`
	Extensions AccountUpdateExtensions `json:"extensions"`
}

func (o AccountUpdateOperation) Type() OperationType { return OperationTypeAccountUpdate }

func (o AccountUpdateOperation) MarshalJSON() ([]byte, error) {
	type alias AccountUpdateOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *AccountUpdateOperation) UnmarshalJSON(data []byte) error {
	type alias AccountUpdateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeAccountUpdate, &payload); err != nil {
		return err
	}
	*o = AccountUpdateOperation(payload)
	return nil
}

type AccountWhitelistOperation struct {
	Fee                AssetAmount       `json:"fee"`
	AuthorizingAccount ObjectID          `json:"authorizing_account"`
	AccountToList      ObjectID          `json:"account_to_list"`
	NewListing         uint8             `json:"new_listing"`
	Extensions         []json.RawMessage `json:"extensions"`
}

func (o AccountWhitelistOperation) Type() OperationType { return OperationTypeAccountWhitelist }

func (o AccountWhitelistOperation) MarshalJSON() ([]byte, error) {
	type alias AccountWhitelistOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *AccountWhitelistOperation) UnmarshalJSON(data []byte) error {
	type alias AccountWhitelistOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeAccountWhitelist, &payload); err != nil {
		return err
	}
	*o = AccountWhitelistOperation(payload)
	return nil
}

type AccountUpgradeOperation struct {
	Fee                     AssetAmount       `json:"fee"`
	AccountToUpgrade        ObjectID          `json:"account_to_upgrade"`
	UpgradeToLifetimeMember bool              `json:"upgrade_to_lifetime_member"`
	Extensions              []json.RawMessage `json:"extensions"`
}

func (o AccountUpgradeOperation) Type() OperationType { return OperationTypeAccountUpgrade }

func (o AccountUpgradeOperation) MarshalJSON() ([]byte, error) {
	type alias AccountUpgradeOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *AccountUpgradeOperation) UnmarshalJSON(data []byte) error {
	type alias AccountUpgradeOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeAccountUpgrade, &payload); err != nil {
		return err
	}
	*o = AccountUpgradeOperation(payload)
	return nil
}

type AccountTransferOperation struct {
	Fee        AssetAmount       `json:"fee"`
	AccountID  ObjectID          `json:"account_id"`
	NewOwner   ObjectID          `json:"new_owner"`
	Extensions []json.RawMessage `json:"extensions"`
}

func (o AccountTransferOperation) Type() OperationType { return OperationTypeAccountTransfer }

func (o AccountTransferOperation) MarshalJSON() ([]byte, error) {
	type alias AccountTransferOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *AccountTransferOperation) UnmarshalJSON(data []byte) error {
	type alias AccountTransferOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeAccountTransfer, &payload); err != nil {
		return err
	}
	*o = AccountTransferOperation(payload)
	return nil
}

// AssetCreateOperation creates a user-defined asset.
type AssetCreateOperation struct {
	Fee                AssetAmount       `json:"fee"`
	Issuer             ObjectID          `json:"issuer"`
	Symbol             string            `json:"symbol"`
	Precision          uint8             `json:"precision"`
	CommonOptions      AssetOptions      `json:"common_options"`
	BitassetOpts       *BitAssetOptions  `json:"bitasset_opts,omitempty"`
	IsPredictionMarket bool              `json:"is_prediction_market"`
	Extensions         []json.RawMessage `json:"extensions"`
}

func (o AssetCreateOperation) Type() OperationType { return OperationTypeAssetCreate }

func (o AssetCreateOperation) MarshalJSON() ([]byte, error) {
	type alias AssetCreateOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *AssetCreateOperation) UnmarshalJSON(data []byte) error {
	type alias AssetCreateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeAssetCreate, &payload); err != nil {
		return err
	}
	*o = AssetCreateOperation(payload)
	return nil
}

type AssetUpdateOperation struct {
	Fee           AssetAmount           `json:"fee"`
	Issuer        ObjectID              `json:"issuer"`
	AssetToUpdate ObjectID              `json:"asset_to_update"`
	NewIssuer     *ObjectID             `json:"new_issuer,omitempty"`
	NewOptions    AssetOptions          `json:"new_options"`
	Extensions    AssetUpdateExtensions `json:"extensions"`
}

func (o AssetUpdateOperation) Type() OperationType { return OperationTypeAssetUpdate }

func (o AssetUpdateOperation) MarshalJSON() ([]byte, error) {
	type alias AssetUpdateOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *AssetUpdateOperation) UnmarshalJSON(data []byte) error {
	type alias AssetUpdateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeAssetUpdate, &payload); err != nil {
		return err
	}
	*o = AssetUpdateOperation(payload)
	return nil
}

type AssetUpdateBitassetOperation struct {
	Fee           AssetAmount       `json:"fee"`
	Issuer        ObjectID          `json:"issuer"`
	AssetToUpdate ObjectID          `json:"asset_to_update"`
	NewOptions    BitAssetOptions   `json:"new_options"`
	Extensions    []json.RawMessage `json:"extensions"`
}

func (o AssetUpdateBitassetOperation) Type() OperationType { return OperationTypeAssetUpdateBitasset }

func (o AssetUpdateBitassetOperation) MarshalJSON() ([]byte, error) {
	type alias AssetUpdateBitassetOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *AssetUpdateBitassetOperation) UnmarshalJSON(data []byte) error {
	type alias AssetUpdateBitassetOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeAssetUpdateBitasset, &payload); err != nil {
		return err
	}
	*o = AssetUpdateBitassetOperation(payload)
	return nil
}

type AssetUpdateFeedProducersOperation struct {
	Fee              AssetAmount       `json:"fee"`
	Issuer           ObjectID          `json:"issuer"`
	AssetToUpdate    ObjectID          `json:"asset_to_update"`
	NewFeedProducers []ObjectID        `json:"new_feed_producers"`
	Extensions       []json.RawMessage `json:"extensions"`
}

func (o AssetUpdateFeedProducersOperation) Type() OperationType {
	return OperationTypeAssetUpdateFeedProducers
}

func (o AssetUpdateFeedProducersOperation) MarshalJSON() ([]byte, error) {
	type alias AssetUpdateFeedProducersOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *AssetUpdateFeedProducersOperation) UnmarshalJSON(data []byte) error {
	type alias AssetUpdateFeedProducersOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeAssetUpdateFeedProducers, &payload); err != nil {
		return err
	}
	*o = AssetUpdateFeedProducersOperation(payload)
	return nil
}

type AssetFundFeePoolOperation struct {
	Fee         AssetAmount       `json:"fee"`
	FromAccount ObjectID          `json:"from_account"`
	AssetID     ObjectID          `json:"asset_id"`
	Amount      int64             `json:"amount"`
	Extensions  []json.RawMessage `json:"extensions"`
}

func (o AssetFundFeePoolOperation) Type() OperationType { return OperationTypeAssetFundFeePool }

func (o AssetFundFeePoolOperation) MarshalJSON() ([]byte, error) {
	type alias AssetFundFeePoolOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *AssetFundFeePoolOperation) UnmarshalJSON(data []byte) error {
	type alias AssetFundFeePoolOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeAssetFundFeePool, &payload); err != nil {
		return err
	}
	*o = AssetFundFeePoolOperation(payload)
	return nil
}

type AssetSettleOperation struct {
	Fee        AssetAmount       `json:"fee"`
	Account    ObjectID          `json:"account"`
	Amount     AssetAmount       `json:"amount"`
	Extensions []json.RawMessage `json:"extensions"`
}

func (o AssetSettleOperation) Type() OperationType { return OperationTypeAssetSettle }

func (o AssetSettleOperation) MarshalJSON() ([]byte, error) {
	type alias AssetSettleOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *AssetSettleOperation) UnmarshalJSON(data []byte) error {
	type alias AssetSettleOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeAssetSettle, &payload); err != nil {
		return err
	}
	*o = AssetSettleOperation(payload)
	return nil
}

type AssetGlobalSettleOperation struct {
	Fee           AssetAmount       `json:"fee"`
	Issuer        ObjectID          `json:"issuer"`
	AssetToSettle ObjectID          `json:"asset_to_settle"`
	SettlePrice   Price             `json:"settle_price"`
	Extensions    []json.RawMessage `json:"extensions"`
}

func (o AssetGlobalSettleOperation) Type() OperationType { return OperationTypeAssetGlobalSettle }

func (o AssetGlobalSettleOperation) MarshalJSON() ([]byte, error) {
	type alias AssetGlobalSettleOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *AssetGlobalSettleOperation) UnmarshalJSON(data []byte) error {
	type alias AssetGlobalSettleOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeAssetGlobalSettle, &payload); err != nil {
		return err
	}
	*o = AssetGlobalSettleOperation(payload)
	return nil
}

type AssetPublishFeedOperation struct {
	Fee        AssetAmount                `json:"fee"`
	Publisher  ObjectID                   `json:"publisher"`
	AssetID    ObjectID                   `json:"asset_id"`
	Feed       PriceFeed                  `json:"feed"`
	Extensions AssetPublishFeedExtensions `json:"extensions"`
}

func (o AssetPublishFeedOperation) Type() OperationType { return OperationTypeAssetPublishFeed }

func (o AssetPublishFeedOperation) MarshalJSON() ([]byte, error) {
	type alias AssetPublishFeedOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *AssetPublishFeedOperation) UnmarshalJSON(data []byte) error {
	type alias AssetPublishFeedOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeAssetPublishFeed, &payload); err != nil {
		return err
	}
	*o = AssetPublishFeedOperation(payload)
	return nil
}

type WitnessCreateOperation struct {
	Fee             AssetAmount `json:"fee"`
	WitnessAccount  ObjectID    `json:"witness_account"`
	URL             string      `json:"url"`
	BlockSigningKey PublicKey   `json:"block_signing_key"`
}

func (o WitnessCreateOperation) Type() OperationType { return OperationTypeWitnessCreate }

func (o WitnessCreateOperation) MarshalJSON() ([]byte, error) {
	type alias WitnessCreateOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *WitnessCreateOperation) UnmarshalJSON(data []byte) error {
	type alias WitnessCreateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeWitnessCreate, &payload); err != nil {
		return err
	}
	*o = WitnessCreateOperation(payload)
	return nil
}

type WitnessUpdateOperation struct {
	Fee            AssetAmount `json:"fee"`
	Witness        ObjectID    `json:"witness"`
	WitnessAccount ObjectID    `json:"witness_account"`
	NewURL         *string     `json:"new_url,omitempty"`
	NewSigningKey  *PublicKey  `json:"new_signing_key,omitempty"`
}

func (o WitnessUpdateOperation) Type() OperationType { return OperationTypeWitnessUpdate }

func (o WitnessUpdateOperation) MarshalJSON() ([]byte, error) {
	type alias WitnessUpdateOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *WitnessUpdateOperation) UnmarshalJSON(data []byte) error {
	type alias WitnessUpdateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeWitnessUpdate, &payload); err != nil {
		return err
	}
	*o = WitnessUpdateOperation(payload)
	return nil
}

// OpWrapper mirrors the BitShares op_wrapper serializer.
type OpWrapper struct {
	Op OperationEnvelope `json:"op"`
}

func (o OpWrapper) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Op OperationEnvelope `json:"op"`
	}{Op: o.Op})
}

func (o *OpWrapper) UnmarshalJSON(data []byte) error {
	var payload struct {
		Op OperationEnvelope `json:"op"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	o.Op = payload.Op
	return nil
}

type ProposalCreateOperation struct {
	Fee                 AssetAmount       `json:"fee"`
	FeePayingAccount    ObjectID          `json:"fee_paying_account"`
	ExpirationTime      Time              `json:"expiration_time"`
	ProposedOps         []OpWrapper       `json:"proposed_ops"`
	ReviewPeriodSeconds *uint32           `json:"review_period_seconds,omitempty"`
	Extensions          []json.RawMessage `json:"extensions"`
}

func (o ProposalCreateOperation) Type() OperationType { return OperationTypeProposalCreate }

func (o ProposalCreateOperation) MarshalJSON() ([]byte, error) {
	type alias ProposalCreateOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *ProposalCreateOperation) UnmarshalJSON(data []byte) error {
	type alias ProposalCreateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeProposalCreate, &payload); err != nil {
		return err
	}
	*o = ProposalCreateOperation(payload)
	return nil
}

type ProposalUpdateOperation struct {
	Fee                     AssetAmount       `json:"fee"`
	FeePayingAccount        ObjectID          `json:"fee_paying_account"`
	Proposal                ObjectID          `json:"proposal"`
	ActiveApprovalsToAdd    []ObjectID        `json:"active_approvals_to_add,omitempty"`
	ActiveApprovalsToRemove []ObjectID        `json:"active_approvals_to_remove,omitempty"`
	OwnerApprovalsToAdd     []ObjectID        `json:"owner_approvals_to_add,omitempty"`
	OwnerApprovalsToRemove  []ObjectID        `json:"owner_approvals_to_remove,omitempty"`
	KeyApprovalsToAdd       []PublicKey       `json:"key_approvals_to_add,omitempty"`
	KeyApprovalsToRemove    []PublicKey       `json:"key_approvals_to_remove,omitempty"`
	Extensions              []json.RawMessage `json:"extensions"`
}

func (o ProposalUpdateOperation) Type() OperationType { return OperationTypeProposalUpdate }

func (o ProposalUpdateOperation) MarshalJSON() ([]byte, error) {
	type alias ProposalUpdateOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *ProposalUpdateOperation) UnmarshalJSON(data []byte) error {
	type alias ProposalUpdateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeProposalUpdate, &payload); err != nil {
		return err
	}
	*o = ProposalUpdateOperation(payload)
	return nil
}

type ProposalDeleteOperation struct {
	Fee                 AssetAmount       `json:"fee"`
	FeePayingAccount    ObjectID          `json:"fee_paying_account"`
	UsingOwnerAuthority bool              `json:"using_owner_authority"`
	Proposal            ObjectID          `json:"proposal"`
	Extensions          []json.RawMessage `json:"extensions"`
}

func (o ProposalDeleteOperation) Type() OperationType { return OperationTypeProposalDelete }

func (o ProposalDeleteOperation) MarshalJSON() ([]byte, error) {
	type alias ProposalDeleteOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *ProposalDeleteOperation) UnmarshalJSON(data []byte) error {
	type alias ProposalDeleteOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeProposalDelete, &payload); err != nil {
		return err
	}
	*o = ProposalDeleteOperation(payload)
	return nil
}

type WithdrawPermissionCreateOperation struct {
	Fee                    AssetAmount `json:"fee"`
	WithdrawFromAccount    ObjectID    `json:"withdraw_from_account"`
	AuthorizedAccount      ObjectID    `json:"authorized_account"`
	WithdrawalLimit        AssetAmount `json:"withdrawal_limit"`
	WithdrawalPeriodSec    uint32      `json:"withdrawal_period_sec"`
	PeriodsUntilExpiration uint32      `json:"periods_until_expiration"`
	PeriodStartTime        Time        `json:"period_start_time"`
}

func (o WithdrawPermissionCreateOperation) Type() OperationType {
	return OperationTypeWithdrawPermissionCreate
}

func (o WithdrawPermissionCreateOperation) MarshalJSON() ([]byte, error) {
	type alias WithdrawPermissionCreateOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *WithdrawPermissionCreateOperation) UnmarshalJSON(data []byte) error {
	type alias WithdrawPermissionCreateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeWithdrawPermissionCreate, &payload); err != nil {
		return err
	}
	*o = WithdrawPermissionCreateOperation(payload)
	return nil
}

type WithdrawPermissionUpdateOperation struct {
	Fee                    AssetAmount `json:"fee"`
	WithdrawFromAccount    ObjectID    `json:"withdraw_from_account"`
	AuthorizedAccount      ObjectID    `json:"authorized_account"`
	PermissionToUpdate     ObjectID    `json:"permission_to_update"`
	WithdrawalLimit        AssetAmount `json:"withdrawal_limit"`
	WithdrawalPeriodSec    uint32      `json:"withdrawal_period_sec"`
	PeriodStartTime        Time        `json:"period_start_time"`
	PeriodsUntilExpiration uint32      `json:"periods_until_expiration"`
}

func (o WithdrawPermissionUpdateOperation) Type() OperationType {
	return OperationTypeWithdrawPermissionUpdate
}

func (o WithdrawPermissionUpdateOperation) MarshalJSON() ([]byte, error) {
	type alias WithdrawPermissionUpdateOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *WithdrawPermissionUpdateOperation) UnmarshalJSON(data []byte) error {
	type alias WithdrawPermissionUpdateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeWithdrawPermissionUpdate, &payload); err != nil {
		return err
	}
	*o = WithdrawPermissionUpdateOperation(payload)
	return nil
}

type WithdrawPermissionClaimOperation struct {
	Fee                 AssetAmount     `json:"fee"`
	WithdrawPermission  ObjectID        `json:"withdraw_permission"`
	WithdrawFromAccount ObjectID        `json:"withdraw_from_account"`
	WithdrawToAccount   ObjectID        `json:"withdraw_to_account"`
	AmountToWithdraw    AssetAmount     `json:"amount_to_withdraw"`
	Memo                json.RawMessage `json:"memo,omitempty"`
}

func (o WithdrawPermissionClaimOperation) Type() OperationType {
	return OperationTypeWithdrawPermissionClaim
}

func (o WithdrawPermissionClaimOperation) MarshalJSON() ([]byte, error) {
	type alias WithdrawPermissionClaimOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *WithdrawPermissionClaimOperation) UnmarshalJSON(data []byte) error {
	type alias WithdrawPermissionClaimOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeWithdrawPermissionClaim, &payload); err != nil {
		return err
	}
	*o = WithdrawPermissionClaimOperation(payload)
	return nil
}

type WithdrawPermissionDeleteOperation struct {
	Fee                  AssetAmount `json:"fee"`
	WithdrawFromAccount  ObjectID    `json:"withdraw_from_account"`
	AuthorizedAccount    ObjectID    `json:"authorized_account"`
	WithdrawalPermission ObjectID    `json:"withdrawal_permission"`
}

func (o WithdrawPermissionDeleteOperation) Type() OperationType {
	return OperationTypeWithdrawPermissionDelete
}

func (o WithdrawPermissionDeleteOperation) MarshalJSON() ([]byte, error) {
	type alias WithdrawPermissionDeleteOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *WithdrawPermissionDeleteOperation) UnmarshalJSON(data []byte) error {
	type alias WithdrawPermissionDeleteOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeWithdrawPermissionDelete, &payload); err != nil {
		return err
	}
	*o = WithdrawPermissionDeleteOperation(payload)
	return nil
}

type CommitteeMemberCreateOperation struct {
	Fee                    AssetAmount `json:"fee"`
	CommitteeMemberAccount ObjectID    `json:"committee_member_account"`
	URL                    string      `json:"url"`
}

func (o CommitteeMemberCreateOperation) Type() OperationType {
	return OperationTypeCommitteeMemberCreate
}

func (o CommitteeMemberCreateOperation) MarshalJSON() ([]byte, error) {
	type alias CommitteeMemberCreateOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *CommitteeMemberCreateOperation) UnmarshalJSON(data []byte) error {
	type alias CommitteeMemberCreateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeCommitteeMemberCreate, &payload); err != nil {
		return err
	}
	*o = CommitteeMemberCreateOperation(payload)
	return nil
}

type CommitteeMemberUpdateOperation struct {
	Fee                    AssetAmount `json:"fee"`
	CommitteeMember        ObjectID    `json:"committee_member"`
	CommitteeMemberAccount ObjectID    `json:"committee_member_account"`
	NewURL                 *string     `json:"new_url,omitempty"`
}

func (o CommitteeMemberUpdateOperation) Type() OperationType {
	return OperationTypeCommitteeMemberUpdate
}

func (o CommitteeMemberUpdateOperation) MarshalJSON() ([]byte, error) {
	type alias CommitteeMemberUpdateOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *CommitteeMemberUpdateOperation) UnmarshalJSON(data []byte) error {
	type alias CommitteeMemberUpdateOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeCommitteeMemberUpdate, &payload); err != nil {
		return err
	}
	*o = CommitteeMemberUpdateOperation(payload)
	return nil
}

type CommitteeMemberUpdateGlobalParametersOperation struct {
	Fee           AssetAmount     `json:"fee"`
	NewParameters ChainParameters `json:"new_parameters"`
}

func (o CommitteeMemberUpdateGlobalParametersOperation) Type() OperationType {
	return OperationTypeCommitteeMemberUpdateGlobalParameters
}

func (o CommitteeMemberUpdateGlobalParametersOperation) MarshalJSON() ([]byte, error) {
	type alias CommitteeMemberUpdateGlobalParametersOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *CommitteeMemberUpdateGlobalParametersOperation) UnmarshalJSON(data []byte) error {
	type alias CommitteeMemberUpdateGlobalParametersOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeCommitteeMemberUpdateGlobalParameters, &payload); err != nil {
		return err
	}
	*o = CommitteeMemberUpdateGlobalParametersOperation(payload)
	return nil
}

type BalanceClaimOperation struct {
	Fee              AssetAmount `json:"fee"`
	DepositToAccount ObjectID    `json:"deposit_to_account"`
	BalanceToClaim   ObjectID    `json:"balance_to_claim"`
	BalanceOwnerKey  PublicKey   `json:"balance_owner_key"`
	TotalClaimed     AssetAmount `json:"total_claimed"`
}

func (o BalanceClaimOperation) Type() OperationType { return OperationTypeBalanceClaim }

func (o BalanceClaimOperation) MarshalJSON() ([]byte, error) {
	type alias BalanceClaimOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *BalanceClaimOperation) UnmarshalJSON(data []byte) error {
	type alias BalanceClaimOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeBalanceClaim, &payload); err != nil {
		return err
	}
	*o = BalanceClaimOperation(payload)
	return nil
}

type OverrideTransferOperation struct {
	Fee        AssetAmount       `json:"fee"`
	Issuer     ObjectID          `json:"issuer"`
	From       ObjectID          `json:"from"`
	To         ObjectID          `json:"to"`
	Amount     AssetAmount       `json:"amount"`
	Memo       json.RawMessage   `json:"memo,omitempty"`
	Extensions []json.RawMessage `json:"extensions"`
}

func (o OverrideTransferOperation) Type() OperationType { return OperationTypeOverrideTransfer }

func (o OverrideTransferOperation) MarshalJSON() ([]byte, error) {
	type alias OverrideTransferOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *OverrideTransferOperation) UnmarshalJSON(data []byte) error {
	type alias OverrideTransferOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeOverrideTransfer, &payload); err != nil {
		return err
	}
	*o = OverrideTransferOperation(payload)
	return nil
}

type AssetSettleCancelOperation struct {
	Fee        AssetAmount       `json:"fee"`
	Settlement ObjectID          `json:"settlement"`
	Account    ObjectID          `json:"account"`
	Amount     AssetAmount       `json:"amount"`
	Extensions []json.RawMessage `json:"extensions"`
}

func (o AssetSettleCancelOperation) Type() OperationType { return OperationTypeAssetSettleCancel }

func (o AssetSettleCancelOperation) MarshalJSON() ([]byte, error) {
	type alias AssetSettleCancelOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *AssetSettleCancelOperation) UnmarshalJSON(data []byte) error {
	type alias AssetSettleCancelOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeAssetSettleCancel, &payload); err != nil {
		return err
	}
	*o = AssetSettleCancelOperation(payload)
	return nil
}

type AssetClaimFeesOperation struct {
	Fee           AssetAmount               `json:"fee"`
	Issuer        ObjectID                  `json:"issuer"`
	AmountToClaim AssetAmount               `json:"amount_to_claim"`
	Extensions    *AssetClaimFeesExtensions `json:"extensions,omitempty"`
}

func (o AssetClaimFeesOperation) Type() OperationType { return OperationTypeAssetClaimFees }

func (o AssetClaimFeesOperation) MarshalJSON() ([]byte, error) {
	type alias AssetClaimFeesOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *AssetClaimFeesOperation) UnmarshalJSON(data []byte) error {
	type alias AssetClaimFeesOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeAssetClaimFees, &payload); err != nil {
		return err
	}
	*o = AssetClaimFeesOperation(payload)
	return nil
}

type FBADistributeOperation struct {
	Fee       AssetAmount `json:"fee"`
	AccountID ObjectID    `json:"account_id"`
	FBAID     ObjectID    `json:"fba_id"`
	Amount    int64       `json:"amount"`
}

func (o FBADistributeOperation) Type() OperationType { return OperationTypeFBADistribute }

func (o FBADistributeOperation) MarshalJSON() ([]byte, error) {
	type alias FBADistributeOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *FBADistributeOperation) UnmarshalJSON(data []byte) error {
	type alias FBADistributeOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeFBADistribute, &payload); err != nil {
		return err
	}
	*o = FBADistributeOperation(payload)
	return nil
}

type BidCollateralOperation struct {
	Fee                  AssetAmount       `json:"fee"`
	Bidder               ObjectID          `json:"bidder"`
	AdditionalCollateral AssetAmount       `json:"additional_collateral"`
	DebtCovered          AssetAmount       `json:"debt_covered"`
	Extensions           []json.RawMessage `json:"extensions"`
}

func (o BidCollateralOperation) Type() OperationType { return OperationTypeBidCollateral }

func (o BidCollateralOperation) MarshalJSON() ([]byte, error) {
	type alias BidCollateralOperation
	if o.Extensions == nil {
		o.Extensions = []json.RawMessage{}
	}
	return marshalOperation(o.Type(), alias(o))
}

func (o *BidCollateralOperation) UnmarshalJSON(data []byte) error {
	type alias BidCollateralOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeBidCollateral, &payload); err != nil {
		return err
	}
	*o = BidCollateralOperation(payload)
	return nil
}

type ExecuteBidOperation struct {
	Fee        AssetAmount `json:"fee"`
	Bidder     ObjectID    `json:"bidder"`
	Debt       AssetAmount `json:"debt"`
	Collateral AssetAmount `json:"collateral"`
}

func (o ExecuteBidOperation) Type() OperationType { return OperationTypeExecuteBid }

func (o ExecuteBidOperation) MarshalJSON() ([]byte, error) {
	type alias ExecuteBidOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *ExecuteBidOperation) UnmarshalJSON(data []byte) error {
	type alias ExecuteBidOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeExecuteBid, &payload); err != nil {
		return err
	}
	*o = ExecuteBidOperation(payload)
	return nil
}

func init() {
	RegisterOperationFactory(OperationTypeCallOrderUpdate, func() Operation { return &CallOrderUpdateOperation{} })
	RegisterOperationFactory(OperationTypeFillOrder, func() Operation { return &FillOrderOperation{} })
	RegisterOperationFactory(OperationTypeAccountCreate, func() Operation { return &AccountCreateOperation{} })
	RegisterOperationFactory(OperationTypeAccountUpdate, func() Operation { return &AccountUpdateOperation{} })
	RegisterOperationFactory(OperationTypeAccountWhitelist, func() Operation { return &AccountWhitelistOperation{} })
	RegisterOperationFactory(OperationTypeAccountUpgrade, func() Operation { return &AccountUpgradeOperation{} })
	RegisterOperationFactory(OperationTypeAccountTransfer, func() Operation { return &AccountTransferOperation{} })
	RegisterOperationFactory(OperationTypeAssetCreate, func() Operation { return &AssetCreateOperation{} })
	RegisterOperationFactory(OperationTypeAssetUpdate, func() Operation { return &AssetUpdateOperation{} })
	RegisterOperationFactory(OperationTypeAssetUpdateBitasset, func() Operation { return &AssetUpdateBitassetOperation{} })
	RegisterOperationFactory(OperationTypeAssetUpdateFeedProducers, func() Operation { return &AssetUpdateFeedProducersOperation{} })
	RegisterOperationFactory(OperationTypeAssetFundFeePool, func() Operation { return &AssetFundFeePoolOperation{} })
	RegisterOperationFactory(OperationTypeAssetSettle, func() Operation { return &AssetSettleOperation{} })
	RegisterOperationFactory(OperationTypeAssetGlobalSettle, func() Operation { return &AssetGlobalSettleOperation{} })
	RegisterOperationFactory(OperationTypeAssetPublishFeed, func() Operation { return &AssetPublishFeedOperation{} })
	RegisterOperationFactory(OperationTypeWitnessCreate, func() Operation { return &WitnessCreateOperation{} })
	RegisterOperationFactory(OperationTypeWitnessUpdate, func() Operation { return &WitnessUpdateOperation{} })
	RegisterOperationFactory(OperationTypeProposalCreate, func() Operation { return &ProposalCreateOperation{} })
	RegisterOperationFactory(OperationTypeProposalUpdate, func() Operation { return &ProposalUpdateOperation{} })
	RegisterOperationFactory(OperationTypeProposalDelete, func() Operation { return &ProposalDeleteOperation{} })
	RegisterOperationFactory(OperationTypeWithdrawPermissionCreate, func() Operation { return &WithdrawPermissionCreateOperation{} })
	RegisterOperationFactory(OperationTypeWithdrawPermissionUpdate, func() Operation { return &WithdrawPermissionUpdateOperation{} })
	RegisterOperationFactory(OperationTypeWithdrawPermissionClaim, func() Operation { return &WithdrawPermissionClaimOperation{} })
	RegisterOperationFactory(OperationTypeWithdrawPermissionDelete, func() Operation { return &WithdrawPermissionDeleteOperation{} })
	RegisterOperationFactory(OperationTypeCommitteeMemberCreate, func() Operation { return &CommitteeMemberCreateOperation{} })
	RegisterOperationFactory(OperationTypeCommitteeMemberUpdate, func() Operation { return &CommitteeMemberUpdateOperation{} })
	RegisterOperationFactory(OperationTypeCommitteeMemberUpdateGlobalParameters, func() Operation { return &CommitteeMemberUpdateGlobalParametersOperation{} })
	RegisterOperationFactory(OperationTypeBalanceClaim, func() Operation { return &BalanceClaimOperation{} })
	RegisterOperationFactory(OperationTypeOverrideTransfer, func() Operation { return &OverrideTransferOperation{} })
	RegisterOperationFactory(OperationTypeAssetSettleCancel, func() Operation { return &AssetSettleCancelOperation{} })
	RegisterOperationFactory(OperationTypeAssetClaimFees, func() Operation { return &AssetClaimFeesOperation{} })
	RegisterOperationFactory(OperationTypeFBADistribute, func() Operation { return &FBADistributeOperation{} })
	RegisterOperationFactory(OperationTypeBidCollateral, func() Operation { return &BidCollateralOperation{} })
	RegisterOperationFactory(OperationTypeExecuteBid, func() Operation { return &ExecuteBidOperation{} })
}

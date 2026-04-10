package protocol

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// MarshalBinary serializes a BitShares transaction into the protocol wire format.
func (tx *Transaction) MarshalBinary() ([]byte, error) {
	if tx == nil {
		return nil, fmt.Errorf("nil transaction")
	}

	w := newBinaryWriter()
	w.writeUint16(tx.RefBlockNum)
	w.writeUint32(tx.RefBlockPrefix)
	w.writeUint32(tx.Expiration.UnixSeconds())
	w.writeVarUint64(uint64(len(tx.Operations)))
	for _, op := range tx.Operations {
		raw, err := op.MarshalBinary()
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(raw); err != nil {
			return nil, err
		}
	}
	if !extensionsEmpty(tx.Extensions) {
		return nil, fmt.Errorf("transaction extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return w.Bytes(), nil
}

// UnmarshalBinary decodes a BitShares transaction from the protocol wire format.
func (tx *Transaction) UnmarshalBinary(data []byte) error {
	if tx == nil {
		return fmt.Errorf("nil transaction")
	}

	r := newBinaryReader(data)
	refBlockNum, err := r.readUint16()
	if err != nil {
		return err
	}
	refBlockPrefix, err := r.readUint32()
	if err != nil {
		return err
	}
	expiration, err := r.readUint32()
	if err != nil {
		return err
	}
	count, err := r.readVarUint64()
	if err != nil {
		return err
	}

	ops := make([]OperationEnvelope, 0, count)
	for i := uint64(0); i < count; i++ {
		op, err := readOperationEnvelope(r)
		if err != nil {
			return err
		}
		ops = append(ops, op)
	}

	extCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("transaction extensions are not supported in binary serialization")
	}

	tx.RefBlockNum = refBlockNum
	tx.RefBlockPrefix = refBlockPrefix
	tx.Expiration = NewTimeFromUnix(expiration)
	tx.Operations = ops
	tx.Extensions = nil
	return nil
}

// MarshalBinary serializes a signed transaction into the protocol wire format.
func (tx *SignedTransaction) MarshalBinary() ([]byte, error) {
	if tx == nil {
		return nil, fmt.Errorf("nil signed transaction")
	}

	raw, err := tx.Transaction.MarshalBinary()
	if err != nil {
		return nil, err
	}

	w := newBinaryWriter()
	if _, err := w.Write(raw); err != nil {
		return nil, err
	}
	w.writeVarUint64(uint64(len(tx.Signatures)))
	for _, sig := range tx.Signatures {
		rawSig, err := hex.DecodeString(strings.TrimSpace(sig))
		if err != nil {
			return nil, err
		}
		if len(rawSig) != 65 {
			return nil, fmt.Errorf("invalid signature length %d", len(rawSig))
		}
		if _, err := w.Write(rawSig); err != nil {
			return nil, err
		}
	}
	return w.Bytes(), nil
}

// SigningPayload returns the chain-id-prefixed payload that must be hashed before signing.
func (tx *Transaction) SigningPayload(chainID string) ([]byte, error) {
	if tx == nil {
		return nil, fmt.Errorf("nil transaction")
	}

	chainID = strings.TrimSpace(chainID)
	if chainID == "" {
		return nil, fmt.Errorf("empty chain id")
	}
	chainBytes, err := hex.DecodeString(chainID)
	if err != nil {
		return nil, err
	}
	if len(chainBytes) != 32 {
		return nil, fmt.Errorf("invalid chain id length %d", len(chainBytes))
	}
	txBytes, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return append(chainBytes, txBytes...), nil
}

// SigningDigest hashes the transaction signing payload with SHA-256.
func (tx *Transaction) SigningDigest(chainID string) ([]byte, error) {
	payload, err := tx.SigningPayload(chainID)
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(payload)
	out := make([]byte, len(sum))
	copy(out, sum[:])
	return out, nil
}

// MarshalBinary serializes an operation envelope into the BitShares static-variant wire format.
func (e OperationEnvelope) MarshalBinary() ([]byte, error) {
	if e.Operation == nil {
		return nil, fmt.Errorf("nil operation")
	}

	switch op := e.Operation.(type) {
	case interface{ MarshalBinary() ([]byte, error) }:
		return op.MarshalBinary()
	default:
		return nil, fmt.Errorf("operation type %T does not support binary serialization", e.Operation)
	}
}

// UnmarshalBinary decodes an operation envelope from the BitShares static-variant wire format.
func (e *OperationEnvelope) UnmarshalBinary(data []byte) error {
	if e == nil {
		return fmt.Errorf("nil operation envelope")
	}
	r := newBinaryReader(data)
	op, err := readOperationEnvelope(r)
	if err != nil {
		return err
	}
	e.Operation = op.Operation
	return nil
}

func readOperationEnvelope(r *binaryReader) (OperationEnvelope, error) {
	var out OperationEnvelope
	kind, err := r.readVarUint64()
	if err != nil {
		return out, err
	}

	switch OperationType(kind) {
	case OperationTypeTransfer:
		var op TransferOperation
		if err := op.unmarshalBinaryBody(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeLimitOrderCreate:
		var op LimitOrderCreateOperation
		if err := op.unmarshalBinaryBody(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeLimitOrderCancel:
		var op LimitOrderCancelOperation
		if err := op.unmarshalBinaryBody(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeLimitOrderUpdate:
		var op LimitOrderUpdateOperation
		if err := op.unmarshalBinaryBody(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeCallOrderUpdate:
		var op CallOrderUpdateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeFillOrder:
		var op FillOrderOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAccountCreate:
		var op AccountCreateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAccountUpdate:
		var op AccountUpdateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAccountWhitelist:
		var op AccountWhitelistOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAccountUpgrade:
		var op AccountUpgradeOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAccountTransfer:
		var op AccountTransferOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAssetCreate:
		var op AssetCreateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAssetUpdate:
		var op AssetUpdateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAssetUpdateBitasset:
		var op AssetUpdateBitassetOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAssetUpdateFeedProducers:
		var op AssetUpdateFeedProducersOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAssetFundFeePool:
		var op AssetFundFeePoolOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAssetSettle:
		var op AssetSettleOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAssetGlobalSettle:
		var op AssetGlobalSettleOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAssetPublishFeed:
		var op AssetPublishFeedOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeWitnessCreate:
		var op WitnessCreateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeWitnessUpdate:
		var op WitnessUpdateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeProposalCreate:
		var op ProposalCreateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeProposalUpdate:
		var op ProposalUpdateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeProposalDelete:
		var op ProposalDeleteOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeWithdrawPermissionCreate:
		var op WithdrawPermissionCreateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeWithdrawPermissionUpdate:
		var op WithdrawPermissionUpdateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeWithdrawPermissionClaim:
		var op WithdrawPermissionClaimOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeWithdrawPermissionDelete:
		var op WithdrawPermissionDeleteOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeCommitteeMemberCreate:
		var op CommitteeMemberCreateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeCommitteeMemberUpdate:
		var op CommitteeMemberUpdateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeCommitteeMemberUpdateGlobalParameters:
		var op CommitteeMemberUpdateGlobalParametersOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeVestingBalanceCreate:
		var op VestingBalanceCreateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeVestingBalanceWithdraw:
		var op VestingBalanceWithdrawOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeWorkerCreate:
		var op WorkerCreateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeCustom:
		var op CustomOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAssert:
		var op AssertOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeBalanceClaim:
		var op BalanceClaimOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeOverrideTransfer:
		var op OverrideTransferOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeTransferToBlind:
		var op TransferToBlindOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeBlindTransfer:
		var op BlindTransferOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeTransferFromBlind:
		var op TransferFromBlindOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAssetSettleCancel:
		var op AssetSettleCancelOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAssetClaimFees:
		var op AssetClaimFeesOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeFBADistribute:
		var op FBADistributeOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeBidCollateral:
		var op BidCollateralOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeExecuteBid:
		var op ExecuteBidOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAssetUpdateIssuer:
		var op AssetUpdateIssuerOperation
		if err := op.unmarshalBinaryBody(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAssetClaimPool:
		var op AssetClaimPoolOperation
		if err := op.unmarshalBinaryBody(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAssetIssue:
		var op AssetIssueOperation
		if err := op.unmarshalBinaryBody(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeAssetReserve:
		var op AssetReserveOperation
		if err := op.unmarshalBinaryBody(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeHTLCCreate:
		var op HTLCCreateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeHTLCRedeem:
		var op HTLCRedeemOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeHTLCRedeemed:
		var op HTLCRedeemedOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeHTLCExtend:
		var op HTLCExtendOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeHTLCRefund:
		var op HTLCRefundOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeCustomAuthorityCreate:
		var op CustomAuthorityCreateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeCustomAuthorityUpdate:
		var op CustomAuthorityUpdateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeCustomAuthorityDelete:
		var op CustomAuthorityDeleteOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeTicketCreate:
		var op TicketCreateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeTicketUpdate:
		var op TicketUpdateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeLiquidityPoolCreate:
		var op LiquidityPoolCreateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeLiquidityPoolDelete:
		var op LiquidityPoolDeleteOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeLiquidityPoolDeposit:
		var op LiquidityPoolDepositOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeLiquidityPoolWithdraw:
		var op LiquidityPoolWithdrawOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeLiquidityPoolExchange:
		var op LiquidityPoolExchangeOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeSametFundCreate:
		var op SametFundCreateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeSametFundDelete:
		var op SametFundDeleteOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeSametFundUpdate:
		var op SametFundUpdateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeSametFundBorrow:
		var op SametFundBorrowOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeSametFundRepay:
		var op SametFundRepayOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeCreditOfferCreate:
		var op CreditOfferCreateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeCreditOfferDelete:
		var op CreditOfferDeleteOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeCreditOfferUpdate:
		var op CreditOfferUpdateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeCreditOfferAccept:
		var op CreditOfferAcceptOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeCreditDealRepay:
		var op CreditDealRepayOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeCreditDealExpired:
		var op CreditDealExpiredOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeLiquidityPoolUpdate:
		var op LiquidityPoolUpdateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	case OperationTypeCreditDealUpdate:
		var op CreditDealUpdateOperation
		if err := op.UnmarshalBinaryFrom(r); err != nil {
			return out, err
		}
		out.Operation = &op
	default:
		return out, fmt.Errorf("unsupported operation type %d", kind)
	}

	return out, nil
}

func extensionsEmpty(raw any) bool {
	switch v := raw.(type) {
	case nil:
		return true
	case []json.RawMessage:
		return len(v) == 0
	case json.RawMessage:
		trimmed := strings.TrimSpace(string(v))
		return trimmed == "" || trimmed == "{}" || trimmed == "[]" || trimmed == "null"
	default:
		return false
	}
}

func operationEnvelopeFromRawJSON(data []byte) (OperationEnvelope, error) {
	var body OperationBody
	if err := json.Unmarshal(data, &body); err != nil {
		return OperationEnvelope{}, err
	}
	if factory := newOperation(body.Kind); factory != nil {
		if err := json.Unmarshal(body.Payload, factory); err != nil {
			return OperationEnvelope{}, err
		}
		return OperationEnvelope{Operation: factory}, nil
	}
	return OperationEnvelope{Operation: &RawOperation{OperationBody: body}}, nil
}

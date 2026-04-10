package protocol

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type StealthConfirmation struct {
	OneTimeKey    PublicKey  `json:"one_time_key"`
	To            *PublicKey `json:"to,omitempty"`
	EncryptedMemo string     `json:"encrypted_memo"`
}

func (s StealthConfirmation) MarshalBinaryInto(w *binaryWriter) error {
	if err := s.OneTimeKey.MarshalBinaryInto(w); err != nil {
		return err
	}
	if err := writeOptionalPublicKey(w, s.To); err != nil {
		return err
	}
	raw, err := hex.DecodeString(strings.TrimSpace(s.EncryptedMemo))
	if err != nil {
		return err
	}
	w.writeBytes(raw)
	return nil
}

func (s *StealthConfirmation) UnmarshalBinaryFrom(r *binaryReader) error {
	key, err := readPublicKey(r)
	if err != nil {
		return err
	}
	to, err := readOptionalPublicKey(r)
	if err != nil {
		return err
	}
	memo, err := r.readBytes()
	if err != nil {
		return err
	}
	s.OneTimeKey = key
	s.To = to
	s.EncryptedMemo = hex.EncodeToString(memo)
	return nil
}

type BlindOutput struct {
	Commitment  string               `json:"commitment"`
	RangeProof  string               `json:"range_proof"`
	Owner       Authority            `json:"owner"`
	StealthMemo *StealthConfirmation `json:"stealth_memo,omitempty"`
}

func (b BlindOutput) MarshalBinaryInto(w *binaryWriter) error {
	commitment, err := hex.DecodeString(strings.TrimSpace(b.Commitment))
	if err != nil {
		return err
	}
	if len(commitment) != 33 {
		return fmt.Errorf("blind output commitment must be 33 bytes")
	}
	if err := writeFixedBytes(w, commitment, 33); err != nil {
		return err
	}
	proof, err := hex.DecodeString(strings.TrimSpace(b.RangeProof))
	if err != nil {
		return err
	}
	w.writeBytes(proof)
	if err := b.Owner.MarshalBinaryInto(w); err != nil {
		return err
	}
	if b.StealthMemo == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	return b.StealthMemo.MarshalBinaryInto(w)
}

func (b *BlindOutput) UnmarshalBinaryFrom(r *binaryReader) error {
	commitment, err := readFixedBytes(r, 33)
	if err != nil {
		return err
	}
	proof, err := r.readBytes()
	if err != nil {
		return err
	}
	var owner Authority
	if err := owner.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	present, err := r.readUint8()
	if err != nil {
		return err
	}
	var memo *StealthConfirmation
	if present != 0 {
		var value StealthConfirmation
		if err := value.UnmarshalBinaryFrom(r); err != nil {
			return err
		}
		memo = &value
	}
	b.Commitment = hex.EncodeToString(commitment)
	b.RangeProof = hex.EncodeToString(proof)
	b.Owner = owner
	b.StealthMemo = memo
	return nil
}

type BlindInput struct {
	Commitment string    `json:"commitment"`
	Owner      Authority `json:"owner"`
}

func (b BlindInput) MarshalBinaryInto(w *binaryWriter) error {
	commitment, err := hex.DecodeString(strings.TrimSpace(b.Commitment))
	if err != nil {
		return err
	}
	if len(commitment) != 33 {
		return fmt.Errorf("blind input commitment must be 33 bytes")
	}
	if err := writeFixedBytes(w, commitment, 33); err != nil {
		return err
	}
	return b.Owner.MarshalBinaryInto(w)
}

func (b *BlindInput) UnmarshalBinaryFrom(r *binaryReader) error {
	commitment, err := readFixedBytes(r, 33)
	if err != nil {
		return err
	}
	var owner Authority
	if err := owner.UnmarshalBinaryFrom(r); err != nil {
		return err
	}
	b.Commitment = hex.EncodeToString(commitment)
	b.Owner = owner
	return nil
}

type TransferToBlindOperation struct {
	Fee            AssetAmount   `json:"fee"`
	Amount         AssetAmount   `json:"amount"`
	From           ObjectID      `json:"from"`
	BlindingFactor string        `json:"blinding_factor"`
	Outputs        []BlindOutput `json:"outputs"`
}

func (o TransferToBlindOperation) Type() OperationType { return OperationTypeTransferToBlind }

func (o TransferToBlindOperation) MarshalJSON() ([]byte, error) {
	type alias TransferToBlindOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *TransferToBlindOperation) UnmarshalJSON(data []byte) error {
	type alias TransferToBlindOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeTransferToBlind, &payload); err != nil {
		return err
	}
	*o = TransferToBlindOperation(payload)
	return nil
}

func (o TransferToBlindOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Amount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.From.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	bf, err := hex.DecodeString(strings.TrimSpace(o.BlindingFactor))
	if err != nil {
		return nil, err
	}
	if len(bf) != 32 {
		return nil, fmt.Errorf("blinding_factor must be 32 bytes")
	}
	if err := writeFixedBytes(w, bf, 32); err != nil {
		return nil, err
	}
	w.writeVarUint64(uint64(len(o.Outputs)))
	for i := range o.Outputs {
		if err := o.Outputs[i].MarshalBinaryInto(w); err != nil {
			return nil, err
		}
	}
	return w.Bytes(), nil
}

func (o *TransferToBlindOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *TransferToBlindOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	amount, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	from, err := readObjectID(r)
	if err != nil {
		return err
	}
	bf, err := readFixedBytes(r, 32)
	if err != nil {
		return err
	}
	count, err := r.readVarUint64()
	if err != nil {
		return err
	}
	outputs := make([]BlindOutput, 0, count)
	for i := uint64(0); i < count; i++ {
		var output BlindOutput
		if err := output.UnmarshalBinaryFrom(r); err != nil {
			return err
		}
		outputs = append(outputs, output)
	}
	o.Fee = fee
	o.Amount = amount
	o.From = from
	o.BlindingFactor = hex.EncodeToString(bf)
	o.Outputs = outputs
	return nil
}

type BlindTransferOperation struct {
	Fee     AssetAmount   `json:"fee"`
	Inputs  []BlindInput  `json:"inputs"`
	Outputs []BlindOutput `json:"outputs"`
}

func (o BlindTransferOperation) Type() OperationType { return OperationTypeBlindTransfer }

func (o BlindTransferOperation) MarshalJSON() ([]byte, error) {
	type alias BlindTransferOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *BlindTransferOperation) UnmarshalJSON(data []byte) error {
	type alias BlindTransferOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeBlindTransfer, &payload); err != nil {
		return err
	}
	*o = BlindTransferOperation(payload)
	return nil
}

func (o BlindTransferOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	w.writeVarUint64(uint64(len(o.Inputs)))
	for i := range o.Inputs {
		if err := o.Inputs[i].MarshalBinaryInto(w); err != nil {
			return nil, err
		}
	}
	w.writeVarUint64(uint64(len(o.Outputs)))
	for i := range o.Outputs {
		if err := o.Outputs[i].MarshalBinaryInto(w); err != nil {
			return nil, err
		}
	}
	return w.Bytes(), nil
}

func (o *BlindTransferOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *BlindTransferOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	inCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	inputs := make([]BlindInput, 0, inCount)
	for i := uint64(0); i < inCount; i++ {
		var input BlindInput
		if err := input.UnmarshalBinaryFrom(r); err != nil {
			return err
		}
		inputs = append(inputs, input)
	}
	outCount, err := r.readVarUint64()
	if err != nil {
		return err
	}
	outputs := make([]BlindOutput, 0, outCount)
	for i := uint64(0); i < outCount; i++ {
		var output BlindOutput
		if err := output.UnmarshalBinaryFrom(r); err != nil {
			return err
		}
		outputs = append(outputs, output)
	}
	o.Fee = fee
	o.Inputs = inputs
	o.Outputs = outputs
	return nil
}

type TransferFromBlindOperation struct {
	Fee            AssetAmount  `json:"fee"`
	Amount         AssetAmount  `json:"amount"`
	To             ObjectID     `json:"to"`
	BlindingFactor string       `json:"blinding_factor"`
	Inputs         []BlindInput `json:"inputs"`
}

func (o TransferFromBlindOperation) Type() OperationType { return OperationTypeTransferFromBlind }

func (o TransferFromBlindOperation) MarshalJSON() ([]byte, error) {
	type alias TransferFromBlindOperation
	return marshalOperation(o.Type(), alias(o))
}

func (o *TransferFromBlindOperation) UnmarshalJSON(data []byte) error {
	type alias TransferFromBlindOperation
	var payload alias
	if err := unmarshalOperationBody(data, OperationTypeTransferFromBlind, &payload); err != nil {
		return err
	}
	*o = TransferFromBlindOperation(payload)
	return nil
}

func (o TransferFromBlindOperation) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(o.Type()))
	if err := o.Fee.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.Amount.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	if err := o.To.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	bf, err := hex.DecodeString(strings.TrimSpace(o.BlindingFactor))
	if err != nil {
		return nil, err
	}
	if len(bf) != 32 {
		return nil, fmt.Errorf("blinding_factor must be 32 bytes")
	}
	if err := writeFixedBytes(w, bf, 32); err != nil {
		return nil, err
	}
	w.writeVarUint64(uint64(len(o.Inputs)))
	for i := range o.Inputs {
		if err := o.Inputs[i].MarshalBinaryInto(w); err != nil {
			return nil, err
		}
	}
	return w.Bytes(), nil
}

func (o *TransferFromBlindOperation) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

func (o *TransferFromBlindOperation) UnmarshalBinaryFrom(r *binaryReader) error {
	fee, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	amount, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	to, err := readObjectID(r)
	if err != nil {
		return err
	}
	bf, err := readFixedBytes(r, 32)
	if err != nil {
		return err
	}
	count, err := r.readVarUint64()
	if err != nil {
		return err
	}
	inputs := make([]BlindInput, 0, count)
	for i := uint64(0); i < count; i++ {
		var input BlindInput
		if err := input.UnmarshalBinaryFrom(r); err != nil {
			return err
		}
		inputs = append(inputs, input)
	}
	o.Fee = fee
	o.Amount = amount
	o.To = to
	o.BlindingFactor = hex.EncodeToString(bf)
	o.Inputs = inputs
	return nil
}

func init() {
	RegisterOperationFactory(OperationTypeTransferToBlind, func() Operation { return &TransferToBlindOperation{} })
	RegisterOperationFactory(OperationTypeBlindTransfer, func() Operation { return &BlindTransferOperation{} })
	RegisterOperationFactory(OperationTypeTransferFromBlind, func() Operation { return &TransferFromBlindOperation{} })
}

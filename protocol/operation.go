package protocol

import (
	"encoding/json"
	"fmt"
)

// Operation describes a BitShares protocol operation.
type Operation interface {
	Type() OperationType
}

// OperationBody preserves the operation payload as raw JSON.
type OperationBody struct {
	Kind    OperationType
	Payload json.RawMessage
}

func NewOperationBody(kind OperationType, payload any) (*OperationBody, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &OperationBody{Kind: kind, Payload: raw}, nil
}

func NewRawOperationBody(kind OperationType, payload []byte) *OperationBody {
	return &OperationBody{Kind: kind, Payload: append(json.RawMessage(nil), payload...)}
}

func (o *OperationBody) Type() OperationType {
	if o == nil {
		return 0
	}
	return o.Kind
}

func (o OperationBody) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{uint16(o.Kind), json.RawMessage(o.Payload)})
}

func (o *OperationBody) UnmarshalJSON(data []byte) error {
	var pair []json.RawMessage
	if err := json.Unmarshal(data, &pair); err != nil {
		return err
	}

	if len(pair) != 2 {
		return fmt.Errorf("invalid operation payload")
	}

	var opType uint16
	if err := json.Unmarshal(pair[0], &opType); err != nil {
		return err
	}

	o.Kind = OperationType(opType)
	o.Payload = append(o.Payload[:0], pair[1]...)
	return nil
}

// RawOperation is an alias for an untyped operation payload.
type RawOperation struct {
	OperationBody
}

// OperationEnvelope stores one protocol operation.
type OperationEnvelope struct {
	Operation Operation
}

func (e OperationEnvelope) Type() OperationType {
	if e.Operation == nil {
		return 0
	}
	return e.Operation.Type()
}

func (e OperationEnvelope) MarshalJSON() ([]byte, error) {
	if e.Operation == nil {
		return nil, fmt.Errorf("nil operation")
	}

	return json.Marshal(e.Operation)
}

func (e *OperationEnvelope) UnmarshalJSON(data []byte) error {
	var body OperationBody
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}

	if factory := newOperation(body.Kind); factory != nil {
		if err := json.Unmarshal(body.Payload, factory); err != nil {
			return err
		}
		e.Operation = factory
		return nil
	}

	e.Operation = &RawOperation{OperationBody: body}
	return nil
}

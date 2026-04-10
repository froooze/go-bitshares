package protocol

import "encoding/json"

// Transaction is a generic BitShares transaction envelope.
type Transaction struct {
	RefBlockNum    uint16              `json:"ref_block_num"`
	RefBlockPrefix uint32              `json:"ref_block_prefix"`
	Expiration     Time                `json:"expiration"`
	Operations     []OperationEnvelope `json:"operations"`
	Extensions     []json.RawMessage   `json:"extensions"`
}

// Push appends an operation envelope.
func (tx *Transaction) Push(op Operation) {
	tx.Operations = append(tx.Operations, OperationEnvelope{
		Operation: op,
	})
}

// PushOperation is a compatibility alias for Push.
func (tx *Transaction) PushOperation(op Operation) {
	tx.Push(op)
}

// SignedTransaction carries signatures for a transaction.
type SignedTransaction struct {
	Transaction
	Signatures []string `json:"signatures"`
}

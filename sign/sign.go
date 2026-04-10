package sign

import (
	"errors"

	"github.com/froooze/go-bitshares/protocol"
)

// Signer can transform a transaction into a signed transaction.
type Signer interface {
	Sign(tx *protocol.Transaction) (*protocol.SignedTransaction, error)
}

// NoopSigner copies the transaction into a signed transaction without adding signatures.
type NoopSigner struct{}

func (NoopSigner) Sign(tx *protocol.Transaction) (*protocol.SignedTransaction, error) {
	if tx == nil {
		return nil, errors.New("nil transaction")
	}

	return &protocol.SignedTransaction{Transaction: *tx}, nil
}

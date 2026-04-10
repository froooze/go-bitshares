package sign

import (
	"fmt"
	"strings"

	"github.com/froooze/go-bitshares/ecc"
	"github.com/froooze/go-bitshares/protocol"
)

// TransactionSigner signs BitShares transactions using one or more private keys.
type TransactionSigner struct {
	ChainID string
	Keys    []*ecc.PrivateKey
}

// Sign applies recoverable compact signatures to the supplied transaction.
// Duplicate public keys are signed only once.
func (s TransactionSigner) Sign(tx *protocol.Transaction) (*protocol.SignedTransaction, error) {
	if tx == nil {
		return nil, fmt.Errorf("nil transaction")
	}
	if strings.TrimSpace(s.ChainID) == "" {
		return nil, fmt.Errorf("chain id is required")
	}
	if len(s.Keys) == 0 {
		return nil, fmt.Errorf("no signing keys configured")
	}

	digest, err := tx.SigningDigest(s.ChainID)
	if err != nil {
		return nil, err
	}

	signatures := make([]string, 0, len(s.Keys))
	seen := make(map[string]struct{}, len(s.Keys))
	for _, key := range s.Keys {
		if key == nil {
			continue
		}
		pub := key.PublicKey()
		if pub == nil {
			continue
		}
		keyID := pub.String()
		if _, ok := seen[keyID]; ok {
			continue
		}
		seen[keyID] = struct{}{}

		sig, err := key.SignCompact(digest)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, sig.Hex())
	}

	if len(signatures) == 0 {
		return nil, fmt.Errorf("no valid signing keys configured")
	}

	return &protocol.SignedTransaction{
		Transaction: *tx,
		Signatures:  signatures,
	}, nil
}

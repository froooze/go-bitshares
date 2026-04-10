package protocol

import (
	"encoding/json"
	"testing"

	"github.com/froooze/go-bitshares/ecc"
)

func TestAccountCreateBinaryRoundTrip(t *testing.T) {
	pub := ecc.PrivateKeyFromSeed([]byte("core-operations-test")).PublicKey().String()
	op := AccountCreateOperation{
		Fee:             AssetAmount{Amount: 1, AssetID: MustParseObjectID("1.3.0")},
		Registrar:       MustParseObjectID("1.2.1"),
		Referrer:        MustParseObjectID("1.2.2"),
		ReferrerPercent: 1000,
		Name:            "alice",
		Owner: Authority{
			WeightThreshold: 1,
			KeyAuths: map[PublicKey]uint16{
				MustPublicKey(pub): 1,
			},
		},
		Active: Authority{
			WeightThreshold: 1,
			AccountAuths: map[ObjectID]uint16{
				MustParseObjectID("1.2.3"): 1,
			},
		},
		Options: AccountOptions{
			MemoKey:       MustPublicKey(pub),
			VotingAccount: MustParseObjectID("1.2.5"),
			NumWitness:    1,
			NumCommittee:  1,
			Votes:         []VoteID{{Type: 1, ID: 2}},
		},
	}

	raw, err := op.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary() error = %v", err)
	}

	env, err := readOperationEnvelope(newBinaryReader(raw))
	if err != nil {
		t.Fatalf("readOperationEnvelope() error = %v", err)
	}
	decoded, ok := env.Operation.(*AccountCreateOperation)
	if !ok {
		t.Fatalf("unexpected decoded type %T", env.Operation)
	}
	if decoded.Name != op.Name || decoded.ReferrerPercent != op.ReferrerPercent {
		t.Fatalf("unexpected decoded account create operation: %#v", decoded)
	}

	js, err := json.Marshal(op)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	var decodedJSON AccountCreateOperation
	if err := json.Unmarshal(js, &decodedJSON); err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}
	if decodedJSON.Name != op.Name || decodedJSON.Options.MemoKey != op.Options.MemoKey {
		t.Fatalf("unexpected JSON decoded account create operation: %#v", decodedJSON)
	}
}

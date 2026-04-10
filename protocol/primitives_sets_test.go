package protocol

import (
	"testing"

	"github.com/froooze/go-bitshares/ecc"
)

func TestWriteObjectIDSetCanonicalizesDuplicates(t *testing.T) {
	t.Parallel()

	w := newBinaryWriter()
	err := writeObjectIDSet(w, []ObjectID{
		MustParseObjectID("1.2.2"),
		MustParseObjectID("1.2.1"),
		MustParseObjectID("1.2.2"),
	})
	if err != nil {
		t.Fatalf("writeObjectIDSet() error = %v", err)
	}

	got, err := readObjectIDSet(newBinaryReader(w.Bytes()))
	if err != nil {
		t.Fatalf("readObjectIDSet() error = %v", err)
	}
	if gotLen, want := len(got), 2; gotLen != want {
		t.Fatalf("len(got) = %d, want %d", gotLen, want)
	}
	if got0, want := got[0].String(), "1.2.1"; got0 != want {
		t.Fatalf("got[0] = %q, want %q", got0, want)
	}
	if got1, want := got[1].String(), "1.2.2"; got1 != want {
		t.Fatalf("got[1] = %q, want %q", got1, want)
	}
}

func TestWritePublicKeySetCanonicalizesDuplicates(t *testing.T) {
	t.Parallel()

	keyA := MustPublicKey(ecc.PrivateKeyFromSeed([]byte("set-key-a")).PublicKey().String())
	keyB := MustPublicKey(ecc.PrivateKeyFromSeed([]byte("set-key-b")).PublicKey().String())

	w := newBinaryWriter()
	err := writePublicKeySet(w, []PublicKey{keyB, keyA, keyB})
	if err != nil {
		t.Fatalf("writePublicKeySet() error = %v", err)
	}

	got, err := readPublicKeySet(newBinaryReader(w.Bytes()))
	if err != nil {
		t.Fatalf("readPublicKeySet() error = %v", err)
	}
	if gotLen, want := len(got), 2; gotLen != want {
		t.Fatalf("len(got) = %d, want %d", gotLen, want)
	}
	if compareStrings(got[0].String(), got[1].String()) >= 0 {
		t.Fatalf("expected public keys to be sorted, got %q then %q", got[0].String(), got[1].String())
	}
}

func TestWriteAddressSetCanonicalizesDuplicates(t *testing.T) {
	t.Parallel()

	keyA := ecc.PrivateKeyFromSeed([]byte("set-address-a")).PublicKey()
	keyB := ecc.PrivateKeyFromSeed([]byte("set-address-b")).PublicKey()
	addrA := MustAddress(keyA.ToAddressString())
	addrB := MustAddress(keyB.ToAddressString())

	w := newBinaryWriter()
	err := writeAddressSet(w, []Address{addrB, addrA, addrB})
	if err != nil {
		t.Fatalf("writeAddressSet() error = %v", err)
	}

	got, err := readAddressSet(newBinaryReader(w.Bytes()))
	if err != nil {
		t.Fatalf("readAddressSet() error = %v", err)
	}
	if gotLen, want := len(got), 2; gotLen != want {
		t.Fatalf("len(got) = %d, want %d", gotLen, want)
	}
	if compareStrings(got[0].String(), got[1].String()) >= 0 {
		t.Fatalf("expected addresses to be sorted, got %q then %q", got[0].String(), got[1].String())
	}
}

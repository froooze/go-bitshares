package ecc

import (
	"crypto/sha256"
	"testing"
)

func TestPrivateKeyRoundTrip(t *testing.T) {
	t.Parallel()

	SetAddressPrefix("BTS")

	priv := PrivateKeyFromSeed([]byte("correct horse battery staple"))
	if priv == nil {
		t.Fatal("PrivateKeyFromSeed returned nil")
	}

	wif := priv.WIF()
	restored, err := PrivateKeyFromWIF([]byte(wif))
	if err != nil {
		t.Fatalf("PrivateKeyFromWIF() error = %v", err)
	}
	if got, want := restored.WIF(), wif; got != want {
		t.Fatalf("WIF round-trip = %q, want %q", got, want)
	}

	pub := priv.PublicKey()
	parsed, err := PublicKeyFromString(pub.String())
	if err != nil {
		t.Fatalf("PublicKeyFromString() error = %v", err)
	}
	if got, want := parsed.String(), pub.String(); got != want {
		t.Fatalf("public key round-trip = %q, want %q", got, want)
	}
}

func TestBrainKeyDerivation(t *testing.T) {
	t.Parallel()

	priv, err := FromBrainKey([]byte(" twelve words  here "), 3)
	if err != nil {
		t.Fatalf("FromBrainKey() error = %v", err)
	}
	if priv == nil || len(priv.Bytes()) != 32 {
		t.Fatalf("FromBrainKey() returned invalid key: %#v", priv)
	}
}

func TestCompactSignatureRecovery(t *testing.T) {
	t.Parallel()

	priv := PrivateKeyFromSeed([]byte("signature-seed"))
	sum := sha256.Sum256([]byte("bitshares-signature"))

	sig, err := priv.SignCompact(sum[:])
	if err != nil {
		t.Fatalf("SignCompact() error = %v", err)
	}

	recovered, compressed, err := sig.RecoverPublicKey(sum[:])
	if err != nil {
		t.Fatalf("RecoverPublicKey() error = %v", err)
	}
	if !compressed {
		t.Fatal("RecoverPublicKey() = compressed=false, want true")
	}
	if got, want := recovered.String(), priv.PublicKey().String(); got != want {
		t.Fatalf("recovered key = %q, want %q", got, want)
	}
}

func TestMemoEncryptDecrypt(t *testing.T) {
	t.Parallel()

	sender := PrivateKeyFromSeed([]byte("sender"))
	recipient := PrivateKeyFromSeed([]byte("recipient"))
	message := []byte("memo payload")
	nonce := "123456789"

	enc, err := EncryptWithChecksum(sender, recipient.PublicKey(), nonce, message)
	if err != nil {
		t.Fatalf("EncryptWithChecksum() error = %v", err)
	}

	dec, err := DecryptWithChecksum(recipient, sender.PublicKey(), nonce, enc, false)
	if err != nil {
		t.Fatalf("DecryptWithChecksum() error = %v", err)
	}
	if got, want := string(dec), string(message); got != want {
		t.Fatalf("memo decrypt = %q, want %q", got, want)
	}
}

func TestGenerateKeys(t *testing.T) {
	t.Parallel()

	keys, pubs, err := GenerateKeys("alice", []byte("correct horse battery staple"), []string{"active", "memo"}, "BTS")
	if err != nil {
		t.Fatalf("GenerateKeys() error = %v", err)
	}
	if len(keys) != 2 || len(pubs) != 2 {
		t.Fatalf("GenerateKeys() sizes = %d/%d, want 2/2", len(keys), len(pubs))
	}
	if keys["active"] == nil || pubs["active"] == "" {
		t.Fatal("GenerateKeys() missing active key")
	}
}

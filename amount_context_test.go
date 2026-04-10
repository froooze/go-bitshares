package bitshares

import (
	"bytes"
	"errors"
	"testing"

	"github.com/ulikunitz/xz/lzma"
)

func TestRoundAmountExact(t *testing.T) {
	t.Parallel()

	if got, want := RoundAmount(0.29, 2), int64(29); got != want {
		t.Fatalf("RoundAmount(0.29, 2) = %d, want %d", got, want)
	}
}

func TestFindBackupActiveKeyMatchesAuthorizedKey(t *testing.T) {
	t.Parallel()

	auth := Authority{
		KeyAuths: []KeyWeightPair{
			{Key: "BTS1111111111111111111111111111111114T1Anm", Weight: 1},
			{Key: "BTS2222222222222222222222222222222224T1Anm", Weight: 1},
		},
	}
	keys := []backupKeyRecord{
		{PubKey: "BTS2222222222222222222222222222222224T1Anm", EncryptedKey: "deadbeef"},
	}

	entry := findBackupActiveKey(keys, auth)
	if entry == nil {
		t.Fatal("findBackupActiveKey() = nil, want matching key")
	}
	if got, want := entry.PubKey, keys[0].PubKey; got != want {
		t.Fatalf("findBackupActiveKey().PubKey = %q, want %q", got, want)
	}
}

func TestDecompressBackupPayloadRequiresContext(t *testing.T) {
	t.Parallel()

	var compressed bytes.Buffer
	writer, err := lzma.NewWriter(&compressed)
	if err != nil {
		t.Fatalf("lzma.NewWriter() error = %v", err)
	}
	if _, err := writer.Write([]byte(`{"wallet":[]}`)); err != nil {
		t.Fatalf("writer.Write() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close() error = %v", err)
	}

	_, err = decompressBackupPayload(nil, compressed.Bytes())
	if !errors.Is(err, ErrNilContext) {
		t.Fatalf("decompressBackupPayload(nil, ...) error = %v, want %v", err, ErrNilContext)
	}
}

package bitshares

import (
	"bytes"
	"context"
	"testing"

	"github.com/ulikunitz/xz/lzma"
)

func TestDecompressBackupPayload(t *testing.T) {
	t.Parallel()

	payload := []byte(`{"wallet":[{"encryption_key":"00"}],"private_keys":[]}`)
	var compressed bytes.Buffer
	writer, err := lzma.NewWriter(&compressed)
	if err != nil {
		t.Fatalf("lzma.NewWriter() error = %v", err)
	}
	if _, err := writer.Write(payload); err != nil {
		t.Fatalf("writer.Write() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close() error = %v", err)
	}

	decoded, err := decompressBackupPayload(context.Background(), compressed.Bytes())
	if err != nil {
		t.Fatalf("decompressBackupPayload() error = %v", err)
	}
	if got, want := string(decoded), string(payload); got != want {
		t.Fatalf("decompressed payload = %q, want %q", got, want)
	}
}

package bitshares

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestAuthorityUnmarshal(t *testing.T) {
	t.Parallel()

	var auth Authority
	data := []byte(`{
		"weight_threshold": 1,
		"key_auths": [["BTS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5SoWcX7", 1]],
		"account_auths": [["1.2.3", 2]],
		"extensions": []
	}`)
	if err := json.Unmarshal(data, &auth); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if len(auth.KeyAuths) != 1 {
		t.Fatalf("key auths len = %d, want 1", len(auth.KeyAuths))
	}
	if len(auth.AccountAuths) != 1 {
		t.Fatalf("account auths len = %d, want 1", len(auth.AccountAuths))
	}
	gotKey := auth.KeyAuths[0].Key
	wantKey := "BTS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5SoWcX7"
	if gotKey != wantKey {
		t.Fatalf("key auth key = %q (len=%d bytes=% x), want %q (len=%d bytes=% x)", gotKey, len(gotKey), []byte(gotKey), wantKey, len(wantKey), []byte(wantKey))
	}
	if got, want := auth.AccountAuths[0].Account, "1.2.3"; got != want {
		t.Fatalf("account auth = %q, want %q", got, want)
	}
}

func TestLookupChain(t *testing.T) {
	t.Parallel()

	cfg, ok := lookupChain(defaultChains[0].ChainID)
	if !ok {
		t.Fatal("lookupChain() = not found")
	}
	if got, want := cfg.AddressPrefix, "BTS"; got != want {
		t.Fatalf("chain prefix = %q, want %q", got, want)
	}
}

func TestChainClientAccountRequiresContextEvenWhenCached(t *testing.T) {
	t.Parallel()

	client := &ChainClient{
		accounts: map[string]*AccountInfo{
			"alice": {Name: "alice"},
		},
		accountByID: make(map[string]*AccountInfo),
		assets:      make(map[string]*AssetInfo),
		assetByID:   make(map[string]*AssetInfo),
	}

	_, err := client.Account(nil, "alice")
	if !errors.Is(err, ErrNilContext) {
		t.Fatalf("Account(nil, ...) error = %v, want %v", err, ErrNilContext)
	}
}

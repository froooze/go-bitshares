package transport

import (
	"encoding/json"
	"net/url"
	"testing"
)

func TestJoinPath(t *testing.T) {
	t.Parallel()

	got := joinPath(" /api/ ", "v1", "/accounts/", "1.2.3")
	want := "/api/v1/accounts/1.2.3"
	if got != want {
		t.Fatalf("joinPath() = %q, want %q", got, want)
	}
}

func TestResolvePath(t *testing.T) {
	t.Parallel()

	base, err := url.Parse("https://node.example/api")
	if err != nil {
		t.Fatalf("parse base URL: %v", err)
	}

	client := &RESTClient{baseURL: base}
	got := client.resolvePath("accounts/1.2.3?expand=true#section")
	want := "https://node.example/api/accounts/1.2.3?expand=true#section"
	if got != want {
		t.Fatalf("resolvePath() = %q, want %q", got, want)
	}
}

func TestEncodeRequestBody(t *testing.T) {
	t.Parallel()

	raw, err := encodeRequestBody(json.RawMessage(`{"hello":"world"}`))
	if err != nil {
		t.Fatalf("encode raw message: %v", err)
	}
	if string(raw) != `{"hello":"world"}` {
		t.Fatalf("encode raw message = %q", string(raw))
	}

	text, err := encodeRequestBody("plain-text")
	if err != nil {
		t.Fatalf("encode string: %v", err)
	}
	if string(text) != "plain-text" {
		t.Fatalf("encode string = %q", string(text))
	}
}

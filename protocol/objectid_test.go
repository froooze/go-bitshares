package protocol

import "testing"

func TestObjectIDRoundTrip(t *testing.T) {
	original := MustParseObjectID("1.2.345")
	if got := original.String(); got != "1.2.345" {
		t.Fatalf("unexpected string: %s", got)
	}

	parsed, err := ParseObjectID("1.2.345")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if parsed != original {
		t.Fatalf("round trip mismatch: %+v != %+v", parsed, original)
	}
}

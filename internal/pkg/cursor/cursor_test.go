package cursor

import (
	"testing"
	"time"
)

func TestCursorEncodeDecode(t *testing.T) {
	testTime := time.Date(2026, 8, 15, 23, 0, 0, 123456000, time.UTC)
	testID := "prod_abc123"

	encoded := Encode(testTime, testID)
	if encoded == "" {
		t.Fatalf("expected encoded cursor string, got empty")
	}

	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatalf("unexpected error decoding cursor: %v", err)
	}

	if decoded == nil {
		t.Fatalf("expected non-nil cursor")
	}

	if decoded.ID != testID {
		t.Errorf("expected ID %s, got %s", testID, decoded.ID)
	}

	if decoded.CreatedAt.UnixNano() != testTime.UnixNano() {
		t.Errorf("expected timestamp %v, got %v", testTime.UnixNano(), decoded.CreatedAt.UnixNano())
	}
}

func TestCursorEmpty(t *testing.T) {
	decoded, err := Decode("")
	if err != nil {
		t.Fatalf("expected no error for empty cursor, got: %v", err)
	}
	if decoded != nil {
		t.Errorf("expected nil for empty cursor, got: %v", decoded)
	}
}

func TestCursorInvalid(t *testing.T) {
	_, err := Decode("not-a-valid-base64-cursor!!!")
	if err == nil {
		t.Fatalf("expected error for invalid cursor, got nil")
	}
}

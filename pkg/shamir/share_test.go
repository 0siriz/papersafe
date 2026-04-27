package shamir

import (
	"bytes"
	"testing"
)

func TestShareMarshalUnmarshal(t *testing.T) {
	original := Share{
		X: 42,
		Y: []byte{1, 2, 3, 4},
	}

	data, err := original.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded Share
	if err := decoded.UnmarshalBinary(data); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if original.X != decoded.X {
		t.Fatalf("X mismatch: %d != %d", original.X, decoded.X)
	}

	if !bytes.Equal(original.Y, decoded.Y) {
		t.Fatalf("Y mismatch: %v != %v", original.Y, decoded.Y)
	}
}

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

func TestSealAndUnsealRoundtrip(t *testing.T) {
	original := []byte{1, 2, 3, 4}
	key := []byte{9, 9, 9, 9}

	s := Share{
		X: 1,
		Y: append([]byte{}, original...),
	}

	if err := s.Seal(key); err != nil {
		t.Fatalf("seal failed: %v", err)
	}

	if !s.Sealed() {
		t.Fatal("expected share to be sealed")
	}

	if bytes.Equal(original, s.Y) {
		t.Fatal("expected Y to be modified after sealing")
	}

	if err := s.Unseal(key); err != nil {
		t.Fatalf("unseal failed: %v", err)
	}

	if s.Sealed() {
		t.Fatal("expected share to be unsealed")
	}

	if !bytes.Equal(original, s.Y) {
		t.Fatalf("expected original data restored, got %v", s.Y)
	}
}

func TestSealFailsIfAlreadySealed(t *testing.T) {
	s := Share{
		X: 1,
		Y: []byte{1, 2, 3},
	}

	key := []byte{1, 2, 3}

	if err := s.Seal(key); err != nil {
		t.Fatalf("initial seal failed: %v", err)
	}

	if err := s.Seal(key); err == nil {
		t.Fatal("expected error on double seal")
	}
}

func TestUnsealFailsIfNotSealed(t *testing.T) {
	s := Share{
		X: 1,
		Y: []byte{1, 2, 3},
	}

	key := []byte{1, 2, 3}

	if err := s.Unseal(key); err == nil {
		t.Fatal("expected error when unsealing unsealed share")
	}
}

func TestSealKeyLengthMismatch(t *testing.T) {
	s := Share{
		X: 1,
		Y: []byte{1, 2, 3},
	}

	key := []byte{1, 2}

	if err := s.Seal(key); err == nil {
		t.Fatal("expected key length error")
	}
}

func TestSealIsSymmetric(t *testing.T) {
	original := []byte{10, 20, 30}
	key := []byte{7, 7, 7}

	s := Share{
		X: 1,
		Y: append([]byte{}, original...),
	}

	if err := s.seal(key); err != nil {
		t.Fatalf("seal failed: %v", err)
	}

	if err := s.seal(key); err != nil {
		t.Fatalf("second seal failed: %v", err)
	}

	if !bytes.Equal(original, s.Y) {
		t.Fatal("expected double XOR to restore original")
	}
}

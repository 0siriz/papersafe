package shamir

import (
	"bytes"
	"testing"
)

func TestSplitAndCombine(t *testing.T) {
	secret := []byte("supersecret")

	parts := 5
	threshold := 3

	shares, err := Split(secret, parts, threshold)
	if err != nil {
		t.Fatalf("split failed: %v", err)
	}

	recovered, err := Combine(shares[:threshold])
	if err != nil {
		t.Fatalf("combine failed: %v", err)
	}

	if !bytes.Equal(secret, recovered) {
		t.Fatalf("expected %v, got %v", secret, recovered)
	}
}

func TestCombineFailsWithTooFewShares(t *testing.T) {
	secret := []byte("secret")

	shares, err := Split(secret, 3, 2)
	if err != nil {
		t.Fatalf("split failed: %v", err)
	}

	_, err = Combine(shares[:1])
	if err == nil {
		t.Fatal("expected error with too few shares")
	}
}

func TestCombineFailsWithDuplicateShares(t *testing.T) {
	secret := []byte("secret")

	shares, err := Split(secret, 3, 2)
	if err != nil {
		t.Fatalf("split failed: %v", err)
	}

	dup := []Share{shares[0], shares[0]}

	_, err = Combine(dup)
	if err == nil {
		t.Fatal("expected error for duplicate shares")
	}
}

func TestCombineFailsWithMismatchedLengths(t *testing.T) {
	s1 := Share{X: 1, Y: []byte{1, 2}}
	s2 := Share{X: 2, Y: []byte{1}}

	_, err := Combine([]Share{s1, s2})
	if err == nil {
		t.Fatal("expected mismatched length error")
	}
}

func TestSplitErrors(t *testing.T) {
	_, err := Split([]byte("secret"), 2, 3)
	if err == nil {
		t.Fatal("expected ErrInvalidParts")
	}

	_, err = Split([]byte("secret"), 256, 2)
	if err == nil {
		t.Fatal("expected ErrTooManyParts")
	}

	_, err = Split([]byte("secret"), 3, 1)
	if err == nil {
		t.Fatal("expected ErrInvalidThreshold")
	}

	_, err = Split([]byte{}, 3, 2)
	if err == nil {
		t.Fatal("expected ErrEmptySecret")
	}
}

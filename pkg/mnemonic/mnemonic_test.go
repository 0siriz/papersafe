package mnemonic

import (
	"bytes"
	"testing"
)

// helper: generate deterministic 32-byte entropy
func testEntropy() []byte {
	e := make([]byte, 32)
	for i := range 32 {
		e[i] = byte(i)
	}
	return e
}

func TestEntropyToMnemonicValid(t *testing.T) {
	entropy := testEntropy()

	words, err := EntropyToMnemonic(entropy)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(words) != 24 {
		t.Fatalf("expected 24 words, got %d", len(words))
	}

	for _, w := range words {
		if _, ok := wordIndex[w]; !ok {
			t.Fatalf("word not in wordlist: %s", w)
		}
	}
}

func TestEntropyToMnemonicInvalidEntropyLength(t *testing.T) {
	_, err := EntropyToMnemonic([]byte{1, 2, 3})

	if err != ErrInvalidEntropy {
		t.Fatalf("expected ErrInvalidEntropy, got %v", err)
	}
}

func TestMnemonicToEntropyValidRoundTrip(t *testing.T) {
	entropy := testEntropy()

	words, err := EntropyToMnemonic(entropy)
	if err != nil {
		t.Fatalf("failed to generate mnemonic: %v", err)
	}

	result, err := MnemonicToEntropy(words)
	if err != nil {
		t.Fatalf("failed to convert back: %v", err)
	}

	if !bytes.Equal(entropy, result) {
		t.Fatalf("roundtrip mismatch\nexpected: %v\ngot: %v", entropy, result)
	}
}

func TestMnemonicToEntropyInvalidLength(t *testing.T) {
	words := []string{"abandon"} // too short

	_, err := MnemonicToEntropy(words)
	if err != ErrInvalidMnemonic {
		t.Fatalf("expected ErrInvalidMnemonic, got %v", err)
	}
}

func TestMnemonicToEntropyInvalidWord(t *testing.T) {
	entropy := testEntropy()

	words, err := EntropyToMnemonic(entropy)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	words[0] = "notaword"

	_, err = MnemonicToEntropy(words)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if _, ok := err.(*InvalidWordError); !ok {
		t.Fatalf("expected InvalidWordError, got %T", err)
	}
}

func TestMnemonicToEntropyInvalidChecksum(t *testing.T) {
	entropy := testEntropy()

	words, err := EntropyToMnemonic(entropy)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// mutate one word but keep it valid
	for i := range wordlist {
		if wordlist[i] != words[0] {
			words[0] = wordlist[i]
			break
		}
	}

	_, err = MnemonicToEntropy(words)
	if err != ErrInvalidChecksum {
		t.Fatalf("expected ErrInvalidChecksum, got %v", err)
	}
}

func TestDeterministicEntropyToMnemonic(t *testing.T) {
	entropy := testEntropy()

	words1, err := EntropyToMnemonic(entropy)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	words2, err := EntropyToMnemonic(entropy)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	for i := range words1 {
		if words1[i] != words2[i] {
			t.Fatalf("non-deterministic output at index %d", i)
		}
	}
}

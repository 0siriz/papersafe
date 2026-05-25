package mnemonic

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidEntropy is returned when the entropy is not exactly 32 bytes.
	ErrInvalidEntropy = errors.New("entropy must be 32 bytes")
	// ErrInvalidMnemonic is returned when the mnemonic does not contain
	// exactly 24 words.
	ErrInvalidMnemonic = errors.New("mnemonic must be 24 words")
	// ErrInvalidChecksum is returned when the embedded checksum does not
	// match the recomputed SHA-256 checksum.
	ErrInvalidChecksum = errors.New("invalid checksum")
)

// InvalidWordError is returned when a word in the mnemonic is not found in
// the BIP39 wordlist.
type InvalidWordError struct {
	// BadWord is the unrecognized word.
	BadWord string
}

func (e *InvalidWordError) Error() string {
	return fmt.Sprintf("invalid word: %s", e.BadWord)
}

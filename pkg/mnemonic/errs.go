package mnemonic

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidEntropy  = errors.New("entropy must be 32 bytes")
	ErrInvalidMnemonic = errors.New("mnemonic must be 24 words")
	ErrInvalidChecksum = errors.New("invalid checksum")
)

type InvalidWordError struct {
	BadWord string
}

func (e *InvalidWordError) Error() string {
	return fmt.Sprintf("invalid word: %s", e.BadWord)
}

// Package shamir implements Shamir's Secret Sharing over GF(2^8). It
// splits a secret into N shares such that any K (threshold) of them can
// reconstruct the original secret, but fewer than K reveal nothing.
package shamir

import "errors"

var (
	// ErrInvalidParts is returned when the number of parts is less than the
	// threshold.
	ErrInvalidParts = errors.New("parts cannot be less then threshold")
	// ErrTooManyParts is returned when the number of parts exceeds 255.
	ErrTooManyParts = errors.New("parts cannot exceed 255")
	// ErrInvalidThreshold is returned when the threshold is less than 2.
	ErrInvalidThreshold = errors.New("threshold must be at least 2")
	// ErrEmptySecret is returned when attempting to split an empty secret.
	ErrEmptySecret = errors.New("cannot split an empty secret")
	// ErrNotEnoughShares is returned when fewer than two shares are provided
	// for reconstruction.
	ErrNotEnoughShares = errors.New("less than two shares cannot be used to reconstruct the secret")
	// ErrMismatchedLengths is returned when shares have different Y-value
	// lengths.
	ErrMismatchedLengths = errors.New("all shares must be the same length")
	// ErrDuplicateShares is returned when two shares share the same X
	// coordinate.
	ErrDuplicateShares = errors.New("duplicate share detected")
	// ErrInvalidSealCall is returned when sealing an already-sealed share.
	ErrInvalidSealCall = errors.New("cannot seal an already sealed share")
	// ErrInvalidUnsealCall is returned when unsealing an already-unsealed
	// share.
	ErrInvalidUnsealCall = errors.New("cannot unseal an already unsealed share")
	// ErrInvalidSealKeyLength is returned when the seal key has an invalid
	// length.
	ErrInvalidSealKeyLength = errors.New("invalid seal key length")
	// ErrSealedShareCombine is returned when attempting to combine a sealed
	// share.
	ErrSealedShareCombine = errors.New("cannot combine with sealed share")
)

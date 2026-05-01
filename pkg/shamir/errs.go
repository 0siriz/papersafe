package shamir

import "errors"

var (
	ErrInvalidParts         = errors.New("parts cannot be less then threshold")
	ErrTooManyParts         = errors.New("parts cannot exceed 255")
	ErrInvalidThreshold     = errors.New("threshold must be at least 2")
	ErrEmptySecret          = errors.New("cannot split an empty secret")
	ErrNotEnoughShares      = errors.New("less than two shares cannot be used to reconstruct the secret")
	ErrMismatchedLengths    = errors.New("all shares must be the same length")
	ErrDuplicateShares      = errors.New("duplicate share detected")
	ErrInvalidSealCall      = errors.New("cannot seal an already sealed share")
	ErrInvalidUnsealCall    = errors.New("cannot unseal an already unsealed share")
	ErrInvalidSealKeyLength = errors.New("invalid seal key length")
	ErrSealedShareCombine   = errors.New("cannot combine with sealed share")
)

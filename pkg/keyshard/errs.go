package keyshard

import "errors"

var (
	// ErrInvalidSignature is returned when a KeyShard's Ed25519 signature
	// fails verification.
	ErrInvalidSignature = errors.New("keyshard has an invalid signature")
)

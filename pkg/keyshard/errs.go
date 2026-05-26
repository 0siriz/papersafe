package keyshard

import "errors"

var (
	// ErrInvalidSignature is returned when a KeyShard's Ed25519 signature
	// fails verification.
	ErrInvalidSignature = errors.New("keyshard has an invalid signature")

	// ErrPublicKeyMismatch is returned when shards from different quorums
	// are passed to ReconstructQuorum.
	ErrPublicKeyMismatch = errors.New("shards have mismatched public keys")
)

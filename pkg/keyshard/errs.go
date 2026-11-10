package keyshard

import "errors"

var (
	// ErrInvalidSignature is returned when a KeyShard's Ed25519 signature
	// fails verification.
	ErrInvalidSignature = errors.New("keyshard has an invalid signature")

	// ErrPublicKeyMismatch is returned when shards from different quorums
	// are passed to ReconstructQuorum.
	ErrPublicKeyMismatch = errors.New("shards have mismatched public keys")

	// ErrInvalidThreshold is returned when a KeyShard has a threshold outside
	// the supported range.
	ErrInvalidThreshold = errors.New("keyshard has an invalid threshold")

	// ErrNotEnoughShards is returned when fewer shards than the threshold are
	// passed to ReconstructQuorum.
	ErrNotEnoughShards = errors.New("not enough shards to meet threshold")

	// ErrMismatchedThreshold is returned when shards with different
	// thresholds are passed to ReconstructQuorum.
	ErrMismatchedThreshold = errors.New("shards have mismatched thresholds")
)

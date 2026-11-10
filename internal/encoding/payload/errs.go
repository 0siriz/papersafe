package payload

import "errors"

var (
	// ErrInvalidPrefix is returned when the payload does not start with
	// the PAPERSAFE prefix and separator.
	ErrInvalidPrefix = errors.New("payload has an invalid prefix")
	// ErrInvalidDomain is returned when the payload does not contain a
	// known domain name followed by a separator.
	ErrInvalidDomain = errors.New("payload has an invalid domain")
	// ErrInvalidVersion is returned when the payload version is missing or
	// malformed.
	ErrInvalidVersion = errors.New("payload has an invalid version")
	// ErrUnsupportedVersion is returned when the payload version is not
	// supported by this implementation.
	ErrUnsupportedVersion = errors.New("payload version is not supported")
	// ErrInvalidChunk is returned when chunk metadata is malformed or the
	// chunk index is outside the range 1..total.
	ErrInvalidChunk = errors.New("payload has invalid chunk metadata")
	// ErrInvalidShardID is returned when keyshard metadata is malformed or the
	// shard ID is outside the range 1..255.
	ErrInvalidShardID = errors.New("payload has an invalid keyshard ID")
	// ErrInvalidShardParams is returned when the keyshard threshold is outside
	// the supported range.
	ErrInvalidShardParams = errors.New("payload has invalid shard parameters")
)

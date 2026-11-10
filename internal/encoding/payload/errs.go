package payload

import "errors"

var (
	// ErrInvalidPrefix is returned when the payload does not start with
	// the PAPERSAFE prefix and separator.
	ErrInvalidPrefix = errors.New("payload has an invalid prefix")
	// ErrInvalidDomain is returned when the payload does not contain a
	// known domain name followed by a separator.
	ErrInvalidDomain = errors.New("payload has an invalid domain")
	// ErrInvalidChunk is returned when chunk metadata is malformed or the
	// chunk index is outside the range 1..total.
	ErrInvalidChunk = errors.New("payload has invalid chunk metadata")
	// ErrInvalidShardID is returned when keyshard metadata is malformed or the
	// shard ID is outside the range 1..255.
	ErrInvalidShardID = errors.New("payload has an invalid keyshard ID")
)

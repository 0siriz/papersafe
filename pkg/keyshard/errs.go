package keyshard

import "errors"

var (
	ErrInvalidSignature   = errors.New("keyshard has an invalid signature")
	ErrZeroedSignatureKey = errors.New("quorum has a zeroed signature key")
)

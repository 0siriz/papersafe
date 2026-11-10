package payload

import (
	"bytes"
	"encoding/base64"
	"strconv"
)

// keyshardVersion is the version of the KEYSHARD payload wire format.
const keyshardVersion = 1

// Keyshard is a keyshard payload carrying the shard ID, the threshold needed
// to reconstruct the quorum, and the encrypted shard data.
type Keyshard struct {
	// ShardID is the Shamir share index for the shard.
	ShardID byte
	// Threshold is the number of shards required to reconstruct the quorum.
	Threshold byte
	Data      []byte
}

// NewKeyshard creates a keyshard payload carrying the shard ID and threshold.
func NewKeyshard(id, threshold byte, data []byte) Keyshard {
	return Keyshard{
		ShardID:   id,
		Threshold: threshold,
		Data:      data,
	}
}

// Domain returns the keyshard domain.
func (k Keyshard) Domain() Domain {
	return DomainKeyshard
}

// AppendText appends the text encoding of the keyshard payload to b.
func (k Keyshard) AppendText(b []byte) ([]byte, error) {
	if err := k.validate(); err != nil {
		return nil, err
	}

	b = appendHeader(b, DomainKeyshard)
	b = strconv.AppendUint(b, keyshardVersion, 10)
	b = append(b, separator...)
	b = strconv.AppendUint(b, uint64(k.ShardID), 10)
	b = append(b, separator...)
	b = strconv.AppendUint(b, uint64(k.Threshold), 10)
	b = append(b, separator...)
	b = base64.RawURLEncoding.AppendEncode(b, k.Data)

	return b, nil
}

// MarshalText returns the text encoding of the keyshard payload.
func (k Keyshard) MarshalText() ([]byte, error) {
	return k.AppendText(nil)
}

// UnmarshalText decodes a keyshard text encoding into the payload.
func (k *Keyshard) UnmarshalText(b []byte) error {
	domain, rest, err := splitHeader(b)
	if err != nil {
		return err
	}
	if domain != DomainKeyshard {
		return ErrInvalidDomain
	}

	return k.unmarshalBody(rest)
}

func (k *Keyshard) unmarshalBody(rest []byte) error {
	versionPart, rest, found := bytes.Cut(rest, []byte(separator))
	if !found {
		return ErrInvalidVersion
	}

	version, err := strconv.Atoi(string(versionPart))
	if err != nil {
		return ErrInvalidVersion
	}
	if version != keyshardVersion {
		return ErrUnsupportedVersion
	}

	idPart, rest, found := bytes.Cut(rest, []byte(separator))
	if !found {
		return ErrInvalidShardID
	}

	id, err := strconv.ParseUint(string(idPart), 10, 8)
	if err != nil || id == 0 {
		return ErrInvalidShardID
	}

	thresholdPart, dataPart, found := bytes.Cut(rest, []byte(separator))
	if !found {
		return ErrInvalidShardParams
	}

	threshold, err := strconv.ParseUint(string(thresholdPart), 10, 8)
	if err != nil || threshold < 2 {
		return ErrInvalidShardParams
	}

	bufSize := base64.RawURLEncoding.DecodedLen(len(dataPart))
	data, err := base64.RawURLEncoding.AppendDecode(make([]byte, 0, bufSize), dataPart)
	if err != nil {
		return err
	}

	k.ShardID = byte(id)
	k.Threshold = byte(threshold)
	k.Data = data

	return nil
}

func (k Keyshard) validate() error {
	if k.ShardID == 0 {
		return ErrInvalidShardID
	}
	if k.Threshold < 2 {
		return ErrInvalidShardParams
	}

	return nil
}

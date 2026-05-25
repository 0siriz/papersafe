package keyshard

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/0siriz/papersafe/pkg/shamir"
)

const (
	// QuorumKeySize is the byte size of the quorum encryption key.
	QuorumKeySize = 32
	// QuorumSigningKeySize is the byte size of the Ed25519 private key.
	QuorumSigningKeySize = ed25519.PrivateKeySize

	// SecretFlagSize is the byte size of the sealed/unsealed flag in the
	// Shamir secret.
	SecretFlagSize = 1
	// SecretFlagUnsealed marks a quorum secret as unsealed (signing key
	// included in the split).
	SecretFlagUnsealed = 0
	// SecretFlagSealed marks a quorum secret as sealed (signing key excluded
	// from the split).
	SecretFlagSealed = 1
)

// Quorum holds the encryption key and optional signing key for shard
// management. An unsealed quorum can issue new shards and its signing key is
// recoverable from the shards. A sealed quorum excludes the signing key from
// the split secret, so reconstruction cannot issue new shards.
type Quorum struct {
	sealed     bool
	Key        []byte
	SigningKey ed25519.PrivateKey
}

// NewQuorum creates a new unsealed Quorum with a random encryption key and
// signing keypair.
func NewQuorum() (*Quorum, error) {
	key := make([]byte, QuorumKeySize)
	rand.Read(key)

	_, signKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Quorum{
		Key:        key,
		SigningKey: signKey,
	}, nil
}

// IsSealed reports whether the quorum is sealed. A sealed quorum's signing key
// is excluded from the Shamir secret, preventing new shard issuance after
// reconstruction.
func (q *Quorum) IsSealed() bool {
	return q.sealed
}

// SetSealed marks the quorum as sealed. Subsequent calls to MakeKeyshards
// will exclude the signing key from the Shamir secret, while still signing
// each shard with the real key.
func (q *Quorum) SetSealed() {
	q.sealed = true
}

// Destroy zeros the signing key and marks the quorum as sealed. The quorum
// should not be used for further shard creation after this call.
func (q *Quorum) Destroy() {
	zeroize(q.SigningKey)
	q.sealed = true
}

// MakeKeyshards splits the quorum key material into the given number of
// Shamir shares, then encrypts and signs each share into a KeyShard.
//
// If the quorum is sealed, only the encryption key is split; the signing key
// is excluded from the Shamir secret but is still used for signing each shard.
// If unsealed, both the encryption key and signing key are split and
// recoverable from the shards.
func (q *Quorum) MakeKeyshards(parts, threshold int) ([]ShardSet, error) {
	signKey := q.SigningKey

	secret := make([]byte, 0, SecretFlagSize+QuorumKeySize+QuorumSigningKeySize)
	if q.sealed {
		secret = append(secret, SecretFlagSealed)
		secret = append(secret, q.Key...)
	} else {
		secret = append(secret, SecretFlagUnsealed)
		secret = append(secret, q.Key...)
		secret = append(secret, signKey...)
	}

	shares, err := shamir.Split(secret, parts, threshold)
	if err != nil {
		return nil, err
	}

	shardSets := make([]ShardSet, parts)
	for i, share := range shares {
		ss, err := MakeKeyshard(&share, signKey)
		if err != nil {
			return nil, err
		}
		shardSets[i] = *ss
	}

	return shardSets, nil
}

// ReconstructQuorum recovers a Quorum from its shards. If the shards were
// created by an unsealed quorum, both the encryption key and signing key are
// recovered. Sealed shards yield only the encryption key, preventing new shard
// issuance.
func ReconstructQuorum(shardSets []ShardSet) (*Quorum, error) {
	shares := make([]shamir.Share, len(shardSets))

	for i, ss := range shardSets {
		share, err := ss.Open()
		if err != nil {
			return nil, err
		}
		shares[i] = *share
	}

	secret, err := shamir.Combine(shares)
	if err != nil {
		return nil, err
	}

	sealed := secret[0] == SecretFlagSealed

	key := make([]byte, QuorumKeySize)
	copy(key, secret[SecretFlagSize:SecretFlagSize+QuorumKeySize])

	var signKey ed25519.PrivateKey
	if !sealed {
		signKey = make([]byte, QuorumSigningKeySize)
		copy(signKey, secret[SecretFlagSize+QuorumKeySize:])
	}

	return &Quorum{
		Key:        key,
		SigningKey: signKey,
		sealed:     sealed,
	}, nil
}

func zeroize(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

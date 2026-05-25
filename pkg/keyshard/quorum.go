package keyshard

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/0siriz/papersafe/pkg/shamir"
)

const (
	QuorumKeySize        = 32
	QuorumSigningKeySize = ed25519.PrivateKeySize
	SecretSize           = QuorumKeySize + QuorumSigningKeySize
)

type Quorum struct {
	sealed     bool
	Key        []byte
	SigningKey ed25519.PrivateKey
}

func NewQuorum(sealed bool) (*Quorum, error) {
	key := make([]byte, 32)
	rand.Read(key)

	_, signKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	quorum := &Quorum{
		Key:        key,
		SigningKey: signKey,
		sealed:     sealed,
	}

	return quorum, nil
}

func (q *Quorum) IsSealed() bool {
	return q.sealed
}

func (q *Quorum) MakeKeyshards(parts, threshold int) ([]KeyShard, [][]string, error) {
	signKey := q.SigningKey

	if allZero(signKey) {
		return nil, nil, ErrZeroedSignatureKey
	}

	if q.sealed {
		signKey = make([]byte, QuorumSigningKeySize)
	}

	secret := make([]byte, 0, SecretSize)
	secret = append(secret, q.Key...)
	secret = append(secret, signKey...)

	shares, err := shamir.Split(secret, parts, threshold)
	if err != nil {
		return nil, nil, err
	}

	keyshards := make([]KeyShard, parts)
	listOfWords := make([][]string, parts)

	for i, share := range shares {
		keyshard, words, err := MakeKeyshard(&share, signKey)
		if err != nil {
			return nil, nil, err
		}

		keyshards[i] = *keyshard
		listOfWords[i] = words
	}

	return keyshards, listOfWords, nil
}

func ReconstructQuorum(keyshards []KeyShard, listOfWords [][]string) (*Quorum, error) {
	shares := make([]shamir.Share, len(keyshards))

	for i, keyshard := range keyshards {
		share, err := keyshard.ToShare(listOfWords[i])
		if err != nil {
			return nil, err
		}

		shares[i] = *share
	}

	secret, err := shamir.Combine(shares)
	if err != nil {
		return nil, err
	}

	key := make([]byte, QuorumKeySize)
	copy(key, secret[:QuorumKeySize])

	signKey := make([]byte, QuorumSigningKeySize)
	copy(signKey, secret[QuorumKeySize:])

	quorum := &Quorum{
		Key:        key,
		SigningKey: signKey,
		sealed:     allZero(signKey),
	}

	return quorum, nil
}

func allZero(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}

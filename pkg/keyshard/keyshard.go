package keyshard

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/0siriz/papersafe/pkg/mnemonic"
	"github.com/0siriz/papersafe/pkg/shamir"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	IDSize             = 1
	PublicKeySize      = ed25519.PublicKeySize
	NonceSize          = chacha20poly1305.NonceSizeX
	SignatureSize      = ed25519.SignatureSize
	KeyShardStaticSize = IDSize + PublicKeySize + NonceSize + SignatureSize
)

type KeyShard struct {
	PublicKey []byte
	ID        byte
	Nonce     []byte
	Content   []byte
	Signature []byte
}

type ShardSet struct {
	Shard KeyShard
	Words []string
}

func (ss *ShardSet) Open() (*shamir.Share, error) {
	return ss.Shard.ToShare(ss.Words)
}

func (k *KeyShard) AppendBinary(b []byte) ([]byte, error) {
	b = append(b, k.PublicKey...)
	b = append(b, k.ID)
	b = append(b, k.Nonce...)
	b = append(b, k.Content...)
	b = append(b, k.Signature...)

	return b, nil
}

func (k *KeyShard) MarshalBinary() ([]byte, error) {
	b, err := k.AppendBinary(make([]byte, 0, KeyShardStaticSize+len(k.Content)))
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (k *KeyShard) UnmarshalBinary(b []byte) error {
	buf := b

	k.PublicKey = make([]byte, PublicKeySize)
	copy(k.PublicKey, buf[:PublicKeySize])
	buf = buf[PublicKeySize:]

	k.ID = buf[0]
	buf = buf[IDSize:]

	k.Nonce = make([]byte, NonceSize)
	copy(k.Nonce, buf[:NonceSize])
	buf = buf[NonceSize:]

	contentSize := len(b) - KeyShardStaticSize
	k.Content = make([]byte, contentSize)
	copy(k.Content, buf[:contentSize])
	buf = buf[contentSize:]

	k.Signature = make([]byte, SignatureSize)
	copy(k.Signature, buf[:SignatureSize])

	return nil
}

func (k *KeyShard) Verify() bool {
	signData := make([]byte, 0, PublicKeySize+IDSize+NonceSize+len(k.Content))
	signData = append(signData, k.PublicKey...)
	signData = append(signData, k.ID)
	signData = append(signData, k.Nonce...)
	signData = append(signData, k.Content...)

	ok := ed25519.Verify(k.PublicKey, signData, k.Signature)

	return ok
}

func MakeKeyshard(share *shamir.Share, signKey ed25519.PrivateKey) (*ShardSet, error) {
	encKey := make([]byte, chacha20poly1305.KeySize)
	rand.Read(encKey)

	aead, err := chacha20poly1305.NewX(encKey)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, NonceSize)
	rand.Read(nonce)

	pubSignKey := signKey.Public().(ed25519.PublicKey)

	additionalData := make([]byte, 0, PublicKeySize+IDSize)
	additionalData = append(additionalData, pubSignKey...)
	additionalData = append(additionalData, share.X)

	content := aead.Seal(nil, nonce, share.Y, additionalData)

	signData := make([]byte, 0, PublicKeySize+IDSize+NonceSize+len(content))
	signData = append(signData, pubSignKey...)
	signData = append(signData, share.X)
	signData = append(signData, nonce...)
	signData = append(signData, content...)

	signature := ed25519.Sign(signKey, signData)

	keyshard := &KeyShard{
		PublicKey: pubSignKey,
		ID:        share.X,
		Nonce:     nonce,
		Content:   content,
		Signature: signature,
	}

	words, err := mnemonic.EntropyToMnemonic(encKey)
	if err != nil {
		return nil, err
	}

	return &ShardSet{Shard: *keyshard, Words: words}, nil
}

func (k *KeyShard) ToShare(words []string) (*shamir.Share, error) {
	if ok := k.Verify(); !ok {
		return nil, ErrInvalidSignature
	}

	encKey, err := mnemonic.MnemonicToEntropy(words)
	if err != nil {
		return nil, err
	}

	aead, err := chacha20poly1305.NewX(encKey)
	if err != nil {
		return nil, err
	}

	additionalData := make([]byte, 0, PublicKeySize+IDSize)
	additionalData = append(additionalData, k.PublicKey...)
	additionalData = append(additionalData, k.ID)

	cleartext, err := aead.Open(nil, k.Nonce, k.Content, additionalData)
	if err != nil {
		return nil, err
	}

	share := &shamir.Share{
		X: k.ID,
		Y: cleartext,
	}

	return share, nil
}

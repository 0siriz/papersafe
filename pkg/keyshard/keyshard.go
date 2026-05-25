package keyshard

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/0siriz/papersafe/pkg/mnemonic"
	"github.com/0siriz/papersafe/pkg/shamir"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	// IDSize is the byte size of the KeyShard ID field.
	IDSize = 1
	// PublicKeySize is the byte size of an Ed25519 public key.
	PublicKeySize = ed25519.PublicKeySize
	// NonceSize is the byte size of the XChaCha20-Poly1305 nonce.
	NonceSize = chacha20poly1305.NonceSizeX
	// SignatureSize is the byte size of an Ed25519 signature.
	SignatureSize = ed25519.SignatureSize
	// KeyShardStaticSize is the total byte size of the fixed-length fields in
	// a KeyShard (ID + PublicKey + Nonce + Signature), excluding the variable-
	// length Content.
	KeyShardStaticSize = IDSize + PublicKeySize + NonceSize + SignatureSize
)

// KeyShard is an encrypted, signed Shamir share. It contains the encrypted
// Y values of a share, the public key used for verification, the nonce used
// for encryption, and the Ed25519 signature over the shard data.
type KeyShard struct {
	// PublicKey is the Ed25519 public key used to verify the shard signature.
	PublicKey []byte
	// ID is the share index (X coordinate).
	ID byte
	// Nonce is the XChaCha20-Poly1305 nonce used for encryption.
	Nonce []byte
	// Content is the encrypted Y values of the Shamir share.
	Content []byte
	// Signature is the Ed25519 signature over the shard data.
	Signature []byte
}

// ShardSet pairs a KeyShard with the BIP39 mnemonic words that encode its
// decryption key.
type ShardSet struct {
	// Shard is the encrypted key shard.
	Shard KeyShard
	// Words is the 24-word BIP39 mnemonic encoding the shard encryption key.
	Words []string
}

// Open decrypts the shard using the mnemonic words and returns the
// underlying Shamir share. It verifies the shard signature before
// decrypting.
func (ss *ShardSet) Open() (*shamir.Share, error) {
	return ss.Shard.ToShare(ss.Words)
}

// AppendBinary appends the binary encoding of the KeyShard to b.
func (k *KeyShard) AppendBinary(b []byte) ([]byte, error) {
	b = append(b, k.PublicKey...)
	b = append(b, k.ID)
	b = append(b, k.Nonce...)
	b = append(b, k.Content...)
	b = append(b, k.Signature...)

	return b, nil
}

// MarshalBinary returns the binary encoding of the KeyShard.
func (k *KeyShard) MarshalBinary() ([]byte, error) {
	b, err := k.AppendBinary(make([]byte, 0, KeyShardStaticSize+len(k.Content)))
	if err != nil {
		return nil, err
	}

	return b, nil
}

// UnmarshalBinary decodes the binary encoding into the KeyShard.
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

// Verify checks the Ed25519 signature against the shard data using the
// embedded public key. Returns true if the signature is valid.
func (k *KeyShard) Verify() bool {
	signData := make([]byte, 0, PublicKeySize+IDSize+NonceSize+len(k.Content))
	signData = append(signData, k.PublicKey...)
	signData = append(signData, k.ID)
	signData = append(signData, k.Nonce...)
	signData = append(signData, k.Content...)

	ok := ed25519.Verify(k.PublicKey, signData, k.Signature)

	return ok
}

// MakeKeyshard encrypts a Shamir share, signs it, and encodes the
// encryption key as a BIP39 mnemonic. The returned ShardSet pairs the
// encrypted shard with the mnemonic words.
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

// ToShare verifies the shard signature, decrypts the content using the
// mnemonic-encoded key, and returns the underlying Shamir share.
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

package keyshard

import (
	"bytes"
	"crypto/ed25519"
	"testing"

	"github.com/0siriz/papersafe/pkg/shamir"
)

func makeTestShare() shamir.Share {
	return shamir.Share{
		X: 7,
		Y: []byte{1, 2, 3, 4, 5},
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(nil)

	ks := &KeyShard{
		PublicKey: pub,
		ID:        1,
		Nonce:     make([]byte, NonceSize),
		Content:   []byte{1, 2, 3},
		Signature: make([]byte, SignatureSize),
	}

	data, err := ks.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded KeyShard
	if err := decoded.UnmarshalBinary(data); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if !bytes.Equal(ks.PublicKey, decoded.PublicKey) {
		t.Fatal("public key mismatch")
	}

	if ks.ID != decoded.ID {
		t.Fatal("ID mismatch")
	}

	if !bytes.Equal(ks.Nonce, decoded.Nonce) {
		t.Fatal("nonce mismatch")
	}

	if !bytes.Equal(ks.Content, decoded.Content) {
		t.Fatal("content mismatch")
	}

	if !bytes.Equal(ks.Signature, decoded.Signature) {
		t.Fatal("signature mismatch")
	}
}

func TestVerifyValidSignature(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)

	ks := &KeyShard{
		PublicKey: pub,
		ID:        1,
		Nonce:     make([]byte, NonceSize),
		Content:   []byte{9, 9, 9},
	}

	signData := make([]byte, 0)
	signData = append(signData, ks.PublicKey...)
	signData = append(signData, ks.ID)
	signData = append(signData, ks.Nonce...)
	signData = append(signData, ks.Content...)

	ks.Signature = ed25519.Sign(priv, signData)

	if !ks.Verify() {
		t.Fatal("expected signature to verify")
	}
}

func TestVerifyFailsOnTamper(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)

	ks := &KeyShard{
		PublicKey: pub,
		ID:        1,
		Nonce:     make([]byte, NonceSize),
		Content:   []byte{1, 2, 3},
	}

	signData := append([]byte{}, pub...)
	signData = append(signData, ks.ID)
	signData = append(signData, ks.Nonce...)
	signData = append(signData, ks.Content...)

	ks.Signature = ed25519.Sign(priv, signData)

	// Tamper with content
	ks.Content[0] ^= 0xFF

	if ks.Verify() {
		t.Fatal("expected verification to fail after tampering")
	}
}

func TestShareToKeyShardRoundtrip(t *testing.T) {
	_, priv, _ := ed25519.GenerateKey(nil)

	share := makeTestShare()

	ss, err := MakeKeyshard(share, priv)
	if err != nil {
		t.Fatalf("MakeKeyshard failed: %v", err)
	}

	if !ss.Shard.Verify() {
		t.Fatal("expected valid signature")
	}

	recovered, err := ss.Open()
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	if share.X != recovered.X {
		t.Fatal("X mismatch")
	}

	if !bytes.Equal(share.Y, recovered.Y) {
		t.Fatalf("Y mismatch: %v != %v", share.Y, recovered.Y)
	}
}

func TestKeyShardToShareWrongMnemonicFails(t *testing.T) {
	_, priv, _ := ed25519.GenerateKey(nil)

	share := makeTestShare()

	ss, err := MakeKeyshard(share, priv)
	if err != nil {
		t.Fatalf("MakeKeyshard failed: %v", err)
	}

	// Wrong words (same length but different)
	wrongWords := []string{"abandon", "abandon", "abandon", "abandon"}

	_, err = ss.Shard.ToShare(wrongWords)
	if err == nil {
		t.Fatal("expected decryption failure with wrong mnemonic")
	}
}

func TestKeyShardToShareInvalidSignature(t *testing.T) {
	_, priv, _ := ed25519.GenerateKey(nil)

	share := makeTestShare()

	ss, err := MakeKeyshard(share, priv)
	if err != nil {
		t.Fatalf("MakeKeyshard failed: %v", err)
	}

	// Corrupt signature
	ss.Shard.Signature[0] ^= 0xFF

	_, err = ss.Shard.ToShare(ss.Words)
	if err == nil {
		t.Fatal("expected failure due to invalid signature")
	}

	if err != ErrInvalidSignature {
		t.Fatalf("expected ErrInvalidSignature, got %v", err)
	}
}

func TestAEADAdditionalDataBinding(t *testing.T) {
	_, priv, _ := ed25519.GenerateKey(nil)

	share := makeTestShare()

	ss, err := MakeKeyshard(share, priv)
	if err != nil {
		t.Fatalf("MakeKeyshard failed: %v", err)
	}

	// Modify ID → should break AEAD
	ss.Shard.ID ^= 0xFF

	_, err = ss.Shard.ToShare(ss.Words)
	if err == nil {
		t.Fatal("expected failure due to modified additional data")
	}
}

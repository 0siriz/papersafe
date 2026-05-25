// Package keyshard provides encrypted, signed, and mnemonic-encoded Shamir
// shares. Each shard is a Share encrypted with XChaCha20-Poly1305, signed
// with Ed25519, and whose encryption key is encoded as a BIP39 mnemonic for
// offline paper storage.
package keyshard

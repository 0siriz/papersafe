// Command papersafe generates secure offline paper backups of cryptographic
// keys and other sensitive data. It uses Shamir's Secret Sharing to split a
// secret into shards, encrypts and signs each shard, and encodes the shard
// encryption key as a BIP39 mnemonic for offline paper storage.
package main

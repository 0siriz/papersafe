// Package shamir implements Shamir's Secret Sharing over GF(2^8). It
// splits a secret into N shares such that any K (threshold) of them can
// reconstruct the original secret, but fewer than K reveal nothing.
package shamir

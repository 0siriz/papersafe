package mnemonic

import (
	"crypto/sha256"
)

// EntropyToMnemonic converts 32 bytes of entropy into a 24-word BIP39
// mnemonic. The last word encodes an 8-bit SHA-256 checksum of the entropy.
func EntropyToMnemonic(entropy []byte) ([]string, error) {
	if len(entropy) != 32 {
		return nil, ErrInvalidEntropy
	}

	checksum, err := computeChecksum(entropy)
	if err != nil {
		return nil, err
	}

	data := append(entropy, checksum)
	var buffer uint32
	var bitsInBuffer uint

	words := make([]string, 0, 24)
	for _, b := range data {
		buffer = (buffer << 8) | uint32(b)
		bitsInBuffer += 8

		if bitsInBuffer >= 11 {
			bitsInBuffer -= 11
			idx := int(buffer >> bitsInBuffer)
			words = append(words, wordlist[idx])

			buffer &= (1 << bitsInBuffer) - 1
		}
	}

	return words, nil
}

// MnemonicToEntropy converts a 24-word BIP39 mnemonic back into the
// original 32 bytes of entropy, verifying the embedded checksum.
func MnemonicToEntropy(words []string) ([]byte, error) {
	if len(words) != 24 {
		return nil, ErrInvalidMnemonic
	}

	data := make([]byte, 0, 33)
	var buffer uint32
	var bitsInBuffer uint

	for _, w := range words {
		idx, ok := wordIndex[w]
		if !ok {
			return nil, &InvalidWordError{BadWord: w}
		}

		buffer = (buffer << 11) | uint32(idx&0x7FF) // idx is 11 bits
		bitsInBuffer += 11

		for bitsInBuffer >= 8 {
			bitsInBuffer -= 8
			byteVal := byte((buffer >> bitsInBuffer) & 0xFF)
			data = append(data, byteVal)

			buffer &= (1 << bitsInBuffer) - 1
		}
	}

	entropy := data[:32]
	checksum := data[32]

	expectedChecksum, err := computeChecksum(entropy)
	if err != nil {
		return nil, err
	}

	if checksum != expectedChecksum {
		return nil, ErrInvalidChecksum
	}

	return entropy, nil
}

func computeChecksum(entropy []byte) (byte, error) {
	h := sha256.New()
	if _, err := h.Write(entropy); err != nil {
		return 0, err
	}
	return h.Sum(nil)[0], nil
}

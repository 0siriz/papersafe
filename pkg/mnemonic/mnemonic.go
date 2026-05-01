package mnemonic

import (
	"crypto/sha256"
)

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

		buffer = (buffer << 11) | uint32(idx)
		bitsInBuffer += 11

		for bitsInBuffer >= 8 {
			bitsInBuffer -= 8
			byteVal := byte(buffer >> bitsInBuffer)
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

package shamir

import (
	"crypto/rand"
	"math/big"
)

func Split(secret []byte, parts, threshold int) ([]Share, error) {
	if parts < threshold {
		return nil, ErrInvalidParts
	}
	if parts > 255 {
		return nil, ErrTooManyParts
	}
	if threshold < 2 {
		return nil, ErrInvalidThreshold
	}
	if len(secret) == 0 {
		return nil, ErrEmptySecret
	}

	xCoordinates := perm()

	out := make([]Share, parts)
	for idx := range out {
		out[idx].X = xCoordinates[idx]
		out[idx].Y = make([]byte, len(secret))
	}

	for idx, val := range secret {
		p, err := makePolynomial(val, byte(threshold+1))
		if err != nil {
			return nil, err
		}

		for i := range parts {
			x := xCoordinates[i]
			y := p.evaluate(x)
			out[i].Y[idx] = y
		}
	}

	return out, nil
}

func Combine(shares []Share) ([]byte, error) {
	if len(shares) < 2 {
		return nil, ErrNotEnoughShares
	}

	firstShareLen := len(shares[0].Y)
	for i := 1; i < len(shares); i++ {
		if len(shares[i].Y) != firstShareLen {
			return nil, ErrMismatchedLengths
		}
	}

	secret := make([]byte, firstShareLen)

	xSamples := make([]byte, len(shares))
	ySamples := make([]byte, len(shares))

	checkMap := make(map[byte]bool)
	for i, share := range shares {
		samp := share.X
		if exists := checkMap[samp]; exists {
			return nil, ErrDuplicateShares
		}
		checkMap[samp] = true
		xSamples[i] = samp
	}

	for idx := range secret {
		for i, share := range shares {
			ySamples[i] = share.Y[idx]
		}

		val := interpolatePolynomial(xSamples, ySamples, 0)

		secret[idx] = val
	}

	return secret, nil
}

func perm() []byte {
	m := make([]byte, 255)

	for i := range 255 {
		bigJ, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := bigJ.Int64()
		m[i] = m[j]
		m[j] = byte(i + 1)
	}
	return m
}

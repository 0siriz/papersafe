package shamir

import "crypto/rand"

type polynomial struct {
	coefficients []byte
}

func (p *polynomial) evaluate(x byte) byte {
	if x == 0 {
		return p.coefficients[0]
	}

	degree := len(p.coefficients) - 1
	out := p.coefficients[degree]
	for i := degree - 1; i >= 0; i-- {
		c := p.coefficients[i]
		out = add(multiply(out, x), c)
	}
	return out
}

func makePolynomial(intercept, degree byte) (polynomial, error) {
	p := polynomial{
		coefficients: make([]byte, degree+1),
	}

	p.coefficients[0] = intercept

	if _, err := rand.Read(p.coefficients[1:]); err != nil {
		return polynomial{}, err
	}

	return p, nil
}

func interpolatePolynomial(xSamples, ySamples []byte, x byte) byte {
	limit := len(xSamples)
	var result, basis byte
	for i := 0; i < limit; i++ {
		basis = 1
		for j := 0; j < limit; j++ {
			if i == j {
				continue
			}
			num := add(x, xSamples[j])
			denominator := add(xSamples[i], xSamples[j])
			term := div(num, denominator)
			basis = multiply(basis, term)
		}
		group := multiply(ySamples[i], basis)
		result = add(result, group)
	}
	return result
}

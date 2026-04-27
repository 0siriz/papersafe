package shamir

import "testing"

func TestPolynomialEvaluate(t *testing.T) {
	p := polynomial{
		coefficients: []byte{5, 3, 2},
	}

	if got := p.evaluate(0); got != 5 {
		t.Fatalf("expected intercept 5, got %d", got)
	}

	v1 := p.evaluate(2)
	v2 := p.evaluate(2)

	if v1 != v2 {
		t.Fatalf("evaluation not deterministic: %d != %d", v1, v2)
	}
}

func TestInterpolatePolynomial(t *testing.T) {
	xSamples := []byte{1, 2}
	ySamples := []byte{
		add(10, multiply(2, 1)),
		add(10, multiply(2, 2)),
	}

	result := interpolatePolynomial(xSamples, ySamples, 0)

	if result != 10 {
		t.Fatalf("expected 10, got %d", result)
	}
}

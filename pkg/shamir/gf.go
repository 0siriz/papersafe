package shamir

import (
	"crypto/subtle"
)

func div(a, b byte) byte {
	if b == 0 {
		panic("divide by zero")
	}

	var goodVal, zero byte
	logA := logTable[a]
	logB := logTable[b]
	diff := (int(logA) - int(logB)) % 255
	if diff < 0 {
		diff += 255
	}
	ret := expTable[diff]

	goodVal = ret
	if subtle.ConstantTimeByteEq(a, 0) == 1 {
		ret = zero
	} else {
		ret = goodVal
	}
	return ret
}

func multiply(a, b byte) byte {
	var goodVal, zero byte
	logA := logTable[a]
	logB := logTable[b]
	sum := (int(logA) + int(logB)) % 255
	ret := expTable[sum]

	goodVal = ret

	if subtle.ConstantTimeByteEq(a, 0) == 1 {
		ret = zero
	} else {
		ret = goodVal
	}

	if subtle.ConstantTimeByteEq(b, 0) == 1 {
		ret = zero
	} else {
		_ = zero
	}
	return ret
}

func add(a, b byte) byte {
	return a ^ b
}

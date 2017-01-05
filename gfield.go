/**
 * 8-bit Galois Field
 * Copyright 2015, Klaus Post
 * Copyright 2015, Backblaze, Inc.  All rights reserved.
 */

package reedsolomon

func gfAdd(a, b byte) byte {
	return a ^ b
}

func gfSub(a, b byte) byte {
	return a ^ b
}

func gfMul(a, b byte) byte {
	return mulTable[a][b]
}

func gfDiv(a, b byte) (byte, error) {
	if a == 0 {
		return 0, nil
	}
	if b == 0 {
		err := errDividend
		return 0, err
	}
	logA := int(logTable[a])
	logB := int(logTable[b])
	logResult := logA - logB
	if logResult < 0 {
		logResult += 255
	}
	return expTable[logResult],nil
}

func gfExp(a byte, n int) byte {
	if n == 0 {
		return 1
	}
	if a == 0 {
		return 0
	}

	logA := logTable[a]
	logResult := int(logA) * n
	for logResult >= 255 {
		logResult -= 255
	}
	return byte(expTable[logResult])
}

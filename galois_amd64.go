//+build !noasm
//+build !appengine

// Copyright 2015, Klaus Post, see LICENSE for details.

package reedsolomon

// This is what the assembler does :
/*
func galMul(low, high, in, out []byte) {
	for n, input := range in {
		l := input & 0xf
		h := input >> 4
		out[n] = low[l] ^ high[h]
	}
}

func galMulXor(low, high, in, out []byte) {
	for n, input := range in {
		l := input & 0xf
		h := input >> 4
		out[n] ^= low[l] ^ high[h]
	}
}
*/

//go:noescape
func gfMulSSSE3(low, high, in, out []byte)

//go:noescape
func gfMulXorSSSE3(low, high, in, out []byte)

//go:noescape
func gfMulXorAVX2(low, high, in, out []byte)

//go:noescape
func gfMulAVX2(low, high, in, out []byte)

func gfMulRemainAVX2(coeff byte, input, output []byte, size int) {
	var done int
	remain := size -32
	if remain < 0 {
		mt := mulTable[coeff]
		for i := done; i < size; i++ {
			output[i] = mt[input[i]]
		}
	} else {
		gfMulAVX2(mulTableLow[coeff][:], mulTableHigh[coeff][:], input, output)
		done = (size >> 5) << 5
		remain = size - done
		if remain > 0 {
			mt := mulTable[coeff]
			for i := done; i < size; i++ {
				output[i] = mt[input[i]]
			}
		}
	}
}

func gfMulRemainXorAVX2(coeff byte, input, output []byte, size int) {
	var done int
	remain := size - 32
	if remain < 0 {
		mt := mulTable[coeff]
		for i := done; i < size; i++ {
			output[i] ^= mt[input[i]]
		}
	} else {
		gfMulXorAVX2(mulTableLow[coeff][:], mulTableHigh[coeff][:], input, output)
		done = (size >> 5) << 5
		remain = size - done
		if remain > 0 {
			mt := mulTable[coeff]
			for i := done; i < size; i++ {
				output[i] ^= mt[input[i]]
			}
		}
	}
}

func gfSliceMulSSSE3(coeff byte, input, output []byte, size int) {

	var done int
	gfMulSSSE3(mulTableLow[coeff][:], mulTableHigh[coeff][:], input, output)
	done = (size >> 4) << 4
	remain := size - done
	if remain > 0 {
		mt := mulTable[coeff]
		for i := done; i < size; i++ {
			output[i] = mt[input[i]]
		}
	}
}

func gfSliceMulXorSSSE3(coeff byte, input, output []byte, size int) {

	var done int
	gfMulXorSSSE3(mulTableLow[coeff][:], mulTableHigh[coeff][:], input, output)
	done = (size >> 4) << 4
	remain := size - done
	if remain > 0 {
		mt := mulTable[coeff]
		for i := done; i < size; i++ {
			output[i] ^= mt[input[i]]
		}
	}
}

package reedsolomon

//go:noescape
func gfMulXorAVX2(low, high, in, out []byte)

//go:noescape
func gfMulAVX2(low, high, in, out []byte)

func gfMulRemain(coeff byte, input, output []byte, size int) {
	var done int
	if size < 32 {
		mt := mulTable[coeff]
		for i := done; i < size; i++ {
			output[i] = mt[input[i]]
		}
	} else {
		gfMulAVX2(mulTableLow[coeff][:], mulTableHigh[coeff][:], input, output)
		done = (size >> 5) << 5
		remain := size - done
		if remain > 0 {
			mt := mulTable[coeff]
			for i := done; i < size; i++ {
				output[i] = mt[input[i]]
			}
		}
	}
}

func gfMulRemainXor(coeff byte, input, output []byte, size int) {
	var done int
	if size < 32 {
		mt := mulTable[coeff]
		for i := done; i < size; i++ {
			output[i] ^= mt[input[i]]
		}
	} else {
		gfMulXorAVX2(mulTableLow[coeff][:], mulTableHigh[coeff][:], input, output)
		done = (size >> 5) << 5
		remain := size - done
		if remain > 0 {
			mt := mulTable[coeff]
			for i := done; i < size; i++ {
				output[i] ^= mt[input[i]]
			}
		}
	}
}

package reedsolomon

//go:noescape
func gfMulSSSE3(low, high, in, out []byte)

//go:noescape
func gfMulXorSSSE3(low, high, in, out []byte)

//go:noescape
func gfMulXorAVX2(low, high, in, out []byte)

//go:noescape
func gfMulAVX2(low, high, in, out []byte)

func gfVectMul(c byte, in, out []byte) {
	var done int
	if cpuID == 0  {
		gfMulAVX2(mulTableLow[c][:], mulTableHigh[c][:], in, out)
		done = (len(in) >> 5) << 5
	} else if cpuID == 1 {
		gfMulSSSE3(mulTableLow[c][:], mulTableHigh[c][:], in, out)
		done = (len(in) >> 4) << 4
	}
	remain := len(in) - done
	if remain > 0 {
		mt := mulTable[c]
		for i := done; i < len(in); i++ {
			out[i] = mt[in[i]]
		}
	}
}

func gfVectMulXor(c byte, in, out []byte) {
	var done int
	if cpuID == 0 {
		gfMulXorAVX2(mulTableLow[c][:], mulTableHigh[c][:], in, out)
		done = (len(in) >> 5) << 5
	} else if cpuID == 1 {
		gfMulXorSSSE3(mulTableLow[c][:], mulTableHigh[c][:], in, out)
		done = (len(in) >> 4) << 4
	}
	remain := len(in) - done
	if remain > 0 {
		mt := mulTable[c]
		for i := done; i < len(in); i++ {
			out[i] ^= mt[in[i]]
		}
	}
}
package reedsolomon

func (r *rs) nosimdEncode(dp matrix) error {
	// check args
	if len(dp) != r.shards {
		return ErrTooFewShards
	}
	_, err := checkShardSize(dp)
	if err != nil {
		return err
	}
	// encoding
	input := dp[0:r.data]
	output := dp[r.data:]
	nosimdRunner(r.gen, input, output, r.data, r.parity)
	return nil
}

func nosimdRunner(gen, input, output matrix, numIn, numOut int) {
	size := len(input[0])
	start := 0
	unitSize := 16 * 1024 // concurrency unit size（Haswell， Skylake， Kabylake's L1 data cache size is 32KB)
	do := unitSize
	for start < size {
		if start+do <= size {
			nosimdWorker(gen, input, output, numIn, numOut, start, do)
			start = start + do
		} else {
			do = size - start
			nosimdWorker(gen, input, output, numIn, numOut, start, do)
			start = size
		}
	}
}

func nosimdWorker(gen, input, output matrix, numData, numParity, start, do int) {
	end := start + do
	for i := 0; i < numData; i++ {
		in := input[i]
		for oi := 0; oi < numParity; oi++ {
			if i == 0 {
				nosimdGfVectMul(gen[oi][i], in[start:end], output[oi][start:end])
			} else {
				nosimdGfVectMulXor(gen[oi][i], in[start:end], output[oi][start:end])
			}
		}
	}
}

func nosimdGfVectMul(c byte, in, out []byte) {
	mt := mulTable[c]
	for i := 0; i < len(in); i++ {
		out[i] = mt[in[i]]
	}
}

func nosimdGfVectMulXor(c byte, in, out []byte) {
	mt := mulTable[c]
	for i := 0; i < len(in); i++ {
		out[i] ^= mt[in[i]]
	}
}

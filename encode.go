package reedsolomon

import "errors"

// Encode : cauchy_matrix * data_matrix(input) -> parity_matrix(output)
// dp : data_matrix(upper) parity_matrix(lower, empty now)
func (r *rs) Encode(dp matrix) error {
	if len(dp) != r.shards {
		return ErrTooFewShards
	}
	size, err := checkShardSize(dp)
	if err != nil {
		return err
	}
	inMap := make(map[int]int)
	outMap := make(map[int]int)
	for i := 0; i < r.data; i++ {
		inMap[i] = i
	}
	for i := r.data; i < r.shards; i++ {
		outMap[i-r.data] = i
	}
	if r.ins == avx2 {
		encodeRunner(r.gen, dp, r.data, r.parity, size, inMap, outMap)
	} else {
		encodeRunnerS(r.gen, dp, r.data, r.parity, size, inMap, outMap)
	}
	return nil
}

func encodeRunner(gen, dp matrix, numIn, numOut, size int, inMap, outMap map[int]int) {
	start := 0
	unitSize := 16 * 1024 // concurrency unit size（Haswell， Skylake， Kabylake's L1 data cache size is 32KB)
	do := unitSize
	for start < size {
		if start+do <= size {
			encodeWorker(gen, dp, start, do, numIn, numOut, inMap, outMap)
			start = start + do
		} else {
			encodeRemain(start, size, gen, dp, numIn, numOut, inMap, outMap)
			start = size
		}
	}
}

func encodeRunnerS(gen, dp matrix, numIn, numOut, size int, inMap, outMap map[int]int) {
	start := 0
	unitSize := 16 * 1024 // concurrency unit size（Haswell， Skylake， Kabylake's L1 data cache size is 32KB)
	do := unitSize
	for start < size {
		if start+do <= size {
			encodeWorkerS(gen, dp, start, do, numIn, numOut, inMap, outMap)
			start = start + do
		} else {
			encodeRemainS(start, size, gen, dp, numIn, numOut, inMap, outMap)
			start = size
		}
	}
}

func encodeWorker(gen, dp matrix, start, do, numIn, numOut int, inMap, outMap map[int]int) {
	end := start + do
	for i := 0; i < numIn; i++ {
		j := inMap[i]
		in := dp[j]
		for oi := 0; oi < numOut; oi++ {
			k := outMap[oi]
			c := gen[oi][i]
			if i == 0 { // it means don't need to copy parity data for xor
				gfMulAVX2(mulTableLow[c][:], mulTableHigh[c][:], in[start:end], dp[k][start:end])
			} else {
				gfMulXorAVX2(mulTableLow[c][:], mulTableHigh[c][:], in[start:end], dp[k][start:end])
			}
		}
	}
}

func encodeRemain(start, size int, gen, dp matrix, numIn, numOut int, inMap, outMap map[int]int) {
	do := size - start
	for i := 0; i < numIn; i++ {
		j := inMap[i]
		in := dp[j]
		for oi := 0; oi < numOut; oi++ {
			k := outMap[oi]
			c := gen[oi][i]
			if i == 0 {
				gfMulRemain(c, in[start:size], dp[k][start:size], do)
			} else {
				gfMulRemainXor(c, in[start:size], dp[k][start:size], do)
			}
		}
	}
}

func encodeWorkerS(gen, dp matrix, start, do, numIn, numOut int, inMap, outMap map[int]int) {
	end := start + do
	for i := 0; i < numIn; i++ {
		j := inMap[i]
		in := dp[j]
		for oi := 0; oi < numOut; oi++ {
			k := outMap[oi]
			c := gen[oi][i]
			if i == 0 { // it means don't need to copy parity data for xor
				gfMulSSSE3(mulTableLow[c][:], mulTableHigh[c][:], in[start:end], dp[k][start:end])
			} else {
				gfMulXorSSSE3(mulTableLow[c][:], mulTableHigh[c][:], in[start:end], dp[k][start:end])
			}
		}
	}
}

func encodeRemainS(start, size int, gen, dp matrix, numIn, numOut int, inMap, outMap map[int]int) {
	do := size - start
	for i := 0; i < numIn; i++ {
		j := inMap[i]
		in := dp[j]
		for oi := 0; oi < numOut; oi++ {
			k := outMap[oi]
			c := gen[oi][i]
			if i == 0 {
				gfMulRemainS(c, in[start:size], dp[k][start:size], do)
			} else {
				gfMulRemainXorS(c, in[start:size], dp[k][start:size], do)
			}
		}
	}
}

var ErrShardSize = errors.New("reedsolomon: shards size equal 0 or not match")

func checkShardSize(m matrix) (int, error) {
	size := len(m[0])
	if size == 0 {
		return size, ErrShardSize
	}
	for _, v := range m {
		if len(v) != size {
			return 0, ErrShardSize
		}
	}
	return size, nil
}

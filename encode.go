package reedsolomon

import (
	"errors"
	"sync"

	"github.com/klauspost/cpuid"
)

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
	encodeRunner(r.gen, dp, r.data, r.parity, size, inMap, outMap)
	return nil
}

func encodeRunner(gen, dp matrix, numIn, numOut, size int, inMap, outMap map[int]int) {
	pipeline := cpuid.CPU.PhysicalCores
	offsets := make(chan [2]int, pipeline)
	wg := &sync.WaitGroup{}
	wg.Add(pipeline)
	for i := 0; i < pipeline; i++ {
		go encodeWorker(offsets, wg, gen, dp, numIn, numOut, inMap, outMap)
	}
	start := 0
	unitSize := 32768 // concurrency unit size（Haswell， Skylake， Kabylake's L1 data cache size)
	do := unitSize
	for start < size {
		if start+do <= size {
			offset := [2]int{start, do}
			offsets <- offset
			start = start + do
		} else {
			encodeRemain(start, size, gen, dp, numIn, numOut, inMap, outMap)
			start = size
		}
	}
	close(offsets)
	wg.Wait()
}

func encodeWorker(offsets chan [2]int, wg *sync.WaitGroup, gen, dp matrix, numIn, numOut int, inMap, outMap map[int]int) {
	defer wg.Done()
	for offset := range offsets {
		start := offset[0]
		do := offset[1]
		end := start + do
		for i := 0; i < numIn; i++ {
			j := inMap[i]
			in := dp[j]
			for oi := 0; oi < numOut; oi++ {
				k := outMap[oi]
				c := gen[oi][i]
				if i == 0 { // it means don't need to copy data from data for xor
					gfMulAVX2(mulTableLow[c][:], mulTableHigh[c][:], in[start:end], dp[k][start:end])
				} else {
					gfMulXorAVX2(mulTableLow[c][:], mulTableHigh[c][:], in[start:end], dp[k][start:end])
				}
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

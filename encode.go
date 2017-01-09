package reedsolomon

import (
	"errors"
	"runtime"
	"sync"
)

const unitSize = 32768 // concurrency unit size

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
	input := dp[0:r.data]
	output := dp[r.data:]
	encodeRunner(r.gen, input, output, r.data, r.parity, size)
	return nil
}

func encodeRunner(gen, input, output matrix, numIn, numOut, size int) {
	pipeline := runtime.GOMAXPROCS(0) / 2
	offsets := make(chan [2]int, pipeline)
	wg := &sync.WaitGroup{}
	wg.Add(pipeline)
	for i := 0; i < pipeline; i++ {
		go encodeWorker(offsets, wg, gen, input, output, numIn, numOut)
	}
	start := 0
	do := unitSize
	for start < size {
		if start+do <= size {
			offset := [2]int{start, do}
			offsets <- offset
			start = start + do
		} else {
			encodeRemain(start, size, gen, input, output, numIn, numOut)
			start = size
		}
	}
	close(offsets)
	wg.Wait()
}

func encodeRemain(start, size int, gen, input, output matrix, numIn, numOut int) {
	do := size - start
	for i := 0; i < numIn; i++ {
		in := input[i]
		for oi := 0; oi < numOut; oi++ {
			c := gen[oi][i]
			if i == 0 {
				gfMulRemain(c, in[start:size], output[oi][start:size], do)
			} else {
				gfMulRemainXor(c, in[start:size], output[oi][start:size], do)
			}
		}
	}
}

func encodeWorker(offsets chan [2]int, wg *sync.WaitGroup, gen, input, output matrix, numIn, numOut int) {
	defer wg.Done()
	for offset := range offsets {
		start := offset[0]
		do := offset[1]
		end := start + do
		for i := 0; i < numIn; i++ {
			in := input[i]
			for oi := 0; oi < numOut; oi++ {
				c := gen[oi][i]
				if i == 0 {
					gfMulAVX2(mulTableLow[c][:], mulTableHigh[c][:], in[start:end], output[oi][start:end])
				} else {
					gfMulXorAVX2(mulTableLow[c][:], mulTableHigh[c][:], in[start:end], output[oi][start:end])
				}
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

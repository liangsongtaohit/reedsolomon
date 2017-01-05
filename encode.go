package reedsolomon

import (
	"runtime"
	"sync"
)

var cpuINS int                       // avx2 = 0; ssse3 = 1; (!=avx2 && != ssse3) = 2
var pipeline = runtime.GOMAXPROCS(0) // number of goroutines for encoding or decoding
var unitSize int                     // flow size of per pipeline

// Encode : cauchy_matrix * data_matrix(input) -> parity_matrix(output)
// dp : data_matrix + parity_matrix(empty now)
func (r reedSolomon) Encode(dp matrix) error {
	// check args
	if len(dp) != r.shards {
		return errTooFewShards
	}
	size, err := checkShardLen(dp, r.shards)
	if err != nil {
		return err
	}

	// encoding
	input := dp[0:r.data]
	output := dp[r.data:]
	unitSize = calcUnit()
	if size < unitSize {
		unitSize = size
	}
	encodeRunner(r.gen, input, output, r.data, r.parity, size)
	return nil
}

func encodeRunner(gen, input, output matrix, numData, numParity, size int) {
	do := unitSize
	offsets := make(chan [2]int, pipeline)
	wg := &sync.WaitGroup{}
	wg.Add(pipeline)
	for i := 1; i <= pipeline; i++ {
		go encodeWorker(offsets, wg, gen, input, output, numData, numParity)
	}
	start := 0
	for start < size {
		if start+do > size {
			do = size - start
		}
		offset := [2]int{start, do}
		offsets <- offset
		start = start + do
	}
	close(offsets)
	wg.Wait()
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
				if i == 0 {
					gfVectMul(gen[oi][i], in[start:end], output[oi][start:end])
				} else {
					gfVectMulXor(gen[oi][i], in[start:end], output[oi][start:end])
				}
			}
		}
	}
}

func calcUnit() int {
	cpuINS = checkCPUINS()
	l1size := checkl1Size()
	if l1size != -1 {
		if cpuINS == 0 {
			f := (l1size / 32) * 32
			if l1size == f {
				return l1size
			} else {
				return f
			}
		} else if cpuINS == 1 {
			f := (l1size / 16) * 16
			if l1size == f {
				return l1size
			} else {
				return f
			}
		} else {
			return l1size // don't have avx2 or ssse3
		}
	} else { // can't get the cacheL1 data size
		return 32 * 1024
	}
}

package reedsolomon

import (
	"errors"
	"runtime"
	"sync"

	"github.com/klauspost/cpuid"
)

var cpuID int                        // avx2:0; ssse3:1; others: 2
var pipeline = runtime.GOMAXPROCS(0) // number of goroutines for encoding or decoding
var unitSize int                     // flow size of per pipeline

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

func encodeRunner(gen, input, output matrix, numData, numParity, size int) {
	unitSize = unit()
	if size < unitSize {
		unitSize = size
	}
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

func unit() int {
	cpuID = id()
	s := l1Size()
	if s != 1 {
		if cpuID == 0 {
			x := (s / 32) * 32
			if x == s {
				return s
			}
			return x
		}
		if cpuID == 1 {
			x := (s / 16) * 16
			if x == s {
				return s
			}
			return x
		}
		return s
	}
	return 32 * 1024
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

func id() int {
	if avx2() {
		return 0
	}
	if ssse3() {
		return 1
	}
	return 2
}

//go:noescape
func ssse3() bool

//go:noescape
func avx2() bool

func l1Size() int {
	return cpuid.CPU.Cache.L1D
}

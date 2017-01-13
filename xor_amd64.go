package reedsolomon

import (
	"sync"
	"github.com/klauspost/cpuid"
)
func Xor(dp matrix) error {
	size, err := checkShardSize(dp)
	if err != nil {
		return err
	}
	data := len(dp)
	input := dp[:data-1]
	output := dp[0]
	xorRunner(input, output, size)
	return nil
}

func xorRunner(input matrix, output []byte, size int) {
	pipeline := cpuid.CPU.PhysicalCores
	offsets := make(chan [2]int, pipeline)
	wg := &sync.WaitGroup{}
	wg.Add(pipeline)
	for i := 0; i < pipeline; i++ {
		go xorWorker(offsets, wg, input, output)
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
			xorRemain(start, size, input, output)
			start = size
		}
	}
	close(offsets)
	wg.Wait()
}

func xorWorker(offsets chan [2]int, wg *sync.WaitGroup, input matrix, output []byte) {
	defer wg.Done()
	for offset := range offsets {
		start := offset[0]
		do := offset[1]
		end := start + do
		for i := 1; i < len(input); i++ {
			in := input[i]
			xorAVX2(in[start:end], output[start:end])
		}
	}
}

func xorRemain(start, size int, input matrix, output []byte) {
	do := size - start
	for i := 1; i < len(input); i++ {
		in := input[i]
		xorRemainAVX2(in[start:size], output[start:size], do)
	}
}


//go:noescape
func xorAVX2(in, out []byte)

func xorRemainAVX2(input, output []byte, size int) {
	var done int
	if size < 32 {
		for i, v := range output {
			v = v ^ input[i]
			output[i] = v
		}
	} else {
		xorAVX2(input, output)
		done = (size >> 5) << 5
		remain := size - done
		if remain > 0 {
			for i := done; i < size; i++ {
				v0 := output[i]
				v1 := input[i]
				v := v0 ^ v1
				output[i] = v
			}
		}
	}
}
package reedsolomon

import (
	"runtime"
	"sync"
)

func encodeNormal(gen, input, output matrix, numIn, numOut, size, unit int) {
	pipeline := runtime.GOMAXPROCS(0)
	offsets := make(chan [2]int, pipeline)
	wg := &sync.WaitGroup{}
	wg.Add(pipeline)
	for i := 1; i <= pipeline; i++ {
		go encodeNormalWorker(offsets, wg, gen, input, output, numIn, numOut)
	}

	start := 0
	do := unit
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

func encodeNormalWorker(offsets chan [2]int, wg *sync.WaitGroup, gen, input, output matrix, numIn, numOut int) {

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
					galMulSlice(c, in[start:end], output[oi][start:end])
				} else {
					galMulSliceXor(c, in[start:end], output[oi][start:end])
				}
			}
		}
	}
}

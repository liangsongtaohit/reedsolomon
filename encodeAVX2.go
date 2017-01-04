/**
 * Reed-Solomon Coding over in GF(2^8).
 * Primitive Polynomial: x^8 + x^4 + x^3 + x^2 + 1 (0x1d)
 *
 * Copyright 2017, Templexxx
 * Copyright 2015, Klaus Post
 * Copyright 2015, Backblaze, Inc.
 */

package reedsolomon

import (
	"sync"
	"runtime"
)

func encodeAVX2(gen, input, output matrix, numIn, numOut, size, unit int) {
	if size < 32 {
		encodeAVX2Remain(0, size, gen, input, output, numIn, numOut)
	} else {
		pipeline := runtime.GOMAXPROCS(0)
		offsets := make(chan [2]int, pipeline)
		wg := &sync.WaitGroup{}
		wg.Add(pipeline)
		for i := 1; i <= pipeline; i++ {
			go encodeAVX2Worker(offsets, wg, gen, input, output, numIn, numOut)
		}
		start := 0
		do := unit
		for start < size {
			if start + do <= size {
				offset := [2]int{start, do}
				offsets <- offset
				start = start + do
			} else {
				encodeAVX2Remain(start, size, gen, input, output, numIn, numOut)
				start = size
			}
		}
		close(offsets)
		wg.Wait()
	}
}

func encodeAVX2Remain(start, size int, gen, input, output matrix, numIn, numOut int) {
	do := size - start
	for i := 0; i < numIn; i++ {
		in := input[i]
		for oi := 0; oi < numOut; oi++ {
			c := gen[oi][i]
			if i == 0 {
				gfMulRemainAVX2(c, in[start:size], output[oi][start:size], do)
			} else {
				gfMulRemainXorAVX2(c, in[start:size], output[oi][start:size], do)
			}
		}
	}
}

func encodeAVX2Worker(offsets chan [2]int, wg *sync.WaitGroup, gen, input, output matrix, numIn, numOut int) {
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
//func encodeNormal(encodingMatrix, inputs, outputs matrix, numInput, numOutput, size int) {
//
//	offsets := make(chan [2]int, worker)
//	wg := &sync.WaitGroup{}
//	wg.Add(worker)
//	for i := 1; i <= worker; i++ {
//		go encodeNormalWorker(offsets, wg, encodingMatrix, inputs, outputs, numInput, numOutput)
//	}
//
//	start := 0
//	do := sizePerRoutine
//	for start < size {
//		if start+do > size {
//			do = size - start
//		}
//		offset := [2]int{start, do}
//		offsets <- offset
//		start = start + do
//	}
//	close(offsets)
//	wg.Wait()
//}



//func encodeNormalWorker(offsets chan [2]int, wg *sync.WaitGroup, encodingMatrix, inputs, outputs matrix, numInput, numOutput int) {
//
//	defer wg.Done()
//	for offset := range offsets {
//		start := offset[0]
//		do := offset[1]
//		end := start + do
//		for i := 0; i < numInput; i++ {
//			input := inputs[i]
//			for oi := 0; oi < numOutput; oi++ {
//				ei := encodingMatrix[oi][i]
//				if i == 0 {
//					gfMulSliceAVX2(ei, input[start:end], outputs[oi][start:end], do)
//				} else {
//					gfMulSliceXorAVX2(ei, input[start:end], outputs[oi][start:end], do)
//				}
//			}
//		}
//	}
//}

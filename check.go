package reedsolomon

import "github.com/klauspost/cpuid"

// check each row in a matrix is equal or not
func checkShardLen(m matrix, numShards int) (size int, err error) {

	size = len(m[0])
	if size == 0 {
		return size, errShardNoData
	}
	for i := 1; i < numShards; i++ {
		if len(m[i]) != size {
			return size, errShardSize
		}
	}
	return
}

//go:noescape
func CheckSSSE3() bool

//go:noescape
func CheckAVX2() bool

func checkl1Size() int {
	return cpuid.CPU.Cache.L1D
}

func checkCPUINS() int {

	if CheckAVX2() {
		return 0
	}
	if CheckSSSE3() {
		return 1
	}
	return 2
}

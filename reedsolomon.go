/**
 * Reed-Solomon Coding over in GF(2^8).
 * Primitive Polynomial: x^8 + x^4 + x^3 + x^2 + 1 (0x1d)
 *
 * Copyright 2017, Templexxx
 * Copyright 2015, Klaus Post
 * Copyright 2015, Backblaze, Inc.
 */
package reedsolomon

type reedSolomon struct {
	data   int    // Number of data shards, should not be modified.
	parity int    // Number of parity shards, should not be modified.
	shards int    // Total number of shards. Calculated, and should not be modified.
	m      matrix // encoding matrix, identity matrix(upper) + generator matrix(lower)
	gen    matrix // generator matrix, cauchy here
}

// New : create a encoding matrix for encoding, reconstruction
func New(d, p int) (*reedSolomon, error) {
	r := reedSolomon{
		data:   d,
		parity: p,
		shards: d + p,
	}
	// check argument
	if (d <= 0) || (p <= 0) || (r.shards >= 255) {
		return nil, errInvalidNumShards
	}

	// create encoding matrix
	e, err := genEncodeMatrix(r.shards, d)
	if err != nil {
		return nil, err
	}
	r.m = e

	// TODO do I need make a new slice?
	r.gen = make(matrix, p)
	for i := range r.gen {
		r.gen[i] = r.m[d+i]
	}

	return &r, err
}

//var cpuINS int  // avx2 = 0; ssse3 = 1; (!=avx2 && != ssse3) = 2
//var pipeline = runtime.GOMAXPROCS(0)    // number of goroutines for encoding or decoding
//var unitSize int // flow size of per pipeline

// Encode : cauchy_matrix * data_matrix(input) -> parity_matrix(output)
// dp : data_matrix + parity_matrix(empty now)
func (r reedSolomon) Encode(dp matrix, cpuINS int) error {
	// check matrix row number
	if len(dp) != r.shards {
		return errTooFewShards
	}
	// check shard length
	// set the first shard' len as the flag
	size, err := checkShardLen(dp, r.shards)
	if err != nil {
		return err
	}

	// encoding
	input := dp[0:r.data]
	output := dp[r.data:]

	//cpuINS := checkCPUINS()
	//unitSize := calcUnit()
	unitSize := 32768
	if size < unitSize {
		unitSize = size
	}
	if cpuINS == 0 {
	encodeAVX2(r.gen, input, output, r.data, r.parity, size, unitSize)
		//fmt.Println(output)
	//} else {
	//	encodeNormal(r.Cauchy, data, parity, r.NumData, r.NumParity, size)
	//}
	//if cpuINS == 0 {

	//} else if cpuINS == 1 {
	//	encodeSSSE3(r.Cauchy, data, parity, r.NumData, r.NumParity, size)
	} else {
		encodeNormal(r.gen, input, output, r.data, r.parity, size, unitSize)
	}
	return nil
}



//func calcUnit() int {
//	l1size := checkl1Size()
//	if l1size != -1 {
//		if cpuINS == 0 {
//			f := (l1size / 32) * 32
//			if l1size == f {
//				return l1size
//			} else {
//				return f
//			}
//		} else if cpuINS == 1 {
//			f := (l1size / 16) * 16
//			if l1size == f {
//				return l1size
//			} else {
//				return f
//			}
//		} else {
//			return l1size
//		}
//	} else {    // can't get the cacheL1 data size
//		return 32 * 1024
//	}
//}
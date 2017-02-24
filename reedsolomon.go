/**
 * Reed-Solomon Coding over in GF(2^8).
 * Primitive Polynomial: x^8 + x^4 + x^3 + x^2 + 1 (0x1d)
 */

package reedsolomon

import "errors"

type rs struct {
	data   int    // Number of data shards
	parity int    // Number of parity shards
	shards int    // Total number of shards
	m      matrix // encoding matrix, identity matrix(upper) + generator matrix(lower)
	gen    matrix // generator matrix(cauchy matrix)
	ins    int    // Extensions Instruction(avx2 ssse3)
}

const (
	avx2  = 0
	ssse3 = 1
)

var ErrTooFewShards = errors.New("reedsolomon: too few shards given for encoding")
var ErrTooManyShards = errors.New("reedsolomon: too many shards given for encoding")
var ErrNoSupportINS = errors.New("reedsolomon: no avx2 or ssse3")

// New : create a encoding matrix for encoding, reconstruction
func New(d, p int) (*rs, error) {
	err := checkShards(d, p)
	if err != nil {
		return nil, err
	}
	r := rs{
		data:   d,
		parity: p,
		shards: d + p,
	}
	if hasAVX2() {
		r.ins = avx2
	} else if hasSSSE3() {
		r.ins = ssse3
	} else {
		return &r, ErrNoSupportINS
	}
	e := genEncodeMatrix(r.shards, d) // create encoding matrix
	r.m = e
	r.gen = NewMatrix(p, d)
	for i := range r.gen {
		r.gen[i] = r.m[d+i]
	}
	return &r, err
}

func checkShards(d, p int) error {
	if (d <= 0) || (p <= 0) {
		return ErrTooFewShards
	}
	if d+p >= 255 {
		return ErrTooManyShards
	}
	return nil
}

//go:noescape
func hasAVX2() bool

//go:noescape
func hasSSSE3() bool

package reedsolomon

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestEncode(t *testing.T) {
	size := 50000
	r, err := New(10, 3)
	if err != nil {
		t.Fatal(err)
	}
	dp := newMatrix(13, size)
	rand.Seed(0)
	for s := 0; s < 13; s++ {
		fillRandom(dp[s])
	}
	err = r.Encode(dp)
	if err != nil {
		t.Fatal(err)
	}
	badDP := newMatrix(13, 100)
	badDP[0] = make([]byte, 1)
	err = r.Encode(badDP)
	if err != ErrShardSize {
		t.Errorf("expected %v, got %v", ErrShardSize, err)
	}
}

func TestVerifyEncode(t *testing.T) {
	r, err := New(5, 5)
	if err != nil {
		t.Fatal(err)
	}
	shards := [][]byte{
		{0, 1},
		{4, 5},
		{2, 3},
		{6, 7},
		{8, 9},
		{0, 0},
		{0, 0},
		{0, 0},
		{0, 0},
		{0, 0},
	}
	r.Encode(shards)
	if shards[5][0] != 97 || shards[5][1] != 64 {
		t.Fatal("shard 5 mismatch")
	}
	if shards[6][0] != 173 || shards[6][1] != 3 {
		t.Fatal("shard 6 mismatch")
	}
	if shards[7][0] != 218 || shards[7][1] != 14 {
		t.Fatal("shard 7 mismatch")
	}
	if shards[8][0] != 107 || shards[8][1] != 35 {
		t.Fatal("shard 8 mismatch")
	}
	if shards[9][0] != 110 || shards[9][1] != 177 {
		t.Fatal("shard 9 mismatch")
	}
}

func TestReconst(t *testing.T) {
	size := 5111
	r, err := New(10, 3)
	if err != nil {
		t.Fatal(err)
	}
	dp := newMatrix(13, size)
	rand.Seed(0)
	for s := 0; s < 10; s++ {
		fillRandom(dp[s])
	}
	err = r.Encode(dp)
	if err != nil {
		t.Fatal(err)
	}
	// restore encode result
	store := newMatrix(3, size)
	store[0] = dp[0]
	store[1] = dp[4]
	store[2] = dp[12]
	// Reconstruct with all dp present
	var lost []int
	err = r.Reconst(dp, lost, true)
	if err != nil {
		t.Fatal(err)
	}
	// 3 dp "missing"
	lost = append(lost, 4)
	lost = append(lost, 0)
	lost = append(lost, 12)
	err = r.Reconst(dp, lost, true)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(store[0], dp[0]) {
		t.Fatal("reconst data mismatch: dp[0]")
	}
	if !bytes.Equal(store[1], dp[4]) {
		t.Fatal("reconst data mismatch: dp[4]")
	}
	if !bytes.Equal(store[2], dp[12]) {
		t.Fatal("reconst data mismatch: dp[12]")
	}
	// Reconstruct with 9 dp present (should fail)
	lost = append(lost, 11)
	err = r.Reconst(dp, lost, true)
	if err != ErrTooFewShards {
		t.Errorf("expected %v, got %v", ErrTooFewShards, err)
	}
}

func TestASM(t *testing.T) {
	d := 10
	p := 4
	size := 9999
	r, err := New(d, p)
	if err != nil {
		t.Fatal(err)
	}
	// asm
	dp := newMatrix(d+p, size)
	rand.Seed(0)
	for i := 0; i < d; i++ {
		fillRandom(dp[i])
	}
	err = r.Encode(dp)
	if err != nil {
		t.Fatal(err)
	}
	// mulTable
	mDP := newMatrix(d+p, size)
	for i := 0; i < d; i++ {
		mDP[i] = dp[i]
	}
	err = r.noasmEncode(mDP)
	if err != nil {
		t.Fatal(err)
	}
	for i, asm := range dp {
		if !bytes.Equal(asm, mDP[i]) {
			t.Fatal("verify asm failed, no match noasm version")
		}
	}
}

// Benchmark 10 data shards and 4 parity shards with 64KB each.
func BenchmarkEncode10x4x64K(b *testing.B) {
	benchmarkEncode(b, 10, 4, 64*1024)
}

func BenchmarkEncode10x4x128K(b *testing.B) {
	benchmarkEncode(b, 10, 4, 128*1024)
}

func BenchmarkEncode10x4x256K(b *testing.B) {
	benchmarkEncode(b, 10, 4, 256*1024)
}

func BenchmarkEncode10x4x512K(b *testing.B) {
	benchmarkEncode(b, 10, 4, 512*1024)
}

func BenchmarkEncode10x4x1M(b *testing.B) {
	benchmarkEncode(b, 10, 4, 1024*1024)
}

func BenchmarkEncode10x4x16M(b *testing.B) {
	benchmarkEncode(b, 10, 4, 16*1024*1024)
}

func BenchmarkEncode28x4x1M(b *testing.B) {
	benchmarkEncode(b, 28, 4, 1024*1024)
}

func BenchmarkEncode28x4x16M(b *testing.B) {
	benchmarkEncode(b, 28, 4, 16*1024*1024)
}

func benchmarkEncode(b *testing.B, data, parity, size int) {
	r, err := New(data, parity)
	if err != nil {
		b.Fatal(err)
	}
	dp := newMatrix(data+parity, size)
	rand.Seed(0)
	for i := 0; i < data; i++ {
		fillRandom(dp[i])
	}
	b.SetBytes(int64(size * data))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = r.Encode(dp)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncode10x1x1M(b *testing.B) {
	benchmarkEncode(b, 10, 1, 1024*1024)
}

func BenchmarkReconst10x4x1M(b *testing.B) {
	benchmarkReconst(b, 10, 4, 1024*1024)
}

func BenchmarkEncode10x1x16M(b *testing.B) {
	benchmarkEncode(b, 10, 16, 1024*1024)
}

func BenchmarkReconst10x4x16M(b *testing.B) {
	benchmarkReconst(b, 10, 4, 16*1024*1024)
}

func benchmarkReconst(b *testing.B, d, p, size int) {
	r, err := New(d, p)
	if err != nil {
		b.Fatal(err)
	}
	dp := newMatrix(d+p, size)
	rand.Seed(0)
	for s := 0; s < d; s++ {
		fillRandom(dp[s])
	}
	err = r.Encode(dp)
	if err != nil {
		b.Fatal(err)
	}
	var lost []int
	shardsToCorrupt := rand.Intn(p)
	for i := 1; i <= shardsToCorrupt; i++ {
		lost = append(lost, rand.Intn(d+p))
	}
	for _, l := range lost {
		dp[l] = make([]byte, size)
	}
	b.SetBytes(int64(size * d))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = r.Reconst(dp, lost, true)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func fillRandom(p []byte) {
	for i := 0; i < len(p); i += 7 {
		val := rand.Int63()
		for j := 0; i+j < len(p) && j < 7; j++ {
			p[i+j] = byte(val)
			val >>= 8
		}
	}
}

func (r *rs) noasmEncode(dp matrix) error {
	// check args
	if len(dp) != r.shards {
		return ErrTooFewShards
	}
	_, err := checkShardSize(dp)
	if err != nil {
		return err
	}
	// encoding
	input := dp[0:r.data]
	output := dp[r.data:]
	noasmRunner(r.gen, input, output, r.data, r.parity)
	return nil
}

func noasmRunner(gen, input, output matrix, numData, numParity int) {
	for i := 0; i < numData; i++ {
		in := input[i]
		for oi := 0; oi < numParity; oi++ {
			if i == 0 {
				noasmGfVectMul(gen[oi][i], in[:], output[oi][:])
			} else {
				noasmGfVectMulXor(gen[oi][i], in[:], output[oi][:])
			}
		}
	}
}

func noasmGfVectMul(c byte, in, out []byte) {
	mt := mulTable[c]
	for i := 0; i < len(in); i++ {
		out[i] = mt[in[i]]
	}
}

func noasmGfVectMulXor(c byte, in, out []byte) {
	mt := mulTable[c]
	for i := 0; i < len(in); i++ {
		out[i] ^= mt[in[i]]
	}
}

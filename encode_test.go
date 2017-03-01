package reedsolomon

import (
	"bytes"
	"math/rand"
	"testing"
  "fmt"
  "time"
)

func TestEncode(t *testing.T) {
	size := 50000
	r, err := New(10, 3)
	if err != nil {
		t.Fatal(err)
	}
	dp := NewMatrix(13, size)
	rand.Seed(0)
	for s := 0; s < 13; s++ {
		fillRandom(dp[s])
	}
	err = r.Encode(dp)
	if err != nil {
		t.Fatal(err)
	}
	badDP := NewMatrix(13, 100)
	badDP[0] = make([]byte, 1)
	err = r.Encode(badDP)
	if err != ErrShardSize {
		t.Errorf("expected %v, got %v", ErrShardSize, err)
	}
}

// test low, high table work
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

// test avx2, if don't have it will test ssse3
func TestASM(t *testing.T) {
	d := 10
	p := 4
	size := 65 * 1024
	r, err := New(d, p)
	if err != nil {
		t.Fatal(err)
	}
	// asm
	dp := NewMatrix(d+p, size)
	rand.Seed(0)
	for i := 0; i < d; i++ {
		fillRandom(dp[i])
	}
	err = r.Encode(dp)
	if err != nil {
		t.Fatal(err)
	}
	// mulTable
	mDP := NewMatrix(d+p, size)
	for i := 0; i < d; i++ {
		mDP[i] = dp[i]
	}
	err = r.noasmEncode(mDP)
	if err != nil {
		t.Fatal(err)
	}
	for i, asm := range dp {
		if !bytes.Equal(asm, mDP[i]) {
			t.Fatal("verify asm failed, no match noasm version; shards: ", i)
		}
	}
}

func TestSSSE3(t *testing.T) {
	d := 10
	p := 4
	size := 65 * 1024
	r, err := New(d, p)
	if err != nil {
		t.Fatal(err)
	}
	r.ins = ssse3
	// asm
	dp := NewMatrix(d+p, size)
	rand.Seed(0)
	for i := 0; i < d; i++ {
		fillRandom(dp[i])
	}
	err = r.Encode(dp)
	if err != nil {
		t.Fatal(err)
	}
	// mulTable
	mDP := NewMatrix(d+p, size)
	for i := 0; i < d; i++ {
		mDP[i] = dp[i]
	}
	err = r.Encode(mDP)
	if err != nil {
		t.Fatal(err)
	}
	for i, asm := range dp {
		if !bytes.Equal(asm, mDP[i]) {
			t.Fatal("verify asm failed, no match noasm version; shards: ", i)
		}
	}
}

//func BenchmarkEncode10x4x8K(b *testing.B) {
//benchmarkEncode(b, 10, 4, 8*1024)
//}

//func BenchmarkEncode10x4x16K(b *testing.B) {
//benchmarkEncode(b, 10, 4, 16*1024)
//}

//func BenchmarkEncode10x4x32K(b *testing.B) {
//benchmarkEncode(b, 10, 4, 32*1024)
///}

//func BenchmarkEncode10x4x64K(b *testing.B) {
//	benchmarkEncode(b, 10, 4, 64*1024)
//}

//func BenchmarkEncode10x4x128K(b *testing.B) {
//benchmarkEncode(b, 10, 4, 128*1024)
//}

//func BenchmarkEncode10x4x256K(b *testing.B) {
//benchmarkEncode(b, 10, 4, 256*1024)
//}

//func BenchmarkEncode10x4x512K(b *testing.B) {
//benchmarkEncode(b, 10, 4, 512*1024)
//}

//func BenchmarkEncode10x4x1M(b *testing.B) {
	//benchmarkEncode(b, 10, 4, 1024*1024)
//}

//func BenchmarkEncode10x4x4M(b *testing.B) {
	//benchmarkEncode(b, 10, 4, 4*1024*1024)
//}

func BenchmarkEncode10x4x16M(b *testing.B) {
	benchmarkEncode(b, 10, 4, 16*1024*1024)
}

//func BenchmarkEncode17x3x1M(b *testing.B) {
//	benchmarkEncode(b, 17, 3, 1024*1024)
//}

//func BenchmarkEncode17x3x4M(b *testing.B) {
//	benchmarkEncode(b, 17, 3, 4*1024*1024)
//}

//func BenchmarkEncode17x3x16M(b *testing.B) {
//	benchmarkEncode(b, 17, 3, 16*1024*1024)
//}

//func BenchmarkEncode28x4x1M(b *testing.B) {
//	benchmarkEncode(b, 28, 4, 1024*1024)
//}

//func BenchmarkEncode28x4x4M(b *testing.B) {
//	benchmarkEncode(b, 28, 4, 4*1024*1024)
//}
func BenchmarkEncode10x4x4K(b *testing.B) {
  benchmarkEncode(b, 10, 4, 4*1024)
}

func BenchmarkEncode28x4x16M(b *testing.B) {
	benchmarkEncode(b, 28, 4, 16*1024*1024)
}

func BenchmarkEncode14x10x1M(b *testing.B) {
	benchmarkEncode(b, 14, 10, 1024*1024)
}

func BenchmarkEncode14x10x4M(b *testing.B) {
	benchmarkEncode(b, 14, 10, 4*1024*1024)
}

func BenchmarkEncode14x10x16M(b *testing.B) {
	benchmarkEncode(b, 14, 10, 16*1024*1024)
}

func TestISA10x4x16M(t *testing.T) {
  ISAway(t, 10, 4, 16*1024*1024)
}
func ISAway(t *testing.T, data, parity, size int) {
  r, err := New(data, parity)
  if err != nil {
    t.Fatal(err)
  }
  dp := NewMatrix(data+parity, size)
  rand.Seed(0)
  for i := 0; i < data; i++ {
    fillRandom(dp[i])
  }
  now := time.Now()
  for i := 0; i < 30; i++ {
    r.Encode(dp)
  }
  cost := time.Since(now)
  io := data * size * 30 / 1000 / 1000
  fmt.Println(cost, io)
}
func benchmarkEncode(b *testing.B, data, parity, size int) {
	r, err := New(data, parity)
	if err != nil {
		b.Fatal(err)
	}
	dp := NewMatrix(data+parity, size)
	rand.Seed(0)
	for i := 0; i < data; i++ {
		fillRandom(dp[i])
	}
	b.SetBytes(int64(size * data))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Encode(dp)
	}
}

func BenchmarkSSSE3Encode28x4x16M(b *testing.B) {
	benchmarkSSSE3Encode(b, 28, 4, 16*1024*1024)
}

func benchmarkSSSE3Encode(b *testing.B, data, parity, size int) {
	r, err := New(data, parity)
	r.ins = ssse3
	if err != nil {
		b.Fatal(err)
	}
	dp := NewMatrix(data+parity, size)
	rand.Seed(0)
	for i := 0; i < data; i++ {
		fillRandom(dp[i])
	}
	b.SetBytes(int64(size * data))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Encode(dp)
	}
}

func BenchmarkNOASMEncode28x4x16M(b *testing.B) {
	benchmarkNOASMEncode(b, 28, 4, 16*1024*1024)
}

func benchmarkNOASMEncode(b *testing.B, data, parity, size int) {
	r, err := New(data, parity)
	if err != nil {
		b.Fatal(err)
	}
	dp := NewMatrix(data+parity, size)
	rand.Seed(0)
	for i := 0; i < data; i++ {
		fillRandom(dp[i])
	}
	b.SetBytes(int64(size * data))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.noasmEncode(dp)
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

// test no simd asm
func TestNoSIMD(t *testing.T) {
	d := 10
	p := 1
	size := 10
	r, err := New(d, p)
	if err != nil {
		t.Fatal(err)
	}
	// asm
	dp := NewMatrix(d+p, size)
	rand.Seed(0)
	for i := 0; i < d; i++ {
		fillRandom(dp[i])
	}
	err = r.nosimdEncode(dp)
	if err != nil {
		t.Fatal(err)
	}
	// mulTable
	mDP := NewMatrix(d+p, size)
	for i := 0; i < d; i++ {
		mDP[i] = dp[i]
	}
	err = r.noasmEncode(mDP)
	if err != nil {
		t.Fatal(err)
	}
	for i, asm := range dp {
		if !bytes.Equal(asm, mDP[i]) {
			t.Fatal("verify simd failed, no match noasm version; shards: ", i)
		}
	}
}

func BenchmarkNOSIMDncode28x4x16M(b *testing.B) {
	benchmarkNOSIMDEncode(b, 28, 4, 16*1024*1024)
}

func benchmarkNOSIMDEncode(b *testing.B, data, parity, size int) {
	r, err := New(data, parity)
	if err != nil {
		b.Fatal(err)
	}
	dp := NewMatrix(data+parity, size)
	rand.Seed(0)
	for i := 0; i < data; i++ {
		fillRandom(dp[i])
	}
	b.SetBytes(int64(size * data))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.nosimdEncode(dp)
	}
}

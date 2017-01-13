package reedsolomon

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestXorASM(t *testing.T) {
	d := 10
	size := 9999
	// asm
	dp := newMatrix(d+1, size)
	rand.Seed(0)
	for i := 0; i < d; i++ {
		fillRandom(dp[i])
	}
	err := Xor(dp)
	if err != nil {
		t.Fatal(err)
	}
	// normal
	mDP := newMatrix(d+1, size)
	for i := 0; i < d; i++ {
		mDP[i] = dp[i]
	}
	noasmXor(mDP)
	for i, asm := range dp {
		if !bytes.Equal(asm, mDP[i]) {
			t.Fatal("verify asm failed, no match noasm version; shards: ", i)
		}
	}
}

func BenchmarkEncode10x1x16M(b *testing.B) {
	benchmarkEncode(b, 10, 1, 16*1024*1024)
}

func BenchmarkXor10x16M(b *testing.B) {
	benchmarkXor(b, 10,16*1024*1024)
}

func BenchmarkEncode28x1x16M(b *testing.B) {
	benchmarkEncode(b, 28, 1, 16*1024*1024)
}

func BenchmarkXor28x16M(b *testing.B) {
	benchmarkXor(b, 28,  16*1024*1024)
}

func benchmarkXor(b *testing.B, data, size int) {
	dp := newMatrix(data+1, size)
	rand.Seed(0)
	for i := 0; i < data; i++ {
		fillRandom(dp[i])
	}
	b.SetBytes(int64(size * data))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := Xor(dp)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func noasmXor(dp matrix) {

	// encoding
	data := len(dp)
	input := dp[:data-1]
	output := dp[0]
	noasmXorRunner(input, output)
}

func noasmXorRunner(input matrix, output []byte) {
	for i := 1; i < len(input); i++ {
		in := input[i]
		for i, v := range output {
			v = v ^ in[i]
			output[i] = v
		}
	}
}


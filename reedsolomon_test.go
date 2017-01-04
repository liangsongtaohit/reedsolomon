package reedsolomon

import (
	"math/rand"
	"testing"
	//"golang.org/x/tools/go/gcimporter15/testdata"
)

func TestEncoding(t *testing.T) {
	perShard := 50000
	r, err := New(10, 3)
	if err != nil {
		t.Fatal(err)
	}
	shards := make([][]byte, 13)
	for s := range shards {
		shards[s] = make([]byte, perShard)
	}
	
	rand.Seed(0)
	for s := 0; s < 13; s++ {
		fillRandom(shards[s])
	}
	
	err = r.Encode(shards)
	if err != nil {
		t.Fatal(err)
	}
	
	//
	//badShards := make([][]byte, 13)
	//badShards[0] = make([]byte, 1)
	//err = r.Encode(badShards)
	//if err != errShardSize {
	//	t.Errorf("expected %v, got %v", errShardSize, err)
	//}
}

//func TestASMEncode(t *testing.T) {
//	NumData := 10
//	NumParity := 4
//	shardSize := 512 * 1024
//	r, err := New(NumData, NumParity)
//	if err != nil {
//		t.Fatal(err)
//	}
//	shards := make([][]byte, NumData+NumParity)
//	for s := range shards {
//		shards[s] = make([]byte, shardSize)
//	}
//
//	rand.Seed(0)
//	for s := 0; s < NumData; s++ {
//		fillRandom(shards[s])
//	}
//
//	b.SetBytes(int64(shardSize * dataShards))
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		err = r.Encode(shards)
//		if err != nil {
//			b.Fatal(err)
//		}
//	}
//}

func TestOneEncode(t *testing.T) {
	codec, err := New(5, 5)
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
	codec.Encode(shards)
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

func fillRandom(p []byte) {
	for i := 0; i < len(p); i += 7 {
		val := rand.Int63()
		for j := 0; i+j < len(p) && j < 7; j++ {
			p[i+j] = byte(val)
			val >>= 8
		}
	}
}

func benchmarkEncode(b *testing.B, dataShards, parityShards, shardSize int) {
	r, err := New(dataShards, parityShards)
	if err != nil {
		b.Fatal(err)
	}
	shards := make([][]byte, dataShards+parityShards)
	for s := range shards {
		shards[s] = make([]byte, shardSize)
	}
	
	rand.Seed(0)
	for s := 0; s < dataShards; s++ {
		fillRandom(shards[s])
	}
	
	b.SetBytes(int64(shardSize * dataShards))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = r.Encode(shards)
		if err != nil {
			b.Fatal(err)
		}
	}
}



// Benchmark 10 data shards and 4 parity shards with 32KB each.
func BenchmarkEncode10x4x32K(b *testing.B) {
	benchmarkEncode(b, 10, 4, 32*1024)
}

// Benchmark 10 data shards and 4 parity shards with 64KB each.
func BenchmarkEncode10x4x64K(b *testing.B) {
	benchmarkEncode(b, 10, 4, 64*1024)
}

// Benchmark 10 data shards and 4 parity shards with 128KB each.
func BenchmarkEncode10x4x128K(b *testing.B) {
	benchmarkEncode(b, 10, 4, 128*1024)
}

// Benchmark 10 data shards and 4 parity shards with 256KB each.
func BenchmarkEncode10x4x256K(b *testing.B) {
	benchmarkEncode(b, 10, 4, 256*1024)
}

// Benchmark 10 data shards and 4 parity shards with 512KB each.
func BenchmarkEncode10x4x512K(b *testing.B) {
	benchmarkEncode(b, 10, 4, 512*1024)
}

// Benchmark 10 data shards and 4 parity shards with 1MB each.
func BenchmarkEncode10x4x1M(b *testing.B) {
	benchmarkEncode(b, 10, 4, 1024*1024)
}

// Benchmark 10 data shards and 4 parity shards with 16MB each.
func BenchmarkEncode10x4x16M(b *testing.B) {
	benchmarkEncode(b, 10, 4, 16*1024*1024)
}

// Benchmark 28 data shards and 4 parity shards with 32KB each.
func BenchmarkEncode28x4x32K(b *testing.B) {
	benchmarkEncode(b, 28, 4, 32*1024)
}

// Benchmark 28 data shards and 4 parity shards with 64KB each.
func BenchmarkEncode28x4x64K(b *testing.B) {
	benchmarkEncode(b, 28, 4, 64*1024)
}

// Benchmark 28 data shards and 4 parity shards with 128KB each.
func BenchmarkEncode28x4x128K(b *testing.B) {
	benchmarkEncode(b, 28, 4, 128*1024)
}

// Benchmark 28 data shards and 4 parity shards with 256KB each.
func BenchmarkEncode28x4x256K(b *testing.B) {
	benchmarkEncode(b, 28, 4, 256*1024)
}

// Benchmark 28 data shards and 4 parity shards with 512KB each.
func BenchmarkEncode28x4x512K(b *testing.B) {
	benchmarkEncode(b, 28, 4, 512*1024)
}

// Benchmark 28 data shards and 4 parity shards with 1MB each.
func BenchmarkEncode28x4x1M(b *testing.B) {
	benchmarkEncode(b, 28, 4, 1024*1024)
}

// Benchmark 28 data shards and 4 parity shards with 16MB each.
func BenchmarkEncode28x4x16M(b *testing.B) {
	benchmarkEncode(b, 28, 4, 16*1024*1024)
}

// Benchmark 14 data shards and 10 parity shards with 32KB each.
func BenchmarkEncode14x10x32K(b *testing.B) {
	benchmarkEncode(b, 14, 10, 32*1024)
}

// Benchmark 14 data shards and 10 parity shards with 1MB each.
func BenchmarkEncode14x10x1M(b *testing.B) {
	benchmarkEncode(b, 14, 10, 1024*1024)
}

// Benchmark 14 data shards and 10 parity shards with 16MB each.
func BenchmarkEncode14x10x16M(b *testing.B) {
	benchmarkEncode(b, 14, 10, 16*1024*1024)
}


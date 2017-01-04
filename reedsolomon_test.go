package reedsolomon

import (
	"math/rand"
	"testing"
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
	
	err = r.Encode(shards, 0)
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

//func TestAVX2(t *testing.T) {
//
//	low := []uint8{0x0, 0x2, 0x4, 0x6, 0x8, 0xa, 0xc, 0xe, 0x10, 0x12, 0x14, 0x16, 0x18, 0x1a, 0x1c, 0x1e}
//	high := []uint8{0x0, 0x20, 0x40, 0x60, 0x80, 0xa0, 0xc0, 0xe0, 0x1d, 0x3d, 0x5d, 0x7d, 0x9d, 0xbd, 0xdd, 0xfd}
//
//	in := []uint8{11, 112, 113, 114, 12, 116, 117, 118, 119, 110, 211, 212, 213, 214, 215, 20,
//		11, 112, 113, 114, 12, 116, 117, 118, 119, 110, 211, 212, 213, 214, 215, 20}
//	out := make([]byte,32)
//	gfMulAVX2(low, high, in, out)
//	fmt.Println(out)
//}

//func TestASMEncode(t *testing.T) {
//
//
//	NumData := 10
//	NumParity := 4
//	shardSize := 32
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
//	err = r.Encode(shards, 0)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println("avx",shards[NumData:])
//
//	//parityA := make([][]byte, NumParity)
//	//for i := range parityA {
//	//	parityA[i] = shards[NumData+i]
//	//}
//
//	err = r.Encode(shards, 2)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println("noasm", shards[NumData:])
//	fmt.Println("gen", r.gen)
//	fmt.Println("shards", shards)
//	//parityN := shards[NumData:]
//	//for i := 0; i < NumParity; i++ {
//	//	for j := 0; j < shardSize; j++ {
//	//		if parityA[i][j] != parityN[i][j] {
//	//			t.Fatal("asm:",parityA[i][j],"noasm:",parityN[i][j])
//	//		}
//	//	}
//	//}
//}
//
//func TestASMEncode(t *testing.T) {
//
//	NumData := 10
//	NumParity := 4
//	shardSize := 32
//	r, err := New(NumData, NumParity)
//	if err != nil {
//		t.Fatal(err)
//	}
//	shards := make([][]byte, NumData + NumParity)
//	for s := range shards {
//		shards[s] = make([]byte, shardSize)
//	}
//
//	for s := 0; s < NumData; s++ {
//		fillRegular(shards[s], s)
//	}
//	err = r.Encode(shards, 0)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println("avx", shards[NumData:])
//}
//
//func TestNOASMEncode(t *testing.T) {
//
//	NumData := 10
//	NumParity := 4
//	shardSize := 32
//	r, err := New(NumData, NumParity)
//	if err != nil {
//		t.Fatal(err)
//	}
//	shards := make([][]byte, NumData + NumParity)
//	for s := range shards {
//		shards[s] = make([]byte, shardSize)
//	}
//
//	for s := 0; s < NumData; s++ {
//		fillRegular(shards[s], s)
//	}
//	err = r.Encode(shards, 2)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println("noavx", shards[NumData:])
//}
	//
	////parityA := make([][]byte, NumParity)
	////for i := range parityA {
	////	parityA[i] = shards[NumData+i]
	////}
	//
	//err = r.Encode(shards, 2)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println("noasm", shards[NumData:])
	//fmt.Println("gen", r.gen)
	//fmt.Println("shards", shards)
	////parityN := shards[NumData:]
	//for i := 0; i < NumParity; i++ {
	//	for j := 0; j < shardSize; j++ {
	//		if parityA[i][j] != parityN[i][j] {
	//			t.Fatal("asm:",parityA[i][j],"noasm:",parityN[i][j])
	//		}
	//	}
	//}
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
	codec.Encode(shards, 2)
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

func fillRegular(p []byte, s int)  {
	size := len(p)
	for i := 0; i < size; i++ {
		v := i + s
		for v > 255 {
			v = v - 255
		}
		p[i] = byte(v)
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
		err = r.Encode(shards, 0)
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


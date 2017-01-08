# Reed-Solomon

Reed-Solomon Erasure Coding engine in Go, with speeds exceeding 2GB/cpu physics core implemented in pure Go.

 * Coding over in GF(2^8).
 * Primitive Polynomial: x^8 + x^4 + x^3 + x^2 + 1 (0x1d)

It released by  [Klauspost ReedSolomon](https://github.com/klauspost/reedsolom), with some optimizations/changes:

1. Refactoring code
2. Use Cauchy matrix as generator matrix, we can use it directly.Vandermonde matrix need some operation for preserving the 
property that any square subset of rows is invertible(and I think there is a way to optimize inverse matrix's performance, I need some time to make it)
3. There are a tool(tools/gentables.go) for generator Primitive Polynomial and it'i log table, exp table, multiply table,
inverse table etc. We can get more info about how galois field work
4. Use a "pipeline mode" for encoding concurrency. And I found L1Cache Size will be a good choice as the concurrency unit,
it improve performance greatly
5. Go1.7 have added some new instruction, and some are what we need here. The byte codes in asm files are changed to
instructions now
6. Drop inverse matrix cache, it’s a statistical fact that only 2-3% shards need to be repaired.
So I don't think it will improve performance much

# Installation
To get the package use the standard:
```bash
go get github.com/templexxx/reedsolomon
```

# Usage

This section assumes you know the basics of Reed-Solomon encoding. A good start is this [Backblaze blog post](https://www.backblaze.com/blog/reed-solomon/) or [my articles](http://templex.xyz) (more info about this package there).

There are only two public function in the package: Encode and Reconst

Encode : calculate parity of data shards;

Reconst: calculate data or parity from present shards;

# Performance
Performance depends mainly on:
1. number of parity shards
2. number of cores of CPU
3. CPU instruction extension(support AVX2 or SSSE3)
4. unit size of concurrence

Example of performance scaling on MacBook Pro "Core i7" 2.6 15" 2.6 GHz Core i7 (I7-6700HQ) - 4 physical cores, 8 logical cores. The example uses 10 data shards with 4 parity shards.

| DataSize | MB/s   | 
|---------|---------|
| 64KB    | 9724.94 | 
| 128KB   | 18165.53|
| 256KB   | 20071.62| 
| 512KB   | 16042.91| 
| 1MB     |14379.51 |
| 16MB    |12692.38 |

# Links
* [Klauspost ReedSolomon](https://github.com/klauspost/reedsolom)
* [intel ISA-L](https://github.com/01org/isa-l)

# License

This code, as the original [Klauspost ReedSolomon](https://github.com/klauspost/reedsolomon) is published under an MIT license. See LICENSE file for more information.

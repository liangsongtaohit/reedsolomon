# Reed-Solomon

Reed-Solomon Erasure Coding engine in Go, with speeds exceeding more than 3GB/s per physics core implemented in pure Go.

 * Coding over in GF(2^8).
 * Primitive Polynomial: x^8 + x^4 + x^3 + x^2 + 1 (0x1d)

It released by  [Klauspost ReedSolomon](https://github.com/klauspost/reedsolomon), with some optimizations/changes:

1. Only support AVX2. I think SSSE3 maybe out of date
2. Use Cauchy matrix as generator matrix, we can use it directly.Vandermonde matrix need some operation for preserving the 
property that any square subset of rows is invertible(and I think there is a way to optimize inverse matrix's performance, I need some time to make it)
3. There are a tool(tools/gentables.go) for generator Primitive Polynomial and it's log table, exp table, multiply table,
inverse table etc. We can get more info about how galois field work
4. Use a "pipeline mode" for encoding concurrency.
And logic cores number will be the pipeline number, anyway I don't think it's necessary to use hyper-threading
5. 32768 bytes(it's the L1 data cache size of many kinds of CPU) will be the default concurrency unit,
   it improve performance greatly(especially if the data shard's size is large)
6. Go1.7 have added some new instruction, and some are what we need here. The byte codes in asm files are changed to
instructions now
7. Delete inverse matrix cache part, it’s a statistical fact that only 2-3% shards need to be repaired.
So I don't think it will improve performance very much
8. Only 500 lines of codes(test & table not include), it's tiny
9. ...

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
2. number of cores of CPU（linear dependence, n-core cpu performance = one * n）
3. CPU instruction extension(only support AVX2)
4. unit size of concurrence
5. size of shards

Example of performance on my MacBook 2014-mid(i5-4278U 2.6GHz 2 physical cores). The example uses 10 data with 4 parity 16MB per shards.

![alt tag](http://templex.xyz/images/reedsolomon/mybench.jpg)

# Links
* [Klauspost ReedSolomon](https://github.com/klauspost/reedsolom)
* [intel ISA-L](https://github.com/01org/isa-l)

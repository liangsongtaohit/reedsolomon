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
4. Use a "pipeline mode" for encoding concurrency,
and physics cores number will be the pipeline number(it saves power, :D )
5. 32768 bytes(it's the L1 data cache size of many kinds of CPU,small unit is more cache-friendly) will be the default concurrency unit,
   it improve performance greatly(especially if the data shard's size is large).
6. Go1.7 have added some new instruction, and some are what we need here. The byte codes in asm files are changed to
instructions now (unfortunately, I added some new byte codes)
7. Delete inverse matrix cache part, it’s a statistical fact that only 2-3% shards need to be repaired.
So I don't think it will improve performance very much
8. Only 500 lines of codes(test & table not include), it's tiny
9. Instead of copying data, I use maps to save position of data. Reconstruction is almost as fast as encoding now
10. AVX intrinsic instructions are not mixed with any of SSE instructions, so we don't need "VZEROUPPER" to avoid AVX-SSE Transition Penalties,
it seems improve performance.
11. Some of Golang's asm OP codes make me uncomfortable, especially the "MOVQ", so I use byte codes to operate the register lower part sometimes.
I still need time to learn the golang asm more.
12. ...

# Installation
To get the package use the standard:
```bash
go get github.com/templexxx/reedsolomon
```

# Usage

This section assumes you know the basics of Reed-Solomon encoding. A good start is this [Backblaze blog post](https://www.backblaze.com/blog/reed-solomon/) or [my blogs](http://templex.xyz) (more info about this package there).

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
6. speed of memory(waste so much time on read/write mem, :D )

Example of performance on my MacBook 2014-mid(i5-4278U 2.6GHz 2 physical cores). The 16MB per shards.

| Encode/Reconst | data+Parity/data+Lost    | Speed (MB/S) | MacBook (i7-6700HQ)|
|----------------|--------------------------|--------------|--------------------|
|      E         | 10+4       |6629.72  | 14379.51 MB/S|
| E              | 28+4       | 7524.88  | 15542.39 MB/S|
| R              | 10+1       | 15198.09 |
| R              | 10+2       | 10993.94  |
| R              | 10+3       | 8559.67  |
| R              | 10+4      | 5283.62  |
| R              | 28+1 | 16735.21  |
| R              | 28+2 | 12581.73  |
| R              | 28+3 | 9783.60  |
| R              | 28+4 | 7788.79  |

# Links
* [Klauspost ReedSolomon](https://github.com/klauspost/reedsolom)
* [intel ISA-L](https://github.com/01org/isa-l)

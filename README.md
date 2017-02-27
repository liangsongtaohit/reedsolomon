# Reed-Solomon

Reed-Solomon Erasure Code engine in pure Go.

More info in [my blogs](http://www.templex.xyz/blog/101/reedsolomon.html) (in Chinese)

4GB/s per physics core

 * Coding over in GF(2^8).
 * Primitive Polynomial: x^8 + x^4 + x^3 + x^2 + 1 (0x1d)

It released by  [Klauspost ReedSolomon](https://github.com/klauspost/reedsolomon), with some optimizations/changes:

1. Support AVX2 and SSSE3
2. Use Cauchy matrix as generator matrix, we can use it directly.Vandermonde matrix need some operation for preserving the
property that any square subset of rows is invertible(and I think there is a way to optimize inverse matrix's performance, I need some time to make it)
3. There are a tool(tools/gentables.go) for generator Primitive Polynomial and it's log table, exp table, multiply table,
inverse table etc. We can get more info about how galois field work
4. Use a single core to encode. If you want use more, please see the history of "encode.go"(it use a "pipeline mode" for encoding concurrency,
and physics cores number will be the pipeline number)
5. 16*1024 bytes(it's a half of L1 data cache size of many kinds of CPU) will be the default calculate/concurrency unit,
   it improve performance greatly(especially when the data shard's size is large).
6. Go1.7 have added some new instruction, and some are what we need here. The byte sequence in asm files are changed to
instructions now (unfortunately, I added some new bytes)
7. Delete inverse matrix cache part, it’s a statistical fact that only 2-3% shards need to be repaired.
So I don't think it will improve performance very much
8. Only 500 lines of codes(test & table not include), it's tiny
9. Instead of copying data, I use maps to save position of data. Reconstruction is almost as fast as encoding now
10. AVX intrinsic instructions are not mixed with any of SSE instructions, so we don't need "VZEROUPPER" to avoid AVX-SSE Transition Penalties,
it seems improve performance.
11. Some of Golang's asm OP codes make me uncomfortable, especially the "MOVQ", so I use byte codes to operate the register lower part sometimes.
I still need time to learn the golang asm more. (Thanks to [asm2plan9s](https://github.com/fwessels/asm2plan9s))
12. I heared that TEST is faster than CMP, so I use TEST in my codes.But I find they have same speed
13. No R8-R15 register in asm codes, because it need one more byte
14. Only import Golang standard library
15. ...

# Installation
To get the package use the standard:
```bash
go get github.com/templexxx/reedsolomon
```

# Usage

This section assumes you know the basics of Reed-Solomon encoding. A good start is this [Backblaze blog post](https://www.backblaze.com/blog/reed-solomon/).

There are only two public function in the package: Encode, Reconst and NewMatrix

NewMatrix: return a [][]byte for encode and reconst

Encode : calculate parity of data shards;

Reconst: calculate data or parity from present shards;

# Performance
Performance depends mainly on:

1. number of parity shards
2. number of cores of CPU (if you want to use parallel version)
3. CPU instruction extension(AVX2 or SSSE3)
4. unit size of concurrence
5. size of shards
6. speed of memory(waste so much time on read/write mem, :D )

Example of performance on my MacBook 2014-mid(i5-4278U 2.6GHz 2 physical cores). The 16MB per shards.
Single core work here:

| Encode/Reconst | data+Parity/data+Lost    | Speed (MB/S) |
|----------------|-------------------|--------------|
| E              |      10+4       |4408.81  |
| E              |      17+3       | 5450.13  |
| R              |      10+1       | 10635.47 |
| R              |      10+2       | 6963.70  |
| R              |      10+3       | 5415.28  |
| R              |      10+4      | 3469.50 |
| R              |      17+1 | 10772.59  |
| R              |      17+2 | 7159.81  |
| R              |      17+3 | 5335.74  |

# Links
* [Klauspost ReedSolomon](https://github.com/klauspost/reedsolom)
* [intel ISA-L](https://github.com/01org/isa-l)

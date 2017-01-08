// Copyright 2016, TempleX, see LICENSE for details.

// Reference: www.ssrc.ucsc.edu/Papers/plank-fast13.pdf

#include "textflag.h"

// func gfMulAVX2(low, high, in, out []byte)
TEXT ·gfMulAVX2(SB), NOSPLIT, $0

	// table -> ymm
	MOVQ	lowTable+0(FP), AX
	MOVQ	highTable+24(FP), BX
	MOVOU (AX), X0
	MOVOU (BX), X1
	// [...,X0] -> [X0, X0]
    VINSERTI128 $1, X0, Y0, Y0      // low_table -> ymm0
    VINSERTI128 $1, X1, Y1, Y1      // high_table -> ymm1

    MOVQ  in+48(FP), AX     // in_add -> AX
    MOVQ  out+72(FP), BX    // out_add -> BX

    // mask -> ymm
    MOVQ  $15, CX
    MOVQ  CX, X2
    VPBROADCASTB X2, Y3

    // if done
	MOVQ  in_len+56(FP), CX // in_len -> CX
	SHRQ $5, CX	// CX = CX >> 5 (calc 32bytes per loop)
	TESTQ CX, CX	// bitwise AND on two operands,if result is 0 (it means no more data)，ZF flag set 1
	JZ    done	// jump to done if ZF is 0

loop:
	// split data byte into two 4-bit
	VMOVDQU (AX), Y4	// in_data -> ymm4
	VPSRLQ $4, Y4, Y5	// shift in_data's 4 high bit to low -> ymm5
    VPAND Y3, Y5, Y5	// mask AND data_shift -> ymm5 (high data)
	VPAND Y3, Y4, Y4	// mask AND data -> ymm4 (low data)

	// shuffle table
	VPSHUFB Y5, Y1, Y6
	VPSHUFB Y4, Y0, Y7

	// combine low, high 4-bit & output
	VPXOR Y6, Y7, Y8
	VMOVDQU Y8, (BX)

	// prepare next loop
	ADDQ $32, AX
	ADDQ $32, BX
	SUBQ $1, CX	// it will affect ZF
	JNZ  loop

done:
	VZEROUPPER	// avoiding-avx-sse-transition-penalties
	RET

// almost same with gfMulAVX2
// two more steps: 1. get the old out_data 2. update the out_data
// func gfMulXorAVX2(low, high, in, out []byte)
TEXT ·gfMulXorAVX2(SB), NOSPLIT, $0

	MOVQ  lowTable+0(FP), AX
	MOVQ  highTable+24(FP), BX
	MOVOU (AX), X0
	MOVOU (BX), X1
    VINSERTI128 $1, X0, Y0, Y0
    VINSERTI128 $1, X1, Y1, Y1

    MOVQ  in+48(FP), AX
    MOVQ  out+72(FP), BX

    MOVQ  $15, CX
    MOVQ  CX, X2
    VPBROADCASTB X2, Y3

	MOVQ  in_len+56(FP), CX
	SHRQ $5, CX
	TESTQ CX, CX
	JZ    done

loop:
	VMOVDQU (AX), Y4
	VMOVDQU (BX), Y8	// out_data -> Ymm
	VPSRLQ $4, Y4, Y5
    VPAND Y3, Y5, Y5
	VPAND Y3, Y4, Y4

	VPSHUFB Y5, Y1, Y6
	VPSHUFB Y4, Y0, Y7

	VPXOR Y6, Y7, Y6
	VPXOR Y6, Y8, Y8	// update result
	VMOVDQU Y8, (BX)

	ADDQ $32, AX
	ADDQ $32, BX
	SUBQ $1, CX
	JNZ  loop

done:
	VZEROUPPER
	RET

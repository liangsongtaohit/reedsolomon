// Reference: www.ssrc.ucsc.edu/Papers/plank-fast13.pdf

#include "textflag.h"

// func gfMulAVX2(low, high, in, out []byte)
TEXT ·gfMulAVX2(SB), NOSPLIT, $0
	// table -> ymm
	MOVQ    lowTable+0(FP), AX   // it's not intel OP code MOVQ, it's more like MOV
	MOVQ    highTable+24(FP), BX
	VMOVDQU (AX), X0             // 128-bit Intel® AVX instructions operate on the lower 128 bits of the YMM registers and zero the upper 128 bits
	VMOVDQU (BX), X1

	// [0..0,X0] -> [X0, X0]
	VINSERTI128 $1, X0, Y0, Y0 // low_table -> ymm0
	VINSERTI128 $1, X1, Y1, Y1 // high_table -> ymm1
	MOVQ        in+48(FP), AX  // in_add -> AX
	MOVQ        out+72(FP), BX // out_add -> BX

	// mask -> ymm
	BYTE         $0xb1; BYTE $0x0f                                                 // MOV $0x0f, CL
	BYTE         $0xc4; BYTE $0xe3; BYTE $0x69; BYTE $0x20; BYTE $0xd1; BYTE $0x00 // VPINSRB $0x00, ECX, XMM2, XMM2
	VPBROADCASTB X2, Y2                                                            // [1111,1111,1111...1111]

	// if done
	MOVQ  in_len+56(FP), CX // in_len -> CX
	SHRQ  $5, CX            // CX = CX >> 5 (calc 32bytes per loop)
	TESTQ CX, CX            // bitwise AND on two operands,if result is 0 (it means no more data)，ZF flag set 1
	JZ    done              // jump to done if ZF is 0

loop:
	// split data byte into two 4-bit
	VMOVDQU (AX), Y4   // in_data -> ymm4
	VPSRLQ  $4, Y4, Y5 // shift in_data's 4 high bit to low -> ymm5
	VPAND   Y2, Y5, Y5 // mask AND data_shift -> ymm5 (high data)
	VPAND   Y2, Y4, Y4 // mask AND data -> ymm4 (low data)

	// shuffle table
	VPSHUFB Y5, Y1, Y6
	VPSHUFB Y4, Y0, Y7

	// gf add low, high 4-bit & output
	VPXOR   Y6, Y7, Y3
	VMOVDQU Y3, (BX)   // can't use Non-Temporal Hint here, because "out" will be read many times

	// next loop
	ADDQ $32, AX
	ADDQ $32, BX
	SUBQ $1, CX  // it will affect ZF
	JNZ  loop

done:
	RET

// almost same with gfMulAVX2
// two more steps: 1. get the old out_data 2. update the out_data
// func gfMulXorAVX2(low, high, in, out []byte)
TEXT ·gfMulXorAVX2(SB), NOSPLIT, $0
	MOVQ         lowTable+0(FP), AX
	MOVQ         highTable+24(FP), BX
	VMOVDQU      (AX), X0
	VMOVDQU      (BX), X1
	VINSERTI128  $1, X0, Y0, Y0
	VINSERTI128  $1, X1, Y1, Y1
	MOVQ         in+48(FP), AX
	MOVQ         out+72(FP), BX
	BYTE         $0xb1; BYTE $0x0f
	BYTE         $0xc4; BYTE $0xe3; BYTE $0x69; BYTE $0x20; BYTE $0xd1; BYTE $0x00
	VPBROADCASTB X2, Y2
	MOVQ         in_len+56(FP), CX
	SHRQ         $5, CX
	TESTQ        CX, CX
	JZ           done

loop:
	VMOVDQU (AX), Y4
	VMOVDQU (BX), Y3   // out_data -> Ymm
	VPSRLQ  $4, Y4, Y5
	VPAND   Y2, Y5, Y5
	VPAND   Y2, Y4, Y4
	VPSHUFB Y5, Y1, Y6
	VPSHUFB Y4, Y0, Y7
	VPXOR   Y6, Y7, Y6
	VPXOR   Y6, Y3, Y3 // update result
	VMOVDQU Y3, (BX)
	ADDQ    $32, AX
	ADDQ    $32, BX
	SUBQ    $1, CX
	JNZ     loop

done:
	RET

// func avx2() bool
TEXT ·avx2(SB), NOSPLIT, $0
	CMPB runtime·support_avx2(SB), $1
	JE   has
	MOVB $0, ret+0(FP)
	RET

has:
	MOVB $1, ret+0(FP)
	RET


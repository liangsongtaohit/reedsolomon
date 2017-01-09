// Copyright 2016, TempleX, see LICENSE for details.

// Reference: www.ssrc.ucsc.edu/Papers/plank-fast13.pdf

// almost same with gfavx2_amd64.s
// SSE instructions replace AVX2 instructions
// SSE instructions's operands can't be as many as AVX2,
// so you need copy some register sometimes for keep old register clean

#include "textflag.h"

// func gfMulSSSE3(low, high, in, out []byte)
TEXT ·gfMulSSSE3(SB), NOSPLIT, $0

	// table -> xmm
	MOVQ	low+0(FP), AX
	MOVQ	high+24(FP), BX
	MOVOU	(AX), X0
	MOVOU	(BX), X1

	// in out data_add -> reg
	MOVQ	in+48(FP), AX
	MOVQ	out+72(FP), BX

	// prepare the mask
	MOVQ	$15, CX
	MOVQ	CX, X3
	PXOR	X2, X2	// clean up xmm
	PSHUFB	X2, X3  // shuffle the mask according index[0,...,0]

	// ready for loop
	MOVQ	in_len+56(FP), CX
	SHRQ	$4, CX
	TESTQ	CX, CX
	JZ	done

loop:
	// in_add -> AX; out_add -> BX;
    // in_len -> CX;
    // lowTable -> X0; highTable -> X1;
    // mask -> X3

    // split data byte into two 4-bit
	MOVOU	(AX), X4 // in_data -> X4
	MOVOU   X4, X5  // in_data_copy -> X5
	PAND    X3, X4  // in_data_low -> X4
	PSRLQ   $4, X5
	PAND    X3, X5  // in_data_high -> X5

	// shuffle table
	MOVOU   X0, X6   // lowTable_copy -> X6
    MOVOU   X1, X7   // highTable_copy -> X7
	PSHUFB  X4, X6   // lowResult -> X6
	PSHUFB  X5, X7   // highResult -> X7

	// combine low, high 4-bit & output
	PXOR    X6, X7
	MOVOU   X7, (BX)

	// prepare next loop
	ADDQ    $16, AX  // in+=16
	ADDQ    $16, BX  // out+=16
	SUBQ    $1, CX
	JNZ     loop

done:
	RET

// func gfMulXorSSSE3(low, high, in, out []byte)
TEXT ·gfMulXorSSSE3(SB), NOSPLIT, $0

	// table -> xmm
	MOVQ	low+0(FP), AX
	MOVQ	high+24(FP), BX
	MOVOU	(AX), X0
	MOVOU	(BX), X1

	// in out data_add -> reg
	MOVQ	in+48(FP), AX
	MOVQ	out+72(FP), BX

	// prepare the mask
	MOVQ	$15, CX
	MOVQ	CX, X3
	PXOR	X2, X2	// clean up xmm
	PSHUFB	X2, X3  // shuffle the mask according index[0,...,0]

	// ready for loop
	MOVQ	in_len+56(FP), CX
	SHRQ	$4, CX
	TESTQ	CX, CX
	JZ	done

loop:
	// in_add -> AX; out_add -> BX;
    // in_len -> CX;
    // lowTable -> X0; highTable -> X1;
    // mask -> X3

    // split data byte into two 4-bit
	MOVOU	(AX), X4 // in_data -> X4
	MOVOU   (BX), X8 // out_data -> X8
	MOVOU   X4, X5  // in_data_copy -> X5
	PAND    X3, X4  // in_data_low -> X4
	PSRLQ   $4, X5
	PAND    X3, X5  // in_data_high -> X5

	// shuffle table
	MOVOU   X0, X6   // lowTable_copy -> X6
    MOVOU   X1, X7   // highTable_copy -> X7
	PSHUFB  X4, X6   // lowResult -> X6
	PSHUFB  X5, X7   // highResult -> X7

	// combine low, high 4-bit & output
	PXOR    X6, X7
	PXOR    X8, X7  // result_update
	MOVOU   X7, (BX)

	// prepare next loop
	ADDQ    $16, AX  // in+=16
	ADDQ    $16, BX  // out+=16
	SUBQ    $1, CX
	JNZ     loop

done:
	RET

// func ssse3() bool
TEXT ·ssse3(SB), NOSPLIT, $0
	XORQ AX, AX
	INCL AX
	CPUID   // when CPUID excutes with AX set to 01H, feature info is ret in CX and DX
	SHRQ $9, CX     // SSSE3 -> CX[9] = 1
	ANDQ $1, CX
	MOVB CX, ret+0(FP)
	RET

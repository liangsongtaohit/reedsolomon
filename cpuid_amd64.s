#include "textflag.h"

// func CheckAVX2() bool
TEXT ·CheckAVX2(SB),NOSPLIT,$0
	CMPB    runtime·support_avx2(SB), $1
	JE      has
        MOVB    $0, ret+0(FP)
	RET
has:
        MOVB    $1, ret+0(FP)
	RET

// func CheckSSSE3() bool
TEXT ·CheckSSSE3(SB), NOSPLIT, $0
	XORQ AX, AX
	INCL AX
	CPUID   // when CPUID excutes with AX set to 01H, feature info is ret in CX and DX
	SHRQ $9, CX     // SSSE3 -> CX[9] = 1
	ANDQ $1, CX
	MOVB CX, ret+0(FP)
	RET

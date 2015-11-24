// +build amd64

#include "textflag.h"

// func divmod(uint64) (q, r uint64)
TEXT ·divmod(SB),NOSPLIT,$0
	MOVQ	a+0(FP), DI
	MOVQ	DI, AX

	// can't get MOVABSQ $5088756985850910791, CX to work
	MOVQ $0x469ee58469ee5847, CX
	
	SHRQ	$1, AX
	MULQ	CX
	MOVQ	DX, AX
	SHRQ	$3, AX

	// some bullshit i cant use IMULQ $58, AX, DX
	BYTE $0x6b
	BYTE $0xd0
	BYTE $0x3a

	SUBQ	DX, DI
	MOVQ	DI, DX
	MOVQ	AX, q+8(FP)
	MOVQ	DX, r+16(FP)
	RET

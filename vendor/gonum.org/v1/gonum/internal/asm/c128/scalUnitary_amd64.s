// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

#define SRC SI
#define DST SI
#define LEN CX
#define IDX AX
#define TAIL BX
#define ALPHA X0
#define ALPHA_C X1
#define ALPHA2 X10
#define ALPHA_C2 X11

#define MOVDDUP_X2_X3    LONG $0xDA120FF2 // MOVDDUP X2, X3
#define MOVDDUP_X4_X5    LONG $0xEC120FF2 // MOVDDUP X4, X5
#define MOVDDUP_X6_X7    LONG $0xFE120FF2 // MOVDDUP X6, X7
#define MOVDDUP_X8_X9    LONG $0x120F45F2; BYTE $0xC8 // MOVDDUP X8, X9

#define ADDSUBPD_X2_X3    LONG $0xDAD00F66 // ADDSUBPD X2, X3
#define ADDSUBPD_X4_X5    LONG $0xECD00F66 // ADDSUBPD X4, X5
#define ADDSUBPD_X6_X7    LONG $0xFED00F66 // ADDSUBPD X6, X7
#define ADDSUBPD_X8_X9    LONG $0xD00F4566; BYTE $0xC8 // ADDSUBPD X8, X9

// func ScalUnitary(alpha complex128, x []complex128)
TEXT ·ScalUnitary(SB), NOSPLIT, $0
	MOVQ x_base+16(FP), SRC // SRC = &x
	MOVQ x_len+24(FP), LEN  // LEN = len(x)
	CMPQ LEN, $0            // if LEN == 0 { return }
	JE   scal_end

	MOVUPS alpha+0(FP), ALPHA     // ALPHA = { imag(alpha), real(alpha) }
	MOVAPS ALPHA, ALPHA_C
	SHUFPD $0x1, ALPHA_C, ALPHA_C // ALPHA_C = { real(alpha), imag(alpha) }

	XORQ   IDX, IDX          // IDX = 0
	MOVAPS ALPHA, ALPHA2     // Copy ALPHA and ALPHA_C for pipelining
	MOVAPS ALPHA_C, ALPHA_C2
	MOVQ   LEN, TAIL
	SHRQ   $2, LEN           // LEN = floor( n / 4 )
	JZ     scal_tail         // if BX == 0 { goto scal_tail }

scal_loop: // do {
	MOVUPS (SRC)(IDX*8), X2   // X_i = { imag(x[i]), real(x[i]) }
	MOVUPS 16(SRC)(IDX*8), X4
	MOVUPS 32(SRC)(IDX*8), X6
	MOVUPS 48(SRC)(IDX*8), X8

	// X_(i+1) = { real(x[i], real(x[i]) }
	MOVDDUP_X2_X3
	MOVDDUP_X4_X5
	MOVDDUP_X6_X7
	MOVDDUP_X8_X9

	// X_i = { imag(x[i]), imag(x[i]) }
	SHUFPD $0x3, X2, X2
	SHUFPD $0x3, X4, X4
	SHUFPD $0x3, X6, X6
	SHUFPD $0x3, X8, X8

	// X_i     = { real(ALPHA) * imag(x[i]), imag(ALPHA) * imag(x[i])  }
	// X_(i+1) = { imag(ALPHA) * real(x[i]), real(ALPHA) * real(x[i])  }
	MULPD ALPHA_C, X2
	MULPD ALPHA, X3
	MULPD ALPHA_C2, X4
	MULPD ALPHA2, X5
	MULPD ALPHA_C, X6
	MULPD ALPHA, X7
	MULPD ALPHA_C2, X8
	MULPD ALPHA2, X9

	// X_(i+1) = {
	//	imag(result[i]):  imag(ALPHA)*real(x[i]) + real(ALPHA)*imag(x[i]),
	//	real(result[i]):  real(ALPHA)*real(x[i]) - imag(ALPHA)*imag(x[i])
	//  }
	ADDSUBPD_X2_X3
	ADDSUBPD_X4_X5
	ADDSUBPD_X6_X7
	ADDSUBPD_X8_X9

	MOVUPS X3, (DST)(IDX*8)   // x[i] = X_(i+1)
	MOVUPS X5, 16(DST)(IDX*8)
	MOVUPS X7, 32(DST)(IDX*8)
	MOVUPS X9, 48(DST)(IDX*8)
	ADDQ   $8, IDX            // IDX += 8
	DECQ   LEN
	JNZ    scal_loop          // } while --LEN > 0

scal_tail:
	ANDQ $3, TAIL // TAIL = TAIL % 4
	JZ   scal_end // if TAIL == 0 { return }

scal_tail_loop: // do {
	MOVUPS (SRC)(IDX*8), X2 // X_i = { imag(x[i]), real(x[i]) }
	MOVDDUP_X2_X3           // X_(i+1) = { real(x[i], real(x[i]) }
	SHUFPD $0x3, X2, X2     // X_i = { imag(x[i]), imag(x[i]) }
	MULPD  ALPHA_C, X2      // X_i     = { real(ALPHA) * imag(x[i]), imag(ALPHA) * imag(x[i])  }
	MULPD  ALPHA, X3        // X_(i+1) = { imag(ALPHA) * real(x[i]), real(ALPHA) * real(x[i])  }

	// X_(i+1) = {
	//	imag(result[i]):  imag(ALPHA)*real(x[i]) + real(ALPHA)*imag(x[i]),
	//	real(result[i]):  real(ALPHA)*real(x[i]) - imag(ALPHA)*imag(x[i])
	//  }
	ADDSUBPD_X2_X3

	MOVUPS X3, (DST)(IDX*8) // x[i] = X_(i+1)
	ADDQ   $2, IDX          // IDX += 2
	DECQ   TAIL
	JNZ    scal_tail_loop   // }  while --LEN > 0

scal_end:
	RET

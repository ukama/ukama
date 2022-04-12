/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2013, Atmel Corporation
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * - Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the disclaimer below.
 *
 * Atmel's name may not be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * DISCLAIMER: THIS SOFTWARE IS PROVIDED BY ATMEL "AS IS" AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT ARE
 * DISCLAIMED. IN NO EVENT SHALL ATMEL BE LIABLE FOR ANY DIRECT, INDIRECT,
 * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA,
 * OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
 * EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
#include <mon_macros.h>

.section monitor
	.text
	.align

	.global monitor_init
	.global init_loadkernel_args
	.extern monitor_vectors
	.extern svc_mgr_veneer

#define NSACR_CLR_BITS  (NSACR_BIT_NSD32DIS | NSACR_BIT_NSASEDIS)
#define NSACR_SET_BITS  (NSACR_BIT_CP | NSACR_BIT_SMP)

monitor_init:
	/* Configure Non-Secure Access control register */
	mrc	p15, 0, r0, c1, c1, 2	/* Read NSACR Register */
	/* Allow NS access VFP/NEON and Co-processor registers 10-11 */
	ldr	r1, =NSACR_CLR_BITS
	bic	r0, r0, r1
	ldr	r1, =NSACR_SET_BITS
	orr	r0, r0, r1
	mcr	p15, 0, r0, c1, c1, 2

	/* Switch to MON mode */
	mrs	r0, cpsr
	bic	r0, r0, #0x1f
	orr	r0, r0, #Mode_MON
	msr	cpsr_c, r0

	/* set the monitor stack */
	ldr	sp, =MON_STACK_BASE

	/* set the monitor vector base */
	ldr	r2, =monitor_vectors
	mcr	p15, 0, r2, c12, c0, 1

	/* Program the DACR to allow client access to *all* domains */
	ldr	r2, =0x55555555
	mcr	p15, 0, r2, c3, c0, 0

	/*
	 * Configure SCR so that A bit in the CPSR
	 * can be modified when in the Non-secure state
	 */
	mrc	p15, 0, r0, c1, c1, 0
	ldr	r1, =(SCR_AW_BIT)
	orr	r0, r0, r1
	mcr	p15, 0, r0, c1, c1, 0

	/* Set NWd initial CPSR and PC (NWd boot entry address) */
	ldr	r0, =MON_DATA_BASE

	ldr	r1, =INITIAL_NWD_CPSR
	str	r1, [r0, #NWD_CPSR_OFF]
	ldr	r1, =NWD_BOOT_ADDR
	str	r1, [r0, #NWD_PC_OFF]

	/* DEBUG */
#if 0
	/* fill NWd boot address with while(true) loops */
	ldr	r2, =0xeafffffe
	str	r2, [r1]
	str	r2, [r1, #4]
#endif

	/* Set SWd initial CPSR and PC (SWd handler) */
	ldr	r1, =INITIAL_SWD_CPSR
	str	r1, [r0, #SWD_CPSR_OFF]
	ldr	r1, =svc_mgr_veneer
	str	r1, [r0, #SWD_PC_OFF]

	/* Switch back to SVC mode */
	mrs	r0, cpsr
	bic	r0, r0, #0x1f
	orr	r0, r0, #Mode_SVC
	msr	cpsr_c, r0

	/* Return to SWd calling function */
	bx	lr

/*
 * Save r0, r1, r2 to DB for loading kernel.
 * void init_loadkernel_args(unsigned int r0, unsigned int r1,
 *			unsigned int r2, unsigned int r3)
 *
 * r0 : 0
 * r1 : Valid machine type number
 * r2 : physical pointer to the device-tree block
 * r3 : overwrite the PC_OFF
 */
init_loadkernel_args:
	ldr	r5, =MON_DATA_BASE
	str	r0, [r5, #NWD_R0_OFF]
	str	r1, [r5, #NWD_R1_OFF]
	str	r2, [r5, #NWD_R2_OFF]
	str	r3, [r5, #NWD_PC_OFF]
	mov	r4, #DB_R012_VALID
	str	r4, [r5, #NWD_R012_VALID_OFF]

	bx	lr

	.end

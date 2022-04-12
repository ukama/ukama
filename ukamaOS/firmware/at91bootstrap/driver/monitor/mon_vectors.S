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

.section monitor
	 .text
	 .align 8

	.extern monitor_smc_handler
	.global monitor_vectors

/*
 *-----------------------------------------------------------------------------
 * Monitor exception vector table
 *      Most of these exceptions does not do anything other than branching to
 *      itself except smc vector. The SMC vector jumps to smc handler which
 *      takes care of proper switching from NWd to SWd and from SWd to NWd.
 * -----------------------------------------------------------------------------
 */
monitor_vectors:
	b	.			/* Reset: not used by Monitor */
	b	.			/* Undef: not used by Monitor */
	ldr	pc, =monitor_smc_handler
	b	.			/* Prefetch: can be used by Monitor */
	b	.			/* Data abort: can be used by Monitor */
	b	.			/* RESERVED */
	b	.			/* IRQ: can be used by Monitor */
	b	.			/* FIQ: can be used by Monitor */

	.end

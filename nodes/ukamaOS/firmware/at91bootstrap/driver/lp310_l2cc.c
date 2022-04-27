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
#include "hardware.h"
#include "arch/at91_sfr.h"
#include "arch/lp310_l2cc.h"

static inline void write_l2cc(unsigned int offset, const unsigned int value)
{
	writel(value, offset + AT91C_BASE_L2CC);
}

static inline unsigned int read_l2cc(unsigned int offset)
{
	return readl(offset + AT91C_BASE_L2CC);
}

#ifdef CONFIG_SAMA5D2
static void l2cache_configure_ram(void)
{
	writel(0x1, SFR_L2CC_HRAMC + AT91C_BASE_SFR);
}
#else
static void l2cache_configure_ram(void) {}
#endif

void l2cache_prepare(void)
{
	unsigned int cfg;

	l2cache_configure_ram();

	/* disable cache if it hasn't been done yet */
	write_l2cc(L2CC_CR, 0x00);

	/* Auxiliary control register */
	/*
	 * Default values of Way size and Associativity are
	 * implementation defined (given by input pins).
	 * As a result, there should be no need to write them
	 * during this init sequence.
	 */
/*
	cfg = L2CC->L2CC_ACR;
	L2CC->L2CC_ACR = cfg;
*/
	/* L2CC global configuration */
	/* TAG, Data Latency Control: keep default */
/*
	cfg = L2CC->L2CC_TRCR;
	L2CC->L2CC_TRCR = cfg;
	cfg = L2CC->L2CC_DRCR;
	L2CC->L2CC_DRCR = cfg;
*/

	/* Prefetch Control */
	cfg = read_l2cc(L2CC_PCR);
	/* prefetch offset: TODO find proper values */
	cfg |= 0x1;
	cfg |= L2CC_PCR_IDLEN | L2CC_PCR_PDEN | L2CC_PCR_DLEN;
	cfg |= L2CC_PCR_DATPEN | L2CC_PCR_INSPEN;
	write_l2cc(L2CC_PCR, cfg);

	/* Power Control */
	cfg = read_l2cc(L2CC_POWCR);
	cfg |= L2CC_POWCR_STBYEN | L2CC_POWCR_DCKGATEN;
	write_l2cc(L2CC_POWCR, cfg);

	/* invalidate all entries */
	cfg = 0x0000ffff;
	write_l2cc(L2CC_IWR, cfg);
	/* check invalidate operation finished */
	while (read_l2cc(L2CC_IWR) != 0)
		;

	/* disable interrupts */
	write_l2cc(L2CC_IMR, 0);
	write_l2cc(L2CC_ICR, 0x01ff);

/*
	1.  Write to the Auxiliary, Tag RAM Latency, Data RAM Latency,
		Prefetch, and Power Control registers using a read-modify-write
			to set up global configurations:
	associativity, Way Size
	latencies for RAM accesses
	allocation policy
	prefetch and power capabilities.

	2. Secure write to the Invalidate by Way, offset 0x77C, to invalidate all entries in cache:
	   Write 0xFFFF to 0x77C
	   Poll cache maintenance register until invalidate operation is complete.
	3. Write to the Lockdown D and Lockdown I Register 9 if required.
	4. Write to interrupt clear register to clear any residual raw interrupts set.
	5. Write to the Interrupt Mask Register if you want to enable interrupts.
	6. Write to Control Register 1 with the LSB set to 1 to enable the cache.
*/
}

void l2cache_enable(void)
{
	/* enable cache, now! */
	write_l2cc(L2CC_CR, 1);
}

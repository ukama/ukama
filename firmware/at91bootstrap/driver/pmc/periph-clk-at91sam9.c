/* ----------------------------------------------------------------------------
 *         Microchip Microprocessor (MPU) Software Team
 * ----------------------------------------------------------------------------
 * Copyright (C) 2019 Microchip Technology Inc. and its subsidiaries
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * - Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the disclaimer below.
 *
 * Microchip's name may not be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * DISCLAIMER: THIS SOFTWARE IS PROVIDED BY MICROCHIP "AS IS" AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT ARE
 * DISCLAIMED. IN NO EVENT SHALL MICROCHIP BE LIABLE FOR ANY DIRECT, INDIRECT,
 * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA,
 * OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
 * EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
#include "board.h"
#include "clk-common.h"
#include "common.h"
#include "pmc.h"

#define PMC_PCER	0x10	/* Peripheral Clock Enable Register  (0:31 PERI_ID) */
#define PMC_PCDR	0x14	/* Peripheral Clock Disable Register (0:31 PERI_ID) */
#define PMC_PCSR	0x18	/* Peripheral Clock Status Register  (0:31 PERI_ID) */
#define PMC_PCER1	0x100	/* Peripheral Clock Enable Register1  32:63 PERI_ID */
#define PMC_PCDR1	0x104	/* Peripheral Clock Disable Register1 32:63 PERI_ID */
#define PMC_PCSR1	0x108	/* Peripheral Clock Status Register1  32:63 PERI_ID */

int pmc_enable_periph_clock(unsigned int periph_id, int divider)
{
	unsigned int mask = 0x01 << (periph_id % 32);

	if ((periph_id / 32) == 1)
		write_pmc(PMC_PCER1, mask);
	else if ((periph_id / 32) == 0)
		write_pmc(PMC_PCER, mask);
	else
		return -1;

	return 0;
}

int pmc_disable_periph_clock(unsigned int periph_id)
{
	unsigned int mask = 0x01 << (periph_id % 32);

	if ((periph_id / 32) == 1)
		write_pmc(PMC_PCDR1, mask);
	else if ((periph_id / 32) == 0)
		write_pmc(PMC_PCDR, mask);
	else
		return -1;

	return 0;
}

int pmc_periph_clock_enabled(unsigned int periph_id)
{
	unsigned int mask = 0x01 << (periph_id % 32);
	int enabled = 0;

	if ((periph_id / 32) == 1)
		enabled = !!(read_pmc(PMC_PCSR1) & mask);
	else if ((periph_id / 32) == 0)
		enabled = !!(read_pmc(PMC_PCSR) & mask);

	return enabled;
}

unsigned int pmc_periph_clock_get_rate(unsigned int periph_id)
{
	return at91_get_ahb_clock();
}

/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2012, Atmel Corporation
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
#include "board.h"
#include "debug.h"
#include "pmc.h"

#include "arch/at91_pit.h"
#include "arch/at91_pmc/pmc.h"

#define MAX_PIV		0xfffff

static inline int pit_readl(unsigned int reg)
{
	return(readl(AT91C_BASE_PITC + reg));
}

static inline void pit_writel(unsigned int value, unsigned reg)
{
	writel(value, (AT91C_BASE_PITC + reg));
}

/*
 * If MASTER_CLOCK = 132M, the maximum delay is about 520.6 second(8.6767 min).
 * it is long enough for using in bootstrap.
 */
int timer_init(void)
{
	pit_writel((MAX_PIV | AT91C_PIT_PITEN), PIT_MR);

	/* Enable PITC Clock */
#ifdef AT91C_ID_PIT
	pmc_enable_periph_clock(AT91C_ID_PIT, PMC_PERIPH_CLK_DIVIDER_NA);
#else
	pmc_enable_periph_clock(AT91C_ID_SYS, PMC_PERIPH_CLK_DIVIDER_NA);
#endif
	return 0;
}

static unsigned int at91_get_pit_value(void)
{
	return(pit_readl(PIT_PIIR));
}

/* Because the below statement is used in the function:
 *	((MASTER_CLOCK >> 10) * usec) is used,
 * to our 32-bit system. the argu "usec" maximum value is:
 * supposed "MASTER_CLOCK" is 132M.
 *	132000000 / 1024 = 128906
 *	(0xffffffff) / 128906 = 33318.
 * So the maximum delay time is 33318 us.
 */
void udelay(unsigned int usec)
{
	unsigned int base = at91_get_pit_value();
	unsigned int delay;
	unsigned int current;

	/* Since our division function which costs much run time
	 * causes the delay time error.
	 * So here using shifting to implement the division.
	 * to change "1000" to "1024", this cause some inaccuacy,
	 * but it is acceptable.
	 * ((MASTER_CLOCK / 1024) * usec) / (16 * 1024)
	 */
	if (pmc_mck_check_h32mxdiv())
		delay = (((MASTER_CLOCK / 2) >> 10) * usec) >> 14;
	else
		delay = ((MASTER_CLOCK >> 10) * usec) >> 14;

	do {
		current = at91_get_pit_value();
		current -= base;
	} while (current < delay);
}

void mdelay(unsigned int msec)
{
	unsigned int base = at91_get_pit_value();
	unsigned int delay;
	unsigned int current;

	if (pmc_mck_check_h32mxdiv())
		delay = (((MASTER_CLOCK / 2) / 1000) * msec) / 16;
	else
		delay = ((MASTER_CLOCK / 1000) * msec) / 16;

	do {
		current = at91_get_pit_value();
		current -= base;
	} while (current < delay);
}

/* Init a special timer for slow clock switch function */
static int timer1_base;

int start_interval_timer(void)
{
	timer1_base = at91_get_pit_value();

	return 0;
}

/* Because the below statement is used in the function:
 *	((MASTER_CLOCK / 1000) * msec) is used,
 * to our 32-bit system. the argu "msec" maximum value is:
 * supposed "MASTER_CLOCK" is 132M.
 *	(0xffffffff) / 132000 = 32537.
 * So the maximum delay time is 32537 ms.
 */
int wait_interval_timer(unsigned int msec)
{
	unsigned int delay;
	unsigned int current;

	if (pmc_mck_check_h32mxdiv())
		delay = (((MASTER_CLOCK / 2) / 1000) * msec) / 16;
	else
		delay = ((MASTER_CLOCK / 1000) * msec) / 16;

	do {
		current = at91_get_pit_value();
		current -= timer1_base;
	} while (current < delay);

	return 0;
}

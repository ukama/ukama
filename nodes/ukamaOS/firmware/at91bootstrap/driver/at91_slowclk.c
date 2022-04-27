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
#include "arch/at91_slowclk.h"
#include "timer.h"
#include "backup.h"

int slowclk_enable_osc32(void)
{
#if !defined(CONFIG_SAMA5D4) && !defined(CONFIG_SAMA5D2)
	unsigned int reg;

	/*
	 * Enable the 32768 Hz oscillator by setting the bit OSC32EN to 1
	 */
	reg = readl(AT91C_BASE_SCKCR);
	reg |= AT91C_SLCKSEL_OSC32EN;
	writel(reg, AT91C_BASE_SCKCR);
#endif

	/* start a internal timer */
	start_interval_timer();

	return 0;
}

static void slowclk_wait_osc32_stable(void)
{
	/* VDDBU keeps feeding oscilator. No need for wait here. */
	if (backup_resume())
		return;

	/*
	 * Wait 32768 Hz Startup Time for clock stabilization (software loop)
	 * wait about 1s (1300ms)
	 */
	wait_interval_timer(1300);
}

#if !defined(CONFIG_SAMA5D4) && !defined(CONFIG_SAMA5D2) \
	&& !defined(CONFIG_SAM9X60) && !defined(CONFIG_SAMA7G5)
static void slowclk_disable_rc32(void)
{
	unsigned int reg;

	/*
	 * Disable the 32kHz RC oscillator by setting the bit RCEN to 0
	 */
	reg = readl(AT91C_BASE_SCKCR);
	reg &= ~AT91C_SLCKSEL_RCEN;
	writel(reg, AT91C_BASE_SCKCR);
}
#endif /* #if !defined(CONFIG_SAMA5D4) && !defined(CONFIG_SAMA5D2) \
	&& !defined(CONFIG_SAM9X60) && !defined(CONFIG_SAMA7G5)
	*/

static int slowclk_select_osc32(void)
{
	unsigned int reg;

	/*
	 * Switching from internal 32kHz RC oscillator to 32768 Hz oscillator
	 * by setting the bit OSCSEL to 1
	 */
	reg = readl(AT91C_BASE_SCKCR);

	/* Do nothing if aleady set. */
	if (reg & AT91C_SLCKSEL_OSCSEL)
		return 0;

	reg |= AT91C_SLCKSEL_OSCSEL;
	writel(reg, AT91C_BASE_SCKCR);

	/*
	 * Waiting 5 slow clock cycles for internal resynchronization
	 * 5 slow clock cycles = ~153 us (5 / 32768)
	 */
	udelay(153);

	return 0;
}

int slowclk_switch_osc32(void)
{
	slowclk_wait_osc32_stable();

	slowclk_select_osc32();

#if !defined(CONFIG_SAMA5D4) && !defined(CONFIG_SAMA5D2) \
	&& !defined(CONFIG_SAM9X60) && !defined(CONFIG_SAMA7G5)
	slowclk_disable_rc32();
#endif

	return 0;
}

/* Switch from 32768 Hz Crystal Oscillator to Internal 32 kHz RC Oscillator */
int slowclk_switch_rc32(void)
{
	unsigned int reg;

#if !defined(CONFIG_SAMA5D4) && !defined(CONFIG_SAMA5D2) \
	&& !defined(CONFIG_SAM9X60) && !defined(CONFIG_SAMA7G5)
	/* Enable the internal 32 kHz RC oscillator for low power by writing a 1 to the RCEN bit. */
	reg = readl(AT91C_BASE_SCKCR);
	reg |= AT91C_SLCKSEL_RCEN;
	writel(reg, AT91C_BASE_SCKCR);

	/* Wait internal 32 kHz RC startup time for clock stabilization (software loop). */
	/* 500 us */
	udelay(500);
#endif

	/* Switch from 32768 Hz oscillator to internal RC by writing a 0 to the OSCSEL bit. */
	reg = readl(AT91C_BASE_SCKCR);
	reg &= ~AT91C_SLCKSEL_OSCSEL;
	writel(reg, AT91C_BASE_SCKCR);

	/* Wait 5 slow clock cycles for internal resynchronization. */
	/* 5 slow clock cycles = ~153 us (5 / 32768) */
	udelay(153);

#if !defined(CONFIG_SAMA5D4) && !defined(CONFIG_SAMA5D2)
	/* Disable the 32768 Hz oscillator by writing a 0 to the OSC32EN bit. */
	reg = readl(AT91C_BASE_SCKCR);
	reg &= ~AT91C_SLCKSEL_OSC32EN;
	writel(reg, AT91C_BASE_SCKCR);
#endif

	return 0;
}

#ifdef CONFIG_SCLK_BYPASS
static int slowclk_osc32_bypass(void)
{
	unsigned int reg;

	/* Bypass the 32kHz oscillator by using an external clock
	 * Set OSC32BYP=1 and OSC32EN=0 atomically
	 */
	reg = readl(AT91C_BASE_SCKCR);
	/* 32768 Hz oscillator is bypassed, accept an external slow clock on XIN32 */
	reg |= AT91C_SLCKSEL_OSC32BYP;
	/* 32768 Hz oscillator is disabled. */
	reg &= ~AT91C_SLCKSEL_OSC32EN;
	writel(reg, AT91C_BASE_SCKCR);

	/*
	 * Switching from internal 32kHz RC oscillator to 32768 Hz oscillator
	 * by setting the bit OSCSEL to 1
	 */
	reg = readl(AT91C_BASE_SCKCR);
	reg |= AT91C_SLCKSEL_OSCSEL;
	writel(reg, AT91C_BASE_SCKCR);

	/*
	 * Waiting 5 slow clock cycles for internal resynchronization
	 * 5 slow clock cycles = ~153 us (5 / 32768)
	 */
	udelay(153);

#if !defined(CONFIG_SAMA5D4) && !defined(CONFIG_SAMA5D2) \
	&& !defined(CONFIG_SAM9X60) && !defined(CONFIG_SAMA7G5)

	/*
	 * Disable the 32kHz RC oscillator by setting the bit RCEN to 0
	 */
	slowclk_disable_rc32();
#endif /* #if !defined(CONFIG_SAMA5D4) && !defined(CONFIG_SAMA5D2) \
	&& !defined(CONFIG_SAM9X60) && !defined(CONFIG_SAMA7G5) */

	return 0;
}

/* Accept an external slow clock on XIN32 */
int slowclk_switch_osc32_bypass(void)
{
	/*
	 * On Atmel sama5d3x, if no quartz is connected to the 32768 Hz oscillator,
	 * the SoC must not be resetted while OSC32BYP=1 and OSC32EN=1 and OSCSEL=1,
	 * because it would cause the SoC to be blocked in reset.
	 */

	/* switch to internal RC32 ocsillator, in order to be sure to
	 * have a reliable source for slow clock, at any moment during
	 * the following steps
	 */
	slowclk_switch_rc32();

	/* Enable OSC32, as it is needed for power supply of the osc by-pass cell
	 */
	slowclk_enable_osc32();

	slowclk_osc32_bypass();

	return 0;
}
#endif /* CONFIG_SCLK_BYPASS */

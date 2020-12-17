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
#include "svc_mgr.h"
#include "arch/at91_pmc/pmc.h"
#include "pmc.h"
#include "matrix.h"
#include "l2cc.h"
#include "debug.h"
#include "rstc.h"
#include "watchdog.h"

/*
 * svc_mgr_main - C entry point of the secure world when a SMC is processed
 * in Normal World
 */
int svc_mgr_main(struct smc_args_t const *args)
{
	int ret = 0;
	unsigned int pck_mask;
	unsigned int silent = 1;

	dbg_loud("--> svc_mgr_main\n");

	switch (args->r0) {
	case 0x24:
		if (args->r1 == PMC_PLLAR
			|| args->r1 == PMC_MCKR
			|| args->r1 == PMC_PCKR
			|| args->r1 == PMC_PCKR1
			|| args->r1 == PMC_PCKR2)
			ret = pmc_read_reg(args->r1);
		else
			ret = -1;
		break;

	case 0x25:
		if (is_peripheral_secure(args->r1))
			ret = -1;
		else if (is_switching_clock_forbiden(args->r1, args->r2, &silent))
			if (silent)
				ret = 0;
			else
				ret = -1;
		else
			ret = pmc_periph_clk(args->r1, args->r2);
		break;

	case 0x26:
		if (is_sys_clk_secure(args->r1) && is_pck_clk_secure(args->r1))
			ret = -1;
		else
			ret = pmc_sys_clk(args->r1, args->r2);
		break;

	case 0x27:
		if (is_usb_hs_secure())
			ret = -1;
		else
			ret = pmc_uckr_clk(args->r1);
		break;

	case 0x28:
		if (is_usb_host_secure())
			ret = -1;
		else
			ret = pmc_usb_setup();
		break;
	case 0x29:
		cpu_reset();

	case 0x23:
		switch (args->r1) {
		case PMC_PCKR:
			pck_mask = AT91C_PMC_PCK0;
			break;
		case PMC_PCKR1:
			pck_mask = AT91C_PMC_PCK1;
			break;
		case PMC_PCKR2:
			pck_mask = AT91C_PMC_PCK2;
			break;
		default:
			return -1;
		}

		if (is_pck_clk_secure(pck_mask)) {
			ret = -1;
		} else {
			pmc_pck_setup(args->r1, args->r2);
			ret = 0;
		}
		break;

	case 0x42:
		l2cache_enable();
		break;
	case 0x50:
		pmc_smd_setup(args->r1);
		break;

	case 0x60:
		ret = at91_wdt_set_counter(args->r1);
		break;

	case 0x61:
		ret = at91_wdt_reload_counter();
		break;

	default:
		dbg_info("svc mgr error: SMC ID (%d) not defined\n",
							args->r0);
		ret = -1;
		break;
	}

	dbg_loud("<-- svc_mgr_main\n");

	return ret;
}

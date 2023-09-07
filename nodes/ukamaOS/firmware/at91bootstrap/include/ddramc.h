/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2006, Atmel Corporation
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
#ifndef __DDRAMC_H__
#define __DDRAMC_H__

struct ddramc_register {
	unsigned int mdr;
	unsigned int cr;
	unsigned int rtr;
	unsigned int t0pr;
	unsigned int t1pr;
	unsigned int t2pr;
	unsigned int lpr;
	unsigned int lpddr2_lpr;
	unsigned int tim_calr;
	unsigned int cal_mr4r;
};

extern int ddram_initialize(unsigned int base_address,
		unsigned int ram_address,
		struct ddramc_register *ddramc_config);

extern int lpddr2_sdram_initialize(unsigned int base_address,
				unsigned int ram_address,
				struct ddramc_register *ddramc_config);

extern int ddr3_sdram_initialize(unsigned int base_address,
			unsigned int ram_address,
			struct ddramc_register *ddramc_config);

extern int lpddr3_sdram_initialize(unsigned int base_address,
				   unsigned int ram_address,
				   struct ddramc_register *ddramc_config);

extern int lpddr1_sdram_initialize(unsigned int base_address,
				   unsigned int ram_address,
				   struct ddramc_register *ddramc_config);

extern void ddramc_dump_regs(unsigned int base_address);

#endif /* #ifndef __DDRAMC_H__ */

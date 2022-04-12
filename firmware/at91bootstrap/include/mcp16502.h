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
#ifndef __MCP16502_H__
#define __MCP16502_H__

#include "gpio.h"
#include "twi.h"

#define MCP16502_MIN	(0)
#define MCP16502_BUCK1	(MCP16502_MIN)
#define MCP16502_BUCK2	(1)
#define MCP16502_BUCK3	(2)
#define MCP16502_BUCK4	(3)
#define MCP16502_LDO1	(4)
#define MCP16502_LDO2	(5)
#define MCP16502_MAX	(MCP16502_LDO2)

/**
 * struct mcp16502_cfg - MCP16502 configuration descriptor
 * @regulator:	regulator identifier
 * @uV:		regulator voltage (in microvolts)
 * @enable:	regulator enable status
 */
struct mcp16502_cfg {
	unsigned int regulator;
	unsigned int uV;
	unsigned int enable;
};

extern int mcp16502_init(int busid, int addr, const struct pio_desc *lpm_desc,
			 const struct mcp16502_cfg *cfg, unsigned int cfg_no);
extern int mcp16502_regulator_get_voltage(unsigned int regid);
extern int mcp16502_regulator_set_voltage(unsigned int regid, unsigned int uV);
extern int mcp16502_regulator_set_enable(unsigned int regid,
			 unsigned int enable);
extern unsigned int mcp16502_regulator_get_enable(unsigned int regid);

#endif

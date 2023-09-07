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
#include "common.h"
#include "debug.h"
#include "hardware.h"
#include "mcp16502.h"
#include "twi.h"

#define MCP16502_BASE(_i)		(((_i) + 1) << 4)
#define MCP16502_LOW_SEL		(0x0D)
#define MCP16502_VSEL			(0x3F)
#define MCP16502_EN			(1 << 7)

/**
 * struct mcp16502_priv - MCP16502 private data structure
 * @busid: TWI bus ID
 * @addr: TWI device address
 */
static struct mcp16502_priv {
	int busid;
	int addr;
} mcp16502;

/* Regulators constrains. */
static struct regulator {
	unsigned int min_uV;
	unsigned int max_uV;
	unsigned int step_uV;
} regulators[] = {
	[MCP16502_BUCK1] = { 1200000, 3700000, 50000, },
	[MCP16502_BUCK2...MCP16502_BUCK4] = { 600000, 1850000, 25000, },
	[MCP16502_LDO1...MCP16502_LDO2] = { 1200000, 3700000, 50000, },
};

/**
 * mcp16502_regulator_set_voltage()	- set regulator voltage
 *
 * @regid:	regulator identifier
 * @uV:		regulator voltage (in microvolts)
 *
 * Returns:	0 on success, negative number in case of failure
 */
int mcp16502_regulator_set_voltage(unsigned int regid, unsigned int uV)
{
	unsigned char steps, selector, val;
	int ret;

	if (regid < MCP16502_MIN || regid > MCP16502_MAX ||
	    uV < regulators[regid].min_uV || uV > regulators[regid].max_uV)
		return -1;

	ret = twi_read(mcp16502.busid, mcp16502.addr, MCP16502_BASE(regid),
		       1, &val, 1);
	if (ret)
		return ret;

	steps = (uV - regulators[regid].min_uV) / regulators[regid].step_uV;
	selector = (MCP16502_LOW_SEL + steps) & MCP16502_VSEL;
	val &= ~MCP16502_VSEL;
	val |= selector;

	ret = twi_write(mcp16502.busid, mcp16502.addr, MCP16502_BASE(regid),
			1, &val, 1);

	return ret;
}

/**
 * mcp16502_regulator_get_voltage()	- get regulator voltage
 *
 * @regid:	regulator identifier
 *
 * Returns:	voltage value (in microvolts) on success, negative number in
 *		case of failure
 */
int mcp16502_regulator_get_voltage(unsigned int regid)
{
	unsigned char val, steps;
	unsigned int uV;
	int ret;

	if (regid < MCP16502_MIN || regid > MCP16502_MAX)
		return -1;

	ret = twi_read(mcp16502.busid, mcp16502.addr, MCP16502_BASE(regid),
		       1, &val, 1);
	if (ret)
		return ret;

	steps = (val & MCP16502_VSEL) - MCP16502_LOW_SEL;
	uV = regulators[regid].min_uV + regulators[regid].step_uV * steps;

	/* Should not happen. */
	if (uV < regulators[regid].min_uV || uV > regulators[regid].max_uV)
		return -1;

	return uV;
}

/**
 * mcp16502_regulator_set_enable()	- enable/disable regulator
 *
 * @regid:	regulator identifier
 * @enable:	new regulator status (0 - disable, 1 - enable)
 *
 * Returns:	0 on success, negative number in case of failure
 */
int mcp16502_regulator_set_enable(unsigned int regid, unsigned int enable)
{
	unsigned char val;
	int ret;

	if (regid < MCP16502_MIN || regid > MCP16502_MAX)
		return -1;

	ret = twi_read(mcp16502.busid, mcp16502.addr, MCP16502_BASE(regid),
		       1, &val, 1);

	if ((enable && (val & MCP16502_EN)) ||
	    (!enable && !(val & MCP16502_EN)))
		return 0;

	if (enable)
		val |= MCP16502_EN;
	else
		val &= ~MCP16502_EN;

	ret = twi_write(mcp16502.busid, mcp16502.addr, MCP16502_BASE(regid),
			1, &val, 1);

	return ret;
}

/**
 * mcp16502_regulator_get_enable()	- get regulator enable status
 *
 * @regid:	regulator identifier
 *
 * Returns:	regulator status (0 - disable, 1 - enable) or
 *		negative number in case of failure
 */
unsigned int mcp16502_regulator_get_enable(unsigned int regid)
{
	unsigned char val;
	int ret;

	if (regid < MCP16502_MIN || regid > MCP16502_MAX)
		return 0;

	ret = twi_read(mcp16502.busid, mcp16502.addr, MCP16502_BASE(regid),
		       1, &val, 1);
	if (ret)
		return 0;

	return !!(val & MCP16502_EN);
}

/**
 * mcp16502_init()		- init MCP16502 PMIC
 *
 * @dev:	i2c device information
 * @desc:	low power mode pin descriptor
 * @cfgs:	regulators configuration
 * @cfg_no:	number of regulators to be configured
 *
 * Returns:	0 on success, negative number in case of failure
 */
int mcp16502_init(int busid, int addr, const struct pio_desc *lpm_desc,
		  const struct mcp16502_cfg *cfgs, unsigned int cfgs_no)
{
	int i, ret;

	if (busid < 0 || addr < 0)
		return -1;

	mcp16502.busid = busid;
	mcp16502.addr = addr;

	/* Make sure PMIC is in active state. */
	if (lpm_desc)
		pio_configure(lpm_desc);

	/* Setup regulators. */
	for (i = 0; i < cfgs_no && cfgs; i++) {
		ret = mcp16502_regulator_set_voltage(cfgs[i].regulator,
						     cfgs[i].uV);
		if (ret) {
			dbg_very_loud("regulator (%d) set voltage failed\n",
					cfgs[i].regulator);
			return ret;
		}

		ret = mcp16502_regulator_set_enable(cfgs[i].regulator,
						    cfgs[i].enable);
		if (ret) {
			dbg_very_loud("regulator (%d) enable failed\n",
				      cfgs[i].regulator);
			return ret;
		}
	}

	return 0;
}

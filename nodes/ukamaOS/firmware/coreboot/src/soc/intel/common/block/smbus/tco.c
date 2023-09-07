/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2018 Intel Corporation.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <arch/io.h>
#include <device/pci_ops.h>
#include <device/device.h>
#include <device/pci.h>
#include <device/pci_def.h>
#include <intelblocks/pcr.h>
#include <intelblocks/tco.h>
#include <soc/iomap.h>
#include <soc/pci_devs.h>
#include <soc/pcr_ids.h>
#include <soc/pm.h>
#include <soc/smbus.h>

#define PCR_DMI_TCOBASE 0x2778
/* Enable TCO I/O range decode. */
#define TCOEN (1 << 1)

/* SMBUS TCO base address. */
#define TCOBASE		0x50
#define TCOCTL		0x54
#define  TCO_BASE_EN		(1 << 8)

/* Get base address of TCO I/O registers. */
static uint16_t tco_get_bar(void)
{
	return TCO_BASE_ADDRESS;
}

uint16_t tco_read_reg(uint16_t tco_reg)
{
	uint16_t tcobase;

	tcobase = tco_get_bar();

	return inw(tcobase + tco_reg);
}

void tco_write_reg(uint16_t tco_reg, uint16_t value)
{
	uint16_t tcobase;

	tcobase = tco_get_bar();

	outw(value, tcobase + tco_reg);
}

void tco_lockdown(void)
{
	uint16_t tcocnt;

	/* TCO Lock down */
	tcocnt = tco_read_reg(TCO1_CNT);
	tcocnt |= TCO_LOCK;
	tco_write_reg(TCO1_CNT, tcocnt);
}

uint32_t tco_reset_status(void)
{
	uint16_t tco1_sts;
	uint16_t tco2_sts;

	/* TCO Status 2 register */
	tco2_sts = tco_read_reg(TCO2_STS);
	tco2_sts |= TCO_STS_SECOND_TO;
	tco_write_reg(TCO2_STS, tco2_sts);

	/* TCO Status 1 register */
	tco1_sts = tco_read_reg(TCO1_STS);

	return (tco2_sts << 16) | tco1_sts;
}

/* Stop TCO timer */
static void tco_timer_disable(void)
{
	uint16_t tcocnt;

	/* Program TCO timer halt */
	tcocnt = tco_read_reg(TCO1_CNT);
	tcocnt |= TCO_TMR_HLT;
	tco_write_reg(TCO1_CNT, tcocnt);
}

/* Enable TCO BAR using SMBUS TCO base to access TCO related register */
static void tco_enable_bar(void)
{
	uint32_t reg32;
	uint16_t tcobase;
#if defined(__SIMPLE_DEVICE__)
	int devfn = PCH_DEVFN_SMBUS;
	pci_devfn_t dev = PCI_DEV(0, PCI_SLOT(devfn), PCI_FUNC(devfn));
#else
	struct device *dev;
	dev = PCH_DEV_SMBUS;
#endif

	/* Disable TCO in SMBUS Device first before changing Base Address */
	reg32 = pci_read_config32(dev, TCOCTL);
	reg32 &= ~TCO_BASE_EN;
	pci_write_config32(dev, TCOCTL, reg32);

	/* Program TCO Base */
	tcobase = tco_get_bar();
	pci_write_config32(dev, TCOBASE, tcobase);

	/* Enable TCO in SMBUS */
	pci_write_config32(dev, TCOCTL, reg32 | TCO_BASE_EN);

	/*
	* Program "TCO Base Address" PCR[DMI] + 2778h[15:5, 1]
	*/
	pcr_write32(PID_DMI, PCR_DMI_TCOBASE, tcobase | TCOEN);
}

/*
 * Enable TCO BAR using SMBUS TCO base to access TCO related register
 * also disable the timer.
 */
void tco_configure(void)
{
	if (CONFIG(SOC_INTEL_COMMON_BLOCK_TCO_ENABLE_THROUGH_SMBUS))
		tco_enable_bar();

	tco_timer_disable();
}

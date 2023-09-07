/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2007-2009 coresystems GmbH
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

#ifndef _PCI_MMIO_CFG_H
#define _PCI_MMIO_CFG_H

#include <stdint.h>
#include <device/mmio.h>
#include <device/pci_type.h>

#if !defined(__ROMCC__)

/* By not assigning this to CONFIG_MMCONF_BASE_ADDRESS here we
 * prevent some sub-optimal constant folding. */
extern u8 *const pci_mmconf;

/* Using a unique datatype for MMIO writes makes the pointers to _not_
 * qualify for pointer aliasing with any other objects in memory.
 *
 * MMIO offset is a value originally derived from 'struct device *'
 * in ramstage. For the compiler to not discard this MMIO offset value
 * from CPU registers after any MMIO writes, -fstrict-aliasing has to
 * be also set for the build.
 *
 * Bottom 12 bits (4 KiB) are reserved to address the registers of a
 * single PCI function. Declare the bank as a union to avoid some casting
 * in the functions below.
 */
union pci_bank {
	uint8_t reg8[4096];
	uint16_t reg16[4096 / sizeof(uint16_t)];
	uint32_t reg32[4096 / sizeof(uint32_t)];
};

static __always_inline
volatile union pci_bank *pcicfg(pci_devfn_t dev)
{
	return (void *)&pci_mmconf[PCI_DEVFN_OFFSET(dev)];
}

static __always_inline
uint8_t pci_mmio_read_config8(pci_devfn_t dev, uint16_t reg)
{
	return pcicfg(dev)->reg8[reg];
}

static __always_inline
uint16_t pci_mmio_read_config16(pci_devfn_t dev, uint16_t reg)
{
	return pcicfg(dev)->reg16[reg / sizeof(uint16_t)];
}

static __always_inline
uint32_t pci_mmio_read_config32(pci_devfn_t dev, uint16_t reg)
{
	return pcicfg(dev)->reg32[reg / sizeof(uint32_t)];
}

static __always_inline
void pci_mmio_write_config8(pci_devfn_t dev, uint16_t reg, uint8_t value)
{
	pcicfg(dev)->reg8[reg] = value;
}

static __always_inline
void pci_mmio_write_config16(pci_devfn_t dev, uint16_t reg, uint16_t value)
{
	pcicfg(dev)->reg16[reg / sizeof(uint16_t)] = value;
}

static __always_inline
void pci_mmio_write_config32(pci_devfn_t dev, uint16_t reg, uint32_t value)
{
	pcicfg(dev)->reg32[reg / sizeof(uint32_t)] = value;
}

static __always_inline
uint8_t *pci_mmio_config8_addr(pci_devfn_t dev, uint16_t reg)
{
	return (uint8_t *)&pcicfg(dev)->reg8[reg];
}

static __always_inline
uint16_t *pci_mmio_config16_addr(pci_devfn_t dev, uint16_t reg)
{
	return (uint16_t *)&pcicfg(dev)->reg16[reg / sizeof(uint16_t)];
}

static __always_inline
uint32_t *pci_mmio_config32_addr(pci_devfn_t dev, uint16_t reg)
{
	return (uint32_t *)&pcicfg(dev)->reg32[reg / sizeof(uint32_t)];
}

#endif /* !defined(__ROMCC__) */

#if CONFIG(MMCONF_SUPPORT)

#if CONFIG_MMCONF_BASE_ADDRESS == 0
#error "CONFIG_MMCONF_BASE_ADDRESS undefined!"
#endif

/* Avoid name collisions as different stages have different signature
 * for these functions. The _s_ stands for simple, fundamental IO or
 * MMIO variant.
 */

static __always_inline
uint8_t pci_s_read_config8(pci_devfn_t dev, uint16_t reg)
{
	return pci_mmio_read_config8(dev, reg);
}

static __always_inline
uint16_t pci_s_read_config16(pci_devfn_t dev, uint16_t reg)
{
	return pci_mmio_read_config16(dev, reg);
}

static __always_inline
uint32_t pci_s_read_config32(pci_devfn_t dev, uint16_t reg)
{
	return pci_mmio_read_config32(dev, reg);
}

static __always_inline
void pci_s_write_config8(pci_devfn_t dev, uint16_t reg, uint8_t value)
{
	pci_mmio_write_config8(dev, reg, value);
}

static __always_inline
void pci_s_write_config16(pci_devfn_t dev, uint16_t reg, uint16_t value)
{
	pci_mmio_write_config16(dev, reg, value);
}

static __always_inline
void pci_s_write_config32(pci_devfn_t dev, uint16_t reg, uint32_t value)
{
	pci_mmio_write_config32(dev, reg, value);
}

#endif

#endif /* _PCI_MMIO_CFG_H */

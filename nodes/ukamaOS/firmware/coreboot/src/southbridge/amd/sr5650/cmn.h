/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2010 Advanced Micro Devices, Inc.
 * Copyright (C) 2015 Timothy Pearson <tpearson@raptorengineeringinc.com>, Raptor Engineering
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

#ifndef __SR5650_CMN_H__
#define __SR5650_CMN_H__

#include <device/pci_ops.h>

#define NBMISC_INDEX	0x60
#define NBHTIU_INDEX	0x94 /* Note: It is different with RS690, whose HTIU index is 0xA8 */
#define NBMC_INDEX	0xE8
#define NBPCIE_INDEX	0xE0
#define L2CFG_INDEX	0xF0
#define L1CFG_INDEX	0xF8
#define EXT_CONF_BASE_ADDRESS	CONFIG_MMCONF_BASE_ADDRESS
#define	TEMP_MMIO_BASE_ADDRESS	0xC0000000

#define axindxc_reg(reg, mask, val) \
	alink_ax_indx(0, (reg), (mask), (val))

#define AB_INDX   0xCD8
#define AB_DATA   (AB_INDX+4)

#if ENV_PCI_SIMPLE_DEVICE
static inline u32 nb_read_index(pci_devfn_t dev, u32 index_reg, u32 index)
#else
static inline u32 nb_read_index(struct device *dev, u32 index_reg, u32 index)
#endif
{
	pci_write_config32(dev, index_reg, index);
	return pci_read_config32(dev, index_reg + 0x4);
}

#if ENV_PCI_SIMPLE_DEVICE
static inline void nb_write_index(pci_devfn_t dev, u32 index_reg, u32 index,
				  u32 data)
#else
static inline void nb_write_index(struct device *dev, u32 index_reg, u32 index,
				  u32 data)
#endif
{
	pci_write_config32(dev, index_reg, index);
	pci_write_config32(dev, index_reg + 0x4, data);
}

#if ENV_PCI_SIMPLE_DEVICE
static inline u32 nbmisc_read_index(pci_devfn_t nb_dev, u32 index)
#else
static inline u32 nbmisc_read_index(struct device *nb_dev, u32 index)
#endif
{
	return nb_read_index((nb_dev), NBMISC_INDEX, (index));
}

#if ENV_PCI_SIMPLE_DEVICE
static inline void nbmisc_write_index(pci_devfn_t nb_dev, u32 index, u32 data)
#else
static inline void nbmisc_write_index(struct device *nb_dev, u32 index,
				      u32 data)
#endif
{
	nb_write_index((nb_dev), NBMISC_INDEX, ((index) | 0x80), (data));
}

#if ENV_PCI_SIMPLE_DEVICE
static inline void set_nbmisc_enable_bits(pci_devfn_t nb_dev, u32 reg_pos,
					  u32 mask, u32 val)
#else
static inline void set_nbmisc_enable_bits(struct device *nb_dev, u32 reg_pos,
					  u32 mask, u32 val)
#endif
{
	u32 reg_old, reg;
	reg = reg_old = nbmisc_read_index(nb_dev, reg_pos);
	reg &= ~mask;
	reg |= val;
	if (reg != reg_old) {
		nbmisc_write_index(nb_dev, reg_pos, reg);
	}
}

#if ENV_PCI_SIMPLE_DEVICE
static inline u32 htiu_read_index(pci_devfn_t nb_dev, u32 index)
#else
static inline u32 htiu_read_index(struct device *nb_dev, u32 index)
#endif
{
	return nb_read_index((nb_dev), NBHTIU_INDEX, (index));
}

#if ENV_PCI_SIMPLE_DEVICE
static inline void htiu_write_index(pci_devfn_t nb_dev, u32 index, u32 data)
#else
static inline void htiu_write_index(struct device *nb_dev, u32 index, u32 data)
#endif
{
	nb_write_index((nb_dev), NBHTIU_INDEX, ((index) | 0x100), (data));
}

#if ENV_PCI_SIMPLE_DEVICE
static inline u32 nbmc_read_index(pci_devfn_t nb_dev, u32 index)
#else
static inline u32 nbmc_read_index(struct device *nb_dev, u32 index)
#endif
{
	return nb_read_index((nb_dev), NBMC_INDEX, (index));
}

#if ENV_PCI_SIMPLE_DEVICE
static inline void nbmc_write_index(pci_devfn_t nb_dev, u32 index, u32 data)
#else
static inline void nbmc_write_index(struct device *nb_dev, u32 index, u32 data)
#endif
{
	nb_write_index((nb_dev), NBMC_INDEX, ((index) | 1 << 9), (data));
}

#if ENV_PCI_SIMPLE_DEVICE
static inline void set_htiu_enable_bits(pci_devfn_t nb_dev, u32 reg_pos,
					u32 mask, u32 val)
#else
static inline void set_htiu_enable_bits(struct device *nb_dev, u32 reg_pos,
					u32 mask, u32 val)
#endif
{
	u32 reg_old, reg;
	reg = reg_old = htiu_read_index(nb_dev, reg_pos);
	reg &= ~mask;
	reg |= val;
	if (reg != reg_old) {
		htiu_write_index(nb_dev, reg_pos, reg);
	}
}

#if ENV_PCI_SIMPLE_DEVICE
static inline void set_nbcfg_enable_bits(pci_devfn_t nb_dev, u32 reg_pos,
					 u32 mask, u32 val)
#else
static inline void set_nbcfg_enable_bits(struct device *nb_dev, u32 reg_pos,
					 u32 mask, u32 val)
#endif
{
	u32 reg_old, reg;
	reg = reg_old = pci_read_config32(nb_dev, reg_pos);
	reg &= ~mask;
	reg |= val;
	if (reg != reg_old) {
		pci_write_config32(nb_dev, reg_pos, reg);
	}
}

#if ENV_PCI_SIMPLE_DEVICE
static inline void set_nbcfg_enable_bits_8(pci_devfn_t nb_dev, u32 reg_pos,
					   u8 mask, u8 val)
#else
static inline void set_nbcfg_enable_bits_8(struct device *nb_dev, u32 reg_pos,
					   u8 mask, u8 val)
#endif
{
	u8 reg_old, reg;
	reg = reg_old = pci_read_config8(nb_dev, reg_pos);
	reg &= ~mask;
	reg |= val;
	if (reg != reg_old) {
		pci_write_config8(nb_dev, reg_pos, reg);
	}
}

#if ENV_PCI_SIMPLE_DEVICE
static inline void set_nbmc_enable_bits(pci_devfn_t nb_dev, u32 reg_pos,
					u32 mask, u32 val)
#else
static inline void set_nbmc_enable_bits(struct device *nb_dev, u32 reg_pos,
					u32 mask, u32 val)
#endif
{
	u32 reg_old, reg;
	reg = reg_old = nbmc_read_index(nb_dev, reg_pos);
	reg &= ~mask;
	reg |= val;
	if (reg != reg_old) {
		nbmc_write_index(nb_dev, reg_pos, reg);
	}
}

#if ENV_PCI_SIMPLE_DEVICE
static inline void set_pcie_enable_bits(pci_devfn_t dev, u32 reg_pos, u32 mask,
					u32 val)
#else
static inline void set_pcie_enable_bits(struct device *dev, u32 reg_pos,
					u32 mask, u32 val)
#endif
{
	u32 reg_old, reg;
	reg = reg_old = nb_read_index(dev, NBPCIE_INDEX, reg_pos);
	reg &= ~mask;
	reg |= val;
	if (reg != reg_old) {
		nb_write_index(dev, NBPCIE_INDEX, reg_pos, reg);
	}
}

void set_pcie_reset(void);
void set_pcie_dereset(void);

#endif /* __SR5650_CMN_H__ */

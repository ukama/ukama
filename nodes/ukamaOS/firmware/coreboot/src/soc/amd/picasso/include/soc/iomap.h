/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2017 Raptor Engineering, LLC
 * Copyright 2017 Google Inc.
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

#ifndef __SOC_PICASSO_IOMAP_H__
#define __SOC_PICASSO_IOMAP_H__

/* MMIO Ranges */
/* IO_APIC_ADDR defined in arch/x86	0xfec00000 */
#define SPI_BASE_ADDRESS		0xfec10000
#define ESPI_BASE_ADDRESS		0xfec20000

#if CONFIG(HPET_ADDRESS_OVERRIDE)
#error HPET address override is not allowed and must be fixed at 0xfed00000
#endif
#define HPET_BASE_ADDRESS		0xfed00000

/*
 * AcpiMmio blocks are at fixed offsets from FED8_0000h and enabled in PMx04[1].
 * All ranges not specified as supported below may, or may not, be listed in
 * any documentation but should be considered reserved through FED8_1FFFh.
 */
#include <amdblocks/acpimmio_map.h>
#define SUPPORTS_ACPIMMIO_SM_PCI_BASE	1 /* 0xfed80000 */
#define SUPPORTS_ACPIMMIO_SMI_BASE	1 /* 0xfed80100 */
#define SUPPORTS_ACPIMMIO_PMIO_BASE	1 /* 0xfed80300 */
#define SUPPORTS_ACPIMMIO_BIOSRAM_BASE	1 /* 0xfed80500 */
#define SUPPORTS_ACPIMMIO_ACPI_BASE	1 /* 0xfed80800 */
#define SUPPORTS_ACPIMMIO_ASF_BASE	1 /* 0xfed80900 */
#define SUPPORTS_ACPIMMIO_SMBUS_BASE	1 /* 0xfed80a00 */
#define SUPPORTS_ACPIMMIO_IOMUX_BASE	1 /* 0xfed80d00 */
#define SUPPORTS_ACPIMMIO_MISC_BASE	1 /* 0xfed80e00 */
#define SUPPORTS_ACPIMMIO_GPIO0_BASE	1 /* 0xfed81500 */
#define SUPPORTS_ACPIMMIO_GPIO1_BASE	1 /* 0xfed81800 */
#define SUPPORTS_ACPIMMIO_GPIO2_BASE	1 /* 0xfed81700 */
#define SUPPORTS_ACPIMMIO_AOAC_BASE	1 /* 0xfed81e00 */

#define ALINK_AHB_ADDRESS		0xfedc0000

/* Reserved				0xfecd1000-0xfedc3fff */

#define APU_I2C2_BASE			0xfedc4000
#define APU_I2C3_BASE			0xfedc5000
#define APU_I2C4_BASE			0xfedc6000
#define   APU_I2C_MIN_BUS		2
#define   APU_I2C_MAX_BUS		4
#define   APU_I2C_BLOCK_SIZE		0x1000
#define   I2C_BASE_ADDRESS		APU_I2C2_BASE
#define   I2C_DEVICE_SIZE		0x00001000
#define   I2C_DEVICE_COUNT		3

#define APU_DMAC0_BASE			0xfedc7000
#define APU_DMAC1_BASE			0xfedc8000
#define APU_UART0_BASE			0xfedc9000
#define APU_UART1_BASE			0xfedca000
/* Reserved				0xfedcb000 */
#define APU_DMAC2_BASE			0xfedcc000
#define APU_DMAC3_BASE			0xfedcd000
#define APU_UART2_BASE			0xfedce000
#define APU_UART3_BASE			0xfedcf000
/* Reserved				0xfedd0000-0xfedd4fff */
#define APU_EMMC_BASE			0xfedd5000
#define APU_EMMC_CONFIG_BASE		0xfedd5800

#define FLASH_BASE_ADDR			((0xffffffff - CONFIG_ROM_SIZE) + 1)

/* I/O Ranges */
#define ACPI_SMI_CTL_PORT		0xb2
#define PICASSO_ACPI_IO_BASE	CONFIG_PICASSO_ACPI_IO_BASE
#define  ACPI_PM_EVT_BLK	(PICASSO_ACPI_IO_BASE + 0x00)     /* 4 bytes */
#define  ACPI_PM1_STS		(ACPI_PM_EVT_BLK + 0x00)	  /* 2 bytes */
#define  ACPI_PM1_EN		(ACPI_PM_EVT_BLK + 0x02)	  /* 2 bytes */
#define  ACPI_PM1_CNT_BLK	(PICASSO_ACPI_IO_BASE + 0x04)     /* 2 bytes */
#define  ACPI_CPU_CONTROL	(PICASSO_ACPI_IO_BASE + 0x08)     /* 6 bytes */
#define  ACPI_GPE0_BLK		(PICASSO_ACPI_IO_BASE + 0x10)     /* 8 bytes */
#define  ACPI_GPE0_STS		(ACPI_GPE0_BLK + 0x00)		  /* 4 bytes */
#define  ACPI_GPE0_EN		(ACPI_GPE0_BLK + 0x04)		  /* 4 bytes */
#define  ACPI_PM_TMR_BLK	(PICASSO_ACPI_IO_BASE + 0x18)     /* 4 bytes */
#define NCP_ERR				0xf0
#define SMB_BASE_ADDR			0xb00
#define PM2_INDEX			0xcd0
#define PM2_DATA			0xcd1
#define BIOSRAM_INDEX			0xcd4
#define BIOSRAM_DATA			0xcd5
#define AB_INDX				0xcd8
#define AB_DATA				(AB_INDX+4)
#define SYS_RESET			0xcf9

/* BiosRam Ranges at 0xfed80500 or I/O 0xcd4/0xcd5 */
#define BIOSRAM_CBMEM_TOP		0xf0 /* 4 bytes */
#define BIOSRAM_UMA_SIZE		0xf4 /* 4 bytes */
#define BIOSRAM_UMA_BASE		0xf8 /* 8 bytes */

#endif /* __SOC_PICASSO_IOMAP_H__ */

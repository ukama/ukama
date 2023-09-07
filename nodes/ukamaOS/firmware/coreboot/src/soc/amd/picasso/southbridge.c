/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2010-2017 Advanced Micro Devices, Inc.
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

#include <console/console.h>
#include <device/mmio.h>
#include <bootstate.h>
#include <cpu/x86/smm.h>
#include <cpu/x86/msr.h>
#include <device/device.h>
#include <device/pci.h>
#include <device/pci_ops.h>
#include <cbmem.h>
#include <amdblocks/amd_pci_util.h>
#include <amdblocks/reset.h>
#include <amdblocks/acpimmio.h>
#include <amdblocks/lpc.h>
#include <amdblocks/acpi.h>
#include <soc/cpu.h>
#include <soc/southbridge.h>
#include <soc/smbus.h>
#include <soc/smi.h>
#include <soc/amd_pci_int_defs.h>
#include <delay.h>
#include <soc/pci_devs.h>
#include <soc/nvs.h>
#include <types.h>

#define FCH_AOAC_UART_FOR_CONSOLE \
		(CONFIG_UART_FOR_CONSOLE == 0 ? FCH_AOAC_DEV_UART0 \
		: CONFIG_UART_FOR_CONSOLE == 1 ? FCH_AOAC_DEV_UART1 \
		: CONFIG_UART_FOR_CONSOLE == 2 ? FCH_AOAC_DEV_UART2 \
		: CONFIG_UART_FOR_CONSOLE == 3 ? FCH_AOAC_DEV_UART3 \
		: -1)
#if FCH_AOAC_UART_FOR_CONSOLE == -1
# error Unsupported UART_FOR_CONSOLE chosen
#endif

/*
 * Table of devices that need their AOAC registers enabled and waited
 * upon (usually about .55 milliseconds). Instead of individual delays
 * waiting for each device to become available, a single delay will be
 * executed.  The console UART is handled separately from this table.
 */
const static int aoac_devs[] = {
	FCH_AOAC_DEV_AMBA,
	FCH_AOAC_DEV_I2C2,
	FCH_AOAC_DEV_I2C3,
	FCH_AOAC_DEV_I2C4,
	FCH_AOAC_DEV_ESPI,
};

/*
 * Table of APIC register index and associated IRQ name. Using IDX_XXX_NAME
 * provides a visible association with the index, therefore helping
 * maintainability of table. If a new index/name is defined in
 * amd_pci_int_defs.h, just add the pair at the end of this table.
 * Order is not important.
 */
const static struct irq_idx_name irq_association[] = {
	{ PIRQ_A,	"INTA#" },
	{ PIRQ_B,	"INTB#" },
	{ PIRQ_C,	"INTC#" },
	{ PIRQ_D,	"INTD#" },
	{ PIRQ_E,	"INTE#" },
	{ PIRQ_F,	"INTF#/GENINT2" },
	{ PIRQ_G,	"INTG#" },
	{ PIRQ_H,	"INTH#" },
	{ PIRQ_MISC,	"Misc" },
	{ PIRQ_MISC0,	"Misc0" },
	{ PIRQ_MISC1,	"Misc1" },
	{ PIRQ_MISC2,	"Misc2" },
	{ PIRQ_SIRQA,	"Ser IRQ INTA" },
	{ PIRQ_SIRQB,	"Ser IRQ INTB" },
	{ PIRQ_SIRQC,	"Ser IRQ INTC" },
	{ PIRQ_SIRQD,	"Ser IRQ INTD" },
	{ PIRQ_SCI,	"SCI" },
	{ PIRQ_SMBUS,	"SMBUS" },
	{ PIRQ_ASF,	"ASF" },
	{ PIRQ_PMON,	"PerMon" },
	{ PIRQ_SD,	"SD" },
	{ PIRQ_SDIO,	"SDIO" },
	{ PIRQ_CIR,	"CIR" },
	{ PIRQ_GPIOA,	"GPIOa" },
	{ PIRQ_GPIOB,	"GPIOb" },
	{ PIRQ_GPIOC,	"GPIOc" },
	{ PIRQ_SATA,	"SATA" },
	{ PIRQ_EMMC,	"eMMC" },
	{ PIRQ_GPP0,	"GPP0" },
	{ PIRQ_GPP1,	"GPP1" },
	{ PIRQ_GPP2,	"GPP2" },
	{ PIRQ_GPP3,	"GPP3" },
	{ PIRQ_GPIO,	"GPIO" },
	{ PIRQ_I2C0,	"I2C0" },
	{ PIRQ_I2C1,	"I2C1" },
	{ PIRQ_I2C2,	"I2C2" },
	{ PIRQ_I2C3,	"I2C3" },
	{ PIRQ_UART0,	"UART0" },
	{ PIRQ_UART1,	"UART1" },
	{ PIRQ_I2C4,	"I2C4" },
	{ PIRQ_I2C5,	"I2C5" },
	{ PIRQ_UART2,	"UART2" },
	{ PIRQ_UART3,	"UART3" },
};

const struct irq_idx_name *sb_get_apic_reg_association(size_t *size)
{
	*size = ARRAY_SIZE(irq_association);
	return irq_association;
}

static void power_on_aoac_device(int dev)
{
	uint8_t byte;

	/* Power on the UART and AMBA devices */
	byte = aoac_read8(AOAC_DEV_D3_CTL(dev));
	byte |= FCH_AOAC_PWR_ON_DEV;
	aoac_write8(AOAC_DEV_D3_CTL(dev), byte);
}

static bool is_aoac_device_enabled(int dev)
{
	uint8_t byte;

	byte = aoac_read8(AOAC_DEV_D3_STATE(dev));
	byte &= (FCH_AOAC_PWR_RST_STATE | FCH_AOAC_RST_CLK_OK_STATE);
	if (byte == (FCH_AOAC_PWR_RST_STATE | FCH_AOAC_RST_CLK_OK_STATE))
		return true;
	else
		return false;
}

static void enable_aoac_console_uart(void)
{
	if (!CONFIG(PICASSO_UART))
		return;

	power_on_aoac_device(FCH_AOAC_UART_FOR_CONSOLE);
}

static bool is_aoac_console_uart_enabled(void)
{
	if (!CONFIG(PICASSO_UART))
		return true;

	return is_aoac_device_enabled(FCH_AOAC_UART_FOR_CONSOLE);
}

void enable_aoac_devices(void)
{
	bool status;
	int i;

	for (i = 0; i < ARRAY_SIZE(aoac_devs); i++)
		power_on_aoac_device(aoac_devs[i]);
	enable_aoac_console_uart();

	/* Wait for AOAC devices to indicate power and clock OK */
	do {
		udelay(100);
		status = true;
		for (i = 0; i < ARRAY_SIZE(aoac_devs); i++)
			status &= is_aoac_device_enabled(aoac_devs[i]);
		status &= is_aoac_console_uart_enabled();
	} while (!status);
}

static void sb_enable_lpc(void)
{
	u8 byte;

	/* Enable LPC controller */
	byte = pm_io_read8(PM_LPC_GATING);
	byte |= PM_LPC_ENABLE;
	pm_io_write8(PM_LPC_GATING, byte);
}

static void sb_enable_cf9_io(void)
{
	uint32_t reg = pm_read32(PM_DECODE_EN);

	pm_write32(PM_DECODE_EN, reg | CF9_IO_EN);
}

static void sb_enable_legacy_io(void)
{
	uint32_t reg = pm_read32(PM_DECODE_EN);

	pm_write32(PM_DECODE_EN, reg | LEGACY_IO_EN);
}

void sb_clk_output_48Mhz(void)
{
	u32 ctrl;

	ctrl = misc_read32(MISC_CLK_CNTL1);
	ctrl |= BP_X48M0_OUTPUT_EN;
	misc_write32(MISC_CLK_CNTL1, ctrl);
}

static uintptr_t sb_init_spi_base(void)
{
	uintptr_t base;

	/* Make sure the base address is predictable */
	base = lpc_get_spibase();

	if (base)
		return base;

	lpc_set_spibase(SPI_BASE_ADDRESS, SPI_ROM_ENABLE);
	return SPI_BASE_ADDRESS;
}

void sb_set_spi100(u16 norm, u16 fast, u16 alt, u16 tpm)
{
	uintptr_t base = sb_init_spi_base();
	write16((void *)(base + SPI100_SPEED_CONFIG),
				(norm << SPI_NORM_SPEED_NEW_SH) |
				(fast << SPI_FAST_SPEED_NEW_SH) |
				(alt << SPI_ALT_SPEED_NEW_SH) |
				(tpm << SPI_TPM_SPEED_NEW_SH));
	write16((void *)(base + SPI100_ENABLE), SPI_USE_SPI100);
}

void sb_disable_4dw_burst(void)
{
	uintptr_t base = sb_init_spi_base();
	write16((void *)(base + SPI100_HOST_PREF_CONFIG),
			read16((void *)(base + SPI100_HOST_PREF_CONFIG))
					& ~SPI_RD4DW_EN_HOST);
}

void sb_read_mode(u32 mode)
{
	uintptr_t base = sb_init_spi_base();
	write32((void *)(base + SPI_CNTRL0),
			(read32((void *)(base + SPI_CNTRL0))
					& ~SPI_READ_MODE_MASK) | mode);
}

static void fch_smbus_init(void)
{
	pm_write8(SMB_ASF_IO_BASE, SMB_BASE_ADDR >> 8);
	smbus_write8(SMBTIMING, SMB_SPEED_400KHZ);
	/* Clear all SMBUS status bits */
	smbus_write8(SMBHSTSTAT, SMBHST_STAT_CLEAR);
	smbus_write8(SMBSLVSTAT, SMBSLV_STAT_CLEAR);
	asf_write8(SMBHSTSTAT, SMBHST_STAT_CLEAR);
	asf_write8(SMBSLVSTAT, SMBSLV_STAT_CLEAR);
}

/* Before console init */
void fch_pre_init(void)
{
	/* Turn on LPC in case the PSP didn't use it.  However, ensure all
	 * decoding is cleared as the PSP may have enabled decode paths. */
	sb_enable_lpc();
	lpc_disable_decodes();

	if (CONFIG(POST_IO) && (CONFIG_POST_IO_PORT == 0x80)
					&& CONFIG(PICASSO_LPC_IOMUX))
		lpc_enable_port80();
	lpc_enable_spi_prefetch();
	sb_init_spi_base();
	sb_disable_4dw_burst();
	sb_set_spi100(SPI_SPEED_33M, SPI_SPEED_33M,
			SPI_SPEED_16M, SPI_SPEED_16M);
	enable_acpimmio_decode();
	fch_smbus_init();
	sb_enable_cf9_io();
	sb_enable_legacy_io();
	enable_aoac_devices();
	sb_reset_i2c_slaves();
	if (CONFIG(PICASSO_UART))
		set_uart_config(CONFIG_UART_FOR_CONSOLE);
}

static void print_num_status_bits(int num_bits, uint32_t status,
				  const char *const bit_names[])
{
	int i;

	if (!status)
		return;

	for (i = num_bits - 1; i >= 0; i--) {
		if (status & (1 << i)) {
			if (bit_names[i])
				printk(BIOS_DEBUG, "%s ", bit_names[i]);
			else
				printk(BIOS_DEBUG, "BIT%d ", i);
		}
	}
}

static void sb_print_pmxc0_status(void)
{
	/* PMxC0 S5/Reset Status shows the source of previous reset. */
	uint32_t pmxc0_status = pm_read32(PM_RST_STATUS);

	static const char *const pmxc0_status_bits[32] = {
		[0] = "ThermalTrip",
		[1] = "FourSecondPwrBtn",
		[2] = "Shutdown",
		[3] = "ThermalTripFromTemp",
		[4] = "RemotePowerDownFromASF",
		[5] = "ShutDownFan0",
		[16] = "UserRst",
		[17] = "SoftPciRst",
		[18] = "DoInit",
		[19] = "DoReset",
		[20] = "DoFullReset",
		[21] = "SleepReset",
		[22] = "KbReset",
		[23] = "LtReset",
		[24] = "FailBootRst",
		[25] = "WatchdogIssueReset",
		[26] = "RemoteResetFromASF",
		[27] = "SyncFlood",
		[28] = "HangReset",
		[29] = "EcWatchdogRst",
	};

	printk(BIOS_DEBUG, "PMxC0 STATUS: 0x%x ", pmxc0_status);
	print_num_status_bits(ARRAY_SIZE(pmxc0_status_bits), pmxc0_status,
			      pmxc0_status_bits);
	printk(BIOS_DEBUG, "\n");
}

/* After console init */
void fch_early_init(void)
{
	sb_print_pmxc0_status();
	i2c_soc_early_init();
}

void sb_enable(struct device *dev)
{
	printk(BIOS_DEBUG, "%s\n", __func__);
}

static void sb_init_acpi_ports(void)
{
	u32 reg;
	msr_t cst_addr;

	/* We use some of these ports in SMM regardless of whether or not
	 * ACPI tables are generated. Enable these ports indiscriminately.
	 */

	pm_write16(PM_EVT_BLK, ACPI_PM_EVT_BLK);
	pm_write16(PM1_CNT_BLK, ACPI_PM1_CNT_BLK);
	pm_write16(PM_TMR_BLK, ACPI_PM_TMR_BLK);
	pm_write16(PM_GPE0_BLK, ACPI_GPE0_BLK);

	/* CpuControl is in \_PR.CP00, 6 bytes */
	cst_addr.hi = 0;
	cst_addr.lo = ACPI_CPU_CONTROL;
	wrmsr(CSTATE_BASE_REG, cst_addr);

	if (CONFIG(HAVE_SMI_HANDLER)) {
		/* APMC - SMI Command Port */
		pm_write16(PM_ACPI_SMI_CMD, APM_CNT);
		configure_smi(SMITYPE_SMI_CMD_PORT, SMI_MODE_SMI);

		/* SMI on SlpTyp requires sending SMI before completion
		 * response of the I/O write.  The BKDG also specifies
		 * clearing ForceStpClkRetry for SMI trapping.
		 */
		reg = pm_read32(PM_PCI_CTRL);
		reg |= FORCE_SLPSTATE_RETRY;
		pm_write32(PM_PCI_CTRL, reg);

		/* Disable SlpTyp feature */
		reg = pm_read8(PM_RST_CTRL1);
		reg &= ~SLPTYPE_CONTROL_EN;
		pm_write8(PM_RST_CTRL1, reg);

		configure_smi(SMITYPE_SLP_TYP, SMI_MODE_SMI);
	} else {
		pm_write16(PM_ACPI_SMI_CMD, 0);
	}

	/* Decode ACPI registers and enable standard features */
	pm_write8(PM_ACPI_CONF, PM_ACPI_DECODE_STD |
				PM_ACPI_GLOBAL_EN |
				PM_ACPI_RTC_EN_EN |
				PM_ACPI_TIMER_EN_EN);
}

static int get_index_bit(uint32_t value, uint16_t limit)
{
	uint16_t i;
	uint32_t t;

	if (limit >= TOTAL_BITS(uint32_t))
		return -1;

	/* get a mask of valid bits. Ex limit = 3, set bits 0-2 */
	t = (1 << limit) - 1;
	if ((value & t) == 0)
		return -1;
	t = 1;
	for (i = 0; i < limit; i++) {
		if (value & t)
			break;
		t <<= 1;
	}
	return i;
}

static void set_nvs_sws(void *unused)
{
	struct soc_power_reg *sws;
	struct global_nvs_t *gnvs;
	int index;

	sws = cbmem_find(CBMEM_ID_POWER_STATE);
	if (sws == NULL)
		return;
	gnvs = cbmem_find(CBMEM_ID_ACPI_GNVS);
	if (gnvs == NULL)
		return;

	index = get_index_bit(sws->pm1_sts & sws->pm1_en, PM1_LIMIT);
	if (index < 0)
		gnvs->pm1i = ~0ULL;
	else
		gnvs->pm1i = index;

	index = get_index_bit(sws->gpe0_sts & sws->gpe0_en, GPE0_LIMIT);
	if (index < 0)
		gnvs->gpei = ~0ULL;
	else
		gnvs->gpei = index;
}

BOOT_STATE_INIT_ENTRY(BS_OS_RESUME, BS_ON_ENTRY, set_nvs_sws, NULL);

void southbridge_init(void *chip_info)
{
	sb_init_acpi_ports();
	acpi_clear_pm1_status();
}

static void set_sb_final_nvs(void)
{
	struct global_nvs_t *gnvs = cbmem_find(CBMEM_ID_ACPI_GNVS);
	if (gnvs == NULL)
		return;

	gnvs->aoac.ic2e = is_aoac_device_enabled(FCH_AOAC_DEV_I2C2);
	gnvs->aoac.ic3e = is_aoac_device_enabled(FCH_AOAC_DEV_I2C3);
	gnvs->aoac.ic4e = is_aoac_device_enabled(FCH_AOAC_DEV_I2C4);
	gnvs->aoac.ut0e = is_aoac_device_enabled(FCH_AOAC_DEV_UART0);
	gnvs->aoac.ut1e = is_aoac_device_enabled(FCH_AOAC_DEV_UART1);
	gnvs->aoac.ut2e = is_aoac_device_enabled(FCH_AOAC_DEV_UART2);
	gnvs->aoac.ut3e = is_aoac_device_enabled(FCH_AOAC_DEV_UART3);
	gnvs->aoac.espi = 1;
}

void southbridge_final(void *chip_info)
{
	uint8_t restored_power = PM_S5_AT_POWER_RECOVERY;

	if (CONFIG(MAINBOARD_POWER_RESTORE))
		restored_power = PM_RESTORE_S0_IF_PREV_S0;
	pm_write8(PM_RTC_SHADOW, restored_power);

	set_sb_final_nvs();
}

/*
 * Update the PCI devices with a valid IRQ number
 * that is set in the mainboard PCI_IRQ structures.
 */
static void set_pci_irqs(void *unused)
{
	/* Write PCI_INTR regs 0xC00/0xC01 */
	write_pci_int_table();

	/* Write IRQs for all devicetree enabled devices */
	write_pci_cfg_irqs();
}

/*
 * Hook this function into the PCI state machine
 * on entry into BS_DEV_ENABLE.
 */
BOOT_STATE_INIT_ENTRY(BS_DEV_ENABLE, BS_ON_ENTRY, set_pci_irqs, NULL);

void save_uma_size(uint32_t size)
{
	biosram_write32(BIOSRAM_UMA_SIZE, size);
}

void save_uma_base(uint64_t base)
{
	biosram_write32(BIOSRAM_UMA_BASE, (uint32_t) base);
	biosram_write32(BIOSRAM_UMA_BASE + 4, (uint32_t) (base >> 32));
}

uint32_t get_uma_size(void)
{
	return biosram_read32(BIOSRAM_UMA_SIZE);
}

uint64_t get_uma_base(void)
{
	uint64_t base;
	base = biosram_read32(BIOSRAM_UMA_BASE);
	base |= ((uint64_t)(biosram_read32(BIOSRAM_UMA_BASE + 4)) << 32);
	return base;
}

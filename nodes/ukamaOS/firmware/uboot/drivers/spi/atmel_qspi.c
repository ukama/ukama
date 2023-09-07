/*
 * Copyright (C) 2017 Atmel Corporation
 *
 * SPDX-License-Identifier:	GPL-2.0+
 */
#include <common.h>
#include <clk.h>
#include <dm.h>
#include <fdtdec.h>
#include <errno.h>
#include <spi.h>
#include <asm/io.h>
#include <mach/clk.h>
#include "atmel_qspi.h"

DECLARE_GLOBAL_DATA_PTR;

static void atmel_qspi_memcpy_fromio(void *dst, unsigned long src, size_t len)
{
	u8 *d = (u8 *)dst;

	while (len--) {
		*d++ = readb(src);
		src++;
	}
}

static void atmel_qspi_memcpy_toio(unsigned long dst, const void *src,
				   size_t len)
{
	const u8 *s = (const u8 *)src;

	while (len--) {
		writeb(*s, dst);
		dst++;
		s++;
	}
}

static int atmel_qspi_set_ifr_tfrtype(u8 flags, u32 *ifr)
{
	u32 ifr_tfrtype;

	switch (flags & SPI_FCMD_TYPE) {
	case SPI_FCMD_READ:
		ifr_tfrtype = QSPI_IFR_TFRTYPE_READ_MEMORY;
		break;

	case SPI_FCMD_WRITE:
		ifr_tfrtype = QSPI_IFR_TFRTYPE_WRITE_MEMORY;
		break;

	case SPI_FCMD_ERASE:
	case SPI_FCMD_WRITE_REG:
		ifr_tfrtype = QSPI_IFR_TFRTYPE_WRITE;
		break;

	case SPI_FCMD_READ_REG:
		ifr_tfrtype = QSPI_IFR_TFRTYPE_READ;
		break;

	default:
		return -EINVAL;
	}

	*ifr = (*ifr & ~QSPI_IFR_TFRTYPE) | ifr_tfrtype;
	return 0;
}

static int atmel_qpsi_set_ifr_width(enum spi_flash_protocol proto, u32 *ifr)
{
	u32 ifr_width;

	switch (proto) {
	case SPI_FPROTO_1_1_1:
		ifr_width = QSPI_IFR_WIDTH_SINGLE_BIT_SPI;
		break;

	case SPI_FPROTO_1_1_2:
		ifr_width = QSPI_IFR_WIDTH_DUAL_OUTPUT;
		break;

	case SPI_FPROTO_1_2_2:
		ifr_width = QSPI_IFR_WIDTH_DUAL_IO;
		break;

	case SPI_FPROTO_2_2_2:
		ifr_width = QSPI_IFR_WIDTH_DUAL_CMD;
		break;

	case SPI_FPROTO_1_1_4:
		ifr_width = QSPI_IFR_WIDTH_QUAD_OUTPUT;
		break;

	case SPI_FPROTO_1_4_4:
		ifr_width = QSPI_IFR_WIDTH_QUAD_IO;
		break;

	case SPI_FPROTO_4_4_4:
		ifr_width = QSPI_IFR_WIDTH_QUAD_CMD;
		break;

	default:
		return -EINVAL;
	}

	*ifr = (*ifr & ~QSPI_IFR_WIDTH) | ifr_width;
	return 0;
}

static int atmel_qspi_xfer(struct udevice *dev, unsigned int bitlen,
			   const void *dout, void *din, unsigned long flags)
{
	/* This controller can only be used with SPI NOR flashes. */
	return -EINVAL;
}

static int atmel_qspi_set_speed(struct udevice *bus, uint hz)
{
	struct atmel_qspi_priv *aq = dev_get_priv(bus);
	u32 scr, scbr, mask, new_value;

	/* Compute the QSPI baudrate */
	scbr = DIV_ROUND_UP(aq->bus_clk_rate, hz);
	if (scbr > 0)
		scbr--;

	new_value = QSPI_SCR_SCBR_(scbr);
	mask = QSPI_SCR_SCBR;

	scr = qspi_readl(aq, QSPI_SCR);
	if ((scr & mask) == new_value)
		return 0;

	scr = (scr & ~mask) | new_value;
	qspi_writel(aq, QSPI_SCR, scr);

	return 0;
}

static int atmel_qspi_set_mode(struct udevice *bus, uint mode)
{
	struct atmel_qspi_priv *aq = dev_get_priv(bus);
	u32 scr, mask, new_value;

	new_value = (QSPI_SCR_CPOL_((mode & SPI_CPOL) != 0) |
		     QSPI_SCR_CPHA_((mode & SPI_CPHA) != 0));
	mask = (QSPI_SCR_CPOL | QSPI_SCR_CPHA);

	scr = qspi_readl(aq, QSPI_SCR);
	if ((scr & mask) == new_value)
		return 0;

	scr = (scr & ~mask) | new_value;
	qspi_writel(aq, QSPI_SCR, scr);

	return 0;
}

static bool
atmel_qspi_is_flash_command_supported(struct udevice *dev,
				      const struct spi_flash_command *cmd)
{
	return true;
}

static int atmel_qspi_exec_flash_command(struct udevice *dev,
					 const struct spi_flash_command *cmd)
{
	struct udevice *bus = dev_get_parent(dev);
	struct atmel_qspi_priv *aq = dev_get_priv(bus);
	unsigned int iar, icr, ifr;
	unsigned int offset;
	unsigned int imr, sr;
	unsigned long memaddr;
	int err;

	iar = 0;
	icr = 0;
	ifr = 0;

	err = atmel_qspi_set_ifr_tfrtype(cmd->flags, &ifr);
	if (err)
		return err;

	err = atmel_qpsi_set_ifr_width(cmd->proto, &ifr);
	if (err)
		return err;

	/* Compute instruction parameters */
	icr |= QSPI_ICR_INST_(cmd->inst);
	ifr |= QSPI_IFR_INSTEN;

	/* Compute address parameters. */
	switch (cmd->addr_len) {
	case 4:
		ifr |= QSPI_IFR_ADDRL_32_BIT;
		/*break;*/ /* fall through the 24bit (3 byte) address case */
	case 3:
		iar = cmd->data_len ? 0 : cmd->addr;
		ifr |= QSPI_IFR_ADDREN;
		offset = cmd->addr;
		break;
	case 0:
		offset = 0;
		break;
	default:
		return -EINVAL;
	}

	/* Compute option parameters. */
	if (cmd->num_mode_cycles) {
		unsigned int mode_cycle_bits, mode_bits;

		icr |= QSPI_ICR_OPT_(cmd->mode);
		ifr |= QSPI_IFR_OPTEN;

		switch (ifr & QSPI_IFR_WIDTH) {
		case QSPI_IFR_WIDTH_SINGLE_BIT_SPI:
		case QSPI_IFR_WIDTH_DUAL_OUTPUT:
		case QSPI_IFR_WIDTH_QUAD_OUTPUT:
			mode_cycle_bits = 1;
			break;
		case QSPI_IFR_WIDTH_DUAL_IO:
		case QSPI_IFR_WIDTH_DUAL_CMD:
			mode_cycle_bits = 2;
			break;
		case QSPI_IFR_WIDTH_QUAD_IO:
		case QSPI_IFR_WIDTH_QUAD_CMD:
			mode_cycle_bits = 4;
			break;
		default:
			return -EINVAL;
		}

		mode_bits = cmd->num_mode_cycles * mode_cycle_bits;
		switch (mode_bits) {
		case 1:
			ifr |= QSPI_IFR_OPTL_1BIT;
			break;

		case 2:
			ifr |= QSPI_IFR_OPTL_2BIT;
			break;

		case 4:
			ifr |= QSPI_IFR_OPTL_4BIT;
			break;

		case 8:
			ifr |= QSPI_IFR_OPTL_8BIT;
			break;

		default:
			return -EINVAL;
		}
	}

	/* Set the number of dummy cycles. */
	if (cmd->num_wait_states)
		ifr |= QSPI_IFR_NBDUM_(cmd->num_wait_states);

	/* Set data enable. */
	if (cmd->data_len)
		ifr |= QSPI_IFR_DATAEN;

	/* Clear pending interrupts. */
	(void)qspi_readl(aq, QSPI_SR);

	/* Set QSPI Instruction Frame registers. */
	qspi_writel(aq, QSPI_IAR, iar);
	qspi_writel(aq, QSPI_ICR, icr);
	qspi_writel(aq, QSPI_IFR, ifr);

	/* Skip to the final steps if there is no data. */
	if (!cmd->data_len)
		goto no_data;

	/* Dummy read of QSPI_IFR to synchronize APB and AHB accesses. */
	(void)qspi_readl(aq, QSPI_IFR);

	/* Stop here for Continuous Read. */
	memaddr = (unsigned long)(aq->membase + offset);
	if (cmd->tx_data)
		/* Write data. */
		atmel_qspi_memcpy_toio(memaddr, cmd->tx_data, cmd->data_len);
	else if (cmd->rx_data)
		/* Read data. */
		atmel_qspi_memcpy_fromio(cmd->rx_data, memaddr, cmd->data_len);

	/* Release the chip-select. */
	qspi_writel(aq, QSPI_CR, QSPI_CR_LASTXFER);

no_data:
	/* Poll INSTruction End and Chip Select Rise flags. */
	imr = QSPI_SR_INSTRE | QSPI_SR_CSR;
	sr = 0;
	while (sr != (QSPI_SR_INSTRE | QSPI_SR_CSR))
		sr |= qspi_readl(aq, QSPI_SR) & imr;

	return 0;
}


static const struct dm_spi_ops atmel_qspi_ops = {
	.xfer				= atmel_qspi_xfer,
	.set_speed			= atmel_qspi_set_speed,
	.set_mode			= atmel_qspi_set_mode,
	.is_flash_command_supported	= atmel_qspi_is_flash_command_supported,
	.exec_flash_command		= atmel_qspi_exec_flash_command,
};

static int atmel_qspi_enable_clk(struct udevice *bus)
{
	struct atmel_qspi_priv *aq = dev_get_priv(bus);
	struct clk clk;
	ulong clk_rate;
	int ret;

	ret = clk_get_by_index(bus, 0, &clk);
	if (ret)
		return -EINVAL;

	ret = clk_enable(&clk);
	if (ret)
		goto free_clock;

	clk_rate = clk_get_rate(&clk);
	if (!clk_rate) {
		ret = -EINVAL;
		goto free_clock;
	}

	aq->bus_clk_rate = clk_rate;

free_clock:
	clk_free(&clk);

	return ret;
}

static int atmel_qspi_probe(struct udevice *bus)
{
	const struct atmel_qspi_platdata *plat = dev_get_platdata(bus);
	struct atmel_qspi_priv *aq = dev_get_priv(bus);
	u32 mr;
	int ret;

	ret = atmel_qspi_enable_clk(bus);
	if (ret)
		return ret;

	aq->regbase = plat->regbase;
	aq->membase = plat->membase;

	/* Reset the QSPI controler */
	qspi_writel(aq, QSPI_CR, QSPI_CR_SWRST);

	/* Set the QSPI controller in Serial Memory Mode */
	mr = (QSPI_MR_NBBITS_8_BIT |
	      QSPI_MR_SMM_MEMORY |
	      QSPI_MR_CSMODE_LASTXFER);
	qspi_writel(aq, QSPI_MR, mr);

	/* Enable the QSPI controller */
	qspi_writel(aq, QSPI_CR, QSPI_CR_QSPIEN);

	return 0;
}

static int atmel_qspi_ofdata_to_platdata(struct udevice *bus)
{
	struct atmel_qspi_platdata *plat = dev_get_platdata(bus);
	const void *blob = gd->fdt_blob;
	int node = dev_of_offset (bus);
	u32 data[4];
	int ret;

	ret = fdtdec_get_int_array(blob, node, "reg", data, ARRAY_SIZE(data));
	if (ret) {
		printf("Error: Can't get base addresses (ret=%d)!\n", ret);
		return -ENODEV;
	}
	plat->regbase = (void *)data[0];
	plat->membase = (void *)data[2];

	return 0;
}

static const struct udevice_id atmel_qspi_ids[] = {
	{ .compatible = "atmel,sama5d2-qspi" },
	{ }
};

U_BOOT_DRIVER(atmel_qspi) = {
	.name		= "atmel_qspi",
	.id		= UCLASS_SPI,
	.of_match	= atmel_qspi_ids,
	.ops		= &atmel_qspi_ops,
	.ofdata_to_platdata = atmel_qspi_ofdata_to_platdata,
	.platdata_auto_alloc_size = sizeof(struct atmel_qspi_platdata),
	.priv_auto_alloc_size = sizeof(struct atmel_qspi_priv),
	.probe		= atmel_qspi_probe,
};

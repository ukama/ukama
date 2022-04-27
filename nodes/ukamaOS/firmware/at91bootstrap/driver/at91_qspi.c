/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2015, Atmel Corporation
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
#include "hardware.h"
#include "board.h"
#include "pmc.h"
#include "div.h"
#include "string.h"
#include "arch/at91_qspi.h"
#include "spi_flash/spi_nor.h"
#include "debug.h"
#include "timer.h"

#ifndef CONFIG_SYS_BASE_QSPI
#error "CONFIG_SYS_BASE_QSPI is not set"
#endif

#ifndef CONFIG_SYS_BASE_QSPI_MEM
#error "CONFIG_SYS_BASE_QSPI_MEM is not set"
#endif

struct qspi_priv {
	u32		reg_base;
	void		*mem;
};

static unsigned int qspi_readl(struct qspi_priv *qspi, u32 reg)
{
	return readl(qspi->reg_base + reg);
}

static void qspi_writel(struct qspi_priv *qspi, u32 reg, u32 value)
{
	writel(value, qspi->reg_base + reg);
}

static int qspi_set_freq(void *priv, u32 clock)
{
	struct qspi_priv *qspi = priv;
	unsigned int scbr;
	unsigned int reg;

	if (clock == 0) {
		scbr = 0;
	} else {
		scbr = div(MASTER_CLOCK + clock - 1, clock);
		if (scbr >= 1)
			scbr--;
	}

	reg = qspi_readl(qspi, QSPI_SCR);
	reg = (reg & ~QSPI_SCR_SCBR) | QSPI_SCR_SCBR_(scbr);
	qspi_writel(qspi, QSPI_SCR, reg);

	return 0;
}

static int qspi_set_mode(void *priv, u8 mode)
{
	struct qspi_priv *qspi = priv;
	unsigned int reg, val, msk;

	reg = qspi_readl(qspi, QSPI_SCR);

	msk = (QSPI_SCR_CPHA | QSPI_SCR_CPOL);
	switch (mode) {
	case SPI_MODE0:
		val = 0;
		break;

	case SPI_MODE1:
		val = QSPI_SCR_CPHA;
		break;

	case SPI_MODE2:
		val = QSPI_SCR_CPOL;
		break;

	case SPI_MODE3:
		val = (QSPI_SCR_CPOL | QSPI_SCR_CPHA);
		break;

	default:
		return -1;
	}

	if ((reg & msk) != val) {
		reg = (reg & ~msk) | val;
		qspi_writel(qspi, QSPI_SCR, reg);
	}

	return 0;
}

static int qspi_init(void *priv)
{
	struct qspi_priv *qspi = priv;

	at91_qspi_hw_init();

	qspi_writel(qspi, QSPI_CR, QSPI_CR_QSPIDIS);
	qspi_writel(qspi, QSPI_CR, QSPI_CR_SWRST);

	qspi_writel(qspi, QSPI_MR, QSPI_MR_SMM_MEMORY);
	qspi_writel(qspi, QSPI_SCR, 0);

	qspi_writel(qspi, QSPI_CR, QSPI_CR_QSPIEN);

	return 0;
}

static int qspi_cleanup(void *priv)
{
	struct qspi_priv *qspi = priv;

	qspi_writel(qspi, QSPI_CR, QSPI_CR_QSPIDIS);
	qspi_writel(qspi, QSPI_CR, QSPI_CR_SWRST);

	return 0;
}

static int qspi_init_ifr(const struct spi_flash_command *cmd,
			 unsigned int *ifr)
{
	*ifr = 0;

	switch (cmd->flags & SFLASH_TYPE_MASK) {
	case SFLASH_TYPE_READ:
		*ifr |= QSPI_IFR_TFRTYPE_READ_MEMORY;
		break;

	case SFLASH_TYPE_WRITE:
		*ifr |= QSPI_IFR_TFRTYPE_WRITE_MEMORY;
		break;

	case SFLASH_TYPE_READ_REG:
		*ifr |= QSPI_IFR_TFRTYPE_READ;
		break;

	case SFLASH_TYPE_WRITE_REG:
	case SFLASH_TYPE_ERASE:
		*ifr |= QSPI_IFR_TFRTYPE_WRITE;
		break;

	default:
		return -1;
	}

	switch (cmd->proto) {
	case SFLASH_PROTO_1_1_1:
		*ifr |= QSPI_IFR_WIDTH_SINGLE_BIT_SPI;
		break;

	case SFLASH_PROTO_1_1_2:
		*ifr |= QSPI_IFR_WIDTH_DUAL_OUTPUT;
		break;

	case SFLASH_PROTO_1_2_2:
		*ifr |= QSPI_IFR_WIDTH_DUAL_IO;
		break;

	case SFLASH_PROTO_2_2_2:
		*ifr |= QSPI_IFR_WIDTH_DUAL_CMD;
		break;

	case SFLASH_PROTO_1_1_4:
		*ifr |= QSPI_IFR_WIDTH_QUAD_OUTPUT;
		break;

	case SFLASH_PROTO_1_4_4:
		*ifr |= QSPI_IFR_WIDTH_QUAD_IO;
		break;

	case SFLASH_PROTO_4_4_4:
		*ifr |= QSPI_IFR_WIDTH_QUAD_CMD;
		break;

	default:
		return -1;
	}

	return 0;
}

static int qspi_exec(void *priv, const struct spi_flash_command *cmd)
{
	struct qspi_priv *qspi = priv;
	unsigned int iar, icr, ifr;
	unsigned int offset;
	unsigned int sr, imr;
	unsigned int timeout = 1000000;

	iar = 0;
	icr = 0;
	ifr = 0;

	/* Init ifr. */
	if (qspi_init_ifr(cmd, &ifr))
		return -1;

	/* Compute instruction parameters. */
	icr |= QSPI_ICR_INST(cmd->inst);
	ifr |= QSPI_IFR_INSTEN;

	/* Compute address parameters. */
	switch (cmd->addr_len) {
	case 4:
		ifr |= QSPI_IFR_ADDRL_32_BIT;
		//break; /* fall through the 24bit (3 byte) address case */
	case 3:
		iar = cmd->data_len ? 0 : cmd->addr;
		ifr |= QSPI_IFR_ADDREN;
		offset = cmd->addr;
		break;
	case 0:
		offset = 0;
		break;
	default:
		return -1;
	}

	/* Compute option parameters. */
	if (cmd->num_mode_cycles) {
		unsigned int mode_cycle_bits, mode_bits;

		icr |= QSPI_ICR_OPT(cmd->mode);
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
			return -1;
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
			return -1;
		}
	}

	/* Set the number of dummy cycles. */
	if (cmd->num_wait_states)
		ifr |= QSPI_IFR_NBDUM_(cmd->num_wait_states);

	/* Set data enable. */
	if (cmd->data_len) {
		ifr |= QSPI_IFR_DATAEN;

		/* Special case for Continuous Read Mode. */
		if (!cmd->tx_data && !cmd->rx_data)
			ifr |= QSPI_IFR_CRM;
	}

	/* Clear pending interrupts. */
	(void)qspi_readl(qspi, QSPI_SR);

	/* Set QSPI Instruction Frame registers. */
	qspi_writel(qspi, QSPI_IAR, iar);
#ifdef CONFIG_CPU_HAS_QSPI_RICR_WICR
	if ((cmd->flags & SFLASH_TYPE_MASK) == SFLASH_TYPE_READ ||
	    (cmd->flags & SFLASH_TYPE_MASK) == SFLASH_TYPE_READ_REG)
		qspi_writel(qspi, QSPI_RICR, icr);
	else
		qspi_writel(qspi, QSPI_WICR, icr);
#else
	qspi_writel(qspi, QSPI_ICR, icr);
#endif
	qspi_writel(qspi, QSPI_IFR, ifr);

	/* Skip to the final steps if there is no data. */
	if (!cmd->data_len)
		goto no_data;

	/* Dummy read of QSPI_IFR to synchronize APB and AHB accesses. */
	(void)qspi_readl(qspi, QSPI_IFR);

	/* Stop here for Continuous Read. */
	if (cmd->tx_data)
		/* Write data. */
		memcpy(qspi->mem + offset, cmd->tx_data, cmd->data_len);
	else if (cmd->rx_data)
		/* Read data. */
		memcpy(cmd->rx_data, qspi->mem + offset, cmd->data_len);
	else
		/* Stop here for continuous read */
		return 0;

	/* Release the chip-select. */
	qspi_writel(qspi, QSPI_CR, QSPI_CR_LASTXFER);

no_data:
	/* Poll INSTruction End and Chip Select Rise flags. */
	imr = (QSPI_SR_INSTRE | QSPI_SR_CSR);
	sr = 0;
	do {
		udelay(1);
		sr |= qspi_readl(qspi, QSPI_SR) & imr;
	} while ((!!(--timeout)) && (sr != imr));

	if (!timeout)
		dbg_info("Timeout waiting Instruction End! sr = %08x\n", sr);

	return 0;
}

static const struct spi_ops qspi_ops = {
	.init		= qspi_init,
	.cleanup	= qspi_cleanup,
	.set_freq	= qspi_set_freq,
	.set_mode	= qspi_set_mode,
	.exec		= qspi_exec,
};

int qspi_loadimage(struct image_info *image)
{
	const struct spi_flash_hwcaps hwcaps = {
		.mask = (SFLASH_HWCAPS_READ_MASK |
			 SFLASH_HWCAPS_PP_MASK),
	};
	struct spi_flash flash;
	struct qspi_priv qspi;
	int ret;

	memset(&qspi, 0, sizeof(qspi));
	qspi.reg_base = CONFIG_SYS_BASE_QSPI;
	qspi.mem = (void *)CONFIG_SYS_BASE_QSPI_MEM;

	memset(&flash, 0, sizeof(flash));
	flash.ops = &qspi_ops;
	spi_flash_set_priv(&flash, &qspi);

	/* Init the SPI controller. */
	ret = spi_flash_init(&flash);
	if (ret) {
		dbg_info("SF: Fail to initialize spi\n");
		return -1;
	}

	/* Probe the SPI flash memory. */
	ret = spi_nor_probe(&flash, &hwcaps);
	if (ret) {
		dbg_info("SF: Fail to probe SPI flash\n");
		spi_flash_cleanup(&flash);
		return -1;
	}

	return spi_flash_loadimage(&flash, image);
}

int qspi_xip(struct spi_flash *flash, void **mem)
{
	struct qspi_priv *qspi = spi_flash_get_priv(flash);
	int ret;

	if (flash->enable_0_4_4) {
		ret = flash->enable_0_4_4(flash, true);
		if (ret)
			return ret;
	}

	*mem = qspi->mem;
	return spi_flash_read(flash, 0, 1, NULL);
}

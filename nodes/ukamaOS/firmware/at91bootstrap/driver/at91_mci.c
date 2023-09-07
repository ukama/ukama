/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2012, Atmel Corporation
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
#include "common.h"
#include "hardware.h"
#include "board.h"
#include "arch/at91_mci.h"
#include "mci_media.h"
#include "div.h"
#include "debug.h"
#include "pmc.h"

#define DEFAULT_SD_BLOCK_LEN		512
#define CONFIG_SYS_DEFAULT_CLK		400000

static inline unsigned int mci_readl(unsigned int reg)
{
	return readl((void *)CONFIG_SYS_BASE_MCI + reg);
}

static inline void mci_writel(unsigned int reg, unsigned int value)
{
	writel((value), (void *)CONFIG_SYS_BASE_MCI + reg);
}

static int at91_mci_set_clock_blklen(unsigned int clock,
					unsigned int blklen)
{
	unsigned int clkdiv;
	unsigned int reg;

	clkdiv = div((at91_get_ahb_clock() + clock), clock);

	blklen &= 0xfffc;

	reg = mci_readl(MCI_BLKR);
	reg &= ~(((0x1 << 16) - 1) << 16);
	reg |= (blklen << 16);
	mci_writel(MCI_BLKR, reg);

#if defined(CONFIG_AT91SAM9X5) || defined(CONFIG_AT91SAM9N12) \
	|| defined(CONFIG_SAMA5D3X)
	clkdiv -= 2;
	reg = mci_readl(MCI_MR);
	reg &= ~((0x01 << 8) - 1);
	reg |= (clkdiv >> 1);
	reg |= (clkdiv & 1) ? AT91C_MCI_CLKODD : 0;

	mci_writel(MCI_MR, reg);
#else
	clkdiv = clkdiv / 2 - 1;
	if (clkdiv > 0xff)
		clkdiv = 0xff;

	reg = mci_readl(MCI_MR);
	reg &= ~((0x01 << 8) - 1);
	reg &= ~(((0x1 << 16) - 1) << 16);
	reg |= AT91C_MCI_MRBLKLEN(blklen);
	reg |= AT91C_MCI_CLKDIV(clkdiv);

	mci_writel(MCI_MR, reg);
#endif

	return 0;
}

static int at91_mci_init(struct sd_card *sdcard)
{
	unsigned int clock = CONFIG_SYS_DEFAULT_CLK;
	unsigned int blklen = DEFAULT_SD_BLOCK_LEN;
	int ret;

	/* software reset */
	mci_writel(MCI_CR, AT91C_MCI_SWRST);

	/* disable mci and disable power save mode */
	mci_writel(MCI_CR, AT91C_MCI_PWSDIS | AT91C_MCI_MCIDIS);

	/* enable Read Proof and Write Proof   */
	mci_writel(MCI_MR, AT91C_MCI_RDPROOF_ENABLE | AT91C_MCI_WRPROOF_ENABLE);

	/* select Slot A and set bus width 1 bit*/
	mci_writel(MCI_SDCR, AT91C_MCI_SCDSEL_SLOTA
				| AT91C_MCI_SCDBUS_1BIT);

	/* set the Data Timeout Register */
	mci_writel(MCI_DTOR, 0x7f);

	/* disable Interrupt */
	mci_writel(MCI_IDR, ~0L);

	/* set mci clock and blocklen */
	ret = at91_mci_set_clock_blklen(clock, blklen);
	if (ret)
		return ret;

	/* enable mci */
	mci_writel(MCI_CR, AT91C_MCI_MCIEN);

	return 0;
}

static int at91_mci_set_clock(struct sd_card *sdcard, unsigned int clock)
{
	unsigned int blklen;
	int ret;

	blklen = (mci_readl(MCI_BLKR) >> 16) & 0xffff;

	ret = at91_mci_set_clock_blklen(clock, blklen);
	if (ret)
		return ret;

	return 0;
}

static int at91_mci_set_bus_width(struct sd_card *sdcard, unsigned int buswidth)
{
	unsigned int reg;

	reg = mci_readl(MCI_SDCR);
	reg &= ~(AT91C_MCI_SCDBUS);

	switch (buswidth) {
	case 8:
		reg |=  AT91C_MCI_SCDBUS_8BIT;
		break;

	case 4:
		reg |= AT91C_MCI_SCDBUS_4BIT;
		break;

	case 1:
		reg |= AT91C_MCI_SCDBUS_1BIT;
		break;

	default:
		return -1;
	}

	mci_writel(MCI_SDCR, reg);

	return 0;
}

/*
 * Refer to the at91sam9g20 datasheet:
 * Figure 35-10. Read Function Flow Diagram
 */
static int at91_mci_read_data(unsigned int *data)
{
	unsigned int status;
	unsigned int error_check = (AT91C_MCI_DCRCE
					| AT91C_MCI_DTOE
					| AT91C_MCI_OVRE);
	/*
	 * Read status register MCI_SR
	 * Wait for RXRDY or error bits
	 */
	do {
		status = mci_readl(MCI_SR);
	} while ((!(status & AT91C_MCI_RXRDY))
			&& (!(status & error_check))); 

	if (status & error_check) {
		dbg_loud("Error to read data, sr: %x\n", status);
		return -1;
	}

	if (status & AT91C_MCI_RXRDY)
		*data = mci_readl(MCI_RDR);

	return 0;
}

static int at91_mci_read_block_data(unsigned int *data,
			unsigned int blocks,
			unsigned int bytes_to_read,
			unsigned int block_len)
{
	unsigned int block;
	unsigned int count;
	unsigned int words_to_read = bytes_to_read >> 2;
	unsigned int words_of_block = block_len >> 2;
	unsigned int tmp;
	int timeout = 10000;
	int ret;

	for (block = 0; block < blocks; block++) {
		for (count = 0; count < words_to_read; count++, data++) {
			ret = at91_mci_read_data((unsigned int *)data);
			if (ret)
				return ret;
		}

		/* Read the no useful data the block */
		for (; count < words_of_block; count++) {
			ret = at91_mci_read_data(&tmp);
			if (ret)
				return ret;
		}
	}

	while ((mci_readl(MCI_SR) & AT91C_MCI_DTIP) && (--timeout))
		;

	if (!timeout) {
		dbg_loud("Data Transfer in Progress.\n");
		return -1;
	}

	return 0;
}

static int at91_mci_write_data(unsigned int *data)
{
	unsigned int status;

	/*
	 * Read status register MCI_SR
	 * Wait for TXRDY or error bits
	 */
	do {
		status = mci_readl(MCI_SR);
	} while (!(status & AT91C_MCI_TXRDY));

	mci_writel(MCI_TDR, *data);

	return 0;
}

static int at91_mci_write_block_data(unsigned int *data,
			unsigned int blocks,
			unsigned int bytes_to_write,
			unsigned int block_len)
{
	unsigned int block;
	unsigned int count;
	unsigned int words_to_write = bytes_to_write >> 2;
	unsigned int words_of_block = block_len >> 2;
	unsigned int tmp = 0;
	int timeout = 10000;
	int ret;

	/* write the valid data of the block */
	for (block = 0; block < blocks; block++) {
		for (count = 0; count < words_to_write; count++, data++) {
			ret = at91_mci_write_data((unsigned int *)data);
			if (ret)
				return ret;
		}

		/* write the no useful data the block */
		for (; count < words_of_block; count++) {
			ret = at91_mci_write_data(&tmp);
			if (ret)
				return ret;
		}
	}

	while ((mci_readl(MCI_SR) & AT91C_MCI_DTIP) && (--timeout))
		;

	if (!timeout) {
		dbg_loud("Data Transfer in Progress.\n");
		return -1;
	}

	return 0;
}

static int at91_mci_send_command(struct sd_command *command, struct sd_data *data)
{
	unsigned int cmdreg;
	unsigned int error_check, status;
	unsigned int block_len = DEFAULT_SD_BLOCK_LEN;
	int ret = 0;

	error_check = AT91C_MCI_RINDE | AT91C_MCI_RDIRE | AT91C_MCI_RENDE
					| AT91C_MCI_RTOE | AT91C_MCI_DTOE;

	cmdreg = AT91C_MCI_CMDNB(command->cmd);

	cmdreg |= AT91C_MCI_MAXLAT_64;

	if ((command->resp_type == SD_RESP_TYPE_R1)
		|| (command->resp_type == SD_RESP_TYPE_R5)
		|| (command->resp_type == SD_RESP_TYPE_R6)
		|| (command->resp_type == SD_RESP_TYPE_R7)) {
		cmdreg |= AT91C_MCI_RSPTYP_48;
		error_check |= AT91C_MCI_RCRCE;
	} else if (command->resp_type == SD_RESP_TYPE_R1B) {
		cmdreg |= AT91C_MCI_RSPTYP_R1B;
		error_check |= AT91C_MCI_RCRCE;
	} else if (command->resp_type == SD_RESP_TYPE_R2) {
		cmdreg |= AT91C_MCI_RSPTYP_136;
		error_check |= AT91C_MCI_RCRCE;
	} else if ((command->resp_type == SD_RESP_TYPE_R3)
		|| (command->resp_type == SD_RESP_TYPE_R4)) {
		cmdreg |= AT91C_MCI_RSPTYP_48;
	} else {
		cmdreg |= AT91C_MCI_RSPTYP_NO;
	}

	if (data) {
		cmdreg |= AT91C_MCI_TRCMD_START;
		cmdreg |= (data->direction == SD_DATA_DIR_RD) ? AT91C_MCI_TRDIR_READ : 0;
		cmdreg |= (data->blocks > 1) ? AT91C_MCI_TRTYP_MULTIPLE : 0;
	}

	if (command->cmd == SD_CMD_READ_MULTIPLE_BLOCK) {
		if (data)
			mci_writel(MCI_BLKR, AT91C_MCI_BLKLEN(data->blocksize)
					| AT91C_MCI_BCNT(data->blocks));
	};

	/* Set the Command Argument Register */
	mci_writel(MCI_ARGR, command->argu);
	/* Set the Command Register */
	mci_writel(MCI_CMDR, cmdreg);

	/* Wait for the command ready status */
	do {
		status = mci_readl(MCI_SR);
	} while (!(status & AT91C_MCI_CMDRDY));

	/* Check error bits in the status */
	if (status & AT91C_MCI_RTOE) {
		dbg_loud("Cmd: %d Response Time-out\n", command->cmd);
		return ERROR_TIMEOUT;
	}

	if (status & error_check) {
		dbg_loud("Cmd: %d, error check: %x, status: %x\n", command->cmd, error_check, status);
		return ERROR_COMM;
	}

	if (command->resp_type == SD_RESP_TYPE_R2) {
		command->resp[0] = mci_readl(MCI_RSPR);
		command->resp[1] = mci_readl(MCI_RSPR1);
		command->resp[2] = mci_readl(MCI_RSPR2);
		command->resp[3] = mci_readl(MCI_RSPR3);
	} else {
		command->resp[0] = mci_readl(MCI_RSPR);
	}

	if (data) {
		if (data->direction == SD_DATA_DIR_RD)
			ret = at91_mci_read_block_data((unsigned int *)data->buff, data->blocks,
						data->blocksize, block_len);
		else
			ret = at91_mci_write_block_data((unsigned int *)data->buff, data->blocks,
						data->blocksize, block_len);
	}

	return ret;
}

static struct sd_host at91_mci_host;

static struct host_ops at91_mci_ops = {
	.init = at91_mci_init,
	.send_command = at91_mci_send_command,
	.set_clock = at91_mci_set_clock,
	.set_bus_width = at91_mci_set_bus_width,
};

int sdcard_register_at91_mci(struct sd_card *sdcard)
{
	sdcard->host = &at91_mci_host;
	sdcard->host->ops = &at91_mci_ops;

	sdcard->host->caps_voltages = SD_OCR_VDD_32_33 | SD_OCR_VDD_33_34;
	sdcard->host->caps_bus_width = BUS_WIDTH_1_BIT | BUS_WIDTH_4_BIT;

#ifdef CONFIG_CPU_HAS_HSMCI0
	sdcard->host->caps_high_speed = 1;
	sdcard->host->caps_bus_width |= BUS_WIDTH_8_BIT;
#endif

	return 0;
}

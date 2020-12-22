// SPDX-License-Identifier: GPL-2.0+
/*
 * SPI Flash Core
 *
 * Copyright (C) 2015 Jagan Teki <jteki@openedev.com>
 * Copyright (C) 2013 Jagannadha Sutradharudu Teki, Xilinx Inc.
 * Copyright (C) 2010 Reinhard Meyer, EMK Elektronik
 * Copyright (C) 2008 Atmel Corporation
 */

#include <common.h>
#include <errno.h>
#include <malloc.h>
#include <mapmem.h>
#include <spi.h>
#include <spi_flash.h>
#include <linux/log2.h>
#include <linux/sizes.h>
#include <dma.h>

#include "sf_internal.h"

static int read_sr(struct spi_flash *flash, u8 *rs)
{
	struct spi_flash_command cmd;
	int ret;

	spi_flash_command_init(&cmd, CMD_READ_STATUS, 0, SPI_FCMD_READ_REG);
	cmd.data_len = 1;
	cmd.rx_data = rs;
	ret = spi_flash_read_common(flash, &cmd);
	if (ret < 0) {
		debug("SF: fail to read status register\n");
		return ret;
	}

	return 0;
}

static int read_fsr(struct spi_flash *flash, u8 *fsr)
{
	struct spi_flash_command cmd;
	int ret;

	spi_flash_command_init(&cmd, CMD_FLAG_STATUS, 0, SPI_FCMD_READ_REG);
	cmd.data_len = 1;
	cmd.rx_data = fsr;
	ret = spi_flash_read_common(flash, &cmd);
	if (ret < 0) {
		debug("SF: fail to read flag status register\n");
		return ret;
	}

	return 0;
}

static int write_sr(struct spi_flash *flash, u8 ws)
{
	struct spi_flash_command cmd;
	int ret;

	spi_flash_command_init(&cmd, CMD_WRITE_STATUS, 0, SPI_FCMD_WRITE_REG);
	cmd.data_len = 1;
	cmd.tx_data = &ws;
	ret = spi_flash_write_common(flash, &cmd);
	if (ret < 0) {
		debug("SF: fail to write status register\n");
		return ret;
	}

	return 0;
}

#if defined(CONFIG_SPI_FLASH_SPANSION) || defined(CONFIG_SPI_FLASH_WINBOND) || \
    defined(CONFIG_SPI_FLASH_SST)
static int read_cr(struct spi_flash *flash, u8 *rc)
{
	struct spi_flash_command cmd;
	int ret;

	spi_flash_command_init(&cmd, CMD_READ_CONFIG, 0, SPI_FCMD_READ_REG);
	cmd.data_len = 1;
	cmd.rx_data = rc;
	ret = spi_flash_read_common(flash, &cmd);
	if (ret < 0) {
		debug("SF: fail to read config register\n");
		return ret;
	}

	return 0;
}

static int write_cr(struct spi_flash *flash, u8 wc)
{
	struct spi_flash_command cmd;
	u8 data[2];
	int ret;

	ret = read_sr(flash, &data[0]);
	if (ret < 0)
		return ret;

	spi_flash_command_init(&cmd, CMD_WRITE_STATUS, 0, SPI_FCMD_WRITE_REG);
	cmd.data_len = 2;
	cmd.tx_data = data;
	data[1] = wc;
	ret = spi_flash_write_common(flash, &cmd);
	if (ret) {
		debug("SF: fail to write config register\n");
		return ret;
	}

	return 0;
}
#endif

#ifdef CONFIG_SPI_FLASH_BAR
/*
 * This "clean_bar" is necessary in a situation when one was accessing
 * spi flash memory > 16 MiB by using Bank Address Register's BA24 bit.
 *
 * After it the BA24 bit shall be cleared to allow access to correct
 * memory region after SW reset (by calling "reset" command).
 *
 * Otherwise, the BA24 bit may be left set and then after reset, the
 * ROM would read/write/erase SPL from 16 MiB * bank_sel address.
 */
static int clean_bar(struct spi_flash *flash)
{
	u8 cmd, bank_sel = 0;

	if (flash->bank_curr == 0)
		return 0;
	cmd = flash->bank_write_cmd;
	flash->bank_curr = 0;

	return spi_flash_write_common(flash, &cmd, 1, &bank_sel, 1);
}

static int write_bar(struct spi_flash *flash, u32 offset)
{
	struct spi_flash_command cmd;
	u8 bank_sel;
	int ret;

	bank_sel = offset / (SPI_FLASH_16MB_BOUN << flash->shift);
	if (bank_sel == flash->bank_curr)
		goto bar_end;

	spi_flash_command_init(&cmd, flash->bank_write_cmd, 0,
			       SPI_FCMD_WRITE_REG);
	cmd.data_len = 1;
	cmd.tx_data = &bank_sel;
	ret = spi_flash_write_common(flash, &cmd);
	if (ret < 0) {
		debug("SF: fail to write bank register\n");
		return ret;
	}

bar_end:
	flash->bank_curr = bank_sel;
	return flash->bank_curr;
}

static int read_bar(struct spi_flash *flash, const struct spi_flash_info *info)
{
	struct spi_flash_command cmd;
	u8 curr_bank = 0;
	int ret;

	if (flash->size <= SPI_FLASH_16MB_BOUN)
		goto bar_end;

	switch (JEDEC_MFR(info)) {
	case SPI_FLASH_CFI_MFR_SPANSION:
		flash->bank_read_cmd = CMD_BANKADDR_BRRD;
		flash->bank_write_cmd = CMD_BANKADDR_BRWR;
		break;
	default:
		flash->bank_read_cmd = CMD_EXTNADDR_RDEAR;
		flash->bank_write_cmd = CMD_EXTNADDR_WREAR;
	}

	spi_flash_command_init(&cmd, flash->bank_read_cmd, 0,
			       SPI_FCMD_READ_REG);
	cmd.data_len = 1;
	cmd.rx_data = &curr_bank;
	ret = spi_flash_read_common(flash, &cmd);
	if (ret) {
		debug("SF: fail to read bank addr register\n");
		return ret;
	}

bar_end:
	flash->bank_curr = curr_bank;
	return 0;
}
#endif

#ifdef CONFIG_SPI_FLASH_4BAIS
static u8 spi_flash_convert_opcode(u8 opcode, const u8 table[][2], size_t size)
{
	size_t i;

	for (i = 0; i < size; i++)
		if (table[i][0] == opcode)
			return table[i][1];

	/* No conversion found, keep input op code. */
	return opcode;
}

static u8 spi_flash_convert_3to4_read(u8 opcode)
{
	static const u8 spi_flash_3to4_read[][2] = {
		{CMD_READ_ARRAY_SLOW,		CMD_READ_ARRAY_SLOW_4B},
		{CMD_READ_ARRAY_FAST,		CMD_READ_ARRAY_FAST_4B},
		{CMD_READ_DUAL_OUTPUT_FAST,	CMD_READ_DUAL_OUTPUT_FAST_4B},
		{CMD_READ_DUAL_IO_FAST,		CMD_READ_DUAL_IO_FAST_4B},
		{CMD_READ_QUAD_OUTPUT_FAST,	CMD_READ_QUAD_OUTPUT_FAST_4B},
		{CMD_READ_QUAD_IO_FAST,		CMD_READ_QUAD_IO_FAST_4B},
	};

	return spi_flash_convert_opcode(opcode, spi_flash_3to4_read,
					ARRAY_SIZE(spi_flash_3to4_read));
}

static u8 spi_flash_convert_3to4_write(u8 opcode)
{
	static const u8 spi_flash_3to4_write[][2] = {
		{CMD_PAGE_PROGRAM,		CMD_PAGE_PROGRAM_4B},
		{CMD_PAGE_PROGRAM_1_1_4,	CMD_PAGE_PROGRAM_1_1_4_4B},
		{CMD_PAGE_PROGRAM_1_4_4,	CMD_PAGE_PROGRAM_1_4_4_4B},
	};

	return spi_flash_convert_opcode(opcode, spi_flash_3to4_write,
					ARRAY_SIZE(spi_flash_3to4_write));
}

static u8 spi_flash_convert_3to4_erase(u8 opcode)
{
	static const u8 spi_flash_3to4_erase[][2] = {
		{CMD_ERASE_4K,	CMD_ERASE_4K_4B},
		{CMD_ERASE_64K,	CMD_ERASE_64K_4B},
	};

	return spi_flash_convert_opcode(opcode, spi_flash_3to4_erase,
					ARRAY_SIZE(spi_flash_3to4_erase));
}

static void spi_flash_set_4byte_addr_opcodes(struct spi_flash *flash,
					     const struct spi_flash_info *info)
{
	flash->read_cmd = spi_flash_convert_3to4_read(flash->read_cmd);
	flash->write_cmd = spi_flash_convert_3to4_write(flash->write_cmd);
	flash->erase_cmd = spi_flash_convert_3to4_erase(flash->erase_cmd);
	flash->addr_len = SPI_FLASH_4B_ADDR_LEN;
}
#endif

#ifdef CONFIG_SF_DUAL_FLASH
static void spi_flash_dual(struct spi_flash *flash, u32 *addr)
{
	switch (flash->dual_flash) {
	case SF_DUAL_STACKED_FLASH:
		if (*addr >= (flash->size >> 1)) {
			*addr -= flash->size >> 1;
			flash->flags |= SNOR_F_USE_UPAGE;
		} else {
			flash->flags &= ~SNOR_F_USE_UPAGE;
		}
		break;
	case SF_DUAL_PARALLEL_FLASH:
		*addr >>= flash->shift;
		break;
	default:
		debug("SF: Unsupported dual_flash=%d\n", flash->dual_flash);
		break;
	}
}
#endif

static int spi_flash_sr_ready(struct spi_flash *flash)
{
	u8 sr;
	int ret;

	ret = read_sr(flash, &sr);
	if (ret < 0)
		return ret;

	return !(sr & STATUS_WIP);
}

static int spi_flash_fsr_ready(struct spi_flash *flash)
{
	u8 fsr;
	int ret;

	ret = read_fsr(flash, &fsr);
	if (ret < 0)
		return ret;

	return fsr & STATUS_PEC;
}

static int spi_flash_ready(struct spi_flash *flash)
{
	int sr, fsr;

	sr = spi_flash_sr_ready(flash);
	if (sr < 0)
		return sr;

	fsr = 1;
	if (flash->flags & SNOR_F_USE_FSR) {
		fsr = spi_flash_fsr_ready(flash);
		if (fsr < 0)
			return fsr;
	}

	return sr && fsr;
}

static int spi_flash_wait_till_ready(struct spi_flash *flash,
				     unsigned long timeout)
{
	unsigned long timebase;
	int ret;

	timebase = get_timer(0);

	while (get_timer(timebase) < timeout) {
		ret = spi_flash_ready(flash);
		if (ret < 0)
			return ret;
		if (ret)
			return 0;
	}

	printf("SF: Timeout!\n");

	return -ETIMEDOUT;
}

int spi_flash_write_common(struct spi_flash *flash,
			   const struct spi_flash_command *cmd)
{
	struct spi_slave *spi = flash->spi;
	unsigned long timeout = SPI_FLASH_PROG_TIMEOUT;
	int ret;

	if (cmd->tx_data == NULL)
		timeout = SPI_FLASH_PAGE_ERASE_TIMEOUT;

	ret = spi_claim_bus(spi);
	if (ret) {
		debug("SF: unable to claim SPI bus\n");
		return ret;
	}

	ret = spi_flash_cmd_write_enable(flash);
	if (ret < 0) {
		debug("SF: enabling write failed\n");
		return ret;
	}

	ret = spi_flash_cmd_write(flash, cmd);
	if (ret < 0) {
		debug("SF: write cmd failed\n");
		return ret;
	}

	ret = spi_flash_wait_till_ready(flash, timeout);
	if (ret < 0) {
		debug("SF: write %s timed out\n",
		      timeout == SPI_FLASH_PROG_TIMEOUT ?
			"program" : "page erase");
		return ret;
	}

	spi_release_bus(spi);

	return ret;
}

int spi_flash_cmd_erase_ops(struct spi_flash *flash, u32 offset, size_t len)
{
	struct spi_flash_command cmd;
	u32 erase_size, erase_addr;
	int ret = -1;

	erase_size = flash->erase_size;
	if (offset % erase_size || len % erase_size) {
		printf("SF: Erase offset/length not multiple of erase size\n");
		return -1;
	}

	if (flash->flash_is_locked) {
		if (flash->flash_is_locked(flash, offset, len) > 0) {
			printf("offset 0x%x is protected and cannot be erased\n",
			       offset);
			return -EINVAL;
		}
	}

	spi_flash_command_init(&cmd, flash->erase_cmd, flash->addr_len,
			       SPI_FCMD_ERASE);
	while (len) {
		erase_addr = offset;

#ifdef CONFIG_SF_DUAL_FLASH
		if (flash->dual_flash > SF_SINGLE_FLASH)
			spi_flash_dual(flash, &erase_addr);
#endif
#ifdef CONFIG_SPI_FLASH_BAR
		ret = write_bar(flash, erase_addr);
		if (ret < 0)
			return ret;
#endif
		cmd.addr = erase_addr;

		debug("SF: erase %2x %2x %2x %2x (%x)\n", cmd.inst,
		      (cmd.addr >> 16) & 0xff,
		      (cmd.addr >>  8) & 0xff,
		      (cmd.addr >>  0) & 0xff,
		      erase_addr);

		ret = spi_flash_write_common(flash, &cmd);
		if (ret < 0) {
			debug("SF: erase failed\n");
			break;
		}

		offset += erase_size;
		len -= erase_size;
	}

#ifdef CONFIG_SPI_FLASH_BAR
	ret = clean_bar(flash);
#endif

	return ret;
}

int spi_flash_cmd_write_ops(struct spi_flash *flash, u32 offset,
		size_t len, const void *buf)
{
	struct spi_slave *spi = flash->spi;
	struct spi_flash_command cmd;
	unsigned long byte_addr, page_size;
	u32 write_addr;
	size_t chunk_len, actual;
	int ret = -1;

	page_size = flash->page_size;

	if (flash->flash_is_locked) {
		if (flash->flash_is_locked(flash, offset, len) > 0) {
			printf("offset 0x%x is protected and cannot be written\n",
			       offset);
			return -EINVAL;
		}
	}

	spi_flash_command_init(&cmd, flash->write_cmd, flash->addr_len,
			       SPI_FCMD_WRITE);
	cmd.proto = flash->write_proto;
	for (actual = 0; actual < len; actual += chunk_len) {
		write_addr = offset;

#ifdef CONFIG_SF_DUAL_FLASH
		if (flash->dual_flash > SF_SINGLE_FLASH)
			spi_flash_dual(flash, &write_addr);
#endif
#ifdef CONFIG_SPI_FLASH_BAR
		ret = write_bar(flash, write_addr);
		if (ret < 0)
			return ret;
#endif
		byte_addr = offset % page_size;
		chunk_len = min(len - actual, (size_t)(page_size - byte_addr));

		if (spi->max_write_size)
			chunk_len = min(chunk_len,
					spi->max_write_size - sizeof(cmd));

		cmd.addr = write_addr;
		cmd.data_len = chunk_len;
		cmd.tx_data = buf + actual;

		debug("SF: 0x%p => cmd = { 0x%02x 0x%02x%02x%02x } chunk_len = %zu\n",
		      buf + actual, cmd.inst,
		      (cmd.addr >> 16) & 0xff,
		      (cmd.addr >>  8) & 0xff,
		      (cmd.addr >>  0) & 0xff,
		      chunk_len);

		ret = spi_flash_write_common(flash, &cmd);
		if (ret < 0) {
			debug("SF: write failed\n");
			break;
		}

		offset += chunk_len;
	}

#ifdef CONFIG_SPI_FLASH_BAR
	ret = clean_bar(flash);
#endif

	return ret;
}

int spi_flash_read_common(struct spi_flash *flash,
			  const struct spi_flash_command *cmd)
{
	struct spi_slave *spi = flash->spi;
	int ret;

	ret = spi_claim_bus(spi);
	if (ret) {
		debug("SF: unable to claim SPI bus\n");
		return ret;
	}

	ret = spi_flash_cmd_read(flash, cmd);
	if (ret < 0) {
		debug("SF: read cmd failed\n");
		return ret;
	}

	spi_release_bus(spi);

	return ret;
}

/*
 * TODO: remove the weak after all the other spi_flash_copy_mmap
 * implementations removed from drivers
 */
void __weak spi_flash_copy_mmap(void *data, void *offset, size_t len)
{
#ifdef CONFIG_DMA
	if (!dma_memcpy(data, offset, len))
		return;
#endif
	memcpy(data, offset, len);
}

int spi_flash_cmd_read_ops(struct spi_flash *flash, u32 offset,
		size_t len, void *data)
{
	struct spi_slave *spi = flash->spi;
	struct spi_flash_command cmd;
	u32 remain_len, read_len, read_addr;
	int bank_sel = 0;
	int ret = -1;

	/* Handle memory-mapped SPI */
	if (flash->memory_map) {
		ret = spi_claim_bus(spi);
		if (ret) {
			debug("SF: unable to claim SPI bus\n");
			return ret;
		}
		spi_xfer(spi, 0, NULL, NULL, SPI_XFER_MMAP);
		spi_flash_copy_mmap(data, flash->memory_map + offset, len);
		spi_xfer(spi, 0, NULL, NULL, SPI_XFER_MMAP_END);
		spi_release_bus(spi);
		return 0;
	}

	spi_flash_command_init(&cmd, flash->read_cmd, flash->addr_len,
			       SPI_FCMD_READ);
	cmd.proto = flash->read_proto;
	cmd.num_wait_states = flash->dummy_byte * 8;
	while (len) {
		read_addr = offset;

#ifdef CONFIG_SF_DUAL_FLASH
		if (flash->dual_flash > SF_SINGLE_FLASH)
			spi_flash_dual(flash, &read_addr);
#endif
#ifdef CONFIG_SPI_FLASH_BAR
		ret = write_bar(flash, read_addr);
		if (ret < 0)
			return ret;
		bank_sel = flash->bank_curr;
#endif
		remain_len = ((SPI_FLASH_16MB_BOUN << flash->shift) *
				(bank_sel + 1)) - offset;
		if (len < remain_len)
			read_len = len;
		else
			read_len = remain_len;

		if (spi->max_read_size)
			read_len = min(read_len, spi->max_read_size);

		cmd.addr = read_addr;
		cmd.data_len = read_len;
		cmd.rx_data = data;
		ret = spi_flash_read_common(flash, &cmd);
		if (ret < 0) {
			debug("SF: read failed\n");
			break;
		}

		offset += read_len;
		len -= read_len;
		data += read_len;
	}

#ifdef CONFIG_SPI_FLASH_BAR
	ret = clean_bar(flash);
#endif
	return ret;
}

#ifdef CONFIG_SPI_FLASH_SST
static bool sst26_process_bpr(u32 bpr_size, u8 *cmd, u32 bit, enum lock_ctl ctl)
{
	switch (ctl) {
		case SST26_CTL_LOCK:
			cmd[bpr_size - (bit / 8) - 1] |= BIT(bit % 8);
			break;
		case SST26_CTL_UNLOCK:
			cmd[bpr_size - (bit / 8) - 1] &= ~BIT(bit % 8);
			break;
		case SST26_CTL_CHECK:
			return !!(cmd[bpr_size - (bit / 8) - 1] & BIT(bit % 8));
	}

	return false;
}

/*
 * sst26wf016/sst26wf032/sst26wf064 have next block protection:
 * 4x   - 8  KByte blocks - read & write protection bits - upper addresses
 * 1x   - 32 KByte blocks - write protection bits
 * rest - 64 KByte blocks - write protection bits
 * 1x   - 32 KByte blocks - write protection bits
 * 4x   - 8  KByte blocks - read & write protection bits - lower addresses
 *
 * We'll support only per 64k lock/unlock so lower and upper 64 KByte region
 * will be treated as single block.
 */

/*
 * Lock, unlock or check lock status of the flash region of the flash (depending
 * on the lock_ctl value)
 */
static int sst26_lock_ctl(struct spi_flash *flash, u32 ofs, size_t len, enum lock_ctl ctl)
{
	u32 i, bpr_ptr, rptr_64k, lptr_64k, bpr_size;
	bool lower_64k = false, upper_64k = false;
	u8 bpr_buff[SST26_MAX_BPR_REG_LEN] = {};
	struct spi_flash_command cmd;
	int ret;

	/* Check length and offset for 64k alignment */
	if ((ofs & (SZ_64K - 1)) || (len & (SZ_64K - 1)))
		return -EINVAL;

	if (ofs + len > flash->size)
		return -EINVAL;

	/* SST26 family has only 16 Mbit, 32 Mbit and 64 Mbit IC */
	if (flash->size != SZ_2M &&
	    flash->size != SZ_4M &&
	    flash->size != SZ_8M)
		return -EINVAL;

	bpr_size = 2 + (flash->size / SZ_64K / 8);

        spi_flash_command_init(&cmd, SST26_CMD_READ_BPR, 0, SPI_FCMD_READ_REG);
        cmd.data_len = bpr_size;
        cmd.rx_data = bpr_buff;
        ret = spi_flash_read_common(flash, &cmd);

	if (ret < 0) {
		printf("SF: fail to read block-protection register\n");
		return ret;
	}

	rptr_64k = min_t(u32, ofs + len , flash->size - SST26_BOUND_REG_SIZE);
	lptr_64k = max_t(u32, ofs, SST26_BOUND_REG_SIZE);

	upper_64k = ((ofs + len) > (flash->size - SST26_BOUND_REG_SIZE));
	lower_64k = (ofs < SST26_BOUND_REG_SIZE);

	/* Lower bits in block-protection register are about 64k region */
	bpr_ptr = lptr_64k / SZ_64K - 1;

	/* Process 64K blocks region */
	while (lptr_64k < rptr_64k) {
		if (sst26_process_bpr(bpr_size, bpr_buff, bpr_ptr, ctl))
			return EACCES;

		bpr_ptr++;
		lptr_64k += SZ_64K;
	}

	/* 32K and 8K region bits in BPR are after 64k region bits */
	bpr_ptr = (flash->size - 2 * SST26_BOUND_REG_SIZE) / SZ_64K;

	/* Process lower 32K block region */
	if (lower_64k)
		if (sst26_process_bpr(bpr_size, bpr_buff, bpr_ptr, ctl))
			return EACCES;

	bpr_ptr++;

	/* Process upper 32K block region */
	if (upper_64k)
		if (sst26_process_bpr(bpr_size, bpr_buff, bpr_ptr, ctl))
			return EACCES;

	bpr_ptr++;

	/* Process lower 8K block regions */
	for (i = 0; i < SST26_BPR_8K_NUM; i++) {
		if (lower_64k)
			if (sst26_process_bpr(bpr_size, bpr_buff, bpr_ptr, ctl))
				return EACCES;

		/* In 8K area BPR has both read and write protection bits */
		bpr_ptr += 2;
	}

	/* Process upper 8K block regions */
	for (i = 0; i < SST26_BPR_8K_NUM; i++) {
		if (upper_64k)
			if (sst26_process_bpr(bpr_size, bpr_buff, bpr_ptr, ctl))
				return EACCES;

		/* In 8K area BPR has both read and write protection bits */
		bpr_ptr += 2;
	}

	/* If we check region status we don't need to write BPR back */
	if (ctl == SST26_CTL_CHECK)
		return 0;

	spi_flash_command_init(&cmd, SST26_CMD_WRITE_BPR, 0, SPI_FCMD_WRITE_REG);
	cmd.data_len = bpr_size;
        cmd.tx_data = bpr_buff;
	ret = spi_flash_write_common(flash, &cmd);

	if (ret < 0) {
		printf("SF: fail to write block-protection register\n");
		return ret;
	}

	return 0;
}

static int sst26_unlock(struct spi_flash *flash, u32 ofs, size_t len)
{
	return sst26_lock_ctl(flash, ofs, len, SST26_CTL_UNLOCK);
}

static int sst26_lock(struct spi_flash *flash, u32 ofs, size_t len)
{
	return sst26_lock_ctl(flash, ofs, len, SST26_CTL_LOCK);
}

/*
 * Returns EACCES (positive value) if region is locked, 0 if region is unlocked,
 * and negative on errors.
 */
static int sst26_is_locked(struct spi_flash *flash, u32 ofs, size_t len)
{
	/*
	 * is_locked function is used for check before reading or erasing flash
	 * region, so offset and length might be not 64k allighned, so adjust
	 * them to be 64k allighned as sst26_lock_ctl works only with 64k
	 * allighned regions.
	 */
	ofs -= ofs & (SZ_64K - 1);
	len = len & (SZ_64K - 1) ? (len & ~(SZ_64K - 1)) + SZ_64K : len;

	return sst26_lock_ctl(flash, ofs, len, SST26_CTL_CHECK);
}

static int sst_byte_write(struct spi_flash *flash, u32 offset, const void *buf)
{
	struct spi_flash_command cmd;
	int ret;
	u8 sr = 0xFFu;

	spi_flash_command_init(&cmd, CMD_SST_BP, flash->addr_len,
			       SPI_FCMD_WRITE);
	cmd.addr = offset;
	cmd.data_len = 1;
	cmd.tx_data = buf;

#ifdef DEBUG
	read_sr(flash, &sr);
#endif
	debug("BP[%02x]: 0x%p => cmd = { 0x%02x 0x%06x }\n",
	      sr, buf, cmd.inst, offset);

	ret = spi_flash_cmd_write_enable(flash);
	if (ret)
		return ret;

	ret = spi_flash_cmd_write(flash, &cmd);
	if (ret)
		return ret;

	return spi_flash_wait_till_ready(flash, SPI_FLASH_PROG_TIMEOUT);
}

int sst_write_wp(struct spi_flash *flash, u32 offset, size_t len,
		const void *buf)
{
	struct spi_slave *spi = flash->spi;
	struct spi_flash_command cmd;
	size_t actual;
	int ret;
	u8 sr = 0xFFu;

	ret = spi_claim_bus(spi);
	if (ret) {
		debug("SF: Unable to claim SPI bus\n");
		return ret;
	}

	/* If the data is not word aligned, write out leading single byte */
	actual = offset % 2;
	if (actual) {
		ret = sst_byte_write(flash, offset, buf);
		if (ret)
			goto done;
	}
	offset += actual;

	ret = spi_flash_cmd_write_enable(flash);
	if (ret)
		goto done;

	spi_flash_command_init(&cmd, CMD_SST_AAI_WP, flash->addr_len,
			       SPI_FCMD_WRITE);
	cmd.addr = offset;
	cmd.data_len = 2;

	for (; actual < len - 1; actual += 2) {
#ifdef DEBUG
		read_sr(flash, &sr);
#endif
		debug("WP[%02x]: 0x%p => cmd = { 0x%02x 0x%06x }\n",
		      sr, buf + actual, cmd.inst, offset);

		cmd.tx_data = buf + actual;
		ret = spi_flash_cmd_write(flash, &cmd);
		if (ret) {
			debug("SF: sst word program failed\n");
			break;
		}

		ret = spi_flash_wait_till_ready(flash, SPI_FLASH_PROG_TIMEOUT);
		if (ret)
			break;

		cmd.addr_len = 0;
		offset += 2;
	}

	if (!ret)
		ret = spi_flash_cmd_write_disable(flash);

	/* If there is a single trailing byte, write it out */
	if (!ret && actual != len)
		ret = sst_byte_write(flash, offset, buf + actual);

 done:
	debug("SF: sst: program %s %zu bytes @ 0x%zx\n",
	      ret ? "failure" : "success", len, offset - actual);

	spi_release_bus(spi);
	return ret;
}

int sst_write_bp(struct spi_flash *flash, u32 offset, size_t len,
		const void *buf)
{
	struct spi_slave *spi = flash->spi;
	size_t actual;
	int ret;

	ret = spi_claim_bus(spi);
	if (ret) {
		debug("SF: Unable to claim SPI bus\n");
		return ret;
	}

	for (actual = 0; actual < len; actual++) {
		ret = sst_byte_write(flash, offset, buf + actual);
		if (ret) {
			debug("SF: sst byte program failed\n");
			break;
		}
		offset++;
	}

	if (!ret)
		ret = spi_flash_cmd_write_disable(flash);

	debug("SF: sst: program %s %zu bytes @ 0x%zx\n",
	      ret ? "failure" : "success", len, offset - actual);

	spi_release_bus(spi);
	return ret;
}
#endif

#if defined(CONFIG_SPI_FLASH_STMICRO) || defined(CONFIG_SPI_FLASH_SST)
static void stm_get_locked_range(struct spi_flash *flash, u8 sr, loff_t *ofs,
				 u64 *len)
{
	u8 mask = SR_BP2 | SR_BP1 | SR_BP0;
	int shift = ffs(mask) - 1;
	int pow;

	if (!(sr & mask)) {
		/* No protection */
		*ofs = 0;
		*len = 0;
	} else {
		pow = ((sr & mask) ^ mask) >> shift;
		*len = flash->size >> pow;
		*ofs = flash->size - *len;
	}
}

/*
 * Return 1 if the entire region is locked, 0 otherwise
 */
static int stm_is_locked_sr(struct spi_flash *flash, loff_t ofs, u64 len,
			    u8 sr)
{
	loff_t lock_offs;
	u64 lock_len;

	stm_get_locked_range(flash, sr, &lock_offs, &lock_len);

	return (ofs + len <= lock_offs + lock_len) && (ofs >= lock_offs);
}

/*
 * Check if a region of the flash is (completely) locked. See stm_lock() for
 * more info.
 *
 * Returns 1 if entire region is locked, 0 if any portion is unlocked, and
 * negative on errors.
 */
int stm_is_locked(struct spi_flash *flash, u32 ofs, size_t len)
{
	int status;
	u8 sr;

	status = read_sr(flash, &sr);
	if (status < 0)
		return status;

	return stm_is_locked_sr(flash, ofs, len, sr);
}

/*
 * Lock a region of the flash. Compatible with ST Micro and similar flash.
 * Supports only the block protection bits BP{0,1,2} in the status register
 * (SR). Does not support these features found in newer SR bitfields:
 *   - TB: top/bottom protect - only handle TB=0 (top protect)
 *   - SEC: sector/block protect - only handle SEC=0 (block protect)
 *   - CMP: complement protect - only support CMP=0 (range is not complemented)
 *
 * Sample table portion for 8MB flash (Winbond w25q64fw):
 *
 *   SEC  |  TB   |  BP2  |  BP1  |  BP0  |  Prot Length  | Protected Portion
 *  --------------------------------------------------------------------------
 *    X   |   X   |   0   |   0   |   0   |  NONE         | NONE
 *    0   |   0   |   0   |   0   |   1   |  128 KB       | Upper 1/64
 *    0   |   0   |   0   |   1   |   0   |  256 KB       | Upper 1/32
 *    0   |   0   |   0   |   1   |   1   |  512 KB       | Upper 1/16
 *    0   |   0   |   1   |   0   |   0   |  1 MB         | Upper 1/8
 *    0   |   0   |   1   |   0   |   1   |  2 MB         | Upper 1/4
 *    0   |   0   |   1   |   1   |   0   |  4 MB         | Upper 1/2
 *    X   |   X   |   1   |   1   |   1   |  8 MB         | ALL
 *
 * Returns negative on errors, 0 on success.
 */
int stm_lock(struct spi_flash *flash, u32 ofs, size_t len)
{
	u8 status_old, status_new;
	u8 mask = SR_BP2 | SR_BP1 | SR_BP0;
	u8 shift = ffs(mask) - 1, pow, val;
	int ret;

	ret = read_sr(flash, &status_old);
	if (ret < 0)
		return ret;

	/* SPI NOR always locks to the end */
	if (ofs + len != flash->size) {
		/* Does combined region extend to end? */
		if (!stm_is_locked_sr(flash, ofs + len, flash->size - ofs - len,
				      status_old))
			return -EINVAL;
		len = flash->size - ofs;
	}

	/*
	 * Need smallest pow such that:
	 *
	 *   1 / (2^pow) <= (len / size)
	 *
	 * so (assuming power-of-2 size) we do:
	 *
	 *   pow = ceil(log2(size / len)) = log2(size) - floor(log2(len))
	 */
	pow = ilog2(flash->size) - ilog2(len);
	val = mask - (pow << shift);
	if (val & ~mask)
		return -EINVAL;

	/* Don't "lock" with no region! */
	if (!(val & mask))
		return -EINVAL;

	status_new = (status_old & ~mask) | val;

	/* Only modify protection if it will not unlock other areas */
	if ((status_new & mask) <= (status_old & mask))
		return -EINVAL;

	write_sr(flash, status_new);

	return 0;
}

/*
 * Unlock a region of the flash. See stm_lock() for more info
 *
 * Returns negative on errors, 0 on success.
 */
int stm_unlock(struct spi_flash *flash, u32 ofs, size_t len)
{
	uint8_t status_old, status_new;
	u8 mask = SR_BP2 | SR_BP1 | SR_BP0;
	u8 shift = ffs(mask) - 1, pow, val;
	int ret;

	ret = read_sr(flash, &status_old);
	if (ret < 0)
		return ret;

	/* Cannot unlock; would unlock larger region than requested */
	if (stm_is_locked_sr(flash, ofs - flash->erase_size, flash->erase_size,
			     status_old))
		return -EINVAL;
	/*
	 * Need largest pow such that:
	 *
	 *   1 / (2^pow) >= (len / size)
	 *
	 * so (assuming power-of-2 size) we do:
	 *
	 *   pow = floor(log2(size / len)) = log2(size) - ceil(log2(len))
	 */
	pow = ilog2(flash->size) - order_base_2(flash->size - (ofs + len));
	if (ofs + len == flash->size) {
		val = 0; /* fully unlocked */
	} else {
		val = mask - (pow << shift);
		/* Some power-of-two sizes are not supported */
		if (val & ~mask)
			return -EINVAL;
	}

	status_new = (status_old & ~mask) | val;

	/* Only modify protection if it will not lock other areas */
	if ((status_new & mask) >= (status_old & mask))
		return -EINVAL;

	write_sr(flash, status_new);

	return 0;
}
#endif

#ifdef CONFIG_SPI_FLASH_MACRONIX
static int macronix_quad_enable(struct spi_flash *flash)
{
	u8 qeb_status;
	int ret;

	ret = read_sr(flash, &qeb_status);
	if (ret < 0)
		return ret;

	if (qeb_status & STATUS_QEB_MXIC)
		return 0;

	ret = write_sr(flash, qeb_status | STATUS_QEB_MXIC);
	if (ret < 0)
		return ret;

	/* read SR and check it */
	ret = read_sr(flash, &qeb_status);
	if (!(ret >= 0 && (qeb_status & STATUS_QEB_MXIC))) {
		printf("SF: Macronix SR Quad bit not clear\n");
		return -EINVAL;
	}

	return ret;
}
#endif

#if defined(CONFIG_SPI_FLASH_SPANSION) || defined(CONFIG_SPI_FLASH_WINBOND) ||\
    defined(CONFIG_SPI_FLASH_SST)
static int spansion_quad_enable(struct spi_flash *flash)
{
	u8 qeb_status;
	int ret;

	ret = read_cr(flash, &qeb_status);
	if (ret < 0)
		return ret;

	if (qeb_status & STATUS_QEB_WINSPAN)
		return 0;

	ret = write_cr(flash, qeb_status | STATUS_QEB_WINSPAN);
	if (ret < 0)
		return ret;

	/* read CR and check it */
	ret = read_cr(flash, &qeb_status);
	if (!(ret >= 0 && (qeb_status & STATUS_QEB_WINSPAN))) {
		printf("SF: Spansion CR Quad bit not clear\n");
		return -EINVAL;
	}

	return ret;
}
#endif

static const struct spi_flash_info *spi_flash_read_id(struct spi_flash *flash)
{
	int				tmp;
	u8				id[SPI_FLASH_MAX_ID_LEN];
	const struct spi_flash_info	*info;

	tmp = spi_flash_cmd(flash, CMD_READ_ID, id, SPI_FLASH_MAX_ID_LEN);
	if (tmp < 0) {
		printf("SF: error %d reading JEDEC ID\n", tmp);
		return ERR_PTR(tmp);
	}

	info = spi_flash_ids;
	for (; info->name != NULL; info++) {
		if (info->id_len) {
			if (!memcmp(info->id, id, info->id_len))
				return info;
		}
	}

	printf("SF: unrecognized JEDEC id bytes: %02x, %02x, %02x\n",
	       id[0], id[1], id[2]);
	return ERR_PTR(-ENODEV);
}

static int set_quad_mode(struct spi_flash *flash,
			 const struct spi_flash_info *info)
{
	switch (JEDEC_MFR(info)) {
#ifdef CONFIG_SPI_FLASH_MACRONIX
	case SPI_FLASH_CFI_MFR_MACRONIX:
		return macronix_quad_enable(flash);
#endif
#if defined(CONFIG_SPI_FLASH_SPANSION) || defined(CONFIG_SPI_FLASH_WINBOND) ||\
    defined(CONFIG_SPI_FLASH_SST)
	case SPI_FLASH_CFI_MFR_SPANSION:
	case SPI_FLASH_CFI_MFR_WINBOND:
	case SPI_FLASH_CFI_MFR_SST:
		return spansion_quad_enable(flash);
#endif
#ifdef CONFIG_SPI_FLASH_STMICRO
	case SPI_FLASH_CFI_MFR_STMICRO:
		debug("SF: QEB is volatile for %02x flash\n", JEDEC_MFR(info));
		return 0;
#endif
	default:
		printf("SF: Need set QEB func for %02x flash\n",
		       JEDEC_MFR(info));
		return -1;
	}
}

#if CONFIG_IS_ENABLED(OF_CONTROL)
int spi_flash_decode_fdt(struct spi_flash *flash)
{
#ifdef CONFIG_DM_SPI_FLASH
	fdt_addr_t addr;
	fdt_size_t size;

	addr = dev_read_addr_size(flash->dev, "memory-map", &size);
	if (addr == FDT_ADDR_T_NONE) {
		debug("%s: Cannot decode address\n", __func__);
		return 0;
	}

	if (flash->size > size) {
		debug("%s: Memory map must cover entire device\n", __func__);
		return -1;
	}
	flash->memory_map = map_sysmem(addr, size);
#endif

	return 0;
}
#endif /* CONFIG_IS_ENABLED(OF_CONTROL) */

int spi_flash_scan(struct spi_flash *flash)
{
	struct spi_slave *spi = flash->spi;
	const struct spi_flash_info *info = NULL;
	bool above_16MB;
	int ret;

	info = spi_flash_read_id(flash);
	if (IS_ERR_OR_NULL(info))
		return -ENOENT;

	/*
	 * Flash powers up read-only, so clear BP# bits.
	 *
	 * Note on some flash (like Macronix), QE (quad enable) bit is in the
	 * same status register as BP# bits, and we need preserve its original
	 * value during a reboot cycle as this is required by some platforms
	 * (like Intel ICH SPI controller working under descriptor mode).
	 */
	if (JEDEC_MFR(info) == SPI_FLASH_CFI_MFR_ATMEL ||
	   (JEDEC_MFR(info) == SPI_FLASH_CFI_MFR_SST) ||
	   (JEDEC_MFR(info) == SPI_FLASH_CFI_MFR_MACRONIX)) {
		u8 sr = 0;

		if (JEDEC_MFR(info) == SPI_FLASH_CFI_MFR_MACRONIX) {
			read_sr(flash, &sr);
			sr &= STATUS_QEB_MXIC;
		}
		write_sr(flash, sr);
	}

	flash->name = info->name;
	flash->memory_map = spi->memory_map;

	if (info->flags & SST_WR)
		flash->flags |= SNOR_F_SST_WR;

#ifndef CONFIG_DM_SPI_FLASH
	flash->write = spi_flash_cmd_write_ops;
#if defined(CONFIG_SPI_FLASH_SST)
	if (flash->flags & SNOR_F_SST_WR) {
		if (spi->mode & SPI_TX_BYTE)
			flash->write = sst_write_bp;
		else
			flash->write = sst_write_wp;
	}
#endif
	flash->erase = spi_flash_cmd_erase_ops;
	flash->read = spi_flash_cmd_read_ops;
#endif

#if defined(CONFIG_SPI_FLASH_STMICRO) || defined(CONFIG_SPI_FLASH_SST)
	/* NOR protection support for STmicro/Micron chips and similar */
	if (JEDEC_MFR(info) == SPI_FLASH_CFI_MFR_STMICRO ||
	    JEDEC_MFR(info) == SPI_FLASH_CFI_MFR_SST) {
		flash->flash_lock = stm_lock;
		flash->flash_unlock = stm_unlock;
		flash->flash_is_locked = stm_is_locked;
	}
#endif

/* sst26wf series block protection implementation differs from other series */
#if defined(CONFIG_SPI_FLASH_SST)
	if (JEDEC_MFR(info) == SPI_FLASH_CFI_MFR_SST && info->id[1] == 0x26) {
		flash->flash_lock = sst26_lock;
		flash->flash_unlock = sst26_unlock;
		flash->flash_is_locked = sst26_is_locked;
	}
#endif

	/* Compute the flash size */
	flash->shift = (flash->dual_flash & SF_DUAL_PARALLEL_FLASH) ? 1 : 0;
	flash->page_size = info->page_size;
	/*
	 * The Spansion S25FS512S, S25FL032P and S25FL064P have 256b pages,
	 * yet use the 0x4d00 Extended JEDEC code. The rest of the Spansion
	 * flashes with the 0x4d00 Extended JEDEC code have 512b pages.
	 * All of the others have 256b pages.
	 */
	if (JEDEC_EXT(info) == 0x4d00) {
		if ((JEDEC_ID(info) != 0x0215) &&
		    (JEDEC_ID(info) != 0x0216) &&
		    (JEDEC_ID(info) != 0x0220))
			flash->page_size = 512;
	}
	flash->page_size <<= flash->shift;
	flash->sector_size = info->sector_size << flash->shift;
	flash->size = flash->sector_size * info->n_sectors << flash->shift;
#ifdef CONFIG_SF_DUAL_FLASH
	if (flash->dual_flash & SF_DUAL_STACKED_FLASH)
		flash->size <<= 1;
#endif

#ifdef CONFIG_SPI_FLASH_USE_4K_SECTORS
	/* Compute erase sector and command */
	if (info->flags & SECT_4K) {
		flash->erase_cmd = CMD_ERASE_4K;
		flash->erase_size = 4096 << flash->shift;
	} else
#endif
	if (info->flags & SECT_4K_ONLY) {
		flash->erase_cmd = CMD_ERASE_4K;
		flash->erase_size = 4096 << flash->shift;
	} else {
		flash->erase_cmd = CMD_ERASE_64K;
		flash->erase_size = flash->sector_size;
	}

	/* Now erase size becomes valid sector size */
	flash->sector_size = flash->erase_size;

	/* Look for read commands */
	flash->read_cmd = CMD_READ_ARRAY_FAST;
	flash->read_proto = SPI_FPROTO_1_1_1;
	if (spi->mode & SPI_RX_SLOW) {
		flash->read_cmd = CMD_READ_ARRAY_SLOW;
	} else if (spi->mode & SPI_RX_QUAD && info->flags & RD_QUAD) {
		flash->read_cmd = CMD_READ_QUAD_OUTPUT_FAST;
		flash->read_proto = SPI_FPROTO_1_1_4;
	} else if (spi->mode & SPI_RX_DUAL && info->flags & RD_DUAL) {
		flash->read_cmd = CMD_READ_DUAL_OUTPUT_FAST;
		flash->read_proto = SPI_FPROTO_1_1_2;
	}

	/* Look for write commands */
	if (info->flags & WR_QPP && spi->mode & SPI_TX_QUAD) {
		flash->write_cmd = CMD_PAGE_PROGRAM_1_1_4;
		flash->write_proto = SPI_FPROTO_1_1_4;
	} else {
		/* Go for default supported write cmd */
		flash->write_cmd = CMD_PAGE_PROGRAM;
		flash->write_proto = SPI_FPROTO_1_1_1;
	}

	/* Set the quad enable bit - only for quad commands */
	if (spi_flash_protocol_get_data_nbits(flash->read_proto) == 4 ||
	    spi_flash_protocol_get_data_nbits(flash->write_proto) == 4) {
		ret = set_quad_mode(flash, info);
		if (ret) {
			debug("SF: Fail to set QEB for %02x\n",
			      JEDEC_MFR(info));
			return -EINVAL;
		}
	}

	/* Read dummy_byte: dummy byte is determined based on the
	 * dummy cycles of a particular command.
	 * Fast commands - dummy_byte = dummy_cycles/8
	 * I/O commands- dummy_byte = (dummy_cycles * no.of lines)/8
	 * For I/O commands except cmd[0] everything goes on no.of lines
	 * based on particular command but incase of fast commands except
	 * data all go on single line irrespective of command.
	 */
	switch (flash->read_cmd) {
	case CMD_READ_QUAD_IO_FAST:
		flash->dummy_byte = 2;
		break;
	case CMD_READ_ARRAY_SLOW:
		flash->dummy_byte = 0;
		break;
	default:
		flash->dummy_byte = 1;
	}

#ifdef CONFIG_SPI_FLASH_STMICRO
	if (info->flags & E_FSR)
		flash->flags |= SNOR_F_USE_FSR;
#endif

	/* Set the address length */
	flash->addr_len = SPI_FLASH_3B_ADDR_LEN;

	above_16MB = ((flash->dual_flash == SF_SINGLE_FLASH) &&
		      (flash->size > SPI_FLASH_16MB_BOUN)) ||
		     ((flash->dual_flash > SF_SINGLE_FLASH) &&
		      (flash->size > SPI_FLASH_16MB_BOUN << 1));

	/*
	 * replace the selected 3-byte address op codes with the associated
	 * 4-byte address op codes, if needed (flash->size > 16 MiB)
	 */
#ifdef CONFIG_SPI_FLASH_4BAIS
	if (above_16MB) {
		if (info->flags & NO_4BAIS) {
			puts("SF: Warning - Only lower 16MiB accessible,");
			puts(" 4-byte address instruction set not supported\n");
		} else {
			spi_flash_set_4byte_addr_opcodes(flash, info);
		}
	}
#endif

	/* Configure the BAR - discover bank cmds and read current bank */
#ifdef CONFIG_SPI_FLASH_BAR
	ret = read_bar(flash, info);
	if (ret < 0)
		return ret;
#endif

#if CONFIG_IS_ENABLED(OF_CONTROL) && !CONFIG_IS_ENABLED(OF_PLATDATA)
	ret = spi_flash_decode_fdt(flash);
	if (ret) {
		debug("SF: FDT decode error\n");
		return -EINVAL;
	}
#endif

#ifndef CONFIG_SPL_BUILD
	printf("SF: Detected %s with page size ", flash->name);
	print_size(flash->page_size, ", erase size ");
	print_size(flash->erase_size, ", total ");
	print_size(flash->size, "");
	if (flash->memory_map)
		printf(", mapped at %p", flash->memory_map);
	puts("\n");
#endif

#if !defined(CONFIG_SPI_FLASH_BAR) && !defined(CONFIG_SPI_FLASH_4BAIS)
	if (above_16MB) {
		puts("SF: Warning - Only lower 16MiB accessible,");
		puts(" Full access #define CONFIG_SPI_FLASH_BAR");
		puts(" or CONFIG_SPI_FLASH_4BAIS\n");
	}
#endif

	return 0;
}

/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2017, Atmel Corporation
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

#include "spi_flash/spi_flash.h"
#include "spi_flash/sfdp.h"
#include "debug.h"
#include "board.h"
#include "gpio.h"
#include "timer.h"
#include "div.h"
#include "fdt.h"

int spi_flash_read_reg(struct spi_flash *flash, u8 inst, u8 *buf, size_t len)
{
	struct spi_flash_command cmd;

	spi_flash_command_init(&cmd, inst, 0, SFLASH_TYPE_READ_REG);
	cmd.proto = flash->reg_proto;
	cmd.data_len = len;
	cmd.rx_data = buf;
	return spi_flash_exec(flash, &cmd);
}

int spi_flash_write_reg(struct spi_flash *flash, u8 inst,
			const u8 *buf, size_t len)
{
	struct spi_flash_command cmd;

	spi_flash_command_init(&cmd, inst, 0, SFLASH_TYPE_WRITE_REG);
	cmd.proto = flash->reg_proto;
	cmd.data_len = len;
	cmd.tx_data = buf;
	return spi_flash_exec(flash, &cmd);
}

static int spi_flash_is_ready(struct spi_flash *flash)
{
	u8 sr, fsr;
	int ret;

	ret = spi_flash_read_sr(flash, &sr);
	if (ret < 0)
		return ret;

	if (sr & SR_WIP)
		return 0;

	if (!(flash->flags & SFLASH_FLG_HAS_FSR))
		return 1;

	ret = spi_flash_read_reg(flash, SFLASH_INST_READ_FSR, &fsr, 1);
	if (ret < 0)
		return ret;

	return (fsr & FSR_READY) != 0;
}

int spi_flash_wait_till_ready_timeout(struct spi_flash *flash,
				      unsigned long timeout)
{
	unsigned long delay = 1000UL; /* 1ms */
	unsigned long loop = div(timeout + delay - 1, delay);
	int ret;

	if (!loop)
		loop = 1;
	while (loop-- > 0) {
		ret = spi_flash_is_ready(flash);
		if (ret < 0)
			return ret;
		if (ret)
			return 0;

		udelay(delay);
	}

	return -1;
}

int spi_flash_hwcaps2cmd(u32 hwcaps)
{
	switch (hwcaps) {
	case SFLASH_HWCAPS_READ:		return SFLASH_CMD_READ;
	case SFLASH_HWCAPS_READ_FAST:		return SFLASH_CMD_READ_FAST;
	case SFLASH_HWCAPS_READ_1_1_2:		return SFLASH_CMD_READ_1_1_2;
	case SFLASH_HWCAPS_READ_1_2_2:		return SFLASH_CMD_READ_1_2_2;
	case SFLASH_HWCAPS_READ_2_2_2:		return SFLASH_CMD_READ_2_2_2;
	case SFLASH_HWCAPS_READ_1_1_4:		return SFLASH_CMD_READ_1_1_4;
	case SFLASH_HWCAPS_READ_1_4_4:		return SFLASH_CMD_READ_1_4_4;
	case SFLASH_HWCAPS_READ_4_4_4:		return SFLASH_CMD_READ_4_4_4;

	case SFLASH_HWCAPS_PP:			return SFLASH_CMD_PP;
	case SFLASH_HWCAPS_PP_1_1_4:		return SFLASH_CMD_PP_1_1_4;
	case SFLASH_HWCAPS_PP_1_4_4:		return SFLASH_CMD_PP_1_4_4;
	case SFLASH_HWCAPS_PP_4_4_4:		return SFLASH_CMD_PP_4_4_4;
	}

	return -1;
}

static int spi_flash_select_read(struct spi_flash *flash,
				 const struct spi_flash_parameters *params,
				 u32 shared_hwcaps)
{
	int cmd, best_match = fls(shared_hwcaps & SFLASH_HWCAPS_READ_MASK) - 1;
	const struct spi_flash_read_command *read;

	if (best_match < 0)
		return -1;

	cmd = spi_flash_hwcaps2cmd((0x1UL << best_match));
	if (cmd < 0)
		return -1;

	read = &params->reads[cmd];
	flash->num_mode_cycles = read->num_mode_cycles;
	flash->num_wait_states = read->num_wait_states;
	flash->read_inst = read->inst;
	flash->read_proto = read->proto;
	return 0;
}

static int spi_flash_select_pp(struct spi_flash *flash,
			       const struct spi_flash_parameters *params,
			       u32 shared_hwcaps)
{
	int cmd, best_match = fls(shared_hwcaps & SFLASH_HWCAPS_PP_MASK) - 1;
	const struct spi_flash_pp_command *pp;

	if (best_match < 0)
		return -1;

	cmd = spi_flash_hwcaps2cmd((0x1UL << best_match));
	if (cmd < 0)
		return -1;

	pp = &params->page_programs[cmd];
	flash->write_inst = pp->inst;
	flash->write_proto = pp->proto;
	return 0;
}

int spi_flash_setup(struct spi_flash *flash,
		    const struct spi_flash_parameters *params,
		    const struct spi_flash_hwcaps *hwcaps)
{
	u32 ignored_mask, shared_mask;
	int err;

	/*
	 * Keep only the hardware capabilities supported by both the SPI
	 * controller and the SPI flash memory.
	 */
	shared_mask = hwcaps->mask & params->hwcaps.mask;

	/* SPI n-n-n protocols are not supported yet. */
	ignored_mask = (SFLASH_HWCAPS_READ_2_2_2 |
			SFLASH_HWCAPS_READ_4_4_4 |
			SFLASH_HWCAPS_PP_4_4_4);
	if (shared_mask & ignored_mask)
		shared_mask &= ~ignored_mask;

	/* Select the (Fast) Read command. */
	err = spi_flash_select_read(flash, params, shared_mask);
	if (err) {
		dbg_info("SF: can't select read settings supported by both the SPI controller and memory.\n");
		return err;
	}

	/* Select the Page Program command. */
	err = spi_flash_select_pp(flash, params, shared_mask);
	if (err) {
		dbg_info("SF: can't select write settings supported by both the SPI controller and memory.\n");
		return err;
	}

	/* Enable Quad I/O if needed. */
	if ((spi_flash_protocol_get_data_nbits(flash->read_proto) == 4 ||
	     spi_flash_protocol_get_data_nbits(flash->write_proto) == 4) &&
	    params->quad_enable)
		err = params->quad_enable(flash);

	return err;
}

#if defined(CONFIG_LOAD_LINUX) || defined(CONFIG_LOAD_ANDROID)
static int update_image_length(struct spi_flash *flash,
			       unsigned int offset,
			       unsigned char *dest,
			       unsigned char flag)
{
	unsigned int length = flash->page_size;
	int ret;

	ret = spi_flash_read(flash, offset, length, dest);
	if (ret)
		return -1;

	if (flag == KERNEL_IMAGE)
		return kernel_size(dest);
#ifdef CONFIG_OF_LIBFDT
	else {
		ret = check_dt_blob_valid((void *)dest);
		if (!ret)
			return of_get_dt_total_size((void *)dest);
	}
#endif
	return -1;
}
#endif

#ifdef CONFIG_QSPI_XIP
int qspi_xip(struct spi_flash *flash, void **mem);
#endif

#ifdef CONFIG_DATAFLASH_RECOVERY
int spi_flash_recovery(struct spi_flash *flash)
{
	const struct spi_flash_erase_map *map = &flash->erase_map;
	const struct spi_flash_erase_region *cur_region, *region = NULL;
	const struct spi_flash_erase_command *cur_erase, *erase = NULL;
	int i, ret;

	/*
	 * If Recovery Button is pressed during boot sequence,
	 * erase dataflash page0
	 */
	dbg_info("SF: Press the recovery button (%s) to recovery\n",
		 RECOVERY_BUTTON_NAME);

	if ((pio_get_value(CONFIG_SYS_RECOVERY_BUTTON_PIN)) == 0) {
		dbg_info("SF: The recovery button (%s) has been pressed,\n",
			 RECOVERY_BUTTON_NAME);
		dbg_info("SF: The page 0 is erasing...\n");

		/* Get the erase region for offset 0. */
		for (i = 0; i < map->num_regions; i++) {
			cur_region = &map->regions[i];

			if (!(cur_region->offset & ~SFLASH_CMD_ERASE_MASK)) {
				region = cur_region;
				break;
			}
		}

		if (!region) {
			dbg_info("SF: Can't find erase region for offset 0\n");
			return -1;
		}

		/* Get the smallest erase size of the region. */
		for (i = 0; i < SFLASH_CMD_ERASE_MAX; i++) {
			cur_erase = &map->commands[i];

			if (!(region->offset & (0x1UL << i)))
				continue;

			if (!erase || erase->size > cur_erase->size)
				erase = cur_erase;
		}

		if (!erase) {
			dbg_info("SF: Can't find erase command for offset 0\n");
			return -1;
		}

		/* Erase the smallest sector/block at offset 0. */
		ret = spi_flash_erase(flash, 0, erase->size);
		if (ret) {
			dbg_info("SF: The erasing failed\n");
			return ret;
		}
		dbg_info("SF: The erasing is done\n");

		return 0;
	}

	return -1;
}
#endif /* CONFIG_DATAFLASH_RECOVERY */

int spi_flash_loadimage(struct spi_flash *flash, struct image_info *image)
{
#if defined(CONFIG_LOAD_LINUX) || defined(CONFIG_LOAD_ANDROID)
	int length;
#endif
	int ret = 0;

#ifdef CONFIG_DATAFLASH_RECOVERY
	if (!spi_flash_recovery(flash)) {
		ret = -2;
		goto err_exit;
	}
#endif /* CONFIG_DATAFLASH_RECOVERY */

#ifdef CONFIG_OF_LIBFDT
	length = update_image_length(flash,
				     image->of_offset,
				     image->of_dest,
				     DT_BLOB);
	if (length == -1) {
		ret = -1;
		goto err_exit;
	}

	image->of_length = length;

	dbg_info("SF: dt blob: Copy %x bytes from %x to %x\n",
		 image->of_length, image->of_offset, image->of_dest);
	ret = spi_flash_read(flash,
			     image->of_offset,
			     image->of_length,
			     image->of_dest);
	if (ret) {
		dbg_info("** SF: DT: Serial flash read error**\n");
		ret = -1;
		goto err_exit;
	}
#endif /* CONFIG_OF_LIBFDT */

#ifdef CONFIG_QSPI_XIP
	ret = qspi_xip(flash, (void **)&image->dest);
	if (ret) {
		dbg_info("** SF: XIP error**\n");
		ret = -1;
		goto err_exit;
	}

	image->dest += image->offset;
	return 0;
#else /* CONFIG_QSPI_XIP */

#if defined(CONFIG_LOAD_LINUX) || defined(CONFIG_LOAD_ANDROID)
	length = update_image_length(flash,
				     image->offset,
				     image->dest,
				     KERNEL_IMAGE);
	if (length == -1) {
		ret = -1;
		goto err_exit;
	}

	image->length = length;
#endif

	dbg_info("SF: Copy %x bytes from %x to %x\n",
		 image->length, image->offset, image->dest);
	ret = spi_flash_read(flash,
			     image->offset,
			     image->length,
			     image->dest);
	if (ret) {
		dbg_info("** SF: Serial flash read error**\n");
		ret = -1;
		goto err_exit;
	}
#endif /* !CONFIG_QSPI_XIP */

err_exit:
	spi_flash_cleanup(flash);
	return ret;
}

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

#include "spi_flash/spi_nor.h"
#include "debug.h"
#include "board.h"
#include "timer.h"

static const struct spi_nor_info *spi_nor_read_id(struct spi_flash *flash)
{
	const struct spi_nor_info *info;
	char id_str[8], *p;
	int i, ret;

	ret = spi_flash_read_reg(flash, SFLASH_INST_READ_ID,
				 flash->id, sizeof(flash->id));
	if (ret < 0)
		return NULL;

	p = id_str;
	for (i = 0; i < 3; i++) {
		u8 lo = (flash->id[i] & 0x0f) >> 0;
		u8 hi = (flash->id[i] & 0xf0) >> 4;

		*p++ = (hi < 10) ? ('0' + hi) : ('a' + hi - 10);
		*p++ = (lo < 10) ? ('0' + lo) : ('a' + lo - 10);
	}
	*p = '\0';
	dbg_info("SF: Got Manufacturer and Device ID: %s\n", id_str);

	for (info = spi_nor_ids; info->name; info++) {
		if (!memcmp(info->id, flash->id, info->id_len))
			return info;
	}

	return NULL;
}

static int spi_nor_read(struct spi_flash *flash, size_t from, size_t len,
			void *buf)
{
	struct spi_flash_command cmd;

	spi_flash_command_init(&cmd, flash->read_inst, flash->addr_len,
			       SFLASH_TYPE_READ);
	cmd.proto = flash->read_proto;
	cmd.addr = from;
	cmd.mode = buf ? flash->normal_mode : flash->xip_mode;
	cmd.num_mode_cycles = flash->num_mode_cycles;
	cmd.num_wait_states = flash->num_wait_states;
	cmd.data_len = len;
	cmd.rx_data = buf;
	return spi_flash_exec(flash, &cmd);
}

static int spi_nor_write(struct spi_flash *flash, size_t to, size_t len,
			 const void *buf)
{
	struct spi_flash_command cmd;
	int ret = 0;

	spi_flash_command_init(&cmd, flash->write_inst, flash->addr_len,
			       SFLASH_TYPE_WRITE);
	cmd.proto = flash->write_proto;
	while (len) {
		size_t page_offset, page_remain;

		page_offset = to & (flash->page_size - 1);
		page_remain = min(flash->page_size - page_offset, len);

		cmd.addr = to;
		cmd.data_len = page_remain;
		cmd.tx_data = buf;

		ret = spi_flash_write_enable(flash);
		if (ret < 0)
			break;

		ret = spi_flash_exec(flash, &cmd);
		if (ret < 0)
			break;

		ret = spi_flash_wait_till_ready(flash);
		if (ret < 0)
			break;

		buf += page_remain;
		to += page_remain;
		len -= page_remain;
	}

	return ret;
}

static int spi_nor_erase(struct spi_flash *flash, size_t offset, size_t len)
{
	const struct spi_flash_erase_map *map = &flash->erase_map;
	struct spi_flash_command cmd;
	int ret = 0;

	if (!spi_flash_has_uniform_erase(flash))
		/* TODO: add support to none uniform erase map. */
		return -1;

	spi_flash_command_init(&cmd, 0, flash->addr_len,
			       SFLASH_TYPE_ERASE);
	cmd.proto = flash->erase_proto;
	while (len) {
		const struct spi_flash_erase_command *erase = NULL;
		size_t i;

		for (i = 0; i < SFLASH_CMD_ERASE_MAX; i++) {
			const struct spi_flash_erase_command *e;
			u32 rem;

			if (!(map->uniform_region.offset & (0x1UL << i)))
				continue;

			e = &map->commands[i];
			spi_flash_div_by_erase_size(e, offset, &rem);
			if (rem)
				continue;

			if (e->size <= len &&
			    (!erase || erase->size < e->size))
				erase = e;
		}
		if (!erase)
			return -1;

		ret = spi_flash_write_enable(flash);
		if (ret < 0)
			break;

		cmd.inst = erase->inst;
		cmd.addr = offset;
		ret = spi_flash_exec(flash, &cmd);
		if (ret < 0)
			break;

		ret = spi_flash_wait_till_ready(flash);
		if (ret < 0)
			break;

		offset += erase->size;
		len -= erase->size;
	}

	return ret;
}

static int spi_nor_init_params(struct spi_flash *flash,
			       const struct spi_nor_info *info,
			       struct spi_flash_parameters *params)
{
	struct spi_flash_erase_map *map = &flash->erase_map;
	u32 erase_mask = 0;

	if (info && info->params)
		memcpy(params, info->params, sizeof(*params));
	else
		memset(params, 0, sizeof(*params));

	/* Set SPI NOR sizes. */
	if (info) {
		params->size = info->sector_size * info->n_sectors;
		params->page_size = info->page_size;
	} else {
		flash->page_size = 256;
	}

	/* Set legacy flash parameters as default. */
	if (info && info->params)
		goto set_erase_map;

	/* (Fast) Read settings. */
	params->hwcaps.mask |= SFLASH_HWCAPS_READ;
	spi_flash_set_read_settings(&params->reads[SFLASH_CMD_READ],
				    0, 0, SFLASH_INST_READ,
				    SFLASH_PROTO_1_1_1);
	if (!info || !(info->flags & SNOR_NO_FR)) {
		params->hwcaps.mask |= SFLASH_HWCAPS_READ_FAST;
		spi_flash_set_read_settings(&params->reads[SFLASH_CMD_READ_FAST],
					    0, 8, SFLASH_INST_FAST_READ,
					    SFLASH_PROTO_1_1_1);
	}

	/* Page Program settings. */
	params->hwcaps.mask |= SFLASH_HWCAPS_PP;
	spi_flash_set_pp_settings(&params->page_programs[SFLASH_CMD_PP],
				  SFLASH_INST_PAGE_PROGRAM,
				  SFLASH_PROTO_1_1_1);

	/* Sector Erase settings. */
set_erase_map:
	erase_mask |= (0x1UL << 0);
	spi_flash_set_erase_command(&map->commands[0],
				    info ? info->sector_size : 65536u,
				    SFLASH_INST_ERASE_64K);
	if (info && info->flags & SNOR_SECT_4K_ONLY) {
		spi_flash_set_erase_command(&map->commands[0],
					    4096u,
					    SFLASH_INST_ERASE_4K);
	} else if (info && info->flags & SNOR_SECT_4K) {
		erase_mask |= (0x1UL << 1);
		spi_flash_set_erase_command(&map->commands[1],
					    4096u,
					    SFLASH_INST_ERASE_4K);
	}
	spi_flash_init_uniform_erase_map(map, erase_mask, params->size);

	/* Quad Enable Requirements & 0-4-4 Mode Entry Method. */
	flash->enable_0_4_4 = NULL;
	switch (spi_flash_get_mfr(flash)) {
	case SFLASH_MFR_MICRON:
		params->quad_enable = NULL;
		flash->enable_0_4_4 = micron_enable_0_4_4;
		flash->xip_mode = 0x00u;
		break;

	case SFLASH_MFR_MACRONIX:
		params->quad_enable = macronix_quad_enable;
		break;

	default:
		params->quad_enable = spansion_quad_enable;
		break;
	}

	/* Override the parameters with data read from SFDP tables. */
	if (!info || !(info->flags & SNOR_SKIP_SFDP))
		spi_flash_parse_sfdp(flash, params);

	return 0;
}

static u8 spi_nor_convert_inst(u8 inst, const u8 table[][2], size_t size)
{
	size_t i;

	for (i = 0; i < size; i++)
		if (table[i][0] == inst)
			return table[i][1];

	/* No conversion found, keep input op code. */
	return inst;
}

static inline u8 spi_nor_convert_3to4_read(u8 inst)
{
	static const u8 spi_nor_3to4_read[][2] = {
		{ SFLASH_INST_READ,		SFLASH_INST_READ_4B },
		{ SFLASH_INST_FAST_READ,	SFLASH_INST_FAST_READ_4B },
		{ SFLASH_INST_FAST_READ_1_1_2,	SFLASH_INST_FAST_READ_1_1_2_4B },
		{ SFLASH_INST_FAST_READ_1_2_2,	SFLASH_INST_FAST_READ_1_2_2_4B },
		{ SFLASH_INST_FAST_READ_1_1_4,	SFLASH_INST_FAST_READ_1_1_4_4B },
		{ SFLASH_INST_FAST_READ_1_4_4,	SFLASH_INST_FAST_READ_1_4_4_4B },
	};

	return spi_nor_convert_inst(inst, spi_nor_3to4_read,
				    ARRAY_SIZE(spi_nor_3to4_read));
}

static inline u8 spi_nor_convert_3to4_write(u8 inst)
{
	static const u8 spi_nor_3to4_program[][2] = {
		{ SFLASH_INST_PAGE_PROGRAM,		SFLASH_INST_PAGE_PROGRAM_4B },
		{ SFLASH_INST_PAGE_PROGRAM_1_1_4,	SFLASH_INST_PAGE_PROGRAM_1_1_4_4B },
		{ SFLASH_INST_PAGE_PROGRAM_1_4_4,	SFLASH_INST_PAGE_PROGRAM_1_4_4_4B },
	};

	return spi_nor_convert_inst(inst, spi_nor_3to4_program,
				    ARRAY_SIZE(spi_nor_3to4_program));
}

static inline u8 spi_nor_convert_3to4_erase(u8 inst)
{
	static const u8 spi_nor_3to4_erase[][2] = {
		{ SFLASH_INST_ERASE_4K,		SFLASH_INST_ERASE_4K_4B },
		{ SFLASH_INST_ERASE_32K,	SFLASH_INST_ERASE_32K_4B },
		{ SFLASH_INST_ERASE_64K,	SFLASH_INST_ERASE_64K_4B },
	};

	return spi_nor_convert_inst(inst, spi_nor_3to4_erase,
				    ARRAY_SIZE(spi_nor_3to4_erase));
}

static void spi_nor_set_4bais(struct spi_flash *flash)
{
	struct spi_flash_erase_map *map = &flash->erase_map;
	int i;

	flash->read_inst = spi_nor_convert_3to4_read(flash->read_inst);
	flash->write_inst = spi_nor_convert_3to4_write(flash->write_inst);

	for (i = 0; i < SFLASH_CMD_ERASE_MAX; i++) {
		struct spi_flash_erase_command *cmd = &map->commands[i];

		cmd->inst = spi_nor_convert_3to4_erase(cmd->inst);
	}

	flash->addr_len = 4;
}

static int sst26_unlock_block_protection(struct spi_flash *flash)
{
	int rc;

	rc = spi_flash_write_enable(flash);
	rc = rc ? rc : spi_flash_write_reg(flash, SFLASH_INST_ULBPR, NULL, 0);
	return rc ? rc : spi_flash_wait_till_ready(flash);
}

int spi_nor_probe(struct spi_flash *flash,
		  const struct spi_flash_hwcaps *hwcaps)
{
	struct spi_flash_parameters params;
	const struct spi_nor_info *info;
	int ret;

	/* Check minimum requirement. */
	if (!flash->ops)
		return -1;

	/* Set the SPI mode. */
	ret = spi_flash_set_mode(flash, CONFIG_SYS_SPI_MODE);
	if (ret)
		return -1;

	/* Set the new baudrate. */
	ret = spi_flash_set_freq(flash, CONFIG_SYS_SPI_CLOCK);
	if (ret)
		return -1;

	/*
	 * Reset the SPI flash memory:
	 * exit the 4-byte address and/or the 4-4-4 modes.
	 */
	if (hwcaps->mask & (SFLASH_HWCAPS_READ_4_4_4 |
			    SFLASH_HWCAPS_PP_4_4_4)) {
		flash->reg_proto = SFLASH_PROTO_4_4_4;
		spi_flash_reset(flash);
	}
	flash->reg_proto = SFLASH_PROTO_1_1_1;
	spi_flash_reset(flash);
	udelay(10);

	/* Set default settings. */
	flash->read_proto = SFLASH_PROTO_1_1_1;
	flash->write_proto = SFLASH_PROTO_1_1_1;
	flash->read = spi_nor_read;
	flash->write = spi_nor_write;
	flash->erase = spi_nor_erase;
	flash->flags = 0;
	flash->normal_mode = 0xFFu;
	flash->xip_mode = 0xA5u;

	/* Get the static SPI flash info (might by NULL). */
	info = spi_nor_read_id(flash);
	if (!info)
		goto init_params;

	if (info->flags & SNOR_HAS_FSR)
		flash->flags |= SFLASH_FLG_HAS_FSR;

	if (info->flags & SNOR_SST_ULBPR)
		if (sst26_unlock_block_protection(flash))
			dbg_info("SF: WARNING: SST26 - can't unlock block protection\n");

	/* Parse the Serial Flash Discoverable Parameter tables. */
init_params:
	ret = spi_nor_init_params(flash, info, &params);
	if (ret < 0)
		return ret;

	flash->size = params.size;
	flash->page_size = params.page_size;

	/*
	 * Configure the SPI memory:
	 * - select instructions for (Fast) Read, Page Program and Sector Erase.
	 * - set the number of mode & wait state cycles.
	 * - set the SPI protocols for register and memory accesses.
	 * - set the Quad Enable bit if needed (required by SPI x-y-4 protos).
	 */
	ret = spi_flash_setup(flash, &params, hwcaps);
	if (ret < 0)
		return ret;

	if (info && info->addr_len) {
		flash->addr_len = info->addr_len;
	} else {
		flash->addr_len = 3;
		if (flash->size > 0x01000000u) {
			if (!info || !(info->flags & SNOR_NO_4BAIS))
				spi_nor_set_4bais(flash);
			else
				dbg_info("SF: WARNING: can't read above 16MiB\n");
		}
	}

	return 0;
}

int spansion_quad_enable(struct spi_flash *flash)
{
	u8 sr[2];
	int rc;

	/* Keep the current value of the Status Register. */
	rc = spi_flash_read_sr(flash, &sr[0]);
	if (rc)
		return rc;

	/* Set the Status and Control/Status 2 Registers. */
	sr[1] = CR_QE_SPAN;
	rc = spi_flash_write_enable(flash);
	rc = rc ? rc : spi_flash_write_reg(flash, SFLASH_INST_WRITE_SR, sr, 2);
	return rc ? rc : spi_flash_wait_till_ready(flash);
}

int spansion_new_quad_enable(struct spi_flash *flash)
{
	u8 sr[2];
	int rc;

	/* Check the current Quand Enable bit value. */
	rc = spi_flash_read_reg(flash, SFLASH_INST_READ_CR, &sr[1], 1);
	if (rc)
		return rc;

	if (sr[1] & CR_QE_SPAN)
		return 0;

	/* Keep the current value of the Status Register. */
	rc = spi_flash_read_sr(flash, &sr[0]);
	if (rc)
		return rc;

	/* Set the Status and Control/Status 2 Registers. */
	sr[1] |= CR_QE_SPAN;
	rc = spi_flash_write_enable(flash);
	rc = rc ? rc : spi_flash_write_reg(flash, SFLASH_INST_WRITE_SR, sr, 2);
	rc = rc ? rc : spi_flash_wait_till_ready(flash);
	if (rc)
		return rc;

	/* Read back and check. */
	rc = spi_flash_read_reg(flash, SFLASH_INST_READ_CR, &sr[1], 1);
	if (rc)
		return rc;

	return (sr[1] & CR_QE_SPAN) ? 0 : -1;
}

int macronix_quad_enable(struct spi_flash *flash)
{
	u8 sr;
	int rc;

	/* Check the current Quand Enable bit value. */
	rc = spi_flash_read_sr(flash, &sr);
	if (rc)
		return rc;

	if (sr & SR_QE_MX)
		return 0;

	/* Set the Status and Control/Status 2 Registers. */
	sr |= SR_QE_MX;
	rc = spi_flash_write_enable(flash);
	rc = rc ? rc : spi_flash_write_reg(flash, SFLASH_INST_WRITE_SR, &sr, 1);
	rc = rc ? rc : spi_flash_wait_till_ready(flash);
	if (rc)
		return rc;

	/* Read back and check. */
	rc = spi_flash_read_sr(flash, &sr);
	if (rc)
		return rc;

	return (sr & SR_QE_MX) ? 0 : -1;
}

int sr2_bit7_quad_enable(struct spi_flash *flash)
{
	u8 sr2;
	int rc;

	/* Check the current Quand Enable bit value. */
	rc = spi_flash_read_reg(flash, SFLASH_INST_READ_SR2, &sr2, 1);
	if (rc)
		return rc;

	if (sr2 & SR2_QE_BIT7)
		return 0;

	/* Set the Status and Control/Status 2 Registers. */
	sr2 |= SR2_QE_BIT7;
	rc = spi_flash_write_enable(flash);
	rc = rc ? rc : spi_flash_write_reg(flash, SFLASH_INST_WRITE_SR2, &sr2, 1);
	rc = rc ? rc : spi_flash_wait_till_ready(flash);
	if (rc)
		return rc;

	/* Read back and check. */
	rc = spi_flash_read_reg(flash, SFLASH_INST_READ_SR2, &sr2, 1);
	if (rc)
		return rc;

	return (sr2 & SR2_QE_BIT7) ? 0 : -1;
}

int micron_enable_0_4_4(struct spi_flash *flash, bool enable)
{
	u8 vcr, mask, value;
	int rc;

	mask = VCR_XIP;
	value = (enable) ? 0 : VCR_XIP;

	/* Check the current XIP bit value. */
	rc = spi_flash_read_reg(flash, SFLASH_INST_READ_VCR, &vcr, 1);
	if (rc)
		return rc;

	if ((vcr & mask) == value)
		return 0;

	/* Update the XIP bit value. */
	vcr = (vcr & ~mask) | value;
	rc = spi_flash_write_enable(flash);
	rc = rc ? rc : spi_flash_write_reg(flash, SFLASH_INST_WRITE_VCR, &vcr, 1);
	rc = rc ? rc : spi_flash_wait_till_ready(flash);
	if (rc)
		return rc;

	/* Read back and check. */
	rc = spi_flash_read_reg(flash, SFLASH_INST_READ_VCR, &vcr, 1);
	if (rc)
		return rc;

	return ((vcr & mask) == value) ? 0 : -1;
}

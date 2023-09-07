/*
 * This file is part of the coreboot project.
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
#include <cbmem.h>
#include <device/device.h>
#include <device/dram/ddr4.h>
#include <string.h>
#include <memory_info.h>
#include <smbios.h>
#include <types.h>

typedef enum {
	BLOCK_0, /* Base Configuration and DRAM Parameters */
	BLOCK_1,
	BLOCK_1_L, /* Standard Module Parameters */
	BLOCK_1_H, /* Hybrid Module Parameters */
	BLOCK_2,
	BLOCK_2_L, /* Hybrid Module Extended Function Parameters */
	BLOCK_2_H, /* Manufacturing Information */
	BLOCK_3    /* End user programmable */
} spd_block_type;

typedef struct {
	spd_block_type type;
	uint16_t start;     /* starting offset from beginning of the spd */
	uint16_t len;       /* size of the block */
	uint16_t crc_start; /* offset from start of crc bytes, 0 if none */
} spd_block;

/* 'SPD contents architecture' as per datasheet */
const spd_block spd_blocks[] = {
	{.type = BLOCK_0, 0, 128, 126},   {.type = BLOCK_1, 128, 128, 126},
	{.type = BLOCK_1_L, 128, 64, 0},  {.type = BLOCK_1_H, 192, 64, 0},
	{.type = BLOCK_2_L, 256, 64, 62}, {.type = BLOCK_2_H, 320, 64, 0},
	{.type = BLOCK_3, 384, 128, 0} };

static bool verify_block(const spd_block *block, spd_raw_data spd)
{
	uint16_t crc, spd_crc;

	spd_crc = (spd[block->start + block->crc_start + 1] << 8)
		  | spd[block->start + block->crc_start];
	crc = ddr_crc16(&spd[block->start], block->len - 2);

	return spd_crc == crc;
}

/* Check if given block is 'reserved' for a given module type */
static bool block_exists(spd_block_type type, u8 dimm_type)
{
	bool is_hybrid;

	switch (type) {
	case BLOCK_0: /* fall-through */
	case BLOCK_1: /* fall-through */
	case BLOCK_1_L: /* fall-through */
	case BLOCK_1_H: /* fall-through */
	case BLOCK_2_H: /* fall-through */
	case BLOCK_3: /* fall-through */
		return true;
	case BLOCK_2_L:
		is_hybrid = (dimm_type >> 4) & ((1 << 3) - 1);
		if (is_hybrid)
			return true;
		return false;
	default: /* fall-through */
		return false;
	}
}


/**
 * \brief Decode the raw SPD data
 *
 * Decodes a raw SPD data from a DDR4 DIMM, and organizes it into a
 * @ref dimm_attr structure. The SPD data must first be read in a contiguous
 * array, and passed to this function.
 *
 * @param dimm pointer to @ref dimm_attr structure where the decoded data is to
 *	       be stored
 * @param spd array of raw data previously read from the SPD.
 *
 * @return @ref spd_status enumerator
 *		SPD_STATUS_OK -- decoding was successful
 *		SPD_STATUS_INVALID -- invalid SPD or not a DDR4 SPD
 *		SPD_STATUS_CRC_ERROR -- checksum mismatch
 */
int spd_decode_ddr4(dimm_attr *dimm, spd_raw_data spd)
{
	u8 reg8;
	u8 bus_width, sdram_width;
	u16 cap_per_die_mbit;
	u16 spd_bytes_total, spd_bytes_used;
	const uint16_t spd_bytes_used_table[] = {0, 128, 256, 384, 512};

	/* Make sure that the SPD dump is indeed from a DDR4 module */
	if (spd[2] != SPD_MEMORY_TYPE_DDR4_SDRAM) {
		printk(BIOS_ERR, "Not a DDR4 SPD!\n");
		dimm->dram_type = SPD_MEMORY_TYPE_UNDEFINED;
		return SPD_STATUS_INVALID;
	}

	spd_bytes_total = (spd[0] >> 4) & 0x7;
	spd_bytes_used = spd[0] & 0xf;

	if (!spd_bytes_total || !spd_bytes_used) {
		printk(BIOS_ERR, "SPD failed basic sanity checks\n");
		return SPD_STATUS_INVALID;
	}

	if (spd_bytes_total >= 3)
		printk(BIOS_WARNING, "SPD Bytes Total value is reserved\n");

	spd_bytes_total = 256 << (spd_bytes_total - 1);

	if (spd_bytes_used > 4) {
		printk(BIOS_ERR, "SPD Bytes Used value is reserved\n");
		return SPD_STATUS_INVALID;
	}

	spd_bytes_used = spd_bytes_used_table[spd_bytes_used];

	if (spd_bytes_used > spd_bytes_total) {
		printk(BIOS_ERR, "SPD Bytes Used is greater than SPD Bytes Total\n");
		return SPD_STATUS_INVALID;
	}

	/* Verify CRC of blocks that have them, do not step over 'used' length */
	for (int i = 0; i < ARRAY_SIZE(spd_blocks); i++) {
		/* this block is not checksumed */
		if (spd_blocks[i].crc_start == 0)
			continue;
		/* we shouldn't have this block */
		if (spd_blocks[i].start + spd_blocks[i].len > spd_bytes_used)
			continue;
		/* check if block exists in the current schema */
		if (!block_exists(spd_blocks[i].type, spd[3]))
			continue;
		if (!verify_block(&spd_blocks[i], spd)) {
			printk(BIOS_ERR, "CRC failed for block %d\n", i);
			return SPD_STATUS_CRC_ERROR;
		}
	}

	dimm->dram_type = SPD_MEMORY_TYPE_DDR4_SDRAM;
	dimm->dimm_type = spd[3] & ((1 << 4) - 1);

	reg8 = spd[13] & ((1 << 4) - 1);
	dimm->bus_width = reg8;
	bus_width = 8 << (reg8 & ((1 << 3) - 1));

	reg8 = spd[12] & ((1 << 3) - 1);
	dimm->sdram_width = reg8;
	sdram_width = 4 << reg8;

	reg8 = spd[4] & ((1 << 4) - 1);
	dimm->cap_per_die_mbit = reg8;
	cap_per_die_mbit = (1 << reg8) * 256;

	reg8 = (spd[12] >> 3) & ((1 << 3) - 1);
	dimm->ranks = reg8 + 1;

	if (!bus_width || !sdram_width) {
		printk(BIOS_ERR, "SPD information is invalid");
		dimm->size_mb = 0;
		return SPD_STATUS_INVALID;
	}

	/* seems to be only one, in mV */
	dimm->vdd_voltage = 1200;

	/* calculate size */
	dimm->size_mb = cap_per_die_mbit / 8 * bus_width / sdram_width * dimm->ranks;

	/* make sure we have the manufacturing information block */
	if (spd_bytes_used > 320) {
		dimm->manufacturer_id = (spd[351] << 8) | spd[350];
		memcpy(dimm->part_number, &spd[329], SPD_DDR4_PART_LEN);
		dimm->part_number[SPD_DDR4_PART_LEN] = 0;
		memcpy(dimm->serial_number, &spd[325], sizeof(dimm->serial_number));
	}
	return SPD_STATUS_OK;
}

enum cb_err spd_add_smbios17_ddr4(const u8 channel, const u8 slot, const u16 selected_freq,
				  const dimm_attr *info)
{
	struct memory_info *mem_info;
	struct dimm_info *dimm;

	/*
	 * Allocate CBMEM area for DIMM information used to populate SMBIOS
	 * table 17
	 */
	mem_info = cbmem_find(CBMEM_ID_MEMINFO);
	if (!mem_info) {
		mem_info = cbmem_add(CBMEM_ID_MEMINFO, sizeof(*mem_info));

		printk(BIOS_DEBUG, "CBMEM entry for DIMM info: 0x%p\n", mem_info);
		if (!mem_info)
			return CB_ERR;

		memset(mem_info, 0, sizeof(*mem_info));
	}

	if (mem_info->dimm_cnt >= ARRAY_SIZE(mem_info->dimm)) {
		printk(BIOS_WARNING, "BUG: Too many DIMM infos for %s.\n", __func__);
		return CB_ERR;
	}

	dimm = &mem_info->dimm[mem_info->dimm_cnt];
	if (info->size_mb) {
		dimm->ddr_type = MEMORY_TYPE_DDR4;
		dimm->ddr_frequency = selected_freq;
		dimm->dimm_size = info->size_mb;
		dimm->channel_num = channel;
		dimm->rank_per_dimm = info->ranks;
		dimm->dimm_num = slot;
		memcpy(dimm->module_part_number, info->part_number, SPD_DDR4_PART_LEN);
		dimm->mod_id = info->manufacturer_id;

		switch (info->dimm_type) {
		case SPD_DIMM_TYPE_SO_DIMM:
			dimm->mod_type = SPD_SODIMM;
			break;
		case SPD_DIMM_TYPE_72B_SO_RDIMM:
			dimm->mod_type = SPD_72B_SO_RDIMM;
			break;
		case SPD_DIMM_TYPE_UDIMM:
			dimm->mod_type = SPD_UDIMM;
			break;
		case SPD_DIMM_TYPE_RDIMM:
			dimm->mod_type = SPD_RDIMM;
			break;
		default:
			dimm->mod_type = SPD_UNDEFINED;
			break;
		}

		dimm->bus_width = info->bus_width;
		memcpy(dimm->serial, info->serial_number,
		       MIN(sizeof(dimm->serial), sizeof(info->serial_number)));

		dimm->vdd_voltage = info->vdd_voltage;
		mem_info->dimm_cnt++;
	}

	return CB_SUCCESS;
}

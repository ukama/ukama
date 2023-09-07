/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2016 Intel Corporation.
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

#include <cbfs.h>
#include <console/console.h>
#include <spd_bin.h>
#include <string.h>
#include <device/early_smbus.h>
#include <device/dram/ddr3.h>

static u8 spd_data[CONFIG_DIMM_MAX * CONFIG_DIMM_SPD_SIZE] CAR_GLOBAL;

void dump_spd_info(struct spd_block *blk)
{
	u8 i;

	for (i = 0; i < CONFIG_DIMM_MAX; i++)
		if (blk->spd_array[i] != NULL && blk->spd_array[i][0] != 0) {
			printk(BIOS_DEBUG, "SPD @ 0x%02X\n", blk->addr_map[i]);
			print_spd_info(blk->spd_array[i]);
		}
}

static bool is_memory_type_ddr4(int dram_type)
{
	return (dram_type == SPD_DRAM_DDR4);
}

static const char *spd_get_module_type_string(int dram_type)
{
	switch (dram_type) {
	case SPD_DRAM_DDR3:
		return "DDR3";
	case SPD_DRAM_LPDDR3_INTEL:
	case SPD_DRAM_LPDDR3_JEDEC:
		return "LPDDR3";
	case SPD_DRAM_DDR4:
		return "DDR4";
	}
	return "UNKNOWN";
}

static int spd_get_banks(const uint8_t spd[], int dram_type)
{
	static const int ddr3_banks[4] = { 8, 16, 32, 64 };
	static const int ddr4_banks[10] = { 4, 8, -1, -1, 8, 16, -1, -1, 16, 32 };
	int index = (spd[SPD_DENSITY_BANKS] >> 4) & 0xf;
	switch (dram_type) {
	/* DDR3 and LPDDR3 has the same bank definition */
	case SPD_DRAM_DDR3:
	case SPD_DRAM_LPDDR3_INTEL:
	case SPD_DRAM_LPDDR3_JEDEC:
		if (index >= ARRAY_SIZE(ddr3_banks))
			return -1;
		return ddr3_banks[index];
	case SPD_DRAM_DDR4:
		if (index >= ARRAY_SIZE(ddr4_banks))
			return -1;
		return ddr4_banks[index];
	default:
		return -1;
	}
}

static int spd_get_capmb(const uint8_t spd[])
{
	static const int spd_capmb[10] = { 1, 2, 4, 8, 16, 32, 64, 128, 48, 96 };
	int index = spd[SPD_DENSITY_BANKS] & 0xf;
	if (index >= ARRAY_SIZE(spd_capmb))
		return -1;
	return spd_capmb[index] * 256;
}

static int spd_get_rows(const uint8_t spd[])
{
	static const int spd_rows[7]  = { 12, 13, 14, 15, 16, 17, 18 };
	int index = (spd[SPD_ADDRESSING] >> 3) & 7;
	if (index >= ARRAY_SIZE(spd_rows))
		return -1;
	return spd_rows[index];
}

static int spd_get_cols(const uint8_t spd[])
{
	static const int spd_cols[4]  = { 9, 10, 11, 12 };
	int index = spd[SPD_ADDRESSING] & 7;
	if (index >= ARRAY_SIZE(spd_cols))
		return -1;
	return spd_cols[index];
}

static int spd_get_ranks(const uint8_t spd[], int dram_type)
{
	static const int spd_ranks[8] = { 1, 2, 3, 4, 5, 6, 7, 8 };
	int organ_offset = is_memory_type_ddr4(dram_type) ? DDR4_ORGANIZATION
							: DDR3_ORGANIZATION;
	int index = (spd[organ_offset] >> 3) & 7;
	if (index >= ARRAY_SIZE(spd_ranks))
		return -1;
	return spd_ranks[index];
}

static int spd_get_devw(const uint8_t spd[], int dram_type)
{
	static const int spd_devw[4]  = { 4, 8, 16, 32 };
	int organ_offset = is_memory_type_ddr4(dram_type) ? DDR4_ORGANIZATION
							: DDR3_ORGANIZATION;
	int index = spd[organ_offset] & 7;
	if (index >= ARRAY_SIZE(spd_devw))
		return -1;
	return spd_devw[index];
}

static int spd_get_busw(const uint8_t spd[], int dram_type)
{
	static const int spd_busw[4]  = { 8, 16, 32, 64 };
	int busw_offset = is_memory_type_ddr4(dram_type) ? DDR4_BUS_DEV_WIDTH
							: DDR3_BUS_DEV_WIDTH;
	int index = spd[busw_offset] & 7;
	if (index >= ARRAY_SIZE(spd_busw))
		return -1;
	return spd_busw[index];
}

static void spd_get_name(const uint8_t spd[], char spd_name[], int dram_type)
{
	switch (dram_type) {
	case SPD_DRAM_DDR3:
		memcpy(spd_name, &spd[DDR3_SPD_PART_OFF], DDR3_SPD_PART_LEN);
		spd_name[DDR3_SPD_PART_LEN] = 0;
		break;
	case SPD_DRAM_LPDDR3_INTEL:
	case SPD_DRAM_LPDDR3_JEDEC:
		memcpy(spd_name, &spd[LPDDR3_SPD_PART_OFF],
			LPDDR3_SPD_PART_LEN);
		spd_name[LPDDR3_SPD_PART_LEN] = 0;
		break;
	case SPD_DRAM_DDR4:
		memcpy(spd_name, &spd[DDR4_SPD_PART_OFF], DDR4_SPD_PART_LEN);
		spd_name[DDR4_SPD_PART_LEN] = 0;
		break;
	default:
		break;
	}
}

void print_spd_info(uint8_t spd[])
{
	char spd_name[DDR4_SPD_PART_LEN+1] = { 0 };
	int type  = spd[SPD_DRAM_TYPE];
	int banks = spd_get_banks(spd, type);
	int capmb = spd_get_capmb(spd);
	int rows  = spd_get_rows(spd);
	int cols  = spd_get_cols(spd);
	int ranks = spd_get_ranks(spd, type);
	int devw  = spd_get_devw(spd, type);
	int busw  = spd_get_busw(spd, type);

	/* Module type */
	printk(BIOS_INFO, "SPD: module type is %s\n",
		spd_get_module_type_string(type));
	/* Module Part Number */
	spd_get_name(spd, spd_name, type);

	printk(BIOS_INFO, "SPD: module part is %s\n", spd_name);

	printk(BIOS_INFO,
		"SPD: banks %d, ranks %d, rows %d, columns %d, density %d Mb\n",
		banks, ranks, rows, cols, capmb);
	printk(BIOS_INFO, "SPD: device width %d bits, bus width %d bits\n",
		devw, busw);

	if (capmb > 0 && busw > 0 && devw > 0 && ranks > 0) {
		/* SIZE = DENSITY / 8 * BUS_WIDTH / SDRAM_WIDTH * RANKS */
		printk(BIOS_INFO, "SPD: module size is %u MB (per channel)\n",
			capmb / 8 * busw / devw * ranks);
	}
}

static void update_spd_len(struct spd_block *blk)
{
	u8 i, j = 0;
	for (i = 0 ; i < CONFIG_DIMM_MAX; i++)
		if (blk->spd_array[i] != NULL)
			j |= blk->spd_array[i][SPD_DRAM_TYPE];

	/* If spd used is DDR4, then its length is 512 byte. */
	if (j == SPD_DRAM_DDR4)
		blk->len = SPD_PAGE_LEN_DDR4;
	else
		blk->len = SPD_PAGE_LEN;
}

int get_spd_cbfs_rdev(struct region_device *spd_rdev, u8 spd_index)
{
	struct cbfsf fh;

	uint32_t cbfs_type = CBFS_TYPE_SPD;

	if (cbfs_boot_locate(&fh, "spd.bin", &cbfs_type) < 0)
		return -1;
	cbfs_file_data(spd_rdev, &fh);
	return rdev_chain(spd_rdev, spd_rdev, spd_index * CONFIG_DIMM_SPD_SIZE,
							CONFIG_DIMM_SPD_SIZE);
}

static void smbus_read_spd(u8 *spd, u8 addr)
{
	u16 i;
	u8 step = 1;

	if (CONFIG(SPD_READ_BY_WORD))
		step = sizeof(uint16_t);

	for (i = 0; i < SPD_PAGE_LEN; i += step) {
		if (CONFIG(SPD_READ_BY_WORD))
			((u16*)spd)[i / sizeof(uint16_t)] =
				 smbus_read_word(0, addr, i);
		else
			spd[i] = smbus_read_byte(0, addr, i);
	}
}

static void get_spd(u8 *spd, u8 addr)
{
	if (smbus_read_byte(0, addr, 0) == 0xff) {
		printk(BIOS_INFO, "No memory dimm at address %02X\n",
			addr << 1);
		/* Make sure spd is zeroed if dimm doesn't exist. */
		memset(spd, 0, CONFIG_DIMM_SPD_SIZE);
		return;
	}
	smbus_read_spd(spd, addr);

	/* Check if module is DDR4, DDR4 spd is 512 byte. */
	if (spd[SPD_DRAM_TYPE] == SPD_DRAM_DDR4 &&
		CONFIG_DIMM_SPD_SIZE > SPD_PAGE_LEN) {
		/* Switch to page 1 */
		smbus_write_byte(0, SPD_PAGE_1, 0, 0);
		smbus_read_spd(spd + SPD_PAGE_LEN, addr);
		/* Restore to page 0 */
		smbus_write_byte(0, SPD_PAGE_0, 0, 0);
	}
}

void get_spd_smbus(struct spd_block *blk)
{
	u8 i;
	unsigned char *spd_data_ptr = car_get_var_ptr(&spd_data);

	for (i = 0 ; i < CONFIG_DIMM_MAX; i++) {
		get_spd(spd_data_ptr + i * CONFIG_DIMM_SPD_SIZE,
			blk->addr_map[i]);
		blk->spd_array[i] = spd_data_ptr + i * CONFIG_DIMM_SPD_SIZE;
	}

	update_spd_len(blk);
}

#if CONFIG_DIMM_SPD_SIZE == 128
int read_ddr3_spd_from_cbfs(u8 *buf, int idx)
{
	const int SPD_CRC_HI = 127;
	const int SPD_CRC_LO = 126;

	const char *spd_file;
	size_t spd_file_len = 0;
	size_t min_len = (idx + 1) * CONFIG_DIMM_SPD_SIZE;

	spd_file = cbfs_boot_map_with_leak("spd.bin", CBFS_TYPE_SPD,
						&spd_file_len);
	if (!spd_file)
		printk(BIOS_EMERG, "file [spd.bin] not found in CBFS");
	if (spd_file_len < min_len)
		printk(BIOS_EMERG, "Missing SPD data.");
	if (!spd_file || spd_file_len < min_len)
		return -1;

	memcpy(buf, spd_file + (idx * CONFIG_DIMM_SPD_SIZE),
		CONFIG_DIMM_SPD_SIZE);

	u16 crc = spd_ddr3_calc_crc(buf, CONFIG_DIMM_SPD_SIZE);

	if (((buf[SPD_CRC_LO] == 0) && (buf[SPD_CRC_HI] == 0))
		|| (buf[SPD_CRC_LO] != (crc & 0xff))
		|| (buf[SPD_CRC_HI] != (crc >> 8))) {
		printk(BIOS_WARNING,
			"SPD CRC %02x%02x is invalid, should be %04x\n",
			buf[SPD_CRC_HI], buf[SPD_CRC_LO], crc);
		buf[SPD_CRC_LO] = crc & 0xff;
		buf[SPD_CRC_HI] = crc >> 8;
		u16 i;
		printk(BIOS_WARNING, "\nDisplay the SPD");
		for (i = 0; i < CONFIG_DIMM_SPD_SIZE; i++) {
			if ((i % 16) == 0x00)
				printk(BIOS_WARNING, "\n%02x:  ", i);
			printk(BIOS_WARNING, "%02x ", buf[i]);
		}
		printk(BIOS_WARNING, "\n");
	}
	return 0;
}
#endif

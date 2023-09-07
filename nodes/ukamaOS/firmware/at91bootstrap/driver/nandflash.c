/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2006, Atmel Corporation
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
#include "arch/at91_pio.h"
#include "gpio.h"

#include "debug.h"

#include "nand.h"
#include "pmecc.h"
#include "hamming.h"
#include "timer.h"
#include "fdt.h"
#include "div.h"

#ifdef CONFIG_NANDFLASH_SMALL_BLOCKS
static struct nand_chip nand_ids[] = {
	/* Samsung 32MB 8Bit */
	{0xec75, 0x800, 0x4000,  0x200, 0x10, 0x0},
	{0,}
};
#else
static struct nand_chip nand_ids[] = {
	/* Samsung K9F2G08U0M 256MB */
	{0xecda, 0x800, 0x20000, 0x800, 0x40, 0x0},
	/* Samsung K9F2G08U0A 256MB */
	{0xecaa, 0x800, 0x20000, 0x800, 0x40, 0x0},
	/* Micron MT29F2G16AAB 256MB */
	{0x2cca, 0x800, 0x20000, 0x800, 0x40, 0x1},
	/* Micron MT29F2G08AAC 256MB  */
	{0x2cda, 0x800, 0x20000, 0x800, 0x40, 0x0},
	/* Micron MT29F2G08ABD 256MB */
	{0x2caa, 0x800, 0x20000, 0x800, 0x40, 0x0},
	/* Micron MT29F4G08ABA 512MB */
	{0x2cdc, 0x800, 0x40000, 0x1000, 0xe0, 0x0},
	/* Mircon MT29H8G08ACAH1 1GB */
	{0x2c38, 0x800, 0x80000, 0x1000, 0xe0, 0x0},
#ifndef CONFIG_AT91SAM9260EK
	/* Hynix HY27UF082G2A 256MB */
	{0xadda, 0x800, 0x20000, 0x800, 0x40, 0x0},
	/* Hynix HY27UF162G2A 256MB */
	{0xadca, 0x800, 0x20000, 0x800, 0x40, 0x1},
	/* Hynix HY27UF162G2B 512MB */
	{0xaddc, 0x1000, 0x20000, 0x800, 0x40, 0x0},
	/* EON EN27LN1G08 128MB */
	{0x92f1, 0x400, 0x20000, 0x800, 0x40, 0x0},
#endif
	{0,}
};
#endif

/* ooblayout */
static struct nand_ooblayout nand_oob_layout;

static struct nand_chip nand_chip_default;

/*
 * NAND Commands
 */

/* 8 bits devices */
static void nand_command(unsigned char cmd)
{
	volatile unsigned long ioaddr = (unsigned long)CONFIG_SYS_NAND_BASE
						| CONFIG_SYS_NAND_MASK_CLE;

	writeb(cmd, ioaddr);
}

static void nand_address(unsigned char addr)
{
	volatile unsigned long ioaddr = (unsigned long)CONFIG_SYS_NAND_BASE
						| CONFIG_SYS_NAND_MASK_ALE;

	writeb(addr, ioaddr);
}

static unsigned char read_byte(void)
{
	return(readb((unsigned long)CONFIG_SYS_NAND_BASE));
}

/* 16 bits devices */
static void nand_command16(unsigned char cmd)
{
	volatile unsigned long ioaddr = (unsigned long)CONFIG_SYS_NAND_BASE
						| CONFIG_SYS_NAND_MASK_CLE;

	writew(cmd, ioaddr);
}

static void nand_address16(unsigned char addr)
{
	volatile unsigned long ioaddr = (unsigned long)CONFIG_SYS_NAND_BASE
						| CONFIG_SYS_NAND_MASK_ALE;

	writew(addr, ioaddr);
}

static unsigned short read_word(void)
{
	return(readw((unsigned long)CONFIG_SYS_NAND_BASE));
}

static void nand_wait_ready(void)
{
	unsigned int timeout = 10000;

	nand_command(CMD_STATUS);
	while ((!(read_byte() & STATUS_READY)) && timeout--)
		;
}

static void nand_cs_enable(void)
{
#ifdef CONFIG_SYS_NAND_ENABLE_PIN
	pio_set_value(CONFIG_SYS_NAND_ENABLE_PIN, 0);
#endif
}

static void nand_cs_disable(void)
{
#ifdef CONFIG_SYS_NAND_ENABLE_PIN
	pio_set_value(CONFIG_SYS_NAND_ENABLE_PIN, 1);
#endif
}

#ifdef CONFIG_NANDFLASH_SMALL_BLOCKS
static void config_nand_ooblayout(struct nand_ooblayout *layout,
				struct nand_info *nand,
				struct nand_chip *chip)
{
	layout->badblockpos = 5;
	layout->eccpos[0] = 0;
	layout->eccpos[1] = 1;
	layout->eccpos[2] = 2;

	switch (chip->pagesize) {
	case 256:
		layout->eccbytes = 3;
		break;

	case 512:
		layout->eccbytes = 6;
		layout->eccpos[3] = 3;
		layout->eccpos[4] = 6;
		layout->eccpos[5] = 7;
		break;

	default:
		break;
	}
}
#else
static void config_nand_ooblayout(struct nand_ooblayout *layout,
				struct nand_info *nand,
				struct nand_chip *chip)
{
	unsigned int i;
	unsigned int oobsize = chip->oobsize;
	layout->badblockpos = 0;

#ifdef CONFIG_USE_PMECC
	layout->eccbytes = div(chip->pagesize, nand->ecc_sector_size)
		* get_pmecc_bytes(nand->ecc_sector_size, nand->ecc_err_bits);
#else	/* Use Software ECC */
	switch (chip->pagesize) {
	case 2048:	/* oobsize is 64. */
		layout->eccbytes = 24;
		break;
	default:
		layout->eccbytes = 48;
		/* Software ECC support up to 128-byte oob */
		if (oobsize >= 128) {
			/* oobsize can be 224, 448. */
			dbg_info("NAND: Use software ECC." \
				"But may cannot correct all errors!\n");
			oobsize = 128;
		} else {
			dbg_info("NAND: Error! Software ECC not support pagesize: %d\n",
				chip->pagesize);
			return;
		}
		break;
	}
#endif
	for (i = 0; i < layout->eccbytes; i++)
		layout->eccpos[i] = oobsize - layout->eccbytes + i;
}
#endif /* #ifdef CONFIG_NANDFLASH_SMALL_BLOCKS */

#ifdef CONFIG_USE_ON_DIE_ECC_SUPPORT
static void write_byte(unsigned char data)
{
	writeb(data, (unsigned long)CONFIG_SYS_NAND_BASE);
}

static void nand_set_feature_on_die_ecc(unsigned char is_enable)
{
	unsigned char i;

	nand_cs_enable();

	nand_command(CMD_SET_FEATURE);
	nand_address(0x90);

	udelay(100);
	if (is_enable)
		write_byte(0x08);
	else
		write_byte(0x00);

	for (i = 0; i < 3; i++)
		write_byte(0x00);

	nand_wait_ready();
	nand_cs_disable();
}

static unsigned char nand_get_feature_on_die_ecc(void)
{
	unsigned char buffer[4];
	unsigned char i;

	nand_cs_enable();

	nand_command(CMD_GET_FEATURE);
	nand_address(0x90);
	nand_wait_ready();
	nand_command(CMD_READ_1);

	for (i = 0; i < 4; i++)
		buffer[i] = read_byte();

	nand_cs_disable();

	return buffer[0];
}

#define ENABLE_ECC	0x08

static int nand_set_on_die_ecc(unsigned char is_enable)
{
	unsigned char status;

	status = (nand_get_feature_on_die_ecc() & ENABLE_ECC) ? 1 : 0;
	if (is_enable == status)
		return 0;

	nand_set_feature_on_die_ecc(is_enable);
	status = (nand_get_feature_on_die_ecc() & ENABLE_ECC) ? 1 : 0;
	if (is_enable == status)
		return 0;

	return -1;
}

static int nand_init_on_die_ecc(void)
{
	unsigned char is_enable;

#ifdef CONFIG_ON_DIE_ECC
	is_enable = 1;
#else
	is_enable = 0;
#endif
	if (nand_set_on_die_ecc(is_enable)) {
		dbg_info("WARNING: Fail to %s On-Die ECC\n",
				is_enable ? "enable" : "disable");

	} else {
		dbg_info("NAND: %s On-Die ECC\n",
				is_enable ? "Enable" : "Disable");
	}

	return 0;
}
#endif /* #ifdef CONFIG_USE_ON_DIE_ECC_SUPPORT */

static void nandflash_read_id(unsigned char *manf_id, unsigned char *dev_id)
{
	nand_cs_enable();

	nand_command(CMD_READID);
	nand_address(0x00);

	*manf_id = read_byte();
	*dev_id = read_byte();

	nand_cs_disable();
}

#ifdef CONFIG_ONFI_DETECT_SUPPORT
static unsigned short onfi_crc16(unsigned short crc,
				unsigned char const *p,
				unsigned int len)
{
	int i;

	while (len--) {
		crc ^= *p++ << 8;
		for (i = 0; i < 8; i++)
			crc = (crc << 1) ^ ((crc & 0x8000) ? 0x8005 : 0);
	}

	return crc;
}

#define ONFI_PARAMS_SIZE		256

#define PARAMS_OFFSET_REVISION		4
#define		PARAMS_REVISION_1_0	(0x1 << 1)
#define		PARAMS_REVISION_2_0	(0x1 << 2)
#define		PARAMS_REVISION_2_1	(0x1 << 3)

#define PARAMS_OFFSET_FEATURES		6
#define		PARAMS_FEATURE_BUSWIDTH		(0x1 << 0)
#define		PARAMS_FEATURE_EXTENDED_PARAM	(0x1 << 7)

#define PARAMS_OFFSET_EXT_PARAM_PAGE_LEN	12
#define PARAMS_OFFSET_PARAMETER_PAGE		14
#define PARAMS_OFFSET_PAGESIZE		80
#define PARAMS_OFFSET_OOBSIZE		84
#define PARAMS_OFFSET_BLOCKSIZE		92
#define PARAMS_OFFSET_NBBLOCKS		96
#define PARAMS_OFFSET_ECC_BITS		112
#define PARAMS_OFFSET_CRC		254

#define ONFI_CRC_BASE			0x4F4E

#define ONFI_MAX_SECTIONS		8

#define ONFI_SECTION_TYPE_0		0
#define ONFI_SECTION_TYPE_1		1
#define ONFI_SECTION_TYPE_2		2

static int get_ext_onfi_param(unsigned char *eccbits,
			      unsigned int *eccwordsize,
			      unsigned int len)
{
	unsigned char ext_params[ONFI_PARAMS_SIZE];
	unsigned char *param = ext_params;
	unsigned char *p = ext_params;
	unsigned char *section, *table;
	unsigned short crc;
	unsigned int i;

	for (i = 0; i < len; i++)
		*p++ = read_byte();

	crc = *(unsigned short *)(param);
	if (onfi_crc16(ONFI_CRC_BASE,
		       (unsigned char *)(param + 2), len - 2) != crc) {
		dbg_info("NAND: Failed in the integrity CRC\n");
		return -1;
	}

	if ((param[2] != 'E') || (param[3] != 'P') ||
	    (param[4] != 'P') || (param[5] != 'S')) {
		dbg_info("NAND: No Extended Parameter Page\n");
		return -1;
	}

	section = param + 16;
	table =  param + 32;
	for (i = 0; i < ONFI_MAX_SECTIONS; i++) {
		if ((*section) == ONFI_SECTION_TYPE_2)
			break;
		table += *(section + 1) * 16;
		section += 2;
	}

	if (i == ONFI_MAX_SECTIONS) {
		dbg_info("NAND: Not find the ECC section\n");
		return -1;
	}

	*eccbits = *table;
	*eccwordsize = 0x01 << *(table + 1);

	return 0;
}

static int nandflash_detect_onfi(struct nand_chip *chip)
{
	unsigned char onfi_ind[4];
	unsigned char onfi_params[ONFI_PARAMS_SIZE];
	unsigned char *p = onfi_params;
	unsigned char *param;
	int i, j;
	unsigned short crc;
	unsigned char manf_id, dev_id;
	unsigned short features, revision, ext_page_len;
	unsigned char num_param_page;

	nand_cs_enable();
	nand_command(CMD_READID);
	nand_address(0x20);

	onfi_ind[0] = read_byte();
	onfi_ind[1] = read_byte();
	onfi_ind[2] = read_byte();
	onfi_ind[3] = read_byte();
	nand_cs_disable();

	if ((onfi_ind[0] != 'O')
		|| (onfi_ind[1] != 'N')
		|| (onfi_ind[2] != 'F')
		|| (onfi_ind[3] != 'I')) {
		dbg_info("NAND: ONFI not supported\n");
		return -1;
	}

	dbg_info("NAND: ONFI flash detected\n");

	nand_cs_enable();

	nand_command(CMD_READ_ONFI);
	nand_address(0x00);
	
	nand_wait_ready();
	nand_command(CMD_READ_1);
	
	for (i = 0; i < 3; i++) {
		param = onfi_params;
		/* Read the onfi parameter table */
		for (j = 0; j < ONFI_PARAMS_SIZE; j++)
			*param++ = read_byte();

		crc = *(unsigned short *)(p + PARAMS_OFFSET_CRC);
		if (onfi_crc16(ONFI_CRC_BASE, (unsigned char *)p, 254) == crc)
			break;
	}

	if (i == 3) {
		dbg_info("NAND: ONFI para CRC error!\n");
		nand_cs_disable();
		return -1;
	}

	revision = *(unsigned short *)(p + PARAMS_OFFSET_REVISION);
	features = *(unsigned short *)(p + PARAMS_OFFSET_FEATURES);
	ext_page_len = *(unsigned short *)(p +
					   PARAMS_OFFSET_EXT_PARAM_PAGE_LEN);
	num_param_page = *(unsigned char *)(p + PARAMS_OFFSET_PARAMETER_PAGE);

	chip->numblocks = *(unsigned short *)(p + PARAMS_OFFSET_NBBLOCKS);
	chip->pagesize	= *(unsigned short *)(p + PARAMS_OFFSET_PAGESIZE);
	chip->blocksize = *(unsigned int  *)(p + PARAMS_OFFSET_BLOCKSIZE)
							* chip->pagesize;
	chip->oobsize	= *(unsigned short *)(p + PARAMS_OFFSET_OOBSIZE);
	chip->buswidth	= features & PARAMS_FEATURE_BUSWIDTH;
	chip->eccbits	= *(unsigned char *)(p + PARAMS_OFFSET_ECC_BITS);
	chip->eccwordsize = 512;

	if ((chip->eccbits == 0xff) &&
	    (revision & PARAMS_REVISION_2_1) &&
	    (features & PARAMS_FEATURE_EXTENDED_PARAM)) {
		/* read the redundent parameter pages */
		for (; i < (num_param_page - 1) * ONFI_PARAMS_SIZE; i++)
			read_byte();

		/* read extended parameter pages */
		if (get_ext_onfi_param(&chip->eccbits,
				       &chip->eccwordsize,
				       ext_page_len * 16)) {
			dbg_info("NAND: Failed to get extended parameter table\n");
			return -1;
		}
	}

	nand_cs_disable();

	nandflash_read_id(&manf_id, &dev_id);

	dbg_info("NAND: Manufacturer ID: %x Chip ID: %x\n", manf_id, dev_id);
	dbg_info("NAND: Page Bytes: %d, Spare Bytes: %d\n" \
		 "NAND: ECC Correctability Bits: %d, ECC Sector Bytes: %d\n",
		 chip->pagesize, chip->oobsize,
		 chip->eccbits, chip->eccwordsize);

	return 0;
}
#endif /* #ifdef CONFIG_ONFI_DETECT_SUPPORT */

static int nandflash_detect_non_onfi(struct nand_chip *chip)
{
	unsigned char manf_id, dev_id;
	unsigned int chipid;
	unsigned int i;

	nandflash_read_id(&manf_id, &dev_id);

	chipid = ((unsigned int)manf_id << 8) | dev_id;

	for (i = 0; i < ARRAY_SIZE(nand_ids); i++) {
		if (chipid == nand_ids[i].chip_id)
			break;
	}

	if (i == ARRAY_SIZE(nand_ids)) {
		dbg_info("NAND: Not found Manufacturer ID: %x," \
			"Chip ID: %x\n", manf_id, dev_id);

		return -1;
	}

	dbg_info("NAND: Manufacturer ID: %x Chip ID: %x\n",
						manf_id, dev_id);

	chip->pagesize	= nand_ids[i].pagesize;
	chip->blocksize = nand_ids[i].blocksize;
	chip->oobsize	= nand_ids[i].oobsize;
	chip->buswidth	= nand_ids[i].buswidth;
	chip->numblocks = nand_ids[i].numblocks;

	return 0;
}

static int nand_info_init(struct nand_info *nand, struct nand_chip *chip)
{
	/* number of blocks in device */
	nand->numblocks = chip->numblocks;
	/* number of data bytes in a block */
	nand->blocksize = chip->blocksize;
	/* number of bytes in page area */
	nand->pagesize = chip->pagesize;
	/* number of pages in block */
	nand->pages_block = div(nand->blocksize, nand->pagesize);
	/* number of pages in device */
	nand->pages_device = nand->numblocks * nand->pages_block;
	/* number of bytes in oob area */
	nand->oobsize = chip->oobsize;
	/* Total number of bytes in a sector */
	nand->sectorsize = nand->pagesize + nand->oobsize;
#ifdef CONFIG_USE_PMECC
	if (choose_pmecc_info(nand, chip))
		return -1;
#endif
	/* the layout of the spare area */
	config_nand_ooblayout(&nand_oob_layout, nand, chip);
	nand->ecclayout = &nand_oob_layout;
	/* data bus width (8/16 bits) */
	nand->buswidth = chip->buswidth;
	if (nand->buswidth) {
		nand->ecclayout->badblockpos *= 2;
		nand->command = nand_command16;
		nand->address = nand_address16;
	} else {
		nand->command = nand_command;
		nand->address = nand_address;
	}

	return 0;
}

static void nandflash_reset(void)
{
	nand_cs_enable();

	nand_command(CMD_RESET);

	nand_wait_ready();

	nand_cs_disable();
}

static int nandflash_get_type(struct nand_info *nand)
{
	struct nand_chip *chip = &nand_chip_default;

	nandflash_reset();

#ifdef CONFIG_ONFI_DETECT_SUPPORT
	int ret;

	/* Check if the Nandflash is ONFI compliant */
	ret = nandflash_detect_onfi(chip);
	if (ret == -1) {
		if (nandflash_detect_non_onfi(chip)) {
			dbg_info("NAND: Not find support device!\n");
			return -1;
		}
	}
#else
	if (nandflash_detect_non_onfi(chip)) {
		dbg_info("NAND: Not find support device!\n");
		return -1;
	}
#endif

#ifdef CONFIG_USE_ON_DIE_ECC_SUPPORT
	if (nand_init_on_die_ecc())
		return -1;
#endif

	return nand_info_init(nand, chip);
}

static void write_column_address(struct nand_info *nand,
				unsigned int column_address)
{
	volatile unsigned int page_size = nand->pagesize;

	if (nand->buswidth)
		column_address >>= 1;

	while (page_size > 2) {
		nand->address(column_address & 0xff);

		page_size >>= 8;
		column_address >>= 8;
	}
}

static void write_row_address(struct nand_info *nand,
				unsigned int row_address)
{
	volatile unsigned int num_pages = nand->pages_device;

	while(num_pages) {
		nand->address(row_address & 0xff);

		num_pages >>= 8;
		row_address >>= 8;
	}
}

static int nand_read_status(void)
{
	unsigned int timeout = 1000;
	unsigned char status;

	do {
		nand_command(CMD_STATUS);
		status = read_byte();
		if (status & STATUS_READY)
			break;
	} while (--timeout);

	if (!timeout)
		return -1;

#ifdef CONFIG_ON_DIE_ECC
	if (status & STATUS_ERROR) {
		dbg_info("WARNING: Read On-Die ECC error\n");
		return -1;
	}
#endif

	return 0;
}

#ifdef CONFIG_NANDFLASH_SMALL_BLOCKS
static int nand_read_sector(struct nand_info *nand, 
			unsigned int row_address,
			unsigned char *buffer,
			unsigned int zone_flag)
{
	unsigned int readbytes, i;
	unsigned int column_address;
	unsigned char command;

	switch (zone_flag) {
	case ZONE_DATA:
		readbytes = nand->pagesize;
		column_address = 0x00;
		command = CMD_READ_A0;
		break;

	case ZONE_INFO:
		readbytes = nand->oobsize;
		column_address = nand->pagesize;
		buffer += nand->pagesize;
		command = CMD_READ_C;
		break;

	case ZONE_DATA | ZONE_INFO:
		readbytes = nand->sectorsize;
		column_address = 0x00;
		command = CMD_READ_A0;
		break;

	default:
		return -1;
	}

	nand_cs_enable();

	/* Write specific command, Read from start */
	nand->command(command);

	write_column_address(nand, column_address);
	write_row_address(nand, row_address);

	if (nand_read_status())
		return -1;

	nand->command(CMD_READ_A0);

	/* Read loop */
	if (nand->buswidth) {
		for (i = 0; i < readbytes / 2; i++) {
			*((short *)buffer) = read_word();
			buffer += 2;
		}
	} else {
		for (i = 0; i < readbytes; i++) {
			*buffer = read_byte();
			buffer++;
		}
	}

	nand_cs_disable();

	return 0;
}
#else /* large blocks */
static int nand_read_sector(struct nand_info *nand,
				unsigned int row_address,
				unsigned char *buffer, 
				unsigned int zone_flag)
{
	unsigned int readbytes, i;
	unsigned int column_address;
	int ret = 0;
	unsigned char *pbuf = buffer;

#ifdef CONFIG_USE_PMECC
	unsigned int usepmecc = 0;

	if ((zone_flag & ZONE_DATA) == ZONE_DATA) {
		usepmecc = 1;
		zone_flag = ZONE_DATA | ZONE_INFO;

		pmecc_enable();
	}
#endif	/* #ifdef CONFIG_USE_PMECC */

	switch (zone_flag) {
	case ZONE_DATA:
		readbytes = nand->pagesize;
		column_address = 0x00;
		break;

	case ZONE_INFO:
		readbytes = nand->oobsize;
		pbuf += nand->pagesize;
		column_address = nand->pagesize;
		break;

	case ZONE_DATA | ZONE_INFO:
		readbytes = nand->sectorsize;
		column_address = 0x00;
		break;

	default:
		return -1;
	}

	nand_cs_enable();

	nand->command(CMD_READ_1);

	write_column_address(nand, column_address);
	write_row_address(nand, row_address);

	nand->command(CMD_READ_2);

	if (nand_read_status())
		return -1;

	nand->command(CMD_READ_1);

#ifdef CONFIG_USE_PMECC
	if (usepmecc)
		pmecc_start_data_phase();
#endif
	/* Read loop */
	if (nand->buswidth) {
		for (i = 0; i < readbytes / 2; i++) {
			*((short *)pbuf) = read_word();
			pbuf += 2;
		}
	} else {
		for (i = 0; i < readbytes; i++)
			*pbuf++ = read_byte();

#ifdef CONFIG_USE_PMECC
		if (usepmecc)
			ret = pmecc_process(nand, buffer);
#endif
	}

	nand_cs_disable();

	return ret;
}
#endif /* #ifdef CONFIG_NANDFLASH_SMALL_BLOCKS */

static int nand_check_badblock(struct nand_info *nand,
				unsigned int block,
				unsigned char *buffer)
{
	unsigned int page;
	unsigned int row_address = block * nand->pages_block;

	/*
	 * Read the first page and second page oob zone
	 * to detect if block is bad
	 */
	for (page = 0; page < 2; page++) {
		nand_read_sector(nand, row_address + page, buffer, ZONE_INFO);
		if (*(buffer + nand->pagesize + nand->ecclayout->badblockpos)
			!= 0xff)
			return -1;
	}

	return 0;
}

#ifdef CONFIG_ENABLE_SW_ECC
static void nand_read_ecc(struct nand_ooblayout *ooblayout,
				unsigned char *buffer,
				unsigned char *ecc)
{
	unsigned int i;

	for (i = 0; i < ooblayout->eccbytes; i++)
		ecc[i] = buffer[ooblayout->eccpos[i]];
}
#endif

static int nand_read_page(struct nand_info *nand,
				unsigned int block,
				unsigned int page,
				unsigned int zone_flag,
				unsigned char *buffer)
{
	unsigned int row_address = block * nand->pages_block + page;

#ifndef CONFIG_ENABLE_SW_ECC
	return nand_read_sector(nand, row_address, buffer, ZONE_DATA);
#else

	int retval;
	unsigned char hamming[48], error;

	retval = nand_read_sector(nand, row_address, buffer,
				ZONE_DATA | ZONE_INFO);
	if (retval)
		return -1;

	nand_read_ecc(nand->ecclayout, buffer + nand->pagesize, hamming);

	error = Hamming_Verify256x(buffer, nand->pagesize, hamming);
	if (error && (error != Hamming_ERROR_SINGLEBIT)) {
		dbg_info("NAND: Hamming ECC error!\n");
		return -1;
	}

	return 0;
#endif /* #ifndef CONFIG_ENABLE_SW_ECC */
}

#ifdef CONFIG_NANDFLASH_RECOVERY
static int nand_erase_block0(struct nand_info *nand)
{
	unsigned int row_address = 0;
	unsigned int timeout = 10000;
	unsigned int status;

	nand_cs_enable();

	nand_command(CMD_ERASE_1);
	write_row_address(nand, row_address);
	nand_command(CMD_ERASE_2);

	udelay(2000);

	nand_command(CMD_STATUS);
	while ((!((status = read_byte()) & STATUS_READY)) && --timeout)
		;

	nand_cs_disable();

	if (status & STATUS_ERROR)
		return -1;

	if (timeout == 0)
		return -2;

	return 0;
}

static int nandflash_recovery(struct nand_info *nand)
{
	int ret = -1;

	/*
	 * If Recovery Button is pressed during boot sequence,
	 * erase nandflash block0
	*/
	dbg_info("NAND: Press the recovery button (%s) to recovery\n",
			RECOVERY_BUTTON_NAME);

	if ((pio_get_value(CONFIG_SYS_RECOVERY_BUTTON_PIN)) == 0) {
		dbg_info("NAND: The recovery button (%s) has been "\
				"pressed\n", RECOVERY_BUTTON_NAME);
		dbg_info("NAND: The block 0 is erasing ...\n");

		ret = nand_erase_block0(nand);
		if (ret)
			dbg_info("NAND: The erasing failed\n");
		else
			dbg_info("NAND: The erasing is done\n");
	}

	return ret;
}
#endif /* #ifdef CONFIG_NANDFLASH_RECOVERY */

static int nand_loadimage(struct nand_info *nand,
				unsigned int offset,
				unsigned int length,
				unsigned char *dest)
{
	unsigned char *buffer = dest;
	unsigned int readsize;
	unsigned int block = 0;
	unsigned int page;
	unsigned int start_page = 0;
	unsigned int end_page;
	unsigned int numpages = 0;
	unsigned int offsetpage = 0;
	unsigned int block_remaining = nand->blocksize
				       - mod(offset, nand->blocksize);
	int ret;

	division(offset, nand->blocksize, &block, &start_page);
	start_page = div(start_page, nand->pagesize);

	while (length > 0) {
		/* read a buffer corresponding to a block */
		if (length < block_remaining)
			readsize = length;
		else
			readsize = block_remaining;

		/* adjust the number of pages to read */
		division(readsize, nand->pagesize, &numpages, &offsetpage);
		if (offsetpage)
			numpages++;

		end_page = start_page + numpages;

		/* check the bad block */
		while (1) {
			if (nand_check_badblock(nand,
					block, buffer) != 0) {
				block++; /* skip this block */
				dbg_info("NAND: Bad block:" \
					" #%x\n", block);
			} else
				break;
		}

		/* read pages of a block */
		for (page = start_page; page < end_page; page++) {
			ret = nand_read_page(nand, block, page,
						ZONE_DATA, buffer);
			if (ret)
				return -1;
			else
				buffer += nand->pagesize;
		}
		length -= readsize;

		block++;
		start_page = 0;
		block_remaining = nand->blocksize;
	}

	return 0;
}

#if defined(CONFIG_LOAD_LINUX) || defined(CONFIG_LOAD_ANDROID)
static int update_image_length(struct nand_info *nand,
				unsigned int offset,
				unsigned char *dest,
				unsigned char flag)
{
	unsigned int length = nand->pagesize;
	int ret;

	ret = nand_loadimage(nand, offset, length, dest);
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

int load_nandflash(struct image_info *image)
{
	struct nand_info nand;
	int ret;

	nandflash_hw_init();

	if (nandflash_get_type(&nand))
		return -1;

#ifdef CONFIG_NANDFLASH_RECOVERY
	if (nandflash_recovery(&nand) == 0)
		return -2;
#endif

#ifdef CONFIG_USE_PMECC
	if (init_pmecc(&nand))
		return -1;
#endif

#ifdef CONFIG_ENABLE_SW_ECC
	dbg_info("NAND: Using Software ECC\n");
#endif

#if defined(CONFIG_LOAD_LINUX) || defined(CONFIG_LOAD_ANDROID)
	int length = update_image_length(&nand,
				image->offset, image->dest, KERNEL_IMAGE);
	if (length == -1)
		return -1;

	image->length = length;
#endif

	dbg_info("NAND: Image: Copy %x bytes from %x to %x\n",
			image->length, image->offset, image->dest);

	ret = nand_loadimage(&nand, image->offset, image->length, image->dest);
	if (ret)
		return ret;

#ifdef CONFIG_OF_LIBFDT
	length = update_image_length(&nand,
			image->of_offset, image->of_dest, DT_BLOB);
	if (length == -1)
		return -1;

	image->of_length = length;

	dbg_info("NAND: dt blob: Copy %x bytes from %x to %x\n",
		image->of_length, image->of_offset, image->of_dest);

	ret = nand_loadimage(&nand, image->of_offset,
				image->of_length, image->of_dest);
	if (ret)
		return ret;
#endif

	return 0;
 }

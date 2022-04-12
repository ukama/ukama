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
#include "common.h"
#include "hardware.h"
#include "board.h"
#include "spi.h"
#include "arch/at91_pio.h"
#include "gpio.h"
#include "string.h"
#include "timer.h"
#include "div.h"
#include "fdt.h"
#include "debug.h"

/* Manufacturer Device ID Read */
#define CMD_READ_DEV_ID			0x9f
/* Continuous Array Read */
#define CMD_READ_ARRAY_FAST		0x0b
#define CMD_READ_ARRAY		0x03

/* JEDEC Code */
#define MANUFACTURER_ID_ATMEL		0x1f
#define MANUFACTURER_ID_MICRON		0x20
#define MANUFACTURER_ID_WINBOND		0xef

/* Family Code */
#define DF_FAMILY_AT26F			0x00
#define DF_FAMILY_AT45			0x20
#define DF_FAMILY_AT26DF		0x40	/* AT25DF and AT26DF */

#define DF_FAMILY_N25Q			0xA0
#define DF_FAMILY_M25P			0x20

/* AT45 Density Code */
#define DENSITY_AT45DB011D		0x0C
#define DENSITY_AT45DB021D		0x14
#define DENSITY_AT45DB041D		0x1C
#define DENSITY_AT45DB081D		0x24
#define DENSITY_AT45DB161D		0x2C
#define DENSITY_AT45DB321D		0x34
#define DENSITY_AT45DB642D		0x3C
#define DENSITY_AT45DB1282D		0x10
#define DENSITY_AT45DB2562D		0x18
#define DENSITY_AT45DB5122D		0x20

/* AT45 Status Register Read */
#define CMD_READ_STATUS_AT45		0xd7

/* AT45 status register bits */
#define STATUS_PAGE_SIZE_AT45		(1 << 0)
#define STATUS_READY_AT45		(1 << 7)

struct dataflash_descriptor {
	unsigned char	family;

	unsigned int	pages;		/* page number */
	unsigned int	page_size;	/* page size */
	unsigned int	page_offset;	/* page offset in command */
	unsigned char	is_power_2;	/* = 1: power of 2, = 0: not*/
	unsigned char	is_spinor;	/* = 1: nor flash, = 0: dataflash */
};

static int df_send_command(unsigned char *cmd,
				unsigned char cmd_len,
				unsigned char *data,
				unsigned int data_len)
{
	int i;

	if (!cmd)
		return -1;

	if (!cmd_len)
		return -1;

	if (data_len)
		if (!data)
			return -1;

	at91_spi_cs_activate();

	/* read spi status to clear events */
	at91_spi_read_sr();

	for (i = 0; i < cmd_len; i++) {
		at91_spi_write_data(*cmd++);
		at91_spi_read_spi();
	}

	for (i = 0; i < data_len; i++) {
		at91_spi_write_data(0);
		*data++ = at91_spi_read_spi();
	}
	at91_spi_cs_deactivate();

	return 0;
}

static int dataflash_read_array(struct dataflash_descriptor *df_desc,
				unsigned int offset,
				unsigned int len,
				void *buf)
{
	unsigned char cmd[5];
	unsigned char cmd_len;
	unsigned int address;
	unsigned int page_addr = 0;
	unsigned int byte_addr = 0;
	unsigned int page_shift;
	unsigned int page_size;
	int ret;

	if (!df_desc->is_power_2) {
		page_shift = df_desc->page_offset;
		page_size = df_desc->page_size;

		division(offset, page_size, &page_addr, &byte_addr);

		address = (page_addr << page_shift) + byte_addr;
	} else
		address = offset;

	cmd[0] = CMD_READ_ARRAY_FAST;
	if (df_desc->pages > 16384) {
		cmd[1] = (unsigned char)(address >> 24);
		cmd[2] = (unsigned char)(address >> 16);
		cmd[3] = (unsigned char)(address >> 8);
		cmd[4] = (unsigned char)address;

	} else {
		cmd[1] = (unsigned char)(address >> 16);
		cmd[2] = (unsigned char)(address >> 8);
		cmd[3] = (unsigned char)address;
		cmd[4] = 0x00;
	}

	cmd_len = 5;

	ret = df_send_command(cmd, cmd_len, buf, len);
	if (ret)
		return -1;

	return 0;
}

static int spinor_read_array(struct dataflash_descriptor *df_desc,
				unsigned int offset,
				unsigned int len,
				void *buf)
{
	unsigned char cmd[5];
	unsigned char cmd_len;
	unsigned int address;
	int ret;

	address = offset;

	cmd[0] = CMD_READ_ARRAY_FAST;
	cmd[1] = (unsigned char)(address >> 16);
	cmd[2] = (unsigned char)(address >> 8);
	cmd[3] = (unsigned char)address;
	cmd_len = 5; /* 5th command is for dummy cycle*/

	ret = df_send_command(cmd, cmd_len, buf, len);
	if (ret)
		return -1;

	return 0;
}

static int read_array(struct dataflash_descriptor *df_desc,
				unsigned int offset,
				unsigned int len,
				void *buf)
{
	if (!df_desc->is_spinor)
		return dataflash_read_array(df_desc, offset, len, buf);
	else
		return spinor_read_array(df_desc, offset, len, buf);
}

#if defined(CONFIG_LOAD_LINUX) || defined(CONFIG_LOAD_ANDROID)
static int update_image_length(struct dataflash_descriptor *df_desc,
				unsigned int offset,
				unsigned char *dest,
				unsigned char flag)
{
	unsigned int length = df_desc->page_size;
	int ret;

	ret = read_array(df_desc, offset, length, dest);
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

static unsigned char df_read_status_at45(unsigned char *status)
{
	unsigned char cmd = CMD_READ_STATUS_AT45;
	int ret;

	ret = df_send_command(&cmd, 1, status, 1);
	if (ret)
		return ret;

	return 0;
}

#ifdef CONFIG_DATAFLASH_RECOVERY

/* AT25 Block Erase(4-KBytes) Command*/
#define CMD_ERASE_BLOCK4K_AT25		0x20
/* Write Enable Command */
#define CMD_WRITE_ENABLE_AT25		0x06
/* Status Register Commands */
#define CMD_READ_STATUS_AT25		0x05
#define CMD_WRITE_STATUS_AT25           0x01
/* Page Erase AT45 */
#define CMD_ERASE_PAGE_AT45		0x81

/* AT25 status register bits */
#define STATUS_READY_AT25		(1 << 0)
#define STATUS_WEL_AT25			(1 << 1)
#define STATUS_SWP_AT25			(3 << 2)
#define STATUS_EPE_AT25			(1 << 5)
#define STATUS_SPRL_AT25		(1 << 7)

static unsigned char df_read_status_at25(unsigned char *status)
{
	unsigned char cmd = CMD_READ_STATUS_AT25;
	int ret;

	ret = df_send_command(&cmd, 1, status, 1);
	if (ret)
		return ret;

	return 0;
}

static int at25_cmd_write_enbale(void)
{
	unsigned char cmd;
	int ret;

	cmd = CMD_WRITE_ENABLE_AT25;
	ret = df_send_command(&cmd, 1, NULL, 0);
	if (ret)
		return ret;

	return 0;
}

static int at25_cmd_write_status_register(unsigned char status)
{
	unsigned char cmd[2];
	int ret;

	cmd[0] = CMD_WRITE_STATUS_AT25;
	cmd[1] = status;

	ret = df_send_command(cmd, 2, NULL, 0);
	if (ret)
		return ret;

	return 0;
}

static int at25_unprotect(void)
{
	unsigned char status;
	int ret;

	/* read status register */
	ret = df_read_status_at25(&status);
	if (ret)
		return ret;

	/* check if All sectors are software unprotected
	 * (all Sector Protection Register are 0)
	 */
	if (!(status & STATUS_SWP_AT25))
		return 0;

	/* check if Sector Protection Registers are locked */
	if (status & STATUS_SPRL_AT25) {
		/* Unprotect Sector Potection Registers. */
		ret = at25_cmd_write_enbale();
		if (ret)
			return ret;

		ret = at25_cmd_write_status_register(0);
		if (ret)
			return ret;
	}

	/* a global unprotect command */
	ret = at25_cmd_write_enbale();
	if (ret)
		return ret;

	ret = at25_cmd_write_status_register(0);
	if (ret)
		return ret;

	/* check Status Register SPRL & SWP bits */
	ret = df_read_status_at25(&status);
	if (ret)
		return ret;

	if (status & (STATUS_SPRL_AT25 | STATUS_SWP_AT25)) {
		dbg_info("SF: Unprotect AT25 failed\n");
		return -1;
	}

	return 0;
}

static int dataflash_page0_erase_at25(void)
{
	unsigned char status;
	unsigned char cmd[5];
	unsigned int timeout = 1000;
	int ret;

	ret = at25_unprotect();
	if (ret)
		return ret;

	ret = at25_cmd_write_enbale();
	if (ret)
		return ret;

	/* Erase page0 */
	cmd[0] = CMD_ERASE_BLOCK4K_AT25;
	cmd[1] = 0;
	cmd[2] = 0;
	cmd[3] = 0;

	ret = df_send_command(cmd, 4, NULL, 0);
	if (ret) {
		dbg_info("SF: AT25 page 0 erase failed\n");
		return ret;
	}

	udelay(33000); /* 33 ms: the maximum delay of udelay() */

	do {
		ret = df_read_status_at25(&status);
		if (ret)
			return ret;

		if (!(status & STATUS_READY_AT25))
			break;
	} while (--timeout);

	if (!timeout) {
		dbg_info("SF: AT25 page0 erase timed out\n");
		return -1;
	}

	return 0;
}

static int dataflash_page0_erase_at45(void)
{
	unsigned char status;
	unsigned char cmd[4];
	unsigned int timeout = 1000;
	int ret;

	cmd[0] = CMD_ERASE_PAGE_AT45;
	cmd[1] = 0;
	cmd[2] = 0;
	cmd[3] = 0;

	ret = df_send_command(cmd, 4, NULL, 0);
	if (ret) {
		dbg_info("SF: AT45 page 0 erase failed\n");
		return ret;
	}

	udelay(33000); /* 33 ms: the maximum delay of udelay() */

	do {
		ret = df_read_status_at45(&status);
		if (ret)
			return ret;

		if (status & STATUS_READY_AT45)
			break;
	} while (--timeout);

	if (!(status & STATUS_READY_AT45)) {
		dbg_info("SF: AT45 page0 erase timed out\n");
		return -1;
	}

	return 0;
}

static int dataflash_recovery(struct dataflash_descriptor *df_desc)
{
	int ret;

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

		if ((df_desc->family == DF_FAMILY_AT26F)
			|| (df_desc->family == DF_FAMILY_AT26DF))
			ret = dataflash_page0_erase_at25();
		 else
			ret = dataflash_page0_erase_at45();

		if (ret) {
			dbg_info("SF: The erasing failed\n");
			return ret;
		}
		dbg_info("SF: The erasing is done\n");

		return 0;
	}

	return -1;
}
#endif /* #ifdef CONFIG_DATAFLASH_RECOVERY */

static int df_n25q_desc_init(struct dataflash_descriptor *df_desc)
{
	df_desc->pages = 16384;
	df_desc->page_size = 256;
	df_desc->page_offset = 0;
	df_desc->is_spinor = 1;
	return 0;
}

static int df_at45_desc_init(struct dataflash_descriptor *df_desc)
{
	unsigned char status;
	unsigned char density;
	int ret;

	ret = df_read_status_at45(&status);
	if (ret)
		return ret;

	if (status & STATUS_PAGE_SIZE_AT45)
		df_desc->is_power_2 = 1;
	else
		df_desc->is_power_2 = 0;

	density = status & 0x3c;
	switch (density) {
	case DENSITY_AT45DB011D:
		df_desc->pages = 512;
		df_desc->page_size = 264;
		df_desc->page_offset = 9;
	break;

	case DENSITY_AT45DB021D:
		df_desc->pages = 1024;
		df_desc->page_size = 264;
		df_desc->page_offset = 9;
		break;

	case DENSITY_AT45DB041D:
		df_desc->pages = 2048;
		df_desc->page_size = 264;
		df_desc->page_offset = 9;
		break;

	case DENSITY_AT45DB081D:
		df_desc->pages = 4096;
		df_desc->page_size = 264;
		df_desc->page_offset = 9;
		break;

	case DENSITY_AT45DB161D:
		df_desc->pages = 4096;
		df_desc->page_size = 528;
		df_desc->page_offset = 10;
		break;

	case DENSITY_AT45DB321D:
		df_desc->pages = 8192;
		df_desc->page_size = 528;
		df_desc->page_offset = 10;
		break;

	case DENSITY_AT45DB642D:
		df_desc->pages = 8192;
		df_desc->page_size = 1056;
		df_desc->page_offset = 11;
		break;
/*
	case DENSITY_AT45DB1282D:
		df_desc->pages = 16384;
		df_desc->pages_size = 1056;
		df_desc->page_offset = 11;
		break;

	case DENSITY_AT45DB2562D:
		df_desc->pages = 16384;
		df_desc->page_size = 2112;
		df_desc->page_offset = 12;
		break;

	case DENSITY_AT45DB5122D:
		df_desc->pages = 32768;
		df_desc->page_size = 2112;
		df_desc->page_offset = 12;
		break;
*/
	default:
		return -1;
	}

	return 0;
}

static int df_at25_desc_init(struct dataflash_descriptor *df_desc)
{
	/* AT25DF321 */
	df_desc->is_power_2 = 1;

	df_desc->pages = 16384;
	df_desc->page_size = 256;
	df_desc->page_offset = 0;

	return 0;
}

static int df_desc_init(struct dataflash_descriptor *df_desc, unsigned char vendor,
			unsigned char family)
{
	int ret;

	df_desc->family = family;

	switch ( vendor ) {
		case MANUFACTURER_ID_ATMEL: {

			if ((df_desc->family == DF_FAMILY_AT26F)
				|| (df_desc->family == DF_FAMILY_AT26DF)) {
				ret = df_at25_desc_init(df_desc);
				if (ret)
					return ret;
			} else if (df_desc->family == DF_FAMILY_AT45) {
				ret = df_at45_desc_init(df_desc);
				if (ret)
					return ret;
			} else {
				dbg_info("SF: Unsupported SerialFlash family %x\n", family);
				return -1;
			}
		}
		break;

		case MANUFACTURER_ID_MICRON:{

			if (df_desc->family == DF_FAMILY_M25P) {
				ret = df_n25q_desc_init(df_desc);
				if (ret)
					return ret;
			}else if (df_desc->family == DF_FAMILY_N25Q) {
				ret = df_n25q_desc_init(df_desc);
				if (ret)
					return ret;
			} else {
				dbg_info("SF: Unsupported SerialFlash family %x\n", family);
				return -1;
			}
		}
		break;

		default:
			dbg_info("SF: Unsupported Manufactorer ID %x\n", vendor);
			return -1;
	}

	return 0;
}

static int dataflash_probe_atmel(struct dataflash_descriptor *df_desc)
{
	unsigned char dev_id[5];
	unsigned char cmd = CMD_READ_DEV_ID;
	int ret;

	/* Read device ID */
	ret = df_send_command(&cmd, 1, dev_id, 5);
	if (ret)
		return ret;

#ifdef CONFIG_DEBUG
	unsigned int i;
	unsigned char *p = dev_id;

	dbg_info("SF: Got Manufacturer and Device ID:");
	for (i = 0; i < 5; i++)
		dbg_info(" %x", *p++);
	dbg_info("\n");
#endif

	if (dev_id[0] != MANUFACTURER_ID_ATMEL &&
	    dev_id[0] != MANUFACTURER_ID_WINBOND &&
	    dev_id[0] != MANUFACTURER_ID_MICRON) {
		dbg_info("Not supported spi flash Manufactorer ID: %x\n",
			 dev_id[0]);
		return -1;
	}

	ret = df_desc_init(df_desc, dev_id[0], (dev_id[1] & 0xe0) );
	if (ret)
		return ret;

	return 0;
}

int spi_flash_loadimage(struct image_info *image)
{
	struct dataflash_descriptor	df_descriptor;
	struct dataflash_descriptor	*df_desc = &df_descriptor;
	int ret = 0;

	memset(df_desc, 0, sizeof(*df_desc));

	at91_spi0_hw_init();

	ret = at91_spi_init(AT91C_SPI_PCS_DATAFLASH,
				CONFIG_SYS_SPI_CLOCK, CONFIG_SYS_SPI_MODE);
	if (ret) {
		dbg_info("SF: Fail to initialize spi\n");
		return -1;
	}

	at91_spi_enable();

	ret = dataflash_probe_atmel(df_desc);
	if (ret) {
		dbg_info("SF: Fail to probe atmel spi flash\n");
		ret = -1;
		goto err_exit;
	}

#ifdef CONFIG_DATAFLASH_RECOVERY
	if (!dataflash_recovery(df_desc)) {
		ret = -2;
		goto err_exit;
	}
#endif

#if defined(CONFIG_LOAD_LINUX) || defined(CONFIG_LOAD_ANDROID)
	int length = update_image_length(df_desc,
				image->offset, image->dest, KERNEL_IMAGE);
	if (length == -1)
		return -1;

	image->length = length;
#endif

	dbg_info("SF: Copy %x bytes from %x to %x\n",
			image->length, image->offset, image->dest);

	ret = read_array(df_desc, image->offset, image->length, image->dest);
	if (ret) {
		dbg_info("** SF: Serial flash read error**\n");
		ret = -1;
		goto err_exit;
	}

#ifdef CONFIG_OF_LIBFDT
	length = update_image_length(df_desc,
			image->of_offset, image->of_dest, DT_BLOB);
	if (length == -1)
		return -1;

	image->of_length = length;

	dbg_info("SF: dt blob: Copy %x bytes from %x to %x\n",
		image->of_length, image->of_offset, image->of_dest);

	ret = read_array(df_desc,
		image->of_offset, image->of_length, image->of_dest);
	if (ret) {
		dbg_info("** SF: DT: Serial flash read error**\n");
		ret = -1;
		goto err_exit;
	}
#endif

err_exit:
	at91_spi_disable();
	return ret;
}

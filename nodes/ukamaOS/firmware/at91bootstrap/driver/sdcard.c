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
#include "autoconf.h"
#include "common.h"
#include "hardware.h"
#include "board.h"

#include "string.h"

#include "ff.h"

#include "debug.h"

#define CHUNK_SIZE	0x40000

static int sdcard_loadimage(char *filename, BYTE *dest)
{
	FIL 	file;
	UINT	byte_to_read = CHUNK_SIZE;
	UINT	byte_read;
	FRESULT	fret;
	int	ret;

	fret = f_open(&file, filename, FA_OPEN_EXISTING | FA_READ);
	if (fret != FR_OK) {
		dbg_info("*** FATFS: f_open, filename: [%s]: error\n", filename);
		ret = -1;
		goto open_fail;
	}

	do {
		byte_read = 0;
		fret = f_read(&file, (void *)(dest), byte_to_read, &byte_read);
		dest += byte_to_read;
	} while (byte_read >= byte_to_read);

	if (fret != FR_OK) {
		dbg_info("*** FATFS: f_read: error\n");
		 ret = -1;
		goto read_fail;
	}
	ret = 0;

read_fail:
	fret = f_close(&file);

open_fail:
	return ret;

}

#ifdef CONFIG_OVERRIDE_CMDLINE_FROM_EXT_FILE
static int sdcard_read_cmd(char *cmdline_file, char *cmdline_args)
{
	FIL 	file;
	UINT	byte_to_read = CMDLINE_BUF_LEN;
        UINT    byte_read;
	FRESULT	fret;
        int 	ret;

	fret = f_open(&file, cmdline_file, FA_OPEN_EXISTING | FA_READ);
	if (fret != FR_OK) {
		dbg_info("*** FATFS: f_open, filename: [%s]: error\n", 
                                                                  cmdline_file);
                ret = -1;
		goto open_fail;
	}

	do {
		byte_read = 0;
		fret = f_read(&file, (char *)(cmdline_args), byte_to_read,
                                                                    &byte_read);
	} while (0);

        dbg_info("SD/MMC: kernel arg string: %s\n", cmdline_args);

	if (fret != FR_OK) {
		dbg_info("*** FATFS: cmdline f_read: error\n");
		ret = -1;
		goto read_fail;
	}
        
        ret = 0;

read_fail:
	fret = f_close(&file);

open_fail:
	return ret;

}
#endif

int load_sdcard(struct image_info *image)
{
	FATFS	fs;
	FRESULT	fret;
	int	ret;

#ifdef CONFIG_AT91_MCI
#if defined(CONFIG_AT91_MCI0)
	at91_mci0_hw_init();
#elif defined(CONFIG_AT91_MCI1)
	at91_mci1_hw_init();
#elif defined(CONFIG_AT91_MCI2)
	at91_mci2_hw_init();
#endif
#endif

#ifdef CONFIG_SDHC
	at91_sdhc_hw_init();
#endif

	/* mount fs */
	fret = f_mount(0, &fs);
	if (fret != FR_OK) {
		dbg_info("*** FATFS: f_mount mount error **\n");
		return -1;
	}

	dbg_info("SD/MMC: Image: Read file %s to %x\n",
					image->filename, image->dest);

	ret = sdcard_loadimage(image->filename, image->dest);
	if (ret) {
		(void)f_mount(0, NULL);
		return ret;
	}

#ifdef CONFIG_OF_LIBFDT
	at91_board_set_dtb_name(image->of_filename);

	if (strcmp(CONFIG_OF_OVERRIDE_DTB_NAME, "")) {
                strcpy(image->of_filename, CONFIG_OF_OVERRIDE_DTB_NAME);
        }

	dbg_info("SD/MMC: dt blob: Read file %s to %x\n",
			image->of_filename, image->of_dest);

	ret = sdcard_loadimage(image->of_filename, image->of_dest);
	if (ret) {
		(void)f_mount(0, NULL);
		return ret;
	}

#endif

#ifdef CONFIG_OVERRIDE_CMDLINE_FROM_EXT_FILE
	dbg_info("SD/MMC: kernel arg string: Read file %s\n",
			image->cmdline_file);

	ret = sdcard_read_cmd(image->cmdline_file, image->cmdline_args);
	if (image->cmdline_args) {
		(void)f_mount(0, NULL);
		return ret;
        }

#endif

	/* umount fs */
	fret = f_mount(0, NULL);
	if (fret != FR_OK) {
		dbg_info("*** FATFS: f_mount umount error **\n");
		return -1;
	}

	return 0;
}

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
#ifndef __NAND_H__
#define __NAND_H__

#ifdef CONFIG_AT91SAM9260EK
#define MAX_ECC_BYTES		48  /* maximum bytes of ecc */
#else
#define MAX_ECC_BYTES		512 /* maximum bytes of ecc */
#endif

struct nand_ooblayout {
	unsigned short	badblockpos;
	unsigned short	eccbytes;
	unsigned short	eccpos[MAX_ECC_BYTES];
};

struct nand_chip {
	unsigned short	chip_id; /* Nand Chip ID */
	unsigned short	numblocks;
	unsigned int	blocksize;
	unsigned short	pagesize;
	unsigned short	oobsize;
	unsigned char	buswidth;
	unsigned char	eccbits;
	unsigned int	eccwordsize;
};

struct nand_info {
	unsigned int	blocksize;	/* size of a block */
	unsigned int	pagesize;	/* size of a page */
	unsigned int	oobsize;	/* size of a oob */
	unsigned int	sectorsize;	/* size of a sector */

	unsigned int	numblocks;	/* number of blocks in device */
	unsigned int	pages_device;	/* number of pages in device */
	unsigned int	pages_block;	/* number of pages in block */

	unsigned int	buswidth;	/* data bus width (8/16 bits) */

	void (*command)(unsigned char cmd);
	void (*address)(unsigned char addr);

	struct nand_ooblayout	*ecclayout;
	/* Used by PMECC */
	int			ecc_sector_size;
	int			ecc_err_bits;
};

#define ZONE_DATA			0x01    /* Sector data zone */
#define ZONE_INFO			0x02    /* Sector info zone */

/* Nand flash chip status codes */
#define STATUS_READY			(0x01 << 6)   /* Status code for Ready */
#define STATUS_ERROR			(0x01 << 0)   /* Status code for Error */

/* Nand flash commands */
#define CMD_READ_1			0x00
#define CMD_READ_2			0x30

#define CMD_READID			0x90

#define CMD_WRITE_1			0x80
#define CMD_WRITE_2			0x10

#define CMD_ERASE_1			0x60
#define CMD_ERASE_2			0xD0

#define CMD_STATUS			0x70
#define CMD_RESET			0xFF

/* Nand flash commands (small blocks) */
#define CMD_READ_A0			0x00
#define CMD_READ_A1			0x01
#define CMD_READ_C			0x50

#define CMD_WRITE_A			0x00
#define CMD_WRITE_C			0x50

/* NandFlash ONFI */
#define CMD_READ_ONFI			0xEC

/* Feature Operations */
#define CMD_SET_FEATURE			0xEF
#define CMD_GET_FEATURE			0xEE

#endif /* #ifndef __NAND_H__ */

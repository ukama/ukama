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

#include "spi_flash/sfdp.h"
#include "spi_flash/spi_nor.h"
#include "string.h"

static int spi_flash_read_sfdp(struct spi_flash *flash, size_t from,
			       size_t len, void *buf)
{
	struct spi_flash_command cmd;

	spi_flash_command_init(&cmd, SFLASH_INST_READ_SFDP, 3,
			       SFLASH_TYPE_READ);
	cmd.proto = flash->read_proto;
	cmd.addr = from;
	cmd.num_wait_states = 8;
	cmd.data_len = len;
	cmd.rx_data = buf;
	return spi_flash_exec(flash, &cmd);
}

struct sfdp_parameter_header {
	u8		id_lsb;
	u8		minor;
	u8		major;
	u8		length; /* in double words */
	u8		parameter_table_pointer[3]; /* byte address */
	u8		id_msb;
};

#define SFDP_PARAM_HEADER_ID(p)	((u16)(((p)->id_msb << 8) | (p)->id_lsb))
#define SFDP_PARAM_HEADER_PTP(p)				\
	((u32)(((p)->parameter_table_pointer[2] << 16) |	\
	       ((p)->parameter_table_pointer[1] <<  8) |	\
	       ((p)->parameter_table_pointer[0] <<  0)))


#define SFDP_BFPT_ID		0xff00u	/* Basic Flash Parameter Table */
#define SFDP_4BAIT_ID		0xff84u	/* 4-byte Address Instruction Table */

#define SFDP_SIGNATURE		0x50444653u
#define SFDP_JESD216_MAJOR	1
#define SFDP_JESD216_MINOR	0
#define SFDP_JESD216A_MINOR	5
#define SFDP_JESD216B_MINOR	6

struct sfdp_header {
	u32		signature; /* Ox50444653 <=> "SFDP" */
	u8		minor;
	u8		major;
	u8		nph; /* 0-base number of parameter headers */
	u8		unused;

	/* Basic Flash Parameter Table. */
	struct sfdp_parameter_header	bfpt_header;
};

/* Basic Flash Parameter Table */

/*
 * JESD216B defines a Basic Flash Parameter Table of 16 DWORDs.
 * They are indexed from 1 but C arrays are indexed from 0.
 */
enum sfdp_bfpt_dword {
	BFPT_DWORD1 = 0,
	BFPT_DWORD2,
	BFPT_DWORD3,
	BFPT_DWORD4,
	BFPT_DWORD5,
	BFPT_DWORD6,
	BFPT_DWORD7,
	BFPT_DWORD8,
	BFPT_DWORD9,
	BFPT_DWORD10,
	BFPT_DWORD11,
	BFPT_DWORD12,
	BFPT_DWORD13,
	BFPT_DWORD14,
	BFPT_DWORD15,
	BFPT_DWORD16,

	BFPT_DWORD_MAX
};

/* The first revision of JESB216 defined only 9 DWORDs. */
#define BFPT_DWORD_MAX_JESD216			9

/* 1st DWORD. */
#define BFPT_DWORD1_FAST_READ_1_1_2		(0x1UL << 16)
#define BFPT_DWORD1_ADDRESS_BYTES_MASK		(0x3UL << 17)
#define BFPT_DWORD1_ADDRESS_BYTES_3_ONLY	(0x0UL << 17)
#define BFPT_DWORD1_ADDRESS_BYTES_3_OR_4	(0x1UL << 17)
#define BFPT_DWORD1_ADDRESS_BYTES_4_ONLY	(0x2UL << 17)
#define BFPT_DWORD1_DTR				(0x1UL << 19)
#define BFPT_DWORD1_FAST_READ_1_2_2		(0x1UL << 20)
#define BFPT_DWORD1_FAST_READ_1_4_4		(0x1UL << 21)
#define BFPT_DWORD1_FAST_READ_1_1_4		(0x1UL << 22)

/* 5th DWORD. */
#define BFPT_DWORD5_FAST_READ_2_2_2		(0x1UL << 0)
#define BFPT_DWORD5_FAST_READ_4_4_4		(0x1UL << 4)

/* 11th DWORD. */
#define BFPT_DWORD11_PAGE_SIZE_SHIFT		4
#define BFPT_DWORD11_PAGE_SIZE_MASK		(0xFUL << 4)

/* 15th DWORD. */

/*
 * (from JESD216B)
 * Quad Enable Requirements (QER):
 * - 000b: Device does not have a QE bit. Device detects 1-1-4 and 1-4-4
 *         reads based on instruction. DQ3/HOLD# functions are hold during
 *         instruction pahse.
 * - 001b: QE is bit 1 of status register 2. It is set via Write Status with
 *         two data bytes where bit 1 of the second byte is one.
 *         [...]
 *         Writing only one byte to the status register has the side-effect of
 *         clearing status register 2, including the QE bit. The 100b code is
 *         used if writing one byte to the status register does not modify
 *         status register 2.
 * - 010b: QE is bit 6 of status register 1. It is set via Write Status with
 *         one data byte where bit 6 is one.
 *         [...]
 * - 011b: QE is bit 7 of status register 2. It is set via Write status
 *         register 2 instruction 3Eh with one data byte where bit 7 is one.
 *         [...]
 *         The status register 2 is read using instruction 3Fh.
 * - 100b: QE is bit 1 of status register 2. It is set via Write Status with
 *         two data bytes where bit 1 of the second byte is one.
 *         [...]
 *         In contrast to the 001b code, writing one byte to the status
 *         register does not modify status register 2.
 * - 101b: QE is bit 1 of status register 2. Status register 1 is read using
 *         Read Status instruction 05h. Status register2 is read using
 *         instruction 35h. QE is set via Writ Status instruction 01h with
 *         two data bytes where bit 1 of the second byte is one.
 *         [...]
 */
#define BFPT_DWORD15_QER_MASK			(0x7UL << 20)
#define BFPT_DWORD15_QER_NONE			(0x0UL << 20) /* Micron */
#define BFPT_DWORD15_QER_SR2_BIT1_BUGGY		(0x1UL << 20)
#define BFPT_DWORD15_QER_SR1_BIT6		(0x2UL << 20) /* Macronix */
#define BFPT_DWORD15_QER_SR2_BIT7		(0x3UL << 20)
#define BFPT_DWORD15_QER_SR2_BIT1_NO_RD		(0x4UL << 20)
#define BFPT_DWORD15_QER_SR2_BIT1		(0x5UL << 20) /* Spansion */

/*
 * (from JESD216B)
 * 0-4-4 Mode Entry Method:
 * - xxx1b: Mode Bits[7:0] = A5h Note: QE must be set prior to using the mode
 * - xx1xb: Read the 8-bit volatile configuration register with 85h, set XIP
 *          bit[3] in the data read, and write the modified data using the
 *          instruction 81h, then Mode Bits[7:0] = 01h
 * - x1xxb: Mode Bit[7:0] = AXh
 * - 1xxxb: Reserved
 */
#define BFPT_DWORD15_0_4_4_MASK          (0xFUL << 16)
#define BFPT_DWORD15_0_4_4_A5            (0x1UL << 16)
#define BFPT_DWORD15_0_4_4_MICRON        (0x1UL << 17)
#define BFPT_DWORD15_0_4_4_AX            (0x1UL << 18)

struct sfdp_bfpt {
	u32	dwords[BFPT_DWORD_MAX];
};

/* Fast Read settings. */

static inline void
spi_flash_set_read_settings_from_bfpt(struct spi_flash_read_command *read,
				      u16 half,
				      enum spi_flash_protocol proto)
{
	read->num_mode_cycles = (half >> 5) & 0x07u;
	read->num_wait_states = (half >> 0) & 0x1Fu;
	read->inst = (half >> 8) & 0xFFu;
	read->proto = proto;
}

struct sfdp_bfpt_read {
	/* The Fast Read x-y-z hardware capability in params->hwcaps.mask. */
	u32				hwcaps;

	/*
	 * The <supported_bit> bit in <supported_dword> BFPT DWORD tells us
	 * whether the Fast Read x-y-z command is supported.
	 */
	enum sfdp_bfpt_dword		supported_dword;
	u32				supported_bit;

	/*
	 * The half-word at offset <setting_shift> in <setting_dword> BFPT DWORD
	 * encodes the op code, the number of mode clocks and the number of wait
	 * states to be used by Fast Read x-y-z command.
	 */
	enum sfdp_bfpt_dword		settings_dword;
	int				settings_shift;

	/* The SPI protocol for this Fast Read x-y-z command. */
	enum spi_flash_protocol		proto;
};

static const struct sfdp_bfpt_read sfdp_bfpt_reads[] = {
	/* Fast Read 1-1-2 */
	{
		SFLASH_HWCAPS_READ_1_1_2,
		BFPT_DWORD1, (0x1UL << 16),	/* Supported bit */
		BFPT_DWORD4, 0,			/* Settings */
		SFLASH_PROTO_1_1_2,
	},

	/* Fast Read 1-2-2 */
	{
		SFLASH_HWCAPS_READ_1_2_2,
		BFPT_DWORD1, (0x1UL << 20),	/* Supported bit */
		BFPT_DWORD4, 16,		/* Settings */
		SFLASH_PROTO_1_2_2,
	},

	/* Fast Read 2-2-2 */
	{
		SFLASH_HWCAPS_READ_2_2_2,
		BFPT_DWORD5,  (0x1UL << 0),	/* Supported bit */
		BFPT_DWORD6, 16,		/* Settings */
		SFLASH_PROTO_2_2_2,
	},

	/* Fast Read 1-1-4 */
	{
		SFLASH_HWCAPS_READ_1_1_4,
		BFPT_DWORD1, (0x1UL << 22),	/* Supported bit */
		BFPT_DWORD3, 16,		/* Settings */
		SFLASH_PROTO_1_1_4,
	},

	/* Fast Read 1-4-4 */
	{
		SFLASH_HWCAPS_READ_1_4_4,
		BFPT_DWORD1, (0x1UL << 21),	/* Supported bit */
		BFPT_DWORD3, 0,			/* Settings */
		SFLASH_PROTO_1_4_4,
	},

	/* Fast Read 4-4-4 */
	{
		SFLASH_HWCAPS_READ_4_4_4,
		BFPT_DWORD5, (0x1UL << 4),	/* Supported bit */
		BFPT_DWORD7, 16,		/* Settings */
		SFLASH_PROTO_4_4_4,
	},
};


/* Sector Erase settings. */

static inline void
spi_flash_set_erase_command_from_bfpt(struct spi_flash_erase_command *cmd,
				      u16 half)
{
	u32 size = (half >> 0) & 0xff;
	u8 inst = (half >> 8) & 0xff;

	/* size == 0 means this Erase Type is not supported. */
	if (size)
		size = (0x1UL << size);
	spi_flash_set_erase_command(cmd, size, inst);
}

struct sfdp_bfpt_erase {
	/*
	 * The half-word at offset <shift> in DWORD <dword> encodes the
	 * op code and erase sector size to be used by Sector Erase commands.
	 */
	enum sfdp_bfpt_dword	dword;
	int			shift;
};

static const struct sfdp_bfpt_erase sfdp_bfpt_erases[SFLASH_CMD_ERASE_MAX] = {
	/* Erase Type 1 in DWORD8 bits[15:0] */
	{BFPT_DWORD8, 0},

	/* Erase Type 2 in DWORD8 bits[31:16] */
	{BFPT_DWORD8, 16},

	/* Erase Type 3 in DWORD9 bits[15:0] */
	{BFPT_DWORD9, 0},

	/* Erase Type 4: in DWORD9 bits[31:16] */
	{BFPT_DWORD9, 16},
};

static int spi_flash_parse_bfpt(struct spi_flash *flash,
				const struct sfdp_parameter_header *bfpt_header,
				struct spi_flash_parameters *params)
{
	struct spi_flash_erase_map *map = &flash->erase_map;
	struct sfdp_bfpt bfpt;
	size_t len;
	int i, cmd, err;
	u32 addr, erase_mask;
	u16 half;

	/* JESD216 Basic Flash Parameter Table length is at least 9 DWORDs. */
	if (bfpt_header->length < BFPT_DWORD_MAX_JESD216)
		return -1;

	/* Read the Basic Flash Parameter Table. */
	len = min(sizeof(bfpt), sizeof(u32) * bfpt_header->length);
	addr = SFDP_PARAM_HEADER_PTP(bfpt_header);
	memset(&bfpt, 0, sizeof(bfpt));
	err = spi_flash_read_sfdp(flash,  addr, len, &bfpt);
	if (err)
		return err;

	/* Flash Memory Density (in bits). */
	params->size = bfpt.dwords[BFPT_DWORD2];
	if (params->size & (0x1UL << 31)) {
		params->size &= ~(0x1UL << 31);
		params->size = 1UL << params->size;
	} else {
		params->size++;
	}
	params->size >>= 3; /* Convert to bytes. */

	/* Fast Read settings. */
	for (i = 0; i < ARRAY_SIZE(sfdp_bfpt_reads); i++) {
		const struct sfdp_bfpt_read *rd = &sfdp_bfpt_reads[i];
		struct spi_flash_read_command *read;

		if (!(bfpt.dwords[rd->supported_dword] & rd->supported_bit))
			continue;

		params->hwcaps.mask |= rd->hwcaps;
		cmd = spi_flash_hwcaps2cmd(rd->hwcaps);
		read = &params->reads[cmd];
		half = bfpt.dwords[rd->settings_dword] >> rd->settings_shift;
		spi_flash_set_read_settings_from_bfpt(read, half, rd->proto);
	}

	/* Sector Erase settings. */
	erase_mask = 0;
	for (i = 0; i < SFLASH_CMD_ERASE_MAX; i++) {
		const struct sfdp_bfpt_erase *er = &sfdp_bfpt_erases[i];

		half = bfpt.dwords[er->dword] >> er->shift;
		spi_flash_set_erase_command_from_bfpt(&map->commands[i], half);

		if (map->commands[i].size)
			erase_mask |= (0x1UL << i);
	}
	spi_flash_init_uniform_erase_map(map, erase_mask, params->size);

	/* Stop here if not JESD216 rev A or later. */
	if (bfpt_header->length < BFPT_DWORD_MAX)
		return 0;

	/* Page size: this field specifies 'N' so the page size = 2^N bytes. */
	params->page_size = bfpt.dwords[BFPT_DWORD11];
	params->page_size &= BFPT_DWORD11_PAGE_SIZE_MASK;
	params->page_size >>= BFPT_DWORD11_PAGE_SIZE_SHIFT;
	params->page_size = (0x1UL << params->page_size);

	/* Enable Quad I/O. */
	switch (bfpt.dwords[BFPT_DWORD15] & BFPT_DWORD15_QER_MASK) {
	default:
	case BFPT_DWORD15_QER_NONE:
		params->quad_enable = NULL;
		break;

	case BFPT_DWORD15_QER_SR2_BIT1_BUGGY:
	case BFPT_DWORD15_QER_SR2_BIT1_NO_RD:
		params->quad_enable = spansion_quad_enable;
		break;

	case BFPT_DWORD15_QER_SR1_BIT6:
		params->quad_enable = macronix_quad_enable;
		break;

	case BFPT_DWORD15_QER_SR2_BIT7:
		params->quad_enable = sr2_bit7_quad_enable;
		break;

	case BFPT_DWORD15_QER_SR2_BIT1:
		params->quad_enable = spansion_new_quad_enable;
		break;
	}

	/* 0-4-4 Mode Entry Method. */
	flash->enable_0_4_4 = NULL;
	if (bfpt.dwords[BFPT_DWORD15] & BFPT_DWORD15_0_4_4_MICRON) {
		flash->enable_0_4_4 = micron_enable_0_4_4;
		flash->xip_mode = 0x00u;
	} else if (bfpt.dwords[BFPT_DWORD15] & BFPT_DWORD15_0_4_4_A5) {
		flash->xip_mode = 0xA5u;
	} else if (bfpt.dwords[BFPT_DWORD15] & BFPT_DWORD15_0_4_4_AX) {
		flash->xip_mode = 0xA0u;
	}

	return 0;
}

int spi_flash_parse_sfdp(struct spi_flash *flash,
			 struct spi_flash_parameters *params)
{
	struct sfdp_parameter_header param_header, bfpt_header;
	struct sfdp_header header;
	int i, err;

	/* Get the SFDP header. */
	err = spi_flash_read_sfdp(flash, 0, sizeof(header), &header);
	if (err)
		return err;

	/* Check the SFDP header version. */
	if (header.signature != SFDP_SIGNATURE ||
	    header.major != SFDP_JESD216_MAJOR ||
	    header.minor < SFDP_JESD216_MINOR)
		return -1;

	/*
	 * Verify that the first and only mandatory parameter header is a
	 * Basic Flash Parameter Table header as specified in JESD216.
	 */
	memcpy(&bfpt_header, &header.bfpt_header, sizeof(bfpt_header));
	if (SFDP_PARAM_HEADER_ID(&bfpt_header) != SFDP_BFPT_ID ||
	    bfpt_header.major != SFDP_JESD216_MAJOR)
		return -1;

	/*
	 * Check other parameter headers to get the latest revision of
	 * the basic flash parameter table.
	 */
	for (i = 0; i < header.nph; i++) {
		err = spi_flash_read_sfdp(flash, sizeof(header) +
					  i * sizeof(param_header),
					  sizeof(param_header),
					  &param_header);
		if (err)
			goto exit;

		if (SFDP_PARAM_HEADER_ID(&param_header) == SFDP_BFPT_ID &&
		    param_header.major == SFDP_JESD216_MAJOR &&
		    (param_header.minor > bfpt_header.minor ||
		     (param_header.minor == bfpt_header.minor &&
		      param_header.length > bfpt_header.length)))
			memcpy(&bfpt_header, &param_header, sizeof(bfpt_header));
	}
	err = spi_flash_parse_bfpt(flash, &bfpt_header, params);
	if (err)
		goto exit;

	/* Parse other parameter headers. */
	for (i = 0; i < header.nph; i++) {
		err = spi_flash_read_sfdp(flash,
					  sizeof(header) +
					  i * sizeof(param_header),
					  sizeof(param_header),
					  &param_header);
		if (err)
			goto exit;

		switch (SFDP_PARAM_HEADER_ID(&param_header)) {
		default:
			break;
		}

		if (err)
			goto exit;
	}

exit:
	return err;
}

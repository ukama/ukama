/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2014, Atmel Corporation
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
#include "string.h"
#include "board_hw_info.h"
#include "ds24xx.h"
#include "at24xx.h"
#include "debug.h"

/* Board Type */
#define BOARD_TYPE_CPU		1
#define BOARD_TYPE_EK		2
#define BOARD_TYPE_DM		3

/* One Wire informaion */
#define BOARD_NAME_LEN		12
#define VENDOR_NAME_LEN		10
#define VENDOR_COUNTRY_LEN	2

#define VENDOR_NAME_ALT_LEN	6
#define VENDOR_COUNTRY_ALT_LEN	3

typedef struct hw_info_map {
	unsigned char total_bytes;
	char vendor_name[VENDOR_NAME_LEN];
	char vendor_country[VENDOR_COUNTRY_LEN];
	unsigned char board_name[BOARD_NAME_LEN];
	unsigned char year;
	unsigned char week;
	unsigned char revision_code;
	unsigned char revision_id;
	unsigned char bom_revision;
	unsigned char revision_mapping;
} __attribute__((packed)) hw_info_map_t;

typedef struct hw_info_alt_map {
	unsigned char size_byte;
	unsigned char vendor_name[VENDOR_NAME_ALT_LEN];
	unsigned char vendor_country[VENDOR_COUNTRY_ALT_LEN];
	unsigned char year;
	unsigned char week;
	unsigned char revision_code;
	unsigned char revision_id;
	unsigned char bom_revision;
	unsigned char crc_serial;
	unsigned char ext_config_addr[2];
	unsigned char board_name[BOARD_NAME_LEN];
	unsigned char revision_mapping;
} __attribute__((packed)) hw_info_alt_map_t;

typedef struct board_hw_info {
	char *board_name;
	char *vendor_name;
	unsigned char board_type;
	unsigned char board_id;
	unsigned char revision_code;
	unsigned char revision_id;
	unsigned char bom_revision;
	unsigned char vendor_id;
	unsigned char year;
	unsigned char week;
} board_info_t;

static unsigned int sn;
static unsigned int rev;
static unsigned char buffer[HW_INFO_TOTAL_SIZE];

static struct {
	char *board_name;
	unsigned char board_type;
	unsigned char board_id;
} board_list[] = {
	{"SAM9x5-EK",	BOARD_TYPE_EK,	BOARD_ID_SAM9X5_EK},
	{"SAM9x5-DM",	BOARD_TYPE_DM,	BOARD_ID_SAM9x5_DM},
	{"SAM9G15-CM",	BOARD_TYPE_CPU,	BOARD_ID_SAM9G15_CM},
	{"SAM9G25-CM",	BOARD_TYPE_CPU,	BOARD_ID_SAM9G25_CM},
	{"SAM9G35-CM",	BOARD_TYPE_CPU,	BOARD_ID_SAM9G35_CM},
	{"SAM9X25-CM",	BOARD_TYPE_CPU,	BOARD_ID_SAM9X25_CM},
	{"SAM9X35-CM",	BOARD_TYPE_CPU,	BOARD_ID_SAM9X35_CM},
	{"PDA-DM",	BOARD_TYPE_DM,	BOARD_ID_PDA_DM},
	{"TM4300",	BOARD_TYPE_DM,	BOARD_ID_PDA_DM},
	{"TM4301",	BOARD_TYPE_DM,	BOARD_ID_PDA_DM},
	{"TM7000",	BOARD_TYPE_DM,	BOARD_ID_PDA7_DM},
	{"SAMA5D3x-MB",	BOARD_TYPE_EK,	BOARD_ID_SAMA5D3X_MB},
	{"SAMA5D3x-DM",	BOARD_TYPE_DM,	BOARD_ID_SAMA5D3X_DM},
	{"SAMA5D31-CM",	BOARD_TYPE_CPU,	BOARD_ID_SAMA5D31_CM},
	{"SAMA5D33-CM",	BOARD_TYPE_CPU,	BOARD_ID_SAMA5D33_CM},
	{"SAMA5D34-CM",	BOARD_TYPE_CPU,	BOARD_ID_SAMA5D34_CM},
	{"SAMA5D35-CM",	BOARD_TYPE_CPU,	BOARD_ID_SAMA5D35_CM},
	{"SAMA5D36-CM",	BOARD_TYPE_CPU,	BOARD_ID_SAMA5D36_CM},
	{"SAMA5D36-CP",	BOARD_TYPE_CPU,	BOARD_ID_SAMA5D36_CP},
	{"SAMA5D44-MB",	BOARD_TYPE_EK,	BOARD_ID_SAMA5D4_MB},
	{"SAMA5D4-XULT",BOARD_TYPE_EK,	BOARD_ID_SAMA5D4_MB},
	{"SAMA5D2-XULT",BOARD_TYPE_EK,	BOARD_ID_SAMA5D2_XULT},
	{"SAMA5D2-ICP", BOARD_TYPE_EK,	BOARD_ID_SAMA5D2_ICP},
	{"SAM9X60-EK",	BOARD_TYPE_EK,	BOARD_ID_SAM9X60_EK},
	{0,		0,		0},
};

static struct {
	char *vendor_name;
	char vendor_id;
} vendor_list[] = {
	{"EMBEST",	VENDOR_EMBEST},
	{"FLEX",	VENDOR_FLEX},
	{"RONETIX",	VENDOR_RONETIX},
	{"COGENT",	VENDOR_COGENT},
	{"PDA INC",	VENDOR_PDA},
	{"PDA",		VENDOR_PDA},
	{"ATMEL-RFO",	VENDOR_ATMEL_RFO},
	{"ATMEL RFO",	VENDOR_ATMEL_RFO},
	{"ATMEL-RF0",	VENDOR_ATMEL_RFO},
	{"MCHIP RFO",	VENDOR_MCHIP_RFO},
	{"MCHIP RDC",	VENDOR_MCHIP_RDC},
	{0,		0},
};

static unsigned char normalize_rev_code(const unsigned char c)
{
	if ((c >= 'A') && (c <= 'Z'))
		return c;

	if ((c >= 'a') && (c <= 'z'))
		return c - 0x20;

	if ((c >= '0') && (c <= '9'))
		return c;

	/* by default, return revision 'A' */
	return 'A';
}

static unsigned char normalize_rev_id(const unsigned char c)
{
	if ((c >= '0') && (c <= '9'))
		return c;

	if (c > '9') {
		if (((c >= 'A') && (c <= 'F'))
			|| ((c >= 'a') && (c <= 'f')))
			return '9';
	}

	/* by default, return revision '0' */
	return '0';
}

static unsigned char normalize_rev_id_map_b(const unsigned char c)
{
	return normalize_rev_code(c);
}

static unsigned char normalize_bom_revision(const unsigned char c)
{
	return normalize_rev_id(c);
}

static int parse_board_hw_info(unsigned char *buff,
			       board_info_t *bd_info)
{
	char board_name[BOARD_NAME_LEN + 1];
	char vendor_name[VENDOR_NAME_LEN + 1];
	int i;

	hw_info_map_t *p = (hw_info_map_t *)buff;

	if (p->total_bytes != HW_INFO_TOTAL_SIZE) {
		dbg_loud("HW Info: The total size: %x isn't correct\n",
			 p->total_bytes);
		return -1;
	}

	for (i = 0; i < ARRAY_SIZE(board_list); i++) {
		if (strncmp(board_list[i].board_name,
				(const char *)p->board_name,
				strlen(board_list[i].board_name)) == 0)
			break;
	}

	if (i == ARRAY_SIZE(board_list))
		goto fail_to_search_board;

	bd_info->board_type = board_list[i].board_type;
	bd_info->board_id = board_list[i].board_id;

	bd_info->revision_code = normalize_rev_code(p->revision_code);
	if (p->revision_mapping == 'B') {
		bd_info->revision_id = normalize_rev_id_map_b(p->revision_id);
		bd_info->bom_revision = normalize_bom_revision(p->bom_revision);
	} else {
		bd_info->revision_id = normalize_rev_id(p->revision_id);
	}
	bd_info->board_name = board_list[i].board_name;

	for (i = 0; i < ARRAY_SIZE(vendor_list); i++) {
		if (strncmp(vendor_list[i].vendor_name,
				(const char *)p->vendor_name,
				strlen(vendor_list[i].vendor_name)) == 0)
			break;
	}

	if (i == ARRAY_SIZE(vendor_list))
		goto fail_to_search_vendor;

	bd_info->vendor_id = vendor_list[i].vendor_id;
	bd_info->vendor_name = vendor_list[i].vendor_name;

	return 0;

fail_to_search_board:
	for (i = 0; i < BOARD_NAME_LEN; i++)
		board_name[i] = p->board_name[i];
	board_name[i] = 0;

	dbg_info("Failed to parse the board name: %s\n", board_name);

	return -1;

fail_to_search_vendor:
	for (i = 0; i < VENDOR_NAME_LEN; i++)
		vendor_name[i] = p->vendor_name[i];
	vendor_name[i] = 0;

	dbg_info("Failed to parse the vendor name: %s\n", vendor_name);

	return -1;
}

static int parse_alt_board_hw_info(unsigned char *buff,
				   board_info_t *bd_info)
{
	unsigned int i;
	unsigned char *board, *vendor;
	hw_info_alt_map_t *p = (hw_info_alt_map_t *)buff;

	if (p->size_byte != 0x41) {
		dbg_info("%s: Size of byte is incorrect\n", __func__);
		return -1;
	}

	for (i = 0; i < ARRAY_SIZE(board_list); i++) {
		if (!strncmp((const char *)p->board_name,
			     (const char *)board_list[i].board_name,
			     strlen(board_list[i].board_name)))
			break;
	}

	if (i == ARRAY_SIZE(board_list))
		goto failed;

	bd_info->board_type = board_list[i].board_type;
	bd_info->board_id = board_list[i].board_id;
	bd_info->board_name = board_list[i].board_name;

	for (i = 0; i < ARRAY_SIZE(vendor_list); i++) {
		if (!strncmp((const char *)p->vendor_name,
			     (const char *)vendor_list[i].vendor_name,
			     strlen(vendor_list[i].vendor_name)))
			break;
	}

	if (i == ARRAY_SIZE(vendor_list))
		goto failed;

	bd_info->vendor_id = vendor_list[i].vendor_id;
	bd_info->vendor_name = vendor_list[i].vendor_name;

	bd_info->revision_code = normalize_rev_code(p->revision_code);
	bd_info->revision_id = normalize_rev_id_map_b(p->revision_id);
	bd_info->bom_revision = normalize_bom_revision(p->bom_revision);

	bd_info->year = p->year;
	bd_info->week = p->week;

	return 0;

failed:
	board = p->board_name;
	vendor = p->vendor_name;
	board[BOARD_NAME_LEN] = 0;
	vendor[VENDOR_NAME_ALT_LEN + VENDOR_COUNTRY_ALT_LEN] = 0;

	dbg_info("%s: Failed to parse with board: %s, vendor: %s\n",
		 __func__, board, vendor);

	return -1;
}
static int get_board_hw_info(unsigned char *buff,
			     unsigned char board_sn,
			     board_info_t *bd_info)
{
	int ret;
	unsigned char mapping_revision = buff[30];

	if (mapping_revision == 'C')
		ret = parse_alt_board_hw_info(buff, bd_info);
	else
		ret = parse_board_hw_info(buff, bd_info);

	if (ret)
		return ret;

	dbg_info("  #%d", board_sn);
	if ((mapping_revision == 'B') ||
	    (mapping_revision == 'C')) {
		dbg_info("  %s [%c%c%c]      %s\n",
		bd_info->board_name,
		bd_info->revision_code,
		bd_info->revision_id,
		bd_info->bom_revision,
		bd_info->vendor_name);
	} else {
		dbg_info("  %s [%c%c]      %s\n",
		bd_info->board_name,
		bd_info->revision_code,
		bd_info->revision_id,
		bd_info->vendor_name);
	}

	return 0;
}

/*******************************************************************************
 * SN layout
 *
 *   31  30         25         20         15         10         5          0
 *   -----------------------------------------------------------------------
 *   |   | Vendor ID| Board ID | Vendor ID| Board ID | Vendor ID| Board ID |
 *   -----------------------------------------------------------------------
 *       |         EK          |         DM          |         CPU         |
 *
 * Rev layout
 *
 *   31              24     21     18     15         10         5          0
 *   -----------------------------------------------------------------------
 *   |               |  EK  |  DM  |  CPU |    EK    |    DM    |   CPU    |
 *   -----------------------------------------------------------------------
 *                   |     Revision id    |        Revision Code           |
 *
 *******************************************************************************
 */
/* SN */
#define SN_MASK			0x1F
#define CM_SN_OFFSET		0
#define DM_SN_OFFSET		10
#define EK_SN_OFFSET		20

/* Vendor */
#define VENDOR_MASK		0x1F
#define CM_VENDOR_OFFSET	5
#define DM_VENDOR_OFFSET	15
#define EK_VENDOR_OFFSET	25

/* Revision Code */
#define REV_MASK		0x1F
#define CM_REV_OFFSET		0
#define DM_REV_OFFSET		5
#define EK_REV_OFFSET		10

/* Revision ID */
#define	REV_ID_MASK		0x07
#define CM_REV_ID_OFFSET	15
#define DM_REV_ID_OFFSET	18
#define EK_REV_ID_OFFSET	21

static int construct_sn_rev(board_info_t *bd_info,
			    unsigned int *psn,
			    unsigned int *prev)
{
	int ret = 0;

	switch (bd_info->board_type) {
	case BOARD_TYPE_CPU:
		*psn |= (bd_info->board_id & SN_MASK);
		*psn |= ((bd_info->vendor_id & VENDOR_MASK)
							<< CM_VENDOR_OFFSET);
		*prev |= ((bd_info->revision_code - 'A') & REV_MASK);
		*prev |= (((bd_info->revision_id - '0') & REV_ID_MASK)
							<< CM_REV_ID_OFFSET);
		break;

	case BOARD_TYPE_DM:
		*psn |= ((bd_info->board_id & SN_MASK) << DM_SN_OFFSET);
		*psn |= ((bd_info->vendor_id & VENDOR_MASK)
							<< DM_VENDOR_OFFSET);
		*prev |= (((bd_info->revision_code - 'A') & REV_MASK)
							<< DM_REV_OFFSET);
		*prev |= (((bd_info->revision_id - '0') & REV_ID_MASK)
							<< DM_REV_ID_OFFSET);
		break;

	case BOARD_TYPE_EK:
		*psn |= ((bd_info->board_id & SN_MASK) << EK_SN_OFFSET);
		*psn |= ((bd_info->vendor_id & VENDOR_MASK)
							<< EK_VENDOR_OFFSET);
		*prev |= (((bd_info->revision_code - 'A') & REV_MASK)
							<< EK_REV_OFFSET);
		*prev |= (((bd_info->revision_id - '0') & REV_ID_MASK)
							<< EK_REV_ID_OFFSET);
		break;

	default:
		dbg_info("WARNING: Unknown board type\n");
		ret = -1;
		break;
	}

	return ret;
}

static unsigned int set_default_sn(void)
{
	unsigned int board_id_cm = 0;
	unsigned int board_id_dm = 0;
	unsigned int board_id_ek = 0;
	unsigned int vendor_cm = 0;
	unsigned int vendor_dm = 0;
	unsigned int vendor_ek = 0;

#if defined(CONFIG_AT91SAM9X5EK)
	/* at91sam9x5ek
	 * CPU Module: SAM9X25-CM, EMBEST
	 * Display Module: SAM9x5-DM, FLEX
	 * EK Module: SAM9x5-EK, FLEX
	 */
	board_id_cm = BOARD_ID_SAM9X25_CM;
	board_id_dm = BOARD_ID_SAM9x5_DM;
	board_id_ek = BOARD_ID_SAM9X5_EK;
	vendor_cm = VENDOR_EMBEST;
	vendor_dm = VENDOR_FLEX;
	vendor_ek = VENDOR_FLEX;

#elif defined(CONFIG_SAMA5D3XEK) || defined(CONFIG_SAMA5D3X_CMP)

	/* sama5d3xek
	 * CPU Module: SAMA5D31-CM, EMBEST
	 * Display Module: SAMA5D3x-DM, FLEX
	 * EK Module: SAMA5D3x-MB, FLEX
	 */
	board_id_cm = BOARD_ID_SAMA5D31_CM;
	board_id_dm = BOARD_ID_SAMA5D3X_DM;
	board_id_ek = BOARD_ID_SAMA5D3X_MB;
	vendor_cm = VENDOR_EMBEST;
	vendor_dm = VENDOR_FLEX;
	vendor_ek = VENDOR_FLEX;
#elif defined(CONFIG_SAMA5D4EK) || defined(CONFIG_SAMA5D4_XPLAINED)
	/*
	 * SAMA5D4-EK
	 * Display Module: SAMA5D3x-DM, FLEX
	 * EK Module: SAMA5D4x-MB, FLEX
	 */
	board_id_dm = BOARD_ID_SAMA5D3X_DM;
	board_id_ek = BOARD_ID_SAMA5D4_MB;
	vendor_cm = VENDOR_EMBEST;
	vendor_dm = VENDOR_FLEX;
#elif defined(CONFIG_SAMA5D2_XPLAINED)
	board_id_ek = BOARD_ID_SAMA5D2_XULT;
	vendor_ek = VENDOR_ATMEL_RFO;
#elif defined(CONFIG_SAMA5D2_ICP)
	board_id_ek = BOARD_ID_SAMA5D2_ICP;
	vendor_ek = VENDOR_MCHIP_RFO;
#elif defined(CONFIG_SAM9X60EK) || defined(CONFIG_SAM9X60_DDR2_SIP_EB) || defined(CONFIG_SAM9X60_SDR_SIP_EB)
	/* sam9x60ek ; sam9x60 ddr2 sip eb and sam9x60 sdr sip eb
	 */
	board_id_ek = BOARD_ID_SAM9X60_EK;
	vendor_ek = VENDOR_MCHIP_RDC;
#else
#error "OneWire: No defined board"
#endif

	return (board_id_cm & SN_MASK)
		| ((vendor_cm & VENDOR_MASK) << CM_VENDOR_OFFSET)
		| ((board_id_dm & SN_MASK) << DM_SN_OFFSET)
		| ((vendor_dm & VENDOR_MASK) << DM_VENDOR_OFFSET)
		| ((board_id_ek & SN_MASK) << EK_SN_OFFSET)
		| ((vendor_ek & VENDOR_MASK) << EK_VENDOR_OFFSET);
}

static unsigned int set_default_rev(void)
{
	unsigned int rev_cm;
	unsigned int rev_dm;
	unsigned int rev_ek;
	unsigned int rev_id_cm;
	unsigned int rev_id_dm;
	unsigned int rev_id_ek;

#if defined(CONFIG_AT91SAM9X5EK)
	/* at91sam9x5ek
	 * CPU Module: 'B', '1'
	 * Display Module: 'B', '0'
	 * EK Module: 'B','0'
	 */
	rev_cm = 'B';
	rev_dm = 'B';
	rev_ek = 'B';
	rev_id_cm = '1';
	rev_id_dm = '0';
	rev_id_ek = '0';

#elif defined(CONFIG_SAMA5D3XEK) || defined(CONFIG_SAMA5D3X_CMP)

	/* sama5d3xek
	 * CPU Module: 'D', '4'
	 * Display Module: 'B', '2'
	 * EK Module: 'C','3'
	 */
	rev_cm = 'D';
	rev_dm = 'B';
	rev_ek = 'C';
	rev_id_cm = '4';
	rev_id_dm = '2';
	rev_id_ek = '3';
#elif defined(CONFIG_SAMA5D4EK) || defined(CONFIG_SAMA5D4_XPLAINED)
	/*
	 * SAMA5D4-EK
	 * Display Module: 'B', '2'
	 * EK Module: 'B','3'
	 */
	rev_cm = 'A';
	rev_dm = 'B';
	rev_ek = 'C';
	rev_id_cm = '0';
	rev_id_dm = '2';
	rev_id_ek = '3';

#elif defined(CONFIG_SAMA5D2_XPLAINED)
	rev_cm = 'A';
	rev_dm = 'A';
	rev_ek = 'A';
	rev_id_cm = '1';
	rev_id_dm = '1';
	rev_id_ek = '1';

#elif defined(CONFIG_SAMA5D2_ICP)
	rev_cm = 'A';
	rev_dm = 'A';
	rev_ek = 'A';
	rev_id_cm = '1';
	rev_id_dm = '1';
	rev_id_ek = '1';
#elif defined(CONFIG_SAM9X60EK)
	rev_cm = 'A';
	rev_dm = 'A';
	rev_ek = 'A';
	rev_id_cm = '0';
	rev_id_dm = '0';
	rev_id_ek = '0';
#elif defined(CONFIG_SAM9X60_DDR2_SIP_EB)
	rev_cm = 'A';
	rev_dm = 'A';
	rev_ek = 'A';
	rev_id_cm = '0';
	rev_id_dm = '0';
	rev_id_ek = '0';
#elif defined(CONFIG_SAM9X60_SDR_SIP_EB)
	rev_cm = 'A';
	rev_dm = 'A';
	rev_ek = 'A';
	rev_id_cm = '0';
	rev_id_dm = '0';
	rev_id_ek = '0';
#else
#error "OneWire: No defined board"
#endif

	return ((rev_cm - 'A') & REV_MASK)
		| (((rev_dm - 'A') & REV_MASK) << DM_REV_OFFSET)
		| (((rev_ek - 'A') & REV_MASK) << EK_REV_OFFSET)
		| (((rev_id_cm - '0') & REV_ID_MASK) << CM_REV_ID_OFFSET)
		| (((rev_id_dm - '0') & REV_ID_MASK) << DM_REV_ID_OFFSET)
		| (((rev_id_ek - '0') & REV_ID_MASK) << EK_REV_ID_OFFSET);
}

unsigned int get_sys_sn(void)
{
	return sn;
}

unsigned int get_sys_rev(void)
{
	return rev;
}

unsigned int get_cm_sn(void)
{
	return (sn  >> CM_SN_OFFSET) & SN_MASK;
}

unsigned int get_cm_vendor(void)
{
	return (sn >> CM_VENDOR_OFFSET) & VENDOR_MASK;
}

char get_cm_rev(void)
{
	return 'A' + ((rev >> CM_REV_OFFSET) & REV_MASK);
}

unsigned int get_dm_sn(void)
{
	return (sn >> DM_SN_OFFSET) & SN_MASK;
}

char get_ek_rev(void)
{
	return 'A' + ((rev >> EK_REV_OFFSET) & REV_MASK);
}

unsigned int get_ek_sn(void)
{
	return (sn  >> EK_SN_OFFSET) & SN_MASK;
}

#if defined(CONFIG_LOAD_ONE_WIRE)
static int load_1wire_info(unsigned char *buff, unsigned int size,
			   unsigned int *psn, unsigned int *prev)
{
	board_info_t board_info;
	board_info_t *bd_info = &board_info;
	unsigned int count;
	unsigned int parsing = 0;
	int i;

	memset(bd_info, 0, sizeof(*bd_info));

	dbg_info("1-Wire: Loading 1-Wire information ...\n");

	count = enumerate_all_rom();
	if (!count) {
		dbg_info("WARNING: 1-Wire: No 1-Wire chip found\n");
		return -1;
	}

	dbg_info("1-Wire: BoardName | [Revid] | VendorName\n");

	for (i = 0; i < count; i++) {
		if (ds24xx_read_memory(i, 0, 0, size, buff) < 0) {
			dbg_info("WARNING: 1-Wire: Failed to read from 1-Wire chip!\n");
			return -1;
		}

		dbg_loud("board: #%d: ", i);
#if (BOOTSTRAP_DEBUG_LEVEL >= DEBUG_LOUD)
		dbg_hexdump(buff, size, DUMP_WIDTH_BIT_8);
#endif

		if (get_board_hw_info(buff, i, bd_info))
			continue;

		if (construct_sn_rev(bd_info, psn, prev))
			continue;

		parsing++;
	}

	if (!parsing)
		return -1;

	return 0;
}
#endif /* #if defined(CONFIG_LOAD_ONE_WIRE) */

#if defined(CONFIG_LOAD_EEPROM)
static int load_eeprom_info(unsigned char *buff, unsigned int size,
			    unsigned int boardsn, unsigned int *psn,
			    unsigned int *prev)
{
	board_info_t board_info;
	board_info_t *bd_info = &board_info;

	memset(bd_info, 0, sizeof(*bd_info));

	dbg_info("EEPROM: Loading AT24xx information ...\n");

	dbg_info("EEPROM: BoardName | [Revid] | VendorName\n");

	if (load_ek_at24xx(buff, size))
		return -1;

	if (get_board_hw_info(buff, boardsn, bd_info))
		return -1;

	if (construct_sn_rev(bd_info, psn, prev))
		return -1;

	dbg_info("EEPROM: BoardDate | Year | Week\n");
	dbg_info("EEPROM:             %u    %u\n",
		 2000 + bd_info->year, bd_info->week);

	return 0;
}
#endif /* #if defined(CONFIG_LOAD_EEPROM) */

void load_board_hw_info(void)
{
	unsigned int size = HW_INFO_TOTAL_SIZE;
	int ret;

#if defined(CONFIG_LOAD_ONE_WIRE)
	ret = load_1wire_info(buffer, size, &sn, &rev);
#endif
#if defined(CONFIG_LOAD_EEPROM)
	ret = load_eeprom_info(buffer, size, 0, &sn, &rev);
#endif
	if (ret) {
#if defined(CONFIG_LOAD_ONE_WIRE)
		dbg_info("\n1-Wire: ");
#endif
#if defined(CONFIG_LOAD_EEPROM)
		dbg_info("\nEEPROM: ");
#endif
		dbg_info("Using default information\n");

		sn = set_default_sn();
		rev = set_default_rev();
	}

#ifdef AT91C_BASE_GPBR
	writel(sn, AT91C_BASE_GPBR + 4 * 2);
	writel(rev, AT91C_BASE_GPBR + 4 * 3);
#endif

#if defined(CONFIG_LOAD_ONE_WIRE)
	dbg_info("\n1-Wire: ");
#endif
#if defined(CONFIG_LOAD_EEPROM)
	dbg_info("\nEEPROM: ");
#endif
	dbg_info("Board sn: %x revision: %x\n\n", sn, rev);

	return;
}

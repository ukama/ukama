/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2011, Atmel Corporation
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
#include "board.h"
#include "gpio.h"
#include "arch/at91_pio.h"
#include "debug.h"
#include "timer.h"

/* Commands */
#define ROM_COMMAND_READ		0x33
#define ROM_COMMAND_MATCH		0x55
#define ROM_COMMAND_SEARCH		0xF0
#define ROM_COMMAND_SKIP		0xCC
#define ROM_COMMAND_RESUME		0xA5
#define ROM_COMMAND_OSKIP		0x3C
#define ROM_COMMAND_OMATCH		0x69

#define MEMORY_COMMAND_WSCRATCHPAD	0x0F
#define MEMORY_COMMAND_RSCRATCHPAD	0xAA
#define MEMORY_COMMAND_CSCRATCHPAD	0x55
#define MEMORY_COMMAND_READMEMORY	0xF0

/* Timing */
#define tSLOT				65
#define tRSTL				480
#define tPDH				60
#define tPDL				240
#define tWOL				60
#define tW1L				5
#define tRL				5
#define tMSR				10
#define tPROG				13000

#define MAX_RETRY			10
#define MAX_BUF_LEN			256

#define MAX_ITEMS			5
#define CHIP_ADDR_LEN			8

/* Family Code */
#define FAMILY_CODE_DS2431		0x2D
#define FAMILY_CODE_DS2433		0x23
#define FAMILY_CODE_DS28EC		0x43

static unsigned char device_id_array[MAX_ITEMS][CHIP_ADDR_LEN];
static unsigned char LastDiscrepancy;
static unsigned char LastFamilyDiscrepancy;
static unsigned char LastDeviceFlag;

static unsigned char buf[MAX_BUF_LEN];
static unsigned char cmp[MAX_BUF_LEN];

static inline void set_wire_low()
{
	pio_set_gpio_output(CONFIG_SYS_ONE_WIRE_PIN, 0);
}

static inline void set_wire_input()
{
	pio_set_gpio_input(CONFIG_SYS_ONE_WIRE_PIN, 0);
}

static inline int read_wire_bit()
{
	return pio_get_value(CONFIG_SYS_ONE_WIRE_PIN);
}

static unsigned char dscrc_table[] = {
	  0, 94,188,226, 97, 63,221,131,194,156,126, 32,163,253, 31, 65,
	157,195, 33,127,252,162, 64, 30, 95,  1,227,189, 62, 96,130,220,
	 35,125,159,193, 66, 28,254,160,225,191, 93,  3,128,222, 60, 98,
	190,224,  2, 92,223,129, 99, 61,124, 34,192,158, 29, 67,161,255,
	 70, 24,250,164, 39,121,155,197,132,218, 56,102,229,187, 89,  7,
	219,133,103, 57,186,228,  6, 88, 25, 71,165,251,120, 38,196,154,
	101, 59,217,135,  4, 90,184,230,167,249, 27, 69,198,152,122, 36,
	248,166, 68, 26,153,199, 37,123, 58,100,134,216, 91,  5,231,185,
	140,210, 48,110,237,179, 81, 15, 78, 16,242,172, 47,113,147,205,
	 17, 79,173,243,112, 46,204,146,211,141,111, 49,178,236, 14, 80,
	175,241, 19, 77,206,144,114, 44,109, 51,209,143, 12, 82,176,238,
	 50,108,142,208, 83, 13,239,177,240,174, 76, 18,145,207, 45,115,
	202,148,118, 40,171,245, 23, 73,  8, 86,180,234,105, 55,213,139,
	 87,  9,235,181, 54,104,138,212,149,203, 41,119,244,170, 72, 22,
	233,183, 85, 11,136,214, 52,106, 43,117,151,201, 74, 20,246,168,
	116, 42,200,150, 21, 75,169,247,182,232, 10, 84,215,137,107, 53
};

static unsigned char docrc8(unsigned char crc8, unsigned char value)
{
	return dscrc_table[crc8 ^ value];
}

static int ds24xx_reset(void)
{
	int i;

	set_wire_low();
	udelay(tRSTL);

	set_wire_input();
	udelay(tPDH);

	i = read_wire_bit();
	udelay(tPDL);

	return i ^ 1;
}

static void ds24xx_write_bit(int bit)
{
	if (bit == 1) {
		set_wire_low();
		udelay(tW1L);
		set_wire_input();
		udelay(tSLOT-tW1L);
	} else {
		set_wire_low();
		udelay(tWOL);
		set_wire_input();
		udelay(tSLOT-tWOL);
	}
}

static int ds24xx_read_bit()
{
	int status;

	set_wire_low();
	udelay(tRL);

	set_wire_input();
	udelay(tMSR / 2);

	status = read_wire_bit();
	udelay(tSLOT-tRL-tMSR);

	return status;
}

static void ds24xx_write_byte(unsigned char data)
{
	int i;

	for (i = 0; i < 8; i++) {
		ds24xx_write_bit(data & 1);
		data >>= 1;
	}
}

static unsigned char ds24xx_read_byte()
{
	int i;
	unsigned char result = 0;

	for (i = 0; i < 8; i++) {
		result >>= 1;
		if (ds24xx_read_bit())
			result |= 0x80;
	}

	return result;
}

/* MAXIM App. Note #187 */
static int ds24xx_search_rom()
{
	int id_bit_number;
	int last_zero, rom_byte_number, search_result;
	int id_bit, cmp_id_bit;
	unsigned char rom_byte_mask, search_direction;
	unsigned char crc8 = 0;

	/* initialize for search */
	id_bit_number = 1;
	last_zero = 0;
	rom_byte_number = 0;
	rom_byte_mask = 1;
	search_result = 0;

	/* if the last call was not the last one */
	if (!LastDeviceFlag) {
		/* 1-Wire reset */
		if (!ds24xx_reset()) {
			 /* reset the search*/
			 dbg_info("1-Wire: reset fail\n");
			 LastDiscrepancy = 0;
			 LastDeviceFlag = 0;
			 LastFamilyDiscrepancy = 0;
			 return 0;
		}

		/* issue the search command */
		ds24xx_write_byte(ROM_COMMAND_SEARCH);

		/* loop to do the search */
		do
		{
			/* read a bit and its complement */
			id_bit = ds24xx_read_bit();
			cmp_id_bit = ds24xx_read_bit();

			/* check for no devices on 1-Wire */
			if ((id_bit == 1) && (cmp_id_bit == 1))
				break;
			else {
				/* all devices coupled have 0 or 1 */
				if (id_bit != cmp_id_bit)
					search_direction = id_bit;
				else {
					/*
					 * if this discrepancy if before
					 * the Last Discrepancy on a previous
					 * next then pick the same as last time
					 */
					if (id_bit_number < LastDiscrepancy)
						search_direction =
							((buf[rom_byte_number]
							& rom_byte_mask) > 0);
					else
						/*
						 * if equal to last pick 1,
						 * if not then pick 0
						 */
						search_direction
						= (id_bit_number
							== LastDiscrepancy);

					/*
					 * if 0 was picked then record its
					 * position in LastZero
					 */
					if (search_direction == 0) {
						last_zero = id_bit_number;
						/*
						 * check for Last discrepancy
						 * in family
						 */
						if (last_zero < 9)
							LastFamilyDiscrepancy
								= last_zero;
					}
				}

				/*
				 * set or clear the bit in the ROM byte
				 * rom_byte_number with mask rom_byte_mask
				 */
				if (search_direction == 1)
					buf[rom_byte_number] |= rom_byte_mask;
				else
					buf[rom_byte_number] &= ~rom_byte_mask;

				/* serial number search direction write bit */
				if (search_direction == 1)
					ds24xx_write_bit(1);
				else
					ds24xx_write_bit(0);

				/*
				 * increment the byte counter id_bit_number
				 * and shift the mask rom_byte_mask
				 */
				id_bit_number++;
				rom_byte_mask <<= 1;

				/*
				 * if the mask is 0 then go to new SerialNum
				 * byte rom_byte_number and reset mask
				 */
				if (rom_byte_mask == 0) {
					crc8 = docrc8(crc8,
						buf[rom_byte_number]);
					rom_byte_number++;
					rom_byte_mask = 1;
				}
			}
		} while(rom_byte_number < 8);  /* loop until through all ROM bytes 0-7 */

		/* if the search was successful then */
		if (!((id_bit_number < 65) || (crc8 != 0))) {
			/* search successful so set LastDiscrepancy,LastDeviceFlag,search_result */
			LastDiscrepancy = last_zero;

			/* check for last device */
			if (LastDiscrepancy == 0)
				LastDeviceFlag = 1;

			search_result = 1;
		}
	}

	/* if no device found then reset counters so next 'search' will be like a first */
	if (!search_result || !buf[0]) {
		LastDiscrepancy = 0;
		LastDeviceFlag = 0;
		LastFamilyDiscrepancy = 0;
		search_result = 0;
	}

	return search_result;
}

static int ds24xx_find_first()
{
	/* reset the search state */
	LastDiscrepancy = 0;
	LastDeviceFlag = 0;
	LastFamilyDiscrepancy = 0;

	return ds24xx_search_rom();
}

static int ds24xx_find_next()
{
	/* leave the search state alone */
	return ds24xx_search_rom();
}

unsigned int enumerate_all_rom(void)
{
	int i;
	int result;
	unsigned int cnt = 0;

	dbg_info("1-Wire: ROM Searching ... ");

	result = ds24xx_find_first();
	while (result) {

		/* save device info */
		for (i = 7; i >= 0; i--)
			device_id_array[cnt][i] = buf[i];
		cnt++;

		result = ds24xx_find_next();
	}

	dbg_info("Done, %d 1-Wire chips found\n\n", cnt);

	return cnt;
}

int ds24xx_read_memory(int chip_index, unsigned char addrh,
			unsigned char addrl, int len, unsigned char *p)
{
	int i, round, retries = 0;
	unsigned char *pbuf[2] = {p, cmp};

	switch(device_id_array[chip_index][0]){
	case FAMILY_CODE_DS2431:
	case FAMILY_CODE_DS2433:
	case FAMILY_CODE_DS28EC:
		break;
	default:
		dbg_info("1-Wire: family %x is not supported\n",
					device_id_array[chip_index][0]);
		return -1;
	}

retry:
	for (round = 0; round < 2; round++) {
		if (!ds24xx_reset())
			dbg_info("1-Wire: reset failed\n");

		ds24xx_write_byte(ROM_COMMAND_MATCH);
		for(i = 0; i < 8; i++)
			ds24xx_write_byte(device_id_array[chip_index][i]);

		ds24xx_write_byte(MEMORY_COMMAND_READMEMORY);
		ds24xx_write_byte(addrl);
		ds24xx_write_byte(addrh);

		for(i = 0; i < len; i++)
			pbuf[round][i] = ds24xx_read_byte();
	}


	/* Compare the buffer, if all the same, return 0 */
	for (i = 0; i < len; i++) {
		if (p[i] != cmp[i])
			break;
	}

	if (i == len)
		return 0;

	if (retries++ < MAX_RETRY)
		goto retry;

	return -1;
}

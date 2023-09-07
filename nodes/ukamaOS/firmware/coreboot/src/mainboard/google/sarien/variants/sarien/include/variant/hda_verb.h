/*
 * This file is part of the coreboot project.
 *
 * Copyright 2018 Google LLC
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#ifndef MAINBOARD_HDA_VERB_H
#define MAINBOARD_HDA_VERB_H

#include <device/azalia_device.h>

const u32 cim_verb_data[] = {
	/* coreboot specific header */
	0x10ec0236,	// Codec Vendor / Device ID: Realtek ALC3204
	0xffffffff,	// Subsystem ID
	0x0000001e,	// Number of jacks (NID entries)

	/* Rest Codec First */
	AZALIA_RESET(0x1),
	/* NID 0x01, HDA Codec Subsystem ID Verb Table */
	AZALIA_SUBVENDOR(0x0, 0x102808b8),

	/* Pin Widget Verb Table */
	AZALIA_PIN_CFG(0x0, 0x12, 0x90a60140),
	AZALIA_PIN_CFG(0x0, 0x13, 0x40000000),
	AZALIA_PIN_CFG(0x0, 0x14, 0x90170110),
	AZALIA_PIN_CFG(0x0, 0x18, 0x411111f0),
	AZALIA_PIN_CFG(0x0, 0x19, 0x04a11030),
	AZALIA_PIN_CFG(0x0, 0x1a, 0x411111f0),
	AZALIA_PIN_CFG(0x0, 0x1b, 0x411111f0),
	AZALIA_PIN_CFG(0x0, 0x1d, 0x40700001),
	AZALIA_PIN_CFG(0x0, 0x1e, 0x421212f2),
	AZALIA_PIN_CFG(0x0, 0x21, 0x04211020),

	/* ALC3204 default-1 */
	0x02050040,
	0x02049800,
	0x02050034,
	0x0204023C,
	/* ALC3204 default-2 */
	0x0205003C,
	0x02040354,
	0x0205003C,
	0x02040314,
	/* ALC3204 Speaker output power - 4 ohm 2W (+12dB gain)
	 * + Combo Jack TRS setting */
	0x02050038,
	0x02043901,
	0x02050045,
	0x02045089,
	/* H/W AGC setting-1 */
	0x02050016,
	0x02040C50,
	0x02050012,
	0x0204EBC2,
	/* H/W AGC setting-2 */
	0x02050013,
	0x0204401D,
	0x02050016,
	0x02044E50,
	/* Zero data  + EAPD to verb-control */
	0x02050037,
	0x0204FE15,
	0x02050010,
	0x02040020,
	/* Zero data */
	0x02050030,
	0x02048000,
	0x02050030,
	0x02048000,
	/* ALC3204 default-3 */
	0x05750003,
	0x05740DA3,
	0x02050046,
	0x02040004,
	/* ALC3204 default-4 */
	0x0205001B,
	0x02040A4B,
	0x02050008,
	0x02046A6C,
	/* JD1 */
	0x02050009,
	0x0204E003,
	0x0205000A,
	0x02047770,
	/* Microphone + Array MIC  security Disable  +ADC clock Enable */
	0x0205000D,
	0x0204A020,
	0x02050005,
	0x02040700,
	/* Speaker Enable */
	0x0205000C,
	0x020401EF,
	0x0205000C,
	0x020401EF,
	/* EQ Bypass + EQ HPF cutoff 250Hz */
	0x05350000,
	0x0534201A,
	0x0535001d,
	0x05340800,
	/* EQ-2 */
	0x0535001e,
	0x05340800,
	0x05350003,
	0x05341EF8,
	/* EQ-3 */
	0x05350004,
	0x05340000,
	0x05450000,
	0x05442000,
	/* EQ-4 */
	0x0545001d,
	0x05440800,
	0x0545001e,
	0x05440800,
	/* EQ-5 */
	0x05450003,
	0x05441EF8,
	0x05450004,
	0x05440000,
	/* EQ Update */
	0x05350000,
	0x0534E01A,
	0x05350000,
	0x0534E01A,
};

const u32 pc_beep_verbs[] = {
/* PCBeep pass through to NID14 for ePSA test-1 */
	0x02050036,
	0x02047717,
	0x02050036,
	0x02047717,
/* PCBeep pass through to NID14 for ePSA test-2 */
	0x01470740,
	0x0143B000,
	0x01470C02,
	0x01470C02,
};

AZALIA_ARRAY_SIZES;

#endif

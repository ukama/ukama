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
#ifndef __BOARD_HW_INFO_H__
#define __BOARD_HW_INFO_H__

/* Board id definition */
#define BOARD_ID_SAM9X5_EK	0
#define BOARD_ID_SAM9x5_DM	1
#define	BOARD_ID_SAM9G15_CM	2
#define	BOARD_ID_SAM9G25_CM	3
#define	BOARD_ID_SAM9G35_CM	4
#define	BOARD_ID_SAM9X25_CM	5
#define	BOARD_ID_SAM9X35_CM	6
#define	BOARD_ID_PDA_DM		7
#define	BOARD_ID_SAMA5D3X_MB	8
#define	BOARD_ID_SAMA5D3X_DM	9
#define	BOARD_ID_SAMA5D31_CM	10
#define	BOARD_ID_SAMA5D33_CM	11
#define	BOARD_ID_SAMA5D34_CM	12
#define	BOARD_ID_SAMA5D35_CM	13
#define	BOARD_ID_SAMA5D36_CM	14
#define	BOARD_ID_PDA7_DM	15
#define	BOARD_ID_SAMA5D4_MB	16
#define	BOARD_ID_SAMA5D4_XULT	17
#define	BOARD_ID_SAMA5D36_CP	18
#define	BOARD_ID_SAMA5D2_XULT	19
#define	BOARD_ID_SAMA5D2_ICP	20
#define	BOARD_ID_SAM9X60_EK	21

/* Vendor id definition */
#define VENDOR_EMBEST		1
#define VENDOR_FLEX		2
#define VENDOR_RONETIX		3
#define VENDOR_COGENT		4
#define VENDOR_PDA		5
#define VENDOR_ATMEL_RFO	6
#define VENDOR_MCHIP_RFO	7
#define VENDOR_MCHIP_RDC	8

#define HW_INFO_TOTAL_SIZE	0x20

extern unsigned int get_cm_sn(void);
extern char get_cm_rev(void);
extern unsigned int get_cm_vendor(void);
extern unsigned int get_dm_sn(void);
extern unsigned int get_sys_sn(void);
extern unsigned int get_sys_rev(void);
extern char get_ek_rev(void);
extern unsigned int get_ek_sn(void);

extern void load_board_hw_info(void);

#endif	/* #ifndef __BOARD_HW_INFO_H__ */

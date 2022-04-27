/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2013, Atmel Corporation
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
#ifndef __AT91_SFR_H__
#define __AT91_SFR_H__

/*
 * Register Offset
 */
#define SFR_DDRCFG	0x04	/* DDR Configuration Register */
#define SFR_CCFG_EBICSA	0x04	/* EBI Chip Select Register */
#define SFR_OHCXIICR	0x10	/* OHCI Interrupt Configuration Register */
#define SFR_OHCIISR	0x14	/* OHCI Interrupt Status Register */
/* Reserved */
#define SFR_SECURE	0x28	/* Security Configuration Register */
#define SFR_UTMICKTRIM	0x30	/* UTMI Clock Trimming Register */
#define SFR_EBICFG	0x40	/* EBI Configuration Register */
#define SFR_CAN		0x48	/* CAN Memories Address-based Register */
#define SFR_SN0		0x4C	/* Serial Number 0 Register */
#define SFR_SN1		0x50	/* Serial Number 1 Register */
#define SFR_AICREDIR	0x54
#define SFR_L2CC_HRAMC	0x58

/*
 * Register Fields
 */

/*--- SFR_DDRCFG: (offset: 0x04) DDR Configuration Register ---*/
#define AT91C_DDRCFG_FDQIEN	(0x01 << 16)	/* Force DDR_DQ Input Buffer Always On */
#define AT91C_DDRCFG_FDQSIEN	(0x01 << 17)	/* Force DDF_DQS Input Buffer Always On */

/*--- SFR_EBICSA: (offset: 0x4) EBI Chip Select Register ---*/
#define AT91C_EBI_CS1A		(0x1UL << 1)
#define		AT91C_EBI_CS1A_SMC	(0x0UL << 1)
#define		AT91C_EBI_CS1A_SDRAMC	(0x1UL << 1)
#define AT91C_EBI_CS3A		(0x1UL << 3)
#define		AT91C_EBI_CS3A_SMC	(0x0UL << 3)
#define		AT91C_EBI_CS3A_SM	(0x1UL << 3)
#define AT91C_EBI_CS4A		(0x1UL << 4)
#define		AT91C_EBI_CS4A_SMC	(0x0UL << 4)
#define		AT91C_EBI_CS4A_SM	(0x1UL << 4)
#define AT91C_EBI_CS5A		(0x1UL << 5)
#define		AT91C_EBI_CS5A_SMC	(0x0UL << 5)
#define		AT91C_EBI_CS5A_SM	(0x1UL << 5)
#define AT91C_EBI_DBPUC		(0x1UL << 8)
#define AT91C_EBI_DBPDC		(0x1UL << 9)
#define AT91C_EBI_DRV		(0x1UL << 16)
#define		AT91C_EBI_DRV_LD	(0x0UL << 16)
#define		AT91C_EBI_DRV_HD	(0x1UL << 16)
#define AT91C_EBI_DQIENF	(0x1UL << 20)
#define AT91C_EBI_NFD0_ON_D16	(0x1UL << 24)
#define AT91C_EBI_DDR_MP_EN	(0x1UL << 25)

/*---SFR_EBICFG: (offset: 0x40) EBI Configuration Register ----*/
#define AT91C_EBICFG_DRIVE0	(0x03 << 0)
#define		AT91C_EBICFG_DRIVE0_LOW		(0x00 << 0)
#define		AT91C_EBICFG_DRIVE0_MEDIUM	(0x02 << 0)
#define		AT91C_EBICFG_DRIVE0_HIGH	(0x03 << 0)
#define AT91C_EBICFG_PULL0	(0x03 << 2)
#define		AT91C_EBICFG_PULL0_UP		(0x00 << 2)
#define		AT91C_EBICFG_PULL0_NONE		(0x01 << 2)
#define		AT91C_EBICFG_PULL0_DOWN		(0x03 << 2)
#define AT91C_EBICFG_SCH0	(0x01 << 4)
#define		AT91C_EBICFG_SCH0_OFF		(0x00 << 4)
#define		AT91C_EBICFG_SCH0_ON		(0x01 << 4)
#define AT91C_EBICFG_DRIVE1	(0x03 << 8)
#define		AT91C_EBICFG_DRIVE1_LOW		(0x00 << 8)
#define		AT91C_EBICFG_DRIVE1_MEDIUM	(0x02 << 8)
#define		AT91C_EBICFG_DRIVE1_HIGH	(0x03 << 8)
#define AT91C_EBICFG_PULL1	(0x03 << 10)
#define		AT91C_EBICFG_PULL1_UP		(0x00 << 10)
#define		AT91C_EBICFG_PULL1_NONE		(0x01 << 10)
#define		AT91C_EBICFG_PULL1_DOWN		(0x03 << 10)
#define AT91C_EBICFG_SCH1	(0x01 << 12)
#define		AT91C_EBICFG_SCH1_OFF		(0x00 << 12)
#define		AT91C_EBICFG_SCH1_ON		(0x01 << 12)
#define AT91C_EBICFG_BMS	(0x01 << 16)
#define		AT91C_EBICFG_BMS_ROM		(0x00 << 16)
#define		AT91C_EBICFG_BMS_EBI		(0x01 << 16)

#define AT91C_UTMICKTRIM_FREQ	0x03

/*---SFR_CAN: (offset: 0x48) CAN Memories Address-based Register ----*/
#define AT91C_CAN0_MEM_ADDR	(0xffff << 0)
#define AT91C_CAN0_MEM_ADDR_(addr)	(((addr) & 0xffff) << 0)
#define AT91C_CAN1_MEM_ADDR	(0xffff << 16)
#define AT91C_CAN1_MEM_ADDR_(addr)	(((addr) & 0xffff) << 16)

#endif /* #ifndef __AT91_SFR_H__ */

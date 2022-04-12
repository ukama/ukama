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
#ifndef __AT91_SDRAMC_H__
#define __AT91_SDRAMC_H__

/*** Register offset in AT91S_SDRAMC structure ***/
#define SDRAMC_MR	0x00	/* SDRAMC Mode Register */
#define SDRAMC_TR	0x04	/* SDRAMC Refresh Timer Register */
#define SDRAMC_CR	0x08	/* SDRAMC Configuration Register */
/* reserved 0x0c */
#define SDRAMC_LPR	0x10	/* SDRAMC Low Power Register */
#define SDRAMC_IER	0x14	/* SDRAMC Interrupt Enable Register */
#define SDRAMC_IDR	0x18	/* SDRAMC Interrupt Disable Register */
#define SDRAMC_IMR	0x1c	/* SDRAMC Interrupt Mask Register */
#define SDRAMC_ISR	0x20	/* SDRAMC Interrupt Mask Register */
#define SDRAMC_MDR	0x24	/* SDRAMC Memory Device Register */
#define SDRAMC_CFR1	0x28	/* SDRAMC Configuration Register 1 -- SAM9x60 */

/*-------- SDRAMC_MR : (SDRAMC Offset: 0x0) SDRAMC Mode Register -------*/
#define AT91C_SDRAMC_MODE	(0xF << 0)	/* SDRAMC Command Mode */
#define 	AT91C_SDRAMC_MODE_NORMAL		0x00
#define 	AT91C_SDRAMC_MODE_NOP			0x01
#define 	AT91C_SDRAMC_MODE_PRECHARGE		0x02
#define 	AT91C_SDRAMC_MODE_LOAD_MODE		0x03
#define 	AT91C_SDRAMC_MODE_AUTO_REFRESH		0x04
#define 	AT91C_SDRAMC_MODE_EXT_LOAD_MODE		0x05
#define 	AT91C_SDRAMC_MODE_DEEP_POWER		0x06

/*------- SDRAMC_TR : (SDRAMC Offset: 0x4) SDRAMC Refresh Timer Register --------*/
#define AT91C_SDRAMC_COUNT	(0xFFF << 0)	/* SDRAMC Refresh Timer Count */

/*-------- SDRAMC_CR : (SDRAMC Offset: 0x8) SDRAMC Configuration Register --------*/
#define AT91C_SDRAMC_NC		(0x3 << 0)	/* Number of Column Bits */
#define 	AT91C_SDRAMC_NC_8		0x0
#define 	AT91C_SDRAMC_NC_9		0x1
#define 	AT91C_SDRAMC_NC_10		0x2
#define 	AT91C_SDRAMC_NC_11		0x3
#define AT91C_SDRAMC_NR		(0x3 << 2)	/* Number of Row Bits */
#define 	AT91C_SDRAMC_NR_11		(0x0 << 2)
#define 	AT91C_SDRAMC_NR_12		(0x1 << 2)
#define 	AT91C_SDRAMC_NR_13		(0x2 << 2)
#define AT91C_SDRAMC_NB		(0x1 << 4)	/* Number of Banks */
#define 	AT91C_SDRAMC_NB_2_BANKS		(0x0 << 4)
#define 	AT91C_SDRAMC_NB_4_BANKS		(0x1 << 4)
#define AT91C_SDRAMC_CAS	(0x3 << 5)	/* CAS Latency */
#define 	AT91C_SDRAMC_CAS_1		(0x1 << 5)
#define 	AT91C_SDRAMC_CAS_2		(0x2 << 5)
#define 	AT91C_SDRAMC_CAS_3		(0x3 << 5)
#define AT91C_SDRAMC_DBW	(0x1 << 7)	/* Data Bus Width */
#define 	AT91C_SDRAMC_DBW_32_BITS	(0x0 << 7)
#define 	AT91C_SDRAMC_DBW_16_BITS	(0x1 << 7)
#define AT91C_SDRAMC_TWR	(0xF << 8)	/* Number of Write Recovery Time Cycles */
#define 	AT91C_SDRAMC_TWR_0		(0x0 << 8)
#define 	AT91C_SDRAMC_TWR_1		(0x1 << 8)
#define 	AT91C_SDRAMC_TWR_2		(0x2 << 8)
#define 	AT91C_SDRAMC_TWR_3		(0x3 << 8)
#define 	AT91C_SDRAMC_TWR_4		(0x4 << 8)
#define 	AT91C_SDRAMC_TWR_5		(0x5 << 8)
#define 	AT91C_SDRAMC_TWR_6		(0x6 << 8)
#define 	AT91C_SDRAMC_TWR_7		(0x7 << 8)
#define 	AT91C_SDRAMC_TWR_8		(0x8 << 8)
#define 	AT91C_SDRAMC_TWR_9		(0x9 << 8)
#define 	AT91C_SDRAMC_TWR_10		(0xA << 8)
#define 	AT91C_SDRAMC_TWR_11		(0xB << 8)
#define 	AT91C_SDRAMC_TWR_12		(0xC << 8)
#define 	AT91C_SDRAMC_TWR_13		(0xD << 8)
#define 	AT91C_SDRAMC_TWR_14		(0xE << 8)
#define 	AT91C_SDRAMC_TWR_15		(0xF << 8)
#define AT91C_SDRAMC_TRC	(0xF << 12)	/* Number of Row Cycle Delay Time Cycles */
#define 	AT91C_SDRAMC_TRC_0		(0x0 << 12)
#define 	AT91C_SDRAMC_TRC_1		(0x1 << 12)
#define 	AT91C_SDRAMC_TRC_2		(0x2 << 12)
#define 	AT91C_SDRAMC_TRC_3		(0x3 << 12)
#define 	AT91C_SDRAMC_TRC_4		(0x4 << 12)
#define 	AT91C_SDRAMC_TRC_5		(0x5 << 12)
#define 	AT91C_SDRAMC_TRC_6		(0x6 << 12)
#define 	AT91C_SDRAMC_TRC_7		(0x7 << 12)
#define 	AT91C_SDRAMC_TRC_8		(0x8 << 12)
#define 	AT91C_SDRAMC_TRC_9		(0x9 << 12)
#define 	AT91C_SDRAMC_TRC_10		(0xA << 12)
#define 	AT91C_SDRAMC_TRC_11		(0xB << 12)
#define 	AT91C_SDRAMC_TRC_12		(0xC << 12)
#define 	AT91C_SDRAMC_TRC_13		(0xD << 12)
#define 	AT91C_SDRAMC_TRC_14		(0xE << 12)
#define 	AT91C_SDRAMC_TRC_15		(0xF << 12)
#define AT91C_SDRAMC_TRP	(0xF << 16)	/* Number of Row Precharge Delay Time Cycles */
#define 	AT91C_SDRAMC_TRP_0		(0x0 << 16)
#define 	AT91C_SDRAMC_TRP_1		(0x1 << 16)
#define 	AT91C_SDRAMC_TRP_2		(0x2 << 16)
#define 	AT91C_SDRAMC_TRP_3		(0x3 << 16)
#define 	AT91C_SDRAMC_TRP_4		(0x4 << 16)
#define 	AT91C_SDRAMC_TRP_5		(0x5 << 16)
#define 	AT91C_SDRAMC_TRP_6		(0x6 << 16)
#define 	AT91C_SDRAMC_TRP_7		(0x7 << 16)
#define 	AT91C_SDRAMC_TRP_8		(0x8 << 16)
#define 	AT91C_SDRAMC_TRP_9		(0x9 << 16)
#define 	AT91C_SDRAMC_TRP_10		(0xA << 16)
#define 	AT91C_SDRAMC_TRP_11		(0xB << 16)
#define 	AT91C_SDRAMC_TRP_12		(0xC << 16)
#define 	AT91C_SDRAMC_TRP_13		(0xD << 16)
#define 	AT91C_SDRAMC_TRP_14		(0xE << 16)
#define 	AT91C_SDRAMC_TRP_15		(0xF << 16)
#define AT91C_SDRAMC_TRCD	(0xF << 20)	/* Number of Row to Column Delay Time Cycles */
#define 	AT91C_SDRAMC_TRCD_0		(0x0 << 20)
#define 	AT91C_SDRAMC_TRCD_1		(0x1 << 20)
#define 	AT91C_SDRAMC_TRCD_2		(0x2 << 20)
#define 	AT91C_SDRAMC_TRCD_3		(0x3 << 20)
#define 	AT91C_SDRAMC_TRCD_4		(0x4 << 20)
#define 	AT91C_SDRAMC_TRCD_5		(0x5 << 20)
#define 	AT91C_SDRAMC_TRCD_6		(0x6 << 20)
#define 	AT91C_SDRAMC_TRCD_7		(0x7 << 20)
#define 	AT91C_SDRAMC_TRCD_8		(0x8 << 20)
#define 	AT91C_SDRAMC_TRCD_9		(0x9 << 20)
#define 	AT91C_SDRAMC_TRCD_10		(0xA << 20)
#define 	AT91C_SDRAMC_TRCD_11		(0xB << 20)
#define 	AT91C_SDRAMC_TRCD_12		(0xC << 20)
#define 	AT91C_SDRAMC_TRCD_13		(0xD << 20)
#define 	AT91C_SDRAMC_TRCD_14		(0xE << 20)
#define 	AT91C_SDRAMC_TRCD_15		(0xF << 20)
#define AT91C_SDRAMC_TRAS	(0xF << 24)	/* Number of Active to Precharge Delay Time Cycles */
#define 	AT91C_SDRAMC_TRAS_0		(0x0 << 24)
#define 	AT91C_SDRAMC_TRAS_1		(0x1 << 24)
#define 	AT91C_SDRAMC_TRAS_2		(0x2 << 24)
#define 	AT91C_SDRAMC_TRAS_3		(0x3 << 24)
#define 	AT91C_SDRAMC_TRAS_4		(0x4 << 24)
#define 	AT91C_SDRAMC_TRAS_5		(0x5 << 24)
#define 	AT91C_SDRAMC_TRAS_6		(0x6 << 24)
#define 	AT91C_SDRAMC_TRAS_7		(0x7 << 24)
#define 	AT91C_SDRAMC_TRAS_8		(0x8 << 24)
#define 	AT91C_SDRAMC_TRAS_9		(0x9 << 24)
#define 	AT91C_SDRAMC_TRAS_10		(0xA << 24)
#define 	AT91C_SDRAMC_TRAS_11		(0xB << 24)
#define 	AT91C_SDRAMC_TRAS_12		(0xC << 24)
#define 	AT91C_SDRAMC_TRAS_13		(0xD << 24)
#define 	AT91C_SDRAMC_TRAS_14		(0xE << 24)
#define 	AT91C_SDRAMC_TRAS_15		(0xF << 24)
#define AT91C_SDRAMC_TXS	(0xF << 28)	/* Number of Exit Self Refresh to Active Delay Time Cycles */
#define 	AT91C_SDRAMC_TXSR_0		(0x0 << 28)
#define 	AT91C_SDRAMC_TXSR_1		(0x1 << 28)
#define 	AT91C_SDRAMC_TXSR_2		(0x2 << 28)
#define 	AT91C_SDRAMC_TXSR_3		(0x3 << 28)
#define 	AT91C_SDRAMC_TXSR_4		(0x4 << 28)
#define 	AT91C_SDRAMC_TXSR_5		(0x5 << 28)
#define 	AT91C_SDRAMC_TXSR_6		(0x6 << 28)
#define 	AT91C_SDRAMC_TXSR_7		(0x7 << 28)
#define 	AT91C_SDRAMC_TXSR_8		(0x8 << 28)
#define 	AT91C_SDRAMC_TXSR_9		(0x9 << 28)
#define 	AT91C_SDRAMC_TXSR_10		(0xA << 28)
#define 	AT91C_SDRAMC_TXSR_11		(0xB << 28)
#define 	AT91C_SDRAMC_TXSR_12		(0xC << 28)
#define 	AT91C_SDRAMC_TXSR_13		(0xD << 28)
#define 	AT91C_SDRAMC_TXSR_14		(0xE << 28)
#define 	AT91C_SDRAMC_TXSR_15		(0xF << 28)

/*-------- SDRAMC_LPR : (SDRAMC Offset: 0x10) SDRAMC Low-power Register --------*/
#define AT91C_SDRAMC_LPCB	(0x3 << 0)	/* Low-power Configurations */
#define 	AT91C_SDRAMC_LPCB_DISABLE		(0x0)
#define 	AT91C_SDRAMC_LPCB_SELF_REFRESH		(0x1)
#define 	AT91C_SDRAMC_LPCB_POWER_DOWN		(0x2)
#define 	AT91C_SDRAMC_LPCB_DEEP_POWER_DOWN	(0x3)
#define AT91C_SDRAMC_PASR	(0x7 << 4)	/* Partial Array Self Refresh (only for Low Power SDRAM) */
#define AT91C_SDRAMC_TCSR	(0x3 << 8)	/* Temperature Compensated Self Refresh (only for Low Power SDRAM) */
#define AT91C_SDRAMC_DS		(0x3 << 10)	/* Drive Strenght (only for Low Power SDRAM) */
#define AT91C_SDRAMC_TIMEOUT	(0x3 << 12)	/* Time to define when Low Power Mode is enabled */
#define 	AT91C_SDRAMC_TIMEOUT_0_CLK_CYCLES	(0x0 << 12)
#define 	AT91C_SDRAMC_TIMEOUT_64_CLK_CYCLES	(0x1 << 12)
#define 	AT91C_SDRAMC_TIMEOUT_128_CLK_CYCLES	(0x2 << 12)

/*-------- SDRAMC_IER : (SDRAMC Offset: 0x14) SDRAMC Interrupt Enable Register --------*/
#define AT91C_SDRAMC_RES	(0x1 << 0)	/* (Refresh Error Status */
/*-------- SDRAMC_IDR : (SDRAMC Offset: 0x18) SDRAMC Interrupt Disable Register -----*/
/*-------- SDRAMC_IMR : (SDRAMC Offset: 0x1c) SDRAMC Interrupt Mask Register --------*/
/*-------- SDRAMC_ISR : (SDRAMC Offset: 0x20) SDRAMC Interrupt Status Register ------*/

/*-------- SDRAMC_MDR : (SDRAMC Offset: 0x24) SDRAM Controller Memory Device Register -------*/
#define AT91C_SDRAMC_MD		(0x3 << 0)	/* Memory Device Type */
#define 	AT91C_SDRAMC_MD_SDRAM			0x0
#define 	AT91C_SDRAMC_MD_LOW_POWER_SDRAM		0x1
#define		AT91C_SDRAMC_MD_SHIFT_SAMPLING_1_CYCLE	(0x1 << 4)
#define		AT91C_SDRAMC_MD_SHIFT_SAMPLING_2_CYCLE	(0x2 << 4)
#define		AT91C_SDRAMC_MD_SHIFT_SAMPLING_3_CYCLE	(0x3 << 4)

/*-------- SDRAMC_CFR1 : (SDRAMC Offset: 0x28) SDRAM Controller Memory Configuration Register 1 - SAM9X60 -------*/
#define	AT91C_SDRAMC_CFR1_TMRD	(0xF << 0)
#define 	AT91C_SDRAMC_CFR1_TMRD_0		(0x0 << 0)
#define 	AT91C_SDRAMC_CFR1_TMRD_1		(0x1 << 0)
#define 	AT91C_SDRAMC_CFR1_TMRD_2		(0x2 << 0)
#define 	AT91C_SDRAMC_CFR1_TMRD_3		(0x3 << 0)
#define 	AT91C_SDRAMC_CFR1_TMRD_4		(0x4 << 0)
#define 	AT91C_SDRAMC_CFR1_TMRD_5		(0x5 << 0)
#define 	AT91C_SDRAMC_CFR1_TMRD_6		(0x6 << 0)
#define 	AT91C_SDRAMC_CFR1_TMRD_7		(0x7 << 0)
#define 	AT91C_SDRAMC_CFR1_TMRD_8		(0x8 << 0)
#define 	AT91C_SDRAMC_CFR1_TMRD_9		(0x9 << 0)
#define 	AT91C_SDRAMC_CFR1_TMRD_10		(0xA << 0)
#define 	AT91C_SDRAMC_CFR1_TMRD_11		(0xB << 0)
#define 	AT91C_SDRAMC_CFR1_TMRD_12		(0xC << 0)
#define 	AT91C_SDRAMC_CFR1_TMRD_13		(0xD << 0)
#define 	AT91C_SDRAMC_CFR1_TMRD_14		(0xE << 0)
#define 	AT91C_SDRAMC_CFR1_TMRD_15		(0xF << 0)

#define AT91C_SDRAMC_CFR1_UNAL	(0x1 << 8)
#define AT91C_SDRAMC_CFR1_ADD_DATA_MUX_SUPPORTED	(0x1 << 10)
#define AT91C_SDRAMC_CFR1_ADD_DATA_MUX_UNSUPPORTED	(0x0 << 10)
#define AT91C_SDRAMC_CFR1_CMD_MUX_SUPPORTED		(0x1 << 11)
#define AT91C_SDRAMC_CFR1_CMD_MUX_UNSUPPORTED		(0x0 << 11)
#endif /* #ifndef __AT91_SDRAMC_H__ */

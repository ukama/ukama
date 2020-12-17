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
#ifndef __AT91_MCI_H__
#define __AT91_MCI_H__

#define MCI_CR		0x00	/* MCI Control Register */
#define MCI_MR		0x04	/* MCI Mode Register */
#define MCI_DTOR	0x08	/* MCI Data Timeout Register */
#define MCI_SDCR	0x0C	/* MCI SD/SDIO Card Register */
#define MCI_ARGR	0x10	/* MCI Argument Register */
#define MCI_CMDR	0x14	/* MCI Command Register */
#define MCI_BLKR	0x18	/* MCI Block Register */
#define MCI_CSTOR	0x1C	/* MCI Completion Signal Timeout Register */
#define MCI_RSPR	0x20	/* MCI Response Register */
#define MCI_RSPR1	0x24	/* MCI Response Register1 */
#define MCI_RSPR2	0x28	/* MCI Response Register2 */
#define MCI_RSPR3	0x2C	/* MCI Response Register3 */
#define MCI_RDR		0x30	/* MCI Receive Data Register */
#define MCI_TDR		0x34	/* MCI Transmit Data Register */
#define MCI_SR		0x40	/* MCI Status Register */
#define MCI_IER		0x44	/* MCI Interrupt Enable Register */
#define MCI_IDR		0x48	/* MCI Interrupt Disable Register */
#define MCI_IMR		0x4C	/* MCI Interrupt Mask Register */
#define MCI_DMA		0x50	/* MCI DMA Configuration Register */
#define MCI_CFG		0x54	/* MCI Configuration Register */

#define MCI_WPCR	0xE4	/* MCI Write Protection Control Register */
#define MCI_WPSR	0xE8	/* MCI Write Protection Status Register */

#define MCI_FIFO	0x200	/* MCI FIFO Aperture Register */

/*-------- MCI_CR : (MCI Offset: 0x0) MCI Control Register --------*/
#define AT91C_MCI_MCIEN		(0x1UL << 0)	/* Multimedia Interface Enable*/
#define AT91C_MCI_MCIDIS	(0x1UL << 1)	/* Multimedia Interface Disable */
#define AT91C_MCI_PWSEN		(0x1UL << 2)	/* Power Save Mode Enable */
#define AT91C_MCI_PWSDIS	(0x1UL << 3)	/* Power Save Mode Disable */
#define AT91C_MCI_IOWAITEN	(0x1UL << 4)	/* SDIO Read Wait Enable */
#define AT91C_MCI_IOWAITDIS	(0x1UL << 5)	/* SDIO Read Wait Disable */
#define AT91C_MCI_SWRST		(0x1UL << 7)	/* MCI Software rese */

/*-------- MCI_MR : (MCI Offset: 0x4) MCI Mode Register --------*/
#define AT91C_MCI_CLKDIV(x)	(x << 0)	/* Clock Divider */
#define AT91C_MCI_PWSDIV(x)	(x <<  8)	/* Power Saving Divider */
#define AT91C_MCI_RDPROOF	(0x1UL << 11)	/* Read Proof Enable */
#define 	AT91C_MCI_RDPROOF_DISABLE	(0x0UL << 11)
#define 	AT91C_MCI_RDPROOF_ENABLE	(0x1UL << 11)
#define AT91C_MCI_WRPROOF	(0x1UL << 12)	/* Write Proof Enable */
#define 	AT91C_MCI_WRPROOF_DISABLE	(0x0UL << 12)
#define 	AT91C_MCI_WRPROOF_ENABLE 	(0x1UL << 12)
#define AT91C_MCI_PDCFBYTE	(0x1UL << 13)	/* PDC Force Byte Transfer */
#define 	AT91C_MCI_PDCFBYTE_DISABLE	(0x0UL << 13)
#define 	AT91C_MCI_PDCFBYTE_ENABLE	(0x1UL << 13)
#define AT91C_MCI_PDCPADV	(0x1UL << 14)	/* PDC Padding Value */
#define AT91C_MCI_PDCMODE	(0x1UL << 15)	/* PDC Oriented Mode */
#define 	AT91C_MCI_PDCMODE_DISABLE	(0x0UL << 15)
#define 	AT91C_MCI_PDCMODE_ENABLE	(0x1UL << 15)
#define AT91C_MCI_CLKODD	(0x1UL << 16)	/* CLKODD */
#define AT91C_MCI_MRBLKLEN(x)	(x << 16)	/* Data Block Length */

/*-------- MCI_DTOR : (MCI Offset: 0x8) MCI Data Timeout Register --------*/
#define AT91C_MCI_DTOCYC(x)	(x << 0)	/* Data Timeout Cycle Number */
#define AT91C_MCI_DTOMUL	(0x7UL <<  4)	/* Data Timeout Multiplier */
#define 	AT91C_MCI_DTOMUL_1		(0x0UL << 4)
#define 	AT91C_MCI_DTOMUL_16		(0x1UL << 4)
#define 	AT91C_MCI_DTOMUL_128		(0x2UL << 4)
#define 	AT91C_MCI_DTOMUL_256		(0x3UL << 4)
#define 	AT91C_MCI_DTOMUL_1024		(0x4UL << 4)
#define 	AT91C_MCI_DTOMUL_4096		(0x5UL << 4)
#define 	AT91C_MCI_DTOMUL_65536		(0x6UL << 4)
#define 	AT91C_MCI_DTOMUL_1048576	(0x7UL << 4)

/*-------- MCI_SDCR : (MCI Offset: 0xc) MCI SD Card Register --------*/
#define AT91C_MCI_SCDSEL	(0x3UL << 0)	/* SD Card/SDIO Selector */
#define 	AT91C_MCI_SCDSEL_SLOTA		(0x0UL)
#define 	AT91C_MCI_SCDSEL_SLOTB		(0x1UL)
#define 	AT91C_MCI_SCDSEL_SLOTC		(0x2UL)
#define 	AT91C_MCI_SCDSEL_SLOTD		(0x3UL)
#define AT91C_MCI_SCDBUS	(0x3UL << 6)	/* SDCard/SDIO Bus Width */
#define 	AT91C_MCI_SCDBUS_1BIT		(0x0UL << 6)
#define 	AT91C_MCI_SCDBUS_4BIT		(0x2UL << 6)
#define 	AT91C_MCI_SCDBUS_8BIT		(0x3UL << 6)

/*-------- MCI_CMDR : (MCI Offset: 0x14) MCI Command Register --------*/
#define AT91C_MCI_CMDNB(x)	(x << 0)	/* Command Number */
#define AT91C_MCI_RSPTYP	(0x3UL << 6)	/* Response Type */
#define 	AT91C_MCI_RSPTYP_NO		(0x0UL << 6)
#define 	AT91C_MCI_RSPTYP_48		(0x1UL << 6)
#define 	AT91C_MCI_RSPTYP_136		(0x2UL << 6)
#define 	AT91C_MCI_RSPTYP_R1B		(0x3UL << 6)
#define AT91C_MCI_SPCMD		(0x7UL << 8)	/* Special CMD */
#define 	AT91C_MCI_SPCMD_NONE		(0x0UL << 8)
#define 	AT91C_MCI_SPCMD_INIT		(0x1UL << 8)
#define 	AT91C_MCI_SPCMD_SYNC		(0x2UL << 8)
#define 	AT91C_MCI_SPCMD_CE_ATA		(0x3UL << 8)
#define 	AT91C_MCI_SPCMD_IT_CMD		(0x4UL << 8)
#define 	AT91C_MCI_SPCMD_IT_REP		(0x5UL << 8)
#define 	AT91C_MCI_SPCMD_BOR		(0x6UL << 8)
#define 	AT91C_MCI_SPCMD_EBO		(0x7UL << 8)
#define AT91C_MCI_OPDCMD	(0x1UL << 11)	/* Open Drain Command */
#define 	AT91C_MCI_OPDCMD_PUSHPULL	(0x0UL << 11)
#define 	AT91C_MCI_OPDCMD_OPENDRAIN	(0x1UL << 11)
#define AT91C_MCI_MAXLAT	(0x1UL << 12)	/* Maximum Latency for Command to respond */
#define 	AT91C_MCI_MAXLAT_5		(0x0UL << 12)
#define 	AT91C_MCI_MAXLAT_64		(0x1UL << 12)
#define AT91C_MCI_TRCMD		(0x3UL << 16)	/* Transfer CMD */
#define 	AT91C_MCI_TRCMD_NO		(0x0UL << 16)
#define 	AT91C_MCI_TRCMD_START		(0x1UL << 16)
#define 	AT91C_MCI_TRCMD_STOP		(0x2UL << 16)
#define AT91C_MCI_TRDIR		(0x1UL << 18)	/* Transfer Direction */
#define 	AT91C_MCI_TRDIR_WRITE		(0x0UL << 18)
#define 	AT91C_MCI_TRDIR_READ		(0x1UL << 18)
#define AT91C_MCI_TRTYP		(0x7UL << 19)	/* Transfer Type */
#define 	AT91C_MCI_TRTYP_SINGLE		(0x0UL << 19)
#define 	AT91C_MCI_TRTYP_MULTIPLE	(0x1UL << 19)
#define 	AT91C_MCI_TRTYP_STREAM		(0x2UL << 19)  
#define 	AT91C_MCI_TRTYP_SDIO_BYTE	(0x4UL << 19)
#define 	AT91C_MCI_TRTYP_SDIO_BLOCK	(0x5UL << 19)
#define AT91C_MCI_IOSPCMD	(0x3UL << 24)	/* SDIO Special Command */
#define 	AT91C_MCI_IOSPCMD_NONE		(0x0UL << 24)
#define 	AT91C_MCI_IOSPCMD_SUSPEND	(0x1UL << 24)
#define 	AT91C_MCI_IOSPCMD_RESUME	(0x2UL << 24)
#define AT91C_MCI_ATACS		(0x1UL << 26)	/* ATA with command completion signal */
#define 	AT91C_MCI_ATACS_NORMAL		(0x0UL << 26)
#define 	AT91C_MCI_ATACS_COMPLETION	(0x1UL << 26)
#define AT91C_MCI_BOOTACK	(0x1UL << 27)	/* Boot Operation Acknowledge */

/*-------- MCI_BLKR : (MCI Offset: 0x18) MCI Block Register --------*/
#define AT91C_MCI_BCNT(x)	(x << 0)	/* (MCI) MMC/SDIO Block Count / SDIO Byte Count */
#define AT91C_MCI_BLKLEN(x)	(x << 16)	/* (MCI) Data Block Length */

/*-------- MCI_CSTOR : (MCI Offset: 0x1c) MCI Completion Signal Timeout Register -------*/
#define AT91C_MCI_CSTOCYC(x)	(x << 0)	/* (MCI) Completion Signal Timeout Cycle Number */
#define AT91C_MCI_CSTOMUL	(0x7UL << 4)	/* (MCI) Completion Signal Timeout Multiplier */
#define 	AT91C_MCI_CSTOMUL_1		(0x0UL << 4)
#define 	AT91C_MCI_CSTOMUL_16		(0x1UL << 4)
#define 	AT91C_MCI_CSTOMUL_128		(0x2UL << 4)
#define 	AT91C_MCI_CSTOMUL_256		(0x3UL << 4)
#define 	AT91C_MCI_CSTOMUL_1024		(0x4UL << 4)
#define 	AT91C_MCI_CSTOMUL_4096		(0x5UL << 4)
#define 	AT91C_MCI_CSTOMUL_65536		(0x6UL << 4)
#define 	AT91C_MCI_CSTOMUL_1048576	(0x7UL << 4)

/*-------- MCI_SR : (MCI Offset: 0x40) MCI Status Register --------*/
#define AT91C_MCI_CMDRDY	(0x1UL <<  0)	/* Command Ready flag */
#define AT91C_MCI_RXRDY		(0x1UL <<  1)	/* RX Ready flag */
#define AT91C_MCI_TXRDY		(0x1UL <<  2)	/* TX Ready flag */
#define AT91C_MCI_BLKE		(0x1UL <<  3)	/* Data Block Transfer Ended flag */
#define AT91C_MCI_DTIP		(0x1UL <<  4)	/* Data Transfer in Progress flag */
#define AT91C_MCI_NOTBUSY	(0x1UL <<  5)	/* Data Line Not Busy flag */
#define AT91C_MCI_ENDRX		(0x1UL <<  6)	/* End of RX Buffer flag */
#define AT91C_MCI_ENDTX		(0x1UL <<  7)	/* End of TX Buffer flag */
#define AT91C_MCI_SDIOIRQA	(0x1UL <<  8)	/* SDIO Interrupt for Slot A */
#define AT91C_MCI_SDIOIRQB	(0x1UL <<  9)	/* SDIO Interrupt for Slot B */
#define AT91C_MCI_SDIOIRQC	(0x1UL << 10)	/* SDIO Interrupt for Slot C */
#define AT91C_MCI_SDIOIRQD	(0x1UL << 11)	/* SDIO Interrupt for Slot D */
#define AT91C_MCI_SDIOWAIT	(0x1UL << 12)	/* SDIO Read Wait operation flag */
#define AT91C_MCI_CSRCV		(0x1UL << 13)	/* CE-ATA Completion Signal flag */
#define AT91C_MCI_RXBUFF	(0x1UL << 14)	/* RX Buffer Full flag */
#define AT91C_MCI_TXBUFE	(0x1UL << 15)	/* TX Buffer Empty flag */
#define AT91C_MCI_RINDE		(0x1UL << 16)	/* Response Index Error flag */
#define AT91C_MCI_RDIRE		(0x1UL << 17)	/* Response Direction Error flag */
#define AT91C_MCI_RCRCE		(0x1UL << 18)	/* Response CRC Error flag */
#define AT91C_MCI_RENDE		(0x1UL << 19)	/* Response End Bit Error flag */
#define AT91C_MCI_RTOE		(0x1UL << 20)	/* Response Time-out Error flag */
#define AT91C_MCI_DCRCE		(0x1UL << 21)	/* data CRC Error flag */
#define AT91C_MCI_DTOE		(0x1UL << 22)	/* Data timeout Error flag */
#define AT91C_MCI_CSTOE		(0x1UL << 23)	/* Completion Signal timeout Error flag */
#define AT91C_MCI_BLKOVRE	(0x1UL << 24)	/* DMA Block Overrun Error flag */
#define AT91C_MCI_DMADONE	(0x1UL << 25)	/* DMA Transfer Done flag */
#define AT91C_MCI_FIFOEMPTY	(0x1UL << 26)	/* FIFO Empty flag */
#define AT91C_MCI_XFRDONE	(0x1UL << 27)	/* Transfer Done flag */
#define AT91C_MCI_OVRE		(0x1UL << 30)	/* Overrun flag */
#define AT91C_MCI_UNRE		(0x1UL << 31)	/* Underrun flag */
/*-------- MCI_IER : (MCI Offset: 0x44) MCI Interrupt Enable Register --------*/
/*-------- MCI_IDR : (MCI Offset: 0x48) MCI Interrupt Disable Register -------*/
/*-------- MCI_IMR : (MCI Offset: 0x4c) MCI Interrupt Mask Register ----------*/

/*-------- MCI_DMA : (MCI Offset: 0x50) MCI DMA Configuration Register --------*/
#define AT91C_MCI_OFFSET	(0x3UL <<  0)	/* DMA Write Buffer Offset */
#define AT91C_MCI_CHKSIZE	(0x7UL <<  4)	/* DMA Channel Read/Write Chunk Size */
#define 	AT91C_MCI_CHKSIZE_1		(0x0UL <<  4)
#define 	AT91C_MCI_CHKSIZE_4		(0x1UL <<  4)
#define 	AT91C_MCI_CHKSIZE_8		(0x2UL <<  4)
#define 	AT91C_MCI_CHKSIZE_16		(0x3UL <<  4)
#define 	AT91C_MCI_CHKSIZE_32		(0x4UL <<  4)
#define AT91C_MCI_DMAEN		(0x1UL <<  8)	/* DMA Hardware Handshaking Enable */
#define 	AT91C_MCI_DMAEN_DISABLE		(0x0UL <<  8)
#define 	AT91C_MCI_DMAEN_ENABLE		(0x1UL <<  8)

/*-------- MCI_CFG : (MCI Offset: 0x54) MCI Configuration Register --------*/
#define AT91C_MCI_FIFOMODE	(0x1UL <<  0)	/* MCI Internal FIFO Control Mode */
#define 	AT91C_MCI_FIFOMODE_AMOUNTDATA	(0x0UL)
#define 	AT91C_MCI_FIFOMODE_ONEDATA	(0x1UL)
#define AT91C_MCI_FERRCTRL	(0x1UL <<  4)	/* Flow Error Flag Reset Control Mode */
#define 	AT91C_MCI_FERRCTRL_RWCMD	(0x0UL << 4)
#define 	AT91C_MCI_FERRCTRL_READSR	(0x1UL << 4)
#define AT91C_MCI_HSMODE	(0x1UL <<  8)	/* High Speed Mode */
#define 	AT91C_MCI_HSMODE_DISABLE	(0x0UL << 8)
#define 	AT91C_MCI_HSMODE_ENABLE		(0x1UL << 8)
#define AT91C_MCI_LSYNC		(0x1UL << 12)	/* Synchronize on last block */
#define 	AT91C_MCI_LSYNC_CURRENT		(0x0UL << 12)
#define 	AT91C_MCI_LSYNC_INFINITE	(0x1UL << 12)

#endif	/* #ifndef __AT91_MCI_H__ */

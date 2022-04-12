/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2015, Atmel Corporation
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
#ifndef __AT91_AES_H__
#define __AT91_AES_H__

/**** Register offset in AT91_AES structure ***/
#define AES_CR		0x00	/* Control Register */
#define AES_MR		0x04	/* Mode Register */
/* 0x08 - 0x0C Reserved */
#define AES_IER		0x10	/* Interrupt Enable Register */
#define AES_IDR		0x14	/* Interrupt Disable Register */
#define AES_IMR		0x18	/* Interrupt Mask Register */
#define AES_ISR		0x1C	/* Interrupt Status Register */
#define AES_KEYWR0	0x20	/* Key Word Register 0 */
#define AES_KEYWR1	0x24	/* Key Word Register 1 */
#define AES_KEYWR2	0x28	/* Key Word Register 2 */
#define AES_KEYWR3	0x2C	/* Key Word Register 3 */
#define AES_KEYWR4	0x30	/* Key Word Register 4 */
#define AES_KEYWR5	0x34	/* Key Word Register 5 */
#define AES_KEYWR6	0x38	/* Key Word Register 6 */
#define AES_KEYWR7	0x3C	/* Key Word Register 7 */
#define AES_IDATAR0	0x40	/* Input Data Register 0 */
#define AES_IDATAR1	0x44	/* Input Data Register 1 */
#define AES_IDATAR2	0x48	/* Input Data Register 2 */
#define AES_IDATAR3	0x4C	/* Input Data Register 3 */
#define AES_ODATAR0	0x50	/* Output Data Register 0 */
#define AES_ODATAR1	0x54	/* Output Data Register 1 */
#define AES_ODATAR2	0x58	/* Output Data Register 2 */
#define AES_ODATAR3	0x5C	/* Output Data Register 3 */
#define AES_IVR0	0x60	/* Initialization Vector Register 0 */
#define AES_IVR1	0x64	/* Initialization Vector Register 1 */
#define AES_IVR2	0x68	/* Initialization Vector Register 2 */
#define AES_IVR3	0x6C	/* Initialization Vector Register 3 */
#define AES_AADLENR	0x70	/* Additional Authenticated Data Length Register */
#define AES_CLENR	0x74	/* Plaintex/Ciphertext Length Register */
#define AES_GHASHR0	0x78	/* GCM Intermediate Hash Word Register 0 */
#define AES_GHASHR1	0x7C	/* GCM Intermediate Hash Word Register 1 */
#define AES_GHASHR2	0x80	/* GCM Intermediate Hash Word Register 2 */
#define AES_GHASHR3	0x84	/* GCM Intermediate Hash Word Register 3 */
#define AES_TAGR0	0x88	/* GCM Authentication Tag Word Register 0 */
#define AES_TAGR1	0x8C	/* GCM Authentication Tag Word Register 1 */
#define AES_TAGR2	0x90	/* GCM Authentication Tag Word Register 2 */
#define AES_TAGR3	0x94	/* GCM Authentication Tag Word Register 3 */
#define AES_CTRR	0x98	/* GCM Encryption Counter Value Register */
#define AES_GCMHR0	0x9C	/* GCM H World Register 0 */
#define AES_GCMHR1	0xA0	/* GCM H World Register 1 */
#define AES_GCMHR2	0xA4	/* GCM H World Register 2 */
#define AES_GCMHR3	0xA8	/* GCM H World Register 3 */

/*-------- AES_CR : (Offset: 0x00) Control Register --------*/
#define	AES_CR_START		(0x1UL << 0)	/* Start Processing */
#define	AES_CR_SWRST		(0x1UL << 8)	/* Software Reset */

/*-------- AES_MR : (Offset: 0x04) Mode Register --------*/
#define	AES_MR_CIPHER		(0x1UL << 0)	/* Processing Mode */
#define		AES_MR_CIPHER_DECRYPT	(0x0UL << 0)
#define		AES_MR_CIPHER_ENCRYPT	(0x1UL << 0)
#define	AES_MR_GTAGEN		(0x1UL << 1)	/* GCM Automatic Tag Generation Enable */
#define	AES_MR_DUALBUFF		(0x1UL << 3)	/* Dual Input Buffer */
#define	AES_MR_PROCDLY_Pos	4		/* Processing Delay */
#define	AES_MR_PROCDLY_Msk	(0xFUL << AES_MR_PROCDLY_Pos)
#define	AES_MR_PROCDLY(value)	(((value) << AES_MR_PROCDLY_Pos) & AES_MR_PROCDLY_Msk)
#define	AES_MR_SMOD		(0x3UL << 8)	/* Start Mode */
#define		AES_MR_SMOD_MANUAL_START	(0x0UL << 8)	/* Manual Mode */
#define		AES_MR_SMOD_AUTO_START		(0x1UL << 8)	/* Auto Mode */
#define		AES_MR_SMOD_IDATAR0_START	(0x2UL << 8)	/* AES_IDATAR0 access only Auto Mode */
#define	AES_MR_KEYSIZE		(0x3UL << 10)	/* Key Size */
#define		AES_MR_KEYSIZE_AES128		(0x0UL << 10)
#define		AES_MR_KEYSIZE_AES192		(0x1UL << 10)
#define		AES_MR_KEYSIZE_AES256		(0x2UL << 10)
#define	AES_MR_OPMOD		(0x7UL << 12)	/* Operation Mode */
#define		AES_MR_OPMOD_ECB		(0x0UL << 12)	/* Electronic Code Book */
#define		AES_MR_OPMOD_CBC		(0x1UL << 12)	/* Cipher Block Chaining */
#define		AES_MR_OPMOD_OFB		(0x2UL << 12)	/* Output Feedback */
#define		AES_MR_OPMOD_CFB		(0x3UL << 12)	/* Cipher Feedback */
#define		AES_MR_OPMOD_CTR		(0x4UL << 12)	/* Counter */
#define		AES_MR_OPMOD_GCM		(0x5UL << 12)	/* Galois Counter */
#define	AES_MR_LOD		(0x1UL << 15)	/* Last Output Data Mode */
#define	AES_MR_CFBS		(0x7UL << 16)	/* Cipher Feedback Data Size */
#define		AES_MR_CFBS_128BIT	(0x0UL << 16)
#define		AES_MR_CFBS_64BIT	(0x1UL << 16)
#define		AES_MR_CFBS_32BIT	(0x2UL << 16)
#define		AES_MR_CFBS_16BIT	(0x3UL << 16)
#define		AES_MR_CFBS_8BIT	(0x4UL << 16)
#define	AES_MR_CKEY		(0xFUL << 20)	/* Key */
#define		AES_MR_CKEY_PASSWD	(0xEUL << 20)

/*-------- AES_ISR : (Offset: 0x1C) Interrupt Status Register --------*/
#define	AES_INT_DATRDY		(0x1UL << 0)	/* Data Ready */
#define	AES_INT_URAD		(0x1UL << 8)	/* Unspecified Register Access Detection Status */
#define	AES_INT_URAT		(0xFUL << 12)	/* Unspecified Register Access */
#define		AES_INT_IDR_WR_PROCESSING	(0x0UL << 12)
#define		AES_INT_ODR_RD_PROCESSING	(0x1UL << 12)
#define		AES_INT_MR_WR_PROCESSING	(0x2UL << 12)
#define		AES_INT_ODR_RD_SUBKGEN		(0x3UL << 12)
#define		AES_INT_MR_WR_SUBKGEN		(0x4UL << 12)
#define		AES_INT_WOR_RD_ACCESS		(0x5UL << 12)
#define	AES_INT_TAGRDY		(0x1UL << 16)	/* GCM Tag Ready */

#endif /* #ifndef __AT91_AES_H__ */

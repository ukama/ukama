/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2008, Atmel Corporation
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
 *
 * 2017 Modified for Ebee Smart Technologies, Berlin
 * Author: Jens Thoms Toerring <jt@toerring.de>
 */

#include "hamming.h"

/*
 * Note: if the Linux kernel is build with CONFIG_MTD_NAND_ECC_SMC
 * set (following the Smart Media Specificatio) the first byte of
 * an ECC contains the lower bits of the line parity (LP00 to LP07)
 * and the second byte the upper bits (LP08 to LP15). Otherwise
 * (and that's the most common case) it's the other way round. So
 * if CONFIG_MTD_NAND_ECC_SMC is set in the Linux configuration
 * define LP00_07 as 0 and LP08_15 as 1, otherwise LP00_07 must be
 * set to 1 and LP08_17 to 0.
 *
 * TODO: make this a configuration setting?
 */

#if 1
/* For kernel configured without CONFIG_MTD_NAND_ECC_SMC (normal case) */
#define LP00_07   1
#define LP08_15   0
#else
/* For kernel configured with CONFIG_MTD_NAND_ECC_SMC */
#define LP00_07   0
#define LP08_15   1
#endif

/*
 * Returns the number of bits set in a byte. AT91SAM9260EK
 * may not have enough memory for a look-up table.
 */
#if defined CONFIG_AT91SAM9260EK
static unsigned char bits_in_byte(unsigned char byte)
{
	unsigned char cnt;
	for (cnt = 0; byte; cnt++)
		byte &= byte - 1;
	return cnt;
}
#else
static const unsigned char BitsSetTable256[256] = {
	#define B2(n) n,     n + 1,     n + 1,     n + 2
	#define B4(n) B2(n), B2(n + 1), B2(n + 1), B2(n + 2)
	#define B6(n) B4(n), B4(n + 1), B4(n + 1), B4(n + 2)
	B6(0), B6(1), B6(1), B6(2)
};

static inline unsigned char bits_in_byte(unsigned char byte)
{
	return BitsSetTable256[byte];
}
#endif

/*
 * Returns the number of bits set in an ECC
 */
static inline unsigned char bits_in_ecc(unsigned char ecc[3])
{
	return   bits_in_byte(ecc[0])
	       + bits_in_byte(ecc[1])
	       + bits_in_byte(ecc[2]);
}

/*
 * Returns if the number of bits set in a byte is odd
 */
static inline unsigned char has_odd_bits(unsigned char byte)
{
	return bits_in_byte(byte) & 1;
}

/*
 * Calculates the Hamming ECC of a 256 byte block of
 * data, returned via 'ecc'.
 */
static void compute_ecc_256(const unsigned char data[256],
			    unsigned char ecc[3])
{
	static const unsigned char tbl[16] =
		             {0x55, 0x56, 0x59, 0x5a, 0x65, 0x66, 0x69, 0x6a,
			      0x95, 0x96, 0x99, 0x9a, 0xa5, 0xa6, 0xa9, 0xaa};
	unsigned char cnt = 0;

	ecc[0] = ecc[1] = ecc[2] = 0xff;

	do {
		/* Only data bytes with odd number of set bits influence
		 * the line parity */
		if (has_odd_bits(*data)) {
			ecc[LP00_07] ^= tbl[cnt & 0x0f];
			ecc[LP08_15] ^= tbl[cnt >> 4];
		}

		ecc[2] ^= *data++;           /* column parity */
	} while (++cnt);

	/* Finalize column parity calculation */
	ecc[2] = ~ (  (has_odd_bits(ecc[2] & 0x55) << 2)
		    | (has_odd_bits(ecc[2] & 0xaa) << 3)
		    | (has_odd_bits(ecc[2] & 0x33) << 4)
		    | (has_odd_bits(ecc[2] & 0xcc) << 5)
		    | (has_odd_bits(ecc[2] & 0x0f) << 6)
		    | (has_odd_bits(ecc[2] & 0xf0) << 7));
}

/*
 * Corrects a single bit flip in 'data' using the
 * correction ECC (i.e., the Xor of the ECC from
 * NAND and the ECC calculated from the data).
 */
static inline void correct_bit_256(unsigned char data[256],
				   const unsigned char ecc[3])
{
	/* Position of flipped bit in byte */
	unsigned char bit =   ((ecc[2] >> 5) & 0x04)
		            | ((ecc[2] >> 4) & 0x02)
		            | ((ecc[2] >> 3) & 0x01);

	/* Adress of byte with flipped bit */
	data +=   ((ecc[LP08_15] << 0) & 0x80)
	        | ((ecc[LP08_15] << 1) & 0x40)
	        | ((ecc[LP08_15] << 2) & 0x20)
	        | ((ecc[LP08_15] << 3) & 0x10)
	        | ((ecc[LP00_07] >> 4) & 0x08)
	        | ((ecc[LP00_07] >> 3) & 0x04)
	        | ((ecc[LP00_07] >> 2) & 0x02)
	        | ((ecc[LP00_07] >> 1) & 0x01);

	/* "Un-flip" that bit */
	*data ^= 1 << bit;
}

/*
 * Compares the ECC read from the NAND with those
 * calculated from the data and, if possible, corrects
 * single bit flips.
 */
static unsigned char verify_256(unsigned char data[256],
				const unsigned char ecc_from_nand[3])
{
	/* Calculate ECC of data */
	unsigned char ecc[3];
	compute_ecc_256(data, ecc);

	/* Xor calculated ECC with that read from the NAND */
	ecc[0] ^= ecc_from_nand[0];
	ecc[1] ^= ecc_from_nand[1];
	ecc[2] ^= ecc_from_nand[2];

	/* Number of bits in result tells if the data are ok. A single
	 * bit flip gets repaired. */
	switch (bits_in_ecc(ecc)) {
	case 0:                             /* everything's fine */
		return 0;

	case 11:                            /* correctable single bit flip */
		correct_bit_256(data, ecc);
		return Hamming_ERROR_SINGLEBIT;

	case 1:                             /* flipped bit in ECC */
		return Hamming_ERROR_ECC;

	default:                            /* more than one bit flipped */
		return Hamming_ERROR_MULTIPLEBITS;
	}
}

/*
 * Checks the data in a page using the ECC values read
 * from the NAND. Single bit flips in slabs of 256 bytes
 * get corrected.
 * Returns 0 if no bit flips were detected, otherwise
 * a positive number (with 'Hamming_ERROR_SINGLEBIT'
 * indicating that there was at least one bit flip
 * but the data got corrected).
 */
unsigned char Hamming_Verify256x(unsigned char *data,
				 unsigned int size,
				 const unsigned char *ecc_from_nand)
{
	unsigned char result = 0;

	while (size > 0) {
		unsigned char error = verify_256(data, ecc_from_nand);

		if (error == Hamming_ERROR_SINGLEBIT)
			result = Hamming_ERROR_SINGLEBIT;
		else if (error)
			return error;

		data          += 256;
		ecc_from_nand += 3;
		size          -= 256;
	}

	return result;
}

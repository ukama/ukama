/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/src/ecc.c $                                              */
/*                                                                        */
/* OpenPOWER FFS Project                                                  */
/*                                                                        */
/* Contributors Listed Below - COPYRIGHT 2014,2015                        */
/* [+] International Business Machines Corp.                              */
/*                                                                        */
/*                                                                        */
/* Licensed under the Apache License, Version 2.0 (the "License");        */
/* you may not use this file except in compliance with the License.       */
/* You may obtain a copy of the License at                                */
/*                                                                        */
/*     http://www.apache.org/licenses/LICENSE-2.0                         */
/*                                                                        */
/* Unless required by applicable law or agreed to in writing, software    */
/* distributed under the License is distributed on an "AS IS" BASIS,      */
/* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or        */
/* implied. See the License for the specific language governing           */
/* permissions and limitations under the License.                         */
/*                                                                        */
/* IBM_PROLOG_END_TAG                                                     */

/*
 *   File: ecc.c
 * Author: Shaun Wetzstein <shaun@us.ibm.com>
 *  Descr: SFC ECC functions
 *   Note:
 *   Date: 08/02/12
 *  Descr: Added New ECC function with correctable bit functionality.
 *   Date: 12/04/13
 */

#include <unistd.h>
#include <stdarg.h>
#include <stdlib.h>
#include <stdbool.h>
#include <malloc.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <limits.h>
#include <ctype.h>
#include <endian.h>
#include <assert.h>

#include "ecc.h"
#include "clib/builtin.h"
#include "attribute.h"
#include "misc.h"

/*
 * This is an alternative way to calculate the ECC byte taken
 * from the SFC spec.
 */
uint8_t sfc_ecc2(uint8_t __data[8]) __unused__;
uint8_t sfc_ecc2(uint8_t __data[8])
{
	uint8_t ecc = 0;

	for (int byte = 0; byte < 8; byte++) {
		for (int bit = 0; bit < 8; bit++) {
			static unsigned char m[] =
			    { 0xff, 0x00, 0x00, 0xe8, 0x42, 0x3c, 0x0f, 0x99 };

			unsigned char x =
			    __data[byte] & m[(byte + bit + 1) & 7];
			x = x ^ (x >> 4);
			x = x ^ (x >> 2);
			x = x ^ (x >> 1);

			ecc ^= (x & 1) << bit;
		}
	}

	return ~ecc;
}

uint8_t sfc_ecc(uint8_t __data[8])
{
/* each bit of the ECC data corresponds to a row in this matrix */
	static uint8_t __matrix[8][8] = {
/* 11111111 00000000 00000000 11101000 01000010 00111100 00001111 10011001 */
		{0xff, 0x00, 0x00, 0xe8, 0x42, 0x3c, 0x0f, 0x99},
/* 10011001 11111111 00000000 00000000 11101000 01000010 00111100 00001111 */
		{0x99, 0xff, 0x00, 0x00, 0xe8, 0x42, 0x3c, 0x0f},
/* 00001111 10011001 11111111 00000000 00000000 11101000 01000010 00111100 */
		{0x0f, 0x99, 0xff, 0x00, 0x00, 0xe8, 0x42, 0x3c},
/* 00111100 00001111 10011001 11111111 00000000 00000000 11101000 01000010 */
		{0x3c, 0x0f, 0x99, 0xff, 0x00, 0x00, 0xe8, 0x42},
/* 01000010 00111100 00001111 10011001 11111111 00000000 00000000 11101000 */
		{0x42, 0x3c, 0x0f, 0x99, 0xff, 0x00, 0x00, 0xe8},
/* 11101000 01000010 00111100 00001111 10011001 11111111 00000000 00000000 */
		{0xe8, 0x42, 0x3c, 0x0f, 0x99, 0xff, 0x00, 0x00},
/* 00000000 11101000 01000010 00111100 00001111 10011001 11111111 00000000 */
		{0x00, 0xe8, 0x42, 0x3c, 0x0f, 0x99, 0xff, 0x00},
/* 00000000 00000000 11101000 01000010 00111100 00001111 10011001 11111111 */
		{0x00, 0x00, 0xe8, 0x42, 0x3c, 0x0f, 0x99, 0xff},
	};

	static uint8_t __mask[] = {
/* 10000000 01000000 00100000 00010000 00001000 00000100 00000010 00000001 */
		0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01
	};

	uint8_t __and[8], __ecc = 0;

	for (uint32_t i = 0; i < sizeof(__matrix) / sizeof(*__matrix); i++) {
		int __popcount = 0;

		for (uint32_t __byte = 0; __byte < 8; __byte++) {
			/* compute the AND of the data and ECC matrix */
			__and[__byte] = __data[__byte] & __matrix[i][__byte];
			/* count the number of '1' bits in the result (parity) */
			__popcount += popcount(__and[__byte]);
		}

		/* if the result is odd parity, turn on corresponding ECC bit */
		if ((__popcount) & 1)	/* odd parity? */
			__ecc |= __mask[i];

#ifdef DEBUG
		printf("\nmatrix[");
		print_binary(NULL, __matrix[i], 8);
		printf("]\n  data[");
		print_binary(NULL, __data, 8);
		printf("]\n   and[");
		print_binary(NULL, __and, 8);
		printf("]\n");
		printf("popcount: %d ecc %2.2x\n", __popcount, __ecc);
#endif
	}

	/* the ECC data is inverted such that */
	/* 0xFFFFFFFFffffffff => 0xFF for erased NOR flash */
	return __ecc ^ 0xFF;
}


/*
 * "aaaaaaaa [wwwwwwww_xxxxxxxx e yyyyyyyy_zzzzzzzz e] [........ ........]"
 */
static void __ecc_dump(FILE * __out, uint32_t __addr,
		       void *__restrict __buf, size_t __buf_sz, bool invert)
{
	if (__buf_sz <= 0 || __buf == NULL)
		return;
	if (__out == NULL)
		__out = stdout;

	size_t hex_size = (16 + 2) * 2;	// 16 bytes per doubleword plus ECC byte
	size_t ascii_size = (8 + 1) * 2;	// 8 bytes per doublewod plus ECC byte

	uint8_t hex[hex_size + 1], ascii[hex_size + 1];
	memset(hex, '.', hex_size), memset(ascii, '.', ascii_size);

	void print_line(void) {
		const char *ansi_red = "\033[1;1m\033[1;31m";
		const char *ansi_norm = "\033[0m";

		uint8_t c1 = sfc_ecc(ascii + 0), e1 = ascii[8];
		if (invert == true)
			c1 = ~c1;
		uint8_t c2 = sfc_ecc(ascii + 9), e2 = ascii[17];
		if (invert == true)
			c2 = ~c2;

		for (size_t i = 0; i < ascii_size; i++)
			ascii[i] = isprint(ascii[i]) ? ascii[i] : '.';

		fprintf(__out, "%08x ", __addr);

		fprintf(__out,
			"[%.8s_%.8s %.*s%.2s%.*s %.8s_%.8s %.*s%.2s%.*s] ",
			hex + 0, hex + 8, (c1 != e1) ? (uint32_t)strlen(ansi_red) : 0,
			ansi_red, hex + 16, (c1 != e1) ? (uint32_t)strlen(ansi_norm) : 0,
			ansi_norm, hex + 18, hex + 26,
			(c2 != e2) ? (uint32_t)strlen(ansi_red) : 0, ansi_red, hex + 34,
			(c2 != e2) ? (uint32_t)strlen(ansi_norm) : 0, ansi_norm);

		fprintf(__out, "[%.8s %.8s]\n", ascii + 0, ascii + 9);
	}

	size_t s = __addr % 16, i = 0;
	__addr &= ~0xF;

	for (i = s; i < __buf_sz + s; i++) {
		unsigned char c = ((unsigned char *)__buf)[i - s];

		hex[((i << 1) + 0) % hex_size] = "0123456789abcdef"[c >> 4];
		hex[((i << 1) + 1) % hex_size] = "0123456789abcdef"[c & 0xf];

		ascii[i % ascii_size] = c;

		if (i == 0)
			continue;

		if ((i + 1) % ascii_size == 0) {
			print_line();
			memset(hex, '.', hex_size), memset(ascii, '.',
							   ascii_size);

			__addr += ascii_size;
		}
	}

	if (i % ascii_size)
		print_line();

	return;
}



/* ======================================== */
static uint64_t ecc_matrix[] = {
        //0000000000000000111010000100001000111100000011111001100111111111
        0x0000e8423c0f99ff,
        //0000000011101000010000100011110000001111100110011111111100000000
        0x00e8423c0f99ff00,
        //1110100001000010001111000000111110011001111111110000000000000000
        0xe8423c0f99ff0000,
        //0100001000111100000011111001100111111111000000000000000011101000
        0x423c0f99ff0000e8,
        //0011110000001111100110011111111100000000000000001110100001000010
        0x3c0f99ff0000e842,
        //0000111110011001111111110000000000000000111010000100001000111100
        0x0f99ff0000e8423c,
        //1001100111111111000000000000000011101000010000100011110000001111
        0x99ff0000e8423c0f,
        //1111111100000000000000001110100001000010001111000000111110011001
        0xff0000e8423c0f99
};

static uint8_t syndrome_matrix[] = {
        GD, E7, E6, UE, E5, UE, UE, 47, E4, UE, UE, 37, UE, 35, 39, UE,
        E3, UE, UE, 48, UE, 30, 29, UE, UE, 57, 27, UE, 31, UE, UE, UE,
        E2, UE, UE, 17, UE, 18, 40, UE, UE, 58, 22, UE, 21, UE, UE, UE,
        UE, 16, 49, UE, 19, UE, UE, UE, 23, UE, UE, UE, UE, 20, UE, UE,
        E1, UE, UE, 51, UE, 46,  9, UE, UE, 34, 10, UE, 32, UE, UE, 36,
        UE, 62, 50, UE, 14, UE, UE, UE, 13, UE, UE, UE, UE, UE, UE, UE,
        UE, 61,  8, UE, 41, UE, UE, UE, 11, UE, UE, UE, UE, UE, UE, UE,
        15, UE, UE, UE, UE, UE, UE, UE, UE, UE, 12, UE, UE, UE, UE, UE,
        E0, UE, UE, 55, UE, 45, 43, UE, UE, 56, 38, UE,  1, UE, UE, UE,
        UE, 25, 26, UE,  2, UE, UE, UE, 24, UE, UE, UE, UE, UE, 28, UE,
        UE, 59, 54, UE, 42, UE, UE, 44,  6, UE, UE, UE, UE, UE, UE, UE,
        5, UE, UE, UE, UE, UE, UE, UE, UE, UE, UE, UE, UE, UE, UE, UE,
        UE, 63, 53, UE,  0, UE, UE, UE, 33, UE, UE, UE, UE, UE, UE, UE,
        3, UE, UE, 52, UE, UE, UE, UE, UE, UE, UE, UE, UE, UE, UE, UE,
        7, UE, UE, UE, UE, UE, UE, UE, UE, 60, UE, UE, UE, UE, UE, UE,
        UE, UE, UE, UE,  4, UE, UE, UE, UE, UE, UE, UE, UE, UE, UE, UE,
};


static uint8_t generate_ecc(uint64_t i_data)
{
        uint8_t result = 0;

        for (int i = 0; i < 8; i++)
        {
                result |= __builtin_parityll(ecc_matrix[i] & i_data) << i;
        }
        return result;
}
static uint8_t verify_ecc(uint64_t i_data, uint8_t i_ecc)
{
       return syndrome_matrix[generate_ecc(i_data) ^ i_ecc ];
}
static uint8_t correct_ecc(uint64_t *io_data, uint8_t *io_ecc)
{
        uint8_t bad_bit = verify_ecc(*io_data, *io_ecc);

        if ((bad_bit != GD) && (bad_bit != UE))  // Good is done, UE is hopeless.
        {
                // Determine if the ECC or data part is bad, do bit flip.
                if (bad_bit >= E7)
                {
                        *io_ecc ^= (1 << (bad_bit - E7));
                }
                else
                {
                        *io_data ^=(1ull << (63 - bad_bit));
                }
        }
        return bad_bit;
}

static void inject_ecc(const uint8_t* i_src, size_t i_srcSz,
               uint8_t* o_dst, bool invert)
{
        assert(0 == (i_srcSz % sizeof(uint64_t)));

        for(size_t i = 0, o = 0;
            i < i_srcSz;
            i += sizeof(uint64_t), o += sizeof(uint64_t) + sizeof(uint8_t))
        {
                // Read data word, copy to destination.
                uint64_t data = *(const uint64_t*)(&i_src[i]);

                *(uint64_t*)(&o_dst[o]) = data;
                data = be64toh(data);

                // Calculate ECC, copy to destination.
                uint8_t ecc = generate_ecc(data);
                o_dst[o + sizeof(uint64_t)] = invert ? ~ecc : ecc;
        }
}
static ecc_status_t remove_ecc(uint8_t* io_src, size_t i_srcSz,
                        uint8_t* o_dst, size_t i_dstSz,
                        bool invert)
{
        assert(0 == (i_dstSz % sizeof(uint64_t)));

        ecc_status_t rc = CLEAN;

        for(size_t i = 0, o = 0;
            i < i_srcSz;
            i += sizeof(uint64_t) + sizeof(uint8_t), o += sizeof(uint64_t))
        {
                // Read data and ECC parts.
                uint64_t data = *(uint64_t*)(&io_src[i]);
                data = be64toh(data);

                uint8_t ecc = io_src[i + sizeof(uint64_t)];
                if(invert)
                {
                        ecc = ~ecc;
                }
                // Calculate failing bit and fix data.
                uint8_t bad_bit = correct_ecc(&data, &ecc);

                // Return data to big endian.
                data = htobe64(data);

                // Perform correction and status update.
                if (bad_bit == UE)
                {
                        rc = UNCORRECTABLE;
                }
                else if (bad_bit != GD)
                {
                        if (rc != UNCORRECTABLE)
                        {
                                rc = CORRECTED;
                        }
                        *(uint64_t*)(&io_src[i]) = data;
                        io_src[i + sizeof(uint64_t)] = invert ? ~ecc : ecc;
                }

                // Copy fixed data to destination buffer.
                *(uint64_t*)(&o_dst[o]) = data;
        }
        return rc;
}
/* ========================================= */

static ssize_t __ecc_inject(void *__restrict __dst, size_t __dst_sz,
                       const void *__restrict __src, size_t __src_sz,
                       bool invert)
{
        int __size = sizeof(uint64_t);

        errno = 0;
        if (__src_sz & (__size - 1)) {
                errno = EINVAL;
                return -1;
        }
        if (__dst_sz < (__src_sz + (__src_sz / __size))) {
                errno = ENOBUFS;
                return -1;
        }

        ssize_t rc=0;
        ssize_t c_sz = __src_sz;
        for (; c_sz; c_sz -= __size) {
        rc += __size + 1;
	      }

        inject_ecc(__src, __src_sz, __dst, invert);
        return  (rc);
}

static ssize_t __ecc_remove(void *__restrict __dst, size_t __dst_sz,
                       const void *__restrict __src, size_t __src_sz,
                       bool invert)
{
        int __size = sizeof(uint64_t);

        errno = 0;
        if ((__src_sz % (__size + 1)) != 0) {
                errno = EINVAL;
                return -1;
        }
        if (__dst_sz < (__src_sz - (__src_sz / __size))) {
                errno = ENOBUFS;
                return -1;
        }


        int target_size = ((__src_sz / (sizeof(uint64_t) + 1))*sizeof(uint64_t));
        if( remove_ecc((uint8_t*)__src, __src_sz, __dst, __dst_sz, invert) != CLEAN)
        {
                target_size = 0;
        }
        return target_size;
}

void sfc_ecc_dump(FILE * __out, uint32_t __addr,
		  void *__restrict __buf, size_t __buf_sz)
{
	return __ecc_dump(__out, __addr, __buf, __buf_sz, false);
}

/* ========================================= */
ssize_t sfc_ecc_inject(void *__restrict __dst, size_t __dst_sz,
                       const void *__restrict __src, size_t __src_sz)
{
        return __ecc_inject(__dst, __dst_sz, __src, __src_sz, true);
}
ssize_t sfc_ecc_remove(void *__restrict __dst, size_t __dst_sz,
                       const void *__restrict __src, size_t __src_sz)
{
        return __ecc_remove(__dst, __dst_sz, __src, __src_sz, true);

}
ssize_t p8_ecc_remove_size (void *__restrict __dst, size_t __dst_sz,
		      void *__restrict __src, size_t __src_sz __unused__)
{
        return __ecc_remove(__dst, __dst_sz, __src, __src_sz, false);
}

/* ========================================= */
ssize_t p8_ecc_inject(void *__restrict __dst, size_t __dst_sz,
		      const void *__restrict __src, size_t __src_sz)
{
        return __ecc_inject(__dst, __dst_sz, __src, __src_sz, false);
}

ecc_status_t p8_ecc_remove (void *__restrict __dst, size_t __dst_sz,
		      void *__restrict __src, size_t __src_sz __unused__)
{
        return remove_ecc(__src, __src_sz, __dst, __dst_sz, false);
}

void p8_ecc_dump(FILE * __out, uint32_t __addr,
                 void *__restrict __buf, size_t __buf_sz)
{
        return __ecc_dump(__out, __addr, __buf, __buf_sz, true);
}


/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/test/checksum.c $                                        */
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

#include <stdio.h>
#include <stdlib.h>

#include <clib/checksum.h>
#include <clib/exception.h>

int main(void)
{
	unsigned int i;
	uint8_t *buf = malloc(32);
	uint8_t align[8] =       { 0x01, 0x02, 0x04, 0x08,
				   0x10, 0x20, 0x40, 0x80 };
	uint8_t unaligned1[9] =  { 0x11, 0x11, 0x11, 0x11,
				   0x22, 0x22, 0x22, 0x22, 0x33 };
	uint8_t unaligned2[10] = { 0x11, 0x11, 0x11, 0x11,
				   0x22, 0x22, 0x22, 0x22, 0x33, 0x33 };
	uint8_t unaligned3[11] = { 0x11, 0x11, 0x11, 0x11,
				   0x22, 0x22, 0x22, 0x22, 0x33, 0x33, 0x33 };

	memcpy(buf, align, sizeof(align));
	uint32_t csum = memcpy_checksum(NULL, buf, sizeof(align));
	if (csum != 0x11224488) {
		printf("fail %d a:%08x e:%08x\n", __LINE__, csum, 0x11224488);
		return 1;
	}

	memcpy(buf, unaligned1, sizeof(unaligned1));
	csum = memcpy_checksum(NULL, buf, sizeof(unaligned1));
	if (csum != 0x00333333) {
		printf("fail %d a:%08x e:%08x\n", __LINE__, csum, 0x00333333);
		return 1;
	}

	memcpy(buf, unaligned2, sizeof(unaligned2));
	csum = memcpy_checksum(NULL, buf, sizeof(unaligned2));
	if (csum != 0x00003333) {
		printf("fail %d a:%08x e:%08x\n", __LINE__, csum, 0x00003333);
		return 1;
	}

	memcpy(buf, unaligned3, sizeof(unaligned3));
	csum = memcpy_checksum(NULL, buf, sizeof(unaligned3));
	if (csum != 0x00000033) {
		printf("fail %d a:%08x e:%08x\n", __LINE__, csum, 0x00000033);
		return 1;
	}

	char dst[16];
	csum = memcpy_checksum(dst, align, sizeof(align));
	if (csum != 0x11224488) {
		printf("fail %d a:%08x e:%08x\n", __LINE__, csum, 0x11224488);
		return 1;
	}

	csum = memcpy_checksum(dst, unaligned1, sizeof(unaligned1));
	for (i = 0; i < sizeof(unaligned1); i++) {
		if (dst[i] != unaligned1[i])
			printf("fail %d byte %i a:%08x e:%08x\n",
					__LINE__, i, dst[i], unaligned1[i]);
	}
	if (csum != 0x00333333) {
		printf("fail %d a:%08x e:%08x\n", __LINE__, csum, 0x00333333);
		return 1;
	}

	csum = memcpy_checksum(dst, unaligned2, sizeof(unaligned2));
	for (i = 0; i < sizeof(unaligned2); i++) {
		if (dst[i] != unaligned2[i])
			printf("fail %d byte %i a:%08x e:%08x\n",
					__LINE__, i, dst[i], unaligned2[i]);
	}
	if (csum != 0x00003333) {
		printf("fail %d a:%08x e:%08x\n", __LINE__, csum, 0x00003333);
		return 1;
	}

	csum = memcpy_checksum(dst, unaligned3, sizeof(unaligned3));
	for (i = 0; i < sizeof(unaligned3); i++) {
		if (dst[i] != unaligned3[i])
			printf("fail %d byte %i a:%08x e:%08x\n",
					__LINE__, i, dst[i], unaligned3[i]);
	}
	if (csum != 0x00000033) {
		printf("fail %d a:%08x e:%08x\n", __LINE__, csum, 0x00000033);
		return 1;
	}

#if 0
	exception_t ex;
	try {
		(void)memcpy_checksum(NULL, &align[1], sizeof(align) - 1);
	} else (ex) {
		printf("Hey, what's this in the weeds? It's a baby, awesome!\n");
	} end_try;

#endif

	return 0;
}

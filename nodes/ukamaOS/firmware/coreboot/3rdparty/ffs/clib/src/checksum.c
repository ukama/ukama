/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/src/checksum.c $                                         */
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
#include <stdint.h>

#include "assert.h"
#include "checksum.h"

uint32_t memcpy_checksum(void *__restrict __dst, const void *__restrict __src,
			 size_t __n)
{
	uint8_t sum[4] = { 0, };

	/* assert(((uintptr_t)__src & 3) == 0); */

	size_t i;

	if (__dst == NULL)
		for (i = 0; i < __n; i++)
			sum[i & 3] ^= ((uint8_t *) __src)[i];
	else
		for (i = 0; i < __n; i++)
			sum[i & 3] ^= ((uint8_t *) __src)[i],
			    ((uint8_t *) __dst)[i] = ((uint8_t *) __src)[i];

	return (sum[0] << 24) | (sum[1] << 16) | (sum[2] << 8) | sum[3];
}

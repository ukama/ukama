/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/checksum.h $                                             */
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

/*! @file checksum.h
 *  @brief Simple XOR checksum
 *  @author Mike Kobler <mkobler@us.ibm.com>
 *  @date 2007-2012
 */

#ifndef __CHECKSUM_H__
#define __CHECKSUM_H__

#include <stdint.h>

/*!
 * @brief Copy bytes from the source reference to the destination reference
 *        while computing a 32-bit checksum
 * @param __dst [in] Destination reference
 * @param __src [in] Source reference (must be 4 byte aligned)
 * @param __n [in] Number of bytes to copy / compute
 * @return 32-bit Checksum value
 */
extern uint32_t memcpy_checksum(void *__restrict __dst,
				const void *__restrict __src, size_t __n)
/*! @cond */
__THROW __nonnull((2)) /*! @endcond */ ;

#endif				/* __CHECKSUM_H__ */

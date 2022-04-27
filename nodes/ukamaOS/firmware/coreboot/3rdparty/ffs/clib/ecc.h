/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/ecc.h $                                                  */
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

/*!
 * @file ecc.h
 * @brief FSP-2 & P8 ECC functions
 * @details These functions are used to insert and remove FSP-2 & P8 ECC
 *          bytes. 8-bits of ECC is inserted every 8 bytes of data such as:
 *          XXxxXXxxXXxxXXxxYY (where XXxx is 4 nibbles of data and YY is
 *          2 nibbles of ECC)
 * @author Shaun Wetzstein <shaun@us.ibm.com>
 * @date 2011
 */

#ifndef __ECC_H__
#define __ECC_H__

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
/** Status for the ECC removal function. */
enum ecc_status
    {
      CLEAN=0,          //< No ECC Error was detected.
      CORRECTED=1,      //< ECC error detected and corrected.
      UNCORRECTABLE=2   //< ECC error detected and uncorrectable.
    };
typedef enum  ecc_status ecc_status_t;

enum ecc_bitfields
    {
        GD = 0xff,      //< Good, ECC matches.
        UE = 0xfe,      //< Uncorrectable.
        E0 = 71,        //< Error in ECC bit 0
        E1 = 70,        //< Error in ECC bit 1
        E2 = 69,        //< Error in ECC bit 2
        E3 = 68,        //< Error in ECC bit 3
        E4 = 67,        //< Error in ECC bit 4
        E5 = 66,        //< Error in ECC bit 5
        E6 = 65,        //< Error in ECC bit 6
        E7 = 64         //< Error in ECC bit 7
    };

/*!
 * @brief Compute the 8-bit ECC (SFC) value given an array of 8
 *        unsigned char data values
 * @param data [in] Input data
 * @return 8-bit SFC ECC value
 */
	extern uint8_t sfc_ecc(uint8_t __data[8]);

/*!
 * @brief Copy bytes from the source buffer to the destination buffer while
 *        computing and injecting an 8-bit SFC ECC value for every 8-bytes
 *        of source buffer read
 * @param __dst [in] Destination buffer
 * @param __dst_sz [in] Destination buffer size (in bytes) which must be large
 *        enough to store both the data and the ECC bytes
 * @param __src [in] Source buffer
 * @param __src_sz [in] Source buffer size (in bytes) which must be a multiple
 *        of 8 bytes
 * @return -1 if an error occurs, number of bytes copied (including ECC bytes)
 *         otherwise.
 *         EINVAL if __src_sz is 0 or not a multiple of 8 bytes
 *         ENOBUFS if __dst_sz is not large enough to store the ECC bytes
 *
 */
	extern ssize_t sfc_ecc_inject(void *__restrict __dst, size_t __dst_sz,
				      const void *__restrict __src,
				      size_t __src_sz)
/*! @cond */
	 __nonnull((1, 3)) /*! @endcond */ ;

/*!
 * @brief Copy bytes from the source buffer to the destination buffer while
 *        computing and removing an 8-bit SFC ECC value for every 9-bytes
 *        of source buffer read
 * @param __dst [in] Destination buffer
 * @param __dst_sz [in] Destination buffer size (in bytes) which must be large
 *        enough to store the data (after ECC removal)
 * @param __src [in] Source buffer
 * @param __src_sz [in] Source buffer size (in bytes) which must be a multiple
 *        9 bytes
 * @return -1 if an error occurs, number of bytes copied (excluding ECC bytes)
 *         otherwise.  if an ECC mismatch error occurs the function returns
 *         immediately and the return code indicates the number of bytes processed
 *         prior to the ECC error.
 *         EINVAL if __src_sz is 0 or not a multiple of 9 bytes
 *         ENOBUFS if __dst_sz is not large enough to store the ECC bytes
 */
	extern ssize_t sfc_ecc_remove(void *__restrict __dst, size_t __dst_sz,
				      const void *__restrict __src,
				      size_t __src_sz)
/*! @cond */
	 __nonnull((1, 3)) /*! @endcond */ ;

/*!
 * @brief Hexdump the contents of a memory buffer to an output stream.
 *        This is a buck-standard hexdump except it issolates the ECC value
 *        in a separate column for easy debug.
 * @param __out [in] Output stream
 * @param __addr [in] Starting put to display (in bytes)
 * @param __buf [in] Data buffer
 * @param __buf_sz [in] Data buffer size (in bytes)
 * @return -1 if an error occurs, number of bytes copied (excluding ECC bytes)
 *         otherwise.  if an ECC mismatch error occurs the function highlights the
 *         the corrupted data with red ANSI.
 */
	extern void sfc_ecc_dump(FILE * __out, uint32_t __addr,
				 void *__restrict __buf, size_t __buf_sz)
/*! @cond */
	 __nonnull((1, 3)) /*! @endcond */ ;

/*!
 * @brief Copy bytes from the source buffer to the destination buffer while
 *        computing and injecting an 8-bit P8 ECC value for every 8-bytes
 *        of source buffer read
 * @param __dst [in] Destination buffer
 * @param __dst_sz [in] Destination buffer size (in bytes) which must be large
 *        enough to store both the data and the ECC bytes
 * @param __src [in] Source buffer
 * @param __src_sz [in] Source buffer size (in bytes) which must be a multiple
 *        of 8 bytes
 * @return -1 if an error occurs, number of bytes copied (including ECC bytes)
 *         otherwise.
 *         EINVAL if __src_sz is 0 or not a multiple of 8 bytes
 *         ENOBUFS if __dst_sz is not large enough to store the ECC bytes
 *
 */
	extern ssize_t p8_ecc_inject(void *__restrict __dst, size_t __dst_sz,
				     const void *__restrict __src,
				     size_t __src_sz)
/*! @cond */
	 __nonnull((1, 3)) /*! @endcond */ ;

/*!
 * @brief Copy bytes from the source buffer to the destination buffer while
 *        computing and removing an 8-bit P8 ECC value for every 9-bytes
 *        of source buffer read
 * @param __dst [in] Destination buffer
 * @param __dst_sz [in] Destination buffer size (in bytes) which must be large
 *        enough to store the data (after ECC removal)
 * @param __src [in] Source buffer
 * @param __src_sz [in] Source buffer size (in bytes) which must be a multiple
 *        9 bytes
 * @return 0 - CLEAN for success
 *         1 - CORRECTED error - _src [in] buffer changed
 *         2 - UNCORRECTABLE error.
 */
  extern ecc_status_t p8_ecc_remove(void *__restrict __dst, size_t __dst_sz,
				     void *__restrict __src,
				     size_t __src_sz)

/*! @cond */
	 __nonnull((1, 3)) /*! @endcond */ ;

/*!
 * @brief Copy bytes from the source buffer to the destination buffer while
 *        computing and removing an 8-bit P8 ECC value for every 9-bytes
 *        of source buffer read
 * @param __dst [in] Destination buffer
 * @param __dst_sz [in] Destination buffer size (in bytes) which must be large
 *        enough to store the data (after ECC removal)
 * @param __src [in] Source buffer
 * @param __src_sz [in] Source buffer size (in bytes) which must be a multiple
 *        9 bytes
 * @return -1 if an error occurs, number of bytes copied (excluding ECC bytes)
 *         otherwise.  if an ECC mismatch error occurs the function returns
 *         immediately and the return code indicates zero number of bytes
 *         processed
 *         EINVAL if __src_sz is 0 or not a multiple of 9 bytes
 *         ENOBUFS if __dst_sz is not large enough to store the ECC bytes
 */
  extern ssize_t p8_ecc_remove_size(void *__restrict __dst, size_t __dst_sz,
				     void *__restrict __src, size_t __src_sz)

/*! @cond */
	 __nonnull((1, 3)) /*! @endcond */ ;

/*!
 * @brief Hexdump the contents of a memory buffer to an output stream.
 *        This is a buck-standard hexdump except it issolates the P8 ECC
 *        value in a separate column for easy debug.
 * @param __out [in] Output stream
 * @param __addr [in] Starting put to display (in bytes)
 * @param __buf [in] Data buffer
 * @param __buf_sz [in] Data buffer size (in bytes)
 * @return -1 if an error occurs, number of bytes copied (excluding ECC bytes)
 *         otherwise.  if an ECC mismatch error occurs the function highlights the
 *         the corrupted data with red ANSI.
 */
	extern void p8_ecc_dump(FILE * __out, uint32_t __addr,
				void *__restrict __buf, size_t __buf_sz)
/*! @cond */
	 __nonnull((1, 3)) /*! @endcond */ ;

#ifdef __cplusplus
}
#endif
#endif				/* __ECC_H__ */

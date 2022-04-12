/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: ffs/test/test_libffs.h $                                      */
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
 *    File: test_libffs.c
 *  Author: Shekar Babu S <shekbabu@in.ibm.com>
 *   Descr: unit test tool for api's in libffs.so
 *    Date: 06/26/2012
 */

#include "libffs.h"

#define FFS_ERROR -1
#define PART_OFFSET 0x3F0000

typedef struct ffs_operations{
	const char  * nor_image;  //!< Path to nor image special/regular file
	const char  * part_entry; //!< Logical partition/entry name
	const char  * i_file;	  //!< Input file
	const char  * o_file;	  //!< Output file
	FILE	    * log;	  //!< Log file
	size_t        device_size;//!< Size of the nor flash
	off_t         part_off;	  //!< Offset of partition table
	size_t        blk_sz;	  //!< Block size
	size_t        entry_sz;	  //!< Partition entry size
	off_t         entry_off;  //!< Offset of partition entry
	uint32_t      user;	  //!< Index to user word in any entry
	uint32_t      value;      //!< User word at index in entry
	ffs_type_t    type;	  //!< Partition type
} ffs_ops_t;

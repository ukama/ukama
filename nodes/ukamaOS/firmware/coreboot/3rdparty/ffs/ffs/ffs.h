/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: ffs/ffs.h $                                                   */
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
#ifndef __FFS_H__
#define __FFS_H__

/* Pull in the correct header depending on what is being built */
#if defined(__KERNEL__)
#include <linux/types.h>
#else
#include <stdint.h>
#endif

/* The version of this partition implementation */
#define FFS_VERSION_1	1

/* Magic number for the partition header (ASCII 'PART') */
#define FFS_MAGIC	0x50415254

/* The maximum length of the partition name */
#define PART_NAME_MAX   15

/*
 * Sizes of the data structures
 */
#define FFS_HDR_SIZE   sizeof(struct ffs_hdr)
#define FFS_ENTRY_SIZE sizeof(struct ffs_entry)

/*
* Size of FFS_HDR w/o entry struct in bytes
*/
#define FFS_HDR_SIZE_NO_ENTRY 0x48

/*
 * Sizes of the data structures w/o checksum
 */
#define FFS_HDR_SIZE_CSUM   (FFS_HDR_SIZE - sizeof(uint32_t))
#define FFS_ENTRY_SIZE_CSUM (FFS_ENTRY_SIZE - sizeof(uint32_t))

/* pid of logical partitions/containers */
#define FFS_PID_TOPLEVEL   0xFFFFFFFF

/*
 * Type of image contained w/in partition
 */
enum type {
	FFS_TYPE_DATA      = 1,
	FFS_TYPE_LOGICAL   = 2,
	FFS_TYPE_PARTITION = 3,
};

/*
 * Flag bit definitions
 */
#define FFS_FLAGS_PROTECTED	0x0001
#define FFS_FLAGS_U_BOOT_ENV	0x0002

/*
 * Number of user data words
 */
#define FFS_USER_WORDS 16

/*
 * Define layout of user.data in struct ffs_entry
 */
enum user_data {
	USER_DATA_VOL  = 0,
	USER_DATA_SIZE = 1,
	USER_DATA_CRC  = 2,
};

/**
 * struct ffs_entry - Partition entry
 *
 * @name:	Opaque null terminated string
 * @base:	Starting offset of partition in flash (in hdr.block_size)
 * @size:	Partition size (in hdr.block_size)
 * @pid:	Parent partition entry (FFS_PID_TOPLEVEL for toplevel)
 * @id:		Partition entry ID [1..65536]
 * @type:	Describe type of partition
 * @flags:	Partition attributes (optional)
 * @actual:	Actual partition size (in bytes)
 * @resvd:	Reserved words for future use
 * @user:	User data (optional)
 * @checksum:	Partition entry checksum (includes all above)
 */
struct ffs_entry {
	char     name[PART_NAME_MAX + 1];
	uint32_t base;
	uint32_t size;
	uint32_t pid;
	uint32_t id;
	uint32_t type;
	uint32_t flags;
	uint32_t actual;
	uint32_t resvd[4];
	struct {
		uint32_t data[FFS_USER_WORDS];
	} user;
	uint32_t checksum;
} __attribute__ ((packed));

/**
 * struct ffs_hdr - FSP Flash Structure header
 *
 * @magic:		Eye catcher/corruption detector
 * @version:		Version of the structure
 * @size:		Size of partition table (in block_size)
 * @entry_size:		Size of struct ffs_entry element (in bytes)
 * @entry_count:	Number of struct ffs_entry elements in @entries array
 * @block_size:		Size of block on device (in bytes)
 * @block_count:	Number of blocks on device
 * @resvd:		Reserved words for future use
 * @checksum:		Header checksum
 * @entries:		Pointer to array of partition entries
 */
struct ffs_hdr {
	uint32_t         magic;
	uint32_t         version;
	uint32_t         size;
	uint32_t         entry_size;
	uint32_t         entry_count;
	uint32_t         block_size;
	uint32_t         block_count;
	uint32_t         resvd[4];
	uint32_t         checksum;
	struct ffs_entry entries[];
} __attribute__ ((packed));

#endif /* __FFS_H__ */

/* Copyright 2019 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Secure non-volatile storage data structure definitions
 */

#ifndef VBOOT_REFERENCE_2SECDATA_STRUCT_H_
#define VBOOT_REFERENCE_2SECDATA_STRUCT_H_

#include "2crc8.h"
#include "2sysincludes.h"

/*****************************************************************************/
/* Firmware secure storage space */

#define VB2_SECDATA_FIRMWARE_VERSION 2

struct vb2_secdata_firmware {
	/* Struct version, for backwards compatibility */
	uint8_t struct_version;

	/* Flags; see vb2_secdata_firmware_flags */
	uint8_t flags;

	/* Firmware versions */
	uint32_t fw_versions;

	/* Reserved for future expansion */
	uint8_t reserved[3];

	/* CRC; must be last field in struct */
	uint8_t crc8;
} __attribute__((packed));

/*****************************************************************************/
/* Kernel secure storage space */

/* Kernel space - KERNEL_NV_INDEX, locked with physical presence. */
#define VB2_SECDATA_KERNEL_VERSION 2
#define VB2_SECDATA_KERNEL_UID 0x4752574c  /* 'LWRG' */

struct vb2_secdata_kernel {
	/* Struct version, for backwards compatibility */
	uint8_t struct_version;

	/* Unique ID to detect space redefinition */
	uint32_t uid;

	/* Kernel versions */
	uint32_t kernel_versions;

	/* Reserved for future expansion */
	uint8_t reserved[3];

	/* CRC; must be last field in struct */
	uint8_t crc8;
} __attribute__((packed));

/*****************************************************************************/
/* Firmware management parameters (FWMP) space */

#define VB2_SECDATA_FWMP_VERSION 0x10  /* 1.0 */
#define VB2_SECDATA_FWMP_HASH_SIZE 32  /* enough for SHA-256 */

struct vb2_secdata_fwmp {
	/* CRC-8 of fields following struct_size */
	uint8_t crc8;

	/* Structure size in bytes */
	uint8_t struct_size;

	/* Structure version (4 bits major, 4 bits minor) */
	uint8_t struct_version;

	/* Reserved; ignored by current reader */
	uint8_t reserved0;

	/* Flags; see enum vb2_secdata_fwmp_flags */
	uint32_t flags;

	/* Hash of developer kernel key */
	uint8_t dev_key_hash[VB2_SECDATA_FWMP_HASH_SIZE];
};

/**
 * Generate CRC for FWMP secure storage space.
 *
 * Calculate CRC hash from struct_version onward.  In valid FWMP data, this CRC
 * value should match the crc8 field.
 *
 * @param sec		Pointer to FWMP struct
 * @return 32-bit CRC hash of FWMP data
 */
static inline uint32_t vb2_secdata_fwmp_crc(struct vb2_secdata_fwmp *sec)
{
	int version_offset = offsetof(struct vb2_secdata_fwmp, struct_version);
	return vb2_crc8((void *)sec + version_offset,
			sec->struct_size - version_offset);
}

#endif  /* VBOOT_REFERENCE_2SECDATA_STRUCT_H_ */

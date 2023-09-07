/* Copyright (c) 2014 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Vboot data structures.
 *
 * Note: Many of the structs have pairs of 32-bit fields and reserved fields.
 * This is to be backwards-compatible with older verified boot data which used
 * 64-bit fields (when we thought that hey, UEFI is 64-bit so all our fields
 * should be too).
 *
 * Offsets should be padded to 32-bit boundaries, since some architectures
 * have trouble with accessing unaligned integers.
 */

#ifndef VBOOT_REFERENCE_2STRUCT_H_
#define VBOOT_REFERENCE_2STRUCT_H_

#include "2api.h"
#include "2constants.h"
#include "2crypto.h"
#include "2sysincludes.h"

/* Flags for vb2_shared_data.flags */
enum vb2_shared_data_flags {
	/* User has explicitly and physically requested recovery */
	VB2_SD_FLAG_MANUAL_RECOVERY = (1 << 0),

	/* Developer mode is enabled */
	VB2_SD_FLAG_DEV_MODE_ENABLED = (1 << 1),

	/*
	 * TODO: might be nice to add flags for why dev mode is enabled - via
	 * gbb, virtual dev switch, or forced on for testing.
	 */

	/* Kernel keyblock was verified by signature (not just hash) */
	VB2_SD_FLAG_KERNEL_SIGNED = (1 << 2),

	/* Software sync needs to update EC-RO, EC-RW, or PD-RW respectively */
	VB2_SD_FLAG_ECSYNC_EC_RO = (1 << 3),
	VB2_SD_FLAG_ECSYNC_EC_RW = (1 << 4),
	VB2_SD_FLAG_ECSYNC_PD_RW = (1 << 5),

	/* Software sync says EC / PD running RW */
	VB2_SD_FLAG_ECSYNC_EC_IN_RW = (1 << 6),
	VB2_SD_FLAG_ECSYNC_PD_IN_RW = (1 << 7),

	/* Display is available on this boot */
	VB2_SD_FLAG_DISPLAY_AVAILABLE = (1 << 8),
};

/* Flags for vb2_shared_data.status */
enum vb2_shared_data_status {
	/* Reinitialized NV data due to invalid checksum */
	VB2_SD_STATUS_NV_REINIT = (1 << 0),

	/* NV data has been initialized */
	VB2_SD_STATUS_NV_INIT = (1 << 1),

	/* Secure data initialized */
	VB2_SD_STATUS_SECDATA_FIRMWARE_INIT = (1 << 2),

	/* Chose a firmware slot */
	VB2_SD_STATUS_CHOSE_SLOT = (1 << 3),

	/* Secure data kernel version space initialized */
	VB2_SD_STATUS_SECDATA_KERNEL_INIT = (1 << 4),

	/* FWMP secure data initialized */
	VB2_SD_STATUS_SECDATA_FWMP_INIT = (1 << 5),
};

/* "V2SD" = vb2_shared_data.magic */
#define VB2_SHARED_DATA_MAGIC 0x44533256

/* Current version of vb2_shared_data struct */
#define VB2_SHARED_DATA_VERSION_MAJOR 2
#define VB2_SHARED_DATA_VERSION_MINOR 0

#define VB2_CONTEXT_MAX_SIZE 192

/*
 * Data shared between vboot API calls.  Stored at the start of the work
 * buffer.
 */
struct vb2_shared_data {
	/* Magic number for struct (VB2_SHARED_DATA_MAGIC) */
	uint32_t magic;

	/* Version of this structure */
	uint16_t struct_version_major;
	uint16_t struct_version_minor;

	/* Public fields are stored in the context object */
	struct vb2_context ctx;

	/* Padding for adding future vb2_context fields */
	uint8_t padding[VB2_CONTEXT_MAX_SIZE - sizeof(struct vb2_context)];

	/* Work buffer length in bytes. */
	uint32_t workbuf_size;

	/*
	 * Amount of work buffer used so far.  Verified boot sub-calls use
	 * this to know where the unused work area starts.
	 */
	uint32_t workbuf_used;

	/* Flags; see enum vb2_shared_data_flags */
	uint32_t flags;

	/*
	 * Reason we are in recovery mode this boot (enum vb2_nv_recovery), or
	 * 0 if we aren't.
	 */
	uint32_t recovery_reason;

	/* Firmware slot used last boot (0=A, 1=B) */
	uint32_t last_fw_slot;

	/* Result of last boot (enum vb2_fw_result) */
	uint32_t last_fw_result;

	/* Firmware slot used this boot */
	uint32_t fw_slot;

	/*
	 * Version for this slot (top 16 bits = key, lower 16 bits = firmware).
	 *
	 * TODO: Make this a union to allow getting/setting those versions
	 * separately?
	 */
	uint32_t fw_version;

	/* Version from secdata_firmware (must be <= fw_version to boot). */
	uint32_t fw_version_secdata;

	/*
	 * Status flags for this boot; see enum vb2_shared_data_status.  Status
	 * is "what we've done"; flags above are "decisions we've made".
	 */
	uint32_t status;

	/* Offset from start of this struct to GBB header */
	uint32_t gbb_offset;

	/**********************************************************************
	 * Data from kernel verification stage.
	 *
	 * TODO: shouldn't be part of the main struct, since that needlessly
	 * uses more memory during firmware verification.
	 */

	/*
	 * Version for the current kernel (top 16 bits = key, lower 16 bits =
	 * kernel preamble).
	 *
	 * TODO: Make this a union to allow getting/setting those versions
	 * separately?
	 */
	uint32_t kernel_version;

	/* Version from secdata_kernel (must be <= kernel_version to boot) */
	uint32_t kernel_version_secdata;

	/**********************************************************************
	 * Temporary variables used during firmware verification.  These don't
	 * really need to persist through to the OS, but there's nowhere else
	 * we can put them.
	 */

	/* Offset of preamble from start of vblock */
	uint32_t vblock_preamble_offset;

	/*
	 * Offset and size of packed data key in work buffer.  Size is 0 if
	 * data key is not stored in the work buffer.
	 */
	uint32_t data_key_offset;
	uint32_t data_key_size;

	/*
	 * Offset and size of firmware preamble in work buffer.  Size is 0 if
	 * preamble is not stored in the work buffer.
	 */
	uint32_t preamble_offset;
	uint32_t preamble_size;

	/*
	 * Offset and size of hash context in work buffer.  Size is 0 if
	 * hash context is not stored in the work buffer.
	 */
	uint32_t hash_offset;
	uint32_t hash_size;

	/*
	 * Current tag we're hashing
	 *
	 * For new structs, this is the offset of the vb2_signature struct
	 * in the work buffer.
	 *
	 * TODO: rename to hash_sig_offset when vboot1 structs are deprecated.
	 */
	uint32_t hash_tag;

	/* Amount of data we still expect to hash */
	uint32_t hash_remaining_size;

	/**********************************************************************
	 * Temporary variables used during kernel verification.  These don't
	 * really need to persist through to the OS, but there's nowhere else
	 * we can put them.
	 *
	 * TODO: make a union with the firmware verification temp variables,
	 * or make both of them workbuf-allocated sub-structs, so that we can
	 * overlap them so kernel variables don't bloat firmware verification
	 * stage memory requirements.
	 */

	/*
	 * Vboot1 shared data header.  This data should eventually get folded
	 * directly into the kernel portion of this struct.
	 */
	struct VbSharedDataHeader *vbsd;

	/*
	 * Offset and size of packed kernel key in work buffer.  Size is 0 if
	 * subkey is not stored in the work buffer.  Note that kernel key may
	 * be inside the firmware preamble.
	 */
	uint32_t kernel_key_offset;
	uint32_t kernel_key_size;
} __attribute__((packed));

/****************************************************************************/

/* Signature at start of the GBB
 * Note that if you compile in the signature as is, you are likely to break any
 * tools that search for the signature. */
#define VB2_GBB_SIGNATURE "$GBB"
#define VB2_GBB_SIGNATURE_SIZE 4
#define VB2_GBB_XOR_CHARS "****"
/* TODO: can we write a macro to produce this at compile time? */
#define VB2_GBB_XOR_SIGNATURE { 0x0e, 0x6d, 0x68, 0x68 }

#define VB2_GBB_HWID_DIGEST_SIZE 32

/* VB2 GBB struct version */
#define VB2_GBB_MAJOR_VER 1
#define VB2_GBB_MINOR_VER 2
/* v1.2 - added fields for sha256 digest of the HWID */

struct vb2_gbb_header {
	/* Fields present in version 1.1 */
	uint8_t  signature[VB2_GBB_SIGNATURE_SIZE]; /* VB2_GBB_SIGNATURE */
	uint16_t major_version;   /* See VB2_GBB_MAJOR_VER */
	uint16_t minor_version;   /* See VB2_GBB_MINOR_VER */
	uint32_t header_size;     /* Size of GBB header in bytes */

	/* Flags (see enum vb2_gbb_flag in 2gbb_flags.h) */
	vb2_gbb_flags_t flags;

	/* Offsets (from start of header) and sizes (in bytes) of components */
	uint32_t hwid_offset;		/* HWID */
	uint32_t hwid_size;
	uint32_t rootkey_offset;	/* Root key */
	uint32_t rootkey_size;
	uint32_t bmpfv_offset;		/* BMP FV; deprecated in current FW */
	uint32_t bmpfv_size;
	uint32_t recovery_key_offset;	/* Recovery key */
	uint32_t recovery_key_size;

	/* Added in version 1.2 */
	uint8_t  hwid_digest[VB2_GBB_HWID_DIGEST_SIZE];	/* SHA-256 of HWID */

	/* Pad to match EXPECTED_VB2_GBB_HEADER_SIZE.  Initialize to 0. */
	uint8_t  pad[48];
} __attribute__((packed));

#define EXPECTED_VB2_GBB_HEADER_SIZE 128

/* VB2_GBB_FLAGS_OFFSET exposed in 2constants.h */
_Static_assert(VB2_GBB_FLAGS_OFFSET == offsetof(struct vb2_gbb_header, flags),
	       "VB2_GBB_FLAGS_OFFSET set incorrectly");

/****************************************************************************/

/*
 * Root key hash for Ryu devices only.  Contains the hash of the root key.
 * This will be embedded somewhere inside the RO part of the firmware, so that
 * it can verify the GBB contains only the official root key.
 */
#define RYU_ROOT_KEY_HASH_MAGIC "RtKyHash"
#define RYU_ROOT_KEY_HASH_MAGIC_INVCASE "rTkYhASH"
#define RYU_ROOT_KEY_HASH_MAGIC_SIZE 8

#define RYU_ROOT_KEY_HASH_VERSION_MAJOR 1
#define RYU_ROOT_KEY_HASH_VERSION_MINOR 0

struct vb2_ryu_root_key_hash {
	/* Magic number (RYU_ROOT_KEY_HASH_MAGIC) */
	uint8_t magic[RYU_ROOT_KEY_HASH_MAGIC_SIZE];

	/* Version of this struct */
	uint16_t header_version_major;
	uint16_t header_version_minor;

	/*
	 * Length of this struct, in bytes, including any variable length data
	 * which follows (there is none, yet).
	 */
	uint32_t struct_size;

	/*
	 * SHA-256 hash digest of the entire root key section from the GBB.  If
	 * all 0 bytes, all root keys will be treated as if matching.
	 */
	uint8_t root_key_hash_digest[32];
};

#define EXPECTED_VB2_RYU_ROOT_KEY_HASH_SIZE 48

/****************************************************************************/

/* Packed public key data */
struct vb2_packed_key {
	/* Offset of key data from start of this struct */
	uint32_t key_offset;
	uint32_t reserved0;

	/* Size of key data in bytes (NOT strength of key in bits) */
	uint32_t key_size;
	uint32_t reserved1;

	/* Signature algorithm used by the key (enum vb2_crypto_algorithm) */
	uint32_t algorithm;
	uint32_t reserved2;

	/* Key version */
	uint32_t key_version;
	uint32_t reserved3;

	/* TODO: when redoing this struct, add a text description of the key */
} __attribute__((packed));

#define EXPECTED_VB2_PACKED_KEY_SIZE 32

/****************************************************************************/

/* Signature data (a secure hash, possibly signed) */
struct vb2_signature {
	/* Offset of signature data from start of this struct */
	uint32_t sig_offset;
	uint32_t reserved0;

	/* Size of signature data in bytes */
	uint32_t sig_size;
	uint32_t reserved1;

	/* Size of the data block which was signed in bytes */
	uint32_t data_size;
	uint32_t reserved2;
} __attribute__((packed));

#define EXPECTED_VB2_SIGNATURE_SIZE 24

/****************************************************************************/

#define VB2_KEYBLOCK_MAGIC "CHROMEOS"
#define VB2_KEYBLOCK_MAGIC_SIZE 8

#define VB2_KEYBLOCK_VERSION_MAJOR 2
#define VB2_KEYBLOCK_VERSION_MINOR 1

/*
 * Keyblock flags.
 *
 * The following flags set where the key is valid.  Not used by firmware
 * verification; only kernel verification.
 */
#define VB2_KEYBLOCK_FLAG_DEVELOPER_0 0x1  /* Developer switch off */
#define VB2_KEYBLOCK_FLAG_DEVELOPER_1 0x2  /* Developer switch on */
#define VB2_KEYBLOCK_FLAG_RECOVERY_0 0x4  /* Not recovery mode */
#define VB2_KEYBLOCK_FLAG_RECOVERY_1 0x8  /* Recovery mode */

/*
 * Keyblock, containing the public key used to sign some other chunk of data.
 *
 * This should be followed by:
 *   1) The data_key key data, pointed to by data_key.key_offset.
 *   2) The checksum data for (vb2_keyblock + data_key data), pointed to
 *      by keyblock_checksum.sig_offset.
 *   3) The signature data for (vb2_keyblock + data_key data), pointed to
 *      by keyblock_signature.sig_offset.
 */
struct vb2_keyblock {
	/* Magic number */
	uint8_t magic[VB2_KEYBLOCK_MAGIC_SIZE];

	/* Version of this header format */
	uint32_t header_version_major;
	uint32_t header_version_minor;

	/*
	 * Length of this entire keyblock, including keys, signatures, and
	 * padding, in bytes
	 */
	uint32_t keyblock_size;
	uint32_t reserved0;

	/*
	 * Signature for this keyblock (header + data pointed to by data_key)
	 * For use with signed data keys
	 */
	struct vb2_signature keyblock_signature;

	/*
	 * SHA-512 hash for this keyblock (header + data pointed to by
	 * data_key) For use with unsigned data keys.
	 *
	 * Only supported for kernel keyblocks, not firmware keyblocks.
	 */
	struct vb2_signature keyblock_hash;

	/* Flags for key (VB2_KEYBLOCK_FLAG_*) */
	uint32_t keyblock_flags;
	uint32_t reserved1;

	/* Key to verify the chunk of data */
	struct vb2_packed_key data_key;
} __attribute__((packed));

#define EXPECTED_VB2_KEYBLOCK_SIZE 112

#endif  /* VBOOT_REFERENCE_2STRUCT_H_ */

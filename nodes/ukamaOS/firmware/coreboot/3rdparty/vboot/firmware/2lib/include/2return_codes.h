/* Copyright (c) 2014 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef VBOOT_REFERENCE_2RETURN_CODES_H_
#define VBOOT_REFERENCE_2RETURN_CODES_H_

#include "2sysincludes.h"

/*
 * Functions which return an error all return this type.  This is a 32-bit
 * value rather than an int so it's consistent across different architectures.
 */
typedef uint32_t vb2_error_t;

/*
 * Return codes from verified boot functions.
 *
 * Note that other values may be passed through from vb2ex_*() calls; see
 * the comment for VB2_ERROR_EX below.
 */
enum vb2_return_code {
	/* Success - no error */
	VB2_SUCCESS = 0,

	/*
	 * All vboot2 error codes start at a large offset from zero, to reduce
	 * the risk of overlap with other error codes (TPM, etc.).
	 */
	VB2_ERROR_BASE = 0x10000000,

	/* Unknown / unspecified error */
	VB2_ERROR_UNKNOWN = VB2_ERROR_BASE + 1,

	/* Mock error for testing */
	VB2_ERROR_MOCK,

	/**********************************************************************
	 * ancient pre-vboot1 errors
	 * TODO: deprecate these once they have all moved over to vboot2 style
	 */
	/* Invalid public key passed to a signature verficiation function. */
	VBOOT_PUBLIC_KEY_INVALID = 4,
	/* Preamble internal structure is invalid */
	VBOOT_PREAMBLE_INVALID = 5,
	/* Shared data is invalid. */
	VBOOT_SHARED_DATA_INVALID = 7,
	/* Kernel Preamble does not contain flags */
	VBOOT_KERNEL_PREAMBLE_NO_FLAGS = 8,

	/**********************************************************************
	 * vboot1-style errors
	 * TODO: deprecate these once they have all moved over to vboot2 style
	 */
	/* Unable to initialize shared data */
	VBERROR_INIT_SHARED_DATA              = 0x10001,
	/* Unable to set boot mode state in TPM */
	VBERROR_TPM_SET_BOOT_MODE_STATE       = 0x10006,
	/* Calling firmware needs to perform a reboot. */
	VBERROR_REBOOT_REQUIRED               = 0x10007,
	/* Unable to set up TPM */
	VBERROR_TPM_FIRMWARE_SETUP            = 0x10008,
	/* Unable to read kernel versions from TPM */
	VBERROR_TPM_READ_KERNEL               = 0x10009,
	/* Unable to write kernel versions to TPM */
	VBERROR_TPM_WRITE_KERNEL              = 0x1000B,
	/* Unable to lock kernel versions in TPM */
	VBERROR_TPM_LOCK_KERNEL               = 0x1000C,
	/* Calling firmware requested shutdown via VbExIsShutdownRequested() */
	VBERROR_SHUTDOWN_REQUESTED            = 0x1000D,
	/* Unable to find a boot device on which to look for a kernel */
	VBERROR_NO_DISK_FOUND                 = 0x1000E,
	/* No OS kernel found on any boot device */
	VBERROR_NO_KERNEL_FOUND               = 0x1000F,
	/* All OS kernels found were invalid (corrupt, improperly signed...) */
	VBERROR_INVALID_KERNEL_FOUND          = 0x10010,
	/* VbSelectAndLoadKernel() requested recovery mode */
	VBERROR_LOAD_KERNEL_RECOVERY          = 0x10011,
	/* Other error inside VbSelectAndLoadKernel() */
	VBERROR_LOAD_KERNEL                   = 0x10012,
	/* Invalid Google binary block */
	VBERROR_INVALID_GBB                   = 0x10013,
	/* Invalid bitmap volume */
	VBERROR_INVALID_BMPFV                 = 0x10014,
	/* Invalid screen index */
	VBERROR_INVALID_SCREEN_INDEX          = 0x10015,
	/* VbExBeep() can't make sound in the background */
	VBERROR_NO_BACKGROUND_SOUND           = 0x10019,
	/* Need EC to reboot to read-only code to switch RW slot */
	VBERROR_EC_REBOOT_TO_SWITCH_RW        = 0x1001A,
	/* Need EC to reboot to read-only code */
	VBERROR_EC_REBOOT_TO_RO_REQUIRED      = 0x10022,
	/* No image present */
	VBERROR_NO_IMAGE_PRESENT              = 0x10026,
	/* failed to draw screen */
	VBERROR_SCREEN_DRAW                   = 0x10027,
	/* Error reading FWMP from TPM (note: not present is not an error) */
	VBERROR_TPM_READ_FWMP                 = 0x10029,
	/* Peripheral busy. Cannot upgrade firmware at present. */
	VBERROR_PERIPHERAL_BUSY               = 0x10030,
	/* Error writing VPD */
	VBERROR_VPD_WRITE                     = 0x10032,
	/* Detachable UI internal functions may return the following codes */
	/* No error; return to UI loop */
	VBERROR_KEEP_LOOPING			= 0x30000,

	/**********************************************************************
	 * SHA errors
	 */
	VB2_ERROR_SHA = VB2_ERROR_BASE + 0x010000,

	/* Bad algorithm in vb2_digest_init() */
	VB2_ERROR_SHA_INIT_ALGORITHM,

	/* Bad algorithm in vb2_digest_extend() */
	VB2_ERROR_SHA_EXTEND_ALGORITHM,

	/* Bad algorithm in vb2_digest_finalize() */
	VB2_ERROR_SHA_FINALIZE_ALGORITHM,

	/* Digest size buffer too small in vb2_digest_finalize() */
	VB2_ERROR_SHA_FINALIZE_DIGEST_SIZE,

	/**********************************************************************
	 * RSA errors
	 */
	VB2_ERROR_RSA = VB2_ERROR_BASE + 0x020000,

	/* Padding mismatch in vb2_check_padding() */
	VB2_ERROR_RSA_PADDING,

	/* Bad algorithm in vb2_check_padding() */
	VB2_ERROR_RSA_PADDING_ALGORITHM,

	/* Null param passed to vb2_verify_digest() */
	VB2_ERROR_RSA_VERIFY_PARAM,

	/* Bad algorithm in vb2_verify_digest() */
	VB2_ERROR_RSA_VERIFY_ALGORITHM,

	/* Bad signature length in vb2_verify_digest() */
	VB2_ERROR_RSA_VERIFY_SIG_LEN,

	/* Work buffer too small in vb2_verify_digest() */
	VB2_ERROR_RSA_VERIFY_WORKBUF,

	/* Digest mismatch in vb2_verify_digest() */
	VB2_ERROR_RSA_VERIFY_DIGEST,

	/* Bad size calculation in vb2_check_padding() */
	VB2_ERROR_RSA_PADDING_SIZE,

	/**********************************************************************
	 * NV storage errors
	 */
	VB2_ERROR_NV = VB2_ERROR_BASE + 0x030000,

	/* Bad header in vb2_nv_check_crc() */
	VB2_ERROR_NV_HEADER,

	/* Bad CRC in vb2_nv_check_crc() */
	VB2_ERROR_NV_CRC,

	/**********************************************************************
	 * Secure data storage errors
	 */
	VB2_ERROR_SECDATA = VB2_ERROR_BASE + 0x040000,

	/* Bad CRC in vb2api_secdata_firmware_check() */
	VB2_ERROR_SECDATA_FIRMWARE_CRC,

	/* Bad struct version in vb2api_secdata_firmware_check() */
	VB2_ERROR_SECDATA_FIRMWARE_VERSION,

	/* Invalid param in vb2_secdata_firmware_get();
	   Deprecated with chromium:972956. */
	VB2_ERROR_DEPRECATED_SECDATA_FIRMWARE_GET_PARAM,

	/* Invalid param in vb2_secdata_firmware_set();
	   Deprecated with chromium:972956. */
	VB2_ERROR_DEPRECATED_SECDATA_FIRMWARE_SET_PARAM,

	/* Invalid flags passed to vb2_secdata_firmware_set();
	   Deprecated with chromium:972956. */
	VB2_ERROR_DEPRECATED_SECDATA_FIRMWARE_SET_FLAGS,

	/* Called vb2_secdata_firmware_get() with uninitialized secdata;
	   Deprecated with chromium:972956. */
	VB2_ERROR_DEPRECATED_SECDATA_FIRMWARE_GET_UNINITIALIZED,

	/* Called vb2_secdata_firmware_set() with uninitialized secdata;
	   Deprecated with chromium:972956. */
	VB2_ERROR_DEPRECATED_SECDATA_FIRMWARE_SET_UNINITIALIZED,

	/* Bad CRC in vb2api_secdata_kernel_check() */
	VB2_ERROR_SECDATA_KERNEL_CRC,

	/* Bad struct version in vb2_secdata_kernel_init() */
	VB2_ERROR_SECDATA_KERNEL_VERSION,

	/* Bad uid in vb2_secdata_kernel_init() */
	VB2_ERROR_SECDATA_KERNEL_UID,

	/* Invalid param in vb2_secdata_kernel_get();
	   Deprecated with chromium:972956. */
	VB2_ERROR_DEPRECATED_SECDATA_KERNEL_GET_PARAM,

	/* Invalid param in vb2_secdata_kernel_set();
	   Deprecated with chromium:972956. */
	VB2_ERROR_DEPRECATED_SECDATA_KERNEL_SET_PARAM,

	/* Invalid flags passed to vb2_secdata_kernel_set();
	   Deprecated with chromium:972956. */
	VB2_ERROR_DEPRECATED_SECDATA_KERNEL_SET_FLAGS,

	/* Called vb2_secdata_kernel_get() with uninitialized secdata_kernel;
	   Deprecated with chromium:972956. */
	VB2_ERROR_DEPRECATED_SECDATA_KERNEL_GET_UNINITIALIZED,

	/* Called vb2_secdata_kernel_set() with uninitialized secdata_kernel;
	   Deprecated with chromium:972956. */
	VB2_ERROR_DEPRECATED_SECDATA_KERNEL_SET_UNINITIALIZED,

	/* Bad size in vb2api_secdata_fwmp_check() */
	VB2_ERROR_SECDATA_FWMP_SIZE,

	/* Incomplete structure in vb2api_secdata_fwmp_check() */
	VB2_ERROR_SECDATA_FWMP_INCOMPLETE,

	/* Bad CRC in vb2api_secdata_fwmp_check() */
	VB2_ERROR_SECDATA_FWMP_CRC,

	/* Bad struct version in vb2_secdata_fwmp_check() */
	VB2_ERROR_SECDATA_FWMP_VERSION,

	/**********************************************************************
	 * Common code errors
	 */
	VB2_ERROR_COMMON = VB2_ERROR_BASE + 0x050000,

	/* Buffer is smaller than alignment offset in vb2_align() */
	VB2_ERROR_ALIGN_BIGGER_THAN_SIZE,

	/* Buffer is smaller than request in vb2_align() */
	VB2_ERROR_ALIGN_SIZE,

	/* Parent wraps around in vb2_verify_member_inside() */
	VB2_ERROR_INSIDE_PARENT_WRAPS,

	/* Member wraps around in vb2_verify_member_inside() */
	VB2_ERROR_INSIDE_MEMBER_WRAPS,

	/* Member outside parent in vb2_verify_member_inside() */
	VB2_ERROR_INSIDE_MEMBER_OUTSIDE,

	/* Member data wraps around in vb2_verify_member_inside() */
	VB2_ERROR_INSIDE_DATA_WRAPS,

	/* Member data outside parent in vb2_verify_member_inside() */
	VB2_ERROR_INSIDE_DATA_OUTSIDE,

	/* Unsupported signature algorithm in vb2_unpack_key_buffer() */
	VB2_ERROR_UNPACK_KEY_SIG_ALGORITHM,                      /* 0x150008 */

	/* Bad key size in vb2_unpack_key_buffer() */
	VB2_ERROR_UNPACK_KEY_SIZE,

	/* Bad key alignment in vb2_unpack_key_buffer() */
	VB2_ERROR_UNPACK_KEY_ALIGN,

	/* Bad key array size in vb2_unpack_key_buffer() */
	VB2_ERROR_UNPACK_KEY_ARRAY_SIZE,

	/* Bad algorithm in vb2_verify_data() */
	VB2_ERROR_VDATA_ALGORITHM,

	/* Incorrect signature size for algorithm in vb2_verify_data() */
	VB2_ERROR_VDATA_SIG_SIZE,

	/* Data smaller than length of signed data in vb2_verify_data() */
	VB2_ERROR_VDATA_NOT_ENOUGH_DATA,

	/* Not enough work buffer for digest in vb2_verify_data() */
	VB2_ERROR_VDATA_WORKBUF_DIGEST,

	/* Not enough work buffer for hash temp data in vb2_verify_data() */
	VB2_ERROR_VDATA_WORKBUF_HASHING,                         /* 0x150010 */

	/*
	 * Bad digest size in vb2_verify_data() - probably because algorithm
	 * is bad.
	 */
	VB2_ERROR_VDATA_DIGEST_SIZE,

	/* Unsupported hash algorithm in vb2_unpack_key_buffer() */
	VB2_ERROR_UNPACK_KEY_HASH_ALGORITHM,

	/* Member data overlaps member header */
	VB2_ERROR_INSIDE_DATA_OVERLAP,

	/* Unsupported packed key struct version */
	VB2_ERROR_UNPACK_KEY_STRUCT_VERSION,

	/*
	 * Buffer too small for total, fixed size, or description reported in
	 * common header, or member data checked via
	 * vb21_verify_common_member().
	 */
	VB2_ERROR_COMMON_TOTAL_SIZE,
	VB2_ERROR_COMMON_FIXED_SIZE,
	VB2_ERROR_COMMON_DESC_SIZE,
	VB2_ERROR_COMMON_MEMBER_SIZE,                            /* 0x150018 */

	/*
	 * Total, fixed, description, or member offset/size not a multiple of
	 * 32 bits.
	 */
	VB2_ERROR_COMMON_TOTAL_UNALIGNED,
	VB2_ERROR_COMMON_FIXED_UNALIGNED,
	VB2_ERROR_COMMON_DESC_UNALIGNED,
	VB2_ERROR_COMMON_MEMBER_UNALIGNED,

	/* Common struct description or member data wraps address space */
	VB2_ERROR_COMMON_DESC_WRAPS,
	VB2_ERROR_COMMON_MEMBER_WRAPS,

	/* Common struct description is not null-terminated */
	VB2_ERROR_COMMON_DESC_TERMINATOR,

	/* Member data overlaps previous data */
	VB2_ERROR_COMMON_MEMBER_OVERLAP,                         /* 0x150020 */

	/* Signature bad magic number */
	VB2_ERROR_SIG_MAGIC,

	/* Signature incompatible version */
	VB2_ERROR_SIG_VERSION,

	/* Signature header doesn't fit */
	VB2_ERROR_SIG_HEADER_SIZE,

	/* Signature unsupported algorithm */
	VB2_ERROR_SIG_ALGORITHM,

	/* Signature bad size for algorithm */
	VB2_ERROR_SIG_SIZE,

	/* Wrong amount of data signed */
	VB2_ERROR_VDATA_SIZE,

	/* Digest mismatch */
	VB2_ERROR_VDATA_VERIFY_DIGEST,

	/* Key algorithm doesn't match signature algorithm */
	VB2_ERROR_VDATA_ALGORITHM_MISMATCH,

	/* Bad magic number in vb2_unpack_key_buffer() */
	VB2_ERROR_UNPACK_KEY_MAGIC,

	/* Null public key buffer passed to vb2_unpack_key_buffer() */
	VB2_ERROR_UNPACK_KEY_BUFFER,

	/**********************************************************************
	 * Keyblock verification errors (all in vb2_verify_keyblock())
	 */
	VB2_ERROR_KEYBLOCK = VB2_ERROR_BASE + 0x060000,

	/* Data buffer too small for header */
	VB2_ERROR_KEYBLOCK_TOO_SMALL_FOR_HEADER,

	/* Magic number not present */
	VB2_ERROR_KEYBLOCK_MAGIC,

	/* Header version incompatible */
	VB2_ERROR_KEYBLOCK_HEADER_VERSION,

	/* Data buffer too small for keyblock */
	VB2_ERROR_KEYBLOCK_SIZE,

	/* Signature data offset outside keyblock */
	VB2_ERROR_KEYBLOCK_SIG_OUTSIDE,

	/* Signature signed more data than size of keyblock */
	VB2_ERROR_KEYBLOCK_SIGNED_TOO_MUCH,

	/* Signature signed less data than size of keyblock header */
	VB2_ERROR_KEYBLOCK_SIGNED_TOO_LITTLE,

	/* Signature invalid */
	VB2_ERROR_KEYBLOCK_SIG_INVALID,

	/* Data key outside keyblock */
	VB2_ERROR_KEYBLOCK_DATA_KEY_OUTSIDE,

	/* Data key outside signed part of keyblock */
	VB2_ERROR_KEYBLOCK_DATA_KEY_UNSIGNED,

	/* Signature signed wrong amount of data */
	VB2_ERROR_KEYBLOCK_SIGNED_SIZE,

	/* No signature matching key ID */
	VB2_ERROR_KEYBLOCK_SIG_ID,

	/**********************************************************************
	 * Preamble verification errors (all in vb2_verify_preamble())
	 */
	VB2_ERROR_PREAMBLE = VB2_ERROR_BASE + 0x070000,

	/* Preamble data too small to contain header */
	VB2_ERROR_PREAMBLE_TOO_SMALL_FOR_HEADER,

	/* Header version incompatible */
	VB2_ERROR_PREAMBLE_HEADER_VERSION,

	/* Header version too old */
	VB2_ERROR_PREAMBLE_HEADER_OLD,

	/* Data buffer too small for preamble */
	VB2_ERROR_PREAMBLE_SIZE,

	/* Signature data offset outside preamble */
	VB2_ERROR_PREAMBLE_SIG_OUTSIDE,

	/* Signature signed more data than size of preamble */
	VB2_ERROR_PREAMBLE_SIGNED_TOO_MUCH,

	/* Signature signed less data than size of preamble header */
	VB2_ERROR_PREAMBLE_SIGNED_TOO_LITTLE,

	/* Signature invalid */
	VB2_ERROR_PREAMBLE_SIG_INVALID,

	/* Body signature outside preamble */
	VB2_ERROR_PREAMBLE_BODY_SIG_OUTSIDE,

	/* Kernel subkey outside preamble */
	VB2_ERROR_PREAMBLE_KERNEL_SUBKEY_OUTSIDE,

	/* Bad magic number */
	VB2_ERROR_PREAMBLE_MAGIC,

	/* Hash is signed */
	VB2_ERROR_PREAMBLE_HASH_SIGNED,

	/* Bootloader outside signed portion of body */
	VB2_ERROR_PREAMBLE_BOOTLOADER_OUTSIDE,

	/* Vmlinuz header outside signed portion of body */
	VB2_ERROR_PREAMBLE_VMLINUZ_HEADER_OUTSIDE,

	/**********************************************************************
	 * Misc higher-level code errors
	 */
	VB2_ERROR_MISC = VB2_ERROR_BASE + 0x080000,

	/* Work buffer too small (see vb2api_init and vb2api_reinit) */
	VB2_ERROR_WORKBUF_SMALL,

	/* Work buffer unaligned (see vb2api_init and vb2api_reinit) */
	VB2_ERROR_WORKBUF_ALIGN,

	/* Work buffer too small in GBB-related function */
	VB2_ERROR_GBB_WORKBUF,

	/* Bad magic number in vb2_read_gbb_header() */
	VB2_ERROR_GBB_MAGIC,

	/* Incompatible version in vb2_read_gbb_header() */
	VB2_ERROR_GBB_VERSION,

	/* Old version in vb2_read_gbb_header() */
	VB2_ERROR_GBB_TOO_OLD,

	/* Header size too small in vb2_read_gbb_header() */
	VB2_ERROR_GBB_HEADER_SIZE,

	/* Work buffer too small for root key in vb2_load_fw_keyblock() */
	VB2_ERROR_FW_KEYBLOCK_WORKBUF_ROOT_KEY,

	/* Work buffer too small for header in vb2_load_fw_keyblock() */
	VB2_ERROR_FW_KEYBLOCK_WORKBUF_HEADER,

	/* Work buffer too small for keyblock in vb2_load_fw_keyblock() */
	VB2_ERROR_FW_KEYBLOCK_WORKBUF,

	/* Keyblock version out of range in vb2_load_fw_keyblock() */
	VB2_ERROR_FW_KEYBLOCK_VERSION_RANGE,

	/* Keyblock version rollback in vb2_load_fw_keyblock() */
	VB2_ERROR_FW_KEYBLOCK_VERSION_ROLLBACK,

	/* Missing firmware data key in vb2_load_fw_preamble() */
	VB2_ERROR_FW_PREAMBLE2_DATA_KEY,

	/* Work buffer too small for header in vb2_load_fw_preamble() */
	VB2_ERROR_FW_PREAMBLE2_WORKBUF_HEADER,

	/* Work buffer too small for preamble in vb2_load_fw_preamble() */
	VB2_ERROR_FW_PREAMBLE2_WORKBUF,

	/* Firmware version out of range in vb2_load_fw_preamble() */
	VB2_ERROR_FW_PREAMBLE_VERSION_RANGE,

	/* Firmware version rollback in vb2_load_fw_preamble() */
	VB2_ERROR_FW_PREAMBLE_VERSION_ROLLBACK,

	/* Not enough space in work buffer for resource object */
	VB2_ERROR_READ_RESOURCE_OBJECT_BUF,

	/* Work buffer too small for header in vb2_load_kernel_keyblock() */
	VB2_ERROR_KERNEL_KEYBLOCK_WORKBUF_HEADER,

	/* Work buffer too small for keyblock in vb2_load_kernel_keyblock() */
	VB2_ERROR_KERNEL_KEYBLOCK_WORKBUF,

	/* Keyblock version out of range in vb2_load_kernel_keyblock() */
	VB2_ERROR_KERNEL_KEYBLOCK_VERSION_RANGE,

	/* Keyblock version rollback in vb2_load_kernel_keyblock() */
	VB2_ERROR_KERNEL_KEYBLOCK_VERSION_ROLLBACK,

	/*
	 * Keyblock flags don't match current mode in
	 * vb2_load_kernel_keyblock().
	 */
	VB2_ERROR_KERNEL_KEYBLOCK_DEV_FLAG,
	VB2_ERROR_KERNEL_KEYBLOCK_REC_FLAG,

	/* Missing firmware data key in vb2_load_kernel_preamble() */
	VB2_ERROR_KERNEL_PREAMBLE2_DATA_KEY,

	/* Work buffer too small for header in vb2_load_kernel_preamble() */
	VB2_ERROR_KERNEL_PREAMBLE2_WORKBUF_HEADER,

	/* Work buffer too small for preamble in vb2_load_kernel_preamble() */
	VB2_ERROR_KERNEL_PREAMBLE2_WORKBUF,

	/* Kernel version out of range in vb2_load_kernel_preamble() */
	VB2_ERROR_KERNEL_PREAMBLE_VERSION_RANGE,

	/* Kernel version rollback in vb2_load_kernel_preamble() */
	VB2_ERROR_KERNEL_PREAMBLE_VERSION_ROLLBACK,

	/* Kernel preamble not loaded before calling vb2api_get_kernel_size() */
	VB2_ERROR_API_GET_KERNEL_SIZE_PREAMBLE,

	/* Unable to unpack kernel subkey in vb2_verify_vblock() */
	VB2_ERROR_VBLOCK_KERNEL_SUBKEY,

	/*
	 * Got a self-signed kernel in vb2_verify_vblock(), but need an
	 * officially signed one.
	 */
	VB2_ERROR_VBLOCK_SELF_SIGNED,

	/* Invalid keyblock hash in vb2_verify_vblock() */
	VB2_ERROR_VBLOCK_KEYBLOCK_HASH,

	/* Invalid keyblock in vb2_verify_vblock() */
	VB2_ERROR_VBLOCK_KEYBLOCK,

	/* Wrong developer key hash in vb2_verify_vblock() */
	VB2_ERROR_VBLOCK_DEV_KEY_HASH,

	/* Work buffer too small in vb2_load_partition() */
	VB2_ERROR_LOAD_PARTITION_WORKBUF,

	/* Unable to read vblock in vb2_load_partition() */
	VB2_ERROR_LOAD_PARTITION_READ_VBLOCK,

	/* Unable to verify vblock in vb2_load_partition() */
	VB2_ERROR_LOAD_PARTITION_VERIFY_VBLOCK,

	/* Kernel body offset too large in vb2_load_partition() */
	VB2_ERROR_LOAD_PARTITION_BODY_OFFSET,

	/* Kernel body too big in vb2_load_partition() */
	VB2_ERROR_LOAD_PARTITION_BODY_SIZE,

	/* Unable to read kernel body in vb2_load_partition() */
	VB2_ERROR_LOAD_PARTITION_READ_BODY,

	/* Unable to unpack data key in vb2_load_partition() */
	VB2_ERROR_LOAD_PARTITION_DATA_KEY,

	/* Unable to verify body in vb2_load_partition() */
	VB2_ERROR_LOAD_PARTITION_VERIFY_BODY,

	/* Unable to get EC image hash in ec_sync_phase1() */
	VB2_ERROR_EC_HASH_IMAGE,

	/* Unable to get expected EC image hash in ec_sync_phase1() */
	VB2_ERROR_EC_HASH_EXPECTED,

	/* Expected and image hashes are different size in ec_sync_phase1() */
	VB2_ERROR_EC_HASH_SIZE,

	/* Incompatible version for vb2_shared_data structure being loaded */
	VB2_ERROR_SHARED_DATA_VERSION,

	/* Bad magic number in vb2_shared_data structure */
	VB2_ERROR_SHARED_DATA_MAGIC,

	/* Some part of GBB data is invalid */
	VB2_ERROR_GBB_INVALID,

	/* Invalid parameter */
	VB2_ERROR_INVALID_PARAMETER,

	/* Problem with workbuf validity (see vb2api_init and vb2api_reinit) */
	VB2_ERROR_WORKBUF_INVALID,

	/**********************************************************************
	 * API-level errors
	 */
	VB2_ERROR_API = VB2_ERROR_BASE + 0x090000,

	/* Bad tag in vb2api_init_hash() */
	VB2_ERROR_API_INIT_HASH_TAG,

	/* Preamble not present in vb2api_init_hash() */
	VB2_ERROR_API_INIT_HASH_PREAMBLE,

	/* Work buffer too small in vb2api_init_hash() */
	VB2_ERROR_API_INIT_HASH_WORKBUF,

	/* Missing firmware data key in vb2api_init_hash() */
	VB2_ERROR_API_INIT_HASH_DATA_KEY,

	/* Uninitialized work area in vb2api_extend_hash() */
	VB2_ERROR_API_EXTEND_HASH_WORKBUF,

	/* Too much data hashed in vb2api_extend_hash() */
	VB2_ERROR_API_EXTEND_HASH_SIZE,

	/* Preamble not present in vb2api_check_hash() */
	VB2_ERROR_API_CHECK_HASH_PREAMBLE,

	/* Uninitialized work area in vb2api_check_hash() */
	VB2_ERROR_API_CHECK_HASH_WORKBUF,

	/* Wrong amount of data hashed in vb2api_check_hash() */
	VB2_ERROR_API_CHECK_HASH_SIZE,

	/* Work buffer too small in vb2api_check_hash() */
	VB2_ERROR_API_CHECK_HASH_WORKBUF_DIGEST,

	/* Bad tag in vb2api_check_hash() */
	VB2_ERROR_API_CHECK_HASH_TAG,

	/* Missing firmware data key in vb2api_check_hash() */
	VB2_ERROR_API_CHECK_HASH_DATA_KEY,

	/* Signature size mismatch in vb2api_check_hash() */
	VB2_ERROR_API_CHECK_HASH_SIG_SIZE,

	/* Phase one needs recovery mode */
	VB2_ERROR_API_PHASE1_RECOVERY,

	/* Bad tag in vb2api_check_hash() */
	VB2_ERROR_API_INIT_HASH_ID,

	/* Signature mismatch in vb2api_check_hash() */
	VB2_ERROR_API_CHECK_HASH_SIG,

	/* Invalid enum vb2_pcr_digest requested to vb2api_get_pcr_digest */
	VB2_ERROR_API_PCR_DIGEST,

	/* Buffer size for the digest is too small for vb2api_get_pcr_digest */
	VB2_ERROR_API_PCR_DIGEST_BUF,

	/* Work buffer too small for recovery key in vb2api_kernel_phase1() */
	VB2_ERROR_API_KPHASE1_WORKBUF_REC_KEY,

	/* Firmware preamble not present for vb2api_kernel_phase1() */
	VB2_ERROR_API_KPHASE1_PREAMBLE,

	/* Wrong amount of kernel data in vb2api_verify_kernel_data() */
	VB2_ERROR_API_VERIFY_KDATA_SIZE,

	/* Kernel preamble not present for vb2api_verify_kernel_data() */
	VB2_ERROR_API_VERIFY_KDATA_PREAMBLE,

	/* Insufficient workbuf for hashing in vb2api_verify_kernel_data() */
	VB2_ERROR_API_VERIFY_KDATA_WORKBUF,

	/* Bad data key in vb2api_verify_kernel_data() */
	VB2_ERROR_API_VERIFY_KDATA_KEY,

	/* Phase one passing through secdata's request to reboot */
	VB2_ERROR_API_PHASE1_SECDATA_REBOOT,

	/* Digest buffer passed into vb2api_check_hash incorrect. */
	VB2_ERROR_API_CHECK_DIGEST_SIZE,

	/**********************************************************************
	 * Errors which may be generated by implementations of vb2ex functions.
	 * Implementation may also return its own specific errors, which should
	 * NOT be in the range VB2_ERROR_BASE...VB2_ERROR_MAX to avoid
	 * conflicting with future vboot2 error codes.
	 */
	VB2_ERROR_EX = VB2_ERROR_BASE + 0x0a0000,

	/* Read resource not implemented
	 * Deprecated: use VB2_ERROR_EX_UNIMPLEMENTED (chromium:944804) */
	VB2_ERROR_EX_DEPRECATED_READ_RESOURCE_UNIMPLEMENTED,

	/* Resource index not found */
	VB2_ERROR_EX_READ_RESOURCE_INDEX,

	/* Size of resource not big enough for requested offset and/or size */
	VB2_ERROR_EX_READ_RESOURCE_SIZE,

	/* TPM clear owner failed */
	VB2_ERROR_EX_TPM_CLEAR_OWNER,

	/* TPM clear owner not implemented
	 * Deprecated: use VB2_ERROR_EX_UNIMPLEMENTED (chromium:944804) */
	VB2_ERROR_DEPRECATED_EX_TPM_CLEAR_OWNER_UNIMPLEMENTED,

	/* Hardware crypto engine doesn't support this algorithm (non-fatal) */
	VB2_ERROR_EX_HWCRYPTO_UNSUPPORTED,

	/* TPM does not understand this command */
	VB2_ERROR_EX_TPM_NO_SUCH_COMMAND,

	/* vb2ex function is unimplemented (stubbed in 2lib/2stub.c) */
	VB2_ERROR_EX_UNIMPLEMENTED,

	/**********************************************************************
	 * Errors generated by host library (non-firmware) start here.
	 */
	VB2_ERROR_HOST_BASE = 0x20000000,

	/**********************************************************************
	 * Errors generated by host library misc functions
	 */
	VB2_ERROR_HOST_MISC = VB2_ERROR_HOST_BASE + 0x010000,

	/* Unable to open file in read_file() */
	VB2_ERROR_READ_FILE_OPEN,

	/* Bad size in read_file() */
	VB2_ERROR_READ_FILE_SIZE,

	/* Unable to allocate buffer in read_file() */
	VB2_ERROR_READ_FILE_ALLOC,

	/* Unable to read data in read_file() */
	VB2_ERROR_READ_FILE_DATA,

	/* Unable to open file in write_file() */
	VB2_ERROR_WRITE_FILE_OPEN,

	/* Unable to write data in write_file() */
	VB2_ERROR_WRITE_FILE_DATA,

	/* Unable to convert string to struct vb_id */
	VB2_ERROR_STR_TO_ID,

	/**********************************************************************
	 * Errors generated by host library key functions
	 */
	VB2_ERROR_HOST_KEY = VB2_ERROR_HOST_BASE + 0x020000,

	/* Unable to allocate key  in vb2_private_key_read_pem() */
	VB2_ERROR_READ_PEM_ALLOC,

	/* Unable to open .pem file in vb2_private_key_read_pem() */
	VB2_ERROR_READ_PEM_FILE_OPEN,

	/* Bad RSA data from .pem file in vb2_private_key_read_pem() */
	VB2_ERROR_READ_PEM_RSA,

	/* Unable to set private key description */
	VB2_ERROR_PRIVATE_KEY_SET_DESC,

	/* Bad magic number in vb2_private_key_unpack() */
	VB2_ERROR_UNPACK_PRIVATE_KEY_MAGIC,

	/* Bad common header in vb2_private_key_unpack() */
	VB2_ERROR_UNPACK_PRIVATE_KEY_HEADER,

	/* Bad key data in vb2_private_key_unpack() */
	VB2_ERROR_UNPACK_PRIVATE_KEY_DATA,

	/* Bad struct version in vb2_private_key_unpack() */
	VB2_ERROR_UNPACK_PRIVATE_KEY_STRUCT_VERSION,

	/* Unable to allocate buffer in vb2_private_key_unpack() */
	VB2_ERROR_UNPACK_PRIVATE_KEY_ALLOC,

	/* Unable to unpack RSA key in vb2_private_key_unpack() */
	VB2_ERROR_UNPACK_PRIVATE_KEY_RSA,

	/* Unable to set description in vb2_private_key_unpack() */
	VB2_ERROR_UNPACK_PRIVATE_KEY_DESC,

	/* Bad bare hash key in vb2_private_key_unpack() */
	VB2_ERROR_UNPACK_PRIVATE_KEY_HASH,

	/* Unable to create RSA data in vb2_private_key_write() */
	VB2_ERROR_PRIVATE_KEY_WRITE_RSA,

	/* Unable to allocate packed key buffer in vb2_private_key_write() */
	VB2_ERROR_PRIVATE_KEY_WRITE_ALLOC,

	/* Unable to write file in vb2_private_key_write() */
	VB2_ERROR_PRIVATE_KEY_WRITE_FILE,

	/* Bad algorithm in vb2_private_key_hash() */
	VB2_ERROR_PRIVATE_KEY_HASH,

	/* Unable to determine key size in vb2_public_key_alloc() */
	VB2_ERROR_PUBLIC_KEY_ALLOC_SIZE,

	/* Unable to allocate buffer in vb2_public_key_alloc() */
	VB2_ERROR_PUBLIC_KEY_ALLOC,

	/* Unable to set public key description */
	VB2_ERROR_PUBLIC_KEY_SET_DESC,

	/* Unable to read key data in vb2_public_key_read_keyb() */
	VB2_ERROR_READ_KEYB_DATA,

	/* Wrong amount of data read in vb2_public_key_read_keyb() */
	VB2_ERROR_READ_KEYB_SIZE,

	/* Unable to allocate key buffer in vb2_public_key_read_keyb() */
	VB2_ERROR_READ_KEYB_ALLOC,

	/* Error unpacking RSA arrays in vb2_public_key_read_keyb() */
	VB2_ERROR_READ_KEYB_UNPACK,

	/* Unable to read key data in vb2_packed_key_read() */
	VB2_ERROR_READ_PACKED_KEY_DATA,

	/* Bad key data in vb2_packed_key_read() */
	VB2_ERROR_READ_PACKED_KEY,

	/* Unable to determine key size in vb2_public_key_pack() */
	VB2_ERROR_PUBLIC_KEY_PACK_SIZE,

	/* Bad hash algorithm in vb2_public_key_hash() */
	VB2_ERROR_PUBLIC_KEY_HASH,

	/* Bad key size in vb2_copy_packed_key() */
	VB2_ERROR_COPY_KEY_SIZE,

	/* Unable to convert back to vb1 crypto algorithm */
	VB2_ERROR_VB1_CRYPTO_ALGORITHM,

	/* Unable to allocate packed key */
	VB2_ERROR_PACKED_KEY_ALLOC,

	/* Unable to copy packed key */
	VB2_ERROR_PACKED_KEY_COPY,

	/**********************************************************************
	 * Errors generated by host library signature functions
	 */
	VB2_ERROR_HOST_SIG = VB2_ERROR_HOST_BASE + 0x030000,

	/* Bad hash algorithm in vb2_digest_info() */
	VB2_ERROR_DIGEST_INFO,

	/*
	 * Unable to determine signature size for key algorithm in
	 * vb2_sig_size_for_key().
	 */
	VB2_ERROR_SIG_SIZE_FOR_KEY,

	/* Bad signature size in vb2_sign_data() */
	VB2_SIGN_DATA_SIG_SIZE,

	/* Unable to get digest info in vb2_sign_data() */
	VB2_SIGN_DATA_DIGEST_INFO,

	/* Unable to get digest size in vb2_sign_data() */
	VB2_SIGN_DATA_DIGEST_SIZE,

	/* Unable to allocate digest buffer in vb2_sign_data() */
	VB2_SIGN_DATA_DIGEST_ALLOC,

	/* Unable to initialize digest in vb2_sign_data() */
	VB2_SIGN_DATA_DIGEST_INIT,

	/* Unable to extend digest in vb2_sign_data() */
	VB2_SIGN_DATA_DIGEST_EXTEND,

	/* Unable to finalize digest in vb2_sign_data() */
	VB2_SIGN_DATA_DIGEST_FINALIZE,

	/* RSA encrypt failed in vb2_sign_data() */
	VB2_SIGN_DATA_RSA_ENCRYPT,

	/* Not enough buffer space to hold signature in vb2_sign_object() */
	VB2_SIGN_OBJECT_OVERFLOW,

	/**********************************************************************
	 * Errors generated by host library keyblock functions
	 */
	VB2_ERROR_HOST_KEYBLOCK = VB2_ERROR_HOST_BASE + 0x040000,

	/* Unable to determine signature sizes for vb2_create_keyblock() */
	VB2_KEYBLOCK_CREATE_SIG_SIZE,

	/* Unable to pack data key for vb2_create_keyblock() */
	VB2_KEYBLOCK_CREATE_DATA_KEY,

	/* Unable to allocate buffer in vb2_create_keyblock() */
	VB2_KEYBLOCK_CREATE_ALLOC,

	/* Unable to sign keyblock in vb2_create_keyblock() */
	VB2_KEYBLOCK_CREATE_SIGN,

	/**********************************************************************
	 * Errors generated by host library firmware preamble functions
	 */
	VB2_ERROR_HOST_FW_PREAMBLE = VB2_ERROR_HOST_BASE + 0x050000,

	/* Unable to determine signature sizes for vb2_create_fw_preamble() */
	VB2_FW_PREAMBLE_CREATE_SIG_SIZE,

	/* Unable to allocate buffer in vb2_create_fw_preamble() */
	VB2_FW_PREAMBLE_CREATE_ALLOC,

	/* Unable to sign preamble in vb2_create_fw_preamble() */
	VB2_FW_PREAMBLE_CREATE_SIGN,

	/**********************************************************************
	 * Errors generated by unit test functions
	 */
	VB2_ERROR_UNIT_TEST = VB2_ERROR_HOST_BASE + 0x060000,

	/* Unable to open an input file needed for a unit test */
	VB2_ERROR_TEST_INPUT_FILE,

	/**********************************************************************
	 * Highest non-zero error generated inside vboot library.  Note that
	 * error codes passed through vboot when it calls external APIs may
	 * still be outside this range.
	 */
	VB2_ERROR_MAX = VB2_ERROR_BASE + 0x1fffffff,
};

#endif  /* VBOOT_REFERENCE_2RETURN_CODES_H_ */

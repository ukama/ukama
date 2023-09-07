/* Copyright (c) 2014 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Recovery reasons
 */

#ifndef VBOOT_REFERENCE_2RECOVERY_REASONS_H_
#define VBOOT_REFERENCE_2RECOVERY_REASONS_H_

/* Recovery reason codes */
enum vb2_nv_recovery {

	/**********************************************************************/
	/**** Uncategorized errors ********************************************/

	/* Recovery not requested. */
	VB2_RECOVERY_NOT_REQUESTED = 0x00,

	/*
	 * Recovery requested from legacy utility.  (Prior to the NV storage
	 * spec, recovery mode was a single bitfield; this value is reserved so
	 * that scripts which wrote 1 to the recovery field are distinguishable
	 * from scripts whch use the recovery reasons listed here.
	 */
	VB2_RECOVERY_LEGACY = 0x01,

	/* User manually requested recovery via recovery button */
	VB2_RECOVERY_RO_MANUAL = 0x02,



	/**********************************************************************/
	/**** Firmware verification (RO) errors (and some EC stuff???) ********/

	/* Unspecified RW verification error (when none of 0x10-0x1f fit) */
	VB2_RECOVERY_RO_INVALID_RW = 0x03,

	/* S3 resume failed (deprecated) */
	VB2_RECOVERY_DEPRECATED_RO_S3_RESUME = 0x04,

	/* TPM error in read-only firmware (deprecated, see 0x54+) */
	VB2_RECOVERY_DEPRECATED_RO_TPM_ERROR = 0x05,

	/* Shared data error in read-only firmware */
	VB2_RECOVERY_RO_SHARED_DATA = 0x06,

	/* Test error from S3Resume() (deprecated) */
	VB2_RECOVERY_DEPRECATED_RO_TEST_S3 = 0x07,

	/* Test error from LoadFirmwareSetup() (deprecated) */
	VB2_RECOVERY_DEPRECATED_RO_TEST_LFS = 0x08,

	/* Test error from LoadFirmware() (deprecated) */
	VB2_RECOVERY_DEPRECATED_RO_TEST_LF = 0x09,

	/*
	 * RW firmware failed signature check (neither RW firmware slot was
	 * valid).  Recovery reason is VB2_RECOVERY_RO_INVALID_RW_CHECK_MIN +
	 * the check value for the slot which came closest to validating; see
	 * VBSD_LF_CHECK_* in vboot_struct.h.
	 */
	VB2_RECOVERY_RO_INVALID_RW_CHECK_MIN = 0x10,

	/* Latest tried RW firmware keyblock verification failed */
	VB2_RECOVERY_FW_KEYBLOCK = 0x13,

	/* Latest tried RW firmware key version too old */
	VB2_RECOVERY_FW_KEY_ROLLBACK = 0x14,

	/* Latest tried RW firmware preamble verification failed */
	VB2_RECOVERY_FW_PREAMBLE = 0x16,

	/* Latest tried RW firmware version too old */
	VB2_RECOVERY_FW_ROLLBACK = 0x17,

	/* Latest tried RW firmware body verification failed */
	VB2_RECOVERY_FW_BODY = 0x1b,

	/* Highest reason for failed RW firmware signature check */
	VB2_RECOVERY_RO_INVALID_RW_CHECK_MAX = 0x1f,

	/*
	 * Firmware boot failure outside of verified boot (RAM init, missing
	 * SSD, etc.).
	 */
	VB2_RECOVERY_RO_FIRMWARE = 0x20,

	/*
	 * Recovery mode TPM initialization requires a system reboot.  The
	 * system was already in recovery mode for some other reason when this
	 * happened.
	 */
	VB2_RECOVERY_RO_TPM_REBOOT = 0x21,

	/* EC software sync - other error */
	VB2_RECOVERY_EC_SOFTWARE_SYNC = 0x22,

	/* EC software sync - unable to determine active EC image */
	VB2_RECOVERY_EC_UNKNOWN_IMAGE = 0x23,

	/* EC software sync - error obtaining EC image hash (deprecated) */
	VB2_RECOVERY_DEPRECATED_EC_HASH = 0x24,

	/* EC software sync - error obtaining expected EC image */
	VB2_RECOVERY_EC_EXPECTED_IMAGE = 0x25,

	/* EC software sync - error updating EC */
	VB2_RECOVERY_EC_UPDATE = 0x26,

	/* EC software sync - unable to jump to EC-RW */
	VB2_RECOVERY_EC_JUMP_RW = 0x27,

	/* EC software sync - unable to protect / unprotect EC-RW */
	VB2_RECOVERY_EC_PROTECT = 0x28,

	/* EC software sync - error obtaining expected EC hash */
	VB2_RECOVERY_EC_EXPECTED_HASH = 0x29,

	/* EC software sync - expected EC image doesn't match hash (deprc.) */
	VB2_RECOVERY_DEPRECATED_EC_HASH_MISMATCH = 0x2a,

	/* Firmware secure data initialization error */
	VB2_RECOVERY_SECDATA_FIRMWARE_INIT = 0x2b,

	/* GBB header is bad */
	VB2_RECOVERY_GBB_HEADER = 0x2c,

	/* Unable to clear TPM owner */
	VB2_RECOVERY_TPM_CLEAR_OWNER = 0x2d,

	/* Error determining/updating virtual dev switch */
	VB2_RECOVERY_DEV_SWITCH = 0x2e,

	/* Error determining firmware slot */
	VB2_RECOVERY_FW_SLOT = 0x2f,

	/* Error updating AUX firmware */
	VB2_RECOVERY_AUX_FW_UPDATE = 0x30,

	/* Unspecified/unknown error in read-only firmware */
	VB2_RECOVERY_RO_UNSPECIFIED = 0x3f,



	/**********************************************************************/
	/**** Kernel verification (RW) errors *********************************/

	/*
	 * User manually requested recovery by pressing a key at developer
	 * warning screen (deprecated)
	 */
	VB2_RECOVERY_DEPRECATED_RW_DEV_SCREEN = 0x41,

	/* No OS kernel detected (deprecated, now 0x5b) */
	VB2_RECOVERY_DEPRECATED_RW_NO_OS = 0x42,

	/* OS kernel failed signature check. Since the kernel corrupts itself
	   (DMVERROR) on a verity failure, may also indicate corrupt rootfs. */
	VB2_RECOVERY_RW_INVALID_OS = 0x43,

	/* TPM error in rewritable firmware (deprecated, see 0x54+) */
	VB2_RECOVERY_DEPRECATED_RW_TPM_ERROR = 0x44,

	/* RW firmware in dev mode, but dev switch is off (deprecated) */
	VB2_RECOVERY_DEPRECATED_RW_DEV_MISMATCH = 0x45,

	/* Shared data error in rewritable firmware */
	VB2_RECOVERY_RW_SHARED_DATA = 0x46,

	/* Test error from LoadKernel() (deprecated) */
	VB2_RECOVERY_DEPRECATED_RW_TEST_LK = 0x47,

	/* No bootable disk found (deprecated, see 0x5a) */
	VB2_RECOVERY_DEPRECATED_RW_NO_DISK = 0x48,

	/* Rebooting did not correct TPM_E_FAIL or TPM_E_FAILEDSELFTEST  */
	VB2_RECOVERY_TPM_E_FAIL = 0x49,

	/* TPM setup error in read-only firmware */
	VB2_RECOVERY_RO_TPM_S_ERROR = 0x50,

	/* TPM write error in read-only firmware */
	VB2_RECOVERY_RO_TPM_W_ERROR = 0x51,

	/* TPM lock error in read-only firmware */
	VB2_RECOVERY_RO_TPM_L_ERROR = 0x52,

	/* TPM update error in read-only firmware */
	VB2_RECOVERY_RO_TPM_U_ERROR = 0x53,

	/* TPM read error in rewritable firmware */
	VB2_RECOVERY_RW_TPM_R_ERROR = 0x54,

	/* TPM write error in rewritable firmware */
	VB2_RECOVERY_RW_TPM_W_ERROR = 0x55,

	/* TPM lock error in rewritable firmware */
	VB2_RECOVERY_RW_TPM_L_ERROR = 0x56,

	/* EC software sync unable to get EC image hash */
	VB2_RECOVERY_EC_HASH_FAILED = 0x57,

	/* EC software sync invalid image hash size */
	VB2_RECOVERY_EC_HASH_SIZE = 0x58,

	/* Unspecified error while trying to load kernel */
	VB2_RECOVERY_LK_UNSPECIFIED = 0x59,

	/* No bootable storage device in system */
	VB2_RECOVERY_RW_NO_DISK = 0x5a,

	/* No bootable kernel found on disk */
	VB2_RECOVERY_RW_NO_KERNEL = 0x5b,

	/* BCB related error in RW firmware (deprecated) */
	VB2_RECOVERY_DEPRECATED_RW_BCB_ERROR = 0x5c,

	/* Kernel secure data initialization error */
	VB2_RECOVERY_SECDATA_KERNEL_INIT = 0x5d,

	/* Fastboot mode requested in firmware (deprecated) */
	VB2_RECOVERY_DEPRECATED_FW_FASTBOOT = 0x5e,

	/* Recovery hash space lock error in RO firmware */
	VB2_RECOVERY_RO_TPM_REC_HASH_L_ERROR = 0x5f,

	/* Failed to disable the TPM [prior to running untrusted code] */
	VB2_RECOVERY_TPM_DISABLE_FAILED = 0x60,

	/* Alt FW Failed hash verification */
	VB2_RECOVERY_ALTFW_HASH_FAILED = 0x61,

	/* Unspecified/unknown error in rewritable firmware */
	VB2_RECOVERY_RW_UNSPECIFIED = 0x7f,



	/**********************************************************************/
	/**** OS level (kernel) errors (deprecated) ***************************/

	/*
	 * Note: we want to avoid having the kernel touch vboot NVRAM directly
	 * in the future, so this whole range is essentially deprecated until
	 * further notice.
	 */

	/* DM-verity error (deprecated) */
	VB2_RECOVERY_DEPRECATED_KE_DM_VERITY = 0x81,

	/* Unspecified/unknown error in kernel (deprecated) */
	VB2_RECOVERY_DEPRECATED_KE_UNSPECIFIED = 0xbf,



	/**********************************************************************/
	/**** OS level (userspace) errors *************************************/

	/* Recovery mode test from user-mode */
	VB2_RECOVERY_US_TEST = 0xc1,

	/* Recovery requested by user-mode via BCB (deprecated) */
	VB2_RECOVERY_DEPRECATED_BCB_USER_MODE = 0xc2,

	/* Fastboot mode requested by user-mode (deprecated) */
	VB2_RECOVERY_DEPRECATED_US_FASTBOOT = 0xc3,

	/* User requested recovery for training memory and rebooting. */
	VB2_RECOVERY_TRAIN_AND_REBOOT = 0xc4,

	/* Unspecified/unknown error in user-mode */
	VB2_RECOVERY_US_UNSPECIFIED = 0xff,
};

#endif  /* VBOOT_REFERENCE_2RECOVERY_REASONS_H_ */

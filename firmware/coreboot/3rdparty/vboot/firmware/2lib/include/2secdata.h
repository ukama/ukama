/* Copyright (c) 2014 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Secure non-volatile storage routines
 */

#ifndef VBOOT_REFERENCE_2SECDATA_H_
#define VBOOT_REFERENCE_2SECDATA_H_

#include "2api.h"

/*****************************************************************************/
/* Firmware secure storage space */

/* Which param to get/set for vb2_secdata_firmware_get/set() */
enum vb2_secdata_firmware_param {
	/* Flags; see vb2_secdata_firmware_flags */
	VB2_SECDATA_FIRMWARE_FLAGS = 0,

	/* Firmware versions */
	VB2_SECDATA_FIRMWARE_VERSIONS,
};

/* Flags for firmware space */
enum vb2_secdata_firmware_flags {
	/*
	 * Last boot was developer mode.  TPM ownership is cleared when
	 * transitioning to/from developer mode.  Set/cleared by
	 * vb2_check_dev_switch().
	 */
	VB2_SECDATA_FIRMWARE_FLAG_LAST_BOOT_DEVELOPER = (1 << 0),

	/*
	 * Virtual developer mode switch is on.  Set/cleared by the
	 * keyboard-controlled dev screens in recovery mode.  Cleared by
	 * vb2_check_dev_switch().
	 */
	VB2_SECDATA_FIRMWARE_FLAG_DEV_MODE = (1 << 1),
};

/**
 * Initialize firmware secure storage context and verify its CRC.
 *
 * This must be called before vb2_secdata_firmware_get/set().
 *
 * @param ctx		Context pointer
 * @return VB2_SUCCESS, or non-zero error code if error.
 */
vb2_error_t vb2_secdata_firmware_init(struct vb2_context *ctx);

/**
 * Read a firmware secure storage value.
 *
 * @param ctx		Context pointer
 * @param param		Parameter to read
 * @return Requested parameter value
 */
uint32_t vb2_secdata_firmware_get(struct vb2_context *ctx,
				  enum vb2_secdata_firmware_param param);

/**
 * Write a firmware secure storage value.
 *
 * @param ctx		Context pointer
 * @param param		Parameter to write
 * @param value		New value
 */
void vb2_secdata_firmware_set(struct vb2_context *ctx,
			      enum vb2_secdata_firmware_param param,
			      uint32_t value);

/*****************************************************************************/
/* Kernel secure storage space
 *
 * These are separate functions so that they don't bloat the size of the early
 * boot code which uses the firmware version space functions.
 */

/* Which param to get/set for vb2_secdata_kernel_get/set() */
enum vb2_secdata_kernel_param {
	/* Kernel versions */
	VB2_SECDATA_KERNEL_VERSIONS = 0,
};

/**
 * Initialize kernel secure storage context and verify its CRC.
 *
 * This must be called before vb2_secdata_kernel_get/set().
 *
 * @param ctx		Context pointer
 * @return VB2_SUCCESS, or non-zero error code if error.
 */
vb2_error_t vb2_secdata_kernel_init(struct vb2_context *ctx);

/**
 * Read a kernel secure storage value.
 *
 * @param ctx		Context pointer
 * @param param		Parameter to read
 * @return Requested parameter value
 */
uint32_t vb2_secdata_kernel_get(struct vb2_context *ctx,
				enum vb2_secdata_kernel_param param);

/**
 * Write a kernel secure storage value.
 *
 * @param ctx		Context pointer
 * @param param		Parameter to write
 * @param value		New value
 */
void vb2_secdata_kernel_set(struct vb2_context *ctx,
			    enum vb2_secdata_kernel_param param,
			    uint32_t value);

/*****************************************************************************/
/* Firmware management parameters (FWMP) space */

/* Flags for FWMP space */
enum vb2_secdata_fwmp_flags {
	VB2_SECDATA_FWMP_DEV_DISABLE_BOOT = (1 << 0),
	VB2_SECDATA_FWMP_DEV_DISABLE_RECOVERY = (1 << 1),
	VB2_SECDATA_FWMP_DEV_ENABLE_USB = (1 << 2),
	VB2_SECDATA_FWMP_DEV_ENABLE_LEGACY = (1 << 3),
	VB2_SECDATA_FWMP_DEV_ENABLE_OFFICIAL_ONLY = (1 << 4),
	VB2_SECDATA_FWMP_DEV_USE_KEY_HASH = (1 << 5),
	/* CCD = case-closed debugging on cr50; flag implemented on cr50 */
	VB2_SECDATA_FWMP_DEV_DISABLE_CCD_UNLOCK = (1 << 6),
};

/**
 * Initialize FWMP secure storage context and verify its CRC.
 *
 * This must be called before vb2_secdata_fwmp_get_flag/get_dev_key_hash().
 *
 * @param ctx		Context pointer
 * @return VB2_SUCCESS, or non-zero error code if error.
 */
vb2_error_t vb2_secdata_fwmp_init(struct vb2_context *ctx);

/**
 * Read a FWMP secure storage flag value.
 *
 * It is unsupported to call before successfully running vb2_secdata_fwmp_init.
 * In this case, vboot will fail and exit.
 *
 * @param ctx		Context pointer
 * @param flag		Flag to read
 * @return current flag value (0 or 1)
 */
int vb2_secdata_fwmp_get_flag(struct vb2_context *ctx,
			      enum vb2_secdata_fwmp_flags flag);

/**
 * Return a pointer to FWMP dev key hash.
 *
 * @param ctx		Context pointer
 * @return uint8_t pointer to dev_key_hash field
 */
uint8_t *vb2_secdata_fwmp_get_dev_key_hash(struct vb2_context *ctx);

#endif  /* VBOOT_REFERENCE_2SECDATA_H_ */

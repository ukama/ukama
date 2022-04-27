/* Copyright (c) 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Common functions between firmware and kernel verified boot.
 */

#ifndef VBOOT_REFERENCE_VBOOT_COMMON_H_
#define VBOOT_REFERENCE_VBOOT_COMMON_H_

#include "2api.h"
#include "2struct.h"
#include "vboot_struct.h"

/**
 * Initialize a public key to refer to [key_data].
 */
void PublicKeyInit(struct vb2_packed_key *key,
		   uint8_t *key_data, uint64_t key_size);

/**
 * Copy a public key from [src] to [dest].
 *
 * Returns 0 if success, non-zero if error.
 */
int PublicKeyCopy(struct vb2_packed_key *dest,
		  const struct vb2_packed_key *src);

/**
 * Verify that the Vmlinuz Header is contained inside of the kernel blob.
 *
 * Returns VB2_SUCCESS or VBOOT_PREAMBLE_INVALID on error
 */
vb2_error_t VerifyVmlinuzInsideKBlob(uint64_t kblob, uint64_t kblob_size,
				     uint64_t header, uint64_t header_size);
/**
 * Initialize a verified boot shared data structure.
 *
 * Returns 0 if success, non-zero if error.
 */
vb2_error_t VbSharedDataInit(VbSharedDataHeader *header, uint64_t size);

/**
 * Reserve [size] bytes of the shared data area.  Returns the offset of the
 * reserved data from the start of the shared data buffer, or 0 if error.
 */
uint64_t VbSharedDataReserve(VbSharedDataHeader *header, uint64_t size);

/**
 * Copy the kernel subkey into the shared data.
 *
 * Returns 0 if success, non-zero if error.
 */
vb2_error_t VbSharedDataSetKernelKey(VbSharedDataHeader *header,
				     const struct vb2_packed_key *src);

/**
 * Check whether recovery is allowed or not.
 *
 * The only way to pass this check and proceed to the recovery process is to
 * physically request a recovery (a.k.a. manual recovery). All other recovery
 * requests including manual recovery requested by a (compromised) host will
 * end up with 'broken' screen.
 *
 * @param ctx vboot2 context pointer
 * @return 1: Yes. 0: No or not sure.
 */
int vb2_allow_recovery(struct vb2_context *ctx);

#endif  /* VBOOT_REFERENCE_VBOOT_COMMON_H_ */

/* Copyright (c) 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Common functions between firmware and kernel verified boot.
 * (Firmware portion)
 */

#include "2common.h"
#include "2misc.h"
#include "2rsa.h"
#include "2sha.h"
#include "2sysincludes.h"
#include "utility.h"
#include "vboot_api.h"
#include "vboot_common.h"

void PublicKeyInit(struct vb2_packed_key *key,
		   uint8_t *key_data, uint64_t key_size)
{
	key->key_offset = vb2_offset_of(key, key_data);
	key->key_size = key_size;
	key->algorithm = VB2_ALG_COUNT; /* Key not present yet */
	key->key_version = 0;
}

int PublicKeyCopy(struct vb2_packed_key *dest, const struct vb2_packed_key *src)
{
	if (dest->key_size < src->key_size)
		return 1;

	dest->key_size = src->key_size;
	dest->algorithm = src->algorithm;
	dest->key_version = src->key_version;
	memcpy(vb2_packed_key_data_mutable(dest),
	       vb2_packed_key_data(src),
	       src->key_size);
	return 0;
}

vb2_error_t VerifyVmlinuzInsideKBlob(uint64_t kblob, uint64_t kblob_size,
				     uint64_t header, uint64_t header_size)
{
	uint64_t end = header-kblob;
	if (end > kblob_size)
		return VBOOT_PREAMBLE_INVALID;
	if (UINT64_MAX - end < header_size)
		return VBOOT_PREAMBLE_INVALID;
	if (end + header_size > kblob_size)
		return VBOOT_PREAMBLE_INVALID;

	return VB2_SUCCESS;
}

uint64_t VbSharedDataReserve(VbSharedDataHeader *header, uint64_t size)
{
	if (!header || size > header->data_size - header->data_used) {
		VB2_DEBUG("VbSharedData buffer out of space.\n");
		return 0;  /* Not initialized, or not enough space left. */
	}

	uint64_t offs = header->data_used;
	VB2_DEBUG("VbSharedDataReserve %d bytes at %d\n", (int)size, (int)offs);

	header->data_used += size;
	return offs;
}

vb2_error_t VbSharedDataSetKernelKey(VbSharedDataHeader *header,
				     const struct vb2_packed_key *src)
{
	struct vb2_packed_key *kdest;

	if (!header)
		return VBOOT_SHARED_DATA_INVALID;
	if (!src)
		return VBOOT_PUBLIC_KEY_INVALID;

	kdest = &header->kernel_subkey;

	VB2_DEBUG("Saving kernel subkey to shared data: size %d, algo %d\n",
		  vb2_rsa_sig_size(vb2_crypto_to_signature(src->algorithm)),
		  (int)src->algorithm);

	/* Attempt to allocate space for key, if it hasn't been allocated yet */
	if (!header->kernel_subkey_data_offset) {
		header->kernel_subkey_data_offset =
			VbSharedDataReserve(header, src->key_size);
		if (!header->kernel_subkey_data_offset)
			return VBOOT_SHARED_DATA_INVALID;
		header->kernel_subkey_data_size = src->key_size;
	}

	/* Copy the kernel sign key blob into the destination buffer */
	PublicKeyInit(kdest,
		      (uint8_t *)header + header->kernel_subkey_data_offset,
		      header->kernel_subkey_data_size);

	return PublicKeyCopy(kdest, src);
}

int vb2_allow_recovery(struct vb2_context *ctx)
{
	/* VB2_GBB_FLAG_FORCE_MANUAL_RECOVERY forces this to always return
	   true. */
	if (vb2_get_gbb(ctx)->flags & VB2_GBB_FLAG_FORCE_MANUAL_RECOVERY)
		return 1;

	/*
	 * If EC is in RW, it implies recovery wasn't manually requested.
	 * On some platforms, EC_IN_RW can't be reset by the EC, thus, this may
	 * return false (=RW). That's ok because if recovery is manual, we will
	 * get the right signal and that's the case we care about.
	 */
	if (!VbExTrustEC(0))
		return 0;

	/* Now we confidently check the recovery switch state at boot */
	return !!(vb2_get_sd(ctx)->vbsd->flags & VBSD_BOOT_REC_SWITCH_ON);
}

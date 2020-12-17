/* Copyright (c) 2014 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Externally-callable APIs
 * (Firmware portion)
 */

#include "2api.h"
#include "2common.h"
#include "2misc.h"
#include "2nvstorage.h"
#include "2rsa.h"
#include "2secdata.h"
#include "2sha.h"
#include "2sysincludes.h"
#include "2tpm_bootmode.h"
#include "vb2_common.h"

vb2_error_t vb2api_fw_phase1(struct vb2_context *ctx)
{
	vb2_error_t rv;
	struct vb2_shared_data *sd = vb2_get_sd(ctx);

	/* Initialize NV context */
	vb2_nv_init(ctx);

	/*
	 * Handle caller-requested reboot due to secdata.  Do this before we
	 * even look at secdata.  If we fail because of a reboot loop we'll be
	 * the first failure so will get to set the recovery reason.
	 */
	if (!(ctx->flags & VB2_CONTEXT_SECDATA_WANTS_REBOOT)) {
		/* No reboot requested */
		vb2_nv_set(ctx, VB2_NV_TPM_REQUESTED_REBOOT, 0);
	} else if (vb2_nv_get(ctx, VB2_NV_TPM_REQUESTED_REBOOT)) {
		/*
		 * Reboot requested... again.  Fool me once, shame on you.
		 * Fool me twice, shame on me.  Fail into recovery to avoid
		 * a reboot loop.
		 */
		vb2api_fail(ctx, VB2_RECOVERY_RO_TPM_REBOOT, 0);
	} else {
		/* Reboot requested for the first time */
		vb2_nv_set(ctx, VB2_NV_TPM_REQUESTED_REBOOT, 1);
		return VB2_ERROR_API_PHASE1_SECDATA_REBOOT;
	}

	/* Initialize firmware secure data */
	rv = vb2_secdata_firmware_init(ctx);
	if (rv)
		vb2api_fail(ctx, VB2_RECOVERY_SECDATA_FIRMWARE_INIT, rv);

	/* Load and parse the GBB header */
	rv = vb2_fw_init_gbb(ctx);
	if (rv)
		vb2api_fail(ctx, VB2_RECOVERY_GBB_HEADER, rv);

	/*
	 * Check for recovery.  Note that this function returns void, since any
	 * errors result in requesting recovery.  That's also why we don't
	 * return error from failures in the preceding two steps; those
	 * failures simply cause us to detect recovery mode here.
	 */
	vb2_check_recovery(ctx);

	/* Check for dev switch */
	rv = vb2_check_dev_switch(ctx);
	if (rv && !(ctx->flags & VB2_CONTEXT_RECOVERY_MODE)) {
		/*
		 * Error in dev switch processing, and we weren't already
		 * headed for recovery mode.  Reboot into recovery mode, since
		 * it's too late to handle those errors this boot, and we need
		 * to take a different path through the dev switch checking
		 * code in that case.
		 */
		vb2api_fail(ctx, VB2_RECOVERY_DEV_SWITCH, rv);
		return rv;
	}

	/*
	 * Check for possible reasons to ask the firmware to make display
	 * available.  sd->recovery_reason may have been set above by
	 * vb2_check_recovery.  VB2_SD_FLAG_DEV_MODE_ENABLED may have been set
	 * above by vb2_check_dev_switch.  VB2_NV_DIAG_REQUEST may have been
	 * set during the last boot in recovery mode.
	 */
	if (!(ctx->flags & VB2_CONTEXT_DISPLAY_INIT) &&
	    (vb2_nv_get(ctx, VB2_NV_DISPLAY_REQUEST) ||
	     sd->flags & VB2_SD_FLAG_DEV_MODE_ENABLED ||
	     sd->recovery_reason ||
	     vb2_nv_get(ctx, VB2_NV_DIAG_REQUEST)))
		ctx->flags |= VB2_CONTEXT_DISPLAY_INIT;
	/* Mark display as available for downstream vboot and vboot callers. */
	if (ctx->flags & VB2_CONTEXT_DISPLAY_INIT)
		sd->flags |= VB2_SD_FLAG_DISPLAY_AVAILABLE;

	/* Return error if recovery is needed */
	if (ctx->flags & VB2_CONTEXT_RECOVERY_MODE) {
		/* Always clear RAM when entering recovery mode */
		ctx->flags |= VB2_CONTEXT_CLEAR_RAM;
		return VB2_ERROR_API_PHASE1_RECOVERY;
	}

	return VB2_SUCCESS;
}

vb2_error_t vb2api_fw_phase2(struct vb2_context *ctx)
{
	vb2_error_t rv;

	/*
	 * Use the slot from the last boot if this is a resume.  Do not set
	 * VB2_SD_STATUS_CHOSE_SLOT so the try counter is not decremented on
	 * failure as we are explicitly not attempting to boot from a new slot.
	 */
	if (ctx->flags & VB2_CONTEXT_S3_RESUME) {
		struct vb2_shared_data *sd = vb2_get_sd(ctx);

		/* Set the current slot to the last booted slot */
		sd->fw_slot = vb2_nv_get(ctx, VB2_NV_FW_TRIED);

		/* Set context flag if we're using slot B */
		if (sd->fw_slot)
			ctx->flags |= VB2_CONTEXT_FW_SLOT_B;

		return VB2_SUCCESS;
	}

	/* Always clear RAM when entering developer mode */
	if (ctx->flags & VB2_CONTEXT_DEVELOPER_MODE)
		ctx->flags |= VB2_CONTEXT_CLEAR_RAM;

	/* Check for explicit request to clear TPM */
	rv = vb2_check_tpm_clear(ctx);
	if (rv) {
		vb2api_fail(ctx, VB2_RECOVERY_TPM_CLEAR_OWNER, rv);
		return rv;
	}

	/* Decide which firmware slot to try this boot */
	rv = vb2_select_fw_slot(ctx);
	if (rv) {
		vb2api_fail(ctx, VB2_RECOVERY_FW_SLOT, rv);
		return rv;
	}

	return VB2_SUCCESS;
}

vb2_error_t vb2api_extend_hash(struct vb2_context *ctx,
		       const void *buf,
		       uint32_t size)
{
	struct vb2_shared_data *sd = vb2_get_sd(ctx);
	struct vb2_digest_context *dc = (struct vb2_digest_context *)
		vb2_member_of(sd, sd->hash_offset);

	/* Must have initialized hash digest work area */
	if (!sd->hash_size)
		return VB2_ERROR_API_EXTEND_HASH_WORKBUF;

	/* Don't extend past the data we expect to hash */
	if (!size || size > sd->hash_remaining_size)
		return VB2_ERROR_API_EXTEND_HASH_SIZE;

	sd->hash_remaining_size -= size;

	if (dc->using_hwcrypto)
		return vb2ex_hwcrypto_digest_extend(buf, size);
	else
		return vb2_digest_extend(dc, buf, size);
}

vb2_error_t vb2api_get_pcr_digest(struct vb2_context *ctx,
			  enum vb2_pcr_digest which_digest,
			  uint8_t *dest,
			  uint32_t *dest_size)
{
	const uint8_t *digest;
	uint32_t digest_size;

	switch (which_digest) {
	case BOOT_MODE_PCR:
		digest = vb2_get_boot_state_digest(ctx);
		digest_size = VB2_SHA1_DIGEST_SIZE;
		break;
	case HWID_DIGEST_PCR:
		digest = vb2_get_gbb(ctx)->hwid_digest;
		digest_size = VB2_GBB_HWID_DIGEST_SIZE;
		break;
	default:
		return VB2_ERROR_API_PCR_DIGEST;
	}

	if (digest == NULL || *dest_size < digest_size)
		return VB2_ERROR_API_PCR_DIGEST_BUF;

	memcpy(dest, digest, digest_size);
	if (digest_size < *dest_size)
		memset(dest + digest_size, 0, *dest_size - digest_size);

	*dest_size = digest_size;

	return VB2_SUCCESS;
}

vb2_error_t vb2api_fw_phase3(struct vb2_context *ctx)
{
	vb2_error_t rv;

	/* Verify firmware keyblock */
	rv = vb2_load_fw_keyblock(ctx);
	if (rv) {
		vb2api_fail(ctx, VB2_RECOVERY_RO_INVALID_RW, rv);
		return rv;
	}

	/* Verify firmware preamble */
	rv = vb2_load_fw_preamble(ctx);
	if (rv) {
		vb2api_fail(ctx, VB2_RECOVERY_RO_INVALID_RW, rv);
		return rv;
	}

	return VB2_SUCCESS;
}

vb2_error_t vb2api_init_hash(struct vb2_context *ctx, uint32_t tag,
			     uint32_t *size)
{
	struct vb2_shared_data *sd = vb2_get_sd(ctx);
	const struct vb2_fw_preamble *pre;
	struct vb2_digest_context *dc;
	struct vb2_public_key key;
	struct vb2_workbuf wb;
	vb2_error_t rv;

	vb2_workbuf_from_ctx(ctx, &wb);

	if (tag == VB2_HASH_TAG_INVALID)
		return VB2_ERROR_API_INIT_HASH_TAG;

	/* Get preamble pointer */
	if (!sd->preamble_size)
		return VB2_ERROR_API_INIT_HASH_PREAMBLE;
	pre = (const struct vb2_fw_preamble *)
		vb2_member_of(sd, sd->preamble_offset);

	/* For now, we only support the firmware body tag */
	if (tag != VB2_HASH_TAG_FW_BODY)
		return VB2_ERROR_API_INIT_HASH_TAG;

	/* Allocate workbuf space for the hash */
	if (sd->hash_size) {
		dc = (struct vb2_digest_context *)
			vb2_member_of(sd, sd->hash_offset);
	} else {
		uint32_t dig_size = sizeof(*dc);

		dc = vb2_workbuf_alloc(&wb, dig_size);
		if (!dc)
			return VB2_ERROR_API_INIT_HASH_WORKBUF;

		sd->hash_offset = vb2_offset_of(sd, dc);
		sd->hash_size = dig_size;
		vb2_set_workbuf_used(ctx, sd->hash_offset + dig_size);
	}

	/*
	 * Work buffer now contains:
	 *   - vb2_shared_data
	 *   - packed firmware data key
	 *   - firmware preamble
	 *   - hash data
	 */

	/*
	 * Unpack the firmware data key to see which hashing algorithm we
	 * should use.
	 *
	 * TODO: really, the firmware body should be hashed, and not signed,
	 * because the signature we're checking is already signed as part of
	 * the firmware preamble.  But until we can change the signing scripts,
	 * we're stuck with a signature here instead of a hash.
	 */
	if (!sd->data_key_size)
		return VB2_ERROR_API_INIT_HASH_DATA_KEY;

	rv = vb2_unpack_key_buffer(&key,
				   vb2_member_of(sd, sd->data_key_offset),
				   sd->data_key_size);
	if (rv)
		return rv;

	sd->hash_tag = tag;
	sd->hash_remaining_size = pre->body_signature.data_size;

	if (size)
		*size = pre->body_signature.data_size;

	if (!(pre->flags & VB2_FIRMWARE_PREAMBLE_DISALLOW_HWCRYPTO)) {
		rv = vb2ex_hwcrypto_digest_init(key.hash_alg,
						pre->body_signature.data_size);
		if (!rv) {
			VB2_DEBUG("Using HW crypto engine for hash_alg %d\n",
				  key.hash_alg);
			dc->hash_alg = key.hash_alg;
			dc->using_hwcrypto = 1;
			return VB2_SUCCESS;
		}
		if (rv != VB2_ERROR_EX_HWCRYPTO_UNSUPPORTED)
			return rv;
		VB2_DEBUG("HW crypto for hash_alg %d not supported, using SW\n",
			  key.hash_alg);
	} else {
		VB2_DEBUG("HW crypto forbidden by preamble, using SW\n");
	}

	return vb2_digest_init(dc, key.hash_alg);
}

vb2_error_t vb2api_check_hash_get_digest(struct vb2_context *ctx,
					 void *digest_out,
					 uint32_t digest_out_size)
{
	struct vb2_shared_data *sd = vb2_get_sd(ctx);
	struct vb2_digest_context *dc = (struct vb2_digest_context *)
		vb2_member_of(sd, sd->hash_offset);
	struct vb2_workbuf wb;

	uint8_t *digest;
	uint32_t digest_size = vb2_digest_size(dc->hash_alg);

	struct vb2_fw_preamble *pre;
	struct vb2_public_key key;
	vb2_error_t rv;

	vb2_workbuf_from_ctx(ctx, &wb);

	/* Get preamble pointer */
	if (!sd->preamble_size)
		return VB2_ERROR_API_CHECK_HASH_PREAMBLE;
	pre = vb2_member_of(sd, sd->preamble_offset);

	/* Must have initialized hash digest work area */
	if (!sd->hash_size)
		return VB2_ERROR_API_CHECK_HASH_WORKBUF;

	/* Should have hashed the right amount of data */
	if (sd->hash_remaining_size)
		return VB2_ERROR_API_CHECK_HASH_SIZE;

	/* Allocate the digest */
	digest = vb2_workbuf_alloc(&wb, digest_size);
	if (!digest)
		return VB2_ERROR_API_CHECK_HASH_WORKBUF_DIGEST;

	/* Finalize the digest */
	if (dc->using_hwcrypto)
		rv = vb2ex_hwcrypto_digest_finalize(digest, digest_size);
	else
		rv = vb2_digest_finalize(dc, digest, digest_size);
	if (rv)
		return rv;

	/* The code below is specific to the body signature */
	if (sd->hash_tag != VB2_HASH_TAG_FW_BODY)
		return VB2_ERROR_API_CHECK_HASH_TAG;

	/*
	 * The body signature is currently a *signature* of the body data, not
	 * just its hash.  So we need to verify the signature.
	 */

	/* Unpack the data key */
	if (!sd->data_key_size)
		return VB2_ERROR_API_CHECK_HASH_DATA_KEY;

	rv = vb2_unpack_key_buffer(&key,
				   vb2_member_of(sd, sd->data_key_offset),
				   sd->data_key_size);
	if (rv)
		return rv;

	/*
	 * Check digest vs. signature.  Note that this destroys the signature.
	 * That's ok, because we only check each signature once per boot.
	 */
	rv = vb2_verify_digest(&key, &pre->body_signature, digest, &wb);
	if (rv)
		vb2api_fail(ctx, VB2_RECOVERY_FW_BODY, rv);

	if (digest_out != NULL) {
		if (digest_out_size < digest_size)
			return VB2_ERROR_API_CHECK_DIGEST_SIZE;
		memcpy(digest_out, digest, digest_size);
	}

	return rv;
}

int vb2api_check_hash(struct vb2_context *ctx)
{
	return vb2api_check_hash_get_digest(ctx, NULL, 0);
}

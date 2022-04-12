/* Copyright 2019 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Firmware management parameters (FWMP) APIs
 */

#include "2common.h"
#include "2crc8.h"
#include "2misc.h"
#include "2secdata.h"
#include "2secdata_struct.h"

vb2_error_t vb2api_secdata_fwmp_check(struct vb2_context *ctx, uint8_t *size)
{
	struct vb2_secdata_fwmp *sec =
		(struct vb2_secdata_fwmp *)&ctx->secdata_fwmp;

	/* Verify that at least the minimum size has been read */
	if (*size < VB2_SECDATA_FWMP_MIN_SIZE) {
		VB2_DEBUG("FWMP: missing %d bytes for minimum size\n",
			  VB2_SECDATA_FWMP_MIN_SIZE - *size);
		*size = VB2_SECDATA_FWMP_MIN_SIZE;
		return VB2_ERROR_SECDATA_FWMP_INCOMPLETE;
	}

	/* Verify that struct_size is reasonable */
	if (sec->struct_size < VB2_SECDATA_FWMP_MIN_SIZE ||
	    sec->struct_size > VB2_SECDATA_FWMP_MAX_SIZE) {
		VB2_DEBUG("FWMP: invalid size: %d\n", sec->struct_size);
		return VB2_ERROR_SECDATA_FWMP_SIZE;
	}

	/* Verify that we have read full structure */
	if (*size < sec->struct_size) {
		VB2_DEBUG("FWMP: missing %d bytes\n", sec->struct_size - *size);
		*size = sec->struct_size;
		return VB2_ERROR_SECDATA_FWMP_INCOMPLETE;
	}
	*size = sec->struct_size;

	/* Verify CRC */
	if (sec->crc8 != vb2_secdata_fwmp_crc(sec)) {
		VB2_DEBUG("FWMP: bad CRC\n");
		return VB2_ERROR_SECDATA_FWMP_CRC;
	}

	/* Verify major version is compatible */
	if ((sec->struct_version >> 4) != (VB2_SECDATA_FWMP_VERSION >> 4)) {
		VB2_DEBUG("FWMP: major version incompatible\n");
		return VB2_ERROR_SECDATA_FWMP_VERSION;
	}

	/*
	 * If this were a 1.1+ reader and the source was a 1.0 struct,
	 * we would need to take care of initializing the extra fields
	 * added in 1.1+.  But that's not an issue yet.
	 */
	return VB2_SUCCESS;
}

vb2_error_t vb2_secdata_fwmp_init(struct vb2_context *ctx)
{
	struct vb2_shared_data *sd = vb2_get_sd(ctx);
	struct vb2_secdata_fwmp *sec =
		(struct vb2_secdata_fwmp *)&ctx->secdata_fwmp;
	vb2_error_t rv;

	/* Skip checking if NO_SECDATA_FWMP is set. */
	if (!(ctx->flags & VB2_CONTEXT_NO_SECDATA_FWMP)) {
		rv = vb2api_secdata_fwmp_check(ctx, &sec->struct_size);
		if (rv)
			return rv;
	}

	/* Mark as initialized */
	sd->status |= VB2_SD_STATUS_SECDATA_FWMP_INIT;

	return VB2_SUCCESS;
}

int vb2_secdata_fwmp_get_flag(struct vb2_context *ctx,
			      enum vb2_secdata_fwmp_flags flag)
{
	struct vb2_shared_data *sd = vb2_get_sd(ctx);
	struct vb2_secdata_fwmp *sec =
		(struct vb2_secdata_fwmp *)&ctx->secdata_fwmp;

	if (!(sd->status & VB2_SD_STATUS_SECDATA_FWMP_INIT)) {
		if (ctx->flags & VB2_CONTEXT_RECOVERY_MODE) {
			VB2_DEBUG("Assuming broken FWMP flag %d as 0\n", flag);
			return 0;
		} else {
			VB2_DIE("Must init FWMP before retrieving flag\n");
		}
	}

	if (ctx->flags & VB2_CONTEXT_NO_SECDATA_FWMP)
		return 0;

	return !!(sec->flags & flag);
}

uint8_t *vb2_secdata_fwmp_get_dev_key_hash(struct vb2_context *ctx)
{
	struct vb2_shared_data *sd = vb2_get_sd(ctx);
	struct vb2_secdata_fwmp *sec =
		(struct vb2_secdata_fwmp *)&ctx->secdata_fwmp;

	if (!(sd->status & VB2_SD_STATUS_SECDATA_FWMP_INIT)) {
		if (ctx->flags & VB2_CONTEXT_RECOVERY_MODE) {
			VB2_DEBUG("Assuming broken FWMP dev_key_hash "
				  "as empty\n");
			return NULL;
		} else {
			VB2_DIE("Must init FWMP before getting dev key hash\n");
		}
	}

	if (ctx->flags & VB2_CONTEXT_NO_SECDATA_FWMP)
		return NULL;

	return sec->dev_key_hash;
}

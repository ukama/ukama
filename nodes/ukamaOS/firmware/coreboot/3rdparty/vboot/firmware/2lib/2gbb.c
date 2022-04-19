/* Copyright 2019 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * GBB accessor functions.
 */

#include "2common.h"
#include "2misc.h"

static vb2_error_t vb2_gbb_read_key(struct vb2_context *ctx, uint32_t offset,
				    uint32_t *size,
				    struct vb2_packed_key **keyp,
				    struct vb2_workbuf *wb)
{
	struct vb2_workbuf wblocal = *wb;
	vb2_error_t rv;

	/* Check offset and size. */
	if (offset < sizeof(struct vb2_gbb_header))
		return VB2_ERROR_GBB_INVALID;
	if (*size < sizeof(**keyp))
		return VB2_ERROR_GBB_INVALID;

	/* GBB header might be padded.  Retrieve the vb2_packed_key
	   header so we can find out what the real size is. */
	*keyp = vb2_workbuf_alloc(&wblocal, sizeof(**keyp));
	if (!*keyp)
		return VB2_ERROR_GBB_WORKBUF;
	rv = vb2ex_read_resource(ctx, VB2_RES_GBB, offset, *keyp,
				 sizeof(**keyp));
	if (rv)
		return rv;

	rv = vb2_verify_packed_key_inside(*keyp, *size, *keyp);
	if (rv)
		return rv;

	/* Deal with a zero-size key (used in testing). */
	*size = (*keyp)->key_offset + (*keyp)->key_size;
	*size = VB2_MAX(*size, sizeof(**keyp));

	/* Now that we know the real size of the key, retrieve the key
	   data, and write it on the workbuf, directly after vb2_packed_key. */
	*keyp = vb2_workbuf_realloc(&wblocal, sizeof(**keyp), *size);
	if (!*keyp)
		return VB2_ERROR_GBB_WORKBUF;

	rv = vb2ex_read_resource(ctx, VB2_RES_GBB,
				 offset + sizeof(**keyp),
				 (void *)*keyp + sizeof(**keyp),
				 *size - sizeof(**keyp));
	if (!rv)
		*wb = wblocal;
	return rv;
}

vb2_error_t vb2_gbb_read_root_key(struct vb2_context *ctx,
				  struct vb2_packed_key **keyp, uint32_t *size,
				  struct vb2_workbuf *wb)
{
	struct vb2_gbb_header *gbb = vb2_get_gbb(ctx);
	uint32_t size_in = gbb->rootkey_size;
	vb2_error_t ret = vb2_gbb_read_key(ctx, gbb->rootkey_offset,
					   &size_in, keyp, wb);
	if (size)
		*size = size_in;
	return ret;
}

vb2_error_t vb2_gbb_read_recovery_key(struct vb2_context *ctx,
				      struct vb2_packed_key **keyp,
				      uint32_t *size, struct vb2_workbuf *wb)
{
	struct vb2_gbb_header *gbb = vb2_get_gbb(ctx);
	uint32_t size_in = gbb->recovery_key_size;
	vb2_error_t ret = vb2_gbb_read_key(ctx, gbb->recovery_key_offset,
					   &size_in, keyp, wb);
	if (size)
		*size = size_in;
	return ret;
}

vb2_error_t vb2api_gbb_read_hwid(struct vb2_context *ctx, char *hwid,
				 uint32_t *size)
{
	struct vb2_gbb_header *gbb = vb2_get_gbb(ctx);
	uint32_t i;
	vb2_error_t ret;

	if (gbb->hwid_size == 0) {
		VB2_DEBUG("invalid HWID size %d\n", gbb->hwid_size);
		return VB2_ERROR_GBB_INVALID;
	}

	*size = VB2_MIN(*size, VB2_GBB_HWID_MAX_SIZE);
	*size = VB2_MIN(*size, gbb->hwid_size);

	ret = vb2ex_read_resource(ctx, VB2_RES_GBB, gbb->hwid_offset,
				  hwid, *size);
	if (ret) {
		VB2_DEBUG("read resource failure: %d\n", ret);
		return ret;
	}

	/* Count HWID size, and ensure that it fits in the given buffer. */
	for (i = 0; i < *size; i++) {
		if (hwid[i] == '\0') {
			*size = i + 1;
			break;
		}
	}
	if (hwid[*size - 1] != '\0')
		return VB2_ERROR_INVALID_PARAMETER;

	return VB2_SUCCESS;
}

vb2_gbb_flags_t vb2api_gbb_get_flags(struct vb2_context *ctx)
{
	struct vb2_gbb_header *gbb = vb2_get_gbb(ctx);
	return gbb->flags;
}

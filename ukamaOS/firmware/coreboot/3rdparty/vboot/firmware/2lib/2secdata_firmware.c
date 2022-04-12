/* Copyright (c) 2014 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Secure storage APIs
 */

#include "2common.h"
#include "2crc8.h"
#include "2misc.h"
#include "2secdata.h"
#include "2secdata_struct.h"
#include "2sysincludes.h"

vb2_error_t vb2api_secdata_firmware_check(struct vb2_context *ctx)
{
	struct vb2_secdata_firmware *sec =
		(struct vb2_secdata_firmware *)ctx->secdata_firmware;

	/* Verify CRC */
	if (sec->crc8 != vb2_crc8(sec, offsetof(struct vb2_secdata_firmware,
						crc8))) {
		VB2_DEBUG("secdata_firmware: bad CRC\n");
		return VB2_ERROR_SECDATA_FIRMWARE_CRC;
	}

	/* Verify version */
	if (sec->struct_version < VB2_SECDATA_FIRMWARE_VERSION) {
		VB2_DEBUG("secdata_firmware: version incompatible\n");
		return VB2_ERROR_SECDATA_FIRMWARE_VERSION;
	}

	return VB2_SUCCESS;
}

uint32_t vb2api_secdata_firmware_create(struct vb2_context *ctx)
{
	struct vb2_secdata_firmware *sec =
		(struct vb2_secdata_firmware *)ctx->secdata_firmware;

	/* Clear the entire struct */
	memset(sec, 0, sizeof(*sec));

	/* Set to current version */
	sec->struct_version = VB2_SECDATA_FIRMWARE_VERSION;

	/* Calculate initial CRC */
	sec->crc8 = vb2_crc8(sec, offsetof(struct vb2_secdata_firmware, crc8));

	/* Mark as changed */
	ctx->flags |= VB2_CONTEXT_SECDATA_FIRMWARE_CHANGED;

	return sizeof(*sec);
}

vb2_error_t vb2_secdata_firmware_init(struct vb2_context *ctx)
{
	struct vb2_shared_data *sd = vb2_get_sd(ctx);
	vb2_error_t rv;

	rv = vb2api_secdata_firmware_check(ctx);
	if (rv)
		return rv;

	/* Set status flag */
	sd->status |= VB2_SD_STATUS_SECDATA_FIRMWARE_INIT;

	/* Read this now to make sure crossystem has it even in rec mode */
	sd->fw_version_secdata =
		vb2_secdata_firmware_get(ctx, VB2_SECDATA_FIRMWARE_VERSIONS);

	return VB2_SUCCESS;
}

uint32_t vb2_secdata_firmware_get(struct vb2_context *ctx,
				  enum vb2_secdata_firmware_param param)
{
	struct vb2_shared_data *sd = vb2_get_sd(ctx);
	struct vb2_secdata_firmware *sec =
		(struct vb2_secdata_firmware *)ctx->secdata_firmware;
	const char *msg;

	if (!(sd->status & VB2_SD_STATUS_SECDATA_FIRMWARE_INIT)) {
		msg = "get before init";
		goto fail;
	}

	switch (param) {
	case VB2_SECDATA_FIRMWARE_FLAGS:
		return sec->flags;

	case VB2_SECDATA_FIRMWARE_VERSIONS:
		return sec->fw_versions;

	default:
		msg = "invalid param";
	}

 fail:
	if (!(ctx->flags & VB2_CONTEXT_RECOVERY_MODE))
		VB2_DIE("%s\n", msg);
	VB2_DEBUG("ERROR [%s] ignored in recovery mode\n", msg);
	return 0;
}

void vb2_secdata_firmware_set(struct vb2_context *ctx,
			      enum vb2_secdata_firmware_param param,
			      uint32_t value)
{
	struct vb2_shared_data *sd = vb2_get_sd(ctx);
	struct vb2_secdata_firmware *sec =
		(struct vb2_secdata_firmware *)ctx->secdata_firmware;
	const char *msg;

	if (!(sd->status & VB2_SD_STATUS_SECDATA_FIRMWARE_INIT)) {
		msg = "set before init";
		goto fail;
	}

	/* If not changing the value, just return early */
	if (value == vb2_secdata_firmware_get(ctx, param))
		return;

	switch (param) {
	case VB2_SECDATA_FIRMWARE_FLAGS:
		/* Make sure flags is in valid range */
		if (value > 0xff) {
			msg = "flags out of range";
			goto fail;
		}

		VB2_DEBUG("secdata_firmware flags updated from %#x to %#x\n",
			  sec->flags, value);
		sec->flags = value;
		break;

	case VB2_SECDATA_FIRMWARE_VERSIONS:
		VB2_DEBUG("secdata_firmware versions updated from "
			  "%#x to %#x\n",
			  sec->fw_versions, value);
		sec->fw_versions = value;
		break;

	default:
		msg = "invalid param";
		goto fail;
	}

	/* Regenerate CRC */
	sec->crc8 = vb2_crc8(sec, offsetof(struct vb2_secdata_firmware, crc8));
	ctx->flags |= VB2_CONTEXT_SECDATA_FIRMWARE_CHANGED;
	return;

 fail:
	if (!(ctx->flags & VB2_CONTEXT_RECOVERY_MODE))
		VB2_DIE("%s\n", msg);
	VB2_DEBUG("ERROR [%s] ignored in recovery mode\n", msg);
}

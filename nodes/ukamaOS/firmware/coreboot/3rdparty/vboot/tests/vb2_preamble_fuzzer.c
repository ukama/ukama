// Copyright 2019 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#include <assert.h>

#include "2api.h"
#include "2common.h"
#include "2misc.h"
#include "2rsa.h"
#include "2secdata.h"
#include "vboot_test.h"

static struct vb2_context *ctx;
__attribute__((aligned(VB2_WORKBUF_ALIGN)))
static uint8_t workbuf[VB2_FIRMWARE_WORKBUF_RECOMMENDED_SIZE];

static const uint8_t *mock_preamble;
static size_t mock_preamble_size;

/* Limit exposure of code for which we didn't set up the environment right. */
void vb2api_fail(struct vb2_context *c, uint8_t reason, uint8_t subcode)
{
	return;
}

void vb2_secdata_firmware_set(struct vb2_context *c,
			      enum vb2_secdata_firmware_param param,
			      uint32_t value)
{
	/* prevent abort from uninitialized secdata */
}

vb2_error_t vb2ex_read_resource(struct vb2_context *c,
				enum vb2_resource_index index, uint32_t offset,
				void *buf, uint32_t size)
{
	if (index != VB2_RES_FW_VBLOCK)
		return VB2_ERROR_EX_READ_RESOURCE_INDEX;

	/* The preamble_offset in our mock shared data is 0, so we can assume
	   that offset here is a direct offset into the preamble. */
	if (offset > mock_preamble_size || mock_preamble_size - offset < size)
		return VB2_ERROR_EX_READ_RESOURCE_SIZE;

	memcpy(buf, mock_preamble + offset, size);
	return VB2_SUCCESS;
}

/* Pretend that signature checks always succeed so the fuzzer can cover more. */
vb2_error_t vb2_check_padding(const uint8_t *sig,
			      const struct vb2_public_key *key)
{
	return VB2_SUCCESS;
}

vb2_error_t vb2_safe_memcmp(const void *s1, const void *s2, size_t size)
{
	return VB2_SUCCESS;
}

int LLVMFuzzerTestOneInput(const uint8_t *data, size_t size);
int LLVMFuzzerTestOneInput(const uint8_t *data, size_t size) {
	const size_t datakey_size = 4096;	// enough for all our signatures

	if (size < datakey_size)
		return 0;

	if (vb2api_init(workbuf, sizeof(workbuf), &ctx))
		abort();

	struct vb2_workbuf wb;
	vb2_workbuf_from_ctx(ctx, &wb);

	uint8_t *key = vb2_workbuf_alloc(&wb, datakey_size);
	assert(key);
	memcpy(key, data, datakey_size);

	mock_preamble = data + datakey_size;
	mock_preamble_size = size - datakey_size;

	struct vb2_shared_data *sd = vb2_get_sd(ctx);
	sd->data_key_offset = vb2_offset_of(sd, key);
	sd->data_key_size = datakey_size;
	vb2_set_workbuf_used(ctx, sd->data_key_offset + sd->data_key_size);

	sd->vblock_preamble_offset = 0;
	vb2_load_fw_preamble(ctx);

	return 0;
}

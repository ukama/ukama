// Copyright 2019 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#include <assert.h>

#include "2api.h"
#include "2common.h"
#include "2misc.h"
#include "2rsa.h"
#include "vboot_test.h"

static struct vb2_context *ctx;
__attribute__((aligned(VB2_WORKBUF_ALIGN)))
static uint8_t workbuf[VB2_FIRMWARE_WORKBUF_RECOMMENDED_SIZE];
static struct {
	struct vb2_gbb_header h;
	uint8_t rootkey[4096];
} gbb;

static const uint8_t *mock_keyblock;
static size_t mock_keyblock_size;

/* Limit exposure of code for which we didn't set up the environment right. */
void vb2api_fail(struct vb2_context *c, uint8_t reason, uint8_t subcode)
{
	return;
}

struct vb2_gbb_header *vb2_get_gbb(struct vb2_context *c)
{
	return &gbb.h;
}

vb2_error_t vb2ex_read_resource(struct vb2_context *c,
				enum vb2_resource_index index, uint32_t offset,
				void *buf, uint32_t size)
{
	const void *rbase;
	size_t rsize;

	switch (index) {
	case VB2_RES_GBB:
		rbase = &gbb;
		rsize = sizeof(gbb);
		break;
	case VB2_RES_FW_VBLOCK:
		rbase = mock_keyblock;
		rsize = mock_keyblock_size;
		break;
	default:
		return VB2_ERROR_EX_READ_RESOURCE_INDEX;
	}

	if (offset > rsize || rsize - offset < size)
		return VB2_ERROR_EX_READ_RESOURCE_SIZE;

	memcpy(buf, rbase + offset, size);
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

int LLVMFuzzerTestOneInput(const uint8_t* data, size_t size);
int LLVMFuzzerTestOneInput(const uint8_t* data, size_t size) {
	if (size < sizeof(gbb.rootkey))
		return 0;

	memset(&gbb.h, 0, sizeof(gbb.h));
	gbb.h.rootkey_offset = gbb.rootkey - (uint8_t *)&gbb;
	gbb.h.rootkey_size = sizeof(gbb.rootkey);

	memcpy(gbb.rootkey, data, sizeof(gbb.rootkey));
	mock_keyblock = data + sizeof(gbb.rootkey);
	mock_keyblock_size = size - sizeof(gbb.rootkey);

	if (vb2api_init(workbuf, sizeof(workbuf), &ctx))
		abort();

	vb2_load_fw_keyblock(ctx);

	return 0;
}

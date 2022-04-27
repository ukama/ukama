/* Copyright 2015 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Tests for kernel secure storage library.
 */

#include "2api.h"
#include "2common.h"
#include "2crc8.h"
#include "2misc.h"
#include "2secdata.h"
#include "2secdata_struct.h"
#include "2sysincludes.h"
#include "test_common.h"
#include "vboot_common.h"

static uint8_t workbuf[VB2_FIRMWARE_WORKBUF_RECOMMENDED_SIZE]
	__attribute__ ((aligned (VB2_WORKBUF_ALIGN)));
static struct vb2_context *ctx;
static struct vb2_shared_data *sd;
static struct vb2_secdata_kernel *sec;

static void reset_common_data(void)
{
	memset(workbuf, 0xaa, sizeof(workbuf));
	TEST_SUCC(vb2api_init(workbuf, sizeof(workbuf), &ctx),
		  "vb2api_init failed");

	sd = vb2_get_sd(ctx);

	sec = (struct vb2_secdata_kernel *)ctx->secdata_kernel;
}

static void test_changed(struct vb2_context *c, int changed, const char *why)
{
	if (changed)
		TEST_NEQ(c->flags & VB2_CONTEXT_SECDATA_KERNEL_CHANGED, 0, why);
	else
		TEST_EQ(c->flags & VB2_CONTEXT_SECDATA_KERNEL_CHANGED, 0, why);

	c->flags &= ~VB2_CONTEXT_SECDATA_KERNEL_CHANGED;
};

static void secdata_kernel_test(void)
{
	uint32_t v = 1;
	reset_common_data();

	/* Check size constant */
	TEST_EQ(VB2_SECDATA_KERNEL_SIZE, sizeof(struct vb2_secdata_kernel),
		"Struct size constant");

	/* Blank data is invalid */
	memset(&ctx->secdata_kernel, 0xa6, sizeof(ctx->secdata_kernel));
	TEST_EQ(vb2api_secdata_kernel_check(ctx),
		VB2_ERROR_SECDATA_KERNEL_CRC, "Check blank CRC");
	TEST_EQ(vb2_secdata_kernel_init(ctx),
		VB2_ERROR_SECDATA_KERNEL_CRC, "Init blank CRC");

	/* Ensure zeroed buffers are invalid */
	memset(&ctx->secdata_kernel, 0, sizeof(ctx->secdata_kernel));
	TEST_EQ(vb2_secdata_kernel_init(ctx), VB2_ERROR_SECDATA_KERNEL_VERSION,
		"Zeroed buffer (invalid version)");

	/* Try with bad version */
	TEST_EQ(vb2api_secdata_kernel_create(ctx), VB2_SECDATA_KERNEL_SIZE,
		"Create");
	sec->struct_version -= 1;
	sec->crc8 = vb2_crc8(sec, offsetof(struct vb2_secdata_kernel, crc8));
	TEST_EQ(vb2api_secdata_kernel_check(ctx),
		VB2_ERROR_SECDATA_KERNEL_VERSION, "Check invalid version");
	TEST_EQ(vb2_secdata_kernel_init(ctx),
		VB2_ERROR_SECDATA_KERNEL_VERSION, "Init invalid version");

	/* Create good data */
	vb2api_secdata_kernel_create(ctx);
	TEST_SUCC(vb2api_secdata_kernel_check(ctx), "Check created CRC");
	TEST_SUCC(vb2_secdata_kernel_init(ctx), "Init created CRC");
	TEST_NEQ(sd->status & VB2_SD_STATUS_SECDATA_KERNEL_INIT, 0,
		 "Init set SD status");
	sd->status &= ~VB2_SD_STATUS_SECDATA_KERNEL_INIT;
	test_changed(ctx, 1, "Create changes data");

	/* Now corrupt it */
	ctx->secdata_kernel[2]++;
	TEST_EQ(vb2api_secdata_kernel_check(ctx),
		VB2_ERROR_SECDATA_KERNEL_CRC, "Check invalid CRC");
	TEST_EQ(vb2_secdata_kernel_init(ctx),
		VB2_ERROR_SECDATA_KERNEL_CRC, "Init invalid CRC");

	/* Make sure UID is checked */

	vb2api_secdata_kernel_create(ctx);
	sec->uid++;
	sec->crc8 = vb2_crc8(sec, offsetof(struct vb2_secdata_kernel, crc8));
	TEST_EQ(vb2_secdata_kernel_init(ctx), VB2_ERROR_SECDATA_KERNEL_UID,
		"Init invalid struct UID");

	/* Read/write versions */
	vb2api_secdata_kernel_create(ctx);
	vb2_secdata_kernel_init(ctx);
	ctx->flags = 0;
	v = vb2_secdata_kernel_get(ctx, VB2_SECDATA_KERNEL_VERSIONS);
	TEST_EQ(v, 0, "Versions created 0");
	test_changed(ctx, 0, "Get doesn't change data");
	vb2_secdata_kernel_set(ctx, VB2_SECDATA_KERNEL_VERSIONS, 0x123456ff);
	test_changed(ctx, 1, "Set changes data");
	vb2_secdata_kernel_set(ctx, VB2_SECDATA_KERNEL_VERSIONS, 0x123456ff);
	test_changed(ctx, 0, "Set again doesn't change data");
	v = vb2_secdata_kernel_get(ctx, VB2_SECDATA_KERNEL_VERSIONS);
	TEST_EQ(v, 0x123456ff, "Versions changed");

	/* Invalid field fails */
	TEST_ABORT(vb2_secdata_kernel_get(ctx, -1), "Get invalid");
	TEST_ABORT(vb2_secdata_kernel_set(ctx, -1, 456), "Set invalid");
	test_changed(ctx, 0, "Set invalid field doesn't change data");

	/* Read/write uninitialized data fails */
	sd->status &= ~VB2_SD_STATUS_SECDATA_KERNEL_INIT;
	TEST_ABORT(vb2_secdata_kernel_get(ctx, VB2_SECDATA_KERNEL_VERSIONS),
		   "Get uninitialized");
	test_changed(ctx, 0, "Get uninitialized doesn't change data");
	TEST_ABORT(vb2_secdata_kernel_set(ctx, VB2_SECDATA_KERNEL_VERSIONS,
					  0x123456ff),
		   "Set uninitialized");
	test_changed(ctx, 0, "Set uninitialized doesn't change data");
}

int main(int argc, char* argv[])
{
	secdata_kernel_test();

	return gTestSuccess ? 0 : 255;
}

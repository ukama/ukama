/* Copyright 2019 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Tests for firmware management parameters (FWMP) library.
 */

#include "2common.h"
#include "2misc.h"
#include "2secdata.h"
#include "2secdata_struct.h"
#include "test_common.h"

static uint8_t workbuf[VB2_FIRMWARE_WORKBUF_RECOMMENDED_SIZE]
	__attribute__ ((aligned (VB2_WORKBUF_ALIGN)));
static struct vb2_context *ctx;
static struct vb2_gbb_header gbb;
static struct vb2_shared_data *sd;
static struct vb2_secdata_fwmp *sec;

static void reset_common_data(void)
{
	memset(workbuf, 0xaa, sizeof(workbuf));
	TEST_SUCC(vb2api_init(workbuf, sizeof(workbuf), &ctx),
		  "vb2api_init failed");

	sd = vb2_get_sd(ctx);
	sd->status = VB2_SD_STATUS_SECDATA_FWMP_INIT;

	memset(&gbb, 0, sizeof(gbb));

	sec = (struct vb2_secdata_fwmp *)ctx->secdata_fwmp;
	sec->struct_size = VB2_SECDATA_FWMP_MIN_SIZE;
	sec->struct_version = VB2_SECDATA_FWMP_VERSION;
	sec->flags = 0;
	sec->crc8 = vb2_secdata_fwmp_crc(sec);
}

/* Mocked functions */

struct vb2_gbb_header *vb2_get_gbb(struct vb2_context *c)
{
	return &gbb;
}

static void check_init_test(void)
{
	uint8_t size;

	/* Check size constants */
	TEST_TRUE(sizeof(struct vb2_secdata_fwmp) >= VB2_SECDATA_FWMP_MIN_SIZE,
		  "Struct min size constant");
	TEST_TRUE(sizeof(struct vb2_secdata_fwmp) <= VB2_SECDATA_FWMP_MAX_SIZE,
		  "Struct max size constant");

	/* struct_size too large */
	reset_common_data();
	sec->struct_size = VB2_SECDATA_FWMP_MAX_SIZE + 1;
	sec->crc8 = vb2_secdata_fwmp_crc(sec);
	size = sec->struct_size;
	TEST_EQ(vb2api_secdata_fwmp_check(ctx, &size),
		VB2_ERROR_SECDATA_FWMP_SIZE, "Check struct_size too large");
	TEST_EQ(vb2_secdata_fwmp_init(ctx),
		VB2_ERROR_SECDATA_FWMP_SIZE, "Init struct_size too large");

	/* struct_size too small */
	reset_common_data();
	sec->struct_size = VB2_SECDATA_FWMP_MIN_SIZE - 1;
	sec->crc8 = vb2_secdata_fwmp_crc(sec);
	size = VB2_SECDATA_FWMP_MIN_SIZE;
	TEST_EQ(vb2api_secdata_fwmp_check(ctx, &size),
		VB2_ERROR_SECDATA_FWMP_SIZE, "Check struct_size too small");

	/* Need more data to reach minimum size */
	reset_common_data();
	sec->struct_size = VB2_SECDATA_FWMP_MIN_SIZE - 1;
	sec->crc8 = vb2_secdata_fwmp_crc(sec);
	size = 0;
	TEST_EQ(vb2api_secdata_fwmp_check(ctx, &size),
		VB2_ERROR_SECDATA_FWMP_INCOMPLETE, "Check more to reach MIN");
	TEST_EQ(vb2_secdata_fwmp_init(ctx),
		VB2_ERROR_SECDATA_FWMP_INCOMPLETE, "Init more to reach MIN");

	/* Need more data to reach full size */
	reset_common_data();
	sec->struct_size = VB2_SECDATA_FWMP_MIN_SIZE + 1;
	sec->crc8 = vb2_secdata_fwmp_crc(sec);
	size = VB2_SECDATA_FWMP_MIN_SIZE;
	TEST_EQ(vb2api_secdata_fwmp_check(ctx, &size),
		VB2_ERROR_SECDATA_FWMP_INCOMPLETE, "Check more for full size");

	/* Bad data is invalid */
	reset_common_data();
	memset(&ctx->secdata_fwmp, 0xa6, sizeof(ctx->secdata_fwmp));
	sec->struct_size = VB2_SECDATA_FWMP_MIN_SIZE;
	size = sec->struct_size;
	TEST_EQ(vb2api_secdata_fwmp_check(ctx, &size),
		VB2_ERROR_SECDATA_FWMP_CRC, "Check bad data CRC");
	TEST_EQ(vb2_secdata_fwmp_init(ctx),
		VB2_ERROR_SECDATA_FWMP_CRC, "Init bad data CRC");

	/* Bad CRC with corruption past minimum size */
	reset_common_data();
	sec->struct_size = VB2_SECDATA_FWMP_MIN_SIZE + 1;
	sec->crc8 = vb2_secdata_fwmp_crc(sec);
	size = sec->struct_size;
	*((uint8_t *)sec + sec->struct_size - 1) += 1;
	TEST_EQ(vb2api_secdata_fwmp_check(ctx, &size),
		VB2_ERROR_SECDATA_FWMP_CRC, "Check corruption CRC");
	TEST_EQ(vb2_secdata_fwmp_init(ctx),
		VB2_ERROR_SECDATA_FWMP_CRC, "Init corruption CRC");

	/* Zeroed data is invalid */
	reset_common_data();
	memset(&ctx->secdata_fwmp, 0, sizeof(ctx->secdata_fwmp));
	sec->struct_size = VB2_SECDATA_FWMP_MIN_SIZE;
	size = sec->struct_size;
	TEST_EQ(vb2api_secdata_fwmp_check(ctx, &size),
		VB2_ERROR_SECDATA_FWMP_VERSION, "Check zeroed data CRC");
	TEST_EQ(vb2_secdata_fwmp_init(ctx),
		VB2_ERROR_SECDATA_FWMP_VERSION, "Init zeroed data CRC");

	/* Major version too high */
	reset_common_data();
	sec->struct_version = ((VB2_SECDATA_FWMP_VERSION >> 4) + 1) << 4;
	sec->crc8 = vb2_secdata_fwmp_crc(sec);
	TEST_EQ(vb2api_secdata_fwmp_check(ctx, &size),
		VB2_ERROR_SECDATA_FWMP_VERSION, "Check major too high");
	TEST_EQ(vb2_secdata_fwmp_init(ctx),
		VB2_ERROR_SECDATA_FWMP_VERSION, "Init major too high");

	/* Major version too low */
	reset_common_data();
	sec->struct_version = ((VB2_SECDATA_FWMP_VERSION >> 4) - 1) << 4;
	sec->crc8 = vb2_secdata_fwmp_crc(sec);
	TEST_EQ(vb2api_secdata_fwmp_check(ctx, &size),
		VB2_ERROR_SECDATA_FWMP_VERSION, "Check major too low");
	TEST_EQ(vb2_secdata_fwmp_init(ctx),
		VB2_ERROR_SECDATA_FWMP_VERSION, "Init major too low");

	/* Minor version difference okay */
	reset_common_data();
	sec->struct_version += 1;
	sec->crc8 = vb2_secdata_fwmp_crc(sec);
	TEST_SUCC(vb2api_secdata_fwmp_check(ctx, &size), "Check minor okay");
	TEST_SUCC(vb2_secdata_fwmp_init(ctx), "Init minor okay");

	/* Good FWMP data at minimum size */
	reset_common_data();
	TEST_SUCC(vb2api_secdata_fwmp_check(ctx, &size), "Check good (min)");
	TEST_SUCC(vb2_secdata_fwmp_init(ctx), "Init good (min)");
	TEST_NEQ(sd->status & VB2_SD_STATUS_SECDATA_FWMP_INIT, 0,
		 "Init flag set");

	/* Good FWMP data at minimum + N size */
	reset_common_data();
	sec->struct_size = VB2_SECDATA_FWMP_MIN_SIZE + 1;
	sec->crc8 = vb2_secdata_fwmp_crc(sec);
	size = sec->struct_size;
	TEST_SUCC(vb2api_secdata_fwmp_check(ctx, &size), "Check good (min+N)");
	TEST_SUCC(vb2_secdata_fwmp_init(ctx), "Init good (min+N)");
	TEST_NEQ(sd->status & VB2_SD_STATUS_SECDATA_FWMP_INIT, 0,
		 "Init flag set");

	/* Skip data check when NO_SECDATA_FWMP set */
	reset_common_data();
	memset(&ctx->secdata_fwmp, 0xa6, sizeof(ctx->secdata_fwmp));
	ctx->flags |= VB2_CONTEXT_NO_SECDATA_FWMP;
	TEST_EQ(vb2_secdata_fwmp_init(ctx), 0,
		"Init skip data check when NO_SECDATA_FWMP set");
	TEST_NEQ(sd->status & VB2_SD_STATUS_SECDATA_FWMP_INIT, 0,
		 "Init flag set");
}

static void get_flag_test(void)
{
	/* Successfully returns value */
	reset_common_data();
	sec->flags |= 1;
	TEST_EQ(vb2_secdata_fwmp_get_flag(ctx, 1), 1,
		"Successfully returns flag value");

	/* NO_SECDATA_FWMP */
	reset_common_data();
	sec->flags |= 1;
	ctx->flags |= VB2_CONTEXT_NO_SECDATA_FWMP;
	TEST_EQ(vb2_secdata_fwmp_get_flag(ctx, 1), 0,
		"NO_SECDATA_FWMP forces default flag value");

	/* FWMP hasn't been initialized (recovery mode) */
	reset_common_data();
	sd->status &= ~VB2_SD_STATUS_SECDATA_FWMP_INIT;
	ctx->flags |= VB2_CONTEXT_RECOVERY_MODE;
	TEST_EQ(vb2_secdata_fwmp_get_flag(ctx, 0), 0,
		"non-init in recovery mode forces default flag value");

	/* FWMP hasn't been initialized (normal mode) */
	reset_common_data();
	sd->status &= ~VB2_SD_STATUS_SECDATA_FWMP_INIT;
	TEST_ABORT(vb2_secdata_fwmp_get_flag(ctx, 0),
		   "non-init in normal mode triggers abort");
}

static void get_dev_key_hash_test(void)
{
	/* CONTEXT_NO_SECDATA_FWMP */
	reset_common_data();
	ctx->flags |= VB2_CONTEXT_NO_SECDATA_FWMP;
	TEST_TRUE(vb2_secdata_fwmp_get_dev_key_hash(ctx) == NULL,
		  "NO_SECDATA_FWMP forces NULL pointer");

	/* FWMP hasn't been initialized (recovery mode) */
	reset_common_data();
	sd->status &= ~VB2_SD_STATUS_SECDATA_FWMP_INIT;
	ctx->flags |= VB2_CONTEXT_RECOVERY_MODE;
	TEST_TRUE(vb2_secdata_fwmp_get_dev_key_hash(ctx) == NULL,
		  "non-init in recovery mode forces NULL pointer");

	/* FWMP hasn't been initialized (normal mode) */
	reset_common_data();
	sd->status &= ~VB2_SD_STATUS_SECDATA_FWMP_INIT;
	TEST_ABORT(vb2_secdata_fwmp_get_dev_key_hash(ctx),
		   "non-init in normal mode triggers abort");

	/* Success case */
	reset_common_data();
	TEST_TRUE(vb2_secdata_fwmp_get_dev_key_hash(ctx) ==
		  sec->dev_key_hash, "proper dev_key_hash pointer returned");
}

int main(int argc, char* argv[])
{
	check_init_test();
	get_flag_test();
	get_dev_key_hash_test();

	return gTestSuccess ? 0 : 255;
}

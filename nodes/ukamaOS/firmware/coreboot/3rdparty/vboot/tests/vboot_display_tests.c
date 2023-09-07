/* Copyright (c) 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Tests for firmware display library.
 */

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "2common.h"
#include "2misc.h"
#include "2nvstorage.h"
#include "2struct.h"
#include "2sysincludes.h"
#include "host_common.h"
#include "test_common.h"
#include "vboot_common.h"
#include "vboot_display.h"
#include "vboot_kernel.h"

/* Mock data */
static uint8_t shared_data[VB_SHARED_DATA_MIN_SIZE];
static VbSharedDataHeader *shared = (VbSharedDataHeader *)shared_data;
static char debug_info[4096];
static struct vb2_context *ctx;
static struct vb2_shared_data *sd;
static uint8_t workbuf[VB2_KERNEL_WORKBUF_RECOMMENDED_SIZE];
static uint32_t mock_localization_count;
static uint32_t mock_altfw_mask;

/* Reset mock data (for use before each test) */
static void ResetMocks(void)
{
	mock_localization_count = 3;
	mock_altfw_mask = 3 << 1;	/* This mask selects 1 and 2 */

	TEST_SUCC(vb2api_init(workbuf, sizeof(workbuf), &ctx),
		  "vb2api_init failed");
	vb2_nv_init(ctx);

	sd = vb2_get_sd(ctx);
	sd->vbsd = shared;

	memset(&shared_data, 0, sizeof(shared_data));
	VbSharedDataInit(shared, sizeof(shared_data));

	*debug_info = 0;
}

/* Mocks */
vb2_error_t VbExGetLocalizationCount(uint32_t *count) {

	if (mock_localization_count == 0xffffffff)
		return VB2_ERROR_UNKNOWN;

	*count = mock_localization_count;
	return VB2_SUCCESS;
}

uint32_t VbExGetAltFwIdxMask() {
	return mock_altfw_mask;
}

vb2_error_t VbExDisplayDebugInfo(const char *info_str, int full_info)
{
	strncpy(debug_info, info_str, sizeof(debug_info));
	debug_info[sizeof(debug_info) - 1] = '\0';
	return VB2_SUCCESS;
}

/* Test displaying debug info */
static void DebugInfoTest(void)
{
	int i;

	/* Recovery string should be non-null for any code */
	for (i = 0; i < 0x100; i++)
		TEST_PTR_NEQ(RecoveryReasonString(i), NULL, "Non-null reason");

	/* Display debug info */
	ResetMocks();
	TEST_SUCC(VbDisplayDebugInfo(ctx),
		  "Display debug info");
	TEST_NEQ(*debug_info, '\0', "  Some debug info was displayed");
}

/* Test display key checking */
static void DisplayKeyTest(void)
{
	ResetMocks();
	VbCheckDisplayKey(ctx, 'q', NULL);
	TEST_EQ(*debug_info, '\0', "DisplayKey q = does nothing");

	ResetMocks();
	VbCheckDisplayKey(ctx, '\t', NULL);
	TEST_NEQ(*debug_info, '\0', "DisplayKey tab = display");

	/* Toggle localization */
	ResetMocks();
	vb2_nv_set(ctx, VB2_NV_LOCALIZATION_INDEX, 0);
	VbCheckDisplayKey(ctx, VB_KEY_DOWN, NULL);
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_LOCALIZATION_INDEX), 2,
		"DisplayKey up");
	VbCheckDisplayKey(ctx, VB_KEY_LEFT, NULL);
	vb2_nv_get(ctx, VB2_NV_LOCALIZATION_INDEX);
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_LOCALIZATION_INDEX), 1,
		"DisplayKey left");
	VbCheckDisplayKey(ctx, VB_KEY_RIGHT, NULL);
	vb2_nv_get(ctx, VB2_NV_LOCALIZATION_INDEX);
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_LOCALIZATION_INDEX), 2,
		"DisplayKey right");
	VbCheckDisplayKey(ctx, VB_KEY_UP, NULL);
	vb2_nv_get(ctx, VB2_NV_LOCALIZATION_INDEX);
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_LOCALIZATION_INDEX), 0,
		"DisplayKey up");

	/* Reset localization if localization count is invalid */
	ResetMocks();
	vb2_nv_set(ctx, VB2_NV_LOCALIZATION_INDEX, 1);
	mock_localization_count = 0xffffffff;
	VbCheckDisplayKey(ctx, VB_KEY_UP, NULL);
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_LOCALIZATION_INDEX), 0,
		"DisplayKey invalid");
}

int main(void)
{
	DebugInfoTest();
	DisplayKeyTest();

	return gTestSuccess ? 0 : 255;
}

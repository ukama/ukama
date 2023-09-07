/* Copyright (c) 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Tests for vboot_api_kernel, part 4 - select and load kernel
 */

#include "2api.h"
#include "2common.h"
#include "2ec_sync.h"
#include "2misc.h"
#include "2nvstorage.h"
#include "2secdata.h"
#include "2sysincludes.h"
#include "host_common.h"
#include "load_kernel_fw.h"
#include "secdata_tpm.h"
#include "test_common.h"
#include "tlcl.h"
#include "tss_constants.h"
#include "vboot_audio.h"
#include "vboot_common.h"
#include "vboot_kernel.h"
#include "vboot_struct.h"

/* Mock data */
static uint8_t workbuf[VB2_KERNEL_WORKBUF_RECOMMENDED_SIZE];
static struct vb2_context *ctx;
static struct vb2_context ctx_nvram_backend;
static struct vb2_shared_data *sd;
static VbSelectAndLoadKernelParams kparams;
static uint8_t shared_data[VB_SHARED_DATA_MIN_SIZE];
static VbSharedDataHeader *shared = (VbSharedDataHeader *)shared_data;
static struct vb2_gbb_header gbb;

static uint32_t rkr_version;
static uint32_t new_version;
static struct RollbackSpaceFwmp rfr_fwmp;
static int rkr_retval, rkw_retval, rkl_retval, rfr_retval;
static vb2_error_t vbboot_retval;

static uint32_t mock_switches[8];
static uint32_t mock_switches_count;
static int mock_switches_are_stuck;

/* Reset mock data (for use before each test) */
static void ResetMocks(void)
{
	memset(&kparams, 0, sizeof(kparams));

	memset(&gbb, 0, sizeof(gbb));
	gbb.major_version = VB2_GBB_MAJOR_VER;
	gbb.minor_version = VB2_GBB_MINOR_VER;
	gbb.flags = 0;

	TEST_SUCC(vb2api_init(workbuf, sizeof(workbuf), &ctx),
		  "vb2api_init failed");
	sd = vb2_get_sd(ctx);
	sd->flags |= VB2_SD_FLAG_DISPLAY_AVAILABLE;

	/*
	 * ctx_nvram_backend is only used as an NVRAM backend (see
	 * VbExNvStorageRead and VbExNvStorageWrite), and with
	 * vb2_set_nvdata and nv2_get_nvdata to manually read and tweak
	 * contents.  No other initialization is needed.
	 */
	memset(&ctx_nvram_backend, 0, sizeof(ctx_nvram_backend));
	vb2_nv_init(&ctx_nvram_backend);
	vb2_nv_set(&ctx_nvram_backend, VB2_NV_KERNEL_MAX_ROLLFORWARD,
		   0xffffffff);

	memset(&shared_data, 0, sizeof(shared_data));
	VbSharedDataInit(shared, sizeof(shared_data));

	memset(&rfr_fwmp, 0, sizeof(rfr_fwmp));
	rfr_retval = TPM_SUCCESS;

	rkr_version = new_version = 0x10002;
	rkr_retval = rkw_retval = rkl_retval = VB2_SUCCESS;
	vbboot_retval = VB2_SUCCESS;

	memset(mock_switches, 0, sizeof(mock_switches));
	mock_switches_count = 0;
	mock_switches_are_stuck = 0;
}

/* Mock functions */

vb2_error_t VbExNvStorageRead(uint8_t *buf)
{
	memcpy(buf, ctx_nvram_backend.nvdata,
	       vb2_nv_get_size(&ctx_nvram_backend));
	return VB2_SUCCESS;
}

vb2_error_t VbExNvStorageWrite(const uint8_t *buf)
{
	memcpy(ctx_nvram_backend.nvdata, buf,
	       vb2_nv_get_size(&ctx_nvram_backend));
	return VB2_SUCCESS;
}

uint32_t RollbackKernelRead(uint32_t *version)
{
	*version = rkr_version;
	return rkr_retval;
}

uint32_t RollbackKernelWrite(uint32_t version)
{
	rkr_version = version;
	return rkw_retval;
}

uint32_t RollbackKernelLock(int recovery_mode)
{
	return rkl_retval;
}

uint32_t RollbackFwmpRead(struct RollbackSpaceFwmp *fwmp)
{
	memcpy(fwmp, &rfr_fwmp, sizeof(*fwmp));
	return rfr_retval;
}

vb2_error_t VbTryLoadKernel(struct vb2_context *c, uint32_t get_info_flags)
{
	shared->kernel_version_tpm = new_version;

	if (vbboot_retval == -1)
		return VB2_ERROR_MOCK;

	return vbboot_retval;
}

vb2_error_t VbBootDeveloper(struct vb2_context *c)
{
	shared->kernel_version_tpm = new_version;

	if (vbboot_retval == -2)
		return VB2_ERROR_MOCK;

	return vbboot_retval;
}

vb2_error_t VbBootRecovery(struct vb2_context *c)
{
	shared->kernel_version_tpm = new_version;

	if (vbboot_retval == -3)
		return VB2_ERROR_MOCK;

	return vbboot_retval;
}

vb2_error_t VbBootDiagnostic(struct vb2_context *c)
{
	if (vbboot_retval == -4)
		return VB2_ERROR_MOCK;

	return vbboot_retval;
}

static void test_slk(vb2_error_t retval, int recovery_reason, const char *desc)
{
	TEST_EQ(VbSelectAndLoadKernel(ctx, shared, &kparams), retval, desc);
	TEST_EQ(vb2_nv_get(&ctx_nvram_backend, VB2_NV_RECOVERY_REQUEST),
		recovery_reason, "  recovery reason");
}

uint32_t VbExGetSwitches(uint32_t request_mask)
{
	if (mock_switches_are_stuck)
		return mock_switches[0] & request_mask;
	if (mock_switches_count < ARRAY_SIZE(mock_switches))
		return mock_switches[mock_switches_count++] & request_mask;
	else
		return 0;
}

vb2_error_t vb2ex_tpm_set_mode(enum vb2_tpm_mode mode_val)
{
	return VB2_SUCCESS;
}

/* Tests */

static void VbSlkTest(void)
{
	ResetMocks();
	test_slk(0, 0, "Normal");
	TEST_EQ(rkr_version, 0x10002, "  version");

	/*
	 * If shared->flags doesn't ask for software sync, we won't notice
	 * that error.
	 */
	ResetMocks();
	test_slk(0, 0, "EC sync not done");

	/* Same if shared->flags asks for sync, but it's overridden by GBB */
	ResetMocks();
	shared->flags |= VBSD_EC_SOFTWARE_SYNC;
	gbb.flags |= VB2_GBB_FLAG_DISABLE_EC_SOFTWARE_SYNC;
	test_slk(0, 0, "EC sync disabled by GBB");

	/* Rollback kernel version */
	ResetMocks();
	rkr_retval = 123;
	test_slk(VBERROR_TPM_READ_KERNEL,
		 VB2_RECOVERY_RW_TPM_R_ERROR, "Read kernel rollback");

	ResetMocks();
	new_version = 0x20003;
	test_slk(0, 0, "Roll forward");
	TEST_EQ(rkr_version, 0x20003, "  version");

	ResetMocks();
	vb2_nv_set(&ctx_nvram_backend, VB2_NV_FW_RESULT, VB2_FW_RESULT_TRYING);
	new_version = 0x20003;
	test_slk(0, 0, "Don't roll forward kernel when trying new FW");
	TEST_EQ(rkr_version, 0x10002, "  version");

	ResetMocks();
	vb2_nv_set(&ctx_nvram_backend, VB2_NV_KERNEL_MAX_ROLLFORWARD, 0x30005);
	new_version = 0x40006;
	test_slk(0, 0, "Limit max roll forward");
	TEST_EQ(rkr_version, 0x30005, "  version");

	ResetMocks();
	vb2_nv_set(&ctx_nvram_backend, VB2_NV_KERNEL_MAX_ROLLFORWARD, 0x10001);
	new_version = 0x40006;
	test_slk(0, 0, "Max roll forward can't rollback");
	TEST_EQ(rkr_version, 0x10002, "  version");

	ResetMocks();
	new_version = 0x20003;
	rkw_retval = 123;
	test_slk(VBERROR_TPM_WRITE_KERNEL,
		 VB2_RECOVERY_RW_TPM_W_ERROR, "Write kernel rollback");

	ResetMocks();
	rkl_retval = 123;
	test_slk(VBERROR_TPM_LOCK_KERNEL,
		 VB2_RECOVERY_RW_TPM_L_ERROR, "Lock kernel rollback");

	/* Boot normal */
	ResetMocks();
	vbboot_retval = -1;
	test_slk(VB2_ERROR_MOCK, 0, "Normal boot bad");

	/* Check that NV_DIAG_REQUEST triggers diagnostic UI */
	if (DIAGNOSTIC_UI) {
		ResetMocks();
		mock_switches[1] = VB_SWITCH_FLAG_PHYS_PRESENCE_PRESSED;
		vb2_nv_set(&ctx_nvram_backend, VB2_NV_DIAG_REQUEST, 1);
		vbboot_retval = -4;
		test_slk(VB2_ERROR_MOCK, 0,
			 "Normal boot with diag");
		TEST_EQ(vb2_nv_get(&ctx_nvram_backend, VB2_NV_DIAG_REQUEST),
			0, "  diag not requested");
	}

	/* Boot dev */
	ResetMocks();
	shared->flags |= VBSD_BOOT_DEV_SWITCH_ON;
	vbboot_retval = -2;
	test_slk(VB2_ERROR_MOCK, 0, "Dev boot bad");

	ResetMocks();
	shared->flags |= VBSD_BOOT_DEV_SWITCH_ON;
	new_version = 0x20003;
	test_slk(0, 0, "Dev doesn't roll forward");
	TEST_EQ(rkr_version, 0x10002, "  version");

	/* Boot recovery */
	ResetMocks();
	shared->recovery_reason = 123;
	vbboot_retval = -3;
	test_slk(VB2_ERROR_MOCK, 0, "Recovery boot bad");

	ResetMocks();
	shared->recovery_reason = 123;
	new_version = 0x20003;
	test_slk(0, 0, "Recovery doesn't roll forward");
	TEST_EQ(rkr_version, 0x10002, "  version");

	ResetMocks();
	shared->recovery_reason = 123;
	rkr_retval = rkw_retval = rkl_retval = VB2_ERROR_MOCK;
	test_slk(0, 0, "Recovery ignore TPM errors");

	ResetMocks();
	shared->recovery_reason = VB2_RECOVERY_TRAIN_AND_REBOOT;
	test_slk(VBERROR_REBOOT_REQUIRED, 0, "Recovery train and reboot");

	// todo: rkr/w/l fail ignored if recovery


}

int main(void)
{
	VbSlkTest();

	return gTestSuccess ? 0 : 255;
}

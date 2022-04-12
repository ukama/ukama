/* Copyright (c) 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * High-level firmware wrapper API - entry points for kernel selection
 */

#include "2common.h"
#include "2ec_sync.h"
#include "2misc.h"
#include "2nvstorage.h"
#include "2rsa.h"
#include "2secdata.h"
#include "2sysincludes.h"
#include "load_kernel_fw.h"
#include "secdata_tpm.h"
#include "utility.h"
#include "vb2_common.h"
#include "vboot_api.h"
#include "vboot_common.h"
#include "vboot_kernel.h"
#include "vboot_test.h"

/* Global variables */
static struct RollbackSpaceFwmp fwmp;
static LoadKernelParams lkp;

#ifdef CHROMEOS_ENVIRONMENT
/* Global variable accessors for unit tests */

struct RollbackSpaceFwmp *VbApiKernelGetFwmp(void)
{
	return &fwmp;
}

struct LoadKernelParams *VbApiKernelGetParams(void)
{
	return &lkp;
}

#endif

void vb2_nv_commit(struct vb2_context *ctx)
{
	/* Exit if nothing has changed */
	if (!(ctx->flags & VB2_CONTEXT_NVDATA_CHANGED))
		return;

	ctx->flags &= ~VB2_CONTEXT_NVDATA_CHANGED;
	VbExNvStorageWrite(ctx->nvdata);
}

uint32_t vb2_get_fwmp_flags(void)
{
	return fwmp.flags;
}

vb2_error_t VbTryLoadKernel(struct vb2_context *ctx, uint32_t get_info_flags)
{
	vb2_error_t rv = VBERROR_NO_DISK_FOUND;
	VbDiskInfo* disk_info = NULL;
	uint32_t disk_count = 0;
	uint32_t i;

	lkp.fwmp = &fwmp;
	lkp.disk_handle = NULL;

	/* Find disks */
	if (VB2_SUCCESS != VbExDiskGetInfo(&disk_info, &disk_count,
					   get_info_flags))
		disk_count = 0;

	/* Loop over disks */
	for (i = 0; i < disk_count; i++) {
		VB2_DEBUG("trying disk %d\n", (int)i);
		/*
		 * Sanity-check what we can. FWIW, VbTryLoadKernel() is always
		 * called with only a single bit set in get_info_flags.
		 *
		 * Ensure that we got a partition with only the flags we asked
		 * for.
		 */
		if (disk_info[i].bytes_per_lba < 512 ||
			(disk_info[i].bytes_per_lba &
				(disk_info[i].bytes_per_lba  - 1)) != 0 ||
					16 > disk_info[i].lba_count ||
					get_info_flags != (disk_info[i].flags &
					~VB_DISK_FLAG_EXTERNAL_GPT)) {
			VB2_DEBUG("  skipping: bytes_per_lba=%" PRIu64
				  " lba_count=%" PRIu64 " flags=%#x\n",
				  disk_info[i].bytes_per_lba,
				  disk_info[i].lba_count,
				  disk_info[i].flags);
			continue;
		}
		lkp.disk_handle = disk_info[i].handle;
		lkp.bytes_per_lba = disk_info[i].bytes_per_lba;
		lkp.gpt_lba_count = disk_info[i].lba_count;
		lkp.streaming_lba_count = disk_info[i].streaming_lba_count
						?: lkp.gpt_lba_count;
		lkp.boot_flags |= disk_info[i].flags & VB_DISK_FLAG_EXTERNAL_GPT
				? BOOT_FLAG_EXTERNAL_GPT : 0;

		vb2_error_t new_rv = LoadKernel(ctx, &lkp);
		VB2_DEBUG("LoadKernel() = %#x\n", new_rv);

		/* Stop now if we found a kernel. */
		if (VB2_SUCCESS == new_rv) {
			VbExDiskFreeInfo(disk_info, lkp.disk_handle);
			return VB2_SUCCESS;
		}

		/* Don't update error if we already have a more specific one. */
		if (VBERROR_INVALID_KERNEL_FOUND != rv)
			rv = new_rv;
	}

	/* If we drop out of the loop, we didn't find any usable kernel. */
	if (get_info_flags & VB_DISK_FLAG_FIXED) {
		switch (rv) {
		case VBERROR_INVALID_KERNEL_FOUND:
			vb2api_fail(ctx, VB2_RECOVERY_RW_INVALID_OS, rv);
			break;
		case VBERROR_NO_KERNEL_FOUND:
			vb2api_fail(ctx, VB2_RECOVERY_RW_NO_KERNEL, rv);
			break;
		case VBERROR_NO_DISK_FOUND:
			vb2api_fail(ctx, VB2_RECOVERY_RW_NO_DISK, rv);
			break;
		default:
			vb2api_fail(ctx, VB2_RECOVERY_LK_UNSPECIFIED, rv);
			break;
		}
	}

	/* If we didn't find any good kernels, don't return a disk handle. */
	VbExDiskFreeInfo(disk_info, NULL);

	return rv;
}

/**
 * Reset any NVRAM requests.
 *
 * @param ctx		Vboot context
 * @return 1 if a reboot is required, 0 otherwise.
 */
static int vb2_reset_nv_requests(struct vb2_context *ctx)
{
	int need_reboot = 0;

	if (vb2_nv_get(ctx, VB2_NV_DISPLAY_REQUEST)) {
		VB2_DEBUG("Unset display request (undo display init)\n");
		vb2_nv_set(ctx, VB2_NV_DISPLAY_REQUEST, 0);
		need_reboot = 1;
	}

	if (vb2_nv_get(ctx, VB2_NV_DIAG_REQUEST)) {
		VB2_DEBUG("Unset diagnostic request (undo display init)\n");
		vb2_nv_set(ctx, VB2_NV_DIAG_REQUEST, 0);
		need_reboot = 1;
	}

	return need_reboot;
}

vb2_error_t VbBootNormal(struct vb2_context *ctx)
{
	struct vb2_shared_data *sd = vb2_get_sd(ctx);
	VbSharedDataHeader *shared = sd->vbsd;
	uint32_t max_rollforward = vb2_nv_get(ctx,
					      VB2_NV_KERNEL_MAX_ROLLFORWARD);

	/* Boot from fixed disk only */
	VB2_DEBUG("Entering\n");

	if (vb2_reset_nv_requests(ctx)) {
		VB2_DEBUG("Normal mode: reboot to reset NVRAM requests\n");
		return VBERROR_REBOOT_REQUIRED;
	}

	vb2_error_t rv = VbTryLoadKernel(ctx, VB_DISK_FLAG_FIXED);

	VB2_DEBUG("Checking if TPM kernel version needs advancing\n");

	/*
	 * Special case for when we're trying a slot with new firmware.
	 * Firmware updates also usually change the kernel key, which means
	 * that the new firmware can only boot a new kernel, and the old
	 * firmware in the previous slot can only boot the previous kernel.
	 *
	 * Don't roll-forward the kernel version, because we don't yet know if
	 * the new kernel will successfully boot.
	 */
	if (vb2_nv_get(ctx, VB2_NV_FW_RESULT) == VB2_FW_RESULT_TRYING) {
		VB2_DEBUG("Trying new FW; skip kernel version roll-forward.\n");
		return rv;
	}

	/*
	 * Limit kernel version rollforward if needed.  Can't limit kernel
	 * version to less than the version currently in the TPM.  That is,
	 * we're limiting rollforward, not allowing rollback.
	 */
	if (max_rollforward < shared->kernel_version_tpm_start)
		max_rollforward = shared->kernel_version_tpm_start;

	if (shared->kernel_version_tpm > max_rollforward) {
		VB2_DEBUG("Limiting TPM kernel version roll-forward "
			  "to %#x < %#x\n",
			  max_rollforward, shared->kernel_version_tpm);

		shared->kernel_version_tpm = max_rollforward;
	}

	if (shared->kernel_version_tpm > shared->kernel_version_tpm_start) {
		uint32_t tpm_rv =
			RollbackKernelWrite(shared->kernel_version_tpm);
		if (tpm_rv) {
			VB2_DEBUG("Error writing kernel versions to TPM.\n");
			vb2api_fail(ctx, VB2_RECOVERY_RW_TPM_W_ERROR, tpm_rv);
			return VBERROR_TPM_WRITE_KERNEL;
		}
	}

	return rv;
}

static vb2_error_t vb2_kernel_setup(struct vb2_context *ctx,
				    VbSharedDataHeader *shared,
				    VbSelectAndLoadKernelParams *kparams)
{
	uint32_t tpm_rv;

	/* Translate vboot1 flags back to vboot2 */
	if (shared->recovery_reason)
		ctx->flags |= VB2_CONTEXT_RECOVERY_MODE;
	if (shared->flags & VBSD_BOOT_DEV_SWITCH_ON)
		ctx->flags |= VB2_CONTEXT_DEVELOPER_MODE;

	/*
	 * The following flags are set by depthcharge.
	 *
	 * TODO: Some of these are set at compile-time, so could be #defines
	 * instead of flags.  That would save on firmware image size because
	 * features that won't be used in an image could be compiled out.
	 */
	if (shared->flags & VBSD_EC_SOFTWARE_SYNC)
		ctx->flags |= VB2_CONTEXT_EC_SYNC_SUPPORTED;
	if (shared->flags & VBSD_EC_SLOW_UPDATE)
		ctx->flags |= VB2_CONTEXT_EC_SYNC_SLOW;
	if (shared->flags & VBSD_EC_EFS)
		ctx->flags |= VB2_CONTEXT_EC_EFS;
	if (shared->flags & VBSD_NVDATA_V2)
		ctx->flags |= VB2_CONTEXT_NVDATA_V2;

	VbExNvStorageRead(ctx->nvdata);
	vb2_nv_init(ctx);

	struct vb2_shared_data *sd = vb2_get_sd(ctx);
	struct vb2_gbb_header *gbb = vb2_get_gbb(ctx);
	sd->recovery_reason = shared->recovery_reason;

	/*
	 * Save a pointer to the old vboot1 shared data, since we haven't
	 * finished porting the library to use the new vb2 context and shared
	 * data.
	 *
	 * TODO: replace this with fields directly in vb2 shared data.
	 */
	sd->vbsd = shared;

	/*
	 * If we're in recovery mode just to do memory retraining, all we
	 * need to do is reboot.
	 */
	if (sd->recovery_reason == VB2_RECOVERY_TRAIN_AND_REBOOT) {
		VB2_DEBUG("Reboot after retraining in recovery.\n");
		return VBERROR_REBOOT_REQUIRED;
	}

	/* Fill in params for calls to LoadKernel() */
	memset(&lkp, 0, sizeof(lkp));
	lkp.kernel_buffer = kparams->kernel_buffer;
	lkp.kernel_buffer_size = kparams->kernel_buffer_size;

	/* Clear output params in case we fail */
	kparams->disk_handle = NULL;
	kparams->partition_number = 0;
	kparams->bootloader_address = 0;
	kparams->bootloader_size = 0;
	kparams->flags = 0;
	memset(kparams->partition_guid, 0, sizeof(kparams->partition_guid));

	/* Read kernel version from the TPM.  Ignore errors in recovery mode. */
	tpm_rv = RollbackKernelRead(&shared->kernel_version_tpm);
	if (tpm_rv) {
		VB2_DEBUG("Unable to get kernel versions from TPM\n");
		if (!(ctx->flags & VB2_CONTEXT_RECOVERY_MODE)) {
			vb2api_fail(ctx, VB2_RECOVERY_RW_TPM_R_ERROR, tpm_rv);
			return VBERROR_TPM_READ_KERNEL;
		}
	}

	shared->kernel_version_tpm_start = shared->kernel_version_tpm;

	/* Read FWMP.  Ignore errors in recovery mode. */
	if (gbb->flags & VB2_GBB_FLAG_DISABLE_FWMP) {
		memset(&fwmp, 0, sizeof(fwmp));
		return VB2_SUCCESS;
	}

	tpm_rv = RollbackFwmpRead(&fwmp);
	if (tpm_rv) {
		VB2_DEBUG("Unable to get FWMP from TPM\n");
		if (!(ctx->flags & VB2_CONTEXT_RECOVERY_MODE)) {
			vb2api_fail(ctx, VB2_RECOVERY_RW_TPM_R_ERROR, tpm_rv);
			return VBERROR_TPM_READ_FWMP;
		}
	}

	return VB2_SUCCESS;
}

static vb2_error_t vb2_kernel_phase4(struct vb2_context *ctx,
				     VbSelectAndLoadKernelParams *kparams)
{
	struct vb2_shared_data *sd = vb2_get_sd(ctx);

	/* Save disk parameters */
	kparams->disk_handle = lkp.disk_handle;
	kparams->partition_number = lkp.partition_number;
	kparams->bootloader_address = lkp.bootloader_address;
	kparams->bootloader_size = lkp.bootloader_size;
	kparams->flags = lkp.flags;
	kparams->kernel_buffer = lkp.kernel_buffer;
	kparams->kernel_buffer_size = lkp.kernel_buffer_size;
	memcpy(kparams->partition_guid, lkp.partition_guid,
	       sizeof(kparams->partition_guid));

	/* Lock the kernel versions if not in recovery mode */
	if (!(ctx->flags & VB2_CONTEXT_RECOVERY_MODE)) {
		uint32_t tpm_rv = RollbackKernelLock(sd->recovery_reason);
		if (tpm_rv) {
			VB2_DEBUG("Error locking kernel versions.\n");
			vb2api_fail(ctx, VB2_RECOVERY_RW_TPM_L_ERROR, tpm_rv);
			return VBERROR_TPM_LOCK_KERNEL;
		}
	}

	return VB2_SUCCESS;
}

static void vb2_kernel_cleanup(struct vb2_context *ctx)
{
	vb2_nv_commit(ctx);
}

vb2_error_t VbSelectAndLoadKernel(struct vb2_context *ctx,
				  VbSharedDataHeader *shared,
				  VbSelectAndLoadKernelParams *kparams)
{
	vb2_error_t rv = vb2_kernel_setup(ctx, shared, kparams);
	if (rv)
		goto VbSelectAndLoadKernel_exit;

	VB2_DEBUG("GBB flags are %#x\n", vb2_get_gbb(ctx)->flags);

	/*
	 * Do EC software sync unless we're in recovery mode. This has UI but
	 * it's just a single non-interactive WAIT screen.
	 */
	if (!(ctx->flags & VB2_CONTEXT_RECOVERY_MODE)) {
		rv = ec_sync_all(ctx);
		if (rv)
			goto VbSelectAndLoadKernel_exit;
	}

	/* Select boot path */
	if (ctx->flags & VB2_CONTEXT_RECOVERY_MODE) {
		/* Recovery boot.  This has UI. */
		if (kparams->inflags & VB_SALK_INFLAGS_ENABLE_DETACHABLE_UI)
			rv = VbBootRecoveryMenu(ctx);
		else
			rv = VbBootRecovery(ctx);
	} else if (DIAGNOSTIC_UI && vb2_nv_get(ctx, VB2_NV_DIAG_REQUEST)) {
		vb2_nv_set(ctx, VB2_NV_DIAG_REQUEST, 0);

		/*
		 * Diagnostic boot. This has a UI but only power button
		 * is used for input so no detachable-specific UI is
		 * needed.  This mode is also 1-shot so it's placed
		 * before developer mode.
		 */
		rv = VbBootDiagnostic(ctx);
		/*
		 * The diagnostic menu should either boot a rom, or
		 * return either of reboot or shutdown.  The following
		 * check is a safety precaution.
		 */
		if (!rv)
			rv = VBERROR_REBOOT_REQUIRED;
	} else if (ctx->flags & VB2_CONTEXT_DEVELOPER_MODE) {
		if (kparams->inflags & VB_SALK_INFLAGS_VENDOR_DATA_SETTABLE)
			ctx->flags |= VB2_CONTEXT_VENDOR_DATA_SETTABLE;

		/* Developer boot.  This has UI. */
		if (kparams->inflags & VB_SALK_INFLAGS_ENABLE_DETACHABLE_UI)
			rv = VbBootDeveloperMenu(ctx);
		else
			rv = VbBootDeveloper(ctx);
	} else {
		/* Normal boot */
		rv = VbBootNormal(ctx);
	}

 VbSelectAndLoadKernel_exit:

	if (VB2_SUCCESS == rv)
		rv = vb2_kernel_phase4(ctx, kparams);

	vb2_kernel_cleanup(ctx);

	/* Pass through return value from boot path */
	VB2_DEBUG("Returning %#x\n", rv);
	return rv;
}

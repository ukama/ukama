/* Copyright (c) 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Display functions used in kernel selection.
 */

#include "2common.h"
#include "2misc.h"
#include "2nvstorage.h"
#include "2sha.h"
#include "2sysincludes.h"
#include "utility.h"
#include "vboot_api.h"
#include "vboot_common.h"
#include "vboot_display.h"
#include "vboot_kernel.h"

static uint32_t disp_current_screen = VB_SCREEN_BLANK;
static uint32_t disp_current_index = 0;
static uint32_t disp_disabled_idx_mask = 0;

__attribute__((weak))
vb2_error_t VbExGetLocalizationCount(uint32_t *count) {
	*count = 0;
	return VB2_ERROR_UNKNOWN;
}

__attribute__((weak))
vb2_error_t VbExGetAltFwIdxMask(void) {
	return 0;
}

vb2_error_t VbDisplayScreen(struct vb2_context *ctx, uint32_t screen, int force,
			    const VbScreenData *data)
{
	uint32_t locale;

	/* If requested screen is the same as the current one, we're done. */
	if (disp_current_screen == screen && !force)
		return VB2_SUCCESS;

	/* Keep track of the currently displayed screen */
	disp_current_screen = screen;

	/* Read the locale last saved */
	locale = vb2_nv_get(ctx, VB2_NV_LOCALIZATION_INDEX);

	return VbExDisplayScreen(screen, locale, data);
}

vb2_error_t VbDisplayMenu(struct vb2_context *ctx, uint32_t screen, int force,
			  uint32_t selected_index, uint32_t disabled_idx_mask)
{
	uint32_t locale;
	uint32_t redraw_base_screen = 0;

	/*
	 * If requested screen/selected_index is the same as the current one,
	 * we're done.
	 */
	if (disp_current_screen == screen &&
	    disp_current_index == selected_index &&
	    !force)
		return VB2_SUCCESS;

	/*
	 * If current screen is not the same, make sure we redraw the base
	 * screen as well to avoid having artifacts from the menu.
	 */
	if (disp_current_screen != screen || force)
		redraw_base_screen = 1;

	/*
	 * Keep track of the currently displayed screen and
	 * selected_index
	 */
	disp_current_screen = screen;
	disp_current_index = selected_index;
	disp_disabled_idx_mask = disabled_idx_mask;

	/* Read the locale last saved */
	locale = vb2_nv_get(ctx, VB2_NV_LOCALIZATION_INDEX);

	return VbExDisplayMenu(screen, locale, selected_index,
			       disabled_idx_mask, redraw_base_screen);
}

static void Uint8ToString(char *buf, uint8_t val)
{
	const char *trans = "0123456789abcdef";
	*buf++ = trans[val >> 4];
	*buf = trans[val & 0xF];
}

static void FillInSha1Sum(char *outbuf, struct vb2_packed_key *key)
{
	uint8_t *buf = ((uint8_t *)key) + key->key_offset;
	uint64_t buflen = key->key_size;
	uint8_t digest[VB2_SHA1_DIGEST_SIZE];
	int i;

	vb2_digest_buffer(buf, buflen, VB2_HASH_SHA1, digest, sizeof(digest));
	for (i = 0; i < sizeof(digest); i++) {
		Uint8ToString(outbuf, digest[i]);
		outbuf += 2;
	}
	*outbuf = '\0';
}

const char *RecoveryReasonString(uint8_t code)
{
	switch(code) {
	case VB2_RECOVERY_NOT_REQUESTED:
		return "Recovery not requested";
	case VB2_RECOVERY_LEGACY:
		return "Recovery requested from legacy utility";
	case VB2_RECOVERY_RO_MANUAL:
		return "recovery button pressed";
	case VB2_RECOVERY_RO_INVALID_RW:
		return "RW firmware failed signature check";
	case VB2_RECOVERY_RO_SHARED_DATA:
		return "Shared data error in read-only firmware";
	case VB2_RECOVERY_FW_KEYBLOCK:
		return "RW firmware unable to verify keyblock";
	case VB2_RECOVERY_FW_KEY_ROLLBACK:
		return "RW firmware key version rollback detected";
	case VB2_RECOVERY_FW_PREAMBLE:
		return "RW firmware unable to verify preamble";
	case VB2_RECOVERY_FW_ROLLBACK:
		return "RW firmware version rollback detected";
	case VB2_RECOVERY_FW_BODY:
		return "RW firmware unable to verify firmware body";
	case VB2_RECOVERY_RO_FIRMWARE:
		return "Firmware problem outside of verified boot";
	case VB2_RECOVERY_RO_TPM_REBOOT:
		return "TPM requires a system reboot (should be transient)";
	case VB2_RECOVERY_EC_SOFTWARE_SYNC:
		return "EC software sync error";
	case VB2_RECOVERY_EC_UNKNOWN_IMAGE:
		return "EC software sync unable to determine active EC image";
	case VB2_RECOVERY_EC_EXPECTED_IMAGE:
		return "EC software sync error "
			"obtaining expected EC image from BIOS";
	case VB2_RECOVERY_EC_UPDATE:
		return "EC software sync error updating EC";
	case VB2_RECOVERY_EC_JUMP_RW:
		return "EC software sync unable to jump to EC-RW";
	case VB2_RECOVERY_EC_PROTECT:
		return "EC software sync protection error";
	case VB2_RECOVERY_EC_EXPECTED_HASH:
		return "EC software sync error "
			"obtaining expected EC hash from BIOS";
	case VB2_RECOVERY_SECDATA_FIRMWARE_INIT:
		return "Firmware secure NVRAM (TPM) initialization error";
	case VB2_RECOVERY_GBB_HEADER:
		return "Error parsing GBB header";
	case VB2_RECOVERY_TPM_CLEAR_OWNER:
		return "Error trying to clear TPM owner";
	case VB2_RECOVERY_DEV_SWITCH:
		return "Error reading or updating developer switch";
	case VB2_RECOVERY_FW_SLOT:
		return "Error selecting RW firmware slot";
	case VB2_RECOVERY_AUX_FW_UPDATE:
		return "Error updating AUX firmware";
	case VB2_RECOVERY_RO_UNSPECIFIED:
		return "Unspecified/unknown error in RO firmware";
	case VB2_RECOVERY_RW_INVALID_OS:
		return "OS kernel or rootfs failed signature check";
	case VB2_RECOVERY_RW_SHARED_DATA:
		return "Shared data error in rewritable firmware";
	case VB2_RECOVERY_TPM_E_FAIL:
		return "TPM error that was not fixed by reboot";
	case VB2_RECOVERY_RO_TPM_S_ERROR:
		return "TPM setup error in read-only firmware";
	case VB2_RECOVERY_RO_TPM_W_ERROR:
		return "TPM write error in read-only firmware";
	case VB2_RECOVERY_RO_TPM_L_ERROR:
		return "TPM lock error in read-only firmware";
	case VB2_RECOVERY_RO_TPM_U_ERROR:
		return "TPM update error in read-only firmware";
	case VB2_RECOVERY_RW_TPM_R_ERROR:
		return "TPM read error in rewritable firmware";
	case VB2_RECOVERY_RW_TPM_W_ERROR:
		return "TPM write error in rewritable firmware";
	case VB2_RECOVERY_RW_TPM_L_ERROR:
		return "TPM lock error in rewritable firmware";
	case VB2_RECOVERY_EC_HASH_FAILED:
		return "EC software sync unable to get EC image hash";
	case VB2_RECOVERY_EC_HASH_SIZE:
		return "EC software sync invalid image hash size";
	case VB2_RECOVERY_LK_UNSPECIFIED:
		return "Unspecified error while trying to load kernel";
	case VB2_RECOVERY_RW_NO_DISK:
		return "No bootable storage device in system";
	case VB2_RECOVERY_RW_NO_KERNEL:
		return "No bootable kernel found on disk";
	case VB2_RECOVERY_SECDATA_KERNEL_INIT:
		return "Kernel secure NVRAM (TPM) initialization error";
	case VB2_RECOVERY_RO_TPM_REC_HASH_L_ERROR:
		return "Recovery hash space lock error in RO firmware";
	case VB2_RECOVERY_TPM_DISABLE_FAILED:
		return "Failed to disable TPM before running untrusted code";
	case VB2_RECOVERY_ALTFW_HASH_FAILED:
		return "Verification of alternative firmware payload failed";
	case VB2_RECOVERY_RW_UNSPECIFIED:
		return "Unspecified/unknown error in RW firmware";
	case VB2_RECOVERY_US_TEST:
		return "Recovery mode test from user-mode";
	case VB2_RECOVERY_TRAIN_AND_REBOOT:
		return "User-mode requested DRAM train and reboot";
	case VB2_RECOVERY_US_UNSPECIFIED:
		return "Unspecified/unknown error in user-mode";
	}
	return "Unknown or deprecated error code";
}

#define DEBUG_INFO_SIZE 512

vb2_error_t VbDisplayDebugInfo(struct vb2_context *ctx)
{
	struct vb2_shared_data *sd = vb2_get_sd(ctx);
	struct vb2_gbb_header *gbb = vb2_get_gbb(ctx);
	struct vb2_workbuf wb;
	VbSharedDataHeader *shared = sd->vbsd;
	char buf[DEBUG_INFO_SIZE] = "";
	char sha1sum[VB2_SHA1_DIGEST_SIZE * 2 + 1];
	uint32_t used = 0;
	vb2_error_t ret;
	uint32_t i;

	vb2_workbuf_from_ctx(ctx, &wb);

	/* Add hardware ID */
	{
		char hwid[VB2_GBB_HWID_MAX_SIZE];
		uint32_t size = sizeof(hwid);
		ret = vb2api_gbb_read_hwid(ctx, hwid, &size);
		if (ret)
			strcpy(hwid, "{INVALID}");
		used += StrnAppend(buf + used, "HWID: ",
				   DEBUG_INFO_SIZE - used);
		used += StrnAppend(buf + used, hwid, DEBUG_INFO_SIZE - used);
	}

	/* Add recovery reason and subcode */
	i = vb2_nv_get(ctx, VB2_NV_RECOVERY_SUBCODE);
	used += StrnAppend(buf + used,
			"\nrecovery_reason: 0x", DEBUG_INFO_SIZE - used);
	used += Uint64ToString(buf + used, DEBUG_INFO_SIZE - used,
			       sd->recovery_reason, 16, 2);
	used += StrnAppend(buf + used, " / 0x", DEBUG_INFO_SIZE - used);
	used += Uint64ToString(buf + used, DEBUG_INFO_SIZE - used, i, 16, 2);
	used += StrnAppend(buf + used, "  ", DEBUG_INFO_SIZE - used);
	used += StrnAppend(buf + used,
			RecoveryReasonString(sd->recovery_reason),
			DEBUG_INFO_SIZE - used);

	/* Add VbSharedDataHeader flags */
	used += StrnAppend(buf + used, "\nVbSD.flags: 0x", DEBUG_INFO_SIZE - used);
	used += Uint64ToString(buf + used, DEBUG_INFO_SIZE - used,
			       shared->flags, 16, 8);

	/* Add raw contents of VbNvStorage */
	used += StrnAppend(buf + used, "\nVbNv.raw:", DEBUG_INFO_SIZE - used);
	for (i = 0; i < vb2_nv_get_size(ctx); i++) {
		used += StrnAppend(buf + used, " ", DEBUG_INFO_SIZE - used);
		used += Uint64ToString(buf + used, DEBUG_INFO_SIZE - used,
				       ctx->nvdata[i], 16, 2);
	}

	/* Add dev_boot_usb flag */
	i = vb2_nv_get(ctx, VB2_NV_DEV_BOOT_USB);
	used += StrnAppend(buf + used, "\ndev_boot_usb: ", DEBUG_INFO_SIZE - used);
	used += Uint64ToString(buf + used, DEBUG_INFO_SIZE - used, i, 10, 0);

	/* Add dev_boot_legacy flag */
	i = vb2_nv_get(ctx, VB2_NV_DEV_BOOT_LEGACY);
	used += StrnAppend(buf + used,
			"\ndev_boot_legacy: ", DEBUG_INFO_SIZE - used);
	used += Uint64ToString(buf + used, DEBUG_INFO_SIZE - used, i, 10, 0);

	/* Add dev_default_boot flag */
	i = vb2_nv_get(ctx, VB2_NV_DEV_DEFAULT_BOOT);
	used += StrnAppend(buf + used,
			"\ndev_default_boot: ", DEBUG_INFO_SIZE - used);
	used += Uint64ToString(buf + used, DEBUG_INFO_SIZE - used, i, 10, 0);

	/* Add dev_boot_signed_only flag */
	i = vb2_nv_get(ctx, VB2_NV_DEV_BOOT_SIGNED_ONLY);
	used += StrnAppend(buf + used, "\ndev_boot_signed_only: ",
			DEBUG_INFO_SIZE - used);
	used += Uint64ToString(buf + used, DEBUG_INFO_SIZE - used, i, 10, 0);

	/* Add TPM versions */
	used += StrnAppend(buf + used,
			   "\nTPM: fwver=0x", DEBUG_INFO_SIZE - used);
	used += Uint64ToString(buf + used, DEBUG_INFO_SIZE - used,
			       sd->fw_version_secdata, 16, 8);
	used += StrnAppend(buf + used, " kernver=0x", DEBUG_INFO_SIZE - used);
	used += Uint64ToString(buf + used, DEBUG_INFO_SIZE - used,
			       sd->kernel_version_secdata, 16, 8);

	/* Add GBB flags */
	used += StrnAppend(buf + used,
			   "\ngbb.flags: 0x", DEBUG_INFO_SIZE - used);
	used += Uint64ToString(buf + used, DEBUG_INFO_SIZE - used,
			       gbb->flags, 16, 8);

	/* Add sha1sum for Root & Recovery keys */
	{
		struct vb2_packed_key *key;
		struct vb2_workbuf wblocal = wb;
		ret = vb2_gbb_read_root_key(ctx, &key, NULL, &wblocal);
		if (!ret) {
			FillInSha1Sum(sha1sum, key);
			used += StrnAppend(buf + used, "\ngbb.rootkey: ",
					   DEBUG_INFO_SIZE - used);
			used += StrnAppend(buf + used, sha1sum,
					   DEBUG_INFO_SIZE - used);
		}
	}

	{
		struct vb2_packed_key *key;
		struct vb2_workbuf wblocal = wb;
		ret = vb2_gbb_read_recovery_key(ctx, &key, NULL, &wblocal);
		if (!ret) {
			FillInSha1Sum(sha1sum, key);
			used += StrnAppend(buf + used, "\ngbb.recovery_key: ",
					   DEBUG_INFO_SIZE - used);
			used += StrnAppend(buf + used, sha1sum,
					   DEBUG_INFO_SIZE - used);
		}
	}

	/* If we're in dev-mode, show the kernel subkey that we expect, too. */
	if (0 == sd->recovery_reason) {
		FillInSha1Sum(sha1sum, &shared->kernel_subkey);
		used += StrnAppend(buf + used,
				"\nkernel_subkey: ", DEBUG_INFO_SIZE - used);
		used += StrnAppend(buf + used, sha1sum, DEBUG_INFO_SIZE - used);
	}

	/* Make sure we finish with a newline */
	used += StrnAppend(buf + used, "\n", DEBUG_INFO_SIZE - used);

	/* TODO: add more interesting data:
	 * - Information on current disks */

	buf[DEBUG_INFO_SIZE - 1] = '\0';
	VB2_DEBUG("[TAB] Debug Info:\n%s", buf);
	return VbExDisplayDebugInfo(buf, 1);
}

vb2_error_t VbCheckDisplayKey(struct vb2_context *ctx, uint32_t key,
			    const VbScreenData *data)
{
	uint32_t loc = 0;
	uint32_t count = 0;

	switch (key) {
	case '\t':
		/* Tab = display debug info */
		return VbDisplayDebugInfo(ctx);
	case VB_KEY_ESC:
		/* Force redraw current screen (to clear Tab debug output) */
		return VbDisplayScreen(ctx, disp_current_screen, 1, data);
	case VB_KEY_LEFT:
	case VB_KEY_RIGHT:
	case VB_KEY_UP:
	case VB_KEY_DOWN:
		/* Arrow keys = change localization */
		loc = vb2_nv_get(ctx, VB2_NV_LOCALIZATION_INDEX);
		if (VB2_SUCCESS != VbExGetLocalizationCount(&count))
			loc = 0;  /* No localization count (bad GBB?) */
		else if (VB_KEY_RIGHT == key || VB_KEY_UP == key)
			loc = (loc < count - 1 ? loc + 1 : 0);
		else
			loc = (loc > 0 ? loc - 1 : count - 1);
		VB2_DEBUG("VbCheckDisplayKey() - change localization to %d\n",
			  (int)loc);
		vb2_nv_set(ctx, VB2_NV_LOCALIZATION_INDEX, loc);
		vb2_nv_set(ctx, VB2_NV_BACKUP_NVRAM_REQUEST, 1);

		/*
		 * Non-manual recovery mode is meant to be left via three-finger
		 * salute (into manual recovery mode). Need to commit nvdata
		 * changes immediately.
		 */
		if ((ctx->flags & VB2_CONTEXT_RECOVERY_MODE) &&
		    !vb2_allow_recovery(ctx))
			vb2_nv_commit(ctx);

		/* Force redraw of current screen */
		return VbDisplayScreen(ctx, disp_current_screen, 1, data);
	}

	return VB2_SUCCESS;
}

/* Copyright (c) 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Stub implementations of firmware-provided API functions.
 */

#include <stdint.h>

#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/time.h>

#include "2common.h"
#include "vboot_api.h"
#include "vboot_test.h"

void VbExSleepMs(uint32_t msec)
{
}

vb2_error_t VbExBeep(uint32_t msec, uint32_t frequency)
{
	return VB2_SUCCESS;
}

vb2_error_t VbExDisplayScreen(uint32_t screen_type, uint32_t locale,
			      const VbScreenData *data)
{
	return VB2_SUCCESS;
}

vb2_error_t VbExDisplayMenu(uint32_t screen_type, uint32_t locale,
			    uint32_t selected_index, uint32_t disabled_idx_mask,
			    uint32_t redraw_base)
{
	return VB2_SUCCESS;
}

vb2_error_t VbExDisplayDebugInfo(const char *info_str, int full_info)
{
	return VB2_SUCCESS;
}

uint32_t VbExKeyboardRead(void)
{
	return 0;
}

uint32_t VbExKeyboardReadWithFlags(uint32_t *flags_ptr)
{
	return 0;
}

uint32_t VbExGetSwitches(uint32_t mask)
{
	return 0;
}

uint32_t VbExIsShutdownRequested(void)
{
	return 0;
}

int VbExTrustEC(int devidx)
{
	return 1;
}

vb2_error_t VbExEcRunningRW(int devidx, int *in_rw)
{
	*in_rw = 0;
	return VB2_SUCCESS;
}

vb2_error_t VbExEcJumpToRW(int devidx)
{
	return VB2_SUCCESS;
}

vb2_error_t VbExEcDisableJump(int devidx)
{
	return VB2_SUCCESS;
}

vb2_error_t VbExEcHashImage(int devidx, enum VbSelectFirmware_t select,
			    const uint8_t **hash, int *hash_size)
{
	static const uint8_t fake_hash[32] = {1, 2, 3, 4};

	*hash = fake_hash;
	*hash_size = sizeof(fake_hash);
	return VB2_SUCCESS;
}

vb2_error_t VbExEcGetExpectedImage(int devidx, enum VbSelectFirmware_t select,
				   const uint8_t **image, int *image_size)
{
	static uint8_t fake_image[64] = {5, 6, 7, 8};
	*image = fake_image;
	*image_size = sizeof(fake_image);
	return VB2_SUCCESS;
}

vb2_error_t VbExEcGetExpectedImageHash(int devidx,
				       enum VbSelectFirmware_t select,
				       const uint8_t **hash, int *hash_size)
{
	static const uint8_t fake_hash[32] = {1, 2, 3, 4};

	*hash = fake_hash;
	*hash_size = sizeof(fake_hash);
	return VB2_SUCCESS;
}

vb2_error_t VbExEcUpdateImage(int devidx, enum VbSelectFirmware_t select,
			      const uint8_t *image, int image_size)
{
	return VB2_SUCCESS;
}

vb2_error_t VbExEcProtect(int devidx, enum VbSelectFirmware_t select)
{
	return VB2_SUCCESS;
}

vb2_error_t VbExEcVbootDone(int in_recovery)
{
	return VB2_SUCCESS;
}

vb2_error_t VbExEcBatteryCutOff(void)
{
	return VB2_SUCCESS;
}

vb2_error_t VbExCheckAuxFw(VbAuxFwUpdateSeverity_t *severity)
{
	*severity = VB_AUX_FW_NO_UPDATE;
	return VB2_SUCCESS;
}

vb2_error_t VbExUpdateAuxFw(void)
{
	return VB2_SUCCESS;
}

vb2_error_t VbExLegacy(enum VbAltFwIndex_t altfw_num)
{
	return 1;
}

uint8_t VbExOverrideGptEntryPriority(const GptEntry *e)
{
	return 0;
}

vb2_error_t VbExSetVendorData(const char *vendor_data_value)
{
	return 0;
}

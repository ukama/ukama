/* Copyright (c) 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Display functions used in kernel selection.
 */

#ifndef VBOOT_REFERENCE_VBOOT_DISPLAY_H_
#define VBOOT_REFERENCE_VBOOT_DISPLAY_H_

vb2_error_t VbDisplayScreen(struct vb2_context *ctx, uint32_t screen, int force,
			  const VbScreenData *data);
vb2_error_t VbDisplayMenu(struct vb2_context *ctx,
			uint32_t screen, int force, uint32_t selected_index,
			uint32_t disabled_idx_mask);
vb2_error_t VbDisplayDebugInfo(struct vb2_context *ctx);
vb2_error_t VbCheckDisplayKey(struct vb2_context *ctx, uint32_t key,
			    const VbScreenData *data);

/**
 * Return a description of the recovery reason code.
 */
const char *RecoveryReasonString(uint8_t code);

#endif  /* VBOOT_REFERENCE_VBOOT_DISPLAY_H_ */

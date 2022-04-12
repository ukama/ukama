/* Copyright 2019 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 * 
 * This header is for APIs that are only used by test code.
 */

#ifndef VBOOT_REFERENCE_TEST_API_H_
#define VBOOT_REFERENCE_TEST_API_H_

/****************************************************************************
 * 2rsa.c
 *
 * Internal functions from 2rsa.c that have error conditions we can't trigger
 * from the public APIs.  These include checks for bad algorithms where the
 * next call level up already checks for bad algorithms, etc.
 *
 * These functions aren't in 2rsa.h because they're not part of the public
 * APIs.
 */
struct vb2_public_key;
int vb2_mont_ge(const struct vb2_public_key *key, uint32_t *a);
vb2_error_t vb2_check_padding(const uint8_t *sig,
			      const struct vb2_public_key *key);

/****************************************************************************
 * vboot_api_kernel.c */

struct RollbackSpaceFwmp;
struct RollbackSpaceFwmp *VbApiKernelGetFwmp(void);

struct LoadKernelParams;
struct LoadKernelParams *VbApiKernelGetParams(void);

#endif  /* VBOOT_REFERENCE_TEST_API_H_ */

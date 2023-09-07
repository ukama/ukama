/* Copyright (c) 2012 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Functions for querying, manipulating and locking secure data spaces
 * stored in the TPM NVRAM.
 */

#include "secdata_tpm.h"
#include "tss_constants.h"
#include "utility.h"

vb2_error_t SetVirtualDevMode(int val)
{
	return VB2_SUCCESS;
}

uint32_t RollbackKernelRead(uint32_t *version)
{
	*version = 0;
	return TPM_SUCCESS;
}

uint32_t RollbackKernelWrite(uint32_t version)
{
	return TPM_SUCCESS;
}

uint32_t RollbackKernelLock(int recovery_mode)
{
	return TPM_SUCCESS;
}

uint32_t RollbackFwmpRead(struct RollbackSpaceFwmp *fwmp)
{
	memset(fwmp, 0, sizeof(*fwmp));
	return TPM_SUCCESS;
}

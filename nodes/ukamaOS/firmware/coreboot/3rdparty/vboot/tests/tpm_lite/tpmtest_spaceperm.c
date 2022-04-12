/* Copyright (c) 2011 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/* Test of space permissions retrieval.
 */

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#include "host_common.h"
#include "test_common.h"
#include "tlcl.h"
#include "tlcl_tests.h"

#define PERMPPGL (TPM_NV_PER_PPWRITE | TPM_NV_PER_GLOBALLOCK)
#define PERMPP TPM_NV_PER_PPWRITE

int main(int argc, char** argv) {
	uint32_t perm;

	TlclLibInit();
	TPM_CHECK(TlclStartupIfNeeded());
	TPM_CHECK(TlclContinueSelfTest());
	TPM_CHECK(TlclAssertPhysicalPresence());

	TPM_CHECK(TlclGetPermissions(INDEX0, &perm));
	TEST_NEQ(perm & PERMPPGL, 0, "INDEX0: PERMPPGL is not set");

	TPM_CHECK(TlclGetPermissions(INDEX1, &perm));
	TEST_NEQ(perm & PERMPP, 0, "INDEX1: PERMPP is not set");

	printf("TEST SUCCEEDED\n");
	exit(0);
}

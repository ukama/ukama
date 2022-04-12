/* Copyright (c) 2011 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/* Testing: ownership testing code, ForceClear, and nvram write limit.
 */

#include <stdio.h>

#include "2common.h"
#include "host_common.h"
#include "tlcl.h"
#include "tlcl_tests.h"


int main(int argc, char** argv) {
	uint8_t disable, deactivated;

	TlclLibInit();
	TPM_CHECK(TlclStartupIfNeeded());
	TPM_CHECK(TlclSelfTestFull());
	TPM_CHECK(TlclAssertPhysicalPresence());
	TPM_CHECK(TlclGetFlags(&disable, &deactivated, NULL));
	printf("disable is %d, deactivated is %d\n", disable, deactivated);
	TPM_CHECK(TlclSetEnable());
	TPM_CHECK(TlclSetDeactivated(0));
	TPM_CHECK(TlclGetFlags(&disable, &deactivated, NULL));
	printf("disable is %d, deactivated is %d\n", disable, deactivated);
	if (disable == 1 || deactivated == 1) {
		VB2_DEBUG("failed to enable or activate");
		return 1;
	}
	printf("TEST SUCCEEDED\n");
	return 0;
}

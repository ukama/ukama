/* Copyright (c) 2011 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Common functions used by tests.
 */

#include <stdint.h>
#include <stdio.h>
#include <string.h>

#include "2common.h"
#include "test_common.h"

/* Global test success flag. */
int gTestSuccess = 1;
int gTestAbortArmed = 0;
jmp_buf gTestJmpEnv;

int test_eq(int result, int expected,
	    const char *preamble, const char *desc, const char *comment)
{
	if (result == expected) {
		fprintf(stderr, "%s: %s ... " COL_GREEN "PASSED\n" COL_STOP,
			preamble, comment ? comment : desc);
		return 1;
	} else {
		fprintf(stderr, "%s: %s ... " COL_RED "FAILED\n" COL_STOP,
			preamble, comment ? comment : desc);
		fprintf(stderr, "	Expected: %#x (%d), got: %#x (%d)\n",
			expected, expected, result, result);
		gTestSuccess = 0;
		return 0;
	}
}

int test_neq(int result, int not_expected,
	     const char *preamble, const char *desc, const char *comment)
{
	if (result != not_expected) {
		fprintf(stderr, "%s: %s, %s ... " COL_GREEN "PASSED\n" COL_STOP,
			preamble, desc, comment);
		return 1;
	} else {
		fprintf(stderr, "%s: %s, %s ... " COL_RED "FAILED\n" COL_STOP,
			preamble, desc, comment);
		fprintf(stderr, "	Didn't expect %#x (%d), but got it.\n",
			not_expected, not_expected);
		gTestSuccess = 0;
		return 0;
	}
}

int test_ptr_eq(const void* result, const void* expected,
		const char *preamble, const char *desc, const char *comment)
{
	if (result == expected) {
		fprintf(stderr, "%s: %s, %s ... " COL_GREEN "PASSED\n" COL_STOP,
			preamble, desc, comment);
		return 1;
	} else {
		fprintf(stderr, "%s: %s, %s ... " COL_RED "FAILED\n" COL_STOP,
			preamble, desc, comment);
		fprintf(stderr, "	Expected: 0x%lx, got: 0x%lx\n",
			(long)expected, (long)result);
		gTestSuccess = 0;
		return 0;
	}
}

int test_ptr_neq(const void* result, const void* not_expected,
		 const char *preamble, const char *desc, const char *comment)
{
	if (result != not_expected) {
		fprintf(stderr, "%s: %s, %s ... " COL_GREEN "PASSED\n" COL_STOP,
			preamble, desc, comment);
		return 1;
	} else {
		fprintf(stderr, "%s: %s, %s ... " COL_RED "FAILED\n" COL_STOP,
			preamble, desc, comment);
		fprintf(stderr, "	Didn't expect 0x%lx, but got it\n",
			(long)not_expected);
		gTestSuccess = 0;
		return 0;
	}
}

int test_str_eq(const char* result, const char* expected,
		const char *preamble, const char *desc, const char *comment)
{
	if (!result || !expected) {
		fprintf(stderr, "%s: %s, %s ... " COL_RED "FAILED\n" COL_STOP,
			preamble, desc, comment);
		fprintf(stderr, "	String compare with NULL\n");
		gTestSuccess = 0;
		return 0;
	} else if (!strcmp(result, expected)) {
		fprintf(stderr, "%s: %s, %s ... " COL_GREEN "PASSED\n" COL_STOP,
			preamble, desc, comment);
		return 1;
	} else {
		fprintf(stderr, "%s " COL_RED "FAILED\n" COL_STOP, comment);
		fprintf(stderr, "	Expected: \"%s\", got: \"%s\"\n",
			expected, result);
		gTestSuccess = 0;
		return 0;
	}
}

int test_str_neq(const char* result, const char* not_expected,
		 const char *preamble, const char *desc, const char *comment)
{
	if (!result || !not_expected) {
		fprintf(stderr, "%s: %s, %s ... " COL_RED "FAILED\n" COL_STOP,
			preamble, desc, comment);
		fprintf(stderr, "	String compare with NULL\n");
		gTestSuccess = 0;
		return 0;
	} else if (strcmp(result, not_expected)) {
		fprintf(stderr, "%s: %s, %s ... " COL_GREEN "PASSED\n" COL_STOP,
			preamble, desc, comment);
		return 1;
	} else {
		fprintf(stderr, "%s: %s, %s ... " COL_RED "FAILED\n" COL_STOP,
			preamble, desc, comment);
		fprintf(stderr, "	Didn't expect: \"%s\", but got it\n",
			not_expected);
		gTestSuccess = 0;
		return 0;
	}
}

int test_succ(int result,
	      const char *preamble, const char *desc, const char *comment)
{
	if (result == 0) {
		fprintf(stderr, "%s: %s ... " COL_GREEN "PASSED\n" COL_STOP,
			preamble, comment ? comment : desc);
	} else {
		fprintf(stderr, "%s: %s ... " COL_RED "FAILED\n" COL_STOP,
			preamble, comment ? comment : desc);
		fprintf(stderr, "	Expected SUCCESS, got: %#x (%d)\n",
			result, result);
		gTestSuccess = 0;
	}
	return !result;
}

int test_true(int result,
	      const char *preamble, const char *desc, const char *comment)
{
	if (result) {
		fprintf(stderr, "%s: %s, %s ... " COL_GREEN "PASSED\n" COL_STOP,
			preamble, desc, comment);
	} else {
		fprintf(stderr, "%s: %s, %s ... " COL_RED "FAILED\n" COL_STOP,
			preamble, desc, comment);
		fprintf(stderr, "	Expected TRUE, got 0\n");
		gTestSuccess = 0;
	}
	return result;
}

int test_false(int result,
	       const char *preamble, const char *desc, const char *comment)
{
	if (!result) {
		fprintf(stderr, "%s: %s, %s ... " COL_GREEN "PASSED\n" COL_STOP,
			preamble, desc, comment);
	} else {
		fprintf(stderr, "%s: %s, %s ... " COL_RED "FAILED\n" COL_STOP,
			preamble, desc, comment);
		fprintf(stderr, "	Expected FALSE, got: 0x%lx\n",
			(long)result);
		gTestSuccess = 0;
	}
	return !result;
}

int test_abort(int aborted,
	       const char *preamble, const char *desc, const char *comment)
{
	if (aborted) {
		fprintf(stderr, "%s: %s ... " COL_GREEN "PASSED\n" COL_STOP,
			preamble, comment ? comment : desc);
	} else {
		fprintf(stderr, "%s: %s ... " COL_RED "FAILED\n" COL_STOP,
			preamble, comment ? comment : desc);
		fprintf(stderr, "	Expected ABORT, but did not get it\n");
		gTestSuccess = 0;
	}
	return aborted;
}

void vb2ex_abort(void)
{
	/*
	 * If expecting an abort call, jump back to TEST_ABORT macro.
	 * Otherwise, force exit to ensure the test fails.
	 */
	if (gTestAbortArmed) {
		longjmp(gTestJmpEnv, 1);
	} else {
		fprintf(stderr, COL_RED "Unexpected ABORT encountered, "
			"exiting\n" COL_STOP);
		exit(1);
	}
}

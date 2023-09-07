/* Copyright 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Test utilities.
 */

#ifdef TEST_COVERAGE
#include <signal.h>
#include <stdlib.h>
#endif

#include "console.h"
#include "hooks.h"
#include "host_command.h"
#include "system.h"
#include "task.h"
#include "test_util.h"
#include "util.h"

struct test_util_tag {
	uint8_t error_count;
};

#define TEST_UTIL_SYSJUMP_TAG 0x5455 /* "TU" */
#define TEST_UTIL_SYSJUMP_VERSION 1

int __test_error_count;

/* Weak reference function as an entry point for unit test */
test_mockable void run_test(void) { }

/* Default dummy test init */
test_mockable void test_init(void) { }

/* Default dummy before test */
test_mockable void before_test(void) { }

/* Default dummy after test */
test_mockable void after_test(void) { }

#ifdef TEST_COVERAGE
extern void __gcov_flush(void);

void emulator_flush(void)
{
	__gcov_flush();
}

void test_end_hook(int sig)
{
	emulator_flush();
	exit(0);
}

void register_test_end_hook(void)
{
	signal(SIGTERM, test_end_hook);
}
#else
void emulator_flush(void)
{
}

void register_test_end_hook(void)
{
}
#endif

void test_reset(void)
{
	if (!system_jumped_to_this_image())
		__test_error_count = 0;
}

void test_pass(void)
{
	ccprintf("Pass!\n");
}

void test_fail(void)
{
	ccprintf("Fail!\n");
}

void test_print_result(void)
{
	if (__test_error_count)
		ccprintf("Fail! (%d tests)\n", __test_error_count);
	else
		ccprintf("Pass!\n");
}

int test_get_error_count(void)
{
	return __test_error_count;
}

uint32_t test_get_state(void)
{
	return system_get_scratchpad();
}

test_mockable void test_clean_up(void)
{
}

void test_reboot_to_next_step(enum test_state_t step)
{
	ccprintf("Rebooting to next test step...\n");
	cflush();
	system_set_scratchpad(TEST_STATE_MASK(step));
	system_reset(SYSTEM_RESET_HARD);
}

test_mockable void test_run_step(uint32_t state)
{
}

void test_run_multistep(void)
{
	uint32_t state = test_get_state();

	if (state & TEST_STATE_MASK(TEST_STATE_PASSED)) {
		test_clean_up();
		system_set_scratchpad(0);
		test_pass();
	} else if (state & TEST_STATE_MASK(TEST_STATE_FAILED)) {
		test_clean_up();
		system_set_scratchpad(0);
		test_fail();
	}

	if (state & TEST_STATE_STEP_1 || state == 0) {
		task_wait_event(-1); /* Wait for run_test() */
		test_run_step(TEST_STATE_MASK(TEST_STATE_STEP_1));
	} else {
		test_run_step(state);
	}
}

#ifdef HAS_TASK_HOSTCMD
int test_send_host_command(int command, int version, const void *params,
			   int params_size, void *resp, int resp_size)
{
	struct host_cmd_handler_args args;

	args.version = version;
	args.command = command;
	args.params = params;
	args.params_size = params_size;
	args.response = resp;
	args.response_max = resp_size;
	args.response_size = 0;

	return host_command_process(&args);
}
#endif  /* TASK_HAS_HOSTCMD */

/* Linear congruential pseudo random number generator */
uint32_t prng(uint32_t seed)
{
	return 22695477 * seed + 1;
}

uint32_t prng_no_seed(void)
{
	static uint32_t seed = 0x1234abcd;
	return seed = prng(seed);
}

static void restore_state(void)
{
	const struct test_util_tag *tag;
	int version, size;

	tag = (const struct test_util_tag *)system_get_jump_tag(
		TEST_UTIL_SYSJUMP_TAG, &version, &size);
	if (tag && version == TEST_UTIL_SYSJUMP_VERSION &&
	    size == sizeof(*tag))
		__test_error_count = tag->error_count;
	else
		__test_error_count = 0;
}
DECLARE_HOOK(HOOK_INIT, restore_state, HOOK_PRIO_DEFAULT);

static void preserve_state(void)
{
	struct test_util_tag tag;
	tag.error_count = __test_error_count;
	system_add_jump_tag(TEST_UTIL_SYSJUMP_TAG, TEST_UTIL_SYSJUMP_VERSION,
			    sizeof(tag), &tag);
}
DECLARE_HOOK(HOOK_SYSJUMP, preserve_state, HOOK_PRIO_DEFAULT);

static int command_run_test(int argc, char **argv)
{
	run_test();
	return EC_SUCCESS;
}
DECLARE_CONSOLE_COMMAND(runtest, command_run_test,
			NULL, NULL);

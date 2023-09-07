/* Copyright 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include "common.h"
#include "console.h"
#include "cpu.h"
#include "hooks.h"
#include "host_command.h"
#include "panic.h"
#include "printf.h"
#include "software_panic.h"
#include "system.h"
#include "task.h"
#include "timer.h"
#include "uart.h"
#include "util.h"

/* Panic data goes at the end of RAM. */
static struct panic_data * const pdata_ptr = PANIC_DATA_PTR;

/* Common SW Panic reasons strings */
const char * const panic_sw_reasons[] = {
#ifdef CONFIG_SOFTWARE_PANIC
	"PANIC_SW_DIV_ZERO",
	"PANIC_SW_STACK_OVERFLOW",
	"PANIC_SW_PD_CRASH",
	"PANIC_SW_ASSERT",
	"PANIC_SW_WATCHDOG",
	"PANIC_SW_RNG",
	"PANIC_SW_PMIC_FAULT",
#endif
};

/**
 * Check an interrupt vector as being a valid software panic
 * @param reason	Reason for panic
 * @return 0 if not a valid software panic reason, otherwise non-zero.
 */
int panic_sw_reason_is_valid(uint32_t reason)
{
	return (IS_ENABLED(CONFIG_SOFTWARE_PANIC) &&
		reason >= PANIC_SW_BASE &&
		(reason - PANIC_SW_BASE) < ARRAY_SIZE(panic_sw_reasons));
}

/**
 * Add a character directly to the UART buffer.
 *
 * @param context	Context; ignored.
 * @param c		Character to write.
 * @return 0 if the character was transmitted, 1 if it was dropped.
 */
#ifndef CONFIG_DEBUG_PRINTF
static int panic_txchar(void *context, int c)
{
	if (c == '\n')
		panic_txchar(context, '\r');

	/* Wait for space in transmit FIFO */
	while (!uart_tx_ready())
		;

	/* Write the character directly to the transmit FIFO */
	uart_write_char(c);

	return 0;
}

void panic_puts(const char *outstr)
{
	/* Flush the output buffer */
	uart_flush_output();

	/* Put all characters in the output buffer */
	while (*outstr)
		panic_txchar(NULL, *outstr++);

	/* Flush the transmit FIFO */
	uart_tx_flush();
}

void panic_printf(const char *format, ...)
{
	va_list args;

	/* Flush the output buffer */
	uart_flush_output();

	va_start(args, format);
	vfnprintf(panic_txchar, NULL, format, args);
	va_end(args);

	/* Flush the transmit FIFO */
	uart_tx_flush();
}
#endif

/**
 * Display a message and reboot
 */
void panic_reboot(void)
{
	panic_puts("\n\nRebooting...\n");
	system_reset(0);
}

#ifdef CONFIG_DEBUG_ASSERT_REBOOTS
#ifdef CONFIG_DEBUG_ASSERT_BRIEF
void panic_assert_fail(const char *fname, int linenum)
{
	panic_printf("\nASSERTION FAILURE at %s:%d\n", fname, linenum);
#ifdef CONFIG_SOFTWARE_PANIC
	software_panic(PANIC_SW_ASSERT, linenum);
#else
	panic_reboot();
#endif
}
#else
void panic_assert_fail(const char *msg, const char *func, const char *fname,
		       int linenum)
{
	panic_printf("\nASSERTION FAILURE '%s' in %s() at %s:%d\n",
		     msg, func, fname, linenum);
#ifdef CONFIG_SOFTWARE_PANIC
	software_panic(PANIC_SW_ASSERT, linenum);
#else
	panic_reboot();
#endif
}
#endif
#endif

void panic(const char *msg)
{
	panic_printf("\n** PANIC: %s\n", msg);
	panic_reboot();
}

struct panic_data *panic_get_data(void)
{
	BUILD_ASSERT(sizeof(struct panic_data) <= CONFIG_PANIC_DATA_SIZE);
	return pdata_ptr->magic == PANIC_DATA_MAGIC ? pdata_ptr : NULL;
}

static void panic_init(void)
{
#ifdef CONFIG_HOSTCMD_EVENTS
	struct panic_data *addr = panic_get_data();

	/* Notify host of new panic event */
	if (addr && !(addr->flags & PANIC_DATA_FLAG_OLD_HOSTEVENT)) {
		host_set_single_event(EC_HOST_EVENT_PANIC);
		addr->flags |= PANIC_DATA_FLAG_OLD_HOSTEVENT;
	}
#endif
}
DECLARE_HOOK(HOOK_INIT, panic_init, HOOK_PRIO_LAST);
DECLARE_HOOK(HOOK_CHIPSET_RESET, panic_init, HOOK_PRIO_LAST);

#ifdef CONFIG_CMD_STACKOVERFLOW
static void stack_overflow_recurse(int n)
{
	ccprintf("+%d", n);

	/*
	 * Force task context switch, since that's where we do stack overflow
	 * checking.
	 */
	msleep(10);

	stack_overflow_recurse(n+1);

	/*
	 * Do work after the recursion, or else the compiler uses tail-chaining
	 * and we don't actually consume additional stack.
	 */
	ccprintf("-%d", n);
}
#endif /* CONFIG_CMD_STACKOVERFLOW */

/*****************************************************************************/
/* Console commands */
#ifdef CONFIG_CMD_CRASH
static int command_crash(int argc, char **argv)
{
	if (argc < 2)
		return EC_ERROR_PARAM1;

	if (!strcasecmp(argv[1], "assert")) {
		ASSERT(0);
	} else if (!strcasecmp(argv[1], "divzero")) {
		volatile int zero = 0;

		cflush();
		ccprintf("%08x", (long)1 / zero);
	} else if (!strcasecmp(argv[1], "udivzero")) {
		volatile int zero = 0;

		cflush();
		ccprintf("%08x", (unsigned long)1 / zero);
#ifdef CONFIG_CMD_STACKOVERFLOW
	} else if (!strcasecmp(argv[1], "stack")) {
		stack_overflow_recurse(1);
#endif
	} else if (!strcasecmp(argv[1], "unaligned")) {
		volatile intptr_t unaligned_ptr = 0xcdef;
		cflush();
		ccprintf("%08x", *(volatile int *)unaligned_ptr);
	} else if (!strcasecmp(argv[1], "watchdog")) {
		while (1)
			;
	} else if (!strcasecmp(argv[1], "hang")) {
		interrupt_disable();
		while (1)
			;
	} else {
		return EC_ERROR_PARAM1;
	}

	/* Everything crashes, so shouldn't get back here */
	return EC_ERROR_UNKNOWN;
}
DECLARE_CONSOLE_COMMAND(crash, command_crash,
		"[assert | divzero | udivzero"
#ifdef CONFIG_CMD_STACKOVERFLOW
			" | stack"
#endif
			" | unaligned | watchdog | hang]",
		"Crash the system (for testing)");
#endif

static int command_panicinfo(int argc, char **argv)
{
	if (pdata_ptr->magic == PANIC_DATA_MAGIC) {
		ccprintf("Saved panic data:%s\n",
			 (pdata_ptr->flags & PANIC_DATA_FLAG_OLD_CONSOLE ?
			  "" : " (NEW)"));

		panic_data_print(pdata_ptr);

		/* Data has now been printed */
		pdata_ptr->flags |= PANIC_DATA_FLAG_OLD_CONSOLE;
	} else {
		ccprintf("No saved panic data available.\n");
	}
	return EC_SUCCESS;
}
DECLARE_CONSOLE_COMMAND(panicinfo, command_panicinfo,
			NULL,
			"Print info from a previous panic");

/*****************************************************************************/
/* Host commands */

enum ec_status host_command_panic_info(struct host_cmd_handler_args *args)
{
	if (pdata_ptr->magic == PANIC_DATA_MAGIC) {
		ASSERT(pdata_ptr->struct_size <= args->response_max);
		memcpy(args->response, pdata_ptr, pdata_ptr->struct_size);
		args->response_size = pdata_ptr->struct_size;

		/* Data has now been returned */
		pdata_ptr->flags |= PANIC_DATA_FLAG_OLD_HOSTCMD;
	}

	return EC_RES_SUCCESS;
}
DECLARE_HOST_COMMAND(EC_CMD_GET_PANIC_INFO,
		     host_command_panic_info,
		     EC_VER_MASK(0));

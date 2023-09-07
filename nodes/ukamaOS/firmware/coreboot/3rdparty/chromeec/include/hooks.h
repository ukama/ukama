/* Copyright 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/* System hooks for Chrome EC */

#ifndef __CROS_EC_HOOKS_H
#define __CROS_EC_HOOKS_H

#include "common.h"

enum hook_priority {
	/* Generic values across all hooks */
	HOOK_PRIO_FIRST = 1,       /* Highest priority */
	HOOK_PRIO_DEFAULT = 5000,  /* Default priority */
	HOOK_PRIO_LAST = 9999,     /* Lowest priority */

	/* Specific hook vales for HOOK_INIT */
	/* DMA inits before ADC, I2C, SPI */
	HOOK_PRIO_INIT_DMA = HOOK_PRIO_FIRST + 1,
	/* LPC inits before modules which need memory-mapped I/O */
	HOOK_PRIO_INIT_LPC = HOOK_PRIO_FIRST + 1,
	/* I2C is needed before chipset inits (battery communications). */
	HOOK_PRIO_INIT_I2C = HOOK_PRIO_FIRST + 2,
	/* Chipset inits before modules which need to know its initial state. */
	HOOK_PRIO_INIT_CHIPSET = HOOK_PRIO_FIRST + 3,
	/* Lid switch inits before power button */
	HOOK_PRIO_INIT_LID = HOOK_PRIO_FIRST + 4,
	/* Power button inits before chipset and switch */
	HOOK_PRIO_INIT_POWER_BUTTON = HOOK_PRIO_FIRST + 5,
	/* Init switch states after power button / lid */
	HOOK_PRIO_INIT_SWITCH = HOOK_PRIO_FIRST + 6,
	/* Init fan before PWM */
	HOOK_PRIO_INIT_FAN = HOOK_PRIO_FIRST + 7,
	/* PWM inits before modules which might use it (LEDs) */
	HOOK_PRIO_INIT_PWM = HOOK_PRIO_FIRST + 8,
	/* SPI inits before modules which might use it (sensors) */
	HOOK_PRIO_INIT_SPI = HOOK_PRIO_FIRST + 9,
	/* Extpower inits before modules which might use it (battery, LEDs) */
	HOOK_PRIO_INIT_EXTPOWER = HOOK_PRIO_FIRST + 10,
	/* Init VBOOT hash later, since it depends on deferred functions */
	HOOK_PRIO_INIT_VBOOT_HASH = HOOK_PRIO_FIRST + 11,
	/* Init charge manager before usage in board init */
	HOOK_PRIO_CHARGE_MANAGER_INIT = HOOK_PRIO_FIRST + 12,

	HOOK_PRIO_INIT_ADC = HOOK_PRIO_DEFAULT,

	/* Specific values to lump temperature-related hooks together */
	HOOK_PRIO_TEMP_SENSOR = 6000,
	/* After all sensors have been polled */
	HOOK_PRIO_TEMP_SENSOR_DONE = HOOK_PRIO_TEMP_SENSOR + 1,
};

enum hook_type {
	/*
	 * System initialization.
	 *
	 * Hook routines are called from main(), after all hard-coded inits,
	 * before task scheduling is enabled.
	 */
	HOOK_INIT = 0,

	/*
	 * System clock changed frequency.
	 *
	 * The "pre" frequency hook is called before we change the frequency.
	 * There is no way to cancel.  Hook routines are always called from
	 * a task, so it's OK to lock a mutex here.  However, they may be called
	 * from a deferred task on some platforms so callbacks must make sure
	 * not to do anything that would require some other deferred task to
	 * run.
	 */
	HOOK_PRE_FREQ_CHANGE,
	HOOK_FREQ_CHANGE,

	/*
	 * About to jump to another image.  Modules which need to preserve data
	 * across such a jump should save it here and restore it in HOOK_INIT.
	 *
	 * Hook routines are called from the context which initiates the jump,
	 * WITH INTERRUPTS DISABLED.
	 */
	HOOK_SYSJUMP,

	/*
	 * Initialization for components such as PMU to be done before host
	 * chipset/AP starts up.
	 *
	 * Hook routines are called from the chipset task.
	 */
	HOOK_CHIPSET_PRE_INIT,

	/* System is starting up.  All suspend rails are now on.
	 *
	 * Hook routines are called from the chipset task.
	 */
	HOOK_CHIPSET_STARTUP,

	/*
	 * System is resuming from suspend, or booting and has reached the
	 * point where all voltage rails are on.
	 *
	 * Hook routines are called from the chipset task.
	 */
	HOOK_CHIPSET_RESUME,

	/*
	 * System is suspending, or shutting down; all voltage rails are still
	 * on.
	 *
	 * Hook routines are called from the chipset task.
	 */
	HOOK_CHIPSET_SUSPEND,

	/*
	 * System is shutting down.  All suspend rails are still on.
	 *
	 * Hook routines are called from the chipset task.
	 */
	HOOK_CHIPSET_SHUTDOWN,

	/*
	 * System reset in S0.  All rails are still up.
	 *
	 * Hook routines are called from the chipset task.
	 */
	HOOK_CHIPSET_RESET,

	/*
	 * AC power plugged in or removed.
	 *
	 * Hook routines are called from the TICK task.
	 */
	HOOK_AC_CHANGE,

	/*
	 * Lid opened or closed.  Based on debounced lid state, not raw lid
	 * GPIO input.
	 *
	 * Hook routines are called from the TICK task.
	 */
	HOOK_LID_CHANGE,

	/*
	 * Device in tablet mode (base behind lid).
	 *
	 * Hook routines are called from the TICK task.
	 */
	HOOK_TABLET_MODE_CHANGE,

	/*
	 * Detachable device connected to a base.
	 *
	 * Hook routines are called from the TICK task.
	 */
	HOOK_BASE_ATTACHED_CHANGE,

	/*
	 * Power button pressed or released.  Based on debounced power button
	 * state, not raw GPIO input.
	 *
	 * Hook routines are called from the TICK task.
	 */
	HOOK_POWER_BUTTON_CHANGE,

	/*
	 * Battery state of charge changed
	 *
	 * Hook routines are called from the charger task.
	 */
	HOOK_BATTERY_SOC_CHANGE,

#ifdef CONFIG_CASE_CLOSED_DEBUG_V1
	/*
	 * Case-closed debugging configuration changed.
	 *
	 * Hook routines are called from the TICK, console, or TPM task.
	 */
	HOOK_CCD_CHANGE,
#endif

#ifdef CONFIG_USB_SUSPEND
	/*
	 * Called when there is a change in USB power management status
	 * (suspended or resumed).
	 *
	 * Hook routines are called from HOOKS task.
	 */
	HOOK_USB_PM_CHANGE,
#endif

	/*
	 * Periodic tick, every HOOK_TICK_INTERVAL.
	 *
	 * Hook routines will be called from the TICK task.
	 */
	HOOK_TICK,

	/*
	 * Periodic tick, every second.
	 *
	 * Hook routines will be called from the TICK task.
	 */
	HOOK_SECOND,

	/*
	 * USB PD cc disconnect event.
	 */
	HOOK_USB_PD_DISCONNECT,

	/*
	 * USB PD cc connection event.
	 */
	HOOK_USB_PD_CONNECT,
};

struct hook_data {
	/* Hook processing routine. */
	void (*routine)(void);
	/* Priority; low numbers = higher priority. */
	int priority;
};

/**
 * Call all the hook routines of a specified type.
 *
 * This function must be called from the correct type-specific context (task);
 * see enum hook_type for details.  hook_notify() should NEVER be called from
 * interrupt context unless specifically allowed for a hook type, because hook
 * routines may need to perform task-level calls like usleep() and mutex
 * operations that are not valid in interrupt context.  Instead of calling a
 * hook from interrupt context, use a deferred function.
 *
 * @param type		Type of hook routines to call.
 */
void hook_notify(enum hook_type type);

struct deferred_data {
	/* Deferred function pointer */
	void (*routine)(void);
};

/**
 * Start a timer to call a deferred routine.
 *
 * The routine will be called after at least the specified delay, in the
 * context of the hook task.
 *
 * @param data	The deferred_data struct created by invoking DECLARE_DEFERRED().
 * @param us	Delay in microseconds until routine will be called.  If the
 *		routine is already pending, subsequent calls will change the
 *		delay.  Pass us=0 to call as soon as possible, or -1 to cancel
 *		the deferred call.
 *
 * @return non-zero if error.
 */
int hook_call_deferred(const struct deferred_data *data, int us);

#ifdef CONFIG_COMMON_RUNTIME
/**
 * Register a hook routine.
 *
 * NOTE: Hook routines must be careful not to leave resources locked which may
 * be needed by other hook routines or deferred function calls.  This can cause
 * a deadlock, because most hooks and all deferred functions are called from
 * the same hook task.  For example:
 *
 *   hook1(): lock foo
 *   deferred1(): lock foo, use foo, unlock foo
 *   hook2(): unlock foo
 *
 * In this example, hook1() and hook2() lock and unlock a shared resource foo
 * (for example, a mutex).  If deferred1() attempts to lock the resource, it
 * will stall waiting for the resource to be unlocked.  But the unlock will
 * never happen, because hook2() won't be called by the hook task until
 * deferred1() returns.
 *
 * @param hooktype	Type of hook for routine (enum hook_type)
 * @param routine	Hook routine, with prototype void routine(void)
 * @param priority      Priority for determining when routine is called vs.
 *			other hook routines; should be between HOOK_PRIO_FIRST
 *                      and HOOK_PRIO_LAST, and should be HOOK_PRIO_DEFAULT
 *			unless there's a compelling reason to care about the
 *			order in which hooks are called.
 */
#define DECLARE_HOOK(hooktype, routine, priority)			\
	const struct hook_data __keep __no_sanitize_address		\
	CONCAT4(__hook_, hooktype, _, routine)				\
	__attribute__((section(".rodata." STRINGIFY(hooktype))))	\
	     = {routine, priority}

/**
 * Register a deferred function call.
 *
 * DECLARE_DEFERRED creates a new deferred_data struct with a name constructed
 * by concatenating _data to the name of the routine passed.
 *
 * To call a deferred routine defined as:
 *     DECLARE_DEFERRED(foo)
 * You would call
 *     hook_call_deferred(&foo_data, delay_in_microseconds);
 *
 * NOTE: Deferred function call routines must be careful not to leave resources
 * locked which may be needed by other hook routines or deferred function
 * calls.  This can cause a deadlock, because most hooks and all deferred
 * functions are called from the same hook task.  See DECLARE_HOOK() for an
 * example.
 *
 * @param routine	Function pointer, with prototype void routine(void)
 */
#define DECLARE_DEFERRED(routine)					\
	const struct deferred_data __keep __no_sanitize_address		\
	CONCAT2(routine, _data)						\
	__attribute__((section(".rodata.deferred")))			\
	     = {routine}

#else /* CONFIG_COMMON_RUNTIME */
#define DECLARE_HOOK(t, func, p)				\
	void CONCAT2(unused_hook_, func)(void) { func(); }
#define DECLARE_DEFERRED(func)					\
	void CONCAT2(unused_deferred_, func)(void) { func(); }
#endif /* CONFIG_COMMON_RUNTIME */

#endif  /* __CROS_EC_HOOKS_H */

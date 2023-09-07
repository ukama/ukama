/* Copyright 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/* Keyboard scanner module for Chrome EC */

#include "chipset.h"
#include "clock.h"
#include "common.h"
#include "console.h"
#include "ec_commands.h"
#include "hooks.h"
#include "host_command.h"
#include "keyboard_config.h"
#include "keyboard_protocol.h"
#include "keyboard_raw.h"
#include "keyboard_scan.h"
#include "lid_switch.h"
#include "switch.h"
#include "system.h"
#include "tablet_mode.h"
#include "task.h"
#include "timer.h"
#include "usb_api.h"
#include "util.h"

/* Console output macros */
#define CPUTS(outstr) cputs(CC_KEYSCAN, outstr)
#define CPRINTF(format, args...) cprintf(CC_KEYSCAN, format, ## args)
#define CPRINTS(format, args...) cprints(CC_KEYSCAN, format, ## args)

#ifdef CONFIG_KEYBOARD_DEBUG
#define CPUTS5(outstr) cputs(CC_KEYSCAN, outstr)
#define CPRINTS5(format, args...) cprints(CC_KEYBOARD, format, ## args)
#else
#define CPUTS5(outstr)
#define CPRINTS5(format, args...)
#endif

#define SCAN_TIME_COUNT 32  /* Number of last scan times to track */

/* If we're waiting for a scan to happen, we'll give it this long */
#define SCAN_TASK_TIMEOUT_US	(100 * MSEC)

#ifndef CONFIG_KEYBOARD_POST_SCAN_CLOCKS
/*
 * Default delay in clocks; this was experimentally determined to be long
 * enough to avoid watchdog warnings or I2C errors on a typical notebook
 * config on STM32.
 */
#define CONFIG_KEYBOARD_POST_SCAN_CLOCKS 16000
#endif

#ifndef CONFIG_KEYBOARD_BOARD_CONFIG
/* Use default keyboard scan config, because board didn't supply one */
struct keyboard_scan_config keyscan_config = {
	.output_settle_us = 50,
	.debounce_down_us = 9 * MSEC,
	.debounce_up_us = 30 * MSEC,
	.scan_period_us = 3 * MSEC,
	.min_post_scan_delay_us = 1000,
	.poll_timeout_us = 100 * MSEC,
	.actual_key_mask = {
		0x14, 0xff, 0xff, 0xff, 0xff, 0xf5, 0xff,
		0xa4, 0xff, 0xfe, 0x55, 0xfa, 0xca  /* full set */
	},
};
#endif

/* Boot key list.  Must be in same order as enum boot_key. */
struct boot_key_entry {
	uint8_t mask_index;
	uint8_t mask_value;
};

#ifdef CONFIG_KEYBOARD_BOOT_KEYS
static const struct boot_key_entry boot_key_list[] = {
	{KEYBOARD_COL_ESC, KEYBOARD_MASK_ESC},   /* Esc */
	{KEYBOARD_COL_DOWN, KEYBOARD_MASK_DOWN}, /* Down-arrow */
	{KEYBOARD_COL_LEFT_SHIFT, KEYBOARD_MASK_LEFT_SHIFT}, /* Left-Shift */
};
static uint32_t boot_key_value = BOOT_KEY_NONE;
#endif

uint8_t keyboard_cols = KEYBOARD_COLS_MAX;

/* Debounced key matrix */
static uint8_t __bss_slow debounced_state[KEYBOARD_COLS_MAX];
/* Mask of keys being debounced */
static uint8_t __bss_slow debouncing[KEYBOARD_COLS_MAX];
/* Keys simulated-pressed */
static uint8_t __bss_slow simulated_key[KEYBOARD_COLS_MAX];
#ifdef CONFIG_KEYBOARD_LANGUAGE_ID
static uint8_t __bss_slow keyboard_id[KEYBOARD_IDS];
#endif

/* Times of last scans */
static uint32_t __bss_slow scan_time[SCAN_TIME_COUNT];
/* Current scan_time[] index */
static int __bss_slow scan_time_index;

/* Index into scan_time[] when each key started debouncing */
static uint8_t __bss_slow scan_edge_index[KEYBOARD_COLS_MAX][KEYBOARD_ROWS];

/* Minimum delay between keyboard scans based on current clock frequency */
static uint32_t __bss_slow post_scan_clock_us;

/*
 * Print all keyboard scan state changes?  Off by default because it generates
 * a lot of debug output, which makes the saved EC console data less useful.
 */
static int __bss_slow print_state_changes;

/* Must init to 0 for scanning at boot */
static volatile uint32_t __bss_slow disable_scanning_mask;

/* Constantly incrementing counter of the number of times we polled */
static volatile int kbd_polls;

/* If true, we'll force a keyboard poll */
static volatile int force_poll;

static int keyboard_scan_is_enabled(void)
{
	/* NOTE: this is just an instantaneous glimpse of the variable. */
	return !disable_scanning_mask;
}

void keyboard_scan_enable(int enable, enum kb_scan_disable_masks mask)
{
	/* Access atomically */
	if (enable) {
		atomic_clear((uint32_t *)&disable_scanning_mask, mask);
	} else {
		atomic_or((uint32_t *)&disable_scanning_mask, mask);
		clear_typematic_key();
	}

	/* Let the task figure things out */
	task_wake(TASK_ID_KEYSCAN);
}

/**
 * Print the keyboard state.
 *
 * @param state		State array to print
 * @param msg		Description of state
 */
static void print_state(const uint8_t *state, const char *msg)
{
	int c;

	CPRINTF("[%T KB %s:", msg);
	for (c = 0; c < keyboard_cols; c++) {
		if (state[c])
			CPRINTF(" %02x", state[c]);
		else
			CPUTS(" --");
	}
	CPUTS("]\n");
}

/**
 * Ensure that the keyboard has been scanned.
 *
 * Makes sure that we've fully gone through the keyboard scanning loop at
 * least once.
 */
static void ensure_keyboard_scanned(int old_polls)
{
	uint64_t start_time;

	start_time = get_time().val;

	/*
	 * Ensure we see the poll task run.
	 *
	 * Note that the poll task is higher priority than ours so we know that
	 * while we're running it's not partway through a poll.  That means that
	 * if kbd_polls changes we've gone through a whole cycle.
	 */
	while ((kbd_polls == old_polls) &&
	       (get_time().val - start_time < SCAN_TASK_TIMEOUT_US))
		usleep(keyscan_config.scan_period_us);
}

/**
 * Simulate a keypress.
 *
 * @param row		Row of key
 * @param col		Column of key
 * @param pressed	Non-zero if pressed, zero if released
 */
static void simulate_key(int row, int col, int pressed)
{
	int old_polls;

	if ((simulated_key[col] & BIT(row)) == ((pressed ? 1 : 0) << row))
		return;  /* No change */

	simulated_key[col] ^= BIT(row);

	/* Keep track of polls now that we've got keys simulated */
	old_polls = kbd_polls;

	print_state(simulated_key, "simulated ");

	/* Force a poll even though no keys are pressed */
	force_poll = 1;

	/* Wake the task to handle changes in simulated keys */
	task_wake(TASK_ID_KEYSCAN);

	/*
	 * Make sure that the keyboard task sees the key for long enough.
	 * That means it needs to have run and for enough time.
	 */
	ensure_keyboard_scanned(old_polls);
	usleep(pressed ?
	       keyscan_config.debounce_down_us : keyscan_config.debounce_up_us);
	ensure_keyboard_scanned(kbd_polls);
}

/**
 * Read the raw keyboard matrix state.
 *
 * Used in pre-init, so must not make task-switching-dependent calls; udelay()
 * is ok because it's a spin-loop.
 *
 * @param state		Destination for new state (must be KEYBOARD_COLS_MAX long).
 *
 * @return 1 if at least one key is pressed, else zero.
 */
static int read_matrix(uint8_t *state)
{
	int c;
	uint8_t r;
	int pressed = 0;

	for (c = 0; c < keyboard_cols; c++) {
		/*
		 * Stop if scanning becomes disabled. Note, scanning is enabled
		 * on boot by default.
		 */
		if (!keyboard_scan_is_enabled())
			break;

		/* Select column, then wait a bit for it to settle */
		keyboard_raw_drive_column(c);
		udelay(keyscan_config.output_settle_us);

		/* Read the row state */
		r = keyboard_raw_read_rows();

		/* Add in simulated keypresses */
		r |= simulated_key[c];

		/*
		 * Keep track of what keys appear to be pressed.  Even if they
		 * don't exist in the matrix, they'll keep triggering
		 * interrupts, so we can't leave scanning mode.
		 */
		pressed |= r;

		/* Mask off keys that don't exist on the actual keyboard */
		r &= keyscan_config.actual_key_mask[c];

#ifdef CONFIG_KEYBOARD_TEST
		/* Use simulated keyscan sequence instead if testing active */
		r = keyscan_seq_get_scan(c, r);
#endif

		/* Store the masked state */
		state[c] = r;
	}

	keyboard_raw_drive_column(KEYBOARD_COLUMN_NONE);

	return pressed ? 1 : 0;
}

#ifdef CONFIG_KEYBOARD_LANGUAGE_ID
/**
 * Read the raw keyboard IDs state.
 *
 * Used in pre-init, so must not make task-switching-dependent calls; udelay()
 * is ok because it's a spin-loop.
 *
 * @param id		Destination for keyboard id (must be KEYBOARD_IDS long).
 *
 */
static void read_matrix_id(uint8_t *id)
{
	int c;

	for (c = 0; c < KEYBOARD_IDS; c++) {
		/* Select the ID pin, then wait a bit for it to settle.
		 * Caveat: If a keyboard maker puts ID pins right after scan
		 * columns, we can't support variable column size with a single
		 * image. */
		keyboard_raw_drive_column(KEYBOARD_COLS_MAX + c);
		udelay(keyscan_config.output_settle_us);

		/* Read the row state */
		id[c] = keyboard_raw_read_rows();

		CPRINTS("Keyboard ID%u: 0x%02x", c, id[c]);
	}

	keyboard_raw_drive_column(KEYBOARD_COLUMN_NONE);
}
#endif

#ifdef CONFIG_KEYBOARD_RUNTIME_KEYS
/**
 * Check special runtime key combinations.
 *
 * @param state		Keyboard state to use when checking keys.
 *
 * @return 1 if a special key was pressed, 0 if not
 */
static int check_runtime_keys(const uint8_t *state)
{
	int num_press = 0;
	int c;

#ifdef BOARD_SAMUS
	int16_t chg_override;

	/*
	 * TODO(crosbug.com/p/34850): remove these hot-keys for samus, should
	 * be done at higher level than this.
	 */
	/*
	 * On samus, ctrl + search + 0|1|2 sets the active charge port
	 * by sending the charge override host command. Should only be sent
	 * when chipset is in S0. Note that 'search' and '1' keys are on
	 * the same column.
	 */
	if ((state[KEYBOARD_COL_LEFT_CTRL] == KEYBOARD_MASK_LEFT_CTRL ||
	     state[KEYBOARD_COL_RIGHT_CTRL] == KEYBOARD_MASK_RIGHT_CTRL) &&
	    ((state[KEYBOARD_COL_SEARCH] & KEYBOARD_MASK_SEARCH) ==
						KEYBOARD_MASK_SEARCH) &&
	    chipset_in_state(CHIPSET_STATE_ON)) {
		if (state[KEYBOARD_COL_KEY_0] == KEYBOARD_MASK_KEY_0) {
			/* Charge from neither port */
			chg_override = -2;
			pd_host_command(EC_CMD_PD_CHARGE_PORT_OVERRIDE, 0,
					&chg_override, 2, NULL, 0);
			return 0;
		} else if (state[KEYBOARD_COL_KEY_1] ==
			   (KEYBOARD_MASK_KEY_1 | KEYBOARD_MASK_SEARCH)) {
			/* Charge from port 0 (left side) */
			chg_override = 0;
			pd_host_command(EC_CMD_PD_CHARGE_PORT_OVERRIDE, 0,
					&chg_override, 2, NULL, 0);
			return 0;
		} else if (state[KEYBOARD_COL_KEY_2] == KEYBOARD_MASK_KEY_2) {
			/* Charge from port 1 (right side) */
			chg_override = 1;
			pd_host_command(EC_CMD_PD_CHARGE_PORT_OVERRIDE, 0,
					&chg_override, 2, NULL, 0);
			return 0;
		}
	}
#endif

	/*
	 * All runtime key combos are (right or left ) alt + volume up + (some
	 * key NOT on the same col as alt or volume up )
	 */
	if (state[KEYBOARD_COL_VOL_UP] != KEYBOARD_MASK_VOL_UP)
		return 0;

	if (state[KEYBOARD_COL_RIGHT_ALT] != KEYBOARD_MASK_RIGHT_ALT &&
	    state[KEYBOARD_COL_LEFT_ALT] != KEYBOARD_MASK_LEFT_ALT)
		return 0;

	/*
	 * Count number of columns with keys pressed.  We know two columns are
	 * pressed for volume up and alt, so if only one more key is pressed
	 * there will be exactly 3 non-zero columns.
	 */
	for (c = 0; c < keyboard_cols; c++) {
		if (state[c])
			num_press++;
	}

	if (num_press != 3)
		return 0;

	/* Check individual keys */
	if (state[KEYBOARD_COL_KEY_R] == KEYBOARD_MASK_KEY_R) {
		/* R = reboot */
		CPRINTS("KB warm reboot");
		keyboard_clear_buffer();
		chipset_reset(CHIPSET_RESET_KB_WARM_REBOOT);
		return 1;
	} else if (state[KEYBOARD_COL_KEY_H] == KEYBOARD_MASK_KEY_H) {
		/* H = hibernate */
		CPRINTS("KB hibernate");
		system_hibernate(0, 0);
		return 1;
	}

	return 0;
}
#endif /* CONFIG_KEYBOARD_RUNTIME_KEYS */

/**
 * Check for ghosting in the keyboard state.
 *
 * Assumes that the state has already been masked with the actual key mask, so
 * that coords which don't correspond with actual keys don't trigger ghosting
 * detection.
 *
 * @param state		Keyboard state to check.
 *
 * @return 1 if ghosting detected, else 0.
 */
static int has_ghosting(const uint8_t *state)
{
	int c, c2;

	for (c = 0; c < keyboard_cols; c++) {
		if (!state[c])
			continue;

		for (c2 = c + 1; c2 < keyboard_cols; c2++) {
			/*
			 * A little bit of cleverness here.  Ghosting happens
			 * if 2 columns share at least 2 keys.  So we OR the
			 * columns together and then see if more than one bit
			 * is set.  x&(x-1) is non-zero only if x has more than
			 * one bit set.
			 */
			uint8_t common = state[c] & state[c2];

			if (common & (common - 1))
				return 1;
		}
	}

	return 0;
}

/**
 * Update keyboard state using low-level interface to read keyboard.
 *
 * @param state		Keyboard state to update.
 *
 * @return 1 if any key is still pressed, 0 if no key is pressed.
 */
static int check_keys_changed(uint8_t *state)
{
	int any_pressed = 0;
	int c, i;
	int any_change = 0;
	static uint8_t __bss_slow new_state[KEYBOARD_COLS_MAX];
	uint32_t tnow = get_time().le.lo;

	/* Save the current scan time */
	if (++scan_time_index >= SCAN_TIME_COUNT)
		scan_time_index = 0;
	scan_time[scan_time_index] = tnow;

	/* Read the raw key state */
	any_pressed = read_matrix(new_state);

	/* Ignore if so many keys are pressed that we're ghosting. */
	if (has_ghosting(new_state))
		return any_pressed;

	/* Check for changes between previous scan and this one */
	for (c = 0; c < keyboard_cols; c++) {
		int diff;

		/* Clear debouncing flag, if sufficient time has elapsed. */
		for (i = 0; i < KEYBOARD_ROWS && debouncing[c]; i++) {
			if (!(debouncing[c] & BIT(i)))
				continue;
			if (tnow - scan_time[scan_edge_index[c][i]] <
			    (state[c] ? keyscan_config.debounce_down_us :
					keyscan_config.debounce_up_us))
				continue;  /* Not done debouncing */
			debouncing[c] &= ~BIT(i);
		}

		/* Recognize change in state, unless debounce in effect. */
		diff = (new_state[c] ^ state[c]) & ~debouncing[c];
		if (!diff)
			continue;
		for (i = 0; i < KEYBOARD_ROWS; i++) {
			if (!(diff & BIT(i)))
				continue;
			scan_edge_index[c][i] = scan_time_index;
			any_change = 1;

			/* Inform keyboard module if scanning is enabled */
			if (keyboard_scan_is_enabled()) {
				/* This is no-op for protocols that require a
				 * full keyboard matrix (e.g., MKBP).
				 */
				keyboard_state_changed(
					i, c, !!(new_state[c] & BIT(i)));
			}
		}

		/* For any keyboard events just sent, turn on debouncing. */
		debouncing[c] |= diff;
		/*
		 * Note: In order to "remember" what was last reported
		 * (up or down), the state bits are only updated if the
		 * edge was not suppressed due to debouncing.
		 */
		state[c] ^= diff;
	}

	if (any_change) {

#ifdef CONFIG_KEYBOARD_SUPPRESS_NOISE
		/* Suppress keyboard noise */
		keyboard_suppress_noise();
#endif

		if (print_state_changes)
			print_state(state, "state");

#ifdef CONFIG_KEYBOARD_PRINT_SCAN_TIMES
		/* Print delta times from now back to each previous scan */
		CPRINTF("[%T kb deltaT");
		for (i = 0; i < SCAN_TIME_COUNT; i++) {
			int tnew = scan_time[
				(SCAN_TIME_COUNT + scan_time_index - i) %
				SCAN_TIME_COUNT];
			CPRINTF(" %d", tnow - tnew);
		}
		CPRINTF("]\n");
#endif

#ifdef CONFIG_KEYBOARD_RUNTIME_KEYS
		/* Swallow special keys */
		if (check_runtime_keys(state))
			return 0;
#endif

#ifdef CONFIG_KEYBOARD_PROTOCOL_MKBP
		keyboard_fifo_add(state);
#endif
	}

	kbd_polls++;

	return any_pressed;
}

#ifdef CONFIG_KEYBOARD_BOOT_KEYS
/*
 * Returns mask of the boot keys that are pressed, with at most the keys used
 * for keyboard-controlled reset also pressed.
 */
static uint32_t check_key_list(const uint8_t *state)
{
	uint8_t curr_state[KEYBOARD_COLS_MAX];
	int c;
	uint32_t boot_key_mask = BOOT_KEY_NONE;
	const struct boot_key_entry *k;

	/* Make copy of current debounced state. */
	memcpy(curr_state, state, sizeof(curr_state));

#ifdef KEYBOARD_MASK_PWRBTN
	/*
	 * Check if KSI2 or KSI3 is asserted for all columns due to power
	 * button hold, and ignore it if so.
	 */
	for (c = 0; c < keyboard_cols; c++)
		if ((keyscan_config.actual_key_mask[c] & KEYBOARD_MASK_PWRBTN)
		    && !(curr_state[c] & KEYBOARD_MASK_PWRBTN))
			break;

	if (c == keyboard_cols)
		for (c = 0; c < keyboard_cols; c++)
			curr_state[c] &= ~KEYBOARD_MASK_PWRBTN;
#endif

	curr_state[KEYBOARD_COL_REFRESH] &= ~KEYBOARD_MASK_REFRESH;

	/* Update mask with all boot keys that were pressed. */
	k = boot_key_list;
	for (c = 0; c < ARRAY_SIZE(boot_key_list); c++, k++) {
		if (curr_state[k->mask_index] & k->mask_value) {
			boot_key_mask |= BIT(c);
			curr_state[k->mask_index] &= ~k->mask_value;
		}
	}

	/* If any other key was pressed, ignore all boot keys. */
	for (c = 0; c < keyboard_cols; c++) {
		if (curr_state[c])
			return BOOT_KEY_NONE;
	}

	CPRINTS("KB boot key mask %x", boot_key_mask);
	return boot_key_mask;
}

/**
 * Check what boot key is down, if any.
 *
 * @param state		Keyboard state at boot.
 *
 * @return the key which is down, or BOOT_KEY_NONE if an unrecognized
 * key combination is down or this isn't the right type of boot to look at
 * boot keys.
 */
static uint32_t check_boot_key(const uint8_t *state)
{
	/*
	 * If we jumped to this image, ignore boot keys.  This prevents
	 * re-triggering events in RW firmware that were already processed by
	 * RO firmware.
	 */
	if (system_jumped_to_this_image())
		return BOOT_KEY_NONE;

	/* If reset was not caused by reset pin, refresh must be held down */
	if (!(system_get_reset_flags() & EC_RESET_FLAG_RESET_PIN) &&
	    !(state[KEYBOARD_COL_REFRESH] & KEYBOARD_MASK_REFRESH))
		return BOOT_KEY_NONE;

	return check_key_list(state);
}
#endif

static void keyboard_freq_change(void)
{
	post_scan_clock_us = (CONFIG_KEYBOARD_POST_SCAN_CLOCKS * 1000) /
		(clock_get_freq() / 1000);
}
DECLARE_HOOK(HOOK_FREQ_CHANGE, keyboard_freq_change, HOOK_PRIO_DEFAULT);

/*****************************************************************************/
/* Interface */

struct keyboard_scan_config *keyboard_scan_get_config(void)
{
	return &keyscan_config;
}

#ifdef CONFIG_KEYBOARD_BOOT_KEYS
uint32_t keyboard_scan_get_boot_keys(void)
{
	return boot_key_value;
}
#endif

const uint8_t *keyboard_scan_get_state(void)
{
	return debounced_state;
}

void keyboard_scan_init(void)
{
	/* Configure GPIO */
	keyboard_raw_init();

	/* Tri-state the columns */
	keyboard_raw_drive_column(KEYBOARD_COLUMN_NONE);

	/* Initialize raw state */
	read_matrix(debounced_state);

#ifdef CONFIG_KEYBOARD_LANGUAGE_ID
	/* Check keyboard ID state */
	read_matrix_id(keyboard_id);
#endif

#ifdef CONFIG_KEYBOARD_BOOT_KEYS
	/* Check for keys held down at boot */
	boot_key_value = check_boot_key(debounced_state);

	/*
	 * If any key other than Esc or Left_Shift was pressed, do not trigger
	 * recovery.
	 */
	if (boot_key_value & ~(BOOT_KEY_ESC | BOOT_KEY_LEFT_SHIFT))
		return;

#ifdef CONFIG_HOSTCMD_EVENTS
	if (boot_key_value & BOOT_KEY_ESC) {
		host_set_single_event(EC_HOST_EVENT_KEYBOARD_RECOVERY);
		/*
		 * In recovery mode, we should force clamshell mode in order to
		 * prevent the keyboard from being disabled unintentionally due
		 * to unstable accel readings.
		 *
		 * You get the same effect if motion sensors or a motion sense
		 * task are disabled in RO.
		 */
		if (IS_ENABLED(CONFIG_TABLET_MODE))
			tablet_disable();
		if (boot_key_value & BOOT_KEY_LEFT_SHIFT)
			host_set_single_event(
				EC_HOST_EVENT_KEYBOARD_RECOVERY_HW_REINIT);
	}
#endif
#endif /* CONFIG_KEYBOARD_BOOT_KEYS */
}

void keyboard_scan_task(void *u)
{
	timestamp_t poll_deadline, start;
	int wait_time;
	uint32_t local_disable_scanning = 0;

	print_state(debounced_state, "init state");

	keyboard_raw_task_start();

	/* Set initial clock frequency-based minimum delay between scans */
	keyboard_freq_change();

	while (1) {
		/* Enable all outputs */
		CPRINTS5("KB wait");

		keyboard_raw_enable_interrupt(1);

		/* Wait for scanning enabled and key pressed. */
		while (1) {
			uint32_t new_disable_scanning;

			/* Read it once to get consistent glimpse */
			new_disable_scanning = disable_scanning_mask;

			if (local_disable_scanning != new_disable_scanning)
				CPRINTS("KB disable_scanning_mask changed: "
					"0x%08x", new_disable_scanning);

			if (!new_disable_scanning) {
				/* Enabled now */
				keyboard_raw_drive_column(KEYBOARD_COLUMN_ALL);
			} else if (!local_disable_scanning) {
				/*
				 * Scanning isn't enabled but it was last time
				 * we looked.
				 *
				 * No race here even though we're basing on a
				 * glimpse of disable_scanning_mask since if
				 * someone changes disable_scanning_mask they
				 * are guaranteed to call task_wake() on us
				 * afterward so we'll run the loop again.
				 */
				keyboard_raw_drive_column(KEYBOARD_COLUMN_NONE);
				keyboard_clear_buffer();
			}

			local_disable_scanning = new_disable_scanning;

			/*
			 * Done waiting if scanning is enabled and a key is
			 * already pressed.  This prevents a race between the
			 * user pressing a key and enable_interrupt()
			 * starting to pay attention to edges.
			 */
			if (!local_disable_scanning &&
			    (keyboard_raw_read_rows() || force_poll))
				break;
			else
				task_wait_event(-1);
		}

		/* We're about to poll, so any existing forces are fulfilled */
		force_poll = 0;

		/* Enter polling mode */
		CPRINTS5("KB poll");
		keyboard_raw_enable_interrupt(0);
		keyboard_raw_drive_column(KEYBOARD_COLUMN_NONE);

		/* Busy polling keyboard state. */
		while (keyboard_scan_is_enabled()) {
			start = get_time();

			/* Check for keys down */
			if (check_keys_changed(debounced_state)) {
				poll_deadline.val = start.val
					+ keyscan_config.poll_timeout_us;
			} else if (timestamp_expired(poll_deadline, &start)) {
				break;
			}

			/* Delay between scans */
			wait_time = keyscan_config.scan_period_us -
				(get_time().val - start.val);

			if (wait_time < keyscan_config.min_post_scan_delay_us)
				wait_time =
					keyscan_config.min_post_scan_delay_us;

			if (wait_time < post_scan_clock_us)
				wait_time = post_scan_clock_us;

			usleep(wait_time);
		}
	}
}

#ifdef CONFIG_LID_SWITCH

static void keyboard_lid_change(void)
{
	if (lid_is_open())
		keyboard_scan_enable(1, KB_SCAN_DISABLE_LID_CLOSED);
	else
		keyboard_scan_enable(0, KB_SCAN_DISABLE_LID_CLOSED);
}
DECLARE_HOOK(HOOK_LID_CHANGE, keyboard_lid_change, HOOK_PRIO_DEFAULT);
DECLARE_HOOK(HOOK_INIT, keyboard_lid_change, HOOK_PRIO_INIT_LID + 1);

#endif

#ifdef CONFIG_USB_SUSPEND
static void keyboard_usb_pm_change(void)
{
	/*
	 * If USB interface is suspended, and host is not asking us to do remote
	 * wakeup, we can turn off the key scanning.
	 */
	if (usb_is_suspended() && !usb_is_remote_wakeup_enabled())
		keyboard_scan_enable(0, KB_SCAN_DISABLE_USB_SUSPENDED);
	else
		keyboard_scan_enable(1, KB_SCAN_DISABLE_USB_SUSPENDED);
}
DECLARE_HOOK(HOOK_USB_PM_CHANGE, keyboard_usb_pm_change, HOOK_PRIO_DEFAULT);
#endif

/*****************************************************************************/
/* Host commands */

static enum ec_status
mkbp_command_simulate_key(struct host_cmd_handler_args *args)
{
	const struct ec_params_mkbp_simulate_key *p = args->params;

	/* Only available on unlocked systems */
	if (system_is_locked())
		return EC_RES_ACCESS_DENIED;

	if (p->col >= keyboard_cols || p->row >= KEYBOARD_ROWS)
		return EC_RES_INVALID_PARAM;

	simulate_key(p->row, p->col, p->pressed);

	return EC_RES_SUCCESS;
}
DECLARE_HOST_COMMAND(EC_CMD_MKBP_SIMULATE_KEY,
		     mkbp_command_simulate_key,
		     EC_VER_MASK(0));

#ifdef CONFIG_KEYBOARD_FACTORY_TEST

/* Run keyboard factory testing, scan out KSO/KSI if any shorted. */
int keyboard_factory_test_scan(void)
{
	int i, j, flags;
	uint16_t shorted = 0;
	int port, id;

	/* Disable keyboard scan while testing */
	keyboard_scan_enable(0, KB_SCAN_DISABLE_LID_CLOSED);
	flags = gpio_get_default_flags(GPIO_KBD_KSO2);

	/* Set all of KSO/KSI pins to internal pull-up and input */
	for (i = 0; i < keyboard_factory_scan_pins_used; i++) {

		if (keyboard_factory_scan_pins[i][0] < 0)
			continue;

		port = keyboard_factory_scan_pins[i][0];
		id = keyboard_factory_scan_pins[i][1];

		gpio_set_alternate_function(port, 1 << id, -1);
		gpio_set_flags_by_mask(port, 1 << id,
			GPIO_INPUT | GPIO_PULL_UP);
	}

	/*
	 * Set start pin to output low, then check other pins
	 * going to low level, it indicate the two pins are shorted.
	 */
	for (i = 0; i < keyboard_factory_scan_pins_used; i++) {

		if (keyboard_factory_scan_pins[i][0] < 0)
			continue;

		port = keyboard_factory_scan_pins[i][0];
		id = keyboard_factory_scan_pins[i][1];

		gpio_set_flags_by_mask(port, 1 << id, GPIO_OUT_LOW);

		for (j = 0; j < i; j++) {

			if (keyboard_factory_scan_pins[j][0] < 0)
				continue;

			if (keyboard_raw_is_input_low(
					keyboard_factory_scan_pins[j][0],
					keyboard_factory_scan_pins[j][1])) {
				shorted = i << 8 | j;
				goto done;
			}
		}
		gpio_set_flags_by_mask(port, 1 << id,
			GPIO_INPUT | GPIO_PULL_UP);
	}
done:
	gpio_config_module(MODULE_KEYBOARD_SCAN, 1);
	gpio_set_flags(GPIO_KBD_KSO2, flags);
	keyboard_scan_enable(1, KB_SCAN_DISABLE_LID_CLOSED);

	return shorted;
}

static enum ec_status keyboard_factory_test(struct host_cmd_handler_args *args)
{
	struct ec_response_keyboard_factory_test *r = args->response;

	/* Only available on unlocked systems */
	if (system_is_locked())
		return EC_RES_ACCESS_DENIED;

	if (keyboard_factory_scan_pins_used == 0)
		return EC_RES_INVALID_COMMAND;

	r->shorted = keyboard_factory_test_scan();

	args->response_size = sizeof(*r);

	return EC_RES_SUCCESS;
}

DECLARE_HOST_COMMAND(EC_CMD_KEYBOARD_FACTORY_TEST,
		     keyboard_factory_test,
		     EC_VER_MASK(0));
#endif

#ifdef CONFIG_KEYBOARD_LANGUAGE_ID
int keyboard_get_keyboard_id(void)
{
	int c;
	uint32_t id = 0;

	BUILD_ASSERT(sizeof(id) >= KEYBOARD_IDS);

	for (c = 0; c < KEYBOARD_IDS; c++) {
		/* Check ID ghosting if more than one bit in any KSIs was set */
		if (keyboard_id[c] & (keyboard_id[c] - 1))
			/* ID ghosting is found */
			return KEYBOARD_ID_UNREADABLE;
		else
			id |= keyboard_id[c] << (c * 8);
	}
	return id;
}
#endif

/*****************************************************************************/
/* Console commands */
#ifdef CONFIG_CMD_KEYBOARD
static int command_ksstate(int argc, char **argv)
{
	if (argc > 1) {
		if (!strcasecmp(argv[1], "force")) {
			print_state_changes = 1;
			keyboard_scan_enable(1, -1);
		} else if (!parse_bool(argv[1], &print_state_changes)) {
			return EC_ERROR_PARAM1;
		}
	}

	print_state(debounced_state, "debounced ");
	print_state(debouncing, "debouncing");

	ccprintf("Keyboard scan disable mask: 0x%08x\n",
		 disable_scanning_mask);
	ccprintf("Keyboard scan state printing %s\n",
		 print_state_changes ? "on" : "off");
	return EC_SUCCESS;
}
DECLARE_CONSOLE_COMMAND(ksstate, command_ksstate,
			"ksstate [on | off | force]",
			"Show or toggle printing keyboard scan state");

static int command_keyboard_press(int argc, char **argv)
{
	if (argc == 1) {
		int i, j;

		ccputs("Simulated keys:\n");
		for (i = 0; i < keyboard_cols; ++i) {
			if (simulated_key[i] == 0)
				continue;
			for (j = 0; j < KEYBOARD_ROWS; ++j)
				if (simulated_key[i] & BIT(j))
					ccprintf("\t%d %d\n", i, j);
		}

	} else if (argc == 3 || argc == 4) {
		int r, c, p;
		char *e;

		c = strtoi(argv[1], &e, 0);
		if (*e || c < 0 || c >= keyboard_cols)
			return EC_ERROR_PARAM1;

		r = strtoi(argv[2], &e, 0);
		if (*e || r < 0 || r >= KEYBOARD_ROWS)
			return EC_ERROR_PARAM2;

		if (argc == 3) {
			/* Simulate a press and release */
			simulate_key(r, c, 1);
			simulate_key(r, c, 0);
		} else {
			p = strtoi(argv[3], &e, 0);
			if (*e || p < 0 || p > 1)
				return EC_ERROR_PARAM3;

			simulate_key(r, c, p);
		}
	}

	return EC_SUCCESS;
}
DECLARE_CONSOLE_COMMAND(kbpress, command_keyboard_press,
			"[col row [0 | 1]]",
			"Simulate keypress");
#endif

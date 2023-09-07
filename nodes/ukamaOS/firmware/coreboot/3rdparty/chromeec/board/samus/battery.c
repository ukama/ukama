/* Copyright 2012 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Battery pack vendor provided charging profile
 */

#include "battery_smart.h"
#include "charge_state.h"
#include "console.h"
#include "ec_commands.h"
#include "extpower.h"
#include "i2c.h"
#include "util.h"

static const struct battery_info info = {
	/*
	 * Design voltage
	 *   max    = 8.4V
	 *   normal = 7.4V
	 *   min    = 6.0V
	 */
	.voltage_max    = 8700,
	.voltage_normal = 7400,
	.voltage_min    = 6000,

	/* Pre-charge current: I <= 0.01C */
	.precharge_current  = 64, /* mA */

	/*
	 * Operational temperature range
	 *   0 <= T_charge    <= 50 deg C
	 * -20 <= T_discharge <= 60 deg C
	 */
	.start_charging_min_c = 0,
	.start_charging_max_c = 50,
	.charging_min_c       = 0,
	.charging_max_c       = 50,
	.discharging_min_c    = -20,
	.discharging_max_c    = 60,
};

const struct battery_info *battery_get_info(void)
{
	return &info;
}

#ifdef CONFIG_CHARGER_PROFILE_OVERRIDE

static int fast_charging_allowed = 1;

/*
 * This can override the smart battery's charging profile. To make a change,
 * modify one or more of requested_voltage, requested_current, or state.
 * Leave everything else unchanged.
 *
 * Return the next poll period in usec, or zero to use the default (which is
 * state dependent).
 */
int charger_profile_override(struct charge_state_data *curr)
{
	/* temp in 0.1 deg C */
	int temp_c;
	const struct charger_info *info;

	/* keep track of last temperature range for hysteresis */
	static enum {
		TEMP_LOW,
		TEMP_NORMAL,
		TEMP_HIGH
	} temp_range = TEMP_NORMAL, prev_temp_range = TEMP_NORMAL;

	/* charging voltage to use at high temp */
	static int high_temp_charging_voltage;

	/* custom profile phase at normal temp */
	static int normal_temp_phase;

	/* battery voltage and current and previous voltage and current */
	int batt_voltage, batt_current;
	static int prev_batt_voltage, prev_batt_current;

	/*
	 * Determine temperature range:
	 * Low: Battery is <15C
	 * Normal: Battery is 15-45C
	 * High: Battery is >45C
	 *
	 * Add 0.2 degrees of hysteresis.
	 * If temp reading was bad use last range.
	 */
	if (!(curr->batt.flags & BATT_FLAG_BAD_TEMPERATURE)) {
		temp_c = curr->batt.temperature - 2731;
		if (temp_c < 149)
			temp_range = TEMP_LOW;
		else if (temp_c > 151 && temp_c < 449)
			temp_range = TEMP_NORMAL;
		else if (temp_c > 451)
			temp_range = TEMP_HIGH;
	}

	/*
	 * Treat voltage and current as a pair, if either is bad fall back to
	 * previous reading.
	 */
	if (curr->batt.flags &
	    (BATT_FLAG_BAD_VOLTAGE | BATT_FLAG_BAD_CURRENT)) {
		batt_voltage = prev_batt_voltage;
		batt_current = prev_batt_current;
	} else {
		batt_voltage = prev_batt_voltage = curr->batt.voltage;
		batt_current = prev_batt_current = curr->batt.current;
	}

	/*
	 * If we are not charging or we aren't using fast charging profiles,
	 * then do not override desired current and voltage and reset some
	 * fast charging profile static variables.
	 */
	if (curr->state != ST_CHARGE || !fast_charging_allowed) {
		prev_temp_range = TEMP_NORMAL;
		normal_temp_phase = 0;
		return 0;
	}

	/*
	 * Okay, impose our custom will:
	 * Normal temp:
	 * Phase 0: CC at 9515mA @ 8.3V
	 *          CV at 8.3V until current drops to 4759mA
	 * Phase 1: CC at 4759mA @ 8.7V
	 *          CV at 8.7V
	 *
	 * Low temp:
	 * CC at 2854mA @ 8.7V
	 * CV at 8.7V
	 *
	 * High temp:
	 * If battery voltage < 8.3V then:
	 *	CC at 6660mA @ 8.3V
	 *	CV at 8.3V (when battery is hot we don't go to fully charged)
	 * else:
	 *	CV at just above battery voltage which will essentially
	 *         terminate the charge and allow battery to cool.
	 * Note that if we ever request a voltage below the present battery
	 * voltage, then we will stop the BQ switching, which will power off
	 * the INA and we won't be able to charge again until AC is
	 * disconnected. see crbug.com/p/35491.
	 */
	switch (temp_range) {
	case TEMP_LOW:
		curr->requested_current = 2854;
		curr->requested_voltage = 8700;
		break;
	case TEMP_NORMAL:
		if (normal_temp_phase == 0) {
			curr->requested_current = 9515;
			curr->requested_voltage = 8300;
			if (batt_current <= 4759 && batt_voltage >= 8200)
				normal_temp_phase = 1;
		}
		if (normal_temp_phase == 1) {
			curr->requested_current = 4759;
			curr->requested_voltage = 8700;
		}
		break;
	case TEMP_HIGH:
		/*
		 * First time TEMP_HIGH is used, get the closest voltage
		 * just above the battery voltage. If it is above 8.3V, we
		 * will use that as the target, otherwise we will use 8.3V.
		 */
		if (prev_temp_range != TEMP_HIGH) {
			info = charger_get_info();
			high_temp_charging_voltage = MAX(8300,
				charger_closest_voltage(batt_voltage +
							info->voltage_step));
		}
		curr->requested_current = 6660;
		curr->requested_voltage = high_temp_charging_voltage;
		break;
	}
	prev_temp_range = temp_range;

	return 0;
}

/* Customs options controllable by host command. */
#define PARAM_FASTCHARGE (CS_PARAM_CUSTOM_PROFILE_MIN + 0)

enum ec_status charger_profile_override_get_param(uint32_t param,
						  uint32_t *value)
{
	if (param == PARAM_FASTCHARGE) {
		*value = fast_charging_allowed;
		return EC_RES_SUCCESS;
	}
	return EC_RES_INVALID_PARAM;
}

enum ec_status charger_profile_override_set_param(uint32_t param,
						  uint32_t value)
{
	if (param == PARAM_FASTCHARGE) {
		fast_charging_allowed = value;
		return EC_RES_SUCCESS;
	}
	return EC_RES_INVALID_PARAM;
}

#ifdef CONFIG_CMD_FASTCHARGE
static int command_fastcharge(int argc, char **argv)
{
	if (argc > 1 && !parse_bool(argv[1], &fast_charging_allowed))
		return EC_ERROR_PARAM1;

	ccprintf("fastcharge %s\n", fast_charging_allowed ? "on" : "off");

	return EC_SUCCESS;
}
DECLARE_CONSOLE_COMMAND(fastcharge, command_fastcharge,
			"[on|off]",
			"Get or set fast charging profile");
#endif  /* CONFIG_CMD_FASTCHARGE */

#endif	/* CONFIG_CHARGER_PROFILE_OVERRIDE */

#ifdef CONFIG_BATTERY_REVIVE_DISCONNECT
/*
 * Check if battery is in disconnect state, a state entered by pulling
 * BATT_DISCONN_N low, and clear that state if we have external power plugged
 * and no battery faults are detected. Disconnect state resembles battery
 * shutdown mode, but extra steps must be taken to get the battery out of this
 * mode.
 */
enum battery_disconnect_state battery_get_disconnect_state(void)
{
	uint8_t data[6];
	int rv;
	/*
	 * Take note if we find that the battery isn't in disconnect state,
	 * and always return NOT_DISCONNECTED without probing the battery.
	 * This assumes the battery will not go to disconnect state during
	 * runtime.
	 */
	static int not_disconnected;

	if (not_disconnected)
		return BATTERY_NOT_DISCONNECTED;

	if (extpower_is_present()) {
		/* Check if battery charging + discharging is disabled. */
		rv = sb_read_mfgacc(PARAM_OPERATION_STATUS,
				SB_ALT_MANUFACTURER_ACCESS, data, sizeof(data));
		if (rv)
			return BATTERY_DISCONNECT_ERROR;
		if (~data[3] & (BATTERY_DISCHARGING_DISABLED |
				BATTERY_CHARGING_DISABLED)) {
			not_disconnected = 1;
			return BATTERY_NOT_DISCONNECTED;
		}

		/*
		 * Battery is neither charging nor discharging. Verify that
		 * we didn't enter this state due to a safety fault.
		 */
		rv = sb_read_mfgacc(PARAM_SAFETY_STATUS,
				SB_ALT_MANUFACTURER_ACCESS, data, sizeof(data));
		if (rv || data[2] || data[3] || data[4] || data[5])
			return BATTERY_DISCONNECT_ERROR;
		else
			/* No safety fault -- clear disconnect state. */
			return BATTERY_DISCONNECTED;
	}
	not_disconnected = 1;
	return BATTERY_NOT_DISCONNECTED;
}
#endif /* CONFIG_BATTERY_REVIVE_DISCONNECT */

#define PARAM_CUT_OFF_LOW  0x10
#define PARAM_CUT_OFF_HIGH 0x00

int board_cut_off_battery(void)
{
	int rv;
	uint8_t buf[3];

	buf[0] = SB_MANUFACTURER_ACCESS & 0xff;
	buf[1] = PARAM_CUT_OFF_LOW;
	buf[2] = PARAM_CUT_OFF_HIGH;

	i2c_lock(I2C_PORT_BATTERY, 1);
	rv = i2c_xfer_unlocked(I2C_PORT_BATTERY, BATTERY_ADDR_FLAGS,
			       buf, 3, NULL, 0, I2C_XFER_SINGLE);
	rv |= i2c_xfer_unlocked(I2C_PORT_BATTERY, BATTERY_ADDR_FLAGS,
				buf, 3, NULL, 0, I2C_XFER_SINGLE);
	i2c_lock(I2C_PORT_BATTERY, 0);

	return rv;
}


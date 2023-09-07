/* Copyright 2018 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * TI bq25710 battery charger driver.
 */

#include "battery.h"
#include "battery_smart.h"
#include "bq25710.h"
#include "charge_ramp.h"
#include "charger.h"
#include "common.h"
#include "console.h"
#include "hooks.h"
#include "i2c.h"
#include "task.h"
#include "timer.h"

#ifndef CONFIG_CHARGER_NARROW_VDC
#error "BQ25710 is a NVDC charger, please enable CONFIG_CHARGER_NARROW_VDC."
#endif

/* Sense resistor configurations and macros */
#define DEFAULT_SENSE_RESISTOR 10

#define INPUT_RESISTOR_RATIO \
	((CONFIG_CHARGER_SENSE_RESISTOR_AC) / DEFAULT_SENSE_RESISTOR)
#define REG_TO_INPUT_CURRENT(REG) ((REG + 1) * 50 / INPUT_RESISTOR_RATIO)
#define INPUT_CURRENT_TO_REG(CUR) (((CUR) * INPUT_RESISTOR_RATIO / 50) - 1)

#define CHARGING_RESISTOR_RATIO \
	((CONFIG_CHARGER_SENSE_RESISTOR) / DEFAULT_SENSE_RESISTOR)
#define REG_TO_CHARGING_CURRENT(REG) ((REG) / CHARGING_RESISTOR_RATIO)
#define CHARGING_CURRENT_TO_REG(CUR) ((CUR) * CHARGING_RESISTOR_RATIO)

/* Console output macros */
#define CPRINTF(format, args...) cprintf(CC_CHARGER, format, ## args)

#ifdef CONFIG_CHARGER_BQ25710_IDCHG_LIMIT_MA
/*
 * If this config option is defined, then the bq25710 needs to remain in
 * performance mode when the AP is in S0. Performance mode is active whenever AC
 * power is connected or when the EN_LWPWR bit in ChargeOption0 is clear.
 */
static uint32_t bq25710_perf_mode_req;
static struct mutex bq25710_perf_mode_mutex;
#endif

/* Charger parameters */
static const struct charger_info bq25710_charger_info = {
	.name         = "bq25710",
	.voltage_max  = 19200,
	.voltage_min  = 1024,
	.voltage_step = 16,
	.current_max  = 8128 / CHARGING_RESISTOR_RATIO,
	.current_min  = 64 / CHARGING_RESISTOR_RATIO,
	.current_step = 64 / CHARGING_RESISTOR_RATIO,
	.input_current_max  = 6400 / INPUT_RESISTOR_RATIO,
	.input_current_min  = 50 / INPUT_RESISTOR_RATIO,
	.input_current_step = 50 / INPUT_RESISTOR_RATIO,
};

static inline int raw_read16(int offset, int *value)
{
	return i2c_read16(I2C_PORT_CHARGER, BQ25710_SMBUS_ADDR1_FLAGS,
			  offset, value);
}

static inline int raw_write16(int offset, int value)
{
	return i2c_write16(I2C_PORT_CHARGER, BQ25710_SMBUS_ADDR1_FLAGS,
			   offset, value);
}

#if defined(CONFIG_CHARGE_RAMP_HW) || \
	defined(CONFIG_USB_PD_VBUS_MEASURE_CHARGER)
static int bq25710_get_low_power_mode(int *mode)
{
	int rv;
	int reg;

	rv = raw_read16(BQ25710_REG_CHARGE_OPTION_0, &reg);
	if (rv)
		return rv;

	*mode = !!(reg & BQ25710_CHARGE_OPTION_0_LOW_POWER_MODE);

	return EC_SUCCESS;
}

static int bq25710_set_low_power_mode(int enable)
{
	int rv;
	int reg;

	rv = raw_read16(BQ25710_REG_CHARGE_OPTION_0, &reg);
	if (rv)
		return rv;

#ifdef CONFIG_CHARGER_BQ25710_IDCHG_LIMIT_MA
	mutex_lock(&bq25710_perf_mode_mutex);
	/*
	 * Performance mode means not in low power mode. The bit that controls
	 * this is EN_LWPWR in ChargeOption0. The 'enable' param in this
	 * function is refeerring to low power mode, so enabling low power mode
	 * means disabling performance mode and vice versa.
	 */
	if (enable)
		bq25710_perf_mode_req &= ~(1 << task_get_current());
	else
		bq25710_perf_mode_req |= (1 << task_get_current());
	enable = !bq25710_perf_mode_req;
#endif

	if (enable)
		reg |= BQ25710_CHARGE_OPTION_0_LOW_POWER_MODE;
	else
		reg &= ~BQ25710_CHARGE_OPTION_0_LOW_POWER_MODE;

	rv = raw_write16(BQ25710_REG_CHARGE_OPTION_0, reg);
#ifdef CONFIG_CHARGER_BQ25710_IDCHG_LIMIT_MA
	mutex_unlock(&bq25710_perf_mode_mutex);
#endif
	if (rv)
		return rv;

	return EC_SUCCESS;
}

static int bq25710_adc_start(int adc_en_mask)
{
	int reg;
	int mode;
	int tries_left = 8;

	/* Save current mode to restore same state after ADC read */
	if (bq25710_get_low_power_mode(&mode))
		return EC_ERROR_UNKNOWN;

	/* Exit low power mode so ADC conversion takes typical time */
	if (bq25710_set_low_power_mode(0))
		return EC_ERROR_UNKNOWN;

	/*
	 * Turn on the ADC for one reading. Note that adc_en_mask
	 * maps to bit[7:0] in ADCOption register.
	 */
	reg = (adc_en_mask & BQ25710_ADC_OPTION_EN_ADC_ALL) |
	      BQ25710_ADC_OPTION_ADC_START;
	if (raw_write16(BQ25710_REG_ADC_OPTION, reg))
		return EC_ERROR_UNKNOWN;

	/*
	 * Wait until the ADC operation completes. The spec says typical
	 * conversion time is 10 msec. If low power mode isn't exited first,
	 * then the conversion time jumps to ~60 msec.
	 */
	do {
		msleep(2);
		raw_read16(BQ25710_REG_ADC_OPTION, &reg);
	} while (--tries_left && (reg & BQ25710_ADC_OPTION_ADC_START));

	/* ADC reading attempt complete, go back to low power mode */
	if (bq25710_set_low_power_mode(mode))
		return EC_ERROR_UNKNOWN;

	/* Could not complete read */
	if (reg & BQ25710_ADC_OPTION_ADC_START)
		return EC_ERROR_TIMEOUT;

	return EC_SUCCESS;
}
#endif

static void bq25710_init(void)
{
	int reg;
	int vsys;
	int rv;

	/*
	 * Reset registers to their default settings. There is no reset pin for
	 * this chip so without a full power cycle, some registers may not be at
	 * their default values. Note, need to save the POR value of
	 * MIN_SYSTEM_VOLTAGE register prior to setting the reset so that the
	 * correct value is preserved.
	 */
	rv = raw_read16(BQ25710_REG_MIN_SYSTEM_VOLTAGE, &vsys);
	rv |= raw_read16(BQ25710_REG_CHARGE_OPTION_3, &reg);
	if (!rv) {
		reg |= BQ25710_CHARGE_OPTION_3_RESET_REG;
		/* Set all registers to default values */
		raw_write16(BQ25710_REG_CHARGE_OPTION_3, reg);
		/* Restore VSYS_MIN voltage to POR reset value */
		raw_write16(BQ25710_REG_MIN_SYSTEM_VOLTAGE, vsys);
	}

	if (!raw_read16(BQ25710_REG_PROCHOT_OPTION_1, &reg)) {
		/* Disbale VDPM prochot profile at initialization */
		reg &= ~BQ25710_PROCHOT_PROFILE_VDPM;
		/*
		 * Enable PROCHOT to be asserted with VSYS min detection. Note
		 * that when no battery is present, then VSYS will be set to the
		 * value in register 0x3E (MinSysVoltage) which means that when
		 * no battery is present prochot will continuosly be asserted.
		 */
		reg |= BQ25710_PROCHOT_PROFILE_VSYS;
#ifdef CONFIG_CHARGER_BQ25710_IDCHG_LIMIT_MA
		/*
		 * Set the IDCHG limit who's value is defined in the config
		 * option in mA. Also, enable IDCHG trigger for prochot.
		 */
		reg &= ~BQ25710_PROCHOT_IDCHG_VTH_MASK;
		/*
		 * IDCHG limit is in 512 mA steps. Note there is a 128 mA offset
		 * so the actual IDCHG limit will be the value stored in bits
		 * 15:10 + 128 mA.
		 */
		reg |= ((CONFIG_CHARGER_BQ25710_IDCHG_LIMIT_MA << 1) &
			BQ25710_PROCHOT_IDCHG_VTH_MASK);
		reg |= BQ25710_PROCHOT_PROFILE_IDCHG;
#endif
		raw_write16(BQ25710_REG_PROCHOT_OPTION_1, reg);
	}

	/* Reduce ILIM from default of 150% to 105% */
	if (!raw_read16(BQ25710_REG_PROCHOT_OPTION_0, &reg)) {
		reg &= ~BQ25710_PROCHOT0_ILIM_VTH_MASK;
		raw_write16(BQ25710_REG_PROCHOT_OPTION_0, reg);
	}

	/*
	 * Reduce peak power mode overload and relax cycle time from default 20
	 * msec to the minimum of 5 msec.
	 */
	if (!raw_read16(BQ25710_REG_CHARGE_OPTION_2, &reg)) {
		reg &= ~BQ25710_CHARGE_OPTION_2_TMAX_MASK;
		raw_write16(BQ25710_REG_CHARGE_OPTION_2, reg);
	}
}
DECLARE_HOOK(HOOK_INIT, bq25710_init, HOOK_PRIO_INIT_I2C + 1);

/* Charger interfaces */
const struct charger_info *charger_get_info(void)
{
	return &bq25710_charger_info;
}

int charger_post_init(void)
{
	/*
	 * Note: bq25710 power on reset state is:
	 *	watch dog timer     = 175 sec
	 *	input current limit = ~1/2 maximum setting
	 *	charging voltage    = 0 mV
	 *	charging current    = 0 mA
	 *	discharge on AC     = disabled
	 */

	return EC_SUCCESS;
}

int charger_get_status(int *status)
{
	int rv;
	int option;

	rv = charger_get_option(&option);
	if (rv)
		return rv;

	/* Default status */
	*status = CHARGER_LEVEL_2;

	if (option & BQ25710_CHARGE_OPTION_0_CHRG_INHIBIT)
		*status |= CHARGER_CHARGE_INHIBITED;

	return EC_SUCCESS;
}

int charger_set_mode(int mode)
{
	int rv;
	int option;

	rv = charger_get_option(&option);
	if (rv)
		return rv;

	if (mode & CHARGER_CHARGE_INHIBITED)
		option |= BQ25710_CHARGE_OPTION_0_CHRG_INHIBIT;
	else
		option &= ~BQ25710_CHARGE_OPTION_0_CHRG_INHIBIT;

	return charger_set_option(option);
}

int charger_enable_otg_power(int enabled)
{
	/* This is controlled with the EN_OTG pin. Support not added yet. */
	return EC_ERROR_UNIMPLEMENTED;
}

int charger_set_otg_current_voltage(int output_current, int output_voltage)
{
	/* Add when needed. */
	return EC_ERROR_UNIMPLEMENTED;
}

int charger_is_sourcing_otg_power(int port)
{
	/* Add when needed. */
	return EC_ERROR_UNIMPLEMENTED;
}

int charger_get_current(int *current)
{
	int rv, reg;

	rv = raw_read16(BQ25710_REG_CHARGE_CURRENT, &reg);
	if (!rv)
		*current = REG_TO_CHARGING_CURRENT(reg);

	return rv;
}

int charger_set_current(int current)
{
	return raw_write16(BQ25710_REG_CHARGE_CURRENT,
		CHARGING_CURRENT_TO_REG(current));
}

/* Get/set charge voltage limit in mV */
int charger_get_voltage(int *voltage)
{
	return raw_read16(BQ25710_REG_MAX_CHARGE_VOLTAGE, voltage);
}

int charger_set_voltage(int voltage)
{
	return raw_write16(BQ25710_REG_MAX_CHARGE_VOLTAGE, voltage);
}

/* Discharge battery when on AC power. */
int charger_discharge_on_ac(int enable)
{
	int rv, option;

	rv = charger_get_option(&option);
	if (rv)
		return rv;

	if (enable)
		option |= BQ25710_CHARGE_OPTION_0_EN_LEARN;
	else
		option &= ~BQ25710_CHARGE_OPTION_0_EN_LEARN;

	return charger_set_option(option);
}

int charger_set_input_current(int input_current)
{
	int num_steps = INPUT_CURRENT_TO_REG(input_current);

	return raw_write16(BQ25710_REG_IIN_HOST, num_steps <<
			  BQ25710_CHARGE_IIN_BIT_0FFSET);
}

int charger_get_input_current(int *input_current)
{
	int rv, reg;

	/*
	 * IIN_DPM register reflects the actual input current limit programmed
	 * in the register, either from host or from ICO. After ICO, the
	 * current limit used by DPM regulation may differ from the IIN_HOST
	 * register settings.
	 */
	rv = raw_read16(BQ25710_REG_IIN_DPM, &reg);
	if (!rv)
		*input_current =
			REG_TO_INPUT_CURRENT((reg >>
					      BQ25710_CHARGE_IIN_BIT_0FFSET));

	return rv;
}

int charger_manufacturer_id(int *id)
{
	return raw_read16(BQ25710_REG_MANUFACTURER_ID, id);
}

int charger_device_id(int *id)
{
	return raw_read16(BQ25710_REG_DEVICE_ADDRESS, id);
}

#ifdef CONFIG_USB_PD_VBUS_MEASURE_CHARGER
int charger_get_vbus_voltage(int port)
{
	int reg, rv;

	rv = bq25710_adc_start(BQ25710_ADC_OPTION_EN_ADC_VBUS);
	if (rv)
		goto error;

	/* Read ADC value */
	rv = raw_read16(BQ25710_REG_ADC_VBUS_PSYS, &reg);
	if (rv)
		goto error;

	reg >>= BQ25710_ADC_VBUS_STEP_BIT_OFFSET;
	/*
	 * LSB => 64mV.
	 * Return 0 when VBUS <= 3.2V as ADC can't measure it.
	 */
	return reg ?
	       (reg * BQ25710_ADC_VBUS_STEP_MV + BQ25710_ADC_VBUS_BASE_MV) : 0;

error:
	CPRINTF("Could not read VBUS ADC! Error: %d\n", rv);
	return 0;
}
#endif

int charger_get_option(int *option)
{
	/* There are 4 option registers, but we only need the first for now. */
	return raw_read16(BQ25710_REG_CHARGE_OPTION_0, option);
}

int charger_set_option(int option)
{
	/* There are 4 option registers, but we only need the first for now. */
	return raw_write16(BQ25710_REG_CHARGE_OPTION_0, option);
}

#ifdef CONFIG_CHARGE_RAMP_HW

static void bq25710_chg_ramp_handle(void)
{
	int ramp_curr;

	/*
	 * Once the charge ramp is stable write back the stable ramp
	 * current to the host input current limit register
	 */
	ramp_curr = chg_ramp_get_current_limit();
	if (chg_ramp_is_stable()) {
		if (ramp_curr && !charger_set_input_current(ramp_curr))
			CPRINTF("bq25710: stable ramp current=%d\n", ramp_curr);
	} else {
		CPRINTF("bq25710: ICO stall, ramp current=%d\n", ramp_curr);
	}
	/*
	 * Disable ICO mode. When ICO mode is active the input current limit is
	 * given by the value in register IIN_DPM (0x22)
	 */
	charger_set_hw_ramp(0);
}
DECLARE_DEFERRED(bq25710_chg_ramp_handle);

int charger_set_hw_ramp(int enable)
{
	int option3_reg, option2_reg, rv;

	rv = raw_read16(BQ25710_REG_CHARGE_OPTION_3, &option3_reg);
	if (rv)
		return rv;
	rv = raw_read16(BQ25710_REG_CHARGE_OPTION_2, &option2_reg);
	if (rv)
		return rv;

	if (enable) {
		/*
		 * ICO mode can only be used when a battery is present. If there
		 * is no battery, then enabling ICO mode will lead to VSYS
		 * dropping out.
		 */
		if (!battery_is_present()) {
			CPRINTF("bq25710: no battery, skip ICO enable\n");
			return EC_ERROR_UNKNOWN;
		}

		/* Set InputVoltage register to BC1.2 minimum ramp voltage */
		rv = raw_write16(BQ25710_REG_INPUT_VOLTAGE,
			BQ25710_BC12_MIN_VOLTAGE_MV);
		if (rv)
			return rv;

		/*  Enable ICO algorithm */
		option3_reg |= BQ25710_CHARGE_OPTION_3_EN_ICO_MODE;

		/* 0b: Input current limit is set by BQ25710_REG_IIN_HOST */
		option2_reg &= ~BQ25710_CHARGE_OPTION_2_EN_EXTILIM;

		/* Charge ramp may take up to 2s to settle down */
		hook_call_deferred(&bq25710_chg_ramp_handle_data, (4 * SECOND));
	} else {
		/*  Disable ICO algorithm */
		option3_reg &= ~BQ25710_CHARGE_OPTION_3_EN_ICO_MODE;

		/*
		 * 1b: Input current limit is set by the lower value of
		 * ILIM_HIZ pin and BQ25710_REG_IIN_HOST
		 */
		option2_reg |= BQ25710_CHARGE_OPTION_2_EN_EXTILIM;
	}

	rv = raw_write16(BQ25710_REG_CHARGE_OPTION_2, option2_reg);
	if (rv)
		return rv;
	return raw_write16(BQ25710_REG_CHARGE_OPTION_3, option3_reg);
}

int chg_ramp_is_stable(void)
{
	int reg;

	if (raw_read16(BQ25710_REG_CHARGER_STATUS, &reg))
		return 0;

	return reg & BQ25710_CHARGE_STATUS_ICO_DONE;
}

int chg_ramp_get_current_limit(void)
{
	int reg, rv;

	rv = raw_read16(BQ25710_REG_IIN_DPM, &reg);
	if (rv) {
		CPRINTF("Could not read iin_dpm current limit! Error: %d\n",
			rv);
		return 0;
	}

	return ((reg >> BQ25710_IIN_DPM_BIT_SHIFT) * BQ25710_IIN_DPM_STEP_MA +
		BQ25710_IIN_DPM_STEP_MA);
}
#endif /* CONFIG_CHARGE_RAMP_HW */

#ifdef CONFIG_CHARGER_BQ25710_IDCHG_LIMIT_MA
/* Called on AP S5 -> S3  and S3/S0iX -> S0 transition */
static void bq25710_chipset_startup(void)
{
	bq25710_set_low_power_mode(0);
}
DECLARE_HOOK(HOOK_CHIPSET_STARTUP, bq25710_chipset_startup, HOOK_PRIO_DEFAULT);
DECLARE_HOOK(HOOK_CHIPSET_RESUME, bq25710_chipset_startup, HOOK_PRIO_DEFAULT);


/* Called on AP S0 -> S0iX/S3 or S3 -> S5 transition */
static void bq25710_chipset_suspend(void)
{
	bq25710_set_low_power_mode(1);
}
DECLARE_HOOK(HOOK_CHIPSET_SUSPEND, bq25710_chipset_suspend, HOOK_PRIO_DEFAULT);
DECLARE_HOOK(HOOK_CHIPSET_SHUTDOWN, bq25710_chipset_suspend, HOOK_PRIO_DEFAULT);
#endif

#ifdef CONFIG_CMD_CHARGER_DUMP
static int console_bq25710_dump_regs(int argc, char **argv)
{
	int i;
	int val;

	/* Dump all readable registers on bq25710. */
	static const uint8_t regs[] = {
		BQ25710_REG_CHARGE_OPTION_0,
		BQ25710_REG_CHARGE_CURRENT,
		BQ25710_REG_MAX_CHARGE_VOLTAGE,
		BQ25710_REG_CHARGE_OPTION_1,
		BQ25710_REG_CHARGE_OPTION_2,
		BQ25710_REG_CHARGE_OPTION_3,
		BQ25710_REG_PROCHOT_OPTION_0,
		BQ25710_REG_PROCHOT_OPTION_1,
		BQ25710_REG_ADC_OPTION,
		BQ25710_REG_CHARGER_STATUS,
		BQ25710_REG_PROCHOT_STATUS,
		BQ25710_REG_IIN_DPM,
		BQ25710_REG_ADC_VBUS_PSYS,
		BQ25710_REG_ADC_IBAT,
		BQ25710_REG_ADC_CMPIN_IIN,
		BQ25710_REG_ADC_VSYS_VBAT,
		BQ25710_REG_PROCHOT_OPTION_1,
		BQ25710_REG_OTG_VOLTAGE,
		BQ25710_REG_OTG_CURRENT,
		BQ25710_REG_INPUT_VOLTAGE,
		BQ25710_REG_MIN_SYSTEM_VOLTAGE,
		BQ25710_REG_IIN_HOST,
		BQ25710_REG_MANUFACTURER_ID,
		BQ25710_REG_DEVICE_ADDRESS,
	};

	for (i = 0; i < ARRAY_SIZE(regs); ++i) {
		if (raw_read16(regs[i], &val))
			continue;
		ccprintf("BQ25710 REG 0x%02x:  0x%04x\n", regs[i], val);
	}

	return 0;
}
DECLARE_CONSOLE_COMMAND(charger_dump, console_bq25710_dump_regs,
			"",
			"Dump all charger registers");

#endif /* CONFIG_CMD_CHARGER_DUMP */

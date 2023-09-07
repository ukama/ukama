/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2011 The ChromiumOS Authors.  All rights reserved.
 * Copyright (C) 2016 secunet Security Networks AG
 * Copyright (C) 2019 Protectli
 * Copyright (C) 2019 Libretrend LDA
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <console/console.h>
#include <delay.h>
#include <stddef.h>
#include <superio/hwm5_conf.h>

#include "env_ctrl.h"
#include "env_ctrl_chip.h"

static void extemp_force_idle_status(const u16 base)
{
	u8 reg;
	int retries = 10;

	/* Wait up to 10ms for non-busy state. */
	while (retries > 0) {
		reg = pnp_read_hwm5_index(base, ITE_EC_EXTEMP_STATUS);

		if ((reg & ITE_EC_EXTEMP_STATUS_HOST_BUSY) == 0x0)
			break;

		retries--;

		mdelay(1);
	}

	if (retries == 0 && (reg & ITE_EC_EXTEMP_STATUS_HOST_BUSY) == 0x1) {
		/*
		 * SIO is busy due to unfinished peci transaction.
		 * Re-configure Register 0x8E to terminate processes.
		 */
		pnp_write_hwm5_index(base, ITE_EC_EXTEMP_CONTROL,
			ITE_EC_EXTEMP_CTRL_AUTO_4HZ | ITE_EC_EXTEMP_CTRL_AUTO_START);
	}
}

/*
 * Setup PECI interface
 */
static void enable_peci(const u16 base)
{
	/* Enable PECI interface */
	pnp_write_hwm5_index(base, ITE_EC_INTERFACE_SELECT,
			   ITE_EC_INTERFACE_SEL_PECI | ITE_EC_INTERFACE_SPEED_TOLERANCE);

	/* Setup External Temperature using PECI GetTemp */
	pnp_write_hwm5_index(base, ITE_EC_EXTEMP_ADDRESS, PECI_CLIENT_ADDRESS);
	pnp_write_hwm5_index(base, ITE_EC_EXTEMP_COMMAND, PECI_GETTEMP_COMMAND);
	pnp_write_hwm5_index(base, ITE_EC_EXTEMP_WRITE_LENGTH, PECI_GETTEMP_WRITE_LENGTH);
	pnp_write_hwm5_index(base, ITE_EC_EXTEMP_READ_LENGTH, PECI_GETTEMP_READ_LENGTH);
	pnp_write_hwm5_index(base, ITE_EC_EXTEMP_CONTROL,
		ITE_EC_EXTEMP_CTRL_AUTO_4HZ | ITE_EC_EXTEMP_CTRL_AUTO_START);
}

/*
 * Set up External Temperature to read via PECI or thermal diode/resistor
 * into TMPINx register
 */
static void enable_tmpin(const u16 base, const u8 tmpin,
			 const struct ite_ec_thermal_config *const conf)
{
	u8 reg;

	reg = pnp_read_hwm5_index(base, ITE_EC_ADC_TEMP_CHANNEL_ENABLE);

	switch (conf->mode) {
	case THERMAL_PECI:
		if (reg & ITE_EC_ADC_TEMP_EXT_REPORTS_TO_MASK) {
			printk(BIOS_WARNING, "PECI specified for multiple TMPIN\n");
			return;
		}
		enable_peci(base);
		reg |= ITE_EC_ADC_TEMP_EXT_REPORTS_TO(tmpin);
		break;
	case THERMAL_DIODE:
		reg |= ITE_EC_ADC_TEMP_DIODE_MODE(tmpin);
		break;
	case THERMAL_RESISTOR:
		reg |= ITE_EC_ADC_TEMP_RESISTOR_MODE(tmpin);
		break;
	default:
		printk(BIOS_WARNING,
		       "Unsupported thermal mode 0x%x on TMPIN%d\n",
		       conf->mode, tmpin);
		return;
	}

	pnp_write_hwm5_index(base, ITE_EC_ADC_TEMP_CHANNEL_ENABLE, reg);

	/* Set temperature offsets */
	if (conf->mode != THERMAL_RESISTOR) {
		reg = pnp_read_hwm5_index(base, ITE_EC_BEEP_ENABLE);
		reg |= ITE_EC_TEMP_ADJUST_WRITE_ENABLE;
		pnp_write_hwm5_index(base, ITE_EC_BEEP_ENABLE, reg);
		pnp_write_hwm5_index(base, ITE_EC_TEMP_ADJUST[tmpin-1], conf->offset);
	}

	/* Set temperature limits */
	u8 max = conf->max;
	pnp_write_hwm5_index(base, ITE_EC_HIGH_TEMP_LIMIT(tmpin), max ? max : 127);
	pnp_write_hwm5_index(base, ITE_EC_LOW_TEMP_LIMIT(tmpin), conf->min);

	/* Enable the startup of monitoring operation */
	reg = pnp_read_hwm5_index(base, ITE_EC_CONFIGURATION);
	reg |= ITE_EC_CONFIGURATION_START;
	pnp_write_hwm5_index(base, ITE_EC_CONFIGURATION, reg);
}

static void fan_smartconfig(const u16 base, const u8 fan,
			    const enum ite_ec_fan_mode mode,
			    const struct ite_ec_fan_smartconfig *const conf)
{
	u8 pwm_ctrl;
	u8 pwm_start = 0;
	u8 pwm_auto = 0;

	if (mode == FAN_SMART_SOFTWARE) {
		pwm_ctrl = ITE_EC_FAN_CTL_PWM_MODE_SOFTWARE;

		/* 50% duty cycle by default */
		const u8 duty = conf->pwm_start ? conf->pwm_start : 50;
		if (CONFIG(SUPERIO_ITE_ENV_CTRL_8BIT_PWM))
			pwm_start = ITE_EC_FAN_CTL_PWM_DUTY(duty);
		else
			pwm_ctrl |= ITE_EC_FAN_CTL_PWM_DUTY(duty);
	} else {
		pwm_ctrl = ITE_EC_FAN_CTL_PWM_MODE_AUTOMATIC;
		pwm_ctrl |= ITE_EC_FAN_CTL_TEMPIN(conf->tmpin);

		pwm_start = ITE_EC_FAN_CTL_PWM_START_DUTY(conf->pwm_start);

		if (CONFIG(SUPERIO_ITE_ENV_CTRL_7BIT_SLOPE_REG)) {
			pwm_auto = conf->slope & 0x7f;
		} else {
			pwm_start |= ITE_EC_FAN_CTL_PWM_SLOPE_BIT6(conf->slope);
			pwm_auto = ITE_EC_FAN_CTL_PWM_SLOPE_LOWER(conf->slope);
		}

		if (conf->smoothing)
			pwm_auto |= ITE_EC_FAN_CTL_AUTO_SMOOTHING_EN;

		pnp_write_hwm5_index(base, ITE_EC_FAN_CTL_TEMP_LIMIT_OFF(fan), conf->tmp_off);
		pnp_write_hwm5_index(base, ITE_EC_FAN_CTL_TEMP_LIMIT_START(fan),
			conf->tmp_start);
		/* Full speed above 127°C by default */
		pnp_write_hwm5_index(base, ITE_EC_FAN_CTL_TEMP_LIMIT_FULL(fan),
			conf->tmp_full ? conf->tmp_full : 127);
		pnp_write_hwm5_index(base, ITE_EC_FAN_CTL_DELTA_TEMP(fan),
			ITE_EC_FAN_CTL_DELTA_TEMP_INTRVL(conf->tmp_delta));
	}

	pnp_write_hwm5_index(base, ITE_EC_FAN_CTL_PWM_CONTROL(fan), pwm_ctrl);
	pnp_write_hwm5_index(base, ITE_EC_FAN_CTL_PWM_START(fan), pwm_start);
	pnp_write_hwm5_index(base, ITE_EC_FAN_CTL_PWM_AUTO(fan), pwm_auto);
}

static void enable_fan(const u16 base, const u8 fan,
		       const struct ite_ec_fan_config *const conf)
{
	u8 reg;

	if (conf->mode == FAN_IGNORE ||
	    (CONFIG(SUPERIO_ITE_ENV_CTRL_NO_ONOFF) &&
	     conf->mode <= FAN_MODE_OFF))
		return;

	/* FAN_CTL2 might have its own frequency setting */
	if (CONFIG(SUPERIO_ITE_ENV_CTRL_PWM_FREQ2) && fan == 2) {
		reg = pnp_read_hwm5_index(base, ITE_EC_ADC_TEMP_EXTRA_CHANNEL_ENABLE);
		reg &= ~ITE_EC_FAN_PWM_CLOCK_MASK;
		reg |= ITE_EC_FAN_PWM_DEFAULT_CLOCK;
		pnp_write_hwm5_index(base, ITE_EC_ADC_TEMP_EXTRA_CHANNEL_ENABLE, reg);
	}

	if (conf->mode >= FAN_SMART_SOFTWARE) {
		fan_smartconfig(base, fan, conf->mode, &conf->smart);
	} else {
		reg = pnp_read_hwm5_index(base, ITE_EC_FAN_CTL_MODE);
		if (conf->mode == FAN_MODE_ON)
			reg |= ITE_EC_FAN_CTL_ON(fan);
		else
			reg &= ~ITE_EC_FAN_CTL_ON(fan);
		pnp_write_hwm5_index(base, ITE_EC_FAN_CTL_MODE, reg);
	}

	if (CONFIG(SUPERIO_ITE_ENV_CTRL_FAN16_CONFIG)
	    && conf->mode >= FAN_MODE_ON) {
		reg = pnp_read_hwm5_index(base, ITE_EC_FAN_TAC_COUNTER_ENABLE);
		reg |= ITE_EC_FAN_TAC_16BIT_ENABLE(fan);
		pnp_write_hwm5_index(base, ITE_EC_FAN_TAC_COUNTER_ENABLE, reg);
	}

	if (CONFIG(SUPERIO_ITE_ENV_CTRL_5FANS) && fan > 3) {
		reg = pnp_read_hwm5_index(base, ITE_EC_FAN_SEC_CTL);
		if (conf->mode >= FAN_MODE_ON)
			reg |= ITE_EC_FAN_SEC_CTL_TAC_EN(fan);
		else
			reg &= ~ITE_EC_FAN_SEC_CTL_TAC_EN(fan);
		pnp_write_hwm5_index(base, ITE_EC_FAN_SEC_CTL, reg);
	} else {
		reg = pnp_read_hwm5_index(base, ITE_EC_FAN_MAIN_CTL);
		if (conf->mode >= FAN_MODE_ON)
			reg |= ITE_EC_FAN_MAIN_CTL_TAC_EN(fan);
		else
			reg &= ~ITE_EC_FAN_MAIN_CTL_TAC_EN(fan);

		/* Some ITEs have SmartGuardian always enabled */
		if (!CONFIG(SUPERIO_ITE_ENV_CTRL_NO_ONOFF)) {
			if (conf->mode >= FAN_SMART_SOFTWARE)
				reg |= ITE_EC_FAN_MAIN_CTL_SMART(fan);
			else
				reg &= ~ITE_EC_FAN_MAIN_CTL_SMART(fan);
		}
		pnp_write_hwm5_index(base, ITE_EC_FAN_MAIN_CTL, reg);
	}
}

static void enable_beeps(const u16 base, const struct ite_ec_config *const conf)
{
	u8 reg = 0;
	u8 freq = ITE_EC_BEEP_TONE_DIVISOR(10) | ITE_EC_BEEP_FREQ_DIVISOR(10);

	if (conf->tmpin_beep) {
		reg |= ITE_EC_BEEP_ON_TMP_LIMIT;
		pnp_write_hwm5_index(base, ITE_EC_BEEP_FREQ_DIV_OF_TMPIN, freq);
	}
	if (conf->fan_beep) {
		reg |= ITE_EC_BEEP_ON_FAN_LIMIT;
		pnp_write_hwm5_index(base, ITE_EC_BEEP_FREQ_DIV_OF_FAN, freq);
	}
	if (conf->vin_beep) {
		reg |= ITE_EC_BEEP_ON_VIN_LIMIT;
		pnp_write_hwm5_index(base, ITE_EC_BEEP_FREQ_DIV_OF_VIN, freq);
	}

	if (reg) {
		reg |= pnp_read_hwm5_index(base, ITE_EC_BEEP_ENABLE);
		pnp_write_hwm5_index(base, ITE_EC_BEEP_ENABLE, reg);
	}
}

void ite_ec_init(const u16 base, const struct ite_ec_config *const conf)
{
	size_t i;

	/* Configure 23.43kHz PWM active high output */
	u8 fan_ctl = pnp_read_hwm5_index(base, ITE_EC_FAN_CTL_MODE);
	fan_ctl &= ~ITE_EC_FAN_PWM_CLOCK_MASK;
	fan_ctl |= ITE_EC_FAN_PWM_DEFAULT_CLOCK;
	fan_ctl |= ITE_EC_FAN_CTL_POLARITY_HIGH;
	pnp_write_hwm5_index(base, ITE_EC_FAN_CTL_MODE, fan_ctl);

	/* Enable HWM if configured */
	for (i = 0; i < ITE_EC_TMPIN_CNT; ++i)
		enable_tmpin(base, i + 1, &conf->tmpin[i]);

	/* Enable External Sensor SMBus Host if configured */
	if (conf->smbus_en) {
		pnp_write_hwm5_index(base, ITE_EC_INTERFACE_SELECT,
			pnp_read_hwm5_index(base, ITE_EC_INTERFACE_SELECT) |
				ITE_EC_INTERFACE_SMB_ENABLE);
	}

	/* Enable reading of voltage pins */
	pnp_write_hwm5_index(base, ITE_EC_ADC_VOLTAGE_CHANNEL_ENABLE, conf->vin_mask);

	/* Enable FANx if configured */
	for (i = 0; i < ITE_EC_FAN_CNT; ++i)
		enable_fan(base, i + 1, &conf->fan[i]);

	/* Enable beeps if configured */
	enable_beeps(base, conf);

	/*
	 * System may get wrong temperature data when SIO is in
	 * busy state. Therefore, check the status and terminate
	 * processes if needed.
	 */
	for (i = 0; i < ITE_EC_TMPIN_CNT; ++i)
		if (conf->tmpin[i].mode == THERMAL_PECI)
			extemp_force_idle_status(base);
}

/*
 * This file is part of the coreboot project.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#ifndef DRIVERS_I2C_NCT7802Y_CHIP_H
#define DRIVERS_I2C_NCT7802Y_CHIP_H

#include <stdint.h>

#define NCT7802Y_PECI_CNT	2
#define NCT7802Y_FAN_CNT	3

enum nct7802y_peci_mode {
	PECI_DISABLED = 0,
	PECI_DOMAIN_0,
	PECI_DOMAIN_1,
	PECI_HIGHEST,
};

struct nct7802y_peci_config {
	enum nct7802y_peci_mode mode;
	u8 base_temp;
};

enum nct7802y_fan_mode {
	FAN_IGNORE = 0,
	FAN_MANUAL,
	FAN_SMART,
};

enum nct7802y_fan_smartmode {
	SMART_FAN_DUTY = 0,	/* Target values given in duty cycle %. */
	SMART_FAN_RPM,		/* Target valuse given in RPM. */
};

enum nct7802y_fan_speed {
	FAN_SPEED_NORMAL = 0,	/* RPM values <= 12,750. */
	FAN_SPPED_HIGHSPEED,	/* RPM values <= 25,500. */
};

enum nct7802y_fan_pecierror {
	PECI_ERROR_KEEP = 0,	/* Keep current value. */
	PECI_ERROR_VALUE,	/* Use `pecierror_minduty`. */
	PECI_ERROR_FULLSPEED,	/* Run PWM at 100% duty cycle. */
};

enum nct7802y_temp_source {
	TEMP_SOURCE_REMOTE_1 = 0,
	TEMP_SOURCE_REMOTE_2,
	TEMP_SOURCE_REMOTE_3,
	TEMP_SOURCE_LOCAL,
	TEMP_SOURCE_PECI_0,
	TEMP_SOURCE_PECI_1,
	TEMP_SOURCE_PROGRAMMABLE_0,
	TEMP_SOURCE_PROGRAMMABLE_1,
};

struct nct7802y_fan_smartconfig {
	enum nct7802y_fan_smartmode mode;
	enum nct7802y_fan_speed speed;
	enum nct7802y_temp_source tempsrc;
	struct {
		u8 temp;
		u16 target;
	} table[4];
	u8 critical_temp;
};

struct nct7802y_fan_config {
	enum nct7802y_fan_mode mode;
	union {
		u8 duty_cycle;
		struct nct7802y_fan_smartconfig smart;
	};
};

/* Implements only those parts currently used by coreboot mainboards. */
struct drivers_i2c_nct7802y_config {
	struct nct7802y_peci_config peci[NCT7802Y_PECI_CNT];
	struct nct7802y_fan_config fan[NCT7802Y_FAN_CNT];
	enum nct7802y_fan_pecierror on_pecierror;
	u8 pecierror_minduty;
};

#endif /* DRIVERS_I2C_NCT7802Y_CHIP_H */

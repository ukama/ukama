/* Copyright 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/* Emulator board configuration */

#ifndef __CROS_EC_BOARD_H
#define __CROS_EC_BOARD_H

/* Optional features */
#define CONFIG_EXTPOWER_GPIO
#undef CONFIG_FMAP
#define CONFIG_POWER_BUTTON
#undef CONFIG_WATCHDOG
#define CONFIG_SWITCH
#define CONFIG_INDUCTIVE_CHARGING

#undef CONFIG_CONSOLE_HISTORY
#define CONFIG_CONSOLE_HISTORY 4

#define CONFIG_WP_ACTIVE_HIGH

#define CONFIG_LIBCRYPTOC

#include "gpio_signal.h"

enum temp_sensor_id {
	TEMP_SENSOR_CPU = 0,
	TEMP_SENSOR_BOARD,
	TEMP_SENSOR_CASE,
	TEMP_SENSOR_BATTERY,

	TEMP_SENSOR_COUNT
};

enum adc_channel {
	ADC_CH_CHARGER_CURRENT,
	ADC_AC_ADAPTER_ID_VOLTAGE,

	ADC_CH_COUNT
};

/* Fake test charge suppliers */
enum {
	CHARGE_SUPPLIER_TEST1,
	CHARGE_SUPPLIER_TEST2,
	CHARGE_SUPPLIER_TEST3,
	CHARGE_SUPPLIER_TEST4,
	CHARGE_SUPPLIER_TEST5,
	CHARGE_SUPPLIER_TEST6,
	CHARGE_SUPPLIER_TEST7,
	CHARGE_SUPPLIER_TEST8,
	CHARGE_SUPPLIER_TEST9,
	CHARGE_SUPPLIER_TEST10,
	CHARGE_SUPPLIER_TEST_COUNT
};

/* Standard-current Rp */
#define PD_SRC_VNC           PD_SRC_DEF_VNC_MV
#define PD_SRC_RD_THRESHOLD  PD_SRC_DEF_RD_THRESH_MV

/* delay necessary for the voltage transition on the power supply */
#define PD_POWER_SUPPLY_TURN_ON_DELAY  20000 /* us */
#define PD_POWER_SUPPLY_TURN_OFF_DELAY 20000 /* us */

/* Define typical operating power and max power */
#define PD_OPERATING_POWER_MW 15000
#define PD_MAX_POWER_MW       60000
#define PD_MAX_CURRENT_MA     3000
#define PD_MAX_VOLTAGE_MV     20000

#define PD_MIN_CURRENT_MA     500
#define PD_MIN_POWER_MW       7500

#endif /* __CROS_EC_BOARD_H */

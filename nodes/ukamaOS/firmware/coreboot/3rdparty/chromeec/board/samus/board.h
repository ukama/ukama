/* Copyright 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/* Configuration for Samus mainboard */

#ifndef __CROS_EC_BOARD_H
#define __CROS_EC_BOARD_H

/* Debug features */
#define CONFIG_CONSOLE_CMDHELP
#define CONFIG_TASK_PROFILING

#undef HEY_USE_BUILTIN_CLKRUN

/* Optional features */
#define CONFIG_ACCELGYRO_LSM6DS0
#define CONFIG_ACCEL_KXCJ9
#define CONFIG_ACCEL_STD_REF_FRAME_OLD
#define CONFIG_ALS_ISL29035
#define CONFIG_BOARD_VERSION_GPIO
#define CONFIG_CMD_ACCELS
#define CONFIG_CMD_ACCEL_INFO
#undef  CONFIG_BATTERY_CRITICAL_SHUTDOWN_TIMEOUT
#define CONFIG_BATTERY_CRITICAL_SHUTDOWN_TIMEOUT 60
#define CONFIG_BATTERY_CUT_OFF
#define CONFIG_POWER_COMMON
#define CONFIG_POWER_SHUTDOWN_PAUSE_IN_S5
#define CONFIG_CHIPSET_CAN_THROTTLE
#define CONFIG_I2C
#define CONFIG_I2C_MASTER
#define CONFIG_KEYBOARD_BOARD_CONFIG
#define CONFIG_KEYBOARD_PROTOCOL_8042
#define CONFIG_KEYBOARD_COL2_INVERTED
#define CONFIG_KEYBOARD_SCANCODE_CALLBACK
#define CONFIG_LID_ANGLE
#define CONFIG_LIGHTBAR_POWER_RAILS
#define CONFIG_LOW_POWER_IDLE
#define CONFIG_POWER_BUTTON
#define CONFIG_POWER_BUTTON_X86
/* Note: not CONFIG_BACKLIGHT_LID. It's handled specially for Samus. */
#define CONFIG_BACKLIGHT_REQ_GPIO GPIO_PCH_BL_EN
/* TODO(crosbug.com/p/29467): remove this workaround when possible. */
#define CONFIG_BATTERY_REQUESTS_NIL_WHEN_DEAD
#define CONFIG_CHARGER_PROFILE_OVERRIDE
#define CONFIG_BATTERY_SMART
#define CONFIG_BATTERY_REVIVE_DISCONNECT
#define CONFIG_CHARGER
#define CONFIG_CHARGER_BQ24773
#define CONFIG_CHARGER_ILIM_PIN_DISABLED
#define CONFIG_CHARGER_SENSE_RESISTOR 5
#define CONFIG_CHARGER_SENSE_RESISTOR_AC 10
#define CONFIG_CHARGER_INPUT_CURRENT 320
#define CONFIG_CHARGER_DISCHARGE_ON_AC
#define CONFIG_CHARGER_DISCHARGE_ON_AC_CUSTOM
#define CONFIG_CHARGER_MIN_BAT_PCT_FOR_POWER_ON 1
#define CONFIG_FANS 2
#define CONFIG_FAN_UPDATE_PERIOD 10
#define CONFIG_FPU
#define CONFIG_GESTURE_DETECTION
#define CONFIG_GESTURE_SW_DETECTION
#define CONFIG_GESTURE_SAMPLING_INTERVAL_MS 5
#undef  CONFIG_HIBERNATE_DELAY_SEC
#define CONFIG_HIBERNATE_DELAY_SEC (3600 * 24 * 7)
#define CONFIG_HOSTCMD_PD
#define CONFIG_HOSTCMD_PD_CHG_CTRL
#define CONFIG_HOSTCMD_PD_PANIC
#define CONFIG_PECI_TJMAX 105
#define CONFIG_PWM
#define CONFIG_PWM_KBLIGHT
#define CONFIG_TEMP_SENSOR
#define CONFIG_TEMP_SENSOR_TMP006
#define CONFIG_TEMP_SENSOR_POWER_GPIO GPIO_PP3300_DSW_GATED_EN
#define CONFIG_THROTTLE_AP
#define CONFIG_UART_HOST 2
#define CONFIG_USB_PORT_POWER_SMART
#define CONFIG_USB_PORT_POWER_SMART_DEFAULT_MODE USB_CHARGE_MODE_CDP
#define CONFIG_USB_PORT_POWER_SMART_INVERTED
#define GPIO_USB1_ILIM_SEL GPIO_USB1_ILIM_SEL_L
#define GPIO_USB2_ILIM_SEL GPIO_USB2_ILIM_SEL_L
#define CONFIG_VBOOT_HASH
#define CONFIG_WIRELESS
#define CONFIG_WIRELESS_SUSPEND \
	(EC_WIRELESS_SWITCH_WLAN | EC_WIRELESS_SWITCH_WLAN_POWER)
/* Do we want EC_WIRELESS_SWITCH_WWAN as well? */

#ifndef __ASSEMBLER__

/* I2C ports */
#define I2C_PORT_BACKLIGHT 0
#define I2C_PORT_BATTERY 0
#define I2C_PORT_CHARGER 0
#define I2C_PORT_PD_MCU 0
#define I2C_PORT_ALS 1
#define I2C_PORT_ACCEL 1
#define I2C_PORT_LIGHTBAR 1
#define I2C_PORT_THERMAL 5

/* 13x8 keyboard scanner uses an entire GPIO bank for row inputs */
#define KB_SCAN_ROW_IRQ  LM4_IRQ_GPIOK
#define KB_SCAN_ROW_GPIO LM4_GPIO_K

/* Host connects to keyboard controller module via LPC */
#define HOST_KB_BUS_LPC

/* USB ports managed by the EC */
#define USB_PORT_COUNT 2

#include "gpio_signal.h"

/* x86 signal definitions */
enum x86_signal {
	X86_PGOOD_PP1050 = 0,
	X86_PGOOD_PP1200,
	X86_PGOOD_PP1800,
	X86_PGOOD_VCORE,

	X86_SLP_S0_DEASSERTED,
	X86_SLP_S3_DEASSERTED,
	X86_SLP_S5_DEASSERTED,
	X86_SLP_SUS_DEASSERTED,
	X86_SUSWARN_DEASSERTED,

	/* Number of X86 signals */
	POWER_SIGNAL_COUNT
};

enum adc_channel {
	/* EC internal die temperature in degrees K. */
	ADC_CH_EC_TEMP = 0,
	/* BAT_TEMP */
	ADC_CH_BAT_TEMP,

	ADC_CH_COUNT
};

enum pwm_channel {
	PWM_CH_KBLIGHT,

	/* Number of PWM channels */
	PWM_CH_COUNT
};

enum temp_sensor_id {
	/* CPU die temperature via PECI */
	TEMP_SENSOR_CPU_PECI,
	/* EC internal temperature sensor */
	TEMP_SENSOR_EC_INTERNAL,
	/* TMP006 U40, die/object temperature near battery charger */
	TEMP_SENSOR_I2C_U40_DIE,
	TEMP_SENSOR_I2C_U40_OBJECT,
	/* TMP006 U41, die/object temperature near CPU */
	TEMP_SENSOR_I2C_U41_DIE,
	TEMP_SENSOR_I2C_U41_OBJECT,
	/* TMP006 U42, die/object temperature left side of C-case */
	TEMP_SENSOR_I2C_U42_DIE,
	TEMP_SENSOR_I2C_U42_OBJECT,
	/* TMP006 U43, die/object temperature right side of C-case */
	TEMP_SENSOR_I2C_U43_DIE,
	TEMP_SENSOR_I2C_U43_OBJECT,
	/* TMP006 U115, die/object temperature right side of D-case */
	TEMP_SENSOR_I2C_U115_DIE,
	TEMP_SENSOR_I2C_U115_OBJECT,
	/* TMP006 U116, die/object temperature left side of D-case */
	TEMP_SENSOR_I2C_U116_DIE,
	TEMP_SENSOR_I2C_U116_OBJECT,

	/* Battery temperature sensor */
	TEMP_SENSOR_BATTERY,

	TEMP_SENSOR_COUNT
};

enum sensor_id {
	BASE_ACCEL,
	LID_ACCEL,
	BASE_GYRO,
	SENSOR_COUNT,
};

/* The number of TMP006 sensor chips on the board. */
#define TMP006_COUNT 6

/* Light sensors attached to the EC. */
enum als_id {
	ALS_ISL29035 = 0,

	ALS_COUNT,
};

/* Wireless signals */
#define WIRELESS_GPIO_WLAN GPIO_WLAN_OFF_L
#define WIRELESS_GPIO_WLAN_POWER GPIO_PP3300_WLAN_EN

/* Discharge battery when on AC power for factory test. */
int board_is_discharging_on_ac(void);

/* Reset PD MCU */
void board_reset_pd_mcu(void);

/* Backboost detected interrupt */
void bkboost_det_interrupt(enum gpio_signal signal);

/* Interrupt handler for JTAG clock */
void jtag_interrupt(enum gpio_signal signal);

/* Bit masks for turning on PP5000 rail in G3 */
#define PP5000_IN_G3_AC       BIT(0)
#define PP5000_IN_G3_LIGHTBAR BIT(1)

/* Enable/disable PP5000 rail mask in G3 */
void set_pp5000_in_g3(int mask, int enable);

/* Define for sensor tasks */
#define CONFIG_GESTURE_SENSOR_BATTERY_TAP 0
#define CONFIG_GESTURE_TAP_OUTER_WINDOW_T 200
#define CONFIG_GESTURE_TAP_INNER_WINDOW_T 30
#define CONFIG_GESTURE_TAP_MIN_INTERSTICE_T 120
#define CONFIG_GESTURE_TAP_MAX_INTERSTICE_T 500

#define CONFIG_LID_ANGLE_SENSOR_BASE BASE_ACCEL
#define CONFIG_LID_ANGLE_SENSOR_LID LID_ACCEL

#endif /* !__ASSEMBLER__ */

#endif /* __CROS_EC_BOARD_H */

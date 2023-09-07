/* Copyright 2014 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/* Veyron board configuration */

#ifndef __CROS_EC_BOARD_H
#define __CROS_EC_BOARD_H

/* Increase size of UART TX buffer. */
#undef CONFIG_UART_TX_BUF_SIZE
#define CONFIG_UART_TX_BUF_SIZE 1024

/* Optional features */
#define CONFIG_AP_HANG_DETECT
#define CONFIG_BATTERY_CUT_OFF
#define CONFIG_BATTERY_SMART
#define CONFIG_BOARD_PRE_INIT
#define CONFIG_CHARGER
#define CONFIG_CHARGER_BQ24715
#define CONFIG_CHARGER_DISCHARGE_ON_AC
#define CONFIG_CHIPSET_RK3288
#define CONFIG_EXTPOWER_GPIO
#define CONFIG_FORCE_CONSOLE_RESUME
#define CONFIG_HOST_COMMAND_STATUS
#define CONFIG_HOSTCMD_VBNV_CONTEXT
#define CONFIG_I2C
#define CONFIG_I2C_MASTER
#define CONFIG_I2C_VIRTUAL_BATTERY
#define CONFIG_KEYBOARD_COL2_INVERTED
#define CONFIG_KEYBOARD_PROTOCOL_MKBP
#define CONFIG_LED_COMMON
#define CONFIG_LED_POLICY_STD
#define CONFIG_LED_BAT_ACTIVE_LOW
#define CONFIG_LED_POWER_ACTIVE_LOW
#define CONFIG_LOW_POWER_IDLE
#define CONFIG_LOW_POWER_S0
#define CONFIG_MKBP_USE_GPIO
#define CONFIG_POWER_BUTTON
#define CONFIG_POWER_BUTTON_ACTIVE_STATE 1
#define CONFIG_POWER_COMMON
#define CONFIG_PWM
#define CONFIG_SPI
#define CONFIG_STM_HWTIMER32
#define CONFIG_UART_RX_DMA
#define CONFIG_VBOOT_HASH

#define CONFIG_HOSTCMD_ALIGNED

#undef CONFIG_CONSOLE_CMDHELP
#undef CONFIG_CONSOLE_HISTORY
#undef CONFIG_SOFTWARE_PANIC
#undef CONFIG_WATCHDOG_HELP
#undef CONFIG_CMD_KEYBOARD

/* Not enough RO flash space */
#undef CONFIG_HOSTCMD_GET_UPTIME_INFO

#define CONFIG_HIBERNATE_WAKEUP_PINS (STM32_PWR_CSR_EWUP1 | STM32_PWR_CSR_EWUP6)

#ifndef __ASSEMBLER__

/* 48 MHz SYSCLK clock frequency */
#define CPU_CLOCK 48000000

/* Keyboard output port list */
#define KB_OUT_PORT_LIST GPIO_A, GPIO_B, GPIO_C

/* Single I2C port, where the EC is the master. */
#define I2C_PORT_MASTER 0
#define I2C_PORT_BATTERY I2C_PORT_MASTER
#define I2C_PORT_CHARGER I2C_PORT_MASTER
#define I2C_PORT_VIRTUAL_BATTERY I2C_PORT_MASTER

/* Timer selection */
#define TIM_CLOCK32  2
#define TIM_WATCHDOG 4

#include "gpio_signal.h"

enum power_signal {
	RK_POWER_GOOD = 0,
	RK_SUSPEND_ASSERTED,
	/* Number of power signals */
	POWER_SIGNAL_COUNT
};

enum pwm_channel {
	PWM_CH_POWER_LED = 0,
	/* Number of PWM channels */
	PWM_CH_COUNT
};

/* Charger module */
#define CONFIG_CHARGER_SENSE_RESISTOR 10 /* Charge sense resistor, mOhm */
#define CONFIG_CHARGER_SENSE_RESISTOR_AC 20 /* Input sensor resistor, mOhm */
/* Input current limit for 45W AC adapter:
 * 45W/19V*85%=2013mA, choose the closest charger setting = 2048mA
 */
#define CONFIG_CHARGER_INPUT_CURRENT 2048 /* mA, based on Link HW design */
#define CONFIG_CHARGER_CURRENT_LIMIT 3000 /* PL102 inductor 3.0A(3.8A) */

#endif /* !__ASSEMBLER__ */

#endif /* __CROS_EC_BOARD_H */

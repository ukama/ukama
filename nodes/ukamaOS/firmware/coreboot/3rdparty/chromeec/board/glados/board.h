/* Copyright 2015 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/* Glados board configuration */

#ifndef __CROS_EC_BOARD_H
#define __CROS_EC_BOARD_H

/*
 * Allow dangerous commands.
 * TODO(shawnn): Remove this config before production.
 */
#define CONFIG_SYSTEM_UNLOCKED

/* Optional features */
#define CONFIG_ACCELGYRO_BMI160
#define CONFIG_ACCEL_KX022
#define CONFIG_ADC
#define CONFIG_ALS_OPT3001
#define CONFIG_BATTERY_CUT_OFF
#define CONFIG_BATTERY_PRESENT_GPIO GPIO_BAT_PRESENT_L
#define CONFIG_BATTERY_SMART
#define CONFIG_BOARD_VERSION_GPIO
#define CONFIG_CHARGE_MANAGER
#define CONFIG_CHARGE_RAMP_HW

#define CONFIG_CHARGER

#define CONFIG_CHARGER_DISCHARGE_ON_AC
#define CONFIG_CHARGER_ISL9237
#define CONFIG_CHARGER_ILIM_PIN_DISABLED
#define CONFIG_CHARGER_INPUT_CURRENT 512
#define CONFIG_CHARGER_MIN_BAT_PCT_FOR_POWER_ON 1
#define CONFIG_CHARGER_PROFILE_OVERRIDE
#define CONFIG_CHARGER_SENSE_RESISTOR 10
#define CONFIG_CHARGER_SENSE_RESISTOR_AC 20
#define CONFIG_CMD_CHARGER_ADC_AMON_BMON

#define CONFIG_CHIPSET_SKYLAKE
#define CONFIG_CHIPSET_RESET_HOOK
#define CONFIG_CLOCK_CRYSTAL
#define CONFIG_EXTPOWER_GPIO
#define CONFIG_HOSTCMD_PD
#define CONFIG_HOSTCMD_PD_PANIC
#define CONFIG_I2C
#define CONFIG_I2C_MASTER
#define CONFIG_KEYBOARD_PROTOCOL_8042
#define CONFIG_LED_COMMON
#define CONFIG_LID_ANGLE
#define CONFIG_LID_ANGLE_SENSOR_BASE BASE_ACCEL
#define CONFIG_LID_ANGLE_SENSOR_LID LID_ACCEL
#define CONFIG_LID_SWITCH
#define CONFIG_LOW_POWER_IDLE
#define CONFIG_LTO
#define CONFIG_POWER_BUTTON
#define CONFIG_POWER_BUTTON_X86
#define CONFIG_POWER_COMMON
#define CONFIG_POWER_SIGNAL_INTERRUPT_STORM_DETECT_THRESHOLD 30
/* All data won't fit in data RAM.  So, moving boundary slightly. */
#undef CONFIG_RO_SIZE
#define CONFIG_RO_SIZE (108 * 1024)
#define CONFIG_SCI_GPIO GPIO_PCH_SCI_L
/* We're space constrained on GLaDOS, so reduce the UART TX buffer size. */
#undef CONFIG_UART_TX_BUF_SIZE
#define CONFIG_UART_TX_BUF_SIZE 512
#define CONFIG_USB_CHARGER
#define CONFIG_USB_MUX_PI3USB30532
#define CONFIG_USB_MUX_PS8740
#define CONFIG_USB_POWER_DELIVERY
#define CONFIG_USB_PD_ALT_MODE
#define CONFIG_USB_PD_ALT_MODE_DFP
#define CONFIG_USB_PD_DUAL_ROLE
#define CONFIG_USB_PD_LOGGING
#define CONFIG_USB_PD_PORT_COUNT 2
#define CONFIG_USB_PD_TCPM_TCPCI
#define CONFIG_USB_PD_TRY_SRC
#define CONFIG_USB_PD_VBUS_DETECT_GPIO
#define CONFIG_BC12_DETECT_PI3USB9281
#define CONFIG_BC12_DETECT_PI3USB9281_CHIP_COUNT 2
#define CONFIG_USBC_SS_MUX
#define CONFIG_USBC_SS_MUX_DFP_ONLY
#define CONFIG_USBC_VCONN
#define CONFIG_USBC_VCONN_SWAP
#define CONFIG_VBOOT_HASH
#define CONFIG_VOLUME_BUTTONS

#define CONFIG_SPI_FLASH_PORT 1
#define CONFIG_SPI_FLASH
#define CONFIG_FLASH_SIZE 524288
#define CONFIG_SPI_FLASH_W25X40

#define CONFIG_TEMP_SENSOR
#define CONFIG_TEMP_SENSOR_BD99992GW
#define CONFIG_THERMISTOR_NCP15WB
#define CONFIG_DPTF

/*
 * Enable 1 slot of secure temporary storage to support
 * suspend/resume with read/write memory training.
 */
#define CONFIG_VSTORE
#define CONFIG_VSTORE_SLOT_COUNT 1

#define CONFIG_WATCHDOG_HELP

#define CONFIG_WIRELESS
#define CONFIG_WIRELESS_SUSPEND \
	(EC_WIRELESS_SWITCH_WLAN | EC_WIRELESS_SWITCH_WLAN_POWER)

/* Wireless signals */
#define WIRELESS_GPIO_WLAN GPIO_WLAN_OFF_L
#define WIRELESS_GPIO_WLAN_POWER GPIO_PP3300_WLAN_EN

/* LED signals */
#define GPIO_BAT_LED_RED GPIO_CHARGE_LED_1
#define GPIO_BAT_LED_GREEN GPIO_CHARGE_LED_2

/* I2C ports */
#define I2C_PORT_PMIC MEC1322_I2C0_0
#define I2C_PORT_USB_CHARGER_1 MEC1322_I2C0_1
#define I2C_PORT_USB_MUX MEC1322_I2C0_1
#define I2C_PORT_USB_CHARGER_2 MEC1322_I2C0_0
#define I2C_PORT_PD_MCU MEC1322_I2C1
#define I2C_PORT_TCPC MEC1322_I2C1
#define I2C_PORT_ALS MEC1322_I2C2
#define I2C_PORT_ACCEL MEC1322_I2C2
#define I2C_PORT_BATTERY MEC1322_I2C3
#define I2C_PORT_CHARGER MEC1322_I2C3

/* Thermal sensors read through PMIC ADC interface */
#define I2C_PORT_THERMAL I2C_PORT_PMIC

/* Ambient Light Sensor address */
#define OPT3001_I2C_ADDR_FLAGS OPT3001_I2C_ADDR1_FLAGS

/* Modules we want to exclude */
#undef CONFIG_CMD_HASH
#undef CONFIG_CMD_TEMP_SENSOR
#undef CONFIG_CMD_TIMERINFO
#undef CONFIG_CONSOLE_CMDHELP
#undef CONFIG_CMD_ADC
#undef CONFIG_CMD_ACCELSPOOF
#undef CONFIG_CMD_FASTCHARGE
#undef CONFIG_CMD_GETTIME
#undef CONFIG_CMD_MEM
#ifdef SECTION_IS_RO
#undef CONFIG_CMD_BATTFAKE
#endif

#ifndef __ASSEMBLER__

#include "gpio_signal.h"
#include "registers.h"

/* ADC signal */
enum adc_channel {
	ADC_VBUS,
	ADC_AMON_BMON,
	ADC_PSYS,
	/* Number of ADC channels */
	ADC_CH_COUNT
};

enum temp_sensor_id {
	TEMP_SENSOR_BATTERY,

	/* These temp sensors are only readable in S0 */
	TEMP_SENSOR_AMBIENT,
	TEMP_SENSOR_CHARGER,
	TEMP_SENSOR_DRAM,
	TEMP_SENSOR_WIFI,

	TEMP_SENSOR_COUNT
};

enum sensor_id {
	BASE_ACCEL,
	BASE_GYRO,
	LID_ACCEL,
	SENSOR_COUNT,
};

/* Light sensors */
enum als_id {
	ALS_OPT3001 = 0,

	ALS_COUNT
};

/* TODO: determine the following board specific type-C power constants */
/*
 * delay to turn on the power supply max is ~16ms.
 * delay to turn off the power supply max is about ~180ms.
 */
#define PD_POWER_SUPPLY_TURN_ON_DELAY  30000  /* us */
#define PD_POWER_SUPPLY_TURN_OFF_DELAY 250000 /* us */

/* delay to turn on/off vconn */
#define PD_VCONN_SWAP_DELAY 5000 /* us */

/* Define typical operating power and max power */
#define PD_OPERATING_POWER_MW 15000
#define PD_MAX_POWER_MW       45000
#define PD_MAX_CURRENT_MA     3000

/* Try to negotiate to 20V since i2c noise problems should be fixed. */
#define PD_MAX_VOLTAGE_MV     20000

/* Reset PD MCU */
void board_reset_pd_mcu(void);

#endif /* !__ASSEMBLER__ */

#endif /* __CROS_EC_BOARD_H */

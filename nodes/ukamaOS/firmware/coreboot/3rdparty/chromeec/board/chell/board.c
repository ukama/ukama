/* Copyright 2015 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/* Chell board-specific configuration */

#include "adc_chip.h"
#include "bd99992gw.h"
#include "charge_manager.h"
#include "charge_state.h"
#include "charger.h"
#include "chipset.h"
#include "console.h"
#include "extpower.h"
#include "gpio.h"
#include "hooks.h"
#include "host_command.h"
#include "i2c.h"
#include "keyboard_scan.h"
#include "lid_switch.h"
#include "pi3usb9281.h"
#include "power.h"
#include "power_button.h"
#include "pwm.h"
#include "pwm_chip.h"
#include "spi.h"
#include "switch.h"
#include "system.h"
#include "task.h"
#include "tcpci.h"
#include "temp_sensor.h"
#include "timer.h"
#include "uart.h"
#include "usb_charge.h"
#include "usb_mux.h"
#include "usb_mux/ps874x.h"
#include "usb_pd.h"
#include "usb_pd_tcpm.h"
#include "util.h"

#define CPRINTS(format, args...) cprints(CC_USBCHARGE, format, ## args)
#define CPRINTF(format, args...) cprintf(CC_USBCHARGE, format, ## args)

#define I2C_ADDR_BD99992_FLAGS 0x30

/* Exchange status with PD MCU. */
static void pd_mcu_interrupt(enum gpio_signal signal)
{
#ifdef HAS_TASK_PDCMD
	/* Exchange status with PD MCU to determine interrupt cause */
	host_command_pd_send_status(0);
#endif
}

void vbus0_evt(enum gpio_signal signal)
{
	/* VBUS present GPIO is inverted */
	usb_charger_vbus_change(0, !gpio_get_level(signal));
	task_wake(TASK_ID_PD_C0);
}

void vbus1_evt(enum gpio_signal signal)
{
	/* VBUS present GPIO is inverted */
	usb_charger_vbus_change(1, !gpio_get_level(signal));
	task_wake(TASK_ID_PD_C1);
}

void usb0_evt(enum gpio_signal signal)
{
	task_set_event(TASK_ID_USB_CHG_P0, USB_CHG_EVENT_BC12, 0);
}

void usb1_evt(enum gpio_signal signal)
{
	task_set_event(TASK_ID_USB_CHG_P1, USB_CHG_EVENT_BC12, 0);
}

#include "gpio_list.h"

/* ADC channels */
const struct adc_t adc_channels[] = {
	/* Vbus sensing. Converted to mV, full ADC is equivalent to 30V. */
	[ADC_VBUS] = {"VBUS", 30000, 1024, 0, 1},
	/* Adapter current output or battery discharging current */
	[ADC_AMON_BMON] = {"AMON_BMON", 25000, 3072, 0, 3},
	/* System current consumption */
	[ADC_PSYS] = {"PSYS", 1, 1, 0, 4},

};
BUILD_ASSERT(ARRAY_SIZE(adc_channels) == ADC_CH_COUNT);

/* PWM channels. Must be in the exactly same order as in enum pwm_channel. */
const struct pwm_t pwm_channels[] = {
	/* Use alternate 100kHz clock source, keep active in low-power idle */
	{2, PWM_CONFIG_ALT_CLOCK | PWM_CONFIG_DSLEEP},
};
BUILD_ASSERT(ARRAY_SIZE(pwm_channels) == PWM_CH_COUNT);

const struct i2c_port_t i2c_ports[]  = {
	{"pmic",     MEC1322_I2C0_0, 400,  GPIO_I2C0_0_SCL, GPIO_I2C0_0_SDA},
	{"muxes",    MEC1322_I2C0_1, 400,  GPIO_I2C0_1_SCL, GPIO_I2C0_1_SDA},
	{"pd_mcu",   MEC1322_I2C1,   500,  GPIO_I2C1_SCL,   GPIO_I2C1_SDA},
	{"batt",     MEC1322_I2C3,   100,  GPIO_I2C3_SCL,   GPIO_I2C3_SDA},
};
const unsigned int i2c_ports_used = ARRAY_SIZE(i2c_ports);

const struct tcpc_config_t tcpc_config[CONFIG_USB_PD_PORT_COUNT] = {
	{
		.bus_type = EC_BUS_TYPE_I2C,
		.i2c_info = {
			.port = I2C_PORT_TCPC,
			.addr_flags = CONFIG_TCPC_I2C_BASE_ADDR_FLAGS,
		},
		.drv = &tcpci_tcpm_drv,
	},
	{
		.bus_type = EC_BUS_TYPE_I2C,
		.i2c_info = {
			.port = I2C_PORT_TCPC,
			.addr_flags = CONFIG_TCPC_I2C_BASE_ADDR_FLAGS + 1,
		},
		.drv = &tcpci_tcpm_drv,

	},
};

/* SPI devices */
const struct spi_device_t spi_devices[] = {
	{ CONFIG_SPI_FLASH_PORT, 0, GPIO_PVT_CS0},
};
const unsigned int spi_devices_used = ARRAY_SIZE(spi_devices);

const enum gpio_signal hibernate_wake_pins[] = {
	GPIO_AC_PRESENT,
	GPIO_LID_OPEN,
	GPIO_POWER_BUTTON_L,
};
const int hibernate_wake_pins_used = ARRAY_SIZE(hibernate_wake_pins);

#ifdef CONFIG_KEYBOARD_FACTORY_TEST
/*
 * We have total 28 pins for keyboard connecter, {-1, -1} mean
 * the N/A pin that don't consider it and reserve index 0 area
 * that we don't have pin 0.
 */
const int keyboard_factory_scan_pins[][2] = {
		{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1},
		{12, 6}, {4, 3}, {4, 2}, {0, 2}, {14, 2},
		{4, 0}, {0, 0}, {-1, -1}, {3, 2}, {10, 3},
		{10, 0}, {12, 5}, {-1, -1}, {10, 2}, {-1, -1},
		{0, 1}, {10, 4}, {-1, -1}, {-1, -1}, {0, 4},
		{10, 7}, {10, 6}, {0, 3}, {0, 5},
};

const int keyboard_factory_scan_pins_used =
			ARRAY_SIZE(keyboard_factory_scan_pins);
#endif

struct pi3usb9281_config pi3usb9281_chips[] = {
	{
		.i2c_port = I2C_PORT_USB_CHARGER_1,
		.mux_lock = NULL,
	},
	{
		.i2c_port = I2C_PORT_USB_CHARGER_2,
		.mux_lock = NULL,
	},
};
BUILD_ASSERT(ARRAY_SIZE(pi3usb9281_chips) ==
	     CONFIG_BC12_DETECT_PI3USB9281_CHIP_COUNT);

static int ps874x_tune_mux(int port)
{
	/* Apply same USB EQ settings to both Type-C mux */
	ps874x_tune_usb_eq(port,
			   PS874X_USB_EQ_TX_6_5_DB,
			   PS874X_USB_EQ_RX_14_3_DB);

	return EC_SUCCESS;
}

struct usb_mux usb_muxes[CONFIG_USB_PD_PORT_COUNT] = {
	{
		.port_addr = 0x1A,
		.driver = &ps874x_usb_mux_driver,
		.board_init = &ps874x_tune_mux,
	},
	{
		.port_addr = 0x10,
		.driver = &ps874x_usb_mux_driver,
		.board_init = &ps874x_tune_mux,
	}
};

/**
 * Reset PD MCU
 */
void board_reset_pd_mcu(void)
{
	gpio_set_level(GPIO_PD_RST_L, 0);
	usleep(100);
	gpio_set_level(GPIO_PD_RST_L, 1);
}

const struct temp_sensor_t temp_sensors[] = {
	{"Battery", TEMP_SENSOR_TYPE_BATTERY, charge_get_battery_temp, 0, 4},

	/* These BD99992GW temp sensors are only readable in S0 */
	{"Ambient", TEMP_SENSOR_TYPE_BOARD, bd99992gw_get_val,
		BD99992GW_ADC_CHANNEL_SYSTHERM0, 4},
	{"Charger", TEMP_SENSOR_TYPE_BOARD, bd99992gw_get_val,
		BD99992GW_ADC_CHANNEL_SYSTHERM1, 4},
	{"DRAM", TEMP_SENSOR_TYPE_BOARD, bd99992gw_get_val,
		BD99992GW_ADC_CHANNEL_SYSTHERM2, 4},
	{"Wifi", TEMP_SENSOR_TYPE_BOARD, bd99992gw_get_val,
		BD99992GW_ADC_CHANNEL_SYSTHERM3, 4},
};
BUILD_ASSERT(ARRAY_SIZE(temp_sensors) == TEMP_SENSOR_COUNT);

static void board_pmic_init(void)
{
	/* DISCHGCNT3 - enable 100 ohm discharge on V1.00A */
	i2c_write8(I2C_PORT_PMIC, I2C_ADDR_BD99992_FLAGS, 0x3e, 0x04);

	/*
	 * No need to re-init below settings since they are present on all MP
	 * ROs and PMIC settings are sticky across sysjump
	 */
	if (system_jumped_to_this_image())
		return;

	/* Set CSDECAYEN / VCCIO decays to 0V at assertion of SLP_S0# */
	i2c_write8(I2C_PORT_PMIC, I2C_ADDR_BD99992_FLAGS, 0x30, 0x4a);

	/*
	 * Set V100ACNT / V1.00A Control Register:
	 * Nominal output = 1.0V.
	 */
	i2c_write8(I2C_PORT_PMIC, I2C_ADDR_BD99992_FLAGS, 0x37, 0x1a);

	/*
	 * Set V085ACNT / V0.85A Control Register:
	 * Lower power mode = 0.7V.
	 * Nominal output = 1.0V.
	 */
	i2c_write8(I2C_PORT_PMIC, I2C_ADDR_BD99992_FLAGS, 0x38, 0x7a);

	/* VRMODECTRL - enable low-power mode for VCCIO and V0.85A */
	i2c_write8(I2C_PORT_PMIC, I2C_ADDR_BD99992_FLAGS, 0x3b, 0x18);
}
DECLARE_HOOK(HOOK_INIT, board_pmic_init, HOOK_PRIO_DEFAULT);

/* Initialize board. */
static void board_init(void)
{
	/* Enable PD MCU interrupt */
	gpio_enable_interrupt(GPIO_PD_MCU_INT);

	/* Enable VBUS interrupt */
	gpio_enable_interrupt(GPIO_USB_C0_VBUS_WAKE_L);
	gpio_enable_interrupt(GPIO_USB_C1_VBUS_WAKE_L);

	/* Enable pericom BC1.2 interrupts */
	gpio_enable_interrupt(GPIO_USB_C0_BC12_INT_L);
	gpio_enable_interrupt(GPIO_USB_C1_BC12_INT_L);

	/* Provide AC status to the PCH */
	gpio_set_level(GPIO_PCH_ACOK, extpower_is_present());

	/* Proto board workarounds */
	if (system_get_board_version() == 0) {
		/* Disable interrupt for SLP_S0 */
		gpio_set_flags(GPIO_PCH_SLP_S0_L,
			       GPIO_INPUT | GPIO_PULL_DOWN);

		/* Add internal pullup on PLATFORM_EC_PROCHOT */
		gpio_set_flags(GPIO_PLATFORM_EC_PROCHOT,
			       GPIO_INPUT | GPIO_PULL_UP);
	}
}
DECLARE_HOOK(HOOK_INIT, board_init, HOOK_PRIO_DEFAULT);

/**
 * Buffer the AC present GPIO to the PCH.
 */
static void board_extpower(void)
{
	gpio_set_level(GPIO_PCH_ACOK, extpower_is_present());
}
DECLARE_HOOK(HOOK_AC_CHANGE, board_extpower, HOOK_PRIO_DEFAULT);

/**
 * Set active charge port -- only one port can be active at a time.
 *
 * @param charge_port   Charge port to enable.
 *
 * Returns EC_SUCCESS if charge port is accepted and made active,
 * EC_ERROR_* otherwise.
 */
int board_set_active_charge_port(int charge_port)
{
	/* charge port is a realy physical port */
	int is_real_port = (charge_port >= 0 &&
			    charge_port < CONFIG_USB_PD_PORT_COUNT);
	/* check if we are source vbus on that port */
	int source = gpio_get_level(charge_port == 0 ? GPIO_USB_C0_5V_EN :
						       GPIO_USB_C1_5V_EN);

	if (is_real_port && source) {
		CPRINTS("Skip enable p%d", charge_port);
		return EC_ERROR_INVAL;
	}

	CPRINTS("New chg p%d", charge_port);

	if (charge_port == CHARGE_PORT_NONE) {
		/* Disable both ports */
		gpio_set_level(GPIO_USB_C0_CHARGE_EN_L, 1);
		gpio_set_level(GPIO_USB_C1_CHARGE_EN_L, 1);
	} else {
		/* Make sure non-charging port is disabled */
		gpio_set_level(charge_port ? GPIO_USB_C0_CHARGE_EN_L :
					     GPIO_USB_C1_CHARGE_EN_L, 1);
		/* Enable charging port */
		gpio_set_level(charge_port ? GPIO_USB_C1_CHARGE_EN_L :
					     GPIO_USB_C0_CHARGE_EN_L, 0);
	}

	return EC_SUCCESS;
}

/**
 * Set the charge limit based upon desired maximum.
 *
 * @param port          Port number.
 * @param supplier      Charge supplier type.
 * @param charge_ma     Desired charge limit (mA).
 * @param charge_mv     Negotiated charge voltage (mV).
 */
void board_set_charge_limit(int port, int supplier, int charge_ma,
			    int max_ma, int charge_mv)
{
	charge_set_input_current_limit(MAX(charge_ma,
				   CONFIG_CHARGER_INPUT_CURRENT), charge_mv);
}

/* Called on AP S5 -> S3 transition */
static void board_chipset_startup(void)
{
	gpio_set_level(GPIO_USB1_ENABLE, 1);
	gpio_set_level(GPIO_ENABLE_TOUCHPAD, 1);
}
DECLARE_HOOK(HOOK_CHIPSET_STARTUP, board_chipset_startup, HOOK_PRIO_DEFAULT);

/* Called on AP S3 -> S5 transition */
static void board_chipset_shutdown(void)
{
	gpio_set_level(GPIO_USB1_ENABLE, 0);
	gpio_set_level(GPIO_ENABLE_TOUCHPAD, 0);
}
DECLARE_HOOK(HOOK_CHIPSET_SHUTDOWN, board_chipset_shutdown, HOOK_PRIO_DEFAULT);

/* Called on AP S3 -> S0 transition */
static void board_chipset_resume(void)
{
	gpio_set_level(GPIO_ENABLE_BACKLIGHT, 1);
	gpio_set_level(GPIO_PP1800_DX_AUDIO_EN, 1);
	gpio_set_level(GPIO_PP1800_DX_DMIC_EN, 1);
}
DECLARE_HOOK(HOOK_CHIPSET_RESUME, board_chipset_resume, HOOK_PRIO_DEFAULT);

/* Called on AP S0 -> S3 transition */
static void board_chipset_suspend(void)
{
	gpio_set_level(GPIO_ENABLE_BACKLIGHT, 0);
	gpio_set_level(GPIO_PP1800_DX_AUDIO_EN, 0);
	gpio_set_level(GPIO_PP1800_DX_DMIC_EN, 0);
}
DECLARE_HOOK(HOOK_CHIPSET_SUSPEND, board_chipset_suspend, HOOK_PRIO_DEFAULT);

void board_hibernate(void)
{
	CPRINTS("Triggering PMIC shutdown.");
	uart_flush_output();

	/* Trigger PMIC shutdown. */
	if (i2c_write8(I2C_PORT_PMIC, I2C_ADDR_BD99992_FLAGS, 0x49, 0x01)) {
		/*
		 * If we can't tell the PMIC to shutdown, instead reset
		 * and don't start the AP. Hopefully we'll be able to
		 * communicate with the PMIC next time.
		 */
		CPRINTS("PMIC i2c failed.");
		system_reset(SYSTEM_RESET_LEAVE_AP_OFF);
	}

	/* Await shutdown. */
	while (1)
		;
}

/* Make the pmic re-sequence the power rails under these conditions. */
#define PMIC_RESET_FLAGS \
	(EC_RESET_FLAG_WATCHDOG | EC_RESET_FLAG_SOFT | EC_RESET_FLAG_HARD)
static void board_handle_reboot(void)
{
	int flags;

	if (system_jumped_to_this_image())
		return;

	/* Interrogate current reset flags from previous reboot. */
	flags = system_get_reset_flags();

	if (!(flags & PMIC_RESET_FLAGS))
		return;

	/* Preserve AP off request. */
	if (flags & EC_RESET_FLAG_AP_OFF)
		chip_save_reset_flags(EC_RESET_FLAG_AP_OFF);

	ccprintf("Restarting system with PMIC.\n");
	/* Flush console */
	cflush();

	/* Bring down all rails but RTC rail (including EC power). */
	gpio_set_level(GPIO_PMIC_LDO_EN, 1);
	while (1)
		; /* wait here */
}
DECLARE_HOOK(HOOK_INIT, board_handle_reboot, HOOK_PRIO_FIRST);

/*
 * Various voltage rails will be enabled / disabled by the PMIC when
 * GPIO_PMIC_SLP_SUS_L changes. We need to delay the disable of V0.85A
 * by approximately 25ms in order to allow V1.00A to sufficiently discharge
 * first.
 *
 * Therefore, after GPIO_PMIC_SLP_SUS_L goes high, ignore the state of
 * the V12_EN pin: Keep V0.85A enabled.
 *
 * When GPIO_PMIC_SLP_SUS_L goes low, delay 25ms, and make V12_EN function
 * as normal - this should result in V0.85A discharging immediately after the
 * i2c write completes.
 */
void chipset_set_pmic_slp_sus_l(int level)
{
	static int previous_level;
	int val;

	gpio_set_level(GPIO_PMIC_SLP_SUS_L, level);

	if (previous_level != level) {
		/* Rising edge: Force V0.85A enable. Falling: Pin control. */
		val = level ? 0x80 : 0;
		if (!level)
			msleep(25);

		i2c_write8(I2C_PORT_PMIC, I2C_ADDR_BD99992_FLAGS,
			   0x43, val);
		previous_level = level;
	}
}

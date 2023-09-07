/* Copyright 2019 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/* Hatch family-specific configuration */
#include "atomic.h"
#include "battery_fuel_gauge.h"
#include "charge_manager.h"
#include "charge_state_v2.h"
#include "chipset.h"
#include "console.h"
#include "cros_board_info.h"
#include "driver/ppc/sn5s330.h"
#include "driver/tcpm/anx7447.h"
#include "driver/tcpm/ps8xxx.h"
#include "driver/tcpm/tcpci.h"
#include "driver/tcpm/tcpm.h"
#include "espi.h"
#include "gpio.h"
#include "hooks.h"
#include "i2c.h"
#include "keyboard_scan.h"
#include "power.h"
#include "stdbool.h"
#include "tcpci.h"
#include "timer.h"
#include "usbc_ppc.h"
#include "util.h"

#define CPRINTS(format, args...) cprints(CC_SYSTEM, format, ## args)
#define CPRINTF(format, args...) cprintf(CC_SYSTEM, format, ## args)

#define CPRINTSUSB(format, args...) cprints(CC_USBCHARGE, format, ## args)
#define CPRINTFUSB(format, args...) cprintf(CC_USBCHARGE, format, ## args)

/******************************************************************************/
/* Wake up pins */
const enum gpio_signal hibernate_wake_pins[] = {
	GPIO_LID_OPEN,
	GPIO_ACOK_OD,
	GPIO_POWER_BUTTON_L,
	/* EC_RST_ODL needs to wake device while in PSL hibernate. */
	GPIO_SYS_RESET_L,
};
const int hibernate_wake_pins_used = ARRAY_SIZE(hibernate_wake_pins);

/******************************************************************************/
/* Keyboard scan setting */
struct keyboard_scan_config keyscan_config = {
	/*
	 * F3 key scan cycle completed but scan input is not
	 * charging to logic high when EC start scan next
	 * column for "T" key, so we set .output_settle_us
	 * to 80us from 50us.
	 */
	.output_settle_us = 80,
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

/******************************************************************************/
/* I2C port map configuration */
const struct i2c_port_t i2c_ports[] = {
	{"sensor",  I2C_PORT_SENSOR,  100, GPIO_I2C0_SCL, GPIO_I2C0_SDA},
	{"ppc0",    I2C_PORT_PPC0,    100, GPIO_I2C1_SCL, GPIO_I2C1_SDA},
	{"tcpc1",   I2C_PORT_TCPC1,   100, GPIO_I2C2_SCL, GPIO_I2C2_SDA},
	{"tcpc0",   I2C_PORT_TCPC0,   100, GPIO_I2C3_SCL, GPIO_I2C3_SDA},
	{"power",   I2C_PORT_POWER,   100, GPIO_I2C5_SCL, GPIO_I2C5_SDA},
	{"eeprom",  I2C_PORT_EEPROM,  100, GPIO_I2C7_SCL, GPIO_I2C7_SDA},
};
const unsigned int i2c_ports_used = ARRAY_SIZE(i2c_ports);

/******************************************************************************/
/* Chipset callbacks/hooks */

__attribute__((weak)) bool board_has_kb_backlight(void)
{
	/* Default enable keyboard backlight */
	return true;
}

/* Called on AP S5 -> S3 transition */
static void baseboard_chipset_startup(void)
{
	/* TODD(b/122266850): Need to fill out this hook */
}
DECLARE_HOOK(HOOK_CHIPSET_STARTUP, baseboard_chipset_startup,
	     HOOK_PRIO_DEFAULT);

/* Called on AP S0iX -> S0 transition */
static void baseboard_chipset_resume(void)
{
	if (board_has_kb_backlight())
		gpio_set_level(GPIO_EC_KB_BL_EN, 1);
}
DECLARE_HOOK(HOOK_CHIPSET_RESUME, baseboard_chipset_resume, HOOK_PRIO_DEFAULT);

/* Called on AP S0 -> S0iX transition */
static void baseboard_chipset_suspend(void)
{
	if (board_has_kb_backlight())
		gpio_set_level(GPIO_EC_KB_BL_EN, 0);
}
DECLARE_HOOK(HOOK_CHIPSET_SUSPEND, baseboard_chipset_suspend,
	     HOOK_PRIO_DEFAULT);

/* Called on AP S3 -> S5 transition */
static void baseboard_chipset_shutdown(void)
{

}
DECLARE_HOOK(HOOK_CHIPSET_SHUTDOWN, baseboard_chipset_shutdown,
	     HOOK_PRIO_DEFAULT);

void board_hibernate(void)
{
	int port;

	/*
	 * To support hibernate from ectool, keyboard, and console,
	 * ensure that the AP is fully shutdown before hibernating.
	 */
#ifdef HAS_TASK_CHIPSET
	chipset_force_shutdown(CHIPSET_SHUTDOWN_BOARD_CUSTOM);
#endif

	/*
	 * If VBUS is not being provided by any of the PD ports,
	 * then enable the SNK FET to allow AC to pass through
	 * if it is later connected to ensure that AC_PRESENT
	 * will wake up the EC from this state
	 */
	for (port = 0; port < CONFIG_USB_PD_PORT_COUNT; ++port)
		ppc_vbus_sink_enable(port, 1);

	/*
	 * This seems like a hack, but the AP chipset state machine
	 * needs time to work through the transitions.  Also, it
	 * works.
	 */
	msleep(300);
}

/******************************************************************************/
/* USB-C PPC Configuration */
struct ppc_config_t ppc_chips[CONFIG_USB_PD_PORT_COUNT] = {
	[USB_PD_PORT_TCPC_0] = {
		.i2c_port = I2C_PORT_PPC0,
		.i2c_addr_flags = SN5S330_ADDR0_FLAGS,
		.drv = &sn5s330_drv
	},

	[USB_PD_PORT_TCPC_1] = {
		.i2c_port = I2C_PORT_TCPC1,
		.i2c_addr_flags = SN5S330_ADDR0_FLAGS,
		.drv = &sn5s330_drv
	},
};
unsigned int ppc_cnt = ARRAY_SIZE(ppc_chips);

/* Power Delivery and charging functions */
void baseboard_tcpc_init(void)
{
	/* Enable PPC interrupts. */
	gpio_enable_interrupt(GPIO_USB_C0_PPC_INT_ODL);
	gpio_enable_interrupt(GPIO_USB_C1_PPC_INT_ODL);

	/* Enable TCPC interrupts. */
	gpio_enable_interrupt(GPIO_USB_C0_TCPC_INT_ODL);
	gpio_enable_interrupt(GPIO_USB_C1_TCPC_INT_ODL);

	/* Enable BC 1.2 interrupts */
	gpio_enable_interrupt(GPIO_USB_C0_BC12_INT_ODL);
	gpio_enable_interrupt(GPIO_USB_C1_BC12_INT_ODL);
}
DECLARE_HOOK(HOOK_INIT, baseboard_tcpc_init, HOOK_PRIO_INIT_I2C + 1);

uint16_t tcpc_get_alert_status(void)
{
	uint16_t status = 0;
	int level;

	/*
	 * Check which port has the ALERT line set and ignore if that TCPC has
	 * its reset line active.
	 */
	if (!gpio_get_level(GPIO_USB_C0_TCPC_INT_ODL)) {
		level = !!(tcpc_config[USB_PD_PORT_TCPC_0].flags &
			   TCPC_FLAGS_RESET_ACTIVE_HIGH);
		if (gpio_get_level(GPIO_USB_C0_TCPC_RST) != level)
			status |= PD_STATUS_TCPC_ALERT_0;
	}

	if (!gpio_get_level(GPIO_USB_C1_TCPC_INT_ODL)) {
		level = !!(tcpc_config[USB_PD_PORT_TCPC_1].flags &
			   TCPC_FLAGS_RESET_ACTIVE_HIGH);
		if (gpio_get_level(GPIO_USB_C1_TCPC_RST) != level)
			status |= PD_STATUS_TCPC_ALERT_1;
	}

	return status;
}

static void reset_pd_port(int port, enum gpio_signal reset_gpio,
			  int hold_delay, int finish_delay)
{
	int level = !!(tcpc_config[port].flags & TCPC_FLAGS_RESET_ACTIVE_HIGH);

	gpio_set_level(reset_gpio, level);
	msleep(hold_delay);
	gpio_set_level(reset_gpio, !level);
	if (finish_delay)
		msleep(finish_delay);
}

void board_reset_pd_mcu(void)
{
	/*
	 * TODO(b/130194590): This should be replaced with a common function
	 * once the gpio signal and delays are added to tcpc_config struct.
	 */

	/* Assert reset to TCPC for required delay only if we have a battery. */
	if (battery_is_present() != BP_YES)
		return;

	/* Reset TCPC0 */
	reset_pd_port(USB_PD_PORT_TCPC_0, GPIO_USB_C0_TCPC_RST,
		      BOARD_TCPC_C0_RESET_HOLD_DELAY,
		      BOARD_TCPC_C0_RESET_POST_DELAY);

	/* Reset TCPC1 */
	reset_pd_port(USB_PD_PORT_TCPC_1, GPIO_USB_C1_TCPC_RST,
		      BOARD_TCPC_C1_RESET_HOLD_DELAY,
		      BOARD_TCPC_C1_RESET_POST_DELAY);
}

int board_set_active_charge_port(int port)
{
	int is_valid_port = (port >= 0 &&
			    port < CONFIG_USB_PD_PORT_COUNT);
	int i;

	if (!is_valid_port && port != CHARGE_PORT_NONE)
		return EC_ERROR_INVAL;

	if (port == CHARGE_PORT_NONE) {
		CPRINTSUSB("Disabling all charger ports");

		/* Disable all ports. */
		for (i = 0; i < ppc_cnt; i++) {
			/*
			 * Do not return early if one fails otherwise we can
			 * get into a boot loop assertion failure.
			 */
			if (ppc_vbus_sink_enable(i, 0))
				CPRINTSUSB("Disabling C%d as sink failed.", i);
		}

		return EC_SUCCESS;
	}

	/* Check if the port is sourcing VBUS. */
	if (ppc_is_sourcing_vbus(port)) {
		CPRINTFUSB("Skip enable C%d", port);
		return EC_ERROR_INVAL;
	}

	CPRINTSUSB("New charge port: C%d", port);

	/*
	 * Turn off the other ports' sink path FETs, before enabling the
	 * requested charge port.
	 */
	for (i = 0; i < ppc_cnt; i++) {
		if (i == port)
			continue;

		if (ppc_vbus_sink_enable(i, 0))
			CPRINTSUSB("C%d: sink path disable failed.", i);
	}

	/* Enable requested charge port. */
	if (ppc_vbus_sink_enable(port, 1)) {
		CPRINTSUSB("C%d: sink path enable failed.", port);
		return EC_ERROR_UNKNOWN;
	}

	return EC_SUCCESS;
}

int ppc_get_alert_status(int port)
{
	if (port == USB_PD_PORT_TCPC_0)
		return gpio_get_level(GPIO_USB_C0_PPC_INT_ODL) == 0;
	else
		return gpio_get_level(GPIO_USB_C1_PPC_INT_ODL) == 0;
}

void board_set_charge_limit(int port, int supplier, int charge_ma,
			    int max_ma, int charge_mv)
{
	charge_set_input_current_limit(MAX(charge_ma,
					   CONFIG_CHARGER_INPUT_CURRENT),
				       charge_mv);
}

#ifdef USB_PD_PORT_TCPC_MST
void baseboard_mst_enable_control(enum mst_source src, int level)
{
	static uint32_t mst_input_levels;

	if (level)
		atomic_or(&mst_input_levels, 1 << src);
	else
		atomic_clear(&mst_input_levels, 1 << src);

	gpio_set_level(GPIO_EN_MST, mst_input_levels ? 1 : 0);
}
#endif

/* Enable or disable input devices, based on chipset state */
#ifndef TEST_BUILD
void lid_angle_peripheral_enable(int enable)
{
	/* TODO(b/125936966): Need to add SKU dependency for convertibles */
	if (chipset_in_state(CHIPSET_STATE_ANY_OFF))
		enable = 0;
	keyboard_scan_enable(enable, KB_SCAN_DISABLE_LID_ANGLE);
}
#endif

static uint8_t sku_id;

uint8_t get_board_sku(void)
{
	return sku_id;
}

/* Read CBI from i2c eeprom and initialize variables for board variants */
static void cbi_init(void)
{
	uint32_t val;

	if (cbi_get_sku_id(&val) != EC_SUCCESS || val > UINT8_MAX) {
		CPRINTS("Read SKU Error value :%d", val);
		return;
	}

	sku_id = val;

	CPRINTS("SKU: %d", sku_id);
}
DECLARE_HOOK(HOOK_INIT, cbi_init, HOOK_PRIO_INIT_I2C + 1);

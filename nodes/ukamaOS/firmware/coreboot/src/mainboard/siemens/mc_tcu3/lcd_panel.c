/*
* This file is part of the coreboot project.
 *
 * Copyright (C) 2014-2019 Siemens AG
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

#include <console/console.h>
#include <device/device.h>
#include <drivers/i2c/ptn3460/ptn3460.h>
#include <hwilib.h>
#include <string.h>
#include <types.h>

#include "soc/gpio.h"
#include "lcd_panel.h"

# define MAX_HWI_NAME_LENGTH 20

/** \brief Reads GPIOs used for LCD panel encoding and returns the 4 bit value
 * @param  no parameters
 * @return LCD panel type encoded in 4 bits
 */
static u8 get_lcd_panel_type(void)
{
	u8 lcd_type_gpio;

	lcd_type_gpio =  ((read_ssus_gpio(LCD_TYPE_GPIO_BIT3) << 3) |
			  (read_ssus_gpio(LCD_TYPE_GPIO_BIT2) << 2) |
			  (read_ssus_gpio(LCD_TYPE_GPIO_BIT1) << 1) |
			  (read_ssus_gpio(LCD_TYPE_GPIO_BIT0)));
	/* There is an inverter in this signals so we need to invert them as well */
	return ((~lcd_type_gpio) & 0x0f);
}
/** \brief This function checks which LCD panel type is used with the mainboard
 *         and provides the name of the matching EDID data set in CBFS.
 *  @param Pointer to the filename in CBFS where the EDID data is located
 *  @return CB_SUCCESS on success otherwise CB_ERR
 */
static enum cb_err get_hwi_filename(char *hwi_block)
{
	u8 lcd_type;
	enum cb_err ret = CB_SUCCESS;

	lcd_type = get_lcd_panel_type();
	printk(BIOS_INFO, "LCD: Found panel type %d\n", lcd_type);

	switch (lcd_type) {
	case LCD_PANEL_TYPE_10_INCH:
		strcpy(hwi_block, "hwinfo10.hex");
		break;
	case LCD_PANEL_TYPE_12_INCH:
		strcpy(hwi_block, "hwinfo12.hex");
		break;
	case LCD_PANEL_TYPE_15_INCH:
		strcpy(hwi_block, "hwinfo15.hex");
		break;
	case LCD_PANEL_TYPE_19_INCH:
		strcpy(hwi_block, "hwinfo19.hex");
		break;
	case LCD_PANEL_TYPE_EDID:
		strcpy(hwi_block, "hwinfo.hex");
		break;
	default:
		printk(BIOS_ERR, "LCD: No supported panel found.\n");
		ret = CB_ERR;
		break;
	}
	return ret;
}

/** \brief This function provides EDID data to the driver for DP2LVDS Bridge (PTN3460)
 * @param  edid_data	pointer to EDID data in driver
*/
enum cb_err mb_get_edid(uint8_t edid_data[0x80])
{
	char hwi_block[MAX_HWI_NAME_LENGTH];

	if (get_hwi_filename(hwi_block) != CB_SUCCESS)
		return CB_ERR;

	if (hwilib_find_blocks(hwi_block) != CB_SUCCESS) {
		printk(BIOS_ERR, "LCD: Info block \"%s\" not found!\n", hwi_block);
		return CB_ERR;
	}

	/* Get EDID data from hwinfo block */
	if (hwilib_get_field(Edid, edid_data, PTN_EDID_LEN) != PTN_EDID_LEN) {
		printk(BIOS_ERR, "LCD: No EDID data available in %s\n", hwi_block);
		return CB_ERR;
	}
	return CB_SUCCESS;
}

/** \brief This function provides EDID block [0..6] to the driver for DP2LVDS Bridge (PTN3460)
 * which has to be used.
*/
uint8_t mb_select_edid_table(void)
{
	return 6; /* With this mainboard we use EDID block 6 for emulation in PTN3460. */
}

/** \brief Function to enable mainboard to adjust the config data of PTN3460.
 * @param   *cfg_ptr  Pointer to the PTN config structure to modify.
 * @return  -1 on error; PTN_CFG_MODIFIED if data was modified and needs to be updated.
*/
int mb_adjust_cfg(struct ptn_3460_config *cfg)
{
	char hwi_block[MAX_HWI_NAME_LENGTH];
	uint8_t disp_con = 0, color_depth = 0;
	uint8_t hwid[4], tcu31_hwid[4] = {7, 9, 2, 0};

	if (get_hwi_filename(hwi_block) != CB_SUCCESS)
		return -1;
	if (hwilib_find_blocks(hwi_block) != CB_SUCCESS) {
		printk(BIOS_ERR, "LCD: Info block \"%s\" not found!\n", hwi_block);
		return -1;
	}

	if (hwilib_get_field(PF_DisplCon, &disp_con, sizeof(disp_con)) != sizeof(disp_con)) {
		printk(BIOS_ERR, "LCD: Missing panel features from %s\n", hwi_block);
		return -1;
	}
	if (hwilib_get_field(PF_Color_Depth, &color_depth,
			     sizeof(color_depth)) != sizeof(color_depth)) {
		printk(BIOS_ERR, "LCD: Missing panel features from %s\n", hwi_block);
		return -1;
	}
	/* Set up configuration data according to the hwinfo block we got. */
	cfg->dp_interface_ctrl = 0x00;
	cfg->lvds_interface_ctrl1 = 0x00;
	if (disp_con == PF_DISPLCON_LVDS_DUAL) {
		/* Turn on dual LVDS lane and clock. */
		cfg->lvds_interface_ctrl1 |= 0x0b;
	}
	if (color_depth == PF_COLOR_DEPTH_6BIT) {
		/* Use 18 bits per pixel. */
		cfg->lvds_interface_ctrl1 |= 0x20;
	}
	/* No clock spreading, 300 mV LVDS swing. */
	cfg->lvds_interface_ctrl2 = 0x03;
	/* Swap LVDS even and odd lanes for HW-ID 7.9.2.0 only. */
	if (hwilib_get_field(HWID, hwid, sizeof(hwid)) == sizeof(hwid) &&
	    !(memcmp(hwid, tcu31_hwid, sizeof(hwid)))) {
		/* Swap LVDS even and odd lane. */
		cfg->lvds_interface_ctrl3 = 0x01;
	} else {
		 /* no LVDS lane swap */
		cfg->lvds_interface_ctrl3 = 0x00;
	}
	/* Delay T2 (VDD to LVDS active) by 16 ms. */
	cfg->t2_delay = 1;
	/* 500 ms from LVDS to backlight active. */
	cfg->t3_timing = 10;
	/* 1 second re-power delay. */
	cfg->t12_timing = 20;
	/* 150 ms backlight off to LVDS inactive. */
	cfg->t4_timing = 3;
	/* Delay T5 (LVDS to VDD inactive) by 16 ms. */
	cfg->t5_delay = 1;
	/* Enable backlight control. */
	cfg->backlight_ctrl = 0;

	return PTN_CFG_MODIFIED;
}

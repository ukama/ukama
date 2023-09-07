/*
 * This file is part of the coreboot project.
 *
 * Copyright 2018 MediaTek Inc.
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

#include <assert.h>
#include <bl31.h>
#include <boardid.h>
#include <bootmode.h>
#include <cbfs.h>
#include <console/console.h>
#include <delay.h>
#include <device/device.h>
#include <edid.h>
#include <gpio.h>
#include <soc/ddp.h>
#include <soc/dsi.h>
#include <soc/gpio.h>
#include <soc/mmu_operations.h>
#include <soc/mtcmos.h>
#include <soc/spm.h>
#include <soc/usb.h>
#include <string.h>

#include "gpio.h"
#include "panel.h"

#include <arm-trusted-firmware/include/export/plat/mediatek/common/plat_params_exp.h>

static void configure_emmc(void)
{
	const gpio_t emmc_pin[] = {
		GPIO(MSDC0_DAT0), GPIO(MSDC0_DAT1),
		GPIO(MSDC0_DAT2), GPIO(MSDC0_DAT3),
		GPIO(MSDC0_DAT4), GPIO(MSDC0_DAT5),
		GPIO(MSDC0_DAT6), GPIO(MSDC0_DAT7),
		GPIO(MSDC0_CMD), GPIO(MSDC0_RSTB),
	};

	for (size_t i = 0; i < ARRAY_SIZE(emmc_pin); i++)
		gpio_set_pull(emmc_pin[i], GPIO_PULL_ENABLE, GPIO_PULL_UP);
}

static void configure_usb(void)
{
	setup_usb_host();
}

static void configure_audio(void)
{
	/* Audio PWR*/
	mtcmos_audio_power_on();

	/* SoC I2S */
	gpio_set_mode(GPIO(CAM_RST0), PAD_CAM_RST0_FUNC_I2S2_LRCK);
	gpio_set_mode(GPIO(CAM_PDN1), PAD_CAM_PDN1_FUNC_I2S2_BCK);
	gpio_set_mode(GPIO(CAM_PDN0), PAD_CAM_PDN0_FUNC_I2S2_MCK);
	gpio_set_mode(GPIO(EINT3), PAD_EINT3_FUNC_I2S3_DO);
}

/* Default implementation for boards without panels defined yet. */
struct panel_description __weak *get_panel_description(int panel_id)
{
	printk(BIOS_ERR, "%s: ERROR: No panels defined for board: %s.\n",
	       __func__, CONFIG_MAINBOARD_PART_NUMBER);
	return NULL;
}

/* Set up backlight control pins as output pin and power-off by default */
static void configure_panel_backlight(void)
{
	gpio_output(GPIO(PERIPHERAL_EN13), 0);
	gpio_output(GPIO(DISP_PWM), 0);
}

static void power_on_panel(struct panel_description *panel)
{
	if (panel->power_on) {
		panel->power_on();
		return;
	}

	/* Default power sequence for most panels. */
	gpio_output(GPIO_LCM_RST_1V8, 0);
	gpio_output(GPIO_PPVARP_LCD_EN, 1);
	gpio_output(GPIO_PPVARN_LCD_EN, 1);
	gpio_output(GPIO_PP1800_LCM_EN, 1);
	gpio_output(GPIO_PP3300_LCM_EN, 1);
	mdelay(6);
	gpio_output(GPIO_LCM_RST_1V8, 1);
	mdelay(6);
}

struct panel_description *get_panel_from_cbfs(struct panel_description *desc)
{
	/* The CBFS name will be panel-{MANUFACTURER}-${PANEL_NAME},
	 * where MANUFACTURER is 3 characters and PANEL_NAME is usually
	 * 13 characters.
	 */
	char cbfs_name[64];
	static union {
		u8 raw[4 * 1024];  /* Most panels only need < 2K. */
		struct panel_serializable_data s;
	} buffer;

	if (!desc->name)
		return NULL;

	snprintf(cbfs_name, sizeof(cbfs_name), "panel-%s", desc->name);
	if (cbfs_boot_load_file(cbfs_name, buffer.raw, sizeof(buffer),
				CBFS_TYPE_STRUCT))
		desc->s = &buffer.s;
	else
		printk(BIOS_ERR, "Missing %s in CBFS.\n", cbfs_name);

	return desc->s ? desc : NULL;
}

static struct panel_description *get_active_panel(void)
{
	/* TODO(hungte) Create a dedicated panel_id() in board_id.c */
	int panel_id = sku_id() >> 4;

	struct panel_description *panel = get_panel_description(panel_id);
	if (!panel) {
		printk(BIOS_ERR, "%s: Panel %d is not supported.\n",
		       __func__, panel_id);
		return NULL;
	}
	assert(panel->s);

	const struct edid *edid = &panel->s->edid;
	const char *name = edid->ascii_string;
	if (name[0] == '\0')
		name = "unknown name";
	printk(BIOS_INFO, "%s: Found ID %d: '%s %s' %dx%d@%dHz\n", __func__,
	       panel_id, edid->manufacturer_name, name, edid->mode.ha,
	       edid->mode.va, edid->mode.refresh);
	return panel;
}

static bool configure_display(void)
{
	struct panel_description *panel = get_active_panel();
	if (!panel)
		return false;

	mtcmos_display_power_on();
	mtcmos_protect_display_bus();
	configure_panel_backlight();
	power_on_panel(panel);

	struct edid *edid = &panel->s->edid;
	edid_set_framebuffer_bits_per_pixel(edid, 32, 0);
	mtk_ddp_init();
	u32 mipi_dsi_flags = (MIPI_DSI_MODE_VIDEO |
			      MIPI_DSI_MODE_VIDEO_SYNC_PULSE |
			      MIPI_DSI_MODE_LPM);
	if (mtk_dsi_init(mipi_dsi_flags, MIPI_DSI_FMT_RGB888, 4, edid,
			 panel->s->init) < 0) {
		printk(BIOS_ERR, "%s: Failed in DSI init.\n", __func__);
		return false;
	}
	mtk_ddp_mode_set(edid);
	set_vbe_mode_info_valid(edid, 0);
	set_vbe_framebuffer_orientation(panel->s->orientation);
	return true;
}

static void register_reset_to_bl31(void)
{
	static struct bl_aux_param_gpio param_reset = {
		.h = { .type = BL_AUX_PARAM_MTK_RESET_GPIO },
		.gpio = { .polarity = ARM_TF_GPIO_LEVEL_HIGH },
	};

	param_reset.gpio.index = GPIO_RESET.id;
	register_bl31_aux_param(&param_reset.h);
}

static void mainboard_init(struct device *dev)
{
	if (display_init_required()) {
		printk(BIOS_INFO, "%s: Starting display init.\n", __func__);
		if (!configure_display())
			printk(BIOS_ERR, "%s: Failed to init display.\n",
			       __func__);
	} else {
		printk(BIOS_INFO, "%s: Skipped display init.\n", __func__);
	}

	configure_emmc();
	configure_usb();
	configure_audio();
	if (spm_init())
		printk(BIOS_ERR,
		       "SPM initialization failed, suspend/resume may fail.\n");

	register_reset_to_bl31();
}

static void mainboard_enable(struct device *dev)
{
	dev->ops->init = &mainboard_init;
}

struct chip_operations mainboard_ops = {
	.name = CONFIG_MAINBOARD_PART_NUMBER,
	.enable_dev = mainboard_enable,
};

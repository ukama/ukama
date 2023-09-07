/*
 * This file is part of the coreboot project.
 *
 * Copyright 2017 Intel Corp.
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

#ifndef BASEBOARD_VARIANTS_H
#define BASEBOARD_VARIANTS_H

#include <soc/gpio.h>
#include <soc/meminit.h>
#include <stdint.h>
#include <vendorcode/google/chromeos/chromeos.h>

/* The next set of functions return the gpio table and fill in the number of
 * entries for each table. */
const struct pad_config *variant_base_gpio_table(size_t *num);
const struct pad_config *variant_override_gpio_table(size_t *num);
const struct pad_config *variant_early_gpio_table(size_t *num);
const struct pad_config *variant_early_override_gpio_table(size_t *num);
const struct pad_config *variant_sleep_gpio_table(size_t *num, int slp_typ);

/* Baseboard default swizzle. Can be reused if swizzle is same. */
extern const struct lpddr4_swizzle_cfg baseboard_lpddr4_swizzle;
/* Return LPDDR4 configuration structure. */
const struct lpddr4_cfg *variant_lpddr4_config(void);
/* Return memory SKU for the board. */
size_t variant_memory_sku(void);
/* Return board SKU */
uint32_t get_board_sku(void);
/* Return ChromeOS gpio table and fill in number of entries. */
const struct cros_gpio *variant_cros_gpios(size_t *num);

/* Seed the NHLT tables with the board specific information. */
struct nhlt;
void variant_nhlt_init(struct nhlt *nhlt);

/* Modify devictree settings during ramstage. */
struct device;
void variant_update_devtree(struct device *dev);
/**
 * variant_ext_usb_status() - Get status of externally visible USB ports
 * @port_type: Type of USB port i.e. USB2/USB3
 * @port_id: USB Port ID
 *
 * This function is supplied by the mainboard/variant to SoC's XHCI driver to
 * identify the status of externally visible USB ports.
 *
 * Return: true if the port is present, false if the port is absent.
 */
bool variant_ext_usb_status(unsigned int port_type, unsigned int port_id);

/* Get no touchscreen SKU ID. */
bool no_touchscreen_sku(uint32_t sku_id);

/* allow each variants to customize smi sleep flow. */
void variant_smi_sleep(u8 slp_typ);

#endif /* BASEBOARD_VARIANTS_H */

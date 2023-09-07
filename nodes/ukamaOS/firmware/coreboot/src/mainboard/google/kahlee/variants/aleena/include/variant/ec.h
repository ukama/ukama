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

#include <baseboard/ec.h>

/* Enable EC backed Keyboard Backlight in ACPI */
#define EC_ENABLE_KEYBOARD_BACKLIGHT

/* Enable Tablet switch */
#define EC_ENABLE_TBMC_DEVICE

/*
 * Enable EC sync interrupt via GPIO controller, EC_SYNC_IRQ is defined in
 * variant/gpio.h
 */
#define EC_ENABLE_SYNC_IRQ_GPIO

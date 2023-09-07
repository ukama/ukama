/*
 * This file is part of the coreboot project.
 *
 * Copyright 2019 Google Inc.
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

#include <soc/mt8183.h>
#include <soc/wdt.h>
#include <soc/gpio.h>

void mt8183_early_init(void)
{
	mtk_wdt_init();
	gpio_set_i2c_eh_rsel();
}

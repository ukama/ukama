/*
 * This file is part of the coreboot project.
 *
 * Copyright 2017-present Facebook, Inc.
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

#ifndef __SOC_CAVIUM_CN81XX_TIMER_H__
#define __SOC_CAVIUM_CN81XX_TIMER_H__

#include <stdint.h>
#include <types.h>

/* Watchdog functions */
void watchdog_set(const size_t index, unsigned int timeout_ms);
void watchdog_poke(const size_t index);
void watchdog_disable(const size_t index);
int watchdog_is_running(const size_t index);

/* Timer functions */
void soc_timer_init(void);

#endif	/* __SOC_CAVIUM_CN81XX_TIMER_H__ */

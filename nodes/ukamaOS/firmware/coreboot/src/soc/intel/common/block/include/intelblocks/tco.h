/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2017-2018 Intel Corporation.
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

#ifndef SOC_INTEL_COMMON_BLOCK_TCO_H
#define SOC_INTEL_COMMON_BLOCK_TCO_H

#include <stdint.h>

/*
 * Enable TCO BAR using SMBUS TCO base to access TCO related register
 * also disable the timer.
 */
void tco_configure(void);
/* Lockdown TCO registers before booting to OS */
void tco_lockdown(void);
/*
 * Resets the tco status registers. This function clears the tco_sts register
 * and returns the status bits set.
 */
uint32_t tco_reset_status(void);
uint16_t tco_read_reg(uint16_t tco_reg);
void tco_write_reg(uint16_t tco_reg, uint16_t value);

#endif /* SOC_INTEL_COMMON_BLOCK_TCO_H */

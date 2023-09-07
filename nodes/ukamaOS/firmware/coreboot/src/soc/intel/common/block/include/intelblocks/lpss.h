/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2017 Intel Corporation.
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

#ifndef SOC_INTEL_COMMON_BLOCK_LPSS_H
#define SOC_INTEL_COMMON_BLOCK_LPSS_H

#include <device/device.h>
#include <stdint.h>

/* D0 and D3 enable config */
enum lpss_pwr_state {
	STATE_D0 = 0,
	STATE_D3 = 3
};

/* Gets controller out of reset */
void lpss_reset_release(uintptr_t base);

/*
 * Update clock divider parameters. Clock frequency is dependent on source
 * clock frequency of each IP block. Resulting clock will be src_freq * (M / N).
 */
void lpss_clk_update(uintptr_t base, uint32_t clk_m_val, uint32_t clk_n_val);

/* Check if controller is in reset. */
bool lpss_is_controller_in_reset(uintptr_t base);

/* Set controller power state to D0 or D3*/
void lpss_set_power_state(const struct device *dev, enum lpss_pwr_state state);

/*
 * Handler to get list of LPSS controllers. The SOC is expected to send out a
 * list of pci devfn for all LPSS controllers supported by the SOC.
 */
const pci_devfn_t *soc_lpss_controllers_list(size_t *size);

/* Check if the device is a LPSS controller */
bool is_dev_lpss(const struct device *dev);

#endif	/* SOC_INTEL_COMMON_BLOCK_LPSS_H */

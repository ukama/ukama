/*
 * Copyright (c) 2019, Arm Limited and Contributors. All rights reserved.
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

#include <common/bl_common.h>
#include <plat/arm/common/plat_arm.h>
#include <plat/common/platform.h>
#include <platform_def.h>

/*
 * Table of regions to map using the MMU.
 * Replace or extend the below regions as required
 */

const mmap_region_t plat_arm_mmap[] = {
	ARM_MAP_SHARED_RAM,
	ARM_MAP_NS_DRAM1,
	CORSTONE700_MAP_DEVICE,
	{0}
};

/* Corstone700 only has one always-on power domain and there
 * is no power control present
 */
void __init plat_arm_pwrc_setup(void)
{
}

unsigned int plat_get_syscnt_freq2(void)
{
	return CORSTONE700_TIMER_BASE_FREQUENCY;
}

/*
 * This file is part of the coreboot project.
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation; version 2 of
 * the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#ifndef _CPU_INTEL_COMMON_H
#define _CPU_INTEL_COMMON_H

#include <stdint.h>

void set_vmx_and_lock(void);
void set_feature_ctrl_vmx(void);
void set_feature_ctrl_lock(void);

/*
 * Init CPPC block with MSRs for Intel Enhanced Speed Step Technology.
 * Version 2 is suggested--this function's implementation of version 3
 * may have room for improvment.
 */
struct cppc_config;
void cpu_init_cppc_config(struct cppc_config *config, u32 version);

/*
 * Returns true if it's not thread 0 on a hyperthreading enabled core.
 */
bool intel_ht_sibling(void);

#endif

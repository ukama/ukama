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

#include <stdint.h>
#include <arch/cpu.h>
#include <cpu/x86/msr.h>
#include <cpu/intel/speedstep.h>
#include "model_2065x.h"

/* MSR Documentation based on
 * "Sandy Bridge Processor Family BIOS Writer's Guide (BWG)"
 * Document Number 504790
 * Revision 1.6.0, June 2012 */

void intel_model_2065x_finalize_smm(void)
{
	/* Lock C-State MSR */
	msr_set_bit(MSR_PKG_CST_CONFIG_CONTROL, 15);

	/* Lock AES-NI only if supported */
	if (cpuid_ecx(1) & (1 << 25))
		msr_set_bit(MSR_FEATURE_CONFIG, 0);

	/* Lock TM interrupts - route thermal events to all processors */
	msr_set_bit(MSR_MISC_PWR_MGMT, 22);
}

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

#include <cpu/x86/lapic.h>
#include <cpu/intel/common/common.h>
#include <arch/cpu.h>

/*
 * Return true if running thread does not have the smallest lapic ID
 * within a CPU core.
 */
bool intel_ht_sibling(void)
{
	struct cpuid_result result;
	unsigned int core_ids, apic_ids, threads;

	/* Is Hyper-Threading supported */
	if (!(cpuid_edx(1) & CPUID_FEAURE_HTT))
		return false;

	apic_ids = 1;
	if (cpuid_eax(0) >= 1)
		apic_ids = (cpuid_ebx(1) >> 16) & 0xff;
	if (apic_ids == 0)
		apic_ids = 1;

	core_ids = 1;
	if (cpuid_eax(0) >= 4) {
		result = cpuid_ext(4, 0);
		core_ids += (result.eax >> 26) & 0x3f;
	}

	threads = (apic_ids / core_ids);
	return !!(lapicid() & (threads - 1));
}

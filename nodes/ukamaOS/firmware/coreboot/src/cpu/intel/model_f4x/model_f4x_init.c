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

#include <device/device.h>
#include <cpu/cpu.h>
#include <cpu/x86/lapic.h>
#include <cpu/x86/cache.h>

static void model_f4x_init(struct device *cpu)
{
	/* Turn on caching if we haven't already */
	x86_enable_cache();

	/* Enable the local CPU APICs */
	setup_lapic();
};

static struct device_operations cpu_dev_ops = {
	.init = model_f4x_init,
};

static const struct cpu_device_id cpu_table[] = {
	{ X86_VENDOR_INTEL, 0x0f41 }, /* Xeon */
	{ X86_VENDOR_INTEL, 0x0f43 }, /* Not tested */
	{ X86_VENDOR_INTEL, 0x0f44 }, /* Not tested */
	{ X86_VENDOR_INTEL, 0x0f47 },
	{ X86_VENDOR_INTEL, 0x0f48 }, /* Not tested */
	{ X86_VENDOR_INTEL, 0x0f49 }, /* Not tested */
	{ X86_VENDOR_INTEL, 0x0f4a }, /* Not tested */
	{ 0, 0 },
};

static const struct cpu_driver model_f4x __cpu_driver = {
	.ops      = &cpu_dev_ops,
	.id_table = cpu_table,
};

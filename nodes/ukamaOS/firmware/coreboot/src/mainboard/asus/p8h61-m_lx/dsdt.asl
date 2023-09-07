/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2017 Tristan Corrick <tristan@corrick.kiwi>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <arch/acpi.h>
DefinitionBlock(
	"dsdt.aml",
	"DSDT",
	0x02,		/* DSDT Revision: ACPI v2.0 and up */
	OEM_ID,
	ACPI_TABLE_CREATOR,
	0x20171231	/* OEM Revision */
)
{
	#include "acpi/platform.asl"
	#include <cpu/intel/common/acpi/cpu.asl>
	#include <southbridge/intel/common/acpi/platform.asl>
	#include <southbridge/intel/bd82x6x/acpi/globalnvs.asl>
	#include <southbridge/intel/common/acpi/sleepstates.asl>

	Scope (\_SB)
	{
		Device (PCI0)
		{
		#include <northbridge/intel/sandybridge/acpi/sandybridge.asl>
		#include <drivers/intel/gma/acpi/default_brightness_levels.asl>
		#include <southbridge/intel/bd82x6x/acpi/pch.asl>
		}
	}
}

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
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 */

#include <soc/pcr_ids.h>

Scope (\_SB.PCI0) {

	/*
	 * Clear register 0x1C20/0x4820
	 * Arg0 - PCR Port ID
	 */
	Method(SCSC, 1, Serialized)
	{
		^PCRA (Arg0, 0x1C20, 0x0)
		^PCRA (Arg0, 0x4820, 0x0)
	}

	/* EMMC */
	Device(PEMC) {
		Name(_ADR, 0x001A0000)
		Name (_DDN, "eMMC Controller")
		Name (TEMP, 0)

		OperationRegion(SCSR, PCI_Config, 0x00, 0x100)
		Field(SCSR, WordAcc, NoLock, Preserve) {
			Offset (0x84),	/* PMECTRLSTATUS */
			PMCR, 16,
			Offset (0xA2),	/* PG_CONFIG */
			, 2,
			PGEN, 1,	/* PG_ENABLE */
		}

		Method(_INI) {
			/* Clear register 0x1C20/0x4820 */
			^^SCSC (PID_EMMC)
		}

		Method(_PS0, 0, Serialized) {
			Stall (50) // Sleep 50 us

			Store(0, PGEN) // Disable PG

			/* Clear register 0x1C20/0x4820 */
			^^SCSC (PID_EMMC)

			/* Set Power State to D0 */
			And (PMCR, 0xFFFC, PMCR)
			Store (PMCR, ^TEMP)
		}

		Method(_PS3, 0, Serialized) {
			Store(1, PGEN) // Enable PG

			/* Set Power State to D3 */
			Or (PMCR, 0x0003, PMCR)
			Store (PMCR, ^TEMP)
		}

		Device (CARD)
		{
			Name (_ADR, 0x00000008)
			Method (_RMV, 0, NotSerialized)
			{
				Return (0)
			}
		}
	}

	/* SD CARD */
	Device (SDXC)
	{
		Name (_ADR, 0x00140005)
		Name (_DDN, "SD Controller")
		Name (TEMP, 0)
		Name (DSUU, ToUUID("f6c13ea5-65cd-461f-ab7a-29f7e8d5bd61"))

		OperationRegion (SDPC, PCI_Config, 0x00, 0x100)
		Field (SDPC, WordAcc, NoLock, Preserve)
		{
			Offset (0x84),	/* PMECTRLSTATUS */
			PMCR, 16,
			Offset (0xA2),	/* PG_CONFIG */
			, 2,
			PGEN, 1,	/* PG_ENABLE */
		}

		/* _DSM x86 Device Specific Method
		 * Arg0: UUID Unique function identifier
		 * Arg1: Integer Revision Level
		 * Arg2: Integer Function Index (0 = Return Supported Functions)
		 * Arg3: Package Parameters
		 */
		Method (_DSM, 4)
		{
			If (LEqual (Arg0, ^DSUU)) {
				/* Check the revision */
				If (LGreaterEqual (Arg1, Zero)) {
					/* Switch statement based on the function index. */
					Switch (ToInteger (Arg2)) {
					/*
					 * Function Index 0 the return value is a buffer containing
					 * one bit for each function index, starting with zero.
					 * Bit 0 - Indicates whether there is support for any functions other than function 0.
					 * Bit 1 - Indicates support to clear power control register
					 * Bit 2 - Indicates support to set power control register
					 * Bit 3 - Indicates support to set 1.8V signalling
					 * Bit 4 - Indicates support to set 3.3V signalling
					 * Bit 5 - Indicates support for HS200 mode
					 * Bit 6 - Indicates support for HS400 mode
					 * Bit 9 - Indicates eMMC I/O Driver Strength
					 */
					/*
					 * For SD we have to support functions to
					 * set 1.8V signalling and 3.3V signalling [BIT4, BIT3]
					 */
					Case (0) {
						Return (Buffer () { 0x19 })
					}

					/*
					 * Function Index 3: Set 1.8v signalling.
					 * We put a sleep of 100ms in this method to
					 * work around a known issue with detecting
					 * UHS SD card on PCH. This is to compensate
					 * for the SD VR slowness.
					 */
					Case (3) {
						Sleep (100)
						Return(Buffer () { 0x00 })
					}
					/*
					 * Function Index 4: Set 3.3v signalling.
					 * We put a sleep of 100ms in this method to
					 * work around a known issue with detecting
					 * UHS SD card on PCH. This is to compensate
					 * for the SD VR slowness.
					 */
					Case (4) {
						Sleep (100)
						Return(Buffer () { 0x00 })
					}
					}
				}
			}
			Return(Buffer() { 0x0 })
		}

		Method(_INI)
		{
			/* Clear register 0x1C20/0x4820 */
			^^SCSC (PID_SDX)
		}

		Method (_PS0, 0, Serialized)
		{
			Store (0, PGEN) /* Disable PG */

			/* Clear register 0x1C20/0x4820 */
			^^SCSC (PID_SDX)

			/* Set Power State to D0 */
			And (PMCR, 0xFFFC, PMCR)
			Store (PMCR, ^TEMP)

#if CONFIG(MB_HAS_ACTIVE_HIGH_SD_PWR_ENABLE)
			/* Change pad mode to Native */
			GPMO(SD_PWR_EN_PIN, 0x1)
#endif
		}

		Method (_PS3, 0, Serialized)
		{
			Store (1, PGEN) /* Enable PG */

			/* Set Power State to D3 */
			Or (PMCR, 0x0003, PMCR)
			Store (PMCR, ^TEMP)

#if CONFIG(MB_HAS_ACTIVE_HIGH_SD_PWR_ENABLE)
			/* Change pad mode to GPIO control */
			GPMO(SD_PWR_EN_PIN, 0x0)

			/* Enable Tx Buffer */
			GTXE(SD_PWR_EN_PIN, 0x1)

			/* Drive TX to zero */
			CTXS(SD_PWR_EN_PIN)
#endif
		}

		Device (CARD)
		{
			Name (_ADR, 0x00000008)
			Method (_RMV, 0, NotSerialized)
			{
				Return (1)
			}
		}
	} /* Device (SDXC) */
}

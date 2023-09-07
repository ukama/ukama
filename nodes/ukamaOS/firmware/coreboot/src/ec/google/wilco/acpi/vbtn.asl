/*
 * This file is part of the coreboot project.
 *
 * Copyright 2019 Google LLC
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

/*
 * Intel Virtual Button driver compatible with the driver found in
 * the Linux kernel at drivers/platform/x86/intel-vbtn.c
 *
 * For tablet/laptop and dock/undock events to work the board must
 * select SYSTEM_TYPE_CONVERTIBLE for the SMBIOS enclosure type to
 * indicate the device is a convertible.
 */

Name (FLAP, 0x40) /* Flag indicating device is in laptop mode */

/* Virtual events */
Name (VPPB, 0xc0) /* Power Button press */
Name (VRPB, 0xc1) /* Power Button release */
Name (VPSP, 0xc2) /* Special key press (LEFTMETA in Linux) */
Name (VRSP, 0xc3) /* Special key release (LEFTMETA in Linux) */
Name (VPVU, 0xc4) /* Volume Up press */
Name (VRVU, 0xc5) /* Volume Up release */
Name (VPVD, 0xc6) /* Volume Down press */
Name (VRVD, 0xc7) /* Volume Down release */
Name (VPRL, 0xc8) /* Rotate Lock press */
Name (VRRL, 0xc9) /* Rotate Lock release */
Name (VDOC, 0xca) /* Docked */
Name (VUND, 0xcb) /* Undocked */
Name (VTBL, 0xcc) /* Tablet Mode */
Name (VLAP, 0xcd) /* Laptop Mode */

Device (VBTN)
{
	Name (_HID, "INT33D6")
	Name (_UID, One)
	Name (_DDN, "Intel Virtual Button Driver")

	/*
	 * This method is called at driver probe time and must exist or
	 * the driver will not load.
	 */
	Method (VBDL)
	{
	}

	/*
	 * This method returns flags indicating tablet and dock modes.
	 * It is called at driver probe time so the OS knows what the
	 * state of the device is at boot.
	 */
	Method (VGBS)
	{
		Local0 = Zero

		/* Check EC orientation for tablet mode flag */
		If (R (OTBL)) {
			Printf ("EC reports tablet mode at boot")
		} Else {
			Printf ("EC reports laptop mode at boot")
			Local0 |= ^^FLAP
		}
		Return (Local0)
	}

	Method(_STA, 0)
	{
		Return (0xF)
	}
}

Device (VBTO)
{
	Name (_HID, "INT33D3")
	Name (_CID, "PNP0C60")
	Name (_UID, One)
	Name (_DDN, "Laptop/tablet mode indicator driver")

	Method (_STA, 0)
	{
		Return (0xF)
	}
}

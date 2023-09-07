/*
 * This file is part of the coreboot project.
 *
 * Copyright 2018 Google LLC
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

/* Scope is \_SB.PCI0.LPCB */

Device (SIO)
{
	Name (_UID, 0)
	Name (_ADR, 0)

	Device (COM1)
	{
		Name (_HID, EisaId ("PNP0501"))
		Name (_UID, 1)

		Method (_STA, 0, NotSerialized)
		{
#if CONFIG(DRIVERS_UART_8250IO)
			Return (0x0f)
#else
			Return (Zero)
#endif
		}

		Name (_CRS, ResourceTemplate ()
		{
			IO (Decode16, 0x03F8, 0x3F8, 0x08, 0x08)
			IRQNoFlags () {4}
		})

		Name (_PRS, ResourceTemplate ()
		{
			StartDependentFn (0, 0) {
				IO (Decode16, 0x03F8, 0x3F8, 0x08, 0x08)
				IRQNoFlags () {4}
			}
			EndDependentFn ()
		})
	}

	Device (PS2K)
	{
		Name (_HID, EisaId ("PNP0303"))
		Name (_CID, EisaId ("PNP030B"))
		Name (_UID, 0)

		Method (_STA, 0, NotSerialized)
		{
#ifdef SIO_EC_ENABLE_PS2K
			Return (0x0f)
#else
			Return (Zero)
#endif
		}

		Name (_CRS, ResourceTemplate()
		{
			IO (Decode16, 0x60, 0x60, 0x01, 0x01)
			IO (Decode16, 0x64, 0x64, 0x01, 0x01)
			IRQ (Edge, ActiveHigh, Exclusive) {1}
		})

		Name (_PRS, ResourceTemplate()
		{
			StartDependentFn (0, 0) {
				IO (Decode16, 0x60, 0x60, 0x01, 0x01)
				IO (Decode16, 0x64, 0x64, 0x01, 0x01)
				IRQ (Edge, ActiveHigh, Exclusive) {1}
			}
			EndDependentFn ()
		})
	}

	Device (PS2M)
	{
		Name (_HID, EisaId ("PNP0F13"))
		Name (_UID, 0)

		Method (_STA, 0, NotSerialized)
		{
#ifdef SIO_EC_ENABLE_PS2M
			Return (0x0f)
#else
			Return (Zero)
#endif
		}

		Name (_CRS, ResourceTemplate()
		{
			IRQ (Edge, ActiveHigh, Exclusive) {12}
		})

		Name (_PRS, ResourceTemplate()
		{
			StartDependentFn (0, 0) {
				IRQ (Edge, ActiveHigh, Exclusive) {12}
			}
			EndDependentFn ()
		})
	}
}

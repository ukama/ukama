/*
 * This file is part of the coreboot project.
 *
 * Copyright (c) 2017 Arthur Heymans <arthur@aheymans.xyz>
 * Copyright (c) 2018 Evgeny Zinoviev <me@ch1p.com>
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


Scope(\_SB.PCI0.LPCB.EC)
{
	Field (ERAM, ByteAcc, NoLock, Preserve)
	{
		Offset (0xb0),
				TSL0, 8, /* Battery0 start threshold */
		Offset (0xb1),
				TSH0, 8, /* Battery0 stop threshold */
		Offset (0xb2),
				TSL1, 8, /* Battery1 start threshold */
		Offset (0xb3),
				TSH1, 8, /* Battery1 stop threshold */
	}
}

Scope(\_SB.PCI0.LPCB.EC.BAT0)
{
	/*
	 * Set threshold on battery0
	 *
	 * Arg0: 0: Start threshold
	 *       1: Stop threshold
	 * Arg1: Percentage
	 */
	Method (SETT, 2, NotSerialized)
	{
		if (Arg1 <= 100)
		{
			if (Arg0 == 0)
			{
				TSL0 = Arg1
			}
			if (Arg0 == 1)
			{
				TSH0 = Arg1
			}
		}
	}

	/**
	 * Get threshold on battery0
	 *
	 * Arg0: 0: Start threshold
	 *       1: Stop threshold
	 */
	Method (GETT, 1, NotSerialized)
	{
		if (Arg0 == 0)
		{
			Return (TSL0)
		}
		if (Arg0 == 1)
		{
			Return (TSH0)
		}
		Return (0)
	}
}

Scope(\_SB.PCI0.LPCB.EC.BAT1)
{
	/*
	 * Set threshold on battery1
	 *
	 * Arg0: 0: Start threshold
	 *       1: Stop threshold
	 * Arg1: Percentage
	 */
	Method (SETT, 2, NotSerialized)
	{
		if (Arg1 <= 100)
		{
			if (Arg0 == 0)
			{
				TSL1 = Arg1
			}
			if (Arg0 == 1)
			{
				TSH1 = Arg1
			}
		}
	}

	/**
	 * Get threshold on battery1
	 *
	 * Arg0: 0: Start threshold
	 *       1: Stop threshold
	 */
	Method (GETT, 1, NotSerialized)
	{
		if (Arg0 == 0)
		{
			Return (TSL1)
		}
		if (Arg0 == 1)
		{
			Return (TSH1 & ~0x80)
		}
		Return (0)
	}
}

/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2007-2009 coresystems GmbH
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

/* The _PTS method (Prepare To Sleep) is called before the OS is
 * entering a sleep state. The sleep state number is passed in Arg0
 */

Method(_PTS,1)
{
}

/* The _WAK method is called on system wakeup */

Method(_WAK,1)
{
	/* Not implemented.  */
	Return(Package(){0,0})
}

/* System Bus */

Scope(\_SB)
{
	/* This method is placed on the top level, so we can make sure it's the
	 * first executed _INI method.
	 */
	Method(_INI, 0)
	{
		/* The DTS data in NVS is probably not up to date.
		 * Update temperature values and make sure AP thermal
		 * interrupts can happen
		 */

		/* TRAP(71) */ /* TODO  */

		/* Determine the Operating System and save the value in OSYS.
		 * We have to do this in order to be able to work around
		 * certain windows bugs.
		 *
		 *    OSYS value | Operating System
		 *    -----------+------------------
		 *       2000    | Windows 2000
		 *       2001    | Windows XP(+SP1)
		 *       2002    | Windows XP SP2
		 *       2006    | Windows Vista
		 *       ????    | Windows 7
		 */

		/* Let's assume we're running at least Windows 2000 */
		Store (2000, OSYS)

		If (CondRefOf(_OSI)) {
			If (_OSI("Windows 2001")) {
				Store (2001, OSYS)
			}

			If (_OSI("Windows 2001 SP1")) {
				Store (2001, OSYS)
			}

			If (_OSI("Windows 2001 SP2")) {
				Store (2002, OSYS)
			}

			If (_OSI("Windows 2001.1")) {
				Store (2001, OSYS)
			}

			If (_OSI("Windows 2001.1 SP1")) {
				Store (2001, OSYS)
			}

			If (_OSI("Windows 2006")) {
				Store (2006, OSYS)
			}

			If (_OSI("Windows 2006.1")) {
				Store (2006, OSYS)
			}

			If (_OSI("Windows 2006 SP1")) {
				Store (2006, OSYS)
			}

			If (_OSI("Windows 2009")) {
				Store (2009, OSYS)
			}

			If (_OSI("Windows 2012")) {
				Store (2012, OSYS)
			}
		}
	}
}

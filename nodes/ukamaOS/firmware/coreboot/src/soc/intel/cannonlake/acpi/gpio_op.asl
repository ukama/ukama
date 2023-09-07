/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2018 Intel Corporation.
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

/*
 * Get GPIO Value
 * Arg0 - GPIO Number
 */
Method (GRXS, 1, Serialized)
{
	OperationRegion (PREG, SystemMemory, GADD (Arg0), 4)
	Field (PREG, AnyAcc, NoLock, Preserve)
	{
		VAL0, 32
	}
	And (GPIORXSTATE_MASK, ShiftRight (VAL0, GPIORXSTATE_SHIFT), Local0)

	Return (Local0)
}

/*
 * Get GPIO Tx Value
 * Arg0 - GPIO Number
 */
Method (GTXS, 1, Serialized)
{
	OperationRegion (PREG, SystemMemory, GADD (Arg0), 4)
	Field (PREG, AnyAcc, NoLock, Preserve)
	{
		VAL0, 32
	}
	And (GPIOTXSTATE_MASK, VAL0, Local0)

	Return (Local0)
}

/*
 * Set GPIO Tx Value
 * Arg0 - GPIO Number
 */
Method (STXS, 1, Serialized)
{
	OperationRegion (PREG, SystemMemory, GADD (Arg0), 4)
	Field (PREG, AnyAcc, NoLock, Preserve)
	{
		VAL0, 32
	}
	Or (GPIOTXSTATE_MASK, VAL0, VAL0)
}

/*
 * Clear GPIO Tx Value
 * Arg0 - GPIO Number
 */
Method (CTXS, 1, Serialized)
{
	OperationRegion (PREG, SystemMemory, GADD (Arg0), 4)
	Field (PREG, AnyAcc, NoLock, Preserve)
	{
		VAL0, 32
	}
	And (Not (GPIOTXSTATE_MASK), VAL0, VAL0)
}

/*
 * Set Pad mode
 * Arg0 - GPIO Number
 * Arg1 - Pad mode
 *     0 = GPIO control pad
 *     1 = Native Function 1
 *     2 = Native Function 2
 *     3 = Native Function 3
 */
Method (GPMO, 2, Serialized)
{
	OperationRegion (PREG, SystemMemory, GADD (Arg0), 4)
	Field (PREG, AnyAcc, NoLock, Preserve)
	{
		VAL0, 32
	}
	Store (VAL0, Local0)
	And (Not (GPIOPADMODE_MASK), Local0, Local0)
	And (ShiftLeft (Arg1, GPIOPADMODE_SHIFT, Arg1), GPIOPADMODE_MASK, Arg1)
	Or (Local0, Arg1, VAL0)
}

/*
 * Enable/Disable Tx buffer
 * Arg0 - GPIO Number
 * Arg1 - TxBuffer state
 *     0 = Disable Tx Buffer
 *     1 = Enable Tx Buffer
 */
Method (GTXE, 2, Serialized)
{
	OperationRegion (PREG, SystemMemory, GADD (Arg0), 4)
	Field (PREG, AnyAcc, NoLock, Preserve)
	{
		VAL0, 32
	}

	If (LEqual (Arg1, 1)) {
		And (Not (GPIOTXBUFDIS_MASK), VAL0, VAL0)
	} ElseIf (LEqual (Arg1, 0)){
		Or (GPIOTXBUFDIS_MASK, VAL0, VAL0)
	}
}

/*
 * Enable/Disable Rx buffer
 * Arg0 - GPIO Number
 * Arg1 - RxBuffer state
 *     0 = Disable Rx Buffer
 *     1 = Enable Rx Buffer
 */
Method (GRXE, 2, Serialized)
{
	OperationRegion (PREG, SystemMemory, GADD (Arg0), 4)
	Field (PREG, AnyAcc, NoLock, Preserve)
	{
		VAL0, 32
	}

	If (LEqual (Arg1, 1)) {
		And (Not (GPIORXBUFDIS_MASK), VAL0, VAL0)
	} ElseIf (LEqual (Arg1, 0)){
		Or (GPIORXBUFDIS_MASK, VAL0, VAL0)
	}
}

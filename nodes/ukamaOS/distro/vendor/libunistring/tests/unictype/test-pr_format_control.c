/* DO NOT EDIT! GENERATED AUTOMATICALLY! */
/* Test the Unicode character type functions.
   Copyright (C) 2007-2022 Free Software Foundation, Inc.

   This file is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published
   by the Free Software Foundation, either version 3 of the License,
   or (at your option) any later version.

   This file is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.  */

#include "test-predicate-part1.h"

    { 0x00AD, 0x00AD },
    { 0x180E, 0x180E },
    { 0x200B, 0x200B },
    { 0x2060, 0x2064 },
    { 0x206A, 0x206F },
    { 0x1BCA0, 0x1BCA3 },
    { 0x1D173, 0x1D17A },
    { 0xE0001, 0xE0001 },
    { 0xE0020, 0xE007F }

#define PREDICATE(c) uc_is_property_format_control (c)
#include "test-predicate-part2.h"

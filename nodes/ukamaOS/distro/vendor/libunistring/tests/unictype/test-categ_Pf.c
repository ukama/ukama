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

    { 0x00BB, 0x00BB },
    { 0x2019, 0x2019 },
    { 0x201D, 0x201D },
    { 0x203A, 0x203A },
    { 0x2E03, 0x2E03 },
    { 0x2E05, 0x2E05 },
    { 0x2E0A, 0x2E0A },
    { 0x2E0D, 0x2E0D },
    { 0x2E1D, 0x2E1D },
    { 0x2E21, 0x2E21 }

#define PREDICATE(c) uc_is_general_category (c, UC_CATEGORY_Pf)
#include "test-predicate-part2.h"

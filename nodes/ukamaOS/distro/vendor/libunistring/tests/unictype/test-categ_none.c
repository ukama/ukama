/* Test the Unicode character type functions.
   Copyright (C) 2007-2009 Free Software Foundation, Inc.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.  */

#include <config.h>

#include "unictype.h"

#include <string.h>

#include "macros.h"

int
main ()
{
  /* This test cannot be compiled on platforms on which _UC_CATEGORY_NONE
     is not exported from the libunistring shared library.  For now,
     MSVC is the only platform where this is a problem.  */
#if !defined _MSC_VER

  uc_general_category_t ct = _UC_CATEGORY_NONE;
  unsigned int c;

  for (c = 0; c < 0x110000; c++)
    ASSERT (!uc_is_general_category (c, ct));

#endif

  return 0;
}

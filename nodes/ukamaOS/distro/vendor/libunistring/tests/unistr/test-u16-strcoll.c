/* Test of u16_strcoll() function.
   Copyright (C) 2010-2022 Free Software Foundation, Inc.

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

/* Written by Bruno Haible <bruno@clisp.org>, 2010.  */

#include <config.h>

#include "unistr.h"

#include "macros.h"

#define U_STRCMP u16_strcoll
#include "test-u16-strcmp.h"

int
main ()
{
  /* This test relies on three facts:
     - setlocale is not being called, therefore the locale is the "C" locale.
     - In the "C" locale, strcoll is equivalent to strcmp.
     - In the u16_strcoll implementation, Unicode strings that are not
       convertible to the locale encoding are sorted higher than convertible
       strings and compared according to u16_strcmp.  */

  test_u16_strcmp ();

  return 0;
}

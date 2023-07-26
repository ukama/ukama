/* Test of u16_strmblen() function.
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

int
main ()
{
  int ret;

  /* Test NUL unit input.  */
  {
    static const uint16_t input[] = { 0 };
    ret = u16_strmblen (input);
    ASSERT (ret == 0);
  }

  /* Test ISO 646 unit input.  */
  {
    ucs4_t c;
    uint16_t buf[2];

    for (c = 1; c < 0x80; c++)
      {
        buf[0] = c;
        buf[1] = 0;
        ret = u16_strmblen (buf);
        ASSERT (ret == 1);
      }
  }

  /* Test BMP unit input.  */
  {
    static const uint16_t input[] = { 0x20AC, 0 };
    ret = u16_strmblen (input);
    ASSERT (ret == 1);
  }

  /* Test 2-units character input.  */
  {
    static const uint16_t input[] = { 0xD835, 0xDD1F, 0 };
    ret = u16_strmblen (input);
    ASSERT (ret == 2);
  }

  /* Test incomplete/invalid 1-unit input.  */
  {
    static const uint16_t input[] = { 0xD835, 0 };
    ret = u16_strmblen (input);
    ASSERT (ret == -1);
  }
  {
    static const uint16_t input[] = { 0xDD1F, 0 };
    ret = u16_strmblen (input);
    ASSERT (ret == -1);
  }

  return 0;
}

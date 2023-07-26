/* Test of u32_mbtoucr() function.
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
  ucs4_t uc;
  int ret;

  /* Test NUL unit input.  */
  {
    static const uint32_t input[] = { 0 };
    uc = 0xBADFACE;
    ret = u32_mbtoucr (&uc, input, 1);
    ASSERT (ret == 1);
    ASSERT (uc == 0);
  }

  /* Test ISO 646 unit input.  */
  {
    ucs4_t c;
    uint32_t buf[1];

    for (c = 0; c < 0x80; c++)
      {
        buf[0] = c;
        uc = 0xBADFACE;
        ret = u32_mbtoucr (&uc, buf, 1);
        ASSERT (ret == 1);
        ASSERT (uc == c);
      }
  }

  /* Test BMP unit input.  */
  {
    static const uint32_t input[] = { 0x20AC };
    uc = 0xBADFACE;
    ret = u32_mbtoucr (&uc, input, 1);
    ASSERT (ret == 1);
    ASSERT (uc == 0x20AC);
  }

  /* Test non-BMP unit input.  */
  {
    static const uint32_t input[] = { 0x1D51F };
    uc = 0xBADFACE;
    ret = u32_mbtoucr (&uc, input, 1);
    ASSERT (ret == 1);
    ASSERT (uc == 0x1D51F);
  }

  /* Test incomplete/invalid 1-unit input.  */
  {
    static const uint32_t input[] = { 0x340000 };
    uc = 0xBADFACE;
    ret = u32_mbtoucr (&uc, input, 1);
    ASSERT (ret == -1);
    ASSERT (uc == 0xFFFD);
  }

  return 0;
}

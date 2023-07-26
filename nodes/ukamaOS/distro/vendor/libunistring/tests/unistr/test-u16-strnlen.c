/* Test of u16_strnlen() function.
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

#include "zerosize-ptr.h"
#include "macros.h"

#define UNIT uint16_t
#define U_STRNLEN u16_strnlen
#include "test-strnlen.h"

int
main ()
{
  /* Simple string.  */
  { /* "Grüß Gott. Здравствуйте! x=(-b±sqrt(b²-4ac))/(2a)  日本語,中文,한글" */
    static const uint16_t input[] =
      { 'G', 'r', 0x00FC, 0x00DF, ' ', 'G', 'o', 't', 't', '.', ' ',
        0x0417, 0x0434, 0x0440, 0x0430, 0x0432, 0x0441, 0x0442, 0x0432, 0x0443,
        0x0439, 0x0442, 0x0435, '!', ' ',
        'x', '=', '(', '-', 'b', 0x00B1, 's', 'q', 'r', 't', '(', 'b', 0x00B2,
        '-', '4', 'a', 'c', ')', ')', '/', '(', '2', 'a', ')', ' ', ' ',
        0x65E5, 0x672C, 0x8A9E, ',', 0x4E2D, 0x6587, ',', 0xD55C, 0xAE00, 0
      };
    check (input, SIZEOF (input));
  }

  /* String with characters outside the BMP.  */
  {
    static const uint16_t input[] =
      { '-', '(', 0xD835, 0xDD1E, 0x00D7, 0xD835, 0xDD1F, ')', '=',
        0xD835, 0xDD1F, 0x00D7, 0xD835, 0xDD1E, 0
      };
    check (input, SIZEOF (input));
  }

  return 0;
}

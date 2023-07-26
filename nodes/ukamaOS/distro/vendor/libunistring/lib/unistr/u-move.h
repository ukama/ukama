/* Copy piece of UTF-8/UTF-16/UTF-32 string.
   Copyright (C) 1999, 2002, 2006, 2009-2022 Free Software Foundation, Inc.
   Written by Bruno Haible <bruno@clisp.org>, 2002.

   This file is free software.
   It is dual-licensed under "the GNU LGPLv3+ or the GNU GPLv2+".
   You can redistribute it and/or modify it under either
     - the terms of the GNU Lesser General Public License as published
       by the Free Software Foundation, either version 3, or (at your
       option) any later version, or
     - the terms of the GNU General Public License as published by the
       Free Software Foundation; either version 2, or (at your option)
       any later version, or
     - the same dual license "the GNU LGPLv3+ or the GNU GPLv2+".

   This file is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
   Lesser General Public License and the GNU General Public License
   for more details.

   You should have received a copy of the GNU Lesser General Public
   License and of the GNU General Public License along with this
   program.  If not, see <https://www.gnu.org/licenses/>.  */

#include <string.h>

UNIT *
FUNC (UNIT *dest, const UNIT *src, size_t n)
{
#if 0
  if (dest < src)
    {
      UNIT *destptr = dest;
      const UNIT *srcptr = src;

      for (; n > 0; n--)
        *destptr++ = *srcptr++;
    }
  else if (dest > src)
    {
      UNIT *destptr = dest + n - 1;
      const UNIT *srcptr = src + n - 1;

      for (; n > 0; n--)
        *destptr-- = *srcptr--;
    }
#else
  memmove ((char *) dest, (const char *) src, n * sizeof (UNIT));
#endif
  return dest;
}

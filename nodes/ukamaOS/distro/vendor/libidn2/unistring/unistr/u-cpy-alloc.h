/* Copy piece of UTF-8/UTF-16/UTF-32 string.
   Copyright (C) 1999, 2002, 2006-2007, 2009-2021 Free Software Foundation,
   Inc.
   Written by Bruno Haible <bruno@clisp.org>, 2002.

   This file is free software.
   It is dual-licensed under "the GNU LGPLv3+ or the GNU GPLv2+".
   You can redistribute it and/or modify it under either
     - the terms of the GNU Lesser General Public License as published
       by the Free Software Foundation; either version 3, or (at your
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

#include <stdlib.h>
#include <string.h>

UNIT *
FUNC (const UNIT *s, size_t n)
{
  UNIT *dest;

  dest = (UNIT *) malloc (n > 0 ? n * sizeof (UNIT) : 1);
  if (dest != NULL)
    {
#if 0
      UNIT *destptr = dest;

      for (; n > 0; n--)
        *destptr++ = *s++;
#else
      if (n > 0)
        memcpy ((char *) dest, (const char *) s, n * sizeof (UNIT));
#endif
    }
  return dest;
}

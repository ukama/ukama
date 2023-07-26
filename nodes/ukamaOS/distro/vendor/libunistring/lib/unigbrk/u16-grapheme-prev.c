/* Previous grapheme cluster function.
   Copyright (C) 2010-2022 Free Software Foundation, Inc.
   Written by Ben Pfaff <blp@cs.stanford.edu>, 2010.

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

#include <config.h>

/* Specification.  */
#include "unigbrk.h"

#include "unistr.h"

const uint16_t *
u16_grapheme_prev (const uint16_t *s, const uint16_t *start)
{
  ucs4_t next;

  if (s == start)
    return NULL;

  s = u16_prev (&next, s, start);
  while (s != start)
    {
      const uint16_t *prev_s;
      ucs4_t prev;

      prev_s = u16_prev (&prev, s, start);
      if (prev_s == NULL)
        {
          /* Ill-formed UTF-16 encoding. */
          return start;
        }

      if (uc_is_grapheme_break (prev, next))
        break;

      s = prev_s;
      next = prev;
    }

  return s;
}

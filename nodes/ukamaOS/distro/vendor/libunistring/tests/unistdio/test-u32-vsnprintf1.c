/* Test of u32_vsnprintf() function.
   Copyright (C) 2007-2022 Free Software Foundation, Inc.

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

/* Written by Bruno Haible <bruno@clisp.org>, 2007.  */

#include <config.h>

#include "unistdio.h"

#include <errno.h>
#include <stdarg.h>
#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#include "unistr.h"
#include "xalloc.h"
#include "macros.h"

#include "test-u32-printf1.h"

static uint32_t *
my_xasprintf (const char *format, ...)
{
  va_list args;
  uint32_t buf[1000];
  int retval;
  size_t length;
  uint32_t *result;

  va_start (args, format);
  retval = u32_vsnprintf (buf, sizeof (buf), format, args);
  va_end (args);
  if (retval < 0 || retval >= (int) sizeof (buf))
    return NULL;
  length = u32_strlen (buf);
  result = XNMALLOC (length + 1, uint32_t);
  u32_cpy (result, buf, length + 1);
  return result;
}

static void
test_vsnprintf ()
{
  test_xfunction (my_xasprintf);
}

int
main (int argc, char *argv[])
{
  test_vsnprintf ();

  return 0;
}

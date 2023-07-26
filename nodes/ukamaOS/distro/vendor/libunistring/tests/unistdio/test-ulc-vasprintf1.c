/* Test of ulc_vasprintf() function.
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

#include <stdarg.h>
#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#include "macros.h"

#include "test-ulc-printf1.h"

static char *
my_xasprintf (const char *format, ...)
{
  va_list args;
  char *result;
  int retval;

  va_start (args, format);
  retval = ulc_vasprintf (&result, format, args);
  va_end (args);
  if (retval < 0)
    return NULL;
  ASSERT (result != NULL);
  return result;
}

static void
test_vasprintf ()
{
  test_xfunction (my_xasprintf);
}

int
main (int argc, char *argv[])
{
  test_vasprintf ();
  return 0;
}

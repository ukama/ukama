/* test-locale.c --- Self tests for locale-related (iconv) IDNA processing
   Copyright (C) 2011-2021 Simon Josefsson

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

#include <config.h>

#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <string.h>
#include <stdint.h>

#include <idn2.h>

static int error_count = 0;
static int break_on_error = 1;

_GL_ATTRIBUTE_FORMAT_PRINTF_STANDARD (1, 2)
     static void fail (const char *format, ...)
{
  va_list arg_ptr;

  va_start (arg_ptr, format);
  vfprintf (stderr, format, arg_ptr);
  va_end (arg_ptr);
  error_count++;
  if (break_on_error)
    exit (EXIT_FAILURE);
}

int
main (void)
{
  uint8_t *out;
  int rc;

#if !HAVE_ICONV
  return 77;
#endif

  if ((rc = idn2_lookup_ul ("abc", NULL, 0)) != IDN2_OK)
    {
      fail ("special #5 failed with %d\n", rc);
    }

  /* test libidn compatibility functions */
  if ((rc = idna_to_ascii_lz ("abc", (char **) &out, 0)) != IDN2_OK)
    {
      fail ("special #6 failed with %d\n", rc);
    }
  else
    {
      idn2_free (out);
    }

  if ((rc = idn2_register_ul ("foo", NULL, NULL, 0)) != IDN2_OK)
    fail ("special #6 failed with %d\n", rc);

  return error_count;
}

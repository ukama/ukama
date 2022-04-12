/* test-strerror.c --- Self tests for idn2_strerror* functions
   Copyright (C) 2016 Tim Rühsen

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

#include <idn2.h>

int G_GNUC_IDN2_ATTRIBUTE_CONST
main (void)
{
  int i, failed = 0;

  /* just cover the code paths in idn2_strerror/idn2_strerror_name */
  for (i = -1000; i <= 1000; i++)
    {
      if (!idn2_strerror (i))
	failed++;
      if (!idn2_strerror_name (i))
	failed++;
    }

  return !!failed;
}

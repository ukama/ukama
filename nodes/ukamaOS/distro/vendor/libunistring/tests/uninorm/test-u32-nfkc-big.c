/* Test of Unicode compliance of compatibility normalization of UTF-32 strings.
   Copyright (C) 2009-2022 Free Software Foundation, Inc.

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

/* Written by Bruno Haible <bruno@clisp.org>, 2009.  */

#include <config.h>

#if GNULIB_TEST_UNINORM_U32_NORMALIZE

#include "uninorm.h"

#include <stdlib.h>

#include "unistr.h"
#include "test-u32-normalize-big.h"

static int
check (const uint32_t *c1, size_t c1_length,
       const uint32_t *c2, size_t c2_length,
       const uint32_t *c3, size_t c3_length,
       const uint32_t *c4, size_t c4_length,
       const uint32_t *c5, size_t c5_length)
{
  /* Check c4 == NFKC(c1) == NFKC(c2) == NFKC(c3) == NFKC(c4) == NFKC(c5).  */
  {
    size_t length;
    uint32_t *result;

    result = u32_normalize (UNINORM_NFKC, c1, c1_length, NULL, &length);
    if (!(result != NULL
          && length == c4_length
          && u32_cmp (result, c4, c4_length) == 0))
      return 1;
    free (result);
  }
  {
    size_t length;
    uint32_t *result;

    result = u32_normalize (UNINORM_NFKC, c2, c2_length, NULL, &length);
    if (!(result != NULL
          && length == c4_length
          && u32_cmp (result, c4, c4_length) == 0))
      return 2;
    free (result);
  }
  {
    size_t length;
    uint32_t *result;

    result = u32_normalize (UNINORM_NFKC, c3, c3_length, NULL, &length);
    if (!(result != NULL
          && length == c4_length
          && u32_cmp (result, c4, c4_length) == 0))
      return 3;
    free (result);
  }
  {
    size_t length;
    uint32_t *result;

    result = u32_normalize (UNINORM_NFKC, c4, c4_length, NULL, &length);
    if (!(result != NULL
          && length == c4_length
          && u32_cmp (result, c4, c4_length) == 0))
      return 4;
    free (result);
  }
  {
    size_t length;
    uint32_t *result;

    result = u32_normalize (UNINORM_NFKC, c5, c5_length, NULL, &length);
    if (!(result != NULL
          && length == c4_length
          && u32_cmp (result, c4, c4_length) == 0))
      return 5;
    free (result);
  }
  return 0;
}

int
main (int argc, char *argv[])
{
  struct normalization_test_file file;

  read_normalization_test_file (argv[1], &file);

  test_specific (&file, check);
  test_other (&file, UNINORM_NFKC);

  free_normalization_test_file (&file);

  return 0;
}

#else

#include <stdio.h>

int
main ()
{
  fprintf (stderr, "Skipping test: uninorm/u32-normalize module not included.\n");
  return 77;
}

#endif

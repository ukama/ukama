/* Test of uppercase mapping for UTF-8 strings.
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

#include "unicase.h"

#include <stdlib.h>

#include "unistr.h"
#include "uninorm.h"
#include "macros.h"

static int
check (const uint8_t *input, size_t input_length,
       const char *iso639_language, uninorm_t nf,
       const uint8_t *expected, size_t expected_length)
{
  size_t length;
  uint8_t *result;

  /* Test return conventions with resultbuf == NULL.  */
  result = u8_toupper (input, input_length, iso639_language, nf, NULL, &length);
  if (!(result != NULL))
    return 1;
  if (!(length == expected_length))
    return 2;
  if (!(u8_cmp (result, expected, expected_length) == 0))
    return 3;
  free (result);

  /* Test return conventions with resultbuf too small.  */
  if (expected_length > 0)
    {
      uint8_t *preallocated;

      length = expected_length - 1;
      preallocated = (uint8_t *) malloc (length * sizeof (uint8_t));
      result = u8_toupper (input, input_length, iso639_language, nf, preallocated, &length);
      if (!(result != NULL))
        return 4;
      if (!(result != preallocated))
        return 5;
      if (!(length == expected_length))
        return 6;
      if (!(u8_cmp (result, expected, expected_length) == 0))
        return 7;
      free (result);
      free (preallocated);
    }

  /* Test return conventions with resultbuf large enough.  */
  {
    uint8_t *preallocated;

    length = expected_length;
    preallocated = (uint8_t *) malloc (length * sizeof (uint8_t));
    result = u8_toupper (input, input_length, iso639_language, nf, preallocated, &length);
    if (!(result != NULL))
      return 8;
    if (!(preallocated == NULL || result == preallocated))
      return 9;
    if (!(length == expected_length))
      return 10;
    if (!(u8_cmp (result, expected, expected_length) == 0))
      return 11;
    free (preallocated);
  }

  return 0;
}

int
main ()
{
  { /* Empty string.  */
    ASSERT (check (NULL, 0, NULL, NULL, NULL, 0) == 0);
    ASSERT (check (NULL, 0, NULL, UNINORM_NFC, NULL, 0) == 0);
  }

  /* Simple string.  */
  { /* "Grüß Gott. Здравствуйте! x=(-b±sqrt(b²-4ac))/(2a)  日本語,中文,한글" */
    static const uint8_t input[] =
      { 'G', 'r', 0xC3, 0xBC, 0xC3, 0x9F, ' ', 'G', 'o', 't', 't', '.', ' ',
        0xD0, 0x97, 0xD0, 0xB4, 0xD1, 0x80, 0xD0, 0xB0, 0xD0, 0xB2, 0xD1, 0x81,
        0xD1, 0x82, 0xD0, 0xB2, 0xD1, 0x83, 0xD0, 0xB9, 0xD1, 0x82, 0xD0, 0xB5,
        '!', ' ', 'x', '=', '(', '-', 'b', 0xC2, 0xB1, 's', 'q', 'r', 't', '(',
        'b', 0xC2, 0xB2, '-', '4', 'a', 'c', ')', ')', '/', '(', '2', 'a', ')',
        ' ', ' ', 0xE6, 0x97, 0xA5, 0xE6, 0x9C, 0xAC, 0xE8, 0xAA, 0x9E, ',',
        0xE4, 0xB8, 0xAD, 0xE6, 0x96, 0x87, ',',
        0xED, 0x95, 0x9C, 0xEA, 0xB8, 0x80, '\n'
      };
    static const uint8_t casemapped[] =
      { 'G', 'R', 0xC3, 0x9C, 0x53, 0x53, ' ', 'G', 'O', 'T', 'T', '.', ' ',
        0xD0, 0x97, 0xD0, 0x94, 0xD0, 0xA0, 0xD0, 0x90, 0xD0, 0x92, 0xD0, 0xA1,
        0xD0, 0xA2, 0xD0, 0x92, 0xD0, 0xA3, 0xD0, 0x99, 0xD0, 0xA2, 0xD0, 0x95,
        '!', ' ', 'X', '=', '(', '-', 'B', 0xC2, 0xB1, 'S', 'Q', 'R', 'T', '(',
        'B', 0xC2, 0xB2, '-', '4', 'A', 'C', ')', ')', '/', '(', '2', 'A', ')',
        ' ', ' ', 0xE6, 0x97, 0xA5, 0xE6, 0x9C, 0xAC, 0xE8, 0xAA, 0x9E, ',',
        0xE4, 0xB8, 0xAD, 0xE6, 0x96, 0x87, ',',
        0xED, 0x95, 0x9C, 0xEA, 0xB8, 0x80, '\n'
      };
    ASSERT (check (input, SIZEOF (input), NULL, NULL, casemapped, SIZEOF (casemapped)) == 0);
  }

  /* Case mapping can increase the number of Unicode characters.  */
  { /* LATIN SMALL LETTER N PRECEDED BY APOSTROPHE */
    static const uint8_t input[]      = { 0xC5, 0x89 };
    static const uint8_t casemapped[] = { 0xCA, 0xBC, 0x4E };
    ASSERT (check (input, SIZEOF (input), NULL, NULL, casemapped, SIZEOF (casemapped)) == 0);
  }
  { /* GREEK SMALL LETTER IOTA WITH DIALYTIKA AND TONOS */
    static const uint8_t input[]      = { 0xCE, 0x90 };
    static const uint8_t casemapped[] = { 0xCE, 0x99, 0xCC, 0x88, 0xCC, 0x81 };
    ASSERT (check (input, SIZEOF (input), NULL, NULL, casemapped, SIZEOF (casemapped)) == 0);
  }

  /* Turkish letters i İ ı I */
  { /* LATIN CAPITAL LETTER I */
    static const uint8_t input[]      = { 0x49 };
    static const uint8_t casemapped[] = { 0x49 };
    ASSERT (check (input, SIZEOF (input), NULL, NULL, casemapped, SIZEOF (casemapped)) == 0);
    ASSERT (check (input, SIZEOF (input), "tr", NULL, casemapped, SIZEOF (casemapped)) == 0);
  }
  { /* LATIN SMALL LETTER I */
    static const uint8_t input[]         = { 0x69 };
    static const uint8_t casemapped[]    = { 0x49 };
    static const uint8_t casemapped_tr[] = { 0xC4, 0xB0 };
    ASSERT (check (input, SIZEOF (input), NULL, NULL, casemapped, SIZEOF (casemapped)) == 0);
    ASSERT (check (input, SIZEOF (input), "tr", NULL, casemapped_tr, SIZEOF (casemapped_tr)) == 0);
  }
  { /* LATIN CAPITAL LETTER I WITH DOT ABOVE */
    static const uint8_t input[]      = { 0xC4, 0xB0 };
    static const uint8_t casemapped[] = { 0xC4, 0xB0 };
    ASSERT (check (input, SIZEOF (input), NULL, NULL, casemapped, SIZEOF (casemapped)) == 0);
    ASSERT (check (input, SIZEOF (input), "tr", NULL, casemapped, SIZEOF (casemapped)) == 0);
  }
  { /* LATIN SMALL LETTER DOTLESS I */
    static const uint8_t input[]      = { 0xC4, 0xB1 };
    static const uint8_t casemapped[] = { 0x49 };
    ASSERT (check (input, SIZEOF (input), NULL, NULL, casemapped, SIZEOF (casemapped)) == 0);
    ASSERT (check (input, SIZEOF (input), "tr", NULL, casemapped, SIZEOF (casemapped)) == 0);
  }
  { /* "topkapı" */
    static const uint8_t input[] =
      { 0x74, 0x6F, 0x70, 0x6B, 0x61, 0x70, 0xC4, 0xB1 };
    static const uint8_t casemapped[] =
      { 0x54, 0x4F, 0x50, 0x4B, 0x41, 0x50, 0x49 };
    ASSERT (check (input, SIZEOF (input), NULL, NULL, casemapped, SIZEOF (casemapped)) == 0);
    ASSERT (check (input, SIZEOF (input), "tr", NULL, casemapped, SIZEOF (casemapped)) == 0);
  }

  /* Uppercasing can increase the number of Unicode characters.  */
  { /* "heiß" */
    static const uint8_t input[]      = { 0x68, 0x65, 0x69, 0xC3, 0x9F };
    static const uint8_t casemapped[] = { 0x48, 0x45, 0x49, 0x53, 0x53 };
    ASSERT (check (input, SIZEOF (input), NULL, NULL, casemapped, SIZEOF (casemapped)) == 0);
  }

  /* Case mappings for some characters can depend on the surrounding characters.  */
  { /* "περισσότερες πληροφορίες" */
    static const uint8_t input[] =
      {
        0xCF, 0x80, 0xCE, 0xB5, 0xCF, 0x81, 0xCE, 0xB9, 0xCF, 0x83, 0xCF, 0x83,
        0xCF, 0x8C, 0xCF, 0x84, 0xCE, 0xB5, 0xCF, 0x81, 0xCE, 0xB5, 0xCF, 0x82,
        ' ', 0xCF, 0x80, 0xCE, 0xBB, 0xCE, 0xB7, 0xCF, 0x81, 0xCE, 0xBF,
        0xCF, 0x86, 0xCE, 0xBF, 0xCF, 0x81, 0xCE, 0xAF, 0xCE, 0xB5, 0xCF, 0x82
      };
    static const uint8_t casemapped[] =
      {
        0xCE, 0xA0, 0xCE, 0x95, 0xCE, 0xA1, 0xCE, 0x99, 0xCE, 0xA3, 0xCE, 0xA3,
        0xCE, 0x8C, 0xCE, 0xA4, 0xCE, 0x95, 0xCE, 0xA1, 0xCE, 0x95, 0xCE, 0xA3,
        ' ', 0xCE, 0xA0, 0xCE, 0x9B, 0xCE, 0x97, 0xCE, 0xA1, 0xCE, 0x9F,
        0xCE, 0xA6, 0xCE, 0x9F, 0xCE, 0xA1, 0xCE, 0x8A, 0xCE, 0x95, 0xCE, 0xA3
      };
    ASSERT (check (input, SIZEOF (input), NULL, NULL, casemapped, SIZEOF (casemapped)) == 0);
  }

  /* Case mapping can require subsequent normalization.  */
  { /* LATIN SMALL LETTER J WITH CARON, COMBINING DOT BELOW */
    static const uint8_t input[]                 = { 0xC7, 0xB0, 0xCC, 0xA3 };
    static const uint8_t casemapped[]            = { 0x4A, 0xCC, 0x8C, 0xCC, 0xA3 };
    static const uint8_t casemapped_normalized[] = { 0x4A, 0xCC, 0xA3, 0xCC, 0x8C };
    ASSERT (check (input, SIZEOF (input), NULL, NULL,        casemapped, SIZEOF (casemapped)) == 0);
    ASSERT (check (input, SIZEOF (input), NULL, UNINORM_NFC, casemapped_normalized, SIZEOF (casemapped_normalized)) == 0);
  }

  return 0;
}

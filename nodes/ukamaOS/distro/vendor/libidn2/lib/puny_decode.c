/* punycode.c - punycode encoding/decoding
   Copyright (C) 2011-2021 Simon Josefsson

   Libidn2 is free software: you can redistribute it and/or modify it
   under the terms of either:

     * the GNU Lesser General Public License as published by the Free
       Software Foundation; either version 3 of the License, or (at
       your option) any later version.

   or

     * the GNU General Public License as published by the Free
       Software Foundation; either version 2 of the License, or (at
       your option) any later version.

   or both in parallel, as here.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received copies of the GNU General Public License and
   the GNU Lesser General Public License along with this program.  If
   not, see <http://www.gnu.org/licenses/>.
*/

/*
  Code copied from http://www.nicemice.net/idn/punycode-spec.gz on
  2011-01-04 with SHA-1 a966a8017f6be579d74a50a226accc7607c40133
  labeled punycode-spec 1.0.3 (2006-Mar-23-Thu).  It is modified for
  Libidn2 by Simon Josefsson.  License on the original code:

  punycode-spec 1.0.3 (2006-Mar-23-Thu)
  http://www.nicemice.net/idn/
  Adam M. Costello
  http://www.nicemice.net/amc/

  B. Disclaimer and license

    Regarding this entire document or any portion of it (including
    the pseudocode and C code), the author makes no guarantees and
    is not responsible for any damage resulting from its use.  The
    author grants irrevocable permission to anyone to use, modify,
    and distribute it in any way that does not diminish the rights
    of anyone else to use, modify, and distribute it, provided that
    redistributed derivative works do not contain misleading author or
    version information.  Derivative works need not be licensed under
    similar terms.

  C. Punycode sample implementation

  punycode-sample.c 2.0.0 (2004-Mar-21-Sun)
  http://www.nicemice.net/idn/
  Adam M. Costello
  http://www.nicemice.net/amc/

  This is ANSI C code (C89) implementing Punycode 1.0.x.
*/

#include <config.h>

#include "idn2.h"		/* IDN2_OK, ... */

/* Re-definitions to avoid modifying code below too much. */
#define punycode_uint uint32_t
#define punycode_success IDN2_OK
#define punycode_overflow IDN2_PUNYCODE_OVERFLOW
#define punycode_big_output IDN2_PUNYCODE_BIG_OUTPUT
#define punycode_bad_input IDN2_PUNYCODE_BAD_INPUT
#define punycode_decode _idn2_punycode_decode_internal

/**********************************************************/
/* Implementation (would normally go in its own .c file): */

#include <stdio.h>
#include <string.h>

#include "punycode.h"

/*** Bootstring parameters for Punycode ***/

enum
{ base = 36, tmin = 1, tmax = 26, skew = 38, damp = 700,
  initial_bias = 72, initial_n = 0x80, delimiter = 0x2D
};

/* basic(cp) tests whether cp is a basic code point: */
#define basic(cp) ((cp >= 'a' && cp <= 'z') || (cp >= '0' && cp <='9') || (cp >= 'A' && cp <='Z') || cp == '-' || cp == '_')

/* decode_digit(cp) returns the numeric value of a basic code */
/* point (for use in representing integers) in the range 0 to */
/* base-1, or base if cp does not represent a value.          */

static unsigned
decode_digit (int cp)
{
  if (cp >= 'a' && cp <= 'z')
    return cp - 'a';
  if (cp >= '0' && cp <= '9')
    return cp - '0' + 26;
  if (cp >= 'A' && cp <= 'Z')
    return cp - 'A';

  return 0;
}

/*** Platform-specific constants ***/

/* maxint is the maximum value of a punycode_uint variable: */
static const punycode_uint maxint = -1;
/* Because maxint is unsigned, -1 becomes the maximum value. */

/*** Bias adaptation function ***/

static punycode_uint
adapt (punycode_uint delta, punycode_uint numpoints, int firsttime)
  _GL_ATTRIBUTE_CONST;

     static punycode_uint adapt (punycode_uint delta, punycode_uint numpoints,
				 int firsttime)
{
  punycode_uint k;

  delta = firsttime ? delta / damp : delta >> 1;
  /* delta >> 1 is a faster way of doing delta / 2 */
  delta += delta / numpoints;

  for (k = 0; delta > ((base - tmin) * tmax) / 2; k += base)
    {
      delta /= base - tmin;
    }

  return k + (base - tmin + 1) * delta / (delta + skew);
}

/*** Main decode function ***/

int
punycode_decode (size_t input_length,
		 const char input[],
		 size_t *output_length, punycode_uint output[])
{
  punycode_uint n, out = 0, i, max_out, bias, oldi, w, k, digit, t;
  size_t b = 0, j, in;

  if (!input_length)
    return punycode_bad_input;

  /* Check that all chars are basic */
  for (j = 0; j < input_length; ++j)
    {
      if (!basic (input[j]))
	return punycode_bad_input;
      if (input[j] == delimiter)
	b = j;
    }

  max_out =
    *output_length > maxint ? maxint : (punycode_uint) * output_length;

  if (input[b] == delimiter)
    {
      /* do not accept leading or trailing delimiter
       *   - leading delim must be omitted if there is no ASCII char in u-label
       *   - trailing delim means there where no non-ASCII chars in u-label
       */
      if (!b || b == input_length - 1)
	return punycode_bad_input;

      if (b >= max_out)
	return punycode_big_output;

      /* Check that all chars before last delimiter are basic chars */
      /* and copy the first b code points to the output. */
      for (j = 0; j < b; j++)
	output[out++] = input[j];

      b += 1;			/* advance to non-basic char encoding */
    }

  /* Initialize the state: */
  n = initial_n;
  i = 0;
  bias = initial_bias;

  /* Main decoding loop:  Start just after the last delimiter if any  */
  /* basic code points were copied; start at the beginning otherwise. */
  for (in = b; in < input_length; ++out)
    {

      /* in is the index of the next ASCII code point to be consumed, */
      /* and out is the number of code points in the output array.    */

      /* Decode a generalized variable-length integer into delta,  */
      /* which gets added to i.  The overflow checking is easier   */
      /* if we increase i as we go, then subtract off its starting */
      /* value at the end to obtain delta.                         */
      for (oldi = i, w = 1, k = base;; k += base)
	{
	  if (in >= input_length)
	    return punycode_bad_input;
	  digit = decode_digit (input[in++]);
	  if (digit >= base)
	    return punycode_bad_input;
	  if (digit > (maxint - i) / w)
	    return punycode_overflow;
	  i += digit * w;
	  t = k <= bias /* + tmin */ ? tmin :	/* +tmin not needed */
	    k >= bias + tmax ? tmax : k - bias;
	  if (digit < t)
	    break;
	  if (w > maxint / (base - t))
	    return punycode_overflow;
	  w *= (base - t);
	}

      bias = adapt (i - oldi, out + 1, oldi == 0);

      /* i was supposed to wrap around from out+1 to 0,   */
      /* incrementing n each time, so we'll fix that now: */
      if (i / (out + 1) > maxint - n)
	return punycode_overflow;
      n += i / (out + 1);
      if (n > 0x10FFFF || (n >= 0xD800 && n <= 0xDBFF))
	return punycode_bad_input;
      i %= (out + 1);

      /* Insert n at position i of the output: */

      /* not needed for Punycode: */
      /* if (basic(n)) return punycode_bad_input; */
      if (out >= max_out)
	return punycode_big_output;

      memmove (output + i + 1, output + i, (out - i) * sizeof *output);
      output[i++] = n;
    }

  *output_length = (size_t) out;
  /* cannot overflow because out <= old value of *output_length */
  return punycode_success;
}

/* Create a compatibility symbol if supported.  Hidden references make
   the target symbol hidden, hence the alias.  */
#ifdef HAVE_SYMVER_ALIAS_SUPPORT
__typeof__ (_idn2_punycode_decode_internal) _idn2_punycode_decode
  __attribute__((visibility ("default"),
		 alias ("_idn2_punycode_decode_internal")));
__asm__ (".symver _idn2_punycode_decode, _idn2_punycode_decode@IDN2_0.0.0");
#endif

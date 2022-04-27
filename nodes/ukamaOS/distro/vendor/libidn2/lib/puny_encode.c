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
#define punycode_encode _idn2_punycode_encode_internal

/**********************************************************/
/* Implementation (would normally go in its own .c file): */

#include <string.h>

#include "punycode.h"

/*** Bootstring parameters for Punycode ***/

enum
{ base = 36, tmin = 1, tmax = 26, skew = 38, damp = 700,
  initial_bias = 72, initial_n = 0x80, delimiter = 0x2D
};

/* basic(cp) tests whether cp is a basic code point: */
#define basic(cp) ((punycode_uint)(cp) < 0x80)

/* encode_digit(d,flag) returns the basic code point whose value      */
/* (when used for representing integers) is d, which needs to be in   */
/* the range 0 to base-1.  The lowercase form is used unless flag is  */
/* nonzero, in which case the uppercase form is used.  The behavior   */
/* is undefined if flag is nonzero and digit d has no uppercase form. */

static char
encode_digit (punycode_uint d, int flag)
{
  return d + 22 + 75 * (d < 26) - ((flag != 0) << 5);
  /*  0..25 map to ASCII a..z or A..Z */
  /* 26..35 map to ASCII 0..9         */
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

/*** Main encode function ***/

int
punycode_encode (size_t input_length_orig,
		 const punycode_uint input[],
		 size_t *output_length, char output[])
{
  punycode_uint input_length, n, delta, h, b, bias, j, m, q, k, t;
  size_t out, max_out;

  /* The Punycode spec assumes that the input length is the same type */
  /* of integer as a code point, so we need to convert the size_t to  */
  /* a punycode_uint, which could overflow.                           */

  if (input_length_orig > maxint)
    return punycode_overflow;
  input_length = (punycode_uint) input_length_orig;

  /* Initialize the state: */

  n = initial_n;
  delta = 0;
  out = 0;
  max_out = *output_length;
  bias = initial_bias;

  /* Handle the basic code points: */

  for (j = 0; j < input_length; ++j)
    {
      if (basic (input[j]))
	{
	  if (max_out - out < 2)
	    return punycode_big_output;
	  output[out++] = (char) input[j];
	}
      else if (input[j] > 0x10FFFF
	       || (input[j] >= 0xD800 && input[j] <= 0xDBFF))
	return punycode_bad_input;
    }

  h = b = (punycode_uint) out;
  /* cannot overflow because out <= input_length <= maxint */

  /* h is the number of code points that have been handled, b is the  */
  /* number of basic code points, and out is the number of ASCII code */
  /* points that have been output.                                    */

  if (b > 0)
    output[out++] = delimiter;

  /* Main encoding loop: */

  while (h < input_length)
    {
      /* All non-basic code points < n have been     */
      /* handled already.  Find the next larger one: */

      for (m = maxint, j = 0; j < input_length; ++j)
	{
	  /* if (basic(input[j])) continue; */
	  /* (not needed for Punycode) */
	  if (input[j] >= n && input[j] < m)
	    m = input[j];
	}

      /* Increase delta enough to advance the decoder's    */
      /* <n,i> state to <m,0>, but guard against overflow: */

      if (m - n > (maxint - delta) / (h + 1))
	return punycode_overflow;
      delta += (m - n) * (h + 1);
      n = m;

      for (j = 0; j < input_length; ++j)
	{
	  /* Punycode does not need to check whether input[j] is basic: */
	  if (input[j] < n /* || basic(input[j]) */ )
	    {
	      if (++delta == 0)
		return punycode_overflow;
	    }

	  if (input[j] == n)
	    {
	      /* Represent delta as a generalized variable-length integer: */

	      for (q = delta, k = base;; k += base)
		{
		  if (out >= max_out)
		    return punycode_big_output;
		  t = k <= bias /* + tmin */ ? tmin :	/* +tmin not needed */
		    k >= bias + tmax ? tmax : k - bias;
		  if (q < t)
		    break;
		  output[out++] = encode_digit (t + (q - t) % (base - t), 0);
		  q = (q - t) / (base - t);
		}

	      output[out++] = encode_digit (q, 0);
	      bias = adapt (delta, h + 1, h == b);
	      delta = 0;
	      ++h;
	    }
	}

      ++delta, ++n;
    }

  *output_length = out;
  return punycode_success;
}

/* Create a compatibility symbol if supported.  Hidden references make
   the target symbol hidden, hence the alias.  */
#ifdef HAVE_SYMVER_ALIAS_SUPPORT
__typeof__ (_idn2_punycode_encode_internal) _idn2_punycode_encode
  __attribute__((visibility ("default"),
		 alias ("_idn2_punycode_encode_internal")));
__asm__ (".symver _idn2_punycode_encode, _idn2_punycode_encode@IDN2_0.0.0");
#endif

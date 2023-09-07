/* test-tounicode.c --- Self tests for Libidn2 _to_unicode_ functions.
   Copyright (C) 2017 Tim Rühsen

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

/* Based on GNU Libidn tst_punycode.c */

#include <config.h>

#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <string.h>
#include <locale.h>

#include <uniconv.h>		/* u8_strconv_from_locale */
#include <unistr.h>

#include <idn2.h>

typedef struct
{
  const char *name;
  const char *punycode;
  const uint32_t u32_expected[48];
  int rc_expected;
} test_t;

const test_t test[] = {
  {
   "(A) Arabic (Egyptian)",
   "xn--egbpdaj6bu4bxfgehfvwxn",
   {
    0x0644, 0x064A, 0x0647, 0x0645, 0x0627, 0x0628, 0x062A, 0x0643,
    0x0644, 0x0645, 0x0648, 0x0634, 0x0639, 0x0631, 0x0628, 0x064A,
    0x061F, 0,
    },
   IDN2_OK},
  {
   "(B) Chinese (simplified)",
   "xn--ihqwcrb4cv8a8dqg056pqjye",
   {
    0x4ED6, 0x4EEC, 0x4E3A, 0x4EC0, 0x4E48, 0x4E0D, 0x8BF4, 0x4E2D,
    0x6587, 0},
   IDN2_OK},
  {
   "(C) Chinese (traditional)",
   "xn--ihqwctvzc91f659drss3x8bo0yb",
   {
    0x4ED6, 0x5011, 0x7232, 0x4EC0, 0x9EBD, 0x4E0D, 0x8AAA, 0x4E2D,
    0x6587, 0},
   IDN2_OK},
  {
   "(D) Czech: Pro<ccaron>prost<ecaron>nemluv<iacute><ccaron>esky",
   "xn--Proprostnemluvesky-uyb24dma41a",
   {
    0x0050, 0x0072, 0x006F, 0x010D, 0x0070, 0x0072, 0x006F, 0x0073,
    0x0074, 0x011B, 0x006E, 0x0065, 0x006D, 0x006C, 0x0075, 0x0076,
    0x00ED, 0x010D, 0x0065, 0x0073, 0x006B, 0x0079, 0},
   IDN2_OK},
  {
   "(E) Hebrew:",
   "xn--4dbcagdahymbxekheh6e0a7fei0b",
   {
    0x05DC, 0x05DE, 0x05D4, 0x05D4, 0x05DD, 0x05E4, 0x05E9, 0x05D5,
    0x05D8, 0x05DC, 0x05D0, 0x05DE, 0x05D3, 0x05D1, 0x05E8, 0x05D9,
    0x05DD, 0x05E2, 0x05D1, 0x05E8, 0x05D9, 0x05EA, 0},
   IDN2_OK},
  {
   "(F) Hindi (Devanagari):",
   "xn--i1baa7eci9glrd9b2ae1bj0hfcgg6iyaf8o0a1dig0cd",
   {
    0x092F, 0x0939, 0x0932, 0x094B, 0x0917, 0x0939, 0x093F, 0x0928,
    0x094D, 0x0926, 0x0940, 0x0915, 0x094D, 0x092F, 0x094B, 0x0902,
    0x0928, 0x0939, 0x0940, 0x0902, 0x092C, 0x094B, 0x0932, 0x0938,
    0x0915, 0x0924, 0x0947, 0x0939, 0x0948, 0x0902, 0},
   IDN2_OK},
  {
   "(G) Japanese (kanji and hiragana):",
   "xn--n8jok5ay5dzabd5bym9f0cm5685rrjetr6pdxa",
   {
    0x306A, 0x305C, 0x307F, 0x3093, 0x306A, 0x65E5, 0x672C, 0x8A9E,
    0x3092, 0x8A71, 0x3057, 0x3066, 0x304F, 0x308C, 0x306A, 0x3044,
    0x306E, 0x304B, 0},
   IDN2_OK},
  {
   "(H) Korean (Hangul syllables):",
   "xn--989aomsvi5e83db1d2a355cv1e0vak1dwrv93d5xbh15a0dt30a5jpsd879ccm6fea98c",
   {
    0xC138, 0xACC4, 0xC758, 0xBAA8, 0xB4E0, 0xC0AC, 0xB78C, 0xB4E4,
    0xC774, 0xD55C, 0xAD6D, 0xC5B4, 0xB97C, 0xC774, 0xD574, 0xD55C,
    0xB2E4, 0xBA74, 0xC5BC, 0xB9C8, 0xB098, 0xC88B, 0xC744, 0xAE4C, 0},
   IDN2_OK},
  {
   "(I) Russian (Cyrillic):",
   "xn--b1abfaaepdrnnbgefbadotcwatmq2g4l",
   {
    0x043F, 0x043E, 0x0447, 0x0435, 0x043C, 0x0443, 0x0436, 0x0435,
    0x043E, 0x043D, 0x0438, 0x043D, 0x0435, 0x0433, 0x043E, 0x0432,
    0x043E, 0x0440, 0x044F, 0x0442, 0x043F, 0x043E, 0x0440, 0x0443,
    0x0441, 0x0441, 0x043A, 0x0438, 0},
   IDN2_OK},
  {
   "(J) Spanish: Porqu<eacute>nopuedensimplementehablarenEspa<ntilde>ol",
   "xn--PorqunopuedensimplementehablarenEspaol-fmd56a",
   {
    0x0050, 0x006F, 0x0072, 0x0071, 0x0075, 0x00E9, 0x006E, 0x006F,
    0x0070, 0x0075, 0x0065, 0x0064, 0x0065, 0x006E, 0x0073, 0x0069,
    0x006D, 0x0070, 0x006C, 0x0065, 0x006D, 0x0065, 0x006E, 0x0074,
    0x0065, 0x0068, 0x0061, 0x0062, 0x006C, 0x0061, 0x0072, 0x0065,
    0x006E, 0x0045, 0x0073, 0x0070, 0x0061, 0x00F1, 0x006F, 0x006C, 0},
   IDN2_OK},
  {
   "(K) Vietnamese:",
   "xn--TisaohkhngthchnitingVit-kjcr8268qyxafd2f1b9g",
   {
    0x0054, 0x1EA1, 0x0069, 0x0073, 0x0061, 0x006F, 0x0068, 0x1ECD,
    0x006B, 0x0068, 0x00F4, 0x006E, 0x0067, 0x0074, 0x0068, 0x1EC3,
    0x0063, 0x0068, 0x1EC9, 0x006E, 0x00F3, 0x0069, 0x0074, 0x0069,
    0x1EBF, 0x006E, 0x0067, 0x0056, 0x0069, 0x1EC7, 0x0074, 0},
   IDN2_OK},
  {
   "(L) 3<nen>B<gumi><kinpachi><sensei>",
   "xn--3B-ww4c5e180e575a65lsy2b",
   {
    0x0033, 0x5E74, 0x0042, 0x7D44, 0x91D1, 0x516B, 0x5148, 0x751F, 0},
   IDN2_OK},
  {
   "(M) <amuro><namie>-with-SUPER-MONKEYS",
   "xn---with-SUPER-MONKEYS-pc58ag80a8qai00g7n9n",
   {
    0x5B89, 0x5BA4, 0x5948, 0x7F8E, 0x6075, 0x002D, 0x0077, 0x0069,
    0x0074, 0x0068, 0x002D, 0x0053, 0x0055, 0x0050, 0x0045, 0x0052,
    0x002D, 0x004D, 0x004F, 0x004E, 0x004B, 0x0045, 0x0059, 0x0053, 0},
   IDN2_OK},
  {
   "(N) Hello-Another-Way-<sorezore><no><basho>",
   "xn--Hello-Another-Way--fc4qua05auwb3674vfr0b",
   {
    0x0048, 0x0065, 0x006C, 0x006C, 0x006F, 0x002D, 0x0041, 0x006E,
    0x006F, 0x0074, 0x0068, 0x0065, 0x0072, 0x002D, 0x0057, 0x0061,
    0x0079, 0x002D, 0x305D, 0x308C, 0x305E, 0x308C, 0x306E, 0x5834,
    0x6240, 0},
   IDN2_OK},
  {
   "(O) <hitotsu><yane><no><shita>2",
   "xn--2-u9tlzr9756bt3uc0v",
   {
    0x3072, 0x3068, 0x3064, 0x5C4B, 0x6839, 0x306E, 0x4E0B, 0x0032, 0},
   IDN2_OK},
  {
   "(P) Maji<de>Koi<suru>5<byou><mae>",
   "xn--MajiKoi5-783gue6qz075azm5e",
   {
    0x004D, 0x0061, 0x006A, 0x0069, 0x3067, 0x004B, 0x006F, 0x0069,
    0x3059, 0x308B, 0x0035, 0x79D2, 0x524D, 0},
   IDN2_OK},
  {
   "(Q) <pafii>de<runba>",
   "xn--de-jg4avhby1noc0d",
   {
    0x30D1, 0x30D5, 0x30A3, 0x30FC, 0x0064, 0x0065, 0x30EB, 0x30F3, 0x30D0},
   IDN2_OK},
  {
   "(R) <sono><supiido><de>",
   "xn--d9juau41awczczp",
   {
    0x305D, 0x306E, 0x30B9, 0x30D4, 0x30FC, 0x30C9, 0x3067, 0},
   IDN2_OK},
  {
   "(S) -> $1",
   "xn---> $1-",
   {
    0x002D, 0x003E, 0x0020, 0x0024, 0x0031, 0},
   IDN2_PUNYCODE_BAD_INPUT},
  {
   "(T) -> $1",
   "xn---> $1",
   {
    0x002D, 0x003E, 0x0020, 0x0024, 0x0031, 0},
   IDN2_PUNYCODE_BAD_INPUT},
  {
   "(U) -> $1.00 <-",
   "xn---> $1.00 <--",
   {
    0x002D, 0x003E, 0x0020, 0x0024, 0x0031, 0x002E, 0x0030, 0x0030,
    0x0020, 0x003C, 0x002D, 0},
   IDN2_PUNYCODE_BAD_INPUT},
  {
   "(V) <sono><supiido><de>",
   "xn--d9juau41awczczp.xn--bel-goa.DE",
   {
    0x305D, 0x306E, 0x30B9, 0x30D4, 0x30FC, 0x30C9, 0x3067, 0x002E,
    0x00FC, 0x0062, 0x0065, 0x006c, 0x002e, 0x0044, 0x0045, 0},
   IDN2_OK},
  {
   "(W) <sono><supiido><de>",
   "XN--d9juau41awczczp",
   {
    0x305D, 0x306E, 0x30B9, 0x30D4, 0x30FC, 0x30C9, 0x3067, 0},
   IDN2_OK},
  {
   "(X) utf-8 crash ?",
   "\200bad.com",
   {
    0},
   IDN2_ENCODING_ERROR		/* or IDN2_ICONV_FAIL with idn2_to_unicode_lzlz() due to bad UTF-8 input */
   },

  /* Test vectors from https://bugs.debian.org/610617 */
  {
   "Debian test #1",
   "XN----7SBAABF4DLDYSIEHP4NTB.XN--P1AI",
   {
    0x0441, 0x0430, 0x043c, 0x0430, 0x0440, 0x0441, 0x043a, 0x0430,
    0x044f, 0x002d, 0x043e, 0x0431, 0x043b, 0x0430, 0x0441, 0x0442,
    0x044c, 0x002e, 0x0440, 0x0444, 0},
   IDN2_OK},
  {
   "Debian test #2",
   "xn----7SBAABF4DLDYSIEHP4NTB.XN--P1AI",
   {
    0x0441, 0x0430, 0x043c, 0x0430, 0x0440, 0x0441, 0x043a, 0x0430,
    0x044f, 0x002d, 0x043e, 0x0431, 0x043b, 0x0430, 0x0441, 0x0442,
    0x044c, 0x002e, 0x0440, 0x0444, 0},
   IDN2_OK},
  {
   "Debian test #3",
   "xn----7SBAABF4DLDYSIEHP4NTB.xn--P1AI",
   {
    0x0441, 0x0430, 0x043c, 0x0430, 0x0440, 0x0441, 0x043a, 0x0430,
    0x044f, 0x002d, 0x043e, 0x0431, 0x043b, 0x0430, 0x0441, 0x0442,
    0x044c, 0x002e, 0x0440, 0x0444, 0},
   IDN2_OK},

  /* Test vectors copied from gnutls */
  {
   "GnuTLS test #1",
   "xn--nxasmm1c.com",
   {
    0x03b2, 0x03cc, 0x03bb, 0x03bf, 0x03c2, 0x002e, 0x0063, 0x006f,
    0x006d, 0},
   IDN2_OK},
  {
   "GnuTLS test #2",
   "xn--nxasmq6b.com",
   {
    0x03b2, 0x03cc, 0x03bb, 0x03bf, 0x03c3, 0x002e, 0x0063, 0x006f,
    0x006d, 0},
   IDN2_OK},
  {
   "GnuTLS test #3",
   "xn--fa-hia.de",
   {
    0x0066, 0x0061, 0x00df, 0x002e, 0x0064, 0x0065, 0},
   IDN2_OK},
  {
   "GnuTLS test #4",
   "xn--bcher-kva.de",
   {
    0x0062, 0x00fc, 0x0063, 0x0068, 0x0065, 0x0072, 0x002e, 0x0064,
    0x0065, 0},
   IDN2_OK},
  {
   "GnuTLS test #5",
   "xn--wgv71a119e.jp",
   {
    0x65e5, 0x672c, 0x8a9e, 0x002e, 0x006a, 0x0070, 0},
   IDN2_OK},
  {
   "GnuTLS test #6",
   "xn--fiqu1az03c18t.xn--mxah1amo.com",
   {
    0x7b80, 0x4f53, 0x4e2d, 0x6587, 0x002e, 0x03b5, 0x03be, 0x03c4,
    0x03c1, 0x03b1, 0x002e, 0x0063, 0x006f, 0x006d, 0},
   IDN2_OK},
  {
   "Empty label",
   "xn---.de",
   {
    0x2e, 0x64, 0x65, 0},
   IDN2_PUNYCODE_BAD_INPUT},
  {
   "No ASCII char but delimiter",
   "xn---tda.de",
   {
    0x2e, 0x64, 0x65, 0},
   IDN2_PUNYCODE_BAD_INPUT},

};

static int debug = 1;
static int error_count = 0;
static int break_on_error = 0;

_GL_ATTRIBUTE_FORMAT_PRINTF_STANDARD (1, 2)
     static void fail (const char *format, ...)
{
  va_list arg_ptr;

  va_start (arg_ptr, format);
  vprintf (format, arg_ptr);
  va_end (arg_ptr);

  error_count++;

  if (break_on_error)
    exit (EXIT_FAILURE);
}

static void
ucs4print (const uint32_t * str, size_t len)
{
  size_t i;

  for (i = 0; i < len; i++)
    printf ("U+%04x%s", str[i], (i + 1) % 8 ? " " : "\n");

  if (len % 8)
    printf ("\n");
}

static size_t G_GNUC_IDN2_ATTRIBUTE_PURE
_u32_strlen (const uint32_t * s)
{
  const uint32_t *e;

  for (e = s; *e; e++)
    ;

  return e - s;
}

static size_t G_GNUC_IDN2_ATTRIBUTE_PURE
_u32_strcmp (const uint32_t * s1, const uint32_t * s2)
{
  while (*s1 && *s2 && *s1 == *s2)
    s1++, s2++;

  return *s1 - *s2;
}

static void
_check_4z (const test_t * t, int rc, uint32_t * ucs4, const char *funcname)
{
  if (rc != t->rc_expected
      && !(rc == IDN2_ICONV_FAIL && t->rc_expected == IDN2_ENCODING_ERROR))
    {
      printf ("Test[%u] '%s' failed (got %d, expected %d):\n",
	      (unsigned) (t - test), t->name, rc, t->rc_expected);
      fail ("  %s(): %s\n", funcname, idn2_strerror (rc));
    }
  else if (rc == IDN2_OK)
    {
      if (_u32_strcmp (t->u32_expected, ucs4) != 0)
	{
	  if (debug)
	    {
	      printf ("got:\n");
	      ucs4print (ucs4, _u32_strlen (ucs4));
	      printf ("expected:\n");
	      ucs4print (t->u32_expected, _u32_strlen (t->u32_expected));
	    }

	  fail ("%s() entry %u mismatch\n", funcname, (unsigned) (t - test));
	}
    }
  else if (debug)
    printf ("returned %d expected %d (%s)\n",
	    rc, t->rc_expected, idn2_strerror (t->rc_expected));

  free (ucs4);
}

int
main (void)
{
  static uint32_t abc_u32[] = { 'a', 'b', 'c', 0 };
  static uint32_t xntda_u32[] = { 'x', 'n', '-', '-', 't', 'd', 'a', 0 };

  uint32_t q[128];
  uint32_t *ucs4, *punycode_u32;
  uint8_t *utf8;
  char *utf8_lz;
  const char *encoding;
  size_t outlen, outlen2;
  int rc, skip_lz = 0;
  unsigned i;

  /* Need to set UTF-8 for u8_strconv_from_locale / u8_strconv_to_locale to work.
   * At least on Debian with libunistring 0.9.6+really0.9.3-0.1 and LC_ALL=C valgrind
   * reports Conditional jump or move depends on uninitialised value */
  setlocale (LC_ALL, "C.UTF-8");
  encoding = locale_charset ();

  if (debug)
    printf ("charset=%s\n", encoding);

  if (strcmp (encoding, "UTF-8") != 0)
    skip_lz = 1;

  for (i = 0; i < sizeof (test) / sizeof (test[0]); i++)
    {
      const test_t *t = test + i;

      if (debug)
	printf ("\nPUNYCODE entry %u: %s\n", i, t->name);

      ucs4 = NULL;		/* freed by _check_4z */
      rc = idn2_to_unicode_8z4z (t->punycode, &ucs4, 0);
      _check_4z (t, rc, ucs4, "idn2_to_unicode_8z4z");

      punycode_u32 = u8_to_u32 ((uint8_t *) t->punycode,
				strlen (t->punycode) + 1, NULL, &outlen);

      if (punycode_u32)
	{
	  ucs4 = NULL;		/* freed by _check_4z */
	  rc = idn2_to_unicode_4z4z (punycode_u32, &ucs4, 0);
	  _check_4z (t, rc, ucs4, "idn2_to_unicode_4z4z");

	  outlen2 = sizeof (q) / sizeof (q[0]) - 1;
	  rc = idn2_to_unicode_44i (punycode_u32, outlen - 1, q, &outlen2, 0);
	  ucs4 = u32_cpy_alloc (q, outlen2 + 1);
	  ucs4[outlen2] = 0;
	  _check_4z (t, rc, ucs4, "idn2_to_unicode_44i");

	  free (punycode_u32);
	}

      ucs4 = NULL;
      rc = idn2_to_unicode_8z8z (t->punycode, (char **) &utf8, 0);
      if (rc == IDN2_OK)
	{
	  ucs4 = u8_to_u32 (utf8, u8_strlen (utf8) + 1, NULL, &outlen);
	  free (utf8);
	}
      _check_4z (t, rc, ucs4, "idn2_to_unicode_8z8z");

      if (skip_lz)
	continue;

      ucs4 = NULL;
      rc = idn2_to_unicode_8zlz (t->punycode, &utf8_lz, 0);
      if (rc == IDN2_OK)
	{
	  utf8 = u8_strconv_from_encoding (utf8_lz, encoding, iconveh_error);
	  free (utf8_lz);
	  ucs4 = u8_to_u32 (utf8, u8_strlen (utf8) + 1, NULL, &outlen);
	  free (utf8);
	}
      _check_4z (t, rc, ucs4, "idn2_to_unicode_8zlz");

      /* Since the test punycodes are completely ASCII,
         idn2_to_unicode_8zlz and idn2_to_unicode_lzlz should
         have the same results */
      ucs4 = NULL;
      rc = idn2_to_unicode_lzlz (t->punycode, (char **) &utf8_lz, 0);
      if (rc == IDN2_OK)
	{
	  utf8 = u8_strconv_from_encoding (utf8_lz, encoding, iconveh_error);
	  free (utf8_lz);
	  ucs4 = u8_to_u32 (utf8, u8_strlen (utf8) + 1, NULL, &outlen);
	  free (utf8);
	}
      _check_4z (t, rc, ucs4, "idn2_to_unicode_lzlz");
    }

  /* special checks with NULL values */
  idn2_to_unicode_8z4z (NULL, NULL, 0);
  idn2_to_unicode_8z4z (NULL, &ucs4, 0);
  idn2_to_unicode_8z4z ("abc", NULL, 0);
  idn2_to_unicode_8z4z ("xn--tda", NULL, 0);

  idn2_to_unicode_4z4z (NULL, NULL, 0);
  idn2_to_unicode_4z4z (NULL, &ucs4, 0);
  idn2_to_unicode_4z4z (abc_u32, NULL, 0);
  idn2_to_unicode_4z4z (xntda_u32, NULL, 0);

  idn2_to_unicode_44i (NULL, 0, NULL, NULL, 0);
  idn2_to_unicode_44i (NULL, 0, NULL, &outlen, 0);
  idn2_to_unicode_44i (NULL, 0, q, NULL, 0);
  outlen = 32;
  idn2_to_unicode_44i (NULL, 0, q, &outlen, 0);
  outlen = 0;
  idn2_to_unicode_44i (NULL, 0, (uint32_t *) 123, &outlen, 0);

  idn2_to_unicode_44i (abc_u32, 0, NULL, NULL, 0);
  idn2_to_unicode_44i (abc_u32, 0, NULL, &outlen, 0);
  idn2_to_unicode_44i (abc_u32, 0, q, NULL, 0);
  outlen = 32;
  idn2_to_unicode_44i (abc_u32, 0, q, &outlen, 0);
  outlen = 0;
  idn2_to_unicode_44i (abc_u32, 0, (uint32_t *) 123, &outlen, 0);
  outlen = 0;
  idn2_to_unicode_44i (abc_u32, 3, (uint32_t *) 123, &outlen, 0);
  outlen = 0;
  idn2_to_unicode_44i (abc_u32, 3, (uint32_t *) 123, NULL, 0);

  idn2_to_unicode_8z8z (NULL, NULL, 0);
  idn2_to_unicode_8z8z (NULL, (char **) &utf8, 0);
  idn2_to_unicode_8z8z ("abc", NULL, 0);
  idn2_to_unicode_8z8z ("xn--tda", NULL, 0);

  idn2_to_unicode_8zlz (NULL, NULL, 0);
  idn2_to_unicode_8zlz (NULL, (char **) &utf8, 0);
  idn2_to_unicode_8zlz ("abc", NULL, 0);
  idn2_to_unicode_8zlz ("xn--tda", NULL, 0);

  idn2_to_unicode_lzlz (NULL, NULL, 0);
  idn2_to_unicode_lzlz (NULL, (char **) &utf8, 0);
  idn2_to_unicode_lzlz ("abc", NULL, 0);
  idn2_to_unicode_lzlz ("xn--tda", NULL, 0);

  if (debug && error_count)
    printf ("error_count: %d\n", error_count);

  return !!error_count;
}

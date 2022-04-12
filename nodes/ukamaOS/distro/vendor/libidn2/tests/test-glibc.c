/* test-glibc.c -- API trace test extracted from the glibc AI_IDN tests.
   Copyright (C) 2019 Red Hat, Inc.

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

/* Before changing test expecations in this file, please contact the
   glibc developers on the libc-alpha mailing list to check if these
   changes are benign and will not lead to glibc test suite failures.
   Thanks.  */

#include <config.h>

#include <locale.h>
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <wchar.h>

#include <uniconv.h>

#include <idn2.h>

/* This assumes that wchar_t uses a UTF-16 or UTF-32 encoding.  */
static const wchar_t L_naemchen[] =
  { L'n', 0344, L'm', L'c', L'h', L'e', L'n', 0 };
static const char *naemchen_latin1 = "n\344mchen";
static const char *naemchen_utf8 = "n\xC3\xA4mchen";

static const wchar_t L_shem[] = { 0x05E9, 0x05DD, 0 };

static const char *shem_utf8 = "\xD7\xA9\xD7\x9D";

/* Detected charset.  Note that charset_latin1 covers both ISO-8859-1
   and ISO-8859-15.  */
enum charset_kind
{
  charset_utf8, charset_latin1, charset_neither
};

/* wcsrtombs with a static buffer.  */
static char * __attribute__((malloc)) wcsrtombs_strdup (const wchar_t *input)
{
  char buf[100];
  const wchar_t *src = input;
  mbstate_t state;

  memset (&state, 0, sizeof (state));

  size_t ret = wcsrtombs (buf, &src, sizeof (buf), &state);
  if (ret == (size_t) -1)
    buf[0] = '\0';

  char *result = strdup (buf);
  if (result == NULL)
    {
      puts ("error: memory allocation failure");
      exit (EXIT_FAILURE);
    }

  return result;
}

static const char *locale;

static enum charset_kind
determine_current_charset_kind (void)
{
  const char *lc_string = locale_charset ();
  enum charset_kind expected;

  if (strcmp (lc_string, "UTF-8") == 0)
    expected = charset_utf8;
  else if (strcmp (lc_string, "ISO-8859-1") == 0
	   || strcmp (lc_string, "ISO-8859-15") == 0
	   || strcmp (lc_string, "CP1252") == 0)
    expected = charset_latin1;
  else
    expected = charset_neither;

  char *naemchen_bytes = wcsrtombs_strdup (L_naemchen);
  char *shem_bytes = wcsrtombs_strdup (L_shem);
  enum charset_kind actual;

  if (strcmp (naemchen_bytes, naemchen_utf8) == 0
      && strcmp (shem_bytes, shem_utf8) == 0)
    actual = charset_utf8;
  else if (strcmp (naemchen_bytes, naemchen_latin1) == 0
	   && strcmp (shem_bytes, "") == 0)
    actual = charset_latin1;
  else
    actual = charset_neither;

  free (shem_bytes);
  free (naemchen_bytes);

  if (expected != actual)
    {
      printf ("error: locale %s: expected charset %u (%s), got %u\n",
	      locale, expected, lc_string, actual);
      exit (EXIT_FAILURE);
    }

  return actual;
}

static int errors;

static void
check_success (const char *func, const char *input, const char *expected,
	       int ret, char *actual)
{
  if (ret != 0)
    {
      printf ("error: locale %s: %s: input \"%s\": %d\n",
	      locale, func, input, ret);
      ++errors;
    }
  else
    {
      if (strcmp (actual, expected) != 0)
	{
	  printf ("error: locale %s: %s: input \"%s\": \"%s\"\n",
		  locale, func, input, actual);
	  ++errors;
	}
      idn2_free (actual);
    }
}

static void
check_lookup_ul_success (const char *input, const char *expected)
{
  char *actual = NULL;
  int ret = idn2_lookup_ul (input, &actual, 0);
  check_success ("idn2_lookup_ul", input, expected, ret, actual);
}

static void
check_to_unicode_lzlz_success (const char *input, const char *expected)
{
  char *actual = NULL;
  int ret = idn2_to_unicode_lzlz (input, &actual, 0);
  check_success ("idn2_to_unicode_lzlz", input, expected, ret, actual);
}

static void
check_to_unicode_lzlz_failure (const char *input, int expected)
{
  char *unexpected = NULL;
  int actual = idn2_to_unicode_lzlz (input, &unexpected, 0);

  if (actual == 0)
    {
      printf ("error: idn2_to_unicode_lzlz: locale %s:"
	      "unexpected success for input \"%s\": \"%s\"\n",
	      locale, input, unexpected);
      ++errors;
      idn2_free (unexpected);
    }
  else if (actual != expected)
    {
      printf ("error: idn2_to_unicode_lzlz: locale %s:"
	      "expected failure %d for input \"%s\", actual %d\n",
	      locale, expected, input, actual);
      ++errors;
    }
}

static void
run_utf8_tests (void)
{
  check_lookup_ul_success ("\327\251\327\2351.example",
			   "xn--1-qic9a.example");
  check_lookup_ul_success ("\327\251\327\235.example", "xn--iebx.example");
  check_lookup_ul_success ("both.cname.idn-cname.n\303\244mchen.example",
			   "both.cname.idn-cname.xn--nmchen-bua.example");
  check_lookup_ul_success ("bu\303\237e.example", "xn--bue-6ka.example");
  check_lookup_ul_success ("n\303\244mchen.example",
			   "xn--nmchen-bua.example");
  check_lookup_ul_success ("n\303\244mchen_zwo.example",
			   "xn--nmchen_zwo-q5a.example");
  check_lookup_ul_success ("with.cname.n\303\244mchen.example",
			   "with.cname.xn--nmchen-bua.example");
  check_lookup_ul_success ("With.idn-cname.n\303\244mchen.example",
			   "with.idn-cname.xn--nmchen-bua.example");

  check_to_unicode_lzlz_success ("non-idn-cname.example",
				 "non-idn-cname.example");
  check_to_unicode_lzlz_success ("non-idn.example", "non-idn.example");
  check_to_unicode_lzlz_success ("non-idn-name.example",
				 "non-idn-name.example");
  check_to_unicode_lzlz_success ("xn--1-qic9a.example",
				 "\327\251\327\2351.example");
  check_to_unicode_lzlz_success ("xn--anderes-nmchen-eib.example",
				 "anderes-n\303\244mchen.example");
  check_to_unicode_lzlz_success ("xn--bue-6ka.example",
				 "bu\303\237e.example");
  check_to_unicode_lzlz_success ("xn--iebx.example",
				 "\327\251\327\235.example");
  check_to_unicode_lzlz_success ("xn--nmchen-bua.example",
				 "n\303\244mchen.example");
  check_to_unicode_lzlz_success ("xn--nmchen_zwo-q5a.example",
				 "n\303\244mchen_zwo.example");

  check_to_unicode_lzlz_failure ("xn---.example", IDN2_PUNYCODE_BAD_INPUT);
  check_to_unicode_lzlz_failure ("xn--x.example", IDN2_PUNYCODE_BAD_INPUT);
}

static void
run_latin1_tests (void)
{
  check_lookup_ul_success ("both.cname.idn-cname.n\344mchen.example",
			   "both.cname.idn-cname.xn--nmchen-bua.example");
  check_lookup_ul_success ("bu\337e.example", "xn--bue-6ka.example");
  check_lookup_ul_success ("n\344mchen.example", "xn--nmchen-bua.example");
  check_lookup_ul_success ("n\344mchen_zwo.example",
			   "xn--nmchen_zwo-q5a.example");
  check_lookup_ul_success ("with.cname.n\344mchen.example",
			   "with.cname.xn--nmchen-bua.example");
  check_lookup_ul_success ("With.idn-cname.n\344mchen.example",
			   "with.idn-cname.xn--nmchen-bua.example");

  check_to_unicode_lzlz_success ("non-idn-cname.example",
				 "non-idn-cname.example");
  check_to_unicode_lzlz_success ("non-idn.example", "non-idn.example");
  check_to_unicode_lzlz_success ("non-idn-name.example",
				 "non-idn-name.example");
  check_to_unicode_lzlz_success ("xn--anderes-nmchen-eib.example",
				 "anderes-n\344mchen.example");
  check_to_unicode_lzlz_success ("xn--bue-6ka.example", "bu\337e.example");
  check_to_unicode_lzlz_success ("xn--nmchen-bua.example",
				 "n\344mchen.example");
  check_to_unicode_lzlz_success ("xn--nmchen_zwo-q5a.example",
				 "n\344mchen_zwo.example");

  check_to_unicode_lzlz_failure ("xn--1-qic9a.example", IDN2_ENCODING_ERROR);
  check_to_unicode_lzlz_failure ("xn--iebx.example", IDN2_ENCODING_ERROR);
  check_to_unicode_lzlz_failure ("xn---.example", IDN2_PUNYCODE_BAD_INPUT);
  check_to_unicode_lzlz_failure ("xn--x.example", IDN2_PUNYCODE_BAD_INPUT);
}

static const char *const locale_candidates[] = {
  "C",
  "C.UTF-8",
  "en_US",
  "en_US.utf8",
  "en_US.iso88591",
  "de_DE",
  "de_DE.utf8",
  "de_DE.iso88591",
  "de_DE.iso885915@euro",
  "fr_FR",
  "fr_FR.utf8",
  "fr_FR.iso88591",
  "he_IL.utf8",
  NULL
};

int
main (void)
{
  bool utf8_seen = false;
  bool latin1_seen = false;

  for (size_t i = 0; locale_candidates[i] != NULL; ++i)
    {
      locale = locale_candidates[i];
      if (setlocale (LC_ALL, locale) == NULL)
	continue;

      switch (determine_current_charset_kind ())
	{
	case charset_utf8:
	  run_utf8_tests ();
	  utf8_seen = true;
	  break;
	case charset_latin1:
	  run_latin1_tests ();
	  latin1_seen = true;
	  break;
	case charset_neither:
	  continue;
	}
    }

  if (!utf8_seen)
    {
      /* Mingw64 does not have a UTF-8 locale.  */
#ifndef __MINGW64__
      puts ("error: no UTF-8 locale found");
      ++errors;
#else
      puts ("warning: no UTF-8 support on Mingw");
#endif
    }

  /* Not everyone has a Latin-1 locale installed.  */
  if (!latin1_seen)
    puts ("warning: no Latin-1 locale found");

  if (!(utf8_seen || latin1_seen))
    {
      puts ("error: no usable locales found");
      ++errors;
    }

  if (errors)
    return EXIT_FAILURE;
  else
    return EXIT_SUCCESS;
}

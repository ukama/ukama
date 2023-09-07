/* idn2.c - command line interface to libidn2
   Copyright (C) 2011-2021 Simon Josefsson, Tim Ruehsen

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

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <locale.h>
#include <unistd.h>

#include <idn2.h>

#include <unitypes.h>
#include <uniconv.h>
#include <unistr.h>
#include <unistring/localcharset.h>

/* Gnulib headers. */
#include "error.h"
#include "gettext.h"
#define _(String) dgettext (PACKAGE, String)
#include "progname.h"
#include "version-etc.h"

#include "idn2_cmd.h"

#include "blurbs.h"

#ifdef __cplusplus
extern				// define a global const variable in C++, C doesn't need it.
#endif
const char version_etc_copyright[] =
  /* Do *not* mark this string for translation.  %s is a copyright
     symbol suitable for this locale, and %d is the copyright
     year.  */
  "Copyright %s 2011-%d Simon Josefsson, Tim Ruehsen.";

static void
usage (int status)
{
  if (status != EXIT_SUCCESS)
    fprintf (stderr, _("Try `%s --help' for more information.\n"),
	     program_name);
  else
    {
      printf (_("\
Usage: %s [OPTION]... [STRINGS]...\n\
"), program_name);
      fputs (_("\
Internationalized Domain Name (IDNA2008) convert STRINGS, or standard input.\n\
\n\
"), stdout);
      fputs (_("\
Command line interface to the Libidn2 implementation of IDNA2008.\n\
\n\
All strings are expected to be encoded in the locale charset.\n\
\n\
To process a string that starts with `-', for example `-foo', use `--'\n\
to signal the end of parameters, as in `idn2 --quiet -- -foo'.\n\
\n\
Mandatory arguments to long options are mandatory for short options too.\n\
"), stdout);
      fputs (_("\
  -h, --help                Print help and exit\n\
  -V, --version             Print version and exit\n\
"), stdout);
      fputs (_("\
  -d, --decode              Decode (punycode) domain name\n\
  -l, --lookup              Lookup domain name (default)\n\
  -r, --register            Register label\n\
"), stdout);
      fputs (_("\
  -T, --tr46t               Enable TR46 transitional processing\n\
  -N, --tr46nt              Enable TR46 non-transitional processing\n\
      --no-tr46             Disable TR46 processing\n\
"), stdout);
      fputs (_("\
      --usestd3asciirules   Enable STD3 ASCII rules\n\
      --no-alabelroundtrip  Disable A-label roundtrip for lookups\n\
      --debug               Print debugging information\n\
      --quiet               Silent operation\n\
"), stdout);
      emit_bug_reporting_address ();
    }
  exit (status);
}

static void
hexdump (const char *prefix, const char *str)
{
  uint8_t *u8;
  uint32_t *u32 = NULL;
  size_t u32len;
  size_t i;
  const char *encoding = locale_charset ();

  u8 = u8_strconv_from_encoding (str, encoding, iconveh_error);
  if (u8)
    u32 = u8_to_u32 (u8, strlen ((char *) u8), NULL, &u32len);

  for (i = 0; i < strlen (str); i++)
    fprintf (stderr, "%s[%lu] = 0x%02x\n",
	     prefix, (unsigned long) i, (unsigned) (str[i] & 0xFF));

  if (u8 && strcmp (str, (char *) u8) != 0)
    for (i = 0; i < strlen ((char *) u8); i++)
      fprintf (stderr, "UTF-8 %s[%lu] = 0x%02x\n",
	       prefix, (unsigned long) i, u8[i] & 0xFF);

  if (u8 && u32)
    for (i = 0; i < u32len; i++)
      fprintf (stderr, "UCS-4 %s[%lu] = U+%04x\n",
	       prefix, (unsigned long) i, u32[i]);
}

static struct gengetopt_args_info args_info;

static void
process_input (char *readbuf, int flags)
{
  size_t len = strlen (readbuf);
  char *output;
  const char *tag;
  int rc;

  if (len && readbuf[len - 1] == '\n')
    readbuf[len - 1] = '\0';

  if (strcmp (readbuf, "show w") == 0)
    {
      puts (WARRANTY);
      return;
    }
  else if (strcmp (readbuf, "show c") == 0)
    {
      puts (CONDITIONS);
      return;
    }

  if (args_info.debug_given)
    hexdump ("input", readbuf);

  if (args_info.register_given)
    {
      rc = idn2_register_ul (readbuf, NULL, &output, flags);
      tag = "register";
    }
  else if (args_info.decode_given)
    {
      rc = idn2_to_unicode_lzlz (readbuf, &output, 0);
      tag = "decode";
    }
  else
    {
      rc = idn2_to_ascii_lz (readbuf, &output, flags);
      tag = "toAscii";
    }

  if (rc == IDN2_OK)
    {
      if (args_info.debug_given)
	hexdump ("output", readbuf);

      printf ("%s\n", output);
      free (output);
    }
  else
    error (EXIT_FAILURE, 0, "%s: %s", tag, idn2_strerror (rc));
}

int
main (int argc, char *argv[])
{
  unsigned cmdn;
  int flags = IDN2_NONTRANSITIONAL;

  setlocale (LC_ALL, "");
  set_program_name (argv[0]);
  bindtextdomain (PACKAGE, LOCALEDIR);
  textdomain (PACKAGE);

  if (cmdline_parser (argc, argv, &args_info) != 0)
    return EXIT_FAILURE;

  if (args_info.version_given)
    {
      version_etc (stdout, "idn2", PACKAGE_NAME, VERSION,
		   "Simon Josefsson, Tim Ruehsen", (char *) NULL);
      return EXIT_SUCCESS;
    }

  if (args_info.help_given)
    usage (EXIT_SUCCESS);

  if (!args_info.quiet_given
      && args_info.inputs_num == 0 && isatty (fileno (stdin)))
    fprintf (stderr, "%s %s\n" GREETING, PACKAGE, VERSION);

  if (args_info.debug_given)
    fprintf (stderr, _("Charset: %s\n"), locale_charset ());

  if (!args_info.quiet_given
      && args_info.inputs_num == 0 && isatty (fileno (stdin)))
    fprintf (stderr, "%s", _("Type each input string on a line by itself, "
			     "terminated by a newline character.\n"));

  if (args_info.tr46t_given)
    flags = IDN2_TRANSITIONAL;
  else if (args_info.tr46nt_given)
    flags = IDN2_NONTRANSITIONAL;
  else if (args_info.no_tr46_given)
    flags = IDN2_NO_TR46;

  if (flags && args_info.usestd3asciirules_given)
    flags |= IDN2_USE_STD3_ASCII_RULES;

  if (flags && args_info.no_alabelroundtrip_given)
    flags |= IDN2_NO_ALABEL_ROUNDTRIP;

  for (cmdn = 0; cmdn < args_info.inputs_num; cmdn++)
    process_input (args_info.inputs[cmdn], flags | IDN2_NFC_INPUT);

  if (!cmdn)
    {
      char *buf = NULL;
      size_t bufsize = 0;

      while (getline (&buf, &bufsize, stdin) > 0)
	process_input (buf, flags);

      free (buf);
    }

  if (ferror (stdin))
    error (EXIT_FAILURE, errno, "%s", _("input error"));

  cmdline_parser_free (&args_info);

  return EXIT_SUCCESS;
}

/* test-lookup.c --- Self tests for IDNA processing
   Copyright (C) 2011-2021 Simon Josefsson
   Copyright (C) 2017-2021 Tim Ruehsen

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
#include <ctype.h>
#include <errno.h>

#include <idn2.h>

#include "unistr.h"		/* u32_to_u8, u8_to_u32 */

struct idna
{
  const char *in;
  const char *out;
  int rc;
  int flags;
};

static const struct idna idna[] = {
  /* Corner cases. */
  {"", "", IDN2_OK},
  {".", ".", IDN2_OK},
  {"..", "..", IDN2_OK},	/* XXX should we disallow this? */

  /* U+19DA */
  {"\xe1\xa7\x9a", "xn--pkf", IDN2_DISALLOWED},

  /* U+32FF */
  {"\xe3\x8b\xbf", "xn--nnqt1l", IDN2_OK, IDN2_NONTRANSITIONAL},
  {"\xe3\x8b\xbf", "xn--nnqt1l", IDN2_OK, IDN2_TRANSITIONAL},
  {"\xe3\x8b\xbf", "xn--nnqt1l", IDN2_OK,
   IDN2_NONTRANSITIONAL | IDN2_NFC_INPUT},
  {"\xe3\x8b\xbf", "xn--nnqt1l", IDN2_OK,
   IDN2_TRANSITIONAL | IDN2_NFC_INPUT},
  {"\xe3\x8b\xbf", "xn--nnqt1l", IDN2_OK,
   IDN2_NONTRANSITIONAL | IDN2_NFC_INPUT | IDN2_ALABEL_ROUNDTRIP},
  {"\xe3\x8b\xbf", "xn--nnqt1l", IDN2_OK,
   IDN2_TRANSITIONAL | IDN2_NFC_INPUT | IDN2_ALABEL_ROUNDTRIP},
  {"\xe3\x8b\xbf", "xn--nnqt1l", IDN2_OK,
   IDN2_NONTRANSITIONAL | IDN2_USE_STD3_ASCII_RULES | IDN2_ALABEL_ROUNDTRIP},
  {"\xe3\x8b\xbf", "xn--nnqt1l", IDN2_OK,
   IDN2_TRANSITIONAL | IDN2_USE_STD3_ASCII_RULES | IDN2_ALABEL_ROUNDTRIP},

  /* UTC's test vectors. */
#include "IdnaTest.inc"

  /* Start of contribution from "Abdulrahman I. ALGhadir" <aghadir@citc.gov.sa>. */
  {"\xd8\xa7\xd9\x84\xd9\x85\xd8\xb1\xd9\x83\xd8\xb2\x2d\xd8\xa7\xd9\x84\xd8\xb3\xd8\xb9\xd9\x88\xd8\xaf\xd9\x8a\x2d\xd9\x84\xd9\x85\xd8\xb9\xd9\x84\xd9\x88\xd9\x85\xd8\xa7\xd8\xaa\x2d\xd8\xa7\xd9\x84\xd8\xb4\xd8\xa8\xd9\x83\xd8\xa9", "xn------nzebbbijg6cvanqv7ec6ooadfebehlc1fg8c"},
  {"\xd8\xb5\xd8\xad\xd8\xa7\xd8\xb1\xd9\x89\x2d\xd9\x86\xd8\xaa",
   "xn----ymckjvv7jwa"},
  {"\xd8\xa5\xd8\xaf\xd8\xa7\xd8\xb1\xd8\xa9\xd8\xa7\xd9\x84\xd8\xaa\xd8\xb1\xd8\xa8\xd9\x8a\xd8\xa9\xd9\x88\xd8\xa7\xd9\x84\xd8\xaa\xd8\xb9\xd9\x84\xd9\x8a\xd9\x85\x2d\xd9\x84\xd9\x84\xd8\xa8\xd9\x86\xd8\xa7\xd8\xaa\x2d\xd8\xa8\xd9\x85\xd8\xad\xd8\xa7\xd9\x81\xd8\xb8\xd8\xa9\xd8\xa7\xd9\x84\xd8\xb9\xd9\x84\xd8\xa7", "xn-----4sdiaabbaaeccgchhdd6d1a9bd4pqal6sqcfcbalbsi6a9d2dh"},
  {"\xd8\xa8\xd9\x86\xd9\x83\x2d\xd8\xa7\xd9\x84\xd8\xa8\xd9\x84\xd8\xa7\xd8\xaf", "xn----zmcabc7b5ikbo"},
  {"\xd8\xb4\xd8\xb1\xd9\x83\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd8\xa8\xd9\x84\xd8\xa7\xd8\xaf\x2d\xd9\x84\xd9\x84\xd8\xa7\xd8\xb3\xd8\xaa\xd8\xab\xd9\x85\xd8\xa7\xd8\xb1", "xn-----ctdabadfpk4bslyf5vuabdaz"},
  {"\xd8\xa7\xd9\x86\xd8\xac\xd8\xa7\xd8\xb2", "xn--mgbaoz1h"},
  {"\xd9\x85\xd8\xa4\xd8\xb3\xd8\xb3\xd8\xa9\xd8\xa7\xd9\x84\xd8\xa7\xd9\x86\xd8\xb8\xd9\x85\xd8\xa9\xd8\xa7\xd9\x84\xd9\x85\xd8\xaa\xd8\xb1\xd8\xa7\xd8\xa8\xd8\xb7\xd8\xa9", "xn--jgbgaaagccdh1fra9di4qehidp"},
  {"\xd8\xb3\xd9\x86\xd8\xa7\xd9\x81\xd9\x8a", "xn--mgbx7bsw"},
  {"\xd8\xa3\xd9\x85\xd8\xa7\xd9\x86\xd8\xa9\xd8\xa7\xd9\x84\xd8\xb9\xd8\xa7\xd8\xb5\xd9\x85\xd8\xa9\xd8\xa7\xd9\x84\xd9\x85\xd9\x82\xd8\xaf\xd8\xb3\xd8\xa9", "xn--igbiaaajcb7czbs0c4hwafghdi"},
  {"\xd8\xa7\xd9\x84\xd8\xb4\xd8\xb1\xd9\x83\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd9\x81\xd9\x86\xd9\x8a\xd8\xa9\x2d\xd9\x84\xd8\xaa\xd9\x88\xd8\xb7\xd9\x8a\xd9\x86\x2d\xd8\xa7\xd9\x84\xd8\xaa\xd9\x82\xd9\x86\xd9\x8a\xd8\xa9", "xn------nzebcjcdhc9eubzc8ixafpgde0bff6b9bgh"},
  {"\xd9\x85\xd8\xa4\xd8\xb3\xd8\xb3\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd8\xa7\xd9\x86\xd8\xb8\xd9\x85\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd9\x85\xd8\xaa\xd8\xb1\xd8\xa7\xd8\xa8\xd8\xb7\xd8\xa9", "xn-----1sdnabaicdej5gtaa9ej8sfhjeq"},
  {"\xd9\x85\xd8\xa4\xd8\xb3\xd8\xb3\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd8\xa7\xd9\x86\xd8\xb8\xd9\x85\xd8\xa9\xd8\xa7\xd9\x84\xd9\x85\xd8\xaa\xd8\xb1\xd8\xa7\xd8\xa8\xd8\xb7\xd8\xa9", "xn----smckaaahcddi8fsaa4ej6rehjdq"},
  {"\xd9\x85\xd8\xa4\xd8\xb3\xd8\xb3\xd8\xa9\xd8\xa7\xd9\x84\xd8\xa7\xd9\x86\xd8\xb8\xd9\x85\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd9\x85\xd8\xaa\xd8\xb1\xd8\xa7\xd8\xa8\xd8\xb7\xd8\xa9", "xn----smcjabahccei8fsaa4ei6rfhiep"},
  {"\xd9\x85\xd8\xa4\xd8\xb3\xd8\xb3\xd8\xa9\xd8\xa7\xd9\x84\xd9\x85\xd8\xaf\xd8\xa7\xd8\xb1\xd8\xa7\xd9\x84\xd8\xaa\xd9\x82\xd9\x86\xd9\x8a", "xn--jgbgaahj6arna6osaedgx2e"},
  {"\xd8\xb4\xd8\xb1\xd9\x83\xd8\xa9\x2d\xd8\xa3\xd9\x85\xd9\x88\xd8\xa7\xd9\x84\x2d\xd9\x84\xd9\x84\xd8\xa7\xd8\xb3\xd8\xaa\xd8\xb4\xd8\xa7\xd8\xb1\xd8\xa7\xd8\xaa\x2d\xd8\xa7\xd9\x84\xd9\x85\xd8\xa7\xd9\x84\xd9\x8a\xd8\xa9", "xn------7yeubaabamkhc9gi5akj50azabakbiq1e0d"},
  {"\xd8\xb4\xd8\xb1\xd9\x83\xd8\xa9\x2d\xd8\xa3\xd9\x85\xd9\x88\xd8\xa7\xd9\x84\x2d\xd9\x84\xd9\x84\xd8\xa7\xd8\xb3\xd8\xaa\xd8\xb4\xd8\xa7\xd8\xb1\xd8\xa7\xd8\xaa\x2d\xd8\xa7\xd9\x84\xd9\x85\xd8\xa7\xd9\x84\xd9\x8a\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd9\x85\xd8\xad\xd8\xaf\xd9\x88\xd8\xaf\xd8\xa9", "xn-------g5fybaababokchc4dwbaxi7bqj59a5abakbdlqg8f0a2e"},
  {"\xd9\x85\xd9\x83\xd8\xaa\xd8\xa8\x2d\xd8\xaf\xd8\xa7\xd8\xb1\xd8\xa7\xd9\x84\xd8\xaa\xd9\x85\xd9\x88\xd9\x8a\xd9\x84\x2d\xd9\x84\xd9\x84\xd8\xae\xd8\xaf\xd9\x85\xd8\xa7\xd8\xaa\x2d\xd8\xa7\xd9\x84\xd8\xaa\xd8\xac\xd8\xa7\xd8\xb1\xd9\x8a\xd8\xa9", "xn------ozeabbabsbecc4a0amf7am37a3abbaggkg7f0cq"},
  {"\xd9\x87\xd9\x8a\xd8\xa6\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd8\xa3\xd9\x85\xd8\xb1\x2d\xd8\xa8\xd8\xa7\xd9\x84\xd9\x85\xd8\xb9\xd8\xb1\xd9\x88\xd9\x81\x2d\xd9\x88\xd8\xa7\xd9\x84\xd9\x86\xd9\x87\xd9\x8a\x2d\xd8\xb9\xd9\x86\x2d\xd8\xa7\xd9\x84\xd9\x85\xd9\x86\xd9\x83\xd8\xb1", "xn--------ochtjcbcgi9idf3ke2mwbgffejflwcedu1ac5cxa"},
  {"\xd8\xa7\xd9\x84\xd8\xb3\xd8\xad\xd9\x8a\xd9\x84\xd9\x8a\x2d\xd9\x84\xd9\x84\xd8\xaa\xd8\xac\xd8\xa7\xd8\xb1\xd8\xa9\x2d\xd9\x88\xd8\xa7\xd9\x84\xd8\xa7\xd9\x86\xd9\x85\xd8\xa7\xd8\xa1", "xn-----usdvbbaaoiwj0dwa8wcbahvu4b9ab"},
  {"\xd8\xa7\xd9\x84\xd9\x85\xd8\xa4\xd8\xb3\xd8\xb3\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd8\xb3\xd8\xb9\xd9\x88\xd8\xaf\xd9\x8a\xd8\xa9\x2d\xd9\x84\xd9\x84\xd8\xb7\xd8\xa7\xd9\x82\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd9\x83\xd9\x87\xd8\xb1\xd8\xa8\xd8\xa7\xd8\xa6\xd9\x8a\xd8\xa9", "xn------bzenbcbbakfccf0gvbzaad1gvb6p1ahgfagj8fsa1er"},
  {"\xd8\xb4\xd8\xb1\xd9\x83\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd8\xb5\xd9\x86\xd8\xa7\xd8\xb9\xd8\xa7\xd8\xaa\x2d\xd8\xa7\xd9\x84\xd9\x83\xd9\x87\xd8\xb1\xd8\xa8\xd8\xa7\xd8\xa6\xd9\x8a\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd9\x85\xd8\xaa\xd8\xb9\xd8\xaf\xd8\xaf\xd9\x87", "xn------lzedaabachfjii4fati4cza3gm5tkashi3an3bn3g"},
  {"\xd8\xa7\xd9\x84\xd9\x87\xd9\x8a\xd8\xa6\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd8\xb9\xd8\xa7\xd9\x85\xd8\xa9\x2d\xd9\x84\xd9\x84\xd8\xba\xd8\xb0\xd8\xa7\xd8\xa1\x2d\xd9\x88\xd8\xa7\xd9\x84\xd8\xaf\xd9\x88\xd8\xa7\xd8\xa1", "xn------0yebzgcabcard9gl3mva5peeagn6b5bd8a"},
  {"\xd8\xa7\xd9\x84\xd9\x85\xd8\xaa\xd8\xad\xd8\xaf", "xn--mgbgji6hg"},
  {"\xd8\xad\xd8\xa7\xd8\xb3\xd8\xa8", "xn--mgbcnz"},
  {"\xd9\x85\xd8\xb5\xd8\xb1\xd9\x81\x2d\xd8\xa7\xd9\x84\xd8\xa5\xd9\x86\xd9\x85\xd8\xa7\xd8\xa1", "xn----nmclhb0d1a1h3aehl"},
  {"\xd8\xa7\xd9\x84\xd8\xa3\xd9\x84\xd9\x88\xd9\x83\xd8\xa9",
   "xn--igbhh7hdb2a"},
  {"\xd8\xa5\xd8\xb9\xd9\x85\xd8\xa7\xd8\xb1", "xn--kgbe4a4a1d"},
  {"\xd9\x81\xd8\xa7\xd8\xb1\xd9\x85", "xn--mgbu0cs"},
  {"\xd9\x81\xd9\x86\xd8\xaa\xd9\x88\xd8\xb1\xd9\x8a", "xn--pgbo0culn"},
  {"\xd8\xa7\xd9\x84\xd8\xaa\xd8\xa7\xd8\xac", "xn--mgbaij1j"},
  {"\xd8\xb4\xd8\xb1\xd9\x83\xd8\xa9\x2d\xd8\xb1\xd8\xa7\xd9\x85\xd8\xa7\xd8\xaa\x2d\xd8\xa7\xd9\x84\xd8\xaf\xd9\x88\xd9\x84\xd9\x8a\xd8\xa9\x2d\xd9\x84\xd9\x84\xd8\xaa\xd9\x82\xd9\x86\xd9\x8a\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd9\x85\xd8\xad\xd8\xaf\xd9\x88\xd8\xaf\xd8\xa9", "xn-------05fabckfbcfe2czahavc3dvuka6abcafmr0a0dp7ch"},
  {"\xd9\x85\xd8\xa4\xd8\xb3\xd8\xb3\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd8\xb2\xd9\x86\xd9\x8a\xd8\xaa\xd8\xa7\xd9\x86\x2d\xd8\xa7\xd9\x84\xd8\xaa\xd8\xac\xd8\xa7\xd8\xb1\xd9\x8a\xd8\xa9", "xn-----1sdnabakgfdy0egla60afg2ac9fk"},
  {"\xd8\xb4\xd8\xb1\xd9\x83\xd8\xa9\x2d\xd8\xb7\xd9\x88\xd8\xaf\x2d\xd9\x84\xd8\xa5\xd8\xaf\xd8\xa7\xd8\xb1\xd8\xa9\x2d\xd9\x88\xd8\xaa\xd8\xb3\xd9\x88\xd9\x8a\xd9\x82\x2d\xd8\xa7\xd9\x84\xd8\xb9\xd9\x82\xd8\xa7\xd8\xb1", "xn-------r5fmcakem8ccwhg4af4duc6ldf3al3gjc2d"},
  {"\xd9\x88\xd8\xb2\xd8\xa7\xd8\xb1\xd8\xa9\xd8\xa7\xd9\x84\xd8\xaa\xd8\xac\xd8\xa7\xd8\xb1\xd8\xa9\xd9\x88\xd8\xa7\xd9\x84\xd8\xb5\xd9\x86\xd8\xa7\xd8\xb9\xd8\xa9", "xn--mgbaaaaicceu5cff6c7cxjg1bxl"},
  {"\xd8\xa7\xd9\x84\xd8\xa3\xd9\x88\xd9\x84\xd9\x89\x2d\xd9\x84\xd9\x84\xd8\xaa\xd8\xb7\xd9\x88\xd9\x8a\xd8\xb1", "xn----qmclo9a9a1fbba4bghu"},
  {"\xd8\xb4\xd8\xb1\xd9\x83\xd8\xa9\x2d\xd8\xb9\xd8\xb0\xd9\x8a\xd8\xa8\x2d\xd9\x86\xd8\xaa\x2d\xd8\xb3\xd9\x88\xd9\x84\x2d\xd8\xa7\xd9\x84\xd8\xb3\xd8\xb9\xd9\x88\xd8\xaf\xd9\x8a\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd9\x85\xd8\xad\xd8\xaf\xd9\x88\xd8\xaf\xd8\xa9", "xn--------gdhbchgcf9byaeaep9bcj7ij1u8acg2ao7dgi7bp"},
  {"\xd8\xb4\xd8\xb1\xd9\x83\xd8\xa9\x2d\xd8\xb9\xd8\xb0\xd9\x8a\xd8\xa8\x2d\xd9\x84\xd9\x84\xd9\x83\xd9\x85\xd8\xa8\xd9\x8a\xd9\x88\xd8\xaa\xd8\xb1\x2d\xd9\x88\xd8\xa7\xd9\x84\xd8\xa7\xd8\xaa\xd8\xb5\xd8\xa7\xd9\x84\xd8\xa7\xd8\xaa\x2d\xd8\xa7\xd9\x84\xd9\x85\xd8\xad\xd8\xaf\xd9\x88\xd8\xaf\xd8\xa9", "xn-------25faaabcbilgdc7cwbaesh4dvb8f3mg1aageerq3gdp2ch"},
  {"\xd8\xb4\xd8\xb1\xd9\x83\xd8\xa9\x2d\xd9\x85\xd8\xb1\xd8\xa8\xd8\xb7\x2d\xd8\xb9\xd8\xb0\xd8\xa8\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd9\x85\xd8\xad\xd8\xaf\xd9\x88\xd8\xaf\xd8\xa9", "xn------qzecbdec3bwagjc8b4cwb1n6akl5e"},
  {"\xd9\x87\xd9\x8a\xd8\xa6\xd8\xa9\x2d\xd8\xaa\xd9\x86\xd8\xb8\xd9\x8a\xd9\x85\x2d\xd8\xa7\xd9\x84\xd9\x83\xd9\x87\xd8\xb1\xd8\xa8\xd8\xa7\xd8\xa1\x2d\xd9\x88\xd8\xa7\xd9\x84\xd8\xa7\xd9\x86\xd8\xaa\xd8\xa7\xd8\xac\x2d\xd8\xa7\xd9\x84\xd9\x85\xd8\xb2\xd8\xaf\xd9\x88\xd8\xac", "xn-------64f1ajacaabfkqi9ac0d2a6a1j9kxahgjsioll2bn2bg"},
  {"\xd8\xb4\xd8\xb1\xd9\x83\xd8\xa9\x2d\xd8\xb9\xd8\xb0\xd9\x8a\xd8\xa8\x2d\xd9\x84\xd9\x84\xd8\xae\xd8\xaf\xd9\x85\xd8\xa7\xd8\xaa\x2d\xd8\xa7\xd9\x84\xd8\xb7\xd8\xa8\xd9\x8a\xd8\xa9", "xn------pzebcebhg6bmjl8b8cya1pzaagr2io"},
  {"\xd8\xb7\xd9\x8a\xd8\xb1\xd8\xa7\xd9\x86\x2d\xd9\x86\xd8\xa7\xd8\xb3",
   "xn----ymcb1bnt1ib5a"},
  {"\xd8\xb4\xd8\xb1\xd9\x83\xd8\xa9\x2d\xd8\xb9\xd8\xb0\xd9\x8a\xd8\xa8\x2d\xd8\xa7\xd9\x84\xd8\xaa\xd8\xac\xd8\xa7\xd8\xb1\xd9\x8a\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd9\x85\xd8\xad\xd8\xaf\xd9\x88\xd8\xaf\xd8\xa9", "xn------pzeabcgfcgyr2aaeoj0czg1j3ahy1fsbi"},
  {"\xd9\x86\xd8\xa7\xd8\xb3", "xn--mgby7c"},
  {"\xd8\xb4\xd8\xb1\xd9\x83\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd9\x81\xd9\x86\xd8\xa7\xd8\xb1\x2d\xd9\x84\xd9\x84\xd8\xa5\xd8\xb3\xd8\xaa\xd8\xab\xd9\x85\xd8\xa7\xd8\xb1\x2d\xd8\xa7\xd9\x84\xd8\xaa\xd8\xac\xd8\xa7\xd8\xb1\xd9\x8a", "xn------hzeiacbalqdjs6deff2al1zmb0aeaiws9k"},
  {"\xd8\xa7\xd9\x84\xd9\x81\xd9\x86\xd8\xa7\xd8\xb1\x2d\xd9\x84\xd8\xa3\xd9\x86\xd8\xb8\xd9\x85\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd8\xa8\xd9\x86\xd8\xa7\xd8\xa1", "xn-----usdgsadaih9f3drftbefnlfh"},
  {"\xd8\xa7\xd9\x84\xd9\x81\xd9\x86\xd8\xa7\xd8\xb1\x2d\xd9\x84\xd9\x84\xd8\xa3\xd9\x86\xd8\xb8\xd9\x85\xd8\xa9\x2d\xd8\xa7\xd9\x84\xd9\x83\xd9\x87\xd8\xb1\xd8\xa8\xd8\xa7\xd8\xa6\xd9\x8a\xd8\xa9", "xn-----zsdnbadaihf1gf8g4frbheafrog5a3f"},
  {"\xd8\xb4\xd8\xb1\xd9\x83\xd8\xa9\x2d\xd8\xb9\xd8\xb0\xd9\x8a\xd8\xa8\x2d\xd8\xa8\xd9\x8a\xd8\xb1\xd8\xaf\xd8\xa7\xd9\x86\xd8\xa7\x2d\xd9\x84\xd9\x84\xd9\x85\xd9\x82\xd8\xa7\xd9\x88\xd9\x84\xd8\xa7\xd8\xaa\x2d\xd8\xa7\xd9\x84\xd9\x85\xd8\xad\xd8\xaf\xd9\x88\xd8\xaf\xd8\xa9", "xn-------15fababcbill1cxajaerg2d2g1jma2bacewis7ej9bd"},
  /* End of contribution from "Abdulrahman I. ALGhadir" <aghadir@citc.gov.sa>. */

  /* These comes from http://www.iana.org/domains/root/db see
     gen-idn-tld-tv.pl */
  {"\xe6\xb5\x8b\xe8\xaf\x95", "xn--0zwm56d"},
  {"\xe0\xa4\xaa\xe0\xa4\xb0\xe0\xa5\x80\xe0\xa4\x95\xe0\xa5\x8d\xe0\xa4\xb7\xe0\xa4\xbe", "xn--11b5bs3a9aj6g"},
  {"\xed\x95\x9c\xea\xb5\xad", "xn--3e0b707e"},
  {"\xe0\xa6\xad\xe0\xa6\xbe\xe0\xa6\xb0\xe0\xa6\xa4", "xn--45brj9c"},
  {"\xd0\x98\xd0\xa1\xd0\x9f\xd0\xab\xd0\xa2\xd0\x90\xd0\x9d\xd0\x98\xd0\x95", "xn--80akhbyknj4f", IDN2_DISALLOWED, IDN2_NO_TR46},	/* iana bug */
  {"испытание", "xn--80akhbyknj4f"},	/* corrected */
  {"\xd0\xa1\xd0\xa0\xd0\x91", "xn--90a3ac", IDN2_DISALLOWED, IDN2_NO_TR46},	/* iana bug */
  {"срб", "xn--90a3ac"},	/* corrected */
  {"\xed\x85\x8c\xec\x8a\xa4\xed\x8a\xb8", "xn--9t4b11yi5a"},
  {"\xe0\xae\x9a\xe0\xae\xbf\xe0\xae\x99\xe0\xaf\x8d\xe0\xae\x95\xe0\xae\xaa\xe0\xaf\x8d\xe0\xae\xaa\xe0\xaf\x82\xe0\xae\xb0\xe0\xaf\x8d", "xn--clchc0ea0b2g2a9gcd"},
  {"\xd7\x98\xd7\xa2\xd7\xa1\xd7\x98", "xn--deba0ad"},
  {"\xe4\xb8\xad\xe5\x9b\xbd", "xn--fiqs8s"},
  {"\xe4\xb8\xad\xe5\x9c\x8b", "xn--fiqz9s"},
  {"\xe0\xb0\xad\xe0\xb0\xbe\xe0\xb0\xb0\xe0\xb0\xa4\xe0\xb1\x8d",
   "xn--fpcrj9c3d"},
  {"\xe0\xb6\xbd\xe0\xb6\x82\xe0\xb6\x9a\xe0\xb7\x8f", "xn--fzc2c9e2c"},
  {"\xe6\xb8\xac\xe8\xa9\xa6", "xn--g6w251d"},
  {"\xe0\xaa\xad\xe0\xaa\xbe\xe0\xaa\xb0\xe0\xaa\xa4", "xn--gecrj9c"},
  {"\xe0\xa4\xad\xe0\xa4\xbe\xe0\xa4\xb0\xe0\xa4\xa4", "xn--h2brj9c"},
  {"\xd8\xa2\xd8\xb2\xd9\x85\xd8\xa7\xdb\x8c\xd8\xb4\xdb\x8c",
   "xn--hgbk6aj7f53bba"},
  {"\xe0\xae\xaa\xe0\xae\xb0\xe0\xae\xbf\xe0\xae\x9f\xe0\xaf\x8d\xe0\xae\x9a\xe0\xaf\x88", "xn--hlcj6aya9esc7a"},
  {"\xe9\xa6\x99\xe6\xb8\xaf", "xn--j6w193g"},
  {"\xce\x94\xce\x9f\xce\x9a\xce\x99\xce\x9c\xce\x89", "xn--jxalpdlp", IDN2_DISALLOWED, IDN2_NO_TR46},	/* iana bug */
  {"δοκιμή", "xn--jxalpdlp"},
  {"\xd8\xa5\xd8\xae\xd8\xaa\xd8\xa8\xd8\xa7\xd8\xb1", "xn--kgbechtv"},
  {"\xe5\x8f\xb0\xe6\xb9\xbe", "xn--kprw13d"},
  {"\xe5\x8f\xb0\xe7\x81\xa3", "xn--kpry57d"},
  {"\xd8\xa7\xd9\x84\xd8\xac\xd8\xb2\xd8\xa7\xd8\xa6\xd8\xb1",
   "xn--lgbbat1ad8j"},
  {"\xd8\xb9\xd9\x85\xd8\xa7\xd9\x86", "xn--mgb9awbf"},
  {"\xd8\xa7\xdb\x8c\xd8\xb1\xd8\xa7\xd9\x86", "xn--mgba3a4f16a"},
  {"\xd8\xa7\xd9\x85\xd8\xa7\xd8\xb1\xd8\xa7\xd8\xaa", "xn--mgbaam7a8h"},
  {"\xd8\xa7\xd9\x84\xd8\xa7\xd8\xb1\xd8\xaf\xd9\x86", "xn--mgbayh7gpa"},
  {"\xd8\xa8\xda\xbe\xd8\xa7\xd8\xb1\xd8\xaa", "xn--mgbbh1a71e"},
  {"\xd8\xa7\xd9\x84\xd9\x85\xd8\xba\xd8\xb1\xd8\xa8", "xn--mgbc0a9azcg"},
  {"\xd8\xa7\xd9\x84\xd8\xb3\xd8\xb9\xd9\x88\xd8\xaf\xd9\x8a\xd8\xa9",
   "xn--mgberp4a5d4ar"},
  {"\xe1\x83\x92\xe1\x83\x94", "xn--node"},
  {"\xe0\xb9\x84\xe0\xb8\x97\xe0\xb8\xa2", "xn--o3cw4h"},
  {"\xd8\xb3\xd9\x88\xd8\xb1\xd9\x8a\xd8\xa9", "xn--ogbpf8fl"},
  {"\xd0\xa0\xd0\xa4", "xn--p1ai", IDN2_DISALLOWED, IDN2_NO_TR46},	/* iana bug */
  {"рф", "xn--p1ai"},		/* corrected */
  {"\xd8\xaa\xd9\x88\xd9\x86\xd8\xb3", "xn--pgbs0dh"},
  {"\xe0\xa8\xad\xe0\xa8\xbe\xe0\xa8\xb0\xe0\xa8\xa4", "xn--s9brj9c"},
  {"\xd9\x85\xd8\xb5\xd8\xb1", "xn--wgbh1c"},
  {"\xd9\x82\xd8\xb7\xd8\xb1", "xn--wgbl6a"},
  {"\xe0\xae\x87\xe0\xae\xb2\xe0\xae\x99\xe0\xaf\x8d\xe0\xae\x95\xe0\xaf\x88",
   "xn--xkc2al3hye2a"},
  {"\xe0\xae\x87\xe0\xae\xa8\xe0\xaf\x8d\xe0\xae\xa4\xe0\xae\xbf\xe0\xae\xaf\xe0\xae\xbe", "xn--xkc2dl3a5ee0h"},
  {"\xe6\x96\xb0\xe5\x8a\xa0\xe5\x9d\xa1", "xn--yfro4i67o"},
  {"\xd9\x81\xd9\x84\xd8\xb3\xd8\xb7\xd9\x8a\xd9\x86", "xn--ygbi2ammx"},
  {"\xe3\x83\x86\xe3\x82\xb9\xe3\x83\x88", "xn--zckzah"},
  /* end of IANA strings */

  /* the following comes from IDNA2003 libidn
     with some new variants inspired by the old test vectors */
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xc3\xad\x64\x6e", "example.xn--dn-mja"
   /* 1-1-1 Has an IDN in just the TLD */
   },
  {"\xc3\xab\x78\x2e\xc3\xad\x64\x6e", "xn--x-ega.xn--dn-mja"
   /* 1-1-2 Has an IDN in the TLD and SLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xc3\xa5\xc3\xbe\xc3\xa7",
   "example.xn--5cae2e"
   /* 1-2-1 Latin-1 TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xc4\x83\x62\xc4\x89",
   "example.xn--b-rhat"
   /* 1-2-2 Latin Extended A TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xc8\xa7\xc6\x80\xc6\x88",
   "example.xn--lhaq98b"
   /* 1-2-3 Latin Extended B TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe1\xb8\x81\xe1\xb8\x83\xe1\xb8\x89",
   "example.xn--2fges"
   /* 1-2-4 Latin Extended Additional TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe4\xb8\xbf\xe4\xba\xba\xe5\xb0\xb8",
   "example.xn--xiqplj17a"
   /* 1-3-1 Han TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe3\x81\x8b\xe3\x81\x8c\xe3\x81\x8d",
   "example.xn--u8jcd"
   /* 1-3-2 Hiragana TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe3\x82\xab\xe3\x82\xac\xe3\x82\xad",
   "example.xn--lckcd"
   /* 1-3-3 Katakana TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe1\x84\x80\xe1\x85\xa1\xe1\x86\xa8",
   "example.xn--p39a",
   IDN2_NOT_NFC, IDN2_NO_TR46
   /* 1-3-4 Hangul Jamo TLD */
   /* Don't resolve as example.xn--ypd8qrh */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xea\xb1\xa9\xeb\x93\x86\xec\x80\xba",
   "example.xn--o69aq2nl0j"
   /* 1-3-5 Hangul TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xea\x80\x8a\xea\x80\xa0\xea\x8a\xb8",
   "example.xn--6l7arby7j"
   /* 1-3-6 Yi TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xce\xb1\xce\xb2\xce\xb3",
   "example.xn--mxacd"
   /* 1-3-7 Greek TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe1\xbc\x82\xe1\xbc\xa6\xe1\xbd\x95",
   "example.xn--fng7dpg"
   /* 1-3-8 Greek Extended TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xd0\xb0\xd0\xb1\xd0\xb2",
   "example.xn--80acd"
   /* 1-3-9 Cyrillic TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xd5\xa1\xd5\xa2\xd5\xa3",
   "example.xn--y9acd"
   /* 1-3-10 Armeian TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe1\x83\x90\xe1\x83\x91\xe1\x83\x92",
   "example.xn--lodcd"
   /* 1-3-11 Georgian TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe2\x88\xa1\xe2\x86\xba\xe2\x8a\x82",
   "example.xn--b7gxomk",
   /* 1-4-1 Symbols TLD */
   IDN2_DISALLOWED		/* valid IDNA2003 invalid IDNA2008 */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe0\xa4\x95\xe0\xa4\x96\xe0\xa4\x97",
   "example.xn--11bcd"
   /* 1-5-1 Devanagari TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe0\xa6\x95\xe0\xa6\x96\xe0\xa6\x97",
   "example.xn--p5bcd"
   /* 1-5-2 Bengali TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe0\xa8\x95\xe0\xa8\x96\xe0\xa8\x97",
   "example.xn--d9bcd"
   /* 1-5-3 Gurmukhi TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe0\xaa\x95\xe0\xaa\x96\xe0\xaa\x97",
   "example.xn--0dccd"
   /* 1-5-4 Gujarati TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe0\xac\x95\xe0\xac\x96\xe0\xac\x97",
   "example.xn--ohccd"
   /* 1-5-5 Oriya TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe0\xae\x95\xe0\xae\x99\xe0\xae\x9a",
   "example.xn--clcid"
   /* 1-5-6 Tamil TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe0\xb0\x95\xe0\xb0\x96\xe0\xb0\x97",
   "example.xn--zoccd"
   /* 1-5-7 Telugu TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe0\xb2\x95\xe0\xb2\x96\xe0\xb2\x97",
   "example.xn--nsccd"
   /* 1-5-8 Kannada TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe0\xb4\x95\xe0\xb4\x96\xe0\xb4\x97",
   "example.xn--bwccd"
   /* 1-5-9 Malayalam TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe0\xb6\x9a\xe0\xb6\x9b\xe0\xb6\x9c",
   "example.xn--3zccd"
   /* 1-5-10 Sinhala TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe0\xb8\x81\xe0\xb8\x82\xe0\xb8\x83",
   "example.xn--12ccd"
   /* 1-5-11 Thai TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe0\xba\x81\xe0\xba\x82\xe0\xba\x84",
   "example.xn--p6ccg"
   /* 1-5-12 Lao TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe0\xbd\x80\xe0\xbd\x81\xe0\xbd\x82",
   "example.xn--5cdcd"
   /* 1-5-13 Tibetan TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe1\x80\x80\xe1\x80\x81\xe1\x80\x82",
   "example.xn--nidcd"
   /* 1-5-14 Myanmar TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe1\x9e\x80\xe1\x9e\x81\xe1\x9e\x82",
   "example.xn--i2ecd"
   /* 1-5-15 Khmer TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xe1\xa0\xa0\xe1\xa0\xa1\xe1\xa0\xa2",
   "example.xn--26ecd"
   /* 1-5-16 Mongolian TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xd8\xa7\xd8\xa8\xd8\xa9",
   "example.xn--mgbcd"
   /* 1-6-1 Arabic TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xd7\x90\xd7\x91\xd7\x92",
   "example.xn--4dbcd"
   /* 1-6-2 Hebrew TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xdc\x90\xdc\x91\xdc\x92",
   "example.xn--9mbcd"
   /* 1-6-3 Syriac TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\x61\x62\x63\xe3\x82\xab\xe3\x82\xac\xe3\x82\xad",
   "example.xn--abc-mj4bfg"
   /* 1-7-1 ASCII and non-Latin TLD */
   },
  {"\x65\x78\x61\x6d\x70\x6c\x65\x2e\xc3\xa5\xc3\xbe\xc3\xa7\xe3\x82\xab\xe3\x82\xac\xe3\x82\xad",
   "example.xn--5cae2e328wfag"
   /* 1-7-2 Latin (non-ASCII) and non-Latin TLD */
   },
  {"\xc3\xad\x21\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65", "",
   IDN2_DISALLOWED, IDN2_NO_TR46
   /* 2-3-1-1 Includes ! before Nameprep */
   /* Don't resolve as xn--!dn-qma.example */
   },
  {"\xc3\xad\x24\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65", "",
   IDN2_DISALLOWED, IDN2_NO_TR46
   /* 2-3-1-2 Includes $ before Nameprep */
   /* Don't resolve as xn--$dn-qma.example */
   },
  {"\xc3\xad\x2b\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65", "",
   IDN2_DISALLOWED, IDN2_NO_TR46
   /* 2-3-1-3 Includes + before Nameprep */
   /* Don't resolve as xn--+dn-qma.example */
   },
  {"\x2d\xc3\xad\x31\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn---1dn-vpa.example", IDN2_OK, IDN2_NO_TR46
   /* 2-3-2-1 Leading hyphen before Nameprep */
   /* Don't resolve as xn---1dn-vpa.example */
   /* Valid according to IDNA2008-lookup! */
   },
  {"\xc3\xad\x31\x64\x6e\x2d\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--1dn--upa.example", IDN2_OK, IDN2_NO_TR46
   /* 2-3-2-2 Trailing hyphen before Nameprep */
   /* Don't resolve as xn--1dn--upa.example */
   /* Valid according to IDNA2008-lookup! */
   },
  {"\xc3\xad\xef\xbc\x8b\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65", "",
   IDN2_DISALLOWED, IDN2_NO_TR46
   /* 2-3-3-1 Gets a + after Nameprep */
   /* Don't resolve as xn--dn-mja0331x.example */
   },
  {"\xc3\xad\xe2\x81\xbc\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65", "",
   IDN2_DISALLOWED, IDN2_NO_TR46
   /* 2-3-3-2 Gets a = after Nameprep */
   /* Don't resolve as xn--dn-mja0343a.example */
   },
  {"\xef\xb9\xa3\xc3\xad\x32\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65", "",
   IDN2_DISALLOWED, IDN2_NO_TR46
   /* 2-3-4-1 Leading hyphen after Nameprep */
   /* Don't resolve as xn--2dn-qma32863a.example */
   /* Don't resolve as xn---2dn-vpa.example */
   },
  {"\xc3\xad\x32\x64\x6e\xef\xbc\x8d\x2e\x65\x78\x61\x6d\x70\x6c\x65", "",
   IDN2_DISALLOWED, IDN2_NO_TR46
   /* 2-3-4-2 Trailing hyphen after Nameprep */
   /* Don't resolve as xn--2dn-qma79363a.example */
   /* Don't resolve as xn--2dn--upa.example */
   },
  {"\xc2\xb9\x31\x2e\x65\x78\x61\x6d\x70\x6c\x65", "11.example",
   IDN2_DISALLOWED, IDN2_NO_TR46
   /* 2-4-1 All-ASCII check, Latin */
   },
  {"\xe2\x85\xa5\x76\x69\x2e\x65\x78\x61\x6d\x70\x6c\x65", "vivi.example",
   IDN2_DISALLOWED, IDN2_NO_TR46
   /* 2-4-2 All-ASCII check, symbol */
   },
  {"\xc3\x9f\x73\x73\x2e\x65\x78\x61\x6d\x70\x6c\x65", "xn--ss-fia.example"
   /* 2-4-3 All-ASCII check, sharp S */
   /* Different output in IDNA2008-lookup compared to IDNA2003! */
   },
  {"\x78\x6e\x2d\x2d\xc3\xaf\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65", "",
   IDN2_2HYPHEN, IDN2_NO_TR46
   /* 2-5-1 ACE prefix before Nameprep, body */
   /* Don't resolve as xn--xn--dn-sja.example */
   /* Don't resolve as xn--dn-sja.example */
   },
  {"\xe2\x85\xb9\x6e\x2d\x2d\xc3\xa4\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_2HYPHEN, IDN2_NO_TR46
   /* 2-5-2 ACE prefix before Nameprep, prefix */
   /* Don't resolve as xn--xn--dn-uia.example */
   /* Don't resolve as xn--dn-uia.example */
   },
  {"", ""
   /* 2-8-1 Zero-length label after Nameprep */
   /* Don't resolve as xn--kba.example */
   /* Don't resolve as xn--.example */
   },
  {"\x33\x30\x30\x32\x2d\x74\x65\x73\x74\xe3\x80\x82\xc3\xad\x64\x6e",
   "3002-test.xn--dn-mja",
   IDN2_DISALLOWED, IDN2_NO_TR46
   /* 2-9-1 U+3002 acts as a label separator */
   /* Don't resolve as xn--3002-testdn-wcb2087m.example */
   /* Not valid in IDNA2008! */
   },
  {"\x66\x66\x30\x65\x2d\x74\x65\x73\x74\xef\xbc\x8e\xc3\xad\x64\x6e",
   "ff0e-test.xn--dn-mja",
   IDN2_DISALLOWED, IDN2_NO_TR46
   /* 2-9-2 U+FF0E acts as a label separator */
   /* Don't resolve as xn--ff0e-testdn-wcb45865f.example */
   /* Not valid in IDNA2008! */
   },
  {"\x66\x66\x36\x31\x2d\x74\x65\x73\x74\xef\xbd\xa1\xc3\xad\x64\x6e",
   "ff61-test.xn--dn-mja",
   IDN2_DISALLOWED, IDN2_NO_TR46
   /* 2-9-3 U+FF61 acts as a label separator */
   /* Don't resolve as xn--ff61-testdn-wcb33975f.example */
   /* Not valid in IDNA2008! */
   },
  {"\x30\x30\x61\x64\x6f\x75\x74\xc2\xad\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--00adoutdn-m5a.example", IDN2_DISALLOWED, IDN2_NO_TR46
   /* 4-1-1-1 00adout<00AD><00ED>dn.example -> 00adout<00ED>dn.example */
   /* Don't resolve as xn--00adoutdn-cna81e.example */
   /* Not valid in IDNA2008! */
   },
  {"\x32\x30\x30\x64\x6f\x75\x74\xe2\x80\x8d\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--200doutdn-m5a.example", IDN2_CONTEXTJ
   /* 4-1-1-2 200dout<200D><00ED>dn.example -> 200dout<00ED>dn.example */
   /* Don't resolve as xn--200doutdn-m5a1678f.example */
   /* Not valid in IDNA2008!" */
   },
  /* To find Virama's, use:
     grep -E '^[^;]+;[^;]+;[^;]+;9;' UnicodeData.txt */
  {"\xe0\xa5\x8d\xe2\x80\x8d", "", IDN2_LEADING_COMBINING
   /* U+094D U+200D => U+094D is combining mark */
   },
  {"foo\xe0\xa5\x8d\xe2\x80\x8d", "xn--foo-umh4320a", IDN2_OK
   /* foo U+094D U+200D => OK due to Virama + U+200D. */
   },
  {"foo𐨿\xe2\x80\x8d\x65\x65", "xn--fooee-zt3bn006o", IDN2_OK
   /* foo U+10A3F U+200D ee => OK due to Virama + U+200D. */
   },
  {"foo྄\xe2\x80\x8d\x65\x65", "xn--fooee-c3s855o", IDN2_OK
   /* foo U+0f84 U+200D ee => OK due to Virama + U+200D. */
   },
  {"foo᮪\xe2\x80\x8d\x65\x65", "xn--fooee-hc8as55a", IDN2_OK
   /* foo U+1bAA (Mc) U+200D ee => OK due to Virama + U+200D. */
   },
  {"\x73\x69\x6d\x70\x6c\x65\x63\x61\x70\x44\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--simplecapddn-1fb.example", IDN2_DISALLOWED, IDN2_NO_TR46
   /* 4-1-2-1 simplecap<0044><00ED>dn.example -> simplecap<0064><00ED>dn.example */
   /* Uppercase not valid in IDNA2008! */
   },
  {"\x6c\x61\x74\x69\x6e\x74\x6f\x67\x72\x65\x65\x6b\xc2\xb5\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--latintogreekdn-cmb716i.example", IDN2_DISALLOWED, IDN2_NO_TR46
   /* 4-1-2-2 latintogreek<00B5><00ED>dn.example -> latintogreek<03BC><00ED>dn.example */
   /* Don't resolve as xn--latintogreekdn-cxa01g.example */
   /* B5 not valid in IDNA2008! */
   },
  {"\x6c\x61\x74\x69\x6e\x65\x78\x74\xc3\x87\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--latinextdn-v6a6e.example", IDN2_DISALLOWED, IDN2_NO_TR46
   /* 4-1-2-3 latinext<00C7><00ED>dn.example -> latinext<00E7><00ED>dn.example */
   /* Don't resolve as xn--latinextdn-twa07b.example */
   /* C7 not valid in IDNA2008! */
   },
  {"\x73\x68\x61\x72\x70\x73\xc3\x9f\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--sharpsdn-vya4l.example"
   /* 4-1-2-4 sharps<00DF><00ED>dn.example -> sharpsss<00ED>dn.example */
   /* Don't resolve as xn--sharpsdn-vya4l.example */
   /* Changed in IDNA2008! */
   },
  {"\x74\x75\x72\x6b\x69\x73\x68\x69\xc4\xb0\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--turkishiidn-wcb701e.example", IDN2_DISALLOWED, IDN2_NO_TR46
   /* 4-1-2-5 turkishi<0130><00ED>dn.example -> turkishi<0069><0307><00ED>dn.example */
   /* Don't resolve as xn--turkishidn-r8a71f.example */
   /* U+0130 not valid in IDNA2008! */
   },
  {"\x65\x78\x70\x74\x77\x6f\xc5\x89\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--exptwondn-m5a502c.example", IDN2_DISALLOWED, IDN2_NO_TR46
   /* 4-1-2-6 exptwo<0149><00ED>dn.example -> exptwo<02BC><006E><00ED>dn.example */
   /* Don't resolve as xn--exptwodn-h2a33g.example */
   /* U+0149 not valid in IDNA2008 */
   },
  {"\x61\x64\x64\x66\x6f\x6c\x64\xcf\x92\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--addfolddn-m5a121f.example", IDN2_DISALLOWED, IDN2_NO_TR46
   /* 4-1-2-7 addfold<03D2><00ED>dn.example -> addfold<03C5><00ED>dn.example */
   /* Don't resolve as xn--addfolddn-m5a462f.example */
   /* U+03D2 not valid in IDNA2008 */
   },
  {"\x65\x78\x70\x74\x68\x72\x65\x65\xe1\xbd\x92\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--expthreedn-r8a5844g.example"
   /* 4-1-2-8 expthree<1F52><00ED>dn.example -> expthree<03C5><0313><0300><00ED>dn.example */
   },
  {"\x6e\x6f\x6e\x62\x6d\x70\xf0\x90\x90\x80\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--nonbmpdn-h2a34747d.example", IDN2_DISALLOWED, IDN2_NO_TR46
   /* 4-1-2-9 nonbmp<10400><00ED>dn.example -> nonbmp<10428><00ED>dn.example */
   /* Don't resolve as xn--nonbmpdn-h2a37046d.example */
   /* U+10400 not valid under IDNA2008 */
   },
  {"\x6e\x6f\x6e\x62\x6d\x70\x74\x6f\x61\x73\x63\x69\x69\xf0\x9d\x90\x80\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--nonbmptoasciiadn-msb.example", IDN2_DISALLOWED, IDN2_NO_TR46
   /* 4-1-2-10 nonbmptoascii<1D400><00ED>dn.example -> nonbmptoasciia<00ED>dn.example */
   /* Don't resolve as xn--nonbmptoasciidn-hpb54112i.example */
   /* U+1d400 not valid IDNA2008 */
   },
  {"\x72\x65\x67\x63\x6f\x6d\x62\x65\xcc\x81\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--regcombdn-h4a8b.example", IDN2_NOT_NFC, IDN2_NO_TR46
   /* 4-2-1-1 regcomb<0065><0301><00ED>dn.example -> regcomb<00E9><00ED>dn.example */
   /* Don't resolve as xn--regcombedn-r8a794d.example */
   /* Input not NFC */
   },
  {"regcombéídn.example", "xn--regcombdn-h4a8b.example"
   /* NFKC of previous */
   },
  {"\x63\x6f\x6d\x62\x61\x6e\x64\x63\x61\x73\x65\x45\xcc\x81\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--combandcasedn-lhb4d.example", IDN2_NOT_NFC, IDN2_NO_TR46
   /* 4-2-1-2 combandcase<0045><0301><00ED>dn.example -> combandcase<00E9><00ED>dn.example */
   /* Don't resolve as xn--combandcaseedn-cmb526f.example */
   },
  {"combandcaseÉídn.example",
   "xn--combandcasedn-lhb4d.example", IDN2_DISALLOWED, IDN2_NO_TR46
   /* NFKC of previous, uppercase not IDNA2008-valid */
   },
  {"combandcaseéídn.example",
   "xn--combandcasedn-lhb4d.example"
   /* Lower case of previous */
   },
  {"\x61\x64\x6a\x63\x6f\x6d\x62\xc2\xba\xcc\x81\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--adjcombdn-m5a9d.example", IDN2_DISALLOWED, IDN2_NO_TR46
   /* 4-2-1-3 adjcomb<00BA><0301><00ED>dn.example -> adjcomb<00F3><00ED>dn.example */
   /* Don't resolve as xn--adjcombdn-1qa57cp3r.example */
   /* U+00BA not IDNA2008-valid */
   },
  {"\x65\x78\x74\x63\x6f\x6d\x62\x6f\x63\xcc\x81\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--extcombodn-r8a52a.example", IDN2_NOT_NFC, IDN2_NO_TR46
   /* 4-2-1-4 extcombo<0063><0301><00ED>dn.example -> extcombo<0107><00ED>dn.example */
   /* Don't resolve as xn--extcombocdn-wcb920e.example */
   },
  {"extcomboćídn.example",
   "xn--extcombodn-r8a52a.example"
   /* NFKC of previous */
   },
  {"\x64\x6f\x75\x62\x6c\x65\x64\x69\x61\x63\x31\x75\xcc\x88\xcc\x81\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--doublediac1dn-6ib836a.example", IDN2_NOT_NFC, IDN2_NO_TR46
   /* 4-2-1-5 doublediac1<0075><0308><0301><00ED>dn.example -> doublediac2<01D8><00ED>dn.example */
   /* Don't resolve as xn--doublediac1udn-cmb526fnd.example */
   },
  {"doublediac1ǘídn.example",
   "xn--doublediac1dn-6ib836a.example"
   /* NFKC of previous */
   },
  {"\x64\x6f\x75\x62\x6c\x65\x64\x69\x61\x63\x32\x75\xcc\x81\xcc\x88\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--doublediac2dn-6ib8qs73a.example", IDN2_NOT_NFC, IDN2_NO_TR46
   /* 4-2-1-6 doublediac2<0075><0301><0308><00ED>dn.example -> doublediac2<01D8><00ED>dn.example */
   /* Don't resolve as xn--doublediac2udn-cmb526fod.example */
   },
  {"doublediac2ú̈ídn.example",
   "xn--doublediac2dn-6ib8qs73a.example"
   /* 4-2-1-6 doublediac2<0075><0301><0308><00ED>dn.example -> doublediac2<01D8><00ED>dn.example */
   /* Don't resolve as xn--doublediac2udn-cmb526fod.example */
   },
  {"\x6e\x65\x77\x6e\x6f\x72\x6d\xf0\xaf\xa1\xb4\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--newnormdn-m5a7856x.example", IDN2_NOT_NFC, IDN2_NO_TR46
   /* 4-2-2-1 newnorm<2F874><00ED>dn.example -> newnorm<5F33><00ED>dn.example should not become <5F53> */
   /* Don't resolve as xn--newnormdn-m5a9396x.example */
   /* Don't resolve as xn--newnormdn-m5a9968x.example */
   /* U+2f876 not IDNA2008-valid */
   },
  {"\xe2\x80\x80\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_NOT_NFC, IDN2_NO_TR46
   /* 4-3-1 Spacing */
   /* Don't resolve as xn--dn-mja3392a.example */
   },
  {" ídn.example", "", IDN2_DISALLOWED, IDN2_NO_TR46
   /* NFKC of previous.  U+0020 */
   },
  {"\xdb\x9d\xc3\xad\x64\x6e\x2d\x32\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_DISALLOWED
   /* 4-3-2 Control */
   /* Don't resolve as xn--dn-2-upa332g.example */
   },
  {"\xee\x80\x85\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_DISALLOWED
   /* 4-3-3 Private use */
   /* Don't resolve as xn--dn-mja1659t.example */
   },
  {"\xf3\xb0\x80\x85\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_DISALLOWED
   /* 4-3-4 Private use, non-BMP */
   /* Don't resolve as xn--dn-mja7922x.example */
   },
  {"\xef\xb7\x9d\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_DISALLOWED
   /* 4-3-5 Non-character */
   /* Don't resolve as xn--dn-mja1210x.example */
   },
  {"\xf0\x9f\xbf\xbe\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_DISALLOWED
   /* 4-3-6 Non-character, non-BMP */
   /* Don't resolve as xn--dn-mja7922x.example */
   },
  {"\xef\xbf\xbd\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_DISALLOWED
   /* 4-3-7 Surrogate points */
   /* Don't resolve as xn--dn-mja7922x.example */
   },
  {"\xef\xbf\xba\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_DISALLOWED
   /* 4-3-8 Inappropriate for plain */
   /* Don't resolve as xn--dn-mja5822x.example */
   },
  {"\xe2\xbf\xb5\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_DISALLOWED
   /* 4-3-9 Inappropriate for canonical */
   /* Don't resolve as xn--dn-mja3729b.example */
   },
  {"\xe2\x81\xaa\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_DISALLOWED
   /* 4-3-10 Change display simple */
   /* Don't resolve as xn--dn-mja7533a.example */
   },
  {"\xe2\x80\x8f\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_DISALLOWED
   /* 4-3-11 Change display RTL */
   /* Don't resolve as xn--dn-mja3992a.example */
   },
  {"\xf3\xa0\x80\x81\xf3\xa0\x81\x85\xf3\xa0\x81\x8e\x68\x69\x69\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_DISALLOWED
   /* 4-3-12 Language tags */
   /* Don't resolve as xn--hiidn-km43aaa.example */
   },
  {"\xd8\xa8\x6f\xd8\xb8\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_BIDI
   /* 4-4-1 Arabic RandALCat-LCat-RandALCat */
   /* Don't resolve as xn--o-0mc3c.example */
   },
  {"\xd8\xa8\xd8\xb8\x6f\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_BIDI
   /* 4-4-2 Arabic RandALCat-RandALCat-other */
   /* Don't resolve as xn--o-0mc2c.example */
   },
  {"\x6f\xd8\xa8\xd8\xb8\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_BIDI
   /* 4-4-3 Arabic other-RandALCat-RandALCat */
   /* Don't resolve as xn--o-1mc2c.example */
   },
  {"\xd7\x91\x6f\xd7\xa1\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_BIDI
   /* 4-4-4 Hebrew RandALCat-LCat-RandALCat */
   /* Don't resolve as xn--o-1hc3c.example */
   },
  {"\xd7\x91\xd7\xa1\x6f\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_BIDI
   /* 4-4-5 Hebrew RandALCat-RandALCat-other */
   /* Don't resolve as xn--o-1hc2c.example */
   },
  {"\x6f\xd7\x91\xd7\xa1\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "", IDN2_BIDI
   /* 4-4-6 Hebrew other-RandALCat-RandALCat */
   /* Don't resolve as xn--o-2hc2c.example */
   },
  {"\xc8\xb7\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--dn-mja33k.example"
   /* 5-1-1 Unassigned in BMP; zone editors should reject */
   },
  {"\xf0\x90\x88\x85\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--dn-mja7734x.example", IDN2_UNASSIGNED, IDN2_NO_TR46
   /* 5-1-2 Unassinged outside BMP; zone editors should reject */
   /* Don't resolve as xn--dn-mja7922x.example */
   },
  {"\xc8\xb4\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--dn-mja12k.example"
   /* 5-2-1 Newly assigned in BMP; zone editors should reject */
   },
  {"\xf0\x90\x80\x85\xc3\xad\x64\x6e\x2e\x65\x78\x61\x6d\x70\x6c\x65",
   "xn--dn-mja9232x.example"
   /* 5-2-2 Newly assigned outside of BMP; zone editors should reject */
   /* Don't resolve as xn--dn-mja7922x.example */
   },

  /* Created while writing Libidn2 to trigger certain code paths. */

  {
   "\xe2\x80\x8c", "", IDN2_CONTEXTJ
   /* Standalone contextj U+200C. */
   },
  {
   "foo\xe2\x80\x8c\x65\x65", "", IDN2_CONTEXTJ
   /* Contextj U+200C with surrounding joining characters. */
   },
  {
   "\xdd\x90\xe2\x80\x8c\x65", "", IDN2_CONTEXTJ
   /* Contextj: U+0750 (0xD9 0x84) U+200C e.  U+0750 is D. */
   },
  {
   "\xdd\x90\xe2\x80\x8c\xdd\x90", "xn--3oba901q", IDN2_OK
   /* Contextj: U+0750 (0xDD 0x90) U+200C U+0750.  U+0750 is D. */
   },
  {
   "\xdd\x90\xcc\x80\xe2\x80\x8c\xcc\x80\xcc\x80\xcc\x80\xdd\x90",
   "xn--ksaaaa036cea4345b", IDN2_OK
   /* Contextj: U+750 U+0300 U+200C U+0300 U+0300 U+0300 U+0750.
      U+0300 is T, U+0750 is D. */
   },
  {
   "\xd8\xa8\xe2\x80\x8c\xd8\xa7\xe2\x80\x8c\xd8\xa7", "", IDN2_CONTEXTJ
   /* Contextj: U+0628 U+200C U+0627 U+200C U+0627.
      http://article.gmane.org/gmane.ietf.idnabis/7366 */
   },
  {
   "\xc2\xb7", "xn--uba", IDN2_OK
   /* Contexto: Lookup of U+00B7 should succeed. */
   },
  {
   "\xd0\x94\xd0\xb0\xc1\x80", "", IDN2_ENCODING_ERROR},
  {
   "\xe2\x84\xa6", "", IDN2_NOT_NFC, IDN2_NO_TR46},
#if 0				/* reject or not? */
  {
   "ab--", "", IDN2_2HYPHEN},
#endif
  {
   "--", "--", IDN2_OK, IDN2_NO_TR46},
  {
   "\xcd\x8f", "", IDN2_LEADING_COMBINING, IDN2_NO_TR46
   /* CCC=0 GC=M */
   },
  {
   "\xd2\x88", "", IDN2_LEADING_COMBINING
   /* CCC=0 GC=M */
   },
  {
   "\xcc\x80", "", IDN2_LEADING_COMBINING
   /* CCC!=0 GC=Mn */
   },
  {
   "\xe1\xad\x84", "", IDN2_LEADING_COMBINING
   /* CCC!=0 GC=Mc */
   },
  {
   "", "", IDN2_OK},
  {
   "\xc2\xb8", "", IDN2_DISALLOWED, IDN2_NO_TR46},
  {
   "\xf4\x8f\xbf\xbf", "", IDN2_DISALLOWED},
  {
   "\xe2\x80\x8d", "", IDN2_CONTEXTJ},
  {
   "\xcd\xb8", "", IDN2_UNASSIGNED, IDN2_NO_TR46},
  {
   "\xcd\xb9", "", IDN2_UNASSIGNED, IDN2_NO_TR46},
  {
   "\x72\xc3\xa4\x6b\x73\x6d\xc3\xb6\x72\x67\xc3\xa5\x73",
   "xn--rksmrgs-5wao1o", IDN2_OK},
  {
   "1\xde\x86", "", IDN2_BIDI
   /* Check that bidi rejects leading non-L/R/AL characters in bidi strings */
   },
  {
   "f\xd7\x99", "", IDN2_BIDI
   /* check that ltr string cannot contain R character */
   },
  {
   "-", "-", IDN2_OK, IDN2_NO_TR46},
  {
   "-a", "-a", IDN2_OK, IDN2_NO_TR46},
  {
   "a-", "a-", IDN2_OK, IDN2_NO_TR46},
  {
   "-a", "-a", IDN2_OK, IDN2_NO_TR46},
  {
   "-a-", "-a-", IDN2_OK, IDN2_NO_TR46},
  {
   "foo", "foo", IDN2_OK},
  {"\xe2\x84\xab", "", IDN2_NOT_NFC, IDN2_NO_TR46},
  {"\xe2\x84\xa6", "", IDN2_NOT_NFC, IDN2_NO_TR46},
  {
   /* blåbærgrød.no composed */
   "\x62\x6c\xc3\xa5\x62\xc3\xa6\x72\x67\x72\xc3\xb8\x64\x2e\x6e\x6f",
   "xn--blbrgrd-fxak7p.no", IDN2_OK},
  {
   /* blåbærgrød.no partially decomposed */
   "\x62\x6c\x61\xcc\x8a\x62\xc3\xa6\x72\x67\x72\xc3\xb8\x64\x2e\x6e\x6f", "",
   IDN2_NOT_NFC, IDN2_NO_TR46},
  {
   /* blåbærgrød.no partially decomposed */
   "\x62\x6c\x61\xcc\x8a\x62\xc3\xa6\x72\x67\x72\xc3\xb8\x64\x2e\x6e\x6f",
   "xn--blbrgrd-fxak7p.no", IDN2_OK, IDN2_NFC_INPUT},
  {
   /* blåbærgrød.no partially decomposed */
   "\x62\x6c\x61\xcc\x8a\x62\xc3\xa6\x72\x67\x72\xc3\xb8\x64\x2e\x6e\x6f",
   "xn--blbrgrd-fxak7p.no", IDN2_OK, IDN2_TRANSITIONAL},
  {
   /* blåbærgrød.no partially decomposed */
   "\x62\x6c\x61\xcc\x8a\x62\xc3\xa6\x72\x67\x72\xc3\xb8\x64\x2e\x6e\x6f",
   "xn--blbrgrd-fxak7p.no", IDN2_OK, IDN2_NONTRANSITIONAL},
  {
   /* 0x00FFFFFF, character with 5 bytes UTF-8 representation */
   "\xf8\xbf\xbf\xbf\xbf", "", IDN2_ENCODING_ERROR, IDN2_NONTRANSITIONAL},
  {
   /* bad utf-8 encoding */
   "\x7e\x64\x61\x72\x10\x2f\x2f\xf9\x2b\x71\x60\x79\x7b\x2e\x63\x75\x2b\x61\x65\x72\x75\x65\x56\x66\x7f\x62\xc5\x76\xe5\x00",
   "", IDN2_ENCODING_ERROR, IDN2_NONTRANSITIONAL,
   },
  /* √.com */
  {"\xe2\x88\x9a.com", "xn--19g.com", IDN2_OK, IDN2_TRANSITIONAL},
  /* domains with non-STD3 characters (removed by default when using TR46 transitional/non-transitional */
  {"_443._tcp.example.com", "_443._tcp.example.com", IDN2_OK, IDN2_NO_TR46},
  {"_443._tcp.example.com", "_443._tcp.example.com", IDN2_OK,
   IDN2_TRANSITIONAL},
  {"_443._tcp.example.com", "_443._tcp.example.com", IDN2_OK,
   IDN2_NONTRANSITIONAL},
  {"_443._tcp.example.com", "_443._tcp.example.com", IDN2_OK, 0},
  {"_443._tcp.example.com", "443.tcp.example.com", IDN2_OK,
   IDN2_USE_STD3_ASCII_RULES | IDN2_NONTRANSITIONAL},
  {"_443._tcp.example.com", "443.tcp.example.com", IDN2_OK,
   IDN2_USE_STD3_ASCII_RULES | IDN2_TRANSITIONAL},
  {"_443._tcp.example.com", "_443._tcp.example.com", IDN2_OK, IDN2_USE_STD3_ASCII_RULES | IDN2_NO_TR46},	/* flag is ignored when not using TR46 */
  /* _??? */
  {"_\xc3\xbc", "xn--_-eha", IDN2_DISALLOWED, IDN2_NO_TR46},
  {"_\xc3\xbc", "xn--_-eha", IDN2_OK, IDN2_TRANSITIONAL},
  {"_\xc3\xbc", "xn--_-eha", IDN2_OK, IDN2_NONTRANSITIONAL},
  {"_\xc3\xbc", "xn--tda", IDN2_OK,
   IDN2_USE_STD3_ASCII_RULES | IDN2_NONTRANSITIONAL},
  {"_\xc3\xbc", "xn--tda", IDN2_OK,
   IDN2_USE_STD3_ASCII_RULES | IDN2_TRANSITIONAL},
  {"_\xc3\xbc", "xn--_-eha", IDN2_DISALLOWED, IDN2_USE_STD3_ASCII_RULES | IDN2_NO_TR46},	/* flag is ignored when not using TR46 */
  /* test invalid flags */
  {"_443._tcp.example.com", "_443._tcp.example.com", IDN2_INVALID_FLAGS,
   IDN2_NONTRANSITIONAL | IDN2_TRANSITIONAL},
  {"_443._tcp.example.com", "_443._tcp.example.com", IDN2_INVALID_FLAGS,
   IDN2_NONTRANSITIONAL | IDN2_NO_TR46},
  {"_443._tcp.example.com", "_443._tcp.example.com", IDN2_INVALID_FLAGS,
   IDN2_TRANSITIONAL | IDN2_NO_TR46},
  {"_443._tcp.example.com", "_443._tcp.example.com", IDN2_INVALID_FLAGS,
   IDN2_TRANSITIONAL | IDN2_NONTRANSITIONAL | IDN2_NO_TR46},
  /* underscore and non-ASCII */
  {"\xc3\xa4_x", "xn--_x-uia", IDN2_OK, IDN2_TRANSITIONAL},
  /* failing lookup round-trip */
  {"xn--te_", "", IDN2_ALABEL_ROUNDTRIP_FAILED},
  /* failing lookup round-trip: â˜º -> xn-- o-oia59s (illegal space in output, see https://gitlab.com/libidn/libidn2/issues/78) */
  {"\xc3\xa2\xcb\x9c\xc2\xba", "", IDN2_DISALLOWED, IDN2_NO_TR46},
  {"\xc3\xa2\xcb\x9c\xc2\xba", "", IDN2_ALABEL_ROUNDTRIP_FAILED,
   IDN2_TRANSITIONAL},
  {"\xc3\xa2\xcb\x9c\xc2\xba", "", IDN2_ALABEL_ROUNDTRIP_FAILED,
   IDN2_NONTRANSITIONAL},
  /* long utf-8 input results in good punycode: 퐀쓄쓄쓄쓄쓄쓄쓄쓄쓄쓄쓄쓼쓄쓄쓄쓄쓄쓄쓄쓄쓄㻄쓄쓄럄䄀싂 */
  {"\xec\x80\x90\xed\x93\xec\x84\x93\x84\x93\xec\x84\xec\x84\x93\xec\x93\xec\x84\x93\x84\x93\xec\x84" "\xec\x84\x93\xec\x93\xec\x84\x93\x84\x93\xec\x84\xec\xbc\x93\xec\x93\xec\x84\x93\x84\x93\xec\x84" "\xec\x84\x93\xec\x93\xec\x84\x93\x84\x93\xec\x84\xec\x84\x93\xec\xbb\xe3\x84\x93\x84\x93\xec\x84" "\xeb\x84\x93\xec\x84\xe4\x84\x9f\x82\x8b\xec\x80",
   "xn--p9mx3db62rwgjlncaaaaaaaaaaaaaaaaaaaba41m468u", IDN2_OK,
   IDN2_NONTRANSITIONAL | IDN2_NFC_INPUT},
  /* long utf-8 input results in good punycode:
   * 髦暩晦晦晦獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳筳獳싂.퐀쓄쓄쓄쓄쓄쓄쓄쓄쓄쓄쓄쓼쓄쓄쓄쓄쓄쓄쓄쓄쓄㻄쓄쓄럄䄀싂.뼀猀獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳ⱁ㩁 */
  {"\xe9\xab\xa6\xe6\x9a\xa9\xe6\x99\xa6\xe6\x99\xa6\xe6\x99\xa6\xe7"
   "\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d"
   "\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3"
   "\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7"
   "\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d"
   "\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3"
   "\xe7\xad\xb3\xe7\x8d\xb3\xec\x8b\x82\x2e\xed\x90\x80\xec\x93\x84"
   "\xec\x93\x84\xec\x93\x84\xec\x93\x84\xec\x93\x84\xec\x93\x84\xec"
   "\x93\x84\xec\x93\x84\xec\x93\x84\xec\x93\x84\xec\x93\x84\xec\x93"
   "\xbc\xec\x93\x84\xec\x93\x84\xec\x93\x84\xec\x93\x84\xec\x93\x84"
   "\xec\x93\x84\xec\x93\x84\xec\x93\x84\xec\x93\x84\xe3\xbb\x84\xec"
   "\x93\x84\xec\x93\x84\xeb\x9f\x84\xe4\x84\x80\xec\x8b\x82\x2e\xeb"
   "\xbc\x80\xe7\x8c\x80\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d"
   "\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3"
   "\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7"
   "\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe7\x8d\xb3\xe2\xb1\x81\xe3\xa9\x81",
   "xn--lkvaa9xr87caaaaaaaaaaaaaaaaaaaaaaaaaaa7968dcp2n7tvk.xn--p9mx3db62rwgjlncaaaaaaaaaaaaaaaaaaaba41m468u.xn--bfj606ben8bfnaaaaaaaaaaaaaaaaaa79563b",
   IDN2_OK, IDN2_NONTRANSITIONAL | IDN2_NFC_INPUT},
  {"髦暩晦晦晦獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳筳獳싂.퐀쓄쓄쓄쓄쓄쓄쓄쓄쓄쓄쓄쓼쓄쓄쓄쓄쓄쓄쓄쓄쓄㻄쓄쓄럄䄀싂.뼀猀獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳獳ⱁ㩁",
   "xn--lkvaa9xr87caaaaaaaaaaaaaaaaaaaaaaaaaaa7968dcp2n7tvk.xn--p9mx3db62rwgjlncaaaaaaaaaaaaaaaaaaaba41m468u.xn--bfj606ben8bfnaaaaaaaaaaaaaaaaaa79563b",
   IDN2_OK, IDN2_NONTRANSITIONAL | IDN2_NFC_INPUT},
};

static int ok = 0, failed = 0;
static int break_on_error = 0;

static const char *
_nextField (char **line)
{
  char *s = *line, *e;

  if (!*s)
    return "";

  if (!(e = strpbrk (s, ";#")))
    {
      e = *line += strlen (s);
    }
  else
    {
      *line = e + (*e == ';');
      *e = 0;
    }

  // trim leading and trailing whitespace
  while (isspace (*s))
    s++;
  while (e > s && isspace (e[-1]))
    *--e = 0;

  return s;
}

static int
_scan_file (const char *fname, int (*scan) (char *))
{
  FILE *fp = fopen (fname, "r");
  char *buf = NULL, *linep;
  size_t bufsize = 0;
  ssize_t buflen;
  int ret = 0;

  if (!fp)
    {
      fprintf (stderr, "Failed to open %s (%d)\n", fname, errno);
      return -1;
    }

  while ((buflen = getline (&buf, &bufsize, fp)) >= 0)
    {
      linep = buf;

      while (isspace (*linep))
	linep++;		// ignore leading whitespace

      // strip off \r\n
      while (buflen > 0 && (buf[buflen] == '\n' || buf[buflen] == '\r'))
	buf[--buflen] = 0;

      if (!*linep || *linep == '#')
	continue;		// skip empty lines and comments

      if ((ret = scan (linep)))
	break;
    }

  free (buf);
  fclose (fp);

  return ret;
}

static void
test_homebrewed (void)
{
  uint32_t dummy_u32[4] = { 'a', 'b', 'c', 0 };
  uint8_t *out;
  size_t i;
  int rc;

  for (i = 0; i < sizeof (idna) / sizeof (idna[0]); i++)
    {
      rc = idn2_lookup_u8 ((uint8_t *) idna[i].in, &out, idna[i].flags);
      printf ("%3d  %-25s %-40s %s\n", (int) i, idn2_strerror_name (rc),
	      rc == IDN2_OK ? idna[i].out : "", idna[i].in);

      if (rc != idna[i].rc && rc == IDN2_ENCODING_ERROR)
	{
	  printf ("utc bug\n");
	}
      else if (rc != idna[i].rc && idna[i].rc != -1)
	{
	  failed++;
	  printf ("expected rc %d got rc %d\n", idna[i].rc, rc);
	}
      else if (rc == IDN2_OK && strcmp ((char *) out, idna[i].out) != 0)
	{
	  failed++;
	  printf ("expected: %s\ngot: %s\n", idna[i].out, out);
	}
      else
	ok++;

      if (rc == IDN2_OK)
	idn2_free (out);

      /* Try the IDN2_NO_TR46 flag behavior */
      if (!(idna[i].flags & (IDN2_NONTRANSITIONAL | IDN2_TRANSITIONAL)))
	{
	  rc =
	    idn2_lookup_u8 ((uint8_t *) idna[i].in, &out,
			    idna[i].flags | IDN2_NO_TR46);
	  printf ("%3d  %-25s %-40s %s\n", (int) i, idn2_strerror_name (rc),
		  rc == IDN2_OK ? idna[i].out : "", idna[i].in);

	  if (rc != idna[i].rc && rc == IDN2_ENCODING_ERROR)
	    {
	      printf ("utc bug\n");
	    }
	  else if (rc != idna[i].rc && idna[i].rc != -1)
	    {
	      failed++;
	      printf ("expected rc %d got rc %d\n", idna[i].rc, rc);
	    }
	  else if (rc == IDN2_OK && strcmp ((char *) out, idna[i].out) != 0)
	    {
	      failed++;
	      printf ("expected: %s\ngot: %s\n", idna[i].out, out);
	    }
	  else
	    ok++;

	  if (rc == IDN2_OK)
	    idn2_free (out);
	}

      /* Try whether the default flags behave as NONTRANSITIONAL */
      if (!(idna[i].flags & (IDN2_NO_TR46 | IDN2_TRANSITIONAL)))
	{
	  rc =
	    idn2_lookup_u8 ((uint8_t *) idna[i].in, &out,
			    idna[i].flags | IDN2_NONTRANSITIONAL);
	  printf ("%3d  %-25s %-40s %s\n", (int) i, idn2_strerror_name (rc),
		  rc == IDN2_OK ? idna[i].out : "", idna[i].in);

	  if (rc != idna[i].rc && rc == IDN2_ENCODING_ERROR)
	    {
	      printf ("utc bug\n");
	    }
	  else if (rc != idna[i].rc && idna[i].rc != -1)
	    {
	      failed++;
	      printf ("expected rc %d got rc %d\n", idna[i].rc, rc);
	    }
	  else if (rc == IDN2_OK && strcmp ((char *) out, idna[i].out) != 0)
	    {
	      failed++;
	      printf ("expected: %s\ngot: %s\n", idna[i].out, out);
	    }
	  else
	    ok++;

	  if (rc == IDN2_OK)
	    idn2_free (out);
	}

      if (failed && break_on_error)
	exit (EXIT_FAILURE);
    }

  /* special calls to cover edge cases */
  if ((rc = idn2_lookup_u8 (NULL, NULL, 0)) != IDN2_OK)
    {
      failed++;
      printf ("special #1 failed with %d\n", rc);
    }
  else
    ok++;

  out = (uint8_t *) 0x123;
  if ((rc = idn2_lookup_u8 (NULL, &out, 0)) != IDN2_OK)
    {
      failed++;
      printf ("special #2 failed with %d\n", rc);
    }
  else if (out)
    {
      failed++;
      printf ("special #2 failed with out!=NULL\n");
    }
  else
    ok++;

  if ((rc = idn2_lookup_ul (NULL, NULL, 0)) != IDN2_OK)
    {
      failed++;
      printf ("special #3 failed with %d\n", rc);
    }
  else
    ok++;

  out = (uint8_t *) 0x123;
  if ((rc = idn2_lookup_ul (NULL, (char **) &out, 0)) != IDN2_OK)
    {
      failed++;
      printf ("special #4 failed with %d\n", rc);
    }
  else if (out)
    {
      failed++;
      printf ("special #4 failed with out!=NULL\n");
    }
  else
    ok++;

  if ((rc = idna_to_ascii_8z ("abc", (char **) &out, 0)) != IDN2_OK)
    {
      failed++;
      printf ("special #7 failed with %d\n", rc);
    }
  else
    {
      idn2_free (out);
      ok++;
    }

  if ((rc = idna_to_ascii_4z (dummy_u32, (char **) &out, 0)) != IDN2_OK)
    {
      failed++;
      printf ("special #8 failed with %d\n", rc);
    }
  else
    {
      idn2_free (out);
      ok++;
    }

  if ((rc = idn2_to_ascii_4i2 (dummy_u32, 4, (char **) &out, 0)) != IDN2_OK)
    {
      failed++;
      printf ("special #9 failed with %d\n", rc);
    }
  else
    {
      idn2_free (out);
      ok++;
    }
}

// decode embedded UTF-16/32 sequences
static uint8_t *
_decodeIdnaTest (const uint8_t * src_u8)
{
  size_t it2 = 0, len;
  uint32_t *src;

  // convert UTF-8 to UCS-4 (Unicode))
  if (!(src = u8_to_u32 (src_u8, u8_strlen (src_u8) + 1, NULL, &len)))
    {
      printf ("u8_to_u32(%s) failed (%d)\n", src_u8, errno);
      return NULL;
    }

  // replace escaped UTF-16 incl. surrogates
  for (size_t it = 0; it < len;)
    {
      if (src[it] == '\\' && src[it + 1] == 'u')
	{
	  src[it2] =
	    ((src[it + 2] >=
	      'A' ? src[it + 2] - 'A' + 10 : src[it + 2] - '0') << 12) +
	    ((src[it + 3] >=
	      'A' ? src[it + 3] - 'A' + 10 : src[it + 3] - '0') << 8) +
	    ((src[it + 4] >=
	      'A' ? src[it + 4] - 'A' + 10 : src[it + 4] - '0') << 4) +
	    (src[it + 5] >= 'A' ? src[it + 5] - 'A' + 10 : src[it + 5] - '0');
	  it += 6;

	  if (src[it2] >= 0xD800 && src[it2] <= 0xDBFF)
	    {
	      // high surrogate followed by low surrogate
	      if (src[it] == '\\' && src[it + 1] == 'u')
		{
		  uint32_t low =
		    ((src[it + 2] >=
		      'A' ? src[it + 2] - 'A' + 10 : src[it + 2] -
		      '0') << 12) + ((src[it + 3] >=
				      'A' ? src[it + 3] - 'A' + 10 : src[it +
									 3] -
				      '0') << 8) + ((src[it + 4] >=
						     'A' ? src[it + 4] - 'A' +
						     10 : src[it + 4] -
						     '0') << 4) + (src[it +
								       5] >=
								   'A' ?
								   src[it +
								       5] -
								   'A' +
								   10 : src[it
									    +
									    5]
								   - '0');
		  if (low >= 0xDC00 && low <= 0xDFFF)
		    src[it2] =
		      0x10000 + (src[it2] - 0xD800) * 0x400 + (low - 0xDC00);
		  else
		    printf ("Missing low surrogate\n");
		  it += 6;
		}
	      else
		{
		  it++;
		  printf ("Missing low surrogate\n");
		}
	    }
	  it2++;
	}
      else
	src[it2++] = src[it++];
    }

  // convert UTF-32 to UTF-8
  uint8_t *dst_u8 = u32_to_u8 (src, it2, NULL, &len);
  if (!dst_u8)
    printf ("u32_to_u8(%s) failed (%d)\n", src_u8, errno);

  free (src);
  return dst_u8;
}

static void
_check_toASCII (const char *source, const char *expected, int transitional,
		int expected_toASCII_failure)
{
  int rc;
  char *ace = NULL;

  rc =
    idn2_lookup_u8 ((uint8_t *) source, (uint8_t **) & ace,
		    transitional ? IDN2_TRANSITIONAL : IDN2_NONTRANSITIONAL);

  // printf("n=%d expected=%s t=%d got=%s, expected_failure=%d\n", n, expected, transitional, ace ? ace : "", expected_toASCII_failure);
  if (rc && expected_toASCII_failure)
    {
      printf ("OK\n");
      ok++;
    }
  else if (rc && !transitional && *expected != '[')
    {
      failed++;
      printf ("Failed: _check_toASCII(%s) -> %d (expected 0) %p\n", source,
	      rc, ace);
    }
  else if (rc == 0 && !transitional && *expected != '['
	   && strcmp (expected, ace))
    {
      failed++;
      printf ("Failed: _check_toASCII(%s) -> %s (expected %s) %p\n", source,
	      ace, expected, ace);
    }
  else
    {
      printf ("OK\n");
      ok++;
    }

  if (rc == IDN2_OK)
    idn2_free (ace);
}

#if HAVE_LIBUNISTRING
extern int _libunistring_version;
#endif

static int
test_IdnaTest (char *linep)
{
  char *source;
  const char *type, *toUnicode, *toASCII, *NV8, *org_source;
  int expected_toASCII_failure;

  type = _nextField (&linep);
  org_source = _nextField (&linep);
  toUnicode = _nextField (&linep);
  toASCII = _nextField (&linep);
  NV8 = _nextField (&linep);	// if set, the input should be disallowed for IDNA2008

  // sigh, these Unicode people really mix UTF-8 and UCS-2/4
  // quick and dirty translation of '\uXXXX' found in IdnaTest.txt including surrogate handling
  source = (char *) _decodeIdnaTest ((uint8_t *) org_source);
  if (!source)
    return 0;			// some Unicode sequences can't be encoded into UTF-8, skip them

  if (!*toUnicode)
    toUnicode = source;
  if (!*toASCII)
    toASCII = toUnicode;
  expected_toASCII_failure = NV8 && *NV8;

  printf ("##########%s#%s#%s#%s#%s#\n", type, org_source, toUnicode, toASCII,
	  NV8);

#if HAVE_LIBUNISTRING
  /* 3 tests fail with libunicode <= 0.9.3 - just skip them until we have a newer version installed */
  /* watch out, libunicode changed versioning scheme up from 0.9.4 */
  /* If !HAVE_LIBUNISTRING, we use internal gnulib code which works. */
  if (_libunistring_version <= 9)
    {
      if (!strcmp (toASCII, "xn--8jb.xn--etb875g"))
	{
	  free (source);
	  return 0;
	}
    }
#endif

  if (*type == 'B')
    {
      _check_toASCII (source, toASCII, 1, expected_toASCII_failure);
      _check_toASCII (source, toASCII, 0, expected_toASCII_failure);
    }
  else if (*type == 'T')
    {
      _check_toASCII (source, toASCII, 1, expected_toASCII_failure);
    }
  else if (*type == 'N')
    {
      _check_toASCII (source, toASCII, 0, expected_toASCII_failure);
    }
  else
    {
      printf ("Failed: Unknown type '%s'\n", type);
    }

  free (source);

  if (failed && break_on_error)
    return 1;

  return 0;
}

static void
separator (void)
{
  puts ("-----------------------------------------------------------"
	"-------------------------------------");
}

static void
test_unicode_range (void)
{
  uint32_t i, ucs4[2];
  uint8_t *utf8, *out;
  size_t len;
  int rc;

  /* Unicode range is 0-0x10FFFF, go a bit further */
  for (i = 0; i < 0x11FFFF; i++)
    {
      ucs4[0] = i;
      ucs4[1] = 0;

      utf8 = u32_to_u8 (ucs4, 2, NULL, &len);

      rc = idn2_lookup_u8 (utf8, &out, 0);
      if (rc == IDN2_OK)
	idn2_free (out);

      rc = idn2_lookup_u8 (utf8, &out, IDN2_NFC_INPUT);
      if (rc == IDN2_OK)
	idn2_free (out);

      rc = idn2_lookup_u8 (utf8, &out, IDN2_TRANSITIONAL);
      if (rc == IDN2_OK)
	idn2_free (out);

      rc = idn2_lookup_u8 (utf8, &out, IDN2_NONTRANSITIONAL);
      if (rc == IDN2_OK)
	idn2_free (out);

      rc = idn2_lookup_u8 (utf8, &out, IDN2_NO_TR46);
      if (rc == IDN2_OK)
	idn2_free (out);

      free (utf8);
    }
}

int
main (int argc, const char *argv[])
{
  separator ();
  puts ("                                          IDNA2008 Lookup\n");
  puts ("  #  Result                    ACE output                  "
	"             Unicode input");
  separator ();

  test_homebrewed ();

  separator ();

  // test all IDNA cases from Unicode 9.0.0
  if (_scan_file
      (argc == 1 ? SRCDIR "/IdnaTest.txt" : argv[1], test_IdnaTest))
    return EXIT_FAILURE;

  separator ();

  test_unicode_range ();

  if (failed)
    {
      printf ("Summary: %d out of %d tests failed\n", failed, ok + failed);
      return EXIT_FAILURE;
    }

  printf ("Summary: All %d tests passed\n", ok + failed);

  return EXIT_SUCCESS;
}

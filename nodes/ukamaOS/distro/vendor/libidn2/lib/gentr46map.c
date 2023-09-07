/* gentr46map.c - generate TR46 lookup tables
   Copyright (C) 2016-2021 Tim Ruehsen

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

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <ctype.h>

#include "tr46map.h"

/* We don't link this tool with gnulib, work around any config.h
   redefine's from gnulib. */
#undef free

static size_t _u32_stream_len (uint32_t * src, size_t len);

static size_t _u32_cp_stream_len (const uint8_t * stream, size_t ncp);

#define countof(a) (sizeof(a)/sizeof(*(a)))

typedef struct
{
  uint32_t cp1, cp2;
  unsigned nmappings:5,		/* 0-18, # of uint32_t at <offset> */
    offset:14,			/* 0-16383, byte offset into mapdata */
    flag_index:3;
  uint8_t flags;
} IDNAMap_gen;

static IDNAMap_gen idna_map[10000];
static size_t map_pos;

static uint8_t genmapdata[16384];
static size_t mapdata_pos;

static uint8_t flag_combination[8];
static unsigned flag_combinations;

static NFCQCMap nfcqc_map[140];
static size_t nfcqc_pos;

static char *
_nextField (char **line)
{
  char *s = *line, *e;

  if (!*s)
    return NULL;

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
  char buf[1024], *linep;
  ssize_t buflen;
  int ret = 0;

  if (!fp)
    {
      fprintf (stderr, "Failed to open %s\n", fname);
      return -1;
    }

  while (fgets (buf, sizeof (buf), fp))
    {
      linep = buf;
      buflen = strlen (buf);

      // strip off \r\n
      while (buflen > 0 && (buf[buflen] == '\n' || buf[buflen] == '\r'))
	buf[--buflen] = 0;

      while (isspace (*linep))
	linep++;		// ignore leading whitespace

      if (!*linep || *linep == '#')
	continue;		// skip empty lines and comments

      if ((ret = scan (linep)))
	break;
    }

  fclose (fp);

  return ret;
}

static size_t
_u32_stream_len (uint32_t * src, size_t len)
{
  unsigned it;
  size_t n = 0;

/*
1 byte: 0-0x7f -> 0xxxxxxx
2 bytes: 0x80-0x3fff ->1xxxxxxx 0xxxxxxx
3 bytes: 0x4000-0x1fffff ->1xxxxxxx 1xxxxxxx 0xxxxxxx
4 bytes: 0x200000-0xFFFFFFF -> 1xxxxxxx 1xxxxxxx 1xxxxxxx 0xxxxxxx
5 bytes: 0x10000000->0xFFFFFFFF -> 1xxxxxxx 1xxxxxxx 1xxxxxxx 1xxxxxxx
*/
  for (it = 0; it < len; it++)
    {
      uint32_t cp = src[it];

      if (cp <= 0x7f)
	n += 1;
      else if (cp <= 0x3fff)
	n += 2;
      else if (cp <= 0x1fffff)
	n += 3;
      else if (cp <= 0xFFFFFFF)
	n += 4;
      else
	n += 5;
    }

  return n;
}

static size_t
_u32_to_stream (uint8_t * dst, size_t dst_size, uint32_t * src, size_t len)
{
  unsigned it;
  size_t n = 0;

  n = _u32_stream_len (src, len);

  if (!dst)
    return n;

  if (dst_size < n)
    return 0;

  for (it = 0; it < len; it++)
    {
      uint32_t cp = src[it];

      if (cp <= 0x7f)
	*dst++ = cp & 0x7F;
      else if (cp <= 0x3fff)
	{
	  *dst++ = 0x80 | ((cp >> 7) & 0x7F);
	  *dst++ = cp & 0x7F;
	}
      else if (cp <= 0x1fffff)
	{
	  *dst++ = 0x80 | ((cp >> 14) & 0x7F);
	  *dst++ = 0x80 | ((cp >> 7) & 0x7F);
	  *dst++ = cp & 0x7F;
	}
      else if (cp <= 0xFFFFFFF)
	{
	  *dst++ = 0x80 | ((cp >> 21) & 0x7F);
	  *dst++ = 0x80 | ((cp >> 14) & 0x7F);
	  *dst++ = 0x80 | ((cp >> 7) & 0x7F);
	  *dst++ = cp & 0x7F;
	}
      else
	{
	  *dst++ = 0x80 | ((cp >> 28) & 0x7F);
	  *dst++ = 0x80 | ((cp >> 21) & 0x7F);
	  *dst++ = 0x80 | ((cp >> 14) & 0x7F);
	  *dst++ = 0x80 | ((cp >> 7) & 0x7F);
	  *dst++ = cp & 0x7F;
	}
    }

  return n;
}

/* copy 'n' codepoints from stream 'src' to 'dst' */
static void
_copy_from_stream (uint32_t * dst, const uint8_t * src, size_t n)
{
  uint32_t cp = 0;

  for (; n; src++)
    {
      cp = (cp << 7) | (*src & 0x7F);
      if ((*src & 0x80) == 0)
	{
	  *dst++ = cp;
	  cp = 0;
	  n--;
	}
    }
}

static int
read_IdnaMappings (char *linep)
{
  IDNAMap_gen *map = &idna_map[map_pos];
  char *flag, *codepoint, *mapping;
  int n;

  codepoint = _nextField (&linep);
  flag = _nextField (&linep);
  mapping = _nextField (&linep);

  if ((n = sscanf (codepoint, "%X..%X", &map->cp1, &map->cp2)) == 1)
    {
      map->cp2 = map->cp1;
    }
  else if (n != 2)
    {
      printf ("Failed to scan mapping codepoint '%s'\n", codepoint);
      return -1;
    }

  if (map->cp1 > map->cp2)
    {
      printf ("Invalid codepoint range '%s'\n", codepoint);
      return -1;
    }

  if (map_pos && map->cp1 <= idna_map[map_pos - 1].cp2)
    {
      printf ("Mapping codepoints out of order '%s'\n", codepoint);
      return -1;
    }

  if (!strcmp (flag, "valid"))
    map->flags |= TR46_FLG_VALID;
  else if (!strcmp (flag, "mapped"))
    map->flags |= TR46_FLG_MAPPED;
  else if (!strcmp (flag, "disallowed"))
    map->flags |= TR46_FLG_DISALLOWED;
  else if (!strcmp (flag, "ignored"))
    map->flags |= TR46_FLG_IGNORED;
  else if (!strcmp (flag, "deviation"))
    map->flags |= TR46_FLG_DEVIATION;
  else if (!strcmp (flag, "disallowed_STD3_mapped"))
    map->flags |= TR46_FLG_DISALLOWED_STD3_MAPPED;
  else if (!strcmp (flag, "disallowed_STD3_valid"))
    map->flags |= TR46_FLG_DISALLOWED_STD3_VALID;
  else
    {
      printf ("Unknown flag '%s'\n", flag);
      return -1;
    }

  if (mapping && *mapping)
    {
      uint32_t cp, tmp[20] = { 0 }, tmp2[20] = { 0 };
      int pos;

      while ((n = sscanf (mapping, " %X%n", &cp, &pos)) == 1)
	{
	  if (mapdata_pos >= countof (genmapdata))
	    {
	      printf ("genmapdata too small - increase and retry\n");
	      break;
	    }

	  if (map->nmappings == 0)
	    {
	      map->offset = mapdata_pos;
	      if (map->offset != mapdata_pos)
		printf ("offset overflow (%u)\n", (unsigned) mapdata_pos);
	    }

	  tmp[map->nmappings] = cp;
	  mapdata_pos += _u32_to_stream (genmapdata + mapdata_pos, 5, &cp, 1);
	  map->nmappings++;
	  mapping += pos;
	}

      /* selftest */
      _copy_from_stream (tmp2, genmapdata + map->offset, map->nmappings);
      for (pos = 0; pos < map->nmappings; pos++)
	if (tmp[pos] != tmp2[pos])
	  abort ();
    }
  else if (map->flags &
	   (TR46_FLG_MAPPED | TR46_FLG_DISALLOWED_STD3_MAPPED |
	    TR46_FLG_DEVIATION))
    {
      if (map->cp1 != 0x200C && map->cp1 != 0x200D)	/* ZWNJ and ZWJ */
	printf ("Missing mapping for '%s'\n", codepoint);
    }

  if (map_pos && map->nmappings == 0)
    {
      /* merge with previous if possible */
      IDNAMap_gen *prev = &idna_map[map_pos - 1];
      if (prev->cp2 + 1 == map->cp1
	  && prev->nmappings == 0 && prev->flags == map->flags)
	{
	  prev->cp2 = map->cp2;
	  memset (map, 0, sizeof (*map));	/* clean up */
	  return 0;
	}
    }

  if (++map_pos >= countof (idna_map))
    {
      printf ("Internal map size too small\n");
      return -1;
    }

  return 0;
}

static int
_compare_map (IDNAMap_gen * m1, IDNAMap_gen * m2)
{
  if (m1->cp1 < m2->cp1)
    return -1;
  if (m1->cp1 > m2->cp2)
    return 1;
  return 0;
}

static int
read_NFCQC (char *linep)
{
  NFCQCMap *map = &nfcqc_map[nfcqc_pos];
  char *codepoint, *type, *check;
  int n;

  codepoint = _nextField (&linep);
  type = _nextField (&linep);
  check = _nextField (&linep);

  if (!type || strcmp (type, "NFC_QC"))
    return 0;

  if ((n = sscanf (codepoint, "%X..%X", &map->cp1, &map->cp2)) == 1)
    {
      map->cp2 = map->cp1;
    }
  else if (n != 2)
    {
      printf ("Failed to scan mapping codepoint '%s'\n", codepoint);
      return -1;
    }

  if (map->cp1 > map->cp2)
    {
      printf ("Invalid codepoint range '%s'\n", codepoint);
      return -1;
    }

  if (*check == 'N')
    map->check = 1;
  else if (*check == 'M')
    map->check = 2;
  else
    {
      printf ("NFQQC: Unknown value '%s'\n", check);
      return -1;
    }

  if (++nfcqc_pos >= countof (nfcqc_map))
    {
      printf ("Internal NFCQC map size too small\n");
      return -1;
    }

  return 0;
}

static int
_compare_map_by_maplen (IDNAMap_gen * m1, IDNAMap_gen * m2)
{
  if (m1->nmappings != m2->nmappings)
    return m2->nmappings - m1->nmappings;
  if (m1->cp1 < m2->cp1)
    return -1;
  if (m1->cp1 > m2->cp2)
    return 1;
  return 0;
}

/*
static uint32_t *
_u32_memmem(uint32_t *haystack, size_t hlen, uint32_t *needle, size_t nlen)
{
  uint32_t *p;

  if (nlen == 0)
    return haystack;

  for (p = haystack; hlen >= nlen; p++, hlen--)
    {
      if (*p == *needle && (nlen == 1 || u32_cmp(p, needle, nlen) == 0))
	return p;
    }

  return NULL;
}
*/

static uint8_t *
_u8_memmem (uint8_t * haystack, size_t hlen, uint8_t * needle, size_t nlen)
{
  uint8_t *p;

  if (nlen == 0)
    return haystack;

  for (p = haystack; hlen >= nlen; p++, hlen--)
    {
      if (*p == *needle && (nlen == 1 || memcmp (p, needle, nlen) == 0))
	return p;
    }

  return NULL;
}

static size_t
_u32_cp_stream_len (const uint8_t * stream, size_t ncp)
{
  const uint8_t *end;

  for (end = stream; ncp; end++)
    {
      if ((*end & 0x80) == 0)
	ncp--;
    }

  return end - stream;
}

/* Remove doubled mappings. With Unicode 6.3.0 the mapping data shrinks
 *  from 7272 to 4322 entries of uint32_t (29088 to 17288 bytes).
 * Converting those 4322 uin32_t values to a uint8_t stream, we decrease mapping
 *  table size from 17288 to 9153 bytes.
 */
static void
_compact_idna_map (void)
{
  unsigned it;

  /* sort into 'longest mappings first' */
  qsort (idna_map, map_pos, sizeof (IDNAMap_gen),
	 (int (*)(const void *, const void *)) _compare_map_by_maplen);

  uint8_t *data = calloc (sizeof (uint8_t), mapdata_pos), *p;
  size_t ndata = 0, slen;

  if (data == NULL)
    abort ();

  for (it = 0; it < map_pos; it++)
    {
      IDNAMap_gen *map = idna_map + it;

      if (!map->nmappings)
	continue;

      slen = _u32_cp_stream_len (genmapdata + map->offset, map->nmappings);

      if ((p = _u8_memmem (data, ndata, genmapdata + map->offset, slen)))
	{
	  map->offset = p - data;
	  continue;
	}

      memcpy (data + ndata, genmapdata + map->offset, slen);
      map->offset = ndata;
      ndata += slen;
    }

  memcpy (genmapdata, data, ndata);
  mapdata_pos = ndata;
  free (data);

  /* sort into 'lowest codepoint first' */
  qsort (idna_map, map_pos, sizeof (IDNAMap_gen),
	 (int (*)(const void *, const void *)) _compare_map);
}

static void
_combine_idna_flags (void)
{
  unsigned it, it2;

  /* There are not many different combinations of flags */
  for (it = 0; it < map_pos; it++)
    {
      IDNAMap_gen *map = idna_map + it;
      int found = 0;

      for (it2 = 0; it2 < flag_combinations && !found; it2++)
	{
	  if (flag_combination[it2] == map->flags)
	    {
	      map->flag_index = it2;
	      found = 1;
	    }
	}

      if (!found)
	{
	  if (flag_combinations >= countof (flag_combination))
	    {
	      fprintf (stderr,
		       "flag_combination[] too small - increase and retry\n");
	      exit (EXIT_FAILURE);
	    }
	  map->flag_index = flag_combinations++;
	  flag_combination[map->flag_index] = map->flags;
	}
    }
  for (it = 0; it < map_pos; it++)
    {
      IDNAMap_gen *map = idna_map + it;

      if (map->flags != flag_combination[map->flag_index])
	{
	  fprintf (stderr, "Flags do not for 0x%X-0x%X)\n", map->cp1,
		   map->cp2);
	  exit (EXIT_FAILURE);
	}
    }
}

static int
_print_tr46_map (uint32_t min, uint32_t max, int do_print)
{
  unsigned it;
  int it2, entries = 0;

  for (it = 0; it < map_pos; it++)
    {
      const IDNAMap_gen *map = idna_map + it;
      uint32_t cp2, cp1 = map->cp1, value, range;
      int n;

      if (cp1 < min)
	continue;

      if (cp1 > max)
	break;

      n = (map->cp2 - cp1) / 0x10000;

      for (it2 = 0; it2 <= n; it2++, cp1 = cp2 + 1)
	{
	  entries++;

	  if (it2 == n)
	    cp2 = map->cp2;
	  else
	    cp2 = cp1 + 0xFFFF;

	  if (!do_print)
	    continue;

	  range = cp2 - cp1;
	  value =
	    (((map->nmappings << 14) | map->offset) << 3) | map->flag_index;

	  if (max == 0xFF)
	    printf ("0x%X,0x%X,", cp1 & 0xFF, range & 0xFF);
	  else if (max == 0xFFFF)
	    printf ("0x%X,0x%X,0x%X,0x%X,",
		    (cp1 >> 8) & 0xFF, cp1 & 0xFF,
		    (range >> 8) & 0xFF, range & 0xFF);
	  else if (max == 0xFFFFFF)
	    printf ("0x%X,0x%X,0x%X,0x%X,0x%X,",
		    (cp1 >> 16) & 0xFF, (cp1 >> 8) & 0xFF, cp1 & 0xFF,
		    (range >> 8) & 0xFF, range & 0xFF);

	  printf ("0x%X,0x%X,0x%X,\n",
		  (value >> 16) & 0xFF, (value >> 8) & 0xFF, value & 0xFF);
	}
    }

  if (max == 0xFF)
    return entries * 5;
  if (max == 0xFFFF)
    return entries * 7;
  if (max == 0xFFFFFF)
    return entries * 8;

  return 0;
}

int
main (void)
{
  unsigned it;

  // read IDNA mappings
  if (_scan_file (SRCDIR "/IdnaMappingTable.txt", read_IdnaMappings))
    return 1;

  _compact_idna_map ();
  _combine_idna_flags ();

  // read NFC QuickCheck table
  if (_scan_file (SRCDIR "/DerivedNormalizationProps.txt", read_NFCQC))
    return 1;

  qsort (nfcqc_map, nfcqc_pos, sizeof (NFCQCMap),
	 (int (*)(const void *, const void *)) _compare_map);

  printf ("/* This file is automatically generated.  DO NOT EDIT! */\n\n");
  printf ("#include <stdint.h>\n");
  printf ("#include \"tr46map.h\"\n\n");

  printf ("static const uint8_t idna_flags[%u] =\n{", flag_combinations);
  for (it = 0; it < flag_combinations; it++)
    {
      printf ("0x%X,", flag_combination[it]);
    }
  printf ("};\n\n");

  printf ("static const uint8_t idna_map_8[%d] = {\n",
	  _print_tr46_map (0x0, 0xFF, 0));
  _print_tr46_map (0x0, 0xFF, 1);
  printf ("};\n\n");

  printf ("static const uint8_t idna_map_16[%d] = {\n",
	  _print_tr46_map (0x100, 0xFFFF, 0));
  _print_tr46_map (0x100, 0xFFFF, 1);
  printf ("};\n\n");

  printf ("static const uint8_t idna_map_24[%d] = {\n",
	  _print_tr46_map (0x10000, 0xFFFFFF, 0));
  _print_tr46_map (0x10000, 0xFFFFFF, 1);
  printf ("};\n\n");

  printf ("static const uint8_t mapdata[%u] = {\n", (unsigned) mapdata_pos);
  for (it = 0; it < mapdata_pos; it++)
    {
      printf ("0x%02X,%s", genmapdata[it], it % 16 == 15 ? "\n" : "");
    }
  printf ("};\n\n");

  printf ("static const NFCQCMap nfcqc_map[%u] = {\n", (unsigned) nfcqc_pos);
  for (it = 0; it < nfcqc_pos; it++)
    {
      NFCQCMap *map = nfcqc_map + it;
      printf ("{0x%X,0x%X,%d},\n", map->cp1, map->cp2, map->check);
    }
  printf ("};\n");

  return 0;
}

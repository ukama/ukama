#!/bin/sh
# Convenience script for regenerating all autogeneratable files that are
# omitted from the version control repository. In particular, this script
# also regenerates all aclocal.m4, config.h.in, Makefile.in, configure files
# with new versions of autoconf or automake.
#
# This script requires autoconf-2.65..2.71 and automake-1.16.4 in the PATH.
# It also requires
#   - the gperf program.

# Copyright (C) 2003-2022 Free Software Foundation, Inc.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

# Prerequisite (if not used from a released tarball): either
#   - the GNULIB_SRCDIR environment variable pointing to a gnulib checkout, or
#   - a preceding invocation of './autopull.sh'.
#
# Usage: ./autogen.sh [--skip-gnulib]
#
# Options:
#   --skip-gnulib       Avoid fetching files from Gnulib.
#                       This option is useful
#                       - when you are working from a released tarball (possibly
#                         with modifications), or
#                       - as a speedup, if the set of gnulib modules did not
#                         change since the last time you ran this script.

skip_gnulib=false
while :; do
  case "$1" in
    --skip-gnulib) skip_gnulib=true; shift;;
    *) break ;;
  esac
done

TEXINFO_VERSION=6.5

if test $skip_gnulib = false; then

  # texinfo.tex
  # The most recent snapshot of it is available in the gnulib repository.
  # But this is a snapshot, with all possible dangers.
  # A stable release of it is available through "automake --add-missing --copy",
  # but that is too old (does not support @arrow{}). So take the version which
  # matches the latest stable texinfo release.
  if test ! -f build-aux/texinfo.tex; then
    { wget -q --timeout=5 -O build-aux/texinfo.tex.tmp 'https://git.savannah.gnu.org/gitweb/?p=texinfo.git;a=blob_plain;f=doc/texinfo.tex;hb=refs/tags/texinfo-'"$TEXINFO_VERSION" \
        && mv build-aux/texinfo.tex.tmp build-aux/texinfo.tex; \
    } || rm -f build-aux/texinfo.tex.tmp
  fi

  if test -n "$GNULIB_SRCDIR"; then
    test -d "$GNULIB_SRCDIR" || {
      echo "*** GNULIB_SRCDIR is set but does not point to an existing directory." 1>&2
      exit 1
    }
  else
    GNULIB_SRCDIR=`pwd`/gnulib
    test -d "$GNULIB_SRCDIR" || {
      echo "*** Subdirectory 'gnulib' does not yet exist. Use './gitsub.sh pull' to create it, or set the environment variable GNULIB_SRCDIR." 1>&2
      exit 1
    }
  fi
  # Now it should contain a gnulib-tool.
  GNULIB_TOOL="$GNULIB_SRCDIR/gnulib-tool"
  test -f "$GNULIB_TOOL" || {
    echo "*** gnulib-tool not found." 1>&2
    exit 1
  }
  GNULIB_MODULES='
    unitypes
    unistr/base
    unistr/u8-check
    unistr/u8-chr
    unistr/u8-cmp
    unistr/u8-cmp2
    unistr/u8-cpy
    unistr/u8-cpy-alloc
    unistr/u8-endswith
    unistr/u8-mblen
    unistr/u8-mbsnlen
    unistr/u8-mbtouc
    unistr/u8-mbtoucr
    unistr/u8-mbtouc-unsafe
    unistr/u8-move
    unistr/u8-next
    unistr/u8-prev
    unistr/u8-set
    unistr/u8-startswith
    unistr/u8-stpcpy
    unistr/u8-stpncpy
    unistr/u8-strcat
    unistr/u8-strchr
    unistr/u8-strcmp
    unistr/u8-strcoll
    unistr/u8-strcpy
    unistr/u8-strcspn
    unistr/u8-strdup
    unistr/u8-strlen
    unistr/u8-strmblen
    unistr/u8-strmbtouc
    unistr/u8-strncat
    unistr/u8-strncmp
    unistr/u8-strncpy
    unistr/u8-strnlen
    unistr/u8-strpbrk
    unistr/u8-strrchr
    unistr/u8-strspn
    unistr/u8-strstr
    unistr/u8-strtok
    unistr/u8-to-u16
    unistr/u8-to-u32
    unistr/u8-uctomb
    unistr/u16-check
    unistr/u16-chr
    unistr/u16-cmp
    unistr/u16-cmp2
    unistr/u16-cpy
    unistr/u16-cpy-alloc
    unistr/u16-endswith
    unistr/u16-mblen
    unistr/u16-mbsnlen
    unistr/u16-mbtouc
    unistr/u16-mbtoucr
    unistr/u16-mbtouc-unsafe
    unistr/u16-move
    unistr/u16-next
    unistr/u16-prev
    unistr/u16-set
    unistr/u16-startswith
    unistr/u16-stpcpy
    unistr/u16-stpncpy
    unistr/u16-strcat
    unistr/u16-strchr
    unistr/u16-strcmp
    unistr/u16-strcoll
    unistr/u16-strcpy
    unistr/u16-strcspn
    unistr/u16-strdup
    unistr/u16-strlen
    unistr/u16-strmblen
    unistr/u16-strmbtouc
    unistr/u16-strncat
    unistr/u16-strncmp
    unistr/u16-strncpy
    unistr/u16-strnlen
    unistr/u16-strpbrk
    unistr/u16-strrchr
    unistr/u16-strspn
    unistr/u16-strstr
    unistr/u16-strtok
    unistr/u16-to-u32
    unistr/u16-to-u8
    unistr/u16-uctomb
    unistr/u32-check
    unistr/u32-chr
    unistr/u32-cmp
    unistr/u32-cmp2
    unistr/u32-cpy
    unistr/u32-cpy-alloc
    unistr/u32-endswith
    unistr/u32-mblen
    unistr/u32-mbsnlen
    unistr/u32-mbtouc
    unistr/u32-mbtoucr
    unistr/u32-mbtouc-unsafe
    unistr/u32-move
    unistr/u32-next
    unistr/u32-prev
    unistr/u32-set
    unistr/u32-startswith
    unistr/u32-stpcpy
    unistr/u32-stpncpy
    unistr/u32-strcat
    unistr/u32-strchr
    unistr/u32-strcmp
    unistr/u32-strcoll
    unistr/u32-strcpy
    unistr/u32-strcspn
    unistr/u32-strdup
    unistr/u32-strlen
    unistr/u32-strmblen
    unistr/u32-strmbtouc
    unistr/u32-strncat
    unistr/u32-strncmp
    unistr/u32-strncpy
    unistr/u32-strnlen
    unistr/u32-strpbrk
    unistr/u32-strrchr
    unistr/u32-strspn
    unistr/u32-strstr
    unistr/u32-strtok
    unistr/u32-to-u16
    unistr/u32-to-u8
    unistr/u32-uctomb
    uniconv/base
    uniconv/u8-conv-from-enc
    uniconv/u8-conv-to-enc
    uniconv/u8-strconv-from-enc
    uniconv/u8-strconv-from-locale
    uniconv/u8-strconv-to-enc
    uniconv/u8-strconv-to-locale
    uniconv/u16-conv-from-enc
    uniconv/u16-conv-to-enc
    uniconv/u16-strconv-from-enc
    uniconv/u16-strconv-from-locale
    uniconv/u16-strconv-to-enc
    uniconv/u16-strconv-to-locale
    uniconv/u32-conv-from-enc
    uniconv/u32-conv-to-enc
    uniconv/u32-strconv-from-enc
    uniconv/u32-strconv-from-locale
    uniconv/u32-strconv-to-enc
    uniconv/u32-strconv-to-locale
    unistdio/base
    unistdio/u8-asnprintf
    unistdio/u8-asprintf
    unistdio/u8-snprintf
    unistdio/u8-sprintf
    unistdio/u8-u8-asnprintf
    unistdio/u8-u8-asprintf
    unistdio/u8-u8-snprintf
    unistdio/u8-u8-sprintf
    unistdio/u8-u8-vasnprintf
    unistdio/u8-u8-vasprintf
    unistdio/u8-u8-vsnprintf
    unistdio/u8-u8-vsprintf
    unistdio/u8-vasnprintf
    unistdio/u8-vasprintf
    unistdio/u8-vsnprintf
    unistdio/u8-vsprintf
    unistdio/u16-asnprintf
    unistdio/u16-asprintf
    unistdio/u16-snprintf
    unistdio/u16-sprintf
    unistdio/u16-u16-asnprintf
    unistdio/u16-u16-asprintf
    unistdio/u16-u16-snprintf
    unistdio/u16-u16-sprintf
    unistdio/u16-u16-vasnprintf
    unistdio/u16-u16-vasprintf
    unistdio/u16-u16-vsnprintf
    unistdio/u16-u16-vsprintf
    unistdio/u16-vasnprintf
    unistdio/u16-vasprintf
    unistdio/u16-vsnprintf
    unistdio/u16-vsprintf
    unistdio/u32-asnprintf
    unistdio/u32-asprintf
    unistdio/u32-snprintf
    unistdio/u32-sprintf
    unistdio/u32-u32-asnprintf
    unistdio/u32-u32-asprintf
    unistdio/u32-u32-snprintf
    unistdio/u32-u32-sprintf
    unistdio/u32-u32-vasnprintf
    unistdio/u32-u32-vasprintf
    unistdio/u32-u32-vsnprintf
    unistdio/u32-u32-vsprintf
    unistdio/u32-vasnprintf
    unistdio/u32-vasprintf
    unistdio/u32-vsnprintf
    unistdio/u32-vsprintf
    unistdio/ulc-asnprintf
    unistdio/ulc-asprintf
    unistdio/ulc-fprintf
    unistdio/ulc-snprintf
    unistdio/ulc-sprintf
    unistdio/ulc-vasnprintf
    unistdio/ulc-vasprintf
    unistdio/ulc-vfprintf
    unistdio/ulc-vsnprintf
    unistdio/ulc-vsprintf
    uniname/base
    uniname/uniname
    unictype/base
    unictype/bidiclass-all
    unictype/block-all
    unictype/category-all
    unictype/combining-class-all
    unictype/ctype-alnum
    unictype/ctype-alpha
    unictype/ctype-blank
    unictype/ctype-cntrl
    unictype/ctype-digit
    unictype/ctype-graph
    unictype/ctype-lower
    unictype/ctype-print
    unictype/ctype-punct
    unictype/ctype-space
    unictype/ctype-upper
    unictype/ctype-xdigit
    unictype/decimal-digit
    unictype/digit
    unictype/joininggroup-all
    unictype/joiningtype-all
    unictype/mirror
    unictype/numeric
    unictype/property-all
    unictype/scripts-all
    unictype/syntax-c-ident
    unictype/syntax-c-whitespace
    unictype/syntax-java-ident
    unictype/syntax-java-whitespace
    uniwidth/base
    uniwidth/u8-strwidth
    uniwidth/u8-width
    uniwidth/u16-strwidth
    uniwidth/u16-width
    uniwidth/u32-strwidth
    uniwidth/u32-width
    uniwidth/width
    unigbrk/base
    unigbrk/u8-grapheme-breaks
    unigbrk/u8-grapheme-next
    unigbrk/u8-grapheme-prev
    unigbrk/u16-grapheme-breaks
    unigbrk/u16-grapheme-next
    unigbrk/u16-grapheme-prev
    unigbrk/u32-grapheme-breaks
    unigbrk/u32-grapheme-next
    unigbrk/u32-grapheme-prev
    unigbrk/uc-gbrk-prop
    unigbrk/uc-is-grapheme-break
    unigbrk/ulc-grapheme-breaks
    unigbrk/uc-grapheme-breaks
    uniwbrk/base
    uniwbrk/u8-wordbreaks
    uniwbrk/u16-wordbreaks
    uniwbrk/u32-wordbreaks
    uniwbrk/ulc-wordbreaks
    uniwbrk/wordbreak-property
    unilbrk/base
    unilbrk/u8-possible-linebreaks
    unilbrk/u8-width-linebreaks
    unilbrk/u16-possible-linebreaks
    unilbrk/u16-width-linebreaks
    unilbrk/u32-possible-linebreaks
    unilbrk/u32-width-linebreaks
    unilbrk/ulc-possible-linebreaks
    unilbrk/ulc-width-linebreaks
    uninorm/base
    uninorm/canonical-decomposition
    uninorm/composition
    uninorm/decomposition
    uninorm/filter
    uninorm/nfc
    uninorm/nfd
    uninorm/nfkc
    uninorm/nfkd
    uninorm/u8-normalize
    uninorm/u8-normcmp
    uninorm/u8-normcoll
    uninorm/u8-normxfrm
    uninorm/u16-normalize
    uninorm/u16-normcmp
    uninorm/u16-normcoll
    uninorm/u16-normxfrm
    uninorm/u32-normalize
    uninorm/u32-normcmp
    uninorm/u32-normcoll
    uninorm/u32-normxfrm
    unicase/base
    unicase/empty-prefix-context
    unicase/empty-suffix-context
    unicase/locale-language
    unicase/tolower
    unicase/totitle
    unicase/toupper
    unicase/u8-casecmp
    unicase/u8-casecoll
    unicase/u8-casefold
    unicase/u8-casexfrm
    unicase/u8-ct-casefold
    unicase/u8-ct-tolower
    unicase/u8-ct-totitle
    unicase/u8-ct-toupper
    unicase/u8-is-cased
    unicase/u8-is-casefolded
    unicase/u8-is-lowercase
    unicase/u8-is-titlecase
    unicase/u8-is-uppercase
    unicase/u8-prefix-context
    unicase/u8-suffix-context
    unicase/u8-tolower
    unicase/u8-totitle
    unicase/u8-toupper
    unicase/u16-casecmp
    unicase/u16-casecoll
    unicase/u16-casefold
    unicase/u16-casexfrm
    unicase/u16-ct-casefold
    unicase/u16-ct-tolower
    unicase/u16-ct-totitle
    unicase/u16-ct-toupper
    unicase/u16-is-cased
    unicase/u16-is-casefolded
    unicase/u16-is-lowercase
    unicase/u16-is-titlecase
    unicase/u16-is-uppercase
    unicase/u16-prefix-context
    unicase/u16-suffix-context
    unicase/u16-tolower
    unicase/u16-totitle
    unicase/u16-toupper
    unicase/u32-casecmp
    unicase/u32-casecoll
    unicase/u32-casefold
    unicase/u32-casexfrm
    unicase/u32-ct-casefold
    unicase/u32-ct-tolower
    unicase/u32-ct-totitle
    unicase/u32-ct-toupper
    unicase/u32-is-cased
    unicase/u32-is-casefolded
    unicase/u32-is-lowercase
    unicase/u32-is-titlecase
    unicase/u32-is-uppercase
    unicase/u32-prefix-context
    unicase/u32-suffix-context
    unicase/u32-tolower
    unicase/u32-totitle
    unicase/u32-toupper
    unicase/ulc-casecmp
    unicase/ulc-casecoll
    unicase/ulc-casexfrm
    relocatable-lib-lgpl
  '
  $GNULIB_TOOL --lib=libunistring --source-base=lib --m4-base=gnulib-m4 --tests-base=tests \
    --with-tests --lgpl=3orGPLv2 --makefile-name=Makefile.gnulib --libtool --local-dir=gnulib-local \
    --import $GNULIB_MODULES
  # Change lib/unistr.h to be usable standalone.
  sed -e 's/if GNULIB_[A-Za-z0-9_]* || .*/if 1/g' \
      -e 's/if GNULIB_[A-Za-z0-9_]*/if 1/g' \
      -e 's/HAVE_INLINE/UNISTRING_HAVE_INLINE/g' \
      < lib/unistr.in.h \
      > lib/unistr.in.h.tmp \
  && mv lib/unistr.in.h.tmp lib/unistr.in.h
  # Change lib/unictype.h, lib/uninorm.h, lib/unicase.h for shared libraries on Woe32 systems.
  sed -e 's/extern const uc_general_category_t UC_/extern LIBUNISTRING_DLL_VARIABLE const uc_general_category_t UC_/' \
      -e 's/extern const uc_property_t UC_/extern LIBUNISTRING_DLL_VARIABLE const uc_property_t UC_/' \
      < lib/unictype.in.h \
      > lib/unictype.in.h.tmp \
  && mv lib/unictype.in.h.tmp lib/unictype.in.h
  sed -e 's/extern const struct unicode_normalization_form /extern LIBUNISTRING_DLL_VARIABLE const struct unicode_normalization_form /' \
      < lib/uninorm.in.h \
      > lib/uninorm.in.h.tmp \
  && mv lib/uninorm.in.h.tmp lib/uninorm.in.h
  sed -e 's/extern const casing_/extern LIBUNISTRING_DLL_VARIABLE const casing_/' \
      < lib/unicase.in.h \
      > lib/unicase.in.h.tmp \
  && mv lib/unicase.in.h.tmp lib/unicase.in.h
  $GNULIB_TOOL --copy-file build-aux/ar-lib; chmod a+x build-aux/ar-lib
  $GNULIB_TOOL --copy-file build-aux/config.guess; chmod a+x build-aux/config.guess
  $GNULIB_TOOL --copy-file build-aux/config.sub;   chmod a+x build-aux/config.sub
  $GNULIB_TOOL --copy-file build-aux/declared.sh lib/declared.sh; chmod a+x lib/declared.sh
  $GNULIB_TOOL --copy-file build-aux/run-test; chmod a+x build-aux/run-test
  $GNULIB_TOOL --copy-file build-aux/test-driver.diff
  # If we got no texinfo.tex so far, take the snapshot from gnulib.
  if test ! -f build-aux/texinfo.tex; then
    $GNULIB_TOOL --copy-file build-aux/texinfo.tex
  fi
fi

aclocal -I m4 -I gnulib-m4
autoconf
autoheader && touch config.h.in
# Make sure we get new versions of files brought in by automake.
(cd build-aux && rm -f ar-lib compile depcomp install-sh mdate-sh missing test-driver)
automake --add-missing --copy
patch build-aux/test-driver < build-aux/test-driver.diff
# Get rid of autom4te.cache directory.
rm -rf autom4te.cache

echo "$0: done.  Now you can run './configure'."

## Makefile for the toplevel directory of GNU libunistring.
## Copyright (C) 2009, 2016, 2018 Free Software Foundation, Inc.
##
## This program is free software: you can redistribute it and/or modify
## it under the terms of the GNU General Public License as published by
## the Free Software Foundation; either version 3 of the License, or
## (at your option) any later version.
##
## This program is distributed in the hope that it will be useful,
## but WITHOUT ANY WARRANTY; without even the implied warranty of
## MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
## GNU General Public License for more details.
##
## You should have received a copy of the GNU General Public License
## along with this program.  If not, see <http://www.gnu.org/licenses/>.

## Process this file with automake to produce Makefile.in.

AUTOMAKE_OPTIONS = 1.5 gnu no-dependencies
ACLOCAL = build-aux/fixaclocal @ACLOCAL@
ACLOCAL_AMFLAGS = -I m4 -I gnulib-m4

SUBDIRS = doc gnulib-local lib tests

EXTRA_DIST = \
  version.sh BUGS DEPENDENCIES HACKING autogen.sh \
  build-aux/fixaclocal build-aux/texi2html build-aux/run-test


# Allow users to use "gnulib-tool --update".
EXTRA_DIST += gnulib-m4/gnulib-cache.m4


# Cygwin support.

EXTRA_DIST += woe32dll/export.h


# Windows support.

EXTRA_DIST += INSTALL.windows build-aux/windres-options

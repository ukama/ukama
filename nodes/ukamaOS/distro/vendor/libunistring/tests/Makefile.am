## Makefile for the tests subdirectory of GNU libunistring.
## Copyright (C) 2009 Free Software Foundation, Inc.
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

# Rules generated and collected by gnulib-tool.
include Makefile.gnulib

AUTOMAKE_OPTIONS += 1.11 subdir-objects color-tests

# Ensure that ../lib/config.h is seen before ../config.h.
DEFAULT_INCLUDES = -I. -I$(srcdir) -I$(top_builddir)/lib -I$(top_builddir)

# The test suite uses the 'localcharset' module.
TESTS_ENVIRONMENT += @LOCALCHARSET_TESTS_ENVIRONMENT@


# For debugging memory leaks and memory allocation bugs.
# You should build with --disable-shared when using valgrind.
CHECKER =
#CHECKER = valgrind --tool=memcheck --suppressions=$(srcdir)/../lib/malloca.valgrind --num-callers=20 --leak-check=yes --leak-resolution=high --show-reachable=yes
#CHECKER = valgrind --tool=massif --format=html --depth=10 --alloc-fn=xmalloc --alloc-fn=xrealloc --stacks=no
CHECKER_END_OF_COMMENTS =

# This must be the last thing that gets added to TESTS_ENVIRONMENT.
TESTS_ENVIRONMENT += $(SHELL) $(top_srcdir)/build-aux/run-test '$(CHECKER)'

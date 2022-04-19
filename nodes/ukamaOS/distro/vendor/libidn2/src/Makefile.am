# Copyright (C) 2011-2021 Simon Josefsson

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

AM_CPPFLAGS = -I$(top_builddir)/gl -I$(top_srcdir)/gl
AM_CPPFLAGS += -I$(top_srcdir) -I$(top_builddir) -I$(top_builddir)/lib -I$(srcdir)
AM_CPPFLAGS += -DLOCALEDIR=\"$(localedir)\"
AM_CFLAGS = $(WARN_CFLAGS)

if !HAVE_LIBUNISTRING
AM_CPPFLAGS += -I$(top_srcdir)/unistring/ -I$(top_builddir)/unistring/
endif

bin_PROGRAMS = idn2

idn2_SOURCES = idn2.c blurbs.h
idn2_LDADD = libidn2_cmd.la ../lib/libidn2.la $(LTLIBUNISTRING) $(LIBS)

idn2.c: $(BUILT_SOURCES)

noinst_LTLIBRARIES = libidn2_cmd.la
libidn2_cmd_la_SOURCES =  idn2.ggo idn2_cmd.c idn2_cmd.h
libidn2_cmd_la_LIBADD = ../gl/libgnu.la $(LTLIBUNISTRING)
libidn2_cmd_la_CFLAGS =

if !HAVE_LIBUNISTRING
libidn2_cmd_la_LIBADD += ../unistring/libunistring.la
endif

# pattern rule (%) needed for parallel make (-j)
idn2_cmd%c idn2_cmd%h: idn2.ggo
	gengetopt --unamed-opts --no-handle-version --no-handle-help \
		--set-package="idn2" \
		--input $< --file-name idn2_cmd
	perl -pi -e 's/\[OPTIONS\]/\[OPTION\]/g' idn2_cmd.c
	perl -pi -e 's/\[FILES\]/\[STRING\]/g' idn2_cmd.c

BUILT_SOURCES = idn2_cmd.c idn2_cmd.h
MAINTAINERCLEANFILES = $(BUILT_SOURCES)

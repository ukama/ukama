# Copyright (c) 2017-2021 Tim Ruehsen

# Permission is hereby granted, free of charge, to any person obtaining a
# copy of this software and associated documentation files (the "Software"),
# to deal in the Software without restriction, including without limitation
# the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
# FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
# DEALINGS IN THE SOFTWARE.

AM_CFLAGS = $(WERROR_CFLAGS) $(WARN_CFLAGS)
AM_CPPFLAGS = -I$(top_srcdir)/lib -I../lib -I$(top_srcdir)/gl -I../gl -DSRCDIR=\"$(abs_srcdir)\" -DTEST_RUN
AM_LDFLAGS = -no-install

LDADD = ../lib/libidn2.la ../gl/libgnu.la $(LTLIBUNISTRING)

IDN_TESTS = \
 libidn2_to_ascii_8z_fuzzer$(EXEEXT) \
 libidn2_to_unicode_8z8z_fuzzer$(EXEEXT) \
 libidn2_register_fuzzer$(EXEEXT)

check_PROGRAMS = $(IDN_TESTS)

libidn2_to_ascii_8z_fuzzer_SOURCES = libidn2_to_ascii_8z_fuzzer.c main.c fuzzer.h
libidn2_to_unicode_8z8z_fuzzer_SOURCES = libidn2_to_unicode_8z8z_fuzzer.c main.c fuzzer.h
libidn2_register_fuzzer_SOURCES = libidn2_register_fuzzer.c main.c fuzzer.h

dist-hook:
	find . -name '*.options' -exec cp -v '{}' $(distdir) ';'
	find . -name '*.dict' -exec cp -v '{}' $(distdir) ';'
	find . -name '*.in' -exec cp -vr '{}' $(distdir) ';'
	find . -name '*.repro' -exec cp -vr '{}' $(distdir) ';'

TESTS = $(IDN_TESTS)

clean-local:
	rm -rf *.gc?? *.log lcov *_fuzzer coverage.info *.o

oss-fuzz:
	if test "$$OUT" != ""; then \
	for ccfile in *_fuzzer.c; do \
	fuzzer=$$(basename $$ccfile .c); \
	$$CXX $$CXXFLAGS -I$(top_srcdir)/lib -I$(top_srcdir) \
	"$${fuzzer}.c" -o "$${fuzzer}" \
	../lib/.libs/libidn2.a $${LIB_FUZZING_ENGINE} -Wl,-Bstatic \
	-Wl,-Bdynamic; \
	done; \
	fi

.PHONY: oss-fuzz

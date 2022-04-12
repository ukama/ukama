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

manual_title = Internationalized Domain Names Library

old_NEWS_hash = 71ff7c3bdfcaf6931c1e980d859de147

bootstrap-tools = gnulib,autoconf,automake,libtoolize,make,makeinfo,help2man,gperf,gengetopt,gtkdocize,tar,gzip

local-checks-to-skip += sc_unmarked_diagnostics sc_bindtextdomain # Re-add when we have translation
local-checks-to-skip += sc_immutable_NEWS
local-checks-to-skip += sc_prohibit_strcmp
local-checks-to-skip += sc_copyright_check
local-checks-to-skip += sc_prohibit_gnu_make_extensions

# Ignore gnulib files.
VC_LIST_ALWAYS_EXCLUDE_REGEX = \
  ^(bootstrap|maint.mk|build-aux/gnupload|src/gl/.*|gl/.*|m4/.*|^fuzz/.*.in/.*|^fuzz/.*.repro/.*)$$

# Explicit syntax-check exceptions.
exclude_file_name_regexp--sc_program_name = ^(tests|examples)/.*\.c$$
exclude_file_name_regexp--sc_prohibit_empty_lines_at_EOF = ^doc/reference/version.xml.in$$
exclude_file_name_regexp--sc_trailing_blank = ^(tests/IdnaTest.(txt|inc))|(doc/gdoc)$$
exclude_file_name_regexp--sc_require_config_h = ^examples/.*\.c|lib/gen.*\.c|^fuzz/main.c$$
exclude_file_name_regexp--sc_require_config_h_first = ^examples/.*\.c|lib/gen.*\.c|^fuzz/main.c$$
exclude_file_name_regexp--sc_prohibit_strncpy = ^src/idn2.c$$
exclude_file_name_regexp--sc_GPL_version = ^COPYINGv2|[^/]*.[ch](.in)?$$
exclude_file_name_regexp--sc_prohibit_atoi_atof = ^lib/gentr46map\.c$$
exclude_file_name_regexp--sc_makefile_at_at_check = ^Makefile.am|src/Makefile\.am|tests/Makefile.am$$
exclude_file_name_regexp--sc_prohibit_empty_lines_at_EOF = ^lib/IdnaMappingTable.txt|doc/reference/version.xml.in$$
exclude_file_name_regexp--sc_prohibit_undesirable_word_seq = ^lib/idna-tables-properties.csv$$

my-update-copyright:
	make update-copyright update-copyright-env='UPDATE_COPYRIGHT_HOLDER="Tim Ruehsen" UPDATE_COPYRIGHT_USE_INTERVALS=2'
	make update-copyright update-copyright-env='UPDATE_COPYRIGHT_HOLDER="Simon Josefsson" UPDATE_COPYRIGHT_USE_INTERVALS=2'

aximport:
	for f in m4/ax_*.m4; do \
		wget -O $$f "https://git.savannah.gnu.org/gitweb/?p=autoconf-archive.git;a=blob_plain;f=$$f"; \
	done

INDENT_SOURCES = examples/*.c lib/*.c lib/*.h src/*.c src/*.h tests/*.c

sc_codespell:
	@if `which codespell > /dev/null`; then \
		codespell -L tim,mitre `git ls-files|egrep -v '_fuzzer.in|_fuzzer.repro|gnulib|^lib/.*\.(csv|txt)$$'`; \
	fi

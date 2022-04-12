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

DISTCHECK_CONFIGURE_FLAGS ?= --enable-gtk-doc --enable-gtk-doc-pdf --disable-gcc-warnings

SUBDIRS = gl

if !HAVE_LIBUNISTRING
SUBDIRS += unistring
endif

SUBDIRS += lib src examples tests fuzz po
ACLOCAL_AMFLAGS = -I m4 -I unistring/m4
EXTRA_DIST = m4/gnulib-cache.m4

if ENABLE_DOC
SUBDIRS += doc
endif

EXTRA_DIST += cfg.mk maint.mk CONTRIBUTING.md README.md
EXTRA_DIST += COPYING COPYING.LESSERv3 COPYING.unicode COPYINGv2

check-valgrind:
	$(MAKE) check LOG_COMPILER="valgrind --error-exitcode=301 --leak-check=yes --show-reachable=yes --track-origins=yes --suppressions=suppressions.valgrind"

clean-local:
	rm -rf */*.gc?? */*/*.gc??

pkgconfigdir = $(libdir)/pkgconfig
pkgconfig_DATA = libidn2.pc

include $(top_srcdir)/aminclude_static.am

local-code-coverage-output: code-coverage-capture
	@cat $(PACKAGE_NAME)-$(VERSION)-coverage/index.html|grep headerCovTableEntryHi|head -1|sed 's/^.*>\([0-9]\+\.[0-9]\+\s*%\)<.*$$/ coverage lines: \1/' || true

ABIDW_COMMON = --no-show-locs --no-corpus-path
LIBIDN2_ABI_LAST_FILE = "$(top_srcdir)/contrib/libidn2-latest-$$(uname -m).abi"
ABIGNORE_FILE = "$(top_srcdir)/contrib/libidn2.abignore"
abi-dump-versioned: lib/libidn2.la
	@echo "*************************************************"
	@echo "Generating versioned ABI files of current libidn2"
	@echo "*************************************************"
	@abidw lib/.libs/libidn2.so $(ABIDW_COMMON) --suppressions $(ABIGNORE_FILE) --out-file "$(srcdir)/contrib/libidn2-$(VERSION)-$$(uname -m).abi"

abi-dump-latest: lib/libidn2.la
	@echo "***************************************"
	@echo "Generating ABI files of current libidn2"
	@echo "***************************************"
	@abidw lib/.libs/libidn2.so $(ABIDW_COMMON) --suppressions $(ABIGNORE_FILE) --out-file $(LIBIDN2_ABI_LAST_FILE)
	@rm -f "./contrib/libidn2-latest-$$(uname -m).tmp"

abi-check-latest: lib/libidn2.la
	@echo "Checking whether the latest ABI dump matches"
	@abidiff --suppressions $(ABIGNORE_FILE) lib/.libs/libidn2.so $(LIBIDN2_ABI_LAST_FILE); if test $$? != 0;then \
		echo "*****************************************************"; \
		echo "libidn2 ABI has changed; use 'make abi-dump-latest'  "; \
		echo "and use 'git diff' to check correctness before commit"; \
		echo "*****************************************************"; \
		exit 1; \
	fi
	@echo "********************************"
	@echo "Current release matches ABI dump"
	@echo "********************************"

# Checks the ABI stability between current git version and the previous tagged release.
ABICHECK_COMMON = --no-added-syms
abi-check: lib/libidn2.la
	@for file in $$(eval echo $(srcdir)/contrib/libidn2-*-$$(uname -m).abi);do \
		echo "Comparing libidn2 with $$file"; \
		abidiff $${file} lib/.libs/libidn2.so --suppressions $(ABIGNORE_FILE) $(ABICHECK_COMMON) --hd2 "$(srcdir)/lib/idn2.h"; \
		if test $$? != 0;then \
			echo "****************************************************************************"; \
			echo "ABI check failed; If intentional add suppression in contrib/libidn2.abignore"; \
			echo "****************************************************************************"; \
			exit 1; \
		fi; \
	done
	@echo "********************"
	@echo "ABI checks completed"
	@echo "********************"

# git-version-gen
EXTRA_DIST += $(top_srcdir)/.version
BUILT_SOURCES = $(top_srcdir)/.version
$(top_srcdir)/.version:
	echo $(VERSION) > $@-t && mv $@-t $@
dist-hook: gen-ChangeLog
	echo $(VERSION) > $(distdir)/.tarball-version
	if test -d "$(top_srcdir)/contrib";then \
		$(MAKE) -C $(top_srcdir) abi-check-latest; \
	fi

.PHONY: gen-ChangeLog
gen-ChangeLog:
	$(AM_V_GEN)if test -d .git; then				\
		$(top_srcdir)/build-aux/gitlog-to-changelog >		\
			$(distdir)/cl-t &&				\
		{ rm -f $(distdir)/ChangeLog &&				\
		  mv $(distdir)/cl-t $(distdir)/ChangeLog; }		\
	fi

.PHONY: abi-dump abi-check

LCOV_INFO=libidn2.info
fuzz-coverage: clean
	$(MAKE) -C gl CFLAGS="$(CFLAGS) --coverage" LDFLAGS="$(LDFLAGS) --coverage"
	$(MAKE) -C unistring CFLAGS="$(CFLAGS) --coverage" LDFLAGS="$(LDFLAGS) --coverage"
	$(MAKE) -C lib CFLAGS="$(CFLAGS) --coverage" LDFLAGS="$(LDFLAGS) --coverage"
	$(MAKE) -C fuzz check CFLAGS="$(CFLAGS) --coverage" LDFLAGS="$(LDFLAGS) --coverage"
	lcov --capture --initial --directory lib --directory fuzz --output-file $(LCOV_INFO)
	lcov --capture --directory lib --directory fuzz --output-file $(LCOV_INFO)
#	lcov --remove $(LCOV_INFO) '*/test_linking.c' '*/css_tokenizer.lex' '*/<stdout>' -o $(LCOV_INFO)
	genhtml --prefix . --ignore-errors source $(LCOV_INFO) --legend --title "fuzz coverage" --output-directory=lcov
	@echo
	@echo "You can now view the coverage report with 'xdg-open lcov/index.html'"

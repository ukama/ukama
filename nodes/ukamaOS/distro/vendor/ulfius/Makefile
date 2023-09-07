#
# Ulfius Framework
#
# Makefile used to build all programs
#
# Copyright 2014-2017 Nicolas Mora <mail@babelouest.org>
#
# This program is free software; you can redistribute it and/or
# modify it under the terms of the GNU Lesser General Public License
# as published by the Free Software Foundation;
# version 2.1 of the License.
#
# This library is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.	See the
# GNU GENERAL PUBLIC LICENSE for more details.
#
# You should have received a copy of the GNU General Public
# License along with this library.	If not, see <http://www.gnu.org/licenses/>.
#

LIBULFIUS_LOCATION=./src
EXAMPLES_LOCATION=./example_programs
ifndef UWSCFLAG
UWSC_LOCATION=./tools/uwsc
endif
TESTS_LOCATION=./test

all:
	cd $(LIBULFIUS_LOCATION) && $(MAKE) $*
ifndef UWSCFLAG
	cd $(UWSC_LOCATION) && $(MAKE) $*
endif

debug:
	cd $(LIBULFIUS_LOCATION) && $(MAKE) debug $*
ifndef UWSCFLAG
	cd $(UWSC_LOCATION) && $(MAKE) debug $*
endif

clean:
	cd $(LIBULFIUS_LOCATION) && $(MAKE) clean
	cd $(EXAMPLES_LOCATION) && $(MAKE) clean
ifndef UWSCFLAG
	cd $(UWSC_LOCATION) && $(MAKE) clean
endif
	cd $(TESTS_LOCATION) && $(MAKE) clean
	cd $(TESTS_LOCATION)/autobahn && $(MAKE) clean
	rm -rf doc/html $(TESTS_LOCATION)/cert/server.* $(TESTS_LOCATION)/cert/root* $(TESTS_LOCATION)/cert/client*

examples:
	cd $(EXAMPLES_LOCATION) && $(MAKE) $*

install:
	cd $(LIBULFIUS_LOCATION) && $(MAKE) static-install
ifndef UWSCFLAG
	cd $(UWSC_LOCATION) && $(MAKE) install
endif

uninstall:
	cd $(LIBULFIUS_LOCATION) && $(MAKE) uninstall
ifndef UWSCFLAG
	cd $(UWSC_LOCATION) && $(MAKE) uninstall
endif

check:
	cd $(TESTS_LOCATION) && $(MAKE)

doxygen:
	doxygen doc/doxygen.cfg

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

CURMAKE := $(abspath $(lastword $(MAKEFILE_LIST)))
CURPATH := $(dir $(CURMAKE))
SYSTEMDIR := $(CURPATH)

SYSSUBDIRS = $(shell find ./*/ -maxdepth 1 -type f -name 'Makefile' -printf '%h' | awk -F './' '{ for (i=1; i<=NF; i++) printf $$i" " }')

PB_SUBDIRS = business customer network collector

.PHONY: subdirs $(SYSSUBDIRS) gen $(PB_SUBDIRS:%=gen-%) clean info

subdirs: $(SYSSUBDIRS)

$(SYSSUBDIRS):
	$(MAKE) -C $@

gen: $(PB_SUBDIRS:%=gen-%)

$(PB_SUBDIRS:%=gen-%):
	$(MAKE) -C $(@:gen-%=%) gen

clean:
	@for dir in $(SYSSUBDIRS); do \
		$(MAKE) -C $$dir clean; \
	done

info:
	@echo $(SYSSUBDIRS)

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

CURMAKE := $(abspath $(lastword $(MAKEFILE_LIST)))
CURPATH := $(dir $(CURMAKE))

SUBDIRS := aisgd ctrl emu

.PHONY: all clean $(SUBDIRS)

all: $(SUBDIRS)

$(SUBDIRS):
	$(MAKE) -C $(CURPATH)$@

clean:
	@for dir in $(SUBDIRS); do \
		echo "Cleaning $$dir"; \
		$(MAKE) -C $(CURPATH)$$dir clean; \
	done

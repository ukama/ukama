# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

CURMAKE := $(abspath $(lastword $(MAKEFILE_LIST)))
CURPATH := $(dir $(CURMAKE))
SYSTEMDIR := $(CURRPATH)

SYSSUBDIRS = $(shell find ./*/ -maxdepth 1 -type f -name 'Makefile' -printf '%h'| awk -F './' '{ for (i=1; i<=NF; i++) printf $$i" " }')

.PHONY: subdirs $(SYSSUBDIRS) clean info

subdirs: $(SYSSUBDIRS)
$(SYSSUBDIRS):
	$(MAKE) -C $@

info:
	@echo $(SYSSUBDIRS)

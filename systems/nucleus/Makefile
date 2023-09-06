# Copyright (c) 2023-present, Ukama Inc.
# All rights reserved.

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

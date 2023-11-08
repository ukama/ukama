# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

LIGHT_GREEN='\033[1;32m'
RED = \033[31m
RESET_COLOR='\033[0m'

CURMAKE   := $(abspath $(lastword $(MAKEFILE_LIST)))
CURPATH   := $(dir $(CURMAKE))
SYSTEMDIR := $(CURRPATH)

SYSSUBDIRS := $(shell find ./*/ -maxdepth 1 -type f -name 'Makefile' -printf '%h'| awk -F './' '{ for (i=1; i<=NF; i++) printf $$i" " }')
MOD_DIRS   := $(shell find . -name go.mod -exec dirname {} \;)

LOG_FILE := systems_build.log

.PHONY: report subdirs $(SYSSUBDIRS) clean info mod

report: subdirs
	cat $(LOG_FILE); \
	rm -f $(LOG_FILE); 

subdirs: $(SYSSUBDIRS)

$(SYSSUBDIRS):
	@echo "$(LIGHT_GREEN) Building System: $@ $(RESET_COLOR)"
	cd $@; \
	if $(MAKE); then \
		cd -; \
		echo "System $@ ... $(LIGHT_GREEN)OK $(RESET_COLOR)" >> $(LOG_FILE); \
	else \
		cd -; \
		echo "$(RED) Building System: $@ Failed $(RESET_COLOR)"; \
		echo "System $@ ... $(RED)Failed $(RESET_COLOR)" >> $(LOG_FILE); \
	fi; \

mod:
	@for dir in $(MOD_DIRS); do \
		echo "Running 'go mod tidy' in $${dir}"; \
		(cd $${dir} && go mod tidy); \
	done

info:
	@echo "Systems are: $(SYSSUBDIRS)"

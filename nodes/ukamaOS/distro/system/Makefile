# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2021-present, Ukama Inc.

# Build System apps [Will be replaced by containers]
#ALLDIR := $(sort $(dir $(wildcard */.)))
#FDIR = $(wildcard $(ALLDIR)/Make*)
#FLTRDIR := $(shell find ./*/ -maxdepth 1 -type f -name 'Makefile' -printf '%h\0' | sort -zu | awk -F'./' '{$$1=$$1; print $$4}') 
include ../../config.mk

CUR_MAKE := $(abspath $(lastword $(MAKEFILE_LIST)))
CUR_PATH := $(dir $(CUR_MAKE))

SYS_SUB_DIRS = $(shell find ./*/ -maxdepth 1 -type f -name 'Makefile' | sed 's|/[^/]*$$||' | awk -F './' '{ for (i=1; i<=NF; i++) printf $$i" " }')


.PHONY: subdirs $(SYS_SUB_DIRS) distclean clean info


#Build sub-directories
subdirs: info $(SYS_SUB_DIRS)
$(SYS_SUB_DIRS):
	$(MAKE) -C $@

#Clean sub-directories
clean:
	@echo "Cleaning System build."
	for dir in $(SYS_SUB_DIRS); do \
                $(MAKE) -C $$dir -f Makefile $@; \
        done

#DistClean sub-directories
distclean:
	@echo "DistClean System build."
	for dir in $(SYS_SUB_DIRS); do \
                $(MAKE) -C $$dir -f Makefile $@; \
        done

info:  
	$(info ============ BUILD INFO ============)
	$(info Target Board:   $(TARGET))
	$(info System Dirs:    $(SYS_SUB_DIRS))
	$(info Architecture:   $(ARCH))
	$(info Compiler:       $(CC))
	$(info ====================================)

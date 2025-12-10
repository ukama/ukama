# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2021-present, Ukama Inc.

include config.mk

CUR_MAKE := $(abspath $(lastword $(MAKEFILE_LIST)))
CUR_DIR := $(dir $(CUR_MAKE))

SUB_DIRS := firmware kernel distro

.PHONY: $(SUB_DIRS) clean distclean info

subdirs: info $(SUB_DIRS)

#Build sub-directories
$(SUB_DIRS):
	$(MAKE) -C $@

#Clean
clean:
	@echo Starting cleaning process.
	for dir in $(BUILDSUBDIRS); do \
		$(MAKE) -C $$dir -f Makefile $@; \
	done

#Distclean
distclean:
	@echo Starting distclean process.
	for dir in $(BUILDSUBDIRS); do \
                $(MAKE) -C $$dir -f Makefile $@; \
	done

info:
	@echo "================================="
	@echo " Building Info "
	@echo "---------------------------------"
	@echo " Target Board  : $(TARGET)"
	@echo " Sub Dirs      : $(SUB_DIRS)"
	@echo " Architecture  : $(ARCH)"
	@echo " Compiler      : $(CC)"
	@echo "================================="

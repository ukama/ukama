# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2022-present, Ukama Inc.

include ../config.mk

# Build RootFS for the UkamaDistro
CUR_MAKE := $(abspath $(firstword $(MAKEFILE_LIST)))
CUR_DIR := $(dir $(CUR_MAKE))

#Targets for the makefile
SUB_DIRS := system
PLATFORM := platform
VENDOR := vendor

.PHONY: $(SUB_DIRS) $(PLATFORM) $(VENDOR) clean info

#Build subdirs
subdirs: info 
ifeq ($(ARCH), $(ARCH_ARM))
	@echo "Building in chroot env."	
else
	$(MAKE) -C $(SUB_DIRS)	
endif

$(SUB_DIRS): $(PLATFORM) $(VENDOR)
	$(MAKE) -C $@

$(PLATFORM): $(VENDOR) 
	 $(MAKE) -C $@ XCC=$(CC)

$(VENDOR):
	$(MAKE) -C $@

# Clean
clean:
	@echo "Cleaning Distro."
	for dir in $(SUB_DIRS); do \
		$(MAKE) -C $$dir -f Makefile $@; \
	done
	@echo Cleaning vendor
	$(MAKE) -C $(VENDOR) -f Makefile $@;
	$(MAKE) -C $(PLATFORM) -f Makefile $@;

distclean :
	@echo "DistClean Distro subdirs."
	for dir in $(SUB_DIRS); do \
                $(MAKE) -j$(NPROCS) -C $$dir -f Makefile $@; \
        done
	@echo Cleaning vendor
	$(MAKE) -C $(VENDOR) -f Makefile $@;
	$(MAKE) -C $(PLATFORM) -f Makefile $@;

info:
	$(info ============ BUILD INFO ============)
	$(info Target Board:   $(TARGET))
	$(info Sub Dirs:       $(SUB_DIRS))
	$(info Architecture:   $(ARCH))
	$(info Compiler:       $(CC))
	$(info ====================================)

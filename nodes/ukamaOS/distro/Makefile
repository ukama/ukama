# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2022-present, Ukama Inc.

include ../config.mk

# Build RootFS for the UkamaDistro
DISTROMAKE := $(abspath $(firstword $(MAKEFILE_LIST)))
DISTRODIR := $(dir $(DISTROMAKE))

#Targets for the makefile
DISTROSUBDIRS := system addons utils
PLATFORM := platform


#Final UKAMA ROOTFS path and name
ifndef ROOTFSPATH
ROOTFSDIR := _ukamafs
ROOTFSPATH := $(DISTRODIR)$(ROOTFSDIR)
endif

# Compiler path
#ifneq ($(TARGETBOARD),$(LOCAL))
#$(info	Building musl based gcc for $(TARGETBOARD))
#XGCCPATH := $(DISTRODIR)tools/musl-cross-make/output/bin/
#else
#XGCCPATH = $(shell which gcc | awk 'BEGIN{FS=OFS="/"}{NF--; print}')/
#$(info Skipping CC tools for $(TARGETBOARD))
#endif

.PHONY: subdirs $(DISTROSUBDIRS) $(PLATFORM) vendor clean

#Exporting Varibles
export

#UkamaFS
ukamarfs: subdirs

#Build subdirs
subdirs: $(DISTROSUBDIRS)
$(DISTROSUBDIRS): rootfs $(PLATFORM)
	$(MAKE) -C $@

$(PLATFORM): $(TOOLCHAIN)
	 $(MAKE) -C $@ XCC=$(XGCCPATH)$(XGCC)

#Root FS directory hierarchy
rootfs: vendor
	@echo Creating rootfs directory $(ROOTFSPATH)
	mkdir -p $(ROOTFSPATH)

# vendor
vendor: $(TOOLCHAIN)
	$(MAKE) -C $@

# Clean
clean:
	@echo "Cleaning Distro."
	rm -rf $(ROOTFSPATH)
	for dir in $(DISTROSUBDIRS); do \
		if [ $$dir != $(CROSSCCDIR) ]; then \
			$(MAKE) -C $$dir -f Makefile $@; \
		fi \
	done
	@echo Cleaning vendor
	$(MAKE) -C vendor -f Makefile $@;

distclean :
	@echo "DistClean Distro subdirs."
	rm -rf $(ROOTFSPATH);
	for dir in $(DISTROSUBDIRS); do \
                $(MAKE) -j$(NPROCS) -C $$dir -f Makefile $@; \
        done


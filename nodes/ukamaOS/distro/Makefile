# Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.
include ../config.mk

# Build RootFS for the UkamaDistro
DISTROMAKE := $(abspath $(firstword $(MAKEFILE_LIST)))
DISTRODIR := $(dir $(DISTROMAKE))

#Targets for the makefile
DISTROSUBDIRS := system addons utils
CROSSCCDIR := tools
PLATFORM := platform


#Final UKAMA ROOTFS path and name
ifndef ROOTFSPATH
ROOTFSDIR := _ukamafs
ROOTFSPATH := $(DISTRODIR)$(ROOTFSDIR)
endif

# Compiler path
ifneq ($(TARGETBOARD),$(LOCAL))
$(info	Building musl based gcc for $(TARGETBOARD))
XGCCPATH := $(DISTRODIR)tools/musl-cross-make/output/bin/
TOOLCHAIN := $(CROSSCCDIR)
else
XGCCPATH = $(shell which gcc | awk 'BEGIN{FS=OFS="/"}{NF--; print}')/
$(info Skipping CC tools for $(TARGETBOARD))
endif	


.PHONY: subdirs $(DISTROSUBDIRS) $(PLATFORM) vendor $(TOOLCHAIN) clean

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

$(CROSSCCDIR):
	@echo "Building musl lib based crosscompiler."
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


# Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.

# Build RootFS for the UkamaDistro
DISTROMAKE := $(abspath $(lastword $(MAKEFILE_LIST)))
DISTRODIR := $(dir $(DISTROMAKE))

#Targets for the makefile
CROSSCCDIR := tools
DISTROSUBDIRS := system addons utils

#Final UKAMA ROOTFS path and name
ifndef ROOTFSPATH
ROOTFSDIR := _ukamafs
ROOTFSPATH := $(DISTRODIR)$(ROOTFSDIR)
endif

# XGCC PATH
XGCCPATH := $(DISTRODIR)tools/musl-cross-make/output/bin/

.PHONY: subdirs $(DISTROSUBDIRS) vendor $(CROSSCCDIR) clean

#Exporting Varibles
export

#UkamaFS
ukamarfs: subdirs

#Build subdirs
subdirs: $(DISTROSUBDIRS)
$(DISTROSUBDIRS): rootfs 
	$(MAKE) -C $@

#Root FS directory hierarchy
rootfs: vendor
	@echo Creating rootfs directory $(ROOTFSPATH)
	mkdir -p $(ROOTFSPATH)

# vendor
vendor: $(CROSSCCDIR)
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


# Build RootFS for the UkamaDistro
DISTROMAKE := $(abspath $(lastword $(MAKEFILE_LIST)))
DISTRODIR := $(dir $(DISTROMAKE))

#Targets for the makefile
DISTROSUBDIRS := tools system addons utils  

#Final UKAMA ROOTFS path and name
ifndef ROOTFSPATH
ROOTFSDIR := _ukamafs
ROOTFSPATH := $(DISTRODIR)$(ROOTFSDIR)
endif

# XGCC PATH
XGCCPATH := $(DISTRODIR)tools/musl-cross-make/output/bin/
ifeq (,$(wildcard $(XGCCPATH)))
    $(error XGCCPATH not found at default $(XGCCPATH))
    $(info  Provide XGCCPATH to make .) 
endif

.PHONY: subdirs $(DISTROSUBDIRS) clean

#Exporting Varibles
export

#UkamaFS
ukamarfs: subdirs

#Build subdirs
subdirs: $(DISTROSUBDIRS)
$(DISTROSUBDIRS): rootfs 
	$(MAKE) -C $@

#Root FS directory hierarchy
rootfs:
	@echo Creating rootfs directory $(ROOTFSPATH)
	mkdir -p $(ROOTFSPATH)

# Clean
clean:
	echo "Cleaning Distro."
	rm -rf $(ROOTFSPATH)
	for dir in $(DISTROSUBDIRS); do \
		$(MAKE) -C $$dir -f Makefile $@; \
  	done

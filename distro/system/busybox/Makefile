#Makefile for busybox

include ../../../config.mk

CURMAKE := $(abspath $(firstword $(MAKEFILE_LIST)))
CURPATH := $(dir $(CURMAKE))

#Source directory
SRCDIR = busybox

#Output
ifndef ROOTFSPATH
override ROOTFSPATH = $(CURPATH)_ukamafs
endif

#Check XGCC PATH
ifndef XGCCPATH
$(error XGCCPATH missing.)
endif

override CC = "$(XGCCPATH)/$(XGCC) --static"

#Config for Builds
BBCONFIG = ukama_defconfig

#Overiding  Architecture
ifeq ($(ARCH), $(ARCHX86))
override ARCH = $(ARCHX86_64)
endif

.PHONY: subdirs $(SRCDIR) clean info

subdirs: $(SRCDIR)

$(SRCDIR):
	@echo Building busybox for $(ARCH) using $(CC)
	mkdir -p $(ROOTFSPATH)
	#Config
	$(MAKE) -j$(NPROCS) -C $@ ARCH=$(ARCH) CC=$(CC) $(BBCONFIG)
	#Make
	$(MAKE) -j$(NPROCS) -C $@ ARCH=$(ARCH) CC=$(CC) install CONFIG_PREFIX=$(ROOTFSPATH)

clean:
	@echo Cleaning $(SRCDIR).
	rm -rf $(ROOTFSPATH)
	for dir in $(SRCDIR); do \
                $(MAKE) -j$(NPROCS) -C $$dir -f Makefile $@; \
        done

info: ; 
	$(info [$@] Building $(SRCDIR) $(TARGETBOARD) for $(ARCH) with $(CC) )

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

XCROSSCC = "$(XGCCPATH)$(XCROSS_COMPILER)"

#Config for Builds
BBCONFIG = ${BB_CONFIG}

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
	$(MAKE) -j$(NPROCS) -C $@ ARCH=$(ARCH) CROSS_COMPILE=$(XCROSSCC) $(BBCONFIG)
	#Make
	$(MAKE) -j$(NPROCS) -C $@ ARCH=$(ARCH) CROSS_COMPILE=$(XCROSSCC) install CONFIG_PREFIX=$(ROOTFSPATH)

clean:
	@echo Cleaning $(SRCDIR).
	rm -rf $(ROOTFSPATH)
	for dir in $(SRCDIR); do \
                $(MAKE) -j$(NPROCS) -C $$dir -f Makefile $@; \
        done

distclean:
	@echo DistClean $(SRCDIR).
	rm -rf $(ROOTFSPATH)
	for dir in $(SRCDIR); do \
                $(MAKE) -j$(NPROCS) -C $$dir -f Makefile $@; \
        done


info: ; 
	$(info [$@] Building $(SRCDIR) $(TARGETBOARD) for $(ARCH) with $(CC) )

#Makefile for UkamaDM

include ../../../config.mk

CURMAKE := $(abspath $(firstword $(MAKEFILE_LIST)))
CURPATH := $(dir $(CURMAKE))

#Source directory
SRCDIR = dm

#Output
ifndef ROOTFSPATH
override ROOTFSPATH = $(CURPATH)_ukamafs
endif

.PHONY: subdirs $(SRCDIR) clean info

ifndef XGCCPATH
$(error XGCCPATH missing.)
endif

#BUILD FLAGS
BUILDFLAGS := 

COMP = $(XGCCPATH)$(XGCC)
LINKER = $(XGCCPATH)$(XLD)

subdirs: $(SRCDIR)

$(SRCDIR):
	@echo Building $@ for $(ARCH) 
	$(MAKE) -j$(NPROCS) -C $@ CC=$(COMP) LD=$(LINKER)
	(mkdir -p $(ROOTFSPATH)/sbin/ && cp -rf $@/build/bin/* $(ROOTFSPATH)/sbin/)
clean:
	@echo Cleaning $(SRCDIR).
	rm -rf $(ROOTFSPATH)
	for dir in $(SRCDIR); do \
		$(MAKE) -j$(NPROCS) -C $$dir -f Makefile $@; \
        done

info:
	@echo Building ukamaDM for Ukama


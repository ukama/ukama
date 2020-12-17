#Makefile for grub

include ../../../config.mk

CURMAKE := $(abspath $(firstword $(MAKEFILE_LIST)))
CURPATH := $(dir $(CURMAKE))

#Source directory
SRCDIR = grub

#Output
ifndef ROOTFSPATH
override ROOTFSPATH = $(CURPATH)_ukamafs
endif

.PHONY: subdirs $(SRCDIR) clean info msg

ifndef XGCCPATH
$(error XGCCPATH missing.)
endif

#BUILD FLAGS
BUILDFLAGS := --disable-nls --disable-werror

#Arch
ifeq ($(ARCH), x86)
CROSS = "$(XGCCPATH)/$(XGCC) --static"

subdirs: $(SRCDIR)

grub:
	@echo Building $@ for $(ARCH) 
	(cd $@ && CC=$(CROSS) ./configure $(BUILDFLAGS) --prefix=$(ROOTFSPATH))
	$(MAKE) -j$(NPROCS) -C $@ 
	$(MAKE) -j$(NPROCS) -C $@ install

clean:
	@echo Cleaning $(SRCDIR).
	rm -rf $(ROOTFSPATH)
	for dir in $(SRCDIR); do \
		$(MAKE) -j$(NPROCS) -C $$dir -f Makefile $@; \
        done

info:
	@echo Building grub for Ukama

endif

msg:
	@echo Build not required for $(ARCH)


#Makefile for dropbear

include ../../../config.mk

CURMAKE := $(abspath $(firstword $(MAKEFILE_LIST)))
CURPATH := $(dir $(CURMAKE))

#Source
SRCDIR = dropbear

#Output
ifndef ROOTFSPATH
override ROOTFSPATH = $(CURPATH)_ukamafs
endif

ifndef XGCCPATH
$(error XGCCPATH missing.)
endif

#Overiding  Architecture
ifeq ($(ARCH), $(ARCHX86))
override ARCH = $(ARCHX86_64)
endif

override CC = "$(XGCCPATH)/$(XGCC) --static"

#Build flags
BUILDFLAGS = --disable-zlib

.PHONY: subdirs $(SRCDIR) clean info

subdirs: $(SRCDIR)

$(SRCDIR):
	@echo Building Dropbear SSH for $(ARCH) using $(CC)
	mkdir -p $(ROOTFSPATH)
	#Config
	(cd $(SRCDIR) && ./configure --host=$(ARCH) --prefix=/ $(BUILDFLAGS) CC=$(CC))
	$(MAKE) -j$(NPROCS) -C $@ 
	#Make
	$(MAKE) -j$(NPROCS) -C $@ install DESTDIR=$(ROOTFSPATH)

clean:
	@echo Cleaning $(SRCDIR).
	rm -rf $(ROOTFSPATH)
	for dir in $(SRCDIR); do \
                $(MAKE) -j$(NPROCS) -C $$dir -f Makefile $@; \
        done

info: ; 
	$(info [$@] Building $(SRCDIR) $(TARGETBOARD) for $(ARCH) with $(CC) )

#Makefile for microInit

include ../../../config.mk

CURMAKE := $(abspath $(firstword $(MAKEFILE_LIST)))
CURPATH := $(dir $(CURMAKE))

#Source directory
SRCDIR = microInit/init microInit/rcinit microInit/sysinit

#Output
ifndef ROOTFSPATH
override ROOTFSPATH = $(CURPATH)_ukamafs
endif

.PHONY: subdirs $(SRCDIR) clean info msg

#BUILD FLAGS
BUILDFLAGS := --disable-nls --disable-werror

#Arch
ifeq ($(ARCH), $(ARCHARM))
TARGETFLAG := armv7-unknown-linux-musleabihf
ADDITIONALLINKERFLAG := CARGO_TARGET_ARMV7_UNKNOWN_LINUX_MUSLEABIHF_LINKER=arm-linux-gnueabihf-ld
else
TARGETFLAG := x86_64-unknown-linux-musl
ADDITIONALLINKERFLAG :=
endif

subdirs: $(SRCDIR)

$(SRCDIR):
	@echo Building $@ for $(ARCH) 
	(cd $@ && $(ADDITIONALLINKERFLAG) cargo build --release --target=$(TARGETFLAG))

clean:
	@echo Cleaning $(SRCDIR).
	for dir in $(SRCDIR); do \
		(cd $$dir && cargo clean) \
        done

distclean:
	@echo DistClean $(SRCDIR).
	for dir in $(SRCDIR); do \
                (cd $$dir && cargo clean) \
        done

info:
	@echo Building init for Ukama

msg:
	@echo Build not required for $(ARCH)


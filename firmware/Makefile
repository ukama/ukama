#Make firmware for ukama nodes 

include ../config.mk

CURMAKE := $(abspath $(firstword $(MAKEFILE_LIST)))
CURDIR := $(dir $(CURMAKE))

#Targets based on boards
ANODETARGET = at91bootstrap uboot
CNODETARGET = coreboot 
AT91BOOTSTRAPBIN = at91bootstrap.bin
UBOOTBIN = u-boot.bin
COREBOOTBIN = coreboot.rom

#Output
ifndef ROOTFSPATH
override ROOTFSPATH = $(CURPATH)_ukamafs
endif
ROOTFSKPATH=$(ROOTFSPATH)/boot

# Config for Builds
AT91CONFIG = ukama_anode_eMMC_uboot_defconfig
UBOOTCONFIG = sama5d27_ukama_anode_emmc_defconfig
CBCONFIG = config.ukama_comv1_16mb

#Set build parameters based on targets
ifeq ($(ANODEBOARD), $(TARGETBOARD))
SRCDIRS = $(ANODETARGET)
override CC = arm-linux-gnueabihf-
endif

ifeq ($(CNODEBOARD), $(TARGETBOARD))
override ARCH = $(ARCHX86)
SRCDIRS = $(CNODETARGET)
endif


.PHONY: subdirs $(SRCDIRS) info

subdirs: $(SRCDIRS)

#AT91BOOTSTARP
at91bootstrap:
	@echo Building $@
	mkdir -p $(ROOTFSKPATH)/$@
	$(MAKE) -j$(NPROCS) -C $@ ARCH=$(ARCH) CROSS_COMPILE=$(CC) $(AT91CONFIG)
	$(MAKE) -j$(NPROCS) -C $@ ARCH=$(ARCH) CROSS_COMPILE=$(CC)
	@echo Copy $@/binaries/$(AT91BOOTSTRAPBIN) $(ROOTFSKPATH)/$@/$(AT91BOOTSTRAPBIN)
	(cp -v $@/binaries/$(AT91BOOTSTRAPBIN) $(ROOTFSKPATH)/$@/$(AT91BOOTSTRAPBIN))	

#UBOOT
uboot:
	@echo Building $@
	mkdir -p $(ROOTFSKPATH)/$@
	$(MAKE) -s -j$(NPROCS) -C $@ ARCH=$(ARCH) CROSS_COMPILE=$(CC) $(UBOOTCONFIG)
	$(MAKE) -j$(NPROCS) -C $@ ARCH=$(ARCH) CROSS_COMPILE=$(CC)
	@echo Copy $@/$(UBOOTBIN) $(ROOTFSKPATH)/$@/$(UBOOTBIN)
	(cp -v $@/$(UBOOTBIN) $(ROOTFSKPATH)/$@/$(UBOOTBIN) )

#coreboot
coreboot:
	@echo Building $@
	mkdir -p $(ROOTFSKPATH)/$@
	(cd $@ && cp -v configs/config.ukama_comv1_16mb .config)
	$(MAKE) -j$(NPROCS) -C $@ ARCH=$(ARCH) CROSS_COMPILE=$(CC)
	@echo Copy $@/build/$(COREBOOTBIN) $(ROOTFSKPATH)/$@/$(COREBOOTBIN)
	(cp -v $@/build/$(COREBOOTBIN) $(ROOTFSKPATH)/$@/$(COREBOOTBIN))

grub:
	@echo Building $@
	mkdir -p $(ROOTFSKPATH)/$@
	(cd $@ && ./bootstrap && ./configure --prefix=$(ROOTFSKPATH)/$@)
	$(MAKE) -j$(NPROCS) -C $@ ARCH=$(ARCH) CROSS_COMPILE=$(CC)
	$(MAKE) -j$(NPROCS) -C $@ install
	
clean :
	@echo Cleaning firmware build.
	rm -rf $(ROOTFSKPATH)
	for dir in $(SRCDIRS); do \
		$(MAKE) -j$(NPROCS) -C $$dir -f Makefile $@; \
        done

info: ; 
	$(info [$@] Building $(TARGETBOARD) for $(ARCH) with $(CC) )

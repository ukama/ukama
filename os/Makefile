# Makefile for the LinuxOS of Ukama

include ../config.mk 

CURMAKE := $(abspath $(firstword $(MAKEFILE_LIST)))
CURPATH := $(dir $(CURMAKE))

#Source directory
SRCDIR = linux

#Output
ifndef ROOTFSPATH
override ROOTFSPATH = $(CURPATH)_ukamafs
endif
ROOTFSKPATH=$(ROOTFSPATH)/boot

# Config for Builds
ANODE_KCONFIG = ukama_anode_defconfig
CNODE_KCONFIG = comNode_lk_defconfig

#Set build parameters based on targets
ifeq ($(ANODEBOARD), $(TARGETBOARD))
KIMAGE = zImage
override CC = arm-linux-gnueabihf-
override ARCH = arm
KCONFIG = $(ANODE_KCONFIG)
OS_ARTIFACTS = $(KIMAGE) modules dtbs
endif

ifeq ($(CNODEBOARD), $(TARGETBOARD))
KIMAGE = bzImage
override CC = 
override ARCH = x86_64
KCONFIG = $(CNODE_KCONFIG)
OS_ARTIFACTS = $(KIMAGE) modules
endif

ifeq ($(LOCAL), $(TARGETBOARD))
KIMAGE = bzImage
override CC =
override ARCH = x86_64
KCONFIG = $(CNODE_KCONFIG)
OS_ARTIFACTS = $(KIMAGE) modules
endif


#Targets based on boards
.PHONY: subdirs $(SRCDIR) info

#Kernel Image
$(OS_ARTIFACTS):
	@echo Building $@
	$(MAKE) -j$(NPROCS) -C $(SRCDIR) ARCH=$(ARCH) CROSS_COMPILE=$(CC) $(KCONFIG)
	$(MAKE) -j$(NPROCS) -C $(SRCDIR) ARCH=$(ARCH) CROSS_COMPILE=$(CC) $(OS_ARTIFACTS)
	#Copy Kernel Image
	@echo Copying Kernel Image $(SRCDIR)/arch/${ARCH}/boot/$(KIMAGE)
	(mkdir -p $(ROOTFSKPATH) && cp -v $(SRCDIR)/arch/${ARCH}/boot/$(KIMAGE) $(ROOTFSKPATH)/$(KIMAGE))
	#Install Modules
	$(MAKE) -j$(NPROCS) -C $(SRCDIR) ARCH=$(ARCH) CROSS_COMPILE=$(CC) INSTALL_MOD_PATH=$(ROOTFSPATH) modules_install
ifeq ($(ARCH), $(ARCHARM))
	#Install DTBS
	$(MAKE) -j$(NPROCS) -C $(SRCDIR) ARCH=$(ARCH) CROSS_COMPILE=$(CC) INSTALL_DTBS_PATH=$(ROOTFSPATH)/boot dtbs_install
endif

clean :
	rm -rf $(ROOTFSPATH);
	for dir in $(SRCDIR); do \
                $(MAKE) -j$(NPROCS) -C $$dir -f Makefile $@; \
        done

distclean :
	rm -rf $(ROOTFSPATH);
	for dir in $(SRCDIR); do \
                $(MAKE) -j$(NPROCS) -C $$dir -f Makefile $@; \
        done

info:  
	$(info [$@] Building $(TARGETBOARD) for $(ARCH) with $(CC) )

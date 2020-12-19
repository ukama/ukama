# Build System tools

include ../../config.mk

SUBDIRS := musl-cross-make 

CURMAKE := $(abspath $(firstword $(MAKEFILE_LIST)))
CURPATH := $(dir $(CURMAKE))

#Set build parameters based on targets
ifeq ($(ANODEBOARD), $(TARGETBOARD))
override CC = gcc
override XTARGET = arm-linux-musleabihf
override ARCH = $(ARCHARM)
CONFIGF ="GCCCONFIG =\"--with-arch=armv6k --with-fpu=vfpv2\""
endif

ifeq ($(CNODEBOARD), $(TARGETBOARD))
override CC = gcc
override XTARGET = x86_64-linux-musl
override ARCH = $(ARCHX86_64)
endif

#Add a flag with supported architecures so no package which is specfic for X86 get build for arm.

.PHONY: subdirs $(SUBDIRS) clean info

#Export variables
export

#musl-cross-make
musl-cross-make:
	$(MAKE) -j$(NPROCS) -C $@ TARGET=$(XTARGET) $(CONFIG)
	$(MAKE) -j$(NPROCS) -C $@ TARGET=$(XTARGET) install

#Build sub-directories
subdirs: $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -j$(NPROCS) -C $@ TARGET=$(XTARGET)
	$(MAKE) -j$(NPROCS) -C $@ TARGET=$(XTARGET) install


#Clean sub-directories
clean:
	@echo "Cleaning tools."
	rm -rf musl-cross-make/output
	for dir in $(SUBDIRS); do \
                $(MAKE) -C $$dir -f Makefile $@; \
        done

#DistClean sub-directories
distclean:
	@echo "DistClean tools."
	rm -rf musl-cross-make/output
	for dir in $(SUBDIRS); do \
                $(MAKE) -C $$dir -f Makefile $@; \
        done

info: 
	@echo $(SUBDIRS)



# Build System apps [Will be replaced by containers]
TOOLSUBDIRS = musl-cross-make 

NPROCS:=1
OS:=$(shell uname -s)

#CPUS
ifeq ($(OS), Linux)
        NPROCS:=$(shell grep -c ^processor /proc/cpuinfo)
endif

#if variables are not defined
ifndef ANODEBOARD
override ANODEBOARD = anode
endif

ifndef CNODEBOARD
override CNODEBOARD = cnode
endif

ifndef TARGETBOARD
override TARGETBOARD = $(ANODEBOARD)
endif

#Set build parameters based on targets
ifeq ($(ANODEBOARD), $(TARGETBOARD))
override XTARGET = arm-linux-musleabihf
override ARCH = arm
GCCCONFIG = --with-arch=armv6k --with-fpu=vfpv2
FWSUBDIRS = $(ANODETARGET)
endif

ifeq ($(CNODEBOARD), $(TARGETBOARD))
override XTARGET = x86_64-linux-musl
override ARCH = x86_64
FWSUBDIRS = $(CNODETARGET)
endif

#Add a flag with supported architecures so no package which is specfic for X86 get build for arm.

.PHONY: subdirs $(TOOLSUBDIRS) clean info

#Export variables
export

#Build sub-directories
subdirs: $(TOOLSUBDIRS)
$(TOOLSUBDIRS):
	$(MAKE) -j$(NPROCS) -C $@ TARGET=$(XTARGET)
	$(MAKE) -j$(NPROCS) -C $@ TARGET=$(XTARGET) install

#musl-cross-make
musl-cross-make:
	$(MAKE) -j$(NPROCS) -C $@ TARGET=$(XTARGET)
	$(MAKE) -j$(NPROCS) -C $@ TARGET=$(XTARGET) install

#Clean sub-directories
clean:
	for dir in $(TOOLSUBDIRS); do \
                $(MAKE) -C $$dir -f Makefile $@; \
        done
info: 
	@echo $(TOOLSUBDIR)
	#@echo $(ALLDIR)
	#@echo $(FDIR)
	#@echo $(shell find ./*/ -maxdepth 1 -type f -name 'Makefile' -printf '%h\0' | sort -zu | awk -F './' '{ for (i=1; i<=NF; i++) printf $$i" " }')
	#@echo $(shell find . -type f -name 'Makefile' -printf '%h' | sort -zu | sed -z 's/$$/ /')
	#@echo $(shell find ./*/ -type f -name 'Makefile' -printf '%h'| awk -F './' '{ for (i=1; i<=NF; i++) printf $$i" " }')



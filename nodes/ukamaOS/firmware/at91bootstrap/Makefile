#
# Default config file is in board/$(BOARD_NAME)/*_defconfig
# First, run xxx_defconfig
# Then, `make menuconfig' if needed
#

# Do not:
# o  use make's built-in rules and variables
#    (this increases performance and avoids hard-to-debug behaviour);
# o  print "Entering directory ...";
MAKEFLAGS += -rR --no-print-directory

CONFIG_CONFIG_IN=Kconfig
CONFIG_DEFCONFIG=.defconfig
CONFIG=config

# Detect if the current command interpreter is a Unix-like shell,
# or a non-*nix interpreter such as Command Prompt/PowerShell.
# When given a *nix shell, we'll further assume that all the necessary tools
# are installed, and enable all the rules and features.
# With a non-*nix interpreter we'll trim the recipes, skipping unsupported
# commands and scripts.
ifeq ($(strip $(shell echo $$0)),$$0)
NIX_SHELL=
else
NIX_SHELL=1
CONFIG_SHELL := $(shell which bash)
ifeq ($(CONFIG_SHELL),)
$(error GNU Bash is needed to build Kconfig host tools!)
endif
endif

BINDIR:=binaries
SYMLINK ?= at91bootstrap.bin
SYMLINK_BOOT ?= boot.bin

include	host-utilities/host.mk

ifeq ($(origin BUILD_DATE), undefined)
# Automatically escape '%' symbols when recipes are implemented as batch files
PERCENT := $(if $(findstring %%_dummy_,$(shell echo %%_dummy_)),%,%%)
# see https://reproducible-builds.org/docs/source-date-epoch/#makefile
DATE_FMT = +$(PERCENT)Y-$(PERCENT)m-$(PERCENT)d $(PERCENT)H:$(PERCENT)M:$(PERCENT)S
ifdef SOURCE_DATE_EPOCH
	BUILD_DATE := $(shell $(DATE) -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>$(DEV_NULL) || $(DATE) -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>$(DEV_NULL) || $(DATE) -u "$(DATE_FMT)")
else
	BUILD_DATE := $(shell $(DATE) "$(DATE_FMT)")
endif
endif
VERSION := 3.10.0
REVISION :=
SCMINFO := $(shell (host-utilities/setlocalversion))

ifeq ($(SCMINFO),)
-include scminfo.mk
SCMINFO=$(RECORD_SCMINFO)
endif

TARBALL_NAME=AT91Bootstrap-$(VERSION)$(REV)$(SCMINFO).tar.gz
TARBALL_DIR=../tarball_dir
TARBALL_PREFIX=at91bootstrap-$(VERSION)$(REV)/

# Use 'make V=1' for the verbose mode
ifneq ("$(origin V)", "command line")
Q=@
export Q
endif

# Function brief: Search for existing files, recursively
# param 1: List of directories to search in
# param 2: List of patterns to match
# return: List of paths to the matching files, relative to and including param 1
rwildcard=$(foreach d,$(wildcard $(1:=/*)),$(strip $(call rwildcard,$d,$2) $(filter $(subst *,%,$2),$d)))

noconfig_targets:= menuconfig defconfig $(CONFIG) oldconfig savedefconfig

# Check first if we want to configure at91bootstrap
#
ifeq ($(filter $(noconfig_targets),$(MAKECMDGOALS)),)
-include .config
endif

ifeq ($(CONFIG_HAVE_DOT_CONFIG),)

all: menuconfig

# configuration
# ---------------------------------------------------------------------------

HOSTCFLAGS=$(CFLAGS_FOR_BUILD)
export HOSTCFLAGS

$(CONFIG)/conf: | $(CONFIG)/at91bootstrap-config
	@$(MAKE) CC="$(HOSTCC)" -C $(CONFIG) conf
	-@if [ ! -f .config ]; then \
		cp $(CONFIG_DEFCONFIG) .config; \
	fi

$(CONFIG)/mconf: | $(CONFIG)/at91bootstrap-config
	@$(MAKE) CC="$(HOSTCC)" -C $(CONFIG) conf mconf
	-@if [ ! -f .config ]; then \
		cp $(CONFIG_DEFCONFIG) .config; \
	fi

menuconfig: $(CONFIG)/mconf
	@if ! KCONFIG_AUTOCONFIG=$(CONFIG)/at91bootstrap-config/auto.conf \
		KCONFIG_AUTOHEADER=$(CONFIG)/at91bootstrap-config/autoconf.h \
		$(CONFIG)/mconf $(CONFIG_CONFIG_IN); then \
		test -f .config.cmd || rm -f .config; \
	fi

$(CONFIG): $(CONFIG)/conf
	@KCONFIG_AUTOCONFIG=$(CONFIG)/at91bootstrap-config/auto.conf \
		KCONFIG_AUTOHEADER=$(CONFIG)/at91bootstrap-config/autoconf.h \
		$(CONFIG)/conf $(CONFIG_CONFIG_IN)

oldconfig: $(CONFIG)/conf
	@KCONFIG_AUTOCONFIG=$(CONFIG)/at91bootstrap-config/auto.conf \
		KCONFIG_AUTOHEADER=$(CONFIG)/at91bootstrap-config/autoconf.h \
		$(CONFIG)/conf --oldconfig $(CONFIG_CONFIG_IN)

defconfig: $(CONFIG)/conf
	@KCONFIG_AUTOCONFIG=$(CONFIG)/at91bootstrap-config/auto.conf \
		KCONFIG_AUTOHEADER=$(CONFIG)/at91bootstrap-config/autoconf.h \
		$(CONFIG)/conf --defconfig=.config $(CONFIG_CONFIG_IN)

savedefconfig: $(CONFIG)/conf
	@KCONFIG_AUTOCONFIG=$(CONFIG)/at91bootstrap-config/auto.conf \
		KCONFIG_AUTOHEADER=$(CONFIG)/at91bootstrap-config/autoconf.h \
		$(CONFIG)/conf --savedefconfig=defconfig $(CONFIG_CONFIG_IN)

else #  Have DOT Config

HOSTARCH := $(patsubst arm%,arm,$(shell uname -m))

AS=$(CROSS_COMPILE)gcc
CC=$(CROSS_COMPILE)gcc
LD=$(CROSS_COMPILE)ld
NM= $(CROSS_COMPILE)nm
SIZE=$(CROSS_COMPILE)size
OBJCOPY=$(CROSS_COMPILE)objcopy
OBJDUMP=$(CROSS_COMPILE)objdump

PROJECT := $(strip $(subst ",,$(CONFIG_PROJECT)))
IMG_ADDRESS := $(strip $(subst ",,$(CONFIG_IMG_ADDRESS)))
IMG_SIZE := $(strip $(subst ",,$(CONFIG_IMG_SIZE)))
JUMP_ADDR := $(strip $(subst ",,$(CONFIG_JUMP_ADDR)))
OF_OFFSET := $(strip $(subst ",,$(CONFIG_OF_OFFSET)))
OF_ADDRESS := $(strip $(subst ",,$(CONFIG_OF_ADDRESS)))
BOOTSTRAP_MAXSIZE := $(strip $(subst ",,$(CONFIG_BOOTSTRAP_MAXSIZE)))
MEMORY := $(strip $(subst ",,$(CONFIG_MEMORY)))
IMAGE_NAME:= $(strip $(subst ",,$(CONFIG_IMAGE_NAME)))
CARD_SUFFIX := $(strip $(subst ",,$(CONFIG_CARD_SUFFIX)))
MEM_BANK := $(strip $(subst ",,$(CONFIG_MEM_BANK)))
MEM_BANK2 := $(strip $(subst ",,$(CONFIG_MEM_BANK2)))
LINUX_KERNEL_ARG_STRING := $(strip $(subst ",,$(CONFIG_LINUX_KERNEL_ARG_STRING)))
LINUX_KERNEL_ARG_STRING_FILE := $(strip $(subst ",,$(CONFIG_LINUX_KERNEL_ARG_STRING_FILE)))

# Board definitions
BOARDNAME:=$(strip $(subst ",,$(CONFIG_BOARDNAME)))

ifeq ($(CONFIG_OVERRIDE_MACH_TYPE), y)
MACH_TYPE:=$(strip $(subst ",,$(CONFIG_CUSTOM_MACH_TYPE)))
else
MACH_TYPE:=$(strip $(subst ",,$(CONFIG_MACH_TYPE)))
endif

LINK_ADDR:=$(strip $(subst ",,$(CONFIG_LINK_ADDR)))
DATA_SECTION_ADDR:=$(strip $(subst ",,$(CONFIG_DATA_SECTION_ADDR)))
TOP_OF_MEMORY:=$(strip $(subst ",,$(CONFIG_TOP_OF_MEMORY)))

# CRYSTAL is UNUSED
CRYSTAL:=$(strip $(subst ",,$(CONFIG_CRYSTAL)))

# driver definitions
SPI_CLK:=$(strip $(subst ",,$(CONFIG_SPI_CLK)))
SPI_BOOT:=$(strip $(subst ",,$(CONFIG_SPI_BOOT)))

ifeq ($(REVISION),)
REV:=
else
REV:=-$(strip $(subst ",,$(REVISION)))
endif

ifeq ($(CONFIG_OF_LIBFDT), y)
BLOB:=-dt
else
BLOB:=
endif

ifeq ($(CONFIG_LOAD_LINUX), y)
TARGET_NAME:=linux-$(or $(subst I,i,$(IMAGE_NAME)),image)
endif

ifeq ($(CONFIG_LOAD_ANDROID), y)
TARGET_NAME:=android-$(or $(subst I,i,$(IMAGE_NAME)),image)
endif

ifeq ($(CONFIG_LOAD_UBOOT), y)
TARGET_NAME:=$(or $(subst -,,$(basename $(IMAGE_NAME))),uboot)
endif

ifeq ($(CONFIG_LOAD_64KB), y)
TARGET_NAME:=$(or $(basename $(IMAGE_NAME)),softpack)
endif

ifeq ($(CONFIG_LOAD_1MB), y)
TARGET_NAME:=$(or $(basename $(IMAGE_NAME)),softpack)
endif

ifeq ($(CONFIG_LOAD_4MB), y)
TARGET_NAME:=$(or $(basename $(IMAGE_NAME)),softpack)
endif

ifeq ($(CONFIG_LOAD_NONE), y)
TARGET_NAME:=$(or $(basename $(IMAGE_NAME)),none)
endif

BOOT_NAME=$(BOARDNAME)-$(PROJECT)$(CARD_SUFFIX)boot-$(TARGET_NAME)$(BLOB)-$(VERSION)$(REV)
AT91BOOTSTRAP:=$(BINDIR)/$(BOOT_NAME).bin

ifeq ($(IMAGE),)
IMAGE=$(BOOT_NAME).bin
endif

COBJS-y:= main.o
SOBJS-y:= crt0_gnu.o

# Verify that BOARDNAME is the name of a subdirectory of board/
BOARD_LOCATE=$(if $(wildcard board/$(BOARDNAME)/.),board/$(BOARDNAME))
ifeq ("$(realpath $(BOARD_LOCATE))", "")
# List all vendor subdirectories found under contrib/board/
CONTRIB_VENDORS=$(foreach file,$(wildcard contrib/board/*),$(if $(wildcard $(addsuffix /.,$(file))),$(file)))
BOARD_LOCATE=$(firstword $(wildcard $(addsuffix /$(BOARDNAME),$(CONTRIB_VENDORS))))
ifeq ("$(realpath $(BOARD_LOCATE))", "")
$(error ERROR: *** $(BOARDNAME) board not found!)
endif
endif

COBJS-y += $(BOARD_LOCATE)/$(BOARDNAME).o
INCL = $(BOARD_LOCATE)

include	lib/lib.mk
include	driver/driver.mk
include	contrib/driver/driver.mk
include	fs/src/fat.mk

GC_SECTIONS=--gc-sections

NOSTDINC_FLAGS := -nostdinc -isystem "$(shell $(CC) -print-file-name=include)"

CPPFLAGS=$(NOSTDINC_FLAGS) -ffunction-sections -g -Os -Wall \
	-mno-unaligned-access \
	-fno-stack-protector -fno-common -fno-builtin -fno-jump-tables -fno-pie \
	-I$(INCL) -Icontrib/include -Iinclude -Ifs/include \
	-I$(CONFIG)/at91bootstrap-config \
	-DAT91BOOTSTRAP_VERSION=\"$(VERSION)$(REV)$(SCMINFO)\" -DCOMPILE_TIME="\"$(BUILD_DATE)\""

ASFLAGS=-g -Os -Wall -I$(INCL) -Iinclude -Icontrib/include

include	toplevel_cpp.mk
include	board/board_cpp.mk

ifneq ("$(wildcard $(BOARD_LOCATE)/board.mk)", "")
include $(BOARD_LOCATE)/board.mk
else
$(warning WARNING: *** file: $(BOARD_LOCATE)/board.mk are not found!)
endif

include	driver/driver_cpp.mk

OBJS:= $(SOBJS-y) $(COBJS-y)

ifeq ($(CONFIG_ENTER_NWD), y)
link_script:=elf32-littlearm-tz.lds
else
link_script:=elf32-littlearm.lds
endif

# Linker flags.
#  -Wl,...:     tell GCC to pass this to linker.
#    -Map:      create map file
#    --cref:    add cross reference to map file
#  -lc 	   : 	tells the linker to tie in newlib
#  -lgcc   : 	tells the linker to tie in newlib
LDFLAGS=-nostartfiles -Map=$(BINDIR)/$(BOOT_NAME).map --cref -static
LDFLAGS+=-T $(link_script) $(GC_SECTIONS) -Ttext $(LINK_ADDR)

ifneq ($(DATA_SECTION_ADDR),)
LDFLAGS+=-Tdata $(DATA_SECTION_ADDR)
endif

REMOVE_SECTIONS=-R .note -R .comment -R .note.gnu.build-id

gccversion := $(shell $(CC) -dumpversion)

ifdef YYY   # For other utils
ifeq ($(CC),gcc) 
TARGETS=no-cross-compiler
else
TARGETS=$(AT91BOOTSTRAP) host-utilities .config filesize
endif
endif

TARGETS=$(AT91BOOTSTRAP)

PHONY:=all

all: CheckCrossCompile PrintFlags $(AT91BOOTSTRAP) ChkFileSize ${AT91BOOTSTRAP}.pmecc

CheckCrossCompile:
	$(if $(filter-out arm,$(HOSTARCH)),$(if $(CROSS_COMPILE),, \
		$(error Error: the CROSS_COMPILE environment variable must be defined)))

PrintFlags:
	$(info CC)
	$(info ========)
	$(info $(CC) $(gccversion))
	$(info )
	$(info as FLAGS)
	$(info ========)
	$(info $(ASFLAGS))
	$(info )
	$(info gcc FLAGS)
	$(info =========)
	$(info $(CPPFLAGS))
	$(info )
	$(info ld FLAGS)
	$(info ========)
	$(info $(LDFLAGS))
	$(info )

$(AT91BOOTSTRAP): $(OBJS) | $(BINDIR)
	@echo "  LD        "$(BOOT_NAME).elf
	$(Q)$(LD) $(LDFLAGS) -n -o $(BINDIR)/$(BOOT_NAME).elf $(OBJS)
#	@$(OBJCOPY) --strip-debug --strip-unneeded $(REMOVE_SECTIONS) $(BINDIR)/$(BOOT_NAME).elf -O binary $(BINDIR)/$(BOOT_NAME).bin
	@$(OBJCOPY) --strip-all $(REMOVE_SECTIONS) $(BINDIR)/$(BOOT_NAME).elf -O binary $@
	@ln -sf $(BOOT_NAME).bin ${BINDIR}/${SYMLINK}
	@ln -sf $(BOOT_NAME).bin ${BINDIR}/${SYMLINK_BOOT}

%.o : %.c .config
	@echo "  CC        "$<
	@$(CC) $(CPPFLAGS) -c -o $@ $<

%.o : %.S .config
	@echo "  AS        "$<
	@$(AS) $(ASFLAGS)  -c -o $@  $<

$(AT91BOOTSTRAP).pmecc: $(AT91BOOTSTRAP)
ifeq ($(CONFIG_NANDFLASH), y)
ifeq ($(CONFIG_USE_PMECC), y)
	$(Q)./scripts/addpmecchead.py $(AT91BOOTSTRAP) $(AT91BOOTSTRAP).pmecc $(BOARDNAME)
endif
endif

PHONY+= bootstrap

rebuild: clean all

ChkFileSize: $(AT91BOOTSTRAP)
	@( fsize=`./scripts/get_sram_size.sh $(BINDIR)/$(BOOT_NAME).map`; \
	  if [ $$? -ne 0 ] ; then \
		rm $(BINDIR)/$(BOOT_NAME).bin ;\
		rm ${BINDIR}/${SYMLINK}; \
		rm ${BINDIR}/${SYMLINK_BOOT}; \
		exit 3; \
	  fi ; \
	  echo "Size of $(BOOT_NAME).bin is $$fsize bytes"; \
	  if [ "$$fsize" -gt "$(BOOTSTRAP_MAXSIZE)" ] ; then \
		echo "[Failed***] It's too big to fit into SRAM area. the support maxium size is $(BOOTSTRAP_MAXSIZE)"; \
		rm $(BINDIR)/$(BOOT_NAME).bin ;\
		rm ${BINDIR}/${SYMLINK}; \
		rm ${BINDIR}/${SYMLINK_BOOT}; \
		exit 2;\
	  else \
	  	echo "[Succeeded] It's OK to fit into SRAM area"; \
		stack_space=`expr $(BOOTSTRAP_MAXSIZE) - $$fsize`; \
		echo "[Attention] The space left for stack is $$stack_space bytes"; \
	  fi )
endif  # CONFIG_HAVE_DOT_CONFIG

PHONY+= rebuild

prepare: .prepared | $(CONFIG)/at91bootstrap-config $(BINDIR)

.prepared: $(wildcard .config)
	@echo AT91BOOTSTRAP_PREP_REV = 1 > $@
	@echo AT91BOOTSTRAP_VERSIONNO = $(VERSION) >> $@
	@echo KBUILD_KCONFIG = $(CONFIG_CONFIG_IN) >> $@
	@echo KCONFIG_DEFCONFIG = $(CONFIG_DEFCONFIG) >> $@
	@echo KCONFIG_CONFIG = .config >> $@
	@echo KCONFIG_AUTOCONFIGOUTPUT = $(CONFIG)/at91bootstrap-config >> $@
	@echo KCONFIG_AUTOCONFIG = $(CONFIG)/at91bootstrap-config/auto.conf >> $@
	@echo KCONFIG_AUTOHEADER = $(CONFIG)/at91bootstrap-config/autoconf.h >> $@
	@echo AT91BOOTSTRAP_SCMREVIN = scminfo.mk >> $@
	@echo AT91BOOTSTRAP_BINOUTPUT = $(BINDIR) >> $@
	@echo AT91BOOTSTRAP_BIN = $(BINDIR)/$(SYMLINK) >> $@
ifneq ($(CONFIG_HAVE_DOT_CONFIG),)
	@echo AT91BOOTSTRAP_ELF = $(BINDIR)/$(BOOT_NAME).elf >> $@
	@echo AT91BOOTSTRAP_MAP = $(BINDIR)/$(BOOT_NAME).map >> $@
	@echo AT91BOOTSTRAP_VERSIONEXTD = $(VERSION)$(REV)$(SCMINFO) >> $@
	@echo AT91BOOTSTRAP_TARBALL = $(TARBALL_NAME) >> $@
	@echo AT91BOOTSTRAP_BOARDDIR = $(BOARD_LOCATE) >> $@
endif

$(CONFIG)/at91bootstrap-config $(BINDIR):
	@$(MKDIR) -p $@

%_defconfig:
	@(conf_file=`find ./ -name $@`; \
	if [ "$$conf_file"x != "x" ]; then \
		cp $$conf_file .config; \
	else \
		echo "Error: *** Cannot find file: $@"; \
		exit 2; \
	fi )
	@$(MAKE) defconfig

update:
	cp .config $(BOARD_LOCATE)/$(BOARDNAME)_defconfig

no-cross-compiler:
	@echo
	@echo
	@echo
	@echo "	You should consider using a cross compiler for ARM !"
	@echo "	I.E: CROSS_COMPILE should contains something useful."
	@echo

debug:
	@echo CONFIG=$(CONFIG)
	@echo AS=$(AS)
	@echo CROSS_COMPILE=$(CROSS_COMPILE)

PHONY+=update no-cross-compiler debug

distrib: mrproper
	$(Q)rm -f  $(call rwildcard,.,*.elf *.map)
	$(Q)rm -fr result
	$(Q)rm -fr build
	$(Q)rm -fr ..make.deps.tmp
	$(Q)rm -fr config/conf

config-clean:
	@echo "  CLEAN        "configuration files!
	$(Q)$(MAKE) -C config distclean
	$(Q)rm -fr config/at91bootstrap-config
	$(Q)rm -f  config/.depend

clean:
	@echo "  CLEAN        "obj and misc files!
	$(Q)rm -f $(CONFIG)/.depend
	$(Q)rm -f $(call rwildcard,.,*.o *.srec *~)

distclean: clean config-clean
#	rm -fr $(BINDIR)
	$(Q)rm -fr .config .config.cmd .config.old
	$(Q)rm -fr .auto.deps
	$(Q)rm -f .installed
	$(Q)rm -f ..*.tmp
	$(Q)rm -f .configured
	$(Q)rm -f .prepared

mrproper: distclean
	@echo "  CLEAN        "binary files!
	$(Q)rm -fr $(BINDIR)
	$(Q)rm -fr log

PHONY+=distrib config-clean clean distclean mrproper

tarball:
	@echo "Tar the source code to ${TARBALL_NAME}"
	$(Q)$(MKDIR) -p ${TARBALL_DIR}
	$(Q)git archive --prefix=${TARBALL_PREFIX} HEAD | gzip > ${TARBALL_DIR}/${TARBALL_NAME}
	$(Q)echo "RECORD_SCMINFO=${SCMINFO}" > ${TARBALL_DIR}/scminfo.mk
	$(Q)cd ${TARBALL_DIR}; tar -xzf ${TARBALL_NAME}
	$(Q)cd ${TARBALL_DIR}; cp scminfo.mk ${TARBALL_PREFIX}
	$(Q)cd ${TARBALL_DIR}; tar czf ${TARBALL_NAME} ${TARBALL_PREFIX}
	$(Q)cp ${TARBALL_DIR}/${TARBALL_NAME} .
	$(Q)rm -rf ${TARBALL_DIR}

PHONY+=tarball

.PHONY: $(PHONY)

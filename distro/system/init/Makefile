# Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.

#
# Makefile for init
#

include ../../../config.mk

TARGET_EXEC = init
TARGET_NAME = init
CAPP_BUILD  = ./$(TARGET_NAME)_capp/
CAPP        = ./capp/

UKAMAOS_ROOT = ../../../

VENDOR_DIR   = $(UKAMAOS_ROOT)/distro/vendor/
VENDOR_BUILD = $(VENDOR_DIR)/build/
VENDOR_INC   = $(VENDOR_BUILD)/include/
VENDOR_USR_INC = $(VENDOR_BUILD)/usr/include
SYS_INC = $(XGCCPATH)/../$(HOST)/include
VENDOR_LIB   = $(VENDOR_BUILD)/lib/
VENDOR_LIB64 = $(VENDOR_BUILD)/lib64/

SCRIPTS = $(UKAMAOS_ROOT)/distro/scripts/
ROOTFS  = ./_ukama_minimal_rootfs/

# Packages needed to compile target.
VENDOR_PKG =

# Setting up various compile and link flags.
CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I$(VENDOR_INC)
CFLAGS+=-I$(VENDOR_USR_INC)
CFLAGS+=-I$(SYS_INC)
CFLAGS+=-I./inc
CFLAGS+=-I$(WIMC_DIR)/inc/
CFLAGS+=-D_REENTRANT

LDFLAGS+=-L${VENDOR_LIB}
LDFLAGS+=-L${VENDOR_LIB64}
LDFLAGS+=-static

LIBS=

# Soruce files.
CFILES   = $(wildcard ./src/*.c)
OBJFILES = $(CFILES:.c=.o)

XCC := $(XGCCPATH)$(XGCC)
XLD := $(XGCCPATH)$(XLD)

.PHONY: $(TARGET_EXEC)

lxce: $(TARGET_EXEC)

$(TARGET_EXEC): $(VENDOR_PKG) $(OBJFILES)
	$(XCC) -o $(TARGET_EXEC) $(OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

$(VENDOR_PKG):
	cd $(VENDOR_DIR) && make $@

%.o: %.c
	$(XCC) $(CFLAGS) -c $< -o $@

capp: $(TARGET_EXEC)
	rm -rf $(CAPP_BUILD); mkdir $(CAPP_BUILD)
	cd $(SCRIPTS); ./mk_minimal_rootfs.sh; cd -
	cp -rf $(SCRIPTS)/$(ROOTFS)/* $(CAPP_BUILD)

	mkdir -p $(CAPP_BUILD)/conf $(CAPP_BUILD)/sbin $(CAPP_BUILD)/cert
	cp $(CAPP)/config.json $(CAPP_BUILD)
	cp $(CAPP)/conf/* $(CAPP_BUILD)/conf
	cp $(TARGET_EXEC) $(CAPP_BUILD)/sbin
	tar -czf $(TARGET_NAME)_capp.tar.gz $(CAPP_BUILD)
	rm -rf $(CAPP_BUILD)

	echo "$(TARGET_NAME) capp is available $(TARGET_NAME)_capp.tar.gz Done."

clean:
	rm -rf $(SCRIPTS)/$(ROOTFS)
	rm -rf $(TARGET_EXEC) $(OBJFILES) $(CAPP_BUILD) *.tar.gz

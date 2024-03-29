# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2021-present, Ukama Inc.

include ../../../config.mk

TARGET_EXEC = lxce.d
TARGET_NAME = lxce
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

PLATFORM_DIR     = $(UKAMAOS_ROOT)/distro/platform/
PLATFORM_BUILD   = $(PLATFORM_DIR)/build/
PLATFORM_INC_SYS = $(PLATFORM_DIR)/sys/inc
PLATFORM_INC_LOG = $(PLATFORM_DIR)/log/inc
PLATFORM_LIB_DIR = $(PLATFORM_DIR)/build/

SCRIPTS = $(UKAMAOS_ROOT)/distro/scripts/
ROOTFS  = ./_ukama_minimal_rootfs/

WIMC_DIR     = $(UKAMAOS_ROOT)/distro/system/wimc/

# Packages needed to compile target.
VENDOR_PKG = tomlc ulfius libcap

# Setting up various compile and link flags.
CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I$(VENDOR_INC)
CFLAGS+=-I$(VENDOR_USR_INC)
CFLAGS+=-I$(SYS_INC)
CFLAGS+=-I$(PLATFORM_INC_SYS)
CFLAGS+=-I$(PLATFORM_INC_LOG)
CFLAGS+=-I./inc
CFLAGS+=-I$(WIMC_DIR)/inc/
CFLAGS+=-D_REENTRANT

LDFLAGS+=-L$(VENDOR_LIB)
LDFLAGS+=-L$(VENDOR_LIB64)
LDFLAGS+=-L$(PLATFORM_LIB_DIR)

# Libraries needed.
ULFIUS_LIB=-lulfius -lorcania -lyder -lmicrohttpd -lgnutls -lnettle
ULFIUS_LIB+=-lhogweed -lp11-kit -lz
CURL_LIB=-lcurl -lssl -lcrypto

LIBS = $(ULFIUS_LIB) $(CURL_LIB)
LIBS+=-ltoml
LIBS+=-lpthread
LIBS+=-ljansson
LIBS+=-lcap
LIBS+=-luuid
LIBS+=-lusys

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

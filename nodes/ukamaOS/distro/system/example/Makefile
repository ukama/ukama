# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2021-present, Ukama Inc.

include ../../../config.mk

TARGET_EXEC = example.d
TARGET_NAME = example
CAPP_BUILD  = ./$(TARGET_NAME)_capp/
CAPP        = ./capp/

SCRIPTS = $(UKAMAOS_ROOT)/distro/scripts
ROOTFS  = ./_ukama_minimal_rootfs

# Packages needed to compile target.
VENDOR_PKG = 

# Setting up various compile and link flags.
CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I$(VENDOR_INC)
CFLAGS+=-I./inc
CFLAGS+=-D_REENTRANT
CFLAGS+=-static

LDFLAGS+=$(RPATH_FLAGS)
LDFLAGS+=-L${VENDOR_LIB}
LDFLAGS+=-static

# Libraries needed.
LIBS =

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

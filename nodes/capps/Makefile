# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2021-present, Ukama Inc.

TARGET_EXEC  = capp
UKAMAOS_ROOT = ../ukamaOS

include $(UKAMAOS_ROOT)/config.mk

VENDOR_DIR     = $(UKAMAOS_ROOT)/distro/vendor/
VENDOR_BUILD   = $(VENDOR_DIR)/build/
VENDOR_INC     = $(VENDOR_BUILD)/include/
VENDOR_USR_INC = $(VENDOR_BUILD)/usr/include
SYS_INC        = $(XGCCPATH)/../$(HOST)/include
VENDOR_LIB     = $(VENDOR_BUILD)/lib/
VENDOR_LIB64   = $(VENDOR_BUILD)/lib64/

SCRIPTS = $(UKAMAOS_ROOT)/distro/scripts/
ROOTFS  = ./_ukama_minimal_rootfs/

# Packages needed to compile target.
VENDOR_PKG = tomlc

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

# Libraries needed.
LIBS=-ltoml

# Soruce files.
CFILES   = $(wildcard ./src/*.c)
OBJFILES = $(CFILES:.c=.o)

XCC := $(XGCCPATH)$(XGCC)
XLD := $(XGCCPATH)$(XLD)

.PHONY: $(TARGET_EXEC)

capp: $(TARGET_EXEC)

$(TARGET_EXEC): $(VENDOR_PKG) $(OBJFILES)
	$(XCC) -o $(TARGET_EXEC) $(OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

$(VENDOR_PKG):
	cd $(VENDOR_DIR) && make $@

%.o: %.c
	$(XCC) $(CFLAGS) -c $< -o $@

clean:
	rm -rf $(TARGET_EXEC) $(OBJFILES)

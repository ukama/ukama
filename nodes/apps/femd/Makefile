# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

include ../../ukamaOS/config.mk

TARGET_EXEC = femd.d

UKAMAOS_ROOT = $(NODES_DIR)/ukamaOS

VENDOR_DIR   = $(UKAMAOS_ROOT)/distro/vendor/
VENDOR_BUILD = $(VENDOR_DIR)/build/
VENDOR_INC   = $(VENDOR_BUILD)/include/
VENDOR_LIB   = $(VENDOR_BUILD)/lib/
VENDOR_LIB64 = $(VENDOR_BUILD)/lib64/

PLATFORM_DIR     = $(UKAMAOS_ROOT)/distro/platform/
PLATFORM_BUILD   = $(PLATFORM_DIR)/build/
PLATFORM_INC_SYS = $(PLATFORM_DIR)/sys/inc
PLATFORM_INC_LOG = $(PLATFORM_DIR)/log/inc
PLATFORM_LIB     = $(PLATFORM_DIR)/build/

VENDOR_PKG=
# Skip platform build on macOS for development
ifeq ($(shell uname -s),Darwin)
    PLATFORM=
    # Add jansson lib path for macOS
    LDFLAGS+=-L/opt/homebrew/opt/jansson/lib
else
    PLATFORM=platform
endif

CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I.
CFLAGS+=-I$(VENDOR_INC)
CFLAGS+=-I$(PLATFORM_INC_SYS)
CFLAGS+=-I$(PLATFORM_INC_LOG)
CFLAGS+=-I./inc
CFLAGS+=-D_REENTRANT

# Add jansson include path for macOS
ifeq ($(shell uname -s),Darwin)
    CFLAGS+=-I/opt/homebrew/opt/jansson/include
endif

LDFLAGS+=-L${VENDOR_LIB} 
LDFLAGS+=-L$(VENDOR_LIB64)
LDFLAGS+=-L$(PLATFORM_LIB)

LIBS+=-lpthread
# Skip usys on macOS for development
ifeq ($(shell uname -s),Darwin)
    LIBS+=-ljansson
else
    LIBS+=-lusys
endif

TARGET_CFILES   = $(wildcard ./src/*.c)
TARGET_OBJFILES = $(TARGET_CFILES:.c=.o)

XCC := $(XGCCPATH)$(XGCC)
XLD := $(XGCCPATH)$(XLD)

ifeq ($(XCC),)
    XCC := gcc
endif

.PHONY: all clean $(TARGET_EXEC)

all: version.h static.h $(TARGET_EXEC)

$(TARGET_EXEC): $(TARGET_OBJFILES)
ifneq ($(PLATFORM),)
	$(MAKE) -C $(PLATFORM_DIR)
endif
	$(XCC) -o $(TARGET_EXEC) $(TARGET_OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

version.h:
	$(NODES_DIR)/utils/scripts/generate_version.sh

static.h:
	$(NODES_DIR)/utils/scripts/generate_static.sh

%.o: %.c
	$(XCC) $(CFLAGS) -c $< -o $@

clean:
	rm -f $(TARGET_EXEC) $(TARGET_OBJFILES) version.h static.h
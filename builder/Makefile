# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2022-present, Ukama Inc.

include ../nodes/ukamaOS/config.mk

BUILDER_EXEC = builder
APP_BUILDER_EXEC = app_builder

UKAMAOS_ROOT = ../nodes/ukamaOS/

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

# Packages needed to compile target
VENDOR_PKG = jansson ulfius
PLATFORM   = platform

# Setting up various compile and link flags
CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I$(VENDOR_INC)
CFLAGS+=-I$(PLATFORM_INC_SYS)
CFLAGS+=-I$(PLATFORM_INC_LOG)
CFLAGS+=-I./inc
CFLAGS+=-I./inc/common/
CFLAGS+=-I../common/inc
CFLAGS+=-D_REENTRANT

LDFLAGS+=-L$(VENDOR_LIB) 
LDFLAGS+=-L$(VENDOR_LIB64)
LDFLAGS+=-L$(PLATFORM_LIB)

# Libraries needed
ULFIUS_LIB=-lulfius -lorcania -lyder -lmicrohttpd -lgnutls -lnettle
ULFIUS_LIB+=-lhogweed -lp11-kit -lz
LIBS = $(ULFIUS_LIB)
LIBS+=-ljansson
LIBS+=-lusys
LIBS+=-ltoml

# Soruce files
#CFILES = $(filter-out ./src/*_app.c, $(wildcard ./src/*.c))
CFILES = $(shell find ./src -name "*.c" ! -name "*_app.c")
OBJFILES = $(CFILES:.c=.o)

XCC := $(XGCCPATH)$(XGCC)
XLD := $(XGCCPATH)$(XLD)

.PHONY: $(BUILDER_EXEC)

$(BUILDER_EXEC): $(VENDOR_PKG) $(PLATFORM) $(OBJFILES)
	$(XCC) -o $(BUILDER_EXEC) $(OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

$(PLATFORM):
	$(MAKE) -C $(PLATFORM_DIR)

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@ 

%.o: %.c
	$(XCC) $(CFLAGS) -c $< -o $@

all: builder

clean:
	rm -f $(BUILDER_EXEC) $(OBJFILES)

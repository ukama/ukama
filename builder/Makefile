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
ifeq ($(ALPINE_BUILD),1)
    $(info Excluding vendor library for Alpine build)
	VENDOR_PKG=
else
	VENDOR_PKG = jansson ulfius tomlc
endif
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
CFLAGS+=$(if $(ALPINE_BUILD),-DALPINE_BUILD=$(ALPINE_BUILD))

LDFLAGS+=-L$(VENDOR_LIB) 
LDFLAGS+=-L$(VENDOR_LIB64)
LDFLAGS+=-L$(PLATFORM_LIB)

# Libraries needed
ULFIUS_LIB=-lulfius -lorcania -lmicrohttpd -lgnutls -lnettle
ULFIUS_LIB+=-lhogweed -lp11-kit -lz
ifeq ($(ALPINE_BUILD),1)
    $(info Excluding yder library for Alpine build)
else
    ULFIUS_LIB+=-lyder
endif
LIBS = $(ULFIUS_LIB)
LIBS+=-ljansson
LIBS+=-lusys
LIBS+=-ltoml

# Soruce files
CFILES   = $(shell find ./src -name "*.c" ! -name "*_app.c")
OBJFILES = $(CFILES:.c=.o)

APP_CFILES   = $(wildcard ./src/*_app.c)
APP_OBJFILES = $(APP_CFILES:.c=.o)

.PHONY: all clean

$(BUILDER_EXEC): $(VENDOR_PKG) $(PLATFORM) $(OBJFILES)
	$(CC) -o $(BUILDER_EXEC) $(OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

ifeq ($(ALPINE_BUILD),1)
$(APP_BUILDER_EXEC): $(PLATFORM) $(APP_OBJFILES)
	$(CC) -o $(APP_BUILDER_EXEC) $(APP_OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."
else
$(APP_BUILDER_EXEC): $(VENDOR_PKG) $(PLATFORM) $(APP_OBJFILES)
	$(CC) -o $(APP_BUILDER_EXEC) $(APP_OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."
endif
$(PLATFORM):
	$(MAKE) -C $(PLATFORM_DIR)

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@ 

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

all: builder app_builder

clean:
	rm -f $(BUILDER_EXEC) $(APP_BUILDER_EXEC) $(OBJFILES) $(APP_OBJFILES)

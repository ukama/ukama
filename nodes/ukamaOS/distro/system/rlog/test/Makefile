# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

include ../../../../config.mk

MOCK_SERVICE = mock_service
MOCK_UKAMA   = mock_ukama

UKAMAOS_ROOT = ../../../../

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
VENDOR_PKG = tomlc curl jansson ulfius libuuid
PLATFORM = platform

# Setting up various compile and link flags
CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I$(VENDOR_INC)
CFLAGS+=-I../inc
CFLAGS+=-I../../common/inc
CFLAGS+=-D_REENTRANT
CFLAGS+=-I$(PLATFORM_INC_SYS)
CFLAGS+=-I$(PLATFORM_INC_LOG)

LDFLAGS+=-L${VENDOR_LIB}
LDFLAGS+=-L$(VENDOR_LIB64)
LDFLAGS+=-L$(PLATFORM_LIB)

# Libraries needed
LIBS=-lusys
LIBS+=-ljansson

# Soruce files
CFILES   = $(wildcard ./*.c)
OBJFILES = $(CFILES:.c=.o)

XCC := $(XGCCPATH)$(XGCC)
XLD := $(XGCCPATH)$(XLD)

.PHONY: $(MOCK_SERVICE) $(MOCK_UKAMA)

$(MOCK_SERVICE): $(VENDOR_PKG) $(OBJFILES)
	$(XCC) -o $(MOCK_SERVICE) mock_service.o $(LDFLAGS) $(LIBS)

$(MOCK_UKAMA):  $(VENDOR_PKG) $(OBJFILES)
	$(XCC) -o $(MOCK_UKAMA) mock_ukama.o $(LDFLAGS) $(LIBS)

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@

$(PLATFORM):
	$(MAKE) -C $(PLATFORM_DIR)

%.o: %.c
	$(XCC) $(CFLAGS) -c $< -o $@

all: mock_service mock_ukama

clean:
	rm -f $(MOCK_SERVICE) $(MOCK_UKAMA) $(OBJFILES)

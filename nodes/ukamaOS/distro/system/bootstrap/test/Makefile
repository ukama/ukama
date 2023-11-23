# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2022-present, Ukama Inc.

include ../../../../config.mk

TARGET_EXEC = bootstrap_server

UKAMAOS_ROOT = ../../../../

VENDOR_DIR   = $(UKAMAOS_ROOT)/distro/vendor/
VENDOR_BUILD = $(VENDOR_DIR)/build/
VENDOR_INC   = $(VENDOR_BUILD)/include/
VENDOR_LIB   = $(VENDOR_BUILD)/lib/
VENDOR_LIB64 = $(VENDOR_BUILD)/lib64/

# Packages needed to compile target
VENDOR_PKG = tomlc curl jansson ulfius libuuid

# Setting up various compile and link flags
CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I$(VENDOR_INC)
CFLAGS+=-I../../common/inc
CFLAGS+=-D_REENTRANT

LDFLAGS+=-L${VENDOR_LIB}
LDFLAGS+=-L$(VENDOR_LIB64)

# Libraries needed
ULFIUS_LIB=-lulfius -lorcania -lyder -lmicrohttpd -lgnutls
ULFIUS_LIB+=-lnettle -lhogweed -lp11-kit -lz
CURL_LIB=-lcurl -lssl -lcrypto
LIBS = $(ULFIUS_LIB) $(CURL_LIB)
LIBS+=-ljansson
LIBS+=-luuid
LIBS+=-lpthread
LIBS+=-ltoml

# Soruce files
CFILES   = $(wildcard ./*.c)
OBJFILES = $(CFILES:.c=.o)

XCC := $(XGCCPATH)$(XGCC)
XLD := $(XGCCPATH)$(XLD)

.PHONY: $(TARGET_EXEC) $(AGENT_EXEC)

bootstrap_server: $(TARGET_EXEC)

$(TARGET_EXEC): $(VENDOR_PKG) $(OBJFILES) $(COMM_OBJFILES)
	$(XCC) -o $(TARGET_EXEC) $(OBJFILES) $(COMM_OBJFILES) $(LDFLAGS) $(LIBS)

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@

%.o: %.c
	$(XCC) $(CFLAGS) -c $< -o $@

all: bootstrap_server

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES)

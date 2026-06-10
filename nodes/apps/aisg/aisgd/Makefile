# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

include ../../../ukamaOS/config.mk

TARGET_EXEC = aisgd
UKAMAOS_ROOT = $(NODES_DIR)/ukamaOS
VENDOR_DIR   = $(UKAMAOS_ROOT)/distro/vendor/
VENDOR_BUILD = $(VENDOR_DIR)/build/
VENDOR_INC   = $(VENDOR_BUILD)/include/
VENDOR_LIB   = $(VENDOR_BUILD)/lib/
VENDOR_LIB64 = $(VENDOR_BUILD)/lib64/
PLATFORM_DIR     = $(UKAMAOS_ROOT)/distro/platform/
PLATFORM_INC_SYS = $(PLATFORM_DIR)/sys/inc
PLATFORM_INC_LOG = $(PLATFORM_DIR)/log/inc
PLATFORM_LIB     = $(PLATFORM_DIR)/build/

ifeq ($(ALPINE_BUILD),1)
    VENDOR_PKG=
else
    VENDOR_PKG = tomlc jansson ulfius
endif
PLATFORM = platform

CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-Wextra
CFLAGS+=-O0
CFLAGS+=-I.
CFLAGS+=-I$(VENDOR_INC)
CFLAGS+=-I$(PLATFORM_INC_SYS)
CFLAGS+=-I$(PLATFORM_INC_LOG)
CFLAGS+=-I./inc
CFLAGS+=-D_REENTRANT

COVFLAGS=-fprofile-arcs -ftest-coverage
LDFLAGS+=-L${VENDOR_LIB}
LDFLAGS+=-L$(VENDOR_LIB64)
LDFLAGS+=-L$(PLATFORM_LIB)
LDFLAGS+=--coverage

CURL_LIB  = -lcurl -lssl -lcrypto

ULFIUS_LIB=-lulfius
ULFIUS_LIB+=-lorcania
ULFIUS_LIB+=-lmicrohttpd
ULFIUS_LIB+=-lgnutls
ULFIUS_LIB+=-lnettle
ULFIUS_LIB+=-lhogweed
ULFIUS_LIB+=-lp11-kit
ULFIUS_LIB+=-lz

LIBS = $(ULFIUS_LIB) $(CURL_LIB)
LIBS+=-ltoml
LIBS+=-ljansson
LIBS+=-lpthread
LIBS+=-lusys

TARGET_CFILES   = $(wildcard ./src/*.c)
TARGET_OBJFILES = $(TARGET_CFILES:.c=.o)

.PHONY: all clean test $(TARGET_EXEC)

all: version.h static.h $(TARGET_EXEC)

$(TARGET_EXEC): $(VENDOR_PKG) $(PLATFORM) $(TARGET_OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(TARGET_OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

$(PLATFORM):
	$(MAKE) -C $(PLATFORM_DIR)

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@

version.h:
	$(NODES_DIR)/utils/scripts/generate_version.sh

static.h:
	$(NODES_DIR)/utils/scripts/generate_static.sh

%.o: %.c
	$(CC) $(CFLAGS) $(COVFLAGS) -c $< -o $@

install: $(TARGET_EXEC)
	mkdir -p $(INSTALL_DIR)$(NODE_APP_DIR)
	mkdir -p $(INSTALL_DIR)$(NODE_APP_CONFIG_DIR)/aisgd
	cp -vrf $(TARGET_EXEC) $(INSTALL_DIR)$(NODE_APP_DIR)
	cp -vrf config/config.toml \
		$(INSTALL_DIR)$(NODE_APP_CONFIG_DIR)/aisgd/config.toml

clean:
	rm -f $(TARGET_EXEC) $(TARGET_OBJFILES) version.h static.h
	rm -f *.gcda *.gcno *.gcov
	rm -f src/*.gcda src/*.gcov src/*.gcno

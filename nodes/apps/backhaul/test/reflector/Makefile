# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

include ../../../../ukamaOS/config.mk

TARGET_EXEC = ukama-test-reflector

ARTIFACTS = $(TARGET_EXEC)

# Packages needed to compile target
ifeq ($(ALPINE_BUILD),1)
    $(info Excluding vendor library for Alpine build)
    VENDOR_PKG=
else
    VENDOR_PKG = jansson ulfius
endif

PLATFORM = platform

CFLAGS  = -c
CFLAGS += -g
CFLAGS += -Wall
CFLAGS += -Wextra
CFLAGS += -O0
CFLAGS += -I.
CFLAGS += -I$(VENDOR_INC)
CFLAGS += -I$(PLATFORM_INC_SYS)
CFLAGS += -I$(PLATFORM_INC_LOG)
CFLAGS += -D_REENTRANT

LDFLAGS += $(RPATH_FLAGS)
LDFLAGS += -L$(VENDOR_LIB)
LDFLAGS += -L$(VENDOR_LIB64)
LDFLAGS += -L$(PLATFORM_LIB)

# Libraries needed
CURL_LIB   = -lcurl -lssl -lcrypto
ULFIUS_LIB = -lulfius -lorcania -lmicrohttpd -lgnutls -lnettle
ULFIUS_LIB += -lhogweed -lp11-kit -lz
ifneq ($(ALPINE_BUILD),1)
ULFIUS_LIB += -lyder
endif

LIBS  = $(ULFIUS_LIB) $(CURL_LIB)
LIBS += -ljansson
LIBS += -lpthread
LIBS += -lusys

CFILES   = $(wildcard *.c)
OBJFILES = $(CFILES:.c=.o)

.PHONY: all clean install $(TARGET_EXEC) $(VENDOR_PKG) $(PLATFORM)

all: $(TARGET_EXEC)

$(TARGET_EXEC): $(VENDOR_PKG) $(PLATFORM) $(OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

$(PLATFORM):
	$(MAKE) -C $(PLATFORM_DIR)

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

install: $(ARTIFACTS)
	mkdir -p $(INSTALL_DIR)/usr/local/bin
	cp -vrf $(ARTIFACTS) $(INSTALL_DIR)/usr/local/bin

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES)

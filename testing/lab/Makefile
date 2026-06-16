# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

include ../../nodes/ukamaOS/config.mk

TARGET_EXEC = ukama-lab

UKAMAOS_ROOT = ../../nodes/ukamaOS/
NODES_DIR    = ../../nodes/
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

BIN_DIR = bin
BIN     = $(BIN_DIR)/$(TARGET_EXEC)

# Packages needed to compile target
ifeq ($(ALPINE_BUILD),1)
    $(info Excluding vendor library for Alpine build)
        VENDOR_PKG=
else
        VENDOR_PKG = tomlc jansson ulfius
endif
PLATFORM = platform

# Compile flags
CFLAGS = -c
CFLAGS += -g
CFLAGS += -Wall
CFLAGS += -Wextra
CFLAGS += -Werror
CFLAGS += -Wdeclaration-after-statement
CFLAGS += -O0
CFLAGS += -std=gnu11
CFLAGS += -D_POSIX_C_SOURCE=200809L
CFLAGS += -D_REENTRANT

# Local app headers first.
# This avoids collisions with Ukama platform headers like log.h.
CFLAGS += -I.
CFLAGS += -I./inc
CFLAGS += -I./inc/common/
CFLAGS += -I../common/inc

# External / Ukama platform headers after local headers.
CFLAGS += -I$(VENDOR_INC)
CFLAGS += -I$(PLATFORM_INC_SYS)
CFLAGS += -I$(PLATFORM_INC_LOG)

# Link flags
LDFLAGS += -L$(VENDOR_LIB)
LDFLAGS += -L$(VENDOR_LIB64)
LDFLAGS += -L$(PLATFORM_LIB)

# Libraries needed
CURL_LIB = -lcurl -lssl -lcrypto
ULFIUS_LIB = -lulfius -lorcania -lyder -lmicrohttpd -lgnutls -lnettle
ULFIUS_LIB += -lhogweed -lp11-kit -lz

LIBS = $(ULFIUS_LIB) $(CURL_LIB)
LIBS += -ljansson
LIBS += -lpthread
LIBS += -lusys

# Source files
TARGET_CFILES = $(wildcard ./src/*.c)
TARGET_CFILES += ./cmd/ukama-lab/main.c

TARGET_OBJFILES = $(TARGET_CFILES:.c=.o)

.PHONY: all clean check fmt $(TARGET_EXEC) $(PLATFORM)

all: version.h $(TARGET_EXEC)

$(TARGET_EXEC): $(VENDOR_PKG) $(PLATFORM) $(TARGET_OBJFILES)
	mkdir -p $(BIN_DIR)
	$(CC) -o $(BIN) $(TARGET_OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

$(PLATFORM):
	$(MAKE) -C $(PLATFORM_DIR)

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@

version.h:
	$(NODES_DIR)/utils/scripts/generate_version.sh

check: all
	$(BIN) list-checks >/dev/null
	$(BIN) list-events >/dev/null
	$(BIN) dry-run scenarios/smoke/usage-accumulation.yaml --print-world

fmt:
	clang-format -i inc/*.h src/*.c cmd/ukama-lab/main.c

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -f $(TARGET_OBJFILES)
	rm -rf $(BIN_DIR)
	rm -f version.h

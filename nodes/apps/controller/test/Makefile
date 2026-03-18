# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

include ../../ukamaOS/config.mk

CC          ?= cc
PYTHON      ?= python3
TARGET      = controllerd

CONTROLLER_HOST ?= 127.0.0.1
CONTROLLER_PORT ?= 18021

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

ifeq ($(ALPINE_BUILD),1)
    $(info Excluding vendor library for Alpine build)
    VENDOR_PKG =
else
    VENDOR_PKG = tomlc jansson ulfius libuuid
endif
PLATFORM = platform

CFLAGS  = -c
CFLAGS += -g
CFLAGS += -Wall
CFLAGS += -Wextra
CFLAGS += -O0
CFLAGS += -D_REENTRANT
CFLAGS += -I.
CFLAGS += -I./inc
CFLAGS += -I../common/inc
CFLAGS += -I$(VENDOR_INC)
CFLAGS += -I$(PLATFORM_INC_SYS)
CFLAGS += -I$(PLATFORM_INC_LOG)

LDFLAGS += $(RPATH_FLAGS)
LDFLAGS += -L$(VENDOR_LIB)
LDFLAGS += -L$(VENDOR_LIB64)
LDFLAGS += -L$(PLATFORM_LIB)

CURL_LIB   = -lcurl -lssl -lcrypto
ULFIUS_LIB = -lulfius -lorcania -lmicrohttpd -lgnutls -lnettle
ULFIUS_LIB += -lhogweed -lp11-kit -lz
ifneq ($(ALPINE_BUILD),1)
ULFIUS_LIB += -lyder
endif

LIBS  = $(ULFIUS_LIB) $(CURL_LIB)
LIBS += -ljansson
LIBS += -lpthread
LIBS += -lm
LIBS += -lusys

SRCS = $(wildcard src/*.c)
OBJS = $(SRCS:.c=.o)

.PHONY: all clean check sim-sunny sim-night sim-all api-test $(PLATFORM) $(VENDOR_PKG)

all: $(TARGET)

$(TARGET): $(VENDOR_PKG) $(PLATFORM) $(OBJS) $(COMM_OBJFILES)
	$(CC) -o $@ $(OBJS) $(COMM_OBJFILES) $(LDFLAGS) $(LIBS)
	@echo ""
	@echo "Build complete: ./$(TARGET)"
	@echo "Run daemon: CONTROLLER_SERIAL_PORT=/dev/ttyXXX CONTROLLER_DRIVER=victron ./$(TARGET)"
	@echo "API test:   make -f test/Makefile.test api-test"
	@echo ""

$(PLATFORM):
	$(MAKE) -C $(PLATFORM_DIR)

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

check:
	$(PYTHON) -m py_compile test/vedirect_sim.py
	bash -n test/test_api.sh

sim-sunny:
	cd test && $(PYTHON) vedirect_sim.py --scenario sunny

sim-night:
	cd test && $(PYTHON) vedirect_sim.py --scenario night

sim-all:
	cd test && $(PYTHON) vedirect_sim.py --scenario all

api-test:
	cd test && CONTROLLER_HOST=$(CONTROLLER_HOST) \
CONTROLLER_PORT=$(CONTROLLER_PORT) ./test_api.sh

clean:
	rm -f $(TARGET) $(OBJS)
	rm -f test/*.pyc
	rm -rf test/__pycache__

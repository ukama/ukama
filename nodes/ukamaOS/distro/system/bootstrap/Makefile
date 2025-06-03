# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2022-present, Ukama Inc.

include ../../../config.mk

TARGET_EXEC = bootstrap
TEST_EXEC   = test_runner

VERSION_FILE = version.h
ARTIFACTS    = $(TARGET_EXEC)

APP_NAME = $(CUR_BUILD_DIRNAME)
APP_CONFIG = $(APP_CONFIG_DIR)/$(APP_NAME)/config.toml
TARGET_CONFIG_DIR = $(INSTALL_DIR)$(NODE_APP_CONFIG_DIR)/$(APP_NAME)

VENDOR_DIR   = $(UKAMAOS_ROOT)/distro/vendor
VENDOR_BUILD = $(VENDOR_DIR)/build
VENDOR_INC   = $(VENDOR_BUILD)/include
VENDOR_LIB   = $(VENDOR_BUILD)/lib
VENDOR_LIB64 = $(VENDOR_BUILD)/lib64

PLATFORM_DIR     = $(UKAMAOS_ROOT)/distro/platform
PLATFORM_BUILD   = $(PLATFORM_DIR)/build
PLATFORM_INC_SYS = $(PLATFORM_DIR)/sys/inc
PLATFORM_INC_LOG = $(PLATFORM_DIR)/log/inc
PLATFORM_LIB     = $(PLATFORM_DIR)/build

# Packages needed to compile target
ifeq ($(ALPINE_BUILD),1)
    $(info Excluding vendor library for Alpine build)
        VENDOR_PKG=
else
	VENDOR_PKG = tomlc jansson ulfius libuuid
endif
PLATFORM   = platform

# Setting up various compile and link flags
CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I.
CFLAGS+=-I$(VENDOR_INC)
CFLAGS+=-I./inc
CFLAGS+=-I./inc/common
CFLAGS+=-I../common/inc
CFLAGS+=-D_REENTRANT
CFLAGS+=-I$(PLATFORM_INC_SYS)
CFLAGS+=-I$(PLATFORM_INC_LOG)

# Coverage flags
COVFLAGS=-fprofile-arcs -ftest-coverage

LDFLAGS+=-L${VENDOR_LIB} 
LDFLAGS+=-L$(VENDOR_LIB64)
LDFLAGS+=-L$(PLATFORM_LIB)
LDFLAGS+=--coverage

LDFLAGS+=-Wl,-rpath,$(PLATFORM_LIB)
LDFLAGS+=-Wl,-rpath,$(VENDOR_LIB)
LDFLAGS+=-Wl,-rpath,$(VENDOR_LIB64)

# Libraries
ULFIUS_LIB_DEP=-lorcania
ifneq ($(ALPINE_BUILD),1)
ULFIUS_LIB_DEP+=-lyder
endif
ULFIUS_LIB_DEP+=-lmicrohttpd
ULFIUS_LIB_DEP+=-lgnutls
ULFIUS_LIB_DEP+=-lnettle
ULFIUS_LIB_DEP+=-lhogweed
ULFIUS_LIB_DEP+=-lp11-kit
ULFIUS_LIB_DEP+=-lz
CURL_LIB_DEP=-lcurl
CURL_LIB_DEP+=-lssl
CURL_LIB_DEP+=-lcrypto

LIBS+=-lpthread
LIBS+=-lrt
LIBS+=-lm
LIBS+=$(ULFIUS_LIB_DEP)
LIBS+=$(CURL_LIB_DEP)
LIBS+=-ltoml
LIBS+=-ljansson
LIBS+=-lpthread
LIBS+=-lusys
LIBS+=-lmicrohttpd
LIBS+=-lulfius

# Source files
TARGET_CFILES   = $(wildcard ./src/*.c)
TARGET_OBJFILES = $(TARGET_CFILES:.c=.o)

# skip main.c
TEST_TARGET_CFILES = $(filter-out ./src/main.c, $(wildcard ./src/*.c))
TEST_CFILES    = $(wildcard ./test/test_*.c)
UNITY_CFILES   = $(VENDOR_DIR)/Unity/src/unity.c

TEST_TARGET_OBJFILES = $(TEST_TARGET_CFILES:.c=.o)
UNITY_OBJFILES = $(UNITY_CFILES:.c=.o)
TEST_OBJFILES  = $(TEST_CFILES:.c=.o)

.PHONY: all clean test coverage

all: version.h $(TARGET_EXEC)

$(TARGET_EXEC): $(VENDOR_PKG) $(PLATFORM) $(TARGET_OBJFILES) $(COMM_OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(TARGET_OBJFILES) $(COMM_OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@

$(PLATFORM):
	$(MAKE) -C $(PLATFORM_DIR)

$(VERSION_FILE):
	$(NODES_DIR)/utils/scripts/generate_version.sh

%.o: %.c
	$(CC) $(CFLAGS) $(COVFLAGS) -I$(VENDOR_DIR)/Unity/src -c $< -o $@

test: version.h $(TEST_OBJFILES) $(TEST_TARGET_OBJFILES) $(UNITY_OBJFILES)
	$(CC) -o $(TEST_EXEC) $(TEST_OBJFILES) $(TEST_TARGET_OBJFILES) $(UNITY_OBJFILES) $(LDFLAGS) $(LIBS)
	./$(TEST_EXEC)

coverage: test
	gcov $(TARGET_CFILES)

clean:
	rm -f $(TARGET_EXEC) $(TARGET_OBJFILES) $(TEST_EXEC) $(TEST_OBJFILES) \
	version.h *.gcda *.gcno *.gcov src/*.gcda src/*.gcov src/*.gcno \
	test/*.gcda test/*.gcno test/*.gcov

install: $(ARTIFACTS) $(VERSION_FILE)
	mkdir -p $(TARGET_CONFIG_DIR)
	cp -vrf $(ARTIFACTS) $(INSTALL_DIR)$(NODE_APP_DIR)
	cp -vrf $(VERSION_FILE) $(TARGET_CONFIG_DIR)
	if [ "$(APP_CONFIG)" != "" ]; then \
		cp -vrf $(APP_CONFIG) $(TARGET_CONFIG_DIR);\
	fi

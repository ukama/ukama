# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2021-present, Ukama Inc.

include ../../../config.mk

TARGET_EXEC = mesh.d
TEST_EXEC   = test_runner

VERSION_FILE = version.h
ARTIFACTS = $(TARGET_EXEC)

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
	VENDOR_PKG = tomlc jansson ulfius libuuid amqp
endif
PLATFORM = platform

# Setting up various compile and link flags
CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I.
CFLAGS+=-I$(VENDOR_INC)
CFLAGS+=-I$(PLATFORM_INC_SYS)
CFLAGS+=-I$(PLATFORM_INC_LOG)
CFLAGS+=-I./inc
CFLAGS+=-I./inc/common/
CFLAGS+=-I../common/inc
CFLAGS+=-D_REENTRANT

# Coverage flags
COVFLAGS=-fprofile-arcs -ftest-coverage

LDFLAGS+=-L${VENDOR_LIB} 
LDFLAGS+=-L$(VENDOR_LIB64)
LDFLAGS+=-L$(PLATFORM_LIB)
LDFLAGS+=--coverage

LDFLAGS+=-Wl,-rpath,$(PLATFORM_LIB)
LDFLAGS+=-Wl,-rpath,$(VENDOR_LIB)
LDFLAGS+=-Wl,-rpath,$(VENDOR_LIB64)

# Libraries needed
CURL_LIB=-lcurl -lssl -lcrypto
ULFIUS_LIB=-lulfius -lorcania -lmicrohttpd -lgnutls -lnettle
ULFIUS_LIB+=-lhogweed -lp11-kit -lz
ifneq ($(ALPINE_BUILD),1)
ULFIUS_LIB+=-lyder
endif
LIBS = $(ULFIUS_LIB) $(CURL_LIB)
LIBS+=-ltoml
LIBS+=-ljansson
LIBS+=-lpthread
LIBS+=-luuid
LIBS+=-lusys

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

# Conditionally add -DUNIT_TEST to CFLAGS if the test target is being built
ifneq (,$(findstring test,$(MAKECMDGOALS)))
CFLAGS += -DUNIT_TEST
endif

all: version.h static.h $(TARGET_EXEC)

$(TARGET_EXEC): $(VENDOR_PKG) $(PLATFORM) $(TARGET_OBJFILES) $(COMM_OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(TARGET_OBJFILES) $(COMM_OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@

$(PLATFORM):
	$(MAKE) -C $(PLATFORM_DIR)

$(VERSION_FILE):
	$(NODES_DIR)/utils/scripts/generate_version.sh

static.h:
	$(NODES_DIR)/utils/scripts/generate_static.sh

%.o: %.c
	$(CC) $(CFLAGS) $(COVFLAGS) -I$(VENDOR_DIR)/Unity/src -c $< -o $@

test: version.h static.h $(TEST_OBJFILES) $(TEST_TARGET_OBJFILES) $(UNITY_OBJFILES)
	$(CC) -o $(TEST_EXEC) $(TEST_OBJFILES) $(TEST_TARGET_OBJFILES) $(UNITY_OBJFILES) $(LDFLAGS) $(LIBS)
	./$(TEST_EXEC)

coverage: test
	gcov $(TARGET_CFILES)

install: $(ARTIFACTS) $(VERSION_FILE)
	mkdir -p $(TARGET_CONFIG_DIR)
	cp -vrf $(ARTIFACTS) $(INSTALL_DIR)$(NODE_APP_DIR)
	cp -vrf $(VERSION_FILE) $(TARGET_CONFIG_DIR)
	if [ "$(APP_CONFIG)" != "" ]; then \
		cp -vrf $(APP_CONFIG) $(TARGET_CONFIG_DIR);\
	fi

clean:
	rm -f $(TARGET_EXEC) $(TARGET_OBJFILES) $(TEST_EXEC) $(TEST_OBJFILES) \
	version.h static.h *.gcda *.gcno *.gcov src/*.gcda src/*.gcov src/*.gcno \
	test/*.gcda test/*.gcno test/*.gcov

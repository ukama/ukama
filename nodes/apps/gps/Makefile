# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

include ../../ukamaOS/config.mk

TARGET_EXEC = gps.d
TEST_EXEC   = test_runner

UKAMAOS_ROOT = $(NODES_DIR)/ukamaOS

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
	VENDOR_PKG = tomlc jansson ulfius
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

# Libraries needed
CURL_LIB=-lcurl -lssl -lcrypto
ULFIUS_LIB=-lulfius -lorcania -lyder -lmicrohttpd -lgnutls -lnettle
ULFIUS_LIB+=-lhogweed -lp11-kit -lz
LIBS = $(ULFIUS_LIB) $(CURL_LIB)
LIBS+=-ljansson
LIBS+=-lpthread
LIBS+=-lusys

# Soruce files
TARGET_CFILES   = $(wildcard ./src/*.c)
TARGET_OBJFILES = $(TARGET_CFILES:.c=.o)

# skip main.c
TEST_TARGET_CFILES = $(filter-out ./src/main.c, $(wildcard ./src/*.c))
TEST_CFILES    = $(wildcard ./test/test_*.c)
UNITY_CFILES   = $(VENDOR_DIR)/Unity/src/unity.c

TEST_TARGET_OBJFILES = $(TEST_TARGET_CFILES:.c=.o)
UNITY_OBJFILES = $(UNITY_CFILES:.c=.o)
TEST_OBJFILES  = $(TEST_CFILES:.c=.o)

.PHONY: all clean test $(TARGET_EXEC)

# Conditionally add -DUNIT_TEST to CFLAGS if the test target is being built
ifneq (,$(findstring test,$(MAKECMDGOALS)))
CFLAGS += -DUNIT_TEST
endif

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

test: version.h static.h $(TEST_OBJFILES) $(TEST_TARGET_OBJFILES) $(UNITY_OBJFILES)
	$(CC) -o $(TEST_EXEC) $(TEST_OBJFILES) $(TEST_TARGET_OBJFILES) \
		$(UNITY_OBJFILES) $(LDFLAGS) $(LIBS)
	./$(TEST_EXEC)

coverage: test
	gcov $(TARGET_CFILES)

%.o: %.c
	$(CC) $(CFLAGS) $(COVFLAGS) -I$(VENDOR_DIR)/Unity/src -c $< -o $@

clean:
	rm -f $(TARGET_EXEC) $(TARGET_OBJFILES) $(TEST_EXEC) $(TEST_OBJFILES) \
		version.h static.h *.gcda *.gcno *.gcov src/*.gcda src/*.gcov src/*.gcno \
		test/*.gcda test/*.gcno test/*.gcov

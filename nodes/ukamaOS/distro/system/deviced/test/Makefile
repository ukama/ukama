# Unit tests for device.d using off-the-shelf Unity (Unity.c/Unity.h).
#
# Usage:
#   make -C test UNITY_DIR=/path/to/Unity run
#
# UNITY_DIR must contain Unity.c and Unity.h (common layout: Unity/src/Unity.c, Unity/src/Unity.h)

include ../../../../config.mk

UNITY_DIR ?= ../../../vendor/Unity/
UNITY_SRC ?= $(UNITY_DIR)/src

CC ?= gcc

TEST_BIN := deviced_tests
PRELOAD_SO := libdeviced_testwrap.so

# Try to reuse build flags from parent when running inside Ukama tree.
# The parent Makefile includes ../../../config.mk. We don't replicate that here.

CFLAGS += -g -O0 -Wall
CFLAGS += -I../inc -I../inc/common -I../src -I./mocksrv
CFLAGS += -I$(UNITY_SRC)
CFLAGS += -I$(VENDOR_INC)
CFLAGS += -I$(PLATFORM_INC_SYS)
CFLAGS += -I$(PLATFORM_INC_LOG)

LDFLAGS += -L.
LDFLAGS += $(RPATH_FLAGS)
LDFLAGS += -L${VENDOR_LIB}
LDFLAGS += -L$(VENDOR_LIB64)
LDFLAGS += -L$(PLATFORM_LIB)

LIBS += -lulfius -lorcania -lmicrohttpd -lgnutls -lnettle -lhogweed -lp11-kit -lz
LIBS += -lcurl -lssl -lcrypto
LIBS += -ljansson -lpthread -lusys

SRCS := \
	./test_deviced.c \
	./mocksrv/mock_services.c \
	$(UNITY_SRC)/unity.c

.PHONY: all run clean build_deviced preload

all: run

build_deviced:
	$(MAKE) -C .. clean all

preload: $(PRELOAD_SO)

$(PRELOAD_SO): ./preload/deviced_testwrap.c
	$(CC) -shared -fPIC -o $@ $<

$(TEST_BIN): $(SRCS) preload
	@if [ -z "$(UNITY_DIR)" ]; then \
		echo "ERROR: UNITY_DIR is required (path to Unity repo)."; \
		exit 1; \
	fi
	$(CC) $(CFLAGS) -o $@ $(SRCS) $(LDFLAGS) $(LIBS)

run: build_deviced $(TEST_BIN)
	DEVICED_TEST_PRELOAD=$(CURDIR)/$(PRELOAD_SO) \
	DEVICED_TEST_BIN=$(abspath ../device.d) \
	./$(TEST_BIN)

clean:
	rm -f $(TEST_BIN) $(PRELOAD_SO)

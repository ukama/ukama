# Local controller.d build for testing with the VE.Direct simulator.
# Prefers pkg-config on both macOS and Linux; falls back to Homebrew paths on macOS.
#
# Usage:
#   cd nodes/apps/controller
#   make -f test/Makefile.test
#   CONTROLLER_SERIAL_PORT=/dev/ttys027 ./controllerd

CC      ?= cc
PYTHON  ?= python3
PKG_CONFIG ?= pkg-config
BREW    ?= brew
TARGET  = controllerd

UNAME_S := $(shell uname -s)
HB      ?= $(shell $(BREW) --prefix 2>/dev/null || echo /opt/homebrew)

PLATFORM_ROOT = ../../ukamaOS/distro/platform
CONTROLLER_HOST ?= 127.0.0.1
CONTROLLER_PORT ?= 8095

PKG_DEPS = libulfius jansson libcurl
PKG_CONFIG_READY := $(shell command -v $(PKG_CONFIG) >/dev/null 2>&1 && $(PKG_CONFIG) --exists $(PKG_DEPS) >/dev/null 2>&1 && echo yes || echo no)

CFLAGS  = -g -Wall -Wextra -O0 -D_REENTRANT
CFLAGS += -I.
CFLAGS += -I./inc
CFLAGS += -I../common/inc
CFLAGS += -I$(PLATFORM_ROOT)/sys/inc
CFLAGS += -I$(PLATFORM_ROOT)/log/inc

LDFLAGS =
LIBS    =

ifeq ($(PKG_CONFIG_READY),yes)
	CFLAGS += $(shell $(PKG_CONFIG) --cflags $(PKG_DEPS))
	LIBS   += $(shell $(PKG_CONFIG) --libs $(PKG_DEPS))
else ifeq ($(UNAME_S),Darwin)
	CFLAGS += -I$(HB)/include
	LDFLAGS += -L$(HB)/lib
	LIBS   += -lulfius -lorcania -lyder -lmicrohttpd -lgnutls -lnettle -lhogweed -lp11-kit -lz
	LIBS   += -ljansson
	LIBS   += -lcurl
else
	LIBS   += -lulfius -lorcania -ljansson -lcurl
endif

ifeq ($(UNAME_S),Darwin)
	CFLAGS += -I$(HB)/include
	LDFLAGS += -L$(HB)/lib
	CFLAGS += -I./test/compat
	CFLAGS += -Dpthread_spinlock_t=pthread_mutex_t
endif

LIBS   += -lpthread
LIBS   += -lm

SRCS    = $(wildcard src/*.c)
SRCS   += test/usys_stubs.c

OBJS    = $(SRCS:.c=.o)

.PHONY: all clean check sim-sunny sim-night sim-all api-test

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) -o $@ $^ $(LDFLAGS) $(LIBS)
	@echo ""
	@echo "  Build complete: ./$(TARGET)"
	@echo ""
	@echo "  Host OS:           $(UNAME_S)"
	@echo "  pkg-config:        $(PKG_CONFIG_READY)"
	@echo "  Script checks:     make -f test/Makefile.test check"
	@echo "  Start simulator:   make -f test/Makefile.test sim-sunny"
	@echo "  Run daemon:        CONTROLLER_SERIAL_PORT=/dev/ttys0XX CONTROLLER_DRIVER=victron ./$(TARGET)"
	@echo "  Test API:          make -f test/Makefile.test api-test"
	@echo ""

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
	cd test && CONTROLLER_HOST=$(CONTROLLER_HOST) CONTROLLER_PORT=$(CONTROLLER_PORT) ./test_api.sh

clean:
	rm -f $(TARGET) $(OBJS)
	rm -f test/*.pyc
	rm -rf test/__pycache__

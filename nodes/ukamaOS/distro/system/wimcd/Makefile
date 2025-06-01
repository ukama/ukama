# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2021-present, Ukama Inc.

include ../../../config.mk

TARGET_EXEC  = wimc.d
AGENT_EXEC   = agent
VERSION_FILE = version.h

ARTIFACTS = $(TARGET_EXEC) $(AGENT_EXEC)

APP_NAME = $(CUR_BUILD_DIRNAME)
APP_CONFIG =
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

# Packages needed to compile target.
ifeq ($(ALPINE_BUILD),1)
    $(info Excluding vendor library for Alpine build)
        VENDOR_PKG=
else
	VENDOR_PKG = tomlc ulfius jansson
endif
PLATFORM = platform

# Setting up various compile and link flags.
CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I.
CFLAGS+=-I$(VENDOR_INC)
CFLAGS+=-I$(PLATFORM_INC_SYS)
CFLAGS+=-I$(PLATFORM_INC_LOG)
CFLAGS+=-I./inc
CFLAGS+=-I./inc/common
CFLAGS+=-I../common/inc
CFLAGS+=-D_REENTRANT

LDFLAGS+=-L${VENDOR_LIB} 
LDFLAGS+=-L$(VENDOR_LIB64)
LDFLAGS+=-L$(PLATFORM_LIB)

LDFLAGS+=-Wl,-rpath,$(PLATFORM_LIB)
LDFLAGS+=-Wl,-rpath,$(VENDOR_LIB)
LDFLAGS+=-Wl,-rpath,$(VENDOR_LIB64)

# Libraries needed.
ULFIUS_LIB=-lulfius -lorcania -lmicrohttpd -lgnutls -lnettle
ULFIUS_LIB+=-lhogweed -lp11-kit -lz
ifneq ($(ALPINE_BUILD),1)
ULFIUS_LIB+=-lyder
endif
CURL_LIB=-lcurl -lssl -lcrypto

LIBS = $(ULFIUS_LIB) $(CURL_LIB)
LIBS+=-ltoml
LIBS+=-lpthread
LIBS+=-ljansson
LIBS+=-lrt
LIBS+=-luuid
LIBS+=-lsqlite3
LIBS+=-lusys

# Soruce files.
CFILES   = $(wildcard ./src/*.c)
OBJFILES = $(CFILES:.c=.o)

COMM_CFILES = $(wildcard ./common/*.c)
COMM_OBJFILES = $(COMM_CFILES:.c=.o)

AGENT_CFILES = $(wildcard ./agents/*.c)
AGENT_OBJFILES = $(AGENT_CFILES:.c=.o)

.PHONY: $(TARGET_EXEC) $(AGENT_EXEC)

all: version.h wimc agent

wimc: $(TARGET_EXEC)

agent: $(AGENT_EXEC)

$(TARGET_EXEC): $(VENDOR_PKG) $(PLATFORM) $(OBJFILES) $(COMM_OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(OBJFILES) $(COMM_OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

$(AGENT_EXEC): $(VENDOR_PKG) $(AGENT_OBJFILES) $(COMM_OBJFILES)
	$(CC) -o $(AGENT_EXEC) $(AGENT_OBJFILES) $(COMM_OBJFILES) \
		$(LDFLAGS) $(LIBS)
	echo "Done."

$(VENDOR_PKG):
	cd $(VENDOR_DIR) && make $@

$(PLATFORM):
	$(MAKE) -C $(PLATFORM_DIR)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

$(VERSION_FILE):
	$(NODES_DIR)/utils/scripts/generate_version.sh

clean:
	rm -f $(TARGET_EXEC) $(AGENT_EXEC) $(OBJFILES) $(COMM_OBJFILES) version.h
	rm -f $(AGENT_OBJFILES)

install: $(ARTIFACTS) $(VERSION_FILE)
	mkdir -p $(TARGET_CONFIG_DIR)
	cp -vrf $(ARTIFACTS) $(INSTALL_DIR)$(NODE_APP_DIR)
	cp -vrf $(VERSION_FILE) $(TARGET_CONFIG_DIR)
	if [ "$(APP_CONFIG)" != "" ]; then \
		cp -vrf $(APP_CONFIG) $(TARGET_CONFIG_DIR);\
	fi


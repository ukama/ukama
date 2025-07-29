# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

include ../../../config.mk

TARGET_EXEC = lookout.d

VERSION_FILE = version.h
ARTIFACTS = $(TARGET_EXEC)

APP_NAME = $(CUR_BUILD_DIRNAME)
TARGET_CONFIG_DIR = $(INSTALL_DIR)$(NODE_APP_CONFIG_DIR)/$(APP_NAME)
APP_CONFIG = 

# Packages needed to compile target
ifeq ($(ALPINE_BUILD),1)
    $(info Excluding vendor library for Alpine build)
        VENDOR_PKG=
else
	VENDOR_PKG = tomlc jansson ulfius libuuid
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
CFLAGS+=-I./inc/common
CFLAGS+=-I../common/inc
CFLAGS+=-D_REENTRANT

LDFLAGS+=$(RPATH_FLAGS)
LDFLAGS+=-L${VENDOR_LIB} 
LDFLAGS+=-L$(VENDOR_LIB64)
LDFLAGS+=-L$(PLATFORM_LIB)

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
LIBS+=-lusys

# Soruce files
CFILES   = $(wildcard ./src/*.c)
OBJFILES = $(CFILES:.c=.o)

.PHONY: $(TARGET_EXEC)

all: version.h lookout

lookout: $(TARGET_EXEC)

$(TARGET_EXEC): $(VENDOR_PKG) $(PLATFORM) $(OBJFILES) $(COMM_OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(OBJFILES) $(COMM_OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

$(PLATFORM):
	$(MAKE) -C $(PLATFORM_DIR)

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@ 

$(VERSION_FILE):
	$(NODES_DIR)/utils/scripts/generate_version.sh

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

install: $(ARTIFACTS) $(VERSION_FILE)
	mkdir -p $(TARGET_CONFIG_DIR)
	cp -vrf $(ARTIFACTS) $(INSTALL_DIR)$(NODE_APP_DIR)
	cp -vrf $(VERSION_FILE) $(TARGET_CONFIG_DIR)
	if [ "$(APP_CONFIG)" != "" ]; then \
		cp -vrf $(APP_CONFIG) $(TARGET_CONFIG_DIR);\
	fi

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES) version.h

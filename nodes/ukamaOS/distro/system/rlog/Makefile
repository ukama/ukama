# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

include ../../../config.mk

TARGET_EXEC = rlog.d

VERSION_FILE = version.h
ARTIFACTS = $(TARGET_EXEC)

APP_NAME = $(CUR_BUILD_DIRNAME)
APP_CONFIG =
TARGET_CONFIG_DIR = $(INSTALL_DIR)$(NODE_APP_CONFIG_DIR)/$(APP_NAME)

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
CFLAGS+=-I./inc
CFLAGS+=-I./inc/common/
CFLAGS+=-I../common/inc
CFLAGS+=-D_REENTRANT
CFLAGS+=-I$(PLATFORM_INC_SYS)
CFLAGS+=-I$(PLATFORM_INC_LOG)

LDFLAGS+=-L${VENDOR_LIB} 
LDFLAGS+=-L$(VENDOR_LIB64)
LDFLAGS+=-L$(PLATFORM_LIB)
LDFLAGS+=$(RPATH_FLAGS)

# Libraries needed
CURL_LIB=-lcurl -lssl -lcrypto -lgnutls 
ULFIUS_LIB=-lulfius -lorcania -lmicrohttpd -lgnutls -lnettle
ULFIUS_LIB+=-lhogweed -lp11-kit -lz
ifneq ($(ALPINE_BUILD),1)
ULFIUS_LIB+=-lyder
endif
LIBS =$(ULFIUS_LIB)
LIBS+=$(CURL_LIB)
LIBS+=-ltoml
LIBS+=-ljansson
LIBS+=-lpthread
LIBS+=-lusys

# Soruce files
CFILES   = $(wildcard ./src/*.c)
OBJFILES = $(CFILES:.c=.o)

.PHONY: $(TARGET_EXEC) $(AGENT_EXEC)

all: version.h rlog.d
#	cd test; make ALPINE_BUILD=$(ALPINE_BUILD) all; cd -

rlog.d: $(TARGET_EXEC)

$(TARGET_EXEC): $(VENDOR_PKG) $(PLATFORM) $(OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(OBJFILES) $(COMM_OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@

$(PLATFORM):
	$(MAKE) -C $(PLATFORM_DIR)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

$(VERSION_FILE):
	$(NODES_DIR)/utils/scripts/generate_version.sh

install: $(ARTIFACTS) $(VERSION)
	mkdir -p $(TARGET_CONFIG_DIR)
	cp -vrf $(ARTIFACTS) $(INSTALL_DIR)$(NODE_APP_DIR)
	cp -vrf $(VERSION_FILE) $(TARGET_CONFIG_DIR)
	if [ "$(APP_CONFIG)" != "" ]; then \
		cp -vrf $(APP_CONFIG) $(TARGET_CONFIG_DIR);\
	fi

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES) version.h
	cd test; make clean; cd -

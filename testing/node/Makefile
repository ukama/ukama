# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

# some generic defs.
CC   = gcc
ARCH = X86_64
XGCC = gcc
XLD  = ld
XGXX = g++
HOST = $(shell gcc -dumpmachine)

VERSION=$(shell git describe --always --dirty=-dirty)
REPO_ROOT := $(abspath $(CURDIR)/../..)
STAGE_DIR := /tmp/virtnode
PKG_NAME  := ukama_${VERSION}.tgz

#VERSION = v0.0.1

OPENSSLTARGET = linux-generic64

TARGET_EXEC = virtualNode

# Setting up various compile and link flags
CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I./inc -I../../nodes/ukamaOS/distro/vendor/build/include
CFLAGS+=-D_REENTRANT

# Libs
LIBS=-ltoml
LIBS+=-ljansson

LDFLAGS+=-L../../nodes/ukamaOS/distro/vendor/build/lib
LDFLAGS+=-L../../nodes/ukamaOS/distro/vendor/build/lib64

# Soruce files
CFILES   = $(wildcard ./src/*.c)
OBJFILES = $(CFILES:.c=.o)

.PHONY: $(TARGET_EXEC)

srvc_router: $(TARGET_EXEC)

$(TARGET_EXEC): $(OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(OBJFILES) $(COMM_OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

sourcetgz: $(TARGET_EXEC)
	rm -rf $(STAGE_DIR)
	rm -f ukama_*.tgz

	( set -e; \
	  mkdir -p $(STAGE_DIR); \
	  rsync -a --info=progress2  \
	    --exclude='**/build' \
	    --exclude='**/out' \
	    --exclude='**/.cache' \
	    --exclude='**/musl-cross-make' \
	    --exclude='**/*.o' \
	    --exclude='**/*.a' \
	    --exclude='**/*.so' \
	    $(REPO_ROOT)/ \
	    $(STAGE_DIR)/ukama; \
	  cd $(STAGE_DIR) && \
	  tar -cf $(PKG_NAME) ukama; \
	  mv $(PKG_NAME) $(CURDIR)/$(PKG_NAME) \
	)

container:  sourcetgz
	@echo Building container
	podman build -t testing/virtualnode:${VERSION} .

all: $(TARGET_EXEC)

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES)

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

UKAMAOS_ROOT := ../../nodes/ukamaOS

VENDOR_DIR   := $(UKAMAOS_ROOT)/distro/vendor
VENDOR_BUILD := $(VENDOR_DIR)/build
VENDOR_INC   := $(VENDOR_BUILD)/include
VENDOR_LIB   := $(VENDOR_BUILD)/lib
VENDOR_LIB64 := $(VENDOR_BUILD)/lib64

PLATFORM_DIR := $(UKAMAOS_ROOT)/distro/platform

# Packages to build (targets in vendor Makefile)
VENDOR_PKG := tomlc jansson ulfius
PLATFORM   := platform

CC := gcc

VERSION     := $(shell git describe --always --dirty=-dirty)
REPO_ROOT   := $(abspath $(CURDIR)/../..)
STAGE_DIR   := /tmp/virtnode
PKG_NAME    := ukama_$(VERSION).tgz
GITMETA_DIR := $(STAGE_DIR)/ukama/.gitmeta

TARGET_EXEC := virtualNode

CFLAGS  := -c -g -Wall -O0 -D_REENTRANT
CFLAGS  += -I./inc -I$(VENDOR_INC)

LIBS    := -ltoml -ljansson
LDFLAGS := -L$(VENDOR_LIB) -L$(VENDOR_LIB64)

# Sources
CFILES   := $(wildcard ./src/*.c)
OBJFILES := $(CFILES:.c=.o)

# Docker build context + boards staging
DOCKER_CONTEXT := $(CURDIR)
BOARDS_SRC     := $(REPO_ROOT)/nodes/configs/boards
BOARDS_STAGE   := $(DOCKER_CONTEXT)/boards

.PHONY: all clean container sourcetgz stage-boards unstage-boards \
        print-version print-local-image $(PLATFORM) $(VENDOR_PKG)

stage-boards:
	@echo "Staging boards into build context..."
	@test -d "$(BOARDS_SRC)" || (echo "ERROR: BOARDS_SRC not found: $(BOARDS_SRC)"; exit 1)
	rm -rf "$(BOARDS_STAGE)"
	mkdir -p "$(BOARDS_STAGE)"
	rsync -a --delete "$(BOARDS_SRC)/" "$(BOARDS_STAGE)/"

unstage-boards:
	@echo "Removing staged boards from build context..."
	rm -rf "$(BOARDS_STAGE)"

all: $(TARGET_EXEC)

srvc_router: $(TARGET_EXEC)

$(TARGET_EXEC): $(VENDOR_PKG) $(PLATFORM) $(OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(OBJFILES) $(LDFLAGS) $(LIBS)
	@echo "Done."

$(PLATFORM):
	$(MAKE) -C $(PLATFORM_DIR)

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@

%.o: %.c
	$(CC) $(CFLAGS) $< -o $@

sourcetgz: $(TARGET_EXEC)
	rm -rf $(STAGE_DIR)
	rm -f ukama_*.tgz

	( set -e; \
	  mkdir -p $(STAGE_DIR); \
	  rsync -a --info=progress2 \
	    --exclude='.git' \
	    --exclude='**/build' \
	    --exclude='**/out' \
	    --exclude='**/.cache' \
	    --exclude='**/musl-cross-make' \
	    --exclude='**/*.o' \
	    --exclude='**/*.a' \
	    --exclude='**/*.so' \
	    $(REPO_ROOT)/ \
	    $(STAGE_DIR)/ukama; \
	  rm -rf $(GITMETA_DIR); \
	  git clone --no-checkout --filter=blob:none --tags --depth=500 $(REPO_ROOT) $(GITMETA_DIR); \
	  printf "gitdir: .gitmeta/.git\n" > $(STAGE_DIR)/ukama/.git; \
	  cd $(STAGE_DIR) && \
	  tar -cf - ukama | pigz -1 > $(PKG_NAME); \
	  mv $(PKG_NAME) $(CURDIR)/$(PKG_NAME); \
	)

container: sourcetgz stage-boards
	@echo "Building container"
	podman build -t testing/virtualnode:$(VERSION) .
	@$(MAKE) unstage-boards

print-version:
	@echo $(VERSION)

print-local-image:
	@echo testing/virtualnode:$(VERSION)

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES)
	rm -rf $(STAGE_DIR)
	rm -f $(PKG_NAME)
	rm -rf build/
	rm -rf $(BOARDS_STAGE)

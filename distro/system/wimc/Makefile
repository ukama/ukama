#
# Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.
#
# This source code is licensed under the XXX-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.
#

#
# Makefile for WIMC.d
#

UKAMAOS=/home/kashif/work/ukama/ukamaOS
COMMON_DIR=$(UKAMAOS)/distro/system/common
PREFIX=/usr/local
ULFIUS=$(COMMON_DIR)/third_party/ulfius/src
ULFIUS_INC=$(COMMON_DIR)/third_party/ulfius/include
COMMON_INC=$(COMMON_DIR)/inc
BUILD=./build
SRC=./src
PROVIDER_SRC=./provider
SRCS=$(wildcard $(SRC)/*.c)
PROVIDER_SRCS=$(wildcard $(PROVIDER_SRC)/*.c)
OBJS=$(patsubst $(SRC)/%.c,$(BUILD)/%.o,$(SRCS))
EXE=wimc.d
PROVIDER=mock_provider

CC=gcc
CFLAGS+=-c -Wall -I$(ULFIUS_INC) -I$(COMMON_INC) -I./inc -D_REENTRANT -g -O0
LDFLAGS=-L$(ULFIUS)
LIBS=-lc -lulfius -lorcania -lsqlite3 -ljansson
PROVIDER_LIBS=-lulfius -lsqlite3

.PHONY: all run clean

all: $(EXE)

libulfius:
	cd $(ULFIUS) && $(MAKE) debug CURLFLAG=1 GNUTLSFLAG=1 JANSSONFLAG=1 WEBSOCKETFLAG=1 YDERFLAG=1 UWSCFLAG=1

$(EXE): $(BUILD) libulfius $(OBJS)
	$(CC) $(OBJS) -o $(BUILD)/$@ $(LDFLAGS) $(LIBS)

$(PROVIDER): $(BUILD)
	$(CC) -g -o $(BUILD)/$(PROVIDER) $(PROVIDER_SRCS) $(PROVIDER_LIBS)
	echo "Binary available at: $(BUILD)/$(PROVIDER)"

$(BUILD)/%.o: $(SRC)/%.c $(OBJ) $(PROVIDER_SRC)/%.c $(PROVIDER_OBJS)
	$(CC) $(CFLAGS) -c $< -o $@

$(BUILD):
	mkdir -p $@

run: $(EXE)
	$<

clean:
	rm -rf $(BUILD);
	cd $(ULFIUS) && make clean

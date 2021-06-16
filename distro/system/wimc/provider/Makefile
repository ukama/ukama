#
# Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.
#
# This source code is licensed under the XXX-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.
#

#
# Makefile for Mock Provider
#

UKAMAOS=/home/kashif/work/ukama/ukamaOS
COMMON_DIR=$(UKAMAOS)/distro/system/common
SRC_DIR=$(UKAMAOS)/distro/system/wimc/provider
ULFIUS=$(COMMON_DIR)/third_party/ulfius/src
ULFIUS_INC=$(COMMON_DIR)/third_party/ulfius/include
COMMON_INC=$(COMMON_DIR)/inc

BUILD=./build
SRC=./
SRCS=$(wildcard $(SRC)/*.c)
OBJS=$(patsubst $(SRC)/%.c,$(BUILD)/%.o,$(SRCS))
EXE=mock_provider

CC=gcc
CFLAGS+=-c -Wall -I$(ULFIUS_INC) -I$(COMMON_INC) -I./inc -D_REENTRANT -g -O0
LDFLAGS=-L$(SRC_DIR)/$(BUILD)
LIBS=-lulfius -lsqlite3 -ljansson -lorcania -lmicrohttpd -lcurl

.PHONY: all run clean

all: $(EXE)

libulfius:
	cd $(ULFIUS) && $(MAKE) static GNUTLSFLAG=1 WEBSOCKETFLAG=1 YDERFLAG=1 UWSCFLAG=1 && cp libulfius.a $(SRC_DIR)/$(BUILD)

$(EXE): $(BUILD) libulfius $(OBJS)
	$(CC) $(OBJS) -o $(BUILD)/$@ $(LDFLAGS) $(LIBS)
	echo "Binary available at: $(BUILD)/$(EXE)"

$(BUILD)/%.o: $(SRC)/%.c $(OBJ)
	$(CC) $(CFLAGS) $(LDFLAGS) -c $< -o $@

$(BUILD):
	mkdir -p $@

run: $(EXE)
	$<

clean:
	rm -rf $(BUILD);
	cd $(ULFIUS) && make clean

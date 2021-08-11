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
SRC_DIR=$(UKAMAOS)/distro/system/wimc/provider
COMMON_INC=../../common/inc
BUILD=./build
SRC=./
SRCS=$(wildcard $(SRC)/*.c)
OBJS=$(patsubst $(SRC)/%.c,$(BUILD)/%.o,$(SRCS))
EXE=mock_server

CC=gcc
CFLAGS+=-c -Wall -I../inc -I$(COMMON_INC) -D_REENTRANT -g -O0
LIBS=-lulfius

.PHONY: all run clean

all: $(EXE)

$(EXE): $(BUILD) $(OBJS)
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

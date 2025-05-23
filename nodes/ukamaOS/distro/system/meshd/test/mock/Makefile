# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2021-present, Ukama Inc.

UKAMAOS=/home/kashif/work/ukama/ukamaOS
SRC_DIR=$(UKAMAOS)/distro/system/wimc/provider
COMMON_INC=../../common/inc
BUILD=./build
SRC=./
SRCS=$(wildcard $(SRC)/*.c)
OBJS=$(patsubst $(SRC)/%.c,$(BUILD)/%.o,$(SRCS))
AMQP_INC=-../../../vendor/amqp/include
LIBRABBITMQ_INC=../../../vendor/amqp/librabbitmq/
MOCK_OBJS=mock_server.o
AMQP_OBJS=amqp_listen.o link.pb-c.o
MOCK=mock_server
AMQP=amqp_listen

CC=gcc
CFLAGS+=-c -Wall -I../inc -I$(COMMON_INC) -I$(AMQP_INC) -I$(LIBRABBITMQ_INC) \
	-D_REENTRANT -g -O0 `pkg-config --cflags 'libprotobuf-c >= 1.0.0'`
LIBS=-lulfius -lrabbitmq `pkg-config --libs 'libprotobuf-c >= 1.0.0'`

.PHONY: all run clean

all:
	echo "Define a target: $(MOCK) $(AMQP)"

$(MOCK): $(BUILD) $(MOCK_OBJS)
	$(CC) $(OBJS) -o $(BUILD)/$@ $(LDFLAGS) $(LIBS)
	echo "Binary available at: $(BUILD)/$(MOCK)"

$(AMQP): $(BUILD) $(AMQP_OBJS)
	$(CC) $(AMQP_OBJS) -o $(BUILD)/$@ $(LDFLAGS) $(LIBS)
	echo "Binary available at: $(BUILD)/$(AMQP)"
	rm *.o

$(BUILD)/%.o: $(SRC)/%.c $(OBJ)
	$(CC) $(CFLAGS) $(LDFLAGS) -c $< -o $@

$(BUILD):
	mkdir -p $@

clean:
	rm -rf $(BUILD)

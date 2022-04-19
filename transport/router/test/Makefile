#
# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.
#
# This source code is licensed under the XXX-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.
#

#
# Makefile for Mock Micro-service
#

SRC=./
SRCS=$(wildcard $(SRC)/*.c)
OBJS=$(patsubst $(SRC)/%.c,./%.o,$(SRCS))
MOCK_OBJS=mock_service.o
MOCK=mock_service

CC=gcc
CFLAGS+=-c -Wall -I../inc -D_REENTRANT -g
LIBS=-lulfius -lcurl

.PHONY: all

all: $(MOCK)

$(MOCK): $(MOCK_OBJS)
	$(CC) $(OBJS) -o $@ $(LDFLAGS) $(LIBS)
	echo "Binary available at: $(MOCK)"

%.o: $(SRC)/%.c $(OBJ)
	$(CC) $(CFLAGS) $(LDFLAGS) -c $< -o $@

clean:
	rm -rf *.o ${MOCK}

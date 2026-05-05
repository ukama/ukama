# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

include ../../../ukamaOS/config.mk

TARGET_EXEC = switchemu.d

VERSION ?= latest

CC ?= gcc
CFLAGS += -c
CFLAGS += -g
CFLAGS += -Wall
CFLAGS += -Wextra
CFLAGS += -Werror
CFLAGS += -O0
CFLAGS += -std=c11
CFLAGS += -D_POSIX_C_SOURCE=200809L
CFLAGS += -I.
CFLAGS += -I./inc
CFLAGS += -D_REENTRANT

LDFLAGS +=
LIBS += -lpthread
LIBS += -lcurl

TARGET_CFILES = $(filter-out ./src/log.c, $(wildcard ./src/*.c))
TARGET_OBJFILES = $(TARGET_CFILES:.c=.o)

.PHONY: all clean $(TARGET_EXEC)

all: version.h $(TARGET_EXEC)

$(TARGET_EXEC): $(TARGET_OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(TARGET_OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

version.h:
	$(NODES_DIR)/utils/scripts/generate_version.sh

%.o: %.c version.h
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -f $(TARGET_EXEC) $(TARGET_OBJFILES) version.h

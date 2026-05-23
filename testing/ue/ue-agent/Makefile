# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

TARGET_EXEC = ue-agent

CFLAGS = -g -Wall -Wextra -O0 -I./inc -D_REENTRANT
LIBS = -lpthread -lcurl

CFILES = $(wildcard src/*.c)
OBJFILES = $(CFILES:.c=.o)

all: $(TARGET_EXEC)

$(TARGET_EXEC): $(OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(OBJFILES) $(LIBS)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES)

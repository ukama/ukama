# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

#
# Makefile for virtual Node
#

# some generic defs.
CC   = gcc
ARCH = X86_64
XGCC = gcc
XLD  = ld
XGXX = g++
HOST = $(shell gcc -dumpmachine)

OPENSSLTARGET = linux-generic64

TARGET_EXEC = virtualNode

# Setting up various compile and link flags
CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I./inc
CFLAGS+=-D_REENTRANT

# Libs
LIBS=-ltoml

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

all: $(TARGET_EXEC)

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES)

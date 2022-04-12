# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

#
# Makefile for the srvc_router
#

# some generic defs.
CC   = gcc
ARCH = X86_64
XGCC = gcc
XLD  = ld
XGXX = g++
HOST = $(shell gcc -dumpmachine)

OPENSSLTARGET = linux-generic64

TARGET_EXEC = srvc_router

# Setting up various compile and link flags
CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I./inc
CFLAGS+=-D_REENTRANT

# Libs
LIBS=-lulfius
LIBS+=-lorcania
LIBS+=-lcurl
LIBS+=-ljansson
LIBS+=-luuid

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

all: srvc_router

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES)

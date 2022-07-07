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

#VERSION = v0.0.1

OPENSSLTARGET = linux-generic64

TARGET_EXEC = virtualNode

# Setting up various compile and link flags
CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I./inc -I../../nodes/ukamaOS/distro/vendor/build/include
CFLAGS+=-D_REENTRANT

# Libs
LIBS=-ltoml
LIBS+=-ljansson

LDFLAGS+=-L../../nodes/ukamaOS/distro/vendor/build/lib
LDFLAGS+=-L../../nodes/ukamaOS/distro/vendor/build/lib64

# Soruce files
CFILES   = $(wildcard ./src/*.c)
OBJFILES = $(CFILES:.c=.o)

#VERSION=$(shell git rev-parse --short HEAD)
VERSION=$(shell git describe --always --dirty=-dirty)

.PHONY: $(TARGET_EXEC)

srvc_router: $(TARGET_EXEC)

$(TARGET_EXEC): $(OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(OBJFILES) $(COMM_OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

	 
sourcetgz: $(TARGET_EXEC)
	rm -rf /tmp/virtnode/
	rm -rf ukamaOS_*.tgz 
	(mkdir -p /tmp/virtnode/ukamaOS/distro && \
	rsync -av --exclude='musl-cross-make' ../../nodes/ukamaOS/distro/ /tmp/virtnode/ukamaOS/distro && \
	rsync -av ../../nodes/ukamaOS/config.mk /tmp/virtnode/ukamaOS/config.mk &&\
	sync && tar -zcvf /tmp/virtnode/ukamaOS.tgz /tmp/virtnode/ukamaOS && \
	mv /tmp/virtnode/ukamaOS.tgz ukamaOS_${VERSION}.tgz)

container:  sourcetgz
	@echo Building container
	buildah bud -f . -t testing/virtualnode:${VERSION}

all: $(TARGET_EXEC)

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES)

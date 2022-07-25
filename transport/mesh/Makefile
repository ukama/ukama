# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

#
# Makefile for the Mesh (server)
#

# some generic defs.
CC   = gcc
ARCH = X86_64
XGCC = gcc
XLD  = ld
XGXX = g++

HOST = $(shell gcc -dumpmachine)

TARGET_EXEC = mesh

UKAMA_OS     = ../../nodes/ukamaOS/
VENDOR_DIR   = $(UKAMA_OS)/distro/vendor/
VENDOR_BUILD = $(VENDOR_DIR)/build/
VENDOR_INC   = $(VENDOR_BUILD)/include/
VENDOR_LIB   = $(VENDOR_BUILD)/lib/
VENDOR_LIB64 = $(VENDOR_BUILD)/lib64/

# Packages needed to compile target
VENDOR_PKG = ulfius tomlc jansson amqp protobuf-c

# Setting up various compile and link flags.
CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I$(VENDOR_INC)
CFLAGS+=-I$(VENDOR_INC)/rabbitmq-c/
CFLAGS+=-I./inc
CFLAGS+=-D_REENTRANT
CFLAGS+=-DTEST_EMBED_CERT

LDFLAGS+=-L${VENDOR_LIB} 
LDFLAGS+=-L${VENDOR_LIB64}

# Libraries needed.
ULFIUS_LIB=-lmicrohttpd -lorcania -lyder -lgnutls -lnettle
ULFIUS_LIB+=-lhogweed -lp11-kit -lz -lulfius
CURL_LIB=-lcurl -lssl -lcrypto
LIBS = $(ULFIUS_LIB) $(CURL_LIB)
LIBS+=-ljansson
LIBS+=-luuid
LIBS+=-lpthread
LIBS+=-ltoml
LIBS+=-lprotobuf-c
LIBS+=-lrabbitmq

# Soruce files.
CFILES   = $(wildcard ./src/*.c)
OBJFILES = $(CFILES:.c=.o)

.PHONY: $(TARGET_EXEC)

mesh: $(TARGET_EXEC)

$(TARGET_EXEC): $(VENDOR_PKG) $(OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

all: $(TARGET_EXEC)

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES);

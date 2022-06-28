# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

#
# Makefile for the Mesh
#

# some generic defs.
CC   = gcc
ARCH = X86_64
XGCC = gcc
XLD  = ld
XGXX = g++
HOST = $(shell gcc -dumpmachine)

TARGET_EXEC = mesh

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

#LDFLAGS+=-L${VENDOR_LIB} 
#LDFLAGS+=-L${VENDOR_LIB64}

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

$(TARGET_EXEC): $(OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

all: $(TARGET_EXEC)

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES);


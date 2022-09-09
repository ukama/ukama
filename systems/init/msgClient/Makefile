# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

#
# Makefile for the system msgClient
#
# Note: Make sure to set GRPC_INSTALL_DIR env var.
#

include ../../../nodes/ukamaOS/config.mk

# some generic defs.
CC   = gcc
ARCH = X86_64
XGCC = gcc
XLD  = ld
XGXX = g++
GRPC = $(GRPC_INSTALL_DIR)

HOST = $(shell gcc -dumpmachine)

TARGET_EXEC = msgClient
EVENTS_GRPC_SLIB = libevents_grpc.so

UKAMA_OS     = ../../../nodes/ukamaOS/
VENDOR_DIR   = $(UKAMA_OS)/distro/vendor/
VENDOR_BUILD = $(VENDOR_DIR)/build/
VENDOR_INC   = $(VENDOR_BUILD)/include/
VENDOR_LIB   = $(VENDOR_BUILD)/lib/
VENDOR_LIB64 = $(VENDOR_BUILD)/lib64/

# Packages needed to compile target
VENDOR_PKG=orcania
VENDOR_PKG+=yder
VENDOR_PKG+=libmicrohttpd
VENDOR_PKG+=ulfius
VENDOR_PKG+=tomlc
VENDOR_PKG+=jansson
VENDOR_PKG+=amqp
VENDOR_PKG+=protobuf-c

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

CPP_FLAGS=-fpic
CPP_FLAGS+=-shared
CPP_FLAGS+=-I$(GRPC)/include
CPP_FLAGS+=-I./inc
CPP_FLAGS+=-I$(GRPC)/third_party/abseil-cpp/
CPP_FLAGS+=-std=c++11

LDFLAGS+=-L$(VENDOR_LIB)
LDFLAGS+=-L$(VENDOR_LIB64)
LDFLAGS+=-L./
LDFLAGS+=-L/home/kashif/work/ukama/grpc_install/lib

CPP_LDFLAGS+=-L/usr/local/lib -pthread -Wl,--no-as-needed -ldl

# Libraries needed.
ULFIUS_LIB=-lulfius -lorcania -lyder -lmicrohttpd -lgnutls -lnettle
ULFIUS_LIB+=-lhogweed -lp11-kit -lz -lulfius
CURL_LIB=-lcurl -lssl -lcrypto
LIBS = $(ULFIUS_LIB) $(CURL_LIB)
LIBS+=-ljansson
LIBS+=-luuid
LIBS+=-lpthread
LIBS+=-ltoml
LIBS+=-lprotobuf-c
LIBS+=-lrabbitmq
LIBS+=-lgrpc++
LIBS+=$(EVENTS_GRPC_SLIB)

# Soruce files.
CFILES    = $(wildcard ./src/*.c)
OBJFILES  = $(CFILES:.c=.o)
CPP_FILES = $(wildcard ./src/cpp/*.cc)

.PHONY: $(TARGET_EXEC)

mesh: $(TARGET_EXEC)

$(TARGET_EXEC): $(VENDOR_PKG) $(EVENTS_GRPC_SLIB) $(OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@

$(EVENTS_GRPC_SLIB):
	$(XGXX) $(CPP_FLAGS) $(CPP_FILES) -o $@ $(CPP_LDFLAGS)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

all: $(TARGET_EXEC)

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES) $(EVENTS_GRPC_SLIB);

# protoc --grpc_out=. --plugin=protoc-gen-grpc=/home/kashif/work/ukama/grpc_install/bin/grpc_cpp_plugin  events.proto 

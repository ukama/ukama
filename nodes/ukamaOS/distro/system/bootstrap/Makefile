# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

#
# Makefile for the bootstrap client
#

include ../../../config.mk

TARGET_EXEC = bootstrap

UKAMAOS_ROOT = ../../../

VENDOR_DIR   = $(UKAMAOS_ROOT)/distro/vendor/
VENDOR_BUILD = $(VENDOR_DIR)/build/
VENDOR_INC   = $(VENDOR_BUILD)/include/
VENDOR_LIB   = $(VENDOR_BUILD)/lib/
VENDOR_LIB64 = $(VENDOR_BUILD)/lib64/

# Packages needed to compile target
VENDOR_PKG = tomlc curl jansson ulfius

# Setting up various compile and link flags
CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I$(VENDOR_INC)
CFLAGS+=-I./inc
CFLAGS+=-I./inc/common/
CFLAGS+=-I../common/inc
CFLAGS+=-D_REENTRANT

LDFLAGS+=-L${VENDOR_LIB} 
LDFLAGS+=-L$(VENDOR_LIB64)

# Libraries needed
CURL_LIB=-lcurl -lssl -lcrypto

LIBS = $(ULFIUS_LIB) $(CURL_LIB)
LIBS+=-ltoml
LIBS+=-ljansson

# Soruce files
CFILES   = $(wildcard ./src/*.c)
OBJFILES = $(CFILES:.c=.o)

XCC := $(XGCCPATH)$(XGCC)
XLD := $(XGCCPATH)$(XLD)

.PHONY: $(TARGET_EXEC) $(AGENT_EXEC)

bootstrap: $(TARGET_EXEC)

$(TARGET_EXEC): $(VENDOR_PKG) $(OBJFILES) $(COMM_OBJFILES)
	$(XCC) -o $(TARGET_EXEC) $(OBJFILES) $(COMM_OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@ 

%.o: %.c
	$(XCC) $(CFLAGS) -c $< -o $@

all: bootstrap

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES)

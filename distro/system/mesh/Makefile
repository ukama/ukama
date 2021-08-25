# Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.

#
# Simple Makefile for the Mesh.d
#

include ../../../config.mk

TARGET_EXEC=mesh.d
PREFIX=/usr/local
MBEDTLS_LOCATION=third_party/mbedtls/
TOML_LOCATION=third_party/tomlc99
MBEDTLS_INCLUDE=third_party/mbedtls/include

COMMON_INC=../common/inc

VENDORDIR=../../vendor
VENDOR_INCLUDE=$(VENDORDIR)/build/include
VENDOR_LIB=$(VENDORDIR)/build/lib

ULFIUS_LIB= -lulfius -lorcania -lyder -lmicrohttpd -lgnutls -lnettle -lhogweed -lp11-kit -lunistring -lz
CURL_LIB= -lcurl -lssl -lcrypto

LIBS=$(ULFIUS_LIB) $(CURL_LIB) -ljansson -luuid -lpthread -ltoml 
CFLAGS+=-c -Wall -I./inc -I$(MBEDTLS_INCLUDE) -I$(COMMON_INC)  \
	-I$(VENDOR_INCLUDE) -D_REENTRANT -g -O0 -DTEST_EMBED_CERT

LDFLAGS+=-L${VENDOR_LIB}


CFILES   = $(wildcard ./src/*.c)
OBJFILES = $(CFILES:.c=.o)

XCC = $(XGCCPATH)$(XGCC)
XLD = $(XGCCPATH)$(XLD)

.PHONY: $(TARGET_EXEC)

mesh: $(TARGET_EXEC)

$(TARGET_EXEC): $(OBJFILES)
	$(XCC) -o $(TARGET_EXEC) $(OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

%.o: %.c
	$(XCC) $(CFLAGS) -c $< -o $@

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES);
	cd ./third_party/mbedtls/ && make clean && cd -;
	cd $(TOML_LOCATION) && make clean



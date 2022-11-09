# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

#
# Makefile for the initClient
#

include ../../../nodes/ukamaOS/config.mk

# some generic defs.
CC   = gcc
ARCH = X86_64
XGCC = gcc
XLD  = ld
XGXX = g++

HOST = $(shell gcc -dumpmachine)

# Targets
TARGET_EXEC      = initClient
TARGET_EXEC_TEST = initClientTest

UKAMA_OS     = ../../../nodes/ukamaOS/
VENDOR_DIR   = $(UKAMA_OS)/distro/vendor/
VENDOR_BUILD = $(VENDOR_DIR)/build/
VENDOR_INC   = $(VENDOR_BUILD)/include/
VENDOR_LIB   = $(VENDOR_BUILD)/lib/
VENDOR_LIB64 = $(VENDOR_BUILD)/lib64/
UNITY_ROOT   = $(UKAMA_OS)/distro/tools/unity/

# Source paths
UNITY_ROOT = $(UKAMA_OS)/distro/tools/unity
SRC_DIR    = src

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
CFLAGS+=-I$(UNITY_ROOT)/src/
CFLAGS+=-I./inc
CFLAGS+=-D_REENTRANT
CFLAGS+=-DTEST_EMBED_CERT

LDFLAGS+=-L$(VENDOR_LIB)
LDFLAGS+=-L$(VENDOR_LIB64)

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

# Memory check and input flags
MEMCHECK_REPORT := memcheck.report
MEMCHECK := valgrind
MEMCHECK_FLAGS+=--log-file=$(MEMCHECK_REPORT)
MEMCHECK_FLAGS+=--track-origins=yes
MEMCHECK_FLAGS+=--leak-check=full
MEMCHECK_FLAGS+=--show-leak-kinds=all

# Source files.
CFILES   = $(wildcard ./src/*.c)
OBJFILES = $(CFILES:.c=.o)

CFILES_TEST   := $(wildcard ./src/*.c ./test/*.c $(UNITY_ROOT)/src/*.c)
CFILES_TEST   := $(filter-out ./src/main.c, $(CFILES_TEST))
OBJFILES_TEST := $(CFILES_TEST:.c=.o)

# Main target for building
all: $(TARGET_EXEC)
	@echo Done.

# CLANG FORMAT
define clang-format
	@echo "Formatting $1";
	$(shell clang-format-9 -i $1;)
	$(shell clang-tidy-10 -checks='-*,readability-identifier-naming' \
		    -config="{CheckOptions: [ \
 { key: readability-identifier-naming.NamespaceCase, value: camelBack },\
 { key: readability-identifier-naming.ClassCase, value: CamelCase  },\
 { key: readability-identifier-naming.StructCase, value: CamelCase  },\
 { key: readability-identifier-naming.FunctionCase, value: lower_case },\
 { key: readability-identifier-naming.VariableCase, value: camelBack },\
 { key: readability-identifier-naming.TypedefCase, value: CamelCase },\
 { key: readability-identifier-naming.GlobalConstantCase, value: camelBack },\
 { key: readability-braces-around-statements.ShortStatementLines, value: 0}\
	    ]}" --fix "$1" -- $(INCFLAGS) )
endef

# Format code to ukama style
formatcodestyle:
	$(foreach src, $(SOURCES), $(call clang-format,$(src)))
	@echo "Source file Done."

# Check ukama code style
checkcodestyle:
	$(eval CWD = $(shell pwd))
	echo $(CWD)
	@for src in $(SOURCES) ; do \
		echo "Checking format for $(CWD)/$$src"; \
		dif=`clang-format-9 "$(CWD)/$$src" | diff "$(CWD)/$$src" - | wc -l`; \
		if [ `echo $$dif` != `echo 0` ]; then \
			echo "clang-format: Fail";\
			echo "Err: $$src $$dif Lines to be modified."; \
			echo "Execute \" make format-style\".";\
			exit 1;\
		else\
			echo "clang-format: Pass.";\
		fi ;\
		echo "clang-tidy: ";\
		clang-tidy-10 --checks='-*,readability-identifier-naming' \
			-config="{CheckOptions: [ \
  { key: readability-identifier-naming.NamespaceCase, value: camelBack },\
  { key: readability-identifier-naming.ClassCase, value: CamelCase  },\
  { key: readability-identifier-naming.StructCase, value: CamelCase  },\
  { key: readability-identifier-naming.FunctionCase, value: lower_case },\
  { key: readability-identifier-naming.VariableCase, value: camelBack },\
  { key: readability-identifier-naming.TypedefCase, value: CamelCase },\
  { key: readability-identifier-naming.GlobalConstantCase, value: camelBack },\
  { key: readability-braces-around-statements.ShortStatementLines, value: 0}\
         ]}" --quiet "$$src" -- $(INCFLAGS); \
	done
	@echo "Ukama Coding Style check is ... PASS"

initClient: $(TARGET_EXEC)
test: $(TARGET_EXEC_TEST)

$(TARGET_EXEC): $(VENDOR_PKG) $(OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

$(TARGET_EXEC_TEST): $(VENDOR_PKG) $(OBJFILES_TEST)
	$(CC) -o $(TARGET_EXEC_TEST) $(OBJFILES_TEST) $(LDFLAGS) $(LIBS)

$(VENDOR_PKG):
#	$(MAKE) -C $(VENDOR_DIR) $@

# Build Object files
%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

# Memory check
memcheck: $(TARGET_EXEC)
	$(MEMCHECK) $(MEMCHECK_FLAGS) ./$(BUILD_DIR)/$(TARGET_EXEC)

# Clean all build files
clean:
	rm -rf $(TARGET_EXEC) $(TARGET_EXEC_TEST) $(OBJFILES) $(OBJFILES_TEST) \
		memcheck.report
	@echo Cleaned.

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2021-present, Ukama Inc.

include ../../../config.mk

# Targets
TARGETBIN:=node.d
INVENTORY:=genInventory
SCHEMA:=genSchema
UTEST:=nodeTest

VERSION_FILE = version.h
ARTIFACTS = $(TARGETBIN)

APP_NAME = $(CUR_BUILD_DIRNAME)
APP_CONFIG = 
TARGET_CONFIG_DIR = $(INSTALL_DIR)$(NODE_APP_CONFIG_DIR)/$(APP_NAME)

DISTRO_DIR:=$(UKAMAOS_ROOT)/distro

# Vendor
VENDOR_DIR:=$(DISTRO_DIR)/vendor
VENDOR_HEADERS_DIR:=$(VENDOR_DIR)/build/include
VENDOR_LIBS_DIR:=$(VENDOR_DIR)/build/lib
VENDOR_LIBS64_DIR:=$(VENDOR_DIR)/build/lib64

# Unit test
UNITY_ROOT:=$(VENDOR_DIR)/Unity

# Platform
PLATFORM_DIR:=$(DISTRO_DIR)/platform
PLATFORM_HEADERS_DIRS:=$(PLATFORM_DIR)/sys/inc
PLATFORM_HEADERS_DIRS+=$(PLATFORM_DIR)/log/inc
PLATFORM_LIB_DIR:=$(PLATFORM_DIR)/build/

INC_HEADERS+=$(VENDOR_DIR)/headers/musl/
INC_HEADERS+=$(DISTRO_DIR)/system/common/inc

# JSON lib
JSON_LIB:=jansson

# HTTP server lib
HTTP_SERVER_LIB:=ulfius

# Source paths
CORE_DIRS=core utils
SVC_DIRS=service
INVT_DIRS=mfg/inventory mfg/common
SCH_DIR=mfg/schema mfg/common

SVC_SRC_DIRS:=$(CORE_DIRS) $(SVC_DIRS)
IVNT_SRC_DIRS:=$(CORE_DIRS) $(INVT_DIRS)
SCH_SRC_DIRS:=$(CORE_DIRS) $(SCH_DIR)

# Includes
INC_DIRS=inc 
INC_DIRS+=$(UNITY_ROOT)/src
INC_DIRS+=$(PLATFORM_HEADERS_DIRS)
INC_DIRS+=$(VENDOR_HEADERS_DIR)
INC_DIRS+=$(INC_HEADERS)
TEST_DIRS:=test $(UNITY_ROOT)/src

# Build
BUILD_DIR:=build

# VENDOR_PKG
ifeq ($(ALPINE_BUILD),1)
    $(info Excluding vendor library for Alpine build)
        VENDOR_PKG=
else
	VENDOR_PKG = tomlc jansson ulfius
endif

# Platform Package
PLATFORM_PKG = usys

# Libraries
ULFIUS_LIB_DEP=-lorcania 
ifneq ($(ALPINE_BUILD),1)
ULFIUS_LIB_DEP+=-lyder
endif
ULFIUS_LIB_DEP+=-lmicrohttpd
ULFIUS_LIB_DEP+=-lgnutls
ULFIUS_LIB_DEP+=-lnettle
ULFIUS_LIB_DEP+=-lhogweed 
ULFIUS_LIB_DEP+=-lp11-kit
ULFIUS_LIB_DEP+=-lz
CURL_LIB_DEP=-lcurl 
CURL_LIB_DEP+=-lssl
CURL_LIB_DEP+=-lcrypto

LIBS+=-lpthread
LIBS+=-lrt
LIBS+=-lm
LIBS+=-lusys
LIBS+=-l$(JSON_LIB)
LIBS+=-l$(HTTP_SERVER_LIB)
LIBS+=$(ULFIUS_LIB_DEP)
LIBS+=$(CURL_LIB_DEP)


# Compiler flags
CFLAGS+=-g
CFLAGS+=-O0
CFLAGS+=-Wall
CFLAGS+=-Wno-unused-variable
CFLAGS+=-fPIC
CFLAGS+=-DHAVE_SYS_TIME_H

LDFLAGS+=$(LDPATH)
LDFLAGS+=-L$(PLATFORM_LIB_DIR) 
LDFLAGS+=-L$(VENDOR_LIBS_DIR)
LDFLAGS+=-L$(VENDOR_LIBS64_DIR)
LDFLAGS+=$(RPATH_FLAGS)

# Memory check and input flags
MEMCHECK_REPORT := $(BUILD_DIR)/memcheck.report
MEMCHECK := valgrind
MEMCHECK_FLAGS+=--log-file=$(MEMCHECK_REPORT)
MEMCHECK_FLAGS+=--track-origins=yes
MEMCHECK_FLAGS+=--leak-check=full 
MEMCHECK_FLAGS+=--show-leak-kinds=all

# Source extensions
SRC_EXTS := c

# Header extensions
HDR_EXTS := h 

# Includes
INC := $(foreach dir, $(INC_DIRS), $(foreach ext, $(HDR_EXTS), $(wildcard $(dir)/*.$(ext))))

# List of all recognized files found in the specified directories for test
#CFILES := $(foreach dir, $(SRC_DIRS), $(foreach ext, $(SRC_EXTS), $(wildcard $(dir)/*.$(ext))))

# Noded Service source files 
SVC_CFILES := $(shell find $(SVC_SRC_DIRS) -name '*.c')
SVC_OBJFILES := $(foreach ext, $(SRC_EXTS), $(patsubst %.$(ext), $(BUILD_DIR)/%.o, $(filter %.$(ext), $(SVC_CFILES))))

# Inventory source files
INVT_CFILES := $(shell find $(IVNT_SRC_DIRS) -name '*.c')
INVT_OBJFILES := $(foreach ext, $(SRC_EXTS), $(patsubst %.$(ext), $(BUILD_DIR)/%.o, $(filter %.$(ext), $(INVT_CFILES))))

# Schema Source files
SCH_CFILES := $(shell find $(SCH_SRC_DIRS) -name '*.c')
SCH_OBJFILES := $(foreach ext, $(SRC_EXTS), $(patsubst %.$(ext), $(BUILD_DIR)/%.o, $(filter %.$(ext), $(SCH_CFILES))))

# Unit Test 
TEST_CFILES := $(foreach dir, $(TEST_DIRS), $(foreach ext, $(SRC_EXTS), $(wildcard $(dir)/*.$(ext))))
TEST_OBJFILES := $(foreach ext, $(SRC_EXTS), $(patsubst %.$(ext), $(BUILD_DIR)/%.o, $(filter %.$(ext), $(TEST_CFILES))))

INC_FILES := $(INC_DIRS:%=-I%)

$(info CC:: $(CC))
$(info LDFLAGS :: $(LDFLAGS))

$(info CFILES :: $(SVC_CFILES))
$(info OBJFILES :: $(SVC_OBJFILES))
$(info INC :: $(INC_FILES))

.PHONY: $(TARGETBIN) $(UTEST) $(BUILD) formatcodestyle checkcodestyle memcheck clean version.h

# Main target for building
all: $(TARGETBIN) $(INVENTORY) $(SCHEMA) version.h 
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
	$(foreach src, $(SVC_CFILES), $(call clang-format,$(src)))	
	@echo "Source file Done."

# Check ukama code style
checkcodestyle:
	$(eval CWD = $(shell pwd))
	echo $(CWD)
	@for src in $(SVC_CFILES) ; do \
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
	@echo "Ukama Coding Style check pass..!!"

$(VENDOR_PKG):
	cd $(VENDOR_DIR) && $(MAKE) $@

$(PLATFORM_PKG):
	cd $(PLATFORM_DIR) && $(MAKE)

# Build Target
$(TARGETBIN): $(SVC_OBJFILES)
	@echo "Building $(TARGETBIN)" 
	$(CC) -o $(BUILD_DIR)/$@ $^ $(LDFLAGS) $(LIBS)

# Build Registry utility
$(INVENTORY): $(INVT_OBJFILES) 
	@echo "Building $(INVENTORY)" 
	$(CC) -o $(BUILD_DIR)/$@ $^ $(LDFLAGS) $(LIBS)

# Build Schema utility
$(SCHEMA): $(SCH_OBJFILES)
	@echo "Building $(SCHEMA)" 
	$(CC) -o $(BUILD_DIR)/$@ $^ $(LDFLAGS) $(LIBS)

# Build Unit test binary
$(UTEST):$(TEST_OBJFILES) $(TARGETBIN)
	$(CC) -o $(BUILD_DIR)/$@ $(TEST_OBJFILES) $(LDFLAGS) $(LIBS) -lusys
	@echo CC: $(BUILD_DIR)/$@ 

# Build Object files 
$(BUILD_DIR)/%.o: %.c $(VENDOR_PKG) $(PLATFORM_PKG) 
	mkdir -p $(dir $@)
	$(CC) $(CFLAGS) $(INC_FILES) -c $< -o $@

# Memory check
memcheck: $(UTEST)
	$(MEMCHECK) $(MEMCHECK_FLAGS) ./$(BUILD_DIR)/$(UTEST)

# Clean all build files
clean:
	rm -rf $(BUILD_DIR) version.h
	@echo Cleaned.

version.h:
	$(NODES_DIR)/utils/scripts/generate_version.sh

install: $(ARTIFACTS) version.h
	mkdir -p $(TARGET_CONFIG_DIR)
	cp -vrf build/$(ARTIFACTS) $(INSTALL_DIR)$(NODE_APP_DIR)
	cp -vrf $(VERSION_FILE) $(TARGET_CONFIG_DIR)
	if [ "$(APP_CONFIG)" != "" ]; then \
		cp -vrf $(APP_CONFIG) $(TARGET_CONFIG_DIR);\
	fi

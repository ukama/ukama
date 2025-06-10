# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2022-present, Ukama Inc.

include ../../config.mk

UKAMAOS_ROOT = ../../

# Targets
TAR_LIB := usys
TARGET_LIB := lib$(TAR_LIB).so
TARGET_STATIC_LIB := lib$(TAR_LIB).a 
UTEST_BIN := platform
DISTRO_DIR:=../

# Vendor
VENDOR_DIR   = $(UKAMAOS_ROOT)/distro/vendor
VENDOR_BUILD = $(VENDOR_DIR)/build
VENDOR_INC   = $(VENDOR_BUILD)/include
VENDOR_LIB   = $(VENDOR_BUILD)/lib
VENDOR_LIB64 = $(VENDOR_BUILD)/lib64

#VENDOR_INC+=$(VENDOR_DIR)/headers/musl/

# Source paths
UNITY_ROOT := $(VENDOR_DIR)/Unity
SRC_DIRS := sys/src log/src
INC_DIRS := sys/inc log/inc $(UNITY_ROOT)/src $(VENDOR_INC)
BUILD_DIR := build
TEST_DIRS := test $(UNITY_ROOT)/src

# Libraries
CURL_LIB=-lcurl -lssl -lcrypto
ULFIUS_LIB=-lulfius -lorcania -lmicrohttpd -lgnutls -lnettle
ULFIUS_LIB+=-lhogweed -lp11-kit -lz
ifeq ($(ALPINE_BUILD),1)
    $(info Excluding yder library for Alpine build)
else
    ULFIUS_LIB+=-lyder
endif

LIBS = $(ULFIUS_LIB) $(CURL_LIB)
LIBS+=-lpthread
LIBS+=-ljansson
LIBS+=-lrt
LIBS+=-lm

# Compiler flags
CFLAGS+=-g
CFLAGS+=-O0
CFLAGS+=-Wall
CFLAGS+=-Wno-unused-variable
CFLAGS+=-fPIC
CFLAGS+=-DHAVE_SYS_TIME_H
CFLAGS+=-DLOG_USE_COLOR

LDFLAGS+=$(LDPATH)
LDFLAGS+=-L$(CURDIR)/$(BUILD_DIR)
LDFLAGS+=-L${VENDOR_LIB}
LDFLAGS+=-L$(VENDOR_LIB64)

LDFLAGS+=-Wl,-rpath,$(PLATFORM_LIB)
LDFLAGS+=-Wl,-rpath,$(VENDOR_LIB)
LDFLAGS+=-Wl,-rpath,$(VENDOR_LIB64)

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

# List of all recognized files found in the specified directories for test
CFILES := $(foreach dir, $(SRC_DIRS), $(foreach ext, $(SRC_EXTS), $(wildcard $(dir)/*.$(ext))))
OBJFILES := $(foreach ext, $(SRC_EXTS), $(patsubst %.$(ext), $(BUILD_DIR)/%.o, $(filter %.$(ext), $(CFILES))))
INC := $(foreach dir, $(INC_DIRS), $(foreach ext, $(HDR_EXTS), $(wildcard $(dir)/*.$(ext))))

# Unit Test 
TEST_CFILES := $(foreach dir, $(TEST_DIRS), $(foreach ext, $(SRC_EXTS), $(wildcard $(dir)/*.$(ext))))
TEST_OBJFILES := $(foreach ext, $(SRC_EXTS), $(patsubst %.$(ext), $(BUILD_DIR)/%.o, $(filter %.$(ext), $(TEST_CFILES))))

INC_FILES := $(INC_DIRS:%=-I%)

$(info XCC:: $(XCC))
$(info LDFLAGS :: $(LDFLAGS))

$(info CFILES :: $(CFILES))
$(info OBJFILES :: $(OBJFILES))
$(info INC :: $(INC_FILES))

.PHONY: install clean memcheck
# Main target for building
all: $(TARGET_STATIC_LIB) $(TARGET_LIB)
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
	@echo "Ukama Coding Style check pass..!!"

# Build Target lib
$(TARGET_STATIC_LIB): $(OBJFILES)
	@echo "Building static lib"
	ar rcs  $(BUILD_DIR)/$@ $^

# Build Target lib 
$(TARGET_LIB): $(OBJFILES)
	@echo "Building dynamic lib" 
	$(CC) -o $(BUILD_DIR)/$@ $^ $(LDFLAGS) $(LIBS) -shared

# Build Unit test binary
$(UTEST_BIN):$(TEST_OBJFILES) $(TARGET_LIB)
	$(CC) -o $(BUILD_DIR)/$@ $(TEST_OBJFILES) $(LDFLAGS) $(LIBS) -l$(TAR_LIB)
	@echo CC: $(BUILD_DIR)/$@ 

# Build Object files 
$(BUILD_DIR)/%.o: %.c 
	mkdir -p $(dir $@)
	$(CC) $(CFLAGS) $(INC_FILES) -c $< -o $@

# Memory check
memcheck: $(UTEST_BIN)
	$(MEMCHECK) $(MEMCHECK_FLAGS) ./$(BUILD_DIR)/$(UTEST_BIN)

install: $(TARGET_STATIC_LIB) $(TARGET_LIB)
	mkdir -p $(VENDOR_BUILD)
	mkdir -p $(VENDOR_INC)
	@echo Installing platfrom lib
	cp -rvf ./sys/inc/*.h $(VENDOR_INC)
	cp -rvf ./log/inc/*.h $(VENDOR_INC)
	cp -rvf ./build/libusys.* $(VENDOR_LIB)

# Clean all build files
clean:
	rm -rf $(BUILD_DIR)
	@echo Cleaned.


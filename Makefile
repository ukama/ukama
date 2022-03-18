# Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.

#
# Makefile for the platform library
#

include ../../../config.mk

# Targets
TARGET:=noded
UTEST:=nodeTest


DISTRO_DIR:=../..

# Unit test
UNITY_ROOT:=$(DISTRO_DIR)/tools/unity

# Vendor
VENDOR_DIR:=$(DISTRO_DIR)/vendors
VENDOR_HEADERS_DIR:=$(VENDOR_DIR)/build/include
VENDOR_LIBS_DIR:=$(VENDOR_DIR)build/libs

# Platform
PLATFORM_DIR:=$(DISTRO_DIR)/platform
PLATFORM_HEADERS_DIRS:=$(PLATFORM_DIR)/sys/inc
PLATFORM_HEADERS_DIRS+=$(PLATFORM_DIR)/log/inc
PLATFORM_LIB_DIR:=$(PLATFORM_DIR)/build/
PLATFORM_LIB:=usys 

# JSON lib
JSON_LIB:=jansson

# Source paths
SRC_DIRS=src 
SRC_DIRS+=utils/src
SRC_DIRS+=src/ledger
SRC_DIRS+=src/ledger
INC_DIRS=inc 
INC_DIRS+=utils
INC_DIRS+=$(UNITY_ROOT)/src
INC_DIRS+=$(PLATFORM_HEADERS_DIRS)
INC_DIRS+=$(VENDOR_HEADERS_DIR)

TEST_DIRS:=test $(UNITY_ROOT)/src

# Build
BUILD_DIR:=build


# Compilers and flags
ifdef XGCCPATH
XCC=$(XGCCPATH)$(XGCC)
XLD=$(XGCCPATH)$(XLD)
else
XCC=gcc
endif

# Libraries
LIBS+=-lpthread
LIBS+=-lrt
LIBS+=-lm
LIBS+=-l$(PLATFORM_LIB)
LIBS+=-l$(JSON_LIB)


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
#CFILES := $(foreach dir, $(SRC_DIRS), $(foreach ext, $(SRC_EXTS), $(wildcard $(dir)/*.$(ext))))
CFILES := $(shell find $(SOURCEDIR) -name '*.c')
OBJFILES := $(foreach ext, $(SRC_EXTS), $(patsubst %.$(ext), $(BUILD_DIR)/%.o, $(filter %.$(ext), $(CFILES))))
INC := $(foreach dir, $(INC_DIRS), $(foreach ext, $(HDR_EXTS), $(wildcard $(dir)/*.$(ext))))

# Unit Test 
TEST_CFILES := $(foreach dir, $(TEST_DIRS), $(foreach ext, $(SRC_EXTS), $(wildcard $(dir)/*.$(ext))))
TEST_OBJFILES := $(foreach ext, $(SRC_EXTS), $(patsubst %.$(ext), $(BUILD_DIR)/%.o, $(filter %.$(ext), $(TEST_CFILES))))

INC_FILES := $(INC_DIRS:%=-I%)

$(info CC:: $(CC))
$(info LDFLAGS :: $(LDFLAGS))

$(info CFILES :: $(CFILES))
$(info OBJFILES :: $(OBJFILES))
$(info INC :: $(INC_FILES))

.PHONY: $(TARGET) $(UTEST) $(BUILD) formatcodestyle checkcodestyle memcheck clean

# Main target for building
all: $(TARGET)
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

# Build Target
$(TARGET): $(OBJFILES)
	@echo "Building $(TARGET)" 
	$(XCC) -o $(BUILD_DIR)/$@ $^ $(LDFLAGS) $(LIBS)

# Build Unit test binary
$(UTEST):$(TEST_OBJFILES) $(TARGET)
	$(XCC) -o $(BUILD_DIR)/$@ $(TEST_OBJFILES) $(LDFLAGS) $(LIBS) -l$(PLATFORM_LIB)
	@echo CC: $(BUILD_DIR)/$@ 

# Build Object files 
$(BUILD_DIR)/%.o: %.c 
	mkdir -p $(dir $@)
	$(XCC) $(CFLAGS) $(INC_FILES) -c $< -o $@

# Memory check
memcheck: $(UTEST)
	$(MEMCHECK) $(MEMCHECK_FLAGS) ./$(BUILD_DIR)/$(UTEST)
		
# Clean all build files
clean:
	rm -rf $(BUILD_DIR)
	@echo Cleaned.





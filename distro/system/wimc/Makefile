#
# Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.
#

#
# Makefile for WIMC.d
#
include ../../../config.mk

XCC=$(XGCCPATH)$(XGCC)
XLD=$(XGCCPATH)$(XLD)

COMMON_DIR=$(CURDIR)/../common
COMMON_INC=$(COMMON_DIR)/inc
BUILD_DIR:=build

SRC_EXTS:=c
HDR_EXTS:=h

# source directories
SRC_DIR=common src
COMMON_SRC_DIR=common
PROVIDER_SRC_DIR=provider common
AGENTS_SRC_DIR=agents common

VENDORDIR=../../vendor
VENDOR_INCLUDE=$(VENDORDIR)/build/include
VENDOR_LIB=$(VENDORDIR)/build/lib

# inlcude directories
INC_DIRS:= $(ULFIUS_INC) $(COMMON_INC) $(CURDIR)/inc $(CURDIR)/inc/common $(VENDOR_INCLUDE)

SRCS:=$(foreach dir, $(SRC_DIR), $(foreach ext, $(SRC_EXTS), $(wildcard $(dir)/*.$(ext))))
OBJS:=$(foreach ext, $(SRC_EXTS), $(patsubst %.$(ext), $(BUILD_DIR)/%.o, $(filter %.$(ext), $(SRCS))))


P_SRCS:=$(foreach dir, $(PROVIDER_SRC_DIR), $(foreach ext, $(SRC_EXTS), $(wildcard $(dir)/*.$(ext))))
P_OBJS:=$(foreach ext, $(SRC_EXTS), $(patsubst %.$(ext), $(BUILD_DIR)/%.o, $(filter %.$(ext), $(P_SRCS))))

A_SRCS:=$(foreach dir, $(AGENTS_SRC_DIR), $(foreach ext, $(SRC_EXTS), $(wildcard $(dir)/*.$(ext))))
A_OBJS:=$(foreach ext, $(SRC_EXTS), $(patsubst %.$(ext), $(BUILD_DIR)/%.o, $(filter %.$(ext), $(A_SRCS))))


#INCLUDES := $(foreach dir, $(INC_DIRS), $(foreach ext, $(HDR_EXTS), $(wildcard $(dir)/*.$(ext))))

# Build  targets
EXE=wimc.d
PROVIDER=mock_provider
AGENT=agent

CFLAGS+=-g -Wall -D_REENTRANT

ULFIUS_LIB= -lulfius -lorcania -lyder -lmicrohttpd -lgnutls -lnettle -lhogweed -lp11-kit -lunistring -lz
CURL_LIB= -lcurl -lssl -lcrypto
LDFLAG=-L$(VENDOR_LIB)
LIBS=-lc $(ULFIUS_LIB) $(CURL_LIB) -ljansson -lsqlite3 -luuid
PROVIDER_LIBS=$(ULFIUS_LIB) -lsqlite3 -ljansson -luuid
AGENTS_LIBS=$(ULFIUS_LIB) $(CURL_LIB) -ljansson -luuid -lpthread -lrt

INCFLAGS:=$(INC_DIRS:%=-I%)
DEPFLAGS := -MMD -MP

# Make sure Make do not delete included dependencies files
.PRECIOUS: $(BUILDDIR)/%.d

# Function to compile using $(CC) : (files: .c)
define compilecc
	@mkdir -p $(dir $1)
	@$(XCC) -c $2 -o $1 $(CFLAGS) $(INCFLAGS) -MT $1 -MF $(BUILD_DIR)/$3.Td $(DEPFLAGS)
	@mv -f $(BUILD_DIR)/$3.Td $(BUILD_DIR)/$3.d && touch $1
	@echo CC: $1
endef


$(BUILD_DIR)/%.o:%.c $(BUILD_DIR)/%.d 
	$(call compilecc,$@,$<,$*)

$(BUILD_DIR)/%.d: ;

.PHONY: all run clean


all: $(EXE) 


$(PROVIDER): $(P_OBJS)
	$(XCC) -I$(COMMON_INC) -g -o $(BUILD_DIR)/$(PROVIDER) $(PROVIDER_SRCS) $(LDFLAG) $(PROVIDER_LIBS)
	echo "Binary available at: $(BUILD_DIR)/$(PROVIDER)"

$(AGENT): $(A_OBJS)
	$(XCC) -g -o $(BUILD_DIR)/$(AGENT) $(A_OBJS) $(LDFLAG) $(AGENTS_LIBS)
	echo "Binary available at: $(BUILD_DIR)/$(AGENT)"

$(EXE): $(OBJS)  
	$(XCC) -o $@ $(OBJS) $(LINKFLAG) $(CFLAGS) $(LDFLAG) $(LIBS)
	echo CC: $@

run: $(EXE)
	$<

clean:
	rm -rf $(BUILD_DIR);
	cd $(ULFIUS) && make clean





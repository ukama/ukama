#
# Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.
#
# This source code is licensed under the XXX-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.
#

#
# Makefile for WIMC.d
#
COMMON_DIR=$(CURDIR)/../common
PREFIX=/usr/local
ULFIUS=$(COMMON_DIR)/third_party/ulfius/src
ULFIUS_INC=$(COMMON_DIR)/third_party/ulfius/include
COMMON_INC=$(COMMON_DIR)/inc
BUILD_DIR:=build

SRC_EXTS:=c
HDR_EXTS:=h

# source directories
SRC_DIR=common src
COMMON_SRC_DIR=common
PROVIDER_SRC_DIR=provider common
AGENTS_SRC_DIR=agents common

# inlcude directories
INC_DIRS:= $(ULFIUS_INC) $(COMMON_INC) $(CURDIR)/inc $(CURDIR)/inc/common

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

CC=gcc
CFLAGS+=-g -Wall -D_REENTRANT
#LDFLAGS=-L$(ULFIUS)
LIBS=-lc -lulfius -lorcania -lsqlite3 -ljansson -lcurl -luuid
PROVIDER_LIBS=-lulfius -lsqlite3 -ljansson -luuid
AGENTS_LIBS=-lulfius -lorcania -ljansson -lcurl -luuid -lpthread -lrt

INCFLAGS:=$(INC_DIRS:%=-I%)
DEPFLAGS := -MMD -MP

# Make sure Make do not delete included dependencies files
.PRECIOUS: $(BUILDDIR)/%.d

# Function to compile using $(CC) : (files: .c)
define compilecc
	@mkdir -p $(dir $1)
	@$(CC) -c $2 -o $1 $(CFLAGS) $(INCFLAGS) -MT $1 -MF $(BUILD_DIR)/$3.Td $(DEPFLAGS)
	@mv -f $(BUILD_DIR)/$3.Td $(BUILD_DIR)/$3.d && touch $1
	@echo CC: $1
endef


$(BUILD_DIR)/%.o:%.c $(BUILD_DIR)/%.d 
	$(call compilecc,$@,$<,$*)

$(BUILD_DIR)/%.d: ;

.PHONY: all run clean


all: $(EXE) 

#libulfius:
#	cd $(ULFIUS) && $(MAKE) debug GNUTLSFLAG=1 WEBSOCKETFLAG=1 YDERFLAG=1 UWSCFLAG=1

$(PROVIDER): $(P_OBJS)
	@$(CC) -I$(ULFIUS_INC) -I$(COMMON_INC) -g -o $(BUILD_DIR)/$(PROVIDER) $(PROVIDER_SRCS) -L$(ULFIUS) $(PROVIDER_LIBS)
	echo "Binary available at: $(BUILD_DIR)/$(PROVIDER)"

$(AGENT): $(A_OBJS)
	@$(CC) -I$(ULFIUS_INC) -I$(COMMON_INC) -g -o $(BUILD_DIR)/$(AGENT) $(A_OBJS) -L$(ULFIUS) $(AGENTS_LIBS)
	echo "Binary available at: $(BUILD_DIR)/$(AGENT)"

$(EXE): $(OBJS)  
	@$(CC) -o $@ $(OBJS) $(LINKFLAG) $(CFLAGS) $(LDFLAGS) $(LIBS)
	echo CC: $@

run: $(EXE)
	$<

clean:
	rm -rf $(BUILD_DIR);
	cd $(ULFIUS) && make clean





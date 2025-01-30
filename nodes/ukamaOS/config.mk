# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2021-present, Ukama Inc.

#
# Global configuration file for all Makefiles.
#

CURMAKE = $(abspath $(firstword $(MAKEFILE_LIST)))
CURPATH = $(dir $(CURMAKE))

CUR_BUILD_DIRNAME = $(notdir $(patsubst %/,%,$(CURPATH)))

NODES_DIR = $(shell echo $(CURPATH) | sed 's|\(.*nodes\)/.*|\1|')
UKAMAOS_ROOT = $(NODES_DIR)/ukamaOS

OS = $(shell uname -s)
NPROCS = 1
ifeq ($(OS), Linux)
        NPROCS = $(shell grep -c ^processor /proc/cpuinfo)
endif
export NPROCS

# Build system
BUILD = x86_64-unknown-linux-gnu

# HOST system
HOST = x86_64-unknown-linux-gnu

override CC 	= gcc
ARCH   	    	= $(x86_64)
OPENSSLTARGET   = linux-generic32
GCCPATH 	= /usr/bin

#Supported architectures
ARCH_ARM    = arm
ARCH_X86    = x86
ARCH_X86_64 = x86_64
ARCH_ARM64  = aarch64

override AMPLIFIER_NODE = amplifier
override TOWER_NODE     = tower
override ACCESS_NODE    = access
override LOCAL          = linux

ifndef TARGET
	override TARGET_BOARD = LOCAL
	export TARGET=$(LOCAL)
else
	override TARGET_BOARD = $(TARGET)
	export TARGET
endif


# Setup paths for configs
APP_CONFIG_DIR = $(NODES_DIR)/configs/capps
NODE_APP_CONFIG_DIR = /conf

# Setup paths for apps
NODE_APP_DIR = /sbin/

ifeq ($(AMPLIFIER_NODE), $(TARGET_BOARD))
	override ARCH   = $(ARCH_ARM)
	HOST            = armv6-alpine-linux-musleabihf
endif

ifeq ($(TOWER_NODE), $(TARGET_BOARD))
	override ARCH   = $(ARCH_X86_64)
	HOST            = x86_64-linux-musl
	OPENSSLTARGET   = linux-generic64
endif

ifeq ($(ACCESS_NODE), $(TARGET_BOARD))
	override ARCH   = $(ARCH_ARM64)
	OPENSSLTARGET   = linux-aarch64
endif

export

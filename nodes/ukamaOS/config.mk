# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2021-present, Ukama Inc.

#
# Global configuration file for all Makefiles.
#

# Get current file and path
CUR_MAKE := $(abspath $(firstword $(MAKEFILE_LIST)))
CUR_PATH := $(dir $(CUR_MAKE))

# Extract directory name for build (useful for folder-based builds)
CUR_BUILD_DIRNAME := $(notdir $(patsubst %/,%,$(CUR_PATH)))

# Define root directory for nodes and UKAMA OS
NODES_DIR := $(shell echo $(CUR_PATH) | sed 's|\(.*nodes\)/.*|\1|')
UKAMAOS_ROOT := $(NODES_DIR)/ukamaOS

VENDOR_DIR   := $(UKAMAOS_ROOT)/distro/vendor
VENDOR_BUILD := $(VENDOR_DIR)/build
VENDOR_INC   := $(VENDOR_BUILD)/include
VENDOR_LIB   := $(VENDOR_BUILD)/lib
VENDOR_LIB64 := $(VENDOR_BUILD)/lib64

PLATFORM_DIR     := $(UKAMAOS_ROOT)/distro/platform/
PLATFORM_BUILD   := $(PLATFORM_DIR)/build/
PLATFORM_INC_SYS := $(PLATFORM_DIR)/sys/inc
PLATFORM_INC_LOG := $(PLATFORM_DIR)/log/inc
PLATFORM_LIB     := $(PLATFORM_DIR)/build/

# used by various Makefile
ARCH_ARM := arm
ARCH_X86 := x86
ARCH_X86_64 := x86_64
ARCH_ARM64 := aarch64

# OS and Processor configuration
OS := $(shell uname -s)
NPROCS := 1
ifeq ($(OS), Linux)
    NPROCS := $(shell grep -c ^processor /proc/cpuinfo)
endif
export NPROCS

# Build and host systems (assumed x86_64 for now)
BUILD := x86_64-unknown-linux-gnu
HOST := x86_64-unknown-linux-gnu

# Set default compiler and paths
override CC := gcc
OPENSSLTARGET := linux-generic32
GCCPATH := /usr/bin
TARGET ?= local
ARCH := $(ARCH_X86_64)
export TARGET

# Paths for application configurations
APP_CONFIG_DIR := $(NODES_DIR)/configs/capps
NODE_APP_CONFIG_DIR := /conf

# Paths for application binaries
NODE_APP_DIR := /sbin/

# Conditional assignments based on TARGET_BOARD
ifeq ($(TARGET),amplifier)
    ARCH := $(ARCH_ARM)
    HOST := armv6-alpine-linux-musleabihf
endif

ifeq ($(TARGET),tower)
    ARCH := $(ARCH_X86_64)
    HOST := x86_64-linux-musl
    OPENSSLTARGET := linux-generic64
endif

ifeq ($(TARGET),access)
    ARCH := $(ARCH_ARM64)
    OPENSSLTARGET := linux-aarch64
endif

ifeq ($(TARGET),local)
	ARCH := $(ARCH_X86_64)
	OPENSSLTARGET := linux-generic32
endif

# BUILD_MODE (debug|release) -> set RPATH_FLAGS for Makefile

# default to debug:
BUILD_MODE ?= debug
export BUILD_MODE

# pick your rpath directories
ifeq ($(BUILD_MODE),release)
  # single hard‐coded path on “make BUILD_MODE=release”
  RPATH_PATHS := /ukama/apps/lib
else
  # debug mode: use the three per‐app vars (must be defined *before* include)
  RPATH_PATHS := $(PLATFORM_LIB) $(VENDOR_LIB) $(VENDOR_LIB64)
endif

# turn them into -Wl,-rpath,<dir> flags
RPATH_FLAGS := $(foreach D,$(RPATH_PATHS),-Wl,-rpath,$(D))
export RPATH_FLAGS

# Export updated variables
export

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
ARCH := $(x86_64)
OPENSSLTARGET := linux-generic32
GCCPATH := /usr/bin

# Supported architectures
ARCH_ARM := arm
ARCH_X86 := x86
ARCH_X86_64 := x86_64
ARCH_ARM64 := aarch64

# Define nodes and local board
override AMPLIFIER_NODE := amplifier
override TOWER_NODE := tower
override ACCESS_NODE := access
override LOCAL := linux

# Set TARGET_BOARD and TARGET if not already set
ifndef TARGET
	override TARGET := $(shell echo $(LOCAL) | tr '[:upper:]' '[:lower:]')
else
	override TARGET := $(shell echo $(TARGET) | tr '[:upper:]' '[:lower:]')
endif
override TARGET_BOARD := $(TARGET)
export TARGET

# Paths for application configurations
APP_CONFIG_DIR := $(NODES_DIR)/configs/capps
NODE_APP_CONFIG_DIR := /conf

# Paths for application binaries
NODE_APP_DIR := /sbin/

# Conditional assignments based on TARGET_BOARD
ifeq ($(TARGET_BOARD), $(AMPLIFIER_NODE))
    override ARCH := $(ARCH_ARM)
    HOST := armv6-alpine-linux-musleabihf
endif

ifeq ($(TARGET_BOARD), $(TOWER_NODE))
    override ARCH := $(ARCH_X86_64)
    HOST := x86_64-linux-musl
    OPENSSLTARGET := linux-generic64
endif

ifeq ($(TARGET_BOARD), $(ACCESS_NODE))
    override ARCH := $(ARCH_ARM64)
    OPENSSLTARGET := linux-aarch64
endif

# Export updated variables
export


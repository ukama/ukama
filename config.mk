# Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.

#
# Global configuration file for all Makefiles.
#

CURMAKE = $(abspath $(firstword $(MAKEFILE_LIST)))
CURPATH = $(dir $(CURMAKE))

#OS
OS = $(shell uname -s)
NPROCS = 1
ifeq ($(OS), Linux)
        NPROCS = $(shell grep -c ^processor /proc/cpuinfo)
endif

# Build system
BUILD = x86_64-unknown-linux-gnu

#Supported architectures
ARCH_ARM    = arm
ARCH_X86    = x86
ARCH_X86_64 = x86_64

#if variables are not defined
#amplifier Node
override A_NODE = anode
override C_NODE = cnode
override LOCAL  = linux

#TARGET
ifndef TARGET
	override TARGET = $(LOCAL)
endif

#Kernelheaders
KERNEL_HEADERS = $(CURPATH)/distro/helpers/kernelheaders/usr/include

# Setup various compilier and linker options for various targets.

ifeq ($(A_NODE), $(TARGET))
	override CC     = arm-linux-gnueabihf-gcc
	override ARCH   = $(ARCH_ARM)
	XCROSS_COMPILER = arm-linux-musleabihf-
	XGCC            = $(XCROSS_COMPILER)gcc
	XLD             = $(XCROSS_COMPILER)ld
	XGXX            = $(XCROSS_COMPILER)g++
	HOST            = arm-linux-musleabihf
	OPENSSLTARGET   = linux-generic32
endif

ifeq ($(C_NODE), $(TARGET))
	override CC     =
	override ARCH   = $(ARCH_X86_64)
	XCROSS_COMPILER = x86_64-linux-musl-
	XGCC            = $(XCROSS_COMPILER)gcc
	XLD             = $(XCROSS_COMPILER)ld
	XGXX            = $(XCROSS_COMPILER)g++
	HOST            = x86_64-linux-musl
	OPENSSLTARGET   = linux-generic64
endif

ifeq ($(LOCAL), $(TARGET))
	override CC     = gcc
	override ARCH   = $(ARCH_X86_64)
	XGCC            = gcc
	XLD             = ld
	XGXX            = g++
	HOST            =
	OPENSSLTARGET   = linux-generic64
endif

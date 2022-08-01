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
override ANODEBOARD = anode
override CNODEBOARD = cnode
override HNODEBOARD = hnode
override LOCAL  = linux

#TARGET
ifndef TARGET
	override TARGETBOARD = $(LOCAL)
else
	override TARGETBOARD = $(TARGET)
endif

#Kernelheaders
KERNEL_HEADERS = $(CURPATH)/distro/helpers/kernelheaders/usr/include

# Setup various compilier and linker options for various targets.

ifeq ($(ANODEBOARD), $(TARGETBOARD))
	override CC     = arm-linux-gnueabihf-gcc
	override ARCH   = $(ARCH_ARM)
	XCROSS_COMPILER = arm-linux-musleabihf-
	XGCC            = $(XCROSS_COMPILER)gcc
	XLD             = $(XCROSS_COMPILER)ld
	XGXX            = $(XCROSS_COMPILER)g++
	HOST            = arm-linux-musleabihf
	OPENSSLTARGET   = linux-generic32
endif

ifeq ($(CNODEBOARD), $(TARGETBOARD))
	override CC     =
	override ARCH   = $(ARCH_X86_64)
	XCROSS_COMPILER = x86_64-linux-musl-
	XGCC            = $(XCROSS_COMPILER)gcc
	XLD             = $(XCROSS_COMPILER)ld
	XGXX            = $(XCROSS_COMPILER)g++
	HOST            = x86_64-linux-musl
	OPENSSLTARGET   = linux-generic64
endif

ifeq ($(LOCAL), $(TARGETBOARD))
	override CC     = gcc
	override ARCH   = $(ARCH_X86_64)
	XGCC            = gcc
	XLD             = ld
	XGXX            = g++
	HOST            = $(shell gcc -dumpmachine)
	OPENSSLTARGET   = linux-generic64
	GCCPATH        = $(shell which gcc)
	XGCCPATH       = $(dir $(GCCPATH))
endif

export

# Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.


#
# Global configuration file for all makefiles.
#

CURMAKE := $(abspath $(firstword $(MAKEFILE_LIST)))
CURPATH := $(dir $(CURMAKE))

#OS
OS:=$(shell uname -s)

#CPUS
NPROCS:=1
ifeq ($(OS), Linux)
        NPROCS:=$(shell grep -c ^processor /proc/cpuinfo)
endif

# Build system
BUILD := x86_64-unknown-linux-gnu

#Supported architectures
ARCHARM := arm
ARCHX86	:= x86
ARCHX86_64 := x86_64

#if variables are not defined
#aNode
ifndef ANODEBOARD
override ANODEBOARD = anode
endif

#cNode
ifndef CNODEBOARD
override CNODEBOARD = cnode
endif

#TARGETNODE
ifndef TARGETBOARD
override TARGETBOARD = $(ANODEBOARD)
endif

#Kernelheaders
KERNELHEADERS = $(CURPATH)/distro/helpers/kernelheaders/usr/include

#XGCCPATH
#defined in dostro makefile

#GCC Compiler for Firmware aNode
ifeq ($(ANODEBOARD), $(TARGETBOARD))
override CC = arm-linux-gnueabihf-gcc
override ARCH = $(ARCHARM)
XCROSS_COMPILER := arm-linux-musleabihf-
XGCC := $(XCROSS_COMPILER)gcc
XLD := $(XCROSS_COMPILER)ld
XGXX := $(XCROSS_COMPILER)g++
HOST := arm-linux-musleabihf
OPENSSLTARGET := linux-generic32
endif

#GCC Compiler for Firmware cNode
ifeq ($(CNODEBOARD), $(TARGETBOARD))
override CC =
override ARCH = $(ARCHX86)
XCROSS_COMPILER := x86_64-linux-musl-
XGCC := $(XCROSS_COMPILER)gcc
XLD := $(XCROSS_COMPILER)ld
XGXX := $(XCROSS_COMPILER)g++
HOST := x86_64-linux-musl
OPENSSLTARGET := linux-generic64
endif



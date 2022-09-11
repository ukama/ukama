# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

#
# Makefile for the system msgClient
#
# Note: Make sure to set GRPC_INSTALL_DIR env var.
#

include ../../../nodes/ukamaOS/config.mk

# some generic defs.
CC   = gcc
ARCH = X86_64
XGCC = gcc
XLD  = ld
XGXX = g++
GRPC = /home/vishal/.local/grpc

HOST = $(shell gcc -dumpmachine)

TARGET_EXEC = msgClient
EVENTS_GRPC_SLIB = libevents_grpc.so

UKAMA_OS     = ../../../nodes/ukamaOS/
VENDOR_DIR   = $(UKAMA_OS)/distro/vendor/
VENDOR_BUILD = $(VENDOR_DIR)/build/
VENDOR_INC   = $(VENDOR_BUILD)/include/
VENDOR_LIB   = $(VENDOR_BUILD)/lib/
VENDOR_LIB64 = $(VENDOR_BUILD)/lib64/

# Packages needed to compile target
VENDOR_PKG=orcania
VENDOR_PKG+=yder
VENDOR_PKG+=libmicrohttpd
VENDOR_PKG+=ulfius
VENDOR_PKG+=tomlc
VENDOR_PKG+=jansson
VENDOR_PKG+=amqp
VENDOR_PKG+=protobuf

# Setting up various compile and link flags.
CFLAGS=-c
CFLAGS+=-g
CFLAGS+=-Wall
CFLAGS+=-O0
CFLAGS+=-I$(VENDOR_INC)
CFLAGS+=-I$(VENDOR_INC)/rabbitmq-c/
CFLAGS+=-I./inc
CFLAGS+=-D_REENTRANT
CFLAGS+=-DTEST_EMBED_CERT

CPP_FLAGS=-fpic
CPP_FLAGS+=-shared
CPP_FLAGS+=-I$(GRPC)/include
CPP_FLAGS+=-I./inc
CPP_FLAGS+=-I$(GRPC)/third_party/abseil-cpp/
CPP_FLAGS+=-std=c++11

LDFLAGS+=-L$(VENDOR_LIB)
LDFLAGS+=-L$(VENDOR_LIB64)
LDFLAGS+=-L./
LDFLAGS+=-L$(GRPC)/lib

CPP_LDFLAGS+=-L/usr/local/lib -L$(GRPC)/lib -pthread -Wl,--no-as-needed -ldl -lprotobuf -lgrpc++ -lgrpc -lssl -lgpr 

# Libraries needed.
ULFIUS_LIB=-lulfius -lorcania -lyder -lmicrohttpd -lgnutls -lnettle
ULFIUS_LIB+=-lhogweed -lp11-kit -lz -lulfius
CURL_LIB=-lcurl -lcrypto
LIBS = $(ULFIUS_LIB) $(CURL_LIB)
LIBS+=-ljansson
LIBS+=-luuid
LIBS+=-lpthread
LIBS+=-ltoml
LIBS+=-lrabbitmq
LIBS+=-lgrpc++
LIBS+=-levents_grpc -lgrpc++ -lgrpc -lgpr
LIBS+=-labsl_bad_any_cast_impl
LIBS+=-labsl_bad_optional_access
LIBS+=-labsl_bad_variant_access
LIBS+=-labsl_base
LIBS+=-labsl_city
LIBS+=-labsl_civil_time
LIBS+=-labsl_cord_internal
LIBS+=-labsl_cord
LIBS+=-labsl_cordz_functions
LIBS+=-labsl_cordz_handle
LIBS+=-labsl_cordz_info
LIBS+=-labsl_cordz_sample_token
LIBS+=-labsl_debugging_internal
LIBS+=-labsl_demangle_internal
LIBS+=-labsl_examine_stack
LIBS+=-labsl_exponential_biased
LIBS+=-labsl_failure_signal_handler
LIBS+=-labsl_flags_commandlineflag_internal
LIBS+=-labsl_flags_commandlineflag
LIBS+=-labsl_flags_config
LIBS+=-labsl_flags_internal
LIBS+=-labsl_flags_marshalling
LIBS+=-labsl_flags_parse
LIBS+=-labsl_flags_private_handle_accessor
LIBS+=-labsl_flags_program_name
LIBS+=-labsl_flags_reflection
LIBS+=-labsl_flags
LIBS+=-labsl_flags_usage_internal
LIBS+=-labsl_flags_usage
LIBS+=-labsl_graphcycles_internal
LIBS+=-labsl_hash
LIBS+=-labsl_hashtablez_sampler
LIBS+=-labsl_int128
LIBS+=-labsl_leak_check_disable
LIBS+=-labsl_leak_check
LIBS+=-labsl_log_severity
LIBS+=-labsl_low_level_hash
LIBS+=-labsl_malloc_internal
LIBS+=-labsl_periodic_sampler
LIBS+=-labsl_random_distributions
LIBS+=-labsl_random_internal_distribution_test_util
LIBS+=-labsl_random_internal_platform
LIBS+=-labsl_random_internal_pool_urbg
LIBS+=-labsl_random_internal_randen_hwaes_impl
LIBS+=-labsl_random_internal_randen_hwaes
LIBS+=-labsl_random_internal_randen_slow
LIBS+=-labsl_random_internal_randen
LIBS+=-labsl_random_internal_seed_material
LIBS+=-labsl_random_seed_gen_exception
LIBS+=-labsl_random_seed_sequences
LIBS+=-labsl_raw_hash_set
LIBS+=-labsl_raw_logging_internal
LIBS+=-labsl_scoped_set_env
LIBS+=-labsl_spinlock_wait
LIBS+=-labsl_stacktrace
LIBS+=-labsl_statusor
LIBS+=-labsl_status
LIBS+=-labsl_strerror
LIBS+=-labsl_str_format_internal
LIBS+=-labsl_strings_internal
LIBS+=-labsl_strings
LIBS+=-labsl_symbolize
LIBS+=-labsl_synchronization
LIBS+=-labsl_throw_delegate
LIBS+=-labsl_time
LIBS+=-labsl_time_zone
LIBS+=-laddress_sorting
LIBS+=-lcrypto
LIBS+=-lgpr
LIBS+=-lgrpc++_alts
LIBS+=-lgrpc++_error_details
LIBS+=-lgrpc_plugin_support
LIBS+=-lgrpcpp_channelz
LIBS+=-lgrpc++_reflection
LIBS+=-lgrpc++
LIBS+=-lgrpc
LIBS+=-lgrpc++_unsecure
LIBS+=-lgrpc_unsecure
LIBS+=-lre2
LIBS+=-lssl
LIBS+=-lupb

# Soruce files.
CFILES    = $(wildcard ./src/*.c)
OBJFILES  = $(CFILES:.c=.o)
CPP_FILES = $(wildcard ./src/cpp/*.cc)

.PHONY: $(TARGET_EXEC)

mesh: $(TARGET_EXEC)

$(TARGET_EXEC): $(VENDOR_PKG) $(EVENTS_GRPC_SLIB) $(OBJFILES)
	$(CC) -o $(TARGET_EXEC) $(OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

$(VENDOR_PKG):
	#$(MAKE) -C $(VENDOR_DIR) $@

$(EVENTS_GRPC_SLIB):
	$(XGXX) $(CPP_FLAGS) $(CPP_FILES) -o $@ $(CPP_LDFLAGS)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

all: $(TARGET_EXEC)

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES) $(EVENTS_GRPC_SLIB);

# protoc --grpc_out=. --plugin=protoc-gen-grpc=/home/kashif/work/ukama/grpc_install/bin/grpc_cpp_plugin  events.proto 

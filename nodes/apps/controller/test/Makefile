# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# For pre-set cities:
# make sim-realtime SIM_CITY="Goma, DRC"
# make sim-realtime SIM_CITY="Nairobi, Kenya"

# Random city:
# make sim-realtime-api SIM_CITY="Skardu, Pakistan" SIM_TZ_OFFSET=5

# Exact coordinates:
# make sim-realtime-coords \
#    SIM_LAT=-1.6792 \
#    SIM_LON=29.2228 \
#    SIM_TZ_OFFSET=2

# for quick test:
# make sim-sunny
# make sim-night
# make sim-all

TEST_DIR  := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
ROOT_DIR  := $(abspath $(TEST_DIR)/..)
UKAMA_DIR := $(abspath $(ROOT_DIR)/../..)

include $(UKAMA_DIR)/ukamaOS/config.mk

CC     ?= cc
PYTHON ?= python3
TARGET := controllerd

CONTROLLER_HOST ?= 127.0.0.1
CONTROLLER_PORT ?= 18021

SCENARIO         ?= sunny
SIM_CITY         ?= Goma, DRC
SIM_LAT          ?=
SIM_LON          ?=
SIM_TZ_OFFSET    ?=
SIM_PV_COUNT     ?= 4
SIM_PV_WATTS     ?= 450
SIM_PV_DERATE    ?= 0.82
SIM_BATTERY_V    ?= 48
SIM_BATTERY_AH   ?= 200
SIM_SOC_INIT     ?= 0.72
SIM_LOAD_DAY_W   ?= 180
SIM_LOAD_NIGHT_W ?= 120
SIM_CLOUD_FACTOR ?= 0.95
SIM_AMBIENT_DAY_C   ?= 31
SIM_AMBIENT_NIGHT_C ?= 22

VENDOR_DIR   := $(UKAMAOS_ROOT)/distro/vendor
VENDOR_BUILD := $(VENDOR_DIR)/build
VENDOR_INC   := $(VENDOR_BUILD)/include
VENDOR_LIB   := $(VENDOR_BUILD)/lib
VENDOR_LIB64 := $(VENDOR_BUILD)/lib64

PLATFORM_DIR     := $(UKAMAOS_ROOT)/distro/platform
PLATFORM_BUILD   := $(PLATFORM_DIR)/build
PLATFORM_INC_SYS := $(PLATFORM_DIR)/sys/inc
PLATFORM_INC_LOG := $(PLATFORM_DIR)/log/inc
PLATFORM_LIB     := $(PLATFORM_DIR)/build

ifeq ($(ALPINE_BUILD),1)
    $(info Excluding vendor library for Alpine build)
    VENDOR_PKG :=
else
    VENDOR_PKG := tomlc jansson ulfius libuuid
endif

PLATFORM := platform

CFLAGS  := -c
CFLAGS  += -g
CFLAGS  += -Wall
CFLAGS  += -Wextra
CFLAGS  += -O0
CFLAGS  += -D_REENTRANT
CFLAGS  += -I$(ROOT_DIR)
CFLAGS  += -I$(ROOT_DIR)/inc
CFLAGS  += -I$(ROOT_DIR)/../common/inc
CFLAGS  += -I$(VENDOR_INC)
CFLAGS  += -I$(PLATFORM_INC_SYS)
CFLAGS  += -I$(PLATFORM_INC_LOG)

LDFLAGS += $(RPATH_FLAGS)
LDFLAGS += -L$(VENDOR_LIB)
LDFLAGS += -L$(VENDOR_LIB64)
LDFLAGS += -L$(PLATFORM_LIB)

CURL_LIB   := -lcurl -lssl -lcrypto
ULFIUS_LIB := -lulfius -lorcania -lmicrohttpd -lgnutls -lnettle
ULFIUS_LIB += -lhogweed -lp11-kit -lz

ifneq ($(ALPINE_BUILD),1)
ULFIUS_LIB += -lyder
endif

LIBS  := $(ULFIUS_LIB) $(CURL_LIB)
LIBS  += -ljansson
LIBS  += -lpthread
LIBS  += -lm
LIBS  += -lusys

SRCS := $(wildcard $(ROOT_DIR)/src/*.c)
OBJS := $(SRCS:.c=.o)

.PHONY: all clean check api-test
.PHONY: sim-scenario sim-sunny sim-night sim-all
.PHONY: sim-realtime sim-realtime-api sim-realtime-coords
.PHONY: $(PLATFORM) $(VENDOR_PKG)

all: $(TARGET)

$(TARGET): $(VENDOR_PKG) $(PLATFORM) $(OBJS) $(COMM_OBJFILES)
	$(CC) -o $@ $(OBJS) $(COMM_OBJFILES) $(LDFLAGS) $(LIBS)
	@echo ""
	@echo "Build complete: ./$(TARGET)"
	@echo ""
	@echo "Run daemon:"
	@echo "  CONTROLLER_DRIVER=victron \\"
	@echo "  CONTROLLER_SERIAL_PORT=/dev/ttyXXX \\"
	@echo "  ./$(TARGET)"
	@echo ""
	@echo "Simulation:"
	@echo "  make sim-realtime"
	@echo "  make sim-realtime SIM_CITY=\"Nairobi, Kenya\""
	@echo "  make sim-realtime-api SIM_CITY=\"Skardu, Pakistan\" \\"
	@echo "       SIM_TZ_OFFSET=5"
	@echo "  make sim-realtime-coords SIM_LAT=-1.6792 \\"
	@echo "       SIM_LON=29.2228 SIM_TZ_OFFSET=2"
	@echo "  make sim-sunny"
	@echo "  make sim-all"
	@echo ""
	@echo "API test:"
	@echo "  make api-test"
	@echo ""

$(PLATFORM):
	$(MAKE) -C $(PLATFORM_DIR)

$(VENDOR_PKG):
	$(MAKE) -C $(VENDOR_DIR) $@

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

check:
	$(PYTHON) -m py_compile $(TEST_DIR)vedirect_sim.py
	bash -n $(TEST_DIR)test_api.sh

sim-scenario:
	cd $(TEST_DIR) && \
	$(PYTHON) vedirect_sim.py --mode scenario --scenario $(SCENARIO)

sim-sunny:
	cd $(TEST_DIR) && \
	$(PYTHON) vedirect_sim.py --mode scenario --scenario sunny

sim-night:
	cd $(TEST_DIR) && \
	$(PYTHON) vedirect_sim.py --mode scenario --scenario night

sim-all:
	cd $(TEST_DIR) && \
	$(PYTHON) vedirect_sim.py --mode scenario --scenario all

sim-realtime:
	cd $(TEST_DIR) && \
	$(PYTHON) vedirect_sim.py --mode realtime \
		--city "$(SIM_CITY)" \
		--pv-count $(SIM_PV_COUNT) \
		--pv-watts $(SIM_PV_WATTS) \
		--pv-derate $(SIM_PV_DERATE) \
		--battery-v $(SIM_BATTERY_V) \
		--battery-ah $(SIM_BATTERY_AH) \
		--soc-init $(SIM_SOC_INIT) \
		--load-day-w $(SIM_LOAD_DAY_W) \
		--load-night-w $(SIM_LOAD_NIGHT_W) \
		--cloud-factor $(SIM_CLOUD_FACTOR) \
		--ambient-day-c $(SIM_AMBIENT_DAY_C) \
		--ambient-night-c $(SIM_AMBIENT_NIGHT_C)

sim-realtime-api:
	cd $(TEST_DIR) && \
	$(PYTHON) vedirect_sim.py --mode realtime \
		--city "$(SIM_CITY)" \
		--use-api \
		$(if $(SIM_TZ_OFFSET),--tz-offset $(SIM_TZ_OFFSET)) \
		--pv-count $(SIM_PV_COUNT) \
		--pv-watts $(SIM_PV_WATTS) \
		--pv-derate $(SIM_PV_DERATE) \
		--battery-v $(SIM_BATTERY_V) \
		--battery-ah $(SIM_BATTERY_AH) \
		--soc-init $(SIM_SOC_INIT) \
		--load-day-w $(SIM_LOAD_DAY_W) \
		--load-night-w $(SIM_LOAD_NIGHT_W) \
		--cloud-factor $(SIM_CLOUD_FACTOR) \
		--ambient-day-c $(SIM_AMBIENT_DAY_C) \
		--ambient-night-c $(SIM_AMBIENT_NIGHT_C)

sim-realtime-coords:
	@if [ -z "$(SIM_LAT)" ] || [ -z "$(SIM_LON)" ] || \
	     [ -z "$(SIM_TZ_OFFSET)" ]; then \
		echo "SIM_LAT, SIM_LON, and SIM_TZ_OFFSET are required."; \
		echo "Example:"; \
		echo "  make sim-realtime-coords SIM_LAT=-1.6792 \\"; \
		echo "       SIM_LON=29.2228 SIM_TZ_OFFSET=2"; \
		exit 1; \
	fi
	cd $(TEST_DIR) && \
	$(PYTHON) vedirect_sim.py --mode realtime \
		--lat $(SIM_LAT) \
		--lon $(SIM_LON) \
		--tz-offset $(SIM_TZ_OFFSET) \
		--pv-count $(SIM_PV_COUNT) \
		--pv-watts $(SIM_PV_WATTS) \
		--pv-derate $(SIM_PV_DERATE) \
		--battery-v $(SIM_BATTERY_V) \
		--battery-ah $(SIM_BATTERY_AH) \
		--soc-init $(SIM_SOC_INIT) \
		--load-day-w $(SIM_LOAD_DAY_W) \
		--load-night-w $(SIM_LOAD_NIGHT_W) \
		--cloud-factor $(SIM_CLOUD_FACTOR) \
		--ambient-day-c $(SIM_AMBIENT_DAY_C) \
		--ambient-night-c $(SIM_AMBIENT_NIGHT_C)

api-test:
	cd $(TEST_DIR) && \
	CONTROLLER_HOST=$(CONTROLLER_HOST) \
	CONTROLLER_PORT=$(CONTROLLER_PORT) \
	./test_api.sh

clean:
	rm -f $(TARGET) $(OBJS)
	rm -f $(TEST_DIR)*.pyc
	rm -rf $(TEST_DIR)__pycache__

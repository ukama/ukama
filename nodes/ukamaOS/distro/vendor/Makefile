# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2021-present, Ukama Inc.

include ../../config.mk

CUR_MAKE = $(abspath $(firstword $(MAKEFILE_LIST)))
CUR_PATH = $(dir $(CUR_MAKE))

export DEP_DIR := $(CUR_PATH)build

PKG_CFG_PATH=${DEP_DIR}/lib/pkgconfig:${DEP_DIR}/lib64/pkgconfig

ifndef BUILD_CLEANUP
	BUILD_CLEANUP=0
endif

define post_build_cleanup
	$(if $(filter ${BUILD_CLEANUP},1), @echo Cleanup disabled, echo Cleaning $(1) && rm -rf $(1))
endef

define print_msg
	@echo -n "Building $(1) ........... "
endef

#BUILD FLAGS
BUILDFLAGS =

BUILD_DIR = tmpbuild
LIGHT_GREEN='\033[1;32m'
NO_COLOR='\033[0m'


INC_FLAGS  = -I.
INC_FLAGS += -I$(DEP_DIR)/include

LIBS =

LDFLAGS = -L$(DEP_DIR)/lib -L$(DEP_DIR)/lib64

ifeq ($(ARCH), $(ARCH_X86))
	export CFLAGS = $(INC_FLAGS) -fPIC
else
	export CFLAGS = $(INC_FLAGS)
endif

$(shell mkdir -p $(DEP_DIR)/lib)

# List of libs required to be build
LIST = orcania
LIST += yder
LIST += jansson
LIST += ulfius
LIST += prometheus-client
LIST += tomlc
LIST += amqp
LIST += libcap
LIST += libuuid

.PHONY: $(LIST)

build: $(LIST)

yder:
	$(call print_msg,"yder")
ifeq (,$(wildcard ${DEP_DIR}/lib/libyder.*))
	$(MAKE) -C $@/src CC=$(CC) clean 
	$(MAKE) -C $@/src static-install CC=$(CC) Y_DISABLE_JOURNALD=1 \
	       	DESTDIR=$(DEP_DIR)
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

orcania:
	$(call print_msg,"orcania")
ifeq (,$(wildcard ${DEP_DIR}/lib/liborcania.*))
	$(MAKE) -C $@/src CC=$(CC) clean
	$(MAKE) -C $@/src static-install CC=$(CC) DESTDIR=$(DEP_DIR)
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t'  ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

ulfius: orcania yder
	$(call print_msg,"ulfius")
ifeq (,$(wildcard ${DEP_DIR}/lib/libulfius.*))
	$(MAKE) -C $@/src CC=$(CC) clean
	$(MAKE) -C $@/src static-install CC=$(CC) DESTDIR=$(DEP_DIR)
	$(MAKE) -C $@ install CC=$(CC) DESTDIR=$(DEP_DIR)
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t'  ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

jansson:
	$(call print_msg,"jansson")
ifeq (,$(wildcard ${DEP_DIR}/lib/libjansson.*))
	(mkdir -p $@/$(BUILD_DIR) && cd $@/$(BUILD_DIR) && \
        cmake ../ -DCMAKE_C_COMPILER=${CC} -DJANSSON_BUILD_DOCS=OFF -DCMAKE_INSTALL_PREFIX=$(DEP_DIR) && \
	$(MAKE) install )
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

prometheus-client:
	$(call print_msg,"prometheus-client")
ifeq (,$(wildcard ${DEP_DIR}/lib/libprom*.*))
	@echo "Building prom"
	mkdir -p $@/prom/build
	cd $@/prom/build && \
		cmake ../ \
			-DCMAKE_C_COMPILER=${CC} \
			-DCMAKE_C_FLAGS="${CFLAGS}" \
			-DCMAKE_EXE_LINKER_FLAGS="${LD_FLAGS}" \
			-DCMAKE_TOOLCHAIN_FILE=$(CUR_PATH)cmake_tool_chain_file.cmake \
			-DCMAKE_INSTALL_PREFIX=/usr && \
		make && make install DESTDIR=$(DEP_DIR)

	@echo "Building promhttp"
	mkdir -p $@/promhttp/build
	cd $@/promhttp/build && \
		cmake ../ \
			-DCMAKE_C_COMPILER=${CC} \
			-DCMAKE_C_FLAGS="${CFLAGS}" \
			-DCMAKE_EXE_LINKER_FLAGS="${LD_FLAGS}" \
			-DCMAKE_TOOLCHAIN_FILE=$(CUR_PATH)cmake_tool_chain_file.cmake \
			-DCMAKE_INSTALL_PREFIX=/usr && \
		make && make install DESTDIR=$(DEP_DIR)

	$(call post_build_cleanup,$@/prom/build)
	$(call post_build_cleanup,$@/promhttp/build)
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

tomlc:
	$(call print_msg,"tomlc")
ifeq (,$(wildcard ${DEP_DIR}/lib/libtoml.*))
	$(MAKE) -C $@ CC=$(CC) clean
	$(MAKE) -C $@ CC=$(CC) 
	$(MAKE) -C $@ install prefix=$(DEP_DIR)
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

libuuid:
	$(call print_msg,"libuuid")
ifeq (,$(wildcard ${DEP_DIR}/lib/libuuid.*))
	cd $@ && \
	find . -name 'libtool' -exec rm -f {} + && \
	find . -name 'ltmain.sh' -exec rm -f {} + && \
	find . -name 'aclocal.m4' -exec rm -f {} + && \
	find . -name 'autom4te.cache' -exec rm -rf {} + && \
	find . -name 'configure' -exec rm -f {} + && \
	find . -name 'config.log' -exec rm -f {} + && \
	find . -name 'config.status' -exec rm -f {} + && \
	find . -name 'Makefile.in' -exec rm -f {} + && \
	autoreconf --force --install --verbose && cd -
	mkdir -p $@/${BUILD_DIR} && cd $@ && \
	aclocal && autoconf && cd ${BUILD_DIR} && \
	CC=${CC} CFLAGS="-I${DEP_DIR}/include" LDFLAGS="-L${DEP_DIR}/lib -static" \
		../configure --prefix=${DEP_DIR} --host=${HOST} \
		--build=${BUILD} --enable-static && \
	make && make install
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

amqp:
	$(call print_msg,"amqp")
ifeq (,$(wildcard ${DEP_DIR}/lib/librabbitmq.*))
	(mkdir -p $@/${BUILD_DIR} && cd $@/${BUILD_DIR} && \
	cmake ../ -DCMAKE_C_COMPILER=${CC} \
		-DCMAKE_TOOLCHAIN_FILE=$(CUR_PATH)cmake_tool_chain_file.cmake \
		-DCMAKE_INSTALL_PREFIX=$(DEP_DIR) \
		-DBUILD_STATIC_LIBS=ON \
		-DENABLE_SSL_SUPPORT=OFF && \
        cmake --build . --target install)
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

libcap:
	$(call print_msg,"libcap")
ifeq (,$(wildcard ${DEP_DIR}/libcap.*))
	$(MAKE) -C $@ CC=$(CC) clean
	$(MAKE) -C $@ CC=$(CC) BUILD_CC=gcc GOLANG=no install DESTDIR=$(DEP_DIR)
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

clean:
	echo "Cleaning Dependencies."
	rm -rf $(DEP_DIR)

info :
	@echo "Build info for vendor:"
	@cat package.list

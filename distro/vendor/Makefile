# Copyright (c) 2021-present, Ukama Inc.
#  All rights reserved.

# Makefile for vendor code based and libarires

include ../../config.mk

CURMAKE = $(abspath $(firstword $(MAKEFILE_LIST)))
CURPATH = $(dir $(CURMAKE))

export DEPDIR = $(CURPATH)/build

ifndef XGCCPATH
	$(error XGCCPATH is missing)
endif

ifeq ($(TARGETBOARD), $(LOCAL))
PKG_CFG_PATH=${DEPDIR}/lib64/pkgconfig
else
PKG_CFG_PATH=${DEPDIR}/lib/pkgconfig
endif

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

export XCXX = $(XGCCPATH)$(XGXX)
export XCC  = $(XGCCPATH)$(XGCC)
export LD   = $(XGCCPATH)$(XLD)

INC_FLAGS  = -I.
INC_FLAGS += -I$(DEPDIR)/include

ifeq ($(TARGET), "anode")
	INC_FLAGS+= -I$(XGCCPATH)../arm-linux-musleabihf/include
endif

LIBS =

LDFLAGS = -L$(DEPDIR)/lib -L$(DEPDIR)/lib64
ifeq ($(TARGET), "anode")
	LDFLAGS += -L$(XGCCPATH)/../arm-linux-musleabihf/lib
endif
export LDFLAGS

ifeq ($(ARCH), $(ARCHX86))
	export CFLAGS = $(INC_FLAGS) -fPIC
else
	export CFLAGS = $(INC_FLAGS)
endif

$(shell mkdir -p $(DEPDIR)/lib)

# List of libs required to be build
LIST = zlib
LIST += orcania
LIST += yder
LIST += libmicrohttpd
LIST += jansson
LIST += ulfius
LIST += prometheus-client
LIST += tomlc
LIST += openssl
LIST += curl
LIST += libuuid
LIST += nettle
LIST += gmp
LIST += p11-kit
LIST += gnutls
LIST += sqlite
LIST += amqp
LIST += protobuf
LIST += protobuf-c
LIST += libunistring
LIST += libcap

.PHONY: $(LIST) iptables iproute2 

build: $(LIST)

zlib:
	$(call print_msg,"zlib")
ifeq (,$(wildcard ${DEPDIR}/lib/libz.a))
	(mkdir -p $@/$(BUILD_DIR) && cd $@/$(BUILD_DIR) && sudo ../configure --static --prefix=$(DEPDIR) && \
		$(MAKE) -j${NPROCS} CC=$(XCC) CFLAGS="${CFLAGS} -fPIC" && \
		$(MAKE) -j${NPROCS} install)
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

yder:
	$(call print_msg,"yder")
ifeq (,$(wildcard ${DEPDIR}/lib/libyder.a))
	$(MAKE) -C $@/src CC=$(XCC) clean 
	$(MAKE) -C $@/src static-install CC=$(XCC) Y_DISABLE_JOURNALD=1 \
	       	DESTDIR=$(DEPDIR)
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

orcania:
	$(call print_msg,"orcania")
ifeq (,$(wildcard ${DEPDIR}/lib/liborcania.a))
	$(MAKE) -C $@/src CC=$(XCC) clean
	$(MAKE) -C $@/src static-install CC=$(XCC) DESTDIR=$(DEPDIR)
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t'  ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

libmicrohttpd:
	$(call print_msg,"libmicrohttpd")
ifeq (,$(wildcard ${DEPDIR}/lib/libmicrohttpd.*))
	(mkdir -p $@/$(BUILD_DIR) && cd $@ && ./bootstrap && cd $(BUILD_DIR) && \
		../configure CC=$(XCC) --host=$(HOST) \
		--disable-curl --disable-https --prefix=$(DEPDIR) && \
		$(MAKE) -C ./src install)
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e  ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

ulfius: orcania yder jansson gnutls curl libmicrohttpd libunistring
	$(call print_msg,"ulfius")
ifeq (,$(wildcard ${DEPDIR}/lib/libulfius.a))
	$(MAKE) -C $@/src CC=$(XCC) clean
	$(MAKE) -C $@ install CC=$(XCC) DESTDIR=$(DEPDIR)
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t'  ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

jansson:
	$(call print_msg,"jansson")
ifeq (,$(wildcard ${DEPDIR}/lib/libjansson.*))
	(mkdir -p $@/$(BUILD_DIR) && cd $@ && autoreconf -i && cd $(BUILD_DIR) && \
		../configure CC=$(XCC) --host=$(HOST) \
		--enable-static --prefix=$(DEPDIR) && \
	$(MAKE) && $(MAKE) install)
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t'  ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

prometheus-client:
	$(call print_msg,"prometheus-client")
ifeq (,$(wildcard ${DEPDIR}/lib/libprom.*))
	mkdir -p $@/prom/build && \
	cd $@/prom/build && \
	cmake ../ -DCMAKE_C_COMPILER=${XCC} -DCMAKE_C_FLAGS="${CFLAGS}" \
		-DCMAKE_LD_FLAGS="${LD_FLAGS}" \
		-DCMAKE_TOOLCHAIN_FILE=$(CURPATH)cmake_tool_chain_file.cmake && \
	make install DESTDIR=$(DEPDIR)
	@echo "Building promhttp"
	mkdir -p $@/promhttp/build && \
	cd $@/promhttp/build && \
	cmake ../ -DCMAKE_C_COMPILER=${XCC} -DCMAKE_C_FLAGS="${CFLAGS}" \
		-DCMAKE_LD_FLAGS="${LD_FLAGS}" \
		-DCMAKE_TOOLCHAIN_FILE=$(CURPATH)cmake_tool_chain_file.cmake && \
	make install DESTDIR=$(DEPDIR)
	$(call post_build_cleanup,$@/prom/build)
	$(call post_build_cleanup,$@/promhttp/build)
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

tomlc:
	$(call print_msg,"tomlc")
ifeq (,$(wildcard ${DEPDIR}/lib/libtoml.a))
	$(MAKE) -C $@ CC=$(XCC) clean
	$(MAKE) -C $@ CC=$(XCC) 
	$(MAKE) -C $@ install prefix=$(DEPDIR)
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

libuuid:
	$(call print_msg,"libuuid")
ifeq (,$(wildcard ${DEPDIR}/lib/libuuid.*))
	mkdir -p $@/${BUILD_DIR} && cd $@ && \
	aclocal && autoconf && cd ${BUILD_DIR} && \
	CC=${XCC} CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib" \
		../configure --prefix=${DEPDIR} --host=${HOST} \
		--build=${BUILD} && \
	make && make install
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

openssl: zlib
	$(call print_msg,"openssl")
ifeq (,$(wildcard ${DEPDIR}/lib/libssl.*))
	mkdir -p $@/${BUILD_DIR} && cd $@/${BUILD_DIR} && \
	CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib \
	       -latomic" ../Configure ${OPENSSLTARGET} --release \
		--prefix=${DEPDIR} --openssldir=${DEPDIR} \
		--with-zlib-lib=${DEPDIR}/lib \
		--with-zlib-include=${DEPDIR}/include && \
	make -j${NPROCS} && make -j${NPROCS} install
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

curl: openssl
	$(call print_msg,"curl")
ifeq (,$(wildcard ${DEPDIR}/lib/libcurl.*))
	mkdir -p $@/${BUILD_DIR} && cd $@ && autoreconf -fi && \
	cd ./${BUILD_DIR} && \
	CC="${XCC}" LDFLAGS="-L${DEPDIR}/lib -L${DEPDIR}/lib64" \
		../configure --prefix=${DEPDIR} --includedir="${DEPDIR}/include" \
		LIBS="-latomic -lhogweed -lnettle" \
		--disable-libcurl-option --with-openssl=${DEPDIR} \
		--with-zlib=${DEPDIR} --host=${HOST} && \
	make -j${NPROCS} && make -j${NPROCS} install
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

nettle:
	$(call print_msg,"nettle")
ifeq (,$(wildcard ${DEPDIR}/lib/libnettle.*))
	(mkdir -p $@/${BUILD_DIR} && cd $@/${BUILD_DIR} && \
	CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib \
		-latomic" ../configure --prefix=${DEPDIR} --host=${HOST} \
		--enable-mini-gmp && \
	make -j${NPROCS} && make -j${NPROCS} install)
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

gmp:
	$(call print_msg,"gmp")
ifeq (,$(wildcard ${DEPDIR}/lib/libgmp.so))
	mkdir -p $@/${BUILD_DIR} && cd $@/${BUILD_DIR} && \
	CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib \
		-latomic" ../configure --prefix=${DEPDIR} --host=${HOST} && \
        make -j${NPROCS} && make -j${NPROCS} install
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

libunistring:
	$(call print_msg,"libunistring")
ifeq (,$(wildcard ${DEPDIR}/lib/libunistring.*))
	mkdir -p $@/${BUILD_DIR} && cd $@ && \
	cd ${BUILD_DIR} && CC="${XCC}" CFLAGS="-I${DEPDIR}/include" \
		LDFLAGS="-L${DEPDIR}/lib" ../configure --prefix=${DEPDIR} \
		--host=${HOST} && \
	make -j${NPROCS} && make -j${NPROCS} install
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

#libidn2:
#	$(call print_msg,"libidn2")
#	mkdir -p $@/${BUILD_DIR} && cd $@ && autoreconf -fi && \
#	CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib" \
#		./configure --prefix=${DEPDIR} --disable-doc --host=${HOST} && \
#	make -j${NPROCS} && make -j${NPROCS} install
#	$(call post_build_cleanup,$@/${BUILD_DIR})

#p11-kit:
#	$(call print_msg,"p11-kit")
#	mkdir -p $@/${BUILD_DIR} && cd $@/${BUILD_DIR} && \
#	CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib" \
#		../configure --prefix=${DEPDIR} --host=${HOST} --without-libffi \
#		--without-libtasn1 --without-systemd && \
#	make -j${NPROCS} && make -j${NPROCS} install
#	$(call post_build_cleanup,$@/${BUILD_DIR})

p11-kit:
	$(call print_msg,"p11-kit")
ifeq (,$(wildcard ${DEPDIR}/lib/libcrypto.*))
	mkdir -p $@/${BUILD_DIR} && cd $@/${BUILD_DIR} && \
	CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib" \
		../configure --prefix=${DEPDIR} --host=${HOST} --without-libffi \
		--without-libtasn1 --without-systemd && \
	make -j${NPROCS} && make -j${NPROCS} install
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

gnutls: p11-kit openssl nettle gmp
	$(call print_msg,"gnutls")
ifeq (,$(wildcard ${DEPDIR}/lib/libgnutls.*))
	mkdir -p $@/${BUILD_DIR} && cd $@/${BUILD_DIR} && \
	CC="${XCC}" CXX="${XCXX}" CFLAGS="-I${DEPDIR}/include" \
		LDFLAGS="-L${DEPDIR}/lib" ../configure --prefix=${DEPDIR} \
		--exec-prefix=${DEPDIR} GMP_LIBS=${DEPDIR}/lib \
		--host=${HOST} --with-included-libtasn1 \
		--disable-openssl-compatibility \
		--with-included-unistring \
		--disable-valgrind-tests --disable-guile --enable-local-libopts \
		--without-idn --disable-doc --without-p11-kit \
		PKG_CONFIG_PATH=${PKG_CFG_PATH} && \
	make -j${NPROCS} && make -j${NPROCS} install
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

sqlite:
	$(call print_msg,"sqlite")
ifeq (,$(wildcard ${DEPDIR}/lib/libsqlite3.*))
	mkdir -p $@/${BUILD_DIR} && cd $@/${BUILD_DIR} && \
	CC="${XCC}" CFLAGS="-I${DEPDIR}/include" \
		LDFLAGS="-L${DEPDIR}/lib" ../configure --prefix=${DEPDIR} \
		--host=${HOST} --disable-readline --enable-static && \
	make -j${NPROCS} && make -j${NPROCS} install
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

protobuf:
	$(call print_msg,"protobuf")
ifeq (,$(wildcard ${DEPDIR}/lib/libprotobuf.*))
	(mkdir -p $@/${BUILD_DIR} && cd $@ && ./autogen.sh && cd ${BUILD_DIR} && \
	CC="${XCC}" CXX="${XCXX}" CFLAGS="-I${DEPDIR}/include" \
		LDFLAGS="-L${DEPDIR}/lib" ../configure --prefix=${DEPDIR} \
		--enable-cross-compile --with-protoc=protoc --host=${HOST} && \
	$(MAKE) -j${NPROCS} && \
	$(MAKE) -j${NPROCS} install)
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

protobuf-c: protobuf
	$(call print_msg,"protobuf-c")
ifeq (,$(wildcard ${DEPDIR}/lib/libprotobuf-c.*))
	(mkdir -p $@/${BUILD_DIR} && cd $@ && ./autogen.sh && cd ${BUILD_DIR} && \
	CC="${XCC}" CXX="${XCXX}" CFLAGS="-I${DEPDIR}/include" \
		CXXFLAGS="-I${DEPDIR}/include" \
		LDFLAGS="-L${DEPDIR}/lib" ../configure --prefix=${DEPDIR} \
		--host=${HOST} --disable-protoc && \
	$(MAKE) -j${NPROCS} && \
        $(MAKE) -j${NPROCS} install)
	$(call post_build_cleanup,$@/${BUILD_DIR})
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

amqp:
	$(call print_msg,"amqp")
ifeq (,$(wildcard ${DEPDIR}/lib/librabbitmq.*))
	(mkdir -p $@/${BUILD_DIR} && cd $@/${BUILD_DIR} && \
	cmake ../ -DCMAKE_C_COMPILER=${XCC} \
		-DCMAKE_TOOLCHAIN_FILE=$(CURPATH)cmake_tool_chain_file.cmake \
		-DCMAKE_INSTALL_PREFIX=$(DEPDIR) \
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
ifeq (,$(wildcard ${DEPDIR}/lib64/libcap.*))
	$(MAKE) -C $@ CC=$(XCC) clean
	$(MAKE) -C $@ CC=$(XCC) BUILD_CC=gcc GOLANG=no install DESTDIR=$(DEPDIR)
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}
else
	@echo -e '\t' ${LIGHT_GREEN}Skipping${NO_COLOR}
endif

dhcpcd:
	$(call print_msg,"dhcpcd")
	(cd $@ && \
	CC="${XCC}" ./configure --small --disable-inet6 --disable-embedded \
		--libexecdir=/lib/dhcpcd --dbdir=/var/lib/dhcpcd && \
	$(MAKE) CC="${XCC}"  && \
	$(MAKE) CC="${XCC}" DESTDIR=${DEPDIR} install && \
	cd ./src && $(MAKE) clean && cd ..)
	$(call post_build_cleanup,$@/src/build)
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}

sysctl:
	$(call print_msg,"sysctl")
	(cd $@ && \
	$(MAKE) clean && $(MAKE) CC="${XCC}" && \
	cp ./sysctl ${DEPDIR}/sbin && \
	cp 00-linuxkit.conf ${DEPDIR}/etc && \
	$(MAKE) clean)
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}

iptables:
	$(call print_msg,"iptables")
	(mkdir -p $@/${BUILD_DIR} && cd $@ && autoreconf -fi && cd ${BUILD_DIR} && \
        CC=${XCC} ../configure --prefix=/ --host=${HOST} && \
        $(MAKE) && $(MAKE) install DESTDIR=${DEPDIR} && \
        $(MAKE) clean)
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}

iproute2:
	$(call print_msg,"iproute2")
	@echo "DESTDIR:"${DESTDIR}
	(cd $@  && \
        CC=${XCC} ./configure && \
        $(MAKE) && $(MAKE) install DESTDIR=${DEPDIR} && \
        $(MAKE) clean)
	@echo -e '\t' ${LIGHT_GREEN}Done${NO_COLOR}

clean:
	echo "Cleaning Dependencies."
	rm -rf $(DEPDIR)

info :
	@echo "Build info for vendor:"
	@cat package.list

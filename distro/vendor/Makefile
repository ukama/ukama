# Copyright (c) 2021-present, Ukama Inc.
#  All rights reserved.

# Makefile for vendor code based and libarires

include ../../config.mk

CURMAKE = $(abspath $(firstword $(MAKEFILE_LIST)))
CURPATH = $(dir $(CURMAKE))

export DEPDIR = $(CURPATH)/build

#ifneq ($(TARGET), "linux")
#	ifndef XGCC_PATH
#		$(error XGCC_PATH is missing)
#	endif
#endif

ifndef BUILD_CLEANUP
	BUILD_CLEANUP=0
endif

define post_build_cleanup
	ifeq ($(BUILD_CLEANUP), 0)
		@echo Cleanup disabled
	else
		@echo Cleaning $(1) && rm -rf $(1))
	endif
endef

define print_msg
	@echo Building $(1) ...........
endef

#BUILD FLAGS
BUILDFLAGS =

BUILD_DIR = build

export XCXX = $(XGCC_PATH)$(XGXX)
export XCC  = $(XGCC_PATH)$(XGCC)
export LD   = $(XGCC_PATH)$(XLD)

INC_FLAGS = -I.
INC_FLAGS+= -I$(DEPDIR)/include

ifeq ($(TARGET), "anode")
	INC_FLAGS+= -I$(XGCC_PATH)../arm-linux-musleabihf/include
endif

LIBS =

LD_FLAGS = -L$(DEPDIR)/lib
ifeq ($(TARGET), "anode")
	LD_FLAGS += -L$(XGCC_PATH)../arm-linux-musleabihf/lib
endif
export LD_FLAGS

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
LIST += libunistring
LIST += p11-kit
LIST += gnutls
LIST += sqlite
LIST += amqp
LIST += protobuf-c

.PHONY: $(LIST)

build: $(LIST)

zlib:
	$(call print_msg,"zlib")
	(cd $@ && ./configure --static --prefix=$(DEPDIR) && \
	$(MAKE) -j${NPROCS} CC=$(XCC) CFLAGS="${CFLAGS}" && \
	$(MAKE) -j${NPROCS} install)

yder:
	$(call print_msg,"yder")
	$(MAKE) -C $@/src static-install CC=$(XCC) Y_DISABLE_JOURNALD=1 \
	       	DESTDIR=$(DEPDIR)

orcania:
	$(call print_msg,"orcania")
	$(MAKE) -C $@/src static-install CC=$(XCC) DESTDIR=$(DEPDIR)

libmicrohttpd:
	$(call print_msg,"libmicrohttpd")
	(cd $@ && ./bootstrap && ./configure CC=$(XCC) --host=$(HOST) \
		--disable-curl --disable-https --prefix=$(DEPDIR) && \
	$(MAKE) -C $@/src install)

ulfius: orcania yder jansson gnutls curl libmicrohttpd
	$(call print_msg,"ulfius")
	$(MAKE) -C $@ install CC=$(XCC) DESTDIR=$(DEPDIR)

jansson:
	$(call print_msg,"jansson")
	(cd $@ && autoreconf -i && ./configure CC=$(XCC) --host=$(HOST) \
		--enable-static --prefix=$(DEPDIR) && \
	$(MAKE) -C $@ install)

prometheus-client:
	$(call print_msg,"prometheus-client")
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
	$(call postbuild_cleanup,$@/prom/build)
	$(call postbuild_cleanup,$@/promhttp/build)

tomlc:
	$(call print_msg,"tomlc")
	($(MAKE) -C $@ CC=$(XCC) && \
	$(MAKE) -C $@ install prefix=$(DEPDIR))

libuuid:
	$(call print_msg,"libuuid")
	mkdir -p $@/${BUILD_DIR} && cd $@ && \
	autoreconf -fi && cd ${BUILD_DIR} && \
	CC=${XCC} CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib" \
		../configure --prefix=${DEPDIR} --host=${HOST} \
		--build=${BUILD} && \
	make && make install
	$(call postbuild_cleanup,$@/${BUILD_DIR})

openssl: zlib
	$(call print_msg,"openssl")
	mkdir -p $@/${BUILD_DIR} && cd $@/${BUILD_DIR} && \
	CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib \
	       -latomic" ../Configure ${OPENSSLTARGET} --release \
		--prefix=${DEPDIR} --openssldir=${DEPDIR} \
		--with-zlib-lib=${DEPDIR}/lib \
		--with-zlib-include=${DEPDIR}/include && \
	make -j${NPROCS} && make -j${NPROCS} install
	$(call postbuild_cleanup,$@/${BUILD_DIR})

curl: openssl
	$(call print_msg,"curl")
	mkdir -p $@/${BUILD_DIR} && cd $@ && autoreconf -fi && \
	cd ${BUILD_DIR} && \
	CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib \
		-latomic" ../configure --prefix=${DEPDIR} \
		--disable-libcurl-option --with-openssl=${DEPDIR} \
		--with-zlib=${DEPDIR} --host=${HOST} && \
	make -j${NPROCS} && make -j${NPROCS} install
	$(call postbuild_cleanup,$@/${BUILD_DIR})

nettle:
	$(call print_msg,"nettle")
	mkdir -p $@/${BUILD_DIR} && cd $@/${BUILD_DIR} && \
	CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib \
		-latomic" ../configure --prefix=${DEPDIR} --host=${HOST} \
		--enable-mini-gmp &&
	make -j${NPROCS} && make -j${NPROCS} install
	$(call postbuild_cleanup,$@/${BUILD_DIR})

gmp:
	$(call print_msg,"gmp")
	mkdir -p $@/${BUILD_DIR} && cd $@/${BUILD_DIR} && \
	CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib \
		-latomic" ../configure --prefix=${DEPDIR} --host=${HOST} && \
        make -j${NPROCS} && make -j${NPROCS} install
	$(call postbuild_cleanup,$@/${BUILD_DIR})

libunistring:
	$(call print_msg,"libunistring")
	mkdir -p $@/${BUILD_DIR} && cd $@ && autoreconf -fi && \
	cd ${BUILD_DIR} && CC="${XCC}" CFLAGS="-I${DEPDIR}/include" \
		LDFLAGS="-L${DEPDIR}/lib" ../configure --prefix=${DEPDIR} \
		--host=${HOST} &&
	make -j${NPROCS} && make -j${NPROCS} install
	$(call postbuild_cleanup,$@/${BUILD_DIR})

libidn2: libunistring
	$(call print_msg,"libidn2")
	mkdir -p $@/${BUILD_DIR} && cd $@ && autoreconf -fi && \
	CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib" \
		./configure --prefix=${DEPDIR} --disable-doc --host=${HOST} && \
	make -j${NPROCS} && make -j${NPROCS} install
	$(call postbuild_cleanup,$@/${BUILD_DIR})

p11-kit:
	$(call print_msg,"p11-kit")
	mkdir -p $@/${BUILD_DIR} && cd $@/${BUILD_DIR} && \
	CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib" \
		../configure --prefix=${DEPDIR} --host=${HOST} --without-libffi \
		--without-libtasn1 -without-systemd && \
	make -j${NPROCS} && make -j${NPROCS} install
	$(call postbuild_cleanup,$@/${BUILD_DIR})

gnutls: p11-kit openssl nettle gmp libunistring
	$(call print_msg,"gnutls")
	mkdir -p $@/${BUILD_DIR} && cd $@/${BUILD_DIR} && \
	CC="${XCC}" CXX="${XCXX}" CFLAGS="-I${DEPDIR}/include" \
		LDFLAGS="-L${DEPDIR}/lib" ../configure --prefix=${DEPDIR} \
		--exec-prefix=${DEPDIR} GMP_LIBS=${DEPDIR}/lib \
		--host=${HOST} --with-included-libtasn1 \
		--disable-openssl-compatibility \
		--disable-valgrind-tests --disable-guile --enable-local-libopts \
		--without-idn --disable-doc \
		PKG_CONFIG_PATH=${DEPDIR}/lib/pkgconfig && \
	make -j${NPROCS} && make -j${NPROCS} install
	$(call postbuild_cleanup,$@/${BUILD_DIR})

sqlite:
	$(call print_msg,"sqlite")
	mkdir -p $@/${BUILD_DIR} && cd $@/${BUILD_DIR} && \
	CC="${XCC} --static" CFLAGS="-I${DEPDIR}/include" \
		LDFLAGS="-L${DEPDIR}/lib" ../configure --prefix=${DEPDIR} \
		--host=${HOST} --disable-readline && \
	make -j${NPROCS} && make -j${NPROCS} install
	$(call postbuild_cleanup,$@/${BUILD_DIR})

protobuf-c:
	$(call print_msg,"protobuf-c")
	(cd $@ && ./autogen.sh && ./configure --prefix=$(DEPDIR) && \
        $(MAKE) -j${NPROCS} && \
        $(MAKE) -j${NPROCS} install)

amqp:
	$(call print_msg,"amqp")
	(mkdir -p $@/${BUILD_DIR} && cd $@/${BUILD_DIR} && \
        cmake -DCMAKE_INSTALL_PREFIX=$(DEPDIR) -DBUILD_SHARED_LIBS=ON \
                -DBUILD_STATIC_LIBS=ON -DENABLE_SSL_SUPPORT=OFF .. && \
        cmake --build ../ --target install)

clean:
	echo "Cleaning Dependencies."
	rm -rf $(DEPDIR)
	for dir in $(LIST); do	\
		$(MAKE) -C $$dir -f Makefile $@; \
	done

info :
	@echo "Build info for vendor:"
	@cat package.list

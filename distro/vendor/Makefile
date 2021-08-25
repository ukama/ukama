# Copyright (c) 2021-present, Ukama Inc.
#  All rights reserved.

# Makefile fo vendor libs

include ../../config.mk

CURMAKE := $(abspath $(firstword $(MAKEFILE_LIST)))
CURPATH := $(dir $(CURMAKE))

export DEPDIR := $(CURPATH)build

ifndef XGCCPATH
$(error XGCCPATH missing.)
endif

ifndef DISCLEANUP
	DISCLEANUP=0
endif

define postbuild_cleanup
	$(if $(filter ${DISCLEANUP},1), @echo Cleanup disabled, echo Cleaning $(1) && rm -rf $(1))
endef

#BUILD FLAGS
BUILDFLAGS :=

BUILDDIR:=builddir

export XCXX= $(XGCCPATH)$(XGXX)
export XCC = $(XGCCPATH)$(XGCC)
export LD = $(XGCCPATH)$(XLD)

INCFLAGS=-I$(DEPDIR)/include -I. -I$(XGCCPATH)../arm-linux-musleabihf/include

LIBS=

export LDFLAGS = -L$(DEPDIR)/lib -L$(XGCCPATH)../arm-linux-musleabihf/lib

ifeq ($(ARCH), $(ARCHX86))
export CFLAGS = $(INCFLAGS) -fPIC
else
export CFLAGS = $(INCFLAGS)
endif
$(shell mkdir -p $(DEPDIR)/lib)

# List of libs required to be build
DEPLIST:= zlib orcania yder libmicrohttpd jansson ulfius prometheus-client \
       	tomlc openssl curl libuuid nettle gmp libunistring p11-kit gnutls sqlite

.PHONY: $(DEPLIST)

build: $(DEPLIST) 

zlib:
	@echo "Building zlib"
	(cd $@ && ./configure --static --prefix=$(DEPDIR) \
		&& make -j${NPROCS} CC=$(XCC) CFLAGS="${CFLAGS}" && \
	       	make -j${NPROCS} install)

yder:
	$(MAKE) -C $@/src static-install CC=$(XCC) Y_DISABLE_JOURNALD=1 \
	       	DESTDIR=$(DEPDIR)

orcania:
	$(MAKE) -C $@/src static-install CC=$(XCC) DESTDIR=$(DEPDIR)

libmicrohttpd:
	(cd $@ && ./bootstrap && ./configure CC=$(XCC) --host=$(HOST) \
		--disable-curl --disable-https --prefix=$(DEPDIR))
	$(MAKE) -C $@/src install 

ulfius: orcania yder jansson gnutls curl libmicrohttpd
	$(MAKE) -C $@ install CC=$(XCC) DESTDIR=$(DEPDIR)

jansson: 
	(cd $@ && autoreconf -i && ./configure CC=$(XCC) --host=$(HOST) \
		--enable-static --prefix=$(DEPDIR))
	$(MAKE) -C $@ install

prometheus-client:
	mkdir -p $@/prom/build && \
		cd $@/prom/build && \
		cmake ../ -DCMAKE_C_COMPILER=${XCC} -DCMAKE_C_FLAGS="${CFLAGS}" \
	       	-DCMAKE_LD_FLAGS="${LDFLAGS}" \
		-DCMAKE_TOOLCHAIN_FILE=$(CURPATH)cmake_tool_chain_file.cmake && \
		make install DESTDIR=$(DEPDIR)
	@echo "Building promhttp"
	mkdir -p $@/promhttp/build && \
		cd $@/promhttp/build && \
		cmake ../ -DCMAKE_C_COMPILER=${XCC} -DCMAKE_C_FLAGS="${CFLAGS}" \
		-DCMAKE_LD_FLAGS="${LDFLAGS}" \
		-DCMAKE_TOOLCHAIN_FILE=$(CURPATH)cmake_tool_chain_file.cmake && \
	       	make install DESTDIR=$(DEPDIR)
	$(call postbuild_cleanup,$@/prom/build)
	$(call postbuild_cleanup,$@/promhttp/build)
tomlc:
	$(MAKE) -C $@ CC=${XCC} && \
	$(MAKE) -C $@ install prefix=$(DEPDIR)

libuuid:
	mkdir -p $@/${BUILDDIR} && cd $@ && autoreconf -fi && cd ${BUILDDIR} && \
		CC=${XCC} CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib" \
		../configure --prefix=${DEPDIR} --host=${HOST} --build=${BUILD} \
		make && make install
	$(call postbuild_cleanup,$@/${BUILDDIR})

openssl: zlib
	mkdir -p $@/${BUILDDIR} && cd $@/${BUILDDIR} && \
	       CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib \
	       -latomic" \
	       ../Configure ${OPENSSLTARGET} --release --prefix=${DEPDIR} \
	       --openssldir=${DEPDIR} --with-zlib-lib=${DEPDIR}/lib \
	       --with-zlib-include=${DEPDIR}/include && \
       	       make -j${NPROCS} && make -j${NPROCS} install 
	$(call postbuild_cleanup,$@/${BUILDDIR})

curl: openssl
	mkdir -p $@/${BUILDDIR} && cd $@ && autoreconf -fi && cd ${BUILDDIR} && \
		CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib -latomic" \
	       	../configure --prefix=${DEPDIR} --disable-libcurl-option \
		--with-openssl=${DEPDIR} --with-zlib=${DEPDIR} --host=${HOST} && \
		make -j${NPROCS} && make -j${NPROCS} install 
	$(call postbuild_cleanup,$@/${BUILDDIR})

nettle:
	 mkdir -p $@/${BUILDDIR} && cd $@/${BUILDDIR} && \
		 CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib -latomic" \
		 ../configure --prefix=${DEPDIR} --host=${HOST} --enable-mini-gmp && \
		 make -j${NPROCS} && make -j${NPROCS} install 
	$(call postbuild_cleanup,$@/${BUILDDIR})

gmp:
	 mkdir -p $@/${BUILDDIR} && cd $@/${BUILDDIR} && \
                 CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib -latomic" \
		 ../configure --prefix=${DEPDIR} --host=${HOST} && \
                 make -j${NPROCS} && make -j${NPROCS} install
	$(call postbuild_cleanup,$@/${BUILDDIR})

libunistring: 
	mkdir -p $@/${BUILDDIR} && cd $@ && autoreconf -fi && cd ${BUILDDIR} && \
                 CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib" \
		 ../configure --prefix=${DEPDIR} --host=${HOST} && \
                 make -j${NPROCS} && make -j${NPROCS} install
	$(call postbuild_cleanup,$@/${BUILDDIR})

libidn2: libunistring
	mkdir -p $@/${BUILDDIR} && cd $@ && autoreconf -fi && \
                 CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib" \
		 ./configure --prefix=${DEPDIR} --disable-doc --host=${HOST} && \
                 make -j${NPROCS} && make -j${NPROCS} install
	$(call postbuild_cleanup,$@/${BUILDDIR})

p11-kit:
	mkdir -p $@/${BUILDDIR} && cd $@/${BUILDDIR} && \
		CC="${XCC}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib" \
		../configure --prefix=${DEPDIR} --host=${HOST} --without-libffi \
		--without-libtasn1 -without-systemd && \
		make -j${NPROCS} && make -j${NPROCS} install
	$(call postbuild_cleanup,$@/${BUILDDIR})

gnutls: p11-kit openssl nettle gmp libunistring
	mkdir -p $@/${BUILDDIR} && cd $@/${BUILDDIR} && \
		CC="${XCC}" CXX="${XCXX}" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib" \
		../configure --prefix=${DEPDIR} --exec-prefix=${DEPDIR} GMP_LIBS=${DEPDIR}/lib \
		--host=${HOST} --with-included-libtasn1 --disable-openssl-compatibility \
		--disable-valgrind-tests --disable-guile --enable-local-libopts \
		--without-idn --disable-doc PKG_CONFIG_PATH=${DEPDIR}/lib/pkgconfig && \
                make -j${NPROCS} && make -j${NPROCS} install
	$(call postbuild_cleanup,$@/${BUILDDIR})

sqlite:
	mkdir -p $@/${BUILDDIR} && cd $@/${BUILDDIR} && \
		CC="${XCC} --static" CFLAGS="-I${DEPDIR}/include" LDFLAGS="-L${DEPDIR}/lib" \
	       	../configure --prefix=${DEPDIR} --host=${HOST} --disable-readline && \
		make -j${NPROCS} && make -j${NPROCS} install
	$(call postbuild_cleanup,$@/${BUILDDIR})

clean:
	echo "Cleaning Dependencies."
	rm -rf $(DEPDIR)
	for dir in $(DEPLIST); do	\
		$(MAKE) -C $$dir -f Makefile $@; \
	done

info :
	@echo "Build info for vendor:"
	@cat package.list

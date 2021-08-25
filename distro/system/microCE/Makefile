#
# Makefile for the micro container engine.
#

include ../../../config.mk

COMMON_INC=../common/inc

VENDOR_INCLUDE=../vendor/build/include
VENDOR_LIB=../vendor/build/lib

ULFIUS_LIB= -lulfius -lorcania -ljansson -lyder -lmicrohttpd -lgnutls -lnettle -lhogweed -lp11-kit -lunistring -lz
CURL_LIB= -lcurl -lssl -lcrypto

CFLAGS+=-c -Wall -I$(VENDOR_INCLUDE) -I$(COMMON_INC) -I/home/vishal/cdrive/go/src/github.com/ukama/ukamaOS/distro/tools/musl-cross-make/output/arm-linux-musleabihf/include/ -D_REENTRANT -g -O0
LIBS=$(ULFIUS_LIBS) $(CURL_LIB) -ltoml 

LDFLAGS+=-L$(VENDOR_LIB) $(LIBS)

XCC = $(XGCCPATH)$(XGCC)
XLD = $(XGCCPATH)$(XLD)

#all: microCE
all: info

clean:
	rm -f *.o microCE.d; cd ./third_party/ulfius/src && make clean && cd -; 
	cd $(TOML_LOCATION) && make clean



microCE.o: log.o
	$(XCC) $(CFLAGS) src/microCE.c

log.o:
	$(XCC) $(CFLAGS) src/log.c

microCE: microCE.o 
	$(XCC) -o microCE.d microCE.o log.o $(LIBS)

info:
	@echo "Yet to build."

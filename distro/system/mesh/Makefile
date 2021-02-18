#
# Simple Makefile for the Mesh.d
#

PREFIX=/usr/local
MBEDTLS_LOCATION=third_party/mbedtls/
MBEDTLS_INCLUDE=third_party/mbedtls/include
COMMON_INC=../common/inc
LIBS=-lmbedtls -lmbedx509 -lmbedcrypto

CC=gcc
CFLAGS+=-c -Wall -I./inc -I$(MBEDTLS_INCLUDE) -I$(COMMON_INC)  \
	-D_REENTRANT -g -O0 -DTEST_EMBED_CERT
LDFLAGS+=-L$(MBEDTLS_LOCATION)/library

all: mesh.d

clean:
	rm -f *.o mesh.d; cd ./third_party/mbedtls/ && make clean && cd -; 

libmbedtls:
	cd $(MBEDTLS_LOCATION) && $(MAKE) -C library

mesh.o: log.o
	$(CC) $(CFLAGS) src/mesh.c

log.o:
	$(CC) $(CFLAGS) src/log.c

mesh.d: libmbedtls mesh.o
	$(CC) -o mesh.d mesh.o log.o $(LDFLAGS) $(LIBS)
	echo "Done!"



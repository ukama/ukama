#
# Makefile for the micro container engine.
#

PREFIX=/usr/local
ULFIUS_LOCATION=third_party/ulfius/src
ULFIUS_INCLUDE=third_party/ulfius/include
TOML_LOCATION=third_party/tomlc99
COMMON_INC=../common/inc

CC=gcc
CFLAGS+=-c -Wall -I$(ULFIUS_INCLUDE) -I$(TOML_LOCATION) -I$(COMMON_INC) -D_REENTRANT -g -O0
LIBS=-lc -lulfius -lorcania -L$(ULFIUS_LOCATION) -L$(TOML_LOCATION)

all: microCE

clean:
	rm -f *.o microCE.d; cd ./third_party/ulfius/src && make clean && cd -; 
	cd $(TOML_LOCATION) && make clean

libulfius:
	cd $(ULFIUS_LOCATION) && $(MAKE) debug CURLFLAG=1 GNUTLSFLAG=1 JANSSONFLAG=1 WEBSOCKETFLAG=1 YDERFLAG=1 UWSCFLAG=1

toml:
	cd $(TOML_LOCATION) && $(MAKE) DEBUG=1

microCE.o: log.o
	$(CC) $(CFLAGS) src/microCE.c

log.o:
	$(CC) $(CFLAGS) src/log.c

microCE: toml libulfius microCE.o
	$(CC) -o microCE.d microCE.o log.o $(LIBS) $(TOML_LOCATION)/libtoml.a


#
# Simple Makefile for the Mesh.d
#
TARGET_EXEC=mesh.d
PREFIX=/usr/local
MBEDTLS_LOCATION=third_party/mbedtls/
TOML_LOCATION=third_party/tomlc99
MBEDTLS_INCLUDE=third_party/mbedtls/include
COMMON_INC=../common/inc
LIBS=-lulfius -lpthread $(TOML_LOCATION)/libtoml.a

CC=gcc
CFLAGS+=-c -Wall -I./inc -I$(MBEDTLS_INCLUDE) -I$(COMMON_INC)  \
	-I$(TOML_LOCATION)\
	-D_REENTRANT -g -O0 -DTEST_EMBED_CERT
LDFLAGS+=-L$(MBEDTLS_LOCATION)/library


CFILES   = $(wildcard ./src/*.c)
OBJFILES = $(CFILES:.c=.o)

all: $(TARGET_EXEC)

$(TARGET_EXEC): $(OBJFILES) toml
	$(CC) -o $(TARGET_EXEC) $(OBJFILES) $(LDFLAGS) $(LIBS)
	echo "Done."

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

toml:
	cd $(TOML_LOCATION) && $(MAKE) DEBUG=1

clean:
	rm -f $(TARGET_EXEC) $(OBJFILES);
	cd ./third_party/mbedtls/ && make clean && cd -;
	cd $(TOML_LOCATION) && make clean



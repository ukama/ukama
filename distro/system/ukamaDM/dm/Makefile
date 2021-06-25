#!/usr/bin/make -f

#
# Makefile for the device management bootstrap, client and server.
# Depends on the code in core and common
#

CC=/home/vishal/cdrive/work/git/ukama/ukamaOS/distro/tools/musl-cross-make/output/bin/x86_64-linux-musl-gcc
#CC=/usr/bin/musl-gcc
INCS=-I"./core" -I"./common" -I"./ukamaedr" -I"./servIface" -I"." -I"./config"
CFLAGS=$(INCS) -g
CDEFINES=-DLWM2M_BOOTSTRAP -DLWM2M_CLIENT_MODE -DLWM2M_LITTLE_ENDIAN \
-DLWM2M_SUPPORT_JSON 
BOOTSTRAP_CDEFINES=-DLWM2M_BOOTSTRAP_SERVER_MODE -DLWM2M_LITTLE_ENDIAN -DLWM2M_SUPPORT_JSON
SERVER_CDEFINES=-DLWM2M_LITTLE_ENDIAN -DLWM2M_SERVER_MODE -DLWM2M_SUPPORT_JSON 
SRC=$(wildcard core/*.c core/er-coap-13/*.c common/*.c)
CLIENT_SRC=$(wildcard client/*.c)
EDR_SRC=$(wildcard ukamaedr/src/*.c)
SERVIFACE_SRC=$(wildcard servIface/*.c)
BOOTSTRAP_SRC=$(wildcard bootstrap/*.c)
SERVER_SRC=$(wildcard server/*.c)
PARSER_SRC=$(wildcard config/*.c)
BINDIR=build/bin
LINKFLAG = --static
VERSION = v0.0.1

.PHONY: all client container

%.o: %.c
	$(CC) $(CFLAGS) $(CDEFINES) $(EXTRA_FLAGS) -c -o \
	$(OBJDIR)/$(notdir $@) $^

all: client bootstrap server container
	@echo "All Done."

client: mktree $(SRC) $(CLIENT_SRC) $(EDR_SRC) $(PARSER_SRC)
	$(CC) $(CFLAGS) $(CDEFINES) $(SRC) $(CLIENT_SRC) $(EDR_SRC) $(PARSER_SRC) \
	-o $(BINDIR)/client $(LINKFLAG) $(LIBS)
	@echo "Client executable at: $(BINDIR)/client"

bootstrap: mktree $(SRC) $(BOOTSTRAP_SRC)
	$(CC) $(CFLAGS) $(BOOTSTRAP_CDEFINES) $(SRC) $(BOOTSTRAP_SRC) \
	-o $(BINDIR)/bootstrap $(LINKFLAG) $(LIBS)
	@echo "Bootstrap executable at: $(BINDIR)/bootstrap"

server:  mktree $(SRC) $(SERVER_SRC) $(SERVIFACE_SRC)
	$(CC) $(CFLAGS) $(SERVER_CDEFINES) $(SRC) $(SERVER_SRC) $(SERVIFACE_SRC) \
	-o $(BINDIR)/server $(LINKFLAG) $(LIBS)
	@echo "Bootstrap executable at: $(BINDIR)/server"

mktree:
	mkdir -p $(BINDIR)

container:
	@echo Building container
	docker build . -t ukama/lwclient:${VERSION}
clean:
	rm -rf build

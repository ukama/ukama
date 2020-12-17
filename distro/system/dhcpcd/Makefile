#Makefile for dhccpcd

include ../../../config.mk

CURMAKE := $(abspath $(firstword $(MAKEFILE_LIST)))
CURPATH := $(dir $(CURMAKE))

#Source directory
SRCDIR = dhcpcd

#Output
ifndef ROOTFSPATH
override ROOTFSPATH = $(CURPATH)_ukamafs
endif

#Check XGCC PATH
ifndef XGCCPATH
$(error XGCCPATH missing.)
endif

override CC = "$(XGCCPATH)/$(XGCC) --static"
#override CPPFLAGS = "-I$(KERENELHEADERS)"

#Config for Builds
CONFIGFLAGS = --sysconfdir=/etc --mandir=/usr/share/man \
	      --bindir=/bin --sbindir=/sbin --datadir=/var\
	      --libdir=/usr/lib  --localstatedir=/var\
	      --localstatedir=/var --libexecdir=/usr/lib/dhcp \
	      --dbdir=/var/lib/dhcp --rundir=/run \
	      --enable-ipv6 --without-dev --without-udev --prefix=/

#Overiding  Architecture
ifeq ($(ARCH), $(ARCHX86))
override ARCH = $(ARCHX86_64)
endif

.PHONY: subdirs $(SRCDIR) clean info

subdirs: $(SRCDIR)

$(SRCDIR):
	@echo Building dhcpcd for $(ARCH) using $(CC)
	#Configure
	(cd $(SRCDIR) && CC=$(CC) CPPFLAGS=$(CPPFLAGS) ./configure $(CONFIGFLAGS))
	#Make
	$(MAKE) -j$(NPROCS) -C $@ 
	#Install
	$(MAKE) -j$(NPROCS) -C $@ install DESTDIR=$(ROOTFSPATH)

clean:
	@echo Cleaning $(SRCDIR).
	rm -rf $(ROOTFSPATH)
	for dir in $(SRCDIR); do \
                $(MAKE) -j$(NPROCS) -C $$dir -f Makefile $@; \
        done

info: ; 
	$(info [$@] Building $(SRCDIR) $(TARGETBOARD) for $(ARCH) with $(CC) )

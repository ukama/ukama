# simple make file for the pam_cap module

topdir=$(shell pwd)/..
include ../Make.Rules

all: pam_cap.so
	$(MAKE) testlink

install: all
	mkdir -p -m 0755 $(FAKEROOT)$(LIBDIR)/security
	install -m 0755 pam_cap.so $(FAKEROOT)$(LIBDIR)/security

# Note (as the author of much of the Linux-PAM library, I am confident
# that this next line does *not* require -lpam on it.) If you think it
# does, *verify that it does*, and if you observe that it fails as
# written (and you know why it fails), email me and explain why. Thanks!

pam_cap.so: pam_cap.o
	$(LD) -o pam_cap.so $< $(LIBCAPLIB) $(LDFLAGS)

pam_cap.o: pam_cap.c
	$(CC) $(CFLAGS) $(IPATH) -c $< -o $@

test_pam_cap: test_pam_cap.c pam_cap.c
	$(CC) $(CFLAGS) $(IPATH) -o $@ test_pam_cap.c $(LIBCAPLIB) $(LDFLAGS) --static

testlink: test.c pam_cap.o
	$(CC) $(CFLAGS) -o $@ $+ -lpam -ldl $(LIBCAPLIB) $(LDFLAGS)

test: pam_cap.so
	make testlink

sudotest: test test_pam_cap
	sudo ./test_pam_cap root 0x0 0x0 0x0 config=./capability.conf
	sudo ./test_pam_cap root 0x0 0x0 0x0 config=./sudotest.conf
	sudo ./test_pam_cap alpha 0x0 0x0 0x0 config=./capability.conf
	sudo ./test_pam_cap alpha 0x0 0x1 0x80 config=./sudotest.conf
	sudo ./test_pam_cap beta 0x0 0x1 0x0 config=./sudotest.conf
	sudo ./test_pam_cap gamma 0x0 0x0 0x81 config=./sudotest.conf
	sudo ./test_pam_cap delta 0x41 0x80 0x41 config=./sudotest.conf

clean:
	rm -f *.o *.so testlink test_pam_cap *~

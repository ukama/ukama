#
# defines
#
topdir=$(shell pwd)/..
include ../Make.Rules
#

DEPS=../libcap/libcap.a ../libcap/libpsx.a

all: psx_test psx_test_wrap libcap_psx_test libcap_launch_test

$(DEPS):
	make -C ../libcap all

test: run_psx_test run_libcap_psx_test

sudotest: test run_libcap_launch_test run_libcap_launch_test

install: all

run_psx_test: psx_test psx_test_wrap
	./psx_test
	./psx_test_wrap

psx_test: psx_test.c $(DEPS)
	$(CC) $(CFLAGS) $(IPATH) -DNOWRAP $< -o $@ $(LIBPSXLIB)

psx_test_wrap: psx_test.c $(DEPS)
	$(CC) $(CFLAGS) $(IPATH) $< -o $@ $(LIBPSXLIB) -Wl,-wrap,pthread_create

run_libcap_psx_test: libcap_psx_test
	./libcap_psx_test

libcap_psx_test: libcap_psx_test.c $(DEPS)
	$(CC) $(CFLAGS) $(IPATH) $< -o $@ $(LIBCAPLIB) $(LIBPSXLIB) -Wl,-wrap,pthread_create --static

run_libcap_launch_test: libcap_launch_test libcap_psx_launch_test noop
	sudo ./libcap_launch_test
	sudo ./libcap_psx_launch_test

libcap_launch_test: libcap_launch_test.c $(DEPS)
	$(CC) $(CFLAGS) $(IPATH) $< -o $@ $(LIBCAPLIB) --static

# this varies only slightly from the above insofar as it currently
# only links in the pthreads fork support. TODO() we need to change
# the source to do something interesting with pthreads.
libcap_psx_launch_test: libcap_launch_test.c $(DEPS)
	$(CC) $(CFLAGS) $(IPATH) -DWITH_PTHREADS $< -o $@ $(LIBCAPLIB) $(LIBPSXLIB) -Wl,-wrap,pthread_create --static

noop: noop.c
	$(CC) $(CFLAGS) $< -o $@ --static

clean:
	rm -f psx_test psx_test_wrap libcap_psx_test libcap_launch_test *~
	rm -f libcap_launch_test libcap_psx_launch_test core noop

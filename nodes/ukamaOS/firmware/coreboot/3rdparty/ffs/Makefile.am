AM_CFLAGS = -I${srcdir}/clib/ -std=gnu99 -D_GNU_SOURCE -I${srcdir}/ffs/  -fshort-enums -fPIC -D_FILE_OFFSET_BITS=64
AM_LDFLAGS = -L. -fPIC -Wl,-whole-archive  -Wl,-no-whole-archive

bin_PROGRAMS = ecc/ecc fpart/fpart fcp/fcp

noinst_LIBRARIES = libclib.a libffs.a

libclib_a_SOURCES = \
	clib/src/err.c \
	clib/src/misc.c \
	clib/src/ecc.c \
	clib/src/list.c \
	clib/src/list_iter.c \
	clib/src/tree.c \
	clib/src/tree_iter.c \
	clib/src/value.c \
	clib/src/trace_indent.c \
	clib/src/checksum.c

libffs_a_SOURCES = ffs/src/libffs.c ffs/src/libffs2.c

ecc_ecc_SOURCES = ecc/src/main.c
ecc_ecc_LDADD = libffs.a libclib.a

fpart_fpart_LDADD = libffs.a libclib.a
fpart_fpart_SOURCES = \
	fpart/src/cmd_create.c \
	fpart/src/cmd_add.c \
        fpart/src/cmd_delete.c \
	fpart/src/cmd_list.c \
	fpart/src/cmd_erase.c \
        fpart/src/cmd_trunc.c \
	fpart/src/cmd_user.c \
	fpart/src/command.c \
	fpart/src/main.c

fcp_fcp_SOURCES = \
	fcp/src/cmd_copy.c \
	fcp/src/cmd_user.c \
	fcp/src/misc.c \
	fcp/src/cmd_erase.c \
	fcp/src/cmd_read.c \
	fcp/src/cmd_write.c \
	fcp/src/cmd_list.c \
	fcp/src/cmd_trunc.c \
	fcp/src/main.c
fcp_fcp_LDADD = libffs.a libclib.a

EXTRA_DIST = fpart/fpart.sh LICENSE NOTICE

noinst_HEADERS = \
./clib/align.h \
./clib/attribute.h \
./clib/bb_trace.h \
./clib/checksum.h \
./clib/compare.h \
./clib/cunit/ecc.h \
./clib/cunit/splay.h \
./clib/cunit/tree.h \
./clib/ecc.h \
./clib/err.h \
./clib/hash.h \
./clib/ident.h \
./clib/libclib.h \
./clib/list.h \
./clib/list_iter.h \
./clib/max.h \
./clib/min.h \
./clib/misc.h \
./clib/nargs.h \
./clib/raii.h \
./clib/trace_indent.h \
./clib/tree_iter.h \
./clib/type.h \
./clib/value.h \
./clib/version.h \
./clib/builtin.h \
./clib/queue.h \
./clib/tree.h \
./ecc/src/main.h \
./fcp/src/main.h \
./fcp/src/misc.h \
./ffs/ffs.h \
./ffs/libffs2.h \
./ffs/src/ffs-fsp.h \
./ffs/test/test_libffs.h \
./ffs/libffs.h \
./fpart/src/main.h

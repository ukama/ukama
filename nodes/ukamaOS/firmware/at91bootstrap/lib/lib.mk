# Makefile for AT91Bootstrap lib directory
# DIRS				+= lib

LIB:=lib

COBJS-y		+= $(LIB)/string.o
COBJS-y		+= $(LIB)/eabi_utils.o
COBJS-y		+= $(LIB)/div.o
COBJS-y		+= $(LIB)/consttime_memequal.o

COBJS-$(CONFIG_CRC32)	+= $(LIB)/crc32.o
COBJS-$(CONFIG_OF_LIBFDT) += $(LIB)/fdt.o

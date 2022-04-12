# Makefile for AT91Bootstrap ./fs/src directory
# DIRS				+= ./fs/src

FS_FAT:=fs/src

COBJS-$(CONFIG_FATFS)	+=  $(FS_FAT)/ff.o
COBJS-$(CONFIG_FATFS)	+=  $(FS_FAT)/diskio.o
COBJS-$(CONFIG_FATFS)	+=  $(FS_FAT)/option/ccsbcs.o



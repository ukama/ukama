# SPDX-License-Identifier: GPL-2.0
ccflags-y += -I$(src)/ -DWILC_ASIC_A0 -DWILC_DEBUGFS

wilc-objs := cfg80211.o netdev.o mon.o \
			hif.o wlan_cfg.o debugfs.o \
			wlan.o sysfs.o bt.o power.o

obj-$(CONFIG_WILC_SDIO) += wilc-sdio.o
wilc-sdio-objs += $(wilc-objs)
wilc-sdio-objs += sdio.o

obj-$(CONFIG_WILC_SPI) += wilc-spi.o
wilc-spi-objs += $(wilc-objs)
wilc-spi-objs += spi.o 

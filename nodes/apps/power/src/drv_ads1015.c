/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <errno.h>
#include <fcntl.h>
#include <unistd.h>
#include <string.h>
#include <sys/ioctl.h>
#include <linux/i2c-dev.h>

#include "drv_ads1015.h"
#include "usys_log.h"

/* ADS1015 registers */
#define ADS1015_REG_CONVERSION	0x00
#define ADS1015_REG_CONFIG	0x01

/* config bits */
#define ADS1015_OS_SINGLE	0x8000

/* MUX for single-ended */
static uint16_t mux_single_ended(int ch) {

	switch (ch) {
	case 0:  return 0x4000; /* AIN0 vs GND */
	case 1:  return 0x5000; /* AIN1 vs GND */
	case 2:  return 0x6000; /* AIN2 vs GND */
	case 3:  return 0x7000; /* AIN3 vs GND */
	default: return 0x4000;
	}
}

/* PGA: +/- 4.096V */
#define ADS1015_PGA_4_096	0x0200

/* MODE single-shot */
#define ADS1015_MODE_SINGLE	0x0100

/* Data rate 1600 SPS */
#define ADS1015_DR_1600		0x0080

/* Comparator disabled */
#define ADS1015_COMP_DISABLE	0x0003

static int i2c_set_slave(int fd, uint8_t addr) {

	if (ioctl(fd, I2C_SLAVE, addr) < 0) return -1;
	return 0;
}

static int i2c_write_reg16(int fd, uint8_t reg, uint16_t val) {

	uint8_t b[3];

	b[0] = reg;
	b[1] = (uint8_t)((val >> 8) & 0xFF);
	b[2] = (uint8_t)(val & 0xFF);

	if (write(fd, b, 3) != 3) return -1;
	return 0;
}

static int i2c_read_reg16(int fd, uint8_t reg, uint16_t *outVal) {

	uint8_t r = reg;
	uint8_t b[2];
	ssize_t n;

	if (write(fd, &r, 1) != 1) return -1;

	n = read(fd, b, 2);
	if (n != 2) return -1;

	/* ADS1015 returns MSB first for register reads */
	*outVal = (uint16_t)((b[0] << 8) | b[1]);
	return 0;
}

int drv_ads1015_open(Ads1015 *d, const char *dev, int addr7) {

	int fd;

	if (!d) return -1;

	memset(d, 0, sizeof(*d));
	d->fd = -1;

	if (!dev || !*dev) return -1;
	if (addr7 < 0x03 || addr7 > 0x77) return -1;

	fd = open(dev, O_RDWR);
	if (fd < 0) {
		usys_log_error("ads1015: open(%s) failed: %s", dev, strerror(errno));
		return -1;
	}

	if (i2c_set_slave(fd, (uint8_t)addr7) != 0) {
		usys_log_error("ads1015: ioctl(I2C_SLAVE,0x%02x) failed: %s",
		               addr7, strerror(errno));
		close(fd);
		return -1;
	}

	d->fd = fd;
	strncpy(d->dev, dev, sizeof(d->dev)-1);
	d->dev[sizeof(d->dev)-1] = '\0';
	d->addr = (uint8_t)addr7;
	return 0;
}

void drv_ads1015_close(Ads1015 *d) {

	if (!d) return;

	if (d->fd >= 0) close(d->fd);
	memset(d, 0, sizeof(*d));
	d->fd = -1;
}

/* Convert ADS1015 conversion register to volts for PGA +/-4.096V */
static double conv_to_volts_pga_4v096(uint16_t conv) {

	/*
	 * Conversion register holds a signed result left-justified.
	 * For ADS1015, the meaningful bits are [15:4].
	 */
	int16_t raw = (int16_t)conv;
	raw >>= 4; /* sign-extended 12-bit value in LSB */

	/* LSB for +/-4.096V is 2.0mV (4.096/2048) */
	return (double)raw * (4.096 / 2048.0);
}

int drv_ads1015_read_single_ended(Ads1015 *d, int ch, double *outVolts) {

	uint16_t cfg;
	uint16_t conv;
	double v;

	if (!d || !outVolts) return -1;
	if (d->fd < 0) return -1;
	if (ch < 0 || ch > 3) return -1;

	cfg = ADS1015_OS_SINGLE |
	      mux_single_ended(ch) |
	      ADS1015_PGA_4_096 |
	      ADS1015_MODE_SINGLE |
	      ADS1015_DR_1600 |
	      ADS1015_COMP_DISABLE;

	if (i2c_write_reg16(d->fd, ADS1015_REG_CONFIG, cfg) != 0) return -1;

	/* 1600SPS -> ~0.625ms, sleep a bit more */
	usleep(2000);

	if (i2c_read_reg16(d->fd, ADS1015_REG_CONVERSION, &conv) != 0) return -1;

	v = conv_to_volts_pga_4v096(conv);

	/* single-ended should not be negative; clamp */
	if (v < 0) v = 0;

	*outVolts = v;
	return 0;
}


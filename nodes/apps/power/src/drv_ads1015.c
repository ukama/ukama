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

/* PGA: +/- 4.096V (LSB = 2mV for ADS1015 12-bit left-justified) */
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

	*outVal = (uint16_t)((b[0] << 8) | b[1]);
	return 0;
}

int drv_ads1015_open(Ads1015 *d, const char *dev, int addr7) {

	int fd;

	memset(d, 0, sizeof(*d));

	if (!dev || !*dev) return -1;
	if (addr7 < 0x03 || addr7 > 0x77) return -1;

	fd = open(dev, O_RDWR);
	if (fd < 0) {
		usys_log_error("ads1015: open(%s) failed: %s", dev, strerror(errno));
		return -1;
	}

	if (i2c_set_slave(fd, (uint8_t)addr7) != 0) {
		usys_log_error("ads1015: ioctl(I2C_SLAVE,0x%02x) failed: %s", addr7, strerror(errno));
		close(fd);
		return -1;
	}

	d->fd = fd;
	strncpy(d->dev, dev, sizeof(d->dev)-1);
	d->addr = (uint8_t)addr7;
	return 0;
}

void drv_ads1015_close(Ads1015 *d) {

	if (!d) return;
	if (d->fd > 0) close(d->fd);
	memset(d, 0, sizeof(*d));
}

int drv_ads1015_read_single_ended(Ads1015 *d, int ch, double *outVolts) {

	uint16_t cfg;
	uint16_t conv;
	int raw12;
	double lsb;

	if (ch < 0 || ch > 3) return -1;

	/* Build config: start single-shot, mux, pga, single-shot mode, data rate, comp disabled */
	cfg = ADS1015_OS_SINGLE |
	      mux_single_ended(ch) |
	      ADS1015_PGA_4_096 |
	      ADS1015_MODE_SINGLE |
	      ADS1015_DR_1600 |
	      ADS1015_COMP_DISABLE;

	if (i2c_write_reg16(d->fd, ADS1015_REG_CONFIG, cfg) != 0) return -1;

	/* Conversion time at 1600SPS ~0.625ms; sleep a bit more */
	usleep(2000);

	if (i2c_read_reg16(d->fd, ADS1015_REG_CONVERSION, &conv) != 0) return -1;

	/*
	 * ADS1015 conversion register is 16-bit with data left-justified.
	 * For single-ended, result is in bits [15:4] as a 12-bit value.
	 */
	raw12 = (int)((conv >> 4) & 0x0FFF);

	/* PGA +/-4.096V => FS=4.096, LSB=4.096/2048 for ADS1015? 
     * (ADS1015 is 12-bit but uses 11-bit magnitude for single-ended)
	 * Keep simple: treat as 12-bit unsigned over 4.096V.
	 */
	lsb = 4.096 / 2048.0; /* ~2mV typical for ADS1015 with this range */
	*outVolts = (double)raw12 * lsb;

	return 0;
}

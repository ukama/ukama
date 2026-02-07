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

#include "drv_lm75.h"
#include "usys_log.h"

#define LM75_REG_TEMP	0x00

static int i2c_set_slave(int fd, uint8_t addr) {

	if (ioctl(fd, I2C_SLAVE, addr) < 0) return -1;
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

int drv_lm75_open(Lm75 *d, const char *dev, int addr7) {

	int fd;

	memset(d, 0, sizeof(*d));

	if (!dev || !*dev) return -1;
	if (addr7 < 0x03 || addr7 > 0x77) return -1;

	fd = open(dev, O_RDWR);
	if (fd < 0) {
		usys_log_error("lm75: open(%s) failed: %s", dev, strerror(errno));
		return -1;
	}

	if (i2c_set_slave(fd, (uint8_t)addr7) != 0) {
		usys_log_error("lm75: ioctl(I2C_SLAVE,0x%02x) failed: %s", addr7, strerror(errno));
		close(fd);
		return -1;
	}

	d->fd = fd;
	strncpy(d->dev, dev, sizeof(d->dev)-1);
	d->addr = (uint8_t)addr7;

	return 0;
}

void drv_lm75_close(Lm75 *d) {

	if (!d) return;
	if (d->fd > 0) close(d->fd);
	memset(d, 0, sizeof(*d));
}

int drv_lm75_read_temp_c(Lm75 *d, double *outTempC) {

	uint16_t w;
	int16_t raw;
	double t;

	if (i2c_read_reg16(d->fd, LM75_REG_TEMP, &w) != 0) return -1;

	/* LM75: temperature is a signed value; common format: 9-bit/11-bit depending variant.
	 * Keep simple: interpret as 9-bit with 0.5C LSB (raw in top 9 bits).
	 */
	raw = (int16_t)w;
	raw >>= 7; /* keep sign */
	t = (double)raw * 0.5;

	*outTempC = t;
	return 0;
}

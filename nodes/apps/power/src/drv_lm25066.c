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

#include "drv_lm25066.h"
#include "usys_log.h"

/*
 * PMBus commands (common for LM25066-family, but confirm on real hw):
 * - STATUS_WORD         0x79
 * - READ_VIN            0x88
 * - READ_VOUT           0x8B
 * - READ_TEMPERATURE_1  0x8D
 * - MFR_READ_IIN        0xD1
 * - MFR_READ_PIN        0xD2
 * - READ_DIAGNOSTIC_WORD 0xE1 (mfr specific)
 *
 * Best-judgment conversions:
 *  - Assume these are DIRECT format using TI coefficient tables.
 *  - Expose raw words so you can validate with a DMM and tighten later.
 */

#define LM25066_CMD_STATUS_WORD		0x79
#define LM25066_CMD_READ_VIN		0x88
#define LM25066_CMD_READ_VOUT		0x8B
#define LM25066_CMD_READ_TEMP1		0x8D
#define LM25066_CMD_MFR_READ_IIN	0xD1
#define LM25066_CMD_MFR_READ_PIN	0xD2
#define LM25066_CMD_READ_DIAG_WORD	0xE1

static int i2c_set_slave(int fd, uint8_t addr) {

	if (ioctl(fd, I2C_SLAVE, addr) < 0) {
		return -1;
	}
	return 0;
}

static int i2c_read_word(int fd, uint8_t cmd, uint16_t *outWord) {

	uint8_t w = cmd;
	uint8_t b[2];
	ssize_t n;

	if (write(fd, &w, 1) != 1) return -1;

	n = read(fd, b, 2);
	if (n != 2) return -1;

	/* SMBus/PMBus word is LSB first */
	*outWord = (uint16_t)(b[0] | ((uint16_t)b[1] << 8));
	return 0;
}

static double pow10i(int exp10) {

	double v = 1.0;
	int i;

	if (exp10 == 0) return 1.0;
	if (exp10 > 0) {
		for (i = 0; i < exp10; i++) v *= 10.0;
		return v;
	}

	for (i = 0; i < -exp10; i++) v /= 10.0;
	return v;
}

/* DIRECT conversion: X = (Y * 10^-R - B) / M */
static double direct_to_real(double y, int m, int b, int r) {

	double yScaled = y * pow10i(-r);
	return (yScaled - (double)b) / (double)m;
}

/*
 * Coefficients: best judgment defaults.
 * If you later confirm different coefficients, only change these helpers.
 */

static void coeff_vin(int *m, int *b, int *r) {
	/* best guess */
	*m = 22070;
	*b = -1800;
	*r = -2;
}

static void coeff_vout(int *m, int *b, int *r) {
	/* best guess */
	*m = 22070;
	*b = -1800;
	*r = -2;
}

static void coeff_temp(int *m, int *b, int *r) {
	/* best guess */
	*m = 16000;
	*b = 0;
	*r = -3;
}

static int coeff_iin(const Lm25066 *d, int *m, int *b, int *r) {

	long mm;
	long bb;

	if (!d || d->rsMohm <= 0) return -1;

	/*
	 * Best guess based on common LM25066 coefficient table patterns:
	 * CL=GND: M = 13661 * RS(mΩ), B=-5200, R=-2
	 * CL=VDD: M =  6854 * RS(mΩ), B=-3100, R=-2
	 */
	mm = (d->clHigh ? 6854L : 13661L) * (long)d->rsMohm;
	bb = d->clHigh ? -3100L : -5200L;

	/* Keep it simple: if it doesn't fit int, fail (don't silently rescale). */
	if (mm < -32768L || mm > 32767L) {
		usys_log_error("lm25066: iin coeff M out of range (rs=%dmohm clHigh=%d): %ld",
		               d->rsMohm, d->clHigh, mm);
		return -1;
	}

	*m = (int)mm;
	*b = (int)bb;
	*r = -2;
	return 0;
}

static int coeff_pin(const Lm25066 *d, int *m, int *b, int *r) {

	long mm;
	long bb;

	if (!d || d->rsMohm <= 0) return -1;

	/*
	 * Best guess based on common LM25066 coefficient table patterns:
	 * CL=GND: M = 736 * RS(mΩ), B=-3300, R=-2
	 * CL=VDD: M = 369 * RS(mΩ), B=-1900, R=-2
	 */
	mm = (d->clHigh ? 369L : 736L) * (long)d->rsMohm;
	bb = d->clHigh ? -1900L : -3300L;

	if (mm < -32768L || mm > 32767L) {
		usys_log_error("lm25066: pin coeff M out of range (rs=%dmohm clHigh=%d): %ld",
		               d->rsMohm, d->clHigh, mm);
		return -1;
	}

	*m = (int)mm;
	*b = (int)bb;
	*r = -2;
	return 0;
}

int drv_lm25066_open(Lm25066 *d, const char *dev, int addr7, int clHigh, int rsMohm) {

	int fd;

	if (!d) return -1;

	memset(d, 0, sizeof(*d));
	d->fd = -1;

	if (!dev || !*dev) return -1;
	if (addr7 < 0x03 || addr7 > 0x77) return -1;

	fd = open(dev, O_RDWR);
	if (fd < 0) {
		usys_log_error("lm25066: open(%s) failed: %s", dev, strerror(errno));
		return -1;
	}

	if (i2c_set_slave(fd, (uint8_t)addr7) != 0) {
		usys_log_error("lm25066: ioctl(I2C_SLAVE,0x%02x) failed: %s",
		               addr7, strerror(errno));
		close(fd);
		return -1;
	}

	d->fd = fd;
	strncpy(d->dev, dev, sizeof(d->dev)-1);
	d->dev[sizeof(d->dev)-1] = '\0';
	d->addr = (uint8_t)addr7;
	d->clHigh = clHigh ? 1 : 0;
	d->rsMohm = rsMohm;

	return 0;
}

void drv_lm25066_close(Lm25066 *d) {

	if (!d) return;

	if (d->fd >= 0) close(d->fd);
	memset(d, 0, sizeof(*d));
	d->fd = -1;
}

int drv_lm25066_read_sample(Lm25066 *d, Lm25066Sample *s) {

	uint16_t w;
	int m, b, r;

	if (!d || !s) return -1;
	if (d->fd < 0) return -1;

	memset(s, 0, sizeof(*s));

	s->assumedDirect = 1;

	/* VIN (treat as unsigned unless proven otherwise) */
	if (i2c_read_word(d->fd, LM25066_CMD_READ_VIN, &w) != 0) return -1;
	s->rawVin = w;
	coeff_vin(&m, &b, &r);
	s->vinV = direct_to_real((double)(uint16_t)w, m, b, r);

	/* VOUT (treat as unsigned unless proven otherwise) */
	if (i2c_read_word(d->fd, LM25066_CMD_READ_VOUT, &w) != 0) return -1;
	s->rawVout = w;
	coeff_vout(&m, &b, &r);
	s->voutV = direct_to_real((double)(uint16_t)w, m, b, r);

	/* TEMP1 (likely signed) */
	if (i2c_read_word(d->fd, LM25066_CMD_READ_TEMP1, &w) != 0) return -1;
	s->rawTemp = w;
	coeff_temp(&m, &b, &r);
	s->tempC = direct_to_real((double)(int16_t)w, m, b, r);

	/* STATUS_WORD */
	if (i2c_read_word(d->fd, LM25066_CMD_STATUS_WORD, &w) != 0) return -1;
	s->statusWord = w;

	/* DIAGNOSTIC_WORD */
	if (i2c_read_word(d->fd, LM25066_CMD_READ_DIAG_WORD, &w) != 0) return -1;
	s->diagnosticWord = w;

	/* IIN + PIN only if RS provided and coeff fits */
	if (d->rsMohm > 0) {
		if (i2c_read_word(d->fd, LM25066_CMD_MFR_READ_IIN, &w) == 0) {
			s->rawIin = w;
			if (coeff_iin(d, &m, &b, &r) == 0) {
				s->iinA = direct_to_real((double)(int16_t)w, m, b, r);
			} else {
				s->iinA = 0;
			}
		}

		if (i2c_read_word(d->fd, LM25066_CMD_MFR_READ_PIN, &w) == 0) {
			s->rawPin = w;
			if (coeff_pin(d, &m, &b, &r) == 0) {
				s->pinW = direct_to_real((double)(int16_t)w, m, b, r);
			} else {
				s->pinW = 0;
			}
		}
	}

	return 0;
}

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
 * LM25066 PMBus commands we care about:
 * - READ_VIN            0x88
 * - READ_VOUT           0x8B
 * - READ_TEMPERATURE_1  0x8D
 * - STATUS_WORD         0x79
 * - READ_DIAGNOSTIC_WORD (E1h) (manufacturer specific)
 * - READ_IIN            (datasheet uses MFR_READ_IIN 0xD1)
 * - READ_PIN            (datasheet uses MFR_READ_PIN 0xD2)
 *
 * Coefficients (DIRECT format) per TI Table 41 (RS in mΩ). :contentReference[oaicite:2]{index=2}
 *
 * X = (Y * 10^-R - B) / M
 * where Y is 12-bit adc value encoded in 16-bit "DIRECT" word.
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

/* DIRECT conversion per TI Table 41 */
static double direct_to_real(int16_t y, int m, int b, int r) {

	/* X = (Y * 10^-R - B) / M */
	double yScaled = (double)y * pow10i(-r);
	return (yScaled - (double)b) / (double)m;
}

static void coeff_vin(int *m, int *b, int *r) {

	/* READ_VIN: M=22070, B=-1800, R=-2 */
	*m = 22070;
	*b = -1800;
	*r = -2;
}

static void coeff_vout(int *m, int *b, int *r) {

	/* READ_VOUT: M=22070, B=-1800, R=-2 */
	*m = 22070;
	*b = -1800;
	*r = -2;
}

static void coeff_temp(int *m, int *b, int *r) {

	/* READ_TEMPERATURE_1: M=16000, B=0, R=-3 */
	*m = 16000;
	*b = 0;
	*r = -3;
}

static int coeff_iin(const Lm25066 *d, int *m, int *b, int *r) {

	/*
	 * READ_IIN (MFR_READ_IIN):
	 * - if CL=GND: M = 13661 * RS, B=-5200, R=-2
	 * - if CL=VDD: M =  6854 * RS, B=-3100, R=-2
	 * RS in mΩ. :contentReference[oaicite:3]{index=3}
	 *
	 * NOTE: TI warns m must fit -32768..32767; adjust by changing R if needed.
	 * We'll auto-adjust by increasing R (less negative => divide Y less) while keeping scale.
	 */
	long baseM = (d->clHigh ? 6854L : 13661L) * (long)d->rsMohm;
	long baseB = d->clHigh ? -3100L : -5200L;
	int baseR = -2;

	if (d->rsMohm <= 0) return -1;

	/* normalize m into int16 range by shifting decimal (changing R) */
	while (baseM > 32767L) {
		/* divide M by 10 and make R one less negative (e.g. -2 -> -1) */
		baseM = (baseM + 5) / 10;
		baseB = (baseB + 5) / 10;
		baseR += 1;
	}

	*m = (int)baseM;
	*b = (int)baseB;
	*r = baseR;
	return 0;
}

static int coeff_pin(const Lm25066 *d, int *m, int *b, int *r) {

	/*
	 * READ_PIN (MFR_READ_PIN):
	 * - if CL=GND: M = 736 * RS, B=-3300, R=-2
	 * - if CL=VDD: M = 369 * RS, B=-1900, R=-2
	 * RS in mΩ. :contentReference[oaicite:4]{index=4}
	 */
	long baseM = (d->clHigh ? 369L : 736L) * (long)d->rsMohm;
	long baseB = d->clHigh ? -1900L : -3300L;
	int baseR = -2;

	if (d->rsMohm <= 0) return -1;

	while (baseM > 32767L) {
		baseM = (baseM + 5) / 10;
		baseB = (baseB + 5) / 10;
		baseR += 1;
	}

	*m = (int)baseM;
	*b = (int)baseB;
	*r = baseR;
	return 0;
}

int drv_lm25066_open(Lm25066 *d, const char *dev, int addr7, int clHigh, int rsMohm) {

	int fd;

	memset(d, 0, sizeof(*d));

	if (!dev || !*dev) return -1;
	if (addr7 < 0x03 || addr7 > 0x77) return -1;

	fd = open(dev, O_RDWR);
	if (fd < 0) {
		usys_log_error("lm25066: open(%s) failed: %s", dev, strerror(errno));
		return -1;
	}

	if (i2c_set_slave(fd, (uint8_t)addr7) != 0) {
		usys_log_error("lm25066: ioctl(I2C_SLAVE,0x%02x) failed: %s", addr7, strerror(errno));
		close(fd);
		return -1;
	}

	d->fd = fd;
	strncpy(d->dev, dev, sizeof(d->dev)-1);
	d->addr = (uint8_t)addr7;
	d->clHigh = clHigh ? 1 : 0;
	d->rsMohm = rsMohm;

	return 0;
}

void drv_lm25066_close(Lm25066 *d) {

	if (!d) return;
	if (d->fd > 0) close(d->fd);
	memset(d, 0, sizeof(*d));
}

int drv_lm25066_read_sample(Lm25066 *d, Lm25066Sample *s) {

	uint16_t w;
	int m, b, r;
	int16_t y;

	memset(s, 0, sizeof(*s));

	/* VIN */
	if (i2c_read_word(d->fd, LM25066_CMD_READ_VIN, &w) != 0) return -1;
	y = (int16_t)w;
	coeff_vin(&m, &b, &r);
	s->vinV = direct_to_real(y, m, b, r);

	/* VOUT */
	if (i2c_read_word(d->fd, LM25066_CMD_READ_VOUT, &w) != 0) return -1;
	y = (int16_t)w;
	coeff_vout(&m, &b, &r);
	s->voutV = direct_to_real(y, m, b, r);

	/* TEMP1 */
	if (i2c_read_word(d->fd, LM25066_CMD_READ_TEMP1, &w) != 0) return -1;
	y = (int16_t)w;
	coeff_temp(&m, &b, &r);
	s->tempC = direct_to_real(y, m, b, r);

	/* STATUS_WORD */
	if (i2c_read_word(d->fd, LM25066_CMD_STATUS_WORD, &w) != 0) return -1;
	s->statusWord = w;

	/* DIAGNOSTIC_WORD */
	if (i2c_read_word(d->fd, LM25066_CMD_READ_DIAG_WORD, &w) != 0) return -1;
	s->diagnosticWord = w;

	/* IIN + PIN only if RS provided */
	if (d->rsMohm > 0) {
		if (i2c_read_word(d->fd, LM25066_CMD_MFR_READ_IIN, &w) != 0) return -1;
		y = (int16_t)w;
		if (coeff_iin(d, &m, &b, &r) == 0) s->iinA = direct_to_real(y, m, b, r);

		if (i2c_read_word(d->fd, LM25066_CMD_MFR_READ_PIN, &w) != 0) return -1;
		y = (int16_t)w;
		if (coeff_pin(d, &m, &b, &r) == 0) s->pinW = direct_to_real(y, m, b, r);
	}

	return 0;
}

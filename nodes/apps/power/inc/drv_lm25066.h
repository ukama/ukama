/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef DRV_LM25066_H
#define DRV_LM25066_H

#include <stdint.h>

typedef struct {
	int		fd;
	char	dev[64];
	uint8_t	addr;		/* 7-bit */
	int		clHigh;		/* 0=CL=GND, 1=CL=VDD */
	int		rsMohm;		/* sense resistor in milliohms */
} Lm25066;

typedef struct {
	double	vinV;
	double	voutV;
	double	iinA;
	double	pinW;
	double	tempC;

	uint16_t	statusWord;
	uint16_t	diagnosticWord;
} Lm25066Sample;

int drv_lm25066_open(Lm25066 *d, const char *dev, int addr7, int clHigh, int rsMohm);
void drv_lm25066_close(Lm25066 *d);

int drv_lm25066_read_sample(Lm25066 *d, Lm25066Sample *s);

#endif

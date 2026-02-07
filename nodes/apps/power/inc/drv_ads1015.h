/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef DRV_ADS1015_H
#define DRV_ADS1015_H

#include <stdint.h>

typedef struct {
	int		fd;
	char	dev[64];
	uint8_t	addr;		/* 7-bit */
} Ads1015;

/* returns ADC input voltage (V) for single-ended channel 0..3 */
int drv_ads1015_open(Ads1015 *d, const char *dev, int addr7);
void drv_ads1015_close(Ads1015 *d);
int drv_ads1015_read_single_ended(Ads1015 *d, int ch, double *outVolts);

#endif

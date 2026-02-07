/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef DRV_LM75_H
#define DRV_LM75_H

#include <stdint.h>

typedef struct {
	int		fd;
	char	dev[64];
	uint8_t	addr;
} Lm75;

int drv_lm75_open(Lm75 *d, const char *dev, int addr7);
void drv_lm75_close(Lm75 *d);
int drv_lm75_read_temp_c(Lm75 *d, double *outTempC);

#endif

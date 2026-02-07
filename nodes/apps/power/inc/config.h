/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef CONFIG_H
#define CONFIG_H

#include <stdint.h>

typedef struct {

	char	    *listenAddr;
	uint16_t	listenPort;
	uint32_t	sampleMs;
	char	    *boardName;		/* "tower" or "amp" (informational) */

	/* LM25066 (hot-swap / input power telemetry) */
	char	*lm25066Dev;		/* e.g. /dev/i2c-1 */
	int		lm25066Addr;		/* 7-bit */
	int		lm25066ClHigh;		/* 0=CL=GND, 1=CL=VDD */
	int		lm25066RsMohm;		/* sense resistor in milliohms */

	/* LM75 (board temp) */
	char	*lm75Dev;
	int		lm75Addr;

	/* ADS1015 (generic analogs via shunt/amp, dividers, etc) */
	char	*ads1015Dev;
	int		ads1015Addr;
	/* Simple channel mapping: -1 means unused */
	int		adsChVin;
	int		adsChVpa;
	int		adsChAux;
} Config;

int config_load(Config *cfg);
void config_free(Config *cfg);

#endif

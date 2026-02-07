/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef JSON_TYPES_H
#define JSON_TYPES_H

#include <stdint.h>

typedef struct {
	uint64_t	sampleUnixMs;

	char		board[16];

	/* LM25066 (if present): input telemetry */
	int		    haveLm25066;
	double		inVolts;			/* VIN */
	double		outVolts;			/* VOUT (after hot-swap) */
	double		inAmps;				/* IIN (requires RS) */
	double		inWatts;			/* PIN (requires RS) */
	double		hsTempC;			/* READ_TEMPERATURE_1 (if used) */
	uint16_t	statusWord;
	uint16_t	diagnosticWord;

	/* LM75 (if present): board temp */
	int		    haveLm75;
	double		boardTempC;

	/* ADS1015 (if present): raw ADC voltages (as seen at ADC input) */
	int		    haveAds1015;
	double		adcVin;				/* channel mapped "vin" */
	double		adcVpa;				/* channel mapped "vpa" */
	double		adcAux;				/* channel mapped "aux" */

	/* basic health */
	int		ok;
	char		err[128];
} PowerMetrics;

#endif

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef WEB_CLIENT_H_
#define WEB_CLIENT_H_

#include "config.h"
#include "usys_types.h"

typedef struct {
	char	nearUrl[256];
	char	farUrl[256];
	long	ts;
} ReflectorSet;

typedef struct {
	int		ok;
	double	ttfbMs;
	double	totalMs;
	long	httpCode;
	int		stalled;
} ProbeResult;

typedef struct {
	int		ok;
	double	mbps;
	double	seconds;
	long	httpCode;
} TransferResult;

int wc_init(void);
void wc_cleanup(void);

int wc_fetch_reflectors(Config *config, ReflectorSet *set);

int wc_probe_ping(Config *config, const char *baseUrl, ProbeResult *out);

int wc_download_blob(Config *config, const char *baseUrl, int bytes, TransferResult *out);
int wc_upload_echo(Config *config, const char *baseUrl, int bytes, TransferResult *out);

#endif /* WEB_CLIENT_H_ */

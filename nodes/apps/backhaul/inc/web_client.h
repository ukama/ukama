/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef WEB_CLIENT_H
#define WEB_CLIENT_H

#include "backhauld.h"
#include "config.h"

int  wc_init(void);
void wc_cleanup(void);

/* bootstrap */
int wc_fetch_reflectors(Config *config, ReflectorSet *out);

/* probes */
int wc_probe_ping(Config *config, const char *baseUrl, ProbeResult *out);

/* transfers */
int wc_download_blob(Config *config, const char *baseUrl, int bytes, TransferResult *out);
int wc_upload_echo(Config *config, const char *baseUrl, int bytes, TransferResult *out);

#endif /* WEB_CLIENT_H */

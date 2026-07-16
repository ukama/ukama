/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_HUB_H_
#define ULAB_HUB_H_

#include <stdint.h>
#include <stdio.h>

#include "ulab.h"

typedef struct {
    char url[ULAB_MAX_URL];
    FILE *logf;
    int curl_ready;
} hub_client_t;

int hub_init(hub_client_t *hub,
             const char *url,
             const char *run_dir,
             ulab_error_t *err);

void hub_close(hub_client_t *hub);

int hub_tar_version_exists(hub_client_t *hub,
                           const char *app,
                           const char *version,
                           int *exists,
                           uint64_t *size_bytes,
                           char *artifact_url,
                           size_t artifact_url_len,
                           ulab_error_t *err);

#endif /* ULAB_HUB_H_ */

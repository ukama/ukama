/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <time.h>
#include <stdio.h>

#include "algo_micro_probe.h"
#include "web_client.h"
#include "usys_log.h"

int algo_micro_probe_run(Config *config, MetricsStore *store, void *unused) {
    (void)unused;

    char nearUrl[256] = {0};
    char farUrl[256]  = {0};
    long ts = 0;

    /* Use near as base for micro probe */
    if (!metrics_store_get_reflectors(store, nearUrl, sizeof(nearUrl), farUrl, sizeof(farUrl), &ts)) {
        return STATUS_NOK;
    }

    ProbeResult pr;
    memset(&pr, 0, sizeof(pr));
    wc_probe_ping(config, nearUrl, &pr);

    MicroSample s;
    memset(&s, 0, sizeof(s));
    s.ts = time(NULL);
    s.ttfbMs = pr.ttfbMs;
    s.ok = pr.ok;
    s.stalled = pr.stalled;

    metrics_store_add_micro(store, s);
    return STATUS_OK;
}

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <time.h>

#include "algo_multi_reflector.h"
#include "web_client.h"
#include "usys_log.h"

static void add_probe(MetricsStore *store,
                      int isNear,
                      ProbeResult pr) {

    MicroSample s;
    memset(&s, 0, sizeof(s));
    s.ts = time(NULL);
    s.ttfbMs = pr.ttfbMs;
    s.ok = pr.ok;
    s.stalled = pr.stalled;

    if (isNear) metrics_store_add_near(store, s);
    else        metrics_store_add_far(store, s);
}

int algo_multi_reflector_run(Config *config,
                             MetricsStore *store,
                             void *unused) {

    (void)unused;

    char nearUrl[256] = {0};
    char farUrl[256]  = {0};
    long ts = 0;

    if (!metrics_store_get_reflectors(store, nearUrl, sizeof(nearUrl), farUrl, sizeof(farUrl), &ts)) {
        return STATUS_NOK;
    }

    if (!nearUrl[0] || !farUrl[0]) return STATUS_NOK;

    ProbeResult pr;

    memset(&pr, 0, sizeof(pr));
    wc_probe_ping(config, nearUrl, &pr);
    add_probe(store, 1, pr);

    memset(&pr, 0, sizeof(pr));
    wc_probe_ping(config, farUrl, &pr);
    add_probe(store, 0, pr);

    return STATUS_OK;
}

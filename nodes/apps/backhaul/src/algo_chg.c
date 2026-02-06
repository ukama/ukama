/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <time.h>

#include "algo_chg.h"
#include "web_client.h"
#include "usys_log.h"

static int pick_bytes_for_target(Config *config,
                                 double mbps,
                                 int minBytes,
                                 int maxBytes) {

    if (mbps <= 0.1) return minBytes;

    double mbytes = (mbps / 8.0) * (double)config->chgTargetSec;
    int bytes = (int)(mbytes * 1000000.0);

    if (bytes < minBytes) bytes = minBytes;
    if (bytes > maxBytes) bytes = maxBytes;

    return bytes;
}

int algo_chg_run(Config *config, MetricsStore *store, void *unused) {
    (void)unused;

    char nearUrl[256] = {0};
    char farUrl[256]  = {0};
    long ts = 0;

    if (!metrics_store_get_reflectors(store,
                                      nearUrl,
                                      sizeof(nearUrl),
                                      farUrl,
                                      sizeof(farUrl), &ts)) {
        return STATUS_NOK;
    }
    if (!nearUrl[0]) return STATUS_NOK;

    BackhaulMetrics snap = metrics_store_get_snapshot(store);

    int dlBytes = pick_bytes_for_target(config,
                                        snap.dlGoodputMbps,
                                        config->chgMinBytes,
                                        config->chgMaxBytes);
    int ulBytes = pick_bytes_for_target(config,
                                        snap.ulGoodputMbps,
                                        config->chgMinBytes,
                                        config->chgMaxBytes);

    TransferResult dl, ul;
    memset(&dl, 0, sizeof(dl));
    memset(&ul, 0, sizeof(ul));

    if (config->chgWarmupBytes > 0) {
        TransferResult warm;
        memset(&warm, 0, sizeof(warm));
        wc_download_blob(config, nearUrl, config->chgWarmupBytes, &warm);
    }

    double dlVals[16] = {0};
    double ulVals[16] = {0};
    int n = config->chgSamples;
    if (n > 16) n = 16;
    if (n <= 0) n = 1;

    int dlOk = 0, ulOk = 0;

    for (int i=0; i<n; i++) {
        memset(&dl, 0, sizeof(dl));
        wc_download_blob(config, nearUrl, dlBytes, &dl);
        if (dl.ok) dlVals[dlOk++] = dl.mbps;
    }

    for (int i=0; i<n; i++) {
        memset(&ul, 0, sizeof(ul));
        wc_upload_echo(config, nearUrl, ulBytes, &ul);
        if (ul.ok) ulVals[ulOk++] = ul.mbps;
    }

    double dlMed = 0.0, ulMed = 0.0;

    if (dlOk > 0) {
        for (int i=0; i<dlOk-1; i++) {
            for (int j=i+1; j<dlOk; j++) {
                if (dlVals[j] < dlVals[i]) { double t=dlVals[i]; dlVals[i]=dlVals[j]; dlVals[j]=t; }
            }
        }
        dlMed = (dlOk % 2) ? dlVals[dlOk/2] : (dlVals[dlOk/2 - 1] + dlVals[dlOk/2]) / 2.0;
    }

    if (ulOk > 0) {
        for (int i=0; i<ulOk-1; i++) {
            for (int j=i+1; j<ulOk; j++) {
                if (ulVals[j] < ulVals[i]) { double t=ulVals[i]; ulVals[i]=ulVals[j]; ulVals[j]=t; }
            }
        }
        ulMed = (ulOk % 2) ? ulVals[ulOk/2] : (ulVals[ulOk/2 - 1] + ulVals[ulOk/2]) / 2.0;
    }

    ChgSample cs;
    memset(&cs, 0, sizeof(cs));
    cs.ts = time(NULL);
    cs.ok = (dlOk > 0 || ulOk > 0) ? 1 : 0;
    cs.dlMbps = dlMed;
    cs.ulMbps = ulMed;
    cs.dlSec = dl.seconds;
    cs.ulSec = ul.seconds;

    metrics_store_add_chg(store, cs);

    return STATUS_OK;
}

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <time.h>

#include "classifier.h"
#include "usys_log.h"

static double clamp(double v, double lo, double hi) {
    if (v < lo) return lo;
    if (v > hi) return hi;
    return v;
}

static double absd(double x) {
    return x < 0 ? -x : x;
}

static int count_recent_failures_locked(MetricsStore *store, int limit) {

    int fails = 0;
    int cnt = store->microCount;

    if (cnt <= 0) return 0;

    if (limit > cnt) limit = cnt;

    for (int i=0; i<limit; i++) {
        int idx = (store->microHead - 1 - i);
        while (idx < 0) idx += store->microCap;
        idx = idx % store->microCap;

        if (!store->microSamples[idx].ok) fails++;
    }

    return fails;
}

static double last_chg_dl_locked(MetricsStore *store) {

    if (store->chgCount <= 0) return 0.0;

    int idx = store->chgHead - 1;
    while (idx < 0) idx += store->chgCap;
    idx = idx % store->chgCap;

    return store->chgSamples[idx].dlMbps;
}

static double last_chg_ul_locked(MetricsStore *store) {

    if (store->chgCount <= 0) return 0.0;

    int idx = store->chgHead - 1;
    while (idx < 0) idx += store->chgCap;
    idx = idx % store->chgCap;

    return store->chgSamples[idx].ulMbps;
}

/* cap detection: if last K dl samples are within +/- capStabilityPct of their median */
static int detect_cap_locked(Config *config, MetricsStore *store, double *capOut) {

    int k = store->chgCount;
    if (k < 5) return 0;
    if (k > 10) k = 10;

    double vals[10] = {0};
    for (int i=0; i<k; i++) {
        int idx = (store->chgHead - 1 - i);
        while (idx < 0) idx += store->chgCap;
        idx = idx % store->chgCap;
        vals[i] = store->chgSamples[idx].dlMbps;
    }

    for (int i=0; i<k-1; i++) {
        for (int j=i+1; j<k; j++) {
            if (vals[j] < vals[i]) { double t=vals[i]; vals[i]=vals[j]; vals[j]=t; }
        }
    }

    double med = (k % 2) ? vals[k/2] : (vals[k/2 - 1] + vals[k/2]) / 2.0;
    if (med < 1.0) return 0;

    double pct = (double)config->capStabilityPct;

    for (int i=0; i<k; i++) {
        double diffPct = (absd(vals[i] - med) / med) * 100.0;
        if (diffPct > pct) return 0;
    }

    *capOut = med;
    return 1;
}

void classifier_run(Config *config, MetricsStore *store) {

    if (!config || !store) return;

    BackhaulAggregates a;
    (void)metrics_store_get_aggregates(store, &a);

    pthread_mutex_lock(&store->lock);

    /* publish aggregates into status fields (single source of truth) */
    store->metrics.nearTtfbMedianMs = a.nearMed;
    store->metrics.nearTtfbP95Ms    = a.nearP95;
    store->metrics.nearTtfbP99Ms    = a.nearP99;

    store->metrics.farTtfbMedianMs = a.farMed;
    store->metrics.farTtfbP95Ms    = a.farP95;
    store->metrics.farTtfbP99Ms    = a.farP99;

    store->metrics.probeSuccessRatePct = a.microSuccPct;
    store->metrics.stallRatePct        = a.microStallPct;

    int recentFails = count_recent_failures_locked(store, config->downConsecFails);

    if (recentFails >= config->downConsecFails) {
        store->metrics.consecFails++;
        store->metrics.consecOk = 0;
    } else {
        store->metrics.consecOk++;
        if (store->metrics.consecOk >= config->recoverConsecOk) {
            store->metrics.consecFails = 0;
        }
    }

    if (store->metrics.consecFails > 0) {
        store->metrics.backhaulState = BACKHAUL_STATE_DOWN;
        store->metrics.confidence = 0.9;
        store->metrics.linkGuess = BACKHAUL_LINK_UNKNOWN;
        store->metrics.lastClassifyTs = time(NULL);
        pthread_mutex_unlock(&store->lock);
        return;
    }

    double cap = 0.0;
    int capped = detect_cap_locked(config, store, &cap);
    if (capped) {
        store->metrics.backhaulState = BACKHAUL_STATE_CAPPED;
        store->metrics.capDetectedMbps = cap;
        store->metrics.confidence = 0.85;
    }

    double stallPct = store->metrics.stallRatePct;
    double nearMed  = store->metrics.nearTtfbMedianMs;
    double nearP99  = store->metrics.nearTtfbP99Ms;

    if (!capped) {
        if (stallPct >= 2.0 || nearP99 >= (double)(config->stallThresholdMs * 2)) {
            store->metrics.backhaulState = BACKHAUL_STATE_DEGRADED;
        } else {
            store->metrics.backhaulState = BACKHAUL_STATE_GOOD;
        }
    }

    BackhaulLinkGuess guess = BACKHAUL_LINK_UNKNOWN;
    double conf = 0.4;

    if (nearMed >= 450.0) {
        guess = BACKHAUL_LINK_SAT_GEO_LIKE;
        conf = 0.8;
    } else if (nearMed >= 80.0 && nearMed <= 180.0 && (nearP99 - nearMed) > 80.0) {
        guess = BACKHAUL_LINK_SAT_LEO_LIKE;
        conf = 0.65;
    } else if (nearMed < 50.0) {
        guess = BACKHAUL_LINK_TERRESTRIAL_LIKE;
        conf = 0.7;
    }

    double farMed = store->metrics.farTtfbMedianMs;
    if (farMed > 0.0 && nearMed > 0.0 && farMed > (nearMed * 2.5)) {
        conf = clamp(conf + 0.1, 0.0, 0.95);
    }

    store->metrics.linkGuess = guess;

    if (capped && guess == BACKHAUL_LINK_TERRESTRIAL_LIKE) {
        conf = 0.6;
    }

    store->metrics.confidence = conf;

    store->metrics.dlGoodputMbps = last_chg_dl_locked(store);
    store->metrics.ulGoodputMbps = last_chg_ul_locked(store);

    store->metrics.lastClassifyTs = time(NULL);

    pthread_mutex_unlock(&store->lock);
}

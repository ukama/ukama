/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <time.h>

#include "metrics_store.h"
#include "usys_mem.h"

static void ring_push_micro(MicroSample *arr, int cap, int *head, int *count, MicroSample s) {
    if (cap <= 0) return;
    arr[*head] = s;
    *head = (*head + 1) % cap;
    if (*count < cap) (*count)++;
}

static void ring_push_chg(ChgSample *arr, int cap, int *head, int *count, ChgSample s) {
    if (cap <= 0) return;
    arr[*head] = s;
    *head = (*head + 1) % cap;
    if (*count < cap) (*count)++;
}

int metrics_store_init(MetricsStore *store, int microCap, int multiCap, int chgCap) {

    if (!store) return USYS_FALSE;
    memset(store, 0, sizeof(*store));

    pthread_mutex_init(&store->lock, NULL);

    store->microCap = microCap;
    store->nearCap  = multiCap;
    store->farCap   = multiCap;
    store->chgCap   = chgCap;

    store->microSamples = (MicroSample *)usys_calloc(microCap, sizeof(MicroSample));
    store->nearSamples  = (MicroSample *)usys_calloc(multiCap, sizeof(MicroSample));
    store->farSamples   = (MicroSample *)usys_calloc(multiCap, sizeof(MicroSample));
    store->chgSamples   = (ChgSample *)usys_calloc(chgCap, sizeof(ChgSample));

    if (!store->microSamples || !store->nearSamples || !store->farSamples || !store->chgSamples) {
        metrics_store_free(store);
        return USYS_FALSE;
    }

    memset(&store->metrics, 0, sizeof(store->metrics));
    store->metrics.backhaulState = BACKHAUL_STATE_UNKNOWN;
    store->metrics.linkGuess = BACKHAUL_LINK_UNKNOWN;
    store->metrics.confidence = 0.0;

    store->diagRequest = BACKHAUL_DIAG_NONE;

    memset(store->reflectorNearUrl, 0, sizeof(store->reflectorNearUrl));
    memset(store->reflectorFarUrl, 0, sizeof(store->reflectorFarUrl));
    store->reflectorTs = 0;

    return USYS_TRUE;
}

void metrics_store_free(MetricsStore *store) {

    if (!store) return;

    pthread_mutex_destroy(&store->lock);

    if (store->microSamples) usys_free(store->microSamples);
    if (store->nearSamples)  usys_free(store->nearSamples);
    if (store->farSamples)   usys_free(store->farSamples);
    if (store->chgSamples)   usys_free(store->chgSamples);

    memset(store, 0, sizeof(*store));
}

void metrics_store_set_reflectors(MetricsStore *store,
                                  const char *nearUrl,
                                  const char *farUrl,
                                  long ts) {

    if (!store) return;

    pthread_mutex_lock(&store->lock);

    if (nearUrl) {
        memset(store->reflectorNearUrl, 0, sizeof(store->reflectorNearUrl));
        strncpy(store->reflectorNearUrl, nearUrl, sizeof(store->reflectorNearUrl) - 1);
    }
    if (farUrl) {
        memset(store->reflectorFarUrl, 0, sizeof(store->reflectorFarUrl));
        strncpy(store->reflectorFarUrl, farUrl, sizeof(store->reflectorFarUrl) - 1);
    }
    store->reflectorTs = ts;

    pthread_mutex_unlock(&store->lock);
}

int metrics_store_get_reflectors(MetricsStore *store,
                                char *nearUrl, size_t nearLen,
                                char *farUrl,  size_t farLen,
                                long *ts) {

    if (!store) return USYS_FALSE;

    pthread_mutex_lock(&store->lock);

    if (nearUrl && nearLen > 0) {
        memset(nearUrl, 0, nearLen);
        strncpy(nearUrl, store->reflectorNearUrl, nearLen - 1);
    }
    if (farUrl && farLen > 0) {
        memset(farUrl, 0, farLen);
        strncpy(farUrl, store->reflectorFarUrl, farLen - 1);
    }
    if (ts) *ts = store->reflectorTs;

    int ok = (store->reflectorNearUrl[0] && store->reflectorFarUrl[0]) ? USYS_TRUE : USYS_FALSE;

    pthread_mutex_unlock(&store->lock);
    return ok;
}

void metrics_store_request_diag(MetricsStore *store, BackhaulDiagRequest req) {
    if (!store) return;
    pthread_mutex_lock(&store->lock);
    store->diagRequest = req;
    pthread_mutex_unlock(&store->lock);
}

BackhaulDiagRequest metrics_store_take_diag_request(MetricsStore *store) {
    if (!store) return BACKHAUL_DIAG_NONE;
    pthread_mutex_lock(&store->lock);
    BackhaulDiagRequest r = store->diagRequest;
    store->diagRequest = BACKHAUL_DIAG_NONE;
    pthread_mutex_unlock(&store->lock);
    return r;
}

void metrics_store_add_micro(MetricsStore *store, MicroSample s) {
    pthread_mutex_lock(&store->lock);
    ring_push_micro(store->microSamples, store->microCap, &store->microHead, &store->microCount, s);
    store->metrics.lastMicroTs = s.ts;
    pthread_mutex_unlock(&store->lock);
}

void metrics_store_add_near(MetricsStore *store, MicroSample s) {
    pthread_mutex_lock(&store->lock);
    ring_push_micro(store->nearSamples, store->nearCap, &store->nearHead, &store->nearCount, s);
    store->metrics.lastMultiTs = s.ts;
    pthread_mutex_unlock(&store->lock);
}

void metrics_store_add_far(MetricsStore *store, MicroSample s) {
    pthread_mutex_lock(&store->lock);
    ring_push_micro(store->farSamples, store->farCap, &store->farHead, &store->farCount, s);
    store->metrics.lastMultiTs = s.ts;
    pthread_mutex_unlock(&store->lock);
}

void metrics_store_add_chg(MetricsStore *store, ChgSample s) {
    pthread_mutex_lock(&store->lock);
    ring_push_chg(store->chgSamples, store->chgCap, &store->chgHead, &store->chgCount, s);
    store->metrics.lastChgTs = s.ts;
    pthread_mutex_unlock(&store->lock);
}

void metrics_store_set_diag(MetricsStore *store, const char *name) {
    pthread_mutex_lock(&store->lock);
    store->metrics.lastDiagTs = time(NULL);
    memset(store->metrics.lastDiagName, 0, sizeof(store->metrics.lastDiagName));
    if (name) strncpy(store->metrics.lastDiagName, name, sizeof(store->metrics.lastDiagName) - 1);
    pthread_mutex_unlock(&store->lock);
}

BackhaulMetrics metrics_store_get_snapshot(MetricsStore *store) {
    BackhaulMetrics m;
    memset(&m, 0, sizeof(m));
    pthread_mutex_lock(&store->lock);
    m = store->metrics;
    pthread_mutex_unlock(&store->lock);
    return m;
}

/* --- aggregates --- */

static int cmp_double(const void *a, const void *b) {
    const double da = *(const double *)a;
    const double db = *(const double *)b;
    if (da < db) return -1;
    if (da > db) return 1;
    return 0;
}

static void compute_percentiles_locked(MicroSample *arr, int cap, int head, int count,
                                       double *median, double *p95, double *p99,
                                       double *successPct, double *stallPct) {

    *median = 0; *p95 = 0; *p99 = 0;
    *successPct = 0; *stallPct = 0;

    if (count <= 0) return;

    double *vals = (double *)usys_calloc(count, sizeof(double));
    if (!vals) return;

    int okCount = 0;
    int stallCount = 0;

    for (int i=0; i<count; i++) {
        int idx = (head - count + i);
        while (idx < 0) idx += cap;
        idx = idx % cap;

        vals[i] = arr[idx].ttfbMs;
        if (arr[idx].ok) okCount++;
        if (arr[idx].stalled) stallCount++;
    }

    qsort(vals, count, sizeof(double), cmp_double);

    int mid = count/2;
    *median = (count % 2) ? vals[mid] : (vals[mid-1] + vals[mid]) / 2.0;

    int i95 = (int)((0.95 * (count - 1)) + 0.5);
    int i99 = (int)((0.99 * (count - 1)) + 0.5);

    if (i95 < 0) i95 = 0;
    if (i95 >= count) i95 = count-1;
    if (i99 < 0) i99 = 0;
    if (i99 >= count) i99 = count-1;

    *p95 = vals[i95];
    *p99 = vals[i99];

    *successPct = (100.0 * okCount) / (double)count;
    *stallPct   = (100.0 * stallCount) / (double)count;

    usys_free(vals);
}

int metrics_store_get_aggregates(MetricsStore *store, BackhaulAggregates *out) {

    if (!store || !out) return USYS_FALSE;
    memset(out, 0, sizeof(*out));

    pthread_mutex_lock(&store->lock);

    out->microCount = store->microCount;
    out->nearCount  = store->nearCount;
    out->farCount   = store->farCount;
    out->chgCount   = store->chgCount;

    compute_percentiles_locked(store->microSamples, store->microCap,
                               store->microHead, store->microCount,
                               &out->microMed, &out->microP95, &out->microP99,
                               &out->microSuccPct, &out->microStallPct);

    compute_percentiles_locked(store->nearSamples, store->nearCap,
                               store->nearHead, store->nearCount,
                               &out->nearMed, &out->nearP95, &out->nearP99,
                               &out->nearSuccPct, &out->nearStallPct);

    compute_percentiles_locked(store->farSamples, store->farCap,
                               store->farHead, store->farCount,
                               &out->farMed, &out->farP95, &out->farP99,
                               &out->farSuccPct, &out->farStallPct);

    pthread_mutex_unlock(&store->lock);
    return USYS_TRUE;
}

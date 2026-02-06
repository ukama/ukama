/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

#ifndef METRICS_STORE_H
#define METRICS_STORE_H

#include <pthread.h>
#include <stddef.h>

#include "backhauld.h"
#include "usys_api.h"

typedef enum {
    BACKHAUL_DIAG_NONE = 0,
    BACKHAUL_DIAG_CHG,
    BACKHAUL_DIAG_PARALLEL,
    BACKHAUL_DIAG_BUFFERBLOAT
} BackhaulDiagRequest;

typedef struct {
    /* computed aggregates from ring buffers */
    double microMed, microP95, microP99, microSuccPct, microStallPct;
    double nearMed,  nearP95,  nearP99,  nearSuccPct,  nearStallPct;
    double farMed,   farP95,   farP99,   farSuccPct,   farStallPct;

    int microCount;
    int nearCount;
    int farCount;
    int chgCount;
} BackhaulAggregates;

typedef struct MetricsStore {
    pthread_mutex_t lock;

    /* ring buffers */
    MicroSample *microSamples;
    MicroSample *nearSamples;
    MicroSample *farSamples;
    ChgSample   *chgSamples;

    int microCap, nearCap, farCap, chgCap;
    int microHead, nearHead, farHead, chgHead;
    int microCount, nearCount, farCount, chgCount;

    /* published current metrics */
    BackhaulMetrics metrics;

    /* static reflectors */
    char reflectorNearUrl[256];
    char reflectorFarUrl[256];
    long reflectorTs;

    /* Option-A diagnostics request flag */
    BackhaulDiagRequest diagRequest;
} MetricsStore;

int  metrics_store_init(MetricsStore *store, int microCap, int multiCap, int chgCap);
void metrics_store_free(MetricsStore *store);

void metrics_store_add_micro(MetricsStore *store, MicroSample s);
void metrics_store_add_near(MetricsStore *store, MicroSample s);
void metrics_store_add_far(MetricsStore *store, MicroSample s);
void metrics_store_add_chg(MetricsStore *store, ChgSample s);

void metrics_store_set_diag(MetricsStore *store, const char *name);

BackhaulMetrics metrics_store_get_snapshot(MetricsStore *store);

/* reflectors */
void metrics_store_set_reflectors(MetricsStore *store, const char *nearUrl, const char *farUrl, long ts);
int  metrics_store_get_reflectors(MetricsStore *store,
                                 char *nearUrl, size_t nearLen,
                                 char *farUrl,  size_t farLen,
                                 long *ts);

/* diagnostics Option-A */
void metrics_store_request_diag(MetricsStore *store, BackhaulDiagRequest req);
BackhaulDiagRequest metrics_store_take_diag_request(MetricsStore *store);

/* aggregates helper (used by classifier + json_serdes) */
int metrics_store_get_aggregates(MetricsStore *store, BackhaulAggregates *out);

#endif /* METRICS_STORE_H_ */

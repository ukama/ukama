/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "json_serdes.h"
#include "metrics_store.h"
#include "backhauld.h"

static const char* state_str(BackhaulState s) {
    
    switch (s) {
    case BACKHAUL_STATE_GOOD:     return "GOOD";
    case BACKHAUL_STATE_DEGRADED: return "DEGRADED";
    case BACKHAUL_STATE_DOWN:     return "DOWN";
    case BACKHAUL_STATE_CAPPED:   return "CAPPED";
    default:                      return "UNKNOWN";
    }
}

static const char* link_str(BackhaulLinkGuess g) {

    switch (g) {
    case BACKHAUL_LINK_TERRESTRIAL_LIKE: return "TERRESTRIAL_LIKE";
    case BACKHAUL_LINK_SAT_LEO_LIKE:     return "SAT_LEO_LIKE";
    case BACKHAUL_LINK_SAT_GEO_LIKE:     return "SAT_GEO_LIKE";
    case BACKHAUL_LINK_CELLULAR_LIKE:    return "CELLULAR_LIKE";
    default:                             return "UNKNOWN";
    }
}

json_t* json_backhaul_status(MetricsStore *store) {

    if (!store) return json_object();

    BackhaulMetrics m = metrics_store_get_snapshot(store);

    json_t *o = json_object();

    json_object_set_new(o, "backhaulState", json_string(state_str(m.backhaulState)));
    json_object_set_new(o, "linkGuess",     json_string(link_str(m.linkGuess)));
    json_object_set_new(o, "confidence",    json_real(m.confidence));

    json_object_set_new(o, "dlGoodputMbps",              json_real(m.dlGoodputMbps));
    json_object_set_new(o, "ulGoodputMbps",              json_real(m.ulGoodputMbps));
    json_object_set_new(o, "bufferbloatInflationFactor", json_real(m.bufferbloatInflationFactor));
    json_object_set_new(o, "capDetectedMbps",            json_real(m.capDetectedMbps));

    json_object_set_new(o, "nearTtfbMedianMs", json_real(m.nearTtfbMedianMs));
    json_object_set_new(o, "nearTtfbP95Ms",    json_real(m.nearTtfbP95Ms));
    json_object_set_new(o, "nearTtfbP99Ms",    json_real(m.nearTtfbP99Ms));

    json_object_set_new(o, "farTtfbMedianMs", json_real(m.farTtfbMedianMs));
    json_object_set_new(o, "farTtfbP95Ms",    json_real(m.farTtfbP95Ms));
    json_object_set_new(o, "farTtfbP99Ms",    json_real(m.farTtfbP99Ms));

    json_object_set_new(o, "probeSuccessRatePct", json_real(m.probeSuccessRatePct));
    json_object_set_new(o, "stallRatePct",        json_real(m.stallRatePct));

    json_object_set_new(o, "lastMicroTs",    json_integer(m.lastMicroTs));
    json_object_set_new(o, "lastMultiTs",    json_integer(m.lastMultiTs));
    json_object_set_new(o, "lastChgTs",      json_integer(m.lastChgTs));
    json_object_set_new(o, "lastClassifyTs", json_integer(m.lastClassifyTs));

    json_object_set_new(o, "lastDiagTs",   json_integer(m.lastDiagTs));
    json_object_set_new(o, "lastDiagName", json_string(m.lastDiagName));

    return o;
}

json_t* json_backhaul_metrics(MetricsStore *store) {

    if (!store) return json_object();

    BackhaulAggregates a;
    (void)metrics_store_get_aggregates(store, &a);

    char nearUrl[256] = {0};
    char farUrl[256] = {0};
    long rts = 0;
    (void)metrics_store_get_reflectors(store, nearUrl, sizeof(nearUrl), farUrl, sizeof(farUrl), &rts);

    json_t *o = json_object();

    json_object_set_new(o, "microSampleCount",     json_integer(a.microCount));
    json_object_set_new(o, "multiNearSampleCount", json_integer(a.nearCount));
    json_object_set_new(o, "multiFarSampleCount",  json_integer(a.farCount));
    json_object_set_new(o, "chgSampleCount",       json_integer(a.chgCount));

    json_object_set_new(o, "microTtfbMedianMs",   json_real(a.microMed));
    json_object_set_new(o, "microTtfbP95Ms",      json_real(a.microP95));
    json_object_set_new(o, "microTtfbP99Ms",      json_real(a.microP99));
    json_object_set_new(o, "microSuccessRatePct", json_real(a.microSuccPct));
    json_object_set_new(o, "microStallRatePct",   json_real(a.microStallPct));

    json_object_set_new(o, "nearTtfbMedianMs",   json_real(a.nearMed));
    json_object_set_new(o, "nearTtfbP95Ms",      json_real(a.nearP95));
    json_object_set_new(o, "nearTtfbP99Ms",      json_real(a.nearP99));
    json_object_set_new(o, "nearSuccessRatePct", json_real(a.nearSuccPct));
    json_object_set_new(o, "nearStallRatePct",   json_real(a.nearStallPct));

    json_object_set_new(o, "farTtfbMedianMs",   json_real(a.farMed));
    json_object_set_new(o, "farTtfbP95Ms",      json_real(a.farP95));
    json_object_set_new(o, "farTtfbP99Ms",      json_real(a.farP99));
    json_object_set_new(o, "farSuccessRatePct", json_real(a.farSuccPct));
    json_object_set_new(o, "farStallRatePct",   json_real(a.farStallPct));

    json_object_set_new(o, "reflectorNearUrl", json_string(nearUrl));
    json_object_set_new(o, "reflectorFarUrl",  json_string(farUrl));
    json_object_set_new(o, "reflectorTs",      json_integer(rts));

    /* published status snapshot */
    json_object_set_new(o, "status", json_backhaul_status(store));

    return o;
}

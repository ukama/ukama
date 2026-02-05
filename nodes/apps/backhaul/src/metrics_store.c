/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>

#include "metrics_store.h"
#include "usys_log.h"
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
	if (name) strncpy(store->metrics.lastDiagName, name, sizeof(store->metrics.lastDiagName)-1);
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

	/* oldest -> newest */
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

static const char* state_str(BackhaulState s) {
	switch (s) {
	case BACKHAUL_STATE_GOOD: return "GOOD";
	case BACKHAUL_STATE_DEGRADED: return "DEGRADED";
	case BACKHAUL_STATE_DOWN: return "DOWN";
	case BACKHAUL_STATE_CAPPED: return "CAPPED";
	default: return "UNKNOWN";
	}
}

static const char* link_str(BackhaulLinkGuess g) {
	switch (g) {
	case BACKHAUL_LINK_TERRESTRIAL_LIKE: return "TERRESTRIAL_LIKE";
	case BACKHAUL_LINK_SAT_LEO_LIKE: return "SAT_LEO_LIKE";
	case BACKHAUL_LINK_SAT_GEO_LIKE: return "SAT_GEO_LIKE";
	case BACKHAUL_LINK_CELLULAR_LIKE: return "CELLULAR_LIKE";
	default: return "UNKNOWN";
	}
}

json_t* metrics_store_status_json(MetricsStore *store) {

	json_t *o = json_object();

	pthread_mutex_lock(&store->lock);

	json_object_set_new(o, "backhaulState", json_string(state_str(store->metrics.backhaulState)));
	json_object_set_new(o, "linkGuess", json_string(link_str(store->metrics.linkGuess)));
	json_object_set_new(o, "confidence", json_real(store->metrics.confidence));

	json_object_set_new(o, "dlGoodputMbps", json_real(store->metrics.dlGoodputMbps));
	json_object_set_new(o, "ulGoodputMbps", json_real(store->metrics.ulGoodputMbps));
	json_object_set_new(o, "bufferbloatInflationFactor", json_real(store->metrics.bufferbloatInflationFactor));
	json_object_set_new(o, "capDetectedMbps", json_real(store->metrics.capDetectedMbps));

	json_object_set_new(o, "nearTtfbMedianMs", json_real(store->metrics.nearTtfbMedianMs));
	json_object_set_new(o, "nearTtfbP95Ms", json_real(store->metrics.nearTtfbP95Ms));
	json_object_set_new(o, "nearTtfbP99Ms", json_real(store->metrics.nearTtfbP99Ms));

	json_object_set_new(o, "farTtfbMedianMs", json_real(store->metrics.farTtfbMedianMs));
	json_object_set_new(o, "farTtfbP95Ms", json_real(store->metrics.farTtfbP95Ms));
	json_object_set_new(o, "farTtfbP99Ms", json_real(store->metrics.farTtfbP99Ms));

	json_object_set_new(o, "probeSuccessRatePct", json_real(store->metrics.probeSuccessRatePct));
	json_object_set_new(o, "stallRatePct", json_real(store->metrics.stallRatePct));

	json_object_set_new(o, "lastMicroTs", json_integer(store->metrics.lastMicroTs));
	json_object_set_new(o, "lastMultiTs", json_integer(store->metrics.lastMultiTs));
	json_object_set_new(o, "lastChgTs", json_integer(store->metrics.lastChgTs));
	json_object_set_new(o, "lastClassifyTs", json_integer(store->metrics.lastClassifyTs));

	json_object_set_new(o, "lastDiagTs", json_integer(store->metrics.lastDiagTs));
	json_object_set_new(o, "lastDiagName", json_string(store->metrics.lastDiagName));

	pthread_mutex_unlock(&store->lock);

	return o;
}

json_t* metrics_store_snapshot_json(MetricsStore *store) {

	json_t *o = json_object();

	pthread_mutex_lock(&store->lock);

	/* compute aggregate percentiles live from ring buffers */
	double microMed=0, microP95=0, microP99=0, microSucc=0, microStall=0;
	double nearMed=0, nearP95=0, nearP99=0, nearSucc=0, nearStall=0;
	double farMed=0, farP95=0, farP99=0, farSucc=0, farStall=0;

	compute_percentiles_locked(store->microSamples, store->microCap,
							   store->microHead, store->microCount,
							   &microMed, &microP95, &microP99, &microSucc, &microStall);

	compute_percentiles_locked(store->nearSamples, store->nearCap,
							   store->nearHead, store->nearCount,
							   &nearMed, &nearP95, &nearP99, &nearSucc, &nearStall);

	compute_percentiles_locked(store->farSamples, store->farCap,
							   store->farHead, store->farCount,
							   &farMed, &farP95, &farP99, &farSucc, &farStall);

	/* expose counts too */
	json_object_set_new(o, "microSampleCount", json_integer(store->microCount));
	json_object_set_new(o, "multiNearSampleCount", json_integer(store->nearCount));
	json_object_set_new(o, "multiFarSampleCount", json_integer(store->farCount));
	json_object_set_new(o, "chgSampleCount", json_integer(store->chgCount));

	json_object_set_new(o, "microTtfbMedianMs", json_real(microMed));
	json_object_set_new(o, "microTtfbP95Ms", json_real(microP95));
	json_object_set_new(o, "microTtfbP99Ms", json_real(microP99));
	json_object_set_new(o, "microSuccessRatePct", json_real(microSucc));
	json_object_set_new(o, "microStallRatePct", json_real(microStall));

	json_object_set_new(o, "nearTtfbMedianMs", json_real(nearMed));
	json_object_set_new(o, "nearTtfbP95Ms", json_real(nearP95));
	json_object_set_new(o, "nearTtfbP99Ms", json_real(nearP99));
	json_object_set_new(o, "nearSuccessRatePct", json_real(nearSucc));
	json_object_set_new(o, "nearStallRatePct", json_real(nearStall));

	json_object_set_new(o, "farTtfbMedianMs", json_real(farMed));
	json_object_set_new(o, "farTtfbP95Ms", json_real(farP95));
	json_object_set_new(o, "farTtfbP99Ms", json_real(farP99));
	json_object_set_new(o, "farSuccessRatePct", json_real(farSucc));
	json_object_set_new(o, "farStallRatePct", json_real(farStall));

	/* published metrics snapshot (classifier writes here) */
	json_object_set_new(o, "status", metrics_store_status_json(store));

	pthread_mutex_unlock(&store->lock);

	return o;
}

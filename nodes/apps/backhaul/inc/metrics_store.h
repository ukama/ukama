/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef METRICS_STORE_H_
#define METRICS_STORE_H_

#include <pthread.h>
#include <time.h>

#include "jansson.h"
#include "usys_types.h"

typedef enum {
	BACKHAUL_STATE_UNKNOWN = 0,
	BACKHAUL_STATE_GOOD,
	BACKHAUL_STATE_DEGRADED,
	BACKHAUL_STATE_DOWN,
	BACKHAUL_STATE_CAPPED
} BackhaulState;

typedef enum {
	BACKHAUL_LINK_UNKNOWN = 0,
	BACKHAUL_LINK_TERRESTRIAL_LIKE,
	BACKHAUL_LINK_SAT_LEO_LIKE,
	BACKHAUL_LINK_SAT_GEO_LIKE,
	BACKHAUL_LINK_CELLULAR_LIKE
} BackhaulLinkGuess;

/* one probe sample */
typedef struct {
	double	ttfbMs;
	int		ok;			/* 1=success */
	int		stalled;	/* 1=stall */
	long	ts;			/* epoch seconds */
} MicroSample;

/* CHG sample (goodput) */
typedef struct {
	double	dlMbps;
	double	ulMbps;
	double	dlSec;
	double	ulSec;
	int		ok;
	long	ts;
} ChgSample;

/* Aggregate snapshot (published) */
typedef struct {

	BackhaulState		backhaulState;
	BackhaulLinkGuess	linkGuess;
	double				confidence;

	double	nearTtfbMedianMs;
	double	nearTtfbP95Ms;
	double	nearTtfbP99Ms;

	double	farTtfbMedianMs;
	double	farTtfbP95Ms;
	double	farTtfbP99Ms;

	double	probeSuccessRatePct;
	double	stallRatePct;

	double	dlGoodputMbps;
	double	ulGoodputMbps;

	double	bufferbloatInflationFactor;
	double	capDetectedMbps;

	long	lastMicroTs;
	long	lastMultiTs;
	long	lastChgTs;
	long	lastClassifyTs;

	long	lastDiagTs;
	char	lastDiagName[64];

	/* internal counters for state machine */
	int		consecFails;
	int		consecOk;

} BackhaulMetrics;

typedef struct {
	pthread_mutex_t	lock;

	MicroSample		*microSamples;
	int				microCap;
	int				microHead;
	int				microCount;

	MicroSample		*nearSamples;
	int				nearCap;
	int				nearHead;
	int				nearCount;

	MicroSample		*farSamples;
	int				farCap;
	int				farHead;
	int				farCount;

	ChgSample		*chgSamples;
	int				chgCap;
	int				chgHead;
	int				chgCount;

	BackhaulMetrics	metrics;

} MetricsStore;

int metrics_store_init(MetricsStore *store,
					   int microCap,
					   int multiCap,
					   int chgCap);

void metrics_store_free(MetricsStore *store);

void metrics_store_add_micro(MetricsStore *store, MicroSample s);
void metrics_store_add_near(MetricsStore  *store, MicroSample s);
void metrics_store_add_far(MetricsStore   *store, MicroSample s);
void metrics_store_add_chg(MetricsStore   *store, ChgSample   s);

void metrics_store_set_diag(MetricsStore *store, const char *name);

BackhaulMetrics metrics_store_get_snapshot(MetricsStore *store);
json_t* metrics_store_snapshot_json(MetricsStore *store);
json_t* metrics_store_status_json(MetricsStore *store);

#endif /* METRICS_STORE_H_ */

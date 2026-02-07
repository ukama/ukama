/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef __METRICS_STORE_H__
#define __METRICS_STORE_H__

#include <stdint.h>
#include <pthread.h>

#include "jansson.h"
#include "usys_types.h"
#include "json_types.h"

typedef enum {
	POWER_SEV_OK = 0,
	POWER_SEV_WARN,
	POWER_SEV_CRIT
} PowerSeverity;

typedef struct {
	const char		*name;

	double			v;		/* Volts */
	double			i;		/* Amps */
	double			w;		/* Watts */

	double			v_min;
	double			v_max;
	double			i_max;
	double			w_max;

	PowerSeverity	severity;
	char			reason[96];
} PowerRail;

typedef struct {
	uint64_t	last_sample_ts_ms;
	uint64_t	last_ok_ts_ms;

	double		temp_board_c;

	/* Summary */
	double		total_w;
	double		wh_total;

	/* Rails (future-ready, can remain 0 until you implement ADC/current drivers) */
	PowerRail	rail_in;
	PowerRail	rail_aux;

	PowerSeverity	overall_severity;
	char			overall_reason[128];

	int			last_err;
	char			last_err_str[128];
} PowerSnapshot;

typedef struct {
	pthread_mutex_t	lock;
	PowerSnapshot	snap;
} MetricsStore;

int metrics_store_init(MetricsStore *s);
void metrics_store_free(MetricsStore *s);

void metrics_store_set_error(MetricsStore *s, int err, const char *fmt, ...);
void metrics_store_update(MetricsStore *s, const PowerSnapshot *newSnap);

void metrics_store_get(MetricsStore *s, PowerSnapshot *out);

json_t *metrics_store_to_json(const PowerSnapshot *s);

void power_metrics_from_snapshot(const PowerSnapshot *s,
                                 const char *boardName,
                                 PowerMetrics *m);
#endif /* __METRICS_STORE_H__ */

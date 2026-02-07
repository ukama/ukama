/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <stdarg.h>

#include "metrics_store.h"
#include "json_types.h"
#include "usys_mem.h"
#include "usys_log.h"

static const char *sev_str(PowerSeverity s) {

	switch (s) {
	case POWER_SEV_OK:	return "OK";
	case POWER_SEV_WARN:	return "WARN";
	case POWER_SEV_CRIT:	return "CRIT";
	default:		return "OK";
	}
}

int metrics_store_init(MetricsStore *s) {

	if (!s) return USYS_FALSE;

	memset(s, 0, sizeof(*s));
	pthread_mutex_init(&s->lock, NULL);

	s->snap.rail_in.name = "in";
	s->snap.rail_aux.name = "aux";

	return USYS_TRUE;
}

void metrics_store_free(MetricsStore *s) {

	if (!s) return;
	pthread_mutex_destroy(&s->lock);
}

void metrics_store_set_error(MetricsStore *s, int err, const char *fmt, ...) {

	va_list ap;

	if (!s) return;

	pthread_mutex_lock(&s->lock);
	s->snap.last_err = err;

	if (fmt && *fmt) {
		va_start(ap, fmt);
		vsnprintf(s->snap.last_err_str, sizeof(s->snap.last_err_str), fmt, ap);
		va_end(ap);
	}
	pthread_mutex_unlock(&s->lock);
}

void metrics_store_update(MetricsStore *s, const PowerSnapshot *newSnap) {

	if (!s || !newSnap) return;

	pthread_mutex_lock(&s->lock);
	s->snap = *newSnap;
	pthread_mutex_unlock(&s->lock);
}

void metrics_store_get(MetricsStore *s, PowerSnapshot *out) {

	if (!s || !out) return;

	pthread_mutex_lock(&s->lock);
	*out = s->snap;
	pthread_mutex_unlock(&s->lock);
}

static json_t *rail_to_json(const PowerRail *r) {

	json_t *o = json_object();

	json_object_set_new(o, "name", json_string(r->name ? r->name : ""));
	json_object_set_new(o, "v", json_real(r->v));
	json_object_set_new(o, "i", json_real(r->i));
	json_object_set_new(o, "w", json_real(r->w));

	json_object_set_new(o, "vMin", json_real(r->v_min));
	json_object_set_new(o, "vMax", json_real(r->v_max));
	json_object_set_new(o, "iMax", json_real(r->i_max));
	json_object_set_new(o, "wMax", json_real(r->w_max));

	json_object_set_new(o, "severity", json_string(sev_str(r->severity)));
	json_object_set_new(o, "reason", json_string(r->reason));

	return o;
}

json_t *metrics_store_to_json(const PowerSnapshot *s) {

	json_t *o;
	json_t *rails;

	if (!s) return NULL;

	o = json_object();

	json_object_set_new(o, "tsMs", json_integer((json_int_t)s->last_sample_ts_ms));
	json_object_set_new(o, "lastOkTsMs", json_integer((json_int_t)s->last_ok_ts_ms));

	json_object_set_new(o, "tempBoardC", json_real(s->temp_board_c));

	json_object_set_new(o, "totalW", json_real(s->total_w));
	json_object_set_new(o, "energyWh", json_real(s->wh_total));

	json_object_set_new(o, "severity", json_string(sev_str(s->overall_severity)));
	json_object_set_new(o, "reason", json_string(s->overall_reason));

	json_object_set_new(o, "lastErr", json_integer(s->last_err));
	json_object_set_new(o, "lastErrStr", json_string(s->last_err_str));

	rails = json_array();
	json_array_append_new(rails, rail_to_json(&s->rail_in));
	json_array_append_new(rails, rail_to_json(&s->rail_aux));
	json_object_set_new(o, "rails", rails);

	return o;
}

static int is_reason_available(const char *r) {

	if (!r || !*r) return 0;
	if (strncmp(r, "not available", 13) == 0) return 0;
	return 1;
}

void power_metrics_from_snapshot(const PowerSnapshot *s,
                                 const char *boardName,
                                 PowerMetrics *m) {

	if (!s || !m) return;

	memset(m, 0, sizeof(*m));

	/* timestamp */
	m->sampleUnixMs = s->last_sample_ts_ms;

	/* board name (best effort) */
	if (boardName && *boardName) {
		snprintf(m->board, sizeof(m->board), "%s", boardName);
	} else {
		m->board[0] = '\0';
	}

	/* basic health */
	m->ok = (s->overall_severity == POWER_SEV_CRIT) ? 0 : 1;

	if (s->last_err != 0) {
		snprintf(m->err, sizeof(m->err), "%s", s->last_err_str);
	} else {
		m->err[0] = '\0';
	}

	/* LM75 mapping */
	m->haveLm75 = (s->temp_board_c != 0) ? 1 : 0;
	m->boardTempC = s->temp_board_c;

	/*
	 * LM25066 mapping:
	 * PowerSnapshot doesn't carry LM25066-specific VOUT/temp/status/diag.
	 * But rail_in is the right place to expose input rail telemetry.
	 *
	 * We treat rail_in as "VIN/IIN/PIN best effort".
	 */
	m->haveLm25066 = is_reason_available(s->rail_in.reason) ? 1 : 0;
	m->inVolts = s->rail_in.v;
	m->inAmps  = s->rail_in.i;
	m->inWatts = s->rail_in.w;

	/* Not available in snapshot today */
	m->outVolts = 0;
	m->hsTempC  = 0;
	m->statusWord = 0;
	m->diagnosticWord = 0;

	/*
	 * ADS1015 raw channels:
	 * Snapshot doesn't store raw adc inputs today, so we can't fill adcVin/adcVpa/adcAux.
	 * If you later add raw fields to PowerSnapshot (or store them elsewhere), update here.
	 */
	m->haveAds1015 = is_reason_available(s->rail_aux.reason) ? 1 : 0;
	m->adcVin = 0;
	m->adcVpa = 0;
	m->adcAux = 0;
}

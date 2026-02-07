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

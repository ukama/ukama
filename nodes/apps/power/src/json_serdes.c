/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "json_serdes.h"

static void set_num_if(json_t *o, const char *k, int cond, double v) {

	if (!cond) return;
	json_object_set_new(o, k, json_real(v));
}

json_t *json_serdes_power_metrics_to_json(const PowerMetrics *m) {

	json_t *o = json_object();

	json_object_set_new(o, "sampleUnixMs", json_integer((json_int_t)m->sampleUnixMs));
	json_object_set_new(o, "board", json_string(m->board));
	json_object_set_new(o, "ok", json_boolean(m->ok ? 1 : 0));
	json_object_set_new(o, "err", json_string(m->err));

	if (m->haveLm25066) {
		json_t *p = json_object();
		set_num_if(p, "inVolts", 1, m->inVolts);
		set_num_if(p, "outVolts", 1, m->outVolts);
		set_num_if(p, "inAmps", 1, m->inAmps);
		set_num_if(p, "inWatts", 1, m->inWatts);
		set_num_if(p, "hsTempC", 1, m->hsTempC);
		json_object_set_new(p, "statusWord", json_integer(m->statusWord));
		json_object_set_new(p, "diagnosticWord", json_integer(m->diagnosticWord));
		json_object_set_new(o, "lm25066", p);
	}

	if (m->haveLm75) {
		json_t *t = json_object();
		set_num_if(t, "boardTempC", 1, m->boardTempC);
		json_object_set_new(o, "lm75", t);
	}

	if (m->haveAds1015) {
		json_t *a = json_object();
		set_num_if(a, "adcVin", 1, m->adcVin);
		set_num_if(a, "adcVpa", 1, m->adcVpa);
		set_num_if(a, "adcAux", 1, m->adcAux);
		json_object_set_new(o, "ads1015", a);
	}

	return o;
}

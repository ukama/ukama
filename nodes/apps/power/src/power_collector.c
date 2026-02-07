/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <stdio.h>

#include "power_collector.h"
#include "usys_log.h"

static const double TEMP_WARN_C = 80.0;
static const double TEMP_CRIT_C = 90.0;

static void rail_init(PowerRail *r, const char *name) {

	memset(r, 0, sizeof(*r));
	r->name = name;
	r->severity = POWER_SEV_OK;
	snprintf(r->reason, sizeof(r->reason), "not available");
}

int power_collect_once(PowerCollectorCtx *c, uint64_t now_ms) {

	PowerSnapshot s;
	PowerSnapshot prev;
	int rc;

	if (!c || !c->store) return USYS_FALSE;

	memset(&s, 0, sizeof(s));
	metrics_store_get(c->store, &prev);

	/* preserve energy counter */
	s.wh_total = prev.wh_total;

	rail_init(&s.rail_in, "in");
	rail_init(&s.rail_aux, "aux");

	s.last_sample_ts_ms = now_ms;
	s.last_ok_ts_ms = prev.last_ok_ts_ms;

	/* Read board temperature */
	if (c->lm75_board) {
		double t = 0;

		rc = drv_lm75_read_temp_c(c->lm75_board, &t);
		if (rc == 0) {
			s.temp_board_c = t;
			s.last_ok_ts_ms = now_ms;
			s.last_err = 0;
			s.last_err_str[0] = '\0';
		} else {
			s = prev;
			s.last_sample_ts_ms = now_ms;
			s.last_err = rc;
			snprintf(s.last_err_str, sizeof(s.last_err_str), "lm75 read failed");
			metrics_store_update(c->store, &s);
			return USYS_FALSE;
		}
	} else {
		/* No sensor configured */
		s.temp_board_c = 0;
		s.last_err = -1;
		snprintf(s.last_err_str, sizeof(s.last_err_str), "lm75 not configured");
	}

	/* Derived: no power sensors yet */
	s.total_w = 0;

	/* Overall severity from temperature */
	s.overall_severity = POWER_SEV_OK;
	snprintf(s.overall_reason, sizeof(s.overall_reason), "ok");

	if (s.temp_board_c >= TEMP_CRIT_C) {
		s.overall_severity = POWER_SEV_CRIT;
		snprintf(s.overall_reason, sizeof(s.overall_reason), "board temp critical (%.1fC)", s.temp_board_c);
	} else if (s.temp_board_c >= TEMP_WARN_C) {
		s.overall_severity = POWER_SEV_WARN;
		snprintf(s.overall_reason, sizeof(s.overall_reason), "board temp high (%.1fC)", s.temp_board_c);
	}

	metrics_store_update(c->store, &s);
	return USYS_TRUE;
}

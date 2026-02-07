/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <string.h>
#include <stdio.h>

#include "power_collector.h"
#include "usys_log.h"

static void rail_init(PowerRail *r, const char *name) {

	memset(r, 0, sizeof(*r));
	r->name = name;
	r->severity = POWER_SEV_OK;
	snprintf(r->reason, sizeof(r->reason), "not available");
}

static double clamp0(double v) {
	return (v < 0) ? 0 : v;
}

static void set_err(PowerSnapshot *s, int rc, const char *msg) {

	s->last_err = rc;
	snprintf(s->last_err_str, sizeof(s->last_err_str), "%s", msg ? msg : "error");
}

static void clear_err(PowerSnapshot *s) {

	s->last_err = 0;
	s->last_err_str[0] = '\0';
}

static int ads_read_rail_v(Ads1015 *adc, int ch, double gain, double off, double *outV) {

	double v = 0;

	if (!adc || !outV) return -1;
	if (ch < 0 || ch > 3) return -1;

	if (drv_ads1015_read_single_ended(adc, ch, &v) != 0) return -1;

	/* Vrail = Vadcin * gain + off */
	*outV = (v * gain) + off;
	return 0;
}

/*
 * Current sense model (LTC6102-style) as defined in PowerCal:
 * I = Vout * (RIN/ROUT) / RSENSE
 *
 * NOTE: This requires an ADC channel wired to the current sense output.
 * Your PowerCal includes rsense/rin/rout but does NOT include ADC channel
 * mapping for those current outputs. Until you add those channels, we
 * will not fabricate currents from ADS1015.
 */
static double current_from_vout(double vout, double rin, double rout, double rsense) {

	if (rsense <= 0 || rout <= 0) return 0;
	return (vout * (rin / rout)) / rsense;
}

int power_collect_once(PowerCollectorCtx *c, uint64_t now_ms) {

	PowerSnapshot s;
	PowerSnapshot prev;
	int rc;

	if (!c || !c->store) return USYS_FALSE;

	memset(&s, 0, sizeof(s));
	metrics_store_get(c->store, &prev);

	/* carry over counters/state that must persist */
	s.wh_total = prev.wh_total;

	/* initialize rails that exist in snapshot */
	rail_init(&s.rail_in, "in");
	rail_init(&s.rail_aux, "aux");

	s.last_sample_ts_ms = now_ms;
	s.last_ok_ts_ms = prev.last_ok_ts_ms;

	/* ------------------------------------------------------------
	 * LM75: board temperature
	 * ------------------------------------------------------------ */
	if (c->lm75_board) {
		double t = 0;

		rc = drv_lm75_read_temp_c(c->lm75_board, &t);
		if (rc == 0) {
			s.temp_board_c = t;
			s.last_ok_ts_ms = now_ms;
		} else {
			/* keep previous value; report error but continue */
			s.temp_board_c = prev.temp_board_c;
			set_err(&s, rc, "lm75 read failed");
		}
	} else {
		s.temp_board_c = 0;
		set_err(&s, -1, "lm75 not configured");
	}

	/* ------------------------------------------------------------
	 * LM25066: best source for input rail V/I/P + supply temperature
	 * ------------------------------------------------------------ */
	if (c->lm25066) {
		Lm25066Sample ps;

		rc = drv_lm25066_read_sample(c->lm25066, &ps);
		if (rc == 0) {
			s.rail_in.v = clamp0(ps.vinV);
			s.rail_in.i = clamp0(ps.iinA);
			s.rail_in.w = clamp0(ps.pinW);

			/* If PIN isn't reliable, fall back to V*I */
			if (s.rail_in.w <= 0 && s.rail_in.v > 0 && s.rail_in.i > 0) {
				s.rail_in.w = s.rail_in.v * s.rail_in.i;
			}

			snprintf(s.rail_in.reason, sizeof(s.rail_in.reason), "ok");
			s.last_ok_ts_ms = now_ms;

		} else {
			/* keep previous value; report error but continue */
			s.rail_in = prev.rail_in;
			set_err(&s, rc, "lm25066 read failed");
		}
	}

	/* ------------------------------------------------------------
	 * ADS1015: rail voltages based on PowerCal mapping (if provided)
	 * ------------------------------------------------------------ */
	if (c->ads1015 && c->cal) {
		double v;

		if (c->cal->ch_12v >= 0 &&
		    ads_read_rail_v(c->ads1015, c->cal->ch_12v,
		                    c->cal->v_gain_12v, c->cal->v_off_12v, &v) == 0) {
			s.rail_aux.v = clamp0(v);
			snprintf(s.rail_aux.reason, sizeof(s.rail_aux.reason), "ok(12v)");
			s.last_ok_ts_ms = now_ms;

		} else if (c->cal->ch_5v >= 0 &&
		           ads_read_rail_v(c->ads1015, c->cal->ch_5v,
		                          c->cal->v_gain_5v, c->cal->v_off_5v, &v) == 0) {
			s.rail_aux.v = clamp0(v);
			snprintf(s.rail_aux.reason, sizeof(s.rail_aux.reason), "ok(5v)");
			s.last_ok_ts_ms = now_ms;

		} else if (c->cal->ch_3v3 >= 0 &&
		           ads_read_rail_v(c->ads1015, c->cal->ch_3v3,
		                          c->cal->v_gain_3v3, c->cal->v_off_3v3, &v) == 0) {
			s.rail_aux.v = clamp0(v);
			snprintf(s.rail_aux.reason, sizeof(s.rail_aux.reason), "ok(3v3)");
			s.last_ok_ts_ms = now_ms;

		} else if (c->cal->ch_28v >= 0 &&
		           ads_read_rail_v(c->ads1015, c->cal->ch_28v,
		                          c->cal->v_gain_28v, c->cal->v_off_28v, &v) == 0) {
			s.rail_aux.v = clamp0(v);
			snprintf(s.rail_aux.reason, sizeof(s.rail_aux.reason), "ok(28v)");
			s.last_ok_ts_ms = now_ms;
		}

		(void)current_from_vout;

	} else if (c->ads1015 && !c->cal) {
		set_err(&s, -1, "ads1015 configured but cal missing");
	}

	/* If we successfully updated OK timestamp on this pass, clear error */
	if (s.last_ok_ts_ms == now_ms) {
		clear_err(&s);
	}

	/* Evaluate overall severity (temp + rails) */
	power_eval(&s);

	metrics_store_update(c->store, &s);
	return USYS_TRUE;
}

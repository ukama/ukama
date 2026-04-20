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

static void rail_init(PowerRail *r, const char *name) {

    memset(r, 0, sizeof(*r));
    r->name = name;
    r->severity = POWER_SEV_OK;
    snprintf(r->reason, sizeof(r->reason), "not available");
}

static double clamp0(double v) {

    return (v < 0.0) ? 0.0 : v;
}

static void set_err(PowerSnapshot *s, int rc, const char *msg) {

    s->last_err = rc;
    snprintf(s->last_err_str, sizeof(s->last_err_str), "%s",
             msg ? msg : "error");
}

static void clear_err(PowerSnapshot *s) {

    s->last_err = 0;
    s->last_err_str[0] = '\0';
}

static void integrate_energy_wh(PowerSnapshot *s, const PowerSnapshot *prev) {

    double dtH = 0.0;

    if (!s || !prev) return;
    if (prev->last_sample_ts_ms == 0) return;
    if (s->last_sample_ts_ms <= prev->last_sample_ts_ms) return;

    dtH = (double)(s->last_sample_ts_ms - prev->last_sample_ts_ms) / 3600000.0;
    s->wh_total = prev->wh_total + (s->total_w * dtH);
}

static void fill_mock(PowerCollectorCtx *c, PowerSnapshot *s, uint64_t now_ms) {

    double phase;

    (void)c;

    phase = (double)(now_ms % 10000ULL) / 10000.0;

    s->have_lm75 = 1;
    s->temp_board_c = 41.0 + (phase * 6.0);

    s->have_lm25066 = 1;
    s->lm25066_vin_v = 27.4 + (phase * 0.5);
    s->lm25066_vout_v = 27.0 + (phase * 0.4);
    s->lm25066_iin_a = 0.85 + (phase * 0.2);
    s->lm25066_pin_w = s->lm25066_vin_v * s->lm25066_iin_a;
    s->lm25066_temp_c = 46.0 + (phase * 5.0);
    s->lm25066_status_word = 0;
    s->lm25066_diagnostic_word = 0;
    s->lm25066_assumed_direct = 1;

    s->have_ads1015 = 1;
    s->ads_vin = 1.12 + (phase * 0.03);
    s->ads_vpa = 0.93 + (phase * 0.02);
    s->ads_aux = 2.41 + (phase * 0.04);

    s->rail_in.v = s->lm25066_vin_v;
    s->rail_in.i = s->lm25066_iin_a;
    s->rail_in.w = s->lm25066_pin_w;
    snprintf(s->rail_in.reason, sizeof(s->rail_in.reason), "ok");

    s->rail_aux.v = s->ads_aux;
    s->rail_aux.i = 0.0;
    s->rail_aux.w = 0.0;
    snprintf(s->rail_aux.reason, sizeof(s->rail_aux.reason), "ok");

    s->last_ok_ts_ms = now_ms;
    clear_err(s);
}

int power_collect_once(PowerCollectorCtx *c, uint64_t now_ms) {

    PowerSnapshot s;
    PowerSnapshot prev;
    int rc;

    if (!c || !c->store) return USYS_FALSE;

    memset(&s, 0, sizeof(s));
    metrics_store_get(c->store, &prev);

    s.last_sample_ts_ms = now_ms;
    s.last_ok_ts_ms = prev.last_ok_ts_ms;

    rail_init(&s.rail_in, "in");
    rail_init(&s.rail_aux, "aux");

    if (c->mockMode) {
        fill_mock(c, &s, now_ms);
        s.total_w = clamp0(s.rail_in.w + s.rail_aux.w);
        integrate_energy_wh(&s, &prev);
        power_eval(&s);
        metrics_store_update(c->store, &s);
        return USYS_TRUE;
    }

    if (c->lm75_board) {
        double t = 0.0;

        rc = drv_lm75_read_temp_c(c->lm75_board, &t);
        if (rc == 0) {
            s.have_lm75 = 1;
            s.temp_board_c = t;
            s.last_ok_ts_ms = now_ms;
        } else {
            s.temp_board_c = prev.temp_board_c;
            set_err(&s, rc, "lm75 read failed");
        }
    }

    if (c->lm25066) {
        Lm25066Sample ps;

        rc = drv_lm25066_read_sample(c->lm25066, &ps);
        if (rc == 0) {
            s.have_lm25066 = 1;
            s.lm25066_vin_v = clamp0(ps.vinV);
            s.lm25066_vout_v = clamp0(ps.voutV);
            s.lm25066_iin_a = clamp0(ps.iinA);
            s.lm25066_pin_w = clamp0(ps.pinW);
            s.lm25066_temp_c = ps.tempC;
            s.lm25066_status_word = ps.statusWord;
            s.lm25066_diagnostic_word = ps.diagnosticWord;
            s.lm25066_assumed_direct = ps.assumedDirect;

            s.rail_in.v = s.lm25066_vin_v;
            s.rail_in.i = s.lm25066_iin_a;
            s.rail_in.w = s.lm25066_pin_w;

            if (s.rail_in.w <= 0.0 && s.rail_in.v > 0.0 && s.rail_in.i > 0.0) {
                s.rail_in.w = s.rail_in.v * s.rail_in.i;
                s.lm25066_pin_w = s.rail_in.w;
            }

            snprintf(s.rail_in.reason, sizeof(s.rail_in.reason), "ok");
            s.last_ok_ts_ms = now_ms;
        } else {
            s.rail_in = prev.rail_in;
            set_err(&s, rc, "lm25066 read failed");
        }
    }

    if (c->ads1015 && c->config) {
        double v = 0.0;
        int gotAny = 0;

        if (c->config->adsChVin >= 0 &&
            drv_ads1015_read_single_ended(c->ads1015,
                                          c->config->adsChVin,
                                          &v) == 0) {
            s.have_ads1015 = 1;
            s.ads_vin = clamp0(v);
            gotAny = 1;
        }

        if (c->config->adsChVpa >= 0 &&
            drv_ads1015_read_single_ended(c->ads1015,
                                          c->config->adsChVpa,
                                          &v) == 0) {
            s.have_ads1015 = 1;
            s.ads_vpa = clamp0(v);
            gotAny = 1;
        }

        if (c->config->adsChAux >= 0 &&
            drv_ads1015_read_single_ended(c->ads1015,
                                          c->config->adsChAux,
                                          &v) == 0) {
            s.have_ads1015 = 1;
            s.ads_aux = clamp0(v);
            gotAny = 1;
        }

        if (gotAny) {
            s.rail_aux.v = (s.ads_aux > 0.0) ? s.ads_aux :
                           (s.ads_vpa > 0.0) ? s.ads_vpa : s.ads_vin;
            snprintf(s.rail_aux.reason, sizeof(s.rail_aux.reason), "ok");
            s.last_ok_ts_ms = now_ms;
        }
    }

    if (s.last_ok_ts_ms == now_ms) {
        clear_err(&s);
    } else if (s.last_err == 0) {
        set_err(&s, -1, "no sensors configured or no valid sample");
    }

    s.total_w = clamp0(s.rail_in.w + s.rail_aux.w);
    integrate_energy_wh(&s, &prev);

    power_eval(&s);
    metrics_store_update(c->store, &s);

    return USYS_TRUE;
}

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */
#include <string.h>
#include "api_json.h"

int json_set_ok(UResponse *resp, json_t *body) {
    if (body == NULL) {
        return ulfius_set_string_body_response(resp, HttpStatus_InternalServerError,
                                               "{\"error\":\"null body\"}");
    }
    return ulfius_set_json_body_response(resp, HttpStatus_OK, body);
}

int json_set_err(UResponse *resp, unsigned int code, const char *msg) {
    json_t *err;
    if (msg == NULL) {
        msg = "error";
    }
    err = json_pack("{s:s}", "error", msg);
    if (err == NULL) {
        return ulfius_set_string_body_response(resp, code, "{\"error\":\"error\"}");
    }
    return ulfius_set_json_body_response(resp, code, err);
}

int json_get_bool(json_t *o, const char *key, int *out_val) {

    json_t *v;
    if (out_val == NULL) return 0;
    *out_val = 0;
    if (o == NULL || key == NULL) return 0;

    v = json_object_get(o, key);
    if (v == NULL) return 0;

    if (json_is_boolean(v)) {
        *out_val = json_is_true(v) ? 1 : 0;
        return 1;
    }
    if (json_is_integer(v)) {
        *out_val = (json_integer_value(v) != 0) ? 1 : 0;
        return 1;
    }
    if (json_is_real(v)) {
        *out_val = (json_real_value(v) != 0.0) ? 1 : 0;
        return 1;
    }
    return 0;
}

int json_get_double(json_t *o, const char *key, double *out_val) {
    json_t *v;

    if (out_val == NULL) return 0;
    *out_val = 0.0;
    if (o == NULL || key == NULL) return 0;

    v = json_object_get(o, key);
    if (v == NULL) return 0;

    if (json_is_real(v)) {
        *out_val = json_real_value(v);
        return 1;
    }
    if (json_is_integer(v)) {
        *out_val = (double)json_integer_value(v);
        return 1;
    }
    return 0;
}

int json_get_string(json_t *o, const char *key, const char **out_str) {
    json_t *v;
    if (out_str == NULL) return 0;
    *out_str = NULL;
    if (o == NULL || key == NULL) return 0;

    v = json_object_get(o, key);
    if (v == NULL) return 0;
    if (json_is_string(v)) {
        *out_str = json_string_value(v);
        return 1;
    }
    return 0;
}

json_t* json_gpio_status(const GpioStatus *st, int fem_unit_num) {
    json_t *o;

    if (st == NULL) return NULL;

    o = json_object();
    if (o == NULL) return NULL;

    json_object_set_new(o, "tx_rf_enable",   st->tx_rf_enable ? json_true() : json_false());
    json_object_set_new(o, "rx_rf_enable",   st->rx_rf_enable ? json_true() : json_false());
    json_object_set_new(o, "pa_vds_enable",  st->pa_vds_enable ? json_true() : json_false());
    json_object_set_new(o, "rf_pal_enable",  st->rf_pal_enable ? json_true() : json_false());
    /* inverted: pa_disable = !28v_vds_enable */
    json_object_set_new(o, "28v_vds_enable", st->pa_disable ? json_false() : json_true());
    json_object_set_new(o, "psu_pgood",      st->pg_reg_5v ? json_true() : json_false());
    json_object_set_new(o, "fem_unit",       json_integer(fem_unit_num));

    return o;
}

int json_extract_dac_request(json_t *body, double *carrier, double *peak) {
    int have1, have2;

    if (carrier == NULL || peak == NULL) return 0;
    if (body == NULL) return 0;

    have1 = json_get_double(body, "carrier_voltage", carrier);
    if (!have1) have1 = json_get_double(body, "carrierVoltage", carrier);

    have2 = json_get_double(body, "peak_voltage", peak);
    if (!have2) have2 = json_get_double(body, "peakVoltage", peak);

    return (have1 && have2) ? 1 : 0;
}

int json_extract_adc_thresholds(json_t *body, double *max_rp, double *max_i) {
    int have1, have2;

    if (max_rp == NULL || max_i == NULL) return 0;
    if (body == NULL) return 0;

    have1 = json_get_double(body, "max_reverse_power", max_rp);
    if (!have1) have1 = json_get_double(body, "maxReversePower", max_rp);

    have2 = json_get_double(body, "max_current", max_i);
    if (!have2) have2 = json_get_double(body, "maxCurrent", max_i);

    return (have1 && have2) ? 1 : 0;
}


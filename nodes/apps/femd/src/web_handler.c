/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#include <string.h>
#include <stdlib.h>

#include "femd.h"
#include "web_handler.h"
#include "web_service.h"

#include "api_json.h"
#include "http_status.h"
#include "jserdes.h"

#include "usys_log.h"
#include "usys_types.h"

#include "version.h"

static int respond_text(UResponse *response, int status, const char *msg) {

    ulfius_set_string_body_response(response, status, msg ? msg : "");
    return U_CALLBACK_CONTINUE;
}

static int respond_json_obj(UResponse *response, int status, json_t *json) {

    char *s;

    if (!json) {
        return respond_text(response, HttpStatus_InternalServerError, "serialize");
    }

    s = json_dumps(json, 0);
    if (!s) {
        return respond_text(response, HttpStatus_InternalServerError, "serialize");
    }

    ulfius_set_string_body_response(response, status, s);
    free(s);

    return U_CALLBACK_CONTINUE;
}

static int parse_fem_unit(const URequest *request, FemUnit *unit) {

    const char *id;

    if (!request || !unit) return STATUS_NOK;

    id = u_map_get(request->map_url, "femId");
    if (!id) {
        return STATUS_NOK;
    }

    if (!strcmp(id, "1")) {
        *unit = FEM_UNIT_1;
        return STATUS_OK;
    }

    if (!strcmp(id, "2")) {
        *unit = FEM_UNIT_2;
        return STATUS_OK;
    }

    return STATUS_NOK;
}

static int parse_adc_channel(const URequest *request, int *ch) {

    const char *s;
    char *end;
    long v;

    if (!request || !ch) return STATUS_NOK;

    s = u_map_get(request->map_url, "channel");
    if (!s) {
        return STATUS_NOK;
    }

    end = NULL;
    v = strtol(s, &end, 10);
    if (!end || *end != '\0') return STATUS_NOK;

    *ch = (int)v;
    return STATUS_OK;
}

static json_t* parse_body_json(const URequest *request) {

    json_error_t err;
    const char *body;

    if (!request) return NULL;

    body = request->binary_body;
    if (!body) return NULL;

    return json_loads(body, 0, &err);
}

static void add_err(json_t *errors, const char *scope, const char *msg) {

    if (!errors || !scope || !msg) return;

    json_object_set_new(errors, scope, json_string(msg));
}

/* Build controller metrics.
 */
static json_t *build_controller_metrics(WebCtx *ctx, json_t *errors) {

    json_t *j = json_object();
    if (!j) return NULL;

    /* Minimal, but stable schema */
    json_object_set_new(j, "ok", json_true());

    (void)ctx;
    (void)errors;
    return j;
}

static json_t *build_fem_metrics(WebCtx *ctx, FemUnit unit, int unitNum, json_t *errors) {

    FemSnapshot s;
    GpioStatus st;
    SafetyConfig cfg;

    json_t *jfem = NULL;
    json_t *jgpio = NULL;
    json_t *jsnap = NULL;
    json_t *jthr = NULL;

    int ok = 1;

    jfem = json_object();
    if (!jfem) return NULL;

    json_object_set_new(jfem, "fem_unit", json_integer(unitNum));

    if (!ctx || !ctx->gpio) {
        ok = 0;
        add_err(errors, (unitNum == 1) ?
                "fem1.gpio" : "fem2.gpio", "gpio ctx missing");
    } else {
        memset(&st, 0, sizeof(st));
        if (gpio_read_all(ctx->gpio, unit, &st) != STATUS_OK) {
            ok = 0;
            add_err(errors, (unitNum == 1) ?
                    "fem1.gpio" : "fem2.gpio", "gpio read failed");
        } else {
            jgpio = json_gpio_status(&st, unitNum);
            if (!jgpio) {
                ok = 0;
                add_err(errors, (unitNum == 1) ?
                        "fem1.gpio" : "fem2.gpio", "gpio serialize failed");
            } else {
                json_object_set_new(jfem, "gpio", jgpio); /* ownership transferred */
            }
        }
    }

    if (!ctx || !ctx->snap) {
        ok = 0;
        add_err(errors, (unitNum == 1) ?
                "fem1.snapshot" : "fem2.snapshot", "snap ctx missing");
    } else {
        memset(&s, 0, sizeof(s));
        if (snapshot_get_fem(ctx->snap, unit, &s) != STATUS_OK) {
            ok = 0;
            add_err(errors, (unitNum == 1) ?
                    "fem1.snapshot" : "fem2.snapshot", "no data");
        } else {
            /* Keep your existing rich schema as-is */
            if (json_serialize_fem_snapshot(&jsnap, unit, &s) != USYS_TRUE || !jsnap) {
                ok = 0;
                add_err(errors, (unitNum == 1) ?
                        "fem1.snapshot" : "fem2.snapshot", "snapshot serialize failed");
            } else {
                json_object_set_new(jfem, "snapshot", jsnap); /* ownership transferred */
            }

            if (s.haveTemp) {
                json_object_set_new(jfem, "temperature", json_pack("{s:f}",
                                                                   "temperature", (double)s.tempC));
            }
            if (s.haveAdc) {
                json_object_set_new(jfem, "adc", json_pack("{s:f, s:f, s:f, s:f}",
                                                          "reverse_power", (double)s.reversePowerDbm,
                                                          "forward_power", (double)s.forwardPowerDbm,
                                                          "pa_current",    (double)s.paCurrentA,
                                                          "adc_temp_volts",(double)s.adcTempVolts));
            }
            if (s.haveDac) {
                json_object_set_new(jfem, "dac", json_pack("{s:f, s:f}",
                                                          "carrier_voltage", (double)s.carrierVoltage,
                                                          "peak_voltage",    (double)s.peakVoltage));
            }
        }
    }

    if (!ctx || !ctx->safety) {
        ok = 0;
        add_err(errors, (unitNum == 1) ?
                "fem1.safety" : "fem2.safety", "safety ctx missing");
    } else {
        if (safety_get_config(ctx->safety, &cfg) != STATUS_OK) {
            ok = 0;
            add_err(errors, (unitNum == 1) ?
                    "fem1.safety" : "fem2.safety", "safety get config failed");
        } else {
            jthr = json_pack("{s:f, s:f, s:f}",
                             "max_reverse_power", (double)cfg.thresholds.max_reverse_power_dbm,
                             "max_current",       (double)cfg.thresholds.max_pa_current_a,
                             "max_temperature",   (double)cfg.thresholds.max_temperature_c);
            if (!jthr) {
                ok = 0;
                add_err(errors, (unitNum == 1) ?
                        "fem1.safety" : "fem2.safety", "thresholds serialize failed");
            } else {
                json_object_set_new(jfem, "safety_adc_thresholds", jthr); /* ownership transferred */
            }
        }
    }

    json_object_set_new(jfem, "serial", json_null());

    json_object_set_new(jfem, "ok", json_boolean(ok ? 1 : 0));
    return jfem;
}

int cb_get_metrics(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;

    json_t *root = NULL;
    json_t *errors = NULL;
    json_t *fems = NULL;

    (void)request;

    root = json_object();
    if (!root) {
        return respond_text(response, HttpStatus_InternalServerError, "oom");
    }

    errors = json_object();
    if (!errors) {
        json_decref(root);
        return respond_text(response, HttpStatus_InternalServerError, "oom");
    }

    json_object_set_new(root, "ts_unix_ms", json_null());

    /* controller block */
    json_t *controller = build_controller_metrics(ctx, errors);
    if (!controller) {
        add_err(errors, "controller", "controller metrics oom");
        controller = json_object();
        json_object_set_new(controller, "ok", json_false());
    }
    json_object_set_new(root, "controller", controller);

    /* fems array */
    fems = json_array();
    if (!fems) {
        json_decref(errors);
        json_decref(root);
        return respond_text(response, HttpStatus_InternalServerError, "oom");
    }

    json_t *fem1 = build_fem_metrics(ctx, FEM_UNIT_1, 1, errors);
    if (!fem1) {
        add_err(errors, "fem1", "fem1 metrics oom");
        fem1 = json_pack("{s:i, s:b}", "fem_unit", 1, "ok", 0);
    }

    json_t *fem2 = build_fem_metrics(ctx, FEM_UNIT_2, 2, errors);
    if (!fem2) {
        add_err(errors, "fem2", "fem2 metrics oom");
        fem2 = json_pack("{s:i, s:b}", "fem_unit", 2, "ok", 0);
    }

    json_array_append_new(fems, fem1);
    json_array_append_new(fems, fem2);
    json_object_set_new(root, "fems", fems);

    /* overall ok if no errors */
    int overall_ok = (json_object_size(errors) == 0) ? 1 : 0;
    json_object_set_new(root, "ok", json_boolean(overall_ok));

    if (overall_ok) {
        json_decref(errors);
    } else {
        json_object_set_new(root, "errors", errors);
    }

    respond_json_obj(response, HttpStatus_OK, root);
    json_decref(root);
    return U_CALLBACK_CONTINUE;
}

int cb_default(const URequest *request, UResponse *response, void *user_data) {
    (void)request;
    (void)user_data;
    return respond_text(response, HttpStatus_NotFound, HttpStatusStr(HttpStatus_NotFound));
}

int cb_options_ok(const URequest *request, UResponse *response, void *user_data) {

    const char *allow = (const char *)user_data;

    (void)request;

    if (allow) {
        u_map_put(response->map_header, "Access-Control-Allow-Methods", allow);
    }

    u_map_put(response->map_header, "Access-Control-Allow-Headers", "Content-Type");
    u_map_put(response->map_header, "Access-Control-Allow-Origin", "*");

    ulfius_set_string_body_response(response, HttpStatus_OK, "");
    return U_CALLBACK_CONTINUE;
}

int cb_not_allowed(const URequest *request, UResponse *response, void *user_data) {

    const char *allow = (const char *)user_data;

    (void)request;

    if (allow) {
        u_map_put(response->map_header, "Allow", allow);
        u_map_put(response->map_header, "Access-Control-Allow-Methods", allow);
    }

    return respond_text(response,
                        HttpStatus_MethodNotAllowed,
                        HttpStatusStr(HttpStatus_MethodNotAllowed));
}

int cb_get_version(const URequest *request, UResponse *response, void *user_data) {

    (void)request;
    (void)user_data;

    ulfius_set_string_body_response(response, HttpStatus_OK, VERSION);
    return U_CALLBACK_CONTINUE;
}

int cb_get_ping(const URequest *request, UResponse *response, void *user_data) {

    (void)request;
    (void)user_data;

    return respond_text(response,
                        HttpStatus_OK,
                        HttpStatusStr(HttpStatus_OK));
}

int cb_get_fems(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    FemSnapshot s1, s2;
    json_t *j = NULL;
    json_t *a = NULL;

    (void)request;

    if (!ctx || !ctx->snap) {
        return respond_text(response,
                            HttpStatus_InternalServerError,
                            HttpStatusStr(HttpStatus_InternalServerError));
    }

    memset(&s1, 0, sizeof(s1));
    memset(&s2, 0, sizeof(s2));

    (void)snapshot_get_fem(ctx->snap, FEM_UNIT_1, &s1);
    (void)snapshot_get_fem(ctx->snap, FEM_UNIT_2, &s2);

    j = json_object();
    if (!j) return respond_text(response, HttpStatus_InternalServerError, "oom");

    a = json_array();
    if (!a) { json_decref(j); return respond_text(response, HttpStatus_InternalServerError, "oom"); }

    json_array_append_new(a, json_pack("{s:i, s:b}",
                                       "fem_unit", 1,
                                       "present", s1.present ? 1 : 0));
    json_array_append_new(a, json_pack("{s:i, s:b}",
                                       "fem_unit", 2,
                                       "present", s2.present ? 1 : 0));

    json_object_set_new(j, "fems", a);

    respond_json_obj(response, HttpStatus_OK, j);
    json_decref(j);

    return U_CALLBACK_CONTINUE;
}

int cb_get_fem(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    FemUnit unit;
    FemSnapshot s;
    json_t *j = NULL;

    if (!ctx || !ctx->snap) {
        return respond_text(response, HttpStatus_InternalServerError, "internal");
    }

    if (parse_fem_unit(request, &unit) != STATUS_OK) {
        return respond_text(response, HttpStatus_BadRequest, "bad femId");
    }

    memset(&s, 0, sizeof(s));
    if (snapshot_get_fem(ctx->snap, unit, &s) != STATUS_OK) {
        return respond_text(response, HttpStatus_NotFound, "no data");
    }

    if (json_serialize_fem_snapshot(&j, unit, &s) != USYS_TRUE || !j) {
        return respond_text(response, HttpStatus_InternalServerError, "serialize");
    }

    respond_json_obj(response, HttpStatus_OK, j);
    json_decref(j);

    return U_CALLBACK_CONTINUE;
}

int cb_get_gpio(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    FemUnit unit;
    GpioStatus st;
    json_t *j = NULL;

    if (!ctx || !ctx->gpio) {
        return respond_text(response, HttpStatus_InternalServerError, "internal");
    }

    if (parse_fem_unit(request, &unit) != STATUS_OK) {
        return respond_text(response, HttpStatus_BadRequest, "bad femId");
    }

    memset(&st, 0, sizeof(st));
    if (gpio_read_all(ctx->gpio, unit, &st) != STATUS_OK) {
        return respond_text(response, HttpStatus_ServiceUnavailable, "gpio");
    }

    j = json_gpio_status(&st, (unit == FEM_UNIT_1) ? 1 : 2);
    if (!j) {
        return respond_text(response, HttpStatus_InternalServerError, "serialize");
    }

    respond_json_obj(response, HttpStatus_OK, j);
    json_decref(j);

    return U_CALLBACK_CONTINUE;
}

static int enqueue_gpio_apply(WebCtx *ctx, FemUnit unit, const GpioStatus *desired) {

    Job job;
    uint32_t nowMs;

    memset(&job, 0, sizeof(job));
    job.lane    = (unit == FEM_UNIT_1) ? LaneFem1 : LaneFem2;
    job.femUnit = unit;
    job.cmd     = JobCmdGpioApply;
    job.prio    = JobPrioHi;
    job.arg.gpioApply.gpio = *desired;

    nowMs = snapshot_now_ms();
    return (jobs_enqueue(ctx->jobs, &job, nowMs) != 0) ? STATUS_OK : STATUS_NOK;
}

int cb_put_gpio(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    FemUnit unit;
    json_t *body = NULL;
    GpioStatus desired;
    int v;

    if (!ctx || !ctx->jobs) {
        return respond_text(response, HttpStatus_InternalServerError, "internal");
    }

    if (parse_fem_unit(request, &unit) != STATUS_OK) {
        return respond_text(response, HttpStatus_BadRequest, "bad femId");
    }

    body = parse_body_json(request);
    if (!body) {
        return respond_text(response, HttpStatus_BadRequest, "bad json");
    }

    memset(&desired, 0, sizeof(desired));

    v = 0;
    if (!json_get_bool(body, "tx_rf_enable", &v)) {
        json_decref(body);
        return respond_text(response, HttpStatus_BadRequest, "tx_rf_enable");
    }
    desired.tx_rf_enable = v ? true : false;
    
    v = 0;
    if (!json_get_bool(body, "rx_rf_enable", &v)) {
        json_decref(body);
        return respond_text(response, HttpStatus_BadRequest, "rx_rf_enable");
    }
    desired.rx_rf_enable = v ? true : false;

    v = 0;
    if (!json_get_bool(body, "pa_vds_enable", &v)) {
        json_decref(body);
        return respond_text(response, HttpStatus_BadRequest, "pa_vds_enable");
    }
    desired.pa_vds_enable = v ? true : false;

    v = 0;
    if (!json_get_bool(body, "rf_pal_enable", &v)) {
        json_decref(body);
        return respond_text(response, HttpStatus_BadRequest, "rf_pal_enable");
    }
    desired.rf_pal_enable = v ? true : false;
    
    v = 0;
    if (!json_get_bool(body, "pa_disable", &v)) {
        json_decref(body);
        return respond_text(response, HttpStatus_BadRequest, "pa_disable");
    }
    desired.pa_disable = v ? true : false;

    json_decref(body);

    if (enqueue_gpio_apply(ctx, unit, &desired) != STATUS_OK) {
        return respond_text(response, HttpStatus_ServiceUnavailable, "queue");
    }

    return respond_text(response, HttpStatus_Accepted, "accepted");
}

int cb_patch_gpio(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    FemUnit unit;
    json_t *body = NULL;
    GpioStatus cur;
    int v;

    if (!ctx || !ctx->jobs || !ctx->gpio) {
        return respond_text(response, HttpStatus_InternalServerError, "internal");
    }

    if (parse_fem_unit(request, &unit) != STATUS_OK) {
        return respond_text(response, HttpStatus_BadRequest, "bad femId");
    }

    memset(&cur, 0, sizeof(cur));
    if (gpio_read_all(ctx->gpio, unit, &cur) != STATUS_OK) {
        return respond_text(response, HttpStatus_ServiceUnavailable, "gpio");
    }

    body = parse_body_json(request);
    if (!body) {
        return respond_text(response, HttpStatus_BadRequest, "bad json");
    }

    if (json_get_bool(body, "tx_rf_enable", &v))  cur.tx_rf_enable = v ?  true : false;
    if (json_get_bool(body, "rx_rf_enable", &v))  cur.rx_rf_enable = v ?  true : false;
    if (json_get_bool(body, "pa_vds_enable", &v)) cur.pa_vds_enable = v ? true : false;
    if (json_get_bool(body, "rf_pal_enable", &v)) cur.rf_pal_enable = v ? true : false;
    if (json_get_bool(body, "pa_disable", &v))    cur.pa_disable = v ?    true : false;

    json_decref(body);

    if (enqueue_gpio_apply(ctx, unit, &cur) != STATUS_OK) {
        return respond_text(response, HttpStatus_ServiceUnavailable, "queue");
    }

    return respond_text(response, HttpStatus_Accepted, "accepted");
}

int cb_get_dac(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    FemUnit unit;
    FemSnapshot s;
    json_t *j = NULL;

    if (!ctx || !ctx->snap) {
        return respond_text(response, HttpStatus_InternalServerError, "internal");
    }

    if (parse_fem_unit(request, &unit) != STATUS_OK) {
        return respond_text(response, HttpStatus_BadRequest, "bad femId");
    }

    memset(&s, 0, sizeof(s));
    if (snapshot_get_fem(ctx->snap, unit, &s) != STATUS_OK) {
        return respond_text(response, HttpStatus_NotFound, "no data");
    }

    if (!s.haveDac) {
        return respond_text(response, HttpStatus_NotFound, "no dac");
    }

    j = json_pack("{s:f, s:f}", "carrier_voltage", (double)s.carrierVoltage,
                  "peak_voltage", (double)s.peakVoltage);
    if (!j) {
        return respond_text(response, HttpStatus_InternalServerError, "oom");
    }

    respond_json_obj(response, HttpStatus_OK, j);
    json_decref(j);

    return U_CALLBACK_CONTINUE;
}

int cb_put_dac(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    FemUnit unit;
    json_t *body = NULL;
    double carrier = 0.0;
    double peak = 0.0;
    Job job;
    uint32_t nowMs;

    if (!ctx || !ctx->jobs) {
        return respond_text(response, HttpStatus_InternalServerError, "internal");
    }

    if (parse_fem_unit(request, &unit) != STATUS_OK) {
        return respond_text(response, HttpStatus_BadRequest, "bad femId");
    }

    body = parse_body_json(request);
    if (!body) {
        return respond_text(response, HttpStatus_BadRequest, "bad json");
    }

    if (!json_extract_dac_request(body, &carrier, &peak)) {
        json_decref(body);
        return respond_text(response, HttpStatus_BadRequest, "carrier/peak");
    }
    json_decref(body);

    memset(&job, 0, sizeof(job));
    job.lane    = (unit == FEM_UNIT_1) ? LaneFem1 : LaneFem2;
    job.femUnit = unit;
    job.prio    = JobPrioHi;

    nowMs = snapshot_now_ms();

    job.cmd                 = JobCmdDacSetCarrier;
    job.arg.voltage.voltage = (float)carrier;
    if (jobs_enqueue(ctx->jobs, &job, nowMs) == 0) {
        return respond_text(response, HttpStatus_ServiceUnavailable, "queue");
    }

    job.cmd                 = JobCmdDacSetPeak;
    job.arg.voltage.voltage = (float)peak;
    if (jobs_enqueue(ctx->jobs, &job, nowMs) == 0) {
        return respond_text(response, HttpStatus_ServiceUnavailable, "queue");
    }

    return respond_text(response, HttpStatus_Accepted, "accepted");
}

int cb_get_temp(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    FemUnit unit;
    FemSnapshot s;
    json_t *j = NULL;

    if (!ctx || !ctx->snap) {
        return respond_text(response, HttpStatus_InternalServerError, "internal");
    }

    if (parse_fem_unit(request, &unit) != STATUS_OK) {
        return respond_text(response, HttpStatus_BadRequest, "bad femId");
    }

    memset(&s, 0, sizeof(s));
    if (snapshot_get_fem(ctx->snap, unit, &s) != STATUS_OK) {
        return respond_text(response, HttpStatus_NotFound, "no data");
    }

    if (!s.haveTemp) {
        return respond_text(response, HttpStatus_NotFound, "no temp");
    }

    j = json_pack("{s:f}", "temperature", (double)s.tempC);
    if (!j) {
        return respond_text(response, HttpStatus_InternalServerError, "oom");
    }

    respond_json_obj(response, HttpStatus_OK, j);
    json_decref(j);

    return U_CALLBACK_CONTINUE;
}

int cb_get_adc_all(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    FemUnit unit;
    FemSnapshot s;
    json_t *j = NULL;

    if (!ctx || !ctx->snap) {
        return respond_text(response, HttpStatus_InternalServerError, "internal");
    }

    if (parse_fem_unit(request, &unit) != STATUS_OK) {
        return respond_text(response, HttpStatus_BadRequest, "bad femId");
    }

    memset(&s, 0, sizeof(s));
    if (snapshot_get_fem(ctx->snap, unit, &s) != STATUS_OK) {
        return respond_text(response, HttpStatus_NotFound, "no data");
    }

    if (!s.haveAdc) {
        return respond_text(response, HttpStatus_NotFound, "no adc");
    }

    j = json_pack("{s:f, s:f, s:f, s:f}",
                  "reverse_power", (double)s.reversePowerDbm,
                  "forward_power", (double)s.forwardPowerDbm,
                  "pa_current", (double)s.paCurrentA,
                  "adc_temp_volts", (double)s.adcTempVolts);
    if (!j) {
        return respond_text(response, HttpStatus_InternalServerError, "oom");
    }

    respond_json_obj(response, HttpStatus_OK, j);
    json_decref(j);

    return U_CALLBACK_CONTINUE;
}

int cb_get_adc_chan(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    FemUnit unit;
    int ch = 0;
    FemSnapshot s;
    json_t *j = NULL;

    if (!ctx || !ctx->snap) {
        return respond_text(response, HttpStatus_InternalServerError, "internal");
    }

    if (parse_fem_unit(request, &unit) != STATUS_OK) {
        return respond_text(response, HttpStatus_BadRequest, "bad femId");
    }

    if (parse_adc_channel(request, &ch) != STATUS_OK) {
        return respond_text(response, HttpStatus_BadRequest, "bad channel");
    }

    memset(&s, 0, sizeof(s));
    if (snapshot_get_fem(ctx->snap, unit, &s) != STATUS_OK) {
        return respond_text(response, HttpStatus_NotFound, "no data");
    }

    if (!s.haveAdc) {
        return respond_text(response, HttpStatus_NotFound, "no adc");
    }

    if (ch == 0) {
        j = json_pack("{s:i, s:f}", "channel", ch, "value", (double)s.reversePowerDbm);
    } else if (ch == 1) {
        j = json_pack("{s:i, s:f}", "channel", ch, "value", (double)s.forwardPowerDbm);
    } else if (ch == 2) {
        j = json_pack("{s:i, s:f}", "channel", ch, "value", (double)s.paCurrentA);
    } else if (ch == 3) {
        j = json_pack("{s:i, s:f}", "channel", ch, "value", (double)s.adcTempVolts);
    } else {
        return respond_text(response, HttpStatus_NotFound, "unknown channel");
    }

    if (!j) {
        return respond_text(response, HttpStatus_InternalServerError, "oom");
    }

    respond_json_obj(response, HttpStatus_OK, j);
    json_decref(j);

    return U_CALLBACK_CONTINUE;
}


int cb_get_adc_thr(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    FemUnit unit;
    SafetyConfig cfg;
    json_t *j = NULL;

    if (!ctx || !ctx->safety) {
        return respond_text(response, HttpStatus_InternalServerError, "internal");
    }

    if (parse_fem_unit(request, &unit) != STATUS_OK) {
        return respond_text(response, HttpStatus_BadRequest, "bad femId");
    }

    (void)unit;

    if (safety_get_config(ctx->safety, &cfg) != STATUS_OK) {
        return respond_text(response, HttpStatus_ServiceUnavailable, "safety");
    }

    j = json_pack("{s:f, s:f, s:f}",
                  "max_reverse_power", (double)cfg.thresholds.max_reverse_power_dbm,
                  "max_current",       (double)cfg.thresholds.max_pa_current_a,
                  "max_temperature",   (double)cfg.thresholds.max_temperature_c);
    if (!j) {
        return respond_text(response, HttpStatus_InternalServerError, "oom");
    }

    respond_json_obj(response, HttpStatus_OK, j);
    json_decref(j);

    return U_CALLBACK_CONTINUE;
}

int cb_put_adc_thr(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    FemUnit unit;
    json_t *body = NULL;
    double max_rp = 0.0;
    double max_i  = 0.0;
    SafetyConfig cfg;

    if (!ctx || !ctx->safety) {
        return respond_text(response, HttpStatus_InternalServerError, "internal");
    }

    if (parse_fem_unit(request, &unit) != STATUS_OK) {
        return respond_text(response, HttpStatus_BadRequest, "bad femId");
    }

    body = parse_body_json(request);
    if (!body) {
        return respond_text(response, HttpStatus_BadRequest, "bad json");
    }

    if (!json_extract_adc_thresholds(body, &max_rp, &max_i)) {
        json_decref(body);
        return respond_text(response, HttpStatus_BadRequest, "thresholds");
    }
    json_decref(body);

    if (safety_get_config(ctx->safety, &cfg) != STATUS_OK) {
        return respond_text(response, HttpStatus_ServiceUnavailable, "safety");
    }

    cfg.thresholds.max_reverse_power_dbm = (float)max_rp;
    cfg.thresholds.max_pa_current_a      = (float)max_i;

    (void)unit;

    if (safety_set_config(ctx->safety, &cfg) != STATUS_OK) {
        return respond_text(response, HttpStatus_ServiceUnavailable, "safety");
    }

    return respond_text(response, HttpStatus_Accepted, "accepted");
}

int cb_post_safety_restore(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    FemUnit unit;

    (void)request;

    if (!ctx || !ctx->safety) {
        return respond_text(response, HttpStatus_InternalServerError, "internal");
    }

    if (parse_fem_unit(request, &unit) != STATUS_OK) {
        return respond_text(response, HttpStatus_BadRequest, "bad femId");
    }

    if (safety_force_restore(ctx->safety, unit) != STATUS_OK) {
        return respond_text(response, HttpStatus_ServiceUnavailable, "safety");
    }

    return respond_text(response, HttpStatus_Accepted, "accepted");
}

int cb_get_serial(const URequest *request, UResponse *response, void *user_data) {

    (void)request;
    (void)user_data;

    return respond_text(response, HttpStatus_NotImplemented, "not implemented");
}

int cb_put_serial(const URequest *request, UResponse *response, void *user_data) {

    (void)request;
    (void)user_data;

    return respond_text(response, HttpStatus_NotImplemented, "not implemented");
}

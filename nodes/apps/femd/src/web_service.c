/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */
#include <ulfius.h>
#include <jansson.h>
#include <string.h>
#include <stdlib.h>

#include "usys_log.h"

#include "femd.h"
#include "http_status.h"
#include "version.h"
#include "web_service.h"
#include "config.h"
#include "api_http.h"
#include "api_json.h"
#include "gpio_controller.h"
#include "i2c_controller.h"
#include "safety_monitor.h"

int cb_not_allowed(const URequest *req,
                   UResponse *resp,
                   void *user_data) {
    const char *allow_str;
    allow_str = (const char *)user_data;
    api_set_cors_allow(resp, allow_str);
    return ulfius_set_string_body_response(resp,
                                           HttpStatus_MethodNotAllowed,
                                           HttpStatusStr(HttpStatus_MethodNotAllowed));
}

int cb_options_ok(const URequest *req,
                  UResponse *resp,
                  void *user_data) {
    const char *allow_str;

    allow_str = (const char *)user_data;
    api_set_cors_allow(resp, allow_str);
    return ulfius_set_string_body_response(resp, HttpStatus_OK, "");
}

int cb_default(const URequest *req,
               UResponse *resp,
               void *user_data) {
    api_set_cors_allow(resp, "GET, POST, PUT, PATCH, DELETE, OPTIONS");
    return json_set_err(resp, HttpStatus_NotFound, "Endpoint not found");
}

int cb_get_health(const URequest *req,
                  UResponse *resp,
                  void *user_data) {
    json_t *o;
    api_set_cors_allow(resp, "GET, OPTIONS");

    o = json_pack("{s:s, s:s, s:i}",
                  "service_name", "femd",
                  "version", VERSION,
                  "uptime", 0); /* TODO: real uptime */
    if (o == NULL) return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    return json_set_ok(resp, o);
}

int cb_get_ping(const URequest *req,
                UResponse *resp,
                void *user_data) {
    api_set_cors_allow(resp, "GET, OPTIONS");
    return ulfius_set_string_body_response(resp, HttpStatus_OK, "pong");
}

int cb_get_version(const URequest *req,
                   UResponse *resp,
                   void *user_data) {
    json_t *o;
    api_set_cors_allow(resp, "GET, OPTIONS");

    o = json_pack("{s:s}", "version", VERSION);
    if (o == NULL) return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    return json_set_ok(resp, o);
}

int cb_get_fems(const URequest *req,
                UResponse *resp,
                void *user_data) {
    json_t *arr, *o;
    api_set_cors_allow(resp, "GET, OPTIONS");

    arr = json_array();
    if (arr == NULL) return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    json_array_append_new(arr, json_integer(1));
    json_array_append_new(arr, json_integer(2));

    o = json_object();
    if (o == NULL) {
        json_decref(arr);
        return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    }
    json_object_set_new(o, "fems", arr);
    return json_set_ok(resp, o);
}

int cb_get_fem(const URequest *req,
               UResponse *resp,
               void *user_data) {
    FemUnit unit;
    int bus;
    json_t *o;

    api_set_cors_allow(resp, "GET, OPTIONS");

    if (!api_parse_fem_id(req, &unit))
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid femId");

    bus = i2c_get_bus_for_fem(unit);
    o = json_pack("{s:i, s:i}", "id", (int)unit, "bus", bus);
    if (o == NULL) return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    return json_set_ok(resp, o);
}

int cb_get_gpio(const URequest *req,
                UResponse *resp,
                void *user_data) {
    ServerConfig *cfg;
    FemUnit unit;
    GpioStatus st;
    json_t *o;
    int rc;

    cfg = (ServerConfig*)user_data;
    api_set_cors_allow(resp, "GET, PUT, PATCH, OPTIONS");

    if (!api_parse_fem_id(req, &unit))
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid femId");

    rc = gpio_read_all(cfg->gpioController, unit, &st);
    if (rc != STATUS_OK)
        return json_set_err(resp, HttpStatus_InternalServerError, "Failed to read GPIO status");

    o = json_gpio_status(&st, (int)unit);
    if (o == NULL) return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    return json_set_ok(resp, o);
}

/* internal helper to handle both PUT and PATCH */
static int handle_gpio_write(const URequest *req,
                             UResponse *resp,
                             void *user_data) {
    ServerConfig *cfg;
    FemUnit unit;
    GpioStatus cur, desired;
    json_t *body, *o;
    int rc, any, v;

    cfg = (ServerConfig*)user_data;
    api_set_cors_allow(resp, "GET, PUT, PATCH, OPTIONS");

    if (!api_parse_fem_id(req, &unit))
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid femId");

    body = ulfius_get_json_body_request(req, NULL);
    if (body == NULL)
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid or missing JSON body");

    rc = gpio_read_all(cfg->gpioController, unit, &cur);
    if (rc != STATUS_OK) {
        json_decref(body);
        return json_set_err(resp, HttpStatus_InternalServerError, "Failed to read GPIO status");
    }

    desired = cur;
    any = 0;

    if (json_get_bool(body, "tx_rf_enable", &v)) {
        desired.tx_rf_enable = v ? 1 : 0; any = 1;
    }
    if (json_get_bool(body, "rx_rf_enable", &v)) {
        desired.rx_rf_enable = v ? 1 : 0; any = 1;
    }
    if (json_get_bool(body, "pa_vds_enable", &v)) {
        desired.pa_vds_enable = v ? 1 : 0; any = 1;
    }
    if (json_get_bool(body, "rf_pal_enable", &v)) {
        desired.rf_pal_enable = v ? 1 : 0; any = 1;
    }
    if (json_get_bool(body, "28v_vds_enable", &v)) {
        /* inverted mapping to pa_disable */
        desired.pa_disable = v ? 0 : 1; any = 1;
    }

    json_decref(body);

    if (!any)
        return json_set_err(resp, HttpStatus_BadRequest, "No valid GPIO fields to update");

    rc = gpio_apply(cfg->gpioController, unit, &desired);
    if (rc != STATUS_OK)
        return json_set_err(resp, HttpStatus_InternalServerError, "Failed to apply GPIO update");

    o = json_gpio_status(&desired, (int)unit);
    if (o == NULL) return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    return json_set_ok(resp, o);
}

int cb_put_gpio(const URequest *req,
                UResponse *resp,
                void *user_data) {
    return handle_gpio_write(req, resp, user_data);
}

int cb_patch_gpio(const URequest *req,
                  UResponse *resp,
                  void *user_data) {
    return handle_gpio_write(req, resp, user_data);
}

int cb_get_dac(const URequest *req,
               UResponse *resp,
               void *user_data) {
    ServerConfig *cfg;
    FemUnit unit;
    float carrier, peak;
    int rc;
    json_t *o;

    cfg = (ServerConfig*)user_data;
    api_set_cors_allow(resp, "GET, PUT, OPTIONS");

    if (!api_parse_fem_id(req, &unit))
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid femId");

    rc = dac_get_config(cfg->i2cController, &carrier, &peak);
    if (rc != STATUS_OK)
        return json_set_err(resp, HttpStatus_InternalServerError, "Failed to read DAC configuration");

    o = json_pack("{s:f, s:f, s:i}",
                  "carrier_voltage", (double)carrier,
                  "peak_voltage",    (double)peak,
                  "fem_unit",        (int)unit);
    if (o == NULL) return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    return json_set_ok(resp, o);
}

int cb_put_dac(const URequest *req,
               UResponse *resp,
               void *user_data) {
    ServerConfig *cfg;
    FemUnit unit;
    json_t *body;
    double carrier, peak;
    int ok, rc;
    json_t *o;

    cfg = (ServerConfig*)user_data;
    api_set_cors_allow(resp, "GET, PUT, OPTIONS");

    if (!api_parse_fem_id(req, &unit))
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid femId");

    body = ulfius_get_json_body_request(req, NULL);
    if (body == NULL)
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid or missing JSON body");

    ok = json_extract_dac_request(body, &carrier, &peak);
    json_decref(body);
    if (!ok)
        return json_set_err(resp, HttpStatus_BadRequest, "carrier/peak voltages required");

    if (carrier < 0.0 || carrier > DAC_MAX_CARRIER_VOLTAGE)
        return json_set_err(resp, HttpStatus_BadRequest, "carrier_voltage out of range");
    if (peak < 0.0 || peak > DAC_MAX_PEAK_VOLTAGE)
        return json_set_err(resp, HttpStatus_BadRequest, "peak_voltage out of range");

    rc = dac_init(cfg->i2cController, unit);
    if (rc != STATUS_OK)
        return json_set_err(resp, HttpStatus_InternalServerError, "Failed to initialize DAC");

    rc = dac_set_carrier_voltage(cfg->i2cController, unit, (float)carrier);
    if (rc != STATUS_OK)
        return json_set_err(resp, HttpStatus_InternalServerError, "Failed to set carrier voltage");

    rc = dac_set_peak_voltage(cfg->i2cController, unit, (float)peak);
    if (rc != STATUS_OK)
        return json_set_err(resp, HttpStatus_InternalServerError, "Failed to set peak voltage");

    o = json_pack("{s:s, s:f, s:f, s:i}",
                  "status", "success",
                  "carrier_voltage", carrier,
                  "peak_voltage",    peak,
                  "fem_unit",        (int)unit);
    if (o == NULL) return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    return json_set_ok(resp, o);
}

int cb_get_temp(const URequest *req,
                UResponse *resp,
                void *user_data) {
    ServerConfig *cfg;
    FemUnit unit;
    float tempc;
    int rc;
    json_t *o;

    cfg = (ServerConfig*)user_data;
    api_set_cors_allow(resp, "GET, OPTIONS");

    if (!api_parse_fem_id(req, &unit))
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid femId");

    rc = temp_sensor_init(cfg->i2cController, unit);
    if (rc != STATUS_OK)
        return json_set_err(resp, HttpStatus_InternalServerError,
                            "Failed to initialize temperature sensor");

    rc = temp_sensor_read(cfg->i2cController, unit, &tempc);
    if (rc != STATUS_OK)
        return json_set_err(resp, HttpStatus_InternalServerError,
                            "Failed to read temperature");

    o = json_pack("{s:f, s:s, s:i}", "temperature", (double)tempc, "unit", "celsius",
                  "fem_unit", (int)unit);
    if (o == NULL) return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    return json_set_ok(resp, o);
}

int cb_get_adc_all(const URequest *req,
                   UResponse *resp,
                   void *user_data) {
    ServerConfig *cfg;
    FemUnit unit;
    float rev_dbm, pa_cur_a;
    int rc1, rc2;
    json_t *units, *o;

    cfg = (ServerConfig*)user_data;
    api_set_cors_allow(resp, "GET, OPTIONS");

    if (!api_parse_fem_id(req, &unit))
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid femId");

    rc1 = adc_init(cfg->i2cController, unit);
    if (rc1 != STATUS_OK)
        return json_set_err(resp, HttpStatus_InternalServerError, "Failed to initialize ADC");

    rc1 = adc_read_reverse_power(cfg->i2cController, unit, &rev_dbm);
    rc2 = adc_read_pa_current(cfg->i2cController, unit, &pa_cur_a);
    if (rc1 != STATUS_OK || rc2 != STATUS_OK)
        return json_set_err(resp, HttpStatus_InternalServerError, "Failed to read ADC values");

    units = json_pack("{s:s, s:s}", "reverse_power", "dBm", "pa_current", "A");
    if (units == NULL) {
        return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    }

    o = json_pack("{s:f, s:f, s:i, s:o}",
                  "reverse_power", (double)rev_dbm,
                  "pa_current",    (double)pa_cur_a,
                  "fem_unit",      (int)unit,
                  "units",         units);
    if (o == NULL) {
        json_decref(units);
        return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    }

    return json_set_ok(resp, o);
}

int cb_get_adc_chan(const URequest *req,
                    UResponse *resp,
                    void *user_data) {
    ServerConfig *cfg;
    FemUnit unit;
    int ch_num;
    float volts;
    int rc;
    json_t *o;

    cfg = (ServerConfig*)user_data;
    api_set_cors_allow(resp, "GET, OPTIONS");

    if (!api_parse_fem_id(req, &unit))
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid femId");
    if (!api_parse_channel_id(req, &ch_num))
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid channel");

    rc = adc_init(cfg->i2cController, unit);
    if (rc != STATUS_OK)
        return json_set_err(resp, HttpStatus_InternalServerError, "Failed to initialize ADC");

    rc = adc_read_channel(cfg->i2cController, unit, (ADCChannel)ch_num, &volts);
    if (rc != STATUS_OK)
        return json_set_err(resp, HttpStatus_InternalServerError, "Failed to read ADC channel");

    o = json_pack("{s:i, s:f, s:i}", "channel", ch_num, "voltage", (double)volts, "fem_unit", (int)unit);
    if (o == NULL) return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    return json_set_ok(resp, o);
}

int cb_get_adc_thr(const URequest *req,
                   UResponse *resp,
                   void *user_data) {
    ServerConfig *cfg;
    FemUnit unit;
    json_t *o;

    cfg = (ServerConfig*)user_data;
    api_set_cors_allow(resp, "GET, PUT, OPTIONS");

    if (!api_parse_fem_id(req, &unit))
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid femId");

    o = json_pack("{s:f, s:f, s:b}",
                  "max_reverse_power", (double)cfg->i2cController->adcState.maxReversePower,
                  "max_current",       (double)cfg->i2cController->adcState.maxCurrent,
                  "safety_enabled",    cfg->i2cController->adcState.safetyEnabled ? 1 : 0);
    if (o == NULL) return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    return json_set_ok(resp, o);
}

int cb_put_adc_thr(const URequest *req,
                   UResponse *resp,
                   void *user_data) {
    ServerConfig *cfg;
    FemUnit unit;
    json_t *body;
    double max_rp, max_i;
    int ok, rc;
    json_t *o;

    cfg = (ServerConfig*)user_data;
    api_set_cors_allow(resp, "GET, PUT, OPTIONS");

    if (!api_parse_fem_id(req, &unit)) {
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid femId");
    }

    body = ulfius_get_json_body_request(req, NULL);
    if (body == NULL) {
        return json_set_err(resp, HttpStatus_BadRequest,
                            "Invalid or missing JSON body");
    }

    ok = json_extract_adc_thresholds(body, &max_rp, &max_i);
    json_decref(body);
    if (!ok) {
        return json_set_err(resp, HttpStatus_BadRequest,
                            "max_reverse_power and max_current required");
    }

    rc = adc_set_safety_thresholds(cfg->i2cController, (float)max_rp, (float)max_i);
    if (rc != STATUS_OK) {
        return json_set_err(resp, HttpStatus_InternalServerError,
                            "Failed to set safety thresholds");
    }

    o = json_pack("{s:s, s:f, s:f}",
                  "status", "success",
                  "max_reverse_power", max_rp,
                  "max_current",       max_i);
    if (o == NULL) {
        return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    }

    return json_set_ok(resp, o);
}

int cb_get_serial(const URequest *req,
                  UResponse *resp,
                  void *user_data) {
    ServerConfig *cfg;
    FemUnit unit;
    char serial[32];
    int rc;
    json_t *o;

    cfg = (ServerConfig*)user_data;
    api_set_cors_allow(resp, "GET, PUT, OPTIONS");

    if (!api_parse_fem_id(req, &unit)) {
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid femId");
    }

    rc = eeprom_read_serial(cfg->i2cController, unit, serial, sizeof(serial));
    if (rc != STATUS_OK) {
        return json_set_err(resp, HttpStatus_NotFound, "No serial number found");
    }

    o = json_pack("{s:s, s:i}", "serial", serial, "fem_unit", (int)unit);
    if (o == NULL) {
        return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    }

    return json_set_ok(resp, o);
}

int cb_put_serial(const URequest *req,
                  UResponse *resp,
                  void *user_data) {
    ServerConfig *cfg;
    FemUnit unit;
    json_t *body;
    const char *serial;
    size_t len;
    int rc;
    json_t *o;

    cfg = (ServerConfig*)user_data;
    api_set_cors_allow(resp, "GET, PUT, OPTIONS");

    if (!api_parse_fem_id(req, &unit)) {
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid femId");
    }

    body = ulfius_get_json_body_request(req, NULL);
    if (body == NULL) {
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid or missing JSON body");
    }

    if (!json_get_string(body, "serial", &serial)) {
        json_decref(body);
        return json_set_err(resp, HttpStatus_BadRequest, "serial required");
    }
    len = strlen(serial);
    json_decref(body);

    if (len == 0 || len > 16) {
        return json_set_err(resp, HttpStatus_BadRequest, "serial must be 1..16 chars");
    }

    rc = eeprom_write_serial(cfg->i2cController, unit, serial);
    if (rc != STATUS_OK) {
        return json_set_err(resp, HttpStatus_InternalServerError, "Failed to write serial number");
    }

    o = json_pack("{s:s, s:s, s:i}", "status", "success", "serial", serial, "fem_unit", (int)unit);
    if (o == NULL) {
        return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    }

    return json_set_ok(resp, o);
}

/* POST /v1/fems/:femId/safety/restore  -> force manual restore */
int cb_post_safety_restore(const URequest *req,
                           UResponse *resp,
                           void *user_data) {
    ServerConfig *cfg;
    FemUnit unit;
    int rc;
    json_t *o;

    cfg = (ServerConfig*)user_data;
    api_set_cors_allow(resp, "POST, OPTIONS");

    if (!api_parse_fem_id(req, &unit)) {
        return json_set_err(resp, HttpStatus_BadRequest, "Invalid femId");
    }

    if (cfg == NULL || cfg->safetyMonitor == NULL) {
        return json_set_err(resp, HttpStatus_InternalServerError, "Safety monitor unavailable");
    }

    rc = safety_monitor_restore_pa(cfg->safetyMonitor, unit);
    if (rc != STATUS_OK) {
        return json_set_err(resp, HttpStatus_InternalServerError, "Failed to restore PA");
    }

    o = json_pack("{s:s, s:i}", "status", "restored", "fem_unit", (int)unit);
    if (o == NULL) {
        return json_set_err(resp, HttpStatus_InternalServerError, "Alloc error");
    }
    return json_set_ok(resp, o);
}

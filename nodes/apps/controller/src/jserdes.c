/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <math.h>
#include <string.h>

#include "jserdes.h"
#include "driver.h"
#include "usys_log.h"

json_t *json_serialize_status(const MetricsSnapshot *snap) {
    const ControllerData *d;
    json_t *obj, *solar, *battery, *ctrl;

    if (!snap) return NULL;

    d    = &snap->data;
    obj  = json_object();
    solar = json_object();
    battery = json_object();
    ctrl  = json_object();

    json_object_set_new(obj, "timestamp_ms",      json_integer((json_int_t)d->timestamp_ms));
    json_object_set_new(obj, "comm_ok",           json_boolean(d->comm_ok));
    json_object_set_new(obj, "comm_errors",       json_integer(d->comm_errors));
    json_object_set_new(obj, "charge_state",      json_string(charge_state_str(d->charge_state)));
    json_object_set_new(obj, "error_code",        json_integer(d->error_code));
    json_object_set_new(obj, "error_str",         json_string(error_code_str(d->error_code)));
    json_object_set_new(obj, "overall_severity",  json_string(severity_str(snap->overall_severity)));
    json_object_set_new(obj, "active_alarm_count",json_integer(snap->active_alarm_count));

    if (d->firmware[0])   json_object_set_new(obj, "firmware",   json_string(d->firmware));
    if (d->serial[0])     json_object_set_new(obj, "serial",     json_string(d->serial));
    if (d->product_id[0]) json_object_set_new(obj, "product_id", json_string(d->product_id));

    json_object_set_new(solar, "voltage_v",     json_real(d->pv_voltage_v));
    json_object_set_new(solar, "current_a",     json_real(d->pv_current_a));
    json_object_set_new(solar, "power_w",       json_real(d->pv_power_w));
    json_object_set_new(solar, "yield_today_kwh",  json_real(d->yield_today_kwh));
    json_object_set_new(solar, "yield_total_kwh",  json_real(d->yield_total_kwh));
    json_object_set_new(obj, "solar", solar);

    json_object_set_new(battery, "voltage_v",   json_real(d->batt_voltage_v));
    json_object_set_new(battery, "current_a",   json_real(d->batt_current_a));
    if (d->batt_soc_pct >= 0) {
        json_object_set_new(battery, "soc_pct", json_integer(d->batt_soc_pct));
    }
    json_object_set_new(obj, "battery", battery);

    if (!isnan(d->temperature_c)) {
        json_object_set_new(ctrl, "temperature_c", json_real(d->temperature_c));
    }
    json_object_set_new(ctrl, "efficiency_pct", json_real(snap->efficiency_pct));

    if (d->relay_available) {
        json_object_set_new(ctrl, "relay_on", json_boolean(d->relay_state));
    }
    if (d->load_output_available) {
        json_object_set_new(ctrl, "load_output_on", json_boolean(d->load_output_state));
        json_object_set_new(ctrl, "load_current_a", json_real(d->load_current_a));
    }
    json_object_set_new(obj, "controller", ctrl);

    return obj;
}

json_t *json_serialize_metrics(const MetricsSnapshot *snap, const char *node_id) {
    const ControllerData *d;
    json_t *obj, *metrics;

    if (!snap) return NULL;

    d       = &snap->data;
    obj     = json_object();
    metrics = json_array();

    if (node_id) json_object_set_new(obj, "node_id", json_string(node_id));
    json_object_set_new(obj, "timestamp_ms", json_integer((json_int_t)d->timestamp_ms));

    #define ADD_METRIC(name, val, unit) do { \
        json_t *m = json_object(); \
        json_object_set_new(m, "name", json_string(name)); \
        json_object_set_new(m, "value", json_real(val)); \
        json_object_set_new(m, "unit", json_string(unit)); \
        json_array_append_new(metrics, m); \
    } while (0)

    ADD_METRIC("solar_panel_voltage",  d->pv_voltage_v,        "V");
    ADD_METRIC("solar_panel_current",  d->pv_current_a,        "A");
    ADD_METRIC("solar_panel_power",    d->pv_power_w,          "W");
    ADD_METRIC("solar_yield_today",    d->yield_today_kwh,     "kWh");
    ADD_METRIC("solar_yield_total",    d->yield_total_kwh,     "kWh");
    ADD_METRIC("battery_voltage",      d->batt_voltage_v,      "V");
    ADD_METRIC("battery_current",      d->batt_current_a,      "A");
    ADD_METRIC("mppt_efficiency",      snap->efficiency_pct,   "%");

    if (d->batt_soc_pct >= 0) {
        ADD_METRIC("battery_charge_percentage", (double)d->batt_soc_pct, "%");
    }
    if (!isnan(d->temperature_c)) {
        ADD_METRIC("controller_temperature", d->temperature_c, "C");
    }
    if (d->load_output_available) {
        ADD_METRIC("load_current", d->load_current_a, "A");
    }

    #undef ADD_METRIC

    json_object_set_new(obj, "metrics", metrics);
    return obj;
}

json_t *json_serialize_alarms(const AlarmRecord *alarms, int count) {
    json_t *arr = json_array();

    for (int i = 0; i < count; i++) {
        const AlarmRecord *a = &alarms[i];
        json_t *item = json_object();

        json_object_set_new(item, "type",         json_string(alarm_type_str(a->type)));
        json_object_set_new(item, "severity",     json_string(severity_str(a->severity)));
        json_object_set_new(item, "active",       json_boolean(a->active));
        json_object_set_new(item, "timestamp_ms", json_integer((json_int_t)a->timestamp_ms));
        json_object_set_new(item, "message",      json_string(a->message));

        json_array_append_new(arr, item);
    }

    return arr;
}

json_t *json_serialize_alarm_notification(const Config *config,
                                          AlarmType type,
                                          Severity severity,
                                          const char *message) {
    json_t *obj = json_object();

    json_object_set_new(obj, "service",  json_string("controller.d"));
    json_object_set_new(obj, "node_id",  json_string(config->nodeId ? config->nodeId : "unknown"));
    json_object_set_new(obj, "type",     json_string(alarm_type_str(type)));
    json_object_set_new(obj, "severity", json_string(severity_str(severity)));
    json_object_set_new(obj, "message",  json_string(message ? message : ""));

    return obj;
}

int json_deserialize_voltage_request(json_t *json, double *voltage_v) {
    json_t *val;

    if (!json || !voltage_v) return -1;

    val = json_object_get(json, "voltage_v");
    if (!val || !json_is_number(val)) {
        usys_log_error("jserdes: missing or invalid 'voltage_v' field");
        return -1;
    }

    *voltage_v = json_number_value(val);
    return 0;
}

int json_deserialize_relay_request(json_t *json, bool *state) {
    json_t *val;

    if (!json || !state) return -1;

    val = json_object_get(json, "state");
    if (!val || !json_is_boolean(val)) {
        usys_log_error("jserdes: missing or invalid 'state' field");
        return -1;
    }

    *state = json_boolean_value(val);
    return 0;
}

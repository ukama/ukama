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
#include "json_types.h"
#include "usys_log.h"

void json_log(json_t *json) {

    char *str;

    str = json_dumps(json, 0);
    if (str) {
        log_debug("json str: %s", str);
        free(str);
    }
}


static bool get_json_entry(json_t *json, char *key, json_type type,
                           char **strValue, int *intValue,
                           double *doubleValue) {

    json_t *jEntry;

    if (json == NULL || key == NULL) return USYS_FALSE;

    jEntry = json_object_get(json, key);
    if (jEntry == NULL) {
        log_error("Missing %s key in json", key);
        return USYS_FALSE;
    }

    switch (type) {
    case (JSON_STRING):
        *strValue = strdup(json_string_value(jEntry));
        break;
    case (JSON_INTEGER):
        *intValue = (int)json_integer_value(jEntry);
        break;
    case (JSON_REAL):
        *doubleValue = json_real_value(jEntry);
        break;
    default:
        log_error("Invalid type for json key-value: %d", type);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

json_t *json_serialize_status(const MetricsSnapshot *snap) {
    const ControllerData *d;
    json_t *obj, *solar, *battery, *ctrl;

    if (!snap) return NULL;

    d       = &snap->data;
    obj     = json_object();
    solar   = json_object();
    battery = json_object();
    ctrl    = json_object();

    json_object_set_new(obj, JSON_KEY_TIMESTAMP_MS,
                        json_integer((json_int_t)d->timestamp_ms));
    json_object_set_new(obj, JSON_KEY_COMM_OK,
                        json_boolean(d->comm_ok));
    json_object_set_new(obj, JSON_KEY_COMM_ERRORS,
                        json_integer(d->comm_errors));
    json_object_set_new(obj, JSON_KEY_CHARGE_STATE,
                        json_string(charge_state_str(d->charge_state)));
    json_object_set_new(obj, JSON_KEY_ERROR_CODE,
                        json_integer(d->error_code));
    json_object_set_new(obj, JSON_KEY_ERROR_STR,
                        json_string(error_code_str(d->error_code)));
    json_object_set_new(obj, JSON_KEY_OVERALL_SEVERITY,
                        json_string(severity_str(snap->overall_severity)));
    json_object_set_new(obj, JSON_KEY_ACTIVE_ALARM_COUNT,
                        json_integer(snap->active_alarm_count));

    if (d->firmware[0]) {
        json_object_set_new(obj, JSON_KEY_FIRMWARE,
                            json_string(d->firmware));
    }
    if (d->serial[0]) {
        json_object_set_new(obj, JSON_KEY_SERIAL,
                            json_string(d->serial));
    }
    if (d->product_id[0]) {
        json_object_set_new(obj, JSON_KEY_PRODUCT_ID,
                            json_string(d->product_id));
    }

    json_object_set_new(solar, JSON_KEY_VOLTAGE_V,
                        json_real(d->pv_voltage_v));
    json_object_set_new(solar, JSON_KEY_CURRENT_A,
                        json_real(d->pv_current_a));
    json_object_set_new(solar, JSON_KEY_POWER_W,
                        json_real(d->pv_power_w));
    json_object_set_new(solar, JSON_KEY_YIELD_TODAY_KWH,
                        json_real(d->yield_today_kwh));
    json_object_set_new(solar, JSON_KEY_YIELD_TOTAL_KWH,
                        json_real(d->yield_total_kwh));
    json_object_set_new(obj, JSON_KEY_SOLAR, solar);

    json_object_set_new(battery, JSON_KEY_VOLTAGE_V,
                        json_real(d->batt_voltage_v));
    json_object_set_new(battery, JSON_KEY_CURRENT_A,
                        json_real(d->batt_current_a));
    if (d->batt_soc_pct >= 0) {
        json_object_set_new(battery, JSON_KEY_SOC_PCT,
                            json_integer(d->batt_soc_pct));
    }
    json_object_set_new(obj, JSON_KEY_BATTERY, battery);

    if (!isnan(d->temperature_c)) {
        json_object_set_new(ctrl, JSON_KEY_TEMPERATURE_C,
                            json_real(d->temperature_c));
    }
    json_object_set_new(ctrl, JSON_KEY_EFFICIENCY_PCT,
                        json_real(snap->efficiency_pct));

    if (d->relay_available) {
        json_object_set_new(ctrl, JSON_KEY_RELAY_ON,
                            json_boolean(d->relay_state));
    }
    if (d->load_output_available) {
        json_object_set_new(ctrl, JSON_KEY_LOAD_OUTPUT_ON,
                            json_boolean(d->load_output_state));
        json_object_set_new(ctrl, JSON_KEY_LOAD_CURRENT_A,
                            json_real(d->load_current_a));
    }
    json_object_set_new(obj, JSON_KEY_CONTROLLER, ctrl);

    return obj;
}

json_t *json_serialize_metrics(const MetricsSnapshot *snap,
                               const char *node_id) {
    const ControllerData *d;
    json_t *obj, *metrics;

    if (!snap) return NULL;

    d       = &snap->data;
    obj     = json_object();
    metrics = json_array();

    if (node_id) {
        json_object_set_new(obj, JSON_KEY_NODE_ID,
                            json_string(node_id));
    }
    json_object_set_new(obj, JSON_KEY_TIMESTAMP_MS,
                        json_integer((json_int_t)d->timestamp_ms));

#define ADD_METRIC(metricName, metricValue, metricUnit) do { \
    json_t *m = json_object(); \
    json_object_set_new(m, JSON_KEY_NAME, json_string(metricName)); \
    json_object_set_new(m, JSON_KEY_VALUE, json_real(metricValue)); \
    json_object_set_new(m, JSON_KEY_UNIT, json_string(metricUnit)); \
    json_array_append_new(metrics, m); \
} while (0)

    ADD_METRIC(JSON_METRIC_SOLAR_PANEL_VOLTAGE,
               d->pv_voltage_v, JSON_UNIT_V);
    ADD_METRIC(JSON_METRIC_SOLAR_PANEL_CURRENT,
               d->pv_current_a, JSON_UNIT_A);
    ADD_METRIC(JSON_METRIC_SOLAR_PANEL_POWER,
               d->pv_power_w, JSON_UNIT_W);
    ADD_METRIC(JSON_METRIC_SOLAR_YIELD_TODAY,
               d->yield_today_kwh, JSON_UNIT_KWH);
    ADD_METRIC(JSON_METRIC_SOLAR_YIELD_TOTAL,
               d->yield_total_kwh, JSON_UNIT_KWH);
    ADD_METRIC(JSON_METRIC_BATTERY_VOLTAGE,
               d->batt_voltage_v, JSON_UNIT_V);
    ADD_METRIC(JSON_METRIC_BATTERY_CURRENT,
               d->batt_current_a, JSON_UNIT_A);
    ADD_METRIC(JSON_METRIC_MPPT_EFFICIENCY,
               snap->efficiency_pct, JSON_UNIT_PERCENT);

    if (d->batt_soc_pct >= 0) {
        ADD_METRIC(JSON_METRIC_BATTERY_CHARGE_PERCENTAGE,
                   (double)d->batt_soc_pct, JSON_UNIT_PERCENT);
    }
    if (!isnan(d->temperature_c)) {
        ADD_METRIC(JSON_METRIC_CONTROLLER_TEMPERATURE,
                   d->temperature_c, JSON_UNIT_C);
    }
    if (d->load_output_available) {
        ADD_METRIC(JSON_METRIC_LOAD_CURRENT,
                   d->load_current_a, JSON_UNIT_A);
    }

#undef ADD_METRIC

    json_object_set_new(obj, JSON_KEY_METRICS, metrics);
    return obj;
}

json_t *json_serialize_alarms(const AlarmRecord *alarms, int count) {
    json_t *arr = json_array();

    for (int i = 0; i < count; i++) {
        const AlarmRecord *a = &alarms[i];
        json_t *item = json_object();

        json_object_set_new(item, JSON_KEY_TYPE,
                            json_string(alarm_type_str(a->type)));
        json_object_set_new(item, JSON_KEY_SEVERITY,
                            json_string(severity_str(a->severity)));
        json_object_set_new(item, JSON_KEY_ACTIVE,
                            json_boolean(a->active));
        json_object_set_new(item, JSON_KEY_TIMESTAMP_MS,
                            json_integer((json_int_t)a->timestamp_ms));
        json_object_set_new(item, JSON_KEY_MESSAGE,
                            json_string(a->message));

        json_array_append_new(arr, item);
    }

    return arr;
}

json_t *json_serialize_alarm_notification(const Config *config,
                                          AlarmType type,
                                          Severity severity,
                                          const char *message) {
    json_t *obj = json_object();

    json_object_set_new(obj, JSON_KEY_SERVICE,
                        json_string(JSON_VAL_SERVICE_CONTROLLERD));
    json_object_set_new(obj, JSON_KEY_NODE_ID,
                        json_string(config->nodeId ?
                                    config->nodeId : JSON_VAL_UNKNOWN));
    json_object_set_new(obj, JSON_KEY_TYPE,
                        json_string(alarm_type_str(type)));
    json_object_set_new(obj, JSON_KEY_SEVERITY,
                        json_string(severity_str(severity)));
    json_object_set_new(obj, JSON_KEY_MESSAGE,
                        json_string(message ? message : ""));

    return obj;
}

int json_deserialize_voltage_request(json_t *json, double *voltage_v) {
    json_t *val;

    if (!json || !voltage_v) return -1;

    val = json_object_get(json, JSON_KEY_VOLTAGE_V);
    if (!val || !json_is_number(val)) {
        usys_log_error("jserdes: missing or invalid '%s' field",
                       JSON_KEY_VOLTAGE_V);
        return -1;
    }

    *voltage_v = json_number_value(val);
    return 0;
}

int json_deserialize_relay_request(json_t *json, bool *state) {
    json_t *val;

    if (!json || !state) return -1;

    val = json_object_get(json, JSON_KEY_STATE);
    if (!val || !json_is_boolean(val)) {
        usys_log_error("jserdes: missing or invalid '%s' field",
                       JSON_KEY_STATE);
        return -1;
    }

    *state = json_boolean_value(val);
    return 0;
}

bool json_deserialize_node_info(char **data,
                                int  *iData,
                                char *tag,
                                json_type type,
                                json_t *json) {

    json_t *jNodeInfo;

    if (json == NULL) return USYS_FALSE;

    jNodeInfo = json_object_get(json, JTAG_NODE_INFO);
    if (jNodeInfo == NULL) {
        log_error("Missing mandatory %s from JSON", JTAG_NODE_INFO);
        return USYS_FALSE;
    }

    if (get_json_entry(jNodeInfo, tag, type, data, iData, NULL) == USYS_FALSE) {
        log_error("Error deserializing node info. tag: %s", tag);
        json_log(json);
        *data = NULL;
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

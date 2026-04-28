/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <jansson.h>
#include <stdio.h>
#include <string.h>

#include "jserdes.h"
#include "driver.h"
#include "json_types.h"
#include "utils.h"

static bool load_request_json(const URequest *request, JsonObj **json) {

    JsonErrObj err;

    if (request == NULL || json == NULL ||
        request->binary_body == NULL || request->binary_body_length == 0) {
        return false;
    }

    memset(&err, 0, sizeof(err));

    *json = json_loadb((const char *)request->binary_body,
                       request->binary_body_length,
                       0,
                       &err);

    return (*json != NULL);
}

bool json_deserialize_bool_request(const URequest *request,
                                   const char *key,
                                   bool *value) {

    JsonObj *json;
    JsonObj *entry;

    if (key == NULL || value == NULL) {
        return false;
    }

    json = NULL;
    if (!load_request_json(request, &json)) {
        return false;
    }

    entry = json_object_get(json, key);
    if (entry == NULL || !json_is_boolean(entry)) {
        json_free(&json);
        return false;
    }

    *value = json_is_true(entry);
    json_free(&json);

    return true;
}

bool json_deserialize_int_request(const URequest *request,
                                  const char *key,
                                  int *value) {

    JsonObj *json;
    JsonObj *entry;

    if (key == NULL || value == NULL) {
        return false;
    }

    json = NULL;
    if (!load_request_json(request, &json)) {
        return false;
    }

    entry = json_object_get(json, key);
    if (entry == NULL || !json_is_integer(entry)) {
        json_free(&json);
        return false;
    }

    *value = (int)json_integer_value(entry);
    json_free(&json);

    return true;
}

bool json_deserialize_firmware_stage_request(const URequest *request,
                                             char *path,
                                             size_t pathLen,
                                             char *version,
                                             size_t versionLen,
                                             char *sha256,
                                             size_t sha256Len) {

    JsonObj *json;
    JsonObj *entry;
    const char *value;

    if (path == NULL || version == NULL || sha256 == NULL) {
        return false;
    }

    path[0] = '\0';
    version[0] = '\0';
    sha256[0] = '\0';

    json = NULL;
    if (!load_request_json(request, &json)) {
        return false;
    }

    entry = json_object_get(json, JTAG_PATH);
    if (entry == NULL || !json_is_string(entry)) {
        json_free(&json);
        return false;
    }

    value = json_string_value(entry);
    snprintf(path, pathLen, "%s", value ? value : "");

    entry = json_object_get(json, JTAG_VERSION);
    if (entry && json_is_string(entry)) {
        value = json_string_value(entry);
        snprintf(version, versionLen, "%s", value ? value : "");
    }

    entry = json_object_get(json, JTAG_SHA256);
    if (entry && json_is_string(entry)) {
        value = json_string_value(entry);
        snprintf(sha256, sha256Len, "%s", value ? value : "");
    }

    json_free(&json);

    return true;
}

bool json_serialize_alarm_notification(JsonObj **json,
                                       const char *serviceName,
                                       const SwitchAlarm *alarm,
                                       bool clear) {

    if (json == NULL || serviceName == NULL || alarm == NULL) {
        return false;
    }

    *json = json_object();
    if (*json == NULL) {
        return false;
    }

    json_object_set_new(*json, JTAG_SERVICE_NAME, json_string(serviceName));
    json_object_set_new(*json, JTAG_CLEAR, json_boolean(clear));
    json_object_set_new(*json, JTAG_CODE, json_integer(alarm->code));
    json_object_set_new(*json,
                        JTAG_SEVERITY,
                        json_string(alarm_severity_to_str(alarm->severity)));
    json_object_set_new(*json, JTAG_RESOURCE, json_string(alarm->resource));
    json_object_set_new(*json, JTAG_TEXT, json_string(alarm->text));
    json_object_set_new(*json, JTAG_TIME, json_integer(time(NULL)));

    return true;
}

JsonObj *json_serialize_switch_info(const SwitchdContext *ctx) {

    JsonObj *json;

    json = json_object();
    json_object_set_new(json, "vendor", json_string(ctx->info.vendor));
    json_object_set_new(json, JTAG_MODEL, json_string(ctx->info.model));
    json_object_set_new(json, "serial", json_string(ctx->info.serial));
    json_object_set_new(json,
                        "hardwareVersion",
                        json_string(ctx->info.hardwareVersion));
    json_object_set_new(json,
                        JTAG_SOFTWARE_VERSION,
                        json_string(ctx->info.softwareVersion));
    json_object_set_new(json,
                        "managementAddress",
                        json_string(ctx->info.managementAddress));
    json_object_set_new(json, JTAG_REACHABLE, json_boolean(ctx->info.reachable));
    json_object_set_new(json, JTAG_PORT_COUNT, json_integer(ctx->portCount));
    json_object_set_new(json, "updatedAt", json_integer(ctx->info.updatedAt));

    return json;
}

JsonObj *json_serialize_switch_health(const SwitchdContext *ctx) {

    JsonObj *json;

    json = json_object();
    json_object_set_new(json, JTAG_REACHABLE, json_boolean(ctx->info.reachable));
    json_object_set_new(json,
                        JTAG_STATE,
                        json_string(state_to_str(ctx->state)));
    json_object_set_new(json,
                        "pollFailures",
                        json_integer(ctx->info.pollFailures));
    json_object_set_new(json,
                        "poeTotalPowerWatts",
                        json_real(ctx->kpis.poeTotalPowerWatts));
    json_object_set_new(json,
                        "poeMaxPowerWatts",
                        json_real(ctx->kpis.poeMaxPowerWatts));
    json_object_set_new(json,
                        "systemTemperatureC",
                        json_real(ctx->kpis.systemTemperatureC));
    json_object_set_new(json,
                        "ambientTemperatureC",
                        json_real(ctx->kpis.ambientTemperatureC));
    json_object_set_new(json,
                        "systemPowerWatts",
                        json_real(ctx->kpis.systemPowerWatts));

    return json;
}

JsonObj *json_serialize_switch_capabilities(const SwitchdContext *ctx) {

    JsonObj *json;

    json = json_object();
    json_object_set_new(json,
                        "supportsPortAdmin",
                        json_boolean(ctx->caps.supportsPortAdmin));
    json_object_set_new(json,
                        "supportsPoeControl",
                        json_boolean(ctx->caps.supportsPoeControl));
    json_object_set_new(json,
                        "supportsPoeCycle",
                        json_boolean(ctx->caps.supportsPoeCycle));
    json_object_set_new(json,
                        "supportsPortCounters",
                        json_boolean(ctx->caps.supportsPortCounters));
    json_object_set_new(json,
                        "supportsPowerMetrics",
                        json_boolean(ctx->caps.supportsPowerMetrics));
    json_object_set_new(json,
                        "supportsSystemMetrics",
                        json_boolean(ctx->caps.supportsSystemMetrics));
    json_object_set_new(json,
                        "supportsFirmwareUpdate",
                        json_boolean(ctx->caps.supportsFirmwareUpdate));
    json_object_set_new(json,
                        "supportsSaveConfig",
                        json_boolean(ctx->caps.supportsSaveConfig));
    json_object_set_new(json, "maxPorts", json_integer(ctx->caps.maxPorts));

    return json;
}

JsonObj *json_serialize_switch_kpis(const SwitchdContext *ctx) {

    JsonObj *json;

    json = json_object();
    json_object_set_new(json,
                        "poeTotalPowerWatts",
                        json_real(ctx->kpis.poeTotalPowerWatts));
    json_object_set_new(json,
                        "poeMaxPowerWatts",
                        json_real(ctx->kpis.poeMaxPowerWatts));
    json_object_set_new(json,
                        "systemTemperatureC",
                        json_real(ctx->kpis.systemTemperatureC));
    json_object_set_new(json,
                        "ambientTemperatureC",
                        json_real(ctx->kpis.ambientTemperatureC));
    json_object_set_new(json,
                        "systemPowerWatts",
                        json_real(ctx->kpis.systemPowerWatts));
    json_object_set_new(json,
                        "inputVoltage",
                        json_real(ctx->kpis.inputVoltage));
    json_object_set_new(json,
                        "systemCurrentAmps",
                        json_real(ctx->kpis.systemCurrentAmps));
    json_object_set_new(json,
                        "inputLinkFailureAlarm",
                        json_boolean(ctx->kpis.inputLinkFailureAlarm));
    json_object_set_new(json,
                        "inputPoeFailureAlarm",
                        json_boolean(ctx->kpis.inputPoeFailureAlarm));
    json_object_set_new(json, "updatedAt", json_integer(ctx->kpis.updatedAt));

    return json;
}

JsonObj *json_serialize_port(const SwitchPortState *port) {

    JsonObj *json;

    json = json_object();
    json_object_set_new(json, JTAG_ID, json_integer(port->id));
    json_object_set_new(json, "name", json_string(port->name));
    json_object_set_new(json, "media", json_string(port->media));
    json_object_set_new(json, "present", json_boolean(port->present));
    json_object_set_new(json,
                        "adminState",
                        json_string(port->adminUp ? "up" : "down"));
    json_object_set_new(json,
                        "linkState",
                        json_string(port->linkUp ? "up" : "down"));
    json_object_set_new(json,
                        "poeSupported",
                        json_boolean(port->poeSupported));
    json_object_set_new(json,
                        "poeState",
                        json_string(port->poeEnabled ? "on" : "off"));
    json_object_set_new(json,
                        "poeOperational",
                        json_boolean(port->poeOperational));
    json_object_set_new(json, "poeClass", json_integer(port->poeClass));
    json_object_set_new(json, "speedBps", json_integer(port->speedBps));
    json_object_set_new(json, "powerWatts", json_real(port->powerWatts));
    json_object_set_new(json, "voltage", json_real(port->voltage));
    json_object_set_new(json, "currentAmps", json_real(port->currentAmps));
    json_object_set_new(json, "rxBytes", json_integer(port->rxBytes));
    json_object_set_new(json, "txBytes", json_integer(port->txBytes));
    json_object_set_new(json, "rxPackets", json_integer(port->rxPackets));
    json_object_set_new(json, "txPackets", json_integer(port->txPackets));
    json_object_set_new(json, "rxErrors", json_integer(port->rxErrors));
    json_object_set_new(json, "txErrors", json_integer(port->txErrors));
    json_object_set_new(json, "rxDrops", json_integer(port->rxDrops));
    json_object_set_new(json, "txDrops", json_integer(port->txDrops));
    json_object_set_new(json, "fault", json_string(port->fault));
    json_object_set_new(json, "updatedAt", json_integer(port->updatedAt));

    return json;
}

JsonObj *json_serialize_ports(const SwitchdContext *ctx) {

    JsonObj *json;
    uint32_t i;

    json = json_array();

    for (i = 0; i < ctx->portCount && i < SWITCHD_MAX_PORTS; i++) {
        json_array_append_new(json, json_serialize_port(&ctx->ports[i]));
    }

    return json;
}

JsonObj *json_serialize_firmware(const SwitchdContext *ctx) {

    JsonObj *json;

    json = json_object();
    json_object_set_new(json, JTAG_PATH, json_string(ctx->fw.path));
    json_object_set_new(json,
                        JTAG_VERSION,
                        json_string(ctx->fw.version));
    json_object_set_new(json, JTAG_SHA256, json_string(ctx->fw.sha256));
    json_object_set_new(json,
                        "tftpFilename",
                        json_string(ctx->fw.tftpFilename));
    json_object_set_new(json, "size", json_integer(ctx->fw.size));
    json_object_set_new(json,
                        JTAG_STATE,
                        json_string(fw_state_to_str(ctx->fw.state)));
    json_object_set_new(json,
                        "executeStatus",
                        json_integer(ctx->fw.executeStatus));
    json_object_set_new(json, JTAG_DETAIL, json_string(ctx->fw.detail));
    json_object_set_new(json, "stagedAt", json_integer(ctx->fw.stagedAt));
    json_object_set_new(json, "updatedAt", json_integer(ctx->fw.updatedAt));

    return json;
}

JsonObj *json_serialize_active_alarms(const SwitchdContext *ctx) {

    JsonObj *json;
    JsonObj *alarm;
    size_t i;

    json = json_array();

    for (i = 0; i < ctx->alarmCount; i++) {
        if (!ctx->alarms[i].active) {
            continue;
        }

        alarm = json_object();
        json_object_set_new(alarm, JTAG_CODE, json_integer(ctx->alarms[i].code));
        json_object_set_new(alarm,
                            JTAG_SEVERITY,
                            json_string(alarm_severity_to_str(
                                ctx->alarms[i].severity)));
        json_object_set_new(alarm,
                            JTAG_RESOURCE,
                            json_string(ctx->alarms[i].resource));
        json_object_set_new(alarm, JTAG_TEXT, json_string(ctx->alarms[i].text));
        json_array_append_new(json, alarm);
    }

    return json;
}

static void metric_add(JsonObj *metrics,
                       const char *name,
                       double value,
                       const char *unit) {

    JsonObj *metric;

    if (!metrics || !name || !unit) {
        return;
    }

    metric = json_object();
    json_object_set_new(metric, JSON_KEY_NAME, json_string(name));
    json_object_set_new(metric, JSON_KEY_VALUE, json_real(value));
    json_object_set_new(metric, JSON_KEY_UNIT, json_string(unit));
    json_array_append_new(metrics, metric);
}

static void metric_name_for_port(char *buf,
                                 size_t len,
                                 const SwitchPortState *port,
                                 const char *suffix) {

    char role[SWITCHD_NAME_LEN];
    size_t i;

    if (!buf || len == 0 || !port || !suffix) {
        return;
    }

    snprintf(role, sizeof(role), "%s", port->name[0] ?
             port->name : "unknown");

    for (i = 0; role[i] != '\0'; i++) {
        if (role[i] == '-' || role[i] == ' ' || role[i] == '.') {
            role[i] = '_';
        }
    }

    snprintf(buf,
             len,
             "port_%u_%s_%s",
             port->id,
             role,
             suffix);
}

static void metrics_add_port(JsonObj *metrics, const SwitchPortState *port) {

    char name[128];

    if (!metrics || !port || port->id == 0) {
        return;
    }

    metric_name_for_port(name, sizeof(name), port, "present");
    metric_add(metrics, name, port->present ? 1.0 : 0.0, JSON_UNIT_BOOL);

    metric_name_for_port(name, sizeof(name), port, "admin_up");
    metric_add(metrics, name, port->adminUp ? 1.0 : 0.0, JSON_UNIT_BOOL);

    metric_name_for_port(name, sizeof(name), port, "link_up");
    metric_add(metrics, name, port->linkUp ? 1.0 : 0.0, JSON_UNIT_BOOL);

    metric_name_for_port(name, sizeof(name), port, "speed_bps");
    metric_add(metrics, name, (double)port->speedBps, "bps");

    metric_name_for_port(name, sizeof(name), port, "rx_bytes");
    metric_add(metrics, name, (double)port->rxBytes, "bytes");

    metric_name_for_port(name, sizeof(name), port, "tx_bytes");
    metric_add(metrics, name, (double)port->txBytes, "bytes");

    metric_name_for_port(name, sizeof(name), port, "rx_packets");
    metric_add(metrics, name, (double)port->rxPackets, "packets");

    metric_name_for_port(name, sizeof(name), port, "tx_packets");
    metric_add(metrics, name, (double)port->txPackets, "packets");

    metric_name_for_port(name, sizeof(name), port, "rx_errors");
    metric_add(metrics, name, (double)port->rxErrors, "count");

    metric_name_for_port(name, sizeof(name), port, "tx_errors");
    metric_add(metrics, name, (double)port->txErrors, "count");

    metric_name_for_port(name, sizeof(name), port, "rx_drops");
    metric_add(metrics, name, (double)port->rxDrops, "count");

    metric_name_for_port(name, sizeof(name), port, "tx_drops");
    metric_add(metrics, name, (double)port->txDrops, "count");

    metric_name_for_port(name, sizeof(name), port, "poe_supported");
    metric_add(metrics, name, port->poeSupported ? 1.0 : 0.0, JSON_UNIT_BOOL);

    metric_name_for_port(name, sizeof(name), port, "poe_enabled");
    metric_add(metrics, name, port->poeEnabled ? 1.0 : 0.0, JSON_UNIT_BOOL);

    metric_name_for_port(name, sizeof(name), port, "poe_operational");
    metric_add(metrics, name, port->poeOperational ? 1.0 : 0.0,
               JSON_UNIT_BOOL);

    metric_name_for_port(name, sizeof(name), port, "poe_class");
    metric_add(metrics, name, (double)port->poeClass, "class");

    metric_name_for_port(name, sizeof(name), port, "poe_power_watts");
    metric_add(metrics, name, port->powerWatts, JSON_UNIT_W);

    metric_name_for_port(name, sizeof(name), port, "poe_voltage");
    metric_add(metrics, name, port->voltage, JSON_UNIT_V);

    metric_name_for_port(name, sizeof(name), port, "poe_current_amps");
    metric_add(metrics, name, port->currentAmps, JSON_UNIT_A);
}

JsonObj *json_serialize_metrics(const SwitchdContext *ctx) {

    JsonObj *obj;
    JsonObj *metrics;
    uint32_t i;
    time_t updatedAt;

    if (ctx == NULL) {
        return NULL;
    }

    obj = json_object();
    metrics = json_array();

    updatedAt = ctx->kpis.updatedAt;
    if (ctx->portCount > 0 && ctx->ports[0].updatedAt > updatedAt) {
        updatedAt = ctx->ports[0].updatedAt;
    }

    json_object_set_new(obj,
                        JSON_KEY_TIMESTAMP_MS,
                        json_integer((json_int_t)updatedAt));

    metric_add(metrics,
               JSON_METRIC_POE_TOTAL_POWER_WATTS,
               ctx->kpis.poeTotalPowerWatts,
               JSON_UNIT_W);
    metric_add(metrics,
               JSON_METRIC_POE_MAX_POWER_WATTS,
               ctx->kpis.poeMaxPowerWatts,
               JSON_UNIT_W);
    metric_add(metrics,
               JSON_METRIC_SYSTEM_TEMPERATURE_C,
               ctx->kpis.systemTemperatureC,
               JSON_UNIT_C);
    metric_add(metrics,
               JSON_METRIC_AMBIENT_TEMPERATURE_C,
               ctx->kpis.ambientTemperatureC,
               JSON_UNIT_C);
    metric_add(metrics,
               JSON_METRIC_SYSTEM_POWER_WATTS,
               ctx->kpis.systemPowerWatts,
               JSON_UNIT_W);
    metric_add(metrics,
               JSON_METRIC_INPUT_VOLTAGE,
               ctx->kpis.inputVoltage,
               JSON_UNIT_V);
    metric_add(metrics,
               JSON_METRIC_SYSTEM_CURRENT_AMPS,
               ctx->kpis.systemCurrentAmps,
               JSON_UNIT_A);
    metric_add(metrics,
               JSON_METRIC_INPUT_LINK_FAILURE_ALARM,
               ctx->kpis.inputLinkFailureAlarm ? 1.0 : 0.0,
               JSON_UNIT_BOOL);
    metric_add(metrics,
               JSON_METRIC_INPUT_POE_FAILURE_ALARM,
               ctx->kpis.inputPoeFailureAlarm ? 1.0 : 0.0,
               JSON_UNIT_BOOL);

    for (i = 0; i < ctx->portCount && i < SWITCHD_MAX_PORTS; i++) {
        metrics_add_port(metrics, &ctx->ports[i]);
    }

    json_object_set_new(obj, JSON_KEY_METRICS, metrics);

    return obj;
}

JsonObj *json_serialize_status(const SwitchdContext *ctx) {

    JsonObj *json;
    JsonObj *switchd;
    JsonObj *operation;

    json = json_object();

    switchd = json_object();
    json_object_set_new(switchd,
                        JTAG_STATE,
                        json_string(state_to_str(ctx->state)));
    json_object_set_new(switchd,
                        JTAG_DRIVER,
                        json_string(ctx->driver ? ctx->driver->name : ""));
    json_object_set_new(switchd,
                        JTAG_REACHABLE,
                        json_boolean(ctx->info.reachable));
    json_object_set_new(switchd,
                        JTAG_UPDATE_IN_PROGRESS,
                        json_boolean(ctx->fw.state == SWITCHD_FW_APPLYING ||
                                     ctx->fw.state == SWITCHD_FW_RECONNECTING ||
                                     ctx->fw.state == SWITCHD_FW_VERIFYING));
    json_object_set_new(json, JTAG_SWITCHD, switchd);

    operation = json_object();
    json_object_set_new(operation, JTAG_ID, json_integer(ctx->op.id));
    json_object_set_new(operation,
                        JTAG_TYPE,
                        json_string(op_type_to_str(ctx->op.type)));
    json_object_set_new(operation,
                        JTAG_STATE,
                        json_string(op_state_to_str(ctx->op.state)));
    json_object_set_new(operation,
                        JTAG_PROGRESS,
                        json_integer(ctx->op.progress));
    json_object_set_new(operation,
                        JTAG_DETAIL,
                        json_string(ctx->op.detail));
    json_object_set_new(json, JTAG_OPERATION, operation);

    json_object_set_new(json, JTAG_SWITCH, json_serialize_switch_info(ctx));
    json_object_set_new(json, JTAG_ALARMS, json_serialize_active_alarms(ctx));

    return json;
}

void json_free(JsonObj **json) {

    if (json && *json) {
        json_decref(*json);
        *json = NULL;
    }
}

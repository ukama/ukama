/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <math.h>

#include <ulfius.h>
#include <jansson.h>

#include "alarms.h"
#include "jserdes.h"
#include "http_status.h"
#include "usys_log.h"

#define DEBOUNCE_COUNT 3  /* Number of consecutive samples to trigger alarm */

int alarms_init(AlarmChecker *checker, const Config *config, MetricsStore *store) {
    if (!checker || !config || !store) return -1;

    memset(checker, 0, sizeof(*checker));
    checker->config = config;
    checker->store = store;
    checker->debounce_samples = DEBOUNCE_COUNT;

    return 0;
}

void alarms_check(AlarmChecker *checker, const ControllerData *data) {
    if (!checker || !data || !checker->config || !checker->store) return;

    const Config *cfg = checker->config;

    if (data->batt_voltage_v > 0) {
        if (data->batt_voltage_v < cfg->lowVoltageCrit) {
            checker->low_volt_count++;
            if (checker->low_volt_count >= checker->debounce_samples) {
                char msg[128];
                snprintf(msg, sizeof(msg), "Battery voltage critical: %.2fV < %.2fV",
                         data->batt_voltage_v, cfg->lowVoltageCrit);
                metrics_store_set_alarm(checker->store, ALARM_LOW_BATTERY_VOLTAGE,
                                        SEVERITY_CRITICAL, msg);
                alarms_send_notification(cfg, ALARM_LOW_BATTERY_VOLTAGE,
                                         SEVERITY_CRITICAL, msg);
            }
        } else if (data->batt_voltage_v < cfg->lowVoltageWarn) {
            checker->low_volt_count++;
            if (checker->low_volt_count >= checker->debounce_samples) {
                char msg[128];
                snprintf(msg, sizeof(msg), "Battery voltage low: %.2fV < %.2fV",
                         data->batt_voltage_v, cfg->lowVoltageWarn);
                metrics_store_set_alarm(checker->store, ALARM_LOW_BATTERY_VOLTAGE,
                                        SEVERITY_WARN, msg);
            }
        } else {
            if (checker->low_volt_count > 0) {
                checker->low_volt_count--;
                if (checker->low_volt_count == 0) {
                    metrics_store_clear_alarm(checker->store, ALARM_LOW_BATTERY_VOLTAGE);
                }
            }
        }
    }

    if (!isnan(data->temperature_c)) {
        if (data->temperature_c > cfg->highTempCrit) {
            checker->high_temp_count++;
            if (checker->high_temp_count >= checker->debounce_samples) {
                char msg[128];
                snprintf(msg, sizeof(msg), "Controller temperature critical: %.1f°C > %.1f°C",
                         data->temperature_c, cfg->highTempCrit);
                metrics_store_set_alarm(checker->store, ALARM_HIGH_TEMPERATURE,
                                        SEVERITY_CRITICAL, msg);
                alarms_send_notification(cfg, ALARM_HIGH_TEMPERATURE,
                                         SEVERITY_CRITICAL, msg);
            }
        } else if (data->temperature_c > cfg->highTempWarn) {
            checker->high_temp_count++;
            if (checker->high_temp_count >= checker->debounce_samples) {
                char msg[128];
                snprintf(msg, sizeof(msg), "Controller temperature high: %.1f°C > %.1f°C",
                         data->temperature_c, cfg->highTempWarn);
                metrics_store_set_alarm(checker->store, ALARM_HIGH_TEMPERATURE,
                                        SEVERITY_WARN, msg);
            }
        } else {
            if (checker->high_temp_count > 0) {
                checker->high_temp_count--;
                if (checker->high_temp_count == 0) {
                    metrics_store_clear_alarm(checker->store, ALARM_HIGH_TEMPERATURE);
                }
            }
        }
    }

    if (data->error_code != 0) {
        checker->fault_count++;
        if (checker->fault_count >= checker->debounce_samples) {
            char msg[128];
            snprintf(msg, sizeof(msg), "Controller error: %s (code %u)",
                     error_code_str(data->error_code), data->error_code);
            metrics_store_set_alarm(checker->store, ALARM_CONTROLLER_FAULT,
                                    SEVERITY_CRITICAL, msg);
            alarms_send_notification(cfg, ALARM_CONTROLLER_FAULT, SEVERITY_CRITICAL, msg);
        }
    } else {
        if (checker->fault_count > 0) {
            checker->fault_count--;
            if (checker->fault_count == 0) {
                metrics_store_clear_alarm(checker->store, ALARM_CONTROLLER_FAULT);
            }
        }
    }

    if (data->comm_ok) {
        checker->comm_fail_count = 0;
        metrics_store_clear_alarm(checker->store, ALARM_COMMUNICATION_LOST);
    }
}

void alarms_check_comm_failure(AlarmChecker *checker, bool comm_ok) {
    if (!checker || !checker->config || !checker->store) return;

    if (!comm_ok) {
        checker->comm_fail_count++;
        if (checker->comm_fail_count >= checker->debounce_samples) {
            char msg[128];
            snprintf(msg, sizeof(msg), "Communication lost with charge controller");
            metrics_store_set_alarm(checker->store, ALARM_COMMUNICATION_LOST,
                                    SEVERITY_CRITICAL, msg);
            alarms_send_notification(checker->config, ALARM_COMMUNICATION_LOST,
                                     SEVERITY_CRITICAL, msg);
        }
    } else {
        if (checker->comm_fail_count > 0) {
            checker->comm_fail_count--;
            if (checker->comm_fail_count == 0) {
                metrics_store_clear_alarm(checker->store, ALARM_COMMUNICATION_LOST);
            }
        }
    }
}

int alarms_send_notification(const Config *config, AlarmType type,
                             Severity severity, const char *message) {
    char url[256];
    char *body = NULL;
    json_t *json = NULL;
    struct _u_request req;
    struct _u_response resp;
    int ret = -1;

    if (!config || !config->enableNotify) return 0;

    if (snprintf(url, sizeof(url), "http://%s:%d%s",
                 config->notifyHost, config->notifyPort, config->notifyPath)
        >= (int)sizeof(url)) {
        return -1;
    }

    json = json_serialize_alarm_notification(config, type, severity, message);
    if (!json) {
        usys_log_error("alarms: failed to serialize notification");
        return -1;
    }

    body = json_dumps(json, 0);
    json_decref(json);

    if (!body) {
        usys_log_error("alarms: json_dumps failed");
        return -1;
    }

    ulfius_init_request(&req);
    ulfius_init_response(&resp);

    req.http_url  = strdup(url);
    req.http_verb = strdup("POST");
    u_map_put(req.map_header, "Content-Type", "application/json");
    ulfius_set_string_body_request(&req, body);

    if (ulfius_send_http_request(&req, &resp) != U_OK) {
        usys_log_error("alarms: failed to send notification to %s", url);
    } else if (resp.status != HttpStatus_OK && resp.status != HttpStatus_Accepted) {
        usys_log_error("alarms: notify.d returned status %d", resp.status);
    } else {
        usys_log_info("alarms: notification sent - type=%s severity=%s",
                      alarm_type_str(type), severity_str(severity));
        ret = 0;
    }

    ulfius_clean_response(&resp);
    ulfius_clean_request(&req);
    free(body);

    return ret;
}

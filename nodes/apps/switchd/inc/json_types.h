/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef JSON_TYPES_H_
#define JSON_TYPES_H_

#define JTAG_STATUS             "status"
#define JTAG_RESULT             "result"
#define JTAG_ERROR              "error"
#define JTAG_DETAIL             "detail"
#define JTAG_CLEAR              "clear"
#define JTAG_SOURCE             "source"
#define JTAG_CODE               "code"
#define JTAG_SEVERITY           "severity"
#define JTAG_RESOURCE           "resource"
#define JTAG_TEXT               "text"
#define JTAG_SERVICE_NAME       "serviceName"
#define JTAG_TIME               "time"
#define JTAG_MODULE             "module"
#define JTAG_NAME               "name"
#define JTAG_VALUE              "value"
#define JTAG_UNITS              "units"
#define JTAG_DETAILS            "details"
#define JTAG_VERSION            "version"
#define JTAG_STATE              "state"
#define JTAG_UP                 "up"
#define JTAG_ON                 "on"
#define JTAG_OFF_MS             "offMs"
#define JTAG_PATH               "path"
#define JTAG_SHA256             "sha256"
#define JTAG_PORT_ID            "portId"
#define JTAG_PORTS              "ports"
#define JTAG_SWITCH             "switch"
#define JTAG_SWITCHD            "switchd"
#define JTAG_OPERATION          "operation"
#define JTAG_ALARMS             "alarms"
#define JTAG_FIRMWARE           "firmware"
#define JTAG_CAPABILITIES       "capabilities"
#define JTAG_REACHABLE          "reachable"
#define JTAG_UPDATE_IN_PROGRESS "updateInProgress"
#define JTAG_DRIVER             "driver"
#define JTAG_MODEL              "model"
#define JTAG_SOFTWARE_VERSION   "softwareVersion"
#define JTAG_PORT_COUNT         "portCount"
#define JTAG_ID                 "id"
#define JTAG_TYPE               "type"
#define JTAG_PROGRESS           "progress"

#define ALARM_HIGH        "high"
#define ALARM_INFO        "info"
#define ALARM_WARNING     "warning"
#define ALARM_NODE        "node"
#define ALARM_SWITCH      "switch"
#define ALARM_MODULE_NONE "none"
#define EMPTY_STRING      ""

#define JSON_KEY_NODE_ID              "node_id"
#define JSON_KEY_TIMESTAMP_MS         "timestamp_ms"
#define JSON_KEY_METRICS              "metrics"
#define JSON_KEY_NAME                 "name"
#define JSON_KEY_VALUE                "value"
#define JSON_KEY_UNIT                 "unit"

#define JSON_UNIT_W                   "W"
#define JSON_UNIT_C                   "C"
#define JSON_UNIT_V                   "V"
#define JSON_UNIT_A                   "A"
#define JSON_UNIT_BOOL                "bool"

#define JSON_METRIC_POE_TOTAL_POWER_WATTS    "poe_total_power_watts"
#define JSON_METRIC_POE_MAX_POWER_WATTS      "poe_max_power_watts"
#define JSON_METRIC_SYSTEM_TEMPERATURE_C      "system_temperature_c"
#define JSON_METRIC_AMBIENT_TEMPERATURE_C     "ambient_temperature_c"
#define JSON_METRIC_SYSTEM_POWER_WATTS        "system_power_watts"
#define JSON_METRIC_INPUT_VOLTAGE             "input_voltage"
#define JSON_METRIC_SYSTEM_CURRENT_AMPS       "system_current_amps"
#define JSON_METRIC_INPUT_LINK_FAILURE_ALARM  "input_link_failure_alarm"
#define JSON_METRIC_INPUT_POE_FAILURE_ALARM   "input_poe_failure_alarm"

#endif /* JSON_TYPES_H_ */

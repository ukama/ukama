/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef JSON_TYPES_H
#define JSON_TYPES_H

#define JTAG_NODE_INFO          "nodeInfo"
#define JTAG_NODE_ID            "UUID"
#define JTAG_TYPE               "type"

/* top-level common keys */
#define JSON_KEY_TIMESTAMP_MS         "timestamp_ms"
#define JSON_KEY_NODE_ID              "node_id"
#define JSON_KEY_SERVICE              "service"
#define JSON_KEY_METRICS              "metrics"

/* status/common fields */
#define JSON_KEY_COMM_OK              "comm_ok"
#define JSON_KEY_COMM_ERRORS          "comm_errors"
#define JSON_KEY_CHARGE_STATE         "charge_state"
#define JSON_KEY_ERROR_CODE           "error_code"
#define JSON_KEY_ERROR_STR            "error_str"
#define JSON_KEY_OVERALL_SEVERITY     "overall_severity"
#define JSON_KEY_ACTIVE_ALARM_COUNT   "active_alarm_count"
#define JSON_KEY_FIRMWARE             "firmware"
#define JSON_KEY_SERIAL               "serial"
#define JSON_KEY_PRODUCT_ID           "product_id"

/* nested object names */
#define JSON_KEY_SOLAR                "solar"
#define JSON_KEY_BATTERY              "battery"
#define JSON_KEY_CONTROLLER           "controller"

/* measurement fields */
#define JSON_KEY_VOLTAGE_V            "voltage_v"
#define JSON_KEY_CURRENT_A            "current_a"
#define JSON_KEY_POWER_W              "power_w"
#define JSON_KEY_YIELD_TODAY_KWH      "yield_today_kwh"
#define JSON_KEY_YIELD_TOTAL_KWH      "yield_total_kwh"
#define JSON_KEY_SOC_PCT              "soc_pct"
#define JSON_KEY_TEMPERATURE_C        "temperature_c"
#define JSON_KEY_EFFICIENCY_PCT       "efficiency_pct"
#define JSON_KEY_RELAY_ON             "relay_on"
#define JSON_KEY_LOAD_OUTPUT_ON       "load_output_on"
#define JSON_KEY_LOAD_CURRENT_A       "load_current_a"

/* metric array item keys */
#define JSON_KEY_NAME                 "name"
#define JSON_KEY_VALUE                "value"
#define JSON_KEY_UNIT                 "unit"

/* alarm keys */
#define JSON_KEY_TYPE                 "type"
#define JSON_KEY_SEVERITY             "severity"
#define JSON_KEY_ACTIVE               "active"
#define JSON_KEY_MESSAGE              "message"

/* request keys */
#define JSON_KEY_STATE                "state"

/* fixed service values */
#define JSON_VAL_SERVICE_CONTROLLERD  "controller.d"
#define JSON_VAL_UNKNOWN              "unknown"

/* metric names */
#define JSON_METRIC_SOLAR_PANEL_VOLTAGE       "solar_panel_voltage"
#define JSON_METRIC_SOLAR_PANEL_CURRENT       "solar_panel_current"
#define JSON_METRIC_SOLAR_PANEL_POWER         "solar_panel_power"
#define JSON_METRIC_SOLAR_YIELD_TODAY         "solar_yield_today"
#define JSON_METRIC_SOLAR_YIELD_TOTAL         "solar_yield_total"
#define JSON_METRIC_BATTERY_VOLTAGE           "battery_voltage"
#define JSON_METRIC_BATTERY_CURRENT           "battery_current"
#define JSON_METRIC_MPPT_EFFICIENCY           "mppt_efficiency"
#define JSON_METRIC_BATTERY_CHARGE_PERCENTAGE "battery_charge_percentage"
#define JSON_METRIC_CONTROLLER_TEMPERATURE    "controller_temperature"
#define JSON_METRIC_LOAD_CURRENT              "load_current"

/* units */
#define JSON_UNIT_V                   "V"
#define JSON_UNIT_A                   "A"
#define JSON_UNIT_W                   "W"
#define JSON_UNIT_KWH                 "kWh"
#define JSON_UNIT_PERCENT             "%"
#define JSON_UNIT_C                   "C"

#endif /* JSON_TYPES_H */

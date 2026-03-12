/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "driver.h"
#include "drv_victron.h"
#include "usys_log.h"

/* List of available drivers */
static const ControllerDriver *drivers[] = {
    &victron_driver,
    /* Future drivers:
     * &epever_driver,
     * &meanwell_driver,
     */
    NULL
};

const ControllerDriver *driver_find(const char *name) {
    if (!name) return NULL;

    for (int i = 0; drivers[i] != NULL; i++) {
        if (strcmp(drivers[i]->name, name) == 0) {
            return drivers[i];
        }
    }

    return NULL;
}

void driver_list_available(void) {
    usys_log_info("Available charge controller drivers:");
    for (int i = 0; drivers[i] != NULL; i++) {
        usys_log_info("  - %s: %s", drivers[i]->name, drivers[i]->description);
    }
}

const char *charge_state_str(ChargeState state) {
    switch (state) {
    case CHARGE_STATE_OFF:          return "off";
    case CHARGE_STATE_FAULT:        return "fault";
    case CHARGE_STATE_BULK:         return "bulk";
    case CHARGE_STATE_ABSORPTION:   return "absorption";
    case CHARGE_STATE_FLOAT:        return "float";
    case CHARGE_STATE_STORAGE:      return "storage";
    case CHARGE_STATE_EQUALIZE:     return "equalize";
    case CHARGE_STATE_UNKNOWN:      return "unknown";
    default:                        return "unknown";
    }
}

const char *error_code_str(uint32_t code) {
    switch (code) {
    case VERR_NONE:                      return "none";
    case VERR_BATTERY_VOLTAGE_HIGH:      return "battery_voltage_high";
    case VERR_CHARGER_TEMP_HIGH:         return "charger_temperature_high";
    case VERR_CHARGER_OVERCURRENT:       return "charger_overcurrent";
    case VERR_CHARGER_CURRENT_REVERSED:  return "charger_current_reversed";
    case VERR_BULK_TIME_LIMIT:           return "bulk_time_limit_exceeded";
    case VERR_CURRENT_SENSOR_FAIL:       return "current_sensor_failure";
    case VERR_TERMINALS_OVERHEATED:      return "terminals_overheated";
    case VERR_CONVERTER_ISSUE:           return "converter_issue";
    case VERR_INPUT_VOLTAGE_HIGH:        return "input_voltage_high";
    case VERR_INPUT_CURRENT_HIGH:        return "input_current_high";
    case VERR_INPUT_SHUTDOWN_BATTERY:    return "input_shutdown_battery";
    case VERR_INPUT_SHUTDOWN_CURRENT:    return "input_shutdown_current";
    case VERR_LOST_COMMUNICATION:        return "lost_communication";
    case VERR_SYNC_CHARGING_CONFIG:      return "sync_charging_config_issue";
    case VERR_BMS_CONNECTION_LOST:       return "bms_connection_lost";
    case VERR_NETWORK_MISCONFIGURED:     return "network_misconfigured";
    case VERR_FACTORY_CALIBRATION:       return "factory_calibration_lost";
    case VERR_INVALID_FIRMWARE:          return "invalid_firmware";
    case VERR_USER_SETTINGS_INVALID:     return "user_settings_invalid";
    default:                             return "unknown_error";
    }
}

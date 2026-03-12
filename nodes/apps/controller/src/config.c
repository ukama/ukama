/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>

#include "config.h"
#include "controllerd.h"
#include "usys_api.h"
#include "usys_log.h"
#include "usys_file.h"
#include "usys_services.h"

static char *getenv_dup(const char *name, const char *def) {
    const char *val = getenv(name);
    if (val && *val) {
        return strdup(val);
    }
    return def ? strdup(def) : NULL;
}

static int getenv_int(const char *name, int def) {
    const char *val = getenv(name);
    if (val && *val) {
        return atoi(val);
    }
    return def;
}

static double getenv_double(const char *name, double def) {
    const char *val = getenv(name);
    if (val && *val) {
        return atof(val);
    }
    return def;
}

static bool getenv_bool(const char *name, bool def) {
    const char *val = getenv(name);
    if (val && *val) {
        if (strcasecmp(val, "true") == 0 || strcmp(val, "1") == 0) {
            return true;
        }
        if (strcasecmp(val, "false") == 0 || strcmp(val, "0") == 0) {
            return false;
        }
    }
    return def;
}

int config_load_from_env(Config *config) {
    if (!config) return -1;

    memset(config, 0, sizeof(*config));

    /* Web service configuration */
    config->listenAddr = getenv_dup("CONTROLLER_LISTEN_ADDR", DEF_LISTEN_ADDR);
    config->listenPort = (uint16_t)getenv_int("CONTROLLER_LISTEN_PORT", DEF_LISTEN_PORT);

    /* Sampling configuration */
    config->sampleMs = (uint32_t)getenv_int("CONTROLLER_SAMPLE_MS", DEF_SAMPLE_MS);

    /* Driver configuration */
    config->driverName = getenv_dup("CONTROLLER_DRIVER", "victron");
    config->serialPort = getenv_dup("CONTROLLER_SERIAL_PORT", DEF_SERIAL_PORT);
    config->baudRate = getenv_int("CONTROLLER_BAUD_RATE", DEF_BAUD_RATE);

    /* Notify.d configuration */
    config->notifyPort = usys_find_service_port(SERVICE_NOTIFY);
    if (!config->notifyPort) {
        config->notifyPort = 8080;
    }
    snprintf(config->notifyHost, sizeof(config->notifyHost), "%s",
             getenv("NOTIFY_HOST") ? getenv("NOTIFY_HOST") : DEF_NOTIFY_HOST);
    snprintf(config->notifyPath, sizeof(config->notifyPath), "%s", DEF_NOTIFY_EP);
    config->enableNotify = getenv_bool("CONTROLLER_ENABLE_NOTIFY", true);

    /* Alarm thresholds */
    config->lowVoltageWarn = getenv_double("CONTROLLER_LOW_VOLT_WARN", DEF_LOW_VOLT_WARN);
    config->lowVoltageCrit = getenv_double("CONTROLLER_LOW_VOLT_CRIT", DEF_LOW_VOLT_CRIT);
    config->highTempWarn = getenv_double("CONTROLLER_HIGH_TEMP_WARN", DEF_HIGH_TEMP_WARN);
    config->highTempCrit = getenv_double("CONTROLLER_HIGH_TEMP_CRIT", DEF_HIGH_TEMP_CRIT);

    /* Node identification */
    config->nodeId = getenv_dup("NODE_ID", "unknown");

    /* Validate required fields */
    if (!config->serialPort) {
        usys_log_error("config: CONTROLLER_SERIAL_PORT is required");
        return -1;
    }

    if (!config->driverName) {
        usys_log_error("config: CONTROLLER_DRIVER is required");
        return -1;
    }

    return 0;
}

void config_log(const Config *config) {
    if (!config) return;

    usys_log_info("Configuration:");
    usys_log_info("  Web service: %s:%d", config->listenAddr, config->listenPort);
    usys_log_info("  Sample interval: %d ms", config->sampleMs);
    usys_log_info("  Driver: %s", config->driverName);
    usys_log_info("  Serial port: %s @ %d baud", config->serialPort, config->baudRate);
    usys_log_info("  Notify.d: %s:%d%s (enabled=%s)",
                  config->notifyHost, config->notifyPort, config->notifyPath,
                  config->enableNotify ? "yes" : "no");
    usys_log_info("  Alarm thresholds:");
    usys_log_info("    Low voltage: warn=%.2fV crit=%.2fV",
                  config->lowVoltageWarn, config->lowVoltageCrit);
    usys_log_info("    High temp: warn=%.1f°C crit=%.1f°C",
                  config->highTempWarn, config->highTempCrit);
    usys_log_info("  Node ID: %s", config->nodeId);
}

void config_print_env_help(void) {
    usys_puts("\nEnvironment variables:");
    usys_puts("  CONTROLLER_LISTEN_ADDR    Listen address (default: 0.0.0.0)");
    usys_puts("  CONTROLLER_LISTEN_PORT    Listen port (default: 8095)");
    usys_puts("  CONTROLLER_SAMPLE_MS      Sample interval in ms (default: 1000)");
    usys_puts("  CONTROLLER_DRIVER         Driver name: victron (default: victron)");
    usys_puts("  CONTROLLER_SERIAL_PORT    Serial port path (default: /dev/ttyUSB0)");
    usys_puts("  CONTROLLER_BAUD_RATE      Serial baud rate (default: 19200)");
    usys_puts("  CONTROLLER_ENABLE_NOTIFY  Enable notify.d integration (default: true)");
    usys_puts("  CONTROLLER_LOW_VOLT_WARN  Low voltage warning threshold (default: 46.0V, 48V system)");
    usys_puts("  CONTROLLER_LOW_VOLT_CRIT  Low voltage critical threshold (default: 44.0V, 48V system)");
    usys_puts("  CONTROLLER_HIGH_TEMP_WARN High temperature warning threshold (default: 55°C)");
    usys_puts("  CONTROLLER_HIGH_TEMP_CRIT High temperature critical threshold (default: 65°C)");
    usys_puts("  NOTIFY_HOST               Notify.d host (default: 127.0.0.1)");
    usys_puts("  NODE_ID                   Node identifier");
}

void config_free(Config *config) {
    if (!config) return;

    if (config->listenAddr) {
        free(config->listenAddr);
        config->listenAddr = NULL;
    }
    if (config->driverName) {
        free(config->driverName);
        config->driverName = NULL;
    }
    if (config->serialPort) {
        free(config->serialPort);
        config->serialPort = NULL;
    }
    if (config->nodeId) {
        free(config->nodeId);
        config->nodeId = NULL;
    }
}

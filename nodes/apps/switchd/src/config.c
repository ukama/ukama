/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>

#include "config.h"
#include "switchd.h"

#include "usys_file.h"
#include "usys_api.h"
#include "usys_log.h"
#include "usys_services.h"
#include "usys_string.h"

static void env_str(const char *name,
                    char *dst,
                    size_t len,
                    const char *defVal) {
    const char *value;

    value = getenv(name);
    snprintf(dst, len, "%s", (value && *value) ? value : defVal);
}

static int env_int(const char *name, int defVal) {
    const char *value;

    value = getenv(name);
    return (value && *value) ? atoi(value) : defVal;
}

static bool env_bool(const char *name, bool defVal) {
    const char *value;

    value = getenv(name);
    if (value == NULL || *value == '\0') {
        return defVal;
    }

    return (strcmp(value, "1") == 0 ||
            strcasecmp(value, "true") == 0 ||
            strcasecmp(value, "yes") == 0 ||
            strcasecmp(value, "on") == 0);
}

void config_print_env_help(void) {
    usys_puts("Environment variables:");
    usys_puts("  SWITCHD_LOG_LEVEL              TRACE|DEBUG|INFO|WARNING|ERROR");
    usys_puts("  SWITCHD_DRIVER                 Driver name (default: tycon_snmp)");
    usys_puts("  SWITCHD_HTTP_HOST              Listen host");
    usys_puts("  SWITCHD_HTTP_PORT              Service port (defaults to usys_find_service_port)");
    usys_puts("  SWITCHD_URL_PREFIX             REST prefix (default: /v1)");
    usys_puts("  SWITCHD_SNMP_HOST              Switch management IP/host");
    usys_puts("  SWITCHD_SNMP_PORT              SNMP port");
    usys_puts("  SWITCHD_SNMP_COMMUNITY         SNMP v2c community");
    usys_puts("  SWITCHD_SNMP_VERSION           SNMP version");
    usys_puts("  SWITCHD_SNMP_TIMEOUT_MS        SNMP timeout in ms");
    usys_puts("  SWITCHD_SNMP_RETRIES           SNMP retries");
    usys_puts("  SWITCHD_POLL_STATUS_SEC        Port/status poll interval");
    usys_puts("  SWITCHD_POLL_KPIS_SEC          KPI poll interval");
    usys_puts("  SWITCHD_POLL_INFO_SEC          Switch info poll interval");
    usys_puts("  SWITCHD_ALARM_SCAN_SEC         Alarm scan interval");
    usys_puts("  SWITCHD_COMMAND_TIMEOUT_MS     Command timeout");
    usys_puts("  SWITCHD_FIRMWARE_RECONNECT_SEC Firmware reconnect timeout");
    usys_puts("  SWITCHD_FIRMWARE_VERIFY_SEC    Firmware verify window");
    usys_puts("  SWITCHD_POE_CYCLE_MS           PoE cycle duration in ms");
    usys_puts("  SWITCHD_NOTIFY_URL             notify.d full URL override");
    usys_puts("  SWITCHD_NOTIFY_HOST            notify.d host when URL override is absent");
    usys_puts("  SWITCHD_NOTIFY_PORT            notify.d port override");
    usys_puts("  SWITCHD_NOTIFY_TIMEOUT_MS      notify.d timeout in ms");
    usys_puts("  SWITCHD_TFTP_BIND_IP           TFTP bind IP");
    usys_puts("  SWITCHD_TFTP_PORT              TFTP port");
    usys_puts("  SWITCHD_TFTP_ROOT              TFTP root directory");
    usys_puts("  SWITCHD_STRICT_LINK_ALARMS     true|false");
    usys_puts("  SWITCHD_SAVE_AFTER_WRITE       true|false");
    usys_puts("  SWITCHD_POLICY_PATH            Port policy path");
}

void config_usage(void) {
    usys_puts("Usage: switch.d [options]");
    usys_puts("Options:");
    usys_puts("  -h, --help                    Show help");
    usys_puts("  -l, --logs <LEVEL>            Set log level");
    usys_puts("  -v, --version                 Print version");
    config_print_env_help();
}

int config_load(SwitchdConfig *cfg) {
    int servicePort;
    int notifyPort;
    char notifyHost[64];

    if (cfg == NULL) {
        return SWITCHD_ERR_INVAL;
    }

    memset(cfg, 0, sizeof(*cfg));

    env_str(ENV_SWITCHD_DRIVER,
            cfg->driverName,
            sizeof(cfg->driverName),
            DEF_DRIVER_NAME);
    env_str(ENV_SWITCHD_HTTP_HOST,
            cfg->httpHost,
            sizeof(cfg->httpHost),
            DEF_HTTP_HOST);

    servicePort = usys_find_service_port(SERVICE_SWITCH);
    if (servicePort <= 0) {
        servicePort = 10310;
    }
    cfg->httpPort = env_int(ENV_SWITCHD_HTTP_PORT, servicePort);

    env_str(ENV_SWITCHD_URL_PREFIX,
            cfg->urlPrefix,
            sizeof(cfg->urlPrefix),
            DEF_URL_PREFIX);
    env_str(ENV_SWITCHD_SNMP_HOST,
            cfg->snmpHost,
            sizeof(cfg->snmpHost),
            DEF_SNMP_HOST);
    cfg->snmpPort = env_int(ENV_SWITCHD_SNMP_PORT, 161);
    env_str(ENV_SWITCHD_SNMP_COMMUNITY,
            cfg->snmpCommunity,
            sizeof(cfg->snmpCommunity),
            DEF_SNMP_COMMUNITY);
    cfg->snmpVersion = env_int(ENV_SWITCHD_SNMP_VERSION, 2);
    cfg->snmpTimeoutMs = env_int(ENV_SWITCHD_SNMP_TIMEOUT_MS, 1500);
    cfg->snmpRetries = env_int(ENV_SWITCHD_SNMP_RETRIES, 1);

    cfg->pollStatusSec = env_int(ENV_SWITCHD_POLL_STATUS_SEC, 5);
    cfg->pollKpisSec = env_int(ENV_SWITCHD_POLL_KPIS_SEC, 15);
    cfg->pollInfoSec = env_int(ENV_SWITCHD_POLL_INFO_SEC, 60);
    cfg->alarmScanSec = env_int(ENV_SWITCHD_ALARM_SCAN_SEC, 5);
    cfg->commandTimeoutMs = env_int(ENV_SWITCHD_COMMAND_TIMEOUT_MS, 15000);
    cfg->firmwareReconnectSec = env_int(ENV_SWITCHD_FIRMWARE_RECONNECT_SEC,
                                        90);
    cfg->firmwareVerifySec = env_int(ENV_SWITCHD_FIRMWARE_VERIFY_SEC, 30);
    cfg->poeCycleMs = env_int(ENV_SWITCHD_POE_CYCLE_MS, 3000);

    env_str(ENV_SWITCHD_NOTIFY_HOST,
            notifyHost,
            sizeof(notifyHost),
            DEF_NOTIFY_HOST);
    notifyPort = env_int(ENV_SWITCHD_NOTIFY_PORT,
                         usys_find_service_port(SERVICE_NOTIFY));
    if (notifyPort <= 0) {
        notifyPort = 9090;
    }

    env_str(ENV_SWITCHD_NOTIFY_URL,
            cfg->notifyUrl,
            sizeof(cfg->notifyUrl),
            "");
    if (cfg->notifyUrl[0] == '\0') {
        snprintf(cfg->notifyUrl,
                 sizeof(cfg->notifyUrl),
                 "http://%s:%d%s",
                 notifyHost,
                 notifyPort,
                 DEF_NOTIFY_EP);
    }
    cfg->notifyTimeoutMs = env_int(ENV_SWITCHD_NOTIFY_TIMEOUT_MS, 2000);

    env_str(ENV_SWITCHD_TFTP_BIND_IP,
            cfg->tftpBindIp,
            sizeof(cfg->tftpBindIp),
            DEF_TFTP_BIND_IP);
    cfg->tftpPort = env_int(ENV_SWITCHD_TFTP_PORT, 1069);
    env_str(ENV_SWITCHD_TFTP_ROOT,
            cfg->tftpRoot,
            sizeof(cfg->tftpRoot),
            DEF_TFTP_ROOT);

    cfg->strictLinkAlarms = env_bool(ENV_SWITCHD_STRICT_LINK_ALARMS, false);
    cfg->saveAfterWrite = env_bool(ENV_SWITCHD_SAVE_AFTER_WRITE, true);
    env_str(ENV_SWITCHD_POLICY_PATH,
            cfg->policyPath,
            sizeof(cfg->policyPath),
            DEF_POLICY_PATH);

    usys_log_debug("Loaded %s configuration: http=%s:%d snmp=%s:%d driver=%s",
                   SERVICE_NAME,
                   cfg->httpHost,
                   cfg->httpPort,
                   cfg->snmpHost,
                   cfg->snmpPort,
                   cfg->driverName);

    return SWITCHD_OK;
}

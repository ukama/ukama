/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <stdio.h>
#include <stdlib.h>

#include "config.h"

#include "usys_log.h"
#include "usys_file.h"
#include "usys_services.h"

static int trim_line(char *s) {

    char *p;
    size_t n;

    if (!s) return STATUS_NOK;

    n = strlen(s);
    while (n > 0 && (s[n - 1] == '\n' || s[n - 1] == '\r' || s[n - 1] == ' ' || s[n - 1] == '\t')) {
        s[n - 1] = '\0';
        n--;
    }

    p = s;
    while (*p == ' ' || *p == '\t') p++;
    if (p != s) memmove(s, p, strlen(p) + 1);

    return STATUS_OK;
}

static int parse_kv(char *line, char **k, char **v) {

    char *eq;

    if (!line || !k || !v) return STATUS_NOK;

    eq = strchr(line, '=');
    if (!eq) return STATUS_NOK;

    *eq = '\0';
    *k = line;
    *v = eq + 1;

    (void)trim_line(*k);
    (void)trim_line(*v);

    if (**k == '\0') return STATUS_NOK;
    return STATUS_OK;
}

static int parse_bool(const char *s, bool *out) {

    if (!s || !out) return STATUS_NOK;

    if (!strcmp(s, "1") || !strcasecmp(s, "true") || !strcasecmp(s, "yes") || !strcasecmp(s, "on")) {
        *out = true;
        return STATUS_OK;
    }
    if (!strcmp(s, "0") || !strcasecmp(s, "false") || !strcasecmp(s, "no") || !strcasecmp(s, "off")) {
        *out = false;
        return STATUS_OK;
    }
    return STATUS_NOK;
}

int config_set_defaults(Config *cfg) {

    if (!cfg) return STATUS_NOK;

    memset(cfg, 0, sizeof(*cfg));

    snprintf(cfg->serviceName, sizeof(cfg->serviceName), "%s", SERVICE_NAME);

    cfg->servicePort = usys_find_service_port(SERVICE_NAME);
    if (!cfg->servicePort) {
        usys_log_error("Unable to find service port for: %s", SERVICE_NAME);
        return STATUS_NOK;
    }

    cfg->nodedPort = usys_find_service_port(SERVICE_NODE);
    if (!cfg->nodedPort) {
        usys_log_error("Unable to find service port for: %s", SERVICE_NODE);
        return STATUS_NOK;
    }

    snprintf(cfg->gpioBasePath, sizeof(cfg->gpioBasePath), "%s", "/sys/devices/platform");

    cfg->i2cBusFem1 = 1;
    cfg->i2cBusFem2 = 2;
    cfg->i2cBusCtrl = 0;

    snprintf(cfg->safetyConfigPath, sizeof(cfg->safetyConfigPath), "%s", "/etc/femd/safety.yaml");

    snprintf(cfg->notifyHost, sizeof(cfg->notifyHost), "%s", "127.0.0.1");
    cfg->notifyPort = 8082;
    snprintf(cfg->notifyPath, sizeof(cfg->notifyPath), "%s", "/v1/notify");

    cfg->samplePeriodMs = 1000;
    cfg->safetyPeriodMs = 500;

    cfg->enableWeb    = true;
    cfg->enableSafety = true;
    cfg->enableNotify = true;

    return STATUS_OK;
}

int config_load(Config *cfg, const char *path) {

    FILE *f;
    char line[512];

    if (!cfg || !path) return STATUS_NOK;

    f = fopen(path, "r");
    if (!f) {
        usys_log_warn("config file not found: %s (using defaults)", path);
        return STATUS_OK;
    }

    while (fgets(line, sizeof(line), f)) {
        char *k, *v;

        (void)trim_line(line);
        if (line[0] == '\0' || line[0] == '#') continue;

        if (parse_kv(line, &k, &v) != STATUS_OK) continue;

        if (!strcmp(k, "service_port"))        cfg->servicePort = atoi(v);
        else if (!strcmp(k, "noded_port"))     cfg->nodedPort = atoi(v);
        else if (!strcmp(k, "service_name"))   snprintf(cfg->serviceName,
                                                        sizeof(cfg->serviceName), "%s", v);
        else if (!strcmp(k, "gpio_base_path")) snprintf(cfg->gpioBasePath,
                                                        sizeof(cfg->gpioBasePath), "%s", v);
        else if (!strcmp(k, "i2c_bus_fem1"))   cfg->i2cBusFem1 = atoi(v);
        else if (!strcmp(k, "i2c_bus_fem2"))   cfg->i2cBusFem2 = atoi(v);
        else if (!strcmp(k, "i2c_bus_ctrl"))   cfg->i2cBusCtrl = atoi(v);
        else if (!strcmp(k, "safety_yaml"))    snprintf(cfg->safetyConfigPath,
                                                        sizeof(cfg->safetyConfigPath), "%s", v);
        else if (!strcmp(k, "notify_host"))    snprintf(cfg->notifyHost,
                                                        sizeof(cfg->notifyHost), "%s", v);
        else if (!strcmp(k, "notify_port"))    cfg->notifyPort = atoi(v);
        else if (!strcmp(k, "notify_path"))    snprintf(cfg->notifyPath,
                                                        sizeof(cfg->notifyPath), "%s", v);
        else if (!strcmp(k, "sample_period_ms")) cfg->samplePeriodMs = (uint32_t)atoi(v);
        else if (!strcmp(k, "safety_period_ms")) cfg->safetyPeriodMs = (uint32_t)atoi(v);
        else if (!strcmp(k, "enable_web"))       (void)parse_bool(v, &cfg->enableWeb);
        else if (!strcmp(k, "enable_safety"))    (void)parse_bool(v, &cfg->enableSafety);
        else if (!strcmp(k, "enable_notify"))    (void)parse_bool(v, &cfg->enableNotify);
    }

    fclose(f);
    return STATUS_OK;
}

int config_validate(const Config *cfg) {

    if (!cfg) return STATUS_NOK;

    if (cfg->servicePort <= 0 || cfg->servicePort > 65535) return STATUS_NOK;
    if (cfg->nodedPort <= 0 || cfg->nodedPort > 65535) return STATUS_NOK;

    if (cfg->i2cBusFem1 < 0 || cfg->i2cBusFem2 < 0 || cfg->i2cBusCtrl < 0) return STATUS_NOK;

    if (cfg->samplePeriodMs == 0) return STATUS_NOK;
    if (cfg->safetyPeriodMs == 0) return STATUS_NOK;

    return STATUS_OK;
}

void config_print(const Config *cfg) {

    if (!cfg) return;

    usys_log_info("Config: servicePort=%d nodedPort=%d serviceName=%s",
                  cfg->servicePort, cfg->nodedPort, cfg->serviceName);

    usys_log_info("Config: gpioBasePath=%s i2cFem1=%d i2cFem2=%d i2cCtrl=%d",
                  cfg->gpioBasePath, cfg->i2cBusFem1, cfg->i2cBusFem2, cfg->i2cBusCtrl);

    usys_log_info("Config: safetyYaml=%s samplePeriodMs=%u safetyPeriodMs=%u",
                  cfg->safetyConfigPath, (unsigned)cfg->samplePeriodMs, (unsigned)cfg->safetyPeriodMs);

    usys_log_info("Config: notify=%s:%d%s web=%s safety=%s notify=%s",
                  cfg->notifyHost, cfg->notifyPort, cfg->notifyPath,
                  cfg->enableWeb ? "true" : "false",
                  cfg->enableSafety ? "true" : "false",
                  cfg->enableNotify ? "true" : "false");
}

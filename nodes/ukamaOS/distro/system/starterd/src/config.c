/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "config.h"
#include "starterd.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>

#include "usys_log.h"
#include "usys_file.h"
#include "usys_services.h"

static int cfg_get_int(const char *name, int defVal) {

    const char *v;
    char *end;
    long n;

    v = getenv(name);
    if (!v || !*v) return defVal;

    end = NULL;
    n = strtol(v, &end, 10);
    if (end == v || *end != '\0') return defVal;

    if (n < -2147483648L) n = -2147483648L;
    if (n >  2147483647L) n =  2147483647L;

    return (int)n;
}

static char* cfg_get_str(const char *name, const char *defVal) {

    const char *v;

    v = getenv(name);
    if (v && *v) {
        return strdup(v);
    }

    if (defVal) {
        return strdup(defVal);
    }

    return NULL;
}

static void cfg_trim(char *s) {

    char *p;
    char *e;

    if (!s) return;

    p = s;
    while (*p && isspace((unsigned char)*p)) p++;
    if (p != s) memmove(s, p, strlen(p) + 1);

    e = s + strlen(s);
    while (e > s && isspace((unsigned char)e[-1])) e--;
    *e = '\0';
}

static void cfg_apply_kv(const char *key, const char *val) {

    if (!key || !*key) return;
    if (!val) val = "";

    if (setenv(key, val, 0) != 0) {
        /* ignore */
    }
}

static void cfg_load_env_file_if_enabled(void) {

    const char *path;
    FILE *f;
    char line[1024];
    char *eq;
    char *key;
    char *val;

    path = getenv("STARTERD_CONFIG");
    if (!path || !*path) return;

    f = fopen(path, "r");
    if (!f) {
        usys_log_error("config: unable to open STARTERD_CONFIG=%s", path);
        return;
    }

    while (fgets(line, sizeof(line), f) != NULL) {

        cfg_trim(line);
        if (line[0] == '\0') continue;
        if (line[0] == '#') continue;

        eq = strchr(line, '=');
        if (!eq) continue;

        *eq = '\0';
        key = line;
        val = eq + 1;

        cfg_trim(key);
        cfg_trim(val);

        cfg_apply_kv(key, val);
    }

    fclose(f);
}

static bool cfg_validate(Config *config) {

    if (!config->manifestPath || !*config->manifestPath) {
        usys_log_error("config: manifest path missing");
        return false;
    }

    if (config->httpPort <= 0 || config->httpPort > 65535) {
        usys_log_error("config: invalid http port %d", config->httpPort);
        return false;
    }

    if (config->wimcPort <= 0 || config->wimcPort > 65535) {
        usys_log_error("config: invalid wimc port %d", config->wimcPort);
        return false;
    }

    if (config->commitTimeoutSec <= 0)      config->commitTimeoutSec = 20;
    if (config->pingTimeoutSec <= 0)        config->pingTimeoutSec = 3;
    if (config->termGraceSec <= 0)          config->termGraceSec = 5;
    if (config->restartMaxBackoffSec <= 0)  config->restartMaxBackoffSec = 60;
    if (config->restartStableResetSec <= 0) config->restartStableResetSec = 300;

    return true;
}

bool config_load(Config *config) {

    if (!config) {
        return false;
    }

    memset(config, 0, sizeof(*config));

    cfg_load_env_file_if_enabled();

    config->manifestPath = cfg_get_str("STARTERD_MANIFEST",   STARTERD_DEFAULT_MANIFEST_FILE);
    config->logPath      = cfg_get_str("STARTERD_LOG_PATH",   STARTERD_DEFAULT_LOG_PATH);
    config->readyFile    = cfg_get_str("STARTERD_READY_FILE", STARTERD_DEFAULT_READY_FILE);

    config->appsRoot     = cfg_get_str("STARTERD_APPS_ROOT", "/ukama/apps");
    config->pkgsDir      = cfg_get_str("STARTERD_PKGS_DIR",  "/ukama/apps/pkgs");
    config->stateDir     = cfg_get_str("STARTERD_STATE_DIR", "/ukama/state/starterd");

    config->httpAddr     = cfg_get_str("STARTERD_HTTP_ADDR", "0.0.0.0");

    /* starter.d port from service registry */
    config->httpPort = usys_find_service_port(SERVICE_STARTER);
    if (config->httpPort <= 0) {
        usys_log_error("SERVICE_STARTER port not found in service registry");
        return false;
    }

    config->wimcHost = cfg_get_str("STARTERD_WIMC_HOST", "0.0.0.0");
    config->wimcPort = usys_find_service_port(SERVICE_WIMC);
    if (config->wimcPort <= 0) {
        usys_log_error("SERVICE_WIMC port not found in service registry");
        return false;
    }

    config->wimcPathTemplate =
        cfg_get_str("STARTERD_WIMC_PATH_TEMPLATE", "/v1/apps/%s/%s");

    config->commitTimeoutSec = cfg_get_int("STARTERD_COMMIT_TIMEOUT_SEC", 20);
    config->pingTimeoutSec   = cfg_get_int("STARTERD_PING_TIMEOUT_SEC",   3);
    config->termGraceSec     = cfg_get_int("STARTERD_TERM_GRACE_SEC",     5);

    config->restartMaxBackoffSec  = cfg_get_int("STARTERD_RESTART_MAX_BACKOFF_SEC",  60);
    config->restartStableResetSec = cfg_get_int("STARTERD_RESTART_STABLE_RESET_SEC", 300);

    config->bootSpace = cfg_get_str("STARTERD_BOOT_SPACE", "boot");

    return cfg_validate(config);
}

void config_free(Config *config) {

    if (!config) return;

    free(config->manifestPath);
    free(config->logPath);
    free(config->readyFile);
    free(config->appsRoot);
    free(config->pkgsDir);
    free(config->stateDir);
    free(config->httpAddr);
    free(config->wimcHost);
    free(config->wimcPathTemplate);
    free(config->bootSpace);

    memset(config, 0, sizeof(*config));
}

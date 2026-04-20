/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <ctype.h>

#include "powerd.h"
#include "config.h"
#include "usys_log.h"
#include "usys_mem.h"

static char *dup_env(const char *name, const char *defVal) {

    char *v = getenv(name);

    if (!v || !*v) {
        if (!defVal) return NULL;
        return strdup(defVal);
    }

    return strdup(v);
}

static int parse_u32(const char *s, uint32_t *out) {

    char *end = NULL;
    unsigned long v;

    if (!s || !*s || !out) return -1;

    v = strtoul(s, &end, 0);
    if (!end || *end != '\0') return -1;
    if (v > 0xFFFFFFFFUL) return -1;

    *out = (uint32_t)v;
    return 0;
}

static int parse_i32(const char *s, int *out) {

    char *end = NULL;
    long v;

    if (!s || !*s || !out) return -1;

    v = strtol(s, &end, 0);
    if (!end || *end != '\0') return -1;
    if (v < -2147483648L || v > 2147483647L) return -1;

    *out = (int)v;
    return 0;
}

static int parse_bool(const char *s, int *out) {

    if (!s || !*s || !out) return -1;

    if (!strcasecmp(s, "1") ||
        !strcasecmp(s, "true") ||
        !strcasecmp(s, "yes") ||
        !strcasecmp(s, "on")) {
        *out = 1;
        return 0;
    }

    if (!strcasecmp(s, "0") ||
        !strcasecmp(s, "false") ||
        !strcasecmp(s, "no") ||
        !strcasecmp(s, "off")) {
        *out = 0;
        return 0;
    }

    return -1;
}

static int parse_addr7(const char *s, int *out7) {

    int v;

    if (parse_i32(s, &v) != 0) return -1;
    if (v < 0x03 || v > 0x77) return -1;

    *out7 = v;
    return 0;
}

static void init_ads_map(Config *cfg) {

    cfg->adsChVin = -1;
    cfg->adsChVpa = -1;
    cfg->adsChAux = -1;
}

static void parse_ads_chmap(Config *cfg, const char *s) {

    const char *p = s;

    if (!cfg || !p || !*p) return;

    while (*p) {
        char key[16];
        char val[8];
        int ki = 0;
        int vi = 0;
        int ch = -1;

        while (*p && isspace((unsigned char)*p)) p++;

        while (*p && *p != '=' && *p != ',' && ki < (int)sizeof(key) - 1) {
            key[ki++] = (char)tolower((unsigned char)*p);
            p++;
        }
        key[ki] = '\0';

        if (*p != '=') {
            while (*p && *p != ',') p++;
            if (*p == ',') p++;
            continue;
        }

        p++;

        while (*p && *p != ',' && vi < (int)sizeof(val) - 1) {
            val[vi++] = *p;
            p++;
        }
        val[vi] = '\0';

        if (parse_i32(val, &ch) == 0 && ch >= 0 && ch <= 3) {
            if (!strcmp(key, "vin")) cfg->adsChVin = ch;
            else if (!strcmp(key, "vpa")) cfg->adsChVpa = ch;
            else if (!strcmp(key, "aux")) cfg->adsChAux = ch;
        }

        if (*p == ',') p++;
    }
}

int config_validate_env(Config *cfg) {

    if (!cfg) return -1;

    if (!cfg->listenAddr || !*cfg->listenAddr) {
        usys_log_error("config: invalid listen address");
        return -1;
    }

    if (cfg->sampleMs < 100) {
        usys_log_error("config: sample period too low: %u", cfg->sampleMs);
        return -1;
    }

    if (!cfg->mockMode) {
        if (cfg->lm75Dev && cfg->lm75Addr == 0) {
            usys_log_error("config: POWER_LM75_DEV set but POWER_LM75_ADDR missing");
            return -1;
        }

        if (cfg->lm25066Dev && cfg->lm25066Addr == 0) {
            usys_log_error("config: POWER_LM25066_DEV set but POWER_LM25066_ADDR missing");
            return -1;
        }

        if (cfg->ads1015Dev && cfg->ads1015Addr == 0) {
            usys_log_error("config: POWER_ADS1015_DEV set but POWER_ADS1015_ADDR missing");
            return -1;
        }
    }

    return 0;
}

void config_log(Config *cfg) {

    if (!cfg) return;

    usys_log_info("config: listen=%s:%u", cfg->listenAddr, cfg->listenPort);
    usys_log_info("config: board=%s sampleMs=%u mock=%d",
                  cfg->boardName ? cfg->boardName : "unknown",
                  cfg->sampleMs,
                  cfg->mockMode);

    usys_log_info("config: lm75 dev=%s addr=0x%02x",
                  cfg->lm75Dev ? cfg->lm75Dev : "(disabled)",
                  cfg->lm75Addr);

    usys_log_info("config: lm25066 dev=%s addr=0x%02x clHigh=%d rsMohm=%d",
                  cfg->lm25066Dev ? cfg->lm25066Dev : "(disabled)",
                  cfg->lm25066Addr,
                  cfg->lm25066ClHigh,
                  cfg->lm25066RsMohm);

    usys_log_info("config: ads1015 dev=%s addr=0x%02x chmap vin=%d vpa=%d aux=%d",
                  cfg->ads1015Dev ? cfg->ads1015Dev : "(disabled)",
                  cfg->ads1015Addr,
                  cfg->adsChVin,
                  cfg->adsChVpa,
                  cfg->adsChAux);
}

int config_load_from_env(Config *cfg) {

    char *s = NULL;

    if (!cfg) return -1;

    memset(cfg, 0, sizeof(*cfg));

    cfg->listenAddr = dup_env("POWER_LISTEN_ADDR", "0.0.0.0");
    cfg->boardName  = dup_env("POWER_BOARD", "unknown");

    cfg->listenPort = usys_find_service_port(SERVICE_NAME);
    if (!cfg->listenPort) {
        usys_log_error("config: listening port for %s not found", SERVICE_NAME);
        return -1;
    }

    cfg->sampleMs = 1000;
    s = getenv("POWER_SAMPLE_MS");
    if (s && *s) {
        if (parse_u32(s, &cfg->sampleMs) != 0) {
            usys_log_error("Invalid POWER_SAMPLE_MS: %s", s);
            return -1;
        }
    }

    cfg->mockMode = 0;
    s = getenv("POWER_MOCK");
    if (s && *s) {
        if (parse_bool(s, &cfg->mockMode) != 0) {
            usys_log_error("Invalid POWER_MOCK: %s", s);
            return -1;
        }
    }

    cfg->lm75Dev    = dup_env("POWER_LM75_DEV", NULL);
    cfg->lm25066Dev = dup_env("POWER_LM25066_DEV", NULL);
    cfg->ads1015Dev = dup_env("POWER_ADS1015_DEV", NULL);

    cfg->lm75Addr = 0;
    cfg->lm25066Addr = 0;
    cfg->ads1015Addr = 0;

    cfg->lm25066ClHigh = 0;
    cfg->lm25066RsMohm = 1;

    init_ads_map(cfg);

    s = getenv("POWER_LM75_ADDR");
    if (s && *s) {
        if (parse_addr7(s, &cfg->lm75Addr) != 0) {
            usys_log_error("Invalid POWER_LM75_ADDR: %s", s);
            return -1;
        }
    }

    s = getenv("POWER_LM25066_ADDR");
    if (s && *s) {
        if (parse_addr7(s, &cfg->lm25066Addr) != 0) {
            usys_log_error("Invalid POWER_LM25066_ADDR: %s", s);
            return -1;
        }
    }

    s = getenv("POWER_LM25066_CL_HIGH");
    if (s && *s) {
        if (parse_i32(s, &cfg->lm25066ClHigh) != 0 ||
            (cfg->lm25066ClHigh != 0 && cfg->lm25066ClHigh != 1)) {
            usys_log_error("Invalid POWER_LM25066_CL_HIGH: %s", s);
            return -1;
        }
    }

    s = getenv("POWER_LM25066_RS_MOHM");
    if (s && *s) {
        if (parse_i32(s, &cfg->lm25066RsMohm) != 0 ||
            cfg->lm25066RsMohm < 1 || cfg->lm25066RsMohm > 1000) {
            usys_log_error("Invalid POWER_LM25066_RS_MOHM: %s", s);
            return -1;
        }
    }

    s = getenv("POWER_ADS1015_ADDR");
    if (s && *s) {
        if (parse_addr7(s, &cfg->ads1015Addr) != 0) {
            usys_log_error("Invalid POWER_ADS1015_ADDR: %s", s);
            return -1;
        }
    }

    s = getenv("POWER_ADS1015_CHMAP");
    if (s && *s) {
        parse_ads_chmap(cfg, s);
    }

    if (config_validate_env(cfg) != 0) return -1;

    return 0;
}

void config_print_env_help(void) {

    usys_puts("Environment:");
    usys_puts("  POWER_LISTEN_ADDR=0.0.0.0");
    usys_puts("  POWER_SAMPLE_MS=1000");
    usys_puts("  POWER_BOARD=tower|amp|unknown");
    usys_puts("  POWER_MOCK=0|1");
    usys_puts("");
    usys_puts("  POWER_LM75_DEV=/dev/i2c-1");
    usys_puts("  POWER_LM75_ADDR=0x48");
    usys_puts("");
    usys_puts("  POWER_LM25066_DEV=/dev/i2c-1");
    usys_puts("  POWER_LM25066_ADDR=0x40");
    usys_puts("  POWER_LM25066_CL_HIGH=0|1");
    usys_puts("  POWER_LM25066_RS_MOHM=1");
    usys_puts("");
    usys_puts("  POWER_ADS1015_DEV=/dev/i2c-1");
    usys_puts("  POWER_ADS1015_ADDR=0x48");
    usys_puts("  POWER_ADS1015_CHMAP=vin=0,vpa=1,aux=2");
}

void config_free(Config *cfg) {

    if (!cfg) return;

    usys_free(cfg->listenAddr);
    usys_free(cfg->boardName);
    usys_free(cfg->lm75Dev);
    usys_free(cfg->lm25066Dev);
    usys_free(cfg->ads1015Dev);

    memset(cfg, 0, sizeof(*cfg));
}

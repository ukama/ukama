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

static int parse_u16(const char *s, uint16_t *out) {

	char *end = NULL;
	long v;

	if (!s || !*s) return -1;
	v = strtol(s, &end, 0);
	if (!end || *end) return -1;
	if (v < 0 || v > 65535) return -1;
	*out = (uint16_t)v;
	return 0;
}

static int parse_u32(const char *s, uint32_t *out) {

	char *end = NULL;
	unsigned long v;

	if (!s || !*s) return -1;
	v = strtoul(s, &end, 0);
	if (!end || *end) return -1;
	if (v > 0xFFFFFFFFUL) return -1;
	*out = (uint32_t)v;
	return 0;
}

static int parse_i32(const char *s, int *out) {

	char *end = NULL;
	long v;

	if (!s || !*s) return -1;
	v = strtol(s, &end, 0);
	if (!end || *end) return -1;
	if (v < -2147483648L || v > 2147483647L) return -1;
	*out = (int)v;
	return 0;
}

static int parse_hex_addr(const char *s, int *out7) {

	int v;

	if (!s || !*s) return -1;
	if (parse_i32(s, &v) != 0) return -1;
	/* accept 7-bit address */
	if (v < 0x03 || v > 0x77) return -1;
	*out7 = v;
	return 0;
}

static void init_ads_map(Config *cfg) {

	cfg->adsChVin = -1;
	cfg->adsChVpa = -1;
	cfg->adsChAux = -1;
}

/*
 * POWER_ADS1015_CHMAP example: "vin=0,vpa=1,aux=2"
 * Unknown keys ignored. Values must be 0..3.
 */
static void parse_ads_chmap(Config *cfg, const char *s) {

	const char *p = s;

	if (!p || !*p) return;

	while (*p) {
		char key[16];
		char val[8];
		int ki = 0, vi = 0;
		int ch = -1;

		while (*p && isspace((unsigned char)*p)) p++;

		while (*p && *p != '=' && *p != ',' && ki < (int)sizeof(key)-1) {
			key[ki++] = (char)tolower((unsigned char)*p);
			p++;
		}
		key[ki] = 0;

		if (*p != '=') {
			while (*p && *p != ',') p++;
			if (*p == ',') p++;
			continue;
		}
		p++; /* '=' */

		while (*p && *p != ',' && vi < (int)sizeof(val)-1) {
			val[vi++] = *p;
			p++;
		}
		val[vi] = 0;

		if (parse_i32(val, &ch) == 0 && ch >= 0 && ch <= 3) {
			if (!strcmp(key, "vin")) cfg->adsChVin = ch;
			else if (!strcmp(key, "vpa")) cfg->adsChVpa = ch;
			else if (!strcmp(key, "aux")) cfg->adsChAux = ch;
		}

		if (*p == ',') p++;
	}
}

int config_load_from_env(Config *cfg) {

	char *s;

	memset(cfg, 0, sizeof(*cfg));

	cfg->listenAddr = dup_env("POWER_LISTEN_ADDR", "0.0.0.0");
	cfg->boardName = dup_env("POWER_BOARD", "unknown");

	cfg->listenPort = 8089;
	s = getenv("POWER_LISTEN_PORT");
	if (s && *s) {
		if (parse_u16(s, &cfg->listenPort) != 0) {
			usys_log_error("Invalid POWER_LISTEN_PORT: %s", s);
			return -1;
		}
	}

	cfg->sampleMs = 1000;
	s = getenv("POWER_SAMPLE_MS");
	if (s && *s) {
		if (parse_u32(s, &cfg->sampleMs) != 0 || cfg->sampleMs < 100) {
			usys_log_error("Invalid POWER_SAMPLE_MS: %s", s);
			return -1;
		}
	}

	/* optional devices */
	cfg->lm25066Dev = dup_env("POWER_LM25066_DEV", NULL);
	cfg->lm75Dev    = dup_env("POWER_LM75_DEV", NULL);
	cfg->ads1015Dev = dup_env("POWER_ADS1015_DEV", NULL);

	cfg->lm25066Addr = 0;
	cfg->lm75Addr    = 0;
	cfg->ads1015Addr = 0;

	cfg->lm25066ClHigh = 0;
	cfg->lm25066RsMohm = 0;

	init_ads_map(cfg);

	s = getenv("POWER_LM25066_ADDR");
	if (cfg->lm25066Dev && s && *s) {
		if (parse_hex_addr(s, &cfg->lm25066Addr) != 0) {
			usys_log_error("Invalid POWER_LM25066_ADDR: %s", s);
			return -1;
		}
	}

	s = getenv("POWER_LM25066_CL_HIGH");
	if (cfg->lm25066Dev && s && *s) {
		if (parse_i32(s, &cfg->lm25066ClHigh) != 0 ||
            (cfg->lm25066ClHigh != 0 &&
             cfg->lm25066ClHigh != 1)) {
			usys_log_error("Invalid POWER_LM25066_CL_HIGH: %s", s);
			return -1;
		}
	}

	s = getenv("POWER_LM25066_RS_MOHM");
	if (cfg->lm25066Dev && s && *s) {
		if (parse_i32(s, &cfg->lm25066RsMohm) != 0 ||
            cfg->lm25066RsMohm < 1 ||
            cfg->lm25066RsMohm > 100) {
			usys_log_error("Invalid POWER_LM25066_RS_MOHM: %s", s);
			return -1;
		}
	}

	s = getenv("POWER_LM75_ADDR");
	if (cfg->lm75Dev && s && *s) {
		if (parse_hex_addr(s, &cfg->lm75Addr) != 0) {
			usys_log_error("Invalid POWER_LM75_ADDR: %s", s);
			return -1;
		}
	}

	s = getenv("POWER_ADS1015_ADDR");
	if (cfg->ads1015Dev && s && *s) {
		if (parse_hex_addr(s, &cfg->ads1015Addr) != 0) {
			usys_log_error("Invalid POWER_ADS1015_ADDR: %s", s);
			return -1;
		}
	}

	s = getenv("POWER_ADS1015_CHMAP");
	if (cfg->ads1015Dev && s && *s) {
		parse_ads_chmap(cfg, s);
	}

	return 0;
}

void config_free(Config *cfg) {

	if (!cfg) return;

	usys_free(cfg->listenAddr);
	usys_free(cfg->boardName);

	usys_free(cfg->lm25066Dev);
	usys_free(cfg->lm75Dev);
	usys_free(cfg->ads1015Dev);

	memset(cfg, 0, sizeof(*cfg));
}

void config_print_env_help(void) {

}

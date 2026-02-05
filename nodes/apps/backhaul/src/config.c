/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>

#include "config.h"
#include "backhauld.h"

#include "usys_log.h"
#include "usys_string.h"
#include "usys_mem.h"
#include "usys_file.h"
#include "usys_services.h"

static int env_to_int(const char *name, int def) {

	const char *v = getenv(name);
	if (!v || !*v) return def;

	return atoi(v);
}

static char* env_to_strdup(const char *name, const char *def) {

	const char *v = getenv(name);
	if (!v || !*v) v = def;

	return v ? strdup(v) : NULL;
}

static int is_strict_env_mode(void) {
    const char *v = getenv("BACKHAULD_STRICT_ENV");
    if (!v || !*v) return USYS_FALSE;
    return atoi(v) ? USYS_TRUE : USYS_FALSE;
}

static void log_missing_env(const char *name, int strict, int *missingCount) {
    if (strict) {
        usys_log_error("Missing required ENV: %s", name);
    } else {
        usys_log_warn("Missing required ENV: %s (using default)", name);
    }
    (*missingCount)++;
}

static int env_is_missing(const char *name) {
    const char *v = getenv(name);
    return (!v || !*v);
}

int config_validate_env(Config *config) {

    int missing = 0;
    int strict = is_strict_env_mode();

    if (!config) return USYS_FALSE;

    /*
     * Required for “real” deployments:
     * - If you want to keep bootstrap defaults, remove these from required list.
     * - If reflector endpoints are fetched from bootstrap, bootstrap host/scheme/ep matter.
     */
    if (env_is_missing("BACKHAULD_BOOTSTRAP_HOST")) {
        log_missing_env("BACKHAULD_BOOTSTRAP_HOST", strict, &missing);
    }
    if (env_is_missing("BACKHAULD_BOOTSTRAP_SCHEME")) {
        log_missing_env("BACKHAULD_BOOTSTRAP_SCHEME", strict, &missing);
    }
    if (env_is_missing("BACKHAULD_BOOTSTRAP_EP")) {
        log_missing_env("BACKHAULD_BOOTSTRAP_EP", strict, &missing);
    }

    /*
     * Hard correctness constraints (even if ENV exists).
     * These prevent division-by-zero / nonsense scheduling.
     */
    if (config->microPeriodMs <= 0) {
        usys_log_error("Invalid BACKHAULD_MICRO_PERIOD_MS=%d", config->microPeriodMs);
        missing++;
    }
    if (config->multiPeriodMs <= 0) {
        usys_log_error("Invalid BACKHAULD_MULTI_PERIOD_MS=%d", config->multiPeriodMs);
        missing++;
    }
    if (config->connectTimeoutMs <= 0 || config->totalTimeoutMs <= 0) {
        usys_log_error("Invalid timeouts connect=%d total=%d",
                       config->connectTimeoutMs, config->totalTimeoutMs);
        missing++;
    }
    if (config->totalTimeoutMs < config->connectTimeoutMs) {
        usys_log_error("Invalid timeouts: totalTimeoutMs (%d) < connectTimeoutMs (%d)",
                       config->totalTimeoutMs, config->connectTimeoutMs);
        missing++;
    }
    if (config->pingBytes < 64) {
        usys_log_warn("pingBytes=%d too small, forcing to 64", config->pingBytes);
        config->pingBytes = 64;
    }

    /*
     * Reflector URLs: if user explicitly sets near/far URLs, we accept them.
     * If both are empty, that means “bootstrap discovery” must work.
     */
    if ((!config->reflectorNearUrl || !*config->reflectorNearUrl) &&
        (!config->reflectorFarUrl  || !*config->reflectorFarUrl)) {

        /* If no explicit reflectors, bootstrap becomes effectively required */
        if (!config->bootstrapHost || !*config->bootstrapHost ||
            !config->bootstrapScheme || !*config->bootstrapScheme ||
            !config->bootstrapEp || !*config->bootstrapEp) {
            usys_log_error("No reflector URLs set and bootstrap config missing/empty");
            missing++;
        }
    }

    if (missing && strict) {
        usys_log_error("Configuration invalid: %d issue(s)", missing);
        return USYS_FALSE;
    }

    if (missing) {
        usys_log_warn("Configuration has %d issue(s), continuing due to non-strict mode", missing);
    }

    return USYS_TRUE;
}

void config_log(Config *config) {

    if (!config) return;

    usys_log_info("backhaul.d effective config:");
    usys_log_info("  serviceName               = %s", config->serviceName ? config->serviceName : "");
    usys_log_info("  servicePort               = %d", config->servicePort);

    usys_log_info("  bootstrapScheme           = %s", config->bootstrapScheme ? config->bootstrapScheme : "");
    usys_log_info("  bootstrapHost             = %s", config->bootstrapHost ? config->bootstrapHost : "");
    usys_log_info("  bootstrapEp               = %s", config->bootstrapEp ? config->bootstrapEp : "");

    usys_log_info("  reflectorNearUrl          = %s", config->reflectorNearUrl ? config->reflectorNearUrl : "");
    usys_log_info("  reflectorFarUrl           = %s", config->reflectorFarUrl ? config->reflectorFarUrl : "");
    usys_log_info("  reflectorRefreshSec       = %d", config->reflectorRefreshSec);

    usys_log_info("  microPeriodMs             = %d", config->microPeriodMs);
    usys_log_info("  multiPeriodMs             = %d", config->multiPeriodMs);
    usys_log_info("  chgPeriodSec              = %d", config->chgPeriodSec);
    usys_log_info("  classifyPeriodSec         = %d", config->classifyPeriodSec);

    usys_log_info("  connectTimeoutMs          = %d", config->connectTimeoutMs);
    usys_log_info("  totalTimeoutMs            = %d", config->totalTimeoutMs);
    usys_log_info("  maxRetries                = %d", config->maxRetries);

    usys_log_info("  pingBytes                 = %d", config->pingBytes);
    usys_log_info("  stallThresholdMs          = %d", config->stallThresholdMs);

    usys_log_info("  chgTargetSec              = %d", config->chgTargetSec);
    usys_log_info("  chgWarmupBytes            = %d", config->chgWarmupBytes);
    usys_log_info("  chgMinBytes               = %d", config->chgMinBytes);
    usys_log_info("  chgMaxBytes               = %d", config->chgMaxBytes);
    usys_log_info("  chgSamples                = %d", config->chgSamples);

    usys_log_info("  parallelStreams           = %d", config->parallelStreams);
    usys_log_info("  parallelMaxBytesTotal     = %d", config->parallelMaxBytesTotal);

    usys_log_info("  downConsecFails           = %d", config->downConsecFails);
    usys_log_info("  recoverConsecOk           = %d", config->recoverConsecOk);
    usys_log_info("  capStabilityPct           = %d", config->capStabilityPct);

    usys_log_info("  windowMicroSamples        = %d", config->windowMicroSamples);
    usys_log_info("  windowMultiSamples        = %d", config->windowMultiSamples);
    usys_log_info("  windowChgSamples          = %d", config->windowChgSamples);

    usys_log_info("  strictEnv                 = %d", is_strict_env_mode());
}

int config_load_from_env(Config *config) {

	if (!config) return USYS_FALSE;
	memset(config, 0, sizeof(*config));

	config->serviceName = usys_strdup(SERVICE_NAME);
	config->servicePort = usys_find_service_port(SERVICE_NAME);
	if (!config->servicePort) {
		/* fallback if service registry missing */
		config->servicePort = env_to_int("BACKHAULD_PORT", 9050);
	}

	config->bootstrapHost	= env_to_strdup("BACKHAULD_BOOTSTRAP_HOST", "bootstrap.ukama.com");
	config->bootstrapScheme	= env_to_strdup("BACKHAULD_BOOTSTRAP_SCHEME", "https");
	config->bootstrapEp		= env_to_strdup("BACKHAULD_BOOTSTRAP_EP", "/reflector");

	config->reflectorNearUrl = env_to_strdup("BACKHAULD_REFLECTOR_NEAR_URL", "");
	config->reflectorFarUrl  = env_to_strdup("BACKHAULD_REFLECTOR_FAR_URL", "");

	config->reflectorRefreshSec = env_to_int("BACKHAULD_REFLECTOR_REFRESH_SEC", 3600);

	config->microPeriodMs     = env_to_int("BACKHAULD_MICRO_PERIOD_MS", 10000);
	config->multiPeriodMs     = env_to_int("BACKHAULD_MULTI_PERIOD_MS", 30000);
	config->chgPeriodSec      = env_to_int("BACKHAULD_CHG_PERIOD_SEC", 1800);
	config->classifyPeriodSec = env_to_int("BACKHAULD_CLASSIFY_PERIOD_SEC", 60);

	config->connectTimeoutMs  = env_to_int("BACKHAULD_CONNECT_TIMEOUT_MS", 2000);
	config->totalTimeoutMs    = env_to_int("BACKHAULD_TOTAL_TIMEOUT_MS", 10000);
	config->maxRetries        = env_to_int("BACKHAULD_MAX_RETRIES", 1);

	config->pingBytes         = env_to_int("BACKHAULD_PING_BYTES", 2048);
	config->stallThresholdMs  = env_to_int("BACKHAULD_STALL_THRESHOLD_MS", 2500);

	config->chgTargetSec      = env_to_int("BACKHAULD_CHG_TARGET_SEC", 3);
	config->chgWarmupBytes    = env_to_int("BACKHAULD_CHG_WARMUP_BYTES", 131072);
	config->chgMinBytes       = env_to_int("BACKHAULD_CHG_MIN_BYTES", 524288);
	config->chgMaxBytes       = env_to_int("BACKHAULD_CHG_MAX_BYTES", 8388608);
	config->chgSamples        = env_to_int("BACKHAULD_CHG_SAMPLES", 3);

	config->parallelStreams       = env_to_int("BACKHAULD_PARALLEL_STREAMS", 4);
	config->parallelMaxBytesTotal = env_to_int("BACKHAULD_PARALLEL_MAX_BYTES_TOTAL", 8388608);

	config->downConsecFails   = env_to_int("BACKHAULD_DOWN_CONSEC_FAILS", 6);
	config->recoverConsecOk   = env_to_int("BACKHAULD_RECOVER_CONSEC_OK", 6);
	config->capStabilityPct   = env_to_int("BACKHAULD_CAP_STABILITY_PCT", 7);

	config->windowMicroSamples = env_to_int("BACKHAULD_WINDOW_MICRO_SAMPLES", 60);
	config->windowMultiSamples = env_to_int("BACKHAULD_WINDOW_MULTI_SAMPLES", 60);
	config->windowChgSamples   = env_to_int("BACKHAULD_WINDOW_CHG_SAMPLES", 20);

	return USYS_TRUE;
}

void config_free(Config *config) {

	if (!config) return;

	free(config->serviceName);
	free(config->bootstrapHost);
	free(config->bootstrapScheme);
	free(config->bootstrapEp);

	if (config->reflectorNearUrl) free(config->reflectorNearUrl);
	if (config->reflectorFarUrl) free(config->reflectorFarUrl);

	memset(config, 0, sizeof(*config));
}


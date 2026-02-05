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


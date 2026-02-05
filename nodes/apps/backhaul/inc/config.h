/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef CONFIG_H_
#define CONFIG_H_

#include "usys_types.h"

typedef struct {

	char	*serviceName;
	int		servicePort;

	/* reflector discovery */
	char	*bootstrapHost;
	char	*bootstrapScheme;
	char	*bootstrapEp;

	char	*reflectorNearUrl;	/* optional override */
	char	*reflectorFarUrl;	/* optional override */

	int		reflectorRefreshSec;

	/* cadence */
	int		microPeriodMs;
	int		multiPeriodMs;
	int		chgPeriodSec;
	int		classifyPeriodSec;

	/* timeouts/retries */
	int		connectTimeoutMs;
	int		totalTimeoutMs;
	int		maxRetries;

	/* sizes/algorithm knobs */
	int		pingBytes;
	int		stallThresholdMs;

	int		chgTargetSec;
	int		chgWarmupBytes;
	int		chgMinBytes;
	int		chgMaxBytes;
	int		chgSamples;

	int		parallelStreams;
	int		parallelMaxBytesTotal;

	int		downConsecFails;
	int		recoverConsecOk;
	int		capStabilityPct;

	/* in-memory windows */
	int		windowMicroSamples;
	int		windowMultiSamples;
	int		windowChgSamples;

} Config;

int  config_load_from_env(Config *config);
void config_free(Config *config);
int  config_validate_env(Config *config);
void config_log(Config *config);

#endif /* CONFIG_H_ */

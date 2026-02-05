/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <string.h>
#include <time.h>

#include "backhauld.h"
#include "algo_micro_probe.h"
#include "web_client.h"
#include "usys_log.h"

int algo_micro_probe_run(Config *config, MetricsStore *store, Worker *worker) {

	char nearUrl[256] = {0};
	char farUrl[256] = {0};
	long ts = 0;

	worker_get_reflectors(worker, nearUrl, sizeof(nearUrl), farUrl, sizeof(farUrl), &ts);

	/* fallback: if reflectors missing, try bootstrap host as baseUrl (rare) */
	const char *baseUrl = NULL;
	char fallback[256] = {0};

	if (nearUrl[0]) {
		baseUrl = nearUrl;
	} else {
		snprintf(fallback, sizeof(fallback), "%s://%s",
                 config->bootstrapScheme,
                 config->bootstrapHost);
		baseUrl = fallback;
	}

	ProbeResult pr;
	memset(&pr, 0, sizeof(pr));

	wc_probe_ping(config, baseUrl, &pr);

	MicroSample s;
	memset(&s, 0, sizeof(s));
	s.ts = time(NULL);
	s.ttfbMs = pr.ttfbMs;
	s.ok = pr.ok;
	s.stalled = pr.stalled;

	metrics_store_add_micro(store, s);

	return STATUS_OK;
}

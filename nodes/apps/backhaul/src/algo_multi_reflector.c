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
#include "algo_multi_reflector.h"
#include "web_client.h"
#include "usys_log.h"

static void add_probe(MetricsStore *store, int isNear, ProbeResult pr) {

	MicroSample s;
	memset(&s, 0, sizeof(s));
	s.ts = time(NULL);
	s.ttfbMs = pr.ttfbMs;
	s.ok = pr.ok;
	s.stalled = pr.stalled;

	if (isNear) metrics_store_add_near(store, s);
	else metrics_store_add_far(store, s);
}

int algo_multi_reflector_run(Config *config, MetricsStore *store, Worker *worker) {

	char nearUrl[256] = {0};
	char farUrl[256] = {0};
	long ts = 0;

	worker_get_reflectors(worker, nearUrl, sizeof(nearUrl), farUrl, sizeof(farUrl), &ts);

	if (!nearUrl[0] || !farUrl[0]) {
		/* no reflectors yet; skip to avoid misleading stats */
		return STATUS_NOK;
	}

	ProbeResult pr;

	memset(&pr, 0, sizeof(pr));
	wc_probe_ping(config, nearUrl, &pr);
	add_probe(store, 1, pr);

	memset(&pr, 0, sizeof(pr));
	wc_probe_ping(config, farUrl, &pr);
	add_probe(store, 0, pr);

	return STATUS_OK;
}

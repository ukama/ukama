/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <pthread.h>
#include <string.h>
#include <time.h>

#include "backhauld.h"
#include "algo_diag.h"
#include "web_client.h"
#include "usys_log.h"

typedef struct {
	Config			*config;
	const char		*baseUrl;
	int				bytes;
	TransferResult	result;
} DlThreadArg;

static void* dl_thread(void *arg) {

	DlThreadArg *a = (DlThreadArg *)arg;
	memset(&a->result, 0, sizeof(a->result));
	wc_download_blob(a->config, a->baseUrl, a->bytes, &a->result);
	return NULL;
}

int algo_diag_parallel_run(Config *config, MetricsStore *store, Worker *worker) {

	char nearUrl[256] = {0};
	char farUrl[256] = {0};
	long ts = 0;

	worker_get_reflectors(worker, nearUrl, sizeof(nearUrl), farUrl, sizeof(farUrl), &ts);
	if (!nearUrl[0]) return STATUS_NOK;

	metrics_store_set_diag(store, "parallel");

	int n = config->parallelStreams;
	if (n < 2) n = 2;
	if (n > 16) n = 16;

	int bytesEach = config->parallelMaxBytesTotal / n;
	if (bytesEach < config->chgMinBytes) bytesEach = config->chgMinBytes;

	pthread_t threads[16];
	DlThreadArg args[16];

	for (int i=0; i<n; i++) {
		args[i].config = config;
		args[i].baseUrl = nearUrl;
		args[i].bytes = bytesEach;
		pthread_create(&threads[i], NULL, dl_thread, &args[i]);
	}

	double sumMbps = 0.0;
	int okCount = 0;

	for (int i=0; i<n; i++) {
		pthread_join(threads[i], NULL);
		if (args[i].result.ok) {
			sumMbps += args[i].result.mbps;
			okCount++;
		}
	}

	/* store as a CHG sample (diag), so it shows up naturally */
	ChgSample cs;
	memset(&cs, 0, sizeof(cs));
	cs.ts = time(NULL);
	cs.ok = (okCount > 0) ? 1 : 0;
	cs.dlMbps = (okCount > 0) ? sumMbps : 0.0;
	cs.ulMbps = 0.0;

	metrics_store_add_chg(store, cs);

	return STATUS_OK;
}

int algo_diag_bufferbloat_run(Config *config, MetricsStore *store, Worker *worker) {

	char nearUrl[256] = {0};
	char farUrl[256] = {0};
	long ts = 0;

	worker_get_reflectors(worker, nearUrl, sizeof(nearUrl), farUrl, sizeof(farUrl), &ts);
	if (!nearUrl[0]) return STATUS_NOK;

	metrics_store_set_diag(store, "bufferbloat");

	/* baseline ping */
	ProbeResult base;
	memset(&base, 0, sizeof(base));
	wc_probe_ping(config, nearUrl, &base);

	/* upload pulse */
	int pulseBytes = config->chgMinBytes;
	TransferResult ul;
	memset(&ul, 0, sizeof(ul));
	wc_upload_echo(config, nearUrl, pulseBytes, &ul);

	/* ping during/after load (simple v1: after; can upgrade to concurrent) */
	ProbeResult loaded;
	memset(&loaded, 0, sizeof(loaded));
	wc_probe_ping(config, nearUrl, &loaded);

	double infl = 0.0;
	if (base.ok && loaded.ok && base.ttfbMs > 1.0) {
		infl = loaded.ttfbMs / base.ttfbMs;
	}

	/* write directly into published metrics (lock lives in store) */
	pthread_mutex_lock(&store->lock);
	store->metrics.bufferbloatInflationFactor = infl;
	pthread_mutex_unlock(&store->lock);

	return STATUS_OK;
}

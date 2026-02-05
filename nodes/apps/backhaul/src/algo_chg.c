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
#include "algo_chg.h"
#include "web_client.h"
#include "usys_log.h"

static int pick_bytes_for_target(Config *config, double mbps, int minBytes, int maxBytes) {

	/* If we have a prior guess, aim for target seconds; else minBytes */
	if (mbps <= 0.1) return minBytes;

	double mbytes = (mbps / 8.0) * (double)config->chgTargetSec; /* MB */
	int bytes = (int)(mbytes * 1000000.0);

	if (bytes < minBytes) bytes = minBytes;
	if (bytes > maxBytes) bytes = maxBytes;

	return bytes;
}

int algo_chg_run(Config *config, MetricsStore *store, Worker *worker) {

	char nearUrl[256] = {0};
	char farUrl[256] = {0};
	long ts = 0;

	worker_get_reflectors(worker, nearUrl, sizeof(nearUrl), farUrl, sizeof(farUrl), &ts);

	if (!nearUrl[0]) return STATUS_NOK;

	BackhaulMetrics snap = metrics_store_get_snapshot(store);

	/* choose bytes based on prior result (keeps cost stable) */
	int dlBytes = pick_bytes_for_target(config, snap.dlGoodputMbps, config->chgMinBytes, config->chgMaxBytes);
	int ulBytes = pick_bytes_for_target(config, snap.ulGoodputMbps, config->chgMinBytes, config->chgMaxBytes);

	TransferResult dl, ul;
	memset(&dl, 0, sizeof(dl));
	memset(&ul, 0, sizeof(ul));

	/* warmup (optional, cheap) */
	if (config->chgWarmupBytes > 0) {
		TransferResult warm;
		memset(&warm, 0, sizeof(warm));
		wc_download_blob(config, nearUrl, config->chgWarmupBytes, &warm);
	}

	/* sample N times, take median Mbps (simple + robust) */
	double dlVals[16] = {0};
	double ulVals[16] = {0};
	int n = config->chgSamples;
	if (n > 16) n = 16;
	if (n <= 0) n = 1;

	int dlOk = 0, ulOk = 0;

	for (int i=0; i<n; i++) {
		memset(&dl, 0, sizeof(dl));
		wc_download_blob(config, nearUrl, dlBytes, &dl);
		if (dl.ok) dlVals[dlOk++] = dl.mbps;
	}

	for (int i=0; i<n; i++) {
		memset(&ul, 0, sizeof(ul));
		wc_upload_echo(config, nearUrl, ulBytes, &ul);
		if (ul.ok) ulVals[ulOk++] = ul.mbps;
	}

	/* compute medians (small N: simple sort) */
	double dlMed = 0.0, ulMed = 0.0;

	if (dlOk > 0) {
		for (int i=0; i<dlOk-1; i++) {
			for (int j=i+1; j<dlOk; j++) {
				if (dlVals[j] < dlVals[i]) { double t=dlVals[i]; dlVals[i]=dlVals[j]; dlVals[j]=t; }
			}
		}
		dlMed = (dlOk % 2) ? dlVals[dlOk/2] : (dlVals[dlOk/2 - 1] + dlVals[dlOk/2]) / 2.0;
	}

	if (ulOk > 0) {
		for (int i=0; i<ulOk-1; i++) {
			for (int j=i+1; j<ulOk; j++) {
				if (ulVals[j] < ulVals[i]) { double t=ulVals[i]; ulVals[i]=ulVals[j]; ulVals[j]=t; }
			}
		}
		ulMed = (ulOk % 2) ? ulVals[ulOk/2] : (ulVals[ulOk/2 - 1] + ulVals[ulOk/2]) / 2.0;
	}

	ChgSample cs;
	memset(&cs, 0, sizeof(cs));
	cs.ts = time(NULL);
	cs.ok = (dlOk > 0 || ulOk > 0) ? 1 : 0;
	cs.dlMbps = dlMed;
	cs.ulMbps = ulMed;
	cs.dlSec = dl.seconds;
	cs.ulSec = ul.seconds;

	metrics_store_add_chg(store, cs);

	return STATUS_OK;
}

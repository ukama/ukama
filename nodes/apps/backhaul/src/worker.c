/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "worker.h"
#include "usys_log.h"
#include "usys_mem.h"

#include "algo_micro_probe.h"
#include "algo_multi_reflector.h"
#include "algo_chg.h"
#include "algo_diag.h"
#include "classifier.h"

static int queue_push(Worker *w, Job j) {

	if (w->count >= w->cap) return USYS_FALSE;

	w->queue[w->tail] = j;
	w->tail = (w->tail + 1) % w->cap;
	w->count++;

	return USYS_TRUE;
}

static int queue_pop(Worker *w, Job *j) {

	if (w->count <= 0) return USYS_FALSE;

	*j = w->queue[w->head];
	w->head = (w->head + 1) % w->cap;
	w->count--;

	return USYS_TRUE;
}

int worker_init(Worker *w, Config *config, MetricsStore *store, int queueCap) {

	if (!w || !config || !store || queueCap <= 0) return USYS_FALSE;

	memset(w, 0, sizeof(*w));

	pthread_mutex_init(&w->lock, NULL);
	pthread_cond_init(&w->cond, NULL);

	pthread_mutex_init(&w->refLock, NULL);

	w->queue = (Job *)usys_calloc(queueCap, sizeof(Job));
	if (!w->queue) return USYS_FALSE;

	w->cap = queueCap;
	w->config = config;
	w->store = store;

	memset(w->nearUrl, 0, sizeof(w->nearUrl));
	memset(w->farUrl, 0, sizeof(w->farUrl));
	w->refTs = 0;

	return USYS_TRUE;
}

void worker_free(Worker *w) {

	if (!w) return;

	worker_stop(w);

	pthread_mutex_destroy(&w->lock);
	pthread_cond_destroy(&w->cond);

	pthread_mutex_destroy(&w->refLock);

	if (w->queue) usys_free(w->queue);
	memset(w, 0, sizeof(*w));
}

int worker_enqueue(Worker *w, JobType type) {

	if (!w) return USYS_FALSE;

	pthread_mutex_lock(&w->lock);

	Job j;
	j.type = type;
	j.ts = time(NULL);

	int ok = queue_push(w, j);
	if (ok) pthread_cond_signal(&w->cond);

	pthread_mutex_unlock(&w->lock);

	return ok;
}

void worker_set_reflectors(Worker *w, const char *nearUrl, const char *farUrl, long ts) {

	pthread_mutex_lock(&w->refLock);

	if (nearUrl) {
		memset(w->nearUrl, 0, sizeof(w->nearUrl));
		strncpy(w->nearUrl, nearUrl, sizeof(w->nearUrl)-1);
	}
	if (farUrl) {
		memset(w->farUrl, 0, sizeof(w->farUrl));
		strncpy(w->farUrl, farUrl, sizeof(w->farUrl)-1);
	}
	w->refTs = ts;

	pthread_mutex_unlock(&w->refLock);
}

void worker_get_reflectors(Worker *w, char *nearUrl, size_t nearLen, char *farUrl, size_t farLen, long *ts) {

	pthread_mutex_lock(&w->refLock);

	if (nearUrl && nearLen > 0) {
		memset(nearUrl, 0, nearLen);
		strncpy(nearUrl, w->nearUrl, nearLen-1);
	}
	if (farUrl && farLen > 0) {
		memset(farUrl, 0, farLen);
		strncpy(farUrl, w->farUrl, farLen-1);
	}
	if (ts) *ts = w->refTs;

	pthread_mutex_unlock(&w->refLock);
}

static void run_job(Worker *w, Job j) {

	switch (j.type) {

	case JOB_MICRO_PROBE:
		algo_micro_probe_run(w->config, w->store, w);
		break;

	case JOB_MULTI_REFLECTOR:
		algo_multi_reflector_run(w->config, w->store, w);
		break;

	case JOB_CHG:
		algo_chg_run(w->config, w->store, w);
		break;

	case JOB_CLASSIFY:
		classifier_run(w->config, w->store);
		break;

	case JOB_DIAG_PARALLEL:
		algo_diag_parallel_run(w->config, w->store, w);
		break;

	case JOB_DIAG_BUFFERBLOAT:
		algo_diag_bufferbloat_run(w->config, w->store, w);
		break;

	default:
		break;
	}
}

static void* worker_thread(void *arg) {

	Worker *w = (Worker *)arg;

	while (1) {

		pthread_mutex_lock(&w->lock);

		while (!w->stop && w->count == 0) {
			pthread_cond_wait(&w->cond, &w->lock);
		}

		if (w->stop) {
			pthread_mutex_unlock(&w->lock);
			break;
		}

		Job j;
		int ok = queue_pop(w, &j);

		pthread_mutex_unlock(&w->lock);

		if (ok) {
			run_job(w, j);
		}
	}

	return NULL;
}

int worker_start(Worker *w) {

	if (!w) return USYS_FALSE;

	w->stop = 0;
	if (pthread_create(&w->thread, NULL, worker_thread, w) != 0) {
		return USYS_FALSE;
	}

	return USYS_TRUE;
}

void worker_stop(Worker *w) {

	if (!w) return;

	pthread_mutex_lock(&w->lock);
	w->stop = 1;
	pthread_cond_signal(&w->cond);
	pthread_mutex_unlock(&w->lock);

	if (w->thread) pthread_join(w->thread, NULL);
	w->thread = 0;
}

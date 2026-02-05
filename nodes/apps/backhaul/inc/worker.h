/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef WORKER_H_
#define WORKER_H_

#include <pthread.h>
#include "config.h"
#include "metrics_store.h"

typedef enum {
	JOB_NONE = 0,
	JOB_MICRO_PROBE,
	JOB_MULTI_REFLECTOR,
	JOB_CHG,
	JOB_CLASSIFY,
	JOB_DIAG_PARALLEL,
	JOB_DIAG_BUFFERBLOAT
} JobType;

typedef struct {
	JobType	type;
	long	ts;
} Job;

typedef struct {

	pthread_mutex_t	lock;
	pthread_cond_t	cond;

	Job				*queue;
	int				cap;
	int				head;
	int				tail;
	int				count;

	int				stop;

	pthread_t		thread;

	Config			*config;
	MetricsStore		*store;

	/* current reflector set (shared) */
	pthread_mutex_t	refLock;
	char			nearUrl[256];
	char			farUrl[256];
	long			refTs;

} Worker;

int worker_init(Worker *w, Config *config,
                MetricsStore *store,
                int queueCap);
void worker_free(Worker *w);

int worker_start(Worker *w);
void worker_stop(Worker *w);

int worker_enqueue(Worker *w, JobType type);

void worker_set_reflectors(Worker *w,
                           const char *nearUrl,
                           const char *farUrl, long ts);
void worker_get_reflectors(Worker *w,
                           char *nearUrl, size_t nearLen,
                           char *farUrl, size_t farLen, long *ts);

#endif /* WORKER_H_ */

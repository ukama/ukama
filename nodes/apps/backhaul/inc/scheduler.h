/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SCHEDULER_H_
#define SCHEDULER_H_

#include <pthread.h>
#include "config.h"
#include "worker.h"

typedef struct {

	pthread_t	thread;
	int			stop;

	Config		*config;
	Worker		*worker;
} Scheduler;

int scheduler_start(Scheduler *s, Config *config, Worker *worker);
void scheduler_stop(Scheduler *s);

#endif /* SCHEDULER_H_ */

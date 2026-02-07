/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "metrics_store.h"
#include "usys_log.h"

int metrics_store_init(MetricsStore *s, const char *boardName) {

	memset(s, 0, sizeof(*s));

	if (pthread_mutex_init(&s->lock, NULL) != 0) {
		usys_log_error("metrics_store: mutex init failed");
		return -1;
	}

	memset(&s->latest, 0, sizeof(s->latest));
	s->latest.ok = 0;
	strncpy(s->latest.board, boardName ? boardName : "unknown", sizeof(s->latest.board)-1);

	return 0;
}

void metrics_store_free(MetricsStore *s) {

	if (!s) return;
	pthread_mutex_destroy(&s->lock);
	memset(s, 0, sizeof(*s));
}

void metrics_store_set(MetricsStore *s, const PowerMetrics *m) {

	pthread_mutex_lock(&s->lock);
	s->latest = *m;
	pthread_mutex_unlock(&s->lock);
}

void metrics_store_get(MetricsStore *s, PowerMetrics *out) {

	pthread_mutex_lock(&s->lock);
	*out = s->latest;
	pthread_mutex_unlock(&s->lock);
}

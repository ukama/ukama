/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef WEB_SERVICE_H
#define WEB_SERVICE_H

#include <ulfius.h>

#include "config.h"
#include "metrics_store.h"

typedef struct {
	const Config	*cfg;
	MetricsStore	*store;
} EpCtx;

int web_service_start(struct _u_instance *inst, EpCtx *ctx);
void web_service_stop(struct _u_instance *inst);

#endif

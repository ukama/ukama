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

struct _u_instance;
struct _u_request;
struct _u_response;

typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;

typedef struct {
	const Config *config;
	MetricsStore *store;
} EpCtx;

int start_web_service(Config *config, UInst *inst, EpCtx *epCtx);

void web_service_stop(struct _u_instance *inst);

int web_service_cb_get_ping(const URequest *request,
                            UResponse *response,
                            void *epConfig);

int web_service_cb_get_version(const URequest *request,
                               UResponse *response,
                               void *epConfig);

int web_service_cb_get_power(const URequest *request,
                             UResponse *response,
                             void *epConfig);

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig);

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *user_data);

#endif

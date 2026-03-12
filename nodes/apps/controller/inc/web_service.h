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
#include "driver.h"
#include "metrics_store.h"

struct _u_instance;
struct _u_request;
struct _u_response;

typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;

typedef struct {
    const Config            *config;
    MetricsStore            *store;
    const ControllerDriver  *driver;
    void                    *driver_ctx;
} EpCtx;

int  web_service_start(const Config *config, UInst *inst, EpCtx *ctx);
void web_service_stop(UInst *inst);

int web_service_cb_get_ping(const URequest *request, UResponse *response,
                            void *user_data);
int web_service_cb_get_version(const URequest *request, UResponse *response,
                               void *user_data);
int web_service_cb_get_status(const URequest *request, UResponse *response,
                              void *user_data);
int web_service_cb_get_metrics(const URequest *request, UResponse *response,
                               void *user_data);
int web_service_cb_get_alarms(const URequest *request, UResponse *response,
                              void *user_data);
int web_service_cb_put_absorption(const URequest *request, UResponse *response,
                                  void *user_data);
int web_service_cb_put_float(const URequest *request, UResponse *response,
                             void *user_data);
int web_service_cb_post_relay(const URequest *request, UResponse *response,
                              void *user_data);
int web_service_cb_default(const URequest *request, UResponse *response,
                           void *user_data);
int web_service_cb_not_allowed(const URequest *request, UResponse *response,
                               void *user_data);

#endif /* WEB_SERVICE_H */

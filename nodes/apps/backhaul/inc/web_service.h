/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef WEB_SERVICE_H
#define WEB_SERVICE_H

#include "config.h"
#include "backhauld.h"

/*
 * Forward declarations to avoid dragging ulfius.h into every file.
 * web_service.c will include <ulfius.h>.
 */
struct _u_instance;
struct _u_request;
struct _u_response;

/* Add these typedefs so UInst/URequest/UResponse are real types */
typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;

int start_web_service(Config *config, UInst *serviceInst, EpCtx *epCtx);

/* GET endpoints */
int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *epConfig);

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *epConfig);

int web_service_cb_status(const URequest *request,
                          UResponse *response,
                          void *epConfig);

int web_service_cb_metrics(const URequest *request,
                           UResponse *response,
                           void *epConfig);

/* POST diagnostics endpoints */
int web_service_cb_post_diag_chg(const URequest *request,
                                 UResponse *response,
                                 void *epConfig);

int web_service_cb_post_diag_parallel(const URequest *request,
                                      UResponse *response,
                                      void *epConfig);

int web_service_cb_post_diag_bufferbloat(const URequest *request,
                                         UResponse *response,
                                         void *epConfig);

/* default / method not allowed */
int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig);

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *user_data);

#endif /* WEB_SERVICE_H */

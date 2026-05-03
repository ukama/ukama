/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef CALLBACK_H
#define CALLBACK_H

#include <jansson.h>
#include <pthread.h>
#include <sqlite3.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>
#include <uuid/uuid.h>

#include "agent.h"
#include "err.h"
#include "log.h"
#include "ulfius.h"
#include "wimc.h"

#define WIMC_EP_STATS  "/stats"
#define WIMC_EP_CLIENT "/content/containers/*"
#define WIMC_EP_ADMIN  "/admin"

#define WIMC_EP_CONTAINER     "/content/containers"
#define WIMC_QUERY_KEY_NAME   "name"

#define WIMC_PARAM_CONTAINER_NAME    "name"
#define WIMC_PARAM_CONTAINER_TAG     "tag"
#define WIMC_PARAM_CONTAINER_PATH    "path"
#define WIMC_PARAMS_CONTAINER_STATUS "status"
#define WIMC_PARAMS_FLAGS            "flag"

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *epConfig);

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *epConfig);

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig);

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *user_data);

int web_service_cb_post_app(const URequest *request,
                            UResponse *response,
                            void *epConfig);

int web_service_cb_get_app_status(const URequest *request,
                                  UResponse *response,
                                  void *epConfig);

int web_service_cb_get_status(const URequest *request,
                              UResponse *response,
                              void *epConfig);

int web_service_cb_get_metrics(const URequest *request,
                               UResponse *response,
                               void *epConfig);

int web_service_cb_put_app_stats_update(const struct _u_request *request,
                                        struct _u_response *response,
                                        void *data);

#endif /* CALLBACK_H */

/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef CALLBACK_H
#define CALLBACK_H

#include <string.h>
#include <sqlite3.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <uuid/uuid.h>
#include <unistd.h>
#include <jansson.h>

#include "ulfius.h"
#include "log.h"
#include "err.h"
#include "wimc.h"
#include "agent.h"

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

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig);

int web_service_cb_get_capp(const URequest *request,
                            UResponse *response,
                            void *epConfig);
    
extern int db_read_path(sqlite3 *db, char *name, char *tag, char *path);

#endif /* CALLBACK_H */

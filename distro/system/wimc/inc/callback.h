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
#define WIMC_MAX_NAME_LEN 128
#define WIMC_MAX_TAG_LEN  32
#define WIMC_MAX_PATH_LEN     256

#define WIMC_PARAM_CONTAINER_NAME    "name"
#define WIMC_PARAM_CONTAINER_TAG     "tag"
#define WIMC_PARAM_CONTAINER_PATH    "path"
#define WIMC_PARAMS_CONTAINER_STATUS "status" 
#define WIMC_PARAMS_FLAGS            "flag"

#define TRUE 1
#define FALSE 0

int callback_get_container(const struct _u_request *request,
			   struct _u_response *response, void *user_data);
int callback_post_container(const struct _u_request *request,
			    struct _u_response *response, void *user_data);
int callback_put_container(const struct _u_request *request,
			   struct _u_response *response, void *user_data);
int callback_delete_container(const struct _u_request *request,
			      struct _u_response *response, void *user_data);
int callback_get_stats(const struct _u_request *request,
		       struct _u_response *response, void *user_data);
int callback_post_agent(const struct _u_request *request,
			struct _u_response *response, void *user_data);
int callback_not_allowed(const struct _u_request *request,
			 struct _u_response *response, void *user_data);
int callback_default(const struct _u_request *request,
		     struct _u_response *response, void *user_data);

#endif /* CALLBACK_H */

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef WIMC_METHODS_H
#define WIMC_METHODS_H

#include <stdio.h>
#include <ulfius.h>

#define WIMC_CMD_TRANSFER 1
#define WIMC_CMD_INFO     2
#define WIMC_CMD_INSPECT  3
#define WIMC_CMD_ALL_TAGS 4

#define WIMC_CMD_TRANSFER_STR "transfer"
#define WIMC_CMD_INFO_STR     "info"
#define WIMC_CMD_INSPECT_STR  "inspect"
#define WIMC_CMD_ALL_TAGS_STR "all-tags"
 

/* type of content to download. */
#define WIMC_TYPE_CONTAINER 1  
#define WIMC_TYPE_DATA      2

#define WIMC_TYPE_CONTAINER_STR "containers"
#define WIMC_TYPE_DATA_STR      "data"

#define WIMC_MAX_EP_LEN     1024
#define WIMC_MAX_NAME_LEN   256

#define METHOD_TYPE_GET "GET"
#define REQ_TIMEOUT     20

typedef struct _u_request req_t;
typedef struct _u_response resp_t;

req_t* create_http_request(char *url, char* ep, char *req_type);

#endif /* WIMC_METHODS_H */

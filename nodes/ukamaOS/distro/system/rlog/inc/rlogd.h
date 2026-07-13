/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef RLOGD_H_
#define RLOGD_H_

#include <ulfius.h>

#include "ingest.h"
#include "log_store.h"
#include "usys_services.h"

#define SERVICE_NAME SERVICE_RLOG
#define STATUS_OK    0
#define STATUS_NOK   -1

#define DEF_NODED_HOST    "localhost"
#define DEF_NODED_EP      "/v1/nodeinfo"
#define DEF_NODE_ID       "ukama-unknown"
#define DEF_INGEST_SOCKET "/run/ukama/rlog-ingest.sock"

#define ENV_RLOG_LOG_LEVEL      "RLOG_LOG_LEVEL"
#define ENV_RLOG_LOG_DIR        "RLOG_LOG_DIR"
#define ENV_RLOG_NODE_ID        "RLOG_NODE_ID"
#define ENV_RLOG_INGEST_SOCKET  "RLOG_INGEST_SOCKET"
#define ENV_RLOG_ROTATE_BYTES   "RLOG_ROTATE_BYTES"
#define ENV_RLOG_ROTATE_SECONDS "RLOG_ROTATE_SECONDS"
#define ENV_RLOG_RETAIN_BYTES   "RLOG_RETAIN_BYTES"
#define ENV_RLOG_RETAIN_DAYS    "RLOG_RETAIN_DAYS"
#define ENV_RLOG_BINDING_IP     "RLOG_BINDING_IP"

#define EP_BS       "/"
#define API_VERSION "v1"
#define URL_PREFIX  EP_BS API_VERSION
#define API_RES_EP(resource) EP_BS resource

#define DEF_BINDING_IP "127.0.0.1"

#define MAX_URL_LEN 512
typedef struct _u_instance UInst;
typedef struct _u_request URequest;
typedef struct _u_response UResponse;

typedef struct {
    int level;
    LogStore *store;
    IngestServer *ingest;
} ThreadData;

#endif /* RLOGD_H_ */

/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef WIMC_H
#define WIMC_H

#include <jansson.h>
#include <pthread.h>
#include <stdint.h>
#include <stdio.h>
#include <sqlite3.h>
#include <ulfius.h>
#include <uuid/uuid.h>
#include <sys/types.h>

#include "agent.h"
#include "usys_services.h"
#include "usys_types.h"

#define SERVICE_NAME       SERVICE_WIMC

#define EP_BS              "/"
#define REST_API_VERSION   "v1"
#define URL_PREFIX         EP_BS REST_API_VERSION
#define API_RES_EP(RES)    EP_BS RES

#define DEF_LOG_LEVEL      "TRACE"

#define WIMC_REQ_TYPE_FETCH  "fetch"
#define WIMC_REQ_TYPE_UPDATE "update"
#define WIMC_REQ_TYPE_CANCEL "cancel"

#define WIMC_RESP_TYPE_STATUS     "status"
#define WIMC_RESP_TYPE_RESULT     "result"
#define WIMC_RESP_TYPE_ERROR      "error"
#define WIMC_RESP_TYPE_PROCESSING "processing"

#define WIMC_CMD_TRANSFER 1
#define WIMC_CMD_INFO     2
#define WIMC_CMD_INSPECT  3
#define WIMC_CMD_ALL_TAGS 4

#define WIMC_CMD_TRANSFER_STR "transfer"
#define WIMC_CMD_INFO_STR     "info"
#define WIMC_CMD_INSPECT_STR  "inspect"
#define WIMC_CMD_ALL_TAGS_STR "all-tags"

#define WIMC_TYPE_CONTAINER 1
#define WIMC_TYPE_DATA      2

#define WIMC_TYPE_CONTAINER_STR "containers"
#define WIMC_TYPE_DATA_STR      "data"

#define WIMC_METHOD_TYPE_GET  "GET"
#define WIMC_METHOD_TYPE_POST "POST"

#define WIMC_EP_STATS    "/stats"
#define WIMC_EP_CLIENT   "/content/containers/*"
#define WIMC_EP_PROVIDER "/content/containers"
#define WIMC_EP_TASKS    "/content/tasks"

#define WIMC_EP_HUB_APPS    "v1/hub/app"
#define WIMC_EP_AGENT_UPDATE "v1/agents"

#define WIMC_MAX_NAME_LEN   256
#define WIMC_MAX_PATH_LEN   256
#define WIMC_MAX_ARGS_LEN   1024
#define WIMC_MAX_URL_LEN    1024
#define WIMC_MAX_ERR_STR    1024
#define WIMC_MAX_HUBS       8

#define WIMC_ACTION_FETCH_STR  "fetch"
#define WIMC_ACTION_UPDATE_STR "update"
#define WIMC_ACTION_CANCEL_STR "cancel"

#define WIMC_METHOD_CHUNK_STR "chunk"
#define WIMC_METHOD_TARGZ_STR "tar.gz"
#define WIMC_METHOD_TEST_STR  "test"

#define WIMC_REQ_TYPE_AGENT_STR    "agent"
#define WIMC_REQ_TYPE_PROVIDER_STR "provider"

#define MAX_AGENTS       10
#define DEFAULT_INTERVAL 10

#define WIMC_METHOD_PRIORITY_ENV "WIMC_METHOD_PRIORITY"
#define WIMC_AGENT_PATH_ENV      "WIMC_AGENT_PATH"
#define WIMC_AGENT_EXEC_NAME     "wimc-agent"

#define TRUE  1
#define FALSE 0

#define DEFAULT_SHMEM "shared_memory"
#define AGENT_EXEC    "/usr/bin/casync"
#define DEFAULT_PATH  "/tmp"

#define DEFAULT_APPS_PKGS_PATH "/ukama/apps/pkgs"
#define WIMC_DB_PATH           "/ukama/apps/db/wimc.db"

#define WIMC_HTTP_CONNECT_TIMEOUT_SEC 5L
#define WIMC_HTTP_TIMEOUT_SEC         30L
#define WIMC_MAX_HTTP_RESPONSE_BYTES  (1024 * 1024)
#define WIMC_MIN_FREE_BYTES           (128LL * 1024LL * 1024LL)
#define WIMC_MAX_PACKAGE_BYTES        (512LL * 1024LL * 1024LL)

typedef struct _u_request req_t;
typedef struct _u_response resp_t;

typedef struct {
    char *method;
    char *url;
    char *iURL;
    char *sURL;
} ServiceURL;

typedef struct {
    char *type;
    char *url;
    char *extraInfo;
    char *createdAt;
    int  size;
} ArtifactFormat;

typedef struct {
    char *name;
    char *version;
    int  formatsCount;
    ArtifactFormat **formats;
} Artifact;

typedef enum {
    WREQ_FETCH = 1,
    WREQ_UPDATE,
    WREQ_CANCEL
} WReqType;

typedef enum {
    WRESP_PROCESSING = 1,
    WRESP_UPDATE,
    WRESP_RESULT,
    WRESP_ERROR
} WRespType;

typedef enum {
    TEST = 1,
    CHUNK,
} MethodType;

typedef enum {
    WSTATUS_PEND = 1,
    WSTATUS_START,
    WSTATUS_RUNNING,
    WSTATUS_DONE,
    WSTATUS_ERROR
} TaskStatus;

typedef struct {
    char *name;
    char *tag;
    char *method;
    char *indexURL;
    char *storeURL;
    long expectedSizeBytes;
} WContent;

typedef struct {
    uuid_t   uuid;
    char     *cbURL;
    int      interval;
    WContent *content;
} WFetch;

typedef struct {
    uuid_t uuid;
    int    interval;
    char   *cbURL;
} WUpdate;

typedef struct {
    uuid_t uuid;
} WCancel;

typedef struct {
    WReqType type;
    WFetch   *fetch;
    WUpdate  *update;
    WCancel  *cancel;
} WimcReq;

typedef struct wtask {
    uuid_t   uuid;
    WContent *content;
    Update   *update;
    char     *localPath;
    TransferState state;
    struct wtask *next;
} WTasks;

typedef struct tStats {
    int start;
    int stop;
    int exitStatus;
    uint64_t n_bytes;
    uint64_t n_requests;
    uint64_t n_local_requests;
    uint64_t n_seed_requests;
    uint64_t n_remote_requests;
    uint64_t n_local_bytes;
    uint64_t n_seed_bytes;
    uint64_t n_remote_bytes;
    uint64_t total_requests;
    uint64_t total_bytes;
    uint64_t nsec;
    uint64_t runtime_nsec;
    TaskStatus status;
    char statusStr[WIMC_MAX_ERR_STR];
} TStats;

typedef struct {
    char method[WIMC_MAX_NAME_LEN];
    char service[WIMC_MAX_NAME_LEN];
    char execPath[WIMC_MAX_PATH_LEN];
    pid_t pid;
    int port;
    int running;
    int restartCount;
} ManagedAgent;

typedef struct {
    ManagedAgent agents[MAX_AGENTS];
    int count;
    pthread_mutex_t mutex;
} AgentManager;

typedef struct {
    int     servicePort;
    char    *dbFile;
    char    *hubURL;
    sqlite3 *db;
    int     maxAgents;
    Agent   **agents;
    WTasks  **tasks;

    AgentManager    *agentManager;
    pthread_mutex_t taskMutex;
    pthread_mutex_t dbMutex;
} Config;

typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;
typedef json_t              JsonObj;
typedef json_error_t        JsonErrObj;

#endif /* WIMC_H */

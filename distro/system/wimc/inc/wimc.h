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

#include <stdio.h>
#include <ulfius.h>
#include <sqlite3.h>
#include <uuid/uuid.h>

#include "agent.h"

#define WIMC_REQ_TYPE_FETCH  "fetch"
#define WIMC_REQ_TYPE_UPDATE "update"
#define WIMC_REQ_TYPE_CANCEL "cancel"

#define WIMC_RESP_TYPE_STATUS "status"
#define WIMC_RESP_TYPE_RESULT "result"
#define WIMC_RESP_TYPE_ERROR  "error"
#define WIMC_RESP_TYPE_PROCESSING "processing"


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

#define WIMC_METHOD_TYPE_GET  "GET"
#define WIMC_METHOD_TYPE_POST "POST"

#define WIMC_EP_STATS    "/stats"
#define WIMC_EP_CLIENT   "/content/containers/*"
#define WIMC_EP_PROVIDER "/content/containers"
#define WIMC_EP_TASKS    "/content/tasks"
#define WIMC_EP_ADMIN    "/admin"
#define WIMC_EP_AGENT    "/admin/agent"
#define WIMC_EP_HUB_CAPPS    "capps"
#define WIMC_EP_AGENT_UPDATE "/admin/agent/update"

#define WIMC_MAX_NAME_LEN   256
#define WIMC_MAX_PATH_LEN   256
#define WIMC_MAX_ARGS_LEN   1024
#define WIMC_MAX_URL_LEN    1024
#define WIMC_MAX_ERR_STR    1024

#define WIMC_ACTION_FETCH_STR      "fetch"
#define WIMC_ACTION_UPDATE_STR     "update"
#define WIMC_ACTION_CANCEL_STR     "cancel"

#define WIMC_METHOD_CHUNK_STR      "chunk"
#define WIMC_METHOD_TEST_STR       "test"

#define WIMC_REQ_TYPE_AGENT_STR    "agent"
#define WIMC_REQ_TYPE_PROVIDER_STR "provider"

#define MAX_AGENTS 20
#define DEFAULT_INTERVAL 10 /* 10 seconds. */

#define TRUE 1
#define FALSE 0

#define DEFAULT_SHMEM "shared_memory"
#define AGENT_EXEC "/usr/bin/casync"
#define DEFAULT_PATH "/tmp"

typedef struct _u_request req_t;
typedef struct _u_response resp_t;

typedef struct {
  char *method; /* Mechanisim supported by service at the url. */
  char *url;    /* callback URL for the agent. */
  char *iURL;   /* Index URL - only when method is chunk */
  char *sURL;   /* Chunk store URL - only when method is chunk */
} ServiceURL;

/*
 {
    "name": "cspace",
    "artifacts": [
        {
            "version": "0.0.1",
            "formats": [
                {
                    "type": "tar.gz",
                    "url": "/capps/cspace/0.0.1.tar.gz",
                    "created_at": "2022-02-24T23:53:48Z",
                    "size_bytes": 279435
                },
                {
                    "type": "chunk",
                    "url": "/capps/cspace/0.0.1.caidx",
                    "created_at": "0001-01-01T00:00:00Z",
                    "extra_info": {
                        "chunks": "/chunks/"
                    }
                }
            ]
        }
    ]
}
*/

typedef struct {

  char *type;      /* type of artifact, e.g., tgz, chunk, etc. */
  char *url;       /* to get the artifact from */
  char *extraInfo; /* any additional info. e.g., chunk URL */
  char *createdAt; /* when was it created (and updated) */
  int  size;       /* size of the artifact (in bytes) */
} ArtifactFormat;

typedef struct {

  char *name;        /* Name of the artifact */
  char *version;     /* version/tag */
  int  formatsCount; /* various format for the artifact */

  ArtifactFormat **formats; /* def of the format */
} Artifact;

typedef enum {

  WREQ_FETCH=1,
  WREQ_UPDATE,
  WREQ_CANCEL
} WReqType;

typedef enum {

  WRESP_PROCESSING=1,
  WRESP_UPDATE,
  WRESP_RESULT,
  WRESP_ERROR
} WRespType;

typedef enum {
  TEST=1,
  CHUNK,
} MethodType;

typedef enum {
  WSTATUS_PEND=1,
  WSTATUS_START,
  WSTATUS_RUNNING,
  WSTATUS_DONE,
  WSTATUS_ERROR
} TaskStatus;

typedef struct {

  char *name;        /* to fetch. */
  char *tag;         /* to fetch. */
  char *method;      /* Method to use with provider. */
  char *providerURL; /* service provider URL. */
  char *indexURL;    /* index URL for CA */
  char *storeURL;    /* chunk store */
} WContent;

typedef struct {

  uuid_t   uuid;     /* UID for future transactions */
  char     *cbURL;   /* Callback URL to send update (from Agent to wimc) */
  int      interval; /* update interval */
  WContent *content; /* Content definition */
} WFetch;

typedef struct {

  uuid_t uuid;
  int    interval;
  char   *cbURL;
} WUpdate;

typedef struct {

  uuid_t uuid;
} WCancel;

/* struct to define the request originating from wimc to agent. */
typedef struct {

  WReqType type;
  WFetch   *fetch;
  WUpdate  *update;
  WCancel  *cancel;
} WimcReq;

/* struct to define the contents activites within wimc. Each
 * client request, not found in the local db, results in adding to the
 * WTasks list.
 *
 * Client can query the task (GET) or cancel it (DELETE), both, using UUID
 * as handle.
 */

typedef struct wtask {

  uuid_t   uuid;      /* Unique ID. */
  WContent *content;  /* define the content. */
  Update   *update;   /* define status of the content activity. */
  char     *localPath;/* Path where content is available at */

  TransferState state; /* current state of the task. */

  struct wtask *next; /* pointer to next record. */
} WTasks;

typedef struct tStats {

  int start;
  int stop;
  int exitStatus;  /* as returned by waitpid */

  /* various stats */
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

  char    *clientPort;
  char    *dbFile;
  char    *adminPort;
  char    *hubURL;     /* Hub URL */
  sqlite3 *db;         /* SQLite3 db for various stats */
  int     maxAgents;   /* Max. number of agents allowed. */
  Agent   **agents;    /* Ptr to Agents, needed for http callback func() */
  WTasks  **tasks;     /* Ptr to Tasks, mostly for http cb */
} WimcCfg;

#endif /* WIMC_H */

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

#include "agent.h"

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

#define WIMC_EP_STATS  "/stats"
#define WIMC_EP_CLIENT "/content/containers/*"
#define WIMC_EP_PROVIDER "/content/containers"
#define WIMC_EP_ADMIN  "/admin"
#define WIMC_EP_AGENT  "/admin/agent"

#define WIMC_MAX_URL_LEN     1024
#define WIMC_MAX_NAME_LEN   256

#define WIMC_ACTION_FETCH_STR      "fetch"
#define WIMC_ACTION_UPDATE_STR     "update"
#define WIMC_ACTION_CANCEL_STR     "cancel"

#define WIMC_METHOD_CHUNK_STR      "chunk"
#define WIMC_METHOD_TEST_STR       "test"

#define WIMC_REQ_TYPE_AGENT_STR    "agent"
#define WIMC_REQ_TYPE_PROVIDER_STR "provider"

#define MAX_AGENTS 20

#define TRUE 1
#define FALSE 0

typedef struct _u_request req_t;
typedef struct _u_response resp_t;

typedef struct {
  int  *method; /* Mechanisim supported by service at the url. */
  char *url;    /* callback URL for the agent. */
} AgentCB;

typedef struct {

  char    *clientPort; 
  char    *dbFile;
  char    *adminPort;
  char    *cloud;      /* cloud-based service provider URL. */
  sqlite3 *db;         /* SQLite3 db for various stats */
  int     maxAgents;   /* Max. number of agents allowed. */
  Agent   **agents;     /* Ptr to Agents, needed for http callback func() */
} WimcCfg;


typedef enum {

  AGENT = 1,
  PROVIDER
} WReqType;

typedef enum {

  ACTION_FETCH=1,
  ACTION_UPDATE,
  ACTION_CANCEL
} ActionType;

typedef enum {
  TEST=1,
  CHUNK,
} MethodType;

typedef struct {

  char       *name;
  char       *tag;
  MethodType method;
  char       *providerURL;
} Content;

typedef struct {

  WReqType   type;
  ActionType action;
  int        id;
  char       *callbackURL;
  int        interval;
  Content    *content;
} WimcReq;

#endif /* WIMC_H */

/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef WIMC_AGENT_H
#define WIMC_AGENT_H

#include "err.h"
#include "log.h"

#define TRUE 1
#define FALSE 0

#define METHOD_TEST  1
#define METHOD_CHUNK 2

#define AGENT_REQ_TYPE_REG    "register"
#define AGENT_REQ_TYPE_UNREG  "unregister"
#define AGENT_REQ_TYPE_UPDATE "update"

#define WIMC_AGENT_STATE_REGISTER   1
#define WIMC_AGENT_STATE_ACTIVE     2
#define WIMC_AGENT_STATE_UNREGISTER 3

#define AGENT_TX_STATE_REQUEST_STR "request";
#define AGENT_TX_STATE_FETCH_STR   "fetch"
#define AGENT_TX_STATE_UNPACK_STR  "unpack"
#define AGENT_TX_STATE_DONE_STR    "done"
#define AGENT_TX_STATE_ERR_STR     "error"

#define AGENT_STATE_REGISTER_STR   "register"
#define AGENT_STATE_ACTIVE_STR     "active"
#define AGENT_STATE_UNREGISTER_STR "unregister"
#define AGENT_STATE_INACTIVE_STR   "inactive"

#define AGENT_REQ_TYPE_REG_STR    "register"
#define AGENT_REQ_TYPE_UNREG_STR  "unregister"
#define AGENT_REQ_TYPE_UPDATE_STR "update"

/* Type of request originating from the agent. */
typedef enum {
  REQ_REG = 1,
  REQ_UNREG,
  REQ_UPDATE
} ReqType;

typedef enum {
  CONTENT_CHUNK = 1,
  CONTENT_OCI,
  CONTENT_BINARY 
} ContentType;

typedef enum {
  REQUEST = 1,
  FETCH,
  UNPACK,
  DONE,
  ERR
} TransferState;

typedef enum {
  REGISTER = 1, 
  ACTIVE,
  UNREGISTER,
  INACTIVE
} AgentState;

typedef struct {

  int  id;
  int  totalKB;
  int  transferKB;
  int  transferState;
  char *voidStr;
} Update;

typedef struct {

  char *method;     /* method Agent support, chunk, OCI,ftp, etc. */
  char *url;        /* some url path */
} Register;

typedef struct {
  int id;
}UnRegister;

/* Struct to define the request originating from the agent. */
typedef struct {
  
  ReqType    type;    /* Type of request, ReqType enum */
  Register   *reg;
  UnRegister *unReg;
  Update     *update;
} AgentReq;

/* Struct to define content. */
typedef struct Content_ {

  char *name;   /* Name of content, e.g., container name */
  char *tag;    /* Content tag, e.g., 'latest' */
  int  type;    /* OCI-image, CA-chunk, raw, etc. */
  int  size;    /* Total size of the content, in KB */
  int  done;    /* fetch so far, in KB */
  int  state;   /* Request, fetch, unpack, done, error */
  
  struct Content_ *next; /* Next item */
} AgentContent;

typedef struct _Agent {
  
  int           id;       /* Internal ID. */
  char          *method;  /* Mechanisim supported by the agent */
  char          *url;     /* callback URL for the agent */
  int           state;    /* Register, active, un-register */
  AgentContent  *content; /* Activity*/
} Agent;

#endif /* WIMC_AGENT_H */

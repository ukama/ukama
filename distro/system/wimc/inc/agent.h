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

#define TRUE 1
#define FALSE 0

#define AGENT_REQ_TYPE_REG    "register"
#define AGENT_REQ_TYPE_UNREG  "unregister"
#define AGENT_REQ_TYPE_UPDATE "update"

#define WIMC_AGENT_STATE_REGISTER   1
#define WIMC_AGENT_STATE_ACTIVE     2
#define WIMC_AGENT_STATE_UNREGISTER 3

#define WIMC_AGENT_OK               0x01
#define WIMC_AGENT_ERROR_EXIST      0x02 /* Agent already exists */
#define WIMC_AGENT_ERROR_BAD_METHOD 0x04
#define WIMC_AGENT_ERROR_BAD_URL    0x08
#define WIMC_AGENT_ERROR_MEMORY     0X16 

#define WIMC_AGENT_OK_STR               "OK"
#define WIMC_AGENT_ERROR_EXIST_STR      "Already Registered" 
#define WIMC_AGENT_ERROR_BAD_METHOD_STR "Bad method"
#define WIMC_AGENT_ERROR_BAD_URL_STR    "Invalid ULR"
#define WIMC_AGENT_ERROR_MEMORY_STR     "Internal memory error"

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

typedef struct {

  int          state;    /* Register, active and un-register. */
  AgentContent *content; /* Content related activity */
} AgentWork;

typedef struct _Agent {
  
  int           id;       /* Internal ID. */
  char          *method;  /* Mechanisim supported by the agent */
  char          *url;     /* callback URL for the agent */
  AgentWork     *work;
  struct _Agent *next;    /* Pointer to next. */
} Agent;

#endif /* WIMC_AGENT_H */

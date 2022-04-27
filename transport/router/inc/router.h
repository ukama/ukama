/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef ROUTER_H
#define ROUTER_H

#include <uuid/uuid.h>
#include <ulfius.h>

#define EP_ROUTE   "/routes"
#define EP_STATS   "/stats"
#define EP_SERVICE "/service"

#define DEFAULT_PATTERN_PATH "/service"

#define FALSE   0
#define TRUE    1
#define UUID    2
#define ERROR   3

#define MAX_LEN 1024
#define MAX_POST_BODY_SIZE 2048

#define UKAMA_ERROR_NONE           1
#define UKAMA_ERROR_INVALID_DEST   2
#define UKAMA_ERROR_BAD_DEST       3

typedef struct _u_request  req_t;
typedef struct _u_response resp_t;

typedef struct {

  char *hostName;
  char *port;
} Config;

/* Pattern is list of key value pair. Value can be regex */
typedef struct pattern_ {

  char *key;
  char *value;

  struct pattern_ *next;
} Pattern;

typedef struct patterns_ {

  Pattern *pattern; /* List of key-value pair */
  char    *path;    /* EP */

  struct patterns_ *next;
} Patterns;

/* Various stats */
typedef struct {

  int msgCount;
} Stats;

/* where to forward the message */
typedef struct {

  char *ip;
  int  port;
} Forward;

typedef struct service_ {

  uuid_t   uuid;      /* UUID assigned */
  char     *name;     /* service name */
  Patterns *patterns; /* List of patterns to match */
  Forward  *forward;  /* Where to forward the message */

  struct service_ *next;
} Service;

/*
 * Router definition
 */
typedef struct {

  Config  *config;   /* router configuration */
  Stats   *stats;    /* overall router stats */
  Service *services; /* Services msg routable via this router */
} Router;

#endif /* ROUTER_H */

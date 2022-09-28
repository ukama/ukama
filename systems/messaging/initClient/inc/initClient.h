/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INIT_CLIENT_H
#define INIT_CLIENT_H

#include <ulfius.h>

#define TRUE  1
#define FALSE 0

#define DEF_LOG_LEVEL   "DEBUG"
#define MAX_BUFFER_SIZE 1024
#define MAX_URL_LEN     1024

#define EP_PING    "/ping"
#define EP_ORGS    "orgs"
#define EP_SYSTEMS "systems"

typedef struct _u_instance UInst;
typedef struct _u_request  URequest;
typedef struct _u_response UResponse;

/* Type of request originating from the agent. */
typedef enum {
	REQ_REGISTER = 1,
	REQ_UNREGISTER,
	REQ_UPDATE
} ReqType;

typedef struct {

	char *org;
	char *name;
	char *cert;
	char *ip;
	char *port;
} Register;

typedef struct {

	ReqType  reqType;
	Register *reg;
} Request;

struct Response {

	char *buffer;
	size_t size;
};

#endif /* INIT_CLIENT_H */

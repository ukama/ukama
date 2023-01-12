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

#include "config.h"

#define TRUE  1
#define FALSE 0

#define DEFAULT_LOG_LEVEL  "DEBUG"
#define DEFAULT_SYSTEM_ORG "Ukama"
#define DEFAULT_API_VER    "v1"

#define MAX_BUFFER_SIZE 1024
#define MAX_URL_LEN     1024
#define MAX_UUID_LEN    37

#define REG_STATUS_NONE         0x00
#define REG_STATUS_NO_UUID      0x04
#define REG_STATUS_HAVE_UUID    0x08
#define REG_STATUS_NO_MATCH     0x16
#define REG_STATUS_MATCH        0x32

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
	REQ_UPDATE,
	REQ_QUERY
} ReqType;

typedef struct {

	char *org;
	char *name;
	char *cert;
	char *ip;
	char *port;
} Register;

typedef struct {

	char *systemName;
	char *systemID;
	char *certificate;
	char *ip;
	int  port;
	int  health;
} QueryResponse;

typedef struct {

	ReqType  reqType;
	Register *reg;
} Request;

struct Response {

	char *buffer;
	size_t size;
};

void free_query_response(QueryResponse *response);
int send_request_to_init(ReqType reqType, Config *config, char **response);
int existing_registration(Config *config, char **cacheUUID, char **systemUUID);

#endif /* INIT_CLIENT_H */

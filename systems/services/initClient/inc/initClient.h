/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#ifndef INIT_CLIENT_H
#define INIT_CLIENT_H

#include <ulfius.h>

#include "config.h"

#define TRUE  1
#define FALSE 0

#define REGISTER_TO_GLOBAL_INIT 1
#define REGISTER_TO_LOCAL_INIT  0

#define DEFAULT_LOG_LEVEL  "DEBUG"
#define DEFAULT_SYSTEM_ORG "Ukama"
#define DEFAULT_API_VER    "v1"

#define MAX_BUFFER_SIZE 1024
#define MAX_URL_LEN     1024
#define MAX_UUID_LEN    37

#define REG_STATUS_NONE         	0x00
#define REG_STATUS_NO_UUID      	0x02
#define REG_STATUS_HAVE_UUID    	0x04
#define REG_STATUS_NO_MATCH     	0x08
#define REG_STATUS_MATCH        	0x10
#define REG_STATUS_PARSING_FAILURE  0x20

#define QUERY_OK    0x00
#define QUERY_ERROR 0x01

#define EP_PING    "/ping"
#define EP_SYSTEMS "/v1/systems"
#define ORGS_STR    "orgs"
#define SYSTEMS_STR "systems"

#define INIT_CLIENT_NAME_STR                  "name"
#define INIT_CLIENT_ERROR_INVALID_KEY_STR     "invalid key"
#define INIT_CLIENT_ERROR_INVALID_SYSTEM_NAME "invalid system name"
#define INIT_CLIENT_ORG_NAME_STR			  "org"


typedef struct _u_instance UInst;
typedef struct _u_request  URequest;
typedef struct _u_response UResponse;

typedef enum {
	REQ_REGISTER = 1,
	REQ_UNREGISTER,
	REQ_UPDATE,
	REQ_QUERY,
	REQ_QUERY_SYSTEM
} ReqType;

typedef struct {

	char *org;
	char *name;
	char *cert;
	char *apiGwIp;
	char *apiGwPort;
    char *nodeGwIp;
    char *nodeGwPort;
} Register;

typedef struct {

	char *systemName;
	char *systemID;
	char *certificate;
	char *apiGwIp;
	int  apiGwPort;
    char *nodeGwIp;
    int  nodeGwPort;
	int  health;
} QueryResponse;

typedef struct {

	ReqType  reqType;
	Register *reg;
} Request;

typedef struct {
	char* localUUID;
	char* globalUUID;
} SystemRegistrationId;

struct Response {

	char *buffer;
	size_t size;
};

void free_query_response(QueryResponse *response);
int send_request_to_init(ReqType reqType, Config *config, char* org,
						 char *systemName, char **response, int global );
int existing_registration(Config *config, char **cacheUUID, char **systemUUID,
                          int global);
int get_system_info(Config *config, char *org, char *systemName,
                    char **systemInfo, int global);
int parse_cache_uuid(char *fileName, SystemRegistrationId **sysReg);

#endif /* INIT_CLIENT_H */

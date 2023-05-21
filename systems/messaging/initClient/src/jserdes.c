/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <stdio.h>
#include <string.h>
#include <jansson.h>

#include "initClient.h"
#include "jserdes.h"
#include "log.h"

/* JSON (de)-serialization functions. */

/*
 * serialize_request --
 *
 */
int serialize_request(Request *request, json_t **json) {

	int ret=FALSE;
	Register *reg;
	char *str=NULL;

	if (request->reqType == (ReqType)REQ_REGISTER) {

		*json = json_object();
		if (*json == NULL) {
			return FALSE;
		}

		reg = request->reg;

		json_object_set_new(*json, JSON_IP,   json_string(reg->ip));
		json_object_set_new(*json, JSON_PORT, json_integer(atoi(reg->port)));
		json_object_set_new(*json, JSON_CERTIFICATE, json_string(reg->cert));

		str = json_dumps(*json, 0);
		if (str) {
			log_debug("Registration JSON: %s", str);
			free(str);
		}
		ret = TRUE;
	} else if (request->reqType == (ReqType)REQ_UNREGISTER) {
	  ret = TRUE;
	} else if (request->reqType == (ReqType)REQ_QUERY) {
	  ret = TRUE;
	} else if (request->reqType == (ReqType)REQ_QUERY_SYSTEM) {
		ret = TRUE;
	}

	return ret;
}

/*
 * deserialize_response --
 *
 */
int deserialize_response(ReqType reqType, QueryResponse **queryResponse,
						 char *str) {

	int ret=TRUE;
	json_t *json=NULL;
	json_t *name, *id, *cert, *ip, *port, *health;

	if (str == NULL) return FALSE;

	json = json_loads(str, JSON_DECODE_ANY, NULL);
	if (!json) {
		log_error("Can no load str into JSON object. Str: %s", str);
		return FALSE;
	}

	name   = json_object_get(json, JSON_SYSTEM_NAME);
	id     = json_object_get(json, JSON_SYSTEM_ID);
	cert   = json_object_get(json, JSON_CERTIFICATE);
	ip     = json_object_get(json, JSON_IP);
	port   = json_object_get(json, JSON_PORT);

	if (reqType == (ReqType)REQ_QUERY) {
		health = json_object_get(json, JSON_HEALTH);
	}

	if (!name || !id || !cert || !ip || !port) {
		log_error("Error deserilaizing response");
		ret = FALSE;
		goto failure;
	}

	if (reqType == (ReqType)REQ_QUERY && !health) {
		log_error("Error deserilaizing response");
		ret = FALSE;
		goto failure;
	}

	*queryResponse = (QueryResponse *)calloc(1, sizeof(QueryResponse));
	if (*queryResponse == NULL) {
		log_error("Memory allocation error of size: %ld",
				  sizeof(QueryResponse));
		ret = FALSE;
		goto failure;
	}

	(*queryResponse)->systemName  = strdup(json_string_value(name));
	(*queryResponse)->systemID    = strdup(json_string_value(id));
	(*queryResponse)->certificate = strdup(json_string_value(cert));
	(*queryResponse)->ip          = strdup(json_string_value(ip));
	(*queryResponse)->port        = json_integer_value(port);

	if (reqType == (ReqType)REQ_QUERY) {
		(*queryResponse)->health      = json_integer_value(health);
	}

 failure:
	json_decref(json);
	return ret;
}

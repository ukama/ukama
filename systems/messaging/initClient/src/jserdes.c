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

#include "jserdes.h"
#include "log.h"

/* JSON (de)-serialization functions. */

/*
 * serialize_agent_request --
 *
 */
int serialize_request(Request *request, json_t **json) {

	int ret=FALSE;
	Register *reg;
	char *str=NULL;

	*json = json_object();
	if (*json == NULL) {
		return FALSE;
	}

	if (request->reqType == (ReqType)REQ_REGISTER) {

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
	}

	return ret;
}

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <jansson.h>
#include <uuid/uuid.h>

#include "mesh.h"
#include "jserdes.h"

#include "static.h"

STATIC void add_map_to_request(json_t **json, UMap *map, int mapType) {

	json_t *jMap=NULL, *jArray=NULL;
	int i;

	if (map == NULL) {
		return;
	}

	if (map->nb_values == 0) {
		return;
	}

	json_object_set_new(*json, JSON_MAP, json_object());
	jMap = json_object_get(*json, JSON_MAP);

	if (mapType == MESH_MAP_TYPE_URL) {
		json_object_set_new(jMap, JSON_TYPE,
							json_string(MESH_MAP_TYPE_URL_STR));
	} else if (mapType == MESH_MAP_TYPE_HDR) {
		json_object_set_new(jMap, JSON_TYPE,
							json_string(MESH_MAP_TYPE_HDR_STR));
	} else if (mapType == MESH_MAP_TYPE_POST) {
		json_object_set_new(jMap, JSON_TYPE,
							json_string(MESH_MAP_TYPE_HDR_STR));
	} else if (mapType == MESH_MAP_TYPE_COOKIE) {
		json_object_set_new(jMap, JSON_TYPE,
							json_string(MESH_MAP_TYPE_COOKIE_STR));
	}

	/* For array of key/value pair. */
	json_object_set_new(jMap, JSON_DATA, json_array());
	jArray = json_object_get(jMap, JSON_DATA);

	if (jArray) {

		for (i=0; i < map->nb_values; i++) {

			json_t *entry = json_object();

			json_object_set_new(entry, JSON_KEY, json_string(map->keys[i]));
			json_object_set_new(entry, JSON_VALUE, json_string(map->values[i]));
			json_object_set_new(entry, JSON_LEN,
								json_integer((int)map->lengths[i]));

			json_array_append_new(jArray, entry);
		}
	}
}

int serialize_device_info(json_t **json, NodeInfo *device) {

	json_t *jDevice=NULL;

	if (device == NULL) {
		return FALSE;
	}

	if (device->nodeID == NULL) {
		return FALSE;
	}

	*json = json_object();
	if (*json == NULL) {
		return FALSE;
	}

	/* Add Device info. Currently only is the UUID. */
	json_object_set_new(*json, JSON_NODE_INFO, json_object());
	jDevice = json_object_get(*json, JSON_NODE_INFO);

	if (jDevice == NULL) {
		json_decref(*json);
		*json=NULL;
		return FALSE;
	}

	json_object_set_new(jDevice, JSON_NODE_ID, json_string(device->nodeID));

	return TRUE;
}

int serialize_local_service_response(char **response, Message *message,
                                     int code, int len, char *data) {

    json_t *json, *obj;
    
    /* basic sanity check */
	if (len == 0 || data == NULL || message == NULL)
		return FALSE;

	json = json_object();
	if (json == NULL) {
		return FALSE;
	}

	json_object_set_new(json, JSON_TYPE, json_string(MESH_SERVICE_RESPONSE));
    json_object_set_new(json, JSON_UUID, json_string(message->seqNo));

	/* Add response info. */
	json_object_set_new(json, JSON_MESSAGE, json_object());
	obj = json_object_get(json, JSON_MESSAGE);
    json_object_set_new(obj, JSON_CODE,   json_integer(code));
	json_object_set_new(obj, JSON_LENGTH, json_integer(len));
	json_object_set_new(obj, JSON_DATA,   json_string(data));

    *response = json_dumps(json, 0);
    json_decref(json);

	return TRUE;
}

STATIC void serialize_message_data(URequest *request, char **data) {

    json_t *json, *jRaw;

    json = json_object();
    if (json == NULL) return;
    
	json_object_set_new(json, JSON_PROTOCOL,
						json_string(request->http_protocol));
	json_object_set_new(json, JSON_METHOD, json_string(request->http_verb));
	json_object_set_new(json, JSON_URL, json_string(request->http_url));
	json_object_set_new(json, JSON_PATH, json_string(request->url_path));

	/* Add maps if they exists. */
	if (request->map_url) {
		add_map_to_request(&json, request->map_url, MESH_MAP_TYPE_URL);
	}

	if (request->map_header) {
		add_map_to_request(&json, request->map_header, MESH_MAP_TYPE_HDR);
	}

	if (request->map_cookie) {
		add_map_to_request(&json, request->map_cookie, MESH_MAP_TYPE_COOKIE);
	}

	if (request->map_post_body) {
		add_map_to_request(&json, request->map_post_body, MESH_MAP_TYPE_POST);
	}

	/* And finally add raw binary data. Currently we assume raw is char* */
	if (request->binary_body_length > 0 && request->binary_body != NULL ){
		json_object_set_new(json, JSON_RAW_DATA, json_object());
		jRaw = json_object_get(json, JSON_RAW_DATA);
		json_object_set_new(jRaw, JSON_LENGTH,
							json_integer((int)request->binary_body_length));
		json_object_set_new(jRaw, JSON_DATA,
							json_string((char *)request->binary_body));
	}

    *data = json_dumps(json, 0);
    json_decref(json);
}

int serialize_websocket_message(char **str, URequest *request, char *uuid) {

    json_t *json=NULL, *jDevice=NULL, *jService=NULL;
	json_t *jRequest=NULL;
    char *data=NULL;

	json = json_object();
	if (json == NULL) {
		return FALSE;
	}

	json_object_set_new(json, JSON_TYPE, json_string(MESH_NODE_REQUEST));
	json_object_set_new(json, JSON_UUID,  json_string(uuid));

    serialize_message_data(request, &data);
	json_object_set_new(json, JSON_MESSAGE, json_object());
	jRequest = json_object_get(json, JSON_MESSAGE);
	json_object_set_new(jRequest, JSON_LENGTH, json_integer(strlen(data)));
    json_object_set_new(jRequest, JSON_CODE,   json_integer(0));
	json_object_set_new(jRequest, JSON_DATA,   json_string(data));

    *str = json_dumps(json, 0);

	return TRUE;
}

int deserialize_node_info(NodeInfo **node, json_t *json) {

	json_t *id, *port;

	if (json == NULL && node == NULL)
		return FALSE;
    
    id = json_object_get(json, JSON_NODE_ID);
    port = json_object_get(json, JSON_PORT);

    if (id == NULL || port == NULL) {
        return FALSE;
    }
    
	*node = (NodeInfo *)calloc(1, sizeof(NodeInfo));
	if (*node == NULL) {
        log_error("Error allocating memory of size: %d", sizeof(NodeInfo));
		return FALSE;
    }

    (*node)->nodeID = strdup(json_string_value(id));
    (*node)->port   = strdup(json_string_value(port));

	return TRUE;
}

STATIC int deserialize_service_info(ServiceInfo **service, json_t *json) {

	json_t *name, *port;
  
	if (json == NULL && service == NULL) return FALSE;

    name = json_object_get(json, JSON_NAME);
    port = json_object_get(json, JSON_PORT);

    if (name == NULL || port == NULL) return FALSE;

	*service = (ServiceInfo *)calloc(1, sizeof(ServiceInfo));
	if (*service == NULL) return FALSE;

    (*service)->name = strdup(json_string_value(name));
    (*service)->port = strdup(json_string_value(port));

	return TRUE;
}

STATIC void deserialize_map_array(UMap **map, json_t *json) {

	json_t *jArray;
	json_t *elem, *key, *val, *len;
	int i, size=0;

	*map = (UMap *)calloc(1, sizeof(UMap));
	if (*map==NULL)
		return;

	u_map_init(*map);

	jArray = json_object_get(json, JSON_DATA);

	if (json_is_array(jArray)) {
		size = json_array_size(jArray);

		for (i=0; i<size; i++) {
			elem = json_array_get(jArray, i);

			key = json_object_get(elem, JSON_KEY);
			val = json_object_get(elem, JSON_VALUE);
			len = json_object_get(elem, JSON_LEN);

			u_map_put(*map, json_string_value(key), json_string_value(val));
		}
	}
}

STATIC void deserialize_map(URequest **request, json_t *json) {

	json_t *obj;
	char *str;

	/* Determine the type of map. */
	obj = json_object_get(json, JSON_TYPE);
	if (obj==NULL) {
		return;
	}

	str = json_string_value(obj);

	if (strcasecmp(str, MESH_MAP_TYPE_URL_STR)==0) {
		deserialize_map_array(&(*request)->map_url, json);
	} else if (strcasecmp(str, MESH_MAP_TYPE_HDR_STR)==0) {
		deserialize_map_array(&(*request)->map_header, json);
	} else if (strcasecmp(str, MESH_MAP_TYPE_POST_STR)==0) {
		deserialize_map_array(&(*request)->map_post_body, json);
	} else if (strcasecmp(str, MESH_MAP_TYPE_COOKIE_STR)==0) {
		deserialize_map_array(&(*request)->map_cookie, json);
	}
}

int deserialize_request_info(URequest **request, char *str) {

	json_t *json, *obj, *jRaw;
	int size, i;

	if (str == NULL) return FALSE;

    json = json_loads(str, JSON_DECODE_ANY, NULL);
    if (json == NULL) return FALSE;

	*request = (URequest *)calloc(1, sizeof(URequest));
	if (*request == NULL) {
        log_error("Error allocating memory of size: %d", sizeof(URequest));
		return FALSE;
    }

	obj = json_object_get(json, JSON_PROTOCOL);
	if (obj) {
		(*request)->http_protocol = strdup(json_string_value(obj));
	} else {
		(*request)->http_protocol = strdup("");
	}

	obj = json_object_get(json, JSON_METHOD);
	if (obj) {
		(*request)->http_verb = strdup(json_string_value(obj));
	} else {
		(*request)->http_verb = strdup("");
	}

	obj = json_object_get(json, JSON_URL);
	if (obj) {
		(*request)->http_url = strdup(json_string_value(obj));
	} else {
		(*request)->http_url = strdup("");
	}

	obj = json_object_get(json, JSON_PATH);
	if (obj) {
		(*request)->url_path = strdup(json_string_value(obj));
	} else {
		(*request)->url_path = strdup("");
	}

	/* de-ser the various map. URL, Header, POST, Cookie */
	for (i=0; i < 4; i++) {
		obj = json_object_get(json, JSON_MAP);
		if (obj) {
			deserialize_map(request, obj);
		}
	}

	/* Lastly, de-serialize raw binary data. */
	jRaw = json_object_get(json, JSON_RAW_DATA);
	if (jRaw) {

		obj = json_object_get(jRaw, JSON_LENGTH);
		size = json_integer_value(obj);
		(*request)->binary_body_length = size;

		/* Get the actual data now */
		obj = json_object_get(jRaw, JSON_DATA);

		if (obj) {

			(*request)->binary_body = (void *)calloc(1, size);
			if ((*request)->binary_body == NULL)
				return FALSE;

			memcpy((*request)->binary_body, (void *)json_string_value(obj),
				   size);
		}
	}

	return TRUE;
}

int deserialize_websocket_message(Message **out, json_t *json) {
    json_t *jType, *jSeq, *jMessage;
    json_t *jLength, *jData, *jCode;
    const char *type_s = NULL;
    const char *seq_s  = NULL;
    const char *data_s = NULL;

    if (!out || !json) return FALSE;
    *out = NULL;

    jType    = json_object_get(json, JSON_TYPE);
    jSeq     = json_object_get(json, JSON_UUID);
    jMessage = json_object_get(json, JSON_MESSAGE);

    if (!json_is_string(jType) || !json_is_string(jSeq) || !json_is_object(jMessage)) {
        log_error("Error decoding JSON. Missing/invalid envelope fields.");
        return FALSE;
    }

    jLength = json_object_get(jMessage, JSON_LENGTH);
    jData   = json_object_get(jMessage, JSON_DATA);
    jCode   = json_object_get(jMessage, JSON_CODE); /* OPTIONAL */

    if (!json_is_integer(jLength) || !json_is_string(jData)) {
        log_error("Error decoding JSON. Missing/invalid message fields.");
        return FALSE;
    }

    type_s = json_string_value(jType);
    seq_s  = json_string_value(jSeq);
    data_s = json_string_value(jData);

    Message *m = (Message *)calloc(1, sizeof(*m));
    if (!m) {
        log_error("OOM allocating Message");
        return FALSE;
    }

    m->reqType  = strdup(type_s ? type_s : "");
    m->seqNo    = strdup(seq_s ? seq_s : "");
    m->dataSize = (int)json_integer_value(jLength);
    m->data     = strdup(data_s ? data_s : "");

    if (json_is_integer(jCode)) {
        m->code = (int)json_integer_value(jCode);
    } else {
        m->code = 0;
    }

    if (!m->reqType || !m->seqNo || !m->data) {
        log_error("OOM duplicating Message fields");
        free(m->reqType);
        free(m->seqNo);
        free(m->data);
        free(m);
        return FALSE;
    }

    if ((int)strlen(m->data) != m->dataSize) {
        log_error("Message length mismatch: declared=%d actual=%zu",
                  m->dataSize, strlen(m->data));
        //        clear_message(&m);
        return FALSE;
    }

    *out = m;
    return TRUE;
}

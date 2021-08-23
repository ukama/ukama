/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <stdio.h>
#include <string.h>
#include <jansson.h>
#include <uuid/uuid.h>

#include "mesh.h"
#include "jserdes.h"

/* JSON for the forward request from device to service provider.

   fwd_request -> { type: "request",
                    seq_no: "1234",
                    device_info: {
                          uuid: "uuid"
                    },
                    service_info: {
                          uuid: "service_uuid"
                    },
                    request_info: {
		          protocol: "HTTP/1.1",
			  method: "GET",
			  url: "locahost:3456/",
			  path: "/some/path",
			  map: { type:"url",
			      data:[ {key_name: "name1", key_value:"value1"},
			             {key_name: "name2", key_value:"value2"},
				   ],
			  raw: { length:"1234",
			         data: "xdddgfdg"
			       },
		   }
		 }
*/

/* JSON for the response from the service provider, via mesh.d

    fwd_request -> { type: "response",
                     seq_no: "1234",
                     service_info: {
                          uuid: "service_uuid"
                     },
                     response_info: {
			  raw: { length:"1234",
			         data: "xdddgfdg"
			       },
		    }
		  }
*/

/*
 * add_map_to_request --
 *
 */
static void add_map_to_request(json_t **json, UMap *map, int mapType) {

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
    json_object_set_new(jMap, JSON_TYPE, json_string(MESH_MAP_TYPE_URL_STR));
  } else if (mapType == MESH_MAP_TYPE_HDR) {
    json_object_set_new(jMap, JSON_TYPE, json_string(MESH_MAP_TYPE_HDR_STR));
  } else if (mapType == MESH_MAP_TYPE_POST) {
    json_object_set_new(jMap, JSON_TYPE, json_string(MESH_MAP_TYPE_HDR_STR));
  } else if (mapType == MESH_MAP_TYPE_COOKIE) {
    json_object_set_new(jMap, JSON_TYPE, json_string(MESH_MAP_TYPE_COOKIE_STR));
  }

  /* For array of key/value pair. */
  json_object_set_new(jMap, JSON_DATA, json_array());
  jArray = json_object_get(jMap, JSON_DATA);

  if (jArray) {

    for (i=0; i < map->nb_values; i++) {

      json_t *entry = json_object();

      json_object_set_new(entry, JSON_KEY, json_string(map->keys[i]));
      json_object_set_new(entry, JSON_VALUE, json_string(map->values[i]));
      json_object_set_new(entry, JSON_LEN, json_integer((int)map->lengths[i]));

      json_array_append_new(jArray, entry);
    }
  }
}

/*
 * serialize_response --
 *
 */
int serialize_response(json_t **json, int size, void *data, uuid_t uuid) {

  json_t *jResp=NULL, *jRespInfo=NULL, *jRaw=NULL;
  json_t *jService=NULL;
  char idStr[36+1];
  char *jStr;

  /* basic sanity check */
  if (size == 0 && data == NULL && uuid_is_null(uuid))
    return FALSE;

  *json = json_object();
  if (*json == NULL) {
    return FALSE;
  }

  json_object_set_new(*json, JSON_MESH_FORWARD, json_object());
  jResp = json_object_get(*json, JSON_MESH_FORWARD);

  if (jResp==NULL) {
    json_decref(*json);
    *json=NULL;
    return FALSE;
  }

  json_object_set_new(jResp, JSON_TYPE, json_string(MESH_TYPE_FWD_RESP));
  json_object_set_new(jResp, JSON_SEQ, json_integer(123)); /* xxx */

  /* Service Info. */
  uuid_unparse(uuid, &idStr[0]);
  json_object_set_new(jResp, JSON_SERVICE_INFO, json_object());
  jService = json_object_get(jResp, JSON_SERVICE_INFO);
  json_object_set_new(jService, JSON_ID, json_string(&idStr[0]));

  /* Add response info. */
  json_object_set_new(jResp, JSON_RESPONSE_INFO, json_object());
  jRespInfo = json_object_get(jResp, JSON_RESPONSE_INFO);

  /* Add raw data */
  json_object_set_new(jRespInfo, JSON_RAW_DATA, json_object());
  jRaw = json_object_get(jRespInfo, JSON_RAW_DATA);

  json_object_set_new(jRaw, JSON_LENGTH, json_integer(size));
  json_object_set_new(jRaw, JSON_DATA, json_string((char *)data));

  jStr = json_dumps(*json, 0);
  log_debug("Serialized response: %s", jStr);
  free(jStr);

  return TRUE;
}

/*
 * serialize_forward_request --
 *
 */

int serialize_forward_request(URequest *request, json_t **json,
			      Config *config, uuid_t uuid) {

  int ret=FALSE;
  json_t *jReq=NULL, *jDevice=NULL, *jService=NULL;
  json_t *jRequest=NULL, *jRaw=NULL;
  char idStr[36+1]; /* 36-bytes for UUID + trailing '\0' */

  /* Basic sanity check. */
  if (request==NULL && config==NULL) {
    return ret;
  }

  *json = json_object();
  if (*json == NULL) {
    return ret;
  }

  json_object_set_new(*json, JSON_MESH_FORWARD, json_object());
  jReq = json_object_get(*json, JSON_MESH_FORWARD);

  if (jReq==NULL) {
    return ret;
  }

  json_object_set_new(jReq, JSON_TYPE, json_string(MESH_TYPE_FWD_REQ));
  json_object_set_new(jReq, JSON_SEQ, json_integer(123));

  /* Add Device info. Currently only is the UUID. */
  json_object_set_new(jReq, JSON_DEVICE_INFO, json_object());
  jDevice = json_object_get(jReq, JSON_DEVICE_INFO);

  uuid_unparse(config->uuid, &idStr[0]);
  json_object_set_new(jDevice, JSON_ID, json_string(idStr));

  /* Add service info., service is the one whose request is being forward. */
  uuid_unparse(uuid, &idStr[0]); /* Service UUID. */

  json_object_set_new(jReq, JSON_SERVICE_INFO, json_object());
  jService = json_object_get(jReq, JSON_SERVICE_INFO);

  json_object_set_new(jService, JSON_ID, json_string(&idStr[0]));

  /* Add request info. */
  json_object_set_new(jReq, JSON_REQUEST_INFO, json_object());
  jRequest = json_object_get(jReq, JSON_REQUEST_INFO);

  json_object_set_new(jRequest, JSON_PROTOCOL,
		      json_string(request->http_protocol));
  json_object_set_new(jRequest, JSON_METHOD, json_string(request->http_verb));
  json_object_set_new(jRequest, JSON_URL, json_string(request->http_url));
  json_object_set_new(jRequest, JSON_PATH, json_string(request->url_path));

  /* Add map if they exists. */
  if (request->map_url) {
    add_map_to_request(&jRequest, request->map_url, MESH_MAP_TYPE_URL);
  }

  if (request->map_header) {
    add_map_to_request(&jRequest, request->map_header, MESH_MAP_TYPE_HDR);
  }

  if (request->map_cookie) {
    add_map_to_request(&jRequest, request->map_cookie, MESH_MAP_TYPE_COOKIE);
  }

  if (request->map_post_body) {
    add_map_to_request(&jRequest, request->map_post_body, MESH_MAP_TYPE_POST);
  }

  /* And finally add raw binary data. Currently we assume raw is char* */
  if (request->binary_body_length > 0 && request->binary_body != NULL ){

    json_object_set_new(jRequest, JSON_RAW_DATA, json_object());
    jRaw = json_object_get(jRequest, JSON_RAW_DATA);

    json_object_set_new(jRaw, JSON_LENGTH,
			json_integer((int)request->binary_body_length));
    json_object_set_new(jRaw, JSON_DATA,
			json_string((char *)request->binary_body));
  }

  return TRUE;
}

/*
 * deserialize_device_info --
 *
 */
static int deserialize_device_info(DeviceInfo **device, json_t *json) {

  json_t *obj;

  if (json == NULL && device == NULL)
    return FALSE;

  *device = (DeviceInfo *)calloc(1, sizeof(DeviceInfo));
  if (*device == NULL)
    return FALSE;

  obj = json_object_get(json, JSON_ID);

  if (obj==NULL) {
    free(*device);
    return FALSE;
  }

  uuid_parse(json_string_value(obj), (*device)->uuid);

  return TRUE;
}

/* 
 * deserialize_service_info --
 *
 */

static int deserialize_service_info(ServiceInfo **service, json_t *json) {

  json_t *obj;
  
  if (json == NULL && service == NULL)
    return FALSE;

  *service = (ServiceInfo *)calloc(1, sizeof(ServiceInfo));
  if (*service == NULL)
    return FALSE;

  obj = json_object_get(json, JSON_ID);

  if (obj==NULL) {
    free(*service);
    return FALSE;
  }

  uuid_parse(json_string_value(obj), (*service)->uuid);

  return TRUE;
}

/* 
 * deserialize_map_array --
 *
 */
static void deserialize_map_array(UMap **map, json_t *json) {

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


/*
 * deserialize_map --
 *
 */

static void deserialize_map(URequest **request, json_t *json) {

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

/*
 * deserialize_request_info --
 *
 */

static int deserialize_request_info(URequest **request, json_t *json) {

  json_t *obj, *jRaw;
  int size, i;

  if (json == NULL && request == NULL)
    return FALSE;

  *request = (URequest *)calloc(1, sizeof(URequest));
  if (*request == NULL)
    return FALSE;

  /* Initialize inner struct elements. */
  ulfius_init_request(*request);

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

      memcpy((*request)->binary_body, (void *)json_string_value(obj), size);
    }
  }

  return TRUE;
}

/*
 * deserialize_forward_request --
 *
 */

int deserialize_forward_request(MRequest **request, json_t *json) {

  int ret=FALSE;
  json_t *jFwd, *obj;
  char *jStr;

  if (json == NULL) {
    return FALSE;
  }

  jFwd = json_object_get(json, JSON_MESH_FORWARD);
  if (jFwd == NULL) {
    goto fail;
  }

  *request = (MRequest *)calloc(1, sizeof(MRequest));
  if (*request == NULL) {
    return FALSE;
  }

  obj = json_object_get(jFwd, JSON_TYPE);
  if (obj == NULL) {
    goto fail;
  } else {
    (*request)->reqType = strdup(json_string_value(obj));
  }

  obj = json_object_get(jFwd, JSON_SEQ);
  if (obj == NULL) {
    goto fail;
  } else {
    (*request)->seqNo = json_integer_value(obj);
  }

  obj = json_object_get(jFwd, JSON_DEVICE_INFO);
  ret = deserialize_device_info(&(*request)->deviceInfo, obj);

  obj = json_object_get(jFwd, JSON_SERVICE_INFO);
  ret = deserialize_service_info(&(*request)->serviceInfo, obj);

  obj = json_object_get(jFwd, JSON_REQUEST_INFO);
  ret = deserialize_request_info(&(*request)->requestInfo, obj);

  return TRUE;

 fail:
  jStr = json_dumps(json, 0);
  log_error("Error decoding JSON: %s", jStr);
  free(jStr);
  return FALSE;
}

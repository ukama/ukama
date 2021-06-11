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

#include "jserdes.h"

/*
 * serialize_transfer_request -- 
 *
 */

#ifdef SERVER_MOCK
int serialize_transfer_request(json_t *req, Agent *agent) {


}
#endif

/*
 * deserialize_transfer_response --
 *
 */

int deserialize_transfer_response(json_t *resp, AgentCB **agent) {

  int i=0, j=0, count=0;
  json_t *array=NULL;
  json_t *elem=NULL, *method=NULL, *url=NULL;
  
  if (!resp) return FALSE;
  
  array = json_object_get(resp, JSON_AGENT_URL);

  if (json_is_array(array)) {
    count = json_array_size(array);

    agent = (AgentCB **)calloc(sizeof(AgentCB), count);
    
    for (i=0; i<count; i++) {
      elem = json_array_get(array, i);

      if (elem == NULL) {
	goto failure;
      }
      
      method = json_object_get(elem, JSON_AGENT_METHOD);
      url  = json_object_get(elem, JSON_AGENT_URL);

      if (method && url) {
	agent[i]->method = strdup(json_string_value(method));
	agent[i]->url = strdup(json_string_value(url));
      }
    }
  }

  return TRUE;
  
 failure:
  for (j=0; j<i; j++) {
    free(agent[j]->method);
    free(agent[j]->url);
  }

  if (&agent[0]) free(&agent[0]);
  return FALSE;
}

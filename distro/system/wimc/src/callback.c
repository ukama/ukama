/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Callback functions for various endpoints and REST methods.
 */

#include "callback.h"
#include "tasks.h"
#include "jserdes.h"

static void free_agent_request(AgentReq *req);

/*
 * decode a u_map into a string
 */

static char *print_map(const struct _u_map * map) {

  char * line, * to_return = NULL;
  const char **keys, * value;
  
  int len, i;
  
  if (map != NULL) {
    keys = u_map_enum_keys(map);
    for (i=0; keys[i] != NULL; i++) {
      value = u_map_get(map, keys[i]);
      len = snprintf(NULL, 0, "key is %s, value is %s", keys[i], value);
      line = o_malloc((len+1)*sizeof(char));
      snprintf(line, (len+1), "key is %s, value is %s", keys[i], value);
      if (to_return != NULL) {
        len = o_strlen(to_return) + o_strlen(line) + 1;
        to_return = o_realloc(to_return, (len+1)*sizeof(char));
        if (o_strlen(to_return) > 0) {
          strcat(to_return, "\n");
        }
      } else {
        to_return = o_malloc((o_strlen(line) + 1)*sizeof(char));
        to_return[0] = 0;
      }
      strcat(to_return, line);
      o_free(line);
    }
    return to_return;
  } else {
    return NULL;
  }
}

/* 
 * log_request -- log various parameters for the incoming request. 
 *
 */

static void log_request(const struct _u_request *request) {

  log_debug("Recevied: %s %s %s", request->http_protocol, request->http_verb,
	    request->http_url);
}


/*
 * callback_post_container --
 *
 */

int callback_post_container(const struct _u_request *request,
			   struct _u_response *response,
			   void *user_data) {

  int resCode=200;
  int count=0, i=0;
  uuid_t uuid;

  char idStr[36+1] = {0};
  char path[WIMC_MAX_PATH_LEN] = {0};
  char *errStr=NULL;
  char *respBody=NULL, *name=NULL, *tag=NULL;
  WimcCfg *cfg=NULL;
  ServiceURL **urls=NULL;
  WRespType respType=WRESP_ERROR;

  Agent *agent=NULL;
  char *providerURL=NULL, *indexURL=NULL, *storeURL=NULL;
  
  cfg = (WimcCfg *)user_data;
  urls = (ServiceURL **)calloc(sizeof(ServiceURL *), 1);
  uuid_clear(uuid);

  log_request(request);

  name = (char *)u_map_get(request->map_url, "name");
  tag  = (char *)u_map_get(request->map_url, "tag");

  if (!name || !tag) {
    log_error("Invalid name and tag in POST request for EP: %s. Ignoring.",
	      WIMC_EP_CONTAINER);
    errStr = msprintf("Invalid container name and/or tag.");
    resCode = 400;
    goto reply;
  }

  log_debug("Processing container name: %s with tag: %s", name, tag);

  if (db_read_path(cfg->db, name, tag, &path[0])) {
    /* we have the contents stored locally. Return the location. */
    respType = WRESP_RESULT;
    goto reply;
  }

  /* Step-1: Ask service provider for the link. */
  resCode = get_service_url_from_provider(cfg, name, tag, &urls[0], &count);
  if (resCode != 200) {
    resCode = 404;
    errStr = msprintf("No service provider found");
    goto reply;
  }
  
  /* Step-2: find out which register agent can handle the method
   *         returned by the service provider.
   */
  agent = find_matching_agent(*cfg->agents, *urls, count, &providerURL,
			      &indexURL, &storeURL);
  if (agent == NULL) {
    resCode = 404; /*Not found but might be available in furture. */
    errStr = msprintf("No service provider found");
    goto reply;
  } else {
    uuid_unparse(agent->uuid, &idStr[0]);
    log_debug("Matching agent found. Agent Id: %s Method: %s URL: %s",
	      idStr, agent->method, agent->url);
  }

  /* Step-3: Ask agent to fetch the data. Agent will update the status
   *         of this transfer via special status CB url and UID.
   *         Tasks list is also updated.
   */

  resCode = communicate_with_the_agent(WREQ_FETCH, name, tag, providerURL,
				       indexURL, storeURL, agent, cfg, &uuid);

  /* Step-4: setup client URL where they can monitor the status
   *         of the request. 
   */
  if (resCode == 200) { /* Task is assigned to Agent. Return ID. */
    uuid_unparse(uuid, &idStr[0]);
    respType = WRESP_PROCESSING;
    goto reply;
  } else {
    errStr = msprintf("No resource found. Error");
  }

reply:

  respBody = process_cli_response(respType, &path[0], &idStr[0], NULL, errStr);
  if (respBody) {
    ulfius_set_string_body_response(response, resCode, respBody);
  } else {
    ulfius_set_string_body_response(response, resCode, "");
  }
  
  if (respBody)
    free(respBody);

  if (errStr)
    free(errStr);

  for (i=0; i<count; i++) {
    if (urls[i]->method && urls[i]->url) {
      free(urls[i]->method);
      free(urls[i]->url);
      free(urls[i]->iURL);
      free(urls[i]->sURL);
    }
    free(urls[i]);
  }

  free(urls);

  return U_CALLBACK_CONTINUE;
}

/*
 * callback_put_container --
 *
 */
int callback_put_container(const struct _u_request *request,
			   struct _u_response *response,
			   void *user_data) {
  
  char *post_params = print_map(request->map_post_body);
  char *response_body = msprintf("OK!\n%s", post_params);

  ulfius_set_string_body_response(response, 200, response_body);
  o_free(response_body);
  o_free(post_params);

  return U_CALLBACK_CONTINUE;
}

/*
 * callback_delete_container --
 *
 */
int callback_delete_container(const struct _u_request *request,
			      struct _u_response *response,
			      void *user_data) {
  
  char *post_params = print_map(request->map_post_body);
  char *response_body = msprintf("OK!\n%s", post_params);

  ulfius_set_string_body_response(response, 200, response_body);
  o_free(response_body);
  o_free(post_params);

  return U_CALLBACK_CONTINUE;
}

/*
 * callback_get_stats --
 *
 */
int callback_get_stats(const struct _u_request *request,
		       struct _u_response *response,
		       void *user_data) {
  
  char *post_params = print_map(request->map_post_body);
  char *response_body = msprintf("OK!\n%s", post_params);
  
  ulfius_set_string_body_response(response, 200, response_body);
  o_free(response_body);
  o_free(post_params);
  
  return U_CALLBACK_CONTINUE;
}

/*
 * callback_put_agent_update --
 *
 */
int callback_put_agent_update(const struct _u_request *request,
			      struct _u_response *response,
			      void *user_data) {

  int ret=WIMC_OK, retCode;
  uuid_t uuid;
  char *resBody;
  json_t *jreq=NULL;
  json_error_t jerr;
  AgentReq *req=NULL;

  WimcCfg *cfg = (WimcCfg *)user_data;

  req = (AgentReq *)calloc(sizeof(AgentReq), 1);

  jreq = ulfius_get_json_body_request(request, &jerr);
  if (!jreq) {
    log_error("json error: %s", jerr.text);
  } else {
    deserialize_agent_request(&req, jreq);
  }

  ret = process_agent_update_request(cfg->tasks, req, &uuid, cfg->db);

  if (ret == WIMC_OK) {
    retCode = 200;
  } else if (ret == WIMC_ERROR_BAD_ID){
    retCode = 404;
  } else {
    retCode = 400;
  }

  resBody = msprintf("%s\n", error_to_str(ret));
  ulfius_set_string_body_response(response, retCode, resBody);
  o_free(resBody);
  free_agent_request(req);
  json_decref(jreq);

  return U_CALLBACK_CONTINUE;
}

static void free_agent_request(AgentReq *req) {

  if (req->type == REQ_REG) {
    free(req->reg->method);
    free(req->reg->url);
    free(req->reg);
  } else if (req->type == REQ_UPDATE) {
    if (req->update->voidStr) free(req->update->voidStr);
    free(req->update);
  }

  free(req);
}

/*
 * callback_post_agent --
 *
 */
int callback_post_agent(const struct _u_request *request,
			struct _u_response *response,
			void *user_data) {
  int ret=WIMC_OK, retCode;
  uuid_t uuid;
  char *resBody;
  json_t *jreq=NULL;
  json_error_t jerr;
  AgentReq *req=NULL;
  char idStr[36+1];

  WimcCfg *cfg = (WimcCfg *)user_data;

  req = (AgentReq *)calloc(sizeof(AgentReq), 1);

  jreq = ulfius_get_json_body_request(request, &jerr);
  if (!jreq) {
    log_error("json error: %s", jerr.text);
  } else {
    deserialize_agent_request(&req, jreq);
  }

  ret = process_agent_register_request(cfg->agents, req, &uuid);

  if (ret == WIMC_OK) {
    retCode = 200;
    uuid_unparse(uuid, &idStr[0]);
    resBody = msprintf("%s\n", idStr);
  } else {
    retCode = 400;
    resBody = msprintf("%s\n", error_to_str(ret));
  }

  ulfius_set_string_body_response(response, retCode, resBody);
  o_free(resBody);
  free_agent_request(req);
  json_decref(jreq);

  return U_CALLBACK_CONTINUE;
}

/*
 * callback_get_task --
 *
 */
int callback_get_task(const struct _u_request *request,
		      struct _u_response *response, void *user_data) {

  int statusCode=200;
  uuid_t uuid;
  WimcCfg *cfg=NULL;
  WTasks *task=NULL;
  char *idStr=NULL, *resBody=NULL, *errStr=NULL;
  WRespType respType= WRESP_ERROR;

  cfg = (WimcCfg *)user_data;
  uuid_clear(uuid);

  idStr = (char *)u_map_get(request->map_url, "id");
  if (!idStr) {
    statusCode = 400;
    errStr = msprintf("%s", WIMC_ERROR_BAD_ID_STR);
    goto done;
  }

  if (uuid_parse(idStr, uuid)==-1) {
    log_error("Error parsing the UUID into binary: %s", idStr);
    statusCode = 400;
    errStr = msprintf("%s", WIMC_ERROR_BAD_ID_STR);
    goto done;
  }

  /* Find matching task. */
  task = *(cfg->tasks);
  while (task != NULL) {
    if (uuid_compare(task->uuid, uuid) == 0) {
      break;
    }
    task = task->next;
  }

  if (!task) {
    log_debug("No matching task found: %s", idStr);
    statusCode = 400;
    errStr = msprintf("%s", WIMC_ERROR_BAD_ID_STR);
    goto done;
  }

  /* Set proper response type. */
  switch(task->state) {
  case REQUEST:
  case FETCH:
  case UNPACK:
    respType = WRESP_PROCESSING;
    break;
  case DONE:
    respType = WRESP_RESULT;
    break;
  case ERR:
    respType = WRESP_ERROR;
    break;
  default:
    respType = WRESP_PROCESSING;
    break;
  }
  statusCode = 200;

 done:
  /* DELETE method require additional processing. */
  if (statusCode == 200) {
    if (strcasecmp(request->http_verb, "DELETE")==0) {
      /* Delete the task from the list. */
      delete_from_tasks(cfg->tasks, task);
    }
  }

  resBody = process_cli_response(respType, NULL, idStr, task, errStr);

  ulfius_set_string_body_response(response, statusCode, resBody);
  o_free(resBody);
  o_free(errStr);

  return U_CALLBACK_CONTINUE;
}

/*
 * callback_not_allowed -- 
 *
 */
int callback_not_allowed(const struct _u_request *request,
			 struct _u_response *response, void *user_data) {

  ulfius_set_string_body_response(response, 403, "Operation not allowed\n");
  return U_CALLBACK_CONTINUE;
}

/*
 * callback_default -- default callback for no-match
 *
 */
int callback_default(const struct _u_request *request,
                     struct _u_response *response, void *user_data) {

  ulfius_set_string_body_response(response, 404, "You are clearly high!\n");
  return U_CALLBACK_CONTINUE;
}

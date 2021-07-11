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
 * validate_path --
 *
 */

static int validate_path(char *path) {

  int val=FALSE;
  struct stat sb;
  char buf[256];

  sprintf(buf, "%s/config.json", path); 
  
  if (stat(path, &sb) == 0 && S_ISDIR(sb.st_mode)) {
    /* Now check for config.json file. */
    if (access(&buf[0], F_OK) == 0) {
      log_debug("Config.json found at: %s", path);
      val = TRUE;
    } else {
      log_debug("Valid path but config.json NOT found!");
    }      
  } else {
    log_debug("Invalid path: %s", path);
  }
  
  return val;
}

/*
 * container_name_and_tag -- 
 *
 */
static int container_name_and_tag(char *str, char *name, char *tag) { 

  char *c, *t;

  c = strtok(str, ":");
  t = strtok(NULL, ":");

  if (c == NULL || t == NULL) {
    goto failure;
  }

  /* sanity check. */
  if (strlen(c) > WIMC_MAX_NAME_LEN ||
      strlen(t) > WIMC_MAX_TAG_LEN) {
    goto failure;
  }

  strncpy(name, c, strlen(c));
  strncpy(tag, t, strlen(t));

  return TRUE;
 failure:
  return FALSE;
}

/* 
 * get_key_value -- extract key value pair from the passed string. It is in
 *                  format key=value. If keyName is valid, also compared it 
 *                  against it.
 *
 */

static int get_key_value(char *str, char *key, char *value, char *keyName,
			 int maxKeyLen) {

  int val = FALSE;
  char *k, *v;
  
  k = strtok(str, "=");
  v = strtok(NULL, "=");
  
  if (k == NULL || v == NULL) {
    goto failure;
  }
  
  if ( keyName != NULL && maxKeyLen > 0) { /* Validate. */
    if ((strcasecmp(keyName, k)==0) || (strlen(v) < maxKeyLen)) {
      val = TRUE;
    }
  }
  
  /* We still copy even if the keyname doesn't match so we can log this. */
  strncpy(value, v, strlen(v));
  strncpy(key, k, strlen(k));
  
 failure:
  return val;
}

/*
 * get_post_params -- get params for the POST method.
 *                    param is of format name=container:tag&path=/some/path
 *
 */

static int get_post_params(char *params, char *name, char *tag, char *path) {

  int val=FALSE;
  char *token1, *token2;
  char key[256]={0}, value[256]={0};

  /* Extract the first token: container name and tag. */
  token1 = strtok(params, "&");
  token2 = strtok(NULL, "&"); 
  
  if (token1 == NULL || token2 == NULL) {
    goto failure;
  }
  
  val = get_key_value(token1, &key[0], &value[0], WIMC_PARAM_CONTAINER_NAME,
		      WIMC_MAX_NAME_LEN + WIMC_MAX_PATH_LEN + 1);
  
  if (val == TRUE) {
    /* Extract the container name and its tag. It is in format
     * container_name:tag 
     */
    val = container_name_and_tag(value, name, tag);
    if (!val) {
      log_debug("Invalid container name and/or tag. Ignoring");
      goto failure;
    }
  } else {
    log_debug("Invalid key %s. Expected %s", key, WIMC_PARAM_CONTAINER_NAME);
    goto failure;
  }
  
  /* Process the path. */
  memset(key, 0, sizeof(char)*256);
  memset(value, 0, sizeof(char)*256);
  
  val = get_key_value(token2, &key[0], &value[0], WIMC_PARAM_CONTAINER_PATH,
		      WIMC_MAX_PATH_LEN);
  if (val == TRUE){
    strncpy(path, value, strlen(value));
  } else {
    goto failure;
  }

  log_debug("Valid POST params recevied. Name: %s Tag: %s Path: %s",
	    name, tag, path);

 failure:
  return val;
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

  int val=FALSE, found=FALSE;
  int resCode=200;
  int count=0, i=0, id=0;
  
  char path[WIMC_MAX_PATH_LEN] = {0};
  char *post_params=NULL, *response_body=NULL;
  char *name=NULL, *tag=NULL;
  WimcCfg *cfg=NULL;
  ServiceURL **urls=NULL;

  Agent *agent=NULL;
  char *providerURL=NULL;
  
  cfg = (WimcCfg *)user_data;
  urls = (ServiceURL **)calloc(sizeof(ServiceURL *), 1);

  log_request(request);

  name = (char *)u_map_get(request->map_url, "name");
  tag  = (char *)u_map_get(request->map_url, "tag");

  if (!name || !tag) {
    log_error("Invalid name and tag in POST request for EP: %s. Ignoring.",
	      WIMC_EP_CONTAINER);
    /* XXX Send error on bad name. */
    response_body = msprintf("Invalid container name and/or tag.");
    resCode = 400;
    goto reply;
  }

  log_debug("Processing container name: %s and tag: %s", name, tag);
  
  if (db_read_path(cfg->db, name, tag, &path[0])) {
      /* we have the contents stored locally. Return the location. */
    response_body = msprintf("%s", path);
    goto reply;
  }

  /* Step-1: Ask service provider for the link. */
  resCode = get_service_url_from_provider(cfg, name, tag, &urls[0], &count);
  if (resCode != 200) {
    resCode = 404;
    response_body = msprintf("No service provider found");
    goto reply;
  }
  
  /* Step-2: find out which register agent can handle the method
   *         returned by the service provider.
   */
  agent = find_matching_agent(*cfg->agents, *urls, count, &providerURL);
  if (agent == NULL) {
    resCode = 404; /*Not found but might be available in furture. */
    response_body = msprintf("No service provider found");
    goto reply;
  } else {
    log_debug("Matching agent found. Agent Id: %d Method: %s URL: %s",
	      agent->id, agent->method, agent->url);
  }

  /* Step-3: Ask agent to fetch the data. Agent will update the status
   *         of this transfer via special status CB url and UID.
   *         Tasks list is also updated.
   */

  resCode = communicate_with_the_agent(WREQ_FETCH, name, tag, providerURL,
				       agent, cfg, &id);

  /* Step-4: setup client URL where they can monitor the status
   *         of the request. 
   */
  if (resCode == 200) { /* Task is assigned to Agent. Return ID. */
    response_body = msprintf("%d", id);
  } else {
    response_body = msprintf("No resource found. Error");
  }

reply:

  ulfius_set_string_body_response(response, resCode, response_body);
  
  o_free(response_body);
  o_free(post_params);

  for (i=0; i<count; i++) {
    if (urls[i]->method && urls[i]->url) {
      free(urls[i]->method);
      free(urls[i]->url);
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

static void free_agent_request(AgentReq *req) {

  if (req->type == REQ_REG) {
    free(req->reg->method);
    free(req->reg->url);
    free(req->reg);
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
  int ret=WIMC_OK, retCode, id=0;
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

  ret = process_agent_request(cfg->agents, req, &id);
  
  if (ret == WIMC_OK) {
    retCode = 200;
    resBody = msprintf("%d\n", id);
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

  int ret, statusCode=200;
  long id=0;
  WimcCfg *cfg=NULL;
  WTasks *task=NULL;
  char *idStr=NULL, *resBody=NULL;

  cfg = (WimcCfg *)user_data;

  idStr = (char *)u_map_get(request->map_url, "id");
  if (!idStr) {
    statusCode = 400;
    resBody = msprintf("%s", WIMC_ERROR_BAD_ID_STR);
    goto done;
  }

  id = atol(idStr);

  /* Find matching task. */
  task = *(cfg->tasks);
  while (task != NULL) {
    if (task->id == id) {
      break;
    }
    task = task->next;
  }

  if (task) {
    char *str;
    /* if ID found, serialize the tasks status and send it back, yo. */
    str = process_task_request(task);
    if (!str) {
      statusCode = 400;
      resBody = msprintf("%s", WIMC_ERROR_MEMORY_STR);
    } else {
      statusCode = 200;
      resBody = msprintf("%s", str);
      free(str);
    }
  } else {
    statusCode = 400;
    resBody = msprintf("%s", WIMC_ERROR_BAD_ID_STR);
  }

 done:
  ulfius_set_string_body_response(response, statusCode, resBody);
  o_free(resBody);

  return U_CALLBACK_CONTINUE;
}

/*
 * callback_delete_task --
 *
 */
int callback_delete_task(const struct _u_request *request,
			 struct _u_response *response, void *user_data) {

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

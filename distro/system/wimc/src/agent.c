/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Agent related functions.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <curl/curl.h>
#include <curl/easy.h>

#include "agent.h"
#include "wimc.h"
#include "tasks.h"
#include "jserdes.h"
#include "common/utils.h"

struct Response {
  char *buffer;
  size_t size;
};

static char *create_cb_url_for_agent(char *port);

/*
 * register_agent -- register new agent
 */

int register_agent(Agent **agents, char *method, char *url, uuid_t *uuid) {

  int i;
  char idStr[36+1];
  Agent *ptr = *agents;

  for (i=0; i < MAX_AGENTS; i++) {

    if (uuid_is_null(ptr[i].uuid)==0) { /* have valid agent id. */
      if (strcasecmp(method, ptr[i].method)==0 &&
	  strcasecmp(url, ptr[i].url)==0) {
	uuid_unparse(ptr[i].uuid, &idStr[0]);
	uuid_copy(*uuid, ptr[i].uuid);
	/* An existing entry. */
	log_debug("Found similar agent at id: %s, method %s and url: %s",
		  idStr, ptr[i].method, ptr[i].url);
	return WIMC_ERROR_EXIST;
      }
    } else {
      uuid_generate(ptr[i].uuid);
      ptr[i].method = strndup(method, strlen(method));
      ptr[i].url    = strndup(url, strlen(url));
      ptr[i].state  = WIMC_AGENT_STATE_REGISTER;

      /* Return the ID. */
      uuid_copy(*uuid, ptr[i].uuid);
      return WIMC_OK;
    }
  }

  /* Max. reached */
  log_debug("Max. allowable number of agents reached. Ignoring");
  return WIMC_ERROR_MAX_AGENTS;
}

/*
 * process_agent_request --
 *
 */

int process_agent_request(Agent **agents, AgentReq *req, uuid_t *uuid){

  int ret=WIMC_OK;
  Register *reg;
  char idStr[36+1];
  
  if (req->type == (ReqType)REQ_REG) {

    reg = req->reg;
    
    /* validate the URL. */
    ret = validate_url(reg->url);
    if (ret != WIMC_OK) {
      log_debug("Agent process failed, unreachable URL: %s: %s", reg->url,
		error_to_str(ret));
      goto done;
    }
    
    ret = register_agent(agents, reg->method, reg->url, uuid);
    if (ret != WIMC_OK) {
      goto done;
    }

    uuid_unparse(*uuid, &idStr[0]);
    log_debug("Agent successfully registered. Id: %s Method: %s URL: %s",
	      idStr, reg->method, reg->url);
  } else if (req->type == (ReqType)REQ_UNREG) {
    
  } else if (req->type == (ReqType)REQ_UPDATE) {

  } else {
    log_debug("Invalid Agent request command: %d", req->type);
    ret = WIMC_ERROR_BAD_METHOD;
    goto done;
  }
  
 done:
    return ret;
}

/*
 * find_matching_agent -- return the Agent which matches the given
 *                        method (as returned from service provider). 
 *                        If there are multiple URL in the list, we always 
 *                        send the first match between the agent and provider.
 *                        We can be more intelligent in future.
 */
Agent *find_matching_agent(Agent *agents, void *vURLs, int count,
			   char **providerURL, char **indexURL,
			   char **storeURL) {

  int i, j;
  Agent *ptr = agents;
  ServiceURL *urls = (ServiceURL *)vURLs;

  /* Sanity check. */
  if (agents == NULL)
    return NULL;

  for (i=0; i < MAX_AGENTS; i++) {

    if (uuid_is_null(ptr[i].uuid)==0) { /* have valid agent id. */
      for (j=0; j < count; j++) {
	if (strcasecmp(urls[i].method, ptr[i].method)==0) {
	  *providerURL = urls[i].url;
	  *indexURL    = urls[i].iURL;
	  *storeURL    = urls[i].sURL;
	  return ptr;
	}
      }
    }
  }

  /* NULL if no match found. */
  return NULL;
}

/*
 * cleanup_wimc_request --
 *
 */
static void cleanup_wimc_request(WimcReq *request) {

  if (request->type == (WReqType)WREQ_FETCH) {
    WFetch *fetch = request->fetch;
    WContent *content = fetch->content;

    if (content) {
      free(content->name);
      free(content->tag);
      free(content->method);
      free(content->providerURL);
      free(content->indexURL);
      free(content->storeURL);
      free(content);
    }

    if (fetch) {
      free(fetch->cbURL);
      free(fetch);
    }
  }

  free(request);
}


/*
 * create_cb_url_for_agent --
 *
 */
static char *create_cb_url_for_agent(char *port) {

  char *cbURL = NULL;

  if (port==NULL) {
    return NULL;
  }

  cbURL = (char *)malloc(WIMC_MAX_URL_LEN);
  if (cbURL) {
    sprintf(cbURL, "http://localhost:%s/%s", port, WIMC_EP_AGENT_UPDATE);
  }

  return cbURL;
}

/*
 * response_callback --
 */
static size_t response_callback(void *contents, size_t size, size_t nmemb,
				void *userp) {

  size_t realsize = size * nmemb;
  struct Response *response = (struct Response *)userp;

  response->buffer = realloc(response->buffer, response->size + realsize + 1);

  if(response->buffer == NULL) {
    log_error("Not enough memory to realloc of size: %s",
	      response->size + realsize + 1);
    return 0;
  }

  memcpy(&(response->buffer[response->size]), contents, realsize);
  response->size += realsize;
  response->buffer[response->size] = 0; /* Null terminate. */

  return realsize;
}

/*
 * create_wimc_request --
 *
 */

static WimcReq *create_wimc_request(WReqType reqType, char *name, char *tag,
				    char *providerURL, char *cbURL,
				    char *iURL, char *sURL,
				    char *method, int interval) {

  WimcReq *request=NULL;
  WFetch  *fetch=NULL;
  WContent *content=NULL;

  request = (WimcReq *)calloc(1, sizeof(WimcReq));
  if (request==NULL) {
    goto done;
  }

  if (reqType == (WReqType)WREQ_FETCH) { /* Request to fetch contents. */

    fetch   = (WFetch *)malloc(sizeof(WFetch));
    content = (WContent *)malloc(sizeof(WContent));

    if (!fetch && !content) {
      goto done;
    }

    request->type = WREQ_FETCH;

    uuid_generate(fetch->uuid);
    fetch->cbURL = strdup(cbURL);
    fetch->interval = interval;

    content->name = strdup(name);
    content->tag  = strdup(tag);
    content->providerURL = strdup(providerURL);
    content->method = strdup(method);
    content->indexURL = strdup(iURL);
    content->storeURL = strdup(sURL);

    fetch->content = content;
    request->fetch = fetch;

  } else if (reqType == (WReqType)WREQ_UPDATE) { /* update an existing req */

  }

  return request;

 done:
  if (content) {
    free(content->name);
    free(content->tag);
    free(content->providerURL);
    free(content->indexURL);
    free(content->storeURL);
    free(content);
  }

  if (fetch) {
    free(fetch->cbURL);
    free(fetch);
  }

  if (request) {
    free(request);
  }

  return NULL;
}

/*
 * send_request_to_agent --
 *
 */
static long send_request_to_agent(WReqType reqType, char *agentURL,
				  json_t *json, int *retCode) {

  long code=0;
  CURL *curl=NULL;
  CURLcode res;
  char *json_str;
  struct curl_slist *headers=NULL;
  struct Response response;

  *retCode = 0;

  /* sanity check */
  if (json==NULL && agentURL==NULL) {
    return code;
  }
  
  curl_global_init(CURL_GLOBAL_ALL);
  curl = curl_easy_init();
  if (curl == NULL) {
    return code;
  }

  response.buffer = malloc(1);
  response.size   = 0;
  json_str = json_dumps(json, 0);

  /* Add to the header. */
  headers = curl_slist_append(headers, "Accept: application/json");
  headers = curl_slist_append(headers, "Content-Type: application/json");
  headers = curl_slist_append(headers, "charset: utf-8");

  curl_easy_setopt(curl, CURLOPT_URL, agentURL);

  curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "POST");
  curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
  curl_easy_setopt(curl, CURLOPT_POSTFIELDS, json_str);
  curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
  curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);

  curl_easy_setopt(curl, CURLOPT_USERAGENT, "wimc/0.1");

  res = curl_easy_perform(curl);

  if (res != CURLE_OK) {
    log_error("Error sending request to Agent: %s", curl_easy_strerror(res));
  } else {
    /* get status code. */
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
    //process_response_from_wimc(reqType, code, &response, id);
  }

  free(json_str);
  free(response.buffer);
  curl_slist_free_all(headers);
  curl_easy_cleanup(curl);
  curl_global_cleanup();

  return code;
}

/*
 * communicate_with_agent -- Function WIMC.d uses to communicate with Agent(s).
 *
 */
long communicate_with_the_agent(WReqType reqType, char *name, char *tag,
				char *providerURL, char *indexURL,
				char *storeURL, Agent *agent, WimcCfg *cfg,
				uuid_t *uuid) {

  /* steps are:
   * 0. Generate unique ID for the content request and CB url.
   * 1. create wimc request for agent.
   * 2. serialize the request as JSON object.
   * 3. send the request to the Agent using provided URL.
   * 4. Update the UUID and return HTTP status code.
   */

  int ret=FALSE;
  long code=0;
  char *cbURL=NULL;
  WimcReq *request=NULL;
  json_t *json=NULL;
  int agentRetCode=0;

  /* Some sanity check. */
  if (!agent && !cfg) {
    return code;
  }

  if (reqType == (WReqType)WREQ_FETCH) {
    if (!indexURL && !storeURL && !providerURL) {
      return code;
    }

    if (!name && !tag) {
      return code;
    }
  }

  cbURL = create_cb_url_for_agent(cfg->adminPort);
  if (!cbURL) {
    goto done;
  }

  request = create_wimc_request(reqType, name, tag, providerURL, cbURL,
				indexURL, storeURL, agent->method,
				DEFAULT_INTERVAL);
  if (!request) {
    goto done;
  }

  ret = serialize_wimc_request(request, &json);
  if (!ret) {
    goto done;
  }

  code = send_request_to_agent(reqType, agent->url, json, &agentRetCode);
  if (code == 200) {
    log_debug("Agent command success. CURL return code: %d Agent code: %d",
	      code, agentRetCode);
  } else {
    log_debug("Agent command success. CURL return code: %d Agent code: %d",
	      code, agentRetCode);
  }

  if (code == 200) { /* Add to tasks list. */
    add_to_tasks(cfg->tasks, request);
    uuid_copy(*uuid, request->fetch->uuid);
  }

 done:
  json_decref(json);
  cleanup_wimc_request(request);
  if (cbURL) {
    free(cbURL);
  }

  return code;
}

void clear_agents(Agent *agent) {

  int i;

  for (i=0; i<MAX_AGENTS; i++){
    if (uuid_is_null(agent[i].uuid)==0) {
      free(agent[i].method);
      free(agent[i].url);
    }
  }
}

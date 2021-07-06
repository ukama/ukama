/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Functions to interact with the service provider in the cloud.
 *
 */

#include <sqlite3.h>
#include <jansson.h>
#include <ulfius.h>
#include <curl/curl.h>

#include "wimc.h"
#include "log.h"
#include "provider.h"
#include "jserdes.h"


struct Response {
  char *buffer;
  size_t size;
};

/* Function def. */
static void log_request(req_t *req);
static void log_response(resp_t *resp);
static void process_response_from_provider(WimcCfg *cfg, long statusCode,
					   void *resp);
static void create_provider_url(WimcCfg *cfg, char *url, char *name,
				char *tag);
static size_t response_callback(void *contents, size_t size, size_t nmemb,
				void *userp);
int fetch_content_from_service_provider(WimcCfg *cfg, char *name, char *tag);

/*
 * log_request --
 *
 */
static void log_request(req_t *req) {

}

static void log_response(resp_t *resp) {

}

/*
 * process_response_from_provider --
 *
 */

static void process_response_from_provider(WimcCfg *cfg, long statusCode,
					   void *resp) {

  struct Response *response;
  json_t *json;
  int count=0, i=0, ret=FALSE;
  AgentCB **agents=NULL;
  
  response = (struct Response *)resp;
  agents = (AgentCB **)calloc(sizeof(AgentCB *), 1);
  
  json = json_loads(response->buffer, JSON_DECODE_ANY, NULL);

  if (!json) {
    log_error("Can not load str into JSON object. Str: %s", response->buffer);
    goto done;
  }

  ret = deserialize_provider_response(&agents[0], &count, json);

  if (ret==FALSE) {
    log_error("Deserialization failed for %s", response->buffer);
    goto done;
  }

  if (count==0) {
    log_debug("No matching agent provided by service provider");
    goto done;
  }
  
  for (i=0; i<count; i++) {
    log_debug("Received Agent CB URL from provider. %d: Method: %s URL: %s",
	      i, agents[i]->method, agents[i]->url);

#if 0
	/* lookup registered matching Agents*/
	agentURL = find_matching_agent_url(agent[i]->method, cfg->agents);
#endif
  }

 done:
  
  for (i=0; i<count; i++) {
    if (agents[i]->method && agents[i]->url) {
      free(agents[i]->method);
      free(agents[i]->url);
    }
    free(agents[i]);
  }
  
  free(agents);

  json_decref(json);
  return;
}

/* 
 * create_provider_url -- create EP for the provider.
 *
 */
static void create_provider_url(WimcCfg *cfg, char *url, char *name,
				char *tag) {

  sprintf(url, "%s/%s?name=%s&tag=%s", cfg->cloud, WIMC_EP_PROVIDER, name, tag);

  return;
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
 * send_get_request_to_provider -- 
 *
 */

static int send_request_to_provider(WimcCfg *cfg, char *name, char *tag) {
  
  int ret;
  long code=0;
  char providerEP[WIMC_MAX_URL_LEN] = {0};
  CURL *curl=NULL;
  CURLcode res;
  struct curl_slist *headers=NULL;
  struct Response response;
  
  /* Sanity check. */
  if (!name || !tag) {
    goto done;
  }
     
  create_provider_url(cfg, &providerEP, name, tag);

  curl = curl_easy_init();
  if (curl == NULL) {
    return code;
  }

  response.buffer = (char *)malloc(1);
  response.size   = 0;
  
  curl_easy_setopt(curl, CURLOPT_URL, providerEP);
  
  curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "GET");
  curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
  curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);

  curl_easy_setopt(curl, CURLOPT_USERAGENT, "wimc.d/0.1");
  
  res = curl_easy_perform(curl);

  if (res != CURLE_OK) {
    log_error("Error sending request to provider: %s", curl_easy_strerror(res));
  } else {
    /* get status code. */
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
    process_response_from_provider(cfg, code, &response);
  }

 done:

  free(response.buffer);
  curl_easy_cleanup(curl);
  
  return code;
}

/*
 * fetch_container_from_service_provider -- 
 *
 */ 

int fetch_content_from_service_provider(WimcCfg *cfg, char *name, char *tag) {

  /* Logic is as follows:
   * 1. Issue GET command to the cloud-based service provider for name:tag
   * 2. Provider will either:
   *    a. reject the request with 404 or
   *    b. accept and return remote_cb URL. 
   * 3. Remote_cb URL is then passed to the Agent, along with status_CB.
   * 4. status_cb keep the db updated for the content.
   */

  send_request_to_provider(cfg, name, tag); 

}

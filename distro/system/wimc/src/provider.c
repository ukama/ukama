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
#include "methods.h"


struct Response {
  char *buffer;
  size_t size;
};

/* Function def. */
static int process_response_from_provider(WimcCfg *cfg, long statusCode,
					   void *resp, ServiceURL **urls);
static void create_provider_url(WimcCfg *cfg, char *url, char *name,
				char *tag);
static size_t response_callback(void *contents, size_t size, size_t nmemb,
				void *userp);
int get_service_url_from_provider(WimcCfg *cfg, char *name, char *tag,
				  ServiceURL **urls, int *count);
/*
 * process_response_from_provider --
 *
 */
static int process_response_from_provider(WimcCfg *cfg, long statusCode,
					   void *resp, ServiceURL **urls) {

  struct Response *response;
  json_t *json;
  int count=0, i=0, ret=FALSE;
  
  response = (struct Response *)resp;
  
  json = json_loads(response->buffer, JSON_DECODE_ANY, NULL);

  if (!json) {
    log_error("Can not load str into JSON object. Str: %s", response->buffer);
    goto done;
  }

  ret = deserialize_provider_response(urls, &count, json);

  if (ret==FALSE) {
    log_error("Deserialization failed for %s", response->buffer);
    goto done;
  }

  if (count==0) {
    log_debug("No matching agent provided by service provider");
    goto done;
  }
  
  for (i=0; i<count; i++) {
    log_debug("Received URLs from provider. %d: Method: %s URL: %s index: %s store: %s",
	      i, urls[i]->method, urls[i]->url, urls[i]->iURL, urls[i]->sURL);
  }

 done:

  /*
  for (i=0; i<count; i++) {
    if (urls[i]->method && urls[i]->url) {
      free(urls[i]->method);
      free(urls[i]->url);
    }
    free(urls[i]);
  }

  free(urls);
  */

  json_decref(json);
  return count;
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
 * get_service_url_from_provider -- 
 *
 */
int get_service_url_from_provider(WimcCfg *cfg, char *name, char *tag,
				  ServiceURL **urls, int *count) {
  
  int ret;
  long code=0;
  char providerEP[WIMC_MAX_URL_LEN] = {0};
  CURL *curl=NULL;
  CURLcode res;
  struct Response response;
  
  /* Sanity check. */
  if (!name || !tag) {
    goto done;
  }
     
  create_provider_url(cfg, &providerEP[0], name, tag);

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
    *count = process_response_from_provider(cfg, code, &response, urls);
  }

 done:

  free(response.buffer);
  curl_easy_cleanup(curl);
  
  return code;
}

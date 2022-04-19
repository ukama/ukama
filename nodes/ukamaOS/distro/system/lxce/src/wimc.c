/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Functions to interact with WIMC.d
 *
 */

#include <jansson.h>
#include <ulfius.h>
#include <curl/curl.h>
#include <string.h>

#include "log.h"
#include "wimc.h"
#include "lxce_config.h"

static int deserialzie_wimc_response(json_t *json, char *path);
static int process_response_from_wimc(Config *config, long statusCode,
				      void *resp, char *path);
static void create_wimc_url(char *url, char *host, char *port, char *name,
			    char *tag);
static size_t response_callback(void *contents, size_t size, size_t nmemb,
				void *userp);
static int get_capp_path_from_wimc(Config *config, char *name, char *tag,
				   char *path);
static int get_capp_path_from_wimc(Config *config, char *name, char *tag,
				   char *path);
/*
 * deserialize_wimc_response --
 *
 */
static int deserialzie_wimc_response(json_t *json, char *path) {

  int ret=FALSE;
  char *type, *result;
  json_t *jResp, *obj;

  if (!json) return FALSE;

  jResp = json_object_get(json, JSON_WIMC_RESPONSE);
  if (jResp==NULL) {
    return FALSE;
  }

  obj = json_object_get(jResp, JSON_TYPE);
  if (obj==NULL) {
    log_error("Missing response type");
    return FALSE;
  }
  type = json_string_value(obj);

  obj = json_object_get(jResp, JSON_VOID_STR);
  if (obj==NULL) {
    log_error("Missing str response.");
    return FALSE;
  }
  result = json_string_value(obj);

  if (strcmp(type, WIMC_RESP_TYPE_RESULT)==0) {
    path = strdup(result);
    ret  = TRUE;
  } else if (strcmp(type, WIMC_RESP_TYPE_ERROR)==0) {
    path = NULL;
    ret  = FALSE;
    log_error("WIMC responded with an error: %s", result);
  } else if (strcmp(type, WIMC_RESP_TYPE_PROCESSING)==0) {
    path = NULL;
    ret  = FALSE;
    log_error("WIMC is processing the request."); /*XXX - handle it properly */
  }

  return ret;
}

/*
 * process_response_from_wimc --
 *
 */
static int process_response_from_wimc(Config *config, long statusCode,
				      void *resp, char *path) {

  struct Response *response;
  json_t *json;
  int count=0, i=0, ret=FALSE;

  response = (struct Response *)resp;

  if (statusCode!=200) {
    log_error("Wimc return code: %d", statusCode);
    path = NULL;
    return FALSE;
  }

  json = json_loads(response->buffer, JSON_DECODE_ANY, NULL);

  if (!json) {
    log_error("Can not load str into JSON object. Str: %s", response->buffer);
    path = NULL;
    return FALSE;
  }

  ret = deserialize_wimc_response(json, path);

  if (ret==FALSE) {
    log_error("Deserialization failed for %s", response->buffer);
  }

  return ret;
}

/*
 * create_wimc_url -- create EP for WIMC
 *
 */
static void create_wimc_url(char *url, char *host, char *port, char *name,
			    char *tag) {

  sprintf(url, "%s:%s/%s?name=%s&tag=%s", host, port, WIMC_EP, name, tag);

  log_debug("Wimc request url: %s", url);

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
 * get_capp_path_from_wimc --
 *
 */
static int get_capp_path_from_wimc(Config *config, char *name, char *tag,
				   char *path) {

  int ret;
  long code=0;
  char wimcEP[LXCE_MAX_URL_LEN] = {0};
  CURL *curl=NULL;
  CURLcode res;
  struct Response response;

  /* Sanity check. */
  if (!config) return FALSE;

  if (!config->wimcHost || !config->wimcPort || !name || !tag) {
    goto done;
  }

  create_wimc_url(&wimcEP[0], config->wimcHost, config->wimcPort, name, tag);

  curl = curl_easy_init();
  if (curl == NULL) {
    return code;
  }

  response.buffer = (char *)malloc(1);
  response.size   = 0;

  curl_easy_setopt(curl, CURLOPT_URL, wimcEP);

  curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "POST");
  curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
  curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);

  curl_easy_setopt(curl, CURLOPT_USERAGENT, "lxce/0.1");

  res = curl_easy_perform(curl);

  if (res != CURLE_OK) {
    log_error("Error sending request to wimc: %s", curl_easy_strerror(res));
  } else {
    /* get status code. */
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
    ret = process_response_from_wimc(config, code, &response, path);

    if (ret) {
      log_debug("Wimc respone for name: %s tag: %s path is: %s", name, tag,
		path);
    } else {
      log_error("Failed to get response for name: %s tag: %s. Code: %s",
		name, tag, code);
    }
  }

 done:
  free(response.buffer);
  curl_easy_cleanup(curl);

  return code;
}

/*
 * get_capp_path -- location of the capp referred by name:tag
 *
 */
int get_capp_path(Config *config, char *name, char *tag, char *path) {

  int code;

  /* sanity check */
  if (config == NULL || name == NULL || tag == NULL)
    return FALSE;

  if (config->wimcHost == NULL || config->wimcPort == NULL)
    return FALSE;

  code = get_capp_path_from_wimc(config, name, tag, path);

  if (code!=200 || !path)
    return FALSE;

  return TRUE;
}

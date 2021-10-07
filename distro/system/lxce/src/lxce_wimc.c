/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <jansson.h>
#include <curl/curl.h>
#include <string.h>

#include "lxce_wimc.h"
#include "jserdes.h"
#include "log.h"

struct Response {
  char *buffer;
  size_t size;
};

/*
 * is_valid_path --
 *
 */
static int is_valid_path(char *path) {

  return TRUE;
}

/*
 * process_response_from_wimc -- 
 *
 */
static int process_response_from_wimc(long statusCode, void *resp,
				      char *wimcResp) {

  struct Response *response;
  json_t *json;
  int ret=FALSE;
 
  response = (struct Response *)resp;

  json = json_loads(response->buffer, JSON_DECODE_ANY, NULL);

  if (!json) {
    log_error("Can not load str into JSON object. Str: %s", response->buffer);
    goto done;
  }

  ret = deserialize_wimc_response(wimcResp, json);
  if (ret==FALSE-1) {
    log_error("Deserialization failed for %s", response->buffer);
    goto done;
  } else if (ret==FALSE) {
    log_error("WIMC response with error: %s", wimcResp);
  } else if (ret==TRUE+1) {
    log_error("No path found for the container or is processing");
  } else if (ret==TRUE) {
    log_debug("Received content path from WIMC. Path: %s", wimcResp);
  }
  
 done:
  json_decref(json);
  return ret;
}

/*
 * response_callback --
 *
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
 * create_wimc_url -- create EP for the wimc.
 *
 */
static void create_wimc_url(char *url, char *name, char *tag,
				char *host, char *port) {

  sprintf(url, "http://%s:%s/%s?name=%s&tag=%s",
	  host, port, WIMC_EP, name, tag);

  return;
}

/*
 * get_container_path_from_wimc --
 *
 */
int get_container_path_from_wimc(char *name, char *tag,
				 char *host, char *port,
				 char *path) {
  int ret;
  long code=0;
  char wimcEP[LXCE_MAX_URL_LEN] = {0};
  CURL *curl=NULL;
  CURLcode res;
  struct Response response;

  /* Sanity check */
  if (!name || !tag || !host || !port) {
    return FALSE;
  }

  /* Steps are:
   * 1. create WIMC url with proper EP
   * 2. send CURL request with name and version.
   * 3. process response from WIMC.
   * 4. verify the path and return.
   */

  /* step-1 */
  create_wimc_url(&wimcEP[0], name, tag, host, port);

  curl = curl_easy_init();
  if (curl == NULL) {
    return code;
  }

  response.buffer = (char *)malloc(1);
  response.size   = 0;

  curl_easy_setopt(curl, CURLOPT_URL, wimcEP);

  curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "GET");
  curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
  curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);

  curl_easy_setopt(curl, CURLOPT_USERAGENT, "lxce/0.1");
  
  /* step-2 */
  res = curl_easy_perform(curl);

  if (res != CURLE_OK) {
    log_error("Error sending request to wimc: %s", curl_easy_strerror(res));
  } else {
    /* get status code. */
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
    /* Step-3 */
    ret = process_response_from_wimc(code, &response, path);
    if (ret != TRUE) {
      free(path);
    }
  }

  /* Validate to ensure we can read the path and it has the config.json
   * file at root.
   */
  if (path) {
    if (!is_valid_path(path)) {
      free(path);
      path=NULL;
    }
  }
  
 done:
  free(response.buffer);
  curl_easy_cleanup(curl);

  return code;
}

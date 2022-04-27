/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Methods related to service router forward plane
 *
 */

#include <curl/curl.h>
#include <curl/easy.h>
#include <ulfius.h>

#include "router.h"
#include "log.h"

/*
 * create_forward_request --
 *
 */
static req_t *init_forward_request(char *host, int port, char *method,
				   char *ep) {

  req_t *req = NULL;
  char url[MAX_LEN] = {0};

  req = (req_t *)calloc(1, sizeof(req_t));
  if (!req) {
    log_error("Error allocating memory of size: %lu", sizeof(req_t));
    return NULL;
  }

  sprintf(url, "http://%s:%d/%s", host, port, ep);

  if (ulfius_init_request(req) != U_OK) {
    goto failure;
  }

  ulfius_set_request_properties(req,
				U_OPT_HTTP_VERB, method,
				U_OPT_HTTP_URL, url,
				U_OPT_TIMEOUT, 20);
  return req;

 failure:
  ulfius_clean_request(req);
  free(req);

  return NULL;
}

/*
 * add_url_parameters --
 *
 */
static void add_url_parameters(req_t *req, Pattern *reqPattern) {

  Pattern *ptr=NULL;

  for (ptr=reqPattern; ptr; ptr=ptr->next) {
    ulfius_set_request_properties(req,
				  U_OPT_URL_PARAMETER,
				  ptr->key, ptr->value);
  }
}

/*
 * create_forward_request --
 *
 */
req_t *create_forward_request(Forward *forward, Pattern *reqPattern,
			      const req_t *request) {

  req_t *fRequest=NULL;

  /* Initialize the forward request */
  fRequest = init_forward_request(forward->ip, forward->port,
				  request->http_verb, forward->defaultPath);
  if (!fRequest) {
    log_error("Error init forward request");
    return NULL;
  }

  /* Add any parameter (key/value) to URL header */
  add_url_parameters(fRequest, reqPattern);

  /* Add any JSON data etc. */

  /* Adjust Content-type and other fields. */

  /* close the parameter list */
  ulfius_set_request_properties(fRequest, U_OPT_NONE);

  return fRequest;
}

/*
 * valid_forward_route --
 *
 */
int valid_forward_route(char *host, int port) {

  CURL     *curl;
  CURLcode response;
  char     url[MAX_LEN] = {0};

  if (host == NULL || port == 0) {
    return UKAMA_ERROR_INVALID_DEST;
  }

  curl = curl_easy_init();

  if (curl) {

    sprintf(url, "http://%s:%d/ping", host, port);
    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_HTTPGET, 1L);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 2L); /* 2 second timeout */

    response = curl_easy_perform(curl);

    curl_easy_cleanup(curl);

    if (response != CURLE_OK) {
      return UKAMA_ERROR_BAD_DEST;
    } else {
      return TRUE;
    }
  }

  return FALSE;
}

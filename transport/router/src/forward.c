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
 * add_url_parameters --
 *
 */
static void add_url_parameters(req_t **req, Pattern *reqPattern) {

  Pattern *ptr=NULL;
  int ret = U_OK;

  for (ptr=reqPattern; ptr; ptr=ptr->next) {
    ulfius_set_request_properties(*req,
          U_OPT_URL_PARAMETER,
          ptr->key, ptr->value,
          U_OPT_NONE);
  }

}

/*
 * create_forward_request --
 *
 */
static int  prepare_url(char *host, int port, char* url) {


  if (host == NULL ) return FALSE;

  sprintf(url, "http://%s:%d", host, port);


  return TRUE;
}

/*
 * create_forward_request --
 *
 */
req_t *create_forward_request(Forward *forward, Pattern *reqPattern,
                             const req_t *request) {

  req_t *fRequest = NULL;
  char url[MAX_LEN] = {0};
  char ep[MAX_LEN] = {0};
  char host[MAX_LEN] = {0};

  /* Prepare URL for service */
  if (!prepare_url(forward->ip, forward->port,
                  url)) {
    log_error("Error preparing URL for requested service");
    return NULL;
  }

  /* Preparing Request */
  fRequest = (req_t *)calloc(1, sizeof(req_t));
  if (!fRequest) {
    log_error("Error allocating memory of size: %lu", sizeof(req_t));
    return NULL;
  }

  if (ulfius_init_request(fRequest) == U_OK) {
    if (ulfius_copy_request(fRequest, request) != U_OK) {
      log_error("Internal error copying request.");
      return NULL;
    }
  } else {
    log_error("Internal error initializing request.");
    return NULL;
  }

  /* Update URL */
  ulfius_set_request_properties(fRequest,
        U_OPT_HTTP_URL, url,
        U_OPT_NONE);

  /* End Point */
  if ( forward->defaultPath == NULL) return NULL;

  if (forward->defaultPath[0] == '/') {
      sprintf(ep, "%s", forward->defaultPath);
  } else  {
      sprintf(ep, "/%s", forward->defaultPath);
  }

  /* Update EP */
  ulfius_set_request_properties(fRequest,
         U_OPT_HTTP_URL_APPEND, ep,
         U_OPT_NONE);

  /* Host Header */
  sprintf(host, "%s:%d", forward->ip, forward->port);

  /* Update Header */
  ulfius_set_request_properties(fRequest,
        U_OPT_HEADER_PARAMETER, "Host", host,
        U_OPT_NONE);

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

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

/*
 * HTTP methods related functions.
 *
 */

#include <ulfius.h>

#include "methods.h"

/* 
 * log_response --
 *
 */

void log_response(resp_t *resp){

}

/*
 * create_request -- create REST api request of req_type 
 *
 */

req_t* create_http_request(char *url, char* ep, char *req_type) {
  int ret;
  req_t* req = NULL;

  req = (req_t *)malloc(sizeof(req_t));

  if (req) {
    
    ret = ulfius_init_request(req);
    if (ret != U_OK) {
      goto failure;
    }

    ret = ulfius_set_request_properties(req, 
					U_OPT_HTTP_VERB, req_type,
					U_OPT_HTTP_URL, url,
					/* U_OPT_HTTP_URL_APPEND, ep, */ // XXX
					U_OPT_TIMEOUT, 20, /* XXX */
					U_OPT_NONE);
    if (ret != U_OK) {
      goto failure;
    }
  }

  return req;
  
 failure:
  ulfius_clean_request(req);
  free(req);

  return NULL;
}


/*
 * send_request -- send http request and return the response.
 *
 */

resp_t *send_http_request(req_t *req) {

  int ret;
  resp_t *resp;

  resp = (resp_t *)malloc(sizeof(resp_t));

  if (resp == NULL) {
    return NULL;
  }

  ret = ulfius_init_response(resp);
  if (ret != U_OK) {
    goto failure;
  }
  
  ret = ulfius_send_http_request(req, resp);
  if (ret == U_OK) {
    return resp;
  }
  
 failure:
  ulfius_clean_response(resp);
  free(resp);
  
  return NULL;
}

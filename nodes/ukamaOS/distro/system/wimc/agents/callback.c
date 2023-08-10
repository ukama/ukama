/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <ulfius.h>
#include <string.h>
#include <curl/curl.h>
#include <curl/easy.h>
#include <jansson.h>

#include "agent.h"
#include "wimc.h"
#include "log.h"
#include "err.h"
#include "agent/jserdes.h"
#include "common/utils.h"
#include "http_status.h"

#include "usys_types.h"
#include "usys_mem.h"
#include "usys_log.h"

static void log_json(json_t *json) {

    char *str = NULL;

    if (json == NULL) return;
    
    str = json_dumps(json, 0);
    if (str) {
        usys_log_debug("json str: %s", str);
        usys_free(str);
    }
}

static int validate_post_request(WimcReq *req) {

  WFetch *fetch=NULL;
  WContent *content=NULL;

  fetch = req->fetch;
  
  if (!validate_url(fetch->cbURL) ||
      !validate_url(fetch->content->indexURL) ||
      !validate_url(fetch->content->storeURL))
      return USYS_FALSE;

  return USYS_TRUE;
}

static void free_wimc_request(WimcReq *req) {

    WContent *content;

    if (!req) return;
    if (req->type != WREQ_FETCH) return;

    free(req->fetch->cbURL);
    content = req->fetch->content;

    if (content) {
        usys_free(content->name);
        usys_free(content->tag);
        usys_free(content->method);
        usys_free(content->indexURL);
        usys_free(content->storeURL);
        usys_free(content);
    }

    usys_free(req->fetch);
    usys_free(req);
}

int agent_web_service_cb_post_capp(const struct _u_request *request,
                                   struct _u_response *response,
                                   void *data) {

    int retCode=0;

    json_t       *json = NULL;
    json_error_t jerr;
    WimcReq      *req = NULL;
    char         *wimcURL = NULL;

    wimcURL = (char *)data;

    json = ulfius_get_json_body_request(request, &jerr);
    if (!json) {
        usys_log_error("JSON error for the agent register request: %s",
                       jerr.text);
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }
    log_json(json);
    
    req = (WimcReq *)calloc(1, sizeof(WimcReq));
    if (req == NULL) {
        usys_log_error("Error allocating memory of size: %ld", sizeof(WimcReq));
        retCode = HttpStatus_InternalServerError;
        goto done;
    }

    if (!deserialize_wimc_request(&req, json)) {
        usys_log_error("Error deserializing wimc request");
        retCode = HttpStatus_BadRequest;
        goto done;
    }

    /* setup cbURL */
    req->fetch->cbURL = (char *)calloc(1, WIMC_MAX_URL_LEN);
    if (req->fetch->cbURL != NULL) {
        sprintf(req->fetch->cbURL, "%s/v1/agents/update", wimcURL);
    } else {
        usys_log_error("Error allocating memory: %ld", WIMC_MAX_URL_LEN);
        retCode = HttpStatus_InternalServerError;
        goto done;
    }

    if (!validate_post_request(req)) {
        usys_log_error("Invalid parameters for capp post");
        retCode = HttpStatus_BadRequest;
        goto done;
    }

    retCode = HttpStatus_OK;
    process_capp_fetch_request(req->fetch);

done:
    free_wimc_request(req);
    json_decref(json);

    ulfius_set_string_body_response(response, retCode, HttpStatusStr(retCode));

    return U_CALLBACK_CONTINUE;
}

int agent_web_service_cb_default(const URequest *request,
                                 UResponse *response,
                                 void *data) {

    ulfius_set_string_body_response(response,
                                    HttpStatus_Unauthorized,
                                    HttpStatusStr(HttpStatus_Unauthorized));

    return U_CALLBACK_CONTINUE;
}

int agent_web_service_cb_ping(const URequest *request,
                              UResponse *response,
                              void *data) {

    ulfius_set_string_body_response(response,
                                    HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));

    return U_CALLBACK_CONTINUE;
}

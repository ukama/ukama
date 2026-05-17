/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <curl/easy.h>
#include <jansson.h>
#include <pthread.h>
#include <string.h>
#include <ulfius.h>

#include "agent.h"
#include "agent/jserdes.h"
#include "common/utils.h"
#include "err.h"
#include "http_status.h"
#include "log.h"
#include "wimc.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_services.h"
#include "usys_types.h"

#include "version.h"

extern void process_capp_fetch_request(WFetch *fetch);
extern int agent_job_is_active(const char *name, const char *tag);

static void log_json(json_t *json) {

    char *str;

    if (json == NULL) {
        return;
    }

    str = json_dumps(json, 0);
    if (str != NULL) {
        usys_log_debug("json str: %s", str);
        usys_free(str);
    }
}

static int validate_post_request(WimcReq *req) {

    WFetch *fetch;
    WContent *content;

    if (req == NULL || req->fetch == NULL || req->fetch->content == NULL) {
        return USYS_FALSE;
    }

    fetch   = req->fetch;
    content = fetch->content;

    if (!validate_url(fetch->cbURL) ||
        !validate_url(content->indexURL)) {
        return USYS_FALSE;
    }

    if (strcmp(content->method, WIMC_METHOD_TARGZ_STR) == 0) {
        return USYS_TRUE;
    }

    if (!validate_url(content->storeURL)) {
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

static void free_wimc_request(WimcReq *req) {

    WContent *content;

    if (req == NULL) {
        return;
    }

    if (req->fetch == NULL) {
        usys_free(req);
        return;
    }

    usys_free(req->fetch->cbURL);
    content = req->fetch->content;

    if (content != NULL) {
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

int agent_web_service_cb_post_capp(const URequest *request,
                                   UResponse *response,
                                   void *data) {

    int retCode;
    int wimcPort;
    json_t *json;
    json_error_t jerr;
    WimcReq *req;

    (void)data;

    retCode = HttpStatus_InternalServerError;
    json = NULL;
    req = NULL;

    json = ulfius_get_json_body_request(request, &jerr);
    if (json == NULL) {
        usys_log_error("JSON error for WIMC request: %s", jerr.text);
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }
    log_json(json);

    req = (WimcReq *)calloc(1, sizeof(WimcReq));
    if (req == NULL) {
        usys_log_error("Error allocating WimcReq");
        goto done;
    }

    if (!deserialize_wimc_request(&req, json)) {
        usys_log_error("Error deserializing WIMC request");
        retCode = HttpStatus_BadRequest;
        goto done;
    }

    if (agent_job_is_active(req->fetch->content->name,
                            req->fetch->content->tag)) {
        retCode = HttpStatus_Conflict;
        goto done;
    }

    req->fetch->cbURL = (char *)calloc(1, WIMC_MAX_URL_LEN);
    if (req->fetch->cbURL == NULL) {
        usys_log_error("Error allocating callback URL");
        goto done;
    }

    wimcPort = usys_find_service_port(SERVICE_WIMC);
    if (wimcPort <= 0 ||
        snprintf(req->fetch->cbURL, WIMC_MAX_URL_LEN,
                 "http://localhost:%d/v1/apps/%s/%s/stats",
                 wimcPort,
                 req->fetch->content->name,
                 req->fetch->content->tag) >= WIMC_MAX_URL_LEN) {
        usys_log_error("Unable to build WIMC callback URL");
        goto done;
    }

    if (!validate_post_request(req)) {
        usys_log_error("Invalid parameters for capp post");
        retCode = HttpStatus_BadRequest;
        goto done;
    }

    process_capp_fetch_request(req->fetch);
    retCode = HttpStatus_OK;

 done:
    free_wimc_request(req);
    json_decref(json);

    ulfius_set_string_body_response(response, retCode, HttpStatusStr(retCode));
    return U_CALLBACK_CONTINUE;
}

int agent_web_service_cb_default(const URequest *request,
                                 UResponse *response,
                                 void *data) {

    (void)request;
    (void)data;

    ulfius_set_string_body_response(response,
                                    HttpStatus_NotFound,
                                    HttpStatusStr(HttpStatus_NotFound));
    return U_CALLBACK_CONTINUE;
}

int agent_web_service_cb_ping(const URequest *request,
                              UResponse *response,
                              void *data) {

    (void)request;
    (void)data;

    ulfius_set_empty_body_response(response, HttpStatus_OK);
    return U_CALLBACK_CONTINUE;
}

int agent_web_service_cb_version(const URequest *request,
                                 UResponse *response,
                                 void *data) {

    (void)request;
    (void)data;

    ulfius_set_string_body_response(response, HttpStatus_OK, VERSION);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *user_data) {

    (void)request;
    (void)user_data;

    ulfius_set_string_body_response(response,
                                    HttpStatus_MethodNotAllowed,
                                    HttpStatusStr(HttpStatus_MethodNotAllowed));
    return U_CALLBACK_CONTINUE;
}

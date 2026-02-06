/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <string.h>
#include <stdlib.h>

#include "web_service.h"
#include <ulfius.h>

#include "http_status.h"
#include "json_serdes.h"
#include "metrics_store.h"
#include "version.h"

static void setup_unsupported_methods(UInst *instance,
                                      char *allowedMethod,
                                      char *prefix,
                                      char *resource) {

    if (strcmp(allowedMethod, "GET") != 0) {
        ulfius_add_endpoint_by_val(instance, "GET", prefix, resource, 0,
                                   &web_service_cb_not_allowed, (void *)allowedMethod);
    }

    if (strcmp(allowedMethod, "POST") != 0) {
        ulfius_add_endpoint_by_val(instance, "POST", prefix, resource, 0,
                                   &web_service_cb_not_allowed, (void *)allowedMethod);
    }

    if (strcmp(allowedMethod, "PUT") != 0) {
        ulfius_add_endpoint_by_val(instance, "PUT", prefix, resource, 0,
                                   &web_service_cb_not_allowed, (void *)allowedMethod);
    }   

    if (strcmp(allowedMethod, "DELETE") != 0) {
        ulfius_add_endpoint_by_val(instance, "DELETE", prefix, resource, 0,
                                   &web_service_cb_not_allowed, (void *)allowedMethod);
    }
}

static void setup_webservice_endpoints(UInst *instance, EpCtx *epCtx) {

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX, API_RES_EP("ping"), 0,
                               &web_service_cb_ping, epCtx);
    setup_unsupported_methods(instance, "GET", URL_PREFIX, API_RES_EP("ping"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX, API_RES_EP("version"), 0,
                               &web_service_cb_version, epCtx);
    setup_unsupported_methods(instance, "GET", URL_PREFIX, API_RES_EP("version"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX, API_RES_EP("status"), 0,
                               &web_service_cb_status, epCtx);
    setup_unsupported_methods(instance, "GET", URL_PREFIX, API_RES_EP("status"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX, API_RES_EP("metrics"), 0,
                               &web_service_cb_metrics, epCtx);
    setup_unsupported_methods(instance, "GET", URL_PREFIX, API_RES_EP("metrics"));

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX, API_RES_EP("diagnostics/chg"), 0,
                               &web_service_cb_post_diag_chg, epCtx);
    setup_unsupported_methods(instance, "POST", URL_PREFIX, API_RES_EP("diagnostics/chg"));

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX, API_RES_EP("diagnostics/parallel"), 0,
                               &web_service_cb_post_diag_parallel, epCtx);
    setup_unsupported_methods(instance, "POST", URL_PREFIX, API_RES_EP("diagnostics/parallel"));

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX, API_RES_EP("diagnostics/bufferbloat"), 0,
                               &web_service_cb_post_diag_bufferbloat, epCtx);
    setup_unsupported_methods(instance, "POST", URL_PREFIX, API_RES_EP("diagnostics/bufferbloat"));

    ulfius_set_default_endpoint(instance, &web_service_cb_default, epCtx);
}

int start_web_service(Config *config, UInst *serviceInst, EpCtx *epCtx) {

    if (ulfius_init_instance(serviceInst, config->servicePort, NULL, NULL) != U_OK) {
        return USYS_FALSE;
    }

    u_map_put(serviceInst->default_headers, "Access-Control-Allow-Origin", "*");

    setup_webservice_endpoints(serviceInst, epCtx);

    if (ulfius_start_framework(serviceInst) != U_OK) {
        ulfius_stop_framework(serviceInst);
        ulfius_clean_instance(serviceInst);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

int web_service_cb_ping(const URequest *request, UResponse *response, void *epConfig) {
    (void)request; (void)epConfig;
    ulfius_set_string_body_response(response, HttpStatus_OK, HttpStatusStr(HttpStatus_OK));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_version(const URequest *request, UResponse *response, void *epConfig) {
    (void)request; (void)epConfig;
    ulfius_set_string_body_response(response, HttpStatus_OK, VERSION);
    return U_CALLBACK_CONTINUE;
}

static void respond_json(UResponse *response, int status, json_t *json) {

    char *s = json_dumps(json, 0);
    if (!s) {
        ulfius_set_string_body_response(response, HttpStatus_InternalServerError,
                                        HttpStatusStr(HttpStatus_InternalServerError));
        return;
    }

    ulfius_set_string_body_response(response, status, s);
    u_map_put(response->map_header, "Content-Type", "application/json");

    free(s);
}

int web_service_cb_status(const URequest *request, UResponse *response, void *epConfig) {
    (void)request;
    EpCtx *ctx = (EpCtx *)epConfig;
    json_t *o = json_backhaul_status(ctx->store);
    respond_json(response, HttpStatus_OK, o);
    json_decref(o);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_metrics(const URequest *request, UResponse *response, void *epConfig) {
    (void)request;
    EpCtx *ctx = (EpCtx *)epConfig;
    json_t *o = json_backhaul_metrics(ctx->store);
    respond_json(response, HttpStatus_OK, o);
    json_decref(o);
    return U_CALLBACK_CONTINUE;
}

static int respond_accepted(UResponse *response, const char *name) {
    json_t *o = json_object();
    json_object_set_new(o, "queued", json_true());
    json_object_set_new(o, "name", json_string(name ? name : ""));
    respond_json(response, HttpStatus_Accepted, o);
    json_decref(o);
    return U_CALLBACK_CONTINUE;
}

/* Option-A: request flag only */
int web_service_cb_post_diag_chg(const URequest *request,
                                 UResponse *response,
                                 void *epConfig) {
    (void)request;
    EpCtx *ctx = (EpCtx *)epConfig;
    metrics_store_request_diag(ctx->store, BACKHAUL_DIAG_CHG);
    return respond_accepted(response, "chg");
}

int web_service_cb_post_diag_parallel(const URequest *request,
                                      UResponse *response,
                                      void *epConfig) {
    (void)request;
    EpCtx *ctx = (EpCtx *)epConfig;
    metrics_store_request_diag(ctx->store, BACKHAUL_DIAG_PARALLEL);
    return respond_accepted(response, "parallel");
}

int web_service_cb_post_diag_bufferbloat(const URequest *request,
                                         UResponse *response,
                                         void *epConfig) {
    (void)request;
    EpCtx *ctx = (EpCtx *)epConfig;
    metrics_store_request_diag(ctx->store, BACKHAUL_DIAG_BUFFERBLOAT);
    return respond_accepted(response, "bufferbloat");
}

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig) {
    (void)request; (void)epConfig;
    ulfius_set_string_body_response(response,
                                    HttpStatus_NotFound,
                                    HttpStatusStr(HttpStatus_NotFound));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *user_data) {
    (void)request; (void)user_data;
    ulfius_set_string_body_response(response,
                                    HttpStatus_MethodNotAllowed,
                                    HttpStatusStr(HttpStatus_MethodNotAllowed));
    return U_CALLBACK_CONTINUE;
}

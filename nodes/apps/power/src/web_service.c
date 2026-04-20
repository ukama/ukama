/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <stdlib.h>
#include <arpa/inet.h>
#include <netinet/in.h>

#include "powerd.h"
#include "web_service.h"
#include "metrics_store.h"
#include "http_status.h"
#include "json_types.h"
#include "json_serdes.h"

#include "ulfius.h"
#include "usys_log.h"

#include "version.h"

static int build_bind_addr(const char *ip,
                           uint16_t port,
                           struct sockaddr_in *sa) {

    if (!sa) return USYS_FALSE;

    memset(sa, 0, sizeof(*sa));
    sa->sin_family = AF_INET;
    sa->sin_port = htons(port);

    if (!ip || !*ip || !strcmp(ip, "0.0.0.0")) {
        sa->sin_addr.s_addr = htonl(INADDR_ANY);
        return USYS_TRUE;
    }

    if (inet_pton(AF_INET, ip, &sa->sin_addr) != 1) {
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

static void setup_unsupported_methods(UInst *instance,
                                      char *allowedMethod,
                                      char *prefix,
                                      char *resource) {

    if (strcmp(allowedMethod, "GET") != 0) {
        ulfius_add_endpoint_by_val(instance, "GET", prefix, resource, 0,
                                   &web_service_cb_not_allowed,
                                   (void *)allowedMethod);
    }

    if (strcmp(allowedMethod, "POST") != 0) {
        ulfius_add_endpoint_by_val(instance, "POST", prefix, resource, 0,
                                   &web_service_cb_not_allowed,
                                   (void *)allowedMethod);
    }

    if (strcmp(allowedMethod, "PUT") != 0) {
        ulfius_add_endpoint_by_val(instance, "PUT", prefix, resource, 0,
                                   &web_service_cb_not_allowed,
                                   (void *)allowedMethod);
    }

    if (strcmp(allowedMethod, "DELETE") != 0) {
        ulfius_add_endpoint_by_val(instance, "DELETE", prefix, resource, 0,
                                   &web_service_cb_not_allowed,
                                   (void *)allowedMethod);
    }
}

static int respond_json(UResponse *response, int status, json_t *obj) {

    char *s = NULL;

    if (!obj) {
        ulfius_add_header_to_response(response, "Content-Type", "application/json");
        ulfius_set_string_body_response(response, status, "{}");
        return U_CALLBACK_CONTINUE;
    }

    s = json_dumps(obj, JSON_INDENT(2));
    json_decref(obj);

    if (!s) {
        ulfius_add_header_to_response(response, "Content-Type", "application/json");
        ulfius_set_string_body_response(response, status, "{}");
        return U_CALLBACK_CONTINUE;
    }

    ulfius_add_header_to_response(response, "Content-Type", "application/json");
    ulfius_set_string_body_response(response, status, s);
    free(s);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_get_ping(const URequest *request,
                            UResponse *response,
                            void *epConfig) {

    (void)request;
    (void)epConfig;

    ulfius_set_string_body_response(response,
                                    HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_get_version(const URequest *request,
                               UResponse *response,
                               void *epConfig) {

    (void)request;
    (void)epConfig;

    ulfius_set_string_body_response(response, HttpStatus_OK, VERSION);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_get_metrics(const URequest *request,
                               UResponse *response,
                               void *epConfig) {

    EpCtx *ctx = (EpCtx *)epConfig;
    PowerSnapshot snap;
    PowerMetrics m;
    json_t *o = NULL;

    (void)request;

    if (!ctx || !ctx->store || !ctx->config) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    metrics_store_get(ctx->store, &snap);
    power_metrics_from_snapshot(&snap, ctx->config->boardName, &m);

    o = json_serdes_power_metrics_to_json(&m);
    return respond_json(response, HttpStatus_OK, o);
}

static void setup_webservice_endpoints(UInst *instance, EpCtx *epCtx) {

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX, API_RES_EP("ping"), 0,
                               &web_service_cb_get_ping, epCtx);
    setup_unsupported_methods(instance, "GET", URL_PREFIX, API_RES_EP("ping"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX, API_RES_EP("version"), 0,
                               &web_service_cb_get_version, epCtx);
    setup_unsupported_methods(instance, "GET", URL_PREFIX, API_RES_EP("version"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX, API_RES_EP("metrics"), 0,
                               &web_service_cb_get_metrics, epCtx);
    setup_unsupported_methods(instance, "GET", URL_PREFIX, API_RES_EP("metrics"));

    /* backward compatibility */
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX, API_RES_EP("status"), 0,
                               &web_service_cb_get_metrics, epCtx);
    setup_unsupported_methods(instance, "GET", URL_PREFIX, API_RES_EP("status"));

    ulfius_set_default_endpoint(instance, &web_service_cb_default, epCtx);
}

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig) {

    (void)request;
    (void)epConfig;

    ulfius_set_string_body_response(response,
                                    HttpStatus_NotFound,
                                    HttpStatusStr(HttpStatus_NotFound));
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

int start_web_service(Config *config, UInst *inst, EpCtx *ctx) {

    int rc;
    struct sockaddr_in bindAddr;

    if (!config || !inst || !ctx) {
        usys_log_error("web_service: invalid args");
        return USYS_FALSE;
    }

    if (!build_bind_addr(config->listenAddr, config->listenPort, &bindAddr)) {
        usys_log_error("web_service: invalid bind address '%s'",
                       config->listenAddr ? config->listenAddr : "(null)");
        return USYS_FALSE;
    }

    usys_log_debug("web_service: init addr=%s port=%d",
                   config->listenAddr ? config->listenAddr : "(null)",
                   config->listenPort);

    rc = ulfius_init_instance(inst,
                              config->listenPort,
                              &bindAddr,
                              NULL);
    if (rc != U_OK) {
        usys_log_error("web_service: ulfius_init_instance failed rc=%d",
                       rc);
        return USYS_FALSE;
    }

    if (inst->default_headers) {
        rc = u_map_put(inst->default_headers,
                       "Access-Control-Allow-Origin",
                       "*");
        if (rc != U_OK) {
            usys_log_error("web_service: u_map_put failed rc=%d", rc);
        }
    }

    usys_log_debug("web_service: adding endpoints");
    setup_webservice_endpoints(inst, ctx);

    usys_log_debug("web_service: starting framework");
    rc = ulfius_start_framework(inst);
    if (rc != U_OK) {
        usys_log_error("web_service: ulfius_start_framework failed rc=%d",
                       rc);
        ulfius_clean_instance(inst);
        return USYS_FALSE;
    }

    usys_log_info("web_service: listening on %s:%d",
                  config->listenAddr,
                  config->listenPort);

    return USYS_TRUE;
}

void web_service_stop(struct _u_instance *inst) {

    ulfius_stop_framework(inst);
    ulfius_clean_instance(inst);
}

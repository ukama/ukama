/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "powerd.h"
#include "web_service.h"
#include "metrics_store.h"
#include "http_status.h"
#include "json_types.h"
#include "json_serdes.h"

#include "ulfius.h"
#include "usys_log.h"

#include "version.h"

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

	char *s;

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

int web_service_cb_get_status(const URequest *request,
                              UResponse *response,
                              void *epConfig) {

	(void)request;

	EpCtx *ctx = (EpCtx *)epConfig;
	PowerSnapshot snap;
	PowerMetrics m;

	if (!ctx || !ctx->store) {
		ulfius_set_string_body_response(response,
		                                HttpStatus_InternalServerError,
		                                HttpStatusStr(HttpStatus_InternalServerError));
		return U_CALLBACK_CONTINUE;
	}

    /* 1. Get internal snapshot */
	metrics_store_get(ctx->store, &snap);

    /* 2. Convert snapshot → API metrics */
	power_metrics_from_snapshot(&snap, ctx->config->boardName, &m);

    /* 3. Serialize */
	json_t *o = json_serdes_power_metrics_to_json(&m);
	respond_json(response, HttpStatus_OK, o);
	json_decref(o);

	return U_CALLBACK_CONTINUE;
}

static void setup_webservice_endpoints(UInst *instance, EpCtx *epCtx) {

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX, API_RES_EP("ping"), 0,
                               &web_service_cb_get_ping, epCtx);
    setup_unsupported_methods(instance, "GET", URL_PREFIX, API_RES_EP("ping"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX, API_RES_EP("version"), 0,
                               &web_service_cb_get_version, epCtx);
    setup_unsupported_methods(instance, "GET", URL_PREFIX, API_RES_EP("version"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX, API_RES_EP("status"), 0,
                               &web_service_cb_get_status, epCtx);
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

	if (ulfius_init_instance(inst, config->listenPort, NULL, NULL) != U_OK) {
		return USYS_FALSE;
	}

	u_map_put(inst->default_headers, "Access-Control-Allow-Origin", "*");

	setup_webservice_endpoints(inst, ctx);

	if (ulfius_start_framework(inst) != U_OK) {
		ulfius_stop_framework(inst);
		ulfius_clean_instance(inst);
		return USYS_FALSE;
	}

	usys_log_info("web_service: listening on %s:%d",
	              ctx->config->listenAddr,
	              ctx->config->listenPort);

	return USYS_TRUE;
}

void web_service_stop(struct _u_instance *inst) {

	ulfius_stop_framework(inst);
	ulfius_clean_instance(inst);
}

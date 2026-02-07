/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "web_service.h"
#include "metrics_store.h"

#include "ulfius.h"
#include "usys_log.h"

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

static int web_service_cb_status(const URequest *request,
                                 UResponse *response,
                                 void *epConfig) {

	EpCtx *ctx = (EpCtx *)epConfig;
	PowerSnapshot s;

	(void)request;

	if (!ctx || !ctx->store) {
		ulfius_set_string_body_response(response,
		                               HttpStatus_InternalServerError,
		                               "store not ready");
		return U_CALLBACK_CONTINUE;
	}

	metrics_store_get(ctx->store, &s);
	return respond_json(response, HttpStatus_OK, metrics_store_to_json(&s));
}

static void setup_webservice_endpoints(UInst *inst, EpCtx *ctx) {

	/* Single endpoint only */
	ulfius_add_endpoint_by_val(inst, "GET",  "/v1/status", NULL, 0, &web_service_cb_status, ctx);

	/* Explicit not-allowed for common mistakes */
	ulfius_add_endpoint_by_val(inst, "POST", "/v1/status", NULL, 0, &web_service_cb_not_allowed, ctx);
	ulfius_add_endpoint_by_val(inst, "PUT",  "/v1/status", NULL, 0, &web_service_cb_not_allowed, ctx);

	/* Default */
	ulfius_set_default_endpoint(inst, &web_service_cb_default, ctx);
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

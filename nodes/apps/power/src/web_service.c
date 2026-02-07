/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "web_service.h"
#include "http_status.h"
#include "json_serdes.h"
#include "metrics_store.h"
#include "usys_log.h"
#include "version.h"

static int respond_json(UResponse *response, int code, json_t *o) {

	char *s = json_dumps(o, JSON_COMPACT);
	if (!s) {
		ulfius_set_string_body_response(response, HttpStatus_InternalServerError, "json error");
		return U_CALLBACK_CONTINUE;
	}

	ulfius_add_header_to_response(response, "Content-Type", "application/json");
	ulfius_set_string_body_response(response, code, s);

	free(s);
	return U_CALLBACK_CONTINUE;
}

int web_service_cb_ping(const URequest *request, UResponse *response, void *user_data) {

	(void)request; (void)user_data;

	json_t *o = json_object();
	json_object_set_new(o, "pong", json_true());
	respond_json(response, HttpStatus_OK, o);
	json_decref(o);

	return U_CALLBACK_CONTINUE;
}

int web_service_cb_version(const URequest *request, UResponse *response, void *user_data) {

	(void)request; (void)user_data;

	json_t *o = json_object();
	json_object_set_new(o, "name", json_string("powerd"));
	json_object_set_new(o, "version", json_string(VERSION));
	json_object_set_new(o, "git", json_string(GIT_COMMIT));
	respond_json(response, HttpStatus_OK, o);
	json_decref(o);

	return U_CALLBACK_CONTINUE;
}

int web_service_cb_power(const URequest *request, UResponse *response, void *epConfig) {

	(void)request;
	EpCtx *ctx = (EpCtx *)epConfig;

	PowerMetrics m;
	metrics_store_get(ctx->store, &m);

	json_t *o = json_serdes_power_metrics_to_json(&m);
	respond_json(response, HttpStatus_OK, o);
	json_decref(o);

	return U_CALLBACK_CONTINUE;
}

int web_service_cb_default(const URequest *request, UResponse *response, void *epConfig) {

	(void)request; (void)epConfig;
	ulfius_set_string_body_response(response,
	                               HttpStatus_NotFound,
	                               HttpStatusStr(HttpStatus_NotFound));
	return U_CALLBACK_CONTINUE;
}

int web_service_cb_not_allowed(const URequest *request, UResponse *response, void *user_data) {

	(void)request; (void)user_data;
	ulfius_set_string_body_response(response,
	                               HttpStatus_MethodNotAllowed,
	                               HttpStatusStr(HttpStatus_MethodNotAllowed));
	return U_CALLBACK_CONTINUE;
}

int web_service_start(struct _u_instance *inst, EpCtx *ctx) {

	if (ulfius_init_instance(inst, ctx->cfg->listenPort, ctx->cfg->listenAddr, NULL) != U_OK) {
		usys_log_error("web_service: init instance failed");
		return -1;
	}

	/* GET /v1/ping */
	ulfius_add_endpoint_by_val(inst, "GET", "/v1", "ping", 0, &web_service_cb_ping, ctx);

	/* GET /v1/version */
	ulfius_add_endpoint_by_val(inst, "GET", "/v1", "version", 0, &web_service_cb_version, ctx);

	/* GET /v1/power */
	ulfius_add_endpoint_by_val(inst, "GET", "/v1", "power", 0, &web_service_cb_power, ctx);

	/* default */
	ulfius_set_default_endpoint(inst, &web_service_cb_default, ctx);

	if (ulfius_start_framework(inst) != U_OK) {
		usys_log_error("web_service: start framework failed");
		return -1;
	}

	usys_log_info("web_service: listening on %s:%d", ctx->cfg->listenAddr, ctx->cfg->listenPort);
	return 0;
}

void web_service_stop(struct _u_instance *inst) {

	ulfius_stop_framework(inst);
	ulfius_clean_instance(inst);
}

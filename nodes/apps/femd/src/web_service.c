/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#include "web_service.h"
#include "web_handlers.h"
#include "usys_log.h"

#define API_RES_EP(x) "/" x

static void setup_webservice_endpoints(UInst *instance, WebCtx *ctx) {

    ulfius_add_endpoint_by_val(instance, "GET",  URL_PREFIX, API_RES_EP("ops/:opId"), 0, &web_cb_get_op, ctx);

    ulfius_add_endpoint_by_val(instance, "GET",  URL_PREFIX, API_RES_EP("controller/snapshot"), 0, &web_cb_get_ctrl_snapshot, ctx);
    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX, API_RES_EP("controller/sample"),   0, &web_cb_post_ctrl_sample, ctx);

    ulfius_add_endpoint_by_val(instance, "GET",  URL_PREFIX, API_RES_EP("fems/:femId/snapshot"), 0, &web_cb_get_fem_snapshot, ctx);
    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX, API_RES_EP("fems/:femId/sample"),   0, &web_cb_post_fem_sample, ctx);

    ulfius_set_default_endpoint(instance, &web_cb_default, ctx);
}

int start_web_service(ServerConfig *serverConfig, UInst *serviceInst, WebCtx *ctx) {

    if (!serverConfig || !serverConfig->config || !serviceInst || !ctx) {
        return USYS_FALSE;
    }

    if (ulfius_init_instance(serviceInst,
                             serverConfig->config->servicePort,
                             NULL,
                             NULL) != U_OK) {
        usys_log_error("Error initializing instance for webservice port %d",
                       serverConfig->config->servicePort);
        return USYS_FALSE;
    }

    u_map_put(serviceInst->default_headers, "Access-Control-Allow-Origin", "*");

    setup_webservice_endpoints(serviceInst, ctx);

    if (!start_framework(serviceInst)) {
        usys_log_error("Failed to start webservices on port: %d",
                       serverConfig->config->servicePort);
        return USYS_FALSE;
    }

    usys_log_info("Webservice started on port: %d",
                  serverConfig->config->servicePort);

    return USYS_TRUE;
}

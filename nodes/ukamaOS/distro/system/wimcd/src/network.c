/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */
#include <ulfius.h>
#include <stdlib.h>

#include "wimc.h"
#include "callback.h"
#include "agent.h"

#include "usys_log.h"
#include "usys_types.h"

static int start_framework(UInst *instance) {

    if (ulfius_start_framework(instance) != U_OK) {
        usys_log_error("Error starting the webservice/websocket.");

        ulfius_stop_framework(instance);
        ulfius_clean_instance(instance);

        return USYS_FALSE;
    }

    usys_log_debug("Webservice sucessfully started.");
    return USYS_TRUE;
}

static void setup_webservice_endpoints(Config *config,
                                       UInst *instance) {

    /* capp related end-points */
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &web_service_cb_ping, config);

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("version"), 0,
                               &web_service_cb_version, config);

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("apps/:name/:tag"), 0,
                               &web_service_cb_post_app, config);

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("apps/:name/:tag/status"), 0,
                               &web_service_cb_get_app_status, config);

    /* used by agent to update its stats */
    ulfius_add_endpoint_by_val(instance, "PUT", URL_PREFIX,
                               API_RES_EP("apps/:name/:tag/stats"), 0,
                               &web_service_cb_put_app_stats_update, config);

    ulfius_set_default_endpoint(instance,
                                &web_service_cb_default,
                                config);
}

int start_web_service(Config *config, UInst *serviceInst) {

    if (ulfius_init_instance(serviceInst,
                             config->servicePort,
                             NULL,
                             NULL) != U_OK) {
        usys_log_error("Error initializing instance for webservice port %d",
                       config->servicePort);
        return USYS_FALSE;
    }

    u_map_put(serviceInst->default_headers, "Access-Control-Allow-Origin", "*");

    setup_webservice_endpoints(config, serviceInst);

    if (!start_framework(serviceInst)) {
        usys_log_error("Failed to start webservices on port: %d",
                       config->servicePort);
        return USYS_FALSE;
    }

    usys_log_debug("Webservice started on port: %d", config->servicePort);

    return USYS_TRUE;
}

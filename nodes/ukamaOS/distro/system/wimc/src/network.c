/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <ulfius.h>
#include <stdlib.h>

#include "wimc.h"
#include "callback.h"
#include "agent.h"

#include "usys_log.h"
#include "usys_types.h"


static int start_framework(Config *config, UInst *instance) {

    if (ulfius_start_framework(instance) != U_OK) {
        usys_log_error("Error starting the webservice/websocket.");

        ulfius_stop_framework(instance); /* don't think need this. XXX */
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
                               API_RES_EP("capps/:name/:tag"), 0,
                               &web_service_cb_get_capp, config);

    /* Agent related end-points */
    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("agents/:id"), 0,
                               &web_service_cb_post_agent, config);

    ulfius_add_endpoint_by_val(instance, "DELETE", URL_PREFIX,
                               API_RES_EP("agents/:id"), 0,
                               &web_service_cb_delete_agent, config);

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("agents/update/:id/"), 0,
                               &web_service_cb_post_agent_update, config);

    /* default - 403 */
    ulfius_set_default_endpoint(instance,
                                &web_service_cb_default,
                                config);
}

int start_web_service(Config *config, UInst *serviceInst) {

    if (ulfius_init_instance(serviceInst,
                             atoi(config->servicePort),
                             NULL,
                             NULL) != U_OK) {
        usys_log_error("Error initializing instance for webservice port %s",
                       config->servicePort);
        return USYS_FALSE;
    }

    /* Set few params. */
    u_map_put(serviceInst->default_headers, "Access-Control-Allow-Origin", "*");

    /* setup endpoints and methods callback. */
    setup_webservice_endpoints(config, serviceInst);

    if (!start_framework(config, serviceInst)) {
        usys_log_error("Failed to start webservices on port: %s",
                       config->servicePort);
        return USYS_FALSE;
    }

    usys_log_debug("Webservice started on port: %s", config->servicePort);

    return USYS_TRUE;
}

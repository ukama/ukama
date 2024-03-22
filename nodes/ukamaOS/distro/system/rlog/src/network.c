/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <ulfius.h>
#include <stdlib.h>
#include <string.h>

#include "usys_types.h"
#include "usys_log.h"

#include "rlogd.h"
#include "web_service.h"

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

static void setup_webservice_endpoints(UInst *instance) {

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &web_service_cb_ping, NULL);

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("log/"), 0,
                               &web_service_cb_post_log, NULL);

    ulfius_set_default_endpoint(instance, &web_service_cb_default, NULL);
}

int start_web_service(int port, UInst *serviceInst) {

    if (ulfius_init_instance(serviceInst, port, NULL, NULL) != U_OK) {
        usys_log_error("Error initializing instance for webservice port %d", port);
        return USYS_FALSE;
    }

    u_map_put(serviceInst->default_headers, "Access-Control-Allow-Origin", "*");
    setup_webservice_endpoints(serviceInst);

    if (!start_framework(serviceInst)) {
        usys_log_error("Failed to start webservices on port: %d", port);
        return USYS_FALSE;
    }

    usys_log_debug("Webservice started on port: %d", port);

    return USYS_TRUE;
}

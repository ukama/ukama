/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Network related stuff based on ulfius framework for the Agents.
 */

#include <ulfius.h>
#include <stdlib.h>

#include "log.h"
#include "callback.h"
#include "agent.h"
#include "wimc.h"
#include "agent/network.h"
#include "agent/jserdes.h"
#include "agent/callback.h"

#include "usys_types.h"
#include "usys_log.h"

#define AGENT_EP "/v1/capps"
#define STAT_EP  "/stats"

static void setup_endpoints(MethodType *method, struct _u_instance *instance) {

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &agent_web_service_cb_ping, NULL);

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("capp"), 0,
                               &agent_web_service_cb_post_capp, method);

    ulfius_set_default_endpoint(instance,
                                &agent_web_service_cb_default, NULL);
}

bool start_web_service(char *port,
                      MethodType *method,
                      struct _u_instance *webInstance) {

    if (ulfius_init_instance(webInstance, atoi(port), NULL, NULL) != U_OK) {
        usys_log_error("Error initializing instance for port %s", port);
        return USYS_FALSE;
    }

    u_map_put(webInstance->default_headers, "Access-Control-Allow-Origin", "*");
    webInstance->max_post_body_size = 1024;

    setup_endpoints(method, webInstance);
  
    if (ulfius_start_framework(webInstance) != U_OK) {
        usys_log_error("Failed to start webservices at port:%s", port);

        ulfius_stop_framework(webInstance); 
        ulfius_clean_instance(webInstance);
        
        return USYS_FALSE;
    }

    usys_log_debug("Webservice started on port:%s", port);

    return USYS_TRUE;
}

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
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>

#include "usys_types.h"
#include "usys_log.h"

#include "rlogd.h"
#include "websocket.h"
#include "web_service.h"

static int init_framework(UInst *inst,
                          struct sockaddr_in *bindAddr,
                          int bindPort) {

    if (ulfius_init_instance(inst,
                             bindPort,
                             bindAddr,
                             NULL)!= U_OK) {
        log_error("Error initializing instance for websocket: %d",
                  bindPort);
        return USYS_FALSE;
    }

    u_map_put(inst->default_headers, "Access-Control-Allow-Origin", "*");

    return USYS_TRUE;
}

static int start_framework(UInst *instance) {

    if (ulfius_start_framework(instance) != U_OK) {
        usys_log_error("Error starting the web framework");

        ulfius_stop_framework(instance);
        ulfius_clean_instance(instance);

        return USYS_FALSE;
    }

    usys_log_debug("Webservice sucessfully started.");
    return USYS_TRUE;
}

static void setup_unsupported_methods(UInst *instance,
                                      char *allowedMethod,
                                      char *prefix,
                                      char *resource) {

    if (strcmp(allowedMethod, "GET") != 0) {
        ulfius_add_endpoint_by_val(instance, "GET", prefix,
                                   resource, 0,
                                   &web_service_cb_not_allowed,
                                   (void *)allowedMethod);
    }

    if (strcmp(allowedMethod, "POST") != 0) {
        ulfius_add_endpoint_by_val(instance, "POST", prefix,
                                   resource, 0,
                                   &web_service_cb_not_allowed,
                                   (void *)allowedMethod);
    }

    if (strcmp(allowedMethod, "PUT") != 0) {
        ulfius_add_endpoint_by_val(instance, "PUT", prefix,
                                   resource, 0,
                                   &web_service_cb_not_allowed,
                                   (void *)allowedMethod);
    }

    if (strcmp(allowedMethod, "DELETE") != 0) {
        ulfius_add_endpoint_by_val(instance, "DELETE", prefix,
                                   resource, 0,
                                   &web_service_cb_not_allowed,
                                   (void *)allowedMethod);
    }
}

static void setup_websocket_endpoints(char *nodeID, UInst *instance) {

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &web_socket_cb_ping, NULL);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("ping"));
    
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("logit"), 0,
                               &web_socket_cb_post_log, nodeID);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("logit"));

    ulfius_set_default_endpoint(instance, &web_socket_cb_default, NULL);
}

static void setup_webservice_endpoints(UInst *instance) {

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &web_service_cb_ping, NULL);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("ping"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("version"), 0,
                               &web_service_cb_version, NULL);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("version"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("level"), 0,
                               &web_service_cb_get_level, NULL);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("level"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("output"), 0,
                               &web_service_cb_get_output, NULL);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("output"));

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("level/:level"), 0,
                               &web_service_cb_post_level, NULL);
    setup_unsupported_methods(instance, "POST",
                              URL_PREFIX, API_RES_EP("level/:level"));

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("output/:output"), 0,
                               &web_service_cb_post_output, NULL);
    setup_unsupported_methods(instance, "POST",
                              URL_PREFIX, API_RES_EP("output/:output"));

    ulfius_set_default_endpoint(instance, &web_service_cb_default, NULL);
}

int start_websocket_server(char *nodeID, int port, UInst *websocketInst) {

    struct sockaddr_in bindAddr;

    memset(&bindAddr, 0, sizeof(bindAddr));
    bindAddr.sin_family = AF_INET;
    bindAddr.sin_port   = htons(port);

    if (getenv(ENV_BINDING_IP) == NULL) {
        bindAddr.sin_addr.s_addr = inet_addr(DEF_BINDING_IP);
    } else {
        bindAddr.sin_addr.s_addr = inet_addr(getenv(ENV_BINDING_IP));
    }

    if (init_framework(websocketInst, &bindAddr, port) != USYS_TRUE) {
        log_error("Error initializing websocket framework on port: %d", port);
        return USYS_FALSE;
    }

    /* setup endpoints and methods callback. */
    setup_websocket_endpoints(nodeID, websocketInst);

    if (start_framework(websocketInst) == USYS_FALSE) {
        log_error("Failed to start websocket at remote port %d", port);
        return USYS_FALSE;
    }

    log_debug("Websocket accepting on port: %d", port);

    return USYS_TRUE;
}

int start_web_services(int port, UInst *serviceInst) {


    if (ulfius_init_instance(serviceInst, port, NULL, NULL) != U_OK) {
        usys_log_error("Error initializing for webservice on: %d", port);
        return USYS_FALSE;
    }
    u_map_put(serviceInst->default_headers, "Access-Control-Allow-Origin", "*");

    /* setup endpoints and methods callback. */
    setup_webservice_endpoints(serviceInst);

    if (!start_framework(serviceInst)) {
        usys_log_error("Failed to start webservices for client on port: %d", port);
        return USYS_FALSE;
    }

    usys_log_debug("Webservice started on port: %d", port);

    return USYS_TRUE;
}

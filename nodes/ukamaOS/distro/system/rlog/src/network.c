/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <arpa/inet.h>
#include <netinet/in.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <ulfius.h>

#include "rlogd.h"
#include "usys_log.h"
#include "usys_types.h"
#include "web_service.h"

static int start_framework(UInst *instance) {
    if (ulfius_start_framework(instance) != U_OK) {
        usys_log_error("Error starting rlog admin service");
        ulfius_clean_instance(instance);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

static void setup_unsupported_methods(UInst *instance,
                                      const char *allowedMethod,
                                      const char *prefix,
                                      const char *resource) {
    if (strcmp(allowedMethod, "GET") != 0) {
        ulfius_add_endpoint_by_val(instance, "GET", prefix, resource, 0,
                                   &web_service_cb_not_allowed, NULL);
    }
    if (strcmp(allowedMethod, "POST") != 0) {
        ulfius_add_endpoint_by_val(instance, "POST", prefix, resource, 0,
                                   &web_service_cb_not_allowed, NULL);
    }
    if (strcmp(allowedMethod, "PUT") != 0) {
        ulfius_add_endpoint_by_val(instance, "PUT", prefix, resource, 0,
                                   &web_service_cb_not_allowed, NULL);
    }
    if (strcmp(allowedMethod, "DELETE") != 0) {
        ulfius_add_endpoint_by_val(instance, "DELETE", prefix, resource, 0,
                                   &web_service_cb_not_allowed, NULL);
    }
}

static void setup_endpoints(UInst *instance) {
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &web_service_cb_ping, NULL);
    setup_unsupported_methods(instance, "GET", URL_PREFIX,
                              API_RES_EP("ping"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("version"), 0,
                               &web_service_cb_version, NULL);
    setup_unsupported_methods(instance, "GET", URL_PREFIX,
                              API_RES_EP("version"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("status"), 0,
                               &web_service_cb_status, NULL);
    setup_unsupported_methods(instance, "GET", URL_PREFIX,
                              API_RES_EP("status"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("level"), 0,
                               &web_service_cb_get_level, NULL);
    setup_unsupported_methods(instance, "GET", URL_PREFIX,
                              API_RES_EP("level"));

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("level/:level"), 0,
                               &web_service_cb_post_level, NULL);
    setup_unsupported_methods(instance, "POST", URL_PREFIX,
                              API_RES_EP("level/:level"));

    ulfius_set_default_endpoint(instance,
                                &web_service_cb_default, NULL);
}

int start_web_services(int port, UInst *serviceInst) {
    struct sockaddr_in bindAddress;
    const char *bindingIp;

    bindingIp = getenv(ENV_RLOG_BINDING_IP);
    if (!bindingIp || !*bindingIp) {
        bindingIp = DEF_BINDING_IP;
    }

    memset(&bindAddress, 0, sizeof(bindAddress));
    bindAddress.sin_family = AF_INET;
    bindAddress.sin_port = htons((uint16_t)port);
    if (inet_pton(AF_INET, bindingIp, &bindAddress.sin_addr) != 1) {
        usys_log_error("Invalid rlog binding address: %s", bindingIp);
        return USYS_FALSE;
    }

    if (ulfius_init_instance(serviceInst, port,
                             &bindAddress, NULL) != U_OK) {
        usys_log_error("Unable to initialize rlog admin service on %s:%d",
                       bindingIp, port);
        return USYS_FALSE;
    }

    setup_endpoints(serviceInst);
    if (start_framework(serviceInst) != USYS_TRUE) {
        return USYS_FALSE;
    }

    usys_log_info("rlog admin service listening on %s:%d",
                  bindingIp, port);
    return USYS_TRUE;
}

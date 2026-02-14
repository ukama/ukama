/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <ulfius.h>

#include "web_service.h"
#include "web_handler.h"
#include "usys_log.h"

static int start_framework(UInst *instance) {

    if (ulfius_start_framework(instance) != U_OK) {
        usys_log_error("Error starting webservice.");

        ulfius_stop_framework(instance);
        ulfius_clean_instance(instance);

        return USYS_FALSE;
    }

    usys_log_debug("Webservice successfully started.");
    return USYS_TRUE;
}

static int method_allowed(const char *m, const char **allowed, int n) {

    int i;

    if (!m || !allowed || n <= 0) return 0;

    for (i = 0; i < n; i++) {
        if (allowed[i] && strcmp(m, allowed[i]) == 0) return 1;
    }

    return 0;
}

static void setup_verbs(UInst *instance,
                        const char *prefix,
                        const char *resource,
                        const char **allowed, int allowed_n,
                        const char *allow_header_str,
                        void *user_data) {

    static const char *KNOWN[] = {"GET","POST","PUT","PATCH","DELETE","OPTIONS"};
    int i;
    int known_n = (int)(sizeof(KNOWN) / sizeof(KNOWN[0]));

    ulfius_add_endpoint_by_val(instance, "OPTIONS", prefix, resource, 0,
                               &cb_options_ok, (void*)allow_header_str);

    for (i = 0; i < known_n; i++) {

        const char *verb = KNOWN[i];

        if (strcmp(verb, "OPTIONS") == 0) continue;

        if (!method_allowed(verb, allowed, allowed_n)) {
            ulfius_add_endpoint_by_val(instance, verb, prefix, resource, 0,
                                       &cb_not_allowed,
                                       (void*)allow_header_str);
        }
    }

    (void)user_data;
}

static void setup_webservice_endpoints(ServerConfig *serverConfig, UInst *instance, WebCtx *ctx) {

    const char *allowed_get_opts[]     = {"GET","OPTIONS"};
    const char *allowed_gpio[]         = {"GET","PUT","PATCH","OPTIONS"};
    const char *allowed_get_put_opts[] = {"GET","PUT","OPTIONS"};
    const char *allowed_post_opts[]    = {"POST","OPTIONS"};

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("version"), 0,
                               &cb_get_version, ctx);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("version"),
                allowed_get_opts, 2, "GET, OPTIONS", ctx);

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &cb_get_ping, ctx);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("ping"),
                allowed_get_opts, 2, "GET, OPTIONS", ctx);

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("metrics"), 0,
                               &cb_get_metrics, ctx);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("metrics"),
                allowed_get_opts, 2, "GET, OPTIONS", ctx);
    
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("fems"),
                               0, &cb_get_fems, ctx);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems"),
                allowed_get_opts, 2, "GET, OPTIONS", ctx);

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("fems/:femId"),
                               0, &cb_get_fem, ctx);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId"),
                allowed_get_opts, 2, "GET, OPTIONS", ctx);

    ulfius_add_endpoint_by_val(instance, "GET",   URL_PREFIX,
                               API_RES_EP("fems/:femId/gpio"),
                               0, &cb_get_gpio, ctx);
    ulfius_add_endpoint_by_val(instance, "PUT",   URL_PREFIX,
                               API_RES_EP("fems/:femId/gpio"),
                               0, &cb_put_gpio, ctx);
    ulfius_add_endpoint_by_val(instance, "PATCH", URL_PREFIX,
                               API_RES_EP("fems/:femId/gpio"),
                               0, &cb_patch_gpio, ctx);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId/gpio"),
                allowed_gpio, 4, "GET, PUT, PATCH, OPTIONS", ctx);

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("fems/:femId/dac"), 0, &cb_get_dac, ctx);
    ulfius_add_endpoint_by_val(instance, "PUT", URL_PREFIX,
                               API_RES_EP("fems/:femId/dac"), 0, &cb_put_dac, ctx);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId/dac"),
                allowed_get_put_opts, 3, "GET, PUT, OPTIONS", ctx);

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("fems/:femId/sensors/temperature"),
                               0, &cb_get_temp, ctx);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId/sensors/temperature"),
                allowed_get_opts, 2, "GET, OPTIONS", ctx);

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("fems/:femId/sensors/adc"),
                               0, &cb_get_adc_all, ctx);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId/sensors/adc"),
                allowed_get_opts, 2, "GET, OPTIONS", ctx);

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("fems/:femId/sensors/adc/channels/:channel"),
                               0, &cb_get_adc_chan, ctx);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId/sensors/adc/channels/:channel"),
                allowed_get_opts, 2, "GET, OPTIONS", ctx);

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("fems/:femId/safety/adc-thresholds"),
                               0, &cb_get_adc_thr, ctx);
    ulfius_add_endpoint_by_val(instance, "PUT", URL_PREFIX,
                               API_RES_EP("fems/:femId/safety/adc-thresholds"),
                               0, &cb_put_adc_thr, ctx);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId/safety/adc-thresholds"),
                allowed_get_put_opts, 3, "GET, PUT, OPTIONS", ctx);

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("fems/:femId/safety/restore"),
                               0, &cb_post_safety_restore, ctx);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId/safety/restore"),
                allowed_post_opts, 2, "POST, OPTIONS", ctx);

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("fems/:femId/eeprom/serial"),
                               0, &cb_get_serial, ctx);
    ulfius_add_endpoint_by_val(instance, "PUT", URL_PREFIX,
                               API_RES_EP("fems/:femId/eeprom/serial"),
                               0, &cb_put_serial, ctx);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId/eeprom/serial"),
                allowed_get_put_opts, 3, "GET, PUT, OPTIONS", ctx);

    ulfius_set_default_endpoint(instance, &cb_default, ctx);

    (void)serverConfig;
}

int start_web_service(ServerConfig *serverConfig, UInst *serviceInst, WebCtx *ctx) {

    if (!serverConfig || !serverConfig->config || !serviceInst || !ctx) {
        return USYS_FALSE;
    }

    if (ulfius_init_instance(serviceInst,
                             serverConfig->config->servicePort,
                             NULL,
                             NULL) != U_OK) {
        usys_log_error("Error initializing webservice port %d",
                       serverConfig->config->servicePort);
        return USYS_FALSE;
    }

    u_map_put(serviceInst->default_headers, "Access-Control-Allow-Origin", "*");

    setup_webservice_endpoints(serverConfig, serviceInst, ctx);

    if (!start_framework(serviceInst)) {
        usys_log_error("Failed to start webservice on port: %d",
                       serverConfig->config->servicePort);
        return USYS_FALSE;
    }

    usys_log_debug("Webservice started on port: %d",
                   serverConfig->config->servicePort);

    return USYS_TRUE;
}

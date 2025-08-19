/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#include <ulfius.h>
#include <stdlib.h>
#include <string.h>

#include "femd.h"
#include "config.h"
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

static int method_allowed(const char *m, const char **allowed, int n) {

    int i;
    for (i = 0; i < n; i++) {
        if (strcmp(m, allowed[i]) == 0) {
            return 1;
        }
    }
    return 0;
}

static void setup_verbs(UInst *instance,
                        const char *prefix,
                        const char *resource,
                        const char **allowed, int allowed_n,
                        const char *allow_header_str) {

    static const char *KNOWN[] = {"GET","POST","PUT","PATCH","DELETE","OPTIONS"};
    int i, known_n;

    known_n = sizeof(KNOWN)/sizeof(KNOWN[0]);

    /* OPTIONS always present */
    ulfius_add_endpoint_by_val(instance, "OPTIONS", prefix, resource, 0,
                               &cb_options_ok, (void*)allow_header_str);

    for (i = 0; i < known_n; i++) {

        const char *verb = KNOWN[i];

        if (strcmp(verb, "OPTIONS") == 0) continue;

        if (!method_allowed(verb, allowed, allowed_n)) {
            ulfius_add_endpoint_by_val(instance, verb, prefix, resource, 0,
                                       &web_service_cb_not_allowed,
                                       (void*)allow_header_str);
        }
    }
}

static void setup_webservice_endpoints(ServerConfig *config, UInst *instance) {

    const char *allowed_get_opts[]     = {"GET","OPTIONS"};
    const char *allowed_gpio[]         = {"GET","PUT","PATCH","OPTIONS"};
    const char *allowed_get_put_opts[] = {"GET","PUT","OPTIONS"};

    /* /v1/health */
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("health"), 0,
                               &cb_get_health, config);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("health"),
                allowed_get_opts, 2, "GET, OPTIONS");

    /* /v1/fems */
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("fems"),
                               0, &cb_get_fems, config);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems"),
                allowed_get_opts, 2, "GET, OPTIONS");

    /* /v1/fems/:femId */
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("fems/:femId"),
                               0, &cb_get_fem, config);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId"),
                allowed_get_opts, 2, "GET, OPTIONS");

    /* /v1/fems/:femId/gpio */
    ulfius_add_endpoint_by_val(instance, "GET",   URL_PREFIX,
                               API_RES_EP("fems/:femId/gpio"),
                               0, &cb_get_gpio, config);
    ulfius_add_endpoint_by_val(instance, "PUT",   URL_PREFIX,
                               API_RES_EP("fems/:femId/gpio"),
                               0, &cb_put_gpio, config);
    ulfius_add_endpoint_by_val(instance, "PATCH", URL_PREFIX,
                               API_RES_EP("fems/:femId/gpio"),
                               0, &cb_patch_gpio, config);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId/gpio"),
                allowed_gpio, 4, "GET, PUT, PATCH, OPTIONS");

    /* /v1/fems/:femId/dac */
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("fems/:femId/dac"), 0, &cb_get_dac, config);
    ulfius_add_endpoint_by_val(instance, "PUT", URL_PREFIX,
                               API_RES_EP("fems/:femId/dac"), 0, &cb_put_dac, config);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId/dac"),
                allowed_get_put_opts, 3, "GET, PUT, OPTIONS");

    /* /v1/fems/:femId/sensors/temperature */
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("fems/:femId/sensors/temperature"),
                               0, &cb_get_temp, config);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId/sensors/temperature"),
                allowed_get_opts, 2, "GET, OPTIONS");

    /* /v1/fems/:femId/sensors/adc */
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("fems/:femId/sensors/adc"),
                               0, &cb_get_adc_all, config);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId/sensors/adc"),
                allowed_get_opts, 2, "GET, OPTIONS");

    /* /v1/fems/:femId/sensors/adc/channels/:channel */
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("fems/:femId/sensors/adc/channels/:channel"),
                               0, &cb_get_adc_chan, config);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId/sensors/adc/channels/:channel"),
                allowed_get_opts, 2, "GET, OPTIONS");

    /* /v1/fems/:femId/safety/adc-thresholds */
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("fems/:femId/safety/adc-thresholds"),
                               0, &cb_get_adc_thr, config);
    ulfius_add_endpoint_by_val(instance, "PUT", URL_PREFIX,
                               API_RES_EP("fems/:femId/safety/adc-thresholds"),
                               0, &cb_put_adc_thr, config);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId/safety/adc-thresholds"),
                allowed_get_put_opts, 3, "GET, PUT, OPTIONS");

    /* /v1/fems/:femId/eeprom/serial */
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("fems/:femId/eeprom/serial"),
                               0, &cb_get_serial, config);
    ulfius_add_endpoint_by_val(instance, "PUT", URL_PREFIX,
                               API_RES_EP("fems/:femId/eeprom/serial"),
                               0, &cb_put_serial, config);
    setup_verbs(instance, URL_PREFIX, API_RES_EP("fems/:femId/eeprom/serial"),
                allowed_get_put_opts, 3, "GET, PUT, OPTIONS");

    /* default 404 */
    ulfius_set_default_endpoint(instance, &web_service_cb_default, config);
}

int start_web_service(ServerConfig *serverConfig, UInst *serviceInst) {

    if (ulfius_init_instance(serviceInst,
                             serverConfig->config->servicePort,
                             NULL,
                             NULL) != U_OK) {
		usys_log_error("Error initializing instance for webservice port %d",
                       serverConfig->config->servicePort);

		return USYS_FALSE;
	}

	u_map_put(serviceInst->default_headers, "Access-Control-Allow-Origin", "*");
	setup_webservice_endpoints(serverConfig, serviceInst);

	if (!start_framework(serviceInst)) {
		usys_log_error("Failed to start webservices on port: %d",
                       serverConfig->config->servicePort);
		return USYS_FALSE;
	}

	usys_log_debug("Webservice started on port: %d",
                   serverConfig->config->servicePort);

	return USYS_TRUE;
}

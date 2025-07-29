/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include "web_service.h"
#include "http_status.h"
#include "config.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

#include "version.h"

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *config) {
    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *config) {
    ulfius_set_string_body_response(response, HttpStatus_OK, VERSION);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_health(const URequest *request,
                          UResponse *response,
                          void *config) {
    ulfius_set_string_body_response(response, HttpStatus_OK, "healthy");
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *config) {
    ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                    HttpStatusStr(HttpStatus_NotFound));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *config) {
    ulfius_set_string_body_response(response, HttpStatus_MethodNotAllowed,
                                    HttpStatusStr(HttpStatus_MethodNotAllowed));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_fem_gpio_status(const URequest *request,
                                   UResponse *response,
                                   void *config) {
    ulfius_set_empty_body_response(response, HttpStatus_OK);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_fem_gpio_control(const URequest *request,
                                    UResponse *response,
                                    void *config) {
    ulfius_set_empty_body_response(response, HttpStatus_Accepted);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_fem_dac_control(const URequest *request,
                                   UResponse *response,
                                   void *config) {
    ulfius_set_empty_body_response(response, HttpStatus_Accepted);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_fem_dac_status(const URequest *request,
                                  UResponse *response,
                                  void *config) {
    ulfius_set_empty_body_response(response, HttpStatus_OK);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_fem_temperature(const URequest *request,
                                   UResponse *response,
                                   void *config) {
    ulfius_set_empty_body_response(response, HttpStatus_OK);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_fem_adc(const URequest *request,
                           UResponse *response,
                           void *config) {
    ulfius_set_empty_body_response(response, HttpStatus_OK);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_fem_eeprom(const URequest *request,
                              UResponse *response,
                              void *config) {
    ulfius_set_empty_body_response(response, HttpStatus_OK);
    return U_CALLBACK_CONTINUE;
}

int start_web_service(Config *config, UInst *serviceInst, void *webConfig) {
    if (ulfius_init_instance(serviceInst,
                             config->servicePort,
                             NULL,
                             NULL) != U_OK) {
        usys_log_error("Error initializing instance for webservice port %d",
                       config->servicePort);
        return STATUS_NOK;
    }

    /* Setup endpoints */
    ulfius_add_endpoint_by_val(serviceInst, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &web_service_cb_ping, config);

    ulfius_add_endpoint_by_val(serviceInst, "GET", URL_PREFIX,
                               API_RES_EP("version"), 0,
                               &web_service_cb_version, config);

    ulfius_add_endpoint_by_val(serviceInst, "GET", URL_PREFIX,
                               API_RES_EP("health"), 0,
                               &web_service_cb_health, config);

    ulfius_set_default_endpoint(serviceInst, &web_service_cb_default, config);

    if (ulfius_start_framework(serviceInst) != U_OK) {
        usys_log_error("Error starting the webservice/websocket.");
        ulfius_stop_framework(serviceInst);
        ulfius_clean_instance(serviceInst);
        return STATUS_NOK;
    }

    usys_log_debug("Webservice successfully started on port: %d", 
                   config->servicePort);
    return STATUS_OK;
}
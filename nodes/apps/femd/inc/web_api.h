/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef WEB_API_H
#define WEB_API_H

#include <stdbool.h>
#include <stdint.h>
#include "gpio_controller.h"
#include "i2c_controller.h"
#include "jserdes.h"
#include "http_status.h"

#define WEB_API_DEFAULT_PORT  8080
#define WEB_API_MAX_PAYLOAD   1024
#define WEB_API_MAX_RESPONSE  2048

typedef struct {
    int port;
    bool running;
    GpioController *gpio_controller;
    I2CController *i2c_controller;
} WebAPIServer;

typedef struct {
    char method[16];
    char path[256];
    char body[WEB_API_MAX_PAYLOAD];
    int content_length;
} HTTPRequest;

typedef struct {
    HttpStatusCode status_code;
    char content_type[64];
    char body[WEB_API_MAX_RESPONSE];
    int body_length;
} HTTPResponse;

int web_api_init(WebAPIServer *server, int port, GpioController *gpio_ctrl, I2CController *i2c_ctrl);
int web_api_start(WebAPIServer *server);
void web_api_stop(WebAPIServer *server);
void web_api_cleanup(WebAPIServer *server);

int web_api_handle_request(WebAPIServer *server, const HTTPRequest *request, HTTPResponse *response);
void web_api_set_response(HTTPResponse *response, HttpStatusCode status, const char *content_type, const char *body);
void web_api_set_json_response(HTTPResponse *response, HttpStatusCode status, JsonObj *json);
void web_api_set_error_response(HTTPResponse *response, HttpStatusCode status, const char *error_message);

int api_gpio_get_status(WebAPIServer *server, int fem_unit, HTTPResponse *response);
int api_gpio_set_control(WebAPIServer *server, int fem_unit, const char *gpio_name, bool enable, HTTPResponse *response);

int api_dac_set_voltages(WebAPIServer *server, int fem_unit, float carrier_voltage, float peak_voltage, HTTPResponse *response);
int api_dac_get_config(WebAPIServer *server, int fem_unit, HTTPResponse *response);

int api_temp_read(WebAPIServer *server, int fem_unit, HTTPResponse *response);
int api_temp_set_threshold(WebAPIServer *server, int fem_unit, float threshold, HTTPResponse *response);

int api_adc_read_channel(WebAPIServer *server, int fem_unit, int channel, HTTPResponse *response);
int api_adc_read_all(WebAPIServer *server, int fem_unit, HTTPResponse *response);
int api_adc_set_safety(WebAPIServer *server, float max_reverse_power, float max_current, HTTPResponse *response);

int api_eeprom_write_serial(WebAPIServer *server, int fem_unit, const char *serial, HTTPResponse *response);
int api_eeprom_read_serial(WebAPIServer *server, int fem_unit, HTTPResponse *response);

int parse_fem_unit(const char *path);
JsonObj *parse_request_json(const char *json_str);

#endif /* WEB_API_H */
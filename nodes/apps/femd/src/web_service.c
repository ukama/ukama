/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#include "femd.h"
#include "http_status.h"

#include "version.h"

int cb_not_allowed(const URequest *req,
                   UResponse *resp,
                   void *user_data) {
    ulfius_set_string_body_response(resp,
                                    HttpStatus_MethodNotAllowed,
                                    VERSION);
    return U_CALLBACK_CONTINUE;
}

int cb_options_ok(const URequest *req,
                  UResponse *resp,
                  void *user_data) {
    ulfius_set_string_body_response(resp,
                                    HttpStatus_OK,
                                    VERSION);
    return U_CALLBACK_CONTINUE;
}

int cb_default(const URequest *req,
               UResponse *resp,
               void *user_data) {
    ulfius_set_string_body_response(resp,
                                    HttpStatus_OK,
                                    VERSION);
    return U_CALLBACK_CONTINUE;
}

/* Health & discovery */
int cb_get_health(const URequest *req,
                  UResponse *resp,
                  void *user_data) {
    ulfius_set_string_body_response(resp,
                                    HttpStatus_OK,
                                    VERSION);
    return U_CALLBACK_CONTINUE;
}

int cb_get_ping(const URequest *req,
                UResponse *resp,
                void *user_data) {
    ulfius_set_string_body_response(resp,
                                    HttpStatus_OK,
                                    VERSION);
    return U_CALLBACK_CONTINUE;
}

int cb_get_version(const URequest *req,
                   UResponse *resp,
                   void *user_data) {
    ulfius_set_string_body_response(resp,
                                    HttpStatus_OK,
                                    VERSION);
    return U_CALLBACK_CONTINUE;
}

int cb_get_fems(const URequest *req,
                UResponse *resp,
                void *user_data) {
    ulfius_set_string_body_response(resp,
                                    HttpStatus_OK,
                                    VERSION);
    return U_CALLBACK_CONTINUE;
}

int cb_get_fem(const URequest *req,
               UResponse *resp,
               void *user_data) {
    ulfius_set_string_body_response(resp,
                                    HttpStatus_OK,
                                    VERSION);
    return U_CALLBACK_CONTINUE;
}

/* GPIO */
int cb_get_gpio(const URequest *req,
                UResponse *resp,
                void *user_data) {
    ulfius_set_string_body_response(resp,
                                    HttpStatus_OK,
                                    VERSION);
    return U_CALLBACK_CONTINUE;
}

int cb_put_gpio(const URequest *req,
                UResponse *resp,
                void *user_data) {
    ulfius_set_string_body_response(resp,
                                    HttpStatus_OK,
                                    VERSION);
    return U_CALLBACK_CONTINUE;
}

int cb_patch_gpio(const URequest *req,
                  UResponse *resp,
                  void *user_data) {
    ulfius_set_string_body_response(resp,
                                    HttpStatus_OK,
                                    VERSION);
    return U_CALLBACK_CONTINUE;
}

/* DAC */
int cb_get_dac(const URequest *req,
               UResponse *resp,
               void *user_data) {
    ulfius_set_string_body_response(resp,
                                    HttpStatus_OK,
                                    VERSION);
    return U_CALLBACK_CONTINUE;
}

int cb_put_dac(const URequest *req,
               UResponse *resp,
               void *user_data) {
    ulfius_set_string_body_response(resp,
                                    HttpStatus_OK,
                                    VERSION);
    return U_CALLBACK_CONTINUE;
}

/* Sensors */
int cb_get_temp(const URequest *req,
                UResponse *resp,
                void *user_data) {
  ulfius_set_string_body_response(resp,
                                  HttpStatus_OK,
                                  VERSION);
  return U_CALLBACK_CONTINUE;
}

int cb_get_adc_all(const URequest *req,
                   UResponse *resp,
                   void *user_data) {
  ulfius_set_string_body_response(resp,
                                  HttpStatus_OK,
                                  VERSION);
  return U_CALLBACK_CONTINUE;
}

int cb_get_adc_chan(const URequest *req,
                    UResponse *resp,
                    void *user_data) {
  ulfius_set_string_body_response(resp,
                                  HttpStatus_OK,
                                  VERSION);
  return U_CALLBACK_CONTINUE;
}

/* Safety thresholds */
int cb_get_adc_thr(const URequest *req,
                   UResponse *resp,
                   void *user_data) {
  ulfius_set_string_body_response(resp,
                                  HttpStatus_OK,
                                  VERSION);
  return U_CALLBACK_CONTINUE;
}

int cb_put_adc_thr(const URequest *req,
                   UResponse *resp,
                   void *user_data) {
  ulfius_set_string_body_response(resp,
                                  HttpStatus_OK,
                                  VERSION);
  return U_CALLBACK_CONTINUE;
}

/* EEPROM serial */
int cb_get_serial(const URequest *req,
                  UResponse *resp,
                  void *user_data) {
  ulfius_set_string_body_response(resp,
                                  HttpStatus_OK,
                                  VERSION);
  return U_CALLBACK_CONTINUE;
}

int cb_put_serial(const URequest *req,
                  UResponse *resp,
                  void *user_data) {
  ulfius_set_string_body_response(resp,
                                  HttpStatus_OK,
                                  VERSION);
  return U_CALLBACK_CONTINUE;
}

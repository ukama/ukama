/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef WEB_SERVICE_H_
#define WEB_SERVICE_H_

#include "femd.h"

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *config);

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *config);

int web_service_cb_health(const URequest *request,
                          UResponse *response,
                          void *config);

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *config);

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *config);

int web_service_cb_fem_gpio_status(const URequest *request,
                                   UResponse *response,
                                   void *config);

int web_service_cb_fem_gpio_control(const URequest *request,
                                    UResponse *response,
                                    void *config);

int web_service_cb_fem_dac_control(const URequest *request,
                                   UResponse *response,
                                   void *config);

int web_service_cb_fem_dac_status(const URequest *request,
                                  UResponse *response,
                                  void *config);

int web_service_cb_fem_temperature(const URequest *request,
                                   UResponse *response,
                                   void *config);

int web_service_cb_fem_adc(const URequest *request,
                           UResponse *response,
                           void *config);

int web_service_cb_fem_eeprom(const URequest *request,
                              UResponse *response,
                              void *config);

int start_web_service(Config *config,
                      UInst *serviceInst,
                      void *webConfig);

#endif /* WEB_SERVICE_H_ */
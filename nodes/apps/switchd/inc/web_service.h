/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef WEB_SERVICE_H_
#define WEB_SERVICE_H_

#include "switchd.h"

int web_service_start(SwitchdContext *ctx, UInst *serviceInst);
void web_service_stop(UInst *serviceInst);

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *epConfig);
int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *epConfig);
int web_service_cb_get_metrics(const URequest *request,
                               UResponse *response,
                               void *epConfig);
int web_service_cb_get_switch(const URequest *request,
                              UResponse *response,
                              void *epConfig);
int web_service_cb_get_switch_health(const URequest *request,
                                     UResponse *response,
                                     void *epConfig);
int web_service_cb_get_switch_capabilities(const URequest *request,
                                           UResponse *response,
                                           void *epConfig);
int web_service_cb_get_switch_alarms(const URequest *request,
                                     UResponse *response,
                                     void *epConfig);
int web_service_cb_get_switch_kpis(const URequest *request,
                                   UResponse *response,
                                   void *epConfig);
int web_service_cb_get_ports(const URequest *request,
                             UResponse *response,
                             void *epConfig);
int web_service_cb_get_port(const URequest *request,
                            UResponse *response,
                            void *epConfig);
int web_service_cb_get_port_kpis(const URequest *request,
                                 UResponse *response,
                                 void *epConfig);
int web_service_cb_post_port_admin(const URequest *request,
                                   UResponse *response,
                                   void *epConfig);
int web_service_cb_post_port_poe(const URequest *request,
                                 UResponse *response,
                                 void *epConfig);
int web_service_cb_post_port_poe_cycle(const URequest *request,
                                       UResponse *response,
                                       void *epConfig);
int web_service_cb_get_firmware(const URequest *request,
                                UResponse *response,
                                void *epConfig);
int web_service_cb_get_firmware_status(const URequest *request,
                                       UResponse *response,
                                       void *epConfig);
int web_service_cb_post_firmware_stage(const URequest *request,
                                       UResponse *response,
                                       void *epConfig);
int web_service_cb_post_firmware_apply(const URequest *request,
                                       UResponse *response,
                                       void *epConfig);
int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig);
int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *epConfig);

#endif /* WEB_SERVICE_H_ */

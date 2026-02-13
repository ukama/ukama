/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#ifndef WEB_HANDLER_H
#define WEB_HANDLER_H

#include <ulfius.h>

int cb_default(const URequest *request, UResponse *response, void *user_data);

int cb_options_ok(const URequest *request, UResponse *response, void *user_data);
int cb_not_allowed(const URequest *request, UResponse *response, void *user_data);

int cb_get_health(const URequest *request, UResponse *response, void *user_data);
int cb_get_version(const URequest *request, UResponse *response, void *user_data);
int cb_get_ping(const URequest *request, UResponse *response, void *user_data);

int cb_get_fems(const URequest *request, UResponse *response, void *user_data);
int cb_get_fem(const URequest *request, UResponse *response, void *user_data);

int cb_get_gpio(const URequest *request, UResponse *response, void *user_data);
int cb_put_gpio(const URequest *request, UResponse *response, void *user_data);
int cb_patch_gpio(const URequest *request, UResponse *response, void *user_data);

int cb_get_dac(const URequest *request, UResponse *response, void *user_data);
int cb_put_dac(const URequest *request, UResponse *response, void *user_data);

int cb_get_temp(const URequest *request, UResponse *response, void *user_data);

int cb_get_adc_all(const URequest *request, UResponse *response, void *user_data);
int cb_get_adc_chan(const URequest *request, UResponse *response, void *user_data);

int cb_get_adc_thr(const URequest *request, UResponse *response, void *user_data);
int cb_put_adc_thr(const URequest *request, UResponse *response, void *user_data);

int cb_post_safety_restore(const URequest *request, UResponse *response, void *user_data);

int cb_get_serial(const URequest *request, UResponse *response, void *user_data);
int cb_put_serial(const URequest *request, UResponse *response, void *user_data);

int cb_get_metrics(const URequest *request, UResponse *response, void *user_data);

#endif /* WEB_HANDLER_H */

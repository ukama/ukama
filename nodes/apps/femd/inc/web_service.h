/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef WEB_SERVICE_H_
#define WEB_SERVICE_H_

#include <ulfius.h>

#include "femd.h"

/* Health & discovery */
int cb_get_health   (const URequest *req, UResponse *resp, void *user_data);
int cb_get_fems     (const URequest *req, UResponse *resp, void *user_data);
int cb_get_fem      (const URequest *req, UResponse *resp, void *user_data);

/* GPIO */
int cb_get_gpio     (const URequest *req, UResponse *resp, void *user_data);
int cb_put_gpio     (const URequest *req, UResponse *resp, void *user_data);
int cb_patch_gpio   (const URequest *req, UResponse *resp, void *user_data);

/* DAC */
int cb_get_dac      (const URequest *req, UResponse *resp, void *user_data);
int cb_put_dac      (const URequest *req, UResponse *resp, void *user_data);

/* Sensors */
int cb_get_temp     (const URequest *req, UResponse *resp, void *user_data);
int cb_get_adc_all  (const URequest *req, UResponse *resp, void *user_data);
int cb_get_adc_chan (const URequest *req, UResponse *resp, void *user_data);

/* Safety thresholds */
int cb_get_adc_thr  (const URequest *req, UResponse *resp, void *user_data);
int cb_put_adc_thr  (const URequest *req, UResponse *resp, void *user_data);

/* EEPROM serial */
int cb_get_serial   (const URequest *req, UResponse *resp, void *user_data);
int cb_put_serial   (const URequest *req, UResponse *resp, void *user_data);

/* Generic helpers */
int web_service_cb_not_allowed(const URequest *req, UResponse *resp, void *user_data);
int web_service_cb_default    (const URequest *req, UResponse *resp, void *user_data);
int cb_options_ok             (const URequest *req, UResponse *resp, void *user_data);

#endif /* WEB_SERVICE_H_ */

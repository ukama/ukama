/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */
#ifndef API_JSON_H
#define API_JSON_H

#include <ulfius.h>
#include <jansson.h>

#include "femd.h"
#include "http_status.h"
#include "gpio_controller.h"

int json_set_ok(UResponse *resp, json_t *body);
int json_set_err(UResponse *resp, unsigned int code, const char *msg);

int json_get_bool(json_t *o, const char *key, int *out_val);
int json_get_double(json_t *o, const char *key, double *out_val);
int json_get_string(json_t *o, const char *key, const char **out_str);

/* Serialize GPIO status as legacy payload (incl. fem_unit).
   Caller must decref the returned object. */
json_t* json_gpio_status(const GpioStatus *st, int fem_unit_num);

/* Extract DAC PUT request body.
   Accepts carrier_voltage/peak_voltage OR carrierVoltage/peakVoltage.
   Returns 1 on success; 0 if required keys missing. */
int json_extract_dac_request(json_t *body, double *carrier, double *peak);

/* Extract ADC safety thresholds PUT body.
   Accepts max_reverse_power/max_current OR camelCase.
   Returns 1 on success; 0 if required keys missing. */
int json_extract_adc_thresholds(json_t *body, double *max_rp, double *max_i);

#endif /* API_JSON_H */

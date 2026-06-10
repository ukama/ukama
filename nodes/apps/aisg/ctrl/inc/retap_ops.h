/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef RETAP_OPS_H_
#define RETAP_OPS_H_

#include "retap.h"

typedef struct {
    char productNumber[64];
    char serialNumber[64];
    char hardwareVersion[64];
    char softwareVersion[64];
} RetapInfo;

typedef struct {
    uint8_t codes[RETAP_MAX_ALARMS];
    int count;
} RetapAlarmList;

bool retap_build_get_information(RetapRequest *request);
bool retap_parse_get_information(RetapResponse *response, RetapInfo *info);
bool retap_build_get_alarm_status(RetapRequest *request);
bool retap_parse_alarm_list(RetapResponse *response, RetapAlarmList *alarms);
bool retap_build_clear_active_alarms(RetapRequest *request);
bool retap_build_alarm_subscribe(RetapRequest *request);
bool retap_build_self_test(RetapRequest *request);
bool retap_build_send_configuration_data(RetapRequest *request,
                                         const uint8_t *data,
                                         size_t len);
bool retap_build_calibrate(RetapRequest *request);
bool retap_build_get_tilt(RetapRequest *request);
bool retap_parse_get_tilt(RetapResponse *response, int16_t *tiltTenthsDeg);
bool retap_build_set_tilt(RetapRequest *request, int16_t tiltTenthsDeg);
bool retap_build_get_device_data(RetapRequest *request, uint8_t field);
bool retap_build_reset_software(RetapRequest *request);

#endif /* RETAP_OPS_H_ */

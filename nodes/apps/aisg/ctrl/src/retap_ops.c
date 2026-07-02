/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "retap_ops.h"

#define RETAP_INFO_FIELD_COUNT 4
#define RETAP_INFO_STR_SIZE    64

static bool build_no_payload(RetapRequest *request, uint8_t procedure)
{
    if (request == NULL) {
        return false;
    }

    retap_request_init(request, procedure);

    return true;
}

static bool response_is_ok(RetapResponse *response)
{
    return retap_response_is_ok(response);
}

static bool copy_retap_string(char *dst,
                              size_t dstSize,
                              const uint8_t *src,
                              size_t len)
{
    if (dst == NULL || src == NULL || dstSize == 0) {
        return false;
    }

    snprintf(dst, dstSize, "%.*s", (int)len, src);

    return true;
}

bool retap_build_get_information(RetapRequest *request)
{
    return build_no_payload(request, RETAP_PROC_GET_INFORMATION);
}

bool retap_parse_get_information(RetapResponse *response, RetapInfo *info)
{
    char *fields[RETAP_INFO_FIELD_COUNT];
    size_t pos = 0;
    size_t len;
    int i;

    if (!response_is_ok(response) || info == NULL) {
        return false;
    }

    memset(info, 0, sizeof(RetapInfo));

    fields[0] = info->productNumber;
    fields[1] = info->serialNumber;
    fields[2] = info->hardwareVersion;
    fields[3] = info->softwareVersion;

    for (i = 0; i < RETAP_INFO_FIELD_COUNT && pos < response->dataLen; i++) {
        len = response->data[pos++];
        if (pos + len > response->dataLen) {
            return false;
        }

        if (!copy_retap_string(fields[i],
                               RETAP_INFO_STR_SIZE,
                               &response->data[pos],
                               len)) {
            return false;
        }

        pos += len;
    }

    return true;
}

bool retap_build_get_error_status(RetapRequest *request)
{
    return build_no_payload(request, RETAP_PROC_GET_ERROR_STATUS);
}

bool retap_build_get_alarm_status(RetapRequest *request)
{
    return retap_build_get_error_status(request);
}

bool retap_parse_return_code_list(RetapResponse *response, RetapAlarmList *alarms)
{
    size_t i;

    if (!response_is_ok(response) || alarms == NULL) {
        return false;
    }

    memset(alarms, 0, sizeof(RetapAlarmList));

    for (i = 0; i < response->dataLen && i < RETAP_MAX_ALARMS; i++) {
        alarms->codes[alarms->count++] = response->data[i];
    }

    return true;
}

bool retap_parse_alarm_list(RetapResponse *response, RetapAlarmList *alarms)
{
    return retap_parse_return_code_list(response, alarms);
}

bool retap_build_clear_active_alarms(RetapRequest *request)
{
    return build_no_payload(request, RETAP_PROC_CLEAR_ACTIVE_ALARMS);
}

bool retap_build_alarm_subscribe(RetapRequest *request)
{
    return build_no_payload(request, RETAP_PROC_ALARM_SUBSCRIBE);
}

bool retap_build_self_test(RetapRequest *request)
{
    return build_no_payload(request, RETAP_PROC_SELF_TEST);
}

bool retap_build_send_configuration_data(RetapRequest *request,
                                         const uint8_t *data,
                                         size_t len)
{
    if (request == NULL || data == NULL) {
        return false;
    }

    /*
     * TS 25.463 single-antenna SendConfigurationData carries at most
     * 70 octets of configuration data per elementary procedure. Larger
     * configuration files must be split by the backend and sent one
     * RETAP request at a time, waiting for OK before the next segment.
     */
    if (len == 0 || len > RETAP_CONFIG_SEGMENT_MAX) {
        return false;
    }

    retap_request_init(request, RETAP_PROC_SEND_CONFIG_DATA);
    memcpy(request->data, data, len);
    request->dataLen = len;

    return true;
}

bool retap_build_calibrate(RetapRequest *request)
{
    return build_no_payload(request, RETAP_PROC_CALIBRATE);
}

bool retap_build_get_tilt(RetapRequest *request)
{
    return build_no_payload(request, RETAP_PROC_GET_TILT);
}

bool retap_parse_get_tilt(RetapResponse *response, int16_t *tiltTenthsDeg)
{
    uint16_t raw;

    if (!response_is_ok(response) || tiltTenthsDeg == NULL) {
        return false;
    }

    if (response->dataLen < 2) {
        return false;
    }

    raw = (uint16_t)response->data[0] | ((uint16_t)response->data[1] << 8);
    *tiltTenthsDeg = (int16_t)raw;

    return true;
}

bool retap_build_set_tilt(RetapRequest *request, int16_t tiltTenthsDeg)
{
    uint16_t raw;

    if (request == NULL) {
        return false;
    }

    raw = (uint16_t)tiltTenthsDeg;

    retap_request_init(request, RETAP_PROC_SET_TILT);
    request->data[0] = (uint8_t)(raw & 0xFF);
    request->data[1] = (uint8_t)((raw >> 8) & 0xFF);
    request->dataLen = 2;

    return true;
}

bool retap_build_get_device_data(RetapRequest *request, uint8_t field)
{
    if (request == NULL) {
        return false;
    }

    retap_request_init(request, RETAP_PROC_GET_DEVICE_DATA);
    request->data[0] = field;
    request->dataLen = 1;

    return true;
}

bool retap_build_reset_software(RetapRequest *request)
{
    return build_no_payload(request, RETAP_PROC_RESET_SOFTWARE);
}

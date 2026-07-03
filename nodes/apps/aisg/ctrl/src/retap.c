/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "retap.h"

void retap_request_init(RetapRequest *request, uint8_t procedure)
{
    memset(request, 0, sizeof(RetapRequest));
    request->procedure = procedure;
}

void retap_response_init(RetapResponse *response)
{
    memset(response, 0, sizeof(RetapResponse));
}

bool retap_encode_request(RetapRequest *request,
                          uint8_t *buf,
                          size_t size,
                          size_t *len)
{
    if (request == NULL || buf == NULL || len == NULL) {
        return false;
    }

    if (request->dataLen > RETAP_MAX_PAYLOAD || size < request->dataLen + 3) {
        return false;
    }

    buf[0] = request->procedure;
    buf[1] = (uint8_t)(request->dataLen & 0xFF);
    buf[2] = (uint8_t)((request->dataLen >> 8) & 0xFF);

    if (request->dataLen > 0) {
        memcpy(&buf[3], request->data, request->dataLen);
    }

    *len = request->dataLen + 3;

    return true;
}

bool retap_decode_response(const uint8_t *buf,
                           size_t len,
                           RetapResponse *response)
{
    uint16_t dataLen;
    const uint8_t *data;

    if (buf == NULL || response == NULL || len < 4) {
        return false;
    }

    retap_response_init(response);

    response->procedure = buf[0];
    dataLen = (uint16_t)buf[1] | ((uint16_t)buf[2] << 8);
    if ((size_t)dataLen + 3 > len || dataLen < 1) {
        return false;
    }

    data = &buf[3];
    response->returnCode = data[0];

    if (response->returnCode == RETAP_RETURN_FAIL) {
        if (dataLen < 2) {
            return false;
        }
        response->failureReason = data[1];
        if (dataLen > 2) {
            if ((size_t)dataLen - 2 > RETAP_MAX_PAYLOAD) {
                return false;
            }
            memcpy(response->data, &data[2], (size_t)dataLen - 2);
            response->dataLen = (size_t)dataLen - 2;
        }
        return true;
    }

    if (dataLen > 1) {
        if ((size_t)dataLen - 1 > RETAP_MAX_PAYLOAD) {
            return false;
        }
        memcpy(response->data, &data[1], (size_t)dataLen - 1);
        response->dataLen = (size_t)dataLen - 1;
    }

    return true;
}

CtrlCode retap_failure_to_ctrl_code(uint8_t failureReason)
{
    switch (failureReason) {
    case 0x01: return CtrlCodeFormatError;
    case 0x05: return CtrlCodeBusy;
    case 0x08: return CtrlCodeFormatError;
    case 0x09: return CtrlCodeHardwareError;
    case 0x0B: return CtrlCodeHardwareError;
    case 0x0E: return CtrlCodeNotCalibrated;
    case 0x0F: return CtrlCodeNotConfigured;
    case 0x11: return CtrlCodeHardwareError;
    case 0x12: return CtrlCodeHardwareError;
    case 0x13: return CtrlCodeOutOfRange;
    case 0x19: return CtrlCodeUnsupportedProcedure;
    case 0x1E: return CtrlCodeUnsupportedProcedure;
    case 0x21: return CtrlCodeWorkingSoftwareMissing;
    case 0x23: return CtrlCodeBusy;
    default:   return CtrlCodeHardwareError;
    }
}

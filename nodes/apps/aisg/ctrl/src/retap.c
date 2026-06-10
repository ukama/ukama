/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "retap.h"

void retap_request_init(RetapRequest *request, uint8_t procedure) {
    memset(request, 0, sizeof(RetapRequest));
    request->procedure = procedure;
}

void retap_response_init(RetapResponse *response) {
    memset(response, 0, sizeof(RetapResponse));
}

bool retap_encode_request(RetapRequest *request,
                          uint8_t *buf,
                          size_t size,
                          size_t *len) {
    if (request == NULL || buf == NULL || len == NULL) return false;
    if (size < request->dataLen + 1) return false;

    buf[0] = request->procedure;
    memcpy(&buf[1], request->data, request->dataLen);
    *len = request->dataLen + 1;

    return true;
}

bool retap_decode_response(const uint8_t *buf,
                           size_t len,
                           RetapResponse *response) {
    if (buf == NULL || response == NULL || len < 2) return false;

    retap_response_init(response);
    response->procedure = buf[0];
    response->returnCode = buf[1];

    if (response->returnCode == RETAP_RETURN_FAIL) {
        if (len < 3) return false;
        response->failureReason = buf[2];
        if (len > 3) {
            memcpy(response->data, &buf[3], len - 3);
            response->dataLen = len - 3;
        }
        return true;
    }

    if (len > 2) {
        memcpy(response->data, &buf[2], len - 2);
        response->dataLen = len - 2;
    }

    return true;
}

CtrlCode retap_failure_to_ctrl_code(uint8_t failureReason) {
    switch (failureReason) {
    case 0x01: return CtrlCodeFormatError;
    case 0x02: return CtrlCodeBusy;
    case 0x03: return CtrlCodeHardwareError;
    case 0x04: return CtrlCodeWorkingSoftwareMissing;
    case 0x05: return CtrlCodeNotConfigured;
    case 0x06: return CtrlCodeNotCalibrated;
    case 0x07: return CtrlCodeOutOfRange;
    case 0x08: return CtrlCodeUnsupportedProcedure;
    default:   return CtrlCodeHardwareError;
    }
}

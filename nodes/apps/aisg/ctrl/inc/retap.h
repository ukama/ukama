/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef RETAP_H_
#define RETAP_H_

#include "ctrl.h"
#include "retap_codes.h"
#include "response.h"

typedef struct {
    uint8_t procedure;
    uint8_t data[RETAP_MAX_PAYLOAD];
    size_t dataLen;
} RetapRequest;

typedef struct {
    uint8_t procedure;
    uint8_t returnCode;
    uint8_t failureReason;
    uint8_t data[RETAP_MAX_PAYLOAD];
    size_t dataLen;
} RetapResponse;

void retap_request_init(RetapRequest *request, uint8_t procedure);
void retap_response_init(RetapResponse *response);
bool retap_encode_request(RetapRequest *request,
                          uint8_t *buf,
                          size_t size,
                          size_t *len);
bool retap_decode_response(const uint8_t *buf,
                           size_t len,
                           RetapResponse *response);
CtrlCode retap_failure_to_ctrl_code(uint8_t failureReason);

#endif /* RETAP_H_ */

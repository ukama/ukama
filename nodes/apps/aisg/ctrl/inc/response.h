/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef RESPONSE_H_
#define RESPONSE_H_

#include "ctrl.h"
#include "request.h"

#define CTRL_REASON_LEN                    256

typedef enum {
    CtrlCodeOk = 0,
    CtrlCodeBusy,
    CtrlCodeFormatError,
    CtrlCodeOutOfRange,
    CtrlCodeNotConfigured,
    CtrlCodeNotCalibrated,
    CtrlCodeWorkingSoftwareMissing,
    CtrlCodeHardwareError,
    CtrlCodeUnsupportedProcedure,
    CtrlCodeTransportError,
    CtrlCodeTimeout,
    CtrlCodeInvalidRequest
} CtrlCode;

typedef struct {
    char id[CTRL_REQ_ID_LEN];
    bool ok;
    CtrlCode code;
    char reason[CTRL_REASON_LEN];
    JsonObj *payload;
} CtrlResponse;

void ctrl_response_init(CtrlResponse *response, const char *id);
void ctrl_response_free(CtrlResponse *response);
bool ctrl_response_set_ok(CtrlResponse *response, JsonObj *payload);
bool ctrl_response_set_error(CtrlResponse *response,
                             CtrlCode code,
                             const char *reason);
JsonObj *ctrl_response_to_json(CtrlResponse *response);
const char *ctrl_code_str(CtrlCode code);

#endif /* RESPONSE_H_ */

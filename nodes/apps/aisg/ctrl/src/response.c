/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "response.h"

const char *ctrl_code_str(CtrlCode code)
{
    switch (code) {
    case CtrlCodeOk:
        return "OK";
    case CtrlCodeBusy:
        return "Busy";
    case CtrlCodeFormatError:
        return "FormatError";
    case CtrlCodeOutOfRange:
        return "OutOfRange";
    case CtrlCodeNotConfigured:
        return "NotConfigured";
    case CtrlCodeNotCalibrated:
        return "NotCalibrated";
    case CtrlCodeWorkingSoftwareMissing:
        return "WorkingSoftwareMissing";
    case CtrlCodeHardwareError:
        return "HardwareError";
    case CtrlCodeUnsupportedProcedure:
        return "UnsupportedProcedure";
    case CtrlCodeUnsupportedDeviceType:
        return "UnsupportedDeviceType";
    case CtrlCodeUnsupportedProtocolVersion:
        return "UnsupportedProtocolVersion";
    case CtrlCodeMultipleDevices:
        return "MultipleDevices";
    case CtrlCodeLinkNotConnected:
        return "LinkNotConnected";
    case CtrlCodeFrameReject:
        return "FrameReject";
    case CtrlCodeReceiverNotReady:
        return "ReceiverNotReady";
    case CtrlCodeProtocolError:
        return "ProtocolError";
    case CtrlCodeTransportError:
        return "TransportError";
    case CtrlCodeTimeout:
        return "Timeout";
    case CtrlCodeInvalidRequest:
        return "InvalidRequest";
    default:
        return "Unknown";
    }
}

void ctrl_response_init(CtrlResponse *response, const char *id)
{
    if (response == NULL) {
        return;
    }

    memset(response, 0, sizeof(CtrlResponse));
    snprintf(response->id, sizeof(response->id), "%s", id ? id : "req");
}

void ctrl_response_free(CtrlResponse *response)
{
    if (response == NULL) {
        return;
    }

    json_decref(response->payload);
    memset(response, 0, sizeof(CtrlResponse));
}

bool ctrl_response_set_ok(CtrlResponse *response, JsonObj *payload)
{
    if (response == NULL) {
        json_decref(payload);
        return false;
    }

    json_decref(response->payload);

    response->ok        = true;
    response->code      = CtrlCodeOk;
    response->reason[0] = '\0';
    response->payload   = payload ? payload : json_object();

    return response->payload != NULL;
}

bool ctrl_response_set_error(CtrlResponse *response,
                             CtrlCode code,
                             const char *reason)
{
    if (response == NULL) {
        return false;
    }

    json_decref(response->payload);

    response->ok   = false;
    response->code = code;
    response->payload = json_object();

    snprintf(response->reason,
             sizeof(response->reason),
             "%s",
             reason ? reason : ctrl_code_str(code));

    return response->payload != NULL;
}

JsonObj *ctrl_response_to_json(CtrlResponse *response)
{
    JsonObj *json    = NULL;
    JsonObj *payload = NULL;

    if (response == NULL) {
        return NULL;
    }

    payload = response->payload ? json_deep_copy(response->payload)
                                : json_object();
    if (payload == NULL) {
        return NULL;
    }

    json = json_object();
    if (json == NULL) {
        json_decref(payload);
        return NULL;
    }

    json_object_set_new(json, "id", json_string(response->id));
    json_object_set_new(json, "ok", json_boolean(response->ok));
    json_object_set_new(json,
                        "code",
                        json_string(ctrl_code_str(response->code)));
    json_object_set_new(json, "reason", json_string(response->reason));
    json_object_set_new(json, "payload", payload);

    return json;
}

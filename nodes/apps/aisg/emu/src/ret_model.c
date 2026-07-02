/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "ret_model.h"

static uint16_t vendor_to_u16(const char *vendor)
{
    uint8_t a = 'U';
    uint8_t b = 'K';

    if (vendor != NULL && vendor[0] != '\0') {
        a = (uint8_t)vendor[0];
        if (vendor[1] != '\0') {
            b = (uint8_t)vendor[1];
        }
    }

    return (uint16_t)(((uint16_t)a << 8) | b);
}

void ret_model_reset_l2(RetModel *model)
{
    if (model == NULL) {
        return;
    }

    model->state = RET_L2_NO_ADDRESS;
    model->address = RET_EMU_ADDR_DEFAULT;
    model->primaryNsExpected = 0;
    model->secondaryNs = 0;
    model->has3gppRelease = false;
    model->hasAisgVersion = false;
    model->negotiated3gppRelease = 0;
    model->negotiatedAisgVersion = 0;
}

void ret_model_init(RetModel *model,
                    const char *vendorCode,
                    const char *serial,
                    bool requiresConfig,
                    int16_t initialTiltTenths,
                    int16_t minTiltTenths,
                    int16_t maxTiltTenths)
{
    if (model == NULL) {
        return;
    }

    memset(model, 0, sizeof(*model));

    snprintf(model->vendorCodeStr,
             sizeof(model->vendorCodeStr),
             "%.2s",
             (vendorCode && vendorCode[0]) ? vendorCode : "UK");
    if (model->vendorCodeStr[1] == '\0') {
        model->vendorCodeStr[1] = 'K';
        model->vendorCodeStr[2] = '\0';
    }

    model->vendorCode = vendor_to_u16(model->vendorCodeStr);
    snprintf(model->serial,
             sizeof(model->serial),
             "%s",
             (serial && serial[0]) ? serial : "UKAMA00000000001");

    model->uniqueIdLen = 0;
    model->uniqueId[model->uniqueIdLen++] = (uint8_t)model->vendorCodeStr[0];
    model->uniqueId[model->uniqueIdLen++] = (uint8_t)model->vendorCodeStr[1];
    while (model->serial[model->uniqueIdLen - 2] != '\0' &&
           model->uniqueIdLen < AISG_XID_UNIQUE_ID_MAX) {
        model->uniqueId[model->uniqueIdLen] =
            (uint8_t)model->serial[model->uniqueIdLen - 2];
        model->uniqueIdLen++;
    }

    model->deviceType = RET_EMU_DEVICE_TYPE_SINGLE_RET;
    model->requiresConfig = requiresConfig;
    model->configured = !requiresConfig;
    model->calibrated = !requiresConfig;
    model->tiltTenths = initialTiltTenths;
    model->minTiltTenths = minTiltTenths;
    model->maxTiltTenths = maxTiltTenths;

    if (!model->configured) {
        ret_model_set_error(model, RETAP_RC_NOT_SCALED);
    }
    if (!model->calibrated) {
        ret_model_set_error(model, RETAP_RC_NOT_CALIBRATED);
    }

    ret_model_reset_l2(model);
}

const char *ret_l2_state_name(RetL2State state)
{
    switch (state) {
    case RET_L2_NO_ADDRESS:       return "NO_ADDRESS";
    case RET_L2_ADDRESS_ASSIGNED: return "ADDRESS_ASSIGNED";
    case RET_L2_CONNECTED:        return "CONNECTED";
    default:                      return "UNKNOWN";
    }
}

void ret_model_set_error(RetModel *model, uint8_t code)
{
    size_t i;

    if (model == NULL || code == RETAP_RC_OK) {
        return;
    }

    for (i = 0; i < model->activeErrorCount; i++) {
        if (model->activeErrors[i] == code) {
            return;
        }
    }

    if (model->activeErrorCount < RETAP_MAX_ALARMS) {
        model->activeErrors[model->activeErrorCount++] = code;
    }
}

void ret_model_clear_error(RetModel *model, uint8_t code)
{
    size_t i;

    if (model == NULL) {
        return;
    }

    for (i = 0; i < model->activeErrorCount; i++) {
        if (model->activeErrors[i] == code) {
            memmove(&model->activeErrors[i],
                    &model->activeErrors[i + 1],
                    model->activeErrorCount - i - 1);
            model->activeErrorCount--;
            return;
        }
    }
}

void ret_model_clear_errors(RetModel *model)
{
    if (model == NULL) {
        return;
    }

    model->activeErrorCount = 0;
}

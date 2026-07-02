/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AISG_EMU_RET_MODEL_H_
#define AISG_EMU_RET_MODEL_H_

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "retap_codes.h"
#include "xid.h"

#define RET_EMU_DEVICE_TYPE_SINGLE_RET     0x01
#define RET_EMU_DEVICE_TYPE_TMA            0x02
#define RET_EMU_DEVICE_TYPE_MULTI_RET      0x11
#define RET_EMU_ADDR_BROADCAST             0xFF
#define RET_EMU_ADDR_DEFAULT               0x00
#define RET_EMU_DEFAULT_ASSIGNED_ADDR      0x01
#define RET_EMU_3GPP_RELEASE_ID            0x06
#define RET_EMU_AISG_VERSION               0x02
#define RET_EMU_LINK_TIMEOUT_MS            180000

typedef enum {
    RET_L2_NO_ADDRESS = 0,
    RET_L2_ADDRESS_ASSIGNED,
    RET_L2_CONNECTED
} RetL2State;

typedef struct {
    RetL2State state;

    uint8_t address;
    uint8_t deviceType;
    uint8_t primaryNsExpected;
    uint8_t secondaryNs;

    bool has3gppRelease;
    bool hasAisgVersion;
    uint8_t negotiated3gppRelease;
    uint8_t negotiatedAisgVersion;

    char vendorCodeStr[3];
    uint16_t vendorCode;
    char serial[18];
    uint8_t uniqueId[AISG_XID_UNIQUE_ID_MAX];
    size_t uniqueIdLen;

    bool configured;
    bool calibrated;
    bool alarmSubscribed;
    bool busy;
    bool requiresConfig;

    int16_t tiltTenths;
    int16_t minTiltTenths;
    int16_t maxTiltTenths;

    uint8_t activeErrors[RETAP_MAX_ALARMS];
    size_t activeErrorCount;
} RetModel;

void ret_model_init(RetModel *model,
                    const char *vendorCode,
                    const char *serial,
                    bool requiresConfig,
                    int16_t initialTiltTenths,
                    int16_t minTiltTenths,
                    int16_t maxTiltTenths);
const char *ret_l2_state_name(RetL2State state);
void ret_model_reset_l2(RetModel *model);
void ret_model_set_error(RetModel *model, uint8_t code);
void ret_model_clear_error(RetModel *model, uint8_t code);
void ret_model_clear_errors(RetModel *model);

#endif /* AISG_EMU_RET_MODEL_H_ */

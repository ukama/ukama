/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AISG_V2_H_
#define AISG_V2_H_

#include "serial.h"
#include "retap.h"

#define AISG_ADDR_BROADCAST                0xFF
#define AISG_ADDR_DEFAULT                  0x00
#define AISG_CTRL_I_FRAME                  0x00
#define AISG_CTRL_XID                      0xBF
#define AISG_DEFAULT_TIMEOUT_MS            3000

typedef struct {
    SerialPort *serial;
    uint8_t deviceAddress;
} AisgBus;

typedef struct {
    bool present;
    uint8_t address;
    char model[64];
} AisgDevice;

void aisg_v2_bus_init(AisgBus *bus, SerialPort *serial);
bool aisg_v2_scan(AisgBus *bus, AisgDevice *device);
bool aisg_v2_send_retap(AisgBus *bus,
                        RetapRequest *request,
                        RetapResponse *response);

#endif /* AISG_V2_H_ */

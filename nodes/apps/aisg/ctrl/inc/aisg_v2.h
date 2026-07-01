/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AISG_V2_H_
#define AISG_V2_H_

#include <stdbool.h>
#include <stdint.h>

#include "serial.h"
#include "retap.h"

/*
 * Ukama AISG scope.
 *
 * We implement AISG v2.0 for one single-antenna RET device on one RS485 bus.
 * TS 25.461 / TS 25.462 / TS 25.463 are authoritative where they conflict
 * with AISG v2.0.
 *
 * Supported:
 *   - one physical RET device
 *   - single-antenna RET device type only
 *
 * Not supported:
 *   - multi-antenna RET
 *   - TMA
 *   - multiple RETs on one bus
 *   - daisy-chain discovery/control
 */
#define AISG_SCOPE_NAME                    "single-antenna-ret"

#define AISG_DEVICE_TYPE_SINGLE_RET        0x01
#define AISG_DEVICE_TYPE_TMA               0x02
#define AISG_DEVICE_TYPE_MULTI_RET         0x11
#define AISG_SUPPORTED_DEVICE_TYPE         AISG_DEVICE_TYPE_SINGLE_RET

#define AISG_ADDR_BROADCAST                0xFF
#define AISG_ADDR_DEFAULT                  0x00

/*
 * Phase 0/1 keeps the old transport shape visible and honest.
 * Real TS 25.462 control helpers are introduced in later phases.
 */
#define AISG_CTRL_I_FRAME                  0x00
#define AISG_CTRL_XID                      0xBF

#define AISG_DEFAULT_TIMEOUT_MS            3000
#define AISG_MAX_RX_ATTEMPTS               4

typedef struct {
    SerialPort *serial;
    uint8_t deviceAddress;
} AisgBus;

typedef struct {
    bool present;
    bool unsupported;
    uint8_t address;
    uint8_t deviceType;
    char model[64];
} AisgDevice;

void aisg_v2_bus_init(AisgBus *bus, SerialPort *serial);
bool aisg_v2_scan(AisgBus *bus, AisgDevice *device);
bool aisg_v2_send_retap(AisgBus *bus,
                        RetapRequest *request,
                        RetapResponse *response);

#endif /* AISG_V2_H_ */

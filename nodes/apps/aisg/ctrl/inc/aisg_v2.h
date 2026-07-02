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
#include <stddef.h>
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
#define AISG_ADDR_ASSIGNED                 0x01

#define AISG_3GPP_RELEASE_ID               0x06
#define AISG_PROTOCOL_VERSION              0x02
#define AISG_LINK_TIMEOUT_MS               180000
#define AISG_MIN_TURNAROUND_US             3000
#define AISG_HDLC_DEFAULT_INFO_MAX         78


#define AISG_DEFAULT_TIMEOUT_MS            3000
#define AISG_SCAN_EXTRA_TIMEOUT_MS         250
#define AISG_MAX_RX_ATTEMPTS               4

#define AISG_XID_FI                        0x81
#define AISG_XID_GI_ADDRESSING             0xF0
#define AISG_XID_PI_UNIQUE_ID              0x01
#define AISG_XID_PI_HDLC_ADDRESS           0x02
#define AISG_XID_PI_BIT_MASK               0x03
#define AISG_XID_PI_DEVICE_TYPE            0x04
#define AISG_XID_PI_3GPP_RELEASE           0x05
#define AISG_XID_PI_VENDOR_CODE            0x06
#define AISG_XID_PI_AISG_VERSION           20
#define AISG_XID_UNIQUE_ID_MAX             19
#define AISG_XID_VENDOR_WILDCARD           0xFFFF

typedef enum {
    AISG_L2_NO_ADDRESS = 0,
    AISG_L2_ADDRESS_ASSIGNED,
    AISG_L2_CONNECTED
} AisgL2State;

typedef enum {
    AISG_ERROR_NONE = 0,
    AISG_ERROR_TRANSPORT,
    AISG_ERROR_TIMEOUT,
    AISG_ERROR_MULTIPLE_DEVICES,
    AISG_ERROR_UNSUPPORTED_DEVICE_TYPE,
    AISG_ERROR_UNSUPPORTED_PROTOCOL_VERSION,
    AISG_ERROR_LINK_NOT_CONNECTED,
    AISG_ERROR_FRAME_REJECT,
    AISG_ERROR_RECEIVER_NOT_READY,
    AISG_ERROR_PROTOCOL
} AisgError;

typedef struct {
    SerialPort *serial;
    uint8_t deviceAddress;
    uint8_t ns;
    uint8_t nr;
    AisgL2State state;
    bool has3gppRelease;
    uint8_t negotiated3gppRelease;
    bool hasAisgVersion;
    uint8_t negotiatedAisgVersion;
    size_t maxInfoLen;
    AisgError lastError;
} AisgBus;

typedef struct {
    bool present;
    bool unsupported;
    uint8_t address;
    uint8_t deviceType;
    uint8_t uniqueId[AISG_XID_UNIQUE_ID_MAX];
    size_t uniqueIdLen;
    bool hasVendorCode;
    uint16_t vendorCode;
    char model[64];
} AisgDevice;

void aisg_v2_bus_init(AisgBus *bus, SerialPort *serial);
void aisg_v2_bus_reset_link(AisgBus *bus);
const char *aisg_v2_l2_state_str(AisgL2State state);
const char *aisg_v2_error_str(AisgError error);
bool aisg_v2_scan(AisgBus *bus, AisgDevice *device);
bool aisg_v2_disconnect(AisgBus *bus);
bool aisg_v2_send_retap(AisgBus *bus,
                        RetapRequest *request,
                        RetapResponse *response);

#endif /* AISG_V2_H_ */

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
#define AISG_ADDR_PRIMARY                  0x01

#define AISG_CTRL_I_PF                     0x10
#define AISG_CTRL_XID_PF                   0xBF
#define AISG_CTRL_XID_BASE                 0xAF
#define AISG_CTRL_SNRM_PF                  0x93
#define AISG_CTRL_UA_PF                    0x73
#define AISG_CTRL_UA_BASE                  0x63
#define AISG_CTRL_DM_PF                    0x1F
#define AISG_CTRL_DM_BASE                  0x0F
#define AISG_CTRL_FRMR_BASE                0x87

#define AISG_XID_FI_STD                    0x81
#define AISG_XID_GI_ADDR                   0xF0
#define AISG_XID_PI_UNIQUE_ID              0x01
#define AISG_XID_PI_HDLC_ADDR              0x02
#define AISG_XID_PI_BIT_MASK               0x03
#define AISG_XID_PI_DEVICE_TYPE            0x04
#define AISG_XID_PI_3GPP_RELEASE           0x05
#define AISG_XID_PI_VENDOR_CODE            0x06
#define AISG_XID_PI_AISG_VERSION           0x14

#define AISG_DEVICE_TYPE_RET_SINGLE        0x01
#define AISG_DEVICE_TYPE_RET_MULTI         0x11
#define AISG_DEVICE_TYPE_TMA               0x02

#define AISG_DEFAULT_TIMEOUT_MS            3000
#define AISG_SHORT_TIMEOUT_MS              900
#define AISG_LONG_TIMEOUT_MS               5000
#define AISG_SCAN_HOLD_MS                  350
#define AISG_MAX_UID_LEN                   19

/*
 * AISG/3GPP HDLC link state for one secondary station.  Window size is kept
 * at the mandatory value of 1; sequence numbers are reset by SNRM/UA.
 */
typedef struct {
    SerialPort *serial;
    uint8_t deviceAddress;
    uint8_t txSeq;
    uint8_t rxSeq;
    bool linked;
} AisgBus;

typedef struct {
    bool present;
    bool linked;
    uint8_t address;
    uint8_t deviceType;
    int baud;
    uint8_t uniqueId[AISG_MAX_UID_LEN];
    size_t uniqueIdLen;
    char uniqueIdHex[(AISG_MAX_UID_LEN * 2) + 1];
    char model[64];
} AisgDevice;

void aisg_v2_bus_init(AisgBus *bus, SerialPort *serial);
bool aisg_v2_scan(AisgBus *bus, AisgDevice *device);
bool aisg_v2_send_retap(AisgBus *bus,
                        RetapRequest *request,
                        RetapResponse *response);

#endif /* AISG_V2_H_ */

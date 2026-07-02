/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AISG_XID_H_
#define AISG_XID_H_

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

/*
 * AISG v2 / 3GPP TS 25.462 XID helpers.
 *
 * This file intentionally lives under ctrl/inc for now.  The emulator's
 * protocol-accurate RET mode includes this header directly instead of moving
 * the current ctrl protocol code into a common directory.
 */
#ifndef AISG_XID_FI
#define AISG_XID_FI                        0x81
#endif
#ifndef AISG_XID_GI_ADDRESSING
#define AISG_XID_GI_ADDRESSING             0xF0
#endif
#ifndef AISG_XID_PI_UNIQUE_ID
#define AISG_XID_PI_UNIQUE_ID              0x01
#endif
#ifndef AISG_XID_PI_HDLC_ADDRESS
#define AISG_XID_PI_HDLC_ADDRESS           0x02
#endif
#ifndef AISG_XID_PI_BIT_MASK
#define AISG_XID_PI_BIT_MASK               0x03
#endif
#ifndef AISG_XID_PI_DEVICE_TYPE
#define AISG_XID_PI_DEVICE_TYPE            0x04
#endif
#ifndef AISG_XID_PI_3GPP_RELEASE
#define AISG_XID_PI_3GPP_RELEASE           0x05
#endif
#ifndef AISG_XID_PI_VENDOR_CODE
#define AISG_XID_PI_VENDOR_CODE            0x06
#endif
#ifndef AISG_XID_PI_AISG_VERSION
#define AISG_XID_PI_AISG_VERSION           20
#endif
#ifndef AISG_XID_UNIQUE_ID_MAX
#define AISG_XID_UNIQUE_ID_MAX             19
#endif
#ifndef AISG_XID_VENDOR_WILDCARD
#define AISG_XID_VENDOR_WILDCARD           0xFFFF
#endif

#define AISG_XID_INFO_MIN_LEN              3
#define AISG_XID_SCAN_ID_LEN               19

typedef struct {
    bool hasUniqueId;
    uint8_t uniqueId[AISG_XID_UNIQUE_ID_MAX];
    size_t uniqueIdLen;

    bool hasAddress;
    uint8_t address;

    bool hasMask;
    uint8_t mask[AISG_XID_UNIQUE_ID_MAX];
    size_t maskLen;

    bool hasDeviceType;
    uint8_t deviceType;

    bool has3gppRelease;
    uint8_t release;

    bool hasAisgVersion;
    uint8_t aisgVersion;

    bool hasVendorCode;
    uint16_t vendorCode;
} XidAddressingParams;

void xid_params_init(XidAddressingParams *params);

bool xid_begin_addressing_info(uint8_t *info, size_t size, size_t *off);
bool xid_finish_info(uint8_t *info, size_t off, size_t *len);
bool xid_append_param(uint8_t *buf,
                      size_t size,
                      size_t *off,
                      uint8_t pi,
                      const uint8_t *pv,
                      size_t pvLen);

bool xid_parse_addressing_info(const uint8_t *info,
                               size_t infoLen,
                               XidAddressingParams *params);

bool xid_build_scan_info(uint8_t *info, size_t size, size_t *len);
bool xid_build_assign_info(const uint8_t *uniqueId,
                           size_t uniqueIdLen,
                           uint8_t assignedAddress,
                           uint8_t deviceType,
                           bool hasVendorCode,
                           uint16_t vendorCode,
                           uint8_t *info,
                           size_t size,
                           size_t *len);
bool xid_build_one_octet_info(uint8_t pi,
                              uint8_t value,
                              uint8_t *info,
                              size_t size,
                              size_t *len);
bool xid_build_device_response_info(const uint8_t *uniqueId,
                                    size_t uniqueIdLen,
                                    uint8_t address,
                                    uint8_t deviceType,
                                    bool hasVendorCode,
                                    uint16_t vendorCode,
                                    uint8_t *info,
                                    size_t size,
                                    size_t *len);

bool xid_unique_id_mask_match(const uint8_t *deviceUniqueId,
                              size_t deviceUniqueIdLen,
                              const uint8_t *scanUniqueId,
                              size_t scanUniqueIdLen,
                              const uint8_t *mask,
                              size_t maskLen);
bool xid_assignment_matches(const XidAddressingParams *params,
                            const uint8_t *deviceUniqueId,
                            size_t deviceUniqueIdLen,
                            uint8_t deviceType,
                            uint16_t vendorCode);

#endif /* AISG_XID_H_ */

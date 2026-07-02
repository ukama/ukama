/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "xid.h"

void xid_params_init(XidAddressingParams *params)
{
    if (params == NULL) {
        return;
    }

    memset(params, 0, sizeof(*params));
}

bool xid_begin_addressing_info(uint8_t *info, size_t size, size_t *off)
{
    if (info == NULL || off == NULL || size < AISG_XID_INFO_MIN_LEN) {
        return false;
    }

    *off = 0;
    info[(*off)++] = AISG_XID_FI;
    info[(*off)++] = AISG_XID_GI_ADDRESSING;
    info[(*off)++] = 0x00; /* GL filled by xid_finish_info(). */

    return true;
}

bool xid_finish_info(uint8_t *info, size_t off, size_t *len)
{
    size_t gl;

    if (info == NULL || len == NULL || off < AISG_XID_INFO_MIN_LEN) {
        return false;
    }

    gl = off - AISG_XID_INFO_MIN_LEN;
    if (gl > 255) {
        return false;
    }

    info[2] = (uint8_t)gl;
    *len = off;

    return true;
}

bool xid_append_param(uint8_t *buf,
                      size_t size,
                      size_t *off,
                      uint8_t pi,
                      const uint8_t *pv,
                      size_t pvLen)
{
    if (buf == NULL || off == NULL || pv == NULL) {
        return false;
    }

    if (pvLen == 0 || pvLen > 255) {
        return false;
    }

    if (*off + 2 + pvLen > size) {
        return false;
    }

    buf[(*off)++] = pi;
    buf[(*off)++] = (uint8_t)pvLen;
    memcpy(&buf[*off], pv, pvLen);
    *off += pvLen;

    return true;
}

bool xid_parse_addressing_info(const uint8_t *info,
                               size_t infoLen,
                               XidAddressingParams *params)
{
    size_t pos;
    size_t end;
    uint8_t gl;
    uint8_t pi;
    uint8_t pl;
    const uint8_t *pv;

    if (info == NULL || params == NULL) {
        return false;
    }

    xid_params_init(params);

    if (infoLen < AISG_XID_INFO_MIN_LEN) {
        return false;
    }

    if (info[0] != AISG_XID_FI || info[1] != AISG_XID_GI_ADDRESSING) {
        return false;
    }

    gl = info[2];
    if ((size_t)gl > infoLen - AISG_XID_INFO_MIN_LEN) {
        return false;
    }

    pos = AISG_XID_INFO_MIN_LEN;
    end = AISG_XID_INFO_MIN_LEN + (size_t)gl;

    while (pos < end) {
        if (pos + 2 > end) {
            return false;
        }

        pi = info[pos++];
        pl = info[pos++];

        if (pl == 0 || pos + (size_t)pl > end) {
            return false;
        }

        pv = &info[pos];
        pos += (size_t)pl;

        switch (pi) {
        case AISG_XID_PI_UNIQUE_ID:
            if (pl > AISG_XID_UNIQUE_ID_MAX) {
                return false;
            }
            params->hasUniqueId = true;
            params->uniqueIdLen = pl;
            memcpy(params->uniqueId, pv, pl);
            break;

        case AISG_XID_PI_HDLC_ADDRESS:
            if (pl != 1) {
                return false;
            }
            params->hasAddress = true;
            params->address = pv[0];
            break;

        case AISG_XID_PI_BIT_MASK:
            if (pl > AISG_XID_UNIQUE_ID_MAX) {
                return false;
            }
            params->hasMask = true;
            params->maskLen = pl;
            memcpy(params->mask, pv, pl);
            break;

        case AISG_XID_PI_DEVICE_TYPE:
            if (pl != 1) {
                return false;
            }
            params->hasDeviceType = true;
            params->deviceType = pv[0];
            break;

        case AISG_XID_PI_3GPP_RELEASE:
            if (pl != 1) {
                return false;
            }
            params->has3gppRelease = true;
            params->release = pv[0];
            break;

        case AISG_XID_PI_VENDOR_CODE:
            if (pl != 2) {
                return false;
            }
            params->hasVendorCode = true;
            params->vendorCode = (uint16_t)(((uint16_t)pv[0] << 8) | pv[1]);
            break;

        case AISG_XID_PI_AISG_VERSION:
            if (pl != 1) {
                return false;
            }
            params->hasAisgVersion = true;
            params->aisgVersion = pv[0];
            break;

        default:
            /* Unknown XID parameters are ignored by design. */
            break;
        }
    }

    return pos == end;
}

bool xid_build_scan_info(uint8_t *info, size_t size, size_t *len)
{
    uint8_t uid[AISG_XID_SCAN_ID_LEN];
    uint8_t mask[AISG_XID_SCAN_ID_LEN];
    size_t off;

    if (info == NULL || len == NULL) {
        return false;
    }

    memset(uid, 0, sizeof(uid));
    memset(mask, 0, sizeof(mask));

    if (!xid_begin_addressing_info(info, size, &off)) {
        return false;
    }

    if (!xid_append_param(info,
                          size,
                          &off,
                          AISG_XID_PI_UNIQUE_ID,
                          uid,
                          sizeof(uid))) {
        return false;
    }

    if (!xid_append_param(info,
                          size,
                          &off,
                          AISG_XID_PI_BIT_MASK,
                          mask,
                          sizeof(mask))) {
        return false;
    }

    return xid_finish_info(info, off, len);
}

bool xid_build_assign_info(const uint8_t *uniqueId,
                           size_t uniqueIdLen,
                           uint8_t assignedAddress,
                           uint8_t deviceType,
                           bool hasVendorCode,
                           uint16_t vendorCode,
                           uint8_t *info,
                           size_t size,
                           size_t *len)
{
    uint8_t address[1];
    uint8_t type[1];
    uint8_t vendor[2];
    size_t off;

    if (uniqueId == NULL || uniqueIdLen == 0 ||
        uniqueIdLen > AISG_XID_UNIQUE_ID_MAX ||
        info == NULL || len == NULL) {
        return false;
    }

    address[0] = assignedAddress;
    type[0] = deviceType;
    vendor[0] = (uint8_t)((vendorCode >> 8) & 0xFF);
    vendor[1] = (uint8_t)(vendorCode & 0xFF);

    if (!xid_begin_addressing_info(info, size, &off)) {
        return false;
    }

    if (!xid_append_param(info,
                          size,
                          &off,
                          AISG_XID_PI_UNIQUE_ID,
                          uniqueId,
                          uniqueIdLen)) {
        return false;
    }

    if (!xid_append_param(info,
                          size,
                          &off,
                          AISG_XID_PI_HDLC_ADDRESS,
                          address,
                          sizeof(address))) {
        return false;
    }

    if (!xid_append_param(info,
                          size,
                          &off,
                          AISG_XID_PI_DEVICE_TYPE,
                          type,
                          sizeof(type))) {
        return false;
    }

    if (hasVendorCode &&
        !xid_append_param(info,
                          size,
                          &off,
                          AISG_XID_PI_VENDOR_CODE,
                          vendor,
                          sizeof(vendor))) {
        return false;
    }

    return xid_finish_info(info, off, len);
}

bool xid_build_one_octet_info(uint8_t pi,
                              uint8_t value,
                              uint8_t *info,
                              size_t size,
                              size_t *len)
{
    uint8_t pv[1];
    size_t off;

    if (info == NULL || len == NULL) {
        return false;
    }

    pv[0] = value;

    if (!xid_begin_addressing_info(info, size, &off)) {
        return false;
    }

    if (!xid_append_param(info, size, &off, pi, pv, sizeof(pv))) {
        return false;
    }

    return xid_finish_info(info, off, len);
}

bool xid_build_device_response_info(const uint8_t *uniqueId,
                                    size_t uniqueIdLen,
                                    uint8_t address,
                                    uint8_t deviceType,
                                    bool hasVendorCode,
                                    uint16_t vendorCode,
                                    uint8_t *info,
                                    size_t size,
                                    size_t *len)
{
    uint8_t addr[1];
    uint8_t type[1];
    uint8_t vendor[2];
    size_t off;

    if (uniqueId == NULL || uniqueIdLen == 0 ||
        uniqueIdLen > AISG_XID_UNIQUE_ID_MAX ||
        info == NULL || len == NULL) {
        return false;
    }

    addr[0] = address;
    type[0] = deviceType;
    vendor[0] = (uint8_t)((vendorCode >> 8) & 0xFF);
    vendor[1] = (uint8_t)(vendorCode & 0xFF);

    if (!xid_begin_addressing_info(info, size, &off)) {
        return false;
    }

    if (!xid_append_param(info,
                          size,
                          &off,
                          AISG_XID_PI_UNIQUE_ID,
                          uniqueId,
                          uniqueIdLen)) {
        return false;
    }

    if (!xid_append_param(info,
                          size,
                          &off,
                          AISG_XID_PI_HDLC_ADDRESS,
                          addr,
                          sizeof(addr))) {
        return false;
    }

    if (!xid_append_param(info,
                          size,
                          &off,
                          AISG_XID_PI_DEVICE_TYPE,
                          type,
                          sizeof(type))) {
        return false;
    }

    if (hasVendorCode &&
        !xid_append_param(info,
                          size,
                          &off,
                          AISG_XID_PI_VENDOR_CODE,
                          vendor,
                          sizeof(vendor))) {
        return false;
    }

    return xid_finish_info(info, off, len);
}

static void padded_unique_id(const uint8_t *uniqueId,
                             size_t uniqueIdLen,
                             uint8_t out[AISG_XID_UNIQUE_ID_MAX])
{
    memset(out, 0x20, AISG_XID_UNIQUE_ID_MAX);

    if (uniqueId == NULL || uniqueIdLen == 0) {
        return;
    }

    if (uniqueIdLen > AISG_XID_UNIQUE_ID_MAX) {
        uniqueIdLen = AISG_XID_UNIQUE_ID_MAX;
    }

    memcpy(out, uniqueId, uniqueIdLen);
}

bool xid_unique_id_mask_match(const uint8_t *deviceUniqueId,
                              size_t deviceUniqueIdLen,
                              const uint8_t *scanUniqueId,
                              size_t scanUniqueIdLen,
                              const uint8_t *mask,
                              size_t maskLen)
{
    uint8_t padded[AISG_XID_UNIQUE_ID_MAX];
    size_t i;

    if (deviceUniqueId == NULL || scanUniqueId == NULL || mask == NULL) {
        return false;
    }

    if (scanUniqueIdLen == 0 || scanUniqueIdLen != maskLen ||
        scanUniqueIdLen > AISG_XID_UNIQUE_ID_MAX) {
        return false;
    }

    padded_unique_id(deviceUniqueId, deviceUniqueIdLen, padded);

    for (i = 0; i < scanUniqueIdLen; i++) {
        if ((padded[i] & mask[i]) != scanUniqueId[i]) {
            return false;
        }
    }

    return true;
}

static bool right_match(const uint8_t *lhs,
                        size_t lhsLen,
                        const uint8_t *rhs,
                        size_t rhsLen)
{
    if (lhs == NULL || rhs == NULL || rhsLen > lhsLen) {
        return false;
    }

    return memcmp(lhs + (lhsLen - rhsLen), rhs, rhsLen) == 0;
}

bool xid_assignment_matches(const XidAddressingParams *params,
                            const uint8_t *deviceUniqueId,
                            size_t deviceUniqueIdLen,
                            uint8_t deviceType,
                            uint16_t vendorCode)
{
    if (params == NULL || deviceUniqueId == NULL || deviceUniqueIdLen == 0) {
        return false;
    }

    if (params->hasMask) {
        return false;
    }

    if (!params->hasAddress) {
        return false;
    }

    if (params->hasUniqueId &&
        !right_match(deviceUniqueId,
                     deviceUniqueIdLen,
                     params->uniqueId,
                     params->uniqueIdLen)) {
        return false;
    }

    if (params->hasDeviceType &&
        params->deviceType != deviceType &&
        params->deviceType != 0xFF) {
        return false;
    }

    if (params->hasVendorCode &&
        params->vendorCode != vendorCode &&
        params->vendorCode != AISG_XID_VENDOR_WILDCARD) {
        return false;
    }

    return true;
}

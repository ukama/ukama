/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <unistd.h>

#include "aisg_v2.h"
#include "hdlc.h"
#include "usys_log.h"

#define AISG_TRACE_BYTES_PER_LINE          16

#define XID_INFO_MIN_LEN                   3
#define XID_SCAN_ID_LEN                    AISG_XID_UNIQUE_ID_MAX

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

static const char *l2_state_name(AisgL2State state)
{
    switch (state) {
    case AISG_L2_NO_ADDRESS:       return "NO_ADDRESS";
    case AISG_L2_ADDRESS_ASSIGNED: return "ADDRESS_ASSIGNED";
    case AISG_L2_CONNECTED:        return "CONNECTED";
    default:                       return "UNKNOWN";
    }
}

static const char *ctrl_name(uint8_t ctrl)
{
    if (hdlc_is_i_frame(ctrl)) return "I";
    if (hdlc_is_xid(ctrl))     return "XID";
    if (hdlc_is_snrm(ctrl))    return "SNRM";
    if (hdlc_is_disc(ctrl))    return "DISC";
    if (hdlc_is_ua(ctrl))      return "UA";
    if (hdlc_is_dm(ctrl))      return "DM";
    if (hdlc_is_rr(ctrl))      return "RR";
    if (hdlc_is_rnr(ctrl))     return "RNR";
    if (hdlc_is_frmr(ctrl))    return "FRMR";

    return "CTRL";
}

static void log_hex_bytes(const char *label,
                          const uint8_t *data,
                          size_t len)
{
    char line[(AISG_TRACE_BYTES_PER_LINE * 3) + 1];
    size_t off;
    size_t i;
    size_t n;
    size_t pos;
    int written;

    if (label == NULL || data == NULL) {
        return;
    }

    usys_log_debug("aisg: %s len=%zu", label, len);

    off = 0;
    while (off < len) {
        n = len - off;
        if (n > AISG_TRACE_BYTES_PER_LINE) {
            n = AISG_TRACE_BYTES_PER_LINE;
        }

        pos = 0;
        memset(line, 0, sizeof(line));

        for (i = 0; i < n; i++) {
            written = snprintf(&line[pos],
                               sizeof(line) - pos,
                               "%02X%s",
                               data[off + i],
                               (i + 1 == n) ? "" : " ");
            if (written <= 0) {
                break;
            }

            pos += (size_t)written;
            if (pos >= sizeof(line)) {
                break;
            }
        }

        usys_log_debug("aisg: %s[%04zu..%04zu] %s",
                       label,
                       off,
                       off + n - 1,
                       line);

        off += n;
    }
}

static void log_frame(const char *label, const HdlcFrame *frame)
{
    if (label == NULL || frame == NULL) {
        return;
    }

    usys_log_debug("aisg: %s addr=0x%02X ctrl=0x%02X(%s) pf=%u "
                   "ns=%u nr=%u info_len=%zu",
                   label,
                   frame->address,
                   frame->control,
                   ctrl_name(frame->control),
                   hdlc_poll_final(frame->control) ? 1 : 0,
                   hdlc_ns(frame->control),
                   hdlc_nr(frame->control),
                   frame->infoLen);

    if (frame->infoLen > 0) {
        log_hex_bytes(label, frame->info, frame->infoLen);
    }
}

static bool same_frame(const HdlcFrame *a, const HdlcFrame *b)
{
    if (a == NULL || b == NULL) {
        return false;
    }

    if (a->address != b->address ||
        a->control != b->control ||
        a->infoLen != b->infoLen) {
        return false;
    }

    if (a->infoLen == 0) {
        return true;
    }

    return memcmp(a->info, b->info, a->infoLen) == 0;
}

static bool read_decoded_frame(AisgBus *bus,
                               const HdlcFrame *txFrame,
                               HdlcFrame *rxFrame,
                               int timeoutMs)
{
    uint8_t raw[HDLC_MAX_FRAME];
    size_t rawLen;
    int attempt;

    if (bus == NULL || rxFrame == NULL) {
        return false;
    }

    for (attempt = 0; attempt < AISG_MAX_RX_ATTEMPTS; attempt++) {
        memset(raw, 0, sizeof(raw));
        rawLen = 0;

        if (!serial_read_frame(bus->serial,
                               raw,
                               sizeof(raw),
                               &rawLen,
                               timeoutMs)) {
            usys_log_debug("aisg: RX timeout attempt=%d", attempt + 1);
            return false;
        }

        log_hex_bytes("RX hdlc", raw, rawLen);

        memset(rxFrame, 0, sizeof(*rxFrame));
        if (!hdlc_decode_frame(raw, rawLen, rxFrame)) {
            usys_log_debug("aisg: RX hdlc decode failed attempt=%d",
                           attempt + 1);
            continue;
        }

        log_frame("RX frame", rxFrame);

        if (txFrame != NULL && same_frame(txFrame, rxFrame)) {
            usys_log_debug("aisg: RX rejected local echo attempt=%d",
                           attempt + 1);
            continue;
        }

        return true;
    }

    usys_log_debug("aisg: RX failed after echo/invalid frames attempts=%d",
                   AISG_MAX_RX_ATTEMPTS);

    return false;
}

static bool send_frame(AisgBus *bus,
                       const HdlcFrame *txFrame,
                       HdlcFrame *rxFrame,
                       int timeoutMs)
{
    uint8_t raw[HDLC_MAX_FRAME];
    size_t rawLen;

    if (bus == NULL || bus->serial == NULL || txFrame == NULL) {
        return false;
    }

    log_frame("TX frame", txFrame);

    memset(raw, 0, sizeof(raw));
    rawLen = 0;

    if (!hdlc_encode_frame(txFrame, raw, sizeof(raw), &rawLen)) {
        usys_log_debug("aisg: TX hdlc encode failed");
        return false;
    }

    log_hex_bytes("TX hdlc", raw, rawLen);

    /*
     * TS 25.462 requires at least 3 ms between receive and transmit.
     * Waiting before every command is a small cost and keeps the bus timing
     * conservative while the state machine is still simple.
     */
    usleep(AISG_MIN_TURNAROUND_US);

    if (!serial_write_all(bus->serial, raw, rawLen)) {
        usys_log_debug("aisg: TX serial write failed");
        return false;
    }

    if (rxFrame == NULL) {
        return true;
    }

    return read_decoded_frame(bus, txFrame, rxFrame, timeoutMs);
}

static bool append_xid_param(uint8_t *buf,
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

static bool begin_xid_addressing_info(uint8_t *info,
                                      size_t size,
                                      size_t *off)
{
    if (info == NULL || off == NULL || size < XID_INFO_MIN_LEN) {
        return false;
    }

    *off = 0;
    info[(*off)++] = AISG_XID_FI;
    info[(*off)++] = AISG_XID_GI_ADDRESSING;
    info[(*off)++] = 0x00; /* GL, filled by finish_xid_info(). */

    return true;
}

static bool finish_xid_info(uint8_t *info, size_t off)
{
    size_t gl;

    if (info == NULL || off < XID_INFO_MIN_LEN) {
        return false;
    }

    gl = off - XID_INFO_MIN_LEN;
    if (gl > 255) {
        return false;
    }

    info[2] = (uint8_t)gl;

    return true;
}

static bool build_xid_scan_info(uint8_t *info, size_t size, size_t *len)
{
    uint8_t uid[XID_SCAN_ID_LEN];
    uint8_t mask[XID_SCAN_ID_LEN];
    size_t off;

    if (info == NULL || len == NULL) {
        return false;
    }

    memset(uid, 0, sizeof(uid));
    memset(mask, 0, sizeof(mask));

    if (!begin_xid_addressing_info(info, size, &off)) {
        return false;
    }

    if (!append_xid_param(info,
                          size,
                          &off,
                          AISG_XID_PI_UNIQUE_ID,
                          uid,
                          sizeof(uid))) {
        return false;
    }

    if (!append_xid_param(info,
                          size,
                          &off,
                          AISG_XID_PI_BIT_MASK,
                          mask,
                          sizeof(mask))) {
        return false;
    }

    if (!finish_xid_info(info, off)) {
        return false;
    }

    *len = off;

    return true;
}

static bool build_xid_assign_info(const AisgDevice *device,
                                  uint8_t assignedAddress,
                                  uint8_t *info,
                                  size_t size,
                                  size_t *len)
{
    uint8_t address[1];
    uint8_t deviceType[1];
    size_t off;

    if (device == NULL || info == NULL || len == NULL) {
        return false;
    }

    if (device->uniqueIdLen == 0 ||
        device->uniqueIdLen > AISG_XID_UNIQUE_ID_MAX) {
        return false;
    }

    address[0] = assignedAddress;
    deviceType[0] = AISG_SUPPORTED_DEVICE_TYPE;

    if (!begin_xid_addressing_info(info, size, &off)) {
        return false;
    }

    if (!append_xid_param(info,
                          size,
                          &off,
                          AISG_XID_PI_UNIQUE_ID,
                          device->uniqueId,
                          device->uniqueIdLen)) {
        return false;
    }

    if (!append_xid_param(info,
                          size,
                          &off,
                          AISG_XID_PI_HDLC_ADDRESS,
                          address,
                          sizeof(address))) {
        return false;
    }

    if (!append_xid_param(info,
                          size,
                          &off,
                          AISG_XID_PI_DEVICE_TYPE,
                          deviceType,
                          sizeof(deviceType))) {
        return false;
    }

    if (!finish_xid_info(info, off)) {
        return false;
    }

    *len = off;

    return true;
}

static bool parse_xid_addressing_info(const uint8_t *info,
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

    memset(params, 0, sizeof(*params));

    if (infoLen < XID_INFO_MIN_LEN) {
        return false;
    }

    if (info[0] != AISG_XID_FI || info[1] != AISG_XID_GI_ADDRESSING) {
        return false;
    }

    gl = info[2];
    if ((size_t)gl > infoLen - XID_INFO_MIN_LEN) {
        return false;
    }

    pos = XID_INFO_MIN_LEN;
    end = XID_INFO_MIN_LEN + (size_t)gl;

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

        case AISG_XID_PI_AISG_VERSION:
            if (pl != 1) {
                return false;
            }
            params->hasAisgVersion = true;
            params->aisgVersion = pv[0];
            break;

        case AISG_XID_PI_VENDOR_CODE:
            if (pl != 2) {
                return false;
            }
            params->hasVendorCode = true;
            params->vendorCode = (uint16_t)(((uint16_t)pv[0] << 8) | pv[1]);
            break;

        default:
            usys_log_debug("aisg: XID ignoring unsupported PI=0x%02X len=%u",
                           pi,
                           pl);
            break;
        }
    }

    return pos == end;
}

static void apply_xid_params_to_device(const XidAddressingParams *params,
                                       AisgDevice *device)
{
    if (params == NULL || device == NULL) {
        return;
    }

    if (params->hasUniqueId) {
        device->uniqueIdLen = params->uniqueIdLen;
        memcpy(device->uniqueId, params->uniqueId, params->uniqueIdLen);
    }

    if (params->hasAddress) {
        device->address = params->address;
    }

    if (params->hasDeviceType) {
        device->deviceType = params->deviceType;
    }

    if (params->hasVendorCode) {
        device->hasVendorCode = true;
        device->vendorCode = params->vendorCode;
    }
}

static bool read_extra_scan_response(AisgBus *bus)
{
    HdlcFrame extra;

    memset(&extra, 0, sizeof(extra));

    if (!read_decoded_frame(bus, NULL, &extra, AISG_SCAN_EXTRA_TIMEOUT_MS)) {
        return false;
    }

    usys_log_debug("aisg: scan saw extra response addr=0x%02X ctrl=0x%02X; "
                   "multiple devices/collision unsupported",
                   extra.address,
                   extra.control);

    return true;
}

static bool xid_scan_single(AisgBus *bus, AisgDevice *device)
{
    HdlcFrame tx;
    HdlcFrame rx;
    XidAddressingParams params;
    size_t infoLen;

    if (bus == NULL || device == NULL) {
        return false;
    }

    memset(&tx, 0, sizeof(tx));
    memset(&rx, 0, sizeof(rx));

    tx.address = AISG_ADDR_BROADCAST;
    tx.control = hdlc_xid_ctrl(true);

    if (!build_xid_scan_info(tx.info, sizeof(tx.info), &infoLen)) {
        usys_log_debug("aisg: XID scan build failed");
        return false;
    }
    tx.infoLen = infoLen;

    usys_log_debug("aisg: XID scan start uid_len=%d mask=all-zero scope=%s",
                   XID_SCAN_ID_LEN,
                   AISG_SCOPE_NAME);

    if (!send_frame(bus, &tx, &rx, AISG_DEFAULT_TIMEOUT_MS)) {
        usys_log_debug("aisg: XID scan failed: no valid response");
        return false;
    }

    if (!hdlc_is_xid(rx.control)) {
        usys_log_debug("aisg: XID scan failed: unexpected ctrl=0x%02X",
                       rx.control);
        return false;
    }

    if (!parse_xid_addressing_info(rx.info, rx.infoLen, &params)) {
        usys_log_debug("aisg: XID scan failed: malformed XID response");
        return false;
    }

    if (!params.hasUniqueId || !params.hasAddress || !params.hasDeviceType) {
        usys_log_debug("aisg: XID scan failed: missing required response "
                       "fields uid=%u addr=%u type=%u",
                       params.hasUniqueId ? 1 : 0,
                       params.hasAddress ? 1 : 0,
                       params.hasDeviceType ? 1 : 0);
        return false;
    }

    apply_xid_params_to_device(&params, device);

    if (read_extra_scan_response(bus)) {
        device->unsupported = true;
        device->present = false;
        return false;
    }

    log_hex_bytes("XID unique-id", device->uniqueId, device->uniqueIdLen);
    usys_log_debug("aisg: XID scan response addr=0x%02X device_type=0x%02X "
                   "vendor=%s0x%04X",
                   device->address,
                   device->deviceType,
                   device->hasVendorCode ? "" : "unknown/",
                   device->hasVendorCode ? device->vendorCode : 0);

    if (device->deviceType != AISG_SUPPORTED_DEVICE_TYPE) {
        device->unsupported = true;
        device->present = false;
        usys_log_debug("aisg: unsupported device_type=0x%02X expected=0x%02X",
                       device->deviceType,
                       AISG_SUPPORTED_DEVICE_TYPE);
        return false;
    }

    return true;
}

static bool xid_assign_address(AisgBus *bus, AisgDevice *device)
{
    HdlcFrame tx;
    HdlcFrame rx;
    XidAddressingParams params;
    size_t infoLen;

    if (bus == NULL || device == NULL) {
        return false;
    }

    memset(&tx, 0, sizeof(tx));
    memset(&rx, 0, sizeof(rx));

    tx.address = AISG_ADDR_BROADCAST;
    tx.control = hdlc_xid_ctrl(true);

    if (!build_xid_assign_info(device,
                               AISG_ADDR_ASSIGNED,
                               tx.info,
                               sizeof(tx.info),
                               &infoLen)) {
        usys_log_debug("aisg: XID address assignment build failed");
        return false;
    }
    tx.infoLen = infoLen;

    usys_log_debug("aisg: XID assign address=0x%02X device_type=0x%02X",
                   AISG_ADDR_ASSIGNED,
                   AISG_SUPPORTED_DEVICE_TYPE);

    if (!send_frame(bus, &tx, &rx, AISG_DEFAULT_TIMEOUT_MS)) {
        usys_log_debug("aisg: XID address assignment failed: no response");
        return false;
    }

    if (!hdlc_is_xid(rx.control)) {
        usys_log_debug("aisg: XID address assignment failed: unexpected "
                       "ctrl=0x%02X",
                       rx.control);
        return false;
    }

    if (rx.address != AISG_ADDR_ASSIGNED) {
        usys_log_debug("aisg: XID address assignment failed: response "
                       "addr=0x%02X expected=0x%02X",
                       rx.address,
                       AISG_ADDR_ASSIGNED);
        return false;
    }

    if (!parse_xid_addressing_info(rx.info, rx.infoLen, &params)) {
        usys_log_debug("aisg: XID address assignment failed: malformed response");
        return false;
    }

    if (!params.hasUniqueId || !params.hasDeviceType) {
        usys_log_debug("aisg: XID address assignment failed: missing uid/type");
        return false;
    }

    if (params.uniqueIdLen != device->uniqueIdLen ||
        memcmp(params.uniqueId, device->uniqueId, device->uniqueIdLen) != 0) {
        usys_log_debug("aisg: XID address assignment failed: UID mismatch");
        return false;
    }

    if (params.deviceType != AISG_SUPPORTED_DEVICE_TYPE) {
        usys_log_debug("aisg: XID address assignment failed: unsupported "
                       "device_type=0x%02X",
                       params.deviceType);
        device->unsupported = true;
        return false;
    }

    device->address = AISG_ADDR_ASSIGNED;
    device->deviceType = params.deviceType;
    device->present = true;
    device->unsupported = false;
    snprintf(device->model, sizeof(device->model), "%s", "single-ret");

    bus->deviceAddress = device->address;
    bus->state = AISG_L2_ADDRESS_ASSIGNED;
    bus->ns = 0;
    bus->nr = 0;

    usys_log_debug("aisg: XID address assignment OK addr=0x%02X state=%s",
                   bus->deviceAddress,
                   l2_state_name(bus->state));

    return true;
}

static bool build_xid_one_octet_info(uint8_t pi,
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

    if (!begin_xid_addressing_info(info, size, &off)) {
        return false;
    }

    if (!append_xid_param(info, size, &off, pi, pv, sizeof(pv))) {
        return false;
    }

    if (!finish_xid_info(info, off)) {
        return false;
    }

    *len = off;

    return true;
}

static bool xid_negotiate_one_octet(AisgBus *bus,
                                    const char *name,
                                    uint8_t pi,
                                    uint8_t offered,
                                    uint8_t *accepted)
{
    HdlcFrame tx;
    HdlcFrame rx;
    XidAddressingParams params;
    size_t infoLen;
    bool hasValue;
    uint8_t value;

    if (bus == NULL || name == NULL || accepted == NULL) {
        return false;
    }

    if (bus->state != AISG_L2_ADDRESS_ASSIGNED) {
        usys_log_debug("aisg: %s negotiation rejected in state=%s",
                       name,
                       l2_state_name(bus->state));
        return false;
    }

    memset(&tx, 0, sizeof(tx));
    memset(&rx, 0, sizeof(rx));

    tx.address = bus->deviceAddress;
    tx.control = hdlc_xid_ctrl(true);

    if (!build_xid_one_octet_info(pi,
                                  offered,
                                  tx.info,
                                  sizeof(tx.info),
                                  &infoLen)) {
        usys_log_debug("aisg: %s negotiation build failed", name);
        return false;
    }
    tx.infoLen = infoLen;

    usys_log_debug("aisg: %s negotiation start pi=%u offered=0x%02X addr=0x%02X",
                   name,
                   pi,
                   offered,
                   bus->deviceAddress);

    if (!send_frame(bus, &tx, &rx, AISG_DEFAULT_TIMEOUT_MS)) {
        usys_log_debug("aisg: %s negotiation failed: no response", name);
        return false;
    }

    if (rx.address != bus->deviceAddress) {
        usys_log_debug("aisg: %s negotiation failed: addr=0x%02X expected=0x%02X",
                       name,
                       rx.address,
                       bus->deviceAddress);
        return false;
    }

    if (!hdlc_is_xid(rx.control)) {
        usys_log_debug("aisg: %s negotiation failed: unexpected ctrl=0x%02X",
                       name,
                       rx.control);
        return false;
    }

    if (!parse_xid_addressing_info(rx.info, rx.infoLen, &params)) {
        usys_log_debug("aisg: %s negotiation failed: malformed XID response",
                       name);
        return false;
    }

    hasValue = false;
    value = 0;

    if (pi == AISG_XID_PI_3GPP_RELEASE) {
        hasValue = params.has3gppRelease;
        value = params.release;
    } else if (pi == AISG_XID_PI_AISG_VERSION) {
        hasValue = params.hasAisgVersion;
        value = params.aisgVersion;
    }

    if (!hasValue) {
        usys_log_debug("aisg: %s negotiation failed: PI=%u missing in response",
                       name,
                       pi);
        return false;
    }

    if (value != offered) {
        usys_log_debug("aisg: %s negotiation failed: accepted=0x%02X expected=0x%02X",
                       name,
                       value,
                       offered);
        return false;
    }

    *accepted = value;

    usys_log_debug("aisg: %s negotiation OK accepted=0x%02X", name, value);

    return true;
}

static bool xid_negotiate_3gpp_release(AisgBus *bus)
{
    uint8_t accepted;

    if (!xid_negotiate_one_octet(bus,
                                 "3GPP release",
                                 AISG_XID_PI_3GPP_RELEASE,
                                 AISG_3GPP_RELEASE_ID,
                                 &accepted)) {
        return false;
    }

    bus->has3gppRelease = true;
    bus->negotiated3gppRelease = accepted;

    return true;
}

static bool xid_negotiate_aisg_version(AisgBus *bus)
{
    uint8_t accepted;

    if (!xid_negotiate_one_octet(bus,
                                 "AISG version",
                                 AISG_XID_PI_AISG_VERSION,
                                 AISG_PROTOCOL_VERSION,
                                 &accepted)) {
        return false;
    }

    bus->hasAisgVersion = true;
    bus->negotiatedAisgVersion = accepted;

    return true;
}

static bool l2_establish_link(AisgBus *bus)
{
    HdlcFrame tx;
    HdlcFrame rx;

    if (bus == NULL) {
        return false;
    }

    if (bus->state != AISG_L2_ADDRESS_ASSIGNED) {
        usys_log_debug("aisg: SNRM rejected in state=%s",
                       l2_state_name(bus->state));
        return false;
    }

    memset(&tx, 0, sizeof(tx));
    memset(&rx, 0, sizeof(rx));

    tx.address = bus->deviceAddress;
    tx.control = hdlc_snrm_ctrl(true);
    tx.infoLen = 0;

    usys_log_debug("aisg: link establishment start SNRM addr=0x%02X",
                   bus->deviceAddress);

    if (!send_frame(bus, &tx, &rx, AISG_DEFAULT_TIMEOUT_MS)) {
        usys_log_debug("aisg: link establishment failed: no UA response");
        return false;
    }

    if (rx.address != bus->deviceAddress) {
        usys_log_debug("aisg: link establishment failed: addr=0x%02X expected=0x%02X",
                       rx.address,
                       bus->deviceAddress);
        return false;
    }

    if (hdlc_is_dm(rx.control)) {
        usys_log_debug("aisg: link establishment failed: secondary returned DM");
        return false;
    }

    if (!hdlc_is_ua(rx.control)) {
        usys_log_debug("aisg: link establishment failed: unexpected ctrl=0x%02X",
                       rx.control);
        return false;
    }

    if (!hdlc_poll_final(rx.control)) {
        usys_log_debug("aisg: link establishment failed: UA missing final bit");
        return false;
    }

    if (rx.infoLen != 0) {
        usys_log_debug("aisg: link establishment failed: UA info_len=%zu expected=0",
                       rx.infoLen);
        return false;
    }

    bus->state = AISG_L2_CONNECTED;
    bus->ns = 0;
    bus->nr = 0;

    usys_log_debug("aisg: link establishment OK state=%s addr=0x%02X",
                   l2_state_name(bus->state),
                   bus->deviceAddress);

    return true;
}

void aisg_v2_bus_init(AisgBus *bus, SerialPort *serial)
{
    if (bus == NULL) {
        return;
    }

    memset(bus, 0, sizeof(AisgBus));

    bus->serial        = serial;
    bus->deviceAddress = AISG_ADDR_DEFAULT;
    bus->ns            = 0;
    bus->nr            = 0;
    bus->state         = AISG_L2_NO_ADDRESS;

    usys_log_debug("aisg: init scope=%s supported_device_type=0x%02X state=%s",
                   AISG_SCOPE_NAME,
                   AISG_SUPPORTED_DEVICE_TYPE,
                   l2_state_name(bus->state));
}

bool aisg_v2_scan(AisgBus *bus, AisgDevice *device)
{
    if (bus == NULL || device == NULL) {
        return false;
    }

    memset(device, 0, sizeof(AisgDevice));

    bus->deviceAddress = AISG_ADDR_DEFAULT;
    bus->state = AISG_L2_NO_ADDRESS;
    bus->ns = 0;
    bus->nr = 0;
    bus->has3gppRelease = false;
    bus->negotiated3gppRelease = 0;
    bus->hasAisgVersion = false;
    bus->negotiatedAisgVersion = 0;

    if (!xid_scan_single(bus, device)) {
        usys_log_debug("aisg: scan failed");
        return false;
    }

    if (!xid_assign_address(bus, device)) {
        usys_log_debug("aisg: address assignment failed");
        return false;
    }

    if (!xid_negotiate_3gpp_release(bus)) {
        usys_log_debug("aisg: 3GPP release negotiation failed");
        return false;
    }

    if (!xid_negotiate_aisg_version(bus)) {
        usys_log_debug("aisg: AISG version negotiation failed");
        return false;
    }

    if (!l2_establish_link(bus)) {
        usys_log_debug("aisg: link establishment failed");
        return false;
    }

    return true;
}

static bool l2_validate_i_response(AisgBus *bus,
                                   const HdlcFrame *rx,
                                   uint8_t txNs,
                                   uint8_t txNr)
{
    uint8_t expectedAck;
    uint8_t rxNs;
    uint8_t rxNr;

    if (bus == NULL || rx == NULL) {
        return false;
    }

    if (rx->address != bus->deviceAddress) {
        usys_log_debug("aisg: I-response addr=0x%02X expected=0x%02X",
                       rx->address,
                       bus->deviceAddress);
        return false;
    }

    if (hdlc_is_frmr(rx->control)) {
        usys_log_debug("aisg: I-response rejected: secondary returned FRMR");
        return false;
    }

    if (hdlc_is_rnr(rx->control)) {
        usys_log_debug("aisg: I-response rejected: secondary returned RNR");
        return false;
    }

    if (hdlc_is_rr(rx->control)) {
        usys_log_debug("aisg: I-response rejected: RR acknowledgement carried no RETAP data");
        return false;
    }

    if (!hdlc_is_i_frame(rx->control)) {
        usys_log_debug("aisg: I-response rejected: unexpected ctrl=0x%02X",
                       rx->control);
        return false;
    }

    if (!hdlc_poll_final(rx->control)) {
        usys_log_debug("aisg: I-response rejected: final bit not set");
        return false;
    }

    expectedAck = (uint8_t)((txNs + 1) & 0x07);
    rxNs = hdlc_ns(rx->control);
    rxNr = hdlc_nr(rx->control);

    if (rxNr != expectedAck) {
        usys_log_debug("aisg: I-response rejected: N(R)=%u expected_ack=%u",
                       rxNr,
                       expectedAck);
        return false;
    }

    if (rxNs != txNr) {
        usys_log_debug("aisg: I-response rejected: N(S)=%u expected=%u",
                       rxNs,
                       txNr);
        return false;
    }

    bus->ns = expectedAck;
    bus->nr = (uint8_t)((rxNs + 1) & 0x07);

    usys_log_debug("aisg: I-response sequence OK next_ns=%u next_nr=%u",
                   bus->ns,
                   bus->nr);

    return true;
}

bool aisg_v2_send_retap(AisgBus *bus,
                        RetapRequest *request,
                        RetapResponse *response)
{
    uint8_t retap[RETAP_MAX_ENCODED];
    HdlcFrame tx;
    HdlcFrame rx;
    size_t retapLen;
    uint8_t txNs;
    uint8_t txNr;
    int timeoutMs;

    if (bus == NULL || request == NULL || response == NULL) {
        return false;
    }

    if (bus->state != AISG_L2_CONNECTED) {
        usys_log_debug("aisg: RETAP rejected: link state=%s expected=CONNECTED",
                       l2_state_name(bus->state));
        return false;
    }

    if (!retap_encode_request(request, retap, sizeof(retap), &retapLen)) {
        usys_log_debug("aisg: failed to encode RETAP request procedure=0x%02X",
                       request->procedure);
        return false;
    }

    if (retapLen > HDLC_MAX_INFO) {
        usys_log_debug("aisg: RETAP encoded payload too large len=%zu max=%u",
                       retapLen,
                       HDLC_MAX_INFO);
        return false;
    }

    memset(&tx, 0, sizeof(tx));
    memset(&rx, 0, sizeof(rx));

    txNs = bus->ns;
    txNr = bus->nr;

    tx.address = bus->deviceAddress;
    tx.control = hdlc_i_ctrl(txNs, txNr, true);
    memcpy(tx.info, retap, retapLen);
    tx.infoLen = retapLen;

    timeoutMs = retap_request_timeout_ms(request);

    usys_log_debug("aisg: RETAP TX procedure=0x%02X data_len=%zu timeout_ms=%d ns=%u nr=%u",
                   request->procedure,
                   request->dataLen,
                   timeoutMs,
                   txNs,
                   txNr);

    if (!send_frame(bus, &tx, &rx, timeoutMs)) {
        usys_log_debug("aisg: RETAP transport failed procedure=0x%02X",
                       request->procedure);
        return false;
    }

    if (!l2_validate_i_response(bus, &rx, txNs, txNr)) {
        return false;
    }

    if (!retap_decode_response(rx.info, rx.infoLen, response)) {
        usys_log_debug("aisg: RETAP response decode failed procedure=0x%02X info_len=%zu",
                       request->procedure,
                       rx.infoLen);
        return false;
    }

    if (response->procedure != request->procedure) {
        usys_log_debug("aisg: RETAP response procedure mismatch got=0x%02X expected=0x%02X",
                       response->procedure,
                       request->procedure);
        return false;
    }

    usys_log_debug("aisg: RETAP RX procedure=0x%02X return=0x%02X failure=0x%02X data_len=%zu",
                   response->procedure,
                   response->returnCode,
                   response->failureReason,
                   response->dataLen);

    return true;
}

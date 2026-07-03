/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <ctype.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>

#include "aisg_v2.h"
#include "hdlc.h"
#include "usys_log.h"

#define AISG_TX_FRAME_RETRY                3
#define AISG_READ_FRAME_BUDGET            12
#define AISG_DISCOVERY_ADDR                0x02
#define AISG_PRIMARY_AISG_VERSION          0x02
#define AISG_PRIMARY_3GPP_RELEASE          0x06

typedef struct {
    uint8_t uniqueId[AISG_MAX_UID_LEN];
    size_t uniqueIdLen;
    uint8_t hdlcAddress;
    uint8_t deviceType;
    uint8_t vendorCode[2];
    bool hasUniqueId;
    bool hasAddress;
    bool hasDeviceType;
    bool hasVendorCode;
} XidAddressInfo;

static bool is_xid_control(uint8_t ctrl)
{
    return ctrl == AISG_CTRL_XID_PF || ctrl == AISG_CTRL_XID_BASE;
}

static bool is_ua_control(uint8_t ctrl)
{
    return ctrl == AISG_CTRL_UA_PF || ctrl == AISG_CTRL_UA_BASE;
}

static bool is_dm_control(uint8_t ctrl)
{
    return ctrl == AISG_CTRL_DM_PF || ctrl == AISG_CTRL_DM_BASE;
}

static bool is_i_frame(uint8_t ctrl)
{
    return (ctrl & 0x01) == 0;
}

static uint8_t make_i_control(uint8_t ns, uint8_t nr)
{
    return (uint8_t)(((ns & 0x07) << 1) | AISG_CTRL_I_PF |
                     ((nr & 0x07) << 5));
}

static uint8_t i_control_ns(uint8_t ctrl)
{
    return (uint8_t)((ctrl >> 1) & 0x07);
}

static uint8_t i_control_nr(uint8_t ctrl)
{
    return (uint8_t)((ctrl >> 5) & 0x07);
}

static void hexdump_line(const char *prefix, const uint8_t *data, size_t len)
{
    char line[512];
    size_t off = 0;
    size_t i;

    if (prefix == NULL || data == NULL) {
        return;
    }

    off += (size_t)snprintf(line + off, sizeof(line) - off, "%s len=%zu", prefix, len);
    for (i = 0; i < len && off + 4 < sizeof(line); i++) {
        off += (size_t)snprintf(line + off, sizeof(line) - off, " %02x", data[i]);
    }

    usys_log_trace("%s", line);
}

static void unique_id_to_hex(AisgDevice *device)
{
    size_t i;
    size_t off = 0;

    if (device == NULL) {
        return;
    }

    device->uniqueIdHex[0] = '\0';
    for (i = 0; i < device->uniqueIdLen && off + 2 < sizeof(device->uniqueIdHex); i++) {
        off += (size_t)snprintf(device->uniqueIdHex + off,
                                sizeof(device->uniqueIdHex) - off,
                                "%02X",
                                device->uniqueId[i]);
    }
}

static void unique_id_to_label(const uint8_t *uid,
                               size_t uidLen,
                               char *dst,
                               size_t dstSize)
{
    size_t i;
    bool printable = true;

    if (dst == NULL || dstSize == 0) {
        return;
    }

    dst[0] = '\0';
    if (uid == NULL || uidLen == 0) {
        snprintf(dst, dstSize, "%s", "RET");
        return;
    }

    for (i = 0; i < uidLen; i++) {
        if (!isprint(uid[i])) {
            printable = false;
            break;
        }
    }

    if (printable) {
        snprintf(dst, dstSize, "%.*s", (int)uidLen, (const char *)uid);
        return;
    }

    for (i = 0; i < uidLen && (i * 2 + 2) < dstSize; i++) {
        snprintf(dst + i * 2, dstSize - i * 2, "%02X", uid[i]);
    }
}

static bool send_payload_only(AisgBus *bus,
                              const uint8_t *payload,
                              size_t payloadLen)
{
    uint8_t frame[HDLC_MAX_FRAME];
    size_t frameLen;

    if (bus == NULL || bus->serial == NULL || payload == NULL) {
        return false;
    }

    if (!hdlc_encode(payload, payloadLen, frame, sizeof(frame), &frameLen)) {
        return false;
    }

    hexdump_line("aisg tx payload", payload, payloadLen);
    hexdump_line("aisg tx frame", frame, frameLen);

    if (!serial_flush_rx(bus->serial)) {
        return false;
    }

    if (!serial_write_all(bus->serial, frame, frameLen)) {
        return false;
    }

    return true;
}

static bool frame_matches_echo(const uint8_t *a,
                               size_t aLen,
                               const uint8_t *b,
                               size_t bLen)
{
    return a != NULL && b != NULL && aLen == bLen && memcmp(a, b, aLen) == 0;
}

static bool read_decoded_payload(AisgBus *bus,
                                 uint8_t *payload,
                                 size_t payloadSize,
                                 size_t *payloadLen,
                                 int timeoutMs)
{
    uint8_t frame[HDLC_MAX_FRAME];
    size_t frameLen;

    if (bus == NULL || payload == NULL || payloadLen == NULL) {
        return false;
    }

    if (!serial_read_frame(bus->serial,
                           frame,
                           sizeof(frame),
                           &frameLen,
                           timeoutMs)) {
        return false;
    }

    hexdump_line("aisg rx frame", frame, frameLen);

    if (!hdlc_decode(frame, frameLen, payload, payloadSize, payloadLen)) {
        usys_log_trace("aisg rx invalid HDLC/FCS frame ignored");
        return false;
    }

    hexdump_line("aisg rx payload", payload, *payloadLen);

    return true;
}

static bool transact_payload(AisgBus *bus,
                             const uint8_t *tx,
                             size_t txLen,
                             uint8_t *rx,
                             size_t rxSize,
                             size_t *rxLen,
                             int timeoutMs,
                             bool (*accept)(const uint8_t *candidate,
                                            size_t candidateLen,
                                            void *ctx),
                             void *acceptCtx)
{
    uint8_t candidate[HDLC_MAX_FRAME];
    size_t candidateLen;
    int i;

    if (!send_payload_only(bus, tx, txLen)) {
        return false;
    }

    for (i = 0; i < AISG_READ_FRAME_BUDGET; i++) {
        candidateLen = 0;
        if (!read_decoded_payload(bus,
                                  candidate,
                                  sizeof(candidate),
                                  &candidateLen,
                                  timeoutMs)) {
            return false;
        }

        if (frame_matches_echo(tx, txLen, candidate, candidateLen)) {
            usys_log_trace("aisg rx local echo ignored");
            continue;
        }

        if (accept != NULL && !accept(candidate, candidateLen, acceptCtx)) {
            usys_log_trace("aisg rx unexpected frame ignored addr=0x%02x ctrl=0x%02x",
                           candidateLen > 0 ? candidate[0] : 0,
                           candidateLen > 1 ? candidate[1] : 0);
            continue;
        }

        if (candidateLen > rxSize) {
            return false;
        }

        memcpy(rx, candidate, candidateLen);
        *rxLen = candidateLen;
        return true;
    }

    return false;
}

static bool append_byte(uint8_t *buf, size_t size, size_t *off, uint8_t value)
{
    if (buf == NULL || off == NULL || *off + 1 > size) {
        return false;
    }

    buf[(*off)++] = value;
    return true;
}

static bool append_bytes(uint8_t *buf,
                         size_t size,
                         size_t *off,
                         const uint8_t *src,
                         size_t srcLen)
{
    if (buf == NULL || off == NULL || src == NULL || *off + srcLen > size) {
        return false;
    }

    memcpy(&buf[*off], src, srcLen);
    *off += srcLen;
    return true;
}

static bool append_tlv(uint8_t *buf,
                       size_t size,
                       size_t *off,
                       uint8_t pi,
                       const uint8_t *pv,
                       size_t pvLen)
{
    if (pvLen > 255) {
        return false;
    }

    return append_byte(buf, size, off, pi) &&
           append_byte(buf, size, off, (uint8_t)pvLen) &&
           append_bytes(buf, size, off, pv, pvLen);
}

static bool xid_begin(uint8_t *buf, size_t size, size_t *off, uint8_t addr)
{
    if (buf == NULL || off == NULL) {
        return false;
    }

    *off = 0;
    return append_byte(buf, size, off, addr) &&
           append_byte(buf, size, off, AISG_CTRL_XID_PF) &&
           append_byte(buf, size, off, AISG_XID_FI_STD) &&
           append_byte(buf, size, off, AISG_XID_GI_ADDR) &&
           append_byte(buf, size, off, 0x00);
}

static bool xid_finish(uint8_t *buf, size_t off)
{
    if (buf == NULL || off < 5 || off - 5 > 255) {
        return false;
    }

    buf[4] = (uint8_t)(off - 5);
    return true;
}

static bool build_xid_scan(uint8_t *buf, size_t size, size_t *len)
{
    size_t off;
    uint8_t uid = 0x00;
    uint8_t mask = 0x00;

    if (!xid_begin(buf, size, &off, AISG_ADDR_BROADCAST)) {
        return false;
    }

    /* Single-RET deployment: wildcard one-octet tree-scan branch. */
    if (!append_tlv(buf, size, &off, AISG_XID_PI_UNIQUE_ID, &uid, 1) ||
        !append_tlv(buf, size, &off, AISG_XID_PI_BIT_MASK, &mask, 1)) {
        return false;
    }

    if (!xid_finish(buf, off)) {
        return false;
    }

    *len = off;
    return true;
}

static bool build_xid_assign(uint8_t *buf,
                             size_t size,
                             size_t *len,
                             const XidAddressInfo *info,
                             uint8_t assignedAddress)
{
    size_t off;
    uint8_t rel = AISG_PRIMARY_3GPP_RELEASE;
    uint8_t vendorWild[2] = {0xFF, 0xFF};
    uint8_t typeWild = 0xFF;
    const uint8_t *vendor = vendorWild;
    uint8_t type = typeWild;

    if (info == NULL || !info->hasUniqueId || info->uniqueIdLen == 0) {
        return false;
    }

    if (!xid_begin(buf, size, &off, AISG_ADDR_BROADCAST)) {
        return false;
    }

    if (info->hasVendorCode) {
        vendor = info->vendorCode;
    } else if (info->uniqueIdLen >= 2) {
        vendor = info->uniqueId;
    }

    if (info->hasDeviceType) {
        type = info->deviceType;
    }

    if (!append_tlv(buf, size, &off,
                    AISG_XID_PI_UNIQUE_ID,
                    info->uniqueId,
                    info->uniqueIdLen) ||
        !append_tlv(buf, size, &off,
                    AISG_XID_PI_HDLC_ADDR,
                    &assignedAddress,
                    1) ||
        !append_tlv(buf, size, &off,
                    AISG_XID_PI_VENDOR_CODE,
                    vendor,
                    2) ||
        !append_tlv(buf, size, &off,
                    AISG_XID_PI_DEVICE_TYPE,
                    &type,
                    1) ||
        !append_tlv(buf, size, &off,
                    AISG_XID_PI_3GPP_RELEASE,
                    &rel,
                    1)) {
        return false;
    }

    if (!xid_finish(buf, off)) {
        return false;
    }

    *len = off;
    return true;
}

static bool build_xid_aisg_version(uint8_t *buf,
                                   size_t size,
                                   size_t *len,
                                   uint8_t address)
{
    size_t off;
    uint8_t version = AISG_PRIMARY_AISG_VERSION;

    if (!xid_begin(buf, size, &off, address)) {
        return false;
    }

    if (!append_tlv(buf, size, &off,
                    AISG_XID_PI_AISG_VERSION,
                    &version,
                    1)) {
        return false;
    }

    if (!xid_finish(buf, off)) {
        return false;
    }

    *len = off;
    return true;
}

static bool parse_xid_address_info(const uint8_t *payload,
                                   size_t len,
                                   XidAddressInfo *info)
{
    size_t end;
    size_t pos;
    uint8_t pi;
    uint8_t pl;

    if (payload == NULL || info == NULL || len < 5) {
        return false;
    }

    if (!is_xid_control(payload[1]) ||
        payload[2] != AISG_XID_FI_STD ||
        payload[3] != AISG_XID_GI_ADDR) {
        return false;
    }

    memset(info, 0, sizeof(*info));

    end = 5 + payload[4];
    if (end > len) {
        return false;
    }

    pos = 5;
    while (pos + 2 <= end) {
        pi = payload[pos++];
        pl = payload[pos++];
        if (pos + pl > end) {
            return false;
        }

        switch (pi) {
        case AISG_XID_PI_UNIQUE_ID:
            if (pl > 0 && pl <= AISG_MAX_UID_LEN) {
                memcpy(info->uniqueId, &payload[pos], pl);
                info->uniqueIdLen = pl;
                info->hasUniqueId = true;
            }
            break;
        case AISG_XID_PI_HDLC_ADDR:
            if (pl == 1) {
                info->hdlcAddress = payload[pos];
                info->hasAddress = true;
            }
            break;
        case AISG_XID_PI_DEVICE_TYPE:
            if (pl == 1) {
                info->deviceType = payload[pos];
                info->hasDeviceType = true;
            }
            break;
        case AISG_XID_PI_VENDOR_CODE:
            if (pl == 2) {
                info->vendorCode[0] = payload[pos];
                info->vendorCode[1] = payload[pos + 1];
                info->hasVendorCode = true;
            }
            break;
        default:
            break;
        }

        pos += pl;
    }

    return info->hasUniqueId;
}

static bool accept_xid_any(const uint8_t *candidate, size_t candidateLen, void *ctx)
{
    XidAddressInfo info;
    (void)ctx;

    return parse_xid_address_info(candidate, candidateLen, &info);
}

static bool accept_xid_addr(const uint8_t *candidate, size_t candidateLen, void *ctx)
{
    uint8_t expected = *(uint8_t *)ctx;

    if (candidate == NULL || candidateLen < 2) {
        return false;
    }

    if (candidate[0] != expected) {
        return false;
    }

    return accept_xid_any(candidate, candidateLen, NULL);
}

static bool accept_ua_addr(const uint8_t *candidate, size_t candidateLen, void *ctx)
{
    uint8_t expected = *(uint8_t *)ctx;

    if (candidate == NULL || candidateLen < 2) {
        return false;
    }

    if (candidate[0] != expected) {
        return false;
    }

    if (is_dm_control(candidate[1])) {
        usys_log_trace("aisg secondary returned DM while establishing link");
        return false;
    }

    return is_ua_control(candidate[1]);
}

static bool accept_retap_i_frame(const uint8_t *candidate,
                                 size_t candidateLen,
                                 void *ctx)
{
    AisgBus *bus = ctx;
    uint8_t remoteNs;
    uint8_t remoteNr;

    if (bus == NULL || candidate == NULL || candidateLen < 5) {
        return false;
    }

    if (candidate[0] != bus->deviceAddress || !is_i_frame(candidate[1])) {
        return false;
    }

    remoteNs = i_control_ns(candidate[1]);
    remoteNr = i_control_nr(candidate[1]);

    if (remoteNr != ((bus->txSeq + 1) & 0x07)) {
        usys_log_trace("aisg I-frame ack mismatch expected=%u got=%u",
                       (bus->txSeq + 1) & 0x07,
                       remoteNr);
        return false;
    }

    if (remoteNs != bus->rxSeq) {
        usys_log_trace("aisg I-frame rx sequence mismatch expected=%u got=%u",
                       bus->rxSeq,
                       remoteNs);
        return false;
    }

    return true;
}

static bool send_snrm(AisgBus *bus, uint8_t address)
{
    uint8_t req[2];
    uint8_t resp[HDLC_MAX_FRAME];
    size_t respLen;
    int retry;

    req[0] = address;
    req[1] = AISG_CTRL_SNRM_PF;

    for (retry = 0; retry < AISG_TX_FRAME_RETRY; retry++) {
        if (transact_payload(bus,
                             req,
                             sizeof(req),
                             resp,
                             sizeof(resp),
                             &respLen,
                             AISG_DEFAULT_TIMEOUT_MS,
                             accept_ua_addr,
                             &address)) {
            bus->deviceAddress = address;
            bus->txSeq = 0;
            bus->rxSeq = 0;
            bus->linked = true;
            return true;
        }
    }

    return false;
}

static bool negotiate_aisg_version(AisgBus *bus, uint8_t address)
{
    uint8_t req[HDLC_MAX_FRAME];
    uint8_t resp[HDLC_MAX_FRAME];
    size_t reqLen;
    size_t respLen;

    if (!build_xid_aisg_version(req, sizeof(req), &reqLen, address)) {
        return false;
    }

    if (!transact_payload(bus,
                          req,
                          reqLen,
                          resp,
                          sizeof(resp),
                          &respLen,
                          AISG_DEFAULT_TIMEOUT_MS,
                          accept_xid_addr,
                          &address)) {
        /* Some strict 3GPP RETs do not implement AISG PI=20. Keep going. */
        usys_log_trace("aisg version negotiation did not respond; continuing with 3GPP RETAP");
    }

    return true;
}

static bool scan_once(AisgBus *bus, XidAddressInfo *found)
{
    uint8_t req[HDLC_MAX_FRAME];
    uint8_t resp[HDLC_MAX_FRAME];
    size_t reqLen;
    size_t respLen;

    if (!build_xid_scan(req, sizeof(req), &reqLen)) {
        return false;
    }

    if (!transact_payload(bus,
                          req,
                          reqLen,
                          resp,
                          sizeof(resp),
                          &respLen,
                          AISG_DEFAULT_TIMEOUT_MS,
                          accept_xid_any,
                          NULL)) {
        return false;
    }

    return parse_xid_address_info(resp, respLen, found);
}

static bool assign_address(AisgBus *bus,
                           const XidAddressInfo *found,
                           uint8_t assignedAddress)
{
    uint8_t req[HDLC_MAX_FRAME];
    uint8_t resp[HDLC_MAX_FRAME];
    size_t reqLen;
    size_t respLen;

    if (!build_xid_assign(req, sizeof(req), &reqLen, found, assignedAddress)) {
        return false;
    }

    if (!transact_payload(bus,
                          req,
                          reqLen,
                          resp,
                          sizeof(resp),
                          &respLen,
                          AISG_DEFAULT_TIMEOUT_MS,
                          accept_xid_addr,
                          &assignedAddress)) {
        return false;
    }

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
    bus->txSeq         = 0;
    bus->rxSeq         = 0;
    bus->linked        = false;
}

bool aisg_v2_scan(AisgBus *bus, AisgDevice *device)
{
    static const int bauds[] = {9600, 38400, 115200};
    XidAddressInfo found;
    uint8_t assignedAddress = AISG_DISCOVERY_ADDR;
    size_t i;
    int originalBaud;

    if (bus == NULL || bus->serial == NULL || device == NULL) {
        return false;
    }

    memset(device, 0, sizeof(AisgDevice));
    bus->linked = false;
    bus->deviceAddress = AISG_ADDR_DEFAULT;
    bus->txSeq = 0;
    bus->rxSeq = 0;

    originalBaud = bus->serial->baud;

    for (i = 0; i < sizeof(bauds) / sizeof(bauds[0]); i++) {
        if (bauds[i] != bus->serial->baud && !serial_set_baud(bus->serial, bauds[i])) {
            continue;
        }

        usleep(AISG_SCAN_HOLD_MS * 1000);
        usys_log_trace("aisg scan try baud=%d", bauds[i]);

        memset(&found, 0, sizeof(found));
        if (!scan_once(bus, &found)) {
            continue;
        }

        if (found.hasAddress && found.hdlcAddress >= 1 && found.hdlcAddress <= 254) {
            assignedAddress = found.hdlcAddress;
        }

        if (assignedAddress == AISG_ADDR_DEFAULT || assignedAddress == AISG_ADDR_BROADCAST) {
            assignedAddress = AISG_DISCOVERY_ADDR;
        }

        if (!assign_address(bus, &found, assignedAddress)) {
            if (found.hasAddress && found.hdlcAddress >= 1 && found.hdlcAddress <= 254 &&
                send_snrm(bus, found.hdlcAddress)) {
                assignedAddress = found.hdlcAddress;
            } else {
                assignedAddress = AISG_DISCOVERY_ADDR;
                if (!assign_address(bus, &found, assignedAddress)) {
                    continue;
                }
            }
        }

        negotiate_aisg_version(bus, assignedAddress);

        if (!bus->linked && !send_snrm(bus, assignedAddress)) {
            continue;
        }

        device->present     = true;
        device->linked      = true;
        device->address     = assignedAddress;
        device->deviceType  = found.hasDeviceType ? found.deviceType : AISG_DEVICE_TYPE_RET_SINGLE;
        device->baud        = bus->serial->baud;
        device->uniqueIdLen = found.uniqueIdLen;
        memcpy(device->uniqueId, found.uniqueId, found.uniqueIdLen);
        unique_id_to_hex(device);
        unique_id_to_label(device->uniqueId,
                           device->uniqueIdLen,
                           device->model,
                           sizeof(device->model));

        bus->deviceAddress = assignedAddress;
        bus->linked = true;

        return true;
    }

    serial_set_baud(bus->serial, originalBaud);
    return false;
}

bool aisg_v2_send_retap(AisgBus *bus,
                        RetapRequest *request,
                        RetapResponse *response)
{
    uint8_t retap[RETAP_MAX_PAYLOAD + 3];
    uint8_t tx[RETAP_MAX_PAYLOAD + 8];
    uint8_t rx[HDLC_MAX_FRAME];
    size_t retapLen;
    size_t txLen;
    size_t rxLen;
    size_t rxRetapLen;
    int retry;

    if (bus == NULL || request == NULL || response == NULL ||
        bus->deviceAddress == AISG_ADDR_DEFAULT || !bus->linked) {
        return false;
    }

    if (!retap_encode_request(request, retap, sizeof(retap), &retapLen)) {
        return false;
    }

    for (retry = 0; retry < AISG_TX_FRAME_RETRY; retry++) {
        txLen = 0;
        tx[txLen++] = bus->deviceAddress;
        tx[txLen++] = make_i_control(bus->txSeq, bus->rxSeq);

        memcpy(&tx[txLen], retap, retapLen);
        txLen += retapLen;

        if (!transact_payload(bus,
                              tx,
                              txLen,
                              rx,
                              sizeof(rx),
                              &rxLen,
                              AISG_LONG_TIMEOUT_MS,
                              accept_retap_i_frame,
                              bus)) {
            send_snrm(bus, bus->deviceAddress);
            continue;
        }

        if (rxLen < 5) {
            return false;
        }

        rxRetapLen = rxLen - 2;
        if (!retap_decode_response(&rx[2], rxRetapLen, response)) {
            return false;
        }

        bus->txSeq = (uint8_t)((bus->txSeq + 1) & 0x07);
        bus->rxSeq = (uint8_t)((i_control_ns(rx[1]) + 1) & 0x07);

        return true;
    }

    return false;
}

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "aisg_v2.h"
#include "hdlc.h"
#include "usys_log.h"

#define AISG_TRACE_BYTES_PER_LINE          16

static const char *ctrl_name(uint8_t ctrl)
{
    if (ctrl == AISG_CTRL_XID) {
        return "XID";
    }

    if ((ctrl & 0x01) == 0) {
        return "I";
    }

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

static void log_decoded_payload(const char *label,
                                const uint8_t *payload,
                                size_t payloadLen)
{
    if (label == NULL || payload == NULL) {
        return;
    }

    if (payloadLen < 2) {
        usys_log_debug("aisg: %s decoded payload too short len=%zu",
                       label,
                       payloadLen);
        return;
    }

    usys_log_debug("aisg: %s decoded addr=0x%02X ctrl=0x%02X(%s) info_len=%zu",
                   label,
                   payload[0],
                   payload[1],
                   ctrl_name(payload[1]),
                   payloadLen - 2);
}

static bool same_payload(const uint8_t *a,
                         size_t aLen,
                         const uint8_t *b,
                         size_t bLen)
{
    if (a == NULL || b == NULL) {
        return false;
    }

    if (aLen != bLen) {
        return false;
    }

    return memcmp(a, b, aLen) == 0;
}

static bool read_decoded_response(AisgBus *bus,
                                  const uint8_t *txPayload,
                                  size_t txPayloadLen,
                                  uint8_t *response,
                                  size_t responseSize,
                                  size_t *responseLen)
{
    uint8_t rxFrame[HDLC_MAX_FRAME];
    uint8_t decoded[HDLC_MAX_FRAME];
    size_t rxFrameLen;
    size_t decodedLen;
    int attempt;

    if (bus == NULL ||
        txPayload == NULL ||
        response == NULL ||
        responseLen == NULL) {
        return false;
    }

    for (attempt = 0; attempt < AISG_MAX_RX_ATTEMPTS; attempt++) {
        memset(rxFrame, 0, sizeof(rxFrame));
        rxFrameLen = 0;

        if (!serial_read_frame(bus->serial,
                               rxFrame,
                               sizeof(rxFrame),
                               &rxFrameLen,
                               AISG_DEFAULT_TIMEOUT_MS)) {
            usys_log_debug("aisg: RX timeout attempt=%d", attempt + 1);
            return false;
        }

        log_hex_bytes("RX hdlc", rxFrame, rxFrameLen);

        memset(decoded, 0, sizeof(decoded));
        decodedLen = 0;

        if (!hdlc_decode(rxFrame,
                         rxFrameLen,
                         decoded,
                         sizeof(decoded),
                         &decodedLen)) {
            usys_log_debug("aisg: RX hdlc decode failed attempt=%d",
                           attempt + 1);
            continue;
        }

        log_hex_bytes("RX payload", decoded, decodedLen);
        log_decoded_payload("RX", decoded, decodedLen);

        if (same_payload(txPayload, txPayloadLen, decoded, decodedLen)) {
            usys_log_debug("aisg: RX rejected local echo attempt=%d",
                           attempt + 1);
            continue;
        }

        if (decodedLen > responseSize) {
            usys_log_debug("aisg: RX decoded payload too large len=%zu size=%zu",
                           decodedLen,
                           responseSize);
            return false;
        }

        memcpy(response, decoded, decodedLen);
        *responseLen = decodedLen;

        return true;
    }

    usys_log_debug("aisg: RX failed after echo/invalid frames attempts=%d",
                   AISG_MAX_RX_ATTEMPTS);

    return false;
}

static bool send_payload(AisgBus *bus,
                         const uint8_t *payload,
                         size_t payloadLen,
                         uint8_t *response,
                         size_t responseSize,
                         size_t *responseLen)
{
    uint8_t txFrame[HDLC_MAX_FRAME];
    size_t txFrameLen;

    if (bus == NULL ||
        bus->serial == NULL ||
        payload == NULL ||
        response == NULL ||
        responseLen == NULL) {
        return false;
    }

    log_hex_bytes("TX payload", payload, payloadLen);
    log_decoded_payload("TX", payload, payloadLen);

    memset(txFrame, 0, sizeof(txFrame));
    txFrameLen = 0;

    if (!hdlc_encode(payload,
                     payloadLen,
                     txFrame,
                     sizeof(txFrame),
                     &txFrameLen)) {
        usys_log_debug("aisg: TX hdlc encode failed");
        return false;
    }

    log_hex_bytes("TX hdlc", txFrame, txFrameLen);

    if (!serial_write_all(bus->serial, txFrame, txFrameLen)) {
        usys_log_debug("aisg: TX serial write failed");
        return false;
    }

    return read_decoded_response(bus,
                                 payload,
                                 payloadLen,
                                 response,
                                 responseSize,
                                 responseLen);
}

void aisg_v2_bus_init(AisgBus *bus, SerialPort *serial)
{
    if (bus == NULL) {
        return;
    }

    memset(bus, 0, sizeof(AisgBus));

    bus->serial        = serial;
    bus->deviceAddress = AISG_ADDR_DEFAULT;

    usys_log_debug("aisg: init scope=%s supported_device_type=0x%02X",
                   AISG_SCOPE_NAME,
                   AISG_SUPPORTED_DEVICE_TYPE);
}

bool aisg_v2_scan(AisgBus *bus, AisgDevice *device)
{
    uint8_t req[5];
    uint8_t resp[HDLC_MAX_FRAME];
    size_t respLen;

    if (bus == NULL || device == NULL) {
        return false;
    }

    memset(device, 0, sizeof(AisgDevice));

    /*
     * Phase 0/1 keeps the existing scan shape only so we can observe
     * real hardware behavior and reject local echo.
     *
     * Phase 3 replaces this with standards-correct TS 25.462 XID scan
     * and address assignment.
     */
    req[0] = AISG_ADDR_BROADCAST;
    req[1] = AISG_CTRL_XID;
    req[2] = 0x81;
    req[3] = 0xF0;
    req[4] = 0x00;

    if (!send_payload(bus, req, sizeof(req), resp, sizeof(resp), &respLen)) {
        usys_log_debug("aisg: scan failed: no valid non-echo response");
        return false;
    }

    if (respLen < 2) {
        usys_log_debug("aisg: scan failed: response too short len=%zu",
                       respLen);
        return false;
    }

    /*
     * Temporary Phase 0/1 behavior:
     * mark present only after a valid decoded non-echo HDLC response.
     *
     * Do not infer true device identity yet. Real device type/model parsing
     * lands with standards-correct XID handling in Phase 3.
     */
    device->present    = true;
    device->unsupported = false;
    device->address    = AISG_ADDR_DEFAULT;
    device->deviceType = 0;
    snprintf(device->model, sizeof(device->model), "%s", "unknown");

    bus->deviceAddress = device->address;

    usys_log_debug("aisg: scan saw valid non-echo response addr=0x%02X "
                   "ctrl=0x%02X len=%zu",
                   resp[0],
                   resp[1],
                   respLen);

    return true;
}

bool aisg_v2_send_retap(AisgBus *bus,
                        RetapRequest *request,
                        RetapResponse *response)
{
    uint8_t retap[RETAP_MAX_PAYLOAD + 1];
    uint8_t tx[RETAP_MAX_PAYLOAD + 8];
    uint8_t rx[HDLC_MAX_FRAME];
    uint8_t rxRetap[RETAP_MAX_PAYLOAD + 1];
    size_t retapLen;
    size_t txLen;
    size_t rxLen;
    size_t rxRetapLen;

    if (bus == NULL || request == NULL || response == NULL) {
        return false;
    }

    if (!retap_encode_request(request, retap, sizeof(retap), &retapLen)) {
        usys_log_debug("aisg: failed to encode RETAP request");
        return false;
    }

    txLen = 0;
    tx[txLen++] = bus->deviceAddress;
    tx[txLen++] = AISG_CTRL_I_FRAME;

    memcpy(&tx[txLen], retap, retapLen);
    txLen += retapLen;

    if (!send_payload(bus, tx, txLen, rx, sizeof(rx), &rxLen)) {
        usys_log_debug("aisg: RETAP transport failed");
        return false;
    }

    if (rxLen < 3) {
        usys_log_debug("aisg: RETAP response too short len=%zu", rxLen);
        return false;
    }

    rxRetapLen = rxLen - 2;
    if (rxRetapLen > sizeof(rxRetap)) {
        usys_log_debug("aisg: RETAP response payload too large len=%zu",
                       rxRetapLen);
        return false;
    }

    memcpy(rxRetap, &rx[2], rxRetapLen);

    return retap_decode_response(rxRetap, rxRetapLen, response);
}

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

static bool send_payload(AisgBus *bus,
                         const uint8_t *payload,
                         size_t payloadLen,
                         uint8_t *response,
                         size_t responseSize,
                         size_t *responseLen)
{
    uint8_t tx[HDLC_MAX_FRAME];
    uint8_t rx[HDLC_MAX_FRAME];
    size_t txLen;
    size_t rxLen;

    if (!hdlc_encode(payload, payloadLen, tx, sizeof(tx), &txLen)) {
        return false;
    }

    if (!serial_write_all(bus->serial, tx, txLen)) {
        return false;
    }

    if (!serial_read_frame(bus->serial,
                           rx,
                           sizeof(rx),
                           &rxLen,
                           AISG_DEFAULT_TIMEOUT_MS)) {
        return false;
    }

    return hdlc_decode(rx, rxLen, response, responseSize, responseLen);
}

void aisg_v2_bus_init(AisgBus *bus, SerialPort *serial)
{
    if (bus == NULL) {
        return;
    }

    memset(bus, 0, sizeof(AisgBus));

    bus->serial        = serial;
    bus->deviceAddress = AISG_ADDR_DEFAULT;
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

    req[0] = AISG_ADDR_BROADCAST;
    req[1] = AISG_CTRL_XID;
    req[2] = 0x81;
    req[3] = 0xF0;
    req[4] = 0x00;

    if (!send_payload(bus, req, sizeof(req), resp, sizeof(resp), &respLen)) {
        return false;
    }

    if (respLen < 2) {
        return false;
    }

    device->present = true;
    device->address = AISG_ADDR_DEFAULT;
    snprintf(device->model, sizeof(device->model), "%s", "RET1T1");

    bus->deviceAddress = device->address;

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
        return false;
    }

    txLen = 0;
    tx[txLen++] = bus->deviceAddress;
    tx[txLen++] = AISG_CTRL_I_FRAME;

    memcpy(&tx[txLen], retap, retapLen);
    txLen += retapLen;

    if (!send_payload(bus, tx, txLen, rx, sizeof(rx), &rxLen)) {
        return false;
    }

    if (rxLen < 3) {
        return false;
    }

    rxRetapLen = rxLen - 2;
    memcpy(rxRetap, &rx[2], rxRetapLen);

    return retap_decode_response(rxRetap, rxRetapLen, response);
}

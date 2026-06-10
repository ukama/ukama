/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "hdlc.h"

static uint16_t update_fcs_bit(uint16_t fcs)
{
    if ((fcs & 1) != 0) {
        return (uint16_t)((fcs >> 1) ^ 0x8408);
    }

    return (uint16_t)(fcs >> 1);
}

uint16_t hdlc_fcs16(const uint8_t *data, size_t len)
{
    uint16_t fcs = 0xFFFF;
    size_t i;
    int bit;

    for (i = 0; i < len; i++) {
        fcs ^= data[i];

        for (bit = 0; bit < 8; bit++) {
            fcs = update_fcs_bit(fcs);
        }
    }

    return (uint16_t)~fcs;
}

static bool append_byte(uint8_t byte,
                        uint8_t *frame,
                        size_t frameSize,
                        size_t *off)
{
    if (*off + 1 > frameSize) {
        return false;
    }

    frame[(*off)++] = byte;

    return true;
}

static bool append_escaped(uint8_t byte,
                           uint8_t *frame,
                           size_t frameSize,
                           size_t *off)
{
    if (byte == HDLC_FLAG || byte == HDLC_ESCAPE) {
        if (!append_byte(HDLC_ESCAPE, frame, frameSize, off)) {
            return false;
        }

        return append_byte(byte ^ HDLC_ESCAPE_XOR, frame, frameSize, off);
    }

    return append_byte(byte, frame, frameSize, off);
}

bool hdlc_encode(const uint8_t *payload,
                 size_t payloadLen,
                 uint8_t *frame,
                 size_t frameSize,
                 size_t *frameLen)
{
    uint16_t fcs;
    size_t off = 0;
    size_t i;

    if (payload == NULL || frame == NULL || frameLen == NULL) {
        return false;
    }

    if (!append_byte(HDLC_FLAG, frame, frameSize, &off)) {
        return false;
    }

    for (i = 0; i < payloadLen; i++) {
        if (!append_escaped(payload[i], frame, frameSize, &off)) {
            return false;
        }
    }

    fcs = hdlc_fcs16(payload, payloadLen);

    if (!append_escaped((uint8_t)(fcs & 0xFF), frame, frameSize, &off)) {
        return false;
    }

    if (!append_escaped((uint8_t)((fcs >> 8) & 0xFF),
                        frame,
                        frameSize,
                        &off)) {
        return false;
    }

    if (!append_byte(HDLC_FLAG, frame, frameSize, &off)) {
        return false;
    }

    *frameLen = off;

    return true;
}

static bool decode_byte(uint8_t byte,
                        uint8_t *decoded,
                        size_t decodedSize,
                        size_t *decodedLen,
                        bool *escaped)
{
    if (*escaped) {
        byte = byte ^ HDLC_ESCAPE_XOR;
        *escaped = false;
    } else if (byte == HDLC_ESCAPE) {
        *escaped = true;
        return true;
    }

    if (*decodedLen >= decodedSize) {
        return false;
    }

    decoded[(*decodedLen)++] = byte;

    return true;
}

bool hdlc_decode(const uint8_t *frame,
                 size_t frameLen,
                 uint8_t *payload,
                 size_t payloadSize,
                 size_t *payloadLen)
{
    uint8_t decoded[HDLC_MAX_FRAME];
    size_t decodedLen = 0;
    size_t i;
    bool started = false;
    bool escaped = false;
    uint16_t rxFcs;
    uint16_t calcFcs;

    if (frame == NULL || payload == NULL || payloadLen == NULL) {
        return false;
    }

    for (i = 0; i < frameLen; i++) {
        if (frame[i] == HDLC_FLAG) {
            if (!started) {
                started = true;
                continue;
            }

            break;
        }

        if (!started) {
            continue;
        }

        if (!decode_byte(frame[i],
                         decoded,
                         sizeof(decoded),
                         &decodedLen,
                         &escaped)) {
            return false;
        }
    }

    if (decodedLen < 3) {
        return false;
    }

    rxFcs = (uint16_t)decoded[decodedLen - 2];
    rxFcs |= ((uint16_t)decoded[decodedLen - 1] << 8);

    calcFcs = hdlc_fcs16(decoded, decodedLen - 2);
    if (rxFcs != calcFcs) {
        return false;
    }

    if (decodedLen - 2 > payloadSize) {
        return false;
    }

    memcpy(payload, decoded, decodedLen - 2);
    *payloadLen = decodedLen - 2;

    return true;
}

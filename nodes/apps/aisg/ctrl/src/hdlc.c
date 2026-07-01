/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "hdlc.h"

#define HDLC_CTRL_PF                       0x10
#define HDLC_CTRL_U_MASK_NO_PF             0xEF
#define HDLC_CTRL_S_MASK_NO_NR_PF          0x0F

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

    if (data == NULL && len != 0) {
        return 0;
    }

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
    if (frame == NULL || off == NULL) {
        return false;
    }

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

    if (payloadLen < 2) {
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

bool hdlc_encode_addr_info(uint8_t address,
                           uint8_t control,
                           const uint8_t *info,
                           size_t infoLen,
                           uint8_t *frame,
                           size_t frameSize,
                           size_t *frameLen)
{
    uint8_t payload[HDLC_MAX_INFO + 2];

    if (frame == NULL || frameLen == NULL) {
        return false;
    }

    if (info == NULL && infoLen != 0) {
        return false;
    }

    if (infoLen > HDLC_MAX_INFO) {
        return false;
    }

    payload[0] = address;
    payload[1] = control;

    if (infoLen > 0) {
        memcpy(&payload[2], info, infoLen);
    }

    return hdlc_encode(payload, infoLen + 2, frame, frameSize, frameLen);
}

bool hdlc_encode_frame(const HdlcFrame *decoded,
                       uint8_t *frame,
                       size_t frameSize,
                       size_t *frameLen)
{
    if (decoded == NULL) {
        return false;
    }

    return hdlc_encode_addr_info(decoded->address,
                                 decoded->control,
                                 decoded->info,
                                 decoded->infoLen,
                                 frame,
                                 frameSize,
                                 frameLen);
}

static bool decode_byte(uint8_t byte,
                        uint8_t *decoded,
                        size_t decodedSize,
                        size_t *decodedLen,
                        bool *escaped)
{
    if (decoded == NULL || decodedLen == NULL || escaped == NULL) {
        return false;
    }

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
    bool ended = false;
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

            ended = true;
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

    if (!started || !ended || escaped) {
        return false;
    }

    if (decodedLen < 4) { /* addr + ctrl + fcs-low + fcs-high */
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

bool hdlc_frame_from_payload(const uint8_t *payload,
                             size_t payloadLen,
                             HdlcFrame *decoded)
{
    if (payload == NULL || decoded == NULL) {
        return false;
    }

    if (payloadLen < 2) {
        return false;
    }

    if (payloadLen - 2 > HDLC_MAX_INFO) {
        return false;
    }

    memset(decoded, 0, sizeof(*decoded));

    decoded->address = payload[0];
    decoded->control = payload[1];
    decoded->infoLen = payloadLen - 2;

    if (decoded->infoLen > 0) {
        memcpy(decoded->info, &payload[2], decoded->infoLen);
    }

    return true;
}

bool hdlc_payload_from_frame(const HdlcFrame *decoded,
                             uint8_t *payload,
                             size_t payloadSize,
                             size_t *payloadLen)
{
    if (decoded == NULL || payload == NULL || payloadLen == NULL) {
        return false;
    }

    if (decoded->infoLen > HDLC_MAX_INFO) {
        return false;
    }

    if (decoded->infoLen + 2 > payloadSize) {
        return false;
    }

    payload[0] = decoded->address;
    payload[1] = decoded->control;

    if (decoded->infoLen > 0) {
        memcpy(&payload[2], decoded->info, decoded->infoLen);
    }

    *payloadLen = decoded->infoLen + 2;

    return true;
}

bool hdlc_decode_frame(const uint8_t *frame,
                       size_t frameLen,
                       HdlcFrame *decoded)
{
    uint8_t payload[HDLC_MAX_INFO + 2];
    size_t payloadLen;

    if (decoded == NULL) {
        return false;
    }

    memset(payload, 0, sizeof(payload));
    payloadLen = 0;

    if (!hdlc_decode(frame,
                     frameLen,
                     payload,
                     sizeof(payload),
                     &payloadLen)) {
        return false;
    }

    return hdlc_frame_from_payload(payload, payloadLen, decoded);
}

uint8_t hdlc_i_ctrl(uint8_t ns, uint8_t nr, bool poll)
{
    uint8_t ctrl;

    ctrl = (uint8_t)(((ns & 0x07) << 1) | ((nr & 0x07) << 5));
    if (poll) {
        ctrl |= HDLC_CTRL_PF;
    }

    return ctrl;
}

uint8_t hdlc_rr_ctrl(uint8_t nr, bool poll)
{
    uint8_t ctrl;

    ctrl = (uint8_t)(0x01 | ((nr & 0x07) << 5));
    if (poll) {
        ctrl |= HDLC_CTRL_PF;
    }

    return ctrl;
}

uint8_t hdlc_rnr_ctrl(uint8_t nr, bool poll)
{
    uint8_t ctrl;

    ctrl = (uint8_t)(0x05 | ((nr & 0x07) << 5));
    if (poll) {
        ctrl |= HDLC_CTRL_PF;
    }

    return ctrl;
}

uint8_t hdlc_snrm_ctrl(bool poll)
{
    return (uint8_t)(0x83 | (poll ? HDLC_CTRL_PF : 0x00));
}

uint8_t hdlc_disc_ctrl(bool poll)
{
    return (uint8_t)(0x43 | (poll ? HDLC_CTRL_PF : 0x00));
}

uint8_t hdlc_ua_ctrl(bool final)
{
    return (uint8_t)(0x63 | (final ? HDLC_CTRL_PF : 0x00));
}

uint8_t hdlc_dm_ctrl(bool final)
{
    return (uint8_t)(0x0F | (final ? HDLC_CTRL_PF : 0x00));
}

uint8_t hdlc_frmr_ctrl(bool final)
{
    return (uint8_t)(0x87 | (final ? HDLC_CTRL_PF : 0x00));
}

uint8_t hdlc_xid_ctrl(bool poll)
{
    return (uint8_t)(0xAF | (poll ? HDLC_CTRL_PF : 0x00));
}

bool hdlc_is_i_frame(uint8_t ctrl)
{
    return (ctrl & 0x01) == 0;
}

bool hdlc_is_rr(uint8_t ctrl)
{
    return !hdlc_is_i_frame(ctrl) && ((ctrl & HDLC_CTRL_S_MASK_NO_NR_PF) == 0x01);
}

bool hdlc_is_rnr(uint8_t ctrl)
{
    return !hdlc_is_i_frame(ctrl) && ((ctrl & HDLC_CTRL_S_MASK_NO_NR_PF) == 0x05);
}

bool hdlc_is_snrm(uint8_t ctrl)
{
    return (ctrl & HDLC_CTRL_U_MASK_NO_PF) == 0x83;
}

bool hdlc_is_disc(uint8_t ctrl)
{
    return (ctrl & HDLC_CTRL_U_MASK_NO_PF) == 0x43;
}

bool hdlc_is_ua(uint8_t ctrl)
{
    return (ctrl & HDLC_CTRL_U_MASK_NO_PF) == 0x63;
}

bool hdlc_is_dm(uint8_t ctrl)
{
    return (ctrl & HDLC_CTRL_U_MASK_NO_PF) == 0x0F;
}

bool hdlc_is_frmr(uint8_t ctrl)
{
    return (ctrl & HDLC_CTRL_U_MASK_NO_PF) == 0x87;
}

bool hdlc_is_xid(uint8_t ctrl)
{
    return (ctrl & HDLC_CTRL_U_MASK_NO_PF) == 0xAF;
}

bool hdlc_poll_final(uint8_t ctrl)
{
    return (ctrl & HDLC_CTRL_PF) != 0;
}

uint8_t hdlc_ns(uint8_t ctrl)
{
    if (!hdlc_is_i_frame(ctrl)) {
        return 0;
    }

    return (uint8_t)((ctrl >> 1) & 0x07);
}

uint8_t hdlc_nr(uint8_t ctrl)
{
    return (uint8_t)((ctrl >> 5) & 0x07);
}

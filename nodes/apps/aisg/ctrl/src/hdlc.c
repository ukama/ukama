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

        return append_byte((uint8_t)(byte ^ HDLC_ESCAPE_XOR),
                           frame,
                           frameSize,
                           off);
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

    *frameLen = 0;

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
    if (decoded == NULL || decodedLen == NULL || escaped == NULL) {
        return false;
    }

    if (*escaped) {
        byte = (uint8_t)(byte ^ HDLC_ESCAPE_XOR);
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

    *payloadLen = 0;

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

bool hdlc_payload_from_frame(const HdlcFrame *input,
                             uint8_t *payload,
                             size_t payloadSize,
                             size_t *payloadLen)
{
    size_t len;

    if (input == NULL || payload == NULL || payloadLen == NULL) {
        return false;
    }

    if (input->infoLen > HDLC_MAX_INFO) {
        return false;
    }

    len = input->infoLen + 2;
    if (len > payloadSize) {
        return false;
    }

    payload[0] = input->address;
    payload[1] = input->control;

    if (input->infoLen > 0) {
        memcpy(&payload[2], input->info, input->infoLen);
    }

    *payloadLen = len;

    return true;
}

bool hdlc_frame_from_payload(const uint8_t *payload,
                             size_t payloadLen,
                             HdlcFrame *output)
{
    if (payload == NULL || output == NULL) {
        return false;
    }

    if (payloadLen < 2) {
        return false;
    }

    if (payloadLen - 2 > HDLC_MAX_INFO) {
        return false;
    }

    memset(output, 0, sizeof(*output));

    output->address = payload[0];
    output->control = payload[1];
    output->infoLen = payloadLen - 2;

    if (output->infoLen > 0) {
        memcpy(output->info, &payload[2], output->infoLen);
    }

    return true;
}

bool hdlc_encode_frame(const HdlcFrame *input,
                       uint8_t *frame,
                       size_t frameSize,
                       size_t *frameLen)
{
    uint8_t payload[HDLC_MAX_FRAME];
    size_t payloadLen;

    if (!hdlc_payload_from_frame(input,
                                 payload,
                                 sizeof(payload),
                                 &payloadLen)) {
        return false;
    }

    return hdlc_encode(payload, payloadLen, frame, frameSize, frameLen);
}

bool hdlc_decode_frame(const uint8_t *frame,
                       size_t frameLen,
                       HdlcFrame *output)
{
    uint8_t payload[HDLC_MAX_FRAME];
    size_t payloadLen;

    if (output == NULL) {
        return false;
    }

    if (!hdlc_decode(frame,
                     frameLen,
                     payload,
                     sizeof(payload),
                     &payloadLen)) {
        return false;
    }

    return hdlc_frame_from_payload(payload, payloadLen, output);
}

bool hdlc_encode_addr_info(uint8_t address,
                           uint8_t control,
                           const uint8_t *info,
                           size_t infoLen,
                           uint8_t *frame,
                           size_t frameSize,
                           size_t *frameLen)
{
    HdlcFrame input;

    if (info == NULL && infoLen != 0) {
        return false;
    }

    if (infoLen > HDLC_MAX_INFO) {
        return false;
    }

    memset(&input, 0, sizeof(input));
    input.address = address;
    input.control = control;
    input.infoLen = infoLen;

    if (infoLen > 0) {
        memcpy(input.info, info, infoLen);
    }

    return hdlc_encode_frame(&input, frame, frameSize, frameLen);
}

uint8_t hdlc_i_ctrl(uint8_t ns, uint8_t nr, bool poll)
{
    uint8_t control;

    control = (uint8_t)(((ns & HDLC_SEQ_MASK) << 1) |
                        ((nr & HDLC_SEQ_MASK) << 5));

    if (poll) {
        control |= HDLC_CTRL_PF;
    }

    return control;
}

uint8_t hdlc_rr_ctrl(uint8_t nr, bool poll)
{
    uint8_t control;

    control = (uint8_t)(HDLC_CTRL_RR | ((nr & HDLC_SEQ_MASK) << 5));
    if (poll) {
        control |= HDLC_CTRL_PF;
    }

    return control;
}

uint8_t hdlc_rnr_ctrl(uint8_t nr, bool poll)
{
    uint8_t control;

    control = (uint8_t)(HDLC_CTRL_RNR | ((nr & HDLC_SEQ_MASK) << 5));
    if (poll) {
        control |= HDLC_CTRL_PF;
    }

    return control;
}

static uint8_t u_ctrl(uint8_t base, bool pf)
{
    if (pf) {
        return (uint8_t)(base | HDLC_CTRL_PF);
    }

    return base;
}

uint8_t hdlc_snrm_ctrl(bool poll)
{
    return u_ctrl(HDLC_CTRL_SNRM, poll);
}

uint8_t hdlc_disc_ctrl(bool poll)
{
    return u_ctrl(HDLC_CTRL_DISC, poll);
}

uint8_t hdlc_ua_ctrl(bool final)
{
    return u_ctrl(HDLC_CTRL_UA, final);
}

uint8_t hdlc_dm_ctrl(bool final)
{
    return u_ctrl(HDLC_CTRL_DM, final);
}

uint8_t hdlc_frmr_ctrl(bool final)
{
    return u_ctrl(HDLC_CTRL_FRMR, final);
}

uint8_t hdlc_xid_ctrl(bool poll)
{
    return u_ctrl(HDLC_CTRL_XID, poll);
}

bool hdlc_is_i_frame(uint8_t control)
{
    return ((control & 0x01) == 0);
}

bool hdlc_is_s_frame(uint8_t control)
{
    return ((control & 0x03) == 0x01);
}

bool hdlc_is_u_frame(uint8_t control)
{
    return ((control & 0x03) == 0x03);
}

bool hdlc_is_rr(uint8_t control)
{
    return hdlc_is_s_frame(control) && ((control & 0x0F) == HDLC_CTRL_RR);
}

bool hdlc_is_rnr(uint8_t control)
{
    return hdlc_is_s_frame(control) && ((control & 0x0F) == HDLC_CTRL_RNR);
}

static bool u_is(uint8_t control, uint8_t base)
{
    return hdlc_is_u_frame(control) && ((control & (uint8_t)~HDLC_CTRL_PF) == base);
}

bool hdlc_is_snrm(uint8_t control)
{
    return u_is(control, HDLC_CTRL_SNRM);
}

bool hdlc_is_disc(uint8_t control)
{
    return u_is(control, HDLC_CTRL_DISC);
}

bool hdlc_is_ua(uint8_t control)
{
    return u_is(control, HDLC_CTRL_UA);
}

bool hdlc_is_dm(uint8_t control)
{
    return u_is(control, HDLC_CTRL_DM);
}

bool hdlc_is_frmr(uint8_t control)
{
    return u_is(control, HDLC_CTRL_FRMR);
}

bool hdlc_is_xid(uint8_t control)
{
    return u_is(control, HDLC_CTRL_XID);
}

bool hdlc_pf(uint8_t control)
{
    return ((control & HDLC_CTRL_PF) != 0);
}

bool hdlc_i_ns(uint8_t control, uint8_t *ns)
{
    if (ns == NULL || !hdlc_is_i_frame(control)) {
        return false;
    }

    *ns = (uint8_t)((control >> 1) & HDLC_SEQ_MASK);

    return true;
}

bool hdlc_i_nr(uint8_t control, uint8_t *nr)
{
    if (nr == NULL || !hdlc_is_i_frame(control)) {
        return false;
    }

    *nr = (uint8_t)((control >> 5) & HDLC_SEQ_MASK);

    return true;
}

bool hdlc_s_nr(uint8_t control, uint8_t *nr)
{
    if (nr == NULL || !hdlc_is_s_frame(control)) {
        return false;
    }

    *nr = (uint8_t)((control >> 5) & HDLC_SEQ_MASK);

    return true;
}

const char *hdlc_control_name(uint8_t control)
{
    if (hdlc_is_i_frame(control)) {
        return "I";
    }

    if (hdlc_is_xid(control)) {
        return "XID";
    }

    if (hdlc_is_snrm(control)) {
        return "SNRM";
    }

    if (hdlc_is_disc(control)) {
        return "DISC";
    }

    if (hdlc_is_ua(control)) {
        return "UA";
    }

    if (hdlc_is_dm(control)) {
        return "DM";
    }

    if (hdlc_is_frmr(control)) {
        return "FRMR";
    }

    if (hdlc_is_rr(control)) {
        return "RR";
    }

    if (hdlc_is_rnr(control)) {
        return "RNR";
    }

    return "CTRL";
}

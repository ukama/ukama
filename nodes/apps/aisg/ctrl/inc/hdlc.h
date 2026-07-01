/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef HDLC_H_
#define HDLC_H_

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#define HDLC_FLAG                          0x7E
#define HDLC_ESCAPE                        0x7D
#define HDLC_ESCAPE_XOR                    0x20

/*
 * Maximum encoded frame buffer used by the AISG controller.
 * This includes flags, escaped bytes and FCS.
 */
#define HDLC_MAX_FRAME                     4096

/*
 * Maximum unescaped INFO field retained in HdlcFrame.
 * This is intentionally generous for the current controller buffers.
 */
#define HDLC_MAX_INFO                      (HDLC_MAX_FRAME - 8)

#define HDLC_ADDR_NO_STATION               0x00
#define HDLC_ADDR_BROADCAST                0xFF

#define HDLC_CTRL_PF                       0x10

/* U-frame base values without P/F bit. */
#define HDLC_CTRL_SNRM                     0x83
#define HDLC_CTRL_DISC                     0x43
#define HDLC_CTRL_UA                       0x63
#define HDLC_CTRL_DM                       0x0F
#define HDLC_CTRL_FRMR                     0x87
#define HDLC_CTRL_XID                      0xAF

/* S-frame base values. */
#define HDLC_CTRL_RR                       0x01
#define HDLC_CTRL_RNR                      0x05

#define HDLC_SEQ_MODULO                    8
#define HDLC_SEQ_MASK                      0x07

typedef struct {
    uint8_t address;
    uint8_t control;
    uint8_t info[HDLC_MAX_INFO];
    size_t infoLen;
} HdlcFrame;

/*
 * Phase-0/1 compatibility helpers.
 * payload is ADDR | CONTROL | INFO... and is the exact FCS input.
 */
bool hdlc_encode(const uint8_t *payload,
                 size_t payloadLen,
                 uint8_t *frame,
                 size_t frameSize,
                 size_t *frameLen);
bool hdlc_decode(const uint8_t *frame,
                 size_t frameLen,
                 uint8_t *payload,
                 size_t payloadSize,
                 size_t *payloadLen);

/* Layer-2 helpers used by Phase-2+ code. */
bool hdlc_encode_frame(const HdlcFrame *input,
                       uint8_t *frame,
                       size_t frameSize,
                       size_t *frameLen);
bool hdlc_decode_frame(const uint8_t *frame,
                       size_t frameLen,
                       HdlcFrame *output);
bool hdlc_encode_addr_info(uint8_t address,
                           uint8_t control,
                           const uint8_t *info,
                           size_t infoLen,
                           uint8_t *frame,
                           size_t frameSize,
                           size_t *frameLen);
bool hdlc_payload_from_frame(const HdlcFrame *input,
                             uint8_t *payload,
                             size_t payloadSize,
                             size_t *payloadLen);
bool hdlc_frame_from_payload(const uint8_t *payload,
                             size_t payloadLen,
                             HdlcFrame *output);

uint16_t hdlc_fcs16(const uint8_t *data, size_t len);

uint8_t hdlc_i_ctrl(uint8_t ns, uint8_t nr, bool poll);
uint8_t hdlc_rr_ctrl(uint8_t nr, bool poll);
uint8_t hdlc_rnr_ctrl(uint8_t nr, bool poll);
uint8_t hdlc_snrm_ctrl(bool poll);
uint8_t hdlc_disc_ctrl(bool poll);
uint8_t hdlc_ua_ctrl(bool final);
uint8_t hdlc_dm_ctrl(bool final);
uint8_t hdlc_frmr_ctrl(bool final);
uint8_t hdlc_xid_ctrl(bool poll);

bool hdlc_is_i_frame(uint8_t control);
bool hdlc_is_s_frame(uint8_t control);
bool hdlc_is_u_frame(uint8_t control);
bool hdlc_is_rr(uint8_t control);
bool hdlc_is_rnr(uint8_t control);
bool hdlc_is_snrm(uint8_t control);
bool hdlc_is_disc(uint8_t control);
bool hdlc_is_ua(uint8_t control);
bool hdlc_is_dm(uint8_t control);
bool hdlc_is_frmr(uint8_t control);
bool hdlc_is_xid(uint8_t control);
bool hdlc_pf(uint8_t control);

bool hdlc_i_ns(uint8_t control, uint8_t *ns);
bool hdlc_i_nr(uint8_t control, uint8_t *nr);
bool hdlc_s_nr(uint8_t control, uint8_t *nr);

const char *hdlc_control_name(uint8_t control);

#endif /* HDLC_H_ */

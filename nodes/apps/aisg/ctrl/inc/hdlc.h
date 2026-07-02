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
#define HDLC_MAX_FRAME                     4096
#define HDLC_MAX_INFO                      2048

/*
 * HDLC frame after FCS validation and byte unescaping.
 *
 * Encoded wire frame:
 *   0x7E | address | control | info... | fcs-low | fcs-high | 0x7E
 */
typedef struct {
    uint8_t address;
    uint8_t control;
    uint8_t info[HDLC_MAX_INFO];
    size_t infoLen;
} HdlcFrame;

/* Backward-compatible payload helpers: payload = address | control | info... */
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

/* Frame-native helpers for AISG Layer 2. */
bool hdlc_encode_frame(const HdlcFrame *decoded,
                       uint8_t *frame,
                       size_t frameSize,
                       size_t *frameLen);
bool hdlc_decode_frame(const uint8_t *frame,
                       size_t frameLen,
                       HdlcFrame *decoded);
bool hdlc_encode_addr_info(uint8_t address,
                           uint8_t control,
                           const uint8_t *info,
                           size_t infoLen,
                           uint8_t *frame,
                           size_t frameSize,
                           size_t *frameLen);
bool hdlc_payload_from_frame(const HdlcFrame *decoded,
                             uint8_t *payload,
                             size_t payloadSize,
                             size_t *payloadLen);
bool hdlc_frame_from_payload(const uint8_t *payload,
                             size_t payloadLen,
                             HdlcFrame *decoded);

uint16_t hdlc_fcs16(const uint8_t *data, size_t len);

/* HDLC UNC / AISG Layer-2 control helpers. */
uint8_t hdlc_i_ctrl(uint8_t ns, uint8_t nr, bool poll);
uint8_t hdlc_rr_ctrl(uint8_t nr, bool poll);
uint8_t hdlc_rnr_ctrl(uint8_t nr, bool poll);
uint8_t hdlc_snrm_ctrl(bool poll);
uint8_t hdlc_disc_ctrl(bool poll);
uint8_t hdlc_ua_ctrl(bool final);
uint8_t hdlc_dm_ctrl(bool final);
uint8_t hdlc_frmr_ctrl(bool final);
uint8_t hdlc_xid_ctrl(bool poll);

bool hdlc_is_i_frame(uint8_t ctrl);
bool hdlc_is_rr(uint8_t ctrl);
bool hdlc_is_rnr(uint8_t ctrl);
bool hdlc_is_snrm(uint8_t ctrl);
bool hdlc_is_disc(uint8_t ctrl);
bool hdlc_is_ua(uint8_t ctrl);
bool hdlc_is_dm(uint8_t ctrl);
bool hdlc_is_frmr(uint8_t ctrl);
bool hdlc_is_xid(uint8_t ctrl);
bool hdlc_poll_final(uint8_t ctrl);
uint8_t hdlc_ns(uint8_t ctrl);
uint8_t hdlc_nr(uint8_t ctrl);

#endif /* HDLC_H_ */

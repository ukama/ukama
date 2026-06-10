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
uint16_t hdlc_fcs16(const uint8_t *data, size_t len);

#endif /* HDLC_H_ */

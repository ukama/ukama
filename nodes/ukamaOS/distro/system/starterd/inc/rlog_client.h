/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#pragma once

#include <jansson.h>
#include <stdbool.h>
#include <stdint.h>

typedef struct RlogClient RlogClient;

RlogClient *rlog_client_create(const char *socketPath,
                               const char *producerBootId,
                               int reconnectMs);
void rlog_client_destroy(RlogClient *client);
bool rlog_client_send_frame(RlogClient *client, json_t *frame);
bool rlog_client_send(RlogClient *client,
                      json_t *record,
                      const char *space,
                      const char *app,
                      int pid,
                      uint32_t generation,
                      const char *stream,
                      uint64_t captureSeq,
                      const char *captureId);

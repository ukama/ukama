/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AISG_EMU_RESPONSE_H_
#define AISG_EMU_RESPONSE_H_

#include <stdbool.h>

#include "emu.h"

typedef struct {
    char     id[64];
    bool     ok;
    char     code[64];
    char     reason[128];
    JsonObj *payload;
} EmuResponse;

void emu_response_init(EmuResponse *response, const char *id);
void emu_response_free(EmuResponse *response);

bool emu_response_set_ok(EmuResponse *response, JsonObj *payload);
bool emu_response_set_error(EmuResponse *response,
                            const char *code,
                            const char *reason);

JsonObj *emu_response_to_json(EmuResponse *response);

#endif /* AISG_EMU_RESPONSE_H_ */

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AISG_EMU_MODEL_H_
#define AISG_EMU_MODEL_H_

#include "emu.h"

typedef struct {
    bool present;
    bool configured;
    bool calibrated;
    bool busy;
    bool alarmSubscribed;
    int16_t tiltTenthsDeg;
    JsonObj *alarms;
} EmuModel;

void emu_model_init(EmuModel *model);
void emu_model_free(EmuModel *model);
bool emu_model_load_scenario(EmuModel *model, const char *scenario);
JsonObj *emu_model_status(EmuModel *model);

#endif /* AISG_EMU_MODEL_H_ */

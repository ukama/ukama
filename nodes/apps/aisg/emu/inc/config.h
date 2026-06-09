/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AISG_EMU_CONFIG_H_
#define AISG_EMU_CONFIG_H_

#include "emu.h"

#define AISG_EMU_CONFIG_FILE "config/config.toml"
#define AISG_EMU_SOCKET_PATH "/var/run/aisg-ctrl.sock"
#define AISG_EMU_SERVICE_NAME "aisg-emu"
#define AISG_EMU_DEFAULT_PORT 18112

typedef struct {
    char socketPath[128];
    char scenario[64];
    int servicePort;
} EmuConfig;

void emu_config_init(EmuConfig *config);
bool emu_config_load(EmuConfig *config, const char *file);
void emu_config_free(EmuConfig *config);

#endif /* AISG_EMU_CONFIG_H_ */

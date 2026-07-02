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

#include "usys_services.h"

#define AISG_EMU_CONFIG_FILE  "/ukama/configs/aisg-emu/config.toml"
#define AISG_EMU_SOCKET_PATH  "/var/run/aisg-ctrl.sock"
#define AISG_EMU_RET_PTY_PATH "/tmp/aisg-ret0"
#define AISG_EMU_SERVICE_NAME SERVICE_AISG_EMU

typedef enum {
    EmuModeContract = 0,
    EmuModeRet
} EmuMode;

typedef struct {
    EmuMode mode;
    char modeName[16];

    /* Existing controller-contract emulator. */
    char socketPath[128];
    char scenario[64];
    int  servicePort;

    /* Protocol-accurate single-antenna RET emulator. */
    char retPtyPath[128];
    char retVendorCode[3];
    char retSerial[18];
    int retRequiresConfig;
    int retInitialTiltTenths;
    int retMinTiltTenths;
    int retMaxTiltTenths;
    int retCalibrateDelayMs;
    int retMoveDelayMs;
} EmuConfig;

bool emu_config_init(EmuConfig *config);
bool emu_config_load(EmuConfig *config, const char *file);
void emu_config_free(EmuConfig *config);
bool emu_config_set_mode(EmuConfig *config, const char *mode);
bool emu_config_set_bool(int *dst, const char *value);
int emu_config_tilt_arg_to_tenths(const char *value, int fallback);

#endif /* AISG_EMU_CONFIG_H_ */

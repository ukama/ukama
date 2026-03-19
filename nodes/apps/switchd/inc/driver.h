/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SWITCHD_DRIVER_H
#define SWITCHD_DRIVER_H

#include "types.h"

typedef struct {
    int (*init)(SwitchdContext *ctx);
    void (*cleanup)(SwitchdContext *ctx);
    int (*probe)(SwitchdContext *ctx);
    int (*refresh_info)(SwitchdContext *ctx,
                        SwitchInfo *info,
                        SwitchCapabilities *caps);
    int (*refresh_ports)(SwitchdContext *ctx,
                         SwitchPortState *ports,
                         uint32_t *count);
    int (*refresh_kpis)(SwitchdContext *ctx, SwitchKpis *kpis);
    int (*set_port_admin)(SwitchdContext *ctx, uint32_t portId, bool up);
    int (*set_port_poe)(SwitchdContext *ctx, uint32_t portId, bool on);
    int (*save_config)(SwitchdContext *ctx);
    int (*reboot_switch)(SwitchdContext *ctx);
    int (*firmware_apply)(SwitchdContext *ctx, const SwitchFirmware *fw);
} SwitchDriverOps;

struct SwitchDriver {
    const char *name;
    SwitchDriverOps ops;
    void *priv;
};

int driver_init(SwitchdContext *ctx);
void driver_cleanup(SwitchdContext *ctx);

#endif

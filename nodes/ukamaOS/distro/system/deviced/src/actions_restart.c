/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <unistd.h>
#include <sys/reboot.h>

#include "actions.h"
#include "deviced.h"

static inline bool is_debug_mode(void) {
    return getenv(ENV_DEVICED_DEBUG_MODE) != NULL;
}

int actions_restart_apply(Config *config) {

    (void)config;

    if (is_debug_mode()) {
        return STATUS_OK;
    }

    sync();
    (void)setuid(0);
    reboot(RB_AUTOBOOT);

    return STATUS_OK;
}

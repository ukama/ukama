/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <errno.h>
#include <string.h>
#include <unistd.h>
#include <sys/reboot.h>

#include "actions.h"
#include "deviced.h"

#include "usys_log.h"

static inline bool is_debug_mode(void) {
    return getenv(ENV_DEVICED_DEBUG_MODE) != NULL;
}

int actions_reboot_apply(Config *config) {

    (void)config;

    if (is_debug_mode()) {
        usys_log_info("reboot: debug mode, skipping OS reboot");
        return STATUS_OK;
    }

    sync();

    /* redundant as starter.d is running priviledged,
     * just double check.
     */
    if (setuid(0) != 0) {
        usys_log_error("reboot: setuid(0) failed: %s", strerror(errno));
        return STATUS_NOK;
    }

    if (reboot(RB_AUTOBOOT) != 0) {
        usys_log_error("reboot: reboot syscall failed: %s", strerror(errno));
        return STATUS_NOK;
    }

    return STATUS_OK;
}

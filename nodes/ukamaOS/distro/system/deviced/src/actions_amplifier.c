/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <fcntl.h>
#include <unistd.h>
#include <string.h>

#include "actions.h"
#include "deviced.h"

#include "usys_log.h"

static int _write_sysfs_value(const char *path, const char *value) {

    int fd = -1;
    ssize_t wr = 0;

    if (!path || !*path || !value) return STATUS_NOK;

    fd = open(path, O_WRONLY);
    if (fd < 0) {
        return STATUS_NOK;
    }

    wr = write(fd, value, (size_t)strlen(value));
    close(fd);

    if (wr < 0) return STATUS_NOK;
    return STATUS_OK;
}

static int _set_chain(const char *base, const char *val) {

    int ret = STATUS_NOK;
    char path[256];

    memset(path, 0, sizeof(path));
    snprintf(path, sizeof(path), "%s/tx_enable", base);
    ret = _write_sysfs_value(path, val);
    if (ret != STATUS_OK) return ret;

    memset(path, 0, sizeof(path));
    snprintf(path, sizeof(path), "%s/PA_disable", base);
    ret = _write_sysfs_value(path, val);
    if (ret != STATUS_OK) return ret;

    memset(path, 0, sizeof(path));
    snprintf(path, sizeof(path), "%s/rx_enable", base);
    ret = _write_sysfs_value(path, val);
    if (ret != STATUS_OK) return ret;

    return STATUS_OK;
}

int actions_radio_apply(Config *config, ControlState desired) {

    const char *val = NULL;
    int ret = STATUS_NOK;

    (void)config;

    val = (desired == CONTROL_STATE_ON) ? "1" : "0";

    usys_log_info("radio: %s", (desired == CONTROL_STATE_ON) ? "on" : "off");

    ret = _set_chain("/devices/platform/fema1-gpios", val);
    if (ret != STATUS_OK) return ret;

    ret = _set_chain("/devices/platform/fema2-gpios", val);
    if (ret != STATUS_OK) return ret;

    return STATUS_OK;
}

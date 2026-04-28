/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include "utils/mfg_helper.h"

#include "usys_log.h"
#include "usys_string.h"

/* Board/module names supported by genSchema/genInventory */
static char* boardname[MODULE_TYPE_MAX] = {
    "com",
    "trx",
    "mask",
    "ctrl",
    "fe1",
    "fe2",
    "cm4"
};

int verify_uuid(char* uuid) {

    int ret;
    unsigned int len;

    ret = 0;
    len = usys_strlen(uuid);

    log_trace("MFG:: UUID %s Length %d \n", uuid, len);

    if ((!len) || (len > UUID_MAX_LENGTH)) {
        usys_log_error("MFG:: Error:: UUID length should be greater than 0 "
                       "and less than 24 characters.");
        ret = -1;
    }

    return ret;
}

int verify_board_name(char* name) {

    int ret;
    unsigned int len;
    unsigned short int idx;

    ret = -1;
    len = usys_strlen(name);

    log_trace("MFG:: Name %s Length %d \n", name, len);

    if ((!len) || (len > NAME_MAX_LENGTH)) {
        usys_log_error("MFG:: Error:: Name length should be greater than 0 "
                       "and less than 24 characters.");
        return ret;
    }

    for (idx = 0; idx < MODULE_TYPE_MAX; idx++) {
        if (!usys_strcmp(name, boardname[idx])) {
            ret = 0;
            break;
        }
    }

    if (ret) {
        usys_log_error("MFG:: Error:: Check the module name %s.", name);
    }

    return ret;
}

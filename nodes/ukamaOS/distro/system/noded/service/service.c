/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include "service.h"

#include "inventory.h"
#include "ledger.h"
#include "store.h"
#include "web_service.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

RegisterDeviceCB registerDeviceCB;


/* Initialize inventory on the boot up
 * At boot up all modules have there respective inventory databases.
 * These databases may be exposed as files or had to be read from eeprom
 * using routines provided by store. Store abstract the access to the databases.
 */
int service_at_bootup(char *invtDb, char *propCfg, char *notifServer) {
    int ret = 0;

    ret = web_service_init(notifServer);
    if (ret) {
        return ret;
    }

    ret = ldgr_init(propCfg);
    if (!ret) {
        invt_init(invtDb, &ldgr_register);
    }

    return ret;
}

int service_at_exit() {
    int ret = STATUS_OK;
    ldgr_exit();
    invt_exit();
    return ret;
}

int service_init(char *invtDb, char *propCfg, char* notifServer) {
    int ret = STATUS_OK;
    ret = service_at_bootup(invtDb, propCfg, notifServer);
    return ret;
}

void service() {
    web_service_start();
}

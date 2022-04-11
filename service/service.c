/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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

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
#include "usys_log.h"
#include "usys_types.h"

RegisterDeviceCB registerDeviceCB;

/* Initialize inventory on the bootup
 * At bootup all modules have there respective inventory databases.
 * These databases may be exposed as files or had to be read from eeprom
 * using routines provided by store. Store abstract the access to the databases.
 */
int service_at_bootup(char *pinvtDb) {
    int ret = 0;
    uint8_t count = 1;

    /* Initialize Store*/
    ret = store_init();
    UnitCfg *cfg = usys_zmalloc(sizeof(UnitCfg));
    if (cfg) {
        ret = get_master_db_info(cfg, pinvtDb);
        if (ret) {
            usys_free(cfg);
            cfg = NULL;
            return ret;
        }
    }

    invt_print_unit_cfg(cfg, count);

    /* Register master module first so that unit info and cfg can be accessed.*/
    ret = invt_register_module(cfg);
    /* After registering master module.
     * Access remaining module if any using unit cfg and register those.*/
    if (!ret) {

        /* Check if the database exist or not.*/
        ret = invt_validating_magic_word(cfg->modUuid);
        if (ret) {
            usys_log_warn("Inventory No Database found for module UUID %s "
                            "Name %s.",
                            cfg->modUuid, cfg->modName);
        } else {
            /* Register other modules if any.*/
            /* Caution:: If Module is registering itself it may not have more
             * modules to register.*/
            ret = invt_register_modules(cfg->modUuid, registerDeviceCB );
        }
    }

    if (cfg) {
        usys_free(cfg->eepromCfg);
        cfg->eepromCfg = NULL;
        usys_free(cfg);
        cfg = NULL;
    }

    return ret;
}

int service_at_exit() {
    int ret = USYS_OK;
    return ret;
}

int service_init() {
    int ret = USYS_OK;

    return ret;
}

void service() {
}

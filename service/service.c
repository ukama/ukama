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

#if 0
/* Writing a newly parsed json to DB (sysfs/eeprom) */
int service_create_invt_db(char **puuid, char** name, char** schema, int count) {
    int ret = 0;

    JSONInput jip;
    jip.fname = schema;
    jip.count = count;
    jip.pname = PROPERTYJSON;

    UnitCfg *pcfg = NULL;
    DevI2cCfg *i2c_cfg = NULL;

    /* Initializes for device DB */
    ret = ubsp_devdb_init(jip.pname);
    if (ret) {
        log_error("MFGUTIL:: UBSP DEVDB init failed %d", ret);
        goto cleanup;
    }

    /* Initializes for input DB like cstructs or JSON file */
    ret = ubsp_idb_init(&jip);
    if (ret) {
        log_error("MFGUTIL:: UBSP IDB init failed %d", ret);
        goto cleanup;
    }

    /* Will just initialize the db if NULL is passed*/
    ret = ubsp_ukdb_init(NULL);
    if (ret) {
        log_info("MFGUTIL:: UBSP init failed %d (Expected -1)", ret);
    }

    for(int idx = 0; idx < count; idx++) {
        log_debug("UUID[%d] = %24s Name[%d] = %24s Schema[%d] = %s \n", idx, puuid[idx], idx, name[idx], idx, jip.fname[idx]);

        /* Assumption Module Name in argument should match */
        UnitCfg *udata = (UnitCfg[]){
            { .mod_uuid = "UK-5001-RFC-1101",
                .mod_name = "RF CTRL BOARD",
                .sysfs = "/tmp/sys/bus/i2c/devices/i2c-0/0-0051/eeprom",
                .eeprom_cfg = &(DevI2cCfg){ .bus = 1, .add = 0x50ul } },
                { .mod_uuid = "UK-4001-RFA-1101",
                        .mod_name = "RF BOARD",
                        .sysfs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0052/eeprom",
                        .eeprom_cfg = &(DevI2cCfg){ .bus = 2, .add = 0x50ul } },
                        { .mod_uuid = "UK-1001-COM-1101",
                                .mod_name = "ComV1",
                                .sysfs = "/tmp/sys/bus/i2c/devices/i2c-0/0-0050/eeprom",
                                .eeprom_cfg = &(DevI2cCfg){ .bus = 0, .add = 0x50ul } },
                                { .mod_uuid = "UK-2001-LTE-1101",
                                        .mod_name = "LTE",
                                        .sysfs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0050/eeprom",
                                        .eeprom_cfg = &(DevI2cCfg){ .bus = 1, .add = 0x50ul } },
                                        { .mod_uuid = "UK-3001-MSK-1101",
                                                .mod_name = "MASK",
                                                .sysfs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0051/eeprom",
                                                .eeprom_cfg = &(DevI2cCfg){ .bus = 1, .add = 0x51ul } },
        };

        /* Find and Read unitCfg of the module from above UnitCfg struct */
        for (int iter = 0; iter < MAX_BOARDS; iter++) {
            if (!strcmp(name[idx], udata[iter].mod_name)) {

                pcfg = malloc(sizeof(UnitCfg));
                if (pcfg) {

                    memset(pcfg, '\0', sizeof(UnitCfg));
                    memcpy(pcfg, &udata[iter], sizeof(UnitCfg));

                    if (udata[iter].eeprom_cfg) {

                        i2c_cfg = malloc(sizeof(DevI2cCfg));
                        if (i2c_cfg) {
                            memset(i2c_cfg, '\0', sizeof(DevI2cCfg));
                            memcpy(i2c_cfg, udata[iter].eeprom_cfg,
                                    sizeof(DevI2cCfg));
                        }

                    }

                    pcfg->eeprom_cfg = i2c_cfg;
                    memcpy(pcfg->mod_uuid, puuid[idx], strlen(puuid[idx]));
                    memcpy(pcfg->mod_name, name[idx], strlen(name[idx]));

                    break;

                } else {

                    log_error("MFGUTIL:: Err(%d): Memory exhausted while getting unit"
                            " config from Test data.",
                            ERR_UBSP_MEMORY_EXHAUSTED);
                    goto cleanup;

                }
            }
        }

        /* Register Module */
        ret = ubsp_register_module(pcfg);
        if (!ret) {

            /* Create a EEPROM DB */
            ret = ubsp_create_ukdb(pcfg->mod_uuid);
            if (!ret) {
                log_info("MFGUTIL:: Created registry for module %s.", name[idx]);
            } else {
                log_error("MFGUTIL:: UBSP registry creation failed %d", ret);
                goto cleanup;
            }

        } else {
            log_error("MFGUTIL:: Registering module failed %d", ret);
            goto cleanup;
        }
    }

    /* Cleanup */
    cleanup:
    ubsp_idb_exit();
    ubsp_exit();
    UBSP_FREE(pcfg->eeprom_cfg);
    UBSP_FREE(pcfg);

    return ret;
}


/* Service request for create inevtory database */
int service_req_create_invt_database() {
    /* Args check for schema info.*/
    if ((sidx != uidx) || (sidx != nidx) || (!sidx) || (sidx > MAX_BOARDS)  ) {
        log_error("MFGUTIL:: Name, schema and UUID entries have to match in count.");
        log_error("MFGUTIL:: At least one set of entries or %d set of entries can be made simultaneously.", MAX_BOARDS);
        exit(0);
    }

    /* Input args and their verification */
    for(int idx = 0; idx < uidx;idx++) {
        /* Verify module uuid and name */
        if (verify_uuid(uuid[idx]) || verify_boardname(name[idx]) ) {
            usage();
            exit(0);
        }
        log_trace("UUID[%d] = %24s Name[%d] = %24s Schema[%d] = %s \n", idx, uuid[idx], idx, name[idx], idx, schema[idx]);
    }

    /* Create EEPROM DB */
    int ret = create_db_hook(uuid, name, schema, uidx);
    if (ret) {
        log_error("MFGUTIL:: Error:: Failed to create registry DB for %s device.", name);
    } else {
        log_info("MFGUTIL:: Created registry DB for device.");
        log_info("MFGUTIL:: Copy directory from /tmp/sys");
    }
}
#endif
/* Initialize inventory on the bootup
 * At bootup all modules have there respective inventory databases.
 * These databases may be exposed as files or had to be read from eeprom
 * using routines provided by store. Store abstract the access to the databases.
 */
int service_at_bootup(char *invtDb, char *propCfg) {
    int ret = 0;

    ret = web_service_init();
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
    int ret = USYS_OK;
    ldgr_exit();
    invt_exit();
    return ret;
}

int service_init(char *invtDb, char *propCfg) {
    int ret = USYS_OK;
    ret = service_at_bootup(invtDb, propCfg);
    return ret;
}

void service() {
    web_service_start();
}

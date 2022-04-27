/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "ukdb/db/db.h"

#include "headers/errorcode.h"
#include "headers/ubsp/devices.h"
#include "headers/utils/list.h"
#include "headers/utils/log.h"
#include "ukdb/db/eeprom.h"
#include "ukdb/db/file.h"

const DBFxnTable *db_fxn_table;
static ListInfo moduledb;
static int moduledbflag = 0;

const DBFxnTable eeprom_fxn_table = { .init = eeprom_init,
                                      .read_block = eeprom_read,
                                      .write_block = eeprom_write,
                                      .read_number = eeprom_read_number,
                                      .write_number = eeprom_write_number,
                                      .erase_block = eeprom_erase,
                                      .write_protect = eeprom_protect,
                                      .rename = eeprom_rename,
                                      .cleanup = eeprom_cleanup,
                                      .remove = eeprom_remove };

const DBFxnTable file_fxn_table = { .init = file_init,
                                    .read_block = file_read,
                                    .write_block = file_write,
                                    .read_number = file_read_number,
                                    .write_number = file_write_number,
                                    .erase_block = file_erase,
                                    .write_protect = file_protect,
                                    .rename = file_rename,
                                    .cleanup = file_cleanup,
                                    .remove = file_remove };

static void free_moduledb(void *ip) {
    ListNode *node = (ListNode *)ip;
    if (node) {
        ModuleDBMap *map = node->data;
        UBSP_FREE(map->db_attr);
        UBSP_FREE(map);
        free(node);
    }
}

int compare_moduledb(void *ipt, void *sd) {
    ModuleDBMap *ip = (ModuleDBMap *)ipt;
    ModuleDBMap *op = (ModuleDBMap *)sd;
    int ret = 0;
    /* If module uuid matches it means module is same.*/
    if (!strcmp(ip->mod_uuid, op->mod_uuid)) {
        ret = 1;
    }
    return ret;
}

/* Searching device in the device list*/
static ModuleDBMap *search_module_in_db(char *puuid) {
    ModuleDBMap *smap = NULL;
    ModuleDBMap *fmap = NULL;
    if (!strcmp(puuid, "")) {
        /* If puuid is empty string. This means request for master module
		 * This assumption is made because on-bootup ubsp wouldn't know
		 * what's the master module uuid. Also it's made sure that first module stored
		 * in the list is the master module during the registration process.
		 */
        fmap = list_head(&moduledb, 0);
        if (fmap) {
            log_trace("DB:: Map found for master module UUID %s Name %s "
                      "in Module map.",
                      fmap->mod_uuid, fmap->mod_name);
        }
    } else {
        smap = malloc(sizeof(ModuleDBMap));
        if (smap) {
            memset(smap->mod_uuid, '\0', MAX_NAME_LENGTH);
            memcpy(smap->mod_uuid, puuid, strlen(puuid));
            /*Search return 1 for found.*/
            /* Remember list_search do shallow copy.*/
            fmap = list_search(&moduledb, smap);
            if (fmap) {
                log_trace("DB::Module UUID: %s found.", smap->mod_uuid);
            } else {
                UBSP_FREE(fmap);
                log_warn("DB:: Module UUID: %s not found.", puuid);
            }
            UBSP_FREE(smap);
        }
    }
    return fmap;
}

static int update_module_db(ModuleDBMap *map) {
    int ret = 0;
    if (map) {
        ret = list_update(&moduledb, map);
        if (ret) {
            log_error(
                "Err(%d): DB:: Module UUID %s and Name %s update to module DB failed.",
                ret, map->mod_uuid, map->mod_name);
        } else {
            log_trace(
                "DB:: Module UUID %s and Name %s update to module DB completed.",
                map->mod_uuid, map->mod_name);
        }
    }
    return ret;
}

ListInfo *get_moduledb() {
    if (moduledbflag == 0) {
        list_new(&moduledb, sizeof(ModuleDBMap), free_moduledb,
                 compare_moduledb, NULL);
        moduledbflag = 1;
        log_trace("DB:: Module DB initialized.");
    }
    return &moduledb;
}

int db_init() {
    int ret = 0;
    list_new(&moduledb, sizeof(ModuleDBMap), free_moduledb, compare_moduledb,
             NULL);
    moduledbflag = 1;
    log_trace("DB:: DB layer initialized.");
    return ret;
}

int db_register_module(UnitCfg *pcfg) {
    int ret = -1;
    ModuleDBMap *pmap = NULL;
    DevI2cCfg *icfg = NULL;
    char *fname = NULL;
    if (pcfg) {
        /* Check If module is already registered.*/
        pmap = search_module_in_db(pcfg->mod_uuid);
        if (pmap) {
            log_warn("DB:: Map for Module %s is already existing.",
                     pcfg->mod_uuid);
            ret = 0;
        } else {
            log_debug("DB:: No map for Module %s found. Adding now.",
                      pcfg->mod_uuid);
            pmap = malloc(sizeof(ModuleDBMap));
            if (pmap) {
                memset(pmap->mod_uuid, '\0', MAX_NAME_LENGTH);
                memset(pmap->mod_name, '\0', MAX_NAME_LENGTH);
                memcpy(pmap->mod_uuid, pcfg->mod_uuid, strlen(pcfg->mod_uuid));
                memcpy(pmap->mod_name, pcfg->mod_name, strlen(pcfg->mod_name));
                if ((pcfg->sysfs != NULL) || (!strcmp(pcfg->sysfs, ""))) {
                    fname = malloc(sizeof(char) * MAX_PATH_LENGTH);
                    if (fname) {
                        memset(fname, '\0', MAX_PATH_LENGTH);
                        memcpy(fname, pcfg->sysfs, strlen(pcfg->sysfs));
                    } else {
                        ret = ERR_UBSP_MEMORY_EXHAUSTED;
                        log_error(
                            "Err(%d):: DB Memory exhausted while adding sys file name to db.",
                            ret);
                        goto cleanup;
                    }
                    pmap->db_attr = fname;
                    pmap->fxn_tbl = &file_fxn_table;
                } else {
                    icfg = malloc(sizeof(icfg));
                    if (icfg) {
                        memcpy(icfg, pcfg->eeprom_cfg, sizeof(DevI2cCfg));
                    } else {
                        ret = ERR_UBSP_MEMORY_EXHAUSTED;
                        log_error(
                            "Err(%d):: DB Memory exhausted while adding eeprom cfg to db.",
                            ret);
                        goto cleanup;
                    }
                    pmap->db_attr = icfg;
                    pmap->fxn_tbl = &eeprom_fxn_table;
                }
                list_append(&moduledb, pmap);
                /* Calling Module init.*/
                ret = pmap->fxn_tbl->init(pmap->db_attr);
                if (ret) {
                    log_debug("DB:: Module %s registration success.",
                              pcfg->mod_uuid);
                }
            } else {
                ret = ERR_UBSP_MEMORY_EXHAUSTED;
                log_error(
                    "Err(%d): DB:: Map for Module %s is not added. Memory exhausted.",
                    ret, pcfg->mod_uuid);
            }
        }
        UBSP_FREE(pmap);
    }
    if (ret) {
        log_error("Err: DB:: Module %s registration failed.", pcfg->mod_uuid);
    }
    return ret;
cleanup:
    UBSP_FREE(fname);
    UBSP_FREE(icfg);
    UBSP_FREE(pmap);
    return ret;
}

void db_unregister_all_module() {
    log_debug("DB:: Cleaning all entries from Module DB.");
    list_destroy(&moduledb);
}

int db_unregister_module(char *puuid) {
    int ret = -1;
    int found = 0;
    if (puuid) {
        /* Check If module is already registered.*/
        ModuleDBMap *pmap = search_module_in_db(puuid);
        if (pmap) {
            log_warn("DB:: Map for Module %s found in moduleDB.", puuid);
            list_remove(&moduledb, pmap);
            log_debug("DB:: Module %s un-registration success.", puuid);
            ret = 0;
            UBSP_FREE(pmap);
        } else {
            log_debug("DB:: No map for Module %s found.", puuid);
        }
    } else {
        log_warn("DB:: UnitCfg is invalid for unregistering module.");
    }
    return ret;
}

int db_register_update_module(char *puuid, UnitCfg *pcfg, uint8_t count) {
    int ret = -1;
    DevI2cCfg *icfg = NULL;
    char *fname = NULL;
    ModuleDBMap *pmap = NULL;
    if (pcfg) {
        if (!puuid) {
            return ret;
        }
        /* Check If module is already registered.*/
        pmap = search_module_in_db(puuid);
        if (pmap) {
            log_trace("DB:: Map for Module %s found in moduleDB.", puuid);
            memset(pmap->mod_name, '\0', MAX_NAME_LENGTH);
            memcpy(pmap->mod_name, pcfg->mod_name, strlen(pcfg->mod_name));
            /* Need to free the db_attr as we don know new one could be same or different.*/
            UBSP_FREE(pmap->db_attr);
            if ((pcfg->sysfs != NULL) || (!strcmp(pcfg->sysfs, ""))) {
                fname = malloc(sizeof(char));
                if (fname) {
                    pmap->db_attr = fname;
                } else {
                    ret = ERR_UBSP_MEMORY_EXHAUSTED;
                    log_error(
                        "Err(%d):: DB Memory exhausted while adding sys file name to db.",
                        ret);
                    goto cleanuppmap;
                }
                pmap->fxn_tbl = &file_fxn_table;
            } else {
                icfg = malloc(sizeof(icfg));
                if (icfg) {
                    memcpy(icfg, pcfg->eeprom_cfg, sizeof(DevI2cCfg));
                } else {
                    ret = ERR_UBSP_MEMORY_EXHAUSTED;
                    log_error(
                        "Err(%d):: DB Memory exhausted while adding eeprom cfg to db.",
                        ret);
                    goto cleanuppmap;
                }
                pmap->db_attr = icfg;
                pmap->fxn_tbl = &eeprom_fxn_table;
            }
            ret = update_module_db(pmap);
            if (!ret) {
                log_error("Err(%d):DB:: Updating module uiid %s db failed.",
                          ret, puuid);
                goto cleanup;
            } else {
                log_debug("DB:: Module %s update success.", pcfg->mod_uuid);
            }
            ret = pmap->fxn_tbl->rename(puuid, pcfg->mod_name);
            if (!ret) {
                log_error("Err(%d):DB:: Renaming db for module uuid %s failed.",
                          ret, puuid);
                goto cleanup;
            } else {
                log_debug("DB:: Module %s update success.", pcfg->mod_uuid);
            }
            UBSP_FREE(pmap);
        } else {
            log_debug("DB:: No map for Module %s found.", puuid);
        }
    } else {
        log_warn("DB:: Module UUID is invalid for unregistering module.");
    }
    return ret;

cleanup:
    UBSP_FREE(fname);
    UBSP_FREE(icfg);
cleanuppmap:
    UBSP_FREE(pmap);
    return ret;
}

/* Select the data base based on the UUID of the module.*/
ModuleDBMap *db_choose_module(char *puuid) {
    ModuleDBMap *pmap = NULL;
    if (puuid) {
        /* Check If module is already registered.*/
        pmap = search_module_in_db(puuid);
        if (pmap) {
            log_trace("DB:: Map for Module %s found in moduleDB.", puuid);

        } else {
            log_warn("DB:: No map for Module %s found.", puuid);
        }
    } else {
        log_warn("Err(%d): DB:: Module UUID is invalid",
                 ERR_UBSP_INVALID_POINTER);
    }
    return pmap;
}

int db_rename(char *p_uuid, char *old_name, char *new_name) {
    int ret = -1;
    ModuleDBMap *pmap = db_choose_module(p_uuid);
    if (pmap) {
        if (pmap->fxn_tbl) {
            ret = pmap->fxn_tbl->rename(old_name, new_name);
        }
    }
    return ret;
}

int db_read_block(char *p_uuid, void *data, off_t offset, uint16_t size) {
    int ret = -1;
    ModuleDBMap *pmap = db_choose_module(p_uuid);
    if (pmap) {
        if (pmap->fxn_tbl) {
            ret = pmap->fxn_tbl->read_block(pmap->db_attr, data, offset, size);
        }
        UBSP_FREE(pmap);
    }
    return ret;
}

int db_write_block(char *p_uuid, void *data, off_t offset, uint16_t size) {
    int ret = -1;
    ModuleDBMap *pmap = db_choose_module(p_uuid);
    if (pmap) {
        if (pmap->fxn_tbl) {
            ret = pmap->fxn_tbl->write_block(pmap->db_attr, data, offset, size);
        }
        UBSP_FREE(pmap);
    }
    return ret;
}

int db_erase_block(char *p_uuid, off_t offset, uint16_t size) {
    int ret = -1;
    ModuleDBMap *pmap = db_choose_module(p_uuid);
    if (pmap) {
        if (pmap->fxn_tbl) {
            ret = pmap->fxn_tbl->erase_block(pmap->db_attr, offset, size);
        }
        UBSP_FREE(pmap);
    }
    return ret;
}
int db_read_number(char *p_uuid, void *val, off_t offset, uint16_t count,
                   uint8_t size) {
    int ret = -1;
    ModuleDBMap *pmap = db_choose_module(p_uuid);
    if (pmap) {
        if (pmap->fxn_tbl) {
            ret = pmap->fxn_tbl->read_number(pmap->db_attr, val, offset, count,
                                             size);
        }
        UBSP_FREE(pmap);
    }
    return ret;
}

int db_write_number(char *p_uuid, void *val, off_t offset, uint16_t count,
                    uint8_t size) {
    int ret = -1;
    ModuleDBMap *pmap = db_choose_module(p_uuid);
    if (pmap) {
        if (pmap->fxn_tbl) {
            ret = pmap->fxn_tbl->write_number(pmap->db_attr, val, offset, count,
                                              size);
        }
        UBSP_FREE(pmap);
    }
    return ret;
}

int db_write_protect(char *p_uuid, void *data) {
    int ret = -1;
    ModuleDBMap *pmap = db_choose_module(p_uuid);
    if (pmap) {
        if (pmap->fxn_tbl) {
            ret = pmap->fxn_tbl->write_protect(pmap->db_attr);
        }
        UBSP_FREE(pmap);
    }
    return ret;
}

int db_remove_database(char *p_uuid) {
    int ret = -1;
    ModuleDBMap *pmap = db_choose_module(p_uuid);
    if (pmap) {
        if (pmap->fxn_tbl) {
            ret = pmap->fxn_tbl->remove(pmap->db_attr);
        }
        UBSP_FREE(pmap);
    }
    return ret;
}

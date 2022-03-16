/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "../../inc/store.h"

#include "errorcode.h"
#include "eeprom_wrapper.h"
#include "schema.h"
#include "noded_macros.h"

#include "usys_list.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

static ListInfo modList;
static int modListFlag = 0;

const StoreOperations eepromOps = { .init = eeprom_wrapper_init,
                .readBlock = eeprom_wrapper_read,
                .writeBlock = eeprom_wrapper_write,
                .readNumber = eeprom_wrapper_read_number,
                .writeNumber = eeprom_wrapper_write_number,
                .eraseBlock = eeprom_wrapper_erase,
                .writeProtect = eeprom_wrapper_protect,
                .rename = eeprom_wrapper_rename,
                .cleanup = eeprom_wrapper_cleanup,
                .remove = eeprom_wrapper_remove };

const StoreOperations fileOps = { .init = usys_file_init,
                .readBlock = usys_file_read,
                .writeBlock = usys_file_write,
                .readNumber = usys_file_read_number,
                .writeNumber = usys_file_write_number,
                .eraseBlock = usys_file_erase,
                .writeProtect = usys_file_protect,
                .rename = usys_file_rename,
                .cleanup = usys_file_cleanup,
                .remove = usys_file_remove };

static void free_module_list(void *ip) {
    ListNode *node = (ListNode *)ip;
    if (node) {
        ModuleMap *map = node->data;
        usys_free(map->storeAttr);
        usys_free(map);
        usys_free(node);
    }
}

static int compare_module_list(void *ipt, void *sd) {
    ModuleMap *ip = (ModuleMap *)ipt;
    ModuleMap *op = (ModuleMap *)sd;
    int ret = 0;
    /* If module uuid matches it means module is same.*/
    if (!usys_strcmp(ip->modUuid, op->modUuid)) {
        ret = 1;
    }
    return ret;
}

/* Searching device in the device list*/
static ModuleMap *search_module_in_list(char *puuid) {
    ModuleMap *smap = NULL;
    ModuleMap *fmap = NULL;
    if (!usys_strcmp(puuid, "")) {

        /* If puuid is empty string. This means request for master module
         * This assumption is made because on-bootup ubsp wouldn't know
         * what's the master module uuid. Also it's made sure that first
         * module stored in the list is the master module during the
         * registration process.
         */
        fmap = usys_list_head(&modList, 0);
        if (fmap) {
            usys_log_trace("Map found for master module UUID %s Name %s "
                            "in Module map.",
                            fmap->modUuid, fmap->modName);
        }

    } else {

        smap = usys_zmalloc(sizeof(ModuleMap));
        if (smap) {
            usys_memset(smap->modUuid, '\0', MAX_NAME_LENGTH);
            usys_memcpy(smap->modUuid, puuid, usys_strlen(puuid));

            /*Search return 1 for found.*/
            /* Remember usys_list_search do shallow copy.*/
            fmap = usys_list_search(&modList, smap);
            if (fmap) {
                usys_log_trace("DB::Module UUID: %s found.", smap->modUuid);
            } else {
                usys_free(fmap);
                usys_log_warn("Module UUID: %s not found.", puuid);
            }
            usys_free(smap);
        }

    }
    return fmap;
}

static int update_module_list(ModuleMap *map) {
    int ret = 0;
    if (map) {

        ret = usys_list_update(&modList, map);
        if (ret) {
            usys_log_error(
                            "Module UUID %s and Name %s update "
                            "to module list in store failed.",
                            ret, map->modUuid, map->modName);
        } else {
            usys_log_trace(
                            "Module UUID %s and Name %s update "
                            "to module list in sore completed.",
                            map->modUuid, map->modName);
        }

    }
    return ret;
}

ListInfo *get_module_list() {
    if (modListFlag == 0) {
        usys_list_new(&modList, sizeof(ModuleMap), free_module_list,
                        compare_module_list, NULL);
        modListFlag = 1;
        usys_log_trace("New Module list initialized for store.");
    }
    return &modList;
}

int store_init() {
    int ret = 0;
    usys_list_new(&modList, sizeof(ModuleMap), free_module_list, compare_module_list,
                    NULL);
    modListFlag = 1;
    usys_log_trace("Module list initialized for store.");
    return ret;
}

int store_register_module(UnitCfg *pcfg) {
    int ret = -1;
    ModuleMap *pmap = NULL;
    DevI2cCfg *icfg = NULL;
    char *fname = NULL;
    if (pcfg) {

        /* Check If module is already registered.*/
        pmap = search_module_in_list(pcfg->modUuid);
        if (pmap) {
            usys_log_warn("Map for Module %s is already existing.",
                            pcfg->modUuid);
            ret = 0;
        } else {

            usys_log_debug("No map for Module %s found. Adding now.",
                            pcfg->modUuid);

            pmap = usys_zmalloc(sizeof(ModuleMap));
            if (pmap) {

                usys_memcpy(pmap->modUuid, pcfg->modUuid,
                                usys_strlen(pcfg->modUuid));
                usys_memcpy(pmap->modName, pcfg->modName,
                                usys_strlen(pcfg->modName));

                /* If the eeprom data is exposed as sysfs files */
                if ((pcfg->sysFs != NULL) || (!usys_strcmp(pcfg->sysFs, ""))) {

                    fname = usys_zmalloc(sizeof(char) * MAX_PATH_LENGTH);
                    if (fname) {
                        usys_memcpy(fname, pcfg->sysFs,
                                        usys_strlen(pcfg->sysFs));
                    } else {
                        ret = ERR_NODED_MEMORY_EXHAUSTED;
                        usys_log_error(
                                        "Memory allocation failed. Error: %s",
                                        usys_error(errno));
                        goto cleanup;
                    }

                    pmap->storeAttr = fname;
                    pmap->storeOps = &fileOps;
                } else {

                    /* If the eeprom data has to be read using
                     * user space driver */
                    icfg = usys_zmalloc(sizeof(icfg));
                    if (icfg) {
                        usys_memcpy(icfg, pcfg->eepromCfg, sizeof(DevI2cCfg));
                    } else {
                        ret = ERR_NODED_MEMORY_EXHAUSTED;
                        usys_log_error(
                                        "Memory exhausted while adding eeprom "
                                        "cfg to module map in store. Error: %s",
                                        usys_error(errno));
                        goto cleanup;
                    }
                    pmap->storeAttr = icfg;
                    pmap->storeOps = &eepromOps;

                }

                usys_list_append(&modList, pmap);
                /* Calling Module init.*/
                ret = pmap->storeOps->init(pmap->storeAttr);
                if (ret) {
                    usys_log_debug("Module %s registration success.",
                                    pcfg->modUuid);
                }

            } else {
                ret = ERR_NODED_MEMORY_EXHAUSTED;
                usys_log_error(
                                "Map for Module %s is not added."
                                "Memory exhausted. Error: %s",
                                pcfg->modUuid, usys_error(errno));
            }
        }
        usys_free(pmap);
    }

    if (ret) {
        usys_log_error("Module %s registration failed.", pcfg->modUuid);
    }

    return ret;
    cleanup:
    usys_free(fname);
    usys_free(icfg);
    usys_free(pmap);
    return ret;
}

void store_deregister_all_module() {
    usys_log_debug("Cleaning all entries from module in  module list.");
    usys_list_destroy(&modList);
}

int store_deregister_module(char *puuid) {
    int ret = -1;
    int found = 0;
    if (puuid) {
        /* Check If module is already registered.*/
        ModuleMap *pmap = search_module_in_list(puuid);
        if (pmap) {
            usys_log_warn("Map for Module %s found in module list.", puuid);
            usys_list_remove(&modList, pmap);
            usys_log_debug("Module %s deregistration success.", puuid);
            ret = 0;
            usys_free(pmap);
        } else {
            usys_log_debug("No map for Module %s found.", puuid);
        }
    } else {
        usys_log_warn("UnitCfg is invalid for deregistrating module.");
    }
    return ret;
}

int store_register_update_module(char *puuid, UnitCfg *pcfg, uint8_t count) {
    int ret = -1;
    DevI2cCfg *icfg = NULL;
    char *fname = NULL;
    ModuleMap *pmap = NULL;
    if (pcfg) {
        if (!puuid) {
            return ret;
        }

        /* Check If module is already registered.*/
        pmap = search_module_in_list(puuid);
        if (pmap) {

            usys_log_trace("Map for Module %s found in module list .", puuid);
            usys_memset(pmap->modName, '\0', MAX_NAME_LENGTH);
            usys_memcpy(pmap->modName, pcfg->modName,
                            usys_strlen(pcfg->modName));

            /* Need to free the storeAttr as we don't know new one
             *  could be same or different.*/
            usys_free(pmap->storeAttr);

            /* If the eeprom data is exposed as sysfs files */
            if ((pcfg->sysFs != NULL) || (!usys_strcmp(pcfg->sysFs, ""))) {
                fname = usys_zmalloc(sizeof(char));
                if (fname) {
                    pmap->storeAttr = fname;
                } else {
                    ret = ERR_NODED_MEMORY_EXHAUSTED;
                    usys_log_error(
                                    "Memory exhausted while adding sys file "
                                    "name to store. Error: %s",
                                    usys_error(errno));
                    goto cleanuppmap;
                }
                pmap->storeOps = &fileOps;

            } else {

                /* If the eeprom data needs a custom driver */
                icfg = usys_zmalloc(sizeof(icfg));
                if (icfg) {
                    usys_memcpy(icfg, pcfg->eepromCfg, sizeof(DevI2cCfg));
                } else {
                    ret = ERR_NODED_MEMORY_EXHAUSTED;
                    usys_log_error(
                                    "Memory exhausted while adding eeprom cfg "
                                    "to store. Error: %s",
                                    usys_error(errno));
                    goto cleanuppmap;
                }
                pmap->storeAttr = icfg;
                pmap->storeOps = &eepromOps;

            }

            /* Update module list */
            ret = update_module_list(pmap);
            if (!ret) {
                usys_log_error("Updating module UUID %s to store failed.",
                                ret, puuid);
                goto cleanup;
            } else {
                usys_log_debug("Module %s update success.", pcfg->modUuid);
            }

            /* Rename module name */
            ret = pmap->storeOps->rename(puuid, pcfg->modName);
            if (!ret) {
                usys_log_error("Renaming store for module UUID %s failed.",
                                ret, puuid);
                goto cleanup;
            } else {
                usys_log_debug("Module %s update success.", pcfg->modUuid);
            }

            usys_free(pmap);

        } else {
            usys_log_debug("No map for Module %s found.", puuid);
        }

    } else {
        usys_log_warn("Module UUID is invalid for deregistering module.");
    }
    return ret;

    cleanup:
    usys_free(fname);
    usys_free(icfg);
    cleanuppmap:
    usys_free(pmap);
    return ret;
}

/* Select the data base based on the UUID of the module.*/
ModuleMap *store_choose_module(char *puuid) {
    ModuleMap *pmap = NULL;
    if (puuid) {
        /* Check If module is already registered.*/
        pmap = search_module_in_list(puuid);
        if (pmap) {
            usys_log_trace("Map for Module %s found in module list.", puuid);

        } else {
            usys_log_warn("No map for Module %s found.", puuid);
        }

    } else {
        usys_log_warn("Module UUID is invalid",
                        ERR_NODED_INVALID_POINTER);
    }
    return pmap;
}

int store_rename(char *p_uuid, char *old_name, char *new_name) {
    int ret = -1;
    ModuleMap *pmap = store_choose_module(p_uuid);
    if (pmap) {
        if (pmap->storeOps) {
            ret = pmap->storeOps->rename(old_name, new_name);
        }
    }
    return ret;
}

int store_read_block(char *p_uuid, void *data, off_t offset, uint16_t size) {
    int ret = -1;
    ModuleMap *pmap = store_choose_module(p_uuid);
    if (pmap) {
        if (pmap->storeOps) {
            ret = pmap->storeOps->readBlock(pmap->storeAttr, data, offset, size);
        }
        usys_free(pmap);
    }
    return ret;
}

int store_write_block(char *p_uuid, void *data, off_t offset, uint16_t size) {
    int ret = -1;
    ModuleMap *pmap = store_choose_module(p_uuid);
    if (pmap) {
        if (pmap->storeOps) {
            ret = pmap->storeOps->writeBlock(pmap->storeAttr, data, offset, size);
        }
        usys_free(pmap);
    }
    return ret;
}

int store_erase_block(char *p_uuid, off_t offset, uint16_t size) {
    int ret = -1;
    ModuleMap *pmap = store_choose_module(p_uuid);
    if (pmap) {
        if (pmap->storeOps) {
            ret = pmap->storeOps->eraseBlock(pmap->storeAttr, offset, size);
        }
        usys_free(pmap);
    }
    return ret;
}

int store_read_number(char *p_uuid, void *val, off_t offset, uint16_t count,
                uint8_t size) {
    int ret = -1;
    ModuleMap *pmap = store_choose_module(p_uuid);
    if (pmap) {
        if (pmap->storeOps) {
            ret = pmap->storeOps->readNumber(pmap->storeAttr, val, offset, count,
                            size);
        }
        usys_free(pmap);
    }
    return ret;
}

int store_write_number(char *p_uuid, void *val, off_t offset, uint16_t count,
                uint8_t size) {
    int ret = -1;
    ModuleMap *pmap = store_choose_module(p_uuid);
    if (pmap) {
        if (pmap->storeOps) {
            ret = pmap->storeOps->writeNumber(pmap->storeAttr, val, offset, count,
                            size);
        }
        usys_free(pmap);
    }
    return ret;
}

int store_write_protect(char *p_uuid, void *data) {
    int ret = -1;
    ModuleMap *pmap = store_choose_module(p_uuid);
    if (pmap) {
        if (pmap->storeOps) {
            ret = pmap->storeOps->writeProtect(pmap->storeAttr);
        }
        usys_free(pmap);
    }
    return ret;
}

int store_remove(char *p_uuid) {
    int ret = -1;
    ModuleMap *pmap = store_choose_module(p_uuid);
    if (pmap) {
        if (pmap->storeOps) {
            ret = pmap->storeOps->remove(pmap->storeAttr);
        }
        usys_free(pmap);
    }
    return ret;
}

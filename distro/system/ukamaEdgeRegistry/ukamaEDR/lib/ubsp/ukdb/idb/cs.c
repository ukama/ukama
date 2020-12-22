/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "ukdb/idb/cs.h"

#include "headers/ubsp/devices.h"
#include "headers/errorcode.h"
#include "headers/ubsp/ukdblayout.h"
#include "inc/globalheader.h"
#include "headers/utils/log.h"

#include "mfgdata/cstructs/mfgdata.h"
#include "mfgdata/cstructs/mfgdata1.inc"
#include "mfgdata/cstructs/mfgdata2.inc"
#include "mfgdata/cstructs/mfgdata3.inc"
#include "mfgdata/cstructs/mfgdata4.inc"
#include "mfgdata/cstructs/mfgdata5.inc"

/* g_mukdb would be replaced by some json parser with mfg data.*/
/* Allmost all global variables will be replaces with json provided structs.*/

#define MAX_MODULE_COUNT 5

UKDB g_mukdb;
uint8_t g_idx_count = 0;
MFGData *g_mfgdata;

static void print_mfg_data(char *info, uint16_t size) {
    uint16_t itr = 0;
    while (itr < size) {
        putchar(*(info + itr));
        itr++;
    }
    fflush(stdout);
}
//TODO: Will be replaced by appropriate fxn after parser.*/
/* Fetch the manufacturing data from parser or C structs.*/

static int get_mfgdata_by_uuid(char *p_uuid) {
    int ret = -1;
    /* Assumption when json parser parse the file it always has master file at index 0
	 * or has to tell the index where it is so default section can be updated.*/
    g_mfgdata = (MFGData[]){
        { .uuid = "UK1001-COMV1", .idx_count = 10, .ukdb = ukdb1 },
        { .uuid = "UK2001-LTE", .idx_count = 8, .ukdb = ukdb2 },
        { .uuid = "UK3001-MASK", .idx_count = 8, .ukdb = ukdb3 },
        { .uuid = "UK4001-RFFE", .idx_count = 8, .ukdb = ukdb4 },
        { .uuid = "UK5001-RFCTRL", .idx_count = 10, .ukdb = ukdb5 },
    };

    /*Default section*/
    if ((!p_uuid) || !strcmp(p_uuid, "")) {
        memcpy(&g_mukdb, &g_mfgdata[0].ukdb, sizeof(UKDB));
        g_idx_count = g_mfgdata[0].idx_count;
        log_trace(
            "CS:: MFG Data set to the mfg_data[%d], Module UUID %s with %d entries in Index Table.",
            0, g_mfgdata[0].uuid, g_idx_count);
        ret = 0;
    } else {
        /* Searching for Module MFG data index.*/
        for (uint8_t iter = 0; iter < MAX_MODULE_COUNT; iter++) {
            if (!strcmp(p_uuid, g_mfgdata[iter].uuid)) {
                memcpy(&g_mukdb, &g_mfgdata[iter].ukdb, sizeof(UKDB));
                g_idx_count = g_mfgdata[iter].idx_count;
                log_trace(
                    "CS:: MFG Data set to the mfg_data[%d], Module UUID %s with %d entries in Index Table.",
                    iter, g_mfgdata[iter].uuid, g_idx_count);
                ret = 0;
                break;
            }
        }
    }
    if (ret) {
        ret = ERR_UBSP_DB_MISSING_MODULE;
        log_error("Err(%d) CS:: No module with UUID %s found in MFG Data.",
                  p_uuid);
    }
    return ret;
}

int cs_init(void *data) {
    int ret = 0;
    get_mfgdata_by_uuid("");
    return ret;
}

int cs_parse(void *data) {
    int ret = 0;
    //TODO: Parse json data from MFG DB. and store it in g_ukdb.
    return ret;
}

/* Reads the Header info from the MFG data.*/
void *cs_fetch_header(char *p_uuid, uint16_t *size) {
    int ret = 0;
    ret = get_mfgdata_by_uuid(p_uuid);
    if (ret) {
        return NULL;
    }

    uint16_t sz = sizeof(UKDBHeader);
    UKDBHeader *header = malloc(sz);
    if (header) {
        memcpy(header, &g_mukdb.header, sz);
        *size = sz;
        log_debug(
            "CS:: Reading IDB Header with %d bytes from Mfg data completed.",
            *size);
    } else {
        ret = -1;
        log_debug("CS:: Reading IDB Header with %d bytes from Mfg data failed.",
                  *size);
    }
    return header;
}

/*Reads the Index Table from MFG data*/
void *cs_fetch_index(char *p_uuid, uint16_t *size) {
    int ret = 0;
    ret = get_mfgdata_by_uuid(p_uuid);
    if (ret) {
        return NULL;
    }

    int sz = sizeof(UKDBIdxTuple) * g_idx_count;
    UKDBIdxTuple *idx = malloc(sz);
    if (idx) {
        memcpy(idx, &g_mukdb.indextable[0], sz);
        *size = sz;
        log_debug(
            "CS:: Reading IDB Index with %d bytes from Mfg data completed.",
            *size);
    } else {
        ret = -1;
        log_debug("CS:: Reading IDB Index with %d bytes from Mfg data failed.",
                  *size);
    }
    return idx;
}

void *cs_fetch_unit_info(char *p_uuid, uint16_t *size) {
    int ret = 0;
    ret = get_mfgdata_by_uuid(p_uuid);
    if (ret) {
        return NULL;
    }

    uint16_t sz = sizeof(UnitInfo);
    UnitInfo *info = malloc(sz);
    if (info) {
        memcpy(info, &g_mukdb.unitinfo, sz);
        *size = sz;
        log_debug(
            "CS:: Reading IDB Unit info with %d bytes from Mfg data completed.",
            *size);
    } else {
        ret = -1;
        log_debug(
            "CS:: Reading IDB Unit info with %d bytes from Mfg data failed.",
            *size);
    }
    return info;
}

void *cs_fetch_unit_cfg(char *p_uuid, uint16_t *size, uint8_t count) {
    int ret = 0;
    ret = get_mfgdata_by_uuid(p_uuid);
    if (ret) {
        return NULL;
    }

    int sz = sizeof(UnitCfg) * count;
    UnitCfg *cfg = malloc(sz);
    if (cfg) {
        memcpy(cfg, &g_mukdb.unitcfg[0], sz);
        /* Overwrite EEPROM CFG for each Module */
        for (int iter = 0; iter < count; iter++) {
            DevI2cCfg *icfg = malloc(sizeof(DevI2cCfg));
            if (icfg) {
                memcpy(icfg, g_mukdb.unitcfg[iter].eeprom_cfg,
                       sizeof(DevI2cCfg));
            }
            cfg[iter].eeprom_cfg = icfg;
        }
        *size = sz;
        log_debug(
            "CS:: Reading IDB Unit Config with %d bytes from Mfg data completed.",
            *size);
    } else {
        ret = -1;
        log_debug(
            "CS:: Reading IDB Unit Config with %d bytes from Mfg data failed.",
            *size);
    }
    return cfg;
}
//TODO: Check Module info and config fetch gains as .inc line 374 comment.
// Each module should have separate .inc and for every module we should call it's ukdb_create with uuid.
/* Reads module info and module cfg form db based on the location in file.*/
//TODO: add prototypes in header files and in fxn ptr table of idb.*/
void *cs_fetch_module_info(char *p_uuid, uint16_t *size, uint8_t idx) {
    int ret = 0;
    ModuleInfo *info = NULL;
    ret = get_mfgdata_by_uuid(p_uuid);
    if (ret) {
        goto cleanup;
    }

    /* Fetching Module Info from MFG data*/
    uint16_t sz = sizeof(ModuleInfo);
    info = malloc(sz);
    if (info) {
        memcpy(info, &g_mukdb.modinfo, sizeof(ModuleInfo));
        if (!(VALIDATE_DEVICE_COUNT(info->dev_count))) {
            ret = ERR_UBSP_VALIDATION_FAILURE;
            log_error("Err(%d): CS:: validation for device count (%d) failed.",
                      ret, info->dev_count);
            goto cleanup;
        }
        /* Fetching Module Config from MFG data*/
        uint16_t csz = sizeof(ModuleCfg) * (info->dev_count);
        ModuleCfg *cfg = (ModuleCfg *)malloc(csz);
        if (cfg) {
            memcpy(cfg, g_mukdb.modcfg, csz);
            info->module_cfg = cfg;
            *size = sz;
            log_debug(
                "CS:: Reading IDB Module info with %d bytes from Mfg data completed.",
                *size);
        } else {
            ret = ERR_UBSP_MEMORY_EXHAUSTED;
            log_error(
                "Err(%d): IDB Memory exhausted while reading module cfg of %d bytes.",
                ret, csz);
            goto cleanup;
        }
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_error(
            "Err(%d): IDB Memory exhausted while reading module info of %d bytes.",
            ret, sz);
    }

cleanup:
    if (ret) {
        UBSP_FREE(info);
    }
    return info;
}

void *cs_fetch_module_device_cfg(DeviceClass class, void *hw_attr) {
    uint16_t size = 0;
    SIZE_OF_DEVICE_CFG(size, class);
    void *dev_cfg = malloc(size);
    if (dev_cfg) {
        memset(dev_cfg, '\0', size);
        if (hw_attr) {
            memcpy(dev_cfg, hw_attr, size);
        }
    } else {
        dev_cfg = NULL;
    }
    return dev_cfg;
}

/* Fetch module info by uuid from higher layer which has to also tell number of modules present in unit so that
 * it can go thorugh all of them and compare uuid.*/
void *cs_fetch_module_info_by_uuid(char *p_uuid, uint16_t *size,
                                   uint8_t count) {
    int ret = -1;
    ret = get_mfgdata_by_uuid(p_uuid);
    if (ret) {
        return NULL;
    }

    uint8_t iter = 0;
    ModuleInfo *info;
    int res = -1;
    /* This is done for case for unit module info.*/
    if (strcmp(p_uuid, "") == 0) {
        res = 0;
    } else {
        res = strcmp(g_mukdb.modinfo.uuid, p_uuid);
    }
    if (!res) {
        /* Fetching Module Info from MFG data*/
        uint16_t sz = sizeof(ModuleInfo);
        info = malloc(sz);
        if (info) {
            memset(info, '\0', sz);
            memcpy(info, &g_mukdb.modinfo, sz);

            /* Fetching Module Config from MFG data*/
            uint16_t csz = sizeof(ModuleCfg) * (info->dev_count);
            ModuleCfg *cfg = (ModuleCfg *)malloc(csz);
            if (cfg) {
                memset(cfg, '\0', csz);
                memcpy(cfg, g_mukdb.modcfg, csz);
                for (int iter = 0; iter < info->dev_count; iter++) {
                    /*Fetching device cfg for each device in a module*/
                    cfg[iter].cfg = cs_fetch_module_device_cfg(
                        cfg[iter].dev_class, g_mukdb.modcfg[iter].cfg);
                }
                info->module_cfg = cfg;
                log_debug(
                    "CS:: Reading IDB Module Config for UUID %s with %d bytes from Mfg data completed.",
                    g_mukdb.modinfo.uuid, csz);
            } else {
                ret = ERR_UBSP_MEMORY_EXHAUSTED;
                log_error(
                    "Err(%d): IDB Memory exhausted while reading Module cfg for UIID %s with %d bytes.",
                    ret, g_mukdb.modinfo.uuid, csz);
                if (info) {
                    free(info);
                    info = NULL;
                }
            }
        } else {
            ret = ERR_UBSP_MEMORY_EXHAUSTED;
            log_error(
                "Err(%d): IDB Memory exhausted while reading module info of %d bytes.",
                ret, sz);
        }
        *size = sz;
        log_debug(
            "CS:: Reading IDB Module info for UUID %s with %d bytes for uuid %s from Mfg data completed.",
            g_mukdb.modinfo.uuid, *size, p_uuid);
    }
    return info;
}

/* Reads module config only */
void *cs_fetch_module_cfg(char *p_uuid, uint16_t *size, uint8_t count) {
    int ret = -1;
    ret = get_mfgdata_by_uuid(p_uuid);
    if (ret) {
        return NULL;
    }

    uint16_t sz = 0;
    uint8_t iter = 0;
    ModuleCfg *mcfg;
    for (iter = 0; iter <= count; iter++) {
        int res = strcmp(g_mukdb.modinfo.uuid, p_uuid);
        if (!res) {
            /* Fetching Module Config from MFG data*/
            uint16_t csz = sizeof(ModuleCfg) * (g_mukdb.modinfo.dev_count);
            mcfg = (ModuleCfg *)malloc(csz);
            if (mcfg) {
                memcpy(mcfg, g_mukdb.modcfg, csz);
                *size = sz;
                log_debug(
                    "CS:: Reading IDB Module config with %d bytes from Mfg data completed.",
                    *size);
            } else {
                ret = ERR_UBSP_MEMORY_EXHAUSTED;
                log_error(
                    "Err(%d): IDB Memory exhausted while reading module cfg of %d bytes.",
                    ret, csz);
            }
        }
        log_debug(
            "CS:: Reading IDB Module Config with %d bytes for uuid %s from Mfg data completed.",
            *size, p_uuid);
        break;
    }
    return mcfg;
}

void *cs_fetch_fact_config(char *p_uuid, uint16_t *size) {
    int ret = 0;
    ret = get_mfgdata_by_uuid(p_uuid);
    if (ret) {
        return NULL;
    }

    char *info = malloc(sizeof(char) * (*size));
    if (info) {
        memcpy(info, g_mukdb.factcfg, sizeof(char) * (*size));
        /* *size=strlen(info); */ /*Not required length will be specified in the index*/
        log_debug("CS:: Reading IDB fact Config with %d bytes.\nData::", *size);
        print_mfg_data(info, *size);
    }
    return info;
}

void *cs_fetch_user_config(char *p_uuid, uint16_t *size) {
    int ret = 0;
    ret = get_mfgdata_by_uuid(p_uuid);
    if (ret) {
        return NULL;
    }

    char *info = malloc(sizeof(char) * (*size));
    if (info) {
        memcpy(info, g_mukdb.usercfg, sizeof(char) * (*size));
        /* *size=strlen(info); */ /*Not required length will be specified in the index*/
        log_debug("CS:: Reading IDB User Config with %d bytes.\nData::", *size);
        print_mfg_data(info, *size);
    }
    return info;
}

void *cs_fetch_fact_calib(char *p_uuid, uint16_t *size) {
    int ret = 0;
    ret = get_mfgdata_by_uuid(p_uuid);
    if (ret) {
        return NULL;
    }

    char *info = malloc(sizeof(char) * (*size));
    if (info) {
        memcpy(info, g_mukdb.factcalib, sizeof(char) * (*size));
        /* *size=strlen(info); */ /*Not required length will be specified in the index*/
        log_debug("CS:: Reading IDB fact Calib with %d bytes.\nData::", *size);
        print_mfg_data(info, *size);
    }
    return info;
}

void *cs_fetch_user_calib(char *p_uuid, uint16_t *size) {
    int ret = 0;
    ret = get_mfgdata_by_uuid(p_uuid);
    if (ret) {
        return NULL;
    }

    char *info = malloc(sizeof(char) * (*size));
    if (info) {
        memcpy(info, g_mukdb.usercalib, sizeof(char) * (*size));
        /* *size=strlen(info); */ /*Not required length will be specified in the index*/
        log_debug("CS:: Reading IDB User calib with %d bytes.\\nData::", *size);
        print_mfg_data(info, *size);
    }
    return info;
    ;
}

void *cs_fetch_bs_certs(char *p_uuid, uint16_t *size) {
    int ret = 0;
    ret = get_mfgdata_by_uuid(p_uuid);
    if (ret) {
        return NULL;
    }

    char *info = malloc(sizeof(char) * (*size));
    if (info) {
        memcpy(info, g_mukdb.bscerts, sizeof(char) * (*size));
        /* *size=strlen(info); */ /*Not required length will be specified in the index*/
        log_debug("CS:: Reading IDB BS certs with %d bytes.\nData::", *size);
        print_mfg_data(info, *size);
    }
    return info;
}

void *cs_fetch_lwm2m_certs(char *p_uuid, uint16_t *size) {
    int ret = 0;
    ret = get_mfgdata_by_uuid(p_uuid);
    if (ret) {
        return NULL;
    }

    char *info = malloc(sizeof(char) * (*size));
    if (info) {
        memcpy(info, g_mukdb.lwm2mcerts, sizeof(char) * (*size));
        /* *size=strlen(info); */ /*Not required length will be specified in the index*/
        log_debug("CS:: Reading IDB lwm2m certs with %d bytes. \nData::",
                  *size);
        print_mfg_data(info, *size);
    }
    return info;
}

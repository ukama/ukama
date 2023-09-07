/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "ukdb/idb/jp.h"

#include "headers/ubsp/devices.h"
#include "headers/errorcode.h"
#include "headers/ubsp/ukdblayout.h"
#include "inc/globalheader.h"
#include "ukdb/idb/sparser.h"
#include "headers/utils/log.h"

static void print_mfg_data(char *info, uint16_t size) {
    uint16_t itr = 0;
    while (itr < size) {
        putchar(*(info + itr));
        itr++;
    }
    fflush(stdout);
}

int jp_init(void *ip) {
    int ret = 0;
    ret = parser_schema_init(ip);
    if (ret) {
        log_debug("Err(%d): JP:: Parsing failed.", ret);
    }
    return ret;
}

void jp_exit() {
    parser_schema_exit();
}

void *jp_fetch_header(char *puuid, uint16_t *size) {
    int ret = 0;
    UKDBHeader *header = NULL;
    UKDB *db = parser_get_mfg_data_by_uuid(puuid);
    if (!db) {
        log_error("Err: JP :: No UKDB found for %s module.", puuid);
        return NULL;
    }

    uint16_t sz = sizeof(UKDBHeader);
    header = malloc(sz);
    if (header) {
        memcpy(header, &db->header, sz);
        *size = sz;
        log_debug(
            "PARSER:: Reading IDB Header with %d bytes for Module UUID %s.",
            *size, puuid);
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug(
            "ERR(%d): PARSER:: Reading IDB Header with %d bytes for Module UUID %s.",
            ret, *size, puuid);
    }
    return header;
}

void *jp_fetch_index(char *puuid, uint16_t *size) {
    int ret = 0;
    UKDB *db = parser_get_mfg_data_by_uuid(puuid);
    if (!db) {
        log_error("Err: JP :: No UKDB found for %s module.", puuid);
        return NULL;
    }

    int sz = sizeof(UKDBIdxTuple) * db->header.idx_cur_tpl;
    UKDBIdxTuple *idx = malloc(sz);
    if (idx) {
        memcpy(idx, &db->indextable[0], sz);
        *size = sz;
        log_debug(
            "JP:: Reading IDB Index table with %d bytes for Module UUID %s.",
            *size, puuid);
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug(
            "ERR(%d): JP:: Reading IDB Index table with %d bytes for Module UUID %s.",
            ret, *size, puuid);
    }
    return idx;
}

void *jp_fetch_unit_info(char *puuid, uint16_t *size) {
    int ret = 0;
    UKDB *db = parser_get_mfg_data_by_uuid(puuid);
    if (!db) {
        log_error("Err: JP :: No UKDB found for %s module.", puuid);
        return NULL;
    }

    uint16_t sz = sizeof(UnitInfo);
    UnitInfo *info = malloc(sz);
    if (info) {
        memcpy(info, &db->unitinfo, sz);
        *size = sz;
        log_debug(
            "JP:: Reading IDB Unit Info with %d bytes for Module UUID %s.",
            *size, puuid);
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug(
            "ERR(%d): JP:: Reading IDB Unit Info with %d bytes for Module UUID %s.",
            ret, *size, puuid);
    }
    return info;
}

void *jp_fetch_unit_cfg(char *puuid, uint16_t *size, uint8_t count) {
    int ret = 0;
    UKDB *db = parser_get_mfg_data_by_uuid(puuid);
    if (!db) {
        log_error("Err: JP :: No UKDB found for %s module.", puuid);
        return NULL;
    }

    int sz = sizeof(UnitCfg) * count;
    UnitCfg *cfg = malloc(sz);
    if (cfg) {
        memcpy(cfg, &db->unitcfg[0], sz);
        /* Overwrite EEPROM CFG for each Module */
        for (int iter = 0; iter < count; iter++) {
            DevI2cCfg *icfg = malloc(sizeof(DevI2cCfg));
            if (icfg) {
                memcpy(icfg, db->unitcfg[iter].eeprom_cfg, sizeof(DevI2cCfg));
            }
            cfg[iter].eeprom_cfg = icfg;
        }
        *size = sz;
        log_debug(
            "JP:: Reading IDB Unit Config with %d bytes for Module UUID %s.",
            *size, puuid);
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug(
            "ERR(%d): JP:: Reading IDB Unit Config with %d bytes for Module UUID %s.",
            ret, *size, puuid);
    }
    return cfg;
}

void *jp_fetch_module_device_cfg(DeviceClass class, void *hw_attr) {
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

void *jp_fetch_module_info_by_uuid(char *puuid, uint16_t *size, uint8_t count) {
    int ret = 0;
    UKDB *db = parser_get_mfg_data_by_uuid(puuid);
    if (!db) {
        log_error("Err: JP :: No UKDB found for %s module.", puuid);
        return NULL;
    }

    uint8_t iter = 0;
    ModuleInfo *info;
    int res = -1;
    /* This is done for case for unit module info.*/
    if (strcmp(puuid, "") == 0) {
        res = 0;
    } else {
        res = strcmp(db->modinfo.uuid, puuid);
    }
    if (!res) {
        /* Fetching Module Info from MFG data*/
        uint16_t sz = sizeof(ModuleInfo);
        info = malloc(sz);
        if (info) {
            memset(info, '\0', sz);
            memcpy(info, &db->modinfo, sz);

            /* Fetching Module Config from MFG data*/
            uint16_t csz = sizeof(ModuleCfg) * (info->dev_count);
            ModuleCfg *cfg = (ModuleCfg *)malloc(csz);
            if (cfg) {
                memset(cfg, '\0', csz);
                memcpy(cfg, db->modcfg, csz);
                for (int iter = 0; iter < info->dev_count; iter++) {
                    /*Fetching device cfg for each device in a module*/
                    cfg[iter].cfg = jp_fetch_module_device_cfg(
                        cfg[iter].dev_class, db->modcfg[iter].cfg);
                }
                info->module_cfg = cfg;
                log_debug(
                    "JP:: Reading IDB Module Config with %d bytes for Module UUID %s.",
                    *size, puuid);
            } else {
                ret = ERR_UBSP_MEMORY_EXHAUSTED;
                log_debug(
                    "ERR(%d): JP:: Reading IDB Module Config with %d bytes for Module UUID %s.",
                    ret, *size, puuid);

                UBSP_FREE(info);
            }
        } else {
            ret = ERR_UBSP_MEMORY_EXHAUSTED;
            log_debug(
                "ERR(%d): JP:: Reading IDB Module Info with %d bytes for Module UUID %s.",
                ret, *size, puuid);
        }
        *size = sz;
        log_debug(
            "JP:: Reading IDB Module info with %d bytes for Module UUID %s.",
            *size, puuid);
    }
    return info;
}

/* Reads module config only */
void *jp_fetch_module_cfg(char *puuid, uint16_t *size, uint8_t count) {
    int ret = 0;
    UKDB *db = parser_get_mfg_data_by_uuid(puuid);
    if (!db) {
        log_error("Err: JP :: No UKDB found for %s module.", puuid);
        return NULL;
    }

    uint16_t sz = 0;
    uint8_t iter = 0;
    ModuleCfg *mcfg;
    for (iter = 0; iter <= count; iter++) {
        int res = strcmp(db->modinfo.uuid, puuid);
        if (!res) {
            /* Fetching Module Config from MFG data*/
            uint16_t csz = sizeof(ModuleCfg) * (db->modinfo.dev_count);
            mcfg = (ModuleCfg *)malloc(csz);
            if (mcfg) {
                memcpy(mcfg, db->modcfg, csz);
                *size = sz;
                log_debug(
                    "JP:: Reading IDB Module config with %d bytes for Module UUID %s.",
                    *size, puuid);
            } else {
                ret = ERR_UBSP_MEMORY_EXHAUSTED;
                log_error(
                    "Err(%d): JP:: Reading IDB Module config with %d bytes for Module UUID %s.",
                    ret, *size, puuid);
            }
        }
        log_debug(
            "JP:: Reading IDB Module Config with %d bytes for uuid %s from Mfg data completed.",
            *size, puuid);
        break;
    }
    return mcfg;
}

void *jp_fetch_fact_config(char *puuid, uint16_t *size) {
    int ret = 0;
    UKDB *db = parser_get_mfg_data_by_uuid(puuid);
    if (!db) {
        log_error("Err: JP :: No UKDB found for %s module.", puuid);
        return NULL;
    }

    char *info = malloc(sizeof(char) * (*size));
    if (info) {
        memcpy(info, db->factcfg, sizeof(char) * (*size));
        log_debug(
            "JP:: Reading IDB fact Config with %d bytes for Module UUID %s.\nData::",
            *size, puuid);
        print_mfg_data(info, *size);
    }
    return info;
}

void *jp_fetch_user_config(char *puuid, uint16_t *size) {
    int ret = 0;
    UKDB *db = parser_get_mfg_data_by_uuid(puuid);
    if (!db) {
        log_error("Err: JP :: No UKDB found for %s module.", puuid);
        return NULL;
    }

    char *info = malloc(sizeof(char) * (*size));
    if (info) {
        memcpy(info, db->usercfg, sizeof(char) * (*size));
        log_debug(
            "JP:: Reading IDB User Config with %d bytes for Module UUID %s.\nData::",
            *size, puuid);
        print_mfg_data(info, *size);
    }
    return info;
}

void *jp_fetch_fact_calib(char *puuid, uint16_t *size) {
    int ret = 0;
    UKDB *db = parser_get_mfg_data_by_uuid(puuid);
    if (!db) {
        log_error("Err: JP :: No UKDB found for %s module.", puuid);
        return NULL;
    }

    char *info = malloc(sizeof(char) * (*size));
    if (info) {
        memcpy(info, db->factcalib, sizeof(char) * (*size));
        log_debug(
            "JP:: Reading IDB factory calibration with %d bytes for Module UUID %s.\nData::",
            *size, puuid);
        print_mfg_data(info, *size);
    }
    return info;
}

void *jp_fetch_user_calib(char *puuid, uint16_t *size) {
    int ret = 0;
    UKDB *db = parser_get_mfg_data_by_uuid(puuid);
    if (!db) {
        log_error("Err: JP :: No UKDB found for %s module.", puuid);
        return NULL;
    }

    char *info = malloc(sizeof(char) * (*size));
    if (info) {
        memcpy(info, db->usercalib, sizeof(char) * (*size));
        log_debug(
            "JP:: Reading IDB user calibration with %d bytes for Module UUID %s.\nData::",
            *size, puuid);
        print_mfg_data(info, *size);
    }
    return info;
}

void *jp_fetch_bs_certs(char *puuid, uint16_t *size) {
    int ret = 0;
    UKDB *db = parser_get_mfg_data_by_uuid(puuid);
    if (!db) {
        log_error("Err: JP :: No UKDB found for %s module.", puuid);
        return NULL;
    }

    char *info = malloc(sizeof(char) * (*size));
    if (info) {
        memcpy(info, db->bscerts, sizeof(char) * (*size));
        log_debug(
            "JP:: Reading IDB Bootstrap certificates with %d bytes for Module UUID %s.\nData::",
            *size, puuid);
        print_mfg_data(info, *size);
    }
    return info;
}

void *jp_fetch_lwm2m_certs(char *puuid, uint16_t *size) {
    int ret = 0;
    UKDB *db = parser_get_mfg_data_by_uuid(puuid);
    if (!db) {
        log_error("Err: JP :: No UKDB found for %s module.", puuid);
        return NULL;
    }

    char *info = malloc(sizeof(char) * (*size));
    if (info) {
        memcpy(info, db->lwm2mcerts, sizeof(char) * (*size));
        log_debug(
            "JP:: Reading IDB Lwm2m certificates with %d bytes for Module UUID %s.\nData::",
            *size, puuid);
        print_mfg_data(info, *size);
    }
    return info;
}

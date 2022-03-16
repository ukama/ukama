/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "jdata.h"

#include "usys_api.h"
#include "usys_error.h"
#include "usys_mem.h"
#include "usys_string.h"

static void print_mfg_data(char *info, uint16_t size) {
    uint16_t itr = 0;
    while (itr < size) {
        usys_putchar(*(info + itr));
        itr++;
    }
    usys_fflush(stdout);
}

int jdata_init(void *ip) {
    int ret = 0;
    ret = parser_schema_init(ip);
    if (ret) {
        usys_log_error("Parsing failed. Error: %d", ret);
    }
    return ret;
}

void jdata_exit() {
    parser_schema_exit();
}

void *jdata_fetch_header(char *puuid, uint16_t *size) {
    int ret = 0;
    SchemaHeader *header = NULL;
    StoreSchema *sschema = parser_get_mfg_data_by_uuid(puuid);
    if (!sschema) {
        usys_log_error("No Mfg data found for %s module.", puuid);
        return NULL;
    }

    uint16_t sz = sizeof(SchemaHeader);
    header = usys_zmalloc(sz);
    if (header) {

        usys_memcpy(header, &sschema->header, sz);
        *size = sz;
        usys_log_debug(
            "Reading Mfg data Header with %d bytes for Module UUID %s.",
            *size, puuid);

    } else {

        ret = ERR_NODED_MEMORY_EXHAUSTED;
        usys_log_error(
            "Reading Mfg Data Header with %d bytes for Module UUID %s failed."
            "Error: %s",
            *size, puuid, usys_error(errno));

    }
    return header;
}

void *jdata_fetch_idx(char *puuid, uint16_t *size) {
    int ret = 0;
    StoreSchema *sschema = parser_get_mfg_data_by_uuid(puuid);
    if (!sschema) {
        usys_log_error("No Mfg data found for %s module.", puuid);
        return NULL;
    }

    int sz = sizeof(SchemaIdxTuple) * sschema->header.idxCurTpl;
    SchemaIdxTuple *idx = usys_zmalloc(sz);
    if (idx) {

        usys_memcpy(idx, &sschema->indexTable[0], sz);
        *size = sz;
        usys_log_debug(
            "Reading Mfg data Index table with %d bytes for Module UUID %s.",
            *size, puuid);

    } else {

        ret = ERR_NODED_MEMORY_EXHAUSTED;
        usys_log_error(
            "Reading Mfg data Index table with %d bytes for Module UUID %s "
            "failed.Error %s",
            *size, puuid, usys_error(errno));

    }
    return idx;
}

void *jdata_fetch_unit_info(char *puuid, uint16_t *size) {
    int ret = 0;
    StoreSchema *sschema = parser_get_mfg_data_by_uuid(puuid);
    if (!sschema) {
        usys_log_error("No StoreSchema found for %s module.", puuid);
        return NULL;
    }

    uint16_t sz = sizeof(UnitInfo);
    UnitInfo *info = usys_zmalloc(sz);
    if (info) {

        usys_memcpy(info, &sschema->unitInfo, sz);
        *size = sz;
        usys_log_debug(
            "Reading Mfg data Unit Info with %d bytes for Module UUID %s.",
            *size, puuid);

    } else {
        ret = ERR_NODED_MEMORY_EXHAUSTED;
        usys_log_error(
            "Reading Mfg data Unit Info with %d bytes for Module UUID %s "
            "failed.Error: %s",
            ret, *size, puuid, usys_error(errno));
    }
    return info;
}

void *jdata_fetch_unit_cfg(char *puuid, uint16_t *size, uint8_t count) {
    int ret = 0;
    StoreSchema *sschema = parser_get_mfg_data_by_uuid(puuid);
    if (!sschema) {
        usys_log_error("No StoreSchema found for %s module.", puuid);
        return NULL;
    }

    int sz = sizeof(UnitCfg) * count;
    UnitCfg *cfg = usys_zmalloc(sz);
    if (cfg) {
        usys_memcpy(cfg, &sschema->unitCfg[0], sz);
        /* Overwrite EEPROM CFG for each Module */
        for (int iter = 0; iter < count; iter++) {
            DevI2cCfg *iCfg = usys_zmalloc(sizeof(DevI2cCfg));
            if (iCfg) {
                usys_memcpy(iCfg, sschema->unitCfg[iter].eepromCfg,
                                sizeof(DevI2cCfg));
            }
            cfg[iter].eepromCfg = iCfg;
        }
        *size = sz;
        usys_log_debug(
            "Reading Mfg data Unit Config with %d bytes for Module UUID %s.",
            *size, puuid);
    } else {
        ret = ERR_NODED_MEMORY_EXHAUSTED;
        usys_log_error(
            "Reading Mfg data Unit Config with %d bytes for Module UUID %s "
            "failed.Error: %s",
            *size, puuid, usys_error(errno));
    }
    return cfg;
}

void *jdata_fetch_module_device_cfg(DeviceClass class, void *hwAttr) {
    uint16_t size = 0;
    SIZE_OF_DEVICE_CFG(size, class);
    void *devCfg = usys_zmalloc(size);
    if (devCfg) {
        if (hwAttr) {
            usys_memcpy(devCfg, hwAttr, size);
        }
    } else {
        devCfg = NULL;
    }
    return devCfg;
}

void *jdata_fetch_module_info_by_uuid(char *puuid, uint16_t *size,
                uint8_t count) {
    int ret = 0;
    StoreSchema *sschema = parser_get_mfg_data_by_uuid(puuid);
    if (!sschema) {
        usys_log_error("No StoreSchema found for %s module.", puuid);
        return NULL;
    }

    uint8_t iter = 0;
    ModuleInfo *info;
    int res = -1;
    /* This is done for case for unit module info.*/
    if (usys_strcmp(puuid, "") == 0) {
        res = 0;
    } else {
        res = usys_strcmp(sschema->modInfo.uuid, puuid);
    }
    if (!res) {
        /* Fetching Module Info from MFG data*/
        uint16_t sz = sizeof(ModuleInfo);

        info = usys_zmalloc(sz);
        if (info) {
            usys_memset(info, '\0', sz);
            usys_memcpy(info, &sschema->modInfo, sz);

            /* Fetching Module Config from MFG data*/
            uint16_t csz = sizeof(ModuleCfg) * (info->devCount);
            ModuleCfg *cfg = (ModuleCfg *)usys_zmalloc(csz);
            if (cfg) {

                usys_memset(cfg, '\0', csz);
                usys_memcpy(cfg, sschema->modCfg, csz);

                for (int iter = 0; iter < info->devCount; iter++) {
                    /*Fetching device cfg for each device in a module*/
                    cfg[iter].cfg = jdata_fetch_module_device_cfg(
                        cfg[iter].devClass, sschema->modCfg[iter].cfg);
                }

                info->modCfg = cfg;
                usys_log_debug(
                    "Reading Mfg data Module Config with %d bytes for Module "
                    "UUID %s.",
                    *size, puuid);

            } else {

                ret = ERR_NODED_MEMORY_EXHAUSTED;
                usys_log_debug(
                    "Reading Mfg data Module Config with %d bytes for Module "
                    "UUID %s failed. Error: %s",
                    *size, puuid, usys_error(errno));

                UBSP_FREE(info);
            }
        } else {
            ret = ERR_NODED_MEMORY_EXHAUSTED;
            usys_log_debug(
                "Reading Mfg data Module Info with %d bytes for Module UUID %s"
                " failed. Error: %s",
                *size, puuid, usys_error(errno));
        }
        *size = sz;
        usys_log_debug(
            "Reading Mfg data Module info with %d bytes for Module UUID %s.",
            *size, puuid);
    }
    return info;
}

/* Reads module config only */
void *jdata_fetch_module_cfg(char *puuid, uint16_t *size, uint8_t count) {
    int ret = 0;
    StoreSchema *db = parser_get_mfg_data_by_uuid(puuid);
    if (!db) {
        usys_log_error("No StoreSchema found for %s module.", puuid);
        return NULL;
    }

    uint16_t sz = 0;
    uint8_t iter = 0;
    ModuleCfg *mcfg;
    for (iter = 0; iter <= count; iter++) {
        int res = usys_strcmp(db->modInfo.uuid, puuid);
        if (!res) {
            /* Fetching Module Config from MFG data*/
            uint16_t csz = sizeof(ModuleCfg) * (db->modInfo.devCount);
            mcfg = (ModuleCfg *)usys_zmalloc(csz);
            if (mcfg) {
                usys_memcpy(mcfg, db->modCfg, csz);
                *size = sz;
                usys_log_debug(
                    "Reading Mfg data Module config with %d bytes for Module"
                    "UUID %s.",
                    *size, puuid);
            } else {
                ret = ERR_NODED_MEMORY_EXHAUSTED;
               usys_log_error(
                    "Reading Mfg data Module config with %d bytes for Module"
                    "UUID %s. Error %s",
                    *size, puuid, usys_error(errno));
            }
        }
        usys_log_debug(
            "Reading Module Config with %d bytes for uuid %s from Mfg data"
            " completed.",
            *size, puuid);
        break;
    }
    return mcfg;
}

void *jdata_fetch_fact_cfg(char *puuid, uint16_t *size) {
    int ret = 0;
    StoreSchema *sschema = parser_get_mfg_data_by_uuid(puuid);
    if (!sschema) {
        usys_log_error("No StoreSchema found for %s module.", puuid);
        return NULL;
    }

    char *info = usys_zmalloc(sizeof(char) * (*size));
    if (info) {
        usys_memcpy(info, sschema->factCfg, sizeof(char) * (*size));
        usys_log_debug(
            "Reading Mfg data fact Config with %d bytes for Module UUID %s.\n"
            "Data::",
            *size, puuid);
        print_mfg_data(info, *size);
    }
    return info;
}

void *jdata_fetch_user_cfg(char *puuid, uint16_t *size) {
    int ret = 0;
    StoreSchema *sschema = parser_get_mfg_data_by_uuid(puuid);
    if (!sschema) {
        usys_log_error("No StoreSchema found for %s module.", puuid);
        return NULL;
    }

    char *info = usys_zmalloc(sizeof(char) * (*size));
    if (info) {
        usys_memcpy(info, sschema->userCfg, sizeof(char) * (*size));
        usys_log_debug(
            "Reading Mfg data User Config with %d bytes for Module UUID %s.\n"
            "Data::",
            *size, puuid);
        print_mfg_data(info, *size);
    }
    return info;
}

void *jdata_fetch_fact_calib(char *puuid, uint16_t *size) {
    int ret = 0;
    StoreSchema *sschema = parser_get_mfg_data_by_uuid(puuid);
    if (!sschema) {
        usys_log_error("No StoreSchema found for %s module.", puuid);
        return NULL;
    }

    char *info = usys_zmalloc(sizeof(char) * (*size));
    if (info) {
        usys_memcpy(info, sschema->factCalib, sizeof(char) * (*size));
        usys_log_debug(
            "Reading Mfg data factory calibration with %d bytes for Module "
            "UUID %s.\n Data::",
            *size, puuid);
        print_mfg_data(info, *size);
    }
    return info;
}

void *jdata_fetch_user_calib(char *puuid, uint16_t *size) {
    int ret = 0;
    StoreSchema *sschema = parser_get_mfg_data_by_uuid(puuid);
    if (!sschema) {
        usys_log_error("No StoreSchema found for %s module.", puuid);
        return NULL;
    }

    char *info = usys_zmalloc(sizeof(char) * (*size));
    if (info) {
        usys_memcpy(info, sschema->userCalib, sizeof(char) * (*size));
        usys_log_debug(
            "Reading Mfg data user calibration with %d bytes for Module "
            "UUID %s.\n Data::",
            *size, puuid);
        print_mfg_data(info, *size);
    }
    return info;
}

void *jdata_fetch_bs_certs(char *puuid, uint16_t *size) {
    int ret = 0;
    StoreSchema *sschema = parser_get_mfg_data_by_uuid(puuid);
    if (!sschema) {
        usys_log_error("No StoreSchema found for %s module.", puuid);
        return NULL;
    }

    char *info = usys_zmalloc(sizeof(char) * (*size));
    if (info) {
        usys_memcpy(info, sschema->bsCerts, sizeof(char) * (*size));
        usys_log_debug(
            "Reading Mfg data Bootstrap certificates with %d bytes for Module "
            "UUID %s.\n Data::",
            *size, puuid);
        print_mfg_data(info, *size);
    }
    return info;
}

void *jdata_fetch_cloud_certs(char *puuid, uint16_t *size) {
    int ret = 0;
    StoreSchema *sschema = parser_get_mfg_data_by_uuid(puuid);
    if (!sschema) {
        usys_log_error("No StoreSchema found for %s module.", puuid);
        return NULL;
    }

    char *info = usys_zmalloc(sizeof(char) * (*size));
    if (info) {
        usys_memcpy(info, sschema->cloudCerts, sizeof(char) * (*size));
        usys_log_debug(
            "Reading Mfg data Cloud certificates with %d bytes for Module "
            "UUID %s.\n"
            "Data::",
            *size, puuid);
        print_mfg_data(info, *size);
    }
    return info;
}

/* Read the payload from json parser or any other input.*/
/* This is meant for input to inventory to create a inventory database  */
int jdata_fetch_payload_from_mfgdata(void **data, char *uuid, uint16_t *size,
                                   uint16_t id) {
    int ret = -1;
    switch (id) {
    case FIELD_ID_FACT_CFG: {
        ret = jdata_fetch_fact_config(data, uuid, size);
        break;
    }
    case FIELD_ID_USER_CFG: {
        ret = jdata_fetch_user_config(data, uuid, size);
        break;
    }

    case FIELD_ID_FACT_CALIB: {
        ret = jdata_fetch_fact_calib(data, uuid, size);
        break;
    }

    case FIELD_ID_USER_CALIB: {
        ret = jdata_fetch_user_calib(data, uuid, size);
        break;
    }
    case FIELD_ID_BS_CERTS: {
        ret = jdata_fetch_bs_certs(data, uuid, size);
        break;
    }
    case FIELD_ID_CLOUD_CERTS: {
        ret = jdata_fetch_lwm2m_certs(data, uuid, size);
        break;
    }
    default: {
        ret = ERR_NODED_JSON_PARSER;
        log_error("Invalid Field id supplied by Index entry. "
                        "Error Code: %d", ret);
    }
    }

    if (!data) {
        data = NULL;
        ret = ERR_NODED_DB_MISSING_INFO;
        log_error("JSON parser failed to read info on 0x%x."
                        "Error Code: %d", id, ret);
    }
    return ret;
}

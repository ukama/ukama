/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "inc/ukdb.h"

#include "headers/errorcode.h"
#include "headers/ubsp/ukdblayout.h"
#include "headers/globalheader.h"
#include "inc/devicedb.h"
#include "ukdb/db/db.h"
#include "ukdb/db/file.h"
#include "ukdb/idb/idb.h"
#include "ukdb/idb/jp.h"
#include "utils/crc32.h"
#include "headers/utils/log.h"

static int validate_unit_type(UnitType unit) {
    int ret = 0;
    switch (unit) {
    case E_TNODESDR:
    case E_TNODELTE:
    case E_HNODE:
    case E_ANODE:
    case E_PSNODE:
        ret = 1;
        break;
    default:
        ret = 0;
    }
    return ret;
}

static int validate_unit_info(UnitInfo *uinfo) {
    int ret = 1;
    if (strncmp(uinfo->uuid, "UK", 2)) {
        ret &= 1;
    }
    if (validate_unit_type(uinfo->unit)) {
        ret &= 1;
    }
    return ret;
}

static int validate_unit_cfg(UnitCfg *cfg, char *fname) {
    int ret = 0;
    if (!strcmp(cfg->sysfs, fname)) {
        ret = 1;
    }
    return ret;
}

/* Based on the module id it returns the UnitCfg */
static UnitCfg *override_master_db_info(char *puuid, uint8_t *count) {
    typedef struct UnitData {
        char id[24];
        uint8_t count;
        UnitCfg *cfg;
    } UnitData;

    UnitData *udata = (UnitData[]){
        { .id = "UK1001-COMV1",
          .count = 3,
          .cfg =
              (UnitCfg[]){
                  { .mod_uuid = "UK1001-COMV1",
                    .mod_name = "COMv1",
                    .sysfs = "/tmp/sys/bus/i2c/devices/i2c-0/0-0050/eeprom",
                    .eeprom_cfg = &(DevI2cCfg){ .bus = 0, .add = 0x50ul } },
                  { .mod_uuid = "UK2001-LTE",
                    .mod_name = "LTE",
                    .sysfs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0050/eeprom",
                    .eeprom_cfg = &(DevI2cCfg){ .bus = 1, .add = 0x50ul } },
                  { .mod_uuid = "UK3001-MASK",
                    .mod_name = "MASK",
                    .sysfs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0051/eeprom",
                    .eeprom_cfg = &(DevI2cCfg){ .bus = 1, .add = 0x51ul } },
              } },
        { .id = "UK5001-RFCTRL",
          .count = 2,
          .cfg =
              (UnitCfg[]){
                  { .mod_uuid = "UK5001-RFCTRL",
                    .mod_name = "RF CTRL BOARD",
                    .sysfs = "/tmp/sys/bus/i2c/devices/i2c-0/0-0051/eeprom",
                    .eeprom_cfg = &(DevI2cCfg){ .bus = 1, .add = 0x50ul } },
                  { .mod_uuid = "UK4001-RFFE",
                    .mod_name = "RF BOARD",
                    .sysfs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0052/eeprom",
                    .eeprom_cfg = &(DevI2cCfg){ .bus = 2, .add = 0x50ul } },
              } }
    };

    UnitCfg *pcfg = NULL;
    DevI2cCfg *i2c_cfg = NULL;
    for (int iter = 0; iter < 2; iter++) {
        if (!strcmp(puuid, udata[iter].id)) {
            *count = udata[iter].count;
            pcfg = malloc(sizeof(UnitCfg) * (*count));
            if (pcfg) {
                memset(pcfg, '\0', sizeof(UnitCfg) * (*count));
                memcpy(&pcfg[0], &udata[iter].cfg[0],
                       sizeof(UnitCfg) * (*count));
                for (int iiter = 0; iiter < *count; iiter++) {
                    if (udata[iter].cfg[iiter].eeprom_cfg) {
                        i2c_cfg = malloc(sizeof(DevI2cCfg));
                        if (i2c_cfg) {
                            memset(i2c_cfg, '\0', sizeof(DevI2cCfg));
                            memcpy(i2c_cfg, udata[iter].cfg[iiter].eeprom_cfg,
                                   sizeof(DevI2cCfg));
                        }
                    }
                    pcfg[iiter].eeprom_cfg = i2c_cfg;
                }
                break;
            } else {
                log_error(
                    "Err(%d): UKDB: Memory exhausted while getting unit config from Test data.",
                    ERR_UBSP_MEMORY_EXHAUSTED);
            }
        }
    }
    return pcfg;
}

static int get_master_db_info(UnitCfg *pcfg, char *systemlnkdb) {
    int ret = 0;
    if (!systemlnkdb) {
        return -1;
    }
    /* Read the system-db file*/
    char *systemdb;
    systemdb = file_read_sym_link(systemlnkdb);
    if (systemdb) {
        if (!file_exist(systemdb)) {
            ret = ERR_UBSP_DB_MISSING;
            log_error(
                "Err(%d): UKDB:: Ukama DB for the system is not available at %s.",
                ret, systemdb);
            free(systemdb);
            return ret;
        }

        UnitInfo *uinfo = malloc(sizeof(UnitInfo));
        if (uinfo) {
            if (file_raw_read(systemdb, uinfo, UKDB_UNIT_INFO_OFFSET,
                              sizeof(UnitInfo)) == sizeof(UnitInfo)) {
                if (validate_unit_info(uinfo)) {
                    log_debug("UKDB:: Unit UIID %s Name %s detected.",
                              uinfo->uuid, uinfo->name);
                    /* Read first cfg which belong to master */
                    int sz = sizeof(UnitCfg) + sizeof(DevI2cCfg);
                    void *cfg = malloc(sz);
                    if (cfg) {
                        if (file_raw_read(systemdb, cfg,
                                          UKDB_UNIT_CONFIG_OFFSET, sz) == sz) {
                            if (validate_unit_cfg(cfg, systemdb)) {
                                memcpy(pcfg, cfg, sizeof(UnitCfg));
                                DevI2cCfg *icfg = malloc(sizeof(DevI2cCfg));
                                if (icfg) {
                                    memcpy(icfg, (cfg + sizeof(UnitCfg)),
                                           sizeof(DevI2cCfg));
                                } else {
                                    icfg = NULL;
                                }
                                pcfg->eeprom_cfg = icfg;
                            } else {
                                ret = ERR_UBSP_INVALID_UNIT_CFG;
                                log_error(
                                    "Err(%d): UKDB:: Invalid Unit Config.",
                                    ret);
                            }
                        } else {
                            ret = ERR_UBSP_INVALID_UNIT_CFG;
                            log_error(
                                "Err(%d): UKDB:: Invalid Unit Config read.",
                                ret);
                        }
                        free(cfg);
                    }
                } else {
                    ret = ERR_UBSP_INVALID_UNIT_INFO;
                    log_error("Err(%d): UKDB:: Invalid Unit Info.", ret);
                }
            } else {
                ret = ERR_UBSP_INVALID_UNIT_INFO;
                log_error("Err(%d): UKDB:: Invalid Unit Info read.", ret);
            }
            free(uinfo);
        } else {
            ret = ERR_UBSP_MEMORY_EXHAUSTED;
            log_error("Err(%d): UKDB:: System out of memory.", ret);
        }
        free(systemdb);
    } else {
        ret = ERR_UBSP_DB_LNK_MISSING;
        log_error(
            "Err(%d): UKDB:: Ukama DB link for the systemdb is not available.",
            ret);
    }
    return ret;
}

int ukdb_validating_magicword(char *p_uuid) {
    int ret = 0;
    UKDBMagicWord ukdb_mw;
    /*Validating*/
    if (db_read_block(p_uuid, &ukdb_mw, UKDB_MAGICWORD_OFFSET,
                      sizeof(UKDBMagicWord)) != sizeof(UKDBMagicWord)) {
        log_error("Err(%d): UKDB read error.", ERR_UBSP_R_FAIL);
        return ERR_UBSP_R_FAIL;
    }
    if (ukdb_mw.magic_word == UKDB_MAGICWORD) {
        log_debug("UKDB:: Ukama DB MagicWord  validation passed for module %s.",
                  p_uuid);
    } else {
        ret = ERR_UBSP_MW_ERR;
        log_error(
            "Err(%d): Ukama DB MagicWord validation failed for module %s.", ret,
            p_uuid);
    }
    return ret;
}

int ukdb_write_magicword(char *p_uuid) {
    int ret = 0;
    UKDBMagicWord ukdb_mw;
    ukdb_mw.magic_word = UKDB_MAGICWORD;
    ukdb_mw.resv1 = UKDB_DEFVAL;
    ukdb_mw.resv2 = UKDB_DEFVAL;
    if (db_write_block(p_uuid, &ukdb_mw, UKDB_MAGICWORD_OFFSET,
                       sizeof(UKDBMagicWord)) != sizeof(UKDBMagicWord)) {
        ret = ERR_UBSP_WR_FAIL;
        log_error("Err(%d): UKDB read error.", ret);
    }

    if (ukdb_validating_magicword(p_uuid)) {
        ret = ERR_UBSP_MW_ERR;
    }
    return ret;
}

/* Write UKDB header*/
int ukdb_write_header(char *p_uuid, UKDBHeader *p_header) {
    int ret = 0;
    if (db_write_block(p_uuid, p_header, UKDB_HEADER_OFFSET,
                       UKDB_HEADER_SIZE) != UKDB_HEADER_SIZE) {
        ret = ERR_UBSP_WR_FAIL;
        log_error("Err(%d): UKDB Header write of %d bytes failed at "
                  "offset 0x%x.",
                  ret, UKDB_HEADER_SIZE, UKDB_HEADER_OFFSET);

    } else {
        log_debug("UKDB:: Header write of %d bytes completed at "
                  "offset 0x%x.",
                  UKDB_HEADER_SIZE, UKDB_HEADER_OFFSET);
    }
    return ret;
}

/* Read UK DB header*/
int ukdb_read_header(char *p_uuid, UKDBHeader *p_header) {
    int ret = 0;
    if (db_read_block(p_uuid, p_header, UKDB_HEADER_OFFSET, UKDB_HEADER_SIZE) <
        UKDB_HEADER_SIZE) {
        ret = ERR_UBSP_R_FAIL;
        log_error("Err(%d): UKDB Header read of %d bytes from "
                  "offset 0x%x failed.",
                  ret, UKDB_HEADER_SIZE, UKDB_HEADER_OFFSET);
    } else {
        log_debug("UKDB:: Header read of %d bytes completed from "
                  "offset 0x%x.",
                  UKDB_HEADER_SIZE, UKDB_HEADER_OFFSET);
        ukdp_print_header(p_header);
    }
    return ret;
}

/* Read db version*/
int ukdb_read_dbversion(char *p_uuid, Version *ver) {
    int ret = 0;
    if (db_read_block(p_uuid, ver, UKDB_HEADER_DBVER_OFFSET, sizeof(Version)) <
        sizeof(Version)) {
        ret = ERR_UBSP_R_FAIL;
        log_error("Err(%d): UKDB version read of %d bytes "
                  "from offset 0x%x failed.",
                  ret, sizeof(Version), UKDB_HEADER_DBVER_OFFSET);
    } else {
        log_debug("UKDB:: Version read of v%d.%d.", ver->major, ver->minor);
    }
    return ret;
}
/* Update UKDB version*/
int ukdb_update_dbversion(char *p_uuid, Version ver) {
    int ret = 0;
    if (db_write_block(p_uuid, &ver, UKDB_HEADER_DBVER_OFFSET,
                       sizeof(Version)) != sizeof(Version)) {
        ret = ERR_UBSP_R_FAIL;
        log_error("Err(%d): UKDB version write of %d bytes "
                  "to offset 0x%x failed.",
                  ret, sizeof(Version), UKDB_HEADER_DBVER_OFFSET);
    } else {
        log_debug("UKDB:: Version updated to v%d.%d.", ver.major, ver.minor);
    }
    return ret;
}

/*Update Index table current index count. */
int ukdb_update_current_idx_count(char *p_uuid, uint16_t *idx_count) {
    int ret = 0;
    if (db_write_number(p_uuid, idx_count, UKDB_IDX_CUR_TPL_COUNT_OFFSET, 1,
                        UKDB_IDX_COUNT_SIZE)) {
        ret = ERR_UBSP_WR_FAIL;
        log_error("Err(%d): UKDB write for index count failed.", ret);
    }
    return ret;
}

/*Update Index table current index count. */
int ukdb_read_current_idx_count(char *p_uuid, uint16_t *idx_count) {
    int ret = 0;
    if (db_read_number(p_uuid, idx_count, UKDB_IDX_CUR_TPL_COUNT_OFFSET, 1,
                       UKDB_IDX_COUNT_SIZE)) {
        ret = ERR_UBSP_R_FAIL;
        log_error("Err(%d): UKDB read for index count failed.", ret);
    }

    if (!(VALIDATE_INDEX_COUNT(*idx_count))) {
        ret = ERR_UBSP_VALIDATION_FAILURE;
        log_error("Err(%d): UKDB validation for index count (%d) failed.", ret,
                  *idx_count);
    }

    return ret;
}

/* Erasing the Database.*/
int ukdb_erase_db(char *p_uuid) {
    int ret = 0;
    //TODO: Add this db_erase api.
    if (db_erase_block(p_uuid, UKDB_START_OFFSET, UKDB_END_OFFSET) !=
        UKDB_END_OFFSET) {
        ret = ERR_UBSP_WR_FAIL;
        log_error("Err(%d) : UKDB erase failed.", ret);
    }
    return ret;
}

/* Erase the last index from the Index table.*/
int ukdb_erase_index(char *p_uuid) {
    /* Erase is done in LIFO manner */
    int ret = 0;

    /* Read available index count. */
    uint16_t idx_count = 0;
    if (ukdb_read_current_idx_count(p_uuid, &idx_count)) {
        return ERR_UBSP_R_FAIL;
    }

    int offset = idx_count * UKDB_IDX_COUNT_SIZE;
    if (db_erase_block(p_uuid, offset, UKDB_IDX_TPL_SIZE) ==
        UKDB_IDX_TPL_SIZE) {
        return ERR_UBSP_WR_FAIL;
    }

    /* Update Index count. */
    idx_count--;
    if (ukdb_update_current_idx_count(p_uuid, &idx_count)) {
        ret = ERR_UBSP_WR_FAIL;
    }
    log_debug("UKDB:: Index deleted from index id 0x%x at offset 0x%x.",
              idx_count + 1, offset);

    return ret;
}

/* Update the index at particular location */
int ukdb_update_index_atid(char *p_uuid, UKDBIdxTuple *p_data, uint16_t idx) {
    int ret = 0;
    int offset = idx * UKDB_IDX_COUNT_SIZE;
    if (db_write_block(p_uuid, p_data, offset, UKDB_IDX_TPL_SIZE) !=
        UKDB_IDX_TPL_SIZE) {
        ret = ERR_UBSP_WR_FAIL;
        log_error(
            "Err(%d): UKDB Index update failed for index id 0x%x at offset 0x%x.",
            ret, idx, idx * UKDB_IDX_COUNT_SIZE);
    }
    log_debug("UKDB:: Index updated for index id 0x%x at offset 0x%x.", idx,
              idx * UKDB_IDX_COUNT_SIZE);
    return ret;
}

/* Search for the field id with valid data.*/

int ukdb_search_fieldid(char *p_uuid, UKDBIdxTuple **idx_data, uint16_t *idx,
                        uint16_t fieldid) {
    int ret = ERR_UBSP_DB_MISSING_FIELD;
    /* Read available index count. */
    uint16_t idx_count = 0;
    if (ukdb_read_current_idx_count(p_uuid, &idx_count)) {
        return ERR_UBSP_R_FAIL;
    }
    *idx = 0xFFFF;
    uint16_t iter = 0;
    UKDBIdxTuple *index = (UKDBIdxTuple *)malloc(sizeof(UKDBIdxTuple));
    if (index) {
        for (; iter < idx_count; iter++) {
            ukdb_read_index(p_uuid, index, iter);
            /* Only the fieldid and valid state matters */
            if ((index->fieldid == fieldid)) {
                /* Mark ret as invalid */
                ret = ERR_UBSP_INVALID_FIELD;
                if (index->valid) {
                    *idx = iter;
                    ret = 0;
                    *idx_data = index;
                    break;
                }
            }
        }
        if (ret) {
            UBSP_FREE(index);
        }
    }
    return ret;
}

/*Update the index for a particular fieldid */
int ukdb_update_index_for_fieldid(char *p_uuid, void *p_data,
                                  uint16_t fieldid) {
    int ret = 0;
    uint16_t idx = 0xFF;
    UKDBIdxTuple *idx_data;
    if (ukdb_search_fieldid(p_uuid, &idx_data, &idx, fieldid)) {
        ret = ERR_UBSP_DB_MISSING_FIELD;
        log_warn("Err(%d): No entry in UKDB with field Id 0x%x.", ret, fieldid);
        return ret;
    }

    if (db_write_block(p_uuid, p_data, idx_data->payload_offset,
                       idx_data->payload_size) != UKDB_IDX_TPL_SIZE) {
        log_error(
            "UKDB::Index write failed for index id 0x%x field id 0x%x at offset 0x%x.",
            idx, idx_data->fieldid, idx * UKDB_IDX_COUNT_SIZE);
        ret = ERR_UBSP_WR_FAIL;
    }
    log_debug(
        "UKDB:: Index updated for index id 0x%x field Id 0x%x at offset 0x%x.",
        idx, idx_data->fieldid, idx * UKDB_IDX_COUNT_SIZE);
    return ret;
}

/* Write index to next avialable index slot in UKDB */
int ukdb_write_index(char *p_uuid, UKDBIdxTuple *p_data) {
    int ret = 0;
    /* Read available index count. */
    uint16_t idx_count = 0;
    ret = ukdb_read_current_idx_count(p_uuid, &idx_count);
    if (ret) {
        return ERR_UBSP_R_FAIL;
    }
    log_trace("UKDB:: UKDB Current Index in Header is %d.", idx_count);
    /*Write the index. */
    int offset = UKDB_IDX_TABLE_OFFSET + idx_count * UKDB_IDX_TPL_SIZE;
    if (db_write_block(p_uuid, p_data, offset, UKDB_IDX_TPL_SIZE) !=
        UKDB_IDX_TPL_SIZE) {
        return ERR_UBSP_WR_FAIL;
    }

    /* Update Index count. */
    idx_count++;
    if (ukdb_update_current_idx_count(p_uuid, &idx_count)) {
        ret = ERR_UBSP_WR_FAIL;
    }
    log_trace("UKDB:: Written index 0x%x at 0x%x offset.", idx_count, offset);
    return ret;
}

/* Read the index at nth location from the Index table of the UKDB*/
int ukdb_read_index(char *p_uuid, UKDBIdxTuple *p_data, uint16_t idx) {
    int ret = 0;
    int offset = UKDB_IDX_TABLE_OFFSET + (idx * UKDB_IDX_TPL_SIZE);
    if (db_read_block(p_uuid, p_data, offset, UKDB_IDX_TPL_SIZE) <
        UKDB_IDX_TPL_SIZE) {
        ret = ERR_UBSP_R_FAIL;
        log_error("Err(%d): UKDB Index read failed from id 0x%x.", ret, idx);
    }
    return ret;
}

/* Erase the payload data for field id and mark the index as invalid. */
int ukdb_erase_payload(char *p_uuid, uint16_t fieldid) {
    int ret = 0;
    uint16_t idx = 0x00;
    UKDBIdxTuple *idx_data;
    if (ukdb_search_fieldid(p_uuid, &idx_data, &idx, fieldid)) {
        ret = ERR_UBSP_DB_MISSING_FIELD;
        log_warn("Err(%d): No entry in UKDB with field Id 0x%x.", ret, fieldid);
        return ret;
    }

    //erase is basically writing 0xFF to DB.
    if (db_erase_block(p_uuid, idx_data->payload_offset,
                       idx_data->payload_size) == UKDB_IDX_TPL_SIZE) {
        return ERR_UBSP_WR_FAIL;
    }
    log_debug("UKDB:: payload erased for field Id 0x%x at offset 0x%x.",
              idx_data->fieldid, idx_data->payload_offset);
    idx_data->state = UKDB_FEAT_DISABLED;
    idx_data->valid = FALSE;
    log_debug("UKDB:: Marking index %d for UKDB as invalid.", idx);
    if (ukdb_update_index_atid(p_uuid, idx_data, idx)) {
        ret = ERR_UBSP_WR_FAIL;
    }
    return ret;
}

/* Write payload to the UKDB */
int ukdb_write_payload(char *p_uuid, void *p_data, uint16_t offset,
                       uint16_t size) {
    int ret = 0;
    if (db_write_block(p_uuid, p_data, offset, size) != size) {
        log_error("UKDB:: Payload write failed at offset 0x%x of size 0x%x.",
                  offset, size);
        ret = ERR_UBSP_WR_FAIL;
    }
    log_trace("UKDB:: Wrote %d bytes of payload to 0x%x offset.", size, offset);
    return ret;
}

/* Write module specific payload */
/* TODO :: Each Module need to write info to specific device*/
/* For now it might be just a simple file.*/
int ukdb_write_module_payload(char *p_uuid, void *p_data, uint16_t offset,
                              uint16_t size) {
    int ret = 0;
    //TODO: Check if required otherwise delete and use above fxn.
    if (db_write_block(p_uuid, p_data, offset, size) != size) {
        log_error(
            "UKDB:: Module payload write failed at offset 0x%x of size 0x%x.",
            offset, size);
        ret = ERR_UBSP_WR_FAIL;
    }
    log_trace("UKDB:: Wrote %d bytes of module payload to 0x%x offset.", size,
              offset);
    return ret;
}

/* Update the payload data with the given field id. It also Update the index in Index table for field Id. */
int ukdb_update_payload(char *p_uuid, void *p_data, uint16_t fieldid,
                        uint16_t size, uint8_t state, Version version) {
    int ret = 0;
    /*Find index first.
        Compare the field id's from 0th index to last aailable index.
        If field Id matches read payload offset calculate crc and write payload.
	 */
    uint16_t idx = 0x0;
    UKDBIdxTuple *idx_data;
    if (ukdb_search_fieldid(p_uuid, &idx_data, &idx, fieldid)) {
        log_warn("UKDB:: No entry in UKDB with fieldid 0x%x.", fieldid);
        return ERR_UBSP_DB_MISSING_FIELD;
    }

    /*Write payload for Index table entries*/
    if (ukdb_write_payload(p_uuid, p_data, idx_data->payload_offset, size)) {
        ret = ERR_UBSP_WR_FAIL;
        log_error(
            "Err(%d): UKDB Payload write failed for index id 0x%x fieldid 0x%x at offset 0x%x.",
            ret, idx, idx_data->fieldid, idx_data->payload_offset);
        return ret;
    }

    /* Calculate CRC */
    uint32_t crc = crc_32(p_data, size);

    /* Update Index */
    idx_data->payload_crc = crc;
    idx_data->payload_size = size;
    idx_data->payload_version = version;
    idx_data->state = state;
    idx_data->valid = TRUE;
    if (ukdb_update_index_atid(p_uuid, idx_data, idx)) {
        ret = ERR_UBSP_WR_FAIL;
        log_error(
            "Err(%d): UKDB Index write failed for index id 0x%x field Id 0x%x at offset 0x%x.",
            ret, idx, idx_data->fieldid, idx * UKDB_IDX_COUNT_SIZE);
        return ret;
    }
    log_debug(
        "UKDB:: Payload updated for index id 0x%x field Id 0x%x at offset 0x%x.",
        idx, idx_data->fieldid, idx_data->payload_offset);

    return ret;
}

UnitCfg *ukdb_alloc_unit_cfg(uint8_t count) {
    UnitCfg *cfg = malloc(sizeof(UnitCfg) * count);
    if (cfg) {
        memset(cfg, '\0', sizeof(UnitCfg));
    }
    return cfg;
}

void ukdb_free_unit_cfg(UnitCfg *cfg, uint8_t count) {
    if (cfg) {
        for (int iter = 0; iter < count; iter++) {
            UBSP_FREE(cfg[iter].eeprom_cfg);
        }
        UBSP_FREE(cfg);
    }
}

char *serialize_unitcfg_payload(UnitCfg *ucfg, uint8_t count, uint16_t *size) {
    int offset = 0;
    char *data = NULL;
    *size = (sizeof(UnitCfg) + sizeof(DevI2cCfg)) * count;
    data = malloc(*size);
    if (data) {
        for (int iter = 0; iter < count; iter++) {
            memcpy(data + offset, &ucfg[iter], sizeof(UnitCfg));
            offset = offset + sizeof(UnitCfg);
            memcpy(data + offset, ucfg[iter].eeprom_cfg, sizeof(DevI2cCfg));
            offset = offset + sizeof(DevI2cCfg);
        }
    } else {
        data = NULL;
    }
    return data;
}

/* Write Unit Config and update the index to index table of the DB */
int ukdb_write_unit_cfg_data(char *p_uuid, UKDBIdxTuple *index, uint8_t count) {
    int ret = 0;
    /* Unit Config */
    UnitCfg *ucfg;
    uint16_t size = 0;
    char *payload = NULL;
    ret = idb_fetch_unit_cfg((void *)&ucfg, p_uuid, &size, count);
    if (ucfg) {
        /*Write payload for Index table entries*/
        payload = serialize_unitcfg_payload(ucfg, count, &size);
        if (payload) {
            if (ukdb_write_payload(p_uuid, payload, index->payload_offset,
                                   size)) {
                /* Need to revert back index */
                ukdb_erase_index(p_uuid);
                ret = ERR_UBSP_WR_FAIL;
                goto cleanup;
            }
            /*Update payload size */
            index->payload_size = size;
            /* Add CRC */
            uint32_t crc_val =
                crc_32((const unsigned char *)payload, index->payload_size);
            index->payload_crc = crc_val;
            log_debug("UKDB:: Calculated CRCfor filed Id %d payload is 0x%x",
                      index->fieldid, crc_val);
            log_debug("UKDB:: Payload added for field Id 0x%x", index->fieldid);

        } else {
            ret = ERR_UBSP_MEMORY_EXHAUSTED;
            log_error(
                "Err(%d): UKDB:: Memory exhausted while serializing Unit Cfg.",
                ret);
            goto cleanup;
        }

        /* Write Index table entries */
        if (ukdb_write_index(p_uuid, index)) {
            ret = ERR_UBSP_WR_FAIL;
            goto cleanup;
        }

        log_debug("UKDB:: Index added to the Index table  for field Id 0x%x",
                  index->fieldid);
    } else {
        log_error("Err(%d): Failed to read Unit config from Mfg data.", ret);
    }
cleanup:
    /* Only free the ucfg not the module config.*/
    UBSP_FREE(payload);
    ukdb_free_unit_cfg(ucfg, count);
    return ret;
}

/* Write Unit Info and update the index to index table of the DB */
int ukdb_write_unit_info_data(char *p1_uuid, UKDBIdxTuple *index, char *p_uuid,
                              uint8_t *count) {
    int ret = 0;
    UnitInfo *uinfo;
    uint16_t size = 0;
    ret = idb_fetch_unit_info((void *)&uinfo, p_uuid, &size);
    if (!ret) {
        /* Unit Info */
        if (uinfo) {
            /*Write payload for Index table entries*/
            if (ukdb_write_payload(p1_uuid, uinfo, index->payload_offset,
                                   index->payload_size)) {
                /* Need to revert back index */
                ukdb_erase_index(p1_uuid);
                ret = ERR_UBSP_WR_FAIL;
                goto cleanup;
            }
            log_debug("UKDB:: Unit Info added for unit Id %s.", uinfo->uuid);

            /* CRC*/
            uint32_t crc_val =
                crc_32((const unsigned char *)uinfo, index->payload_size);
            log_debug("UKDB:: Calculated crc for unit id %s is 0x%x",
                      uinfo->uuid, crc_val);
            index->payload_crc = crc_val;

            /* Write Index table entries */
            if (ukdb_write_index(p1_uuid, index)) {
                ret = ERR_UBSP_WR_FAIL;
                goto cleanup;
            }

            /*Assign module count.*/
            *count = uinfo->mod_count;
            if (!VALIDATE_MODULE_COUNT(*count)) {
                ret = ERR_UBSP_VALIDATION_FAILURE;
                log_error(
                    "Err(%d): UKDB:: Validation for module count (%d) failed.",
                    ret, *count);
                goto cleanup;
            }

            log_debug(
                "UKDB:: Index added to the Index table for field Id 0x%x and Unit Id %s.",
                index->fieldid, uinfo->uuid);
        }
    } else {
        log_error("Err(%d): Failed to read Unit info from Mfg data.", ret);
    }
cleanup:
    UBSP_FREE(uinfo);

    return ret;
}

ModuleCfg *ukdb_alloc_module_cfg(uint8_t count) {
    ModuleCfg *cfg = malloc(sizeof(ModuleCfg) * count);
    if (cfg) {
        memset(cfg, '\0', sizeof(ModuleCfg) * count);
    }
    return cfg;
}

/* Free the memory used by moduleCfg */
void ukdb_free_module_cfg(ModuleCfg *cfg, uint8_t count) {
    if (cfg) {
        for (int iter = 0; iter < count; iter++) {
            void *dev_cfg = cfg[iter].cfg;
            UBSP_FREE(dev_cfg);
        }
        UBSP_FREE(cfg);
    }
}

/* Serialize the Module Cfg data */
char *serialize_module_config_data(ModuleCfg *mcfg, uint8_t count,
                                   uint16_t *size) {
    int offset = 0;
    char *data = NULL;
    int psize = 0;
    for (int iter = 0; iter < count; iter++) {
        uint16_t cfg_size = 0;
        SIZE_OF_DEVICE_CFG(cfg_size, mcfg[iter].dev_class);
        psize = psize + cfg_size;
    }
    psize = (sizeof(ModuleCfg) * count) + psize;
    data = malloc(psize);
    if (data) {
        for (int iter = 0; iter < count; iter++) {
            memcpy(data + offset, &mcfg[iter], sizeof(ModuleCfg));
            offset = offset + sizeof(ModuleCfg);
            uint16_t cfg_size = 0;
            SIZE_OF_DEVICE_CFG(cfg_size, mcfg[iter].dev_class);
            memcpy(data + offset, mcfg[iter].cfg, cfg_size);
            offset = offset + cfg_size;
        }
        *size = psize;
    } else {
        data = NULL;
    }
    return data;
}

/* Write Module Config and update the index to index table of the DB */
int ukdb_write_module_cfg_data(char *p_uuid, ModuleInfo *minfo,
                               UKDBIdxTuple *cfg_index) {
    int ret = 0;
    char *payload = NULL;
    uint16_t size = 0;

    /* Serialize the module config data.*/
    payload = serialize_module_config_data(minfo->module_cfg, minfo->dev_count,
                                           &size);
    if (payload) {
        /*Write Module Cfg for the module*/
        if (ukdb_write_module_payload(p_uuid, payload,
                                      cfg_index->payload_offset, size)) {
            /* Need to revert back index */
            ukdb_erase_index(p_uuid);
            ret = ERR_UBSP_WR_FAIL;
            goto cleanup;
        }
        cfg_index->payload_size = size;
        log_debug("UKDB:: Module Config added for module Id %s  ", minfo->uuid);
        uint32_t crc_val =
            crc_32((const unsigned char *)payload, cfg_index->payload_size);
        cfg_index->payload_crc = crc_val;
        log_debug("UKDB:: Calculated CRC for module Id %s is 0x%x", minfo->uuid,
                  crc_val);
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_error(
            "Err(%d): UKDB:: Memory exhausted while serializing Module config.",
            ret);
        goto cleanup;
    }

    /* Write Index table entries */
    if (ukdb_write_index(p_uuid, cfg_index)) {
        ret = ERR_UBSP_WR_FAIL;
        goto cleanup;
    }
    log_debug(
        "UKDB:: Index added to the Index table for field Id 0x%x and module Id %s",
        cfg_index->fieldid, minfo->uuid);

cleanup:
    UBSP_FREE(payload);
    ukdb_free_module_cfg(minfo->module_cfg, minfo->dev_count);
    return ret;
}

/* Write module info to UKDB*/
int ukdb_write_module_info_data(char *p_uuid, UKDBIdxTuple *info_index,
                                UKDBIdxTuple *cfg_index, uint8_t *mcount,
                                uint8_t *idx) {
    int ret = 0;
    ModuleInfo *minfo;
    uint16_t size = 0;
    ret = idb_fetch_module_info_by_uuid((void *)&minfo, p_uuid, &size, *mcount);
    if (!ret) {
        /* Unit Info */
        if (minfo) {
            /*Write payload for Index table entries*/
            if (ukdb_write_module_payload(p_uuid, minfo,
                                          info_index->payload_offset,
                                          info_index->payload_size)) {
                /* Need to revert back index */
                ukdb_erase_index(p_uuid);
                ret = ERR_UBSP_WR_FAIL;
                goto cleanmid;
            }
            log_debug(
                "UKDB:: Added Module Info for module Id %s with %d devices.",
                p_uuid, minfo->dev_count);
            uint32_t crc_val =
                crc_32((const unsigned char *)minfo, info_index->payload_size);
            info_index->payload_crc = crc_val;
            log_debug("UKDB:: Calculated CRC32 for Module %s is 0x%x", p_uuid,
                      crc_val);
            /* Write Index table entries */
            if (ukdb_write_index(p_uuid, info_index)) {
                ret = ERR_UBSP_WR_FAIL;
                goto cleanmid;
            }
            //Updating write index count.
            (*idx)++;
            log_debug(
                "UKDB:: Index added to the Index table for field Id 0x%x module Id %s",
                info_index->fieldid, p_uuid);

            log_debug("UKDB:: Adding Module Config now for module Id %s",
                      p_uuid);
            /* Write Module Cfg as well Module info contains the info for Module cfg also.*/
            //TODO : lot better option will ope after ukdb_create is complete.
            ret = ukdb_write_module_cfg_data(p_uuid, minfo, cfg_index);
        }

    } else {
        log_error("Err(%d): Failed to read Unit info from Mfg data.", ret);
    }
cleanmid:
    UBSP_FREE(minfo);

    return ret;
}

/* Search for the field Id in the Index Table read from MFG data.*/
int get_fieldid_index(UKDBIdxTuple *index, uint16_t fid, uint8_t idx_count,
                      uint8_t *idx_iter) {
    int ret = -1;
    uint8_t iter = 0;
    for (; iter < idx_count; iter++) {
        if (index[iter].fieldid == fid) {
            *idx_iter = iter;
            log_trace(
                "UKDB:: UKDB Field Id 0x%04x found at MFG data location mfg_index_table[%d].",
                fid, iter);
            ret = 0;
            break;
        }
    }
    return ret;
}

/* Set up environment for DB creation process.*/
int ukdb_idb_init(void *data) {
    int ret = 0;
    ret = idb_init(data);
    return ret;
}

/* Write generic payload to db like some files.*/
int ukdb_write_generic_data(UKDBIdxTuple *index, char *p_uuid, uint16_t fid) {
    int ret = 0;
    uint16_t size = index->payload_size;
    char *payload;
    ret = idb_fetch_payload_from_mfgdata((void *)&payload, p_uuid, &size, fid);
    if (payload) {
        /*Write payload for Index table entries*/
        if (ukdb_write_payload(p_uuid, payload, index->payload_offset,
                               index->payload_size)) {
            /* Need to revert back index */
            ukdb_erase_index(p_uuid);
            ret = ERR_UBSP_WR_FAIL;
            goto cleanpayload;
        }
        log_debug("UKDB:: Added payload for field Id 0x%x.", index->fieldid);
        uint32_t crc_val =
            crc_32((const unsigned char *)payload, index->payload_size);
        index->payload_crc = crc_val;
        log_debug("UKDB:: Calculated CRC32 for field Id %d is 0x%x",
                  index->fieldid, crc_val);
        /* Write Index table entries */
        if (ukdb_write_index(p_uuid, index)) {
            ret = ERR_UBSP_WR_FAIL;
            goto cleanpayload;
        }
        log_debug(
            "UKDB:: Index added to the Index table at id 0x%x for field Id 0x%x",
            index->fieldid, index->fieldid);
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug(
            "Err(%d): Memory exhausted for while reading Feild Id 0x%x for UUID %s.",
            index->fieldid, p_uuid);
        return ret;
    }

cleanpayload:
    UBSP_FREE(payload);

    return ret;
}

int ukdb_register_module(UnitCfg *cfg) {
    int ret = 0;
    ret = db_register_module(cfg);
    if (ret) {
        log_error("Err(%d): UKDB registering module uuid %s failed.", ret,
                  cfg->mod_uuid);
    }
    return ret;
}

int ukdb_unregister_module(char *puuid) {
    int ret = 0;
    ret = db_unregister_module(puuid);
    if (ret) {
        log_error("Err(%d): UKDB unregistering module uuid %s failed.", ret,
                  puuid);
    }
    return ret;
}

/* UKDB registering Modules in DB so that these can be accessed later on.*/
int ukdb_register_modules(char *p_uuid) {
    int ret = 0;
    uint16_t size = 0;
    uint8_t mod_count = 0;
    UnitCfg *ucfg;
    char unit_uuid[24] = { '\0' };
    log_debug(
        "UKDB:: Read UKDB Unit Info from module %s for module registration process.",
        p_uuid);
    UnitInfo *uinfo = (UnitInfo *)malloc(sizeof(UnitInfo));
    if (uinfo) {
        ret = ukdb_read_unit_info(p_uuid, uinfo, &size);
        if (!ret) {
            ukdb_print_unit_info(uinfo);
            memcpy(unit_uuid, uinfo->uuid, strlen(uinfo->uuid));
            mod_count = uinfo->mod_count;
            if (uinfo) {
                free(uinfo);
                uinfo = NULL;
            }
            log_debug(
                "UKDB:: Read UKDB Unit Config for %s for registration process.",
                unit_uuid);
            ucfg = (UnitCfg *)ukdb_alloc_unit_cfg(mod_count);
            if (ucfg) {
                ret = ukdb_read_unit_cfg(p_uuid, ucfg, mod_count, &size);
                if (!ret) {
                    ukdb_print_unit_cfg(ucfg, mod_count);
                    log_debug("UKDB:: Registering %d module for Unit %s.",
                              mod_count, unit_uuid);
                    for (uint8_t iter = 0; iter < mod_count; iter++) {
                        /* Register module*/
                        ret = ukdb_register_module(&ucfg[iter]);
                        if (!ret) {
                            /* Validate DB */
                            ret =
                                ukdb_validating_magicword(ucfg[iter].mod_uuid);
                            if (ret) {
                                log_warn(
                                    "UKDB:: No valid database found for module UUID %s Name %s."
                                    " Moving on to next module.",
                                    ucfg[iter].mod_uuid, ucfg[iter].mod_name);
                                //Un-register module
                                //ret =
                                //   ukdb_unregister_module(ucfg[iter].mod_uuid);
                                continue;
                            } else {
                                /* register device*/
                                ret =
                                    ukdb_register_devices(ucfg[iter].mod_uuid);
                            }
                        } else {
                            goto cleanunitcfg;
                        }
                    }
                } else {
                    log_debug("Err(%d) UKDB:: Read Unit Config fail for %s.",
                              unit_uuid);
                    goto cleanunitcfg;
                }
            } else {
                ret = ERR_UBSP_MEMORY_EXHAUSTED;
                log_debug("Err(%d) UKDB:: Read Unit Config fail for %s.",
                          unit_uuid);
                goto cleanunitinfo;
            }
        } else {
            log_debug("Err(%d) UKDB:: Read Unit Info fail for %s.", ret,
                      p_uuid);
            goto cleanunitinfo;
        }
    }

cleanunitcfg:
    ukdb_free_unit_cfg(ucfg, mod_count);

cleanunitinfo:
    UBSP_FREE(uinfo);
    return ret;
}

/* UKDB registering devices in Device DB so that these can be accessed later on.*/
int ukdb_register_devices(char *p_uuid) {
    int ret = 0;
    uint16_t size = 0;
    uint8_t count = 0;
    ModuleCfg *mcfg;
    char name[24] = { '\0' };
    ModuleInfo *minfo = (ModuleInfo *)malloc(sizeof(ModuleInfo));
    if (minfo) {
        log_debug("UKDB:: Read UKDB Module Info from module %s for device "
                  "registration process.",
                  p_uuid);
        ret = ukdb_read_module_info(p_uuid, minfo, &size);
        if (!ret) {
            ukdb_print_module_info(minfo);
            count = minfo->dev_count;
            memcpy(name, minfo->name, strlen(minfo->name));
            if (minfo) {
                free(minfo);
                minfo = NULL;
            }

            /* Read Module Cfg */
            mcfg = (ModuleCfg *)ukdb_alloc_module_cfg(count);
            if (mcfg) {
                log_debug("UKDB:: Read UKDB Module Info from module %s for"
                          "device registration process.",
                          p_uuid);
                size = 0;
                ret = ukdb_read_module_cfg(p_uuid, mcfg, count, &size);
                if (!ret) {
                    ukdb_print_module_cfg(mcfg, count);
                    /* Register devices in devicedb */
                    ret = devdb_register(p_uuid, name, count, mcfg);

                } else {
                    log_debug("Err(%d) UKDB:: Read Module Config for %s"
                              " failed.",
                              p_uuid);
                    goto cleanmcfg;
                }
            } else {
                ret = ERR_UBSP_MEMORY_EXHAUSTED;
                log_debug("Err(%d)::UKDB Read UKDB Module Config for %s "
                          "failed.",
                          p_uuid);
                goto cleanminfo;
            }

        } else {
            ret = ERR_UBSP_MEMORY_EXHAUSTED;
            log_debug("Err(%d)::UKDB Read UKDB Module Config for %s "
                      "failed.",
                      p_uuid);
            goto cleanminfo;
        }
    }

cleanmcfg:
    ukdb_free_module_cfg(mcfg, count);

cleanminfo:
    UBSP_FREE(minfo);

    return ret;
}

int ukdb_pre_create_db_setup(char *puuid) {
    int ret = -1;
    uint8_t count = 0;
    UnitCfg *cfg = NULL;
    log_debug("UKDB:: Only meant for test setup to create DB.");
    cfg = override_master_db_info(puuid, &count);
    if (cfg) {
        ukdb_print_unit_cfg(cfg, count);
        /* Register modules so that unit info and other cfg can be accessed.*/
        for (int iter = 0; iter < count; iter++) {
            ret = ukdb_register_module(&cfg[iter]);
            if (ret) {
                log_debug("Err(%d): UKDB:: Failed to register module %s.", ret,
                          cfg[iter].mod_uuid);
            }
        }
        for (int iter = 0; iter < count; iter++) {
            UBSP_FREE(cfg[iter].eeprom_cfg);
        }
    } else {
        log_error("Err(%d): Failed to get Unit Config for %s from test config.",
                  ret, puuid);
    }
    UBSP_FREE(cfg);
    return ret;
}

/* Cleaning DB info and mappings*/
void ukdb_exit() {
    db_unregister_all_module();
}

void ukdb_idb_exit() {
    idb_exit();
}
/* Initialize ukdb on the bootup
 * At bootup all modules have there respective db.
 */
int ukdb_init(char *psystemdb) {
    int ret = 0;
    uint8_t count = 1;
    /* Initialize module db*/
    ret = db_init();
    UnitCfg *cfg = malloc(sizeof(UnitCfg));
    if (cfg) {
        ret = get_master_db_info(cfg, psystemdb);
        if (ret) {
            free(cfg);
            return ret;
        }
    }
    ukdb_print_unit_cfg(cfg, count);
    /* Register master module first so that unit info and cfg can be accessed.*/
    ret = ukdb_register_module(cfg);
    /* After registering master module.Access remaining module if any using unit cfg and register those.*/
    if (!ret) {
        /* Check if the database exist or not.*/
        ret = ukdb_validating_magicword(cfg->mod_uuid);
        if (ret) {
            log_warn("UKDB:: No Database found for module UUID %s Name %s.",
                     cfg->mod_uuid, cfg->mod_name);
        } else {
            /* Register other modules if any.*/
            /* Caution:: If Module is registering itself it may not have more modules to register.*/
            ret = ukdb_register_modules(cfg->mod_uuid);
        }
    }
    if (cfg) {
        UBSP_FREE(cfg->eeprom_cfg);
        UBSP_FREE(cfg);
    }
    return ret;
}

/* removed the database for the module */
int ukdb_remove_db(char *puuid) {
    int ret = 0;
    ret = db_remove_database(puuid);
    if (!ret) {
        ret = ukdb_unregister_module(puuid);
        if (ret) {
            log_error("Err(%d): UKDB:: Failed to unregister Module %s.", ret,
                      puuid);
        }
    } else {
        log_error("Err(%d): UKDB:: Failed to remove database for Module %s.",
                  ret, puuid);
    }
    return ret;
}

/*Creating UK DB from scratch using some json/static structs.*/
//TODO try reading the MFG data index table and write data based on the fields present in that.

int ukdb_create_db(char *p_uuid) {
    int ret = 0;
    uint16_t size = 0;

    /* Magic Word */
    UKDBMagicWord ukdb_mw = { 0 };
    ret = db_read_block(p_uuid, &ukdb_mw, UKDB_MAGICWORD_OFFSET,
                        sizeof(UKDBMagicWord)) != sizeof(UKDBMagicWord);
    if (ret) {
        log_warn("Err(%d): Problem while reading UKDB.", ret);
        log_warn("UKDB:: Creating new one.");
        //return ERR_UBSP_R_FAIL;
    }

    if (ukdb_mw.magic_word == UKDB_MAGICWORD) {
        log_warn("UKDB:: Ukama DB is present in device. Re-writing it again.");
    } else {
        /* Write Magic Word */
        if (ukdb_write_magicword(p_uuid)) {
            log_error("UKDB:: Ukama DB creation failed.");
            return ERR_UBSP_WR_FAIL;
        }
    }

    /* Update header info. */
    UKDBHeader *header;
    ret = idb_fetch_header((void *)&header, p_uuid, &size);
    if (ret) {
        UBSP_FREE(header);
        return ret;
    } else {
        if (ukdb_write_header(p_uuid, header)) {
            log_error("Err:: Write to UKDB header failed.");
            UBSP_FREE(header);
            return ERR_UBSP_WR_FAIL;
        }
    }
    UBSP_FREE(header);

    /* Populate Index list from the MFG data.*/
    UKDBIdxTuple *index;
    ret = idb_fetch_index((void *)&index, p_uuid, &size);
    if (ret) {
        return ret;
    }
    uint8_t idx_count = size / sizeof(UKDBIdxTuple);
    ukdb_print_index_table(index, idx_count);

    /* For each index table entry */
    uint8_t idx = 0;
    uint8_t mod_count = 0;
    uint8_t idx_iter = 0;

    /* Write Unit info first.*/
    /* Find where does field id sits in the index list read from MFG data. if not found skip to next.*/
    log_trace("UKDB:: Starting Unit Info write for Module UUID %s.", p_uuid);
    ret = get_fieldid_index(index, FIELDID_UNIT_INFO, idx_count, &idx_iter);
    if (!ret) {
        ret = ukdb_write_unit_info_data(p_uuid, &index[idx_iter], p_uuid,
                                        &mod_count);
        if (!ret) {
            idx++;
        } else {
            goto cleanindex;
        }
        log_trace("UKDB:: Completed Unit Info write for Module UUID %s.",
                  p_uuid);
    } else {
        /* This id will use if the Module DB is getting created.
		 * In this case we might not have any Unit info but surely we would have
		 * a module info and module config to write.
		 */
        mod_count = 1;
        log_trace("UKDB:: Unit Info field for Module UUID %s not found.",
                  p_uuid);
    }

    /*Write unit Config. */
    log_trace("UKDB:: Starting Unit Config write for UUID %s with %d modules.",
              p_uuid, mod_count);
    ret = get_fieldid_index(index, FIELDID_UNIT_CONFIG, idx_count, &idx_iter);
    if (!ret) {
        //TODO: Don't even need UnitCFg here. it can be removed.
        ret = ukdb_write_unit_cfg_data(p_uuid, &index[idx_iter], mod_count);
        if (ret) {
            goto cleanindex;
        }
        idx++;
        log_trace(
            "UKDB:: Completed Unit Config write for Module UUID %s with %d modules.",
            p_uuid, mod_count);
    } else {
        log_trace("UKDB:: Unit Config field for Module UUID %s not found.",
                  p_uuid);
    }

    /* Write Module info */
    log_trace("UKDB:: Starting Module Info write for Module UUID %s.", p_uuid);
    ret = get_fieldid_index(index, FIELDID_MODULE_INFO, idx_count, &idx_iter);
    if (!ret) {
        uint8_t mcfg_iter = 0;
        ret = get_fieldid_index(index, FIELDID_MODULE_CONFIG, idx_count,
                                &mcfg_iter);
        if (!ret) {
            /* Module config and Module info both are added here to DB.*/
            ret = ukdb_write_module_info_data(
                p_uuid, &index[idx_iter], &index[mcfg_iter], &mod_count, &idx);
            if (ret) {
                goto cleanindex;
            }
            idx++;
        } else {
            log_debug("UKDB:: Module Config for Module UUID %s not found.",
                      p_uuid);
        }
    } else {
        log_debug("UKDB:: Module Info for Module UUID %s not found.", p_uuid);
    }

    /* All the remaining fields are just a files which are stored as it is.*/
    uint16_t genfield_id = FIELDID_FACT_CONFIG;
    uint8_t temp_idx = idx;
    while (temp_idx < idx_count) {
        ret = get_fieldid_index(index, genfield_id, idx_count, &idx_iter);
        if (!ret) {
            log_debug("UKDB:: Writing data for field Id 0x%x Module UIID %s.",
                      index[idx_iter].fieldid, p_uuid);
            ret =
                ukdb_write_generic_data(&index[idx_iter], p_uuid, genfield_id);
            if (!ret) {
                idx++;
                log_debug("UKDB:: Wrote data for field Id 0x%x Module UIID %s.",
                          index[idx_iter].fieldid, p_uuid);
            } else {
                goto cleanindex;
            }
        }
        temp_idx++;
        genfield_id++;
    }
    log_debug("UKDB:: UKDB created with total of %d entries.", idx);

cleanindex:
    UBSP_FREE(index);
    return ret;
}

/* Validate the crc matches */
int ukdb_validate_payload(void *p_data, uint32_t crc, uint16_t size) {
    int ret = 0;
    uint32_t calc_crc = crc_32(p_data, size);
    if (calc_crc != crc) {
        ret = ERR_UBSP_CRC_FAILURE;
    }
    log_trace("UKDB:: CRC Expected 0x%x Calculated 0x%x.", crc, calc_crc);
    return ret;
}

/* Read payload */
int ukdb_read_payload(char *p_uuid, void *p_data, uint16_t offset,
                      uint16_t size) {
    int ret = 0;
    if (db_read_block(p_uuid, p_data, offset, size) < size) {
        ret = ERR_UBSP_R_FAIL;
        log_error("Err(%d): UKDB Index read failed from offset 0x%x.", ret,
                  offset);
    }
    log_trace("UKDB read %d bytes of payload from 0x%x offset.", size, offset);
    return ret;
}

/* Read payload for field id.*/
int ukdb_read_payload_for_fieldid(char *p_uuid, void *p_payload, uint16_t fid,
                                  uint16_t *size) {
    int ret = -1;
    uint16_t idx = 0;

    UKDBIdxTuple *idx_data;
    ret = ukdb_search_fieldid(p_uuid, &idx_data, &idx, fid);
    if (ret) {
        log_error("Err(%d): UKDB search error for field id 0x%x.", ret, fid);
        return ret;
    }

    if (p_payload) {
        ret = ukdb_read_payload(p_uuid, p_payload, idx_data->payload_offset,
                                idx_data->payload_size);
        if (ret) {
            log_error("Err(%d): Payload read failure for the field id 0x%x.",
                      ret, fid);
            UBSP_FREE(idx_data);
        }
    }
    *size = idx_data->payload_size;
    /* validate index */
    ret = ukdb_validate_payload(p_payload, idx_data->payload_crc, *size);
    if (ret) {
        log_error("Err(%d): CRC failure for the field id 0x%x.", ret, fid);
        *size = 0;
    }

    UBSP_FREE(idx_data);
    return ret;
}

void ukdp_print_header(UKDBHeader *header) {
    log_trace("*************************************************************");
    log_trace(
        "*********************UKDB Header*******************************");
    log_trace("*************************************************************");
    log_trace("  *	 UKDB Version:                  v%d.%d",
              header->dbversion.major, header->dbversion.minor);
    log_trace("  *	 UKDB Table Offset:             0x%x",
              header->idx_tbl_offset);
    log_trace("  *	 UKDB Index Tuple Size:         %d", header->idx_tpl_size);
    log_trace("  *	 UKDB Tuple max count:          %d",
              header->idx_tpl_max_count);
    log_trace("  *	 UKDB Current Tuple Id:         %d", header->idx_cur_tpl);
    log_trace("  *	 UKDB Capability:               0x%x", header->mod_cap);
    log_trace("  *	 UKDB Mode:                     0x%x", header->mod_mode);
    log_trace("  *	 UKDB Device Owner:             0x%x", header->mod_devown);
    log_trace("*************************************************************");
    log_trace("*************************************************************");
}

void ukdb_print_index_table(UKDBIdxTuple *idx_tbl, uint8_t count) {
    log_trace("*************************************************************");
    log_trace(
        "*********************Index Table*******************************");
    log_trace("*************************************************************");
    for (uint8_t iter = 0; iter < count; iter++) {
        log_trace(
            "*************************************************************");
        log_trace("  *	 Field Id:                      0x%04x ",
                  idx_tbl[iter].fieldid);
        log_trace("  *	 Payload Offset:                0x%04x ",
                  idx_tbl[iter].payload_offset);
        log_trace("  *	 Payload Size:                  0x%04x bytes",
                  idx_tbl[iter].payload_size);
        log_trace("  *	 Payload Version:               v%d.%d",
                  idx_tbl[iter].payload_version.major,
                  idx_tbl[iter].payload_version.minor);
        log_trace("  *	 Payload CRC:                   0x%04x",
                  idx_tbl[iter].payload_crc);
        log_trace("  *	 Field State:                   0x%01x",
                  idx_tbl[iter].state);
        log_trace("  *	 Field Valid:                   0x%01x",
                  idx_tbl[iter].valid);
        log_trace(
            "*************************************************************");
    }
    log_trace("*************************************************************");
}

void ukdb_print_unit_info(UnitInfo *p_uinfo) {
    log_trace("*************************************************************");
    log_trace("*********************Unit Info*******************************");
    log_trace("*************************************************************");
    log_trace("  *	 Unit UUID:                    %s", p_uinfo->uuid);
    log_trace("  *	 Unit Name:                    %s", p_uinfo->name);
    log_trace("  *	 Unit Type:                    0x%01x", p_uinfo->unit);
    log_trace("  *	 Unit Part No.:                %s", p_uinfo->partno);
    log_trace("  *	 Unit Skew:                    %s", p_uinfo->skew);
    log_trace("  *	 Unit MAC:                     %s", p_uinfo->mac);
    log_trace("  *	 Unit SWVer:                   v%d.%d",
              p_uinfo->swver.major, p_uinfo->swver.minor);
    log_trace("  *	 Unit PSWVer:                  v%d.%d",
              p_uinfo->swver.major, p_uinfo->swver.minor);
    log_trace("  *	 Unit Assm_Date:               %s", p_uinfo->assm_date);
    log_trace("  *	 Unit OEM_Name:                %s", p_uinfo->oem_name);
    log_trace("  *	 Unit Module Count:            %d", p_uinfo->mod_count);
    log_trace("*************************************************************");
    log_trace("*************************************************************");
}

void ukdb_print_module_info(ModuleInfo *p_minfo) {
    log_trace("*************************************************************");
    log_trace("*******************Module Info*******************************");
    log_trace("*************************************************************");
    log_trace("  *	  Module UUID:              %s", p_minfo->uuid);
    log_trace("  *	  Module Name:              %s", p_minfo->name);
    log_trace("  *	  Module Type:              0x%01x", p_minfo->module);
    log_trace("  *	  Module Part No.:          %s", p_minfo->partno);
    log_trace("  *	  Module HWVer:             %s", p_minfo->hwver);
    log_trace("  *	  Module MAC:               %s", p_minfo->mac);
    log_trace("  *	  Module SWVer:             v%d.%d", p_minfo->swver.major,
              p_minfo->swver.minor);
    log_trace("  *	  Module PSWVer:            v%d.%d", p_minfo->swver.major,
              p_minfo->swver.minor);
    log_trace("  *	  Module MFG_Date:          %s", p_minfo->mfg_date);
    log_trace("  *	  Module MFG_Name:          %s", p_minfo->mfg_name);
    log_trace("  *	  Module Device Count:      %d", p_minfo->dev_count);
    log_trace("*************************************************************");
    log_trace("*************************************************************");
}

void ukdb_print_dev_gpio_cfg(DevGpioCfg *pdev) {
    if (pdev) {
        log_trace("  *   GPIO Number:               0x%x", pdev->gpio_num);
        log_trace("  *   GPIO Direction:            0x%x", pdev->direction);
    }
}

void ukdb_print_dev_i2c_cfg(DevI2cCfg *pdev) {
    if (pdev) {
        log_trace("  *   I2C Bus:                   0x%x", pdev->bus);
        log_trace("  *   Address:                   0x%x", pdev->add);
    }
}

void ukdb_print_dev_spi_cfg(DevSpiCfg *pdev) {
    if (pdev) {
        log_trace("  *   SPI Bus:                   0x%x", pdev->bus);
        ukdb_print_dev_gpio_cfg(&(pdev->cs));
    }
}

void ukdb_print_dev_uart_cfg(DevUartCfg *pdev) {
    if (pdev) {
        log_trace("  *   Uart Number:               0x%x", pdev->uartno);
    }
}

void ukdb_print_dev(void *dev, DeviceClass class) {
    switch (class) {
    case DEV_CLASS_GPIO: {
        ukdb_print_dev_gpio_cfg(dev);
        break;
    }
    case DEV_CLASS_I2C: {
        ukdb_print_dev_i2c_cfg(dev);
        break;
    }
    case DEV_CLASS_SPI: {
        ukdb_print_dev_spi_cfg(dev);
        break;
    }
    case DEV_CLASS_UART: {
        ukdb_print_dev_uart_cfg(dev);
        break;
    }
    default: {
        log_trace("  *   Invalid device found.");
    }
    }
}

void ukdb_print_unit_cfg(UnitCfg *p_ucfg, uint8_t count) {
    uint8_t iter = 0;
    log_trace("*************************************************************");
    log_trace("*******************Unit Config*******************************");
    for (; iter < count; iter++) {
        log_trace(
            "*************************************************************");
        log_trace("  *	  Module UUID:           %s", p_ucfg[iter].mod_uuid);
        log_trace("  *	  Module Name:           %s", p_ucfg[iter].mod_name);
        log_trace("  *    EEPROM SysFS Name:     %s", p_ucfg[iter].sysfs);
        ukdb_print_dev_i2c_cfg(p_ucfg[iter].eeprom_cfg);
        log_trace(
            "*************************************************************");
    }
}

void ukdb_print_module_cfg(ModuleCfg *p_mcfg, uint8_t count) {
    uint8_t iter = 0;
    log_trace("*************************************************************");
    log_trace("*******************Module Config*****************************");
    for (; iter < count; iter++) {
        log_trace(
            "*************************************************************");
        log_trace("  *	  Device Name:              %s", p_mcfg[iter].dev_name);
        log_trace("  *	  Device Disc:              %s", p_mcfg[iter].dev_disc);
        log_trace("  *	  Device Type:              0x%x",
                  p_mcfg[iter].dev_type);
        log_trace("  *	  Device Class:             0x%x",
                  p_mcfg[iter].dev_class);
        log_trace("  *	  Device SysFile:           %s", p_mcfg[iter].sysfile);
        ukdb_print_dev(p_mcfg[iter].cfg, p_mcfg[iter].dev_class);
        log_trace(
            "*************************************************************");
    }
}

/* Read the payload from UKDB for bsp or for application layer*/
int ukdb_read_payload_from_ukdb(char *p_uuid, void *p_data, uint16_t id,
                                uint16_t *size) {
    int ret = -1;
    switch (id) {
    case FIELDID_UNIT_INFO: {
        ret = ukdb_read_payload_for_fieldid(p_uuid, p_data, FIELDID_FACT_CONFIG,
                                            size);
        break;
    }
    case FIELDID_UNIT_CONFIG: ///this won't work
    {
        ret = ukdb_read_payload_for_fieldid(p_uuid, p_data, FIELDID_UNIT_CONFIG,
                                            size);
        break;
    }
    case FIELDID_MODULE_INFO: {
        ret = ukdb_read_payload_for_fieldid(p_uuid, p_data, FIELDID_FACT_CONFIG,
                                            size);
        break;
    }
    case FIELDID_MODULE_CONFIG: {
        ret = ukdb_read_payload_for_fieldid(p_uuid, p_data,
                                            FIELDID_MODULE_CONFIG, size);
        break;
    }
    case FIELDID_FACT_CONFIG: {
        ret = ukdb_read_payload_for_fieldid(p_uuid, p_data, FIELDID_FACT_CONFIG,
                                            size);
        break;
    }
    case FIELDID_USER_CONFIG: {
        ret = ukdb_read_payload_for_fieldid(p_uuid, p_data, FIELDID_USER_CONFIG,
                                            size);
        break;
    }

    case FIELDID_FACT_CALIB: {
        ret = ukdb_read_payload_for_fieldid(p_uuid, p_data, FIELDID_FACT_CALIB,
                                            size);
        break;
    }

    case FIELDID_USER_CALIB: {
        ret = ukdb_read_payload_for_fieldid(p_uuid, p_data, FIELDID_USER_CALIB,
                                            size);
        break;
    }
    case FIELDID_BS_CERTS: {
        ret = ukdb_read_payload_for_fieldid(p_uuid, p_data, FIELDID_BS_CERTS,
                                            size);
        break;
    }
    case FIELDID_LWM2M_CERTS: {
        ret = ukdb_read_payload_for_fieldid(p_uuid, p_data, FIELDID_LWM2M_CERTS,
                                            size);
        break;
    }
    default: {
        ret = ERR_UBSP_DB_MISSING_FIELD;
        log_error("Err(%d): Invalid Field id supplied by Index entry.", ret);
    }
    }

    if (ret) {
        p_data = NULL;
        ret = ERR_UBSP_R_FAIL;
        log_error("Err(%d): UKDB failed to read info on 0x%x.", ret, id);
    }
    return ret;
}

/* This will read unit info and size of the info.*/
int ukdb_read_unit_info(char *p_uuid, UnitInfo *p_info, uint16_t *size) {
    int ret = -1;
    uint16_t unit_fid = FIELDID_UNIT_INFO;
    uint16_t idx = 0;

    UKDBIdxTuple *idx_data;
    ret = ukdb_search_fieldid(p_uuid, &idx_data, &idx, unit_fid);
    if (ret) {
        log_error("Err(%d): UKDB search error for field id 0x%x.", ret,
                  unit_fid);
        ret = ERR_UBSP_DB_MISSING_UNIT_INFO;
        return ret;
    }

    if (p_info) {
        ret = ukdb_read_payload(p_uuid, p_info, idx_data->payload_offset,
                                idx_data->payload_size);
        if (ret) {
            log_error("Err(%d): Payload read failure for the field id 0x%x.",
                      ret, unit_fid);
            UBSP_FREE(idx_data);
        }
        //p_info = info;
        *size = idx_data->payload_size;
    }

    /* validate index */
    ret = ukdb_validate_payload(p_info, idx_data->payload_crc, *size);
    if (ret) {
        log_error("Err(%d): CRC failure for the field id 0x%x.", ret, unit_fid);
    }
    UBSP_FREE(idx_data);
    return ret;
}

int deserialize_unit_cfg_data(UnitCfg **p_ucfg, char *payload, uint8_t count,
                              uint16_t *size) {
    /* || Unit Info 1 | EEPROM CFG  || Unit Info 2 | EEPROM CFG  || */
    int ret = 0;
    int offset = 0;
    for (int iter = 0; iter < count; iter++) {
        /* Copy Unit Cfg first*/
        memcpy(&(*p_ucfg)[iter], payload + offset, sizeof(UnitCfg));
        /* Create a eeprom cfg and assign reference to eeprom_cfg in UnitCfg */
        offset = offset + sizeof(UnitCfg);
        /* Our Unit config assumes all eeprom are on I2C bus*/
        DevI2cCfg *icfg = malloc(sizeof(DevI2cCfg));
        if (icfg) {
            memcpy(icfg, payload + offset, sizeof(DevI2cCfg));
            (*p_ucfg)[iter].eeprom_cfg = icfg;
        } else {
            ret = ERR_UBSP_MEMORY_EXHAUSTED;
        }
        offset = offset + sizeof(DevI2cCfg);
        /* Size returned to reader of Unit config.*/
        *size = *size + sizeof(UnitCfg);
    }
    return ret;
}

int ukdb_read_unit_cfg(char *p_uuid, UnitCfg *p_ucfg, uint8_t count,
                       uint16_t *size) {
    int ret = -1;
    uint16_t fid = FIELDID_UNIT_CONFIG;
    uint16_t idx = 0;
    char *payload = NULL;
    UKDBIdxTuple *idx_data;
    /* Searching for index */
    ret = ukdb_search_fieldid(p_uuid, &idx_data, &idx, fid);
    if (ret) {
        log_error("Err(%d): UKDB search error for field id 0x%x.", ret, fid);
        ret = ERR_UBSP_DB_MISSING_UNIT_CFG;
        return ret;
    }
    /*Reading payload*/
    payload = malloc(sizeof(char) * idx_data->payload_size);
    if (payload) {
        ret = ukdb_read_payload(p_uuid, payload, idx_data->payload_offset,
                                idx_data->payload_size);
        if (ret) {
            log_error("Err(%d): Payload read failure for the field id 0x%x.",
                      ret, fid);
            UBSP_FREE(idx_data);
        }

        /* validate CRC check*/
        ret = ukdb_validate_payload(payload, idx_data->payload_crc,
                                    idx_data->payload_size);
        if (ret) {
            log_error("Err(%d): CRC failure for the field id 0x%x.", ret,
                      p_ucfg);
            goto cleanup;
        }

        /* Deserialize to Unit Config.*/
        ret = deserialize_unit_cfg_data(&p_ucfg, payload, count, size);
        if (ret) {
            ret = ERR_UBSP_DESERIAL_FAIL;
            log_error("Err(%d): Deserialize failure for Unit Config.", ret);
            goto cleanup;
        }
    }

cleanup:
    UBSP_FREE(payload);
    UBSP_FREE(idx_data);
    return ret;
}

/* This will read module info and size of the info.*/
int ukdb_read_module_info(char *p_uuid, ModuleInfo *p_info, uint16_t *size) {
    //TODO
    int ret = -1;
    uint16_t fid = FIELDID_MODULE_INFO;
    uint16_t idx = 0;

    UKDBIdxTuple *idx_data;
    ret = ukdb_search_fieldid(p_uuid, &idx_data, &idx, fid);
    if (ret) {
        log_error("Err(%d): UKDB search error for field id 0x%x.", ret, fid);
        ret = ERR_UBSP_DB_MISSING_MODULE_INFO;
        return ret;
    }

    if (p_info) {
        ret = ukdb_read_payload(p_uuid, p_info, idx_data->payload_offset,
                                idx_data->payload_size);
        if (ret) {
            log_error("Err(%d): Payload read failure for the field id 0x%x.",
                      ret, fid);
            UBSP_FREE(idx_data);
        }
        *size = idx_data->payload_size;
    }

    /* validate index */
    ret = ukdb_validate_payload(p_info, idx_data->payload_crc, *size);
    if (ret) {
        log_error("Err(%d): CRC failure for the field id 0x%x.", ret, fid);
    }
    UBSP_FREE(idx_data);
    return ret;
}

void *ukdb_deserialize_devices(const char *payload, int offset, uint16_t class,
                               int *size) {
    void *dev = NULL;
    switch (class) {
    case DEV_CLASS_GPIO: {
        DevGpioCfg *cfg = malloc(sizeof(DevGpioCfg));
        if (cfg) {
            memcpy(cfg, payload + offset, sizeof(DevGpioCfg));
        } else {
            cfg = NULL;
        }
        dev = cfg;
        *size = sizeof(DevGpioCfg);
        break;
    }
    case DEV_CLASS_I2C: {
        DevI2cCfg *cfg = malloc(sizeof(DevI2cCfg));
        if (cfg) {
            memcpy(cfg, payload + offset, sizeof(DevI2cCfg));
        } else {
            cfg = NULL;
        }
        dev = cfg;
        *size = sizeof(DevI2cCfg);
        break;
    }
    case DEV_CLASS_SPI: {
        DevSpiCfg *cfg = malloc(sizeof(DevSpiCfg));
        if (cfg) {
            memcpy(cfg, payload + offset, sizeof(DevSpiCfg));
        } else {
            cfg = NULL;
        }
        dev = cfg;
        *size = sizeof(DevSpiCfg);
        break;
    }
    case DEV_CLASS_UART: {
        DevUartCfg *cfg = malloc(sizeof(DevUartCfg));
        if (cfg) {
            memcpy(cfg, payload + offset, sizeof(DevUartCfg));
        } else {
            cfg = NULL;
        }
        dev = cfg;
        *size = sizeof(DevUartCfg);
        break;
    }
    default:
        dev = NULL;
        *size = 0;
        log_error("Err(%d): PARSER:: Unkown device type failed to parse.",
                  ERR_UBSP_INVALID_DEVICE_CFG);
    }
    return dev;
}

int deserialize_module_cfg_data(ModuleCfg **p_mcfg, char *payload,
                                uint8_t count, uint16_t *size) {
    /* || Unit Info 1 | EEPROM CFG  || Unit Info 2 | EEPROM CFG  || */
    int ret = 0;
    int offset = 0;
    for (int iter = 0; iter < count; iter++) {
        /* Copy Module Cfg first*/
        memcpy(&(*p_mcfg)[iter], payload + offset, sizeof(ModuleCfg));
        offset = offset + sizeof(ModuleCfg);
        int cfg_size = 0;
        /* Create a device cfg and assign reference to cfg in ModuleCfg */
        void *cfg = ukdb_deserialize_devices(
            payload, offset, (*p_mcfg)[iter].dev_class, &cfg_size);
        if (cfg) {
            (*p_mcfg)[iter].cfg = cfg;
        } else {
            ret = ERR_UBSP_DESERIAL_FAIL;
            log_error(
                "Err(%d):: UKDB: Deserialization failure for module config.",
                ret);
        }
        offset = offset + cfg_size;
        /* Size returned to reader of Unit config.*/
        *size = *size + sizeof(ModuleCfg);
    }
    return ret;
}

/* This will read module config and count of the module*/
int ukdb_read_module_cfg(char *p_uuid, ModuleCfg *p_cfg, uint8_t count,
                         uint16_t *size) {
    int ret = -1;
    uint16_t fid = FIELDID_MODULE_CONFIG;
    uint16_t idx = 0;
    char *payload = NULL;
    UKDBIdxTuple *idx_data;

    ret = ukdb_search_fieldid(p_uuid, &idx_data, &idx, fid);
    if (ret) {
        log_error("Err(%d): UKDB search error for field id 0x%x.", ret, fid);
        ret = ERR_UBSP_DB_MISSING_MODULE_CFG;
        return ret;
    }

    payload = malloc(sizeof(char) * idx_data->payload_size);
    if (payload) {
        /* Read the DB*/
        ret = ukdb_read_payload(p_uuid, payload, idx_data->payload_offset,
                                idx_data->payload_size);
        if (ret) {
            log_error("Err(%d): Payload read failure for the field id 0x%x.",
                      ret, fid);
            goto cleanup;
        }

        /* validate CRC check */
        ret = ukdb_validate_payload(payload, idx_data->payload_crc,
                                    idx_data->payload_size);
        if (ret) {
            log_error("Err(%d): UKDB:: CRC failure for the field id 0x%x.", ret,
                      fid);
            goto cleanup;
        }

        if (p_cfg) {
            /* Deserialize payload to Module Config */
            ret = deserialize_module_cfg_data(&p_cfg, payload, count, size);
            if (!ret) {
                log_debug(
                    "UKDB:: Read Module Info %d bytes for Module %s  with device count %d.",
                    *size, p_uuid, count);
            } else {
                log_error(
                    "Err(%d): UKDB:: Payload deserialize failure for the field id 0x%x.",
                    ret, fid);
                goto cleanup;
            }
        } else {
            ret = ERR_UBSP_INVALID_POINTER;
            log_error(
                "Err(%d): UKDB:: Invalid payload pointer for the field id 0x%x.",
                ret, fid);
            goto cleanup;
        }
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_error(
            "Err(%d): UKDB:: Memory exhausted while reading payload for the field id 0x%x.",
            ret, fid);
        goto cleanup;
    }

cleanup:
    UBSP_FREE(payload);
    UBSP_FREE(idx_data);
    return ret;
}

/* Read both Unit info and it's config. */
int ukdb_read_unit(char *p_uuid, UnitInfo *p_uinfo, UnitCfg *p_ucfg) {
    int ret = 0;
    uint16_t size = 0;
    ret = ukdb_read_unit_info(p_uuid, p_uinfo, &size);
    if (ret) {
        log_error("Err(%d): Unit Info read failure.", ret);
        return ret;
    }

    ukdb_print_unit_info(p_uinfo);

    ret = ukdb_read_unit_cfg(p_uuid, p_ucfg, p_uinfo->mod_count, &size);
    if (ret) {
        log_error("Err(%d): Unit Config read failure.", ret);
    }
    ukdb_print_unit_cfg(p_ucfg, p_uinfo->mod_count);
    return ret;
}

/* Reads Module info and it's config. */
int ukdb_read_module(char *p_uuid, ModuleInfo *p_minfo) {
    int ret = 0;
    uint16_t size = 0;
    ret = ukdb_read_module_info(p_uuid, p_minfo, &size);
    if (ret) {
        log_error("Err(%d): Module Info read failure.", ret);
        return ret;
    }
    ukdb_print_module_info(p_minfo);

    ret = ukdb_read_module_cfg(p_uuid, p_minfo->module_cfg, p_minfo->dev_count,
                               &size);
    if (ret) {
        log_error("Err(%d): Module Config read failure.", ret);
    }
    ukdb_print_module_cfg(p_minfo->module_cfg, p_minfo->dev_count);

    return ret;
}
/* Read fact config. */
int ukdb_read_fact_config(char *p_uuid, void *data, uint16_t *size) {
    int ret = 0;
    ret = ukdb_read_payload_from_ukdb(p_uuid, data, FIELDID_FACT_CONFIG, size);
    if (ret) {
        log_error("Err(%d) UKDB failed to read info on 0x%x.", ret,
                  FIELDID_FACT_CONFIG);
    }
    log_debug("UKDB:: Fact Config Field Id : %d data read with size %d bytes.",
              FIELDID_FACT_CONFIG, *size);
    return ret;
}

/* Read user config. */
int ukdb_read_user_config(char *p_uuid, void *data, uint16_t *size) {
    int ret = 0;
    ret = ukdb_read_payload_from_ukdb(p_uuid, data, FIELDID_USER_CONFIG, size);
    if (ret) {
        log_error("Err(%d) UKDB failed to read info on 0x%x.", ret,
                  FIELDID_USER_CONFIG);
    }
    log_debug("UKDB:: Fact Config Field Id : %d data read with size %d bytes.",
              FIELDID_USER_CONFIG, *size);
    return ret;
}

/* Read fact calib. */
int ukdb_read_fact_calib(char *p_uuid, void *data, uint16_t *size) {
    int ret = 0;
    ret = ukdb_read_payload_from_ukdb(p_uuid, data, FIELDID_FACT_CALIB, size);
    if (ret) {
        log_error("Err(%d) UKDB failed to read info on 0x%x.", ret,
                  FIELDID_FACT_CALIB);
    }
    log_debug(
        "UKDB:: Fact Calibration Field Id : %d data read with size %d bytes.",
        FIELDID_FACT_CALIB, *size);
    return ret;
}

/* Read user calib. */
int ukdb_read_user_calib(char *p_uuid, void *data, uint16_t *size) {
    int ret = 0;
    ret = ukdb_read_payload_from_ukdb(p_uuid, data, FIELDID_USER_CALIB, size);
    if (ret) {
        log_error("Err(%d) UKDB failed to read info on 0x%x.", ret,
                  FIELDID_USER_CALIB);
    }
    log_debug(
        "UKDB:: User Calibration Field Id : %d data read with size %d bytes.",
        FIELDID_USER_CALIB, *size);
    return ret;
}

/* Read bootstrap certs. */
int ukdb_read_bs_certs(char *p_uuid, void *data, uint16_t *size) {
    int ret = 0;
    ret = ukdb_read_payload_from_ukdb(p_uuid, data, FIELDID_BS_CERTS, size);
    if (ret) {
        log_error("Err(%d) UKDB failed to read info on 0x%x.", ret,
                  FIELDID_BS_CERTS);
    }
    log_debug(
        "UKDB:: Bootstrap certs Field Id : %d data read with size %d bytes.",
        FIELDID_BS_CERTS, *size);
    return ret;
}

/* Read lwm2m certs. */
int ukdb_read_lwm2m_certs(char *p_uuid, void *data, uint16_t *size) {
    int ret = 0;
    ret = ukdb_read_payload_from_ukdb(p_uuid, data, FIELDID_LWM2M_CERTS, size);
    if (ret) {
        log_error("Err(%d) UKDB failed to read info on 0x%x.", ret,
                  FIELDID_LWM2M_CERTS);
    }
    log_debug("UKDB:: Lwm2m certs Field Id : %d data read with size %d bytes.",
              FIELDID_LWM2M_CERTS, *size);
    return ret;
}

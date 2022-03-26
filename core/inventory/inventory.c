/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "inventory.h"

#include "errorcode.h"
#include "ledger.h"
#include "mfg.h"
#include "schema.h"
#include "store.h"
#include "utils/crc32.h"

#include "usys_error.h"
#include "usys_file.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

static int validate_unit_type(UnitType unit) {
    int ret = 0;
    switch (unit) {
    case UNIT_TNODESDR:
    case UNIT_TNODELTE:
    case UNIT_HNODE:
    case UNIT_ANODE:
    case UNIT_PSNODE:
        ret = 1;
        break;
    default:
        ret = 0;
    }

    return ret;
}

static int validate_unit_info(UnitInfo *uInfo) {
    /* TODO: Need to put some logic here bsed on our UUID */
    int ret = 1;
    if (usys_strncmp(uInfo->uuid, "UK", 2)) {
        ret &= 1;
    }

    if (validate_unit_type(uInfo->unit)) {
        ret &= 1;
    }
    return ret;
}

static int validate_unit_cfg(UnitCfg *cfg, char *fname) {
    int ret = 0;
    if (!usys_strcmp(cfg->sysFs, fname)) {
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
        { .id = "UK-1001-COM-1101",
          .count = 3,
          .cfg =
              (UnitCfg[]){
                  { .modUuid = "UK-1001-COM-1101",
                    .modName = "COMv1",
                    .sysFs = "/tmp/sys/bus/i2c/devices/i2c-0/0-0050/eeprom",
                    .eepromCfg = &(DevI2cCfg){ .bus = 0, .add = 0x50ul } },
                  { .modUuid = "UK-2001-LTE-1101",
                    .modName = "LTE",
                    .sysFs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0050/eeprom",
                    .eepromCfg = &(DevI2cCfg){ .bus = 1, .add = 0x50ul } },
                  { .modUuid = "UK-3001-MSK-1101",
                    .modName = "MASK",
                    .sysFs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0051/eeprom",
                    .eepromCfg = &(DevI2cCfg){ .bus = 1, .add = 0x51ul } },
              } },
        { .id = "UK-5001-RFC-1101",
          .count = 2,
          .cfg =
              (UnitCfg[]){
                  { .modUuid = "UK-5001-RFC-1101",
                    .modName = "RF CTRL BOARD",
                    .sysFs = "/tmp/sys/bus/i2c/devices/i2c-0/0-0051/eeprom",
                    .eepromCfg = &(DevI2cCfg){ .bus = 1, .add = 0x50ul } },
                  { .modUuid = "UK-4001-RFA-1101",
                    .modName = "RF BOARD",
                    .sysFs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0052/eeprom",
                    .eepromCfg = &(DevI2cCfg){ .bus = 2, .add = 0x50ul } },
              } }
    };

    UnitCfg *pcfg = NULL;
    DevI2cCfg *i2cCfg = NULL;
    for (int iter = 0; iter < 2; iter++) {
        if (!usys_strcmp(puuid, udata[iter].id)) {
            *count = udata[iter].count;
            pcfg = usys_zmalloc(sizeof(UnitCfg) * (*count));
            if (pcfg) {
                usys_memcpy(&pcfg[0], &udata[iter].cfg[0],
                            sizeof(UnitCfg) * (*count));

                for (int iiter = 0; iiter < *count; iiter++) {
                    if (udata[iter].cfg[iiter].eepromCfg) {
                        i2cCfg = usys_zmalloc(sizeof(DevI2cCfg));
                        if (i2cCfg) {
                            usys_memcpy(i2cCfg,
                                        udata[iter].cfg[iiter].eepromCfg,
                                        sizeof(DevI2cCfg));
                        }
                    }

                    pcfg[iiter].eepromCfg = i2cCfg;
                }

                break;

            } else {
                usys_log_error(
                    "Memory exhausted while getting unit config from Testdata.",
                    usys_error(errno));
            }
        }
    }
    return pcfg;
}

/* Read Master inventory data uses raw read.
 *  As store will be init after this */
int invt_get_master_unit_cfg(UnitCfg *pcfg, char *invtLnkDb) {
    int ret = 0;
    if (!invtLnkDb) {
        return -1;
    }

    /* Read the inventory-db symbolic link file*/
    char *invtDb = usys_file_read_sym_link(invtLnkDb);
    if (!invtDb) {
        ret = ERR_NODED_DB_LNK_MISSING;
        usys_log_error("Symbolic link for the inventory Db is not available."
                       "ErrorCode: %d",
                       ret);
        return ret;
    }

    /* Check if symbolic link file exist */
    if (!usys_file_exist(invtDb)) {
        ret = ERR_NODED_DB_MISSING;
        usys_log_error("Inventory Db for the system is not available at %s.",
                       ret, invtDb);
        usys_free(invtDb);
        return ret;
    }

    UnitInfo *uInfo = usys_zmalloc(sizeof(UnitInfo));
    if (!uInfo) {
        ret = ERR_NODED_MEMORY_EXHAUSTED;
        usys_log_error(" Memory allocation failed. Error %s",
                       usys_error(errno));
        usys_free(invtDb);
        return ret;
    }

    /* Read Unit Info */
    if (usys_file_raw_read(invtDb, uInfo, SCH_UNIT_INFO_OFFSET,
                           sizeof(UnitInfo)) == sizeof(UnitInfo)) {
        /* Validate Unit Info */
        if (validate_unit_info(uInfo)) {
            usys_log_debug("Unit UIID %s Name %s detected.", uInfo->uuid,
                           uInfo->name);

            /* Read first cfg which belong to master */
            int sz = sizeof(UnitCfg) + sizeof(DevI2cCfg);
            void *cfg = usys_zmalloc(sz);
            if (cfg) {
                /* Read Unit Config */
                if (usys_file_raw_read(invtDb, cfg, SCH_UNIT_CONFIG_OFFSET,
                                       sz) == sz) {
                    if (validate_unit_cfg(cfg, invtDb)) {
                        usys_memcpy(pcfg, cfg, sizeof(UnitCfg));
                        /* Read EEPROM (Inventory DB) config for Module */
                        DevI2cCfg *icfg = usys_zmalloc(sizeof(DevI2cCfg));
                        if (icfg) {
                            usys_memcpy(icfg, (cfg + sizeof(UnitCfg)),
                                        sizeof(DevI2cCfg));
                        } else {
                            icfg = NULL;
                        }
                        pcfg->eepromCfg = icfg;

                    } else {
                        ret = ERR_NODED_INVALID_UNIT_CFG;
                        usys_log_error("Err(%d): UKDB:: Invalid Unit Config.",
                                       ret);
                    }
                } else {
                    ret = ERR_NODED_READ_UNIT_CFG;
                    usys_log_error("Unable to read Unit Config. Error Code: %d",
                                   ret);
                }
                usys_free(cfg);
            }
        } else {
            ret = ERR_NODED_INVALID_UNIT_INFO;
            usys_log_error("Invalid Unit Info. Error Code: %d.", ret);
        }

    } else {
        ret = ERR_NODED_READ_UNIT_INFO;
        usys_log_error("Unable to read Unit Info. Error Code: %d.", ret);
    }

    usys_free(uInfo);
    usys_free(invtDb);
    return ret;
}

/* Validate Magic word */
int invt_validating_magic_word(char *pUuid) {
    int ret = 0;
    SchemaMagicWord mw;

    /* Read Magic word */
    if (store_read_block(pUuid, &mw, SCH_MAGIC_WORD_OFFSET,
                         sizeof(SchemaMagicWord)) != sizeof(SchemaMagicWord)) {
        ret = ERR_NODED_R_FAIL;
        usys_log_error("Inventory magic word read error. Error Code: %d", ret);
        return ret;
    }

    /* Validating Magic word */
    if (mw.magicWord == SCH_MAGIC_WORD) {
        usys_log_debug("Inventory DB Magic Word validation pass for module %s.",
                       pUuid);
    } else {
        ret = ERR_NODED_MW_ERR;
        usys_log_error("Inventory DB MagicWord validation failed for module %s."
                       " Error Code: %d",
                       ret, pUuid);
    }

    return ret;
}

/* Write magic word */
int invt_write_magic_word(char *pUuid) {
    int ret = 0;
    SchemaMagicWord mw;
    mw.magicWord = SCH_MAGIC_WORD;
    mw.resv1 = SCH_DEFVAL;
    mw.resv2 = SCH_DEFVAL;

    /* Store write magic word */
    if (store_write_block(pUuid, &mw, SCH_MAGIC_WORD_OFFSET,
                          sizeof(SchemaMagicWord)) != sizeof(SchemaMagicWord)) {
        ret = ERR_NODED_WR_FAIL;
        usys_log_error("Store write magic word. Error Code: %d", ret);
        return ret;
    }

    if (invt_validating_magic_word(pUuid)) {
        ret = ERR_NODED_MW_ERR;
    }
    return ret;
}

/* Write Schema header*/
int invt_write_header(char *pUuid, SchemaHeader *pHeader) {
    int ret = 0;
    if (store_write_block(pUuid, pHeader, SCH_HEADER_OFFSET, SCH_HEADER_SIZE) !=
        SCH_HEADER_SIZE) {
        ret = ERR_NODED_WR_FAIL;
        usys_log_error("Store Header write of %d bytes failed at "
                       "offset 0x%x. Error Code : %d",
                       SCH_HEADER_SIZE, SCH_HEADER_OFFSET, ret);

    } else {
        usys_log_debug("Store Header write of %d bytes completed at "
                       "offset 0x%x.",
                       SCH_HEADER_SIZE, SCH_HEADER_OFFSET);
    }
    return ret;
}

/* Read UK DB header*/
int invt_read_header(char *pUuid, SchemaHeader *pHeader) {
    int ret = 0;
    if (store_read_block(pUuid, pHeader, SCH_HEADER_OFFSET, SCH_HEADER_SIZE) <
        SCH_HEADER_SIZE) {
        ret = ERR_NODED_R_FAIL;
        usys_log_error("Store Header read of %d bytes from "
                       "offset 0x%x failed. Error Code: %d",
                       SCH_HEADER_SIZE, SCH_HEADER_OFFSET, ret);
    } else {
        usys_log_debug("Store Header read of %d bytes completed from "
                       "offset 0x%x.",
                       SCH_HEADER_SIZE, SCH_HEADER_OFFSET);
        invt_print_header(pHeader);
    }
    return ret;
}

/* Read Store version*/
int invt_read_schema_version(char *pUuid, Version *ver) {
    int ret = 0;

    if (store_read_block(pUuid, ver, SCH_HEADER_DBVER_OFFSET, sizeof(Version)) <
        sizeof(Version)) {
        ret = ERR_NODED_R_FAIL;
        usys_log_error("Schema version read of %d bytes "
                       "from offset 0x%x failed.",
                       ret, sizeof(Version), SCH_HEADER_DBVER_OFFSET);
    } else {
        usys_log_debug("Schema Version read of v%d.%d.", ver->major,
                       ver->minor);
    }

    return ret;
}

/* Update Schema version*/
int invt_update_schema_version(char *pUuid, Version ver) {
    int ret = 0;

    if (store_write_block(pUuid, &ver, SCH_HEADER_DBVER_OFFSET,
                          sizeof(Version)) != sizeof(Version)) {
        ret = ERR_NODED_R_FAIL;
        usys_log_error("Schema version write of %d bytes "
                       "to offset 0x%x failed. Error Code: %d",
                       sizeof(Version), SCH_HEADER_DBVER_OFFSET, ret);
    } else {
        usys_log_debug("Schema Version updated to v%d.%d.", ver.major,
                       ver.minor);
    }

    return ret;
}

/* Update Index table current index count. */
int invt_update_current_idx_count(char *pUuid, uint16_t *idxCount) {
    int ret = 0;

    if (store_write_number(pUuid, idxCount, SCH_IDX_CUR_TPL_COUNT_OFFSET, 1,
                           SCH_IDX_COUNT_SIZE)) {
        ret = ERR_NODED_WR_FAIL;
        usys_log_error("Store update to index count failed.Error Code: %d",
                       ret);
    }

    return ret;
}

/*Update Index table current index count. */
int invt_read_current_idx_count(char *pUuid, uint16_t *idxCount) {
    int ret = 0;

    if (store_read_number(pUuid, idxCount, SCH_IDX_CUR_TPL_COUNT_OFFSET, 1,
                          SCH_IDX_COUNT_SIZE)) {
        ret = ERR_NODED_R_FAIL;
        usys_log_error("Store read for index count failed.Error Code: %d", ret);
    }

    if (!(VALIDATE_INDEX_COUNT(*idxCount))) {
        ret = ERR_NODED_VALIDATION_FAILURE;
        usys_log_error("Inventory validation for index count (%d) failed."
                       "Error Code: %d",
                       *idxCount, ret);
    }

    return ret;
}

/* Erasing the Store.*/
int invt_erase_db(char *pUuid) {
    int ret = 0;

    if (store_erase_block(pUuid, SCH_START_OFFSET, SCH_END_OFFSET) !=
        SCH_END_OFFSET) {
        ret = ERR_NODED_WR_FAIL;
        usys_log_error("Store erase failed. Error Code: %d", ret);
    }

    return ret;
}

/* Erase the last index from the Index table.*/
int invt_erase_idx(char *pUuid) {
    /* Erase is done in LIFO manner */
    int ret = 0;

    /* Read available index count. */
    uint16_t idxCount = 0;
    if (invt_read_current_idx_count(pUuid, &idxCount)) {
        return ERR_NODED_R_FAIL;
    }

    int offset = SCH_IDX_TABLE_OFFSET + (idxCount * SCH_IDX_TPL_SIZE);
    if (store_erase_block(pUuid, offset, SCH_IDX_TPL_SIZE) ==
        SCH_IDX_TPL_SIZE) {
        return ERR_NODED_WR_FAIL;
    }

    /* Update Index count. */
    idxCount--;
    if (invt_update_current_idx_count(pUuid, &idxCount)) {
        ret = ERR_NODED_WR_FAIL;
    }
    usys_log_debug("Index deleted from index id 0x%x at offset 0x%x.",
                   idxCount + 1, offset);

    return ret;
}

/* Update the index at particular location */
int invt_update_index_at_id(char *pUuid, SchemaIdxTuple *pData, uint16_t idx) {
    int ret = 0;
    int offset = SCH_IDX_TABLE_OFFSET + (idx * SCH_IDX_TPL_SIZE);

    if (store_write_block(pUuid, pData, offset, SCH_IDX_TPL_SIZE) !=
        SCH_IDX_TPL_SIZE) {
        ret = ERR_NODED_WR_FAIL;
        usys_log_error("Store Index update failed for index id 0x%x at"
                       "offset 0x%x. Error code:  %d",
                       idx, idx * SCH_IDX_COUNT_SIZE, ret);
    }

    usys_log_debug("Index updated for index id 0x%x at offset 0x%x.", idx,
                   idx * SCH_IDX_COUNT_SIZE);
    return ret;
}

/* Search for the field id with valid data.*/
int invt_search_field_id(char *pUuid, SchemaIdxTuple **idxData, uint16_t *idx,
                         uint16_t fieldId) {
    int ret = ERR_NODED_DB_MISSING_FIELD;

    /* Read available index count. */
    uint16_t idxCount = 0;
    if (invt_read_current_idx_count(pUuid, &idxCount)) {
        return ERR_NODED_R_FAIL;
    }

    *idx = 0xFFFF;
    uint16_t iter = 0;
    SchemaIdxTuple *index =
        (SchemaIdxTuple *)usys_zmalloc(sizeof(SchemaIdxTuple));
    if (index) {
        for (; iter < idxCount; iter++) {
            invt_read_idx(pUuid, index, iter);

            /* Only the fieldId and valid state matters */
            if ((index->fieldId == fieldId)) {
                /* Mark ret as invalid */
                ret = ERR_NODED_INVALID_FIELD;
                if (index->valid) {
                    *idx = iter;
                    ret = 0;
                    *idxData = index;
                    break;
                }
            }
        }

        if (ret) {
            usys_free(index);
            index = NULL;
        }
    }
    return ret;
}

/* Update the index for a particular fieldId */
int invt_update_index_for_field_id(char *pUuid, void *pData, uint16_t fieldId) {
    int ret = 0;
    uint16_t idx = 0xFF;
    SchemaIdxTuple *idxData;
    /* Search fieldID */
    if (invt_search_field_id(pUuid, &idxData, &idx, fieldId)) {
        ret = ERR_NODED_DB_MISSING_FIELD;
        usys_log_error("No entry in UKDB with field Id 0x%x. Error Code: %d",
                       fieldId, ret);
        return ret;
    }

    /* Store data */
    if (store_write_block(pUuid, pData, idxData->payloadOffset,
                          idxData->payloadSize) != SCH_IDX_TPL_SIZE) {
        ret = ERR_NODED_WR_FAIL;
        usys_log_error("Index write failed for index id 0x%x field id 0x%x at "
                       "offset 0x%x.Error Code: %d",
                       idx, idxData->fieldId, idx * SCH_IDX_COUNT_SIZE, ret);
    }

    usys_log_debug(
        "Index updated for index id 0x%x field Id 0x%x at offset 0x%x.", idx,
        idxData->fieldId, idx * SCH_IDX_COUNT_SIZE);
    return ret;
}

/* Write index to next available index slot in store */
int invt_write_idx(char *pUuid, SchemaIdxTuple *pData) {
    int ret = 0;

    /* Read available index count. */
    uint16_t idxCount = 0;
    ret = invt_read_current_idx_count(pUuid, &idxCount);
    if (ret) {
        return ERR_NODED_R_FAIL;
    }
    usys_log_trace("Inventory Current Index in Header is %d.", idxCount);

    /*Write the index. */
    int offset = SCH_IDX_TABLE_OFFSET + idxCount * SCH_IDX_TPL_SIZE;
    if (store_write_block(pUuid, pData, offset, SCH_IDX_TPL_SIZE) !=
        SCH_IDX_TPL_SIZE) {
        return ERR_NODED_WR_FAIL;
    }

    /* Update Index count. */
    idxCount++;
    if (invt_update_current_idx_count(pUuid, &idxCount)) {
        ret = ERR_NODED_WR_FAIL;
    }

    usys_log_trace("Inventory written index 0x%x at 0x%x offset.", idxCount,
                   offset);
    return ret;
}

/* Read the index at nth location from the Index table of the UKDB*/
int invt_read_idx(char *pUuid, SchemaIdxTuple *pData, uint16_t idx) {
    int ret = 0;
    int offset = SCH_IDX_TABLE_OFFSET + (idx * SCH_IDX_TPL_SIZE);

    if (store_read_block(pUuid, pData, offset, SCH_IDX_TPL_SIZE) <
        SCH_IDX_TPL_SIZE) {
        ret = ERR_NODED_R_FAIL;
        usys_log_error("Store Index read failed from id 0x%x. Error Code: %d",
                       ret, idx);
    }

    return ret;
}

/* Erase the payload data for field id and mark the index as invalid. */
int invt_erase_payload(char *pUuid, uint16_t fieldId) {
    int ret = 0;
    uint16_t idx = 0x00;
    SchemaIdxTuple *idxData;

    if (invt_search_field_id(pUuid, &idxData, &idx, fieldId)) {
        ret = ERR_NODED_DB_MISSING_FIELD;
        usys_log_warn("No entry in UKDB with field Id 0x%x.Error Code: %d", ret,
                      fieldId);
        return ret;
    }

    /*erase is basically writing 0xFF to database. */
    if (store_erase_block(pUuid, idxData->payloadOffset,
                          idxData->payloadSize) == SCH_IDX_TPL_SIZE) {
        return ERR_NODED_WR_FAIL;
    }

    usys_log_debug("Store payload erased for field Id 0x%x at offset 0x%x.",
                   idxData->fieldId, idxData->payloadOffset);

    idxData->state = SCH_FEAT_DISABLED;
    idxData->valid = USYS_FALSE;

    usys_log_debug("Inventory Marking index %d for UKDB as invalid.", idx);

    /* Update index */
    if (invt_update_index_at_id(pUuid, idxData, idx)) {
        ret = ERR_NODED_WR_FAIL;
    }

    return ret;
}

/* Write module specific payload */
/* TODO :: Each Module need to write info to specific device*/
/* For now it might be just a simple file.*/
int invt_write_module_payload(char *pUuid, void *pData, uint16_t offset,
                              uint16_t size) {
    int ret = 0;

    //TODO: Check if required otherwise delete and use above fxn.
    if (store_write_block(pUuid, pData, offset, size) != size) {
        usys_log_error(
            "Inventory Module payload write failed at offset 0x%x of "
            "size 0x%x.",
            offset, size);
        ret = ERR_NODED_WR_FAIL;
    }

    usys_log_trace("Inventory Wrote %d bytes of module payload to 0x%x"
                   " offset.",
                   size, offset);
    return ret;
}

/* Update the payload data with the given field id. It also Update the index
 * in Index table for field Id. */
int invt_update_payload(char *pUuid, void *pData, uint16_t fieldId,
                        uint16_t size) {
    int ret = 0;
    /*Find index first.
        Compare the field id's from 0th index to last aailable index.
        If field Id matches read payload offset calculate crc and write payload.
     */
    uint16_t idx = 0x0;
    SchemaIdxTuple *idxData;
    if (invt_search_field_id(pUuid, &idxData, &idx, fieldId)) {
        usys_log_warn("Inventory No entry in UKDB with fieldId 0x%x.", fieldId);
        return ERR_NODED_DB_MISSING_FIELD;
    }

    /*Write payload for Index table entries*/
    if (invt_write_module_payload(pUuid, pData, idxData->payloadOffset, size)) {
        ret = ERR_NODED_WR_FAIL;
        usys_log_error(
            "Inventory payload write failed for index id 0x%x fieldId 0x%x at"
            " offset 0x%x. Error Code: %d",
            idx, idxData->fieldId, idxData->payloadOffset, ret);
        return ret;
    }

    /* Calculate CRC */
    uint32_t crc = crc_32(pData, size);

    /* Update Index */
    idxData->payloadCrc = crc;
    idxData->payloadSize = size;
    idxData->valid = USYS_TRUE;
    if (invt_update_index_at_id(pUuid, idxData, idx)) {
        ret = ERR_NODED_WR_FAIL;
        usys_log_error(
            "Inventory Index write failed for index id 0x%x field Id "
            "0x%x at offset 0x%x. Error Code: %d",
            idx, idxData->fieldId, idx * SCH_IDX_COUNT_SIZE, ret);
        return ret;
    }

    usys_log_debug("Inventory Payload updated for index id 0x%x field Id "
                   "0x%x at offset 0x%x with crc %u for %d bytes.",
                   idx, idxData->fieldId, idxData->payloadOffset,
                   idxData->payloadCrc, idxData->payloadSize);

    return ret;
}

/* Allocate memory for unit config */
UnitCfg *invt_alloc_unit_cfg(uint8_t count) {
    UnitCfg *cfg = usys_zmalloc(sizeof(UnitCfg) * count);
    return cfg;
}

/* Free allocated memory for unit config */
void invt_free_unit_cfg(UnitCfg *cfg, uint8_t count) {
    if (cfg) {
        for (int iter = 0; iter < count; iter++) {
            usys_free(cfg[iter].eepromCfg);
            cfg[iter].eepromCfg = NULL;
        }
        usys_free(cfg);
        cfg = NULL;
    }
}

/* Serialize unit config */
char *serialize_unitcfg_payload(UnitCfg *ucfg, uint8_t count, uint16_t *size) {
    int offset = 0;
    char *data = NULL;
    *size = (sizeof(UnitCfg) + sizeof(DevI2cCfg)) * count;

    data = usys_zmalloc(*size);
    if (data) {
        for (int iter = 0; iter < count; iter++) {
            usys_memcpy(data + offset, &ucfg[iter], sizeof(UnitCfg));
            offset = offset + sizeof(UnitCfg);
            usys_memcpy(data + offset, ucfg[iter].eepromCfg, sizeof(DevI2cCfg));
            offset = offset + sizeof(DevI2cCfg);
        }
    } else {
        data = NULL;
    }

    return data;
}

/* Write Unit Config and update the index to index table of the DB */
int invt_write_unit_cfg_data(char *pUuid, SchemaIdxTuple *index,
                             uint8_t count) {
    int ret = 0;
    /* Unit Config */
    UnitCfg *ucfg;
    uint16_t size = 0;
    char *payload = NULL;
    ret = mfg_fetch_unit_cfg(&ucfg, pUuid, &size, count);
    if (ucfg) {
        /*Write payload for Index table entries*/
        payload = serialize_unitcfg_payload(ucfg, count, &size);
        if (payload) {
            if (invt_write_module_payload(pUuid, payload, index->payloadOffset,
                                   size)) {
                /* Need to revert back index */
                invt_erase_idx(pUuid);
                ret = ERR_NODED_WR_FAIL;
                goto cleanup;
            }

            /* Update payload size */
            index->payloadSize = size;

            /* Add CRC */
            uint32_t crcVal =
                crc_32((const unsigned char *)payload, index->payloadSize);
            index->payloadCrc = crcVal;

            usys_log_debug("Inventory Calculated CRCfor filed Id %d payload "
                           "is 0x%x",
                           index->fieldId, crcVal);
            usys_log_debug("Inventory Payload added for field Id 0x%x",
                           index->fieldId);

        } else {
            ret = ERR_NODED_MEMORY_EXHAUSTED;
            usys_log_error("Error while allocating memory for Unit Cfg."
                           " Error: %s",
                           usys_error(errno));
            goto cleanup;
        }

        /* Write Index table entries */
        if (invt_write_idx(pUuid, index)) {
            ret = ERR_NODED_WR_FAIL;
            goto cleanup;
        }

        usys_log_debug("Inventory Index added to the Index table for"
                       "field Id 0x%x",
                       index->fieldId);
    } else {
        usys_log_error("Failed to read Unit config from Mfg data. "
                       "Error code: %d",
                       ret);
    }
cleanup:
    /* Only free the ucfg not the module config.*/
    usys_free(payload);
    payload = NULL;
    invt_free_unit_cfg(ucfg, count);
    return ret;
}

/* Write Unit Info and update the index to index table of the DB */
int invt_write_unit_info_data(char *p1Uuid, SchemaIdxTuple *index, char *pUuid,
                              uint8_t *count) {
    int ret = 0;
    UnitInfo *uInfo;
    uint16_t size = 0;
    ret = mfg_fetch_unit_info(&uInfo, pUuid, &size);
    if (!ret) {
        /* Unit Info */
        if (uInfo) {
            /*Write payload for Index table entries*/
            if (invt_write_module_payload(p1Uuid, uInfo, index->payloadOffset,
                                   sizeof(UnitInfo))) {
                /* Need to revert back index */
                invt_erase_idx(p1Uuid);
                ret = ERR_NODED_WR_FAIL;
                goto cleanup;
            }
            usys_log_debug("Inventory Unit Info added for unit Id %s.",
                           uInfo->uuid);

            /* CRC*/
            uint32_t crcVal =
                crc_32((const unsigned char *)uInfo, sizeof(UnitInfo));
            usys_log_debug("Inventory Calculated crc for unit id %s is 0x%x",
                           uInfo->uuid, crcVal);
            index->payloadCrc = crcVal;

            /* Write Index table entries */
            if (invt_write_idx(p1Uuid, index)) {
                ret = ERR_NODED_WR_FAIL;
                goto cleanup;
            }

            /*Assign module count.*/
            *count = uInfo->modCount;
            if (!VALIDATE_MODULE_COUNT(*count)) {
                ret = ERR_NODED_VALIDATION_FAILURE;
                usys_log_error("Inventory Validation for module %d failed.",
                               ret, *count);
                goto cleanup;
            }

            usys_log_debug("Inventory Index added to the Index table for "
                           "field Id 0x%x and Unit Id %s.",
                           index->fieldId, uInfo->uuid);
        }
    } else {
        usys_log_error("Failed to read Unit info from Mfg data."
                       "Error Code: %d",
                       ret);
    }

cleanup:
    usys_free(uInfo);
    uInfo = NULL;
    return ret;
}

ModuleCfg *invt_alloc_module_cfg(uint8_t count) {
    ModuleCfg *cfg = usys_zmalloc(sizeof(ModuleCfg) * count);
    if (cfg) {
        memset(cfg, '\0', sizeof(ModuleCfg) * count);
    }
    return cfg;
}

/* Free the memory used by moduleCfg */
void invt_free_module_cfg(ModuleCfg *cfg, uint8_t count) {
    if (cfg) {
        for (int iter = 0; iter < count; iter++) {
            void *devCfg = cfg[iter].cfg;
            usys_free(devCfg);
            devCfg = NULL;
        }
        usys_free(cfg);
        cfg = NULL;
    }
}

/* Serialize the Module Cfg data */
char *serialize_module_config_data(ModuleCfg *mcfg, uint8_t count,
                                   uint16_t *size) {
    int offset = 0;
    char *data = NULL;
    int psize = 0;

    /* Calculate memory required for HW attributes  */
    for (int iter = 0; iter < count; iter++) {
        uint16_t cfgSize = 0;
        SIZE_OF_DEVICE_CFG(cfgSize, mcfg[iter].devClass);
        psize += cfgSize;
    }

    /* Allocate memory for Module config
     * Sizeof Module config = sizeof(Module Config for each module)
     *                       + sizeof(HWATTR for each module config)*/
    psize = (sizeof(ModuleCfg) * count) + psize;
    data = usys_zmalloc(psize);
    if (data) {
        for (int iter = 0; iter < count; iter++) {
            usys_memcpy(data + offset, &mcfg[iter], sizeof(ModuleCfg));
            offset = offset + sizeof(ModuleCfg);
            uint16_t cfg_size = 0;
            SIZE_OF_DEVICE_CFG(cfg_size, mcfg[iter].devClass);
            usys_memcpy(data + offset, mcfg[iter].cfg, cfg_size);
            offset = offset + cfg_size;
        }

        *size = psize;

    } else {
        data = NULL;
    }

    return data;
}

/* Write Module Config and update the index to index table of the DB */
int invt_write_module_cfg_data(char *pUuid, ModuleInfo *minfo,
                               SchemaIdxTuple *cfgIndex) {
    int ret = 0;
    char *payload = NULL;
    uint16_t size = 0;

    /* Serialize the module config data.*/
    payload =
        serialize_module_config_data(minfo->modCfg, minfo->devCount, &size);
    if (payload) {
        /*Write Module Cfg for the module*/
        if (invt_write_module_payload(pUuid, payload, cfgIndex->payloadOffset,
                                      size)) {
            /* Need to revert back index */
            invt_erase_idx(pUuid);
            ret = ERR_NODED_WR_FAIL;
            goto cleanup;
        }

        cfgIndex->payloadSize = size;
        usys_log_debug("Inventory Module Config added for module Id %s",
                       minfo->uuid);
        uint32_t crcVal = crc_32((const unsigned char *)payload, size);
        cfgIndex->payloadCrc = crcVal;
        usys_log_debug("Inventory Calculated CRC for module Id %s is 0x%x",
                       minfo->uuid, crcVal);
    } else {
        ret = ERR_NODED_MEMORY_EXHAUSTED;
        usys_log_error("Memory error while serializing Module config."
                       "Error: %s",
                       usys_error(errno));
        goto cleanup;
    }

    /* Write Index table entries */
    if (invt_write_idx(pUuid, cfgIndex)) {
        ret = ERR_NODED_WR_FAIL;
        goto cleanup;
    }
    usys_log_debug("Inventory Index added to the Index table for "
                   "field Id 0x%x and module Id %s",
                   cfgIndex->fieldId, minfo->uuid);

cleanup:
    usys_free(payload);
    payload = NULL;
    invt_free_module_cfg(minfo->modCfg, minfo->devCount);
    return ret;
}

/* Write module info to UKDB*/
int invt_write_module_info_data(char *pUuid, SchemaIdxTuple *infoIndex,
                                SchemaIdxTuple *cfgIndex, uint8_t *mcount,
                                uint8_t *idx) {
    int ret = 0;
    ModuleInfo *minfo;
    uint16_t size = 0;

    ret = mfg_fetch_module_info_by_uuid(&minfo, pUuid, &size, *mcount);
    if (ret) {
        usys_log_error("Failed to read Unit info from Mfg data.");
        return (-1);
    }

    /* Unit Info */
    if (minfo) {
        /*Write payload for Index table entries*/
        if (invt_write_module_payload(pUuid, minfo, infoIndex->payloadOffset,
                                      sizeof(ModuleInfo))) {
            /* Need to revert back index */
            invt_erase_idx(pUuid);
            ret = ERR_NODED_WR_FAIL;
            goto cleanmid;
        }
        usys_log_debug("Inventory Added Module Info for module Id %s"
                       " with %d devices.",
                       pUuid, minfo->devCount);

        uint32_t crcVal =
            crc_32((const unsigned char *)minfo, sizeof(ModuleInfo));
        infoIndex->payloadCrc = crcVal;
        usys_log_debug("Inventory Calculated CRC32 for Module %s is 0x%x",
                       pUuid, crcVal);

        /* Write Index table entries */
        if (invt_write_idx(pUuid, infoIndex)) {
            ret = ERR_NODED_WR_FAIL;
            goto cleanmid;
        }

        //Updating write index count.
        (*idx)++;
        usys_log_debug("Inventory Index added to the Index table for "
                       "field Id 0x%x module Id %s",
                       infoIndex->fieldId, pUuid);

        usys_log_debug("Inventory Adding Module Config now for"
                       " module Id %s",
                       pUuid);

        /* Write Module Cfg as well Module info contains the info
         *  for Module cfg also.*/
        //TODO : lot better option will open after invt_create is complete.
        ret = invt_write_module_cfg_data(pUuid, minfo, cfgIndex);
    }

cleanmid:
    usys_free(minfo);
    minfo = NULL;

    return ret;
}

/* Search for the field Id in the Index Table read from MFG data.*/
int invt_get_field_id_idx(SchemaIdxTuple *index, uint16_t fid, uint8_t idxCount,
                       uint8_t *idxIter) {
    int ret = -1;
    uint8_t iter = 0;
    for (; iter < idxCount; iter++) {
        if (index[iter].fieldId == fid) {
            *idxIter = iter;
            usys_log_trace("Inventory Field Id 0x%04x found at MFG data "
                           "location mfg_index_table[%d].",
                           fid, iter);
            ret = 0;
            break;
        }
    }
    return ret;
}

/* Set up environment for DB creation process.*/
int invt_mfg_init(void *data) {
    int ret = 0;
    ret = mfg_init(data);
    return ret;
}

/* Write generic payload to db like some files.*/
int invt_write_generic_data(SchemaIdxTuple *index, char *pUuid, uint16_t fid) {
    int ret = 0;
    uint16_t size = index->payloadSize;
    char *payload;
    ret = mfg_fetch_payload_from_mfg_data((void *)&payload, pUuid, &size, fid);
    if (payload) {
        /*Write payload for Index table entries*/
        if (invt_write_module_payload(pUuid, payload, index->payloadOffset,
                               index->payloadSize)) {
            /* Need to revert back index */
            invt_erase_idx(pUuid);
            ret = ERR_NODED_WR_FAIL;
            goto cleanpayload;
        }

        usys_log_debug("Inventory Added payload for field Id 0x%x.",
                       index->fieldId);
        uint32_t crcVal =
            crc_32((const unsigned char *)payload, index->payloadSize);
        index->payloadCrc = crcVal;
        usys_log_debug("Inventory Calculated CRC32 for field Id %d is 0x%x",
                       index->fieldId, crcVal);
        /* Write Index table entries */
        if (invt_write_idx(pUuid, index)) {
            ret = ERR_NODED_WR_FAIL;
            goto cleanpayload;
        }
        usys_log_debug("Inventory Index added to the Index table at id "
                       "0x%x for field Id 0x%x",
                       index->fieldId, index->fieldId);
    } else {
        ret = ERR_NODED_MEMORY_EXHAUSTED;
        usys_log_debug("Memory error while reading "
                       "Field Id 0x%x for UUID %s. Error %s",
                       index->fieldId, pUuid, usys_error(errno));
        return ret;
    }

cleanpayload:
    usys_free(payload);
    payload = NULL;

    return ret;
}

/* Registering  modules in store */
int invt_register_module(UnitCfg *cfg) {
    int ret = 0;

    ret = store_register_module(cfg);
    if (ret) {
        usys_log_error("Inventory registering module uuid %s failed.", ret,
                       cfg->modUuid);
    }

    return ret;
}

/* deregister module in store */
int invt_deregister_module(char *puuid) {
    int ret = 0;
    ret = store_deregister_module(puuid);
    if (ret) {
        usys_log_error("Inventory deregistering module uuid %s failed.", ret,
                       puuid);
    }
    return ret;
}

/* UKDB registering Modules in DB so that these can be accessed later on.*/
int invt_register_modules(char *pUuid, RegisterDeviceCB registerDev) {
    int ret = 0;
    uint16_t size = 0;
    uint8_t modCount = 0;
    UnitCfg *ucfg;
    char unitUuid[UUID_LENGTH] = { '\0' };
    usys_log_debug("Inventory read Unit Info from module %s for module "
                   "registration process.",
                   pUuid);
    UnitInfo *uInfo = (UnitInfo *)usys_zmalloc(sizeof(UnitInfo));
    if (uInfo) {
        /* Read unit info */
        ret = invt_read_unit_info(pUuid, uInfo, &size);
        if (!ret) {
            invt_print_unit_info(uInfo);
            usys_memcpy(unitUuid, uInfo->uuid, usys_strlen(uInfo->uuid));
            modCount = uInfo->modCount;
            if (uInfo) {
                usys_free(uInfo);
                uInfo = NULL;
            }
            usys_log_debug("Inventory read Unit Config for %s for "
                           "registration process.",
                           unitUuid);

            ucfg = (UnitCfg *)invt_alloc_unit_cfg(modCount);
            if (ucfg) {
                /* Read unit config */
                ret = invt_read_unit_cfg(pUuid, ucfg, modCount, &size);
                if (!ret) {
                    invt_print_unit_cfg(ucfg, modCount);
                    usys_log_debug("Inventory Registering %d module "
                                   "for Unit %s.",
                                   modCount, unitUuid);
                    for (uint8_t iter = 0; iter < modCount; iter++) {
                        /* Register module*/
                        ret = invt_register_module(&ucfg[iter]);
                        if (!ret) {
                            /* Validate DB */
                            ret =
                                invt_validating_magic_word(ucfg[iter].modUuid);
                            if (ret) {
                                usys_log_warn(
                                    "Inventory No valid database found for module"
                                    " UUID %s Name %s. Moving on to next module.",
                                    ucfg[iter].modUuid, ucfg[iter].modName);
                                //Un-register module
                                //ret =
                                //   invt_unregister_module(ucfg[iter].modUuid);
                                continue;
                            } else {
                                /* register device*/
                                ret = invt_register_devices(ucfg[iter].modUuid,
                                                            registerDev);
                            }

                        } else {
                            goto cleanunitcfg;
                        }
                    }
                } else {
                    usys_log_debug("Read Unit Config fail for %s."
                                   "Error Code: %d",
                                   unitUuid, ret);
                    goto cleanunitcfg;
                }
            } else {
                ret = ERR_NODED_MEMORY_EXHAUSTED;
                usys_log_debug("Memory error while reading unit config for %s."
                               "Error Code: %d",
                               unitUuid, ret);
                goto cleanunitinfo;
            }
        } else {
            usys_log_debug("Read Unit Info fail for %s. Error Code: %d", pUuid,
                           ret);
            goto cleanunitinfo;
        }
    }

cleanunitcfg:
    invt_free_unit_cfg(ucfg, modCount);

cleanunitinfo:
    usys_free(uInfo);
    uInfo = NULL;
    return ret;
}

/* UKDB registering devices in Device DB so that these can be accessed later on.*/
int invt_register_devices(char *pUuid, RegisterDeviceCB registerDev) {
    int ret = 0;
    uint16_t size = 0;
    uint8_t count = 0;
    ModuleCfg *mcfg;
    char name[NAME_LENGTH] = { '\0' };
    ModuleInfo *minfo = (ModuleInfo *)usys_zmalloc(sizeof(ModuleInfo));
    if (minfo) {
        usys_log_debug("Inventory Read UKDB Module Info from module %s "
                       "for device registration process.",
                       pUuid);

        /* Read module info */
        ret = invt_read_module_info(pUuid, minfo, &size);
        if (!ret) {
            invt_print_module_info(minfo);
            count = minfo->devCount;
            usys_memcpy(name, minfo->name, usys_strlen(minfo->name));
            if (minfo) {
                usys_free(minfo);
                minfo = NULL;
            }

            /* Read Module Cfg */
            mcfg = (ModuleCfg *)invt_alloc_module_cfg(count);
            if (mcfg) {
                usys_log_debug("Inventory Read UKDB Module Info from module %s"
                               "for device registration process.",
                               pUuid);
                size = 0;

                /* Read module config */
                ret = invt_read_module_cfg(pUuid, mcfg, count, &size);
                if (!ret) {
                    invt_print_module_cfg(mcfg, count);

                    if (registerDev) {
                        /* Register devices in devicedb */
                        ret = registerDev(pUuid, name, count, mcfg);
                    }

                } else {
                    usys_log_debug("Read Module Config for %s failed."
                                   " Error Code: %d.",
                                   pUuid, ret);
                    goto cleanmcfg;
                }

            } else {
                ret = ERR_NODED_MEMORY_EXHAUSTED;
                usys_log_debug("Memory error while reading module Config for %s"
                               "Error: %s",
                               pUuid, usys_error(errno));
                goto cleanminfo;
            }

        } else {
            ret = ERR_NODED_MEMORY_EXHAUSTED;
            usys_log_debug("Err(%d)::UKDB Read UKDB Module Config for %s "
                           "failed.",
                           pUuid);
            goto cleanminfo;
        }
    }

cleanmcfg:
    invt_free_module_cfg(mcfg, count);

cleanminfo:
    usys_free(minfo);
    minfo = NULL;

    return ret;
}

int invt_pre_create_store_setup(char *puuid) {
    int ret = -1;
    uint8_t count = 0;
    UnitCfg *cfg = NULL;
    usys_log_debug("Caution !! Only meant for test setup to create DB.");

    cfg = override_master_db_info(puuid, &count);
    if (cfg) {
        invt_print_unit_cfg(cfg, count);
        /* Register modules so that unit info and other cfg can be accessed.*/
        for (int iter = 0; iter < count; iter++) {
            ret = invt_register_module(&cfg[iter]);
            if (ret) {
                usys_log_debug("Err(%d): UKDB:: Failed to register module %s.",
                               ret, cfg[iter].modUuid);
            }
        }

        for (int iter = 0; iter < count; iter++) {
            usys_free(cfg[iter].eepromCfg);
            cfg[iter].eepromCfg = NULL;
        }

    } else {
        usys_log_error(
            "Err(%d): Failed to get Unit Config for %s from test config.", ret,
            puuid);
    }
    usys_free(cfg);
    cfg = NULL;
    return ret;
}

/* Cleaning module map info from store */
void invt_exit() {
    store_deregister_all_module();
}

/* exit mfg module */
void invt_mfg_exit() {
    mfg_exit();
}

/* Inventory Init */
int invt_init(char *invtDb, RegisterDeviceCB regCb) {
    int ret = 0;
    uint8_t count = 1;

    /* Initialize Store*/
    ret = store_init();
    UnitCfg *cfg = usys_zmalloc(sizeof(UnitCfg));
    if (cfg) {
        ret = invt_get_master_unit_cfg(cfg, invtDb);
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
            ret = invt_register_modules(cfg->modUuid, regCb);
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

/* removed the database for the module */
int invt_remove_db(char *puuid) {
    int ret = 0;

    ret = store_remove(puuid);
    if (!ret) {
        ret = invt_deregister_module(puuid);
        if (ret) {
            usys_log_error(" Failed to unregister Module %s.", ret, puuid);
        }

    } else {
        usys_log_error(" Failed to remove database for Module %s.", ret, puuid);
    }

    return ret;
}

/* Creating inventory database from scratch using some manufacturing data in
 * JSON format. */
int invt_create_db(char *pUuid) {
    int ret = 0;
    uint16_t size = 0;

    /* Magic Word */
    SchemaMagicWord mw = { 0 };
    ret = store_read_block(pUuid, &mw, SCH_MAGIC_WORD_OFFSET,
                           sizeof(SchemaMagicWord)) != sizeof(SchemaMagicWord);
    if (ret) {
        usys_log_warn("Problem while reading Inventory data. Error Code:%d",
                      ret);
        usys_log_warn("Creating new one Inventory database.");
        //return  ERR_NODED_R_FAIL;
    }

    if (mw.magicWord == SCH_MAGIC_WORD) {
        usys_log_warn(
            "Inventory database is present in device. Re-writing it again.");
    } else {
        /* Write Magic Word */
        if (invt_write_magic_word(pUuid)) {
            usys_log_error("Inventory database creation failed.");
            return ERR_NODED_WR_FAIL;
        }
    }

    /* Update header info. */
    SchemaHeader *header;
    ret = mfg_fetch_header(&header, pUuid, &size);
    if (ret) {
        usys_log_error("Failed to read header.");
    }

    if (!header) {
        usys_log_error("Header read is NULL");
        usys_free(header);
        header = NULL;
        return (-1);
    } else {
        if (invt_write_header(pUuid, header)) {
            usys_log_error("Write to Inventory database header failed.");
            usys_free(header);
            header = NULL;
            return ERR_NODED_WR_FAIL;
        }
    }
    usys_free(header);
    header = NULL;

    /* Populate Index list from the MFG data.*/
    SchemaIdxTuple *index;
    ret = mfg_fetch_idx(&index, pUuid, &size);
    if (ret) {
        usys_log_error("Failed to read index.");
    }
    if (!index) {
        usys_log_error("Index read is NULL");
        return -1;
    }
    uint8_t idxCount = size / sizeof(SchemaIdxTuple);
    invt_print_index_table(index, idxCount);

    /* For each index table entry */
    uint8_t idx = 0;
    uint8_t modCount = 0;
    uint8_t idxIter = 0;

    /* Write Unit info first.*/
    /* Find where does field id sits in the index list read from MFG data.
     *  if not found skip to next.*/
    usys_log_trace("Starting Unit Info write for Module UUID %s.", pUuid);

    ret = invt_get_field_id_idx(index, FIELD_ID_UNIT_INFO, idxCount, &idxIter);
    if (!ret) {
        ret =
            invt_write_unit_info_data(pUuid, &index[idxIter], pUuid, &modCount);
        if (!ret) {
            idx++;
        } else {
            goto cleanindex;
        }
        usys_log_trace("Inventory Completed Unit Info write for "
                       "Module UUID %s.",
                       pUuid);
    } else {
        /* This id will use if the Module DB is getting created.
         * In this case we might not have any Unit info but surely we would have
         * a module info and module config to write.
         */
        modCount = 1;
        usys_log_trace("Inventory Unit Info field for Module UUID %s "
                       "not found.",
                       pUuid);
    }

    /*Write unit Config. */
    usys_log_trace("Inventory Starting Unit Config write for UUID %s "
                   "with %d modules.",
                   pUuid, modCount);

    ret = invt_get_field_id_idx(index, FIELD_ID_UNIT_CFG, idxCount, &idxIter);
    if (!ret) {
        //TODO: Don't even need UnitCFg here. it can be removed.
        ret = invt_write_unit_cfg_data(pUuid, &index[idxIter], modCount);
        if (ret) {
            goto cleanindex;
        }
        idx++;
        usys_log_trace("Inventory Completed Unit Config write for "
                       "Module UUID %s with %d modules.",
                       pUuid, modCount);
    } else {
        usys_log_trace("Inventory Unit Config field for Module UUID %s "
                       "not found.",
                       pUuid);
    }

    /* Write Module info */
    usys_log_trace("Inventory Starting Module Info write for "
                   "Module UUID %s.",
                   pUuid);

    ret = invt_get_field_id_idx(index, FIELD_ID_MODULE_INFO, idxCount, &idxIter);
    if (!ret) {
        uint8_t mcfg_iter = 0;
        ret = invt_get_field_id_idx(index, FIELD_ID_MODULE_CFG, idxCount,
                                 &mcfg_iter);
        if (!ret) {
            /* Module config and Module info both are added here to DB.*/
            ret = invt_write_module_info_data(
                pUuid, &index[idxIter], &index[mcfg_iter], &modCount, &idx);
            if (ret) {
                goto cleanindex;
            }
            idx++;

        } else {
            usys_log_debug("Inventory Module Config for Module UUID %s"
                           " not found.",
                           pUuid);
        }
    } else {
        usys_log_debug("Inventory Module Info for Module UUID %s "
                       "not found.",
                       pUuid);
    }

    /* All the remaining fields are just a files which are stored as it is.*/
    uint16_t genfield_id = FIELD_ID_FACT_CFG;
    uint8_t temp_idx = idx;
    while (temp_idx < idxCount) {
        ret = invt_get_field_id_idx(index, genfield_id, idxCount, &idxIter);
        if (!ret) {
            usys_log_debug("Inventory Writing data for field Id 0x%x"
                           " Module UIID %s.",
                           index[idxIter].fieldId, pUuid);

            ret = invt_write_generic_data(&index[idxIter], pUuid, genfield_id);
            if (!ret) {
                idx++;
                usys_log_debug("Inventory Wrote data for field Id 0x%x "
                               "Module UIID %s.",
                               index[idxIter].fieldId, pUuid);
            } else {
                goto cleanindex;
            }
        }
        temp_idx++;
        genfield_id++;
    }
    usys_log_debug("Inventory UKDB created with total of %d entries.", idx);

cleanindex:
    usys_free(index);
    index = NULL;
    return ret;
}

/* Validate the crc matches */
int invt_validate_payload(void *pData, uint32_t crc, uint16_t size) {
    int ret = 0;

    uint32_t calcCrc = crc_32(pData, size);
    if (calcCrc != crc) {
        ret = ERR_NODED_CRC_FAILURE;
    }

    usys_log_trace("Inventory CRC Expected 0x%x Calculated 0x%x.", crc,
                   calcCrc);
    return ret;
}

/* Read payload */
int invt_read_payload(char *pUuid, void *pData, uint16_t offset,
                      uint16_t size) {
    int ret = 0;

    if (store_read_block(pUuid, pData, offset, size) < size) {
        ret = ERR_NODED_R_FAIL;
        usys_log_error("Inventory Index read failed from offset 0x%x."
                       "Error Code: %d",
                       ret, offset);
    }

    usys_log_trace("Inventory read %d bytes of payload from 0x%x offset.", size,
                   offset);
    return ret;
}

/* Read payload for field id.*/
int invt_read_payload_for_field_id(char *pUuid, void **data, uint16_t fid,
                                   uint16_t *size) {
    int ret = -1;
    uint16_t idx = 0;
    void *payload = *data;
    SchemaIdxTuple *idxData;
    ret = invt_search_field_id(pUuid, &idxData, &idx, fid);
    if (ret) {
        usys_log_error("Err(%d): Inventory search error for field id 0x%x.",
                       ret, fid);
        return ret;
    }

    *size = idxData->payloadSize;

    /* Check if memory available */
    if (!payload) {
        payload = usys_zmalloc(*size);
        if (payload) {
            *data = payload;
        } else {
            usys_log_error("Inventory memory failure. Error %s.", ret,
                           usys_error(errno));
            return errno;
        }
    }

    if (payload) {
        ret = invt_read_payload(pUuid, payload, idxData->payloadOffset,
                                idxData->payloadSize);
        if (ret) {
            usys_log_error("Payload read failure for the "
                           "field id 0x%x.Error Code: %d",
                           fid, ret);
            usys_free(idxData);
            idxData = NULL;
        }
    }

    /* validate index */
    ret = invt_validate_payload(payload, idxData->payloadCrc, *size);
    if (ret) {
        usys_log_error("CRC failure for the field id 0x%x. Error Code: %d", fid,
                       ret);
        *size = 0;
    }

    usys_free(idxData);
    idxData = NULL;
    return ret;
}

void invt_print_header(SchemaHeader *header) {
    usys_log_trace(
        "*************************************************************");
    usys_log_trace(
        "*********************Schema Header*******************************");
    usys_log_trace(
        "*************************************************************");
    usys_log_trace("  *	 UKDB Version:                  v%d.%d",
                   header->version.major, header->version.minor);
    usys_log_trace("  *	 UKDB Table Offset:             0x%x",
                   header->idxTblOffset);
    usys_log_trace("  *	 UKDB Index Tuple Size:         %d",
                   header->idxTplSize);
    usys_log_trace("  *	 UKDB Tuple max count:          %d",
                   header->idxTplMaxCount);
    usys_log_trace("  *	 UKDB Current Tuple Id:         %d", header->idxCurTpl);
    usys_log_trace("  *	 UKDB Capability:               0x%x", header->modCap);
    usys_log_trace("  *	 UKDB Mode:                     0x%x", header->modMode);
    usys_log_trace("  *	 UKDB Device Owner:             0x%x",
                   header->modDevOwn);
    usys_log_trace(
        "*************************************************************");
    usys_log_trace(
        "*************************************************************");
}

void invt_print_index_table(SchemaIdxTuple *idx_tbl, uint8_t count) {
    usys_log_trace(
        "*************************************************************");
    usys_log_trace(
        "*********************Index Table*******************************");
    usys_log_trace(
        "*************************************************************");
    for (uint8_t iter = 0; iter < count; iter++) {
        usys_log_trace(
            "*************************************************************");
        usys_log_trace("  *	 Field Id:                      0x%04x ",
                       idx_tbl[iter].fieldId);
        usys_log_trace("  *	 Payload Offset:                0x%04x ",
                       idx_tbl[iter].payloadOffset);
        usys_log_trace("  *	 Payload Size:                  0x%04x bytes",
                       idx_tbl[iter].payloadSize);
        usys_log_trace("  *	 Payload Version:               v%d.%d",
                       idx_tbl[iter].payloadVer.major,
                       idx_tbl[iter].payloadVer.minor);
        usys_log_trace("  *	 Payload CRC:                   0x%04x",
                       idx_tbl[iter].payloadCrc);
        usys_log_trace("  *	 Field State:                   0x%01x",
                       idx_tbl[iter].state);
        usys_log_trace("  *	 Field Valid:                   0x%01x",
                       idx_tbl[iter].valid);
        usys_log_trace(
            "*************************************************************");
    }
    usys_log_trace(
        "*************************************************************");
}

void invt_print_unit_info(UnitInfo *pUnitInfo) {
    usys_log_trace(
        "*************************************************************");
    usys_log_trace(
        "*********************Unit Info*******************************");
    usys_log_trace(
        "*************************************************************");
    usys_log_trace("  *	 Unit UUID:                    %s", pUnitInfo->uuid);
    usys_log_trace("  *	 Unit Name:                    %s", pUnitInfo->name);
    usys_log_trace("  *	 Unit Type:                    0x%01x",
                   pUnitInfo->unit);
    usys_log_trace("  *	 Unit Part No.:                %s", pUnitInfo->partNo);
    usys_log_trace("  *	 Unit Skew:                    %s", pUnitInfo->skew);
    usys_log_trace("  *	 Unit MAC:                     %s", pUnitInfo->mac);
    usys_log_trace("  *	 Unit SWVer:                   v%d.%d",
                   pUnitInfo->swVer.major, pUnitInfo->swVer.minor);
    usys_log_trace("  *	 Unit PSWVer:                  v%d.%d",
                   pUnitInfo->swVer.major, pUnitInfo->swVer.minor);
    usys_log_trace("  *	 Unit Assm_Date:               %s",
                   pUnitInfo->assmDate);
    usys_log_trace("  *	 Unit OEM_Name:                %s", pUnitInfo->oemName);
    usys_log_trace("  *	 Unit Module Count:            %d",
                   pUnitInfo->modCount);
    usys_log_trace(
        "*************************************************************");
    usys_log_trace(
        "*************************************************************");
}

void invt_print_module_info(ModuleInfo *p_minfo) {
    usys_log_trace(
        "*************************************************************");
    usys_log_trace(
        "*******************Module Info*******************************");
    usys_log_trace(
        "*************************************************************");
    usys_log_trace("  *	  Module UUID:              %s", p_minfo->uuid);
    usys_log_trace("  *	  Module Name:              %s", p_minfo->name);
    usys_log_trace("  *	  Module Type:              0x%01x", p_minfo->module);
    usys_log_trace("  *	  Module Part No.:          %s", p_minfo->partNo);
    usys_log_trace("  *	  Module HWVer:             %s", p_minfo->hwVer);
    usys_log_trace("  *	  Module MAC:               %s", p_minfo->mac);
    usys_log_trace("  *	  Module SWVer:             v%d.%d",
                   p_minfo->swVer.major, p_minfo->swVer.minor);
    usys_log_trace("  *	  Module PSWVer:            v%d.%d",
                   p_minfo->swVer.major, p_minfo->swVer.minor);
    usys_log_trace("  *	  Module MFG_Date:          %s", p_minfo->mfgDate);
    usys_log_trace("  *	  Module MFG_Name:          %s", p_minfo->mfgName);
    usys_log_trace("  *	  Module Device Count:      %d", p_minfo->devCount);
    usys_log_trace(
        "*************************************************************");
    usys_log_trace(
        "*************************************************************");
}

void invt_print_dev_gpio_cfg(DevGpioCfg *pdev) {
    if (pdev) {
        usys_log_trace("  *   GPIO Number:               0x%x", pdev->gpioNum);
        usys_log_trace("  *   GPIO Direction:            0x%x",
                       pdev->direction);
    }
}

void invt_print_dev_i2c_cfg(DevI2cCfg *pdev) {
    if (pdev) {
        usys_log_trace("  *   I2C Bus:                   0x%x", pdev->bus);
        usys_log_trace("  *   Address:                   0x%x", pdev->add);
    }
}

void invt_print_dev_spi_cfg(DevSpiCfg *pdev) {
    if (pdev) {
        usys_log_trace("  *   SPI Bus:                   0x%x", pdev->bus);
        invt_print_dev_gpio_cfg(&(pdev->cs));
    }
}

void invt_print_dev_uart_cfg(DevUartCfg *pdev) {
    if (pdev) {
        usys_log_trace("  *   Uart Number:               0x%x", pdev->uartNo);
    }
}

void invt_print_dev(void *dev, DeviceClass class) {
    switch (class) {
    case DEV_CLASS_GPIO: {
        invt_print_dev_gpio_cfg(dev);
        break;
    }
    case DEV_CLASS_I2C: {
        invt_print_dev_i2c_cfg(dev);
        break;
    }
    case DEV_CLASS_SPI: {
        invt_print_dev_spi_cfg(dev);
        break;
    }
    case DEV_CLASS_UART: {
        invt_print_dev_uart_cfg(dev);
        break;
    }
    default: {
        usys_log_trace("  *   Invalid device found.");
    }
    }
}

void invt_print_unit_cfg(UnitCfg *pUnitCfg, uint8_t count) {
    uint8_t iter = 0;
    usys_log_trace(
        "*************************************************************");
    usys_log_trace(
        "*******************Unit Config*******************************");
    for (; iter < count; iter++) {
        usys_log_trace(
            "*************************************************************");
        usys_log_trace("  *	  Module UUID:           %s",
                       pUnitCfg[iter].modUuid);
        usys_log_trace("  *	  Module Name:           %s",
                       pUnitCfg[iter].modName);
        usys_log_trace("  *    EEPROM SysFS Name:     %s",
                       pUnitCfg[iter].sysFs);
        invt_print_dev_i2c_cfg(pUnitCfg[iter].eepromCfg);
        usys_log_trace(
            "*************************************************************");
    }
}

void invt_print_module_cfg(ModuleCfg *p_mcfg, uint8_t count) {
    uint8_t iter = 0;
    usys_log_trace(
        "*************************************************************");
    usys_log_trace(
        "*******************Module Config*****************************");
    for (; iter < count; iter++) {
        usys_log_trace(
            "*************************************************************");
        usys_log_trace("  *	  Device Name:              %s",
                       p_mcfg[iter].devName);
        usys_log_trace("  *	  Device Disc:              %s",
                       p_mcfg[iter].devDesc);
        usys_log_trace("  *	  Device Type:              0x%x",
                       p_mcfg[iter].devType);
        usys_log_trace("  *	  Device Class:             0x%x",
                       p_mcfg[iter].devClass);
        usys_log_trace("  *	  Device SysFile:           %s",
                       p_mcfg[iter].sysFile);
        invt_print_dev(p_mcfg[iter].cfg, p_mcfg[iter].devClass);
        usys_log_trace(
            "*************************************************************");
    }
}

/* Read the payload from store for applications*/
int invt_read_payload_from_store(char *pUuid, void *pData, uint16_t id,
                                 uint16_t *size) {
    int ret = -1;
    switch (id) {
    case FIELD_ID_UNIT_INFO: {
        ret = invt_read_payload_for_field_id(pUuid, pData, FIELD_ID_UNIT_INFO,
                                             size);
        break;
    }
    case FIELD_ID_UNIT_CFG: ///this won't work
    {
        ret = invt_read_payload_for_field_id(pUuid, pData, FIELD_ID_UNIT_CFG,
                                             size);
        break;
    }
    case FIELD_ID_MODULE_INFO: {
        ret = invt_read_payload_for_field_id(pUuid, pData, FIELD_ID_MODULE_INFO,
                                             size);
        break;
    }
    case FIELD_ID_MODULE_CFG: {
        ret = invt_read_payload_for_field_id(pUuid, pData, FIELD_ID_MODULE_CFG,
                                             size);
        break;
    }
    case FIELD_ID_FACT_CFG: {
        ret = invt_read_payload_for_field_id(pUuid, pData, FIELD_ID_FACT_CFG,
                                             size);
        break;
    }
    case FIELD_ID_USER_CFG: {
        ret = invt_read_payload_for_field_id(pUuid, pData, FIELD_ID_USER_CFG,
                                             size);
        break;
    }

    case FIELD_ID_FACT_CALIB: {
        ret = invt_read_payload_for_field_id(pUuid, pData, FIELD_ID_FACT_CALIB,
                                             size);
        break;
    }

    case FIELD_ID_USER_CALIB: {
        ret = invt_read_payload_for_field_id(pUuid, pData, FIELD_ID_USER_CALIB,
                                             size);
        break;
    }
    case FIELD_ID_BS_CERTS: {
        ret = invt_read_payload_for_field_id(pUuid, pData, FIELD_ID_BS_CERTS,
                                             size);
        break;
    }
    case FIELD_ID_CLOUD_CERTS: {
        ret = invt_read_payload_for_field_id(pUuid, pData, FIELD_ID_CLOUD_CERTS,
                                             size);
        break;
    }
    default: {
        ret = ERR_NODED_DB_MISSING_FIELD;
        usys_log_error("Invalid Field id supplied by Index entry."
                       "Error Code %d",
                       ret);
    }
    }

    if (ret) {
        pData = NULL;
        ret = ERR_NODED_R_FAIL;
        usys_log_error("Failed to read info on 0x%x.Error Code %d", id, ret);
    }
    return ret;
}

/* This will read unit info and size of the info.*/
int invt_read_unit_info(char *pUuid, UnitInfo *p_info, uint16_t *size) {
    int ret = -1;
    uint16_t unit_fid = FIELD_ID_UNIT_INFO;
    uint16_t idx = 0;

    SchemaIdxTuple *idxData;
    ret = invt_search_field_id(pUuid, &idxData, &idx, unit_fid);
    if (ret) {
        usys_log_error("Err(%d): UKDB search error for field id 0x%x.", ret,
                       unit_fid);
        ret = ERR_NODED_DB_MISSING_UNIT_INFO;
        return ret;
    }

    if (p_info) {
        ret = invt_read_payload(pUuid, p_info, idxData->payloadOffset,
                                sizeof(UnitInfo));
        if (ret) {
            usys_log_error(
                "Err(%d): Payload read failure for the field id 0x%x.", ret,
                unit_fid);
            usys_free(idxData);
            idxData = NULL;
        }
        //p_info = info;
        *size = sizeof(UnitInfo);
    }

    /* validate index */
    ret = invt_validate_payload(p_info, idxData->payloadCrc, *size);
    if (ret) {
        usys_log_error("Err(%d): CRC failure for the field id 0x%x.", ret,
                       unit_fid);
    }
    usys_free(idxData);
    idxData = NULL;
    return ret;
}

int invt_deserialize_unit_cfg_data(UnitCfg **pUnitCfg, char *payload, uint8_t count,
                              uint16_t *size) {
    /* || Unit Info 1 | EEPROM CFG  || Unit Info 2 | EEPROM CFG  || */
    int ret = 0;
    int offset = 0;
    for (int iter = 0; iter < count; iter++) {
        /* Copy Unit Cfg first*/
        usys_memcpy(&(*pUnitCfg)[iter], payload + offset, sizeof(UnitCfg));
        /* Create a eeprom cfg and assign reference to eeprom_cfg in UnitCfg */
        offset = offset + sizeof(UnitCfg);
        /* Our Unit config assumes all eeprom are on I2C bus*/
        DevI2cCfg *icfg = usys_zmalloc(sizeof(DevI2cCfg));
        if (icfg) {
            usys_memcpy(icfg, payload + offset, sizeof(DevI2cCfg));
            (*pUnitCfg)[iter].eepromCfg = icfg;
        } else {
            ret = ERR_NODED_MEMORY_EXHAUSTED;
        }
        offset = offset + sizeof(DevI2cCfg);
        /* Size returned to reader of Unit config.*/
        *size = *size + sizeof(UnitCfg);
    }
    return ret;
}

int invt_read_unit_cfg(char *pUuid, UnitCfg *pUnitCfg, uint8_t count,
                       uint16_t *size) {
    int ret = -1;
    uint16_t fid = FIELD_ID_UNIT_CFG;
    uint16_t idx = 0;
    char *payload = NULL;
    SchemaIdxTuple *idxData;

    /* Searching for index */
    ret = invt_search_field_id(pUuid, &idxData, &idx, fid);
    if (ret) {
        usys_log_error("Inventory search error for field id 0x%x."
                       "Error Code: %d",
                       fid, ret);
        ret = ERR_NODED_DB_MISSING_UNIT_CFG;
        return ret;
    }

    /*Reading payload*/
    payload = usys_zmalloc(sizeof(char) * idxData->payloadSize);
    if (payload) {
        ret = invt_read_payload(pUuid, payload, idxData->payloadOffset,
                                idxData->payloadSize);
        if (ret) {
            usys_log_error("Payload read failure for the field id 0x%x.",
                           "Error Code: %d", fid, ret);
            usys_free(idxData);
            idxData = NULL;
        }

        /* validate CRC check*/
        ret = invt_validate_payload(payload, idxData->payloadCrc,
                                    idxData->payloadSize);
        if (ret) {
            usys_log_error("CRC failure for the field id 0x%x.", ret, pUnitCfg);
            goto cleanup;
        }

        /* Deserialize to Unit Config.*/
        ret = invt_deserialize_unit_cfg_data(&pUnitCfg, payload, count, size);
        if (ret) {
            ret = ERR_NODED_DESERIAL_FAIL;
            usys_log_error("Deserialize failure for Unit Config."
                           "Error Code: %d",
                           ret);
            goto cleanup;
        }
    }

cleanup:
    usys_free(payload);
    payload = NULL;
    usys_free(idxData);
    idxData = NULL;
    return ret;
}

/* This will read module info and size of the info.*/
int invt_read_module_info(char *pUuid, ModuleInfo *p_info, uint16_t *size) {
    //TODO
    int ret = -1;
    uint16_t fid = FIELD_ID_MODULE_INFO;
    uint16_t idx = 0;

    SchemaIdxTuple *idxData;
    ret = invt_search_field_id(pUuid, &idxData, &idx, fid);
    if (ret) {
        usys_log_error("Inventory search error for field id 0x%x."
                       "Error Code: %d",
                       fid, ret);
        ret = ERR_NODED_DB_MISSING_MODULE_INFO;
        return ret;
    }

    if (p_info) {
        ret = invt_read_payload(pUuid, p_info, idxData->payloadOffset,
                                sizeof(ModuleInfo));
        if (ret) {
            usys_log_error("Payload read failure for the field id 0x%x.",
                           "Error Code: %d", fid, ret);
            usys_free(idxData);
            idxData = NULL;
        }
        *size = sizeof(ModuleInfo);
    }

    /* validate index */
    ret = invt_validate_payload(p_info, idxData->payloadCrc, *size);
    if (ret) {
        usys_log_error("CRC failure for the field id 0x%x.", "Error Code: %d",
                       fid, ret);
    }
    usys_free(idxData);
    idxData = NULL;
    return ret;
}

void *invt_deserialize_devices(const char *payload, int offset, uint16_t class,
                               int *size) {
    void *dev = NULL;
    const char *cfgData = payload + offset;
    switch (class) {
    case DEV_CLASS_GPIO: {
        DevGpioCfg *cfg = usys_zmalloc(sizeof(DevGpioCfg));
        if (cfg) {
            usys_memcpy(cfg, cfgData, sizeof(DevGpioCfg));
        } else {
            cfg = NULL;
        }
        dev = cfg;
        *size = sizeof(DevGpioCfg);
        break;
    }
    case DEV_CLASS_I2C: {
        DevI2cCfg *cfg = usys_zmalloc(sizeof(DevI2cCfg));
        if (cfg) {
            usys_memcpy(cfg, cfgData, sizeof(DevI2cCfg));
        } else {
            cfg = NULL;
        }
        dev = cfg;
        *size = sizeof(DevI2cCfg);
        break;
    }
    case DEV_CLASS_SPI: {
        DevSpiCfg *cfg = usys_zmalloc(sizeof(DevSpiCfg));
        if (cfg) {
            usys_memcpy(cfg, cfgData, sizeof(DevSpiCfg));
        } else {
            cfg = NULL;
        }
        dev = cfg;
        *size = sizeof(DevSpiCfg);
        break;
    }
    case DEV_CLASS_UART: {
        DevUartCfg *cfg = usys_zmalloc(sizeof(DevUartCfg));
        if (cfg) {
            usys_memcpy(cfg, cfgData, sizeof(DevUartCfg));
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
        usys_log_error("Unkown device type failed to parse.Error Code: %d",
                       ERR_NODED_INVALID_DEVICE_CFG);
    }
    return dev;
}

int invt_deserialize_module_cfg_data(ModuleCfg **p_mcfg, char *payload,
                                uint8_t count, uint16_t *size) {
    /* Layout
     *  || Unit Info 1 | EEPROM CFG  || Unit Info 2 | EEPROM CFG  ||
     *  */
    int ret = 0;
    int offset = 0;
    for (int iter = 0; iter < count; iter++) {
        /* Copy Module Cfg first*/
        usys_memcpy(&(*p_mcfg)[iter], payload + offset, sizeof(ModuleCfg));
        offset = offset + sizeof(ModuleCfg);
        int cfg_size = 0;

        /* Create a device cfg and assign reference to cfg in ModuleCfg */
        void *cfg = invt_deserialize_devices(
            payload, offset, (*p_mcfg)[iter].devClass, &cfg_size);
        if (cfg) {
            (*p_mcfg)[iter].cfg = cfg;
        } else {
            ret = ERR_NODED_DESERIAL_FAIL;
            usys_log_error("Deserialization failure for module config."
                           "Error Code: %d",
                           ret);
        }

        offset = offset + cfg_size;
        /* Size returned to reader of Unit config.*/
        *size = *size + sizeof(ModuleCfg);
    }

    return ret;
}

/* This will read module config and count of the module*/
int invt_read_module_cfg(char *pUuid, ModuleCfg *pCfg, uint8_t count,
                         uint16_t *size) {
    int ret = -1;
    uint16_t fid = FIELD_ID_MODULE_CFG;
    uint16_t idx = 0;
    char *payload = NULL;
    SchemaIdxTuple *idxData;

    ret = invt_search_field_id(pUuid, &idxData, &idx, fid);
    if (ret) {
        usys_log_error("Inventory search error for field id 0x%x."
                       "Error Code: %d",
                       fid, ret);
        ret = ERR_NODED_DB_MISSING_MODULE_CFG;
        return ret;
    }

    payload = usys_zmalloc(sizeof(char) * idxData->payloadSize);
    if (payload) {
        /* Read the DB*/
        ret = invt_read_payload(pUuid, payload, idxData->payloadOffset,
                                idxData->payloadSize);
        if (ret) {
            usys_log_error("Payload read failure for the field id 0x%x.",
                           "Error Code: %d", fid, ret);
            goto cleanup;
        }

        /* validate CRC check */
        ret = invt_validate_payload(payload, idxData->payloadCrc,
                                    idxData->payloadSize);
        if (ret) {
            usys_log_error("CRC failure for the field id 0x%x.", ret, fid);
            goto cleanup;
        }

        if (pCfg) {
            /* Deserialize payload to Module Config */
            ret = invt_deserialize_module_cfg_data(&pCfg, payload, count, size);
            if (!ret) {
                usys_log_debug("Read Module Info %d bytes for Module %s "
                               "with device count %d.",
                               *size, pUuid, count);
            } else {
                usys_log_error("Payload deserialize failure for the "
                               "field id 0x%x.Error Code: %d",
                               fid, ret);
                goto cleanup;
            }
        } else {
            ret = ERR_NODED_INVALID_POINTER;
            usys_log_error("Invalid payload pointer for the field id 0x%x."
                           "Error Code: %d",
                           fid, ret);
            goto cleanup;
        }
    } else {
        ret = ERR_NODED_MEMORY_EXHAUSTED;
        usys_log_error("Memory error while reading payload for the "
                       "field id 0x%x. Error: %s",
                       fid, usys_error(errno));
        goto cleanup;
    }

cleanup:
    usys_free(payload);
    payload = NULL;
    usys_free(idxData);
    idxData = NULL;
    return ret;
}

/* Read fact config. */
int invt_read_fact_config(char *pUuid, void *data, uint16_t *size) {
    int ret = 0;

    ret = invt_read_payload_from_store(pUuid, data, FIELD_ID_FACT_CFG, size);
    if (ret) {
        usys_log_error("Inventory failed to read info on 0x%x.Error Code: %d",
                       FIELD_ID_FACT_CFG, ret);
    }
    usys_log_debug("Inventory Fact Config Field Id :%d "
                   "data read with size %d bytes.",
                   FIELD_ID_FACT_CFG, *size);

    return ret;
}

/* Read user config. */
int invt_read_user_config(char *pUuid, void *data, uint16_t *size) {
    int ret = 0;
    ret = invt_read_payload_from_store(pUuid, data, FIELD_ID_USER_CFG, size);
    if (ret) {
        usys_log_error("Inventory failed to read info on 0x%x.Error Code: %d",
                       FIELD_ID_USER_CFG, ret);
    }
    usys_log_debug("Inventory Fact Config Field Id : %d "
                   "data read with size %d bytes.",
                   FIELD_ID_USER_CFG, *size);
    return ret;
}

/* Read fact calib. */
int invt_read_fact_calib(char *pUuid, void *data, uint16_t *size) {
    int ret = 0;
    ret = invt_read_payload_from_store(pUuid, data, FIELD_ID_FACT_CALIB, size);
    if (ret) {
        usys_log_error("Inventory failed to read info on 0x%x.Error Code: %d",
                       FIELD_ID_FACT_CALIB);
    }
    usys_log_debug("Inventory Fact Calibration Field Id : %d "
                   "data read with size %d bytes.",
                   FIELD_ID_FACT_CALIB, *size);
    return ret;
}

/* Read user calib. */
int invt_read_user_calib(char *pUuid, void *data, uint16_t *size) {
    int ret = 0;
    ret = invt_read_payload_from_store(pUuid, data, FIELD_ID_USER_CALIB, size);
    if (ret) {
        usys_log_error("Inventory failed to read info on 0x%x.Error Code: %d",
                       FIELD_ID_USER_CALIB, ret);
    }
    usys_log_debug("Inventory User Calibration Field Id : %d "
                   "data read with size %d bytes.",
                   FIELD_ID_USER_CALIB, *size);
    return ret;
}

/* Read bootstrap certs. */
int invt_read_bs_certs(char *pUuid, void *data, uint16_t *size) {
    int ret = 0;
    ret = invt_read_payload_from_store(pUuid, data, FIELD_ID_BS_CERTS, size);
    if (ret) {
        usys_log_error("Inventory failed to read info on 0x%x.Error Code: %d",
                       FIELD_ID_BS_CERTS, ret);
    }
    usys_log_debug("Inventory Bootstrap certs Field Id : %d data "
                   "read with size %d bytes.",
                   FIELD_ID_BS_CERTS, *size);
    return ret;
}

/* Read cloud certs. */
int invt_read_cloud_certs(char *pUuid, void *data, uint16_t *size) {
    int ret = 0;
    ret = invt_read_payload_from_store(pUuid, data, FIELD_ID_CLOUD_CERTS, size);
    if (ret) {
        usys_log_error("Inventory failed to read info on 0x%x.Error Code: %d",
                       FIELD_ID_CLOUD_CERTS, ret);
    }
    usys_log_debug("Inventory Lwm2m certs Field Id : %d data "
                   "read with size %d bytes.",
                   FIELD_ID_CLOUD_CERTS, *size);
    return ret;
}

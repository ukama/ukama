/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "mfg.h"

#include "errorcode.h"
#include "jdata.h"

#include "usys_log.h"

const MfgOperations *mfgOps;

const MfgOperations *mfgOps =
    &(MfgOperations){ .init = jdata_init,
                      .exit = jdata_exit,
                      .readHeader = jdata_fetch_header,
                      .readIndex = jdata_fetch_idx,
                      .readUnitInfo = jdata_fetch_unit_info,
                      .readUnitCfg = jdata_fetch_unit_cfg,
                      .readModuleInfoByUuid = jdata_fetch_module_info_by_uuid,
                      .readModuleCfg = jdata_fetch_module_cfg,
                      .readFactCfg = jdata_fetch_fact_cfg,
                      .readUserCfg = jdata_fetch_user_cfg,
                      .readFactCalib = jdata_fetch_fact_calib,
                      .readUserCalib = jdata_fetch_user_calib,
                      .readBsCerts = jdata_fetch_bs_certs,
                      .readCloudCerts = jdata_fetch_cloud_certs };

int mfg_init(void *data) {
    int ret = 0;
    if (data) {
        ret = mfgOps->init(data);
        if (ret) {
            usys_log_error("Mfg initialization failed. Error: %d", ret);
            return ret;
        }
    }
    usys_log_debug("Mfg initialization completed.");
    return ret;
}

void mfg_exit() {
    if (mfgOps->exit) {
        mfgOps->exit();
    }
}

int mfg_fetch_header(SchemaHeader **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = mfgOps->readHeader(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int mfg_fetch_idx(SchemaIdxTuple **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = mfgOps->readIndex(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int mfg_fetch_unit_info(UnitInfo **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = mfgOps->readUnitInfo(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int mfg_fetch_unit_cfg(UnitCfg **data, char *uuid, uint16_t *size, uint8_t count) {
    int ret = 0;
    *data = mfgOps->readUnitCfg(uuid, size, count);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int mfg_fetch_module_info(ModuleInfo **data, char *uuid, uint16_t *size,
                          uint8_t idx) {
    int ret = 0;
    *data = mfgOps->readModuleInfo(uuid, size, idx);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int mfg_fetch_module_info_by_uuid(ModuleInfo **data, char *uuid, uint16_t *size,
                                  uint8_t count) {
    int ret = 0;
    *data = mfgOps->readModuleInfoByUuid(uuid, size, count);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int mfg_fetch_module_cfg(void **data, char *uuid, uint16_t *size,
                         uint8_t count) {
    int ret = 0;
    *data = mfgOps->readModuleCfg(uuid, size, count);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int mfg_fetch_fact_cfg(void **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = mfgOps->readFactCfg(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int mfg_fetch_user_cfg(void **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = mfgOps->readUserCfg(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int mfg_fetch_fact_calib(void **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = mfgOps->readFactCalib(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int mfg_fetch_user_calib(void **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = mfgOps->readUserCalib(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int mfg_fetch_bs_certs(void **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = mfgOps->readBsCerts(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int mfg_fetch_cloud_certs(void **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = mfgOps->readCloudCerts(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

/* Read the payload from json parser or any other input provided to mfg.*/
int mfg_fetch_payload_from_mfg_data(void **data, char *uuid, uint16_t *size,
                                    uint16_t id) {
    int ret = -1;
    switch (id) {
    case FIELD_ID_FACT_CFG: {
        ret = mfg_fetch_fact_cfg(data, uuid, size);
        break;
    }
    case FIELD_ID_USER_CFG: {
        ret = mfg_fetch_user_cfg(data, uuid, size);
        break;
    }

    case FIELD_ID_FACT_CALIB: {
        ret = mfg_fetch_fact_calib(data, uuid, size);
        break;
    }

    case FIELD_ID_USER_CALIB: {
        ret = mfg_fetch_user_calib(data, uuid, size);
        break;
    }
    case FIELD_ID_BS_CERTS: {
        ret = mfg_fetch_bs_certs(data, uuid, size);
        break;
    }
    case FIELD_ID_CLOUD_CERTS: {
        ret = mfg_fetch_cloud_certs(data, uuid, size);
        break;
    }
    default: {
        ret = ERR_NODED_JSON_PARSER;
        usys_log_error("Invalid Field id supplied by Index entry.Error %d",
                       ret);
    }
    }

    if (!data) {
        data = NULL;
        ret = ERR_NODED_DB_MISSING_INFO;
        log_error("JSON parser failed to read info on 0x%x.Error %d", id, ret);
    }
    return ret;
}

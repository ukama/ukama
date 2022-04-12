/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "ukdb/idb/idb.h"

#include "headers/errorcode.h"
#include "headers/ubsp/ukdblayout.h"
#include "headers/utils/log.h"
#include "ukdb/idb/cs.h"
#include "ukdb/idb/jp.h"

const IDBFxnTable *idb_fxn_table;

const IDBFxnTable jp_fxn_table = { .init = jp_init,
                                   .exit = jp_exit,
                                   .read_header = jp_fetch_header,
                                   .read_index = jp_fetch_index,
                                   .read_unit_info = jp_fetch_unit_info,
                                   .read_unit_cfg = jp_fetch_unit_cfg,
                                   .read_module_info_by_uuid =
                                       jp_fetch_module_info_by_uuid,
                                   .read_module_cfg = jp_fetch_module_cfg,
                                   .read_fact_config = jp_fetch_fact_config,
                                   .read_user_config = jp_fetch_user_config,
                                   .read_fact_calib = jp_fetch_fact_calib,
                                   .read_user_calib = jp_fetch_user_calib,
                                   .read_bs_certs = jp_fetch_bs_certs,
                                   .read_lwm2m_certs = jp_fetch_lwm2m_certs };

const IDBFxnTable cs_fxn_table = {
    .init = cs_init,
    .exit = NULL,
    .read_header = cs_fetch_header,
    .read_index = cs_fetch_index,
    .read_unit_info = cs_fetch_unit_info,
    .read_unit_cfg = cs_fetch_unit_cfg,
    .read_module_info = cs_fetch_module_info,
    .read_module_info_by_uuid = cs_fetch_module_info_by_uuid,
    .read_module_cfg = cs_fetch_module_cfg,
    .read_fact_config = cs_fetch_fact_config,
    .read_user_config = cs_fetch_user_config,
    .read_fact_calib = cs_fetch_fact_calib,
    .read_user_calib = cs_fetch_user_calib,
    .read_bs_certs = cs_fetch_bs_certs,
    .read_lwm2m_certs = cs_fetch_lwm2m_certs,
};

int idb_init(void *data) {
    int ret = 0;
    if (data) {
        /* Start with the JSON parser */
        idb_fxn_table = &jp_fxn_table;
        ret = idb_fxn_table->init(data);
        if (ret) {
            // TODO:: Only for test purpose. remove for Targets.
            /* If JSON parser init fails switch to C struct */
            idb_fxn_table = &cs_fxn_table;
            ret = idb_fxn_table->init(data);
        }
    } else {
        /* if  no data is available switch to c structs.*/
        idb_fxn_table = &cs_fxn_table;
        ret = idb_fxn_table->init(data);
    }
    log_debug("IDB:: IDB layer initialized.");
    return ret;
}

void idb_exit() {
    if (idb_fxn_table->exit) {
        idb_fxn_table->exit();
    }
}

int idb_fetch_header(void **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = idb_fxn_table->read_header(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int idb_fetch_index(void **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = idb_fxn_table->read_index(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int idb_fetch_unit_info(void **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = idb_fxn_table->read_unit_info(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int idb_fetch_unit_cfg(void **data, char *uuid, uint16_t *size, uint8_t count) {
    int ret = 0;
    *data = idb_fxn_table->read_unit_cfg(uuid, size, count);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int idb_fetch_module_info(void **data, char *uuid, uint16_t *size,
                          uint8_t idx) {
    int ret = 0;
    *data = idb_fxn_table->read_module_info(uuid, size, idx);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int idb_fetch_module_info_by_uuid(void **data, char *uuid, uint16_t *size,
                                  uint8_t count) {
    int ret = 0;
    *data = idb_fxn_table->read_module_info_by_uuid(uuid, size, count);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int idb_fetch_module_cfg(void **data, char *uuid, uint16_t *size,
                         uint8_t count) {
    int ret = 0;
    *data = idb_fxn_table->read_module_cfg(uuid, size, count);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int idb_fetch_fact_config(void **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = idb_fxn_table->read_fact_config(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int idb_fetch_user_config(void **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = idb_fxn_table->read_user_config(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int idb_fetch_fact_calib(void **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = idb_fxn_table->read_fact_calib(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int idb_fetch_user_calib(void **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = idb_fxn_table->read_user_calib(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int idb_fetch_bs_certs(void **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = idb_fxn_table->read_bs_certs(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

int idb_fetch_lwm2m_certs(void **data, char *uuid, uint16_t *size) {
    int ret = 0;
    *data = idb_fxn_table->read_lwm2m_certs(uuid, size);
    if (!data) {
        ret = -1;
    }
    return ret;
}

/* Read the payload from json parser or any other input.*/
/* This is meant for input to UKDB */
int idb_fetch_payload_from_mfgdata(void **data, char *uuid, uint16_t *size,
                                   uint16_t id) {
    int ret = -1;
    switch (id) {
    case FIELDID_FACT_CONFIG: {
        ret = idb_fetch_fact_config(data, uuid, size);
        break;
    }
    case FIELDID_USER_CONFIG: {
        ret = idb_fetch_user_config(data, uuid, size);
        break;
    }

    case FIELDID_FACT_CALIB: {
        ret = idb_fetch_fact_calib(data, uuid, size);
        break;
    }

    case FIELDID_USER_CALIB: {
        ret = idb_fetch_user_calib(data, uuid, size);
        break;
    }
    case FIELDID_BS_CERTS: {
        ret = idb_fetch_bs_certs(data, uuid, size);
        break;
    }
    case FIELDID_LWM2M_CERTS: {
        ret = idb_fetch_lwm2m_certs(data, uuid, size);
        break;
    }
    default: {
        ret = ERR_UBSP_JSON_PARSER;
        log_error("Err(%d): Invalid Field id supplied by Index entry.", ret);
    }
    }

    if (!data) {
        data = NULL;
        ret = ERR_UBSP_DB_MISSING_INFO;
        log_error("ERR(%d): JSON parser failed to read info on 0x%x.", ret, id);
    }
    return ret;
}

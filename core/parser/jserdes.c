/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "jserdes.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

int json_serialize_error(JsonObj** json, int code, const char* str ) {
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!json) {
        return JSON_CREATION_ERR;
    }

    json_object_set_new(*json, JTAG_ERROR, json_object());

    JsonObj* jError = json_object_get(*json, JTAG_ERROR);
    if (jError) {

        json_object_set_new(jError, JTAG_ERROR_CODE, json_integer(code));

        json_object_set_new(jError, JTAG_ERROR_CSTRING, json_string(str));

    } else {
        return JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_version(JsonObj** json, Version* ver) {
    int ret = JSON_ENCODING_OK;

    if (*json) {

        json_object_set_new(*json, JTAG_MAJOR_VERSION, json_integer(ver->major));

        json_object_set_new(*json, JTAG_MINOR_VERSION, json_integer(ver->minor));

    } else {
        return JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_module_cfg(JsonObj** json, ModuleCfg* mCfg, uint8_t count){
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!json) {
        return JSON_CREATION_ERR;
    }

    if (!mCfg) {
        return JSON_NO_VAL_TO_ENCODE;
    }

    json_object_set_new(*json, JTAG_MODULE_CONFIG, json_object());

    JsonObj* jMCfg = json_object_get(*json, JTAG_MODULE_CONFIG);
    if (jMCfg) {

        json_object_set_new(jMCfg, JTAG_NAME, json_string(mCfg->devName));

        json_object_set_new(jMCfg, JTAG_DESCRIPTION,
                        json_string(mCfg->devDesc));

        json_object_set_new(jMCfg, JTAG_TYPE,
                        json_integer(mCfg->devType));

        json_object_set_new(jMCfg, JTAG_CLASS,
                        json_integer(mCfg->devClass));

        json_object_set_new(jMCfg, JTAG_DEV_SYSFS_FILE,
                        json_string(mCfg->sysFile));

    } else {
        return JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_module_info(JsonObj** json, ModuleInfo* mInfo){
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!json) {
        return JSON_CREATION_ERR;
    }

    if (!mInfo) {
        return JSON_NO_VAL_TO_ENCODE;
    }

    json_object_set_new(*json, JTAG_UNIT_INFO, json_object());

    JsonObj* jMInfo = json_object_get(*json, JTAG_UNIT_INFO);
    if (jMInfo) {

        json_object_set_new(jMInfo, JTAG_UUID, json_string(mInfo->uuid));

        json_object_set_new(jMInfo, JTAG_NAME, json_string(mInfo->name));

        json_object_set_new(jMInfo, JTAG_TYPE, json_integer(mInfo->module));

        json_object_set_new(jMInfo, JTAG_PART_NUMBER, json_string(mInfo->partNo));

        json_object_set_new(jMInfo, JTAG_HW_VERSION, json_string(mInfo->hwVer));

        json_object_set_new(jMInfo, JTAG_MAC, json_string(mInfo->mac));

        json_object_set_new(jMInfo, JTAG_PROD_SW_VERSION, json_object());

        JsonObj* jPVer = json_object_get(jMInfo, JTAG_PROD_SW_VERSION);
        if (jPVer) {
            ret = json_serialize_version(&jPVer, &mInfo->pSwVer);
            if (ret != JSON_ENCODING_OK) {
                return ret;
            }
        }

        json_object_set_new(jMInfo, JTAG_SW_VERISION, json_object());

        JsonObj* jSVer = json_object_get(jMInfo, JTAG_SW_VERISION);
        if (jSVer) {
            ret = json_serialize_version(&jSVer, &mInfo->pSwVer);
            if (ret != JSON_ENCODING_OK) {
                return ret;
            }
        }

        json_object_set_new(jMInfo, JTAG_MFG_DATE, json_string(mInfo->mfgDate));

        json_object_set_new(jMInfo, JTAG_MFG_NAME, json_string(mInfo->mfgName));

        json_object_set_new(jMInfo, JTAG_DEVICE_COUNT, json_integer(mInfo->devCount));

    } else {
        return JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_unit_cfg(JsonObj** json, UnitCfg* uCfg, uint8_t count){
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!json) {
        return JSON_CREATION_ERR;
    }

    if (!uCfg) {
        return JSON_NO_VAL_TO_ENCODE;
    }

    json_object_set_new(*json, JTAG_UNIT_CONFIG, json_array());

    JsonObj* jUCfgArr = json_object_get(*json, JTAG_UNIT_CONFIG);
    if (jUCfgArr) {

        for(int iter = 0; iter < count; iter++) {
            json_t* jUCfg = json_object();

            json_object_set_new(jUCfg, JTAG_UUID, json_string(uCfg->modUuid));

            json_object_set_new(jUCfg, JTAG_NAME, json_string(uCfg->modName));

            json_object_set_new(jUCfg, JTAG_INVT_SYSFS_FILE,
                            json_string(uCfg->sysFs));

            /* Add element to array */
            json_array_append(jUCfgArr, jUCfg);
            json_decref(jUCfg);

        }

    } else {
        return JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_unit_info(JsonObj** json, UnitInfo* uInfo){
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!json) {
        return JSON_CREATION_ERR;
    }

    if (!uInfo) {
        return JSON_NO_VAL_TO_ENCODE;
    }

    json_object_set_new(*json, JTAG_UNIT_INFO, json_object());

    JsonObj* jUInfo = json_object_get(*json, JTAG_UNIT_INFO);
    if (jUInfo) {

        json_object_set_new(jUInfo, JTAG_UUID, json_string(uInfo->uuid));

        json_object_set_new(jUInfo, JTAG_NAME, json_string(uInfo->name));

        json_object_set_new(jUInfo, JTAG_TYPE, json_integer(uInfo->unit));

        json_object_set_new(jUInfo, JTAG_PART_NUMBER,
                        json_string(uInfo->partNo));

        json_object_set_new(jUInfo, JTAG_SKEW, json_string(uInfo->skew));

        json_object_set_new(jUInfo, JTAG_MAC, json_string(uInfo->mac));

        json_object_set_new(jUInfo, JTAG_PROD_SW_VERSION, json_object());

        JsonObj* jPVer = json_object_get(jUInfo, JTAG_PROD_SW_VERSION);
        if (jPVer) {
            ret = json_serialize_version(&jPVer, &uInfo->swVer);
            if (ret != JSON_ENCODING_OK) {
                return ret;
            }
        }

        json_object_set_new(jUInfo, JTAG_SW_VERISION, json_object());

        JsonObj* jSVer = json_object_get(jUInfo, JTAG_SW_VERISION);
        if (jSVer) {
            ret = json_serialize_version(&jSVer, &uInfo->pSwVer);
            if (ret != JSON_ENCODING_OK) {
                return ret;
            }
        }

        json_object_set_new(jUInfo, JTAG_ASM_DATE,
                        json_string(uInfo->assmDate));

        json_object_set_new(jUInfo, JTAG_OEM_NAME,
                        json_string(uInfo->oemName));

        json_object_set_new(jUInfo, JTAG_MODULE_COUNT,
                        json_integer(uInfo->modCount));

    } else {
        return JSON_CREATION_ERR;
    }

    return ret;
}

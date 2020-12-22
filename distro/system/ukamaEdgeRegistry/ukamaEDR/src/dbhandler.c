/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "inc/dbhandler.h"

#include "headers/errorcode.h"
#include "headers/globalheader.h"
#include "headers/ubsp/ubsp.h"
#include "inc/regdb.h"
#include "inc/reghelper.h"
#include "inc/registry.h"
#include "headers/utils/log.h"

#include <stdbool.h>

UnitInfo *db_read_unit_info(char *puuid) {
    int ret = 0;
    uint16_t size = 0;
    UnitInfo *uinfo = ubsp_alloc(sizeof(UnitInfo));
    if (uinfo) {
        ret = ubsp_read_unit_info(puuid, uinfo, &size);
        if (!ret) {
            log_debug("DBHANDLER:: Read Unit Info for %s done", puuid);
        } else {
            log_debug("Err(%d): DBHANDLER:: Read Unit Info for %s failed", ret,
                      puuid);
        }
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug("Err(%d): DBHANDLER:: Read Unit Info for %s failed", ret,
                  puuid);
    }

    if (ret) {
        UKAMA_FREE(uinfo);
    }

    return uinfo;
}

UnitCfg *db_read_unit_cfg(char *puuid, uint8_t count) {
    int ret = 0;
    uint16_t size = 0;

    /* Read Unit Cfg */
    UnitCfg *ucfg = ubsp_alloc_unit_cfg(count);
    if (ucfg) {
        size = 0;
        ret = ubsp_read_unit_cfg(puuid, ucfg, count, &size);
        if (!ret) {
            log_debug("DBHANDLER:: Read UKDB Unit Config for %s done", puuid);
        } else {
            log_debug(
                "Err(%d): DBHANDLER:: Read UKDB Unit Config for %s failed", ret,
                puuid);
        }
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug("Err(%d): DBHANDLER:: Read Unit Config for %s failed", ret,
                  puuid);
    }

    if (ret) {
        UKAMA_FREE(ucfg);
    }

    return ucfg;
}

void db_free_unit_cfg(UnitCfg **cfg, uint8_t count) {
    ubsp_free_unit_cfg(*cfg, count);
    *cfg = NULL;
}

ModuleInfo *db_read_module_info(char *puuid) {
    int ret = 0;
    uint16_t size = 0;
    ModuleInfo *minfo = ubsp_alloc(sizeof(ModuleInfo));
    if (minfo) {
        ret = ubsp_read_module_info(puuid, minfo, &size);
        if (!ret) {
            //ukdb_print_module_info(minfo);
            log_debug("DBHANDLER:: Read Module Info for %s", puuid);
        }
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug("Err(%d): DBHANDLER:: Read Module Info for %s failed", ret,
                  puuid);
    }

    if (ret) {
        UKAMA_FREE(minfo);
    }
    return minfo;
}

ModuleCfg *db_read_module_cfg(char *puuid, uint8_t count) {
    int ret = 0;
    uint16_t size = 0;

    /* Read Module Cfg */
    ModuleCfg *mcfg = ubsp_alloc_module_cfg(count);
    if (mcfg) {
        size = 0;
        ret = ubsp_read_module_cfg(puuid, mcfg, count, &size);
        if (!ret) {
            log_debug("DBHANDLER:: Read UKDB Module Config for %s done", puuid);
        } else {
            log_debug("Err(%d): DB Read UKDB Module Config for %s failed", ret,
                      puuid);
        }
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug("Err(%d): DBHANDLER:: Read Module Config for %s failed", ret,
                  puuid);
    }

    if (ret) {
        UKAMA_FREE(mcfg);
    }

    return mcfg;
}

int db_read_inst_data_from_dev(DevObj *obj, PData *sd) {
    int ret = 0;
    /* Send it to UBSP */
    int propid = sd->prop->id;
    ret = ubsp_read_from_prop(obj, &propid, reg_data_value(sd));
    return ret;
}

int db_write_inst_data_to_dev(DevObj *obj, PData *sd) {
    int ret = 0;
    int propid = sd->prop->id;
    ret = ubsp_write_to_prop(obj, &propid, reg_data_value(sd));
    return ret;
}

Property *db_read_dev_property(DevObj *obj, int *pcount) {
    int ret = 0;
    int idx = 0;
    uint16_t count = 0;
    Property *prop = NULL;
    ret = upsb_read_dev_prop_count(obj, &count);
    if (count > 0) {
        prop = malloc(sizeof(Property) * count);
        if (prop) {
            ret = ubsp_read_dev_props(obj, prop);
            if (!ret) {
                print_properties(prop, count);
                log_trace(
                    "DBHANDLER:: %d Property read for Device name %s Disc: %s Module UUID %s.",
                    count, obj->name, obj->disc, obj->mod_UUID);
            } else {
                log_error(
                    "Err(%d): DBHANDLER:: Property read for Device name %s Disc: %s Module UUID %s failed.",
                    ret, obj->name, obj->disc, obj->mod_UUID);
            }
        } else {
            ret = ERR_UBSP_MEMORY_EXHAUSTED;
            log_error(
                "Err(%d): DBHANDLER::Failed in memory allocation for Property read.",
                ret);
        }
    }
    if (ret) {
        UKAMA_FREE(prop);
        *pcount = 0;
    }
    *pcount = count;
    return prop;
}

int db_set_prop_val(PData *pdata, int type, void *val) {
    int ret = 0;
    switch (type) {
    case TYPE_BOOL: {
        pdata->value.boolval = *(bool *)val;
    } break;
    case TYPE_UINT16: {
        pdata->value.sintval = *(uint16_t *)val;
    } break;
    case TYPE_UINT32:
    case TYPE_INT32: {
        pdata->value.intval = *(int *)val;
    } break;
    case TYPE_FLOAT:
    case TYPE_DOUBLE: {
        pdata->value.doubleval = *(double *)val;
    } break;
    case TYPE_STRING: {
        strcpy(pdata->value.stringval, val);
    } break;
    default: {
        ret = -1;
    }
    }
    return ret;
}

int db_read_boolval_prop(DevObj *obj, uint8_t req, bool *mdata, PData *pdata,
                         int *size) {
    int ret = 0;
    if (req) {
        if (pdata && pdata->prop) {
            log_trace("DBHANDLER:: Sending Reading request Resource Id %d "
                      "Property %d for Device %s Disc: %s Module ID %s",
                      pdata->resourceId, pdata->prop->id, obj->name, obj->disc,
                      obj->mod_UUID);
            ret = db_read_inst_data_from_dev(obj, pdata);
        }
        *mdata = pdata->value.boolval;
        *size = sizeof(pdata->value.boolval);
        log_trace("DBHANDLER:: Read resourceID %d = %s for Device Name %s"
                  "Module %s Disc %s",
                  pdata->resourceId, ((*mdata) ? "true" : "false"), obj->name,
                  obj->mod_UUID, obj->disc);
    }
    return ret;
}

int db_read_doubleval_prop(DevObj *obj, uint8_t req, double *mdata,
                           PData *pdata, int *size) {
    int ret = 0;
    if (req) {
        if (pdata && pdata->prop) {
            log_trace("DBHANDLER:: Sending Reading request Resource Id %d "
                      "Property %d for Device %s Disc: %s Module ID %s",
                      pdata->resourceId, pdata->prop->id, obj->name, obj->disc,
                      obj->mod_UUID);
            ret = db_read_inst_data_from_dev(obj, pdata);
        }
        *mdata = pdata->value.doubleval;
        *size = sizeof(pdata->value.doubleval);
        log_trace("DBHANDLER:: Read resourceID %d = %lf for Device Name %s"
                  "Module %s Disc %s",
                  pdata->resourceId, *mdata, obj->name, obj->mod_UUID,
                  obj->disc);
    }
    return ret;
}

int db_read_intval_prop(DevObj *obj, uint8_t req, int *mdata, PData *pdata,
                        int *size) {
    int ret = 0;
    if (req) {
        if (pdata && pdata->prop) {
            log_trace("DBHANDLER:: Sending Reading request Resource Id %d "
                      "Property %d for Device %s Disc: %s Module ID %s",
                      pdata->resourceId, pdata->prop->id, obj->name, obj->disc,
                      obj->mod_UUID);
            ret = db_read_inst_data_from_dev(obj, pdata);
        }
        *mdata = pdata->value.intval;
        *size = sizeof(pdata->value.intval);
        log_trace("DBHANDLER:: Read resourceID %d = %d for Device Name %s"
                  "Module %s Disc %s",
                  pdata->resourceId, *mdata, obj->name, obj->mod_UUID,
                  obj->disc);
    }
    return ret;
}

int db_read_longintval_prop(DevObj *obj, uint8_t req, int64_t *mdata,
                            PData *pdata, int *size) {
    int ret = 0;
    if (req) {
        if (pdata && pdata->prop) {
            log_trace("DBHANDLER:: Sending Reading request Resource Id %d "
                      "Property %d for Device %s Disc: %s Module ID %s",
                      pdata->resourceId, pdata->prop->id, obj->name, obj->disc,
                      obj->mod_UUID);
            ret = db_read_inst_data_from_dev(obj, pdata);
        }
        *mdata = pdata->value.lintval;
        *size = sizeof(pdata->value.lintval);
        log_trace("DBHANDLER:: Read resourceID %d = %d for Device Name %s"
                  "Module %s Disc %s",
                  pdata->resourceId, *mdata, obj->name, obj->mod_UUID,
                  obj->disc);
    }
    return ret;
}

int db_read_shortintval_prop(DevObj *obj, uint8_t req, uint16_t *mdata,
                             PData *pdata, int *size) {
    int ret = 0;
    if (req) {
        if (pdata && pdata->prop) {
            log_trace("DBHANDLER:: Sending Reading request Resource Id %d "
                      "Property %d for Device %s Disc: %s Module ID %s",
                      pdata->resourceId, pdata->prop->id, obj->name, obj->disc,
                      obj->mod_UUID);
            ret = db_read_inst_data_from_dev(obj, pdata);
        }
        *mdata = pdata->value.sintval;
        *size = sizeof(pdata->value.sintval);
        log_trace("DBHANDLER:: Read resourceID %d = %d for Device Name %s"
                  "Module %s Disc %s",
                  pdata->resourceId, *mdata, obj->name, obj->mod_UUID,
                  obj->disc);
    }
    return ret;
}

int db_read_strval_prop(DevObj *obj, uint8_t req, char *mdata, PData *pdata,
                        int *size) {
    int ret = 0;
    if (req) {
        if (pdata && pdata->prop) {
            log_trace("DBHANDLER:: Sending Reading request Resource Id %d "
                      "Property %d for Device %s Disc: %s Module ID %s",
                      pdata->resourceId, pdata->prop->id, obj->name, obj->disc,
                      obj->mod_UUID);
            ret = db_read_inst_data_from_dev(obj, pdata);
        }
        memset(mdata, '\0', MAX_LWM2M_OBJ_STR_LEN);
        memcpy(mdata, pdata->value.stringval, strlen(pdata->value.stringval));
        *size = sizeof(char) * MAX_LWM2M_OBJ_STR_LEN;
        log_trace("DBHANDLER:: Read resourceID %d = %s for Device Name %s"
                  "Module %s Disc %s",
                  pdata->resourceId, mdata, obj->name, obj->mod_UUID,
                  obj->disc);
    }
    return ret;
}

int db_write_boolval_prop(DevObj *obj, uint8_t req, bool mdata, PData *pdata) {
    int ret = 0;
    if (req) {
        if (pdata && pdata->prop) {
            pdata->value.boolval = mdata;
            log_trace("DBHANDLER:: Sending Write request Resource Id %d = %s"
                      "Property %d for Device %s Disc: %s Module ID %s",
                      pdata->resourceId, ((mdata) ? "true" : "false"),
                      obj->name, obj->disc, obj->mod_UUID);
            ret = db_write_inst_data_to_dev(obj, pdata);
        }
    }
    return ret;
}

int db_write_intval_prop(DevObj *obj, uint8_t req, int mdata, PData *pdata) {
    int ret = 0;
    if (req) {
        if (pdata && pdata->prop) {
            pdata->value.intval = mdata;
            log_trace("DBHANDLER:: Sending Write request Resource Id %d = %d"
                      "Property %d for Device %s Disc: %s Module ID %s",
                      pdata->resourceId, mdata, obj->name, obj->disc,
                      obj->mod_UUID);
            ret = db_write_inst_data_to_dev(obj, pdata);
        }
    }
    return ret;
}

int db_write_shortintval_prop(DevObj *obj, uint8_t req, uint16_t mdata,
                              PData *pdata) {
    int ret = 0;
    if (req) {
        if (pdata && pdata->prop) {
            pdata->value.sintval = mdata;
            log_trace("DBHANDLER:: Sending Write request Resource Id %d = %d"
                      "Property %d for Device %s Disc: %s Module ID %s",
                      pdata->resourceId, mdata, obj->name, obj->disc,
                      obj->mod_UUID);
            ret = db_write_inst_data_to_dev(obj, pdata);
        }
    }
    return ret;
}

int db_write_longintval_prop(DevObj *obj, uint8_t req, int64_t mdata,
                             PData *pdata) {
    int ret = 0;
    if (req) {
        if (pdata && pdata->prop) {
            pdata->value.lintval = mdata;
            log_trace("DBHANDLER:: Sending Write request Resource Id %d = %d"
                      "Property %d for Device %s Disc: %s Module ID %s",
                      pdata->resourceId, mdata, obj->name, obj->disc,
                      obj->mod_UUID);
            ret = db_write_inst_data_to_dev(obj, pdata);
        }
    }
    return ret;
}

int db_write_doubleval_prop(DevObj *obj, uint8_t req, double mdata,
                            PData *pdata) {
    int ret = 0;
    if (req) {
        if (pdata && pdata->prop) {
            pdata->value.doubleval = mdata;
            log_trace("DBHANDLER:: Sending Write request Resource Id %d = %lf"
                      "Property %d for Device %s Disc: %s Module ID %s",
                      pdata->resourceId, mdata, obj->name, obj->disc,
                      obj->mod_UUID);
            ret = db_write_inst_data_to_dev(obj, pdata);
        }
    }
    return ret;
}

int db_write_strval_prop(DevObj *obj, uint8_t req, char *mdata, PData *pdata) {
    int ret = 0;
    if (req) {
        if (pdata && pdata->prop) {
            if (mdata) {
                memset(&pdata->value.stringval, '\0', 32);
                memcpy(&pdata->value.stringval, mdata, MIN(31, strlen(mdata)));
                log_trace(
                    "DBHANDLER:: Sending Write request Resource Id %d = %lf"
                    "Property %d for Device %s Disc: %s Module ID %s",
                    pdata->resourceId, *mdata, obj->name, obj->disc,
                    obj->mod_UUID);
                ret = db_write_inst_data_to_dev(obj, pdata);
            } else {
                ret = -1;
            }
        }
    }
    return ret;
}

void db_versiontostr(Version ver, char *str) {
    sprintf(str, "%d.%d", ver.major, ver.minor);
}

int db_register_alarm_callback(DevObj *obj, int *prop, CallBackFxn fn) {
    int ret = 0;
    log_debug(
        "DBHANDLER:: Registering Callback function for Device name %s Disc: %s Module UUID %s.",
        obj->name, obj->disc, obj->mod_UUID);
    ret = ubsp_register_app_cb(obj, prop, fn);
    return ret;
}

int db_enable_alarm(DevObj *obj, int prop) {
    int ret = 0;
    void *data = NULL;
    ret = ubsp_enable_irq(obj, &prop);
    return ret;
}

/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "headers/ubsp/ubsp.h"

#include "headers/errorcode.h"
#include "inc/devicedb.h"
#include "inc/ukdb.h"
#include "headers/utils/log.h"
#include "utils/jcreator.h"

void ubsp_banner() {
    log_trace(
        "\t||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||");
    log_trace(
        "\t||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||");
    log_trace(
        "\t||||||||  |||||  ||||         ||||         ||||         ||||||||");
    log_trace(
        "\t||||||||  |||||  ||||  |||||  ||||  |||||||||||  |||||  ||||||||");
    log_trace(
        "\t||||||||  |||||  ||||  |||||  ||||  |||||||||||  |||||  ||||||||");
    log_trace(
        "\t||||||||  |||||  ||||         ||||         ||||         ||||||||");
    log_trace(
        "\t||||||||  |||||  ||||  |||||  |||||||||||  ||||  |||||||||||||||");
    log_trace(
        "\t||||||||  |||||  ||||  |||||  |||||||||||  ||||  |||||||||||||||");
    log_trace(
        "\t||||||||         ||||         ||||         ||||  |||||||||||||||");
    log_trace(
        "\t||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||");
    log_trace(
        "\t||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||");
}

void *ubsp_alloc(size_t size) {
    void *mem = NULL;
    if (size > 0) {
        mem = malloc(sizeof(char) * size);
        if (mem) {
            memset(mem, '\0', sizeof(char) * size);
        }
    }
    return mem;
}

void ubsp_free(void *mem) {
    if (mem) {
        free(mem);
    }
}

void ubsp_free_unit_cfg(UnitCfg *cfg, uint8_t module_count) {
    if (cfg) {
        ukdb_free_unit_cfg(cfg, module_count);
    }
}

UnitCfg *ubsp_alloc_unit_cfg(uint8_t module_count) {
    return ukdb_alloc_unit_cfg(module_count);
}

void ubsp_free_module_cfg(ModuleCfg *cfg, uint8_t dev_count) {
    if (cfg) {
        ukdb_free_module_cfg(cfg, dev_count);
    }
}

ModuleCfg *ubsp_alloc_module_cfg(uint8_t dev_count) {
    return ukdb_alloc_module_cfg(dev_count);
}

int ubsp_devdb_init(void *data) {
    int ret = 0;
    ret = devdb_init(data);
    if (ret) {
        log_error("UBSP(%d):: Failed while initializing Ukama Device DB.", ret);
    }
    return ret;
}

int ubsp_ukdb_init(char *sys_db_path) {
    int ret = 0;
    ret = ukdb_init(sys_db_path);
    return ret;
}

int ubsp_exit() {
    int ret = 0;
    log_warn("UBSP:: Cleaning UBSP UKDB and DeviceDB.");
    devdb_exit();
    ukdb_exit();
    return ret;
}

int ubsp_register_module(UnitCfg *cfg) {
    int ret = 0;
    ret = ukdb_register_module(cfg);
    if (ret) {
        log_error(
            "UBSP(%d):: Failed while registering new Module uuid %s to DB.",
            ret, cfg->mod_uuid);
    }
    return ret;
}

int ubsp_unregister_module(char *puuid) {
    int ret = 0;
    ret = ukdb_unregister_module(puuid);
    if (ret) {
        log_error("UBSP(%d):: Failed while unregistering Module uuid %s to DB.",
                  ret, puuid);
    }
    return ret;
}

int ubsp_idb_init(void *data) {
    int ret = 0;
    ret = ukdb_idb_init(data);
    if (ret) {
        log_error("UBSP(%d):: Failed while initializing Ukama Device DB.", ret);
    }
    return ret;
}

void ubsp_idb_exit() {
    ukdb_idb_exit();
}

int ubsp_read_header(char *puuid, UKDBHeader *pheader) {
    int ret = -1;
    if (pheader) {
        ret = ukdb_read_header(puuid, pheader);
        if (ret) {
            log_error("UBSP(%d):: Failed to read header for module %s"
                      " Ukama DB.",
                      ret, puuid);
        }
    } else {
        ret = ERR_UBSP_INVALID_POINTER;
    }
    return ret;
}

int ubsp_validating_magicword(char *puuid) {
    int ret = 0;
    ret = ukdb_validating_magicword(puuid);
    if (ret) {
        log_error("UBSP(%d):: Failed to validate magic word for module %s"
                  " Ukama DB.",
                  ret, puuid);
    }
    return ret;
}

int ubsp_read_dbversion(char *puuid, Version *pver) {
    int ret = 0;
    ret = ukdb_read_dbversion(puuid, pver);
    if (ret) {
        log_error("UBSP(%d):: Failed to read DB version for module %s"
                  " Ukama DB.",
                  ret, puuid);
    }
    return ret;
}

int ubsp_update_dbversion(char *puuid, Version ver) {
    int ret = 0;
    ret = ukdb_update_dbversion(puuid, ver);
    if (ret) {
        log_error("UBSP(%d):: Failed to update DB version for module %s"
                  " Ukama DB.",
                  ret, puuid);
    }
    return ret;
}

int ubsp_read_unit_info(char *puuid, UnitInfo *pdata, uint16_t *psize) {
    int ret = 0;
    ret = ukdb_read_unit_info(puuid, pdata, psize);
    if (ret) {
        log_error("UBSP(%d):: Failed to read unit info from %s"
                  " Ukama DB.",
                  ret, puuid);
    }
    return ret;
}

int ubsp_read_unit_cfg(char *puuid, UnitCfg *pucfg, uint8_t count,
                       uint16_t *psize) {
    int ret = 0;
    ret = ukdb_read_unit_cfg(puuid, pucfg, count, psize);
    if (ret) {
        log_error("UBSP(%d):: Failed to read unit config from %s"
                  " Ukama DB.",
                  ret, puuid);
    }
    return ret;
}
#if 0
int ubsp_create_unit_schema(UnitInfo *unit_info, UnitCfg *unit_cfg, char* junit_schema) {
    int ret = 0;
    ret = jcreator_unit_schema(unit_info, unit_cfg, junit_schema);
    if (ret) {
        log_error("UBSP(%d):: Failed to create Unit schema",
                  ret);
    }
    return ret;
}
#endif
int ubsp_read_module_info(char *puuid, ModuleInfo *pinfo, uint16_t *psize) {
    int ret = 0;
    ret = ukdb_read_module_info(puuid, pinfo, psize);
    if (ret) {
        log_error("UBSP(%d):: Failed to read module info from %s"
                  " Ukama DB.",
                  ret, puuid);
    }
    return ret;
}

int ubsp_read_module_cfg(char *puuid, ModuleCfg *pmcfg, uint8_t count,
                         uint16_t *psize) {
    int ret = 0;
    ret = ukdb_read_module_cfg(puuid, pmcfg, count, psize);
    if (ret) {
        log_error("UBSP(%d):: Failed to read module config from %s"
                  " Ukama DB.",
                  ret, puuid);
    }
    return ret;
}
#if 0
int jcreator_module_schema(ModuleInfo *mod_info, char* jmod_schema) {
    int ret = 0;
    ret = jcreator_module_schema( mod_info, jmod_schema);
    if (ret) {
        log_error("UBSP(%d):: Failed to create Module schema",
                  ret);
    }
    return ret;
}
#endif
int ubsp_create_schema(UnitInfo *unit_info, UnitCfg *unit_cfg,
                       ModuleInfo *mod_info, char **j_schema) {
    int ret = 0;
    ret = jcreator_schema(unit_info, unit_cfg, mod_info, j_schema);
    if (ret) {
        log_error("UBSP(%d):: Failed to create Module schema", ret);
    }
    return ret;
}

int ubsp_read_fact_config(char *puuid, void *pdata, uint16_t *psize) {
    int ret = 0;
    ret = ukdb_read_fact_config(puuid, pdata, psize);
    if (ret) {
        log_error("UBSP(%d):: Failed to read fact config from %s"
                  " Ukama DB.",
                  ret, puuid);
    }
    return ret;
}

int ubsp_read_user_config(char *puuid, void *pdata, uint16_t *psize) {
    int ret = 0;
    ret = ukdb_read_user_config(puuid, pdata, psize);
    if (ret) {
        log_error("UBSP(%d):: Failed to read user config from %s"
                  " Ukama DB.",
                  ret, puuid);
    }
    return ret;
}

int ubsp_read_fact_calib(char *puuid, void *pdata, uint16_t *psize) {
    int ret = 0;
    ret = ukdb_read_fact_calib(puuid, pdata, psize);
    if (ret) {
        log_error("UBSP(%d):: Failed to read fact calibration from %s"
                  " Ukama DB.",
                  ret, puuid);
    }
    return ret;
}

int ubsp_read_user_calib(char *puuid, void *pdata, uint16_t *psize) {
    int ret = 0;
    ret = ukdb_read_user_calib(puuid, pdata, psize);
    if (ret) {
        log_error("UBSP(%d):: Failed to read user calibration from %s"
                  " Ukama DB.",
                  ret, puuid);
    }
    return ret;
}

int ubsp_read_bs_certs(char *puuid, void *pdata, uint16_t *psize) {
    int ret = 0;
    ret = ukdb_read_bs_certs(puuid, pdata, psize);
    if (ret) {
        log_error("UBSP(%d):: Failed to read bootstrap certs from %s"
                  " Ukama DB.",
                  ret, puuid);
    }
    return ret;
}

int ubsp_read_lwm2m_certs(char *puuid, void *pdata, uint16_t *psize) {
    int ret = 0;
    ret = ukdb_read_lwm2m_certs(puuid, pdata, psize);
    if (ret) {
        log_error("UBSP(%d):: Failed to read lwm2m certs from %s"
                  " Ukama DB.",
                  ret, puuid);
    }
    return ret;
}

int ubsp_pre_create_ukdb_hook(char *mod_uuid) {
    int ret = 0;
    ret = ukdb_pre_create_db_setup(mod_uuid);
    if (ret) {
        log_error(
            "UBSP(%d):: Failed to execute pre DB creation hook or Module UUID %s",
            ret, mod_uuid);
    }
    return ret;
}

int ubsp_create_ukdb(char *mod_uuid) {
    int ret = 0;
    ret = ukdb_create_db(mod_uuid);
    if (ret) {
        log_error("UBSP(%d):: Failed to create UKDB for Module UUID %s", ret,
                  mod_uuid);
    }
    return ret;
}

int ubsp_remove_ukdb(char *mod_uuid) {
    int ret = 0;
    ret = ukdb_remove_db(mod_uuid);
    if (ret) {
        log_error("UBSP(%d):: Failed to remove UKDB for Module UUID %s", ret,
                  mod_uuid);
    }
    return ret;
}

int ubsp_read_registered_dev_count(DeviceType type, uint16_t *count) {
    int ret = 0;
    ret = devdb_read_reg_dev_count(type, count);
    if (ret) {
        log_error("UBSP(%d):: Failed to read registered device count from"
                  " Device DB of type 0x%x.",
                  ret, type);
    }
    return ret;
}

int ubsp_read_registered_dev(DeviceType type, Device *dev) {
    int ret = 0;
    ret = devdb_read_reg_dev(type, dev);
    if (ret) {
        log_error("UBSP(%d):: Failed to read registered devices from"
                  " Device DB of type 0x%x.",
                  ret, type);
    }
    return ret;
}

int upsb_read_dev_prop_count(DevObj *obj, uint16_t *count) {
    int ret = 0;
    ret = devdb_read_prop_count(obj, count);
    if (ret) {
        log_error("UBSP(%d):: Failed to read property count from"
                  " Device Name %s, Disc: %s Module UUID %s.",
                  ret, obj->name, obj->disc, obj->mod_UUID);
    }
    return ret;
}

int ubsp_read_dev_props(DevObj *obj, Property *prop) {
    int ret = 0;
    ret = devdb_read_prop(obj, prop);
    if (ret) {
        log_error("UBSP(%d):: Failed to read property from"
                  " Device Name %s, Disc: %s Module UUID %s.",
                  ret, obj->name, obj->disc, obj->mod_UUID);
    }
    return ret;
}

int ubsp_read_from_prop(DevObj *obj, int *prop, void *data) {
    int ret = 0;
    ret = devdb_read(obj, prop, data);
    if (ret) {
        log_error("UBSP(%d):: Failed to read from property[%d] of"
                  " Device Name %s, Disc: %s Module UUID %s.",
                  ret, *prop, obj->name, obj->disc, obj->mod_UUID);
    }
    return ret;
}

int ubsp_write_to_prop(DevObj *obj, int *prop, void *data) {
    int ret = 0;
    ret = devdb_write(obj, prop, data);
    if (ret) {
        log_error("UBSP(%d):: Failed to write to property[%d] of"
                  " Device Name %s, Disc: %s Module UUID %s.",
                  ret, *prop, obj->name, obj->disc, obj->mod_UUID);
    }
    return ret;
}

int ubsp_enable(DevObj *obj, void *data) {
    /*TODO Still need to check use case */
    log_debug("UBSP:: TODO task.");
    return 0;
}

int ubsp_disable(DevObj *obj, void *data) {
    /*TODO Still need to check use case */
    log_debug("UBSP:: TODO task.");
    return 0;
}

int ubsp_enable_irq(DevObj *obj, int *prop) {
    int ret = 0;
    /* TODO: remove value*/
    int value = 0;
    ret = devdb_enable_irq(obj, prop, &value);
    if (ret) {
        log_debug("UBSP(%d):: Failed to enable alert for "
                  "Device name %s Disc: %s Module UUID %s ",
                  ret, obj->name, obj->disc, obj->mod_UUID);
    }
    return ret;
}

int ubsp_disable_irq(DevObj *obj, int *prop) {
    int ret = 0;
    /* TODO: remove value */
    int value = 0;
    ret = devdb_disable_irq(obj, prop, &value);
    if (ret) {
        log_debug("UBSP(%d):: Failed to disable alert for "
                  "Device name %s Disc: %s Module UUID %s ",
                  ret, obj->name, obj->disc, obj->mod_UUID);
    }
    return ret;
}

int ubsp_register_app_cb(DevObj *obj, void *prop, CallBackFxn fn) {
    int ret = 0;
    ret = devdb_reg_app_cb(obj, prop, fn);
    if (ret) {
        log_debug("UBSP(%d):: Failed to registering Callback function for "
                  "Device name %s Disc: %s Module UUID %s.",
                  ret, obj->name, obj->disc, obj->mod_UUID);
    }
    return ret;
}

int ubsp_deregister_app_cb(DevObj *obj, void *prop, CallBackFxn fn) {
    int ret = 0;
    ret = devdb_dereg_app_cb(obj, prop, fn);
    if (ret) {
        log_debug("UBSP(%d):: Failed to de-registering Callback function for "
                  "Device name %s Disc: %s Module UUID %s.",
                  ret, obj->name, obj->disc, obj->mod_UUID);
    }
    return ret;
}

/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "utils/jcreator.h"

#include "headers/ubsp/devices.h"
#include "headers/errorcode.h"
#include "inc/globalheader.h"
#include "headers/utils/log.h"
#include "headers/ubsp/property.h"

#include "utils/cJSON.h"

#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int jcreator_db_verison_schema(Version ver, cJSON **version_schema) {
    int ret = 0;
    if (cJSON_AddNumberToObject(*version_schema, "major", ver.major) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        return ret;
    }

    if (cJSON_AddNumberToObject(*version_schema, "minor", ver.minor) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        return ret;
    }
    return ret;
}

int jcreator_gpio_schema(DevGpioCfg *cfg, cJSON **dev_schema) {
    int ret = 0;
    char dir[7] = { '\0' };
    if (cfg) {
        if (cfg->direction) {
            memcpy(&dir, "input", 5);
        } else {
            memcpy(&dir, "output", 6);
        }
        if (cJSON_AddStringToObject(*dev_schema, "direction", dir) == NULL) {
            ret = ERR_UBSP_CRT_JSON_SCHEMA;
            return ret;
        }
        if (cJSON_AddNumberToObject(*dev_schema, "number", cfg->gpio_num) ==
            NULL) {
            ret = ERR_UBSP_CRT_JSON_SCHEMA;
            return ret;
        }
    } else {
        *dev_schema = NULL;
    }
    return ret;
}

int jcreator_i2c_schema(DevI2cCfg *cfg, cJSON **dev_schema) {
    int ret = 0;
    if (cJSON_AddNumberToObject(*dev_schema, "bus", cfg->bus) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        return ret;
    }

    if (cJSON_AddNumberToObject(*dev_schema, "address", cfg->add) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        return ret;
    }
    return ret;
}

int jcreator_spi_schema(DevSpiCfg *cfg, cJSON **dev_schema) {
    int ret = 0;
    cJSON *jcs_schema = cJSON_CreateObject();

    ret = jcreator_gpio_schema(&cfg->cs, &jcs_schema);
    if (ret) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        return ret;
    }
    cJSON_AddItemToObject(*dev_schema, "cs", jcs_schema);

    if (cJSON_AddNumberToObject(*dev_schema, "bus", cfg->bus) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        return ret;
    }

    return ret;
}

int jcreator_uart_schema(DevUartCfg *cfg, cJSON **dev_schema) {
    int ret = 0;
    if (cJSON_AddNumberToObject(*dev_schema, "uartno", cfg->uartno) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        return ret;
    }
    return ret;
}

cJSON *jcreator_dev_schema(uint16_t class, void *pcfg) {
    int ret = 0;
    cJSON *jdev_schema = cJSON_CreateObject();
    switch (class) {
    case DEV_CLASS_GPIO:
        ret = jcreator_gpio_schema(pcfg, &jdev_schema);
        break;
    case DEV_CLASS_I2C:
        ret = jcreator_i2c_schema(pcfg, &jdev_schema);
        break;
    case DEV_CLASS_SPI:
        ret = jcreator_spi_schema(pcfg, &jdev_schema);
        break;
    case DEV_CLASS_UART:
        ret = jcreator_uart_schema(pcfg, &jdev_schema);
        break;
    default:
        ret = ERR_UBSP_INVALID_DEVICE_CFG;
        log_error("Err(%d): JSONCREATOR:: Unkown device type failed to create.",
                  ret);
    }

    if (ret) {
        cJSON_Delete(jdev_schema);
        jdev_schema = NULL;
    }
    return jdev_schema;
}

int jcreator_unitcfg_schema(UnitCfg *unit_cfg, int count, cJSON **junit_cfgs) {
    int ret = 0;
    cJSON *junit_cfg = NULL;
    /* Array of Modules present in unit.*/
    for (int iter = 0; iter < count; iter++) {
        junit_cfg = cJSON_CreateObject();
        if (cJSON_AddStringToObject(junit_cfg, "UUID",
                                    unit_cfg[iter].mod_uuid) == NULL) {
            ret = ERR_UBSP_CRT_JSON_SCHEMA;
            goto cleanup;
        }

        if (cJSON_AddStringToObject(junit_cfg, "name",
                                    unit_cfg[iter].mod_name) == NULL) {
            ret = ERR_UBSP_CRT_JSON_SCHEMA;
            goto cleanup;
        }

        if (cJSON_AddStringToObject(junit_cfg, "dbsysfs",
                                    unit_cfg[iter].sysfs) == NULL) {
            ret = ERR_UBSP_CRT_JSON_SCHEMA;
            goto cleanup;
        }

        cJSON *jeeprom_cfg = cJSON_CreateObject();
        if (jcreator_i2c_schema(unit_cfg[iter].eeprom_cfg, &jeeprom_cfg)) {
            ret = ERR_UBSP_CRT_JSON_SCHEMA;
            goto cleanup;
        }
        cJSON_AddItemToObject(junit_cfg, "devicedb", jeeprom_cfg);

        cJSON_AddItemToArray(*junit_cfgs, junit_cfg);
    }
    if (!ret) {
        log_debug("JCREATOR: Created Unit Config schema.");
    }
cleanup:
    if (ret) {
        log_error("Err(%d):: JCREATOR: Failed to create Unit COnfig schema.",
                  ret);
        cJSON_Delete(junit_cfg);
    }
    return ret;
}

cJSON *jcreator_unit_schema(UnitInfo *unit_info) {
    int ret = 0;
    char *string = NULL;
    cJSON *junit_info = cJSON_CreateObject();

    /* Adding Unit Info */
    /* junit_info = cJSON_AddItemToObject(junit_schema, "unit_info");
    if (junit_info == NULL) {
    	ret = ERR_UBSP_CRT_JSON_SCHEMA;
    	goto cleanup_unitinfo_schema;
    }
    */

    if (cJSON_AddStringToObject(junit_info, "UUID", unit_info->uuid) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    if (cJSON_AddStringToObject(junit_info, "name", unit_info->name) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    if (cJSON_AddNumberToObject(junit_info, "type", unit_info->unit) == NULL) {
        goto cleanup;
    }

    if (cJSON_AddStringToObject(junit_info, "partno", unit_info->partno) ==
        NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    if (cJSON_AddStringToObject(junit_info, "skew", unit_info->skew) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    if (cJSON_AddStringToObject(junit_info, "mac", unit_info->mac) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    cJSON *jsw_version = cJSON_CreateObject();
    if (jcreator_db_verison_schema(unit_info->swver, &jsw_version)) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }
    cJSON_AddItemToObject(junit_info, "sw_version", jsw_version);

    cJSON *jproduction_sw_version = cJSON_CreateObject();
    if (jcreator_db_verison_schema(unit_info->pswver,
                                   &jproduction_sw_version)) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }
    cJSON_AddItemToObject(junit_info, "production_sw_version",
                          jproduction_sw_version);

    if (cJSON_AddStringToObject(junit_info, "assembly_date",
                                unit_info->assm_date) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    if (cJSON_AddStringToObject(junit_info, "OEM_name", unit_info->oem_name) ==
        NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    if (cJSON_AddNumberToObject(junit_info, "module_count",
                                unit_info->mod_count) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    string = cJSON_Print(junit_info);
    if (string == NULL) {
        fprintf(stderr, "Failed to print monitor.\n");
    } else {
        log_debug("JCREATOR: Created Unit Info schema is : \n %s \n.", string);
        free(string);
    }

    if (!ret) {
        log_debug("JCREATOR: Created Unit Info schema.");
    }

cleanup:
    if (ret) {
        log_error("Err(%d):: JCREATOR: Failed to create Unit Info schema.",
                  ret);
        cJSON_Delete(junit_info);
        junit_info = NULL;
    }

    return junit_info;
}

int jcreator_modulecfg_schema(ModuleCfg *mod_cfg, int count,
                              cJSON **jmod_cfgs) {
    int ret = 0;
    cJSON *jmod_cfg = NULL;
    /* Array of Modules present in unit.*/
    for (int iter = 0; iter < count; iter++) {
        jmod_cfg = cJSON_CreateObject();
        if (cJSON_AddStringToObject(jmod_cfg, "name", mod_cfg[iter].dev_name) ==
            NULL) {
            ret = ERR_UBSP_CRT_JSON_SCHEMA;
            goto cleanup;
        }

        if (cJSON_AddStringToObject(jmod_cfg, "description",
                                    mod_cfg[iter].dev_disc) == NULL) {
            ret = ERR_UBSP_CRT_JSON_SCHEMA;
            goto cleanup;
        }

        if (cJSON_AddNumberToObject(jmod_cfg, "type", mod_cfg[iter].dev_type) ==
            NULL) {
            ret = ERR_UBSP_CRT_JSON_SCHEMA;
            goto cleanup;
        }

        if (cJSON_AddNumberToObject(jmod_cfg, "class",
                                    mod_cfg[iter].dev_class) == NULL) {
            ret = ERR_UBSP_CRT_JSON_SCHEMA;
            goto cleanup;
        }

        if (cJSON_AddStringToObject(jmod_cfg, "devsysfs",
                                    mod_cfg[iter].sysfile) == NULL) {
            ret = ERR_UBSP_CRT_JSON_SCHEMA;
            goto cleanup;
        }

        cJSON *jdev_cfg =
            jcreator_dev_schema(mod_cfg[iter].dev_class, mod_cfg[iter].cfg);
        if (jdev_cfg == NULL) {
            ret = ERR_UBSP_CRT_JSON_SCHEMA;
            log_warn("JCREATOR: No dev schema found for ModuleCfg[%d].", iter);
        } else {
            cJSON_AddItemToObject(jmod_cfg, "dev_hwattrs", jdev_cfg);
        }

        cJSON_AddItemToArray(*jmod_cfgs, jmod_cfg);
    }

    if (!ret) {
        log_debug("JCREATOR: Created module config schema.");
    }

cleanup:
    if (ret) {
        log_error("Err(%d):: JCREATOR: Failed to create module config schema.",
                  ret);
        cJSON_Delete(jmod_cfg);
    }
    return ret;
}

cJSON *jcreator_module_schema(ModuleInfo *mod_info) {
    int ret = 0;
    char *string = NULL;
    cJSON *jmod_info = cJSON_CreateObject();

    if (cJSON_AddStringToObject(jmod_info, "UUID", mod_info->uuid) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    if (cJSON_AddStringToObject(jmod_info, "name", mod_info->name) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    if (cJSON_AddNumberToObject(jmod_info, "type", mod_info->module) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    if (cJSON_AddStringToObject(jmod_info, "partno", mod_info->partno) ==
        NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    if (cJSON_AddStringToObject(jmod_info, "hw_version", mod_info->hwver) ==
        NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    if (cJSON_AddStringToObject(jmod_info, "mac", mod_info->mac) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    cJSON *jsw_version = cJSON_CreateObject();
    if (jcreator_db_verison_schema(mod_info->swver, &jsw_version)) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }
    cJSON_AddItemToObject(jmod_info, "sw_version", jsw_version);

    cJSON *jproduction_sw_version = cJSON_CreateObject();
    if (jcreator_db_verison_schema(mod_info->pswver, &jproduction_sw_version)) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }
    cJSON_AddItemToObject(jmod_info, "production_sw_version",
                          jproduction_sw_version);

    if (cJSON_AddStringToObject(jmod_info, "manufacturing_date",
                                mod_info->mfg_date) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    if (cJSON_AddStringToObject(jmod_info, "manufacturer_name",
                                mod_info->mfg_name) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    if (cJSON_AddNumberToObject(jmod_info, "device_count",
                                mod_info->dev_count) == NULL) {
        ret = ERR_UBSP_CRT_JSON_SCHEMA;
        goto cleanup;
    }

    string = cJSON_Print(jmod_info);
    if (string == NULL) {
        fprintf(stderr, "Failed to print monitor.\n");
    } else {
        log_debug("JCREATOR: Created module Info schema is : \n %s \n.",
                  string);
        free(string);
    }

    if (!ret) {
        log_debug("JCREATOR: Created module info schema.");
    }

cleanup:
    if (ret) {
        log_error("Err(%d):: JCREATOR: Failed to create module info schema.",
                  ret);
        cJSON_Delete(jmod_info);
        jmod_info = NULL;
    }

    return jmod_info;
}

int jcreator_schema(UnitInfo *unit_info, UnitCfg *unit_cfg,
                    ModuleInfo *mod_info, char **schema) {
    int ret = 0;
    cJSON *junit_cfgs = NULL;
    cJSON *jmod_cfgs = NULL;
    cJSON *jmod = NULL;
    cJSON *jschema = cJSON_CreateObject();
    uint8_t mod_count = 1; /* At least one module is present */
    if (unit_info && unit_cfg) {
        mod_count = unit_info->mod_count;
        /* Adding Unit Info */
        cJSON *junit_schema = jcreator_unit_schema(unit_info);
        if (junit_schema == NULL) {
            ret = ERR_UBSP_CRT_JSON_SCHEMA;
            goto cleanup;
        }
        cJSON_AddItemToObject(jschema, "unit_info", junit_schema);

        /* Adding Unit Config. */
        junit_cfgs = cJSON_AddArrayToObject(jschema, "unit_config");
        if (junit_cfgs == NULL) {
            ret = ERR_UBSP_CRT_JSON_SCHEMA;
            goto cleanup;
        }

        if (jcreator_unitcfg_schema(unit_cfg, mod_count, &junit_cfgs)) {
            ret = ERR_UBSP_CRT_JSON_SCHEMA;
            goto cleanup;
        }
    }

    if (mod_info) {
        /* This is array of modules */
        cJSON *jmods = cJSON_AddArrayToObject(jschema, "modules");

        /* Array of Modules present in unit.*/
        for (int iter = 0; iter < mod_count; iter++) {
            cJSON *jmod = cJSON_CreateObject();
            /* Adding Module Info */
            cJSON *jmod_schema = jcreator_module_schema(&mod_info[iter]);
            if (jmod_schema == NULL) {
                ret = ERR_UBSP_CRT_JSON_SCHEMA;
                goto cleanup;
            }
            cJSON_AddItemToObject(jmod, "module_info", jmod_schema);

            if (mod_info[iter].module_cfg) {
                /* Adding Module Config. */
                jmod_cfgs = cJSON_AddArrayToObject(jmod, "module_config");
                if (jmod_cfgs == NULL) {
                    ret = ERR_UBSP_CRT_JSON_SCHEMA;
                    goto cleanup;
                }

                if (jcreator_modulecfg_schema(mod_info[iter].module_cfg,
                                              mod_info[iter].dev_count,
                                              &jmod_cfgs)) {
                    ret = ERR_UBSP_CRT_JSON_SCHEMA;
                    goto cleanup;
                }
            }
            cJSON_AddItemToArray(jmods, jmod);
        }
    }

    *schema = cJSON_Print(jschema);
    if (*schema == NULL) {
        fprintf(stderr, "Failed to print monitor.\n");
    } else {
        log_debug("JCREATOR: Created schema is : \n %s \n.", *schema);
    }

    if (ret) {
        log_error("Err(%d):: JCREATOR: Failed to create Schema.", ret);
    }

cleanup:
    cJSON_Delete(jschema);
    return ret;
}

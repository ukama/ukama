/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "mfg_parser.h"

#include "device.h"
#include "errorcode.h"

#include "usys_file.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_types.h"

StoreSchema *mfgStoreSchema[MAX_JSON_SCHEMA] = { '\0' };

static int read_mfg_data(char *fname, char *buff, off_t size) {
    int read_bytes = 0;
    /* Create input file descriptor */
    int fd = usys_open(fname, O_RDONLY, 0644);
    if (fd == -1) {
        usys_log_error("Failed to open file. Error: %s", usys_error(errno));
        return fd;
    }

    off_t off = usys_lseek(fd, 0, SEEK_SET);
    read_bytes = usys_read(fd, buff, size);
    return read_bytes;
}

static int read_mfg_data_size(char *fname) {
    int read_bytes = 0;
    /* Create input file descriptor */
    int fd = usys_open(fname, O_RDONLY, 0644);
    if (fd == -1) {
        usys_log_error("Failed to open file. Error: %s", usys_error(errno));
        return fd;
    }

    off_t off = usys_lseek(fd, 0L, SEEK_END);
    return off;
}

void parser_free_unit_cfg(UnitCfg **cfg, uint8_t count) {
    if (*cfg) {

        /* Device Cfgs */
        for (int iter = 0; iter < count; iter++) {
            usys_free((*cfg)->eepromCfg);
        }

        usys_free(*cfg);
    }
}

void parser_free_mfg_data(StoreSchema **sschema) {
    if (*sschema) {
        usys_free((*sschema)->indexTable);

        /* Unit Cfgs */
        for (int iter = 0; iter < (*sschema)->unitInfo.modCount; iter++) {
            usys_free((*sschema)->unitCfg[iter].eepromCfg);
        }

        usys_free((*sschema)->unitCfg);

        /* Module Cfgs */
        for (int iter = 0; iter < (*sschema)->modInfo.devCount; iter++) {
            usys_free((*sschema)->modCfg[iter].cfg);
        }

        usys_free((*sschema)->modCfg);
        usys_free((*sschema)->factCfg);
        usys_free((*sschema)->userCfg);
        usys_free((*sschema)->factCalib);
        usys_free((*sschema)->userCalib);
        usys_free((*sschema)->bsCerts);
        usys_free((*sschema)->cloudCerts);
        usys_free(*sschema);
    }
}

/* Parser to read integer value from JSON object */
bool parser_read_integer_value(const JsonObj *obj, int *ivalue) {
    bool ret = USYS_FALSE;

    /* Check if object is number */
    if (json_is_number(obj)) {
        *ivalue = json_integer_value(obj) ;
        ret = USYS_TRUE;
    }

    return ret;
}

/* Parser to read integer value from JSON object */
bool parser_read_integer_object(const JsonObj *obj, const char* key,
                int *ivalue) {
    bool ret = USYS_FALSE;

    /* Integer Json Object */
    const JsonObj *jIntObj = json_object_get(obj, key);

    /* Check if object is number */
    if (jIntObj && json_is_number(obj)) {
        *ivalue = json_integer_value(obj) ;
        ret = USYS_TRUE;
    }

    return ret;
}

/* Parser to read string value from JSON object */
bool parser_read_string_value(const JsonObj *obj, char *svalue) {
    bool ret = USYS_FALSE;
    int len = 0;

    /* Check if object is string */
    if (json_is_string(obj)) {
        len = json_string_length(obj);

        svalue = usys_malloc(sizeof(char) * len);
        if (svalue) {
            usys_memset(svalue, '\0', sizeof(char) * len);
            char *str = json_string_value(obj);
            usys_strcpy(svalue, str);
            json_decref(obj);
            ret = USYS_TRUE;
        }

    }

    return ret;
}

/* Parser to read string value from JSON object */
bool parser_read_string_object(const JsonObj *obj, const char* key,
                char **svalue) {
    bool ret = USYS_FALSE;

    /* String Json Object */
    const JsonObj *jStrObj = json_object_get(obj, key);

    /* Check if object is number */
    if (jStrObj && json_is_string(obj)) {
        int length = json_string_length(obj);

        *svalue = usys_malloc(sizeof(char) * length);
        if (*svalue) {
            usys_memset(*svalue, '\0', sizeof(char) * length);
            char *str = json_string_value(obj);
            usys_strcpy(*svalue, str);
            json_decref(obj);
            ret = USYS_TRUE;
        }
    }

    return ret;
}

/* Wrapper on top of parse_read_string */
bool parser_read_string_object_wrapper(const JsonObj *obj, const char* key,
                char* str) {
    bool ret = USYS_FALSE;
    char *tstr;
    if (parser_read_string_object(obj, key, &tstr)) {
        usys_strcpy(str, tstr);
        usys_free(tstr);
    }

    return ret;
}

/* Parser to read boolean value from JSON object */
bool parser_read_boolean_object(const JsonObj *obj, const char* key,
                bool *bvalue) {
    bool ret = USYS_FALSE;

    /* Integer Json Object */
    const JsonObj *jBoolObj = json_object_get(obj, key);

    /* Check if object is number */
    if (jBoolObj && json_is_boolean(obj)) {
        *bvalue = json_boolean_value(obj) ;
        ret = USYS_TRUE;
    }

    return ret;
}

/* Parse version */
Version *parse_version(const JsonObj *jVersion) {
    const JsonObj *jMajor = NULL;
    const JsonObj *jMinor = NULL;

    Version *pversion = usys_zmalloc(sizeof(Version));
    if (pversion) {

        /* Major version */
        if (!parser_read_integer_object(jVersion, JTAG_MAJOR_VERSION,
                        &pversion->major)) {
            usys_free(pversion);
            return NULL;
        }

        /* Minor version */
        if (!parser_read_integer_object(jVersion, JTAG_MINOR_VERSION,
                        &pversion->minor)) {
            usys_free(pversion);
            return NULL;
        }

    }
    return pversion;
}

/* Parse schema header */
SchemaHeader *parse_schema_header(const JsonObj *jHeader) {
    const JsonObj *jVersion = NULL;

    SchemaHeader *pHeader = usys_zmalloc(sizeof(SchemaHeader));
    if (pHeader) {

        jVersion = json_object_get(jHeader, JTAG_VERSION);

        /* Version infor for schema */
        Version *pversion = parse_version(jVersion);
        if (pversion) {
            usys_memcpy(&pHeader->version, pversion, sizeof(Version));
            usys_free(pversion);
        } else {
            goto cleanup;
        }

        /* Header info for schema */

        /* Index table offset */
        if (!parser_read_integer_object(jHeader, JTAG_IDX_TABLE_OFFSET,
                        &pHeader->idxTblOffset)) {
            goto cleanup;
        }

        /* Index tuple size  */
        if (!parser_read_integer_object(jHeader, JTAG_IDX_TUPLE_SIZE,
                        &pHeader->idxTplSize)) {
            goto cleanup;
        }

        /* Index Tuple Max Count */
        if (!parser_read_integer_object(jHeader, JTAG_IDX_TUPLE_MAX_COUNT,
                        &pHeader->idxTplMaxCount)) {
            goto cleanup;
        }

        /* Index Current Tuple */
        if (!parser_read_integer_object(jHeader, JTAG_IDX_CURR_TUPLE,
                        &pHeader->idxCurTpl)) {
            goto cleanup;
        }

        /* Module Capability */
        char* modCap;
        if (!parser_read_string_object(jHeader, JTAG_MODULE_CAPABILITY,
                        &modCap)) {
            goto cleanup;
        } else {
            if (!usys_strcmp(modCap, "AUTONOMOUS")) {
                pHeader->modCap = MOD_CAP_AUTONOMOUS;
            } else {
                pHeader->modCap = MOD_CAP_DEPENDENT;
            }
            usys_free(modCap);
        }

        /* Module Mode */
        char* modMode;
        if (!parser_read_string_object(jHeader, JTAG_MODULE_MODE,
                        &modMode)) {
            goto cleanup;
        } else {
            if (!usys_strcmp(modMode, JTAG_MODULE_MODE_MASTER)) {
                pHeader->modMode = MOD_MODE_MASTER;
            } else {
                pHeader->modMode = MOD_MODE_SLAVE;
            }
            usys_free(modMode);
        }


        /* Module Device Owner */
        char* modDevOwn;
        if (!parser_read_string_object(jHeader, JTAG_MODULE_DEV_OWNER,
                        &modDevOwn)) {
            goto cleanup;
        } else {
            if (!usys_strcmp(modMode, JTAG_DEV_OWNER)) {
                pHeader->modDevOwn = MOD_DEV_OWNER;
            } else {
                pHeader->modDevOwn = MOD_DEV_LENDER;
            }
            usys_free(modDevOwn);
        }

    }
    return pHeader;

    cleanup:
    usys_free(pHeader);
    return pHeader;
}

/* Parse Index Table */
SchemaIdxTuple *parse_schema_idx_table(const JsonObj *jIdxTab, uint8_t count) {
    const JsonObj *jIndexTpl = NULL;
    int iter = 0;
    if (!count) {
        goto cleanup;
    }

    SchemaIdxTuple *pIndexTable = usys_zmalloc(sizeof(SchemaIdxTuple) * count);
    if (pIndexTable) {
        json_array_foreach(jIdxTab, iter, jIndexTpl) {

            /* Field Id */
            if (!parser_read_integer_object(jIndexTpl, JTAG_FIELD_ID,
                            &pIndexTable[iter].fieldId)) {
                goto cleanup;
            }

            /* Payload offset */
            if (!parser_read_integer_object(jIndexTpl, JTAG_PAYLOAD_OFFSET,
                            &pIndexTable[iter].payloadOffset)) {
                goto cleanup;
            }

            /* Payload size */
            if (!parser_read_integer_object(jIndexTpl, JTAG_PAYLOAD_SIZE,
                            &pIndexTable[iter].payloadSize)) {
                goto cleanup;
            }

            /* Payload version */
            if (!parser_read_integer_object(jIndexTpl, JTAG_PAYLOAD_VERSION,
                            &pIndexTable[iter].payloadVer)) {
                goto cleanup;
            }

            /* Payload CRC */
            if (!parser_read_integer_object(jIndexTpl, JTAG_PAYLOAD_CRC,
                            &pIndexTable[iter].payloadCrc)) {
                goto cleanup;
            }

            /* State */
            char* payloadState;
            if (!parser_read_string_object(jIndexTpl, JTAG_STATE,
                            &payloadState)) {
                goto cleanup;
            } else {
                if (!usys_strcmp(payloadState, JTAG_STATE_ENABLED)) {
                    pIndexTable[iter].state = IDX_ENTRY_ENABLED;
                } else {
                    pIndexTable[iter].state = IDX_ENTRY_DISABLED;
                }
                usys_free(payloadState);
            }

            /* Valid */
            if (!parser_read_boolean_object(jIndexTpl, JTAG_VALID,
                            &pIndexTable[iter].valid)) {
                goto cleanup;
            }

        }

        /* Verify if all tuples are parsed */
        if (iter == count) {
            usys_log_debug(
                            "All %d tuples read from the manufacturing data.",
                            count);
        } else {
            usys_log_error(
                            "Error: Only %d tuples read from the manufacturing data expected were %d.",
                            (count - iter), count);
        }

    } else {
        usys_log_error("Memory exhausted. Error: %s",
                        usys_error(errno));
    }

    return pIndexTable;

    cleanup:
    usys_free(pIndexTable);
    pIndexTable = NULL;
    return pIndexTable;
}

/* Parse GPIO Device HW Attributes */
DevGpioCfg *parse_schema_dev_gpio(const JsonObj *jDevSchema) {
    const JsonObj *jGPIOCfg = NULL;

    DevGpioCfg *pDevCfg = usys_zmalloc(sizeof(DevGpioCfg));
    if (pDevCfg) {

        /* Direction */
        char *direct;
        if (!parser_read_string_object(jDevSchema, JTAG_GPIO_DIRECTION,
                        &direct)) {
            goto cleanup;
        } else {
            if (!usys_strcmp(direct, JTAG_GPIO_DIRECTION)) {
                pDevCfg->direction= DEV_GPIO_INPUT;
            } else {
                pDevCfg->direction = DEV_GPIO_OUTPUT;
            }
            usys_free(direct);
        }

        /* Number */
        if (!parser_read_integer_object(jDevSchema, JTAG_GPIO_NUMBER,
                        &pDevCfg->gpioNum)) {
            goto cleanup;
        }

    } else {
        usys_log_error("Memory exhausted while parsing DevGpioCfg. Error: %s",
                        usys_error(errno));
        goto cleanup;
    }

    return pDevCfg;

    cleanup:
    usys_free(pDevCfg);
    pDevCfg = NULL;
    return pDevCfg;
}

/* Parse I2C Device HW Attributes */
DevI2cCfg *parse_schema_dev_i2c(const JsonObj *jDevSchema) {
    const JsonObj *jDevCfg = NULL;
    const JsonObj *jBus = NULL;
    const JsonObj *jAddress = NULL;

    DevI2cCfg *pDevCfg = NULL;
    pDevCfg = usys_zmalloc(sizeof(DevI2cCfg));
    if (pDevCfg) {

        /* Bus */
        if (!parser_read_integer_object(jDevSchema, JTAG_BUS,
                        &pDevCfg->bus)) {
            goto cleanup;
        }

        /* Address */
        if (!parser_read_integer_object(jDevSchema, JTAG_ADDRESS,
                        &pDevCfg->add)) {
            goto cleanup;
        }

    } else {
        usys_log_error("Memory exhausted while parsing DevI2cCfg.Error: %s",
                        usys_error(errno));
        goto cleanup;
    }
    return pDevCfg;

    cleanup:
    usys_free(pDevCfg);
    pDevCfg = NULL;
    return pDevCfg;
}

/* Parse SPI schema */
DevSpiCfg *parse_schema_dev_spi(const JsonObj *jDevSchema) {
    const JsonObj *jDevCfg = NULL;
    const JsonObj *jCS = NULL;

    DevSpiCfg *pDevCfg = usys_zmalloc(sizeof(DevSpiCfg));
    if (pDevCfg) {

        /* Chip select */
        jCS = json_object_get(jDevSchema, JTAG_CHIP_SELECT);
        if (json_is_object(jCS)) {

            DevGpioCfg *pCS = parse_schema_dev_gpio(jCS);
            if (pCS) {
                usys_memcpy(&pDevCfg->cs, pCS, sizeof(DevGpioCfg));
            } else {
                usys_log_error("Failed to parse DevSpiCfg.cs");
                goto cleanup;
            }

        } else {
            usys_log_error("Failed to parse DevSpiCfg.cs.Error: Unexpected json tag");
            goto cleanup;
        }

        /* Bus */
        if (!parser_read_integer_object(jDevSchema, JTAG_BUS,
                        &pDevCfg->bus)) {
            goto cleanup;
        }

    } else {
        usys_log_error("Memory exhausted while parsing DevSpiCfg.Error: %s",
                        usys_error(errno));
        goto cleanup;
    }
    return pDevCfg;

    cleanup:
    usys_free(pDevCfg);
    pDevCfg = NULL;
    return pDevCfg;
}

/* Parse UART schema */
DevUartCfg *parse_schema_dev_uart(const JsonObj *jDevSchema) {
    const JsonObj *jDevCfg = NULL;

    DevUartCfg *pDevCfg = usys_zmalloc(sizeof(DevUartCfg));
    if (pDevCfg) {

        /* UART Number */
        if (!parser_read_integer_object(jDevSchema, JTAG_UART,
                        &pDevCfg->uartNo)) {
            goto cleanup;
        }

    } else {
        usys_log_error("Memory exhausted while parsing DevUartCfg.Error %s",
                        usys_error(errno));
        goto cleanup;
    }
    return pDevCfg;

    cleanup:
    usys_free(pDevCfg);
    return pDevCfg;
}


/* Parse device config */
void *parse_schema_devices(const JsonObj *jSchema, uint16_t class) {
    void *dev = NULL;
    switch (class) {
        case DEV_CLASS_GPIO:
            dev = parse_schema_dev_gpio(jSchema);
            break;
        case DEV_CLASS_I2C:
            dev = parse_schema_dev_i2c(jSchema);
            break;
        case DEV_CLASS_SPI:
            dev = parse_schema_dev_spi(jSchema);
            break;
        case DEV_CLASS_UART:
            dev = parse_schema_dev_uart(jSchema);
            break;
        default:
            usys_log_error("Unkown device type failed to parse. Error %s",
                            ERR_NODED_INVALID_DEVICE_CFG);
    }
    return dev;
}

/* Parse Unit Info */
void *parse_schema_unit_info(const JsonObj *jSchema) {
    const JsonObj *jUnitInfo = NULL;
    const JsonObj *jSWVer = NULL;
    const JsonObj *jProdSWVer = NULL;
    UnitInfo *pUnitInfo = NULL;
    int ret = 0;

    /* Unit Info */
    jUnitInfo = json_object_get(jSchema, JTAG_UNIT_INFO);
    if (json_is_object(jUnitInfo)) {

        pUnitInfo = usys_zmalloc(sizeof(UnitInfo));
        if (pUnitInfo) {

            /* UUID */
            if (!parse_read_string_object_wrapper(jUnitInfo, JTAG_UUID,
                            &pUnitInfo->uuid)) {
                goto cleanup;
            }

            /* Name */
            if (!parser_read_string_object_wrapper(jUnitInfo, JTAG_NAME,
                            &pUnitInfo->name)) {
                goto cleanup;
            }

             /* Type */
             if (!parser_read_integer_object(jUnitInfo, JTAG_TYPE,
                             &pUnitInfo->unit)) {
                 goto cleanup;
             }

             /* Part Number */
             if (!parser_read_string_object_wrapper(jUnitInfo, JTAG_PART_NUMBER,
                             &pUnitInfo->partNo)) {
                 goto cleanup;
             }

             /* Skew */
             if (!parser_read_string_object_wrapper(jUnitInfo, JTAG_SKEW,
                             &pUnitInfo->skew)) {
                 goto cleanup;
             }

             /* MAC */
             if (!parser_read_string_object_wrapper(jUnitInfo, JTAG_MAC,
                             &pUnitInfo->mac)) {
                 goto cleanup;
             }

             /* SW Version */
             jSWVer = json_object_get(jUnitInfo, JTAG_SW_VERISION);
             Version *sVersion = parse_version(jSWVer);
             if (sVersion) {
                 usys_memcpy(&pUnitInfo->swVer, sVersion, sizeof(Version));
                 usys_free(sVersion);
             } else {
                 usys_log_error("Failed to parse UnitInfo.swVersion.");
                 goto cleanup;
             }

             /* Production SW Version */
             jProdSWVer = json_object_get(jUnitInfo, JTAG_PROD_SW_VERSION);
             Version *pVersion = parse_version(jProdSWVer);
             if (pVersion) {
                 usys_memcpy(&pUnitInfo->swVer, pVersion, sizeof(Version));
                 usys_free(pVersion);
             } else {
                 usys_log_error("Failed to parse UnitInfo.prodSwVersion.");
                 goto cleanup;
             }

             /* Assembly Date */
             if (!parser_read_string_object_wrapper(jUnitInfo, JTAG_ASM_DATE,
                             &pUnitInfo->assmDate)) {
                 goto cleanup;
             }

             /* OEM Name */
             if (!parser_read_string_object_wrapper(jUnitInfo, JTAG_OEM_NAME,
                             &pUnitInfo->oemName)) {
                 goto cleanup;
             }

             /* Module Count */
             if (!parser_read_integer_object(jUnitInfo, JTAG_MODULE_COUNT,
                             &pUnitInfo->modCount)) {
                 goto cleanup;
             }

        } else {
            usys_log_error(
                            "Memory exhausted while parsing Unit Info. Error: %s",
                            usys_error(errno));
            goto cleanup;
        }
    } else {
        usys_log_error(
                        "Unexpected JSON object found instead of Unit Info.");

        goto cleanup;
    }
    return pUnitInfo;

    cleanup:
    usys_free(pUnitInfo);
    pUnitInfo = NULL;
    return pUnitInfo;
}

/* Parse Unit Config */
void *parse_schema_unit_config(const JsonObj *jSchema, uint16_t count) {
    const JsonObj *jUnitCfgs = NULL;
    const JsonObj *jUnitCfg = NULL;
    const JsonObj *jDevice = NULL;
    UnitCfg *pUnitCfg = NULL;
    int iter = 0;

    if (!count) {
        goto cleanup;
    }

    jUnitCfgs = json_object_get(jSchema, JTAG_UNIT_CONFIG);
    if (json_is_object(jUnitCfgs)) {

        pUnitCfg = usys_zmalloc(sizeof(UnitCfg) * count);
        if (pUnitCfg) {

            json_array_foreach(jUnitCfgs, iter, jUnitCfg) {
                /* UUID */
                if (!parse_read_string_object_wrapper(jUnitCfg, JTAG_UUID,
                                &pUnitCfg[iter].modUuid)) {
                    goto cleanup;
                }

                /* Name */
                if (!parser_read_string_object_wrapper(jUnitCfg, JTAG_NAME,
                                &pUnitCfg[iter].modName)) {
                    goto cleanup;
                }

                /* SysFs */
                if (!parser_read_string_object_wrapper(jUnitCfg,
                                JTAG_INVT_SYSFS_FILE, &pUnitCfg[iter].sysFs)) {
                    goto cleanup;
                }

                /* EEPROM */
                jDevice = json_object_get(jUnitCfg, JTAG_INVT_DEV_INFO);
                DevI2cCfg *pDev = parse_schema_devices(jDevice, DEV_CLASS_I2C);
                if (pDev) {
                    pUnitCfg[iter].eepromCfg = pDev;
                } else {
                    usys_log_error(
                                    "Failed to parse UnitCfg[%d].eepromCfg",
                                    iter);
                    goto cleanup;
                }
                iter++;
            }

            /* Verify count of device read */
            if (iter == count) {
                usys_log_debug(
                                "All %d modules info read from the"
                                "manufacturing data's unit config.",
                                count);
            } else {
                usys_log_error(
                                "Only %d modules info read from the "
                                "manufacturing data's unit config expected "
                                "were %d.",
                                iter, count);
            }
        }
    } else {
        usys_log_error(
                        "Unexpected JSON object found instead"
                        "of Unit Config.");
        goto cleanup;
    }
    return pUnitCfg;

    cleanup:
    parser_free_unit_cfg(&pUnitCfg, count);
    return pUnitCfg;
}

/* Parse Module Info */
void *parse_schema_module_info(const JsonObj *jSchema) {
    const JsonObj *jModuleInfo = NULL;
    const JsonObj *jSWVer = NULL;
    const JsonObj *jProdSWVer = NULL;

    ModuleInfo *pModuleInfo = NULL;
    int ret = 0;
    /* Module Info */
    jModuleInfo = json_object_get(jSchema, "jModuleInfo");
    if (json_is_object(jModuleInfo)) {

        pModuleInfo = usys_zmalloc(sizeof(ModuleInfo));
        if (pModuleInfo) {

            /* UUID */
            if (!parse_read_string_object_wrapper(jModuleInfo, JTAG_UUID,
                            &pModuleInfo->uuid)) {
                goto cleanup;
            }

            /* Name */
            if (!parser_read_string_object_wrapper(jModuleInfo, JTAG_NAME,
                            &pModuleInfo->name)) {
                goto cleanup;
            }

            /* Type */
            if (!parser_read_integer_object(jModuleInfo, JTAG_TYPE,
                            &pModuleInfo->module)) {
                goto cleanup;
            }

            /* Part Number */
            if (!parser_read_string_object_wrapper(jModuleInfo, JTAG_PART_NUMBER,
                            &pModuleInfo->partNo)) {
                goto cleanup;
            }

            /* HW Version */
            if (!parser_read_string_object_wrapper(jModuleInfo, JTAG_HW_VERSION,
                            &pModuleInfo->hwVer)) {
                goto cleanup;
            }

            /* MAC */
            if (!parser_read_string_object_wrapper(jModuleInfo, JTAG_MAC,
                            &pModuleInfo->mac)) {
                goto cleanup;
            }

            /* SW Version */
            jSWVer = json_object_get(jModuleInfo, JTAG_SW_VERISION);
            Version *sVersion = parse_version(jSWVer);
            if (sVersion) {
                usys_memcpy(&pModuleInfo->swVer, sVersion, sizeof(Version));
                usys_free(sVersion);
            } else {
                usys_log_error("Failed to parse UnitInfo.swVersion.");
                goto cleanup;
            }

            /* Production SW Version */
            jProdSWVer = json_object_get(jModuleInfo, JTAG_PROD_SW_VERSION);
            Version *pVersion = parse_version(jProdSWVer);
            if (pVersion) {
                usys_memcpy(&pModuleInfo->swVer, pVersion, sizeof(Version));
                usys_free(pVersion);
            } else {
                usys_log_error("Failed to parse UnitInfo.prodSwVersion.");
                goto cleanup;
            }

            /* Manufacturing Date  */
            if (!parser_read_string_object_wrapper(jModuleInfo, JTAG_MFG_DATE,
                            &pModuleInfo->mfgDate)) {
                goto cleanup;
            }

            /* Manufacturer Name */
            if (!parser_read_string_object_wrapper(jModuleInfo, JTAG_MFG_NAME,
                            &pModuleInfo->mfgName)) {
                goto cleanup;
            }

            /* Device Count */
            if (!parser_read_integer_object(jModuleInfo, JTAG_DEVICE_COUNT,
                            &pModuleInfo->devCount)) {
                goto cleanup;
            }

        } else {
            usys_log_error(
                            "Memory exhausted while parsing Module Info. "
                            "Error %s",
                            usys_error(errno));
            goto cleanup;
        }

    } else {
        usys_log_error(
                        "Unexpected JSON object found instead"
                        " of Module Info.");
        goto cleanup;
    }
    return pModuleInfo;

    cleanup:
    usys_free(pModuleInfo);
    pModuleInfo = NULL;
    return pModuleInfo;
}

/* Parse module config */
void *parse_schema_module_config(const JsonObj *jSchema, uint16_t count) {
    const JsonObj *jModCfgs = NULL;
    const JsonObj *jModCfg = NULL;
    const JsonObj *jDevice = NULL;
    ModuleCfg *pModCfg = NULL;
    int iter = 0;

    if (!count) {
        goto cleanup;
    }

    jModCfgs = json_object_get(jSchema, JTAG_MODULE_CONFIG);
    if (!json_is_array(jModCfgs)) {
        usys_log_error(
                        "Unexpected JSON object found instead of Module Config.");
        goto cleanup;
    }

    pModCfg = usys_zmalloc(sizeof(ModuleCfg) * count);
    if (pModCfg) {

        json_array_foreach(jModCfgs, iter, jModCfg) {

            /* Name */
            if (!parser_read_string_object_wrapper(jModCfg, JTAG_NAME,
                            &pModCfg[iter].devName)) {
                goto cleanup;
            }

            /* Description */
            if (!parse_read_string_object_wrapper(jModCfg, JTAG_DESCRIPTION,
                            &pModCfg[iter].devDesc)) {
                goto cleanup;
            }

            /* SysFs */
            if (!parser_read_string_object_wrapper(jModCfg,
                            JTAG_DEV_SYSFS_FILE, &pModCfg[iter].sysFile)) {
                goto cleanup;
            }

            /* Device Class */
            if (!parser_read_integer_object(jModCfg, JTAG_CLASS,
                            &pModCfg[iter].devClass)) {
                goto cleanup;
            }

            /* Device Type */
            if (!parser_read_integer_object(jModCfg, JTAG_TYPE,
                            &pModCfg[iter].devType)) {
                goto cleanup;
            }

            /* Device HW attributes */
            jDevice = json_object_get(jModCfg, JTAG_DEV_HW_ATTRS);
            DevI2cCfg *pDev = parse_schema_devices(jDevice,
                            pModCfg[iter].devClass);
            if (pDev) {
                pModCfg[iter].cfg = pDev;
            } else {
                usys_log_error(
                                "Failed to parse UnitCfg[%d].cfg",
                                iter);
                goto cleanup;
            }

            iter++;
        }

        /* Verify count of devices parsed */
        if (iter == count) {
            usys_log_debug(
                            "PARSER:: All %d device info read from the manufacturing data's Module config.",
                            count);
        } else {
            usys_log_error(
                            "PARSER::Only %d device info read from the manufacturing data's Module config expected were %d.",
                            iter, count);
        }

    } else {
        usys_log_error(
                        "Memory exhausted while parsing Module Config. "
                        "Error %s",
                        usys_error(errno));
        goto cleanup;
    }
    return pModCfg;

    cleanup:
    usys_free(pModCfg);
    pModCfg = NULL;
    return pModCfg;
}

/* parse generic file data */
void *parse_schema_generic_file_data(const JsonObj *jSchema, uint16_t *size,
                char *name_key) {
    const JsonObj *jGenData = NULL;
    char *genFileName = NULL;
    char *genData = NULL;
    int ret = 0;

    genFileName = usys_zmalloc(sizeof(char) * PATH_LENGTH);
    if (genFileName) {

        /* Read file name */
        if (!parser_read_string_object_wrapper(jSchema,
                        name_key, &genFileName)) {
            goto cleanup;
        }

        /* Read data size*/
        int rsize = read_mfg_data_size(genFileName);
        if (rsize <= 0) {
            usys_log_error("Failed to read data.",
                            genFileName);
            goto cleanup;
        }

        /* Read data */
        genData = usys_zmalloc(sizeof(char) * rsize);
        if (genData) {

            ret = read_mfg_data(genFileName, genData, rsize);
            if (ret == rsize) {
                usys_log_debug("File %s read data of %d bytes.",
                                genFileName, rsize);
            } else {
                usys_log_error(
                                "File %s read data of %d bytes "
                                "expected %d bytes.",
                                genFileName, ret, rsize);
            }

            *size = ret;

        } else {

            usys_log_error(
                            "Memory exhausted while reading data from file %s.",
                            "Error: %s", genFileName, usys_error(errno));
            goto cleanup;

        }

    } else {
        usys_log_error(
                        "Memory exhausted while reading data file name.",
                        "Error: %s", usys_error(errno));
        goto cleanup;
    }

    cleanup:
    usys_free(genFileName);
    return genData;
}

void *parse_schema_generic_file_data_wrapper(const JsonObj *jSchema, int iter,
                int* payloadSize, char *name_key, bool* status) {
    *status = USYS_FALSE;
    int size = 0;
    char *factCfg =
                    parse_schema_generic_file_data(jSchema, &size,
                                    JTAG_FACTORY_CONFIG);
    if (factCfg) {
        if (size != *payloadSize) {
            usys_log_warn(
                            "Size read for Field id 0x%x is %d bytes "
                            "and size mentioned in index table [%d] is 0x%d bytes.",
                            size, iter, *payloadSize, size);
            usys_log_debug(
                            "Updating index table [%d] size to %d bytes.",
                            iter, size);
            *payloadSize = size;
        }
        *status = USYS_TRUE;
    }
    return factCfg;
}

/* Parse payloads */
int parse_schema_payload(const JsonObj *jSchema, StoreSchema **schema,
                uint16_t id,
                int iter) {
    int ret = 0;
    bool status = USYS_FALSE;
    switch (id) {
        case FIELD_ID_UNIT_INFO: {
            UnitInfo *pUnitInfo = parse_schema_unit_info(jSchema);
            if (pUnitInfo) {
                usys_memcpy(&(*schema)->unitInfo, pUnitInfo, sizeof(UnitInfo));
                usys_free(pUnitInfo);
                pUnitInfo = NULL;
            } else {
                ret = -1;
                goto cleanup;
            }
            break;
        }
        case FIELD_ID_UNIT_CFG: {
            uint16_t modCount = (*schema)->unitInfo.modCount;
            UnitCfg *pUnitCfg = parse_schema_unit_config(schema, modCount);
            if (pUnitCfg) {
                (*schema)->unitCfg = pUnitCfg;
            } else {
                ret = -1;
                goto cleanup;
            }
            break;
        }
        case FIELD_ID_MODULE_INFO: {
            ModuleInfo *pModuleInfo = parse_schema_module_info(jSchema);
            if (pModuleInfo) {
                usys_memcpy(&(*schema)->modInfo, pModuleInfo,
                                sizeof(ModuleInfo));
                usys_free(pModuleInfo);
                pModuleInfo = NULL;
            } else {
                ret = -1;
                goto cleanup;
            }
            break;
        }
        case FIELD_ID_MODULE_CFG: {
            uint16_t devCount = (*schema)->modInfo.devCount;
            ModuleCfg *pModuleCfg =
                            parse_schema_module_config(jSchema, devCount);
            if (pModuleCfg) {
                (*schema)->modCfg = pModuleCfg;
            } else {
                ret = -1;
                goto cleanup;
            }
            break;
        }
        case FIELD_ID_FACT_CFG: {
            (*schema)->factCfg =
                            parse_schema_generic_file_data_wrapper(jSchema,
                                       iter,
                                       &(*schema)->indexTable[iter].payloadSize,
                                       JTAG_FACTORY_CONFIG,
                                       &status);

            if (!status) {
                ret = -1;
                goto cleanup;
            }

            break;
        }
        case FIELD_ID_USER_CFG: {
            (*schema)->userCfg =
                            parse_schema_generic_file_data_wrapper(jSchema,
                                       iter,
                                       &(*schema)->indexTable[iter].payloadSize,
                                       JTAG_USER_CONFIG,
                                       &status);

            if (!status) {
                ret = -1;
                goto cleanup;
            }

            break;
        }
        case FIELD_ID_FACT_CALIB: {
            (*schema)->factCalib =
                            parse_schema_generic_file_data_wrapper(jSchema,
                                       iter,
                                       &(*schema)->indexTable[iter].payloadSize,
                                       JTAG_FACTORY_CALIB,
                                       &status);

            if (!status) {
                ret = -1;
                goto cleanup;
            }

            break;
        }
        case FIELD_ID_USER_CALIB: {
            (*schema)->userCalib =
                            parse_schema_generic_file_data_wrapper(jSchema,
                                       iter,
                                       &(*schema)->indexTable[iter].payloadSize,
                                       JTAG_USER_CALIB,
                                       &status);

            if (!status) {
                ret = -1;
                goto cleanup;
            }

            break;
        }
        case FIELD_ID_BS_CERTS: {
            (*schema)->bsCerts =
                            parse_schema_generic_file_data_wrapper(jSchema,
                                       iter,
                                       &(*schema)->indexTable[iter].payloadSize,
                                       JTAG_BOOTSTRAP_CERTS,
                                       &status);

            if (!status) {
                ret = -1;
                goto cleanup;
            }

            break;
        }
        case FIELD_ID_CLOUD_CERTS: {
            (*schema)->cloudCerts =
                            parse_schema_generic_file_data_wrapper(jSchema,
                                       iter,
                                       &(*schema)->indexTable[iter].payloadSize,
                                       JTAG_CLOUD_CERTS,
                                       &status);

            if (!status) {
                ret = -1;
                goto cleanup;
            }

            break;
        }
        default: {
            ret = ERR_NODED_INVALID_FIELD;
            usys_log_error("Invalid Field id supplied by Index entry.Error %d", ret);
        }
    }

    cleanup:
    return ret;
}

void parser_error(JsonErrObj *jErr, char* msg) {
    if (jErr) {
        usys_log_error("%s. Error: %s ", msg, jErr->text);
    } else {
        usys_log_error("%s. No error info available", msg);
    }
}

int parse_mfg_schema(const char *mfgdata, uint8_t idx) {
    int ret = 0;
    JsonObj *jSchema = NULL;
    JsonErrObj *jErr;
    const JsonObj *jHeader = NULL;
    const JsonObj *jIdxTable = NULL;
    const JsonObj *jUnitInfo = NULL;
    const JsonObj *jUnitCfg = NULL;
    const JsonObj *jModCfg = NULL;
    const JsonObj *jFactConfig = NULL;
    const JsonObj *jUserCfg = NULL;
    const JsonObj *jFactCalib = NULL;
    const JsonObj *jUserCalib = NULL;
    const JsonObj *jBootstrapCerts = NULL;
    const JsonObj *jCloudCerts = NULL;
    StoreSchema *storeSchema = NULL;
    mfgStoreSchema[idx] = usys_zmalloc(sizeof(storeSchema));
    if (mfgStoreSchema[idx]) {
        storeSchema = mfgStoreSchema[idx];
        jSchema = json_loads(mfgdata, JSON_DECODE_ANY, jErr);
        if (!jSchema) {
            parser_error(jErr, "Failed to parse schema");
            ret = ERR_NODED_JSON_PARSER;
            goto cleanup;
        }

        /* Header */
        jHeader = json_object_get(jSchema, JTAG_HEADER);
        if (jHeader) {
            SchemaHeader *pHeader = parse_schema_header(jHeader);
            if (pHeader) {
                usys_memcpy(&storeSchema->header, pHeader,
                                sizeof(SchemaHeader));
                usys_free(pHeader);
            } else {
                ret = ERR_NODED_INVALID_JSON_OBJECT;
                goto cleanup;
            }
        } else {
            ret = -1;
        }

        /* Index Table */
        jIdxTable = json_object_get(jSchema, "index_table");
        if (jIdxTable) {

            SchemaIdxTuple *pIndexTable = parse_schema_index_table(
                            jIdxTable, storeSchema->header.idxCurTpl);
            if (pIndexTable) {
                /* Free me once done*/
                storeSchema->indexTable = pIndexTable;
            } else {
                ret = ERR_NODED_INVALID_JSON_OBJECT;
                goto cleanup;
            }

        } else {
            ret = -1;
        }

        for (int iter = 0; iter < storeSchema->header.idxCurTpl; iter++) {

            uint16_t id = storeSchema->indexTable[iter].fieldId;
            ret = parse_schema_payload(jSchema, &storeSchema, id, iter);
            if (ret) {
                usys_log_error(
                                "Failed parsing for Field Id 0x%x from mfg data"
                                ".Error: %d",
                                id, ret);
                goto cleanup;
            } else {
                usys_log_debug(
                                "Parsing for Field Id 0x%x from mfg data "
                                "completed.",
                                id);
            }

        }
    }

    cleanup:
    json_decref(jSchema);
    if (ret) {
        parser_free_mfg_data(&storeSchema);
    }
    return ret;
}

StoreSchema *parser_get_mfg_data_by_uuid(char *puuid) {
    int ret = 0;
    StoreSchema *sschema = NULL;

    /*Default section*/
    if ((!puuid) || !usys_strcmp(puuid, "")) {
        sschema = mfgStoreSchema[0];
        usys_log_trace(
                        "MFG Data set to the Module UUID %s with %d entries"
                        " in Index Table.",
                        sschema->modInfo.uuid, sschema->header.idxCurTpl);
    } else {

        /* Searching for Module MFG data index.*/
        for (uint8_t iter = 0; iter < MAX_JSON_SCHEMA; iter++) {

            if (mfgStoreSchema[iter]) {
                if (!usys_strcmp(puuid, mfgStoreSchema[iter]->modInfo.uuid)) {
                    sschema = mfgStoreSchema[iter];
                    usys_log_trace(
                                    "PARSER:: MFG Data set to the Module UUID"
                                    "%s with %d entries in Index Table.",
                                    sschema->modInfo.uuid,
                                    sschema->header.idxCurTpl);
                    break;
                }

            }
        }
    }
    return sschema;
}

int parser_schema_init(JSONInput *ip) {
    int ret = 0;
    char *fname = NULL;
    if ((ip->fname) && (ip->count > 0)) {

        for (uint8_t iter = 0; iter < ip->count; iter++) {
            if (!ip->fname[iter]) {
                ret = -1;
                return ret;
            }

            fname = ip->fname[iter];
            usys_log_debug("PARSER:: Starting the parsing of %s.", fname);
            off_t size = read_mfg_data_size(fname);

            char *schemabuff = usys_zmalloc((sizeof(char) * size) + 1);
            if (schemabuff) {

                /* Read mfg data from file */
                ret = read_mfg_data(fname, schemabuff, size);
                if (ret == size) {
                    usys_log_debug(
                                    "File %s read manufacturing data of %d"
                                    " bytes.",
                                    fname, size);

                    /* Parse the mfg data to store schema */
                    ret = parse_mfg_schema(schemabuff, iter);
                    if (ret) {
                        usys_log_error("Err(%d): PARSER:: Parsing failed for %s.",
                                        ret, fname);
                    } else {
                        usys_log_debug("PARSER: Parsing completed for %s.",
                                        fname);
                    }

                }

            }

            usys_free(schemabuff);
        }
    }
    return ret;
}

void parser_schema_exit() {
    for (int iter = 0; iter < MAX_JSON_SCHEMA; iter++) {

        if (mfgStoreSchema[iter]) {
            parser_free_mfg_data(&mfgStoreSchema[iter]);
        }

    }
}

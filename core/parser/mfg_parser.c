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

#include "usys_api.h"
#include "usys_file.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
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

void parser_free_unit_cfg(NodeCfg **cfg, uint8_t count) {
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

/* Parse schema header */
SchemaHeader *parse_schema_header(const JsonObj *jHeader) {
    const JsonObj *jVersion = NULL;

    SchemaHeader *pHeader = usys_zmalloc(sizeof(SchemaHeader));
    if (pHeader) {
        jVersion = json_object_get(jHeader, JTAG_VERSION);

        /* Version info for schema */
        Version *pversion = parse_version(jVersion);
        if (pversion) {
            usys_memcpy(&pHeader->version, pversion, sizeof(Version));
            usys_free(pversion);
        } else {
            goto cleanup;
        }

        /* Header info for schema */

        /* Index table offset */
        uint16_t idxTblOffset = 0;
        if (!parser_read_uint16_object(jHeader, JTAG_IDX_TABLE_OFFSET,
                                       &idxTblOffset)) {
            goto cleanup;
        } else {
            pHeader->idxTblOffset = idxTblOffset;
        }

        /* Index tuple size  */
        uint16_t idxTplSize = 0;
        if (!parser_read_uint16_object(jHeader, JTAG_IDX_TUPLE_SIZE,
                                       &idxTplSize)) {
            goto cleanup;
        } else {
            pHeader->idxTplSize = idxTplSize;
        }

        /* Index Tuple Max Count */
        uint16_t idxTplMaxCount = 0;
        if (!parser_read_uint16_object(jHeader, JTAG_IDX_TUPLE_MAX_COUNT,
                                       &idxTplMaxCount)) {
            goto cleanup;
        } else {
            pHeader->idxTplMaxCount = idxTplMaxCount;
        }

        /* Index Current Tuple */
        uint16_t idxCurTpl = 0;
        if (!parser_read_uint16_object(jHeader, JTAG_IDX_CURR_TUPLE,
                                       &idxCurTpl)) {
            goto cleanup;
        } else {
            pHeader->idxCurTpl = idxCurTpl;
        }

        /* Module Capability */
        char *modCap;
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
        char *modMode;
        if (!parser_read_string_object(jHeader, JTAG_MODULE_MODE, &modMode)) {
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
        char *modDevOwn;
        if (!parser_read_string_object(jHeader, JTAG_MODULE_DEV_OWNER,
                                       &modDevOwn)) {
            goto cleanup;
        } else {
            if (!usys_strcmp(modDevOwn, JTAG_DEV_OWNER)) {
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
    pHeader = NULL;
    return pHeader;
}

/* Parse Index Table */
SchemaIdxTuple *parse_schema_idx_table(const JsonObj *jIdxTab, uint8_t count) {
    const JsonObj *jIndexTpl = NULL;
    const JsonObj *jPayloadVersion = NULL;

    int iter = 0;
    if (!count) {
        goto cleanup;
    }

    SchemaIdxTuple *pIndexTable = usys_zmalloc(sizeof(SchemaIdxTuple) * count);
    if (pIndexTable) {
        json_array_foreach(jIdxTab, iter, jIndexTpl) {
            /* Field Id */
            uint16_t fieldId = 0;
            if (!parser_read_uint16_object(jIndexTpl, JTAG_FIELD_ID,
                                           &fieldId)) {
                goto cleanup;
            } else {
                pIndexTable[iter].fieldId = fieldId;
            }

            /* Payload offset */
            uint16_t payloadOffset = 0;
            if (!parser_read_uint16_object(jIndexTpl, JTAG_PAYLOAD_OFFSET,
                                           &payloadOffset)) {
                goto cleanup;
            } else {
                pIndexTable[iter].payloadOffset = payloadOffset;
            }

            /* Payload size */
            uint16_t payloadSize = 0;
            if (!parser_read_uint16_object(jIndexTpl, JTAG_PAYLOAD_SIZE,
                                           &payloadSize)) {
                goto cleanup;
            } else {
                pIndexTable[iter].payloadSize = payloadSize;
            }

            /* Payload version */
            jPayloadVersion = json_object_get(jIndexTpl, JTAG_PAYLOAD_VERSION);
            Version *sVersion = parse_version(jPayloadVersion);
            if (sVersion) {
                usys_memcpy(&(pIndexTable[iter].payloadVer), sVersion,
                            sizeof(Version));
                usys_free(sVersion);
            } else {
                usys_log_error("Failed to parse Payload[%d].Version.", iter);
                goto cleanup;
            }

            /* Payload CRC */
            uint32_t payloadCrc = 0;
            if (!parser_read_uint32_object(jIndexTpl, JTAG_PAYLOAD_CRC,
                                           &payloadCrc)) {
                goto cleanup;
            } else {
                pIndexTable[iter].payloadCrc = payloadCrc;
            }

            /* State */
            char *payloadState;
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
            usys_log_debug("All %d Index tuples read from the manufacturing "
                           "data.",
                           count);
        } else {
            usys_log_error("Error: Only %d Index tuples read from the "
                           "manufacturing data expected were %d.",
                           (count - iter), count);
        }

    } else {
        usys_log_error("Memory exhausted. Error: %s", usys_error(errno));
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
                pDevCfg->direction = DEV_GPIO_INPUT;
            } else {
                pDevCfg->direction = DEV_GPIO_OUTPUT;
            }
            usys_free(direct);
        }

        /* Number */
        int gpioNum = 0;
        if (!parser_read_integer_object(jDevSchema, JTAG_GPIO_NUMBER,
                                        &gpioNum)) {
            goto cleanup;
        } else {
            pDevCfg->gpioNum = gpioNum;
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
        uint8_t bus = 0;
        if (!parser_read_uint8_object(jDevSchema, JTAG_BUS, &bus)) {
            goto cleanup;
        } else {
            pDevCfg->bus = bus;
        }

        /* Address */
        uint16_t add = 0;
        if (!parser_read_uint16_object(jDevSchema, JTAG_ADDRESS, &add)) {
            goto cleanup;
        } else {
            pDevCfg->add = add;
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
            usys_log_error(
                "Failed to parse DevSpiCfg.cs.Error: Unexpected json tag");
            goto cleanup;
        }

        /* Bus */
        uint8_t bus = 0;
        if (!parser_read_uint8_object(jDevSchema, JTAG_BUS, &bus)) {
            goto cleanup;
        } else {
            pDevCfg->bus = bus;
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
        uint16_t uartNo = 0;
        if (!parser_read_uint16_object(jDevSchema, JTAG_UART, &uartNo)) {
            goto cleanup;
        } else {
            pDevCfg->uartNo = uartNo;
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
void *parse_schema_node_info(const JsonObj *jSchema) {
    const JsonObj *jUnitInfo = NULL;
    const JsonObj *jSWVer = NULL;
    const JsonObj *jProdSWVer = NULL;
    NodeInfo *pUnitInfo = NULL;
    int ret = 0;

    /* Unit Info */
    jUnitInfo = json_object_get(jSchema, JTAG_NODE_INFO);
    if (json_is_object(jUnitInfo)) {
        pUnitInfo = usys_zmalloc(sizeof(NodeInfo));
        if (pUnitInfo) {
            /* UUID */
            if (!parser_read_string_object_wrapper(jUnitInfo, JTAG_UUID,
                                                   pUnitInfo->uuid)) {
                goto cleanup;
            }

            /* Name */
            if (!parser_read_string_object_wrapper(jUnitInfo, JTAG_NAME,
                                                   pUnitInfo->name)) {
                goto cleanup;
            }

            /* Type */
            int uType = 0;
            if (!parser_read_integer_object(jUnitInfo, JTAG_TYPE, &uType)) {
                goto cleanup;
            } else {
                pUnitInfo->unit = (UnitType)uType;
            }

            /* Part Number */
            if (!parser_read_string_object_wrapper(jUnitInfo, JTAG_PART_NUMBER,
                                                   pUnitInfo->partNo)) {
                goto cleanup;
            }

            /* Skew */
            if (!parser_read_string_object_wrapper(jUnitInfo, JTAG_SKEW,
                                                   pUnitInfo->skew)) {
                goto cleanup;
            }

            /* MAC */
            if (!parser_read_string_object_wrapper(jUnitInfo, JTAG_MAC,
                                                   pUnitInfo->mac)) {
                goto cleanup;
            }

            /* SW Version */
            jSWVer = json_object_get(jUnitInfo, JTAG_SW_VERISION);
            Version *sVersion = parse_version(jSWVer);
            if (sVersion) {
                usys_memcpy(&pUnitInfo->swVer, sVersion, sizeof(Version));
                usys_free(sVersion);
            } else {
                usys_log_error("Failed to parse NodeInfo.swVersion.");
                goto cleanup;
            }

            /* Production SW Version */
            jProdSWVer = json_object_get(jUnitInfo, JTAG_PROD_SW_VERSION);
            Version *pVersion = parse_version(jProdSWVer);
            if (pVersion) {
                usys_memcpy(&pUnitInfo->swVer, pVersion, sizeof(Version));
                usys_free(pVersion);
            } else {
                usys_log_error("Failed to parse NodeInfo.prodSwVersion.");
                goto cleanup;
            }

            /* Assembly Date */
            if (!parser_read_string_object_wrapper(jUnitInfo, JTAG_ASM_DATE,
                                                   pUnitInfo->assmDate)) {
                goto cleanup;
            }

            /* OEM Name */
            if (!parser_read_string_object_wrapper(jUnitInfo, JTAG_OEM_NAME,
                                                   pUnitInfo->oemName)) {
                goto cleanup;
            }

            /* Module Count */
            uint8_t modCount = 0;
            if (!parser_read_uint8_object(jUnitInfo, JTAG_MODULE_COUNT,
                                          &modCount)) {
                goto cleanup;
            } else {
                pUnitInfo->modCount = modCount;
            }
        } else {
            usys_log_error(
                "Memory exhausted while parsing Unit Info. Error: %s",
                usys_error(errno));
            goto cleanup;
        }
    } else {
        usys_log_error("Unexpected JSON object found instead of Unit Info.");

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
    const JsonObj *jNodeCfgs = NULL;
    const JsonObj *jNodeCfg = NULL;
    const JsonObj *jDevice = NULL;
    NodeCfg *pNodeCfg = NULL;
    int iter = 0;

    if (!count) {
        goto cleanup;
    }

    jNodeCfgs = json_object_get(jSchema, JTAG_UNIT_CONFIG);
    if (json_is_array(jNodeCfgs)) {
        pNodeCfg = usys_zmalloc(sizeof(NodeCfg) * count);
        if (pNodeCfg) {
            json_array_foreach(jNodeCfgs, iter, jNodeCfg) {
                /* UUID */
                if (!parser_read_string_object_wrapper(
                        jNodeCfg, JTAG_UUID, pNodeCfg[iter].modUuid)) {
                    goto cleanup;
                }

                /* Name */
                if (!parser_read_string_object_wrapper(
                        jNodeCfg, JTAG_NAME, pNodeCfg[iter].modName)) {
                    goto cleanup;
                }

                /* SysFs */
                if (!parser_read_string_object_wrapper(
                        jNodeCfg, JTAG_INVT_SYSFS_FILE, pNodeCfg[iter].sysFs)) {
                    goto cleanup;
                }

                /* EEPROM */
                jDevice = json_object_get(jNodeCfg, JTAG_INVT_DEV_INFO);
                DevI2cCfg *pDev = parse_schema_devices(jDevice, DEV_CLASS_I2C);
                if (pDev) {
                    pNodeCfg[iter].eepromCfg = pDev;
                } else {
                    usys_log_error("Failed to parse NodeCfg[%d].eepromCfg",
                                   iter);
                    goto cleanup;
                }
            }

            /* Verify count of device read */
            if (iter == count) {
                usys_log_debug("All %d modules info read from the"
                               "manufacturing data's unit config.",
                               count);
            } else {
                usys_log_error("Only %d modules info read from the "
                               "manufacturing data's unit config expected "
                               "were %d.",
                               iter, count);
            }
        }
    } else {
        usys_log_error("Unexpected JSON object found instead"
                       "of Unit Config.");
        goto cleanup;
    }
    return pNodeCfg;

cleanup:
    parser_free_unit_cfg(&pNodeCfg, count);
    return pNodeCfg;
}

/* Parse Module Info */
void *parse_schema_module_info(const JsonObj *jSchema) {
    const JsonObj *jModuleInfo = NULL;
    const JsonObj *jSWVer = NULL;
    const JsonObj *jProdSWVer = NULL;

    ModuleInfo *pModuleInfo = NULL;
    int ret = 0;
    /* Module Info */
    jModuleInfo = json_object_get(jSchema, JTAG_MODULE_INFO);
    if (json_is_object(jModuleInfo)) {
        pModuleInfo = usys_zmalloc(sizeof(ModuleInfo));
        if (pModuleInfo) {
            /* UUID */
            if (!parser_read_string_object_wrapper(jModuleInfo, JTAG_UUID,
                                                   pModuleInfo->uuid)) {
                goto cleanup;
            }

            /* Name */
            if (!parser_read_string_object_wrapper(jModuleInfo, JTAG_NAME,
                                                   pModuleInfo->name)) {
                goto cleanup;
            }

            /* Type */
            int modType = 0;
            if (!parser_read_integer_object(jModuleInfo, JTAG_TYPE, &modType)) {
                goto cleanup;
            } else {
                pModuleInfo->module = (ModuleType)modType;
            }

            /* Part Number */
            if (!parser_read_string_object_wrapper(
                    jModuleInfo, JTAG_PART_NUMBER, pModuleInfo->partNo)) {
                goto cleanup;
            }

            /* HW Version */
            if (!parser_read_string_object_wrapper(jModuleInfo, JTAG_HW_VERSION,
                                                   pModuleInfo->hwVer)) {
                goto cleanup;
            }

            /* MAC */
            if (!parser_read_string_object_wrapper(jModuleInfo, JTAG_MAC,
                                                   pModuleInfo->mac)) {
                goto cleanup;
            }

            /* SW Version */
            jSWVer = json_object_get(jModuleInfo, JTAG_SW_VERISION);
            Version *sVersion = parse_version(jSWVer);
            if (sVersion) {
                usys_memcpy(&pModuleInfo->swVer, sVersion, sizeof(Version));
                usys_free(sVersion);
            } else {
                usys_log_error("Failed to parse NodeInfo.swVersion.");
                goto cleanup;
            }

            /* Production SW Version */
            jProdSWVer = json_object_get(jModuleInfo, JTAG_PROD_SW_VERSION);
            Version *pVersion = parse_version(jProdSWVer);
            if (pVersion) {
                usys_memcpy(&pModuleInfo->pSwVer, pVersion, sizeof(Version));
                usys_free(pVersion);
            } else {
                usys_log_error("Failed to parse NodeInfo.prodSwVersion.");
                goto cleanup;
            }

            /* Manufacturing Date  */
            if (!parser_read_string_object_wrapper(jModuleInfo, JTAG_MFG_DATE,
                                                   pModuleInfo->mfgDate)) {
                goto cleanup;
            }

            /* Manufacturer Name */
            if (!parser_read_string_object_wrapper(jModuleInfo, JTAG_MFG_NAME,
                                                   pModuleInfo->mfgName)) {
                goto cleanup;
            }

            /* Device Count */
            uint8_t devCount = 0;
            if (!parser_read_uint8_object(jModuleInfo, JTAG_DEVICE_COUNT,
                                          &devCount)) {
                goto cleanup;
            } else {
                pModuleInfo->devCount = devCount;
            }

        } else {
            usys_log_error("Memory exhausted while parsing Module Info. "
                           "Error %s",
                           usys_error(errno));
            goto cleanup;
        }

    } else {
        usys_log_error("Unexpected JSON object found instead"
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
                                                   pModCfg[iter].devName)) {
                goto cleanup;
            }

            /* Description */
            if (!parser_read_string_object_wrapper(jModCfg, JTAG_DESCRIPTION,
                                                   pModCfg[iter].devDesc)) {
                goto cleanup;
            }

            /* SysFs */
            if (!parser_read_string_object_wrapper(jModCfg, JTAG_DEV_SYSFS_FILE,
                                                   pModCfg[iter].sysFile)) {
                goto cleanup;
            }

            /* Device Class */
            uint16_t devClass = 0;
            if (!parser_read_uint16_object(jModCfg, JTAG_CLASS, &devClass)) {
                goto cleanup;
            } else {
                pModCfg[iter].devClass = devClass;
            }

            /* Device Type */
            uint16_t devType = 0;
            if (!parser_read_uint16_object(jModCfg, JTAG_TYPE, &devType)) {
                goto cleanup;
            } else {
                pModCfg[iter].devType = devType;
            }

            /* Device HW attributes */
            jDevice = json_object_get(jModCfg, JTAG_DEV_HW_ATTRS);
            if (jDevice) {
                DevI2cCfg *pDev =
                    parse_schema_devices(jDevice, pModCfg[iter].devClass);
                if (pDev) {
                    pModCfg[iter].cfg = pDev;
                } else {
                    usys_log_error("Failed to parse NodeCfg[%d].cfg", iter);
                    goto cleanup;
                }
            }
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
        usys_log_error("Memory exhausted while parsing Module Config. "
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
        if (!parser_read_string_object_wrapper(jSchema, name_key,
                                               genFileName)) {
            goto cleanup;
        }

        /* Read data size*/
        int rsize = read_mfg_data_size(genFileName);
        if (rsize <= 0) {
            usys_log_error("Failed to read data.", genFileName);
            goto cleanup;
        }

        /* Read data */
        genData = usys_zmalloc(sizeof(char) * rsize);
        if (genData) {
            ret = read_mfg_data(genFileName, genData, rsize);
            if (ret == rsize) {
                usys_log_debug("File %s read data of %d bytes.", genFileName,
                               rsize);
            } else {
                usys_log_error("File %s read data of %d bytes "
                               "expected %d bytes.",
                               genFileName, ret, rsize);
            }

            *size = ret;

        } else {
            usys_log_error("Memory exhausted while reading data from file %s.",
                           "Error: %s", genFileName, usys_error(errno));
            goto cleanup;
        }

    } else {
        usys_log_error("Memory exhausted while reading data file name.",
                       "Error: %s", usys_error(errno));
        goto cleanup;
    }

cleanup:
    usys_free(genFileName);
    return genData;
}

void *parse_schema_generic_file_data_wrapper(const JsonObj *jSchema, int iter,
                                             int *payloadSize, char *name_key,
                                             bool *status) {
    *status = USYS_FALSE;
    uint16_t size = 0;
    char *factCfg =
        parse_schema_generic_file_data(jSchema, &size, JTAG_FACTORY_CONFIG);
    if (factCfg) {
        if (size != *payloadSize) {
            usys_log_warn(
                "Size read for Field id 0x%x is %d bytes "
                "and size mentioned in index table [%d] is 0x%d bytes.",
                iter, size, *payloadSize, size);
            usys_log_debug("Updating index table [%d] size to %d bytes.", iter,
                           size);
            *payloadSize = size;
        }
        *status = USYS_TRUE;
    }
    return factCfg;
}

/* Parse payloads */
int parse_schema_payload(const JsonObj *jSchema, StoreSchema **schema,
                         uint16_t id, int iter) {
    int ret = 0;
    bool status = USYS_FALSE;
    switch (id) {
    case FIELD_ID_UNIT_INFO: {
        NodeInfo *pUnitInfo = parse_schema_node_info(jSchema);
        if (pUnitInfo) {
            usys_memcpy(&(*schema)->unitInfo, pUnitInfo, sizeof(NodeInfo));
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
        NodeCfg *pNodeCfg = parse_schema_unit_config(jSchema, modCount);
        if (pNodeCfg) {
            (*schema)->unitCfg = pNodeCfg;
        } else {
            ret = -1;
            goto cleanup;
        }
        break;
    }
    case FIELD_ID_MODULE_INFO: {
        ModuleInfo *pModuleInfo = parse_schema_module_info(jSchema);
        if (pModuleInfo) {
            usys_memcpy(&(*schema)->modInfo, pModuleInfo, sizeof(ModuleInfo));
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
        ModuleCfg *pModuleCfg = parse_schema_module_config(jSchema, devCount);
        if (pModuleCfg) {
            (*schema)->modCfg = pModuleCfg;
        } else {
            ret = -1;
            goto cleanup;
        }
        break;
    }
    case FIELD_ID_FACT_CFG: {
        int payloadSize = 0;
        (*schema)->factCfg = parse_schema_generic_file_data_wrapper(
            jSchema, iter, &payloadSize, JTAG_FACTORY_CONFIG, &status);

        if (!status) {
            ret = -1;
            goto cleanup;
        }
        (*schema)->indexTable[iter].payloadSize = payloadSize;

        break;
    }
    case FIELD_ID_USER_CFG: {
        int payloadSize = 0;
        (*schema)->userCfg = parse_schema_generic_file_data_wrapper(
            jSchema, iter, &payloadSize, JTAG_USER_CONFIG, &status);

        if (!status) {
            ret = -1;
            goto cleanup;
        }
        (*schema)->indexTable[iter].payloadSize = payloadSize;

        break;
    }
    case FIELD_ID_FACT_CALIB: {
        int payloadSize = 0;
        (*schema)->factCalib = parse_schema_generic_file_data_wrapper(
            jSchema, iter, &payloadSize, JTAG_FACTORY_CALIB, &status);

        if (!status) {
            ret = -1;
            goto cleanup;
        }
        (*schema)->indexTable[iter].payloadSize = payloadSize;

        break;
    }
    case FIELD_ID_USER_CALIB: {
        int payloadSize = 0;
        (*schema)->userCalib = parse_schema_generic_file_data_wrapper(
            jSchema, iter, &payloadSize, JTAG_USER_CALIB, &status);

        if (!status) {
            ret = -1;
            goto cleanup;
        }
        (*schema)->indexTable[iter].payloadSize = payloadSize;

        break;
    }
    case FIELD_ID_BS_CERTS: {
        int payloadSize = 0;
        (*schema)->bsCerts = parse_schema_generic_file_data_wrapper(
            jSchema, iter, &payloadSize, JTAG_BOOTSTRAP_CERTS, &status);

        if (!status) {
            ret = -1;
            goto cleanup;
        }
        (*schema)->indexTable[iter].payloadSize = payloadSize;

        break;
    }
    case FIELD_ID_CLOUD_CERTS: {
        int payloadSize = 0;
        (*schema)->cloudCerts = parse_schema_generic_file_data_wrapper(
            jSchema, iter, &payloadSize, JTAG_CLOUD_CERTS, &status);

        if (!status) {
            ret = -1;
            goto cleanup;
        }
        (*schema)->indexTable[iter].payloadSize = payloadSize;

        break;
    }
    default: {
        ret = ERR_NODED_INVALID_FIELD;
        usys_log_error("Invalid Field id supplied by Index entry.Error %d",
                       ret);
    }
    }

cleanup:
    return ret;
}

int parse_mfg_schema(const char *mfgdata, uint8_t idx) {
    int ret = 0;
    JsonObj *jSchema = NULL;
    JsonErrObj *jErr = NULL;
    const JsonObj *jHeader = NULL;
    const JsonObj *jIdxTable = NULL;
    StoreSchema *storeSchema = NULL;
    mfgStoreSchema[idx] = usys_zmalloc(sizeof(StoreSchema));
    if (mfgStoreSchema[idx]) {
        storeSchema = mfgStoreSchema[idx];
        jSchema = json_loads(mfgdata, JSON_DECODE_ANY, jErr);
        if (!jSchema) {
            parser_error(jErr, "Failed to parse schema");
            ret = ERR_NODED_JSON_PARSER;
            goto cleanup;
        }

        /* Debug Info */
        char *out = json_dumps(jSchema, (JSON_INDENT(4) | JSON_COMPACT |
                                         JSON_ENCODE_ANY));
        if (out) {
            usys_log_trace("Schema at Idx %d is ::\n %s\n", idx, out);
            usys_free(out);
            out = NULL;
        }

        /* Header */
        jHeader = json_object_get(jSchema, JTAG_HEADER);
        if (jHeader) {
            SchemaHeader *pHeader = parse_schema_header(jHeader);
            if (pHeader) {
                usys_memcpy(&storeSchema->header, pHeader,
                            sizeof(SchemaHeader));
                usys_free(pHeader);
                pHeader = NULL;
            } else {
                ret = ERR_NODED_INVALID_JSON_OBJECT;
                goto cleanup;
            }
        } else {
            ret = ERR_NODED_INVALID_JSON_OBJECT;
            goto cleanup;
        }

        /* Index Table */
        jIdxTable = json_object_get(jSchema, JTAG_INDEX_TABLE);
        if (jIdxTable) {
            SchemaIdxTuple *pIndexTable = parse_schema_idx_table(
                jIdxTable, storeSchema->header.idxCurTpl);
            if (pIndexTable) {
                /* Free me once done*/
                storeSchema->indexTable = pIndexTable;
            } else {
                ret = ERR_NODED_INVALID_JSON_OBJECT;
                goto cleanup;
            }

        } else {
            ret = ERR_NODED_INVALID_JSON_OBJECT;
            goto cleanup;
        }

        for (int iter = 0; iter < storeSchema->header.idxCurTpl; iter++) {
            uint16_t id = storeSchema->indexTable[iter].fieldId;
            ret = parse_schema_payload(jSchema, &storeSchema, id, iter);
            if (ret) {
                usys_log_error("Failed parsing for Field Id 0x%x from mfg data"
                               ".Error: %d",
                               id, ret);
                goto cleanup;
            } else {
                usys_log_debug("Parsing for Field Id 0x%x from mfg data "
                               "completed.",
                               id);
            }
        }
    }

cleanup:
    if (jSchema) {
        json_decref(jSchema);
        jSchema = NULL;
    }
    if (ret) {
        parser_free_mfg_data(&storeSchema);
    }
    return ret;
}

StoreSchema *parser_get_mfg_data_by_uuid(char *puuid) {
    StoreSchema *sschema = NULL;

    /*Default section*/
    if ((!puuid) || !usys_strcmp(puuid, "")) {
        sschema = mfgStoreSchema[0];
        usys_log_trace("MFG Data set to the Module UUID %s with %d entries"
                       " in Index Table.",
                       sschema->modInfo.uuid, sschema->header.idxCurTpl);
    } else {
        /* Searching for Module MFG data index.*/
        for (uint8_t iter = 0; iter < MAX_JSON_SCHEMA; iter++) {
            if (mfgStoreSchema[iter]) {
                if (!usys_strcmp(puuid, mfgStoreSchema[iter]->modInfo.uuid)) {
                    sschema = mfgStoreSchema[iter];
                    usys_log_trace("PARSER:: MFG Data set to the Module UUID"
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
                    usys_log_debug("File %s read manufacturing data of %d"
                                   " bytes.",
                                   fname, size);

                    /* Parse the mfg data to store schema */
                    ret = parse_mfg_schema(schemabuff, iter);
                    if (ret) {
                        usys_log_error("Err(%d):PARSER: Parsing failed for %s.",
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

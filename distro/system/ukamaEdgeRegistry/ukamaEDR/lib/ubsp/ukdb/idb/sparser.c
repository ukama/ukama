/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "ukdb/idb/sparser.h"

#include "headers/ubsp/devices.h"
#include "headers/errorcode.h"
#include "inc/globalheader.h"
#include "headers/utils/log.h"

#include "utils/cJSON.h"

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <errno.h>
#include <fcntl.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>

UKDB *g_pukdb[MAX_JSON_SCHEMA] = { '\0' };

static int read_mfg_data(char *fname, char *buff, off_t size) {
    int read_bytes = 0;
    /* Create input file descriptor */
    int fd = open(fname, O_RDONLY, 0644);
    if (fd == -1) {
        perror("open");
        return fd;
    }
    off_t off = lseek(fd, 0, SEEK_SET);
    read_bytes = read(fd, buff, size);
    return read_bytes;
}

static int read_mfg_data_size(char *fname) {
    int read_bytes = 0;
    /* Create input file descriptor */
    int fd = open(fname, O_RDONLY, 0644);
    if (fd == -1) {
        perror("open");
        return fd;
    }
    off_t off = lseek(fd, 0L, SEEK_END);
    return off;
}

void parser_free_unit_cfg(UnitCfg **cfg, uint8_t count) {
    if (*cfg) {
        /* Device Cfgs */
        for (int iter = 0; iter < count; iter++) {
            UBSP_FREE((*cfg)->eeprom_cfg);
        }
        UBSP_FREE(*cfg);
    }
}

void parser_free_ukdb_mfg_data(UKDB **pukdb) {
    if (*pukdb) {
        UBSP_FREE((*pukdb)->indextable);
        /* Unit Cfgs */
        for (int iter = 0; iter < (*pukdb)->unitinfo.mod_count; iter++) {
            UBSP_FREE((*pukdb)->unitcfg[iter].eeprom_cfg);
        }

        UBSP_FREE((*pukdb)->unitcfg);
        /* Module Cfgs */
        for (int iter = 0; iter < (*pukdb)->modinfo.dev_count; iter++) {
            UBSP_FREE((*pukdb)->modcfg[iter].cfg);
        }
        UBSP_FREE((*pukdb)->modcfg);
        UBSP_FREE((*pukdb)->factcfg);
        UBSP_FREE((*pukdb)->usercfg);
        UBSP_FREE((*pukdb)->factcalib);
        UBSP_FREE((*pukdb)->usercalib);
        UBSP_FREE((*pukdb)->bscerts);
        UBSP_FREE((*pukdb)->lwm2mcerts);
        UBSP_FREE(*pukdb);
    }
}

Version *parse_schema_version(const cJSON *version) {
    const cJSON *major = NULL;
    const cJSON *minor = NULL;
    Version *pversion = malloc(sizeof(Version));
    if (pversion) {
        memset(pversion, 0, sizeof(Version));
        major = cJSON_GetObjectItemCaseSensitive(version, "major");
        if (cJSON_IsNumber(major)) {
            pversion->major = version->valueint;
        } else {
            UBSP_FREE(pversion);
            return pversion;
        }
        minor = cJSON_GetObjectItemCaseSensitive(version, "minor");
        if (cJSON_IsNumber(minor)) {
            pversion->minor = version->valueint;
        } else {
            UBSP_FREE(pversion);
            return pversion;
        }
    }
    return pversion;
}

UKDBHeader *parse_schema_header(const cJSON *header) {
    const cJSON *dbversion = NULL;
    const cJSON *index_tbl_offset = NULL;
    const cJSON *index_tpl_size = NULL;
    const cJSON *index_tuple_max_count = NULL;
    const cJSON *index_current_tuple = NULL;
    const cJSON *module_capability = NULL;
    const cJSON *module_mode = NULL;
    const cJSON *module_device_owner = NULL;
    UKDBHeader *pheader = malloc(sizeof(UKDBHeader));
    if (pheader) {
        memset(pheader, 0, sizeof(UKDBHeader));

        dbversion = cJSON_GetObjectItemCaseSensitive(header, "dbversion");
        Version *pversion = parse_schema_version(dbversion);
        if (pversion) {
            memcpy(&pheader->dbversion, pversion, sizeof(Version));
            UBSP_FREE(pversion);
        } else {
            goto cleanup;
        }

        index_tbl_offset =
            cJSON_GetObjectItemCaseSensitive(header, "index_table_offset");
        if (cJSON_IsNumber(index_tbl_offset)) {
            pheader->idx_tbl_offset = index_tbl_offset->valueint;
        } else {
            goto cleanup;
        }

        index_tpl_size =
            cJSON_GetObjectItemCaseSensitive(header, "index_tuple_size");
        if (cJSON_IsNumber(index_tpl_size)) {
            pheader->idx_tpl_size = index_tpl_size->valueint;
        } else {
            goto cleanup;
        }

        index_tuple_max_count =
            cJSON_GetObjectItemCaseSensitive(header, "index_tuple_max_count");
        if (cJSON_IsNumber(index_tuple_max_count)) {
            pheader->idx_tpl_max_count = index_tuple_max_count->valueint;
        } else {
            goto cleanup;
        }

        index_current_tuple =
            cJSON_GetObjectItemCaseSensitive(header, "index_current_tuple");
        if (cJSON_IsNumber(index_current_tuple)) {
            pheader->idx_cur_tpl = index_current_tuple->valueint;
        } else {
            goto cleanup;
        }

        module_capability =
            cJSON_GetObjectItemCaseSensitive(header, "module_capability");
        if (cJSON_IsString(module_capability)) {
            if (!strcmp(module_capability->valuestring, "AUTONOMOUS")) {
                pheader->mod_cap = 1;
            } else {
                pheader->mod_cap = 0;
            }
        } else {
            goto cleanup;
        }

        module_mode = cJSON_GetObjectItemCaseSensitive(header, "module_mode");
        if (cJSON_IsString(module_mode)) {
            if (!strcmp(module_mode->valuestring, "MASTER")) {
                pheader->mod_mode = 1;
            } else {
                pheader->mod_mode = 0;
            }
        } else {
            goto cleanup;
        }

        module_device_owner =
            cJSON_GetObjectItemCaseSensitive(header, "module_device_owner");
        if (cJSON_IsString(module_device_owner)) {
            if (!strcmp(module_device_owner->valuestring, "DEVICE_OWNER")) {
                pheader->mod_devown = 1;
            } else {
                pheader->mod_devown = 0;
            }
        } else {
            goto cleanup;
        }
    }
    return pheader;

cleanup:
    UBSP_FREE(pheader);
    return pheader;
}

UKDBIdxTuple *parse_schema_index_table(const cJSON *index_table,
                                       uint8_t count) {
    const cJSON *index_tpl = NULL;
    const cJSON *field_id = NULL;
    const cJSON *payload_offset = NULL;
    const cJSON *payload_size = NULL;
    const cJSON *payload_version = NULL;
    const cJSON *dbversion = NULL;
    const cJSON *payload_crc = NULL;
    const cJSON *state = NULL;
    const cJSON *valid = NULL;
    UKDBIdxTuple *pindex_table = NULL;
    int iter = 0;
    if (!count) {
        goto cleanup;
    }
    pindex_table = malloc(sizeof(UKDBIdxTuple) * count);
    if (pindex_table) {
        memset(pindex_table, 0, sizeof(UKDBIdxTuple) * count);
        cJSON_ArrayForEach(index_tpl, index_table) {
            if (iter >= count) {
                log_warn(
                    "PARSER: Expected index tuples are %d but seems like schema has more.",
                    count);
                break;
            }
            field_id = cJSON_GetObjectItemCaseSensitive(index_tpl, "field_id");
            if (cJSON_IsNumber(field_id)) {
                pindex_table[iter].fieldid = field_id->valueint;
            } else {
                goto cleanup;
            }

            payload_offset =
                cJSON_GetObjectItemCaseSensitive(index_tpl, "payload_offset");
            if (cJSON_IsNumber(payload_offset)) {
                pindex_table[iter].payload_offset = payload_offset->valueint;
            } else {
                goto cleanup;
            }

            payload_size =
                cJSON_GetObjectItemCaseSensitive(index_tpl, "payload_size");
            if (cJSON_IsNumber(payload_size)) {
                pindex_table[iter].payload_size = payload_size->valueint;
            } else {
                goto cleanup;
            }

            payload_version =
                cJSON_GetObjectItemCaseSensitive(index_tpl, "payload_version");
            Version *pversion = parse_schema_version(payload_version);
            if (pversion) {
                memcpy(&pindex_table[iter].payload_version, pversion,
                       sizeof(Version));
                UBSP_FREE(pversion);
            } else {
                goto cleanup;
            }

            payload_crc =
                cJSON_GetObjectItemCaseSensitive(index_tpl, "payload_crc");
            if (cJSON_IsNumber(payload_crc)) {
                pindex_table[iter].payload_crc = payload_crc->valueint;
            } else {
                goto cleanup;
            }

            state = cJSON_GetObjectItemCaseSensitive(index_tpl, "state");
            if (cJSON_IsString(state)) {
                if (!strcmp(state->valuestring, "ENABLED")) {
                    pindex_table[iter].state = 1;
                } else {
                    pindex_table[iter].state = 0;
                }
            } else {
                goto cleanup;
            }

            valid = cJSON_GetObjectItemCaseSensitive(index_tpl, "valid");
            if (cJSON_IsString(valid)) {
                if (!strcmp(valid->valuestring, "TRUE")) {
                    pindex_table[iter].valid = 1;
                } else {
                    pindex_table[iter].valid = 0;
                }
            } else {
                goto cleanup;
            }
            iter++;
        }
        if (iter == count) {
            log_debug(
                "PARSER:: All %d tuples read from the manufacturing data.",
                count);
        } else {
            log_error(
                "Err: PARSER:: Only %d tuples read from the manufacturing data expected were %d.",
                (count - iter), count);
        }
    } else {
        log_error("Err(%d):PARSER: Memory exhausted.",
                  ERR_UBSP_MEMORY_EXHAUSTED);
    }
    return pindex_table;
cleanup:
    UBSP_FREE(pindex_table);
    return pindex_table;
}

DevGpioCfg *parse_schema_dev_gpio(const cJSON *schema) {
    const cJSON *gpio_cfg = NULL;
    const cJSON *direction = NULL;
    const cJSON *number = NULL;

    DevGpioCfg *pdev_cfg = NULL;
    pdev_cfg = malloc(sizeof(DevGpioCfg));
    if (pdev_cfg) {
        memset(pdev_cfg, 0, sizeof(DevGpioCfg));
        direction = cJSON_GetObjectItemCaseSensitive(schema, "direction");
        if (cJSON_IsString(direction)) {
            if (!strcmp(direction->valuestring, "input")) {
                pdev_cfg->direction = 0;
            } else {
                pdev_cfg->direction = 1;
            }
        } else {
            log_error("Err: PARSER:: Failed to parse DevGpioCfg.direction");
            goto cleanup;
        }

        number = cJSON_GetObjectItemCaseSensitive(schema, "number");
        if (cJSON_IsNumber(number)) {
            pdev_cfg->gpio_num = number->valueint;
        } else {
            log_error("Err: PARSER:: Failed to parse DevGpioCfg.number");
            goto cleanup;
        }
    } else {
        log_error("Err(%d):PARSER: Memory exhausted while parsing DevGpioCfg.",
                  ERR_UBSP_MEMORY_EXHAUSTED);
        goto cleanup;
    }
    return pdev_cfg;
cleanup:
    UBSP_FREE(pdev_cfg);
    return pdev_cfg;
}

DevI2cCfg *parse_schema_dev_i2c(const cJSON *schema) {
    const cJSON *dev_cfg = NULL;
    const cJSON *bus = NULL;
    const cJSON *address = NULL;

    DevI2cCfg *pdev_cfg = NULL;
    pdev_cfg = malloc(sizeof(DevI2cCfg));
    if (pdev_cfg) {
        memset(pdev_cfg, 0, sizeof(DevI2cCfg));
        bus = cJSON_GetObjectItemCaseSensitive(schema, "bus");
        if (cJSON_IsNumber(bus)) {
            pdev_cfg->bus = bus->valueint;
        } else {
            log_error("Err: PARSER:: Failed to parse DevI2cCfg.bus");
            goto cleanup;
        }

        address = cJSON_GetObjectItemCaseSensitive(schema, "address");
        if (cJSON_IsNumber(address)) {
            pdev_cfg->add = address->valueint;
        } else {
            log_error("Err: PARSER:: Failed to parse DevI2cCfg.address");
            goto cleanup;
        }
    } else {
        log_error("Err(%d):PARSER: Memory exhausted while parsing DevI2cCfg.",
                  ERR_UBSP_MEMORY_EXHAUSTED);
        goto cleanup;
    }
    return pdev_cfg;
cleanup:
    UBSP_FREE(pdev_cfg);
    return pdev_cfg;
}

DevSpiCfg *parse_schema_dev_spi(const cJSON *schema) {
    const cJSON *dev_cfg = NULL;
    const cJSON *bus = NULL;
    const cJSON *cs = NULL;

    DevSpiCfg *pdev_cfg = NULL;
    pdev_cfg = malloc(sizeof(DevSpiCfg));
    if (dev_cfg) {
        memset(pdev_cfg, 0, sizeof(DevSpiCfg));
        cs = cJSON_GetObjectItemCaseSensitive(schema, "cs");
        if (cJSON_IsObject(cs)) {
            DevGpioCfg *pcs = parse_schema_dev_gpio(cs);
            if (cs) {
                memcpy(&pdev_cfg->cs, pcs, sizeof(DevGpioCfg));
            } else {
                log_error("Err: PARSER:: Failed to parse DevSpiCfg.cs");
                goto cleanup;
            }
        } else {
            log_error("Err: PARSER:: Failed to parse DevSpiCfg.cs");
            goto cleanup;
        }

        bus = cJSON_GetObjectItemCaseSensitive(schema, "bus");
        if (cJSON_IsNumber(bus)) {
            pdev_cfg->bus = bus->valueint;
        } else {
            log_error("Err: PARSER:: Failed to parse DevSpiCfg.bus");
            goto cleanup;
        }
    } else {
        log_error("Err(%d):PARSER: Memory exhausted while parsing DevSpiCfg.",
                  ERR_UBSP_MEMORY_EXHAUSTED);
        goto cleanup;
    }
    return pdev_cfg;
cleanup:
    UBSP_FREE(pdev_cfg);
    return pdev_cfg;
}

DevUartCfg *parse_schema_dev_uart(const cJSON *schema) {
    const cJSON *dev_cfg = NULL;
    const cJSON *uartno = NULL;

    DevUartCfg *pdev_cfg = NULL;
    pdev_cfg = malloc(sizeof(DevUartCfg));
    if (dev_cfg) {
        memset(pdev_cfg, 0, sizeof(DevUartCfg));
        uartno = cJSON_GetObjectItemCaseSensitive(schema, "uartno");
        if (cJSON_IsNumber(uartno)) {
            pdev_cfg->uartno = uartno->valueint;
        } else {
            log_error("Err: PARSER:: Failed to parse DevUartCfg.bus");
            goto cleanup;
        }
    } else {
        log_error("Err(%d):PARSER: Memory exhausted while parsing DevUartCfg.",
                  ERR_UBSP_MEMORY_EXHAUSTED);
        goto cleanup;
    }
    return pdev_cfg;
cleanup:
    UBSP_FREE(pdev_cfg);
    return pdev_cfg;
}

void *parse_schema_devices(const cJSON *schema, uint16_t class) {
    void *dev = NULL;
    switch (class) {
    case DEV_CLASS_GPIO:
        dev = parse_schema_dev_gpio(schema);
        break;
    case DEV_CLASS_I2C:
        dev = parse_schema_dev_i2c(schema);
        break;
    case DEV_CLASS_SPI:
        dev = parse_schema_dev_spi(schema);
        break;
    case DEV_CLASS_UART:
        dev = parse_schema_dev_uart(schema);
        break;
    default:
        log_error("Err(%d): PARSER:: Unkown device type failed to parse.",
                  ERR_UBSP_INVALID_DEVICE_CFG);
    }
    return dev;
}

void *parse_schema_unit_info(const cJSON *schema) {
    const cJSON *unit_info = NULL;
    const cJSON *uuid = NULL;
    const cJSON *name = NULL;
    const cJSON *type = NULL;
    const cJSON *partno = NULL;
    const cJSON *skew = NULL;
    const cJSON *mac = NULL;
    const cJSON *sw_version = NULL;
    const cJSON *production_sw_version = NULL;
    const cJSON *assembly_date = NULL;
    const cJSON *oem_name = NULL;
    const cJSON *module_count = NULL;
    UnitInfo *punit_info = NULL;
    int ret = 0;
    /* Unit Info */
    unit_info = cJSON_GetObjectItemCaseSensitive(schema, "unit_info");
    if (cJSON_IsObject(unit_info)) {
        punit_info = malloc(sizeof(UnitInfo));
        if (punit_info) {
            memset(punit_info, '\0', sizeof(UnitInfo));
            uuid = cJSON_GetObjectItemCaseSensitive(unit_info, "UUID");
            if (cJSON_IsString(uuid)) {
                memset(punit_info->uuid, '\0', 24);
                memcpy(punit_info->uuid, uuid->valuestring,
                       strlen(uuid->valuestring));
            } else {
                log_error("Err: PARSER:: Failed to parse UnitInfo.uuid.");
                goto cleanup;
            }

            name = cJSON_GetObjectItemCaseSensitive(unit_info, "name");
            if (cJSON_IsString(name)) {
                memset(punit_info->name, '\0', 24);
                memcpy(punit_info->name, name->valuestring,
                       strlen(name->valuestring));
            } else {
                log_error("Err: PARSER:: Failed to parse UnitInfo.name.");
                goto cleanup;
            }

            type = cJSON_GetObjectItemCaseSensitive(unit_info, "type");
            if (cJSON_IsNumber(type)) {
                punit_info->unit = type->valueint;
            } else {
                log_error("Err: PARSER:: Failed to parse UnitInfo.unit.");
                goto cleanup;
            }

            partno = cJSON_GetObjectItemCaseSensitive(unit_info, "partno");
            if (cJSON_IsString(partno)) {
                memset(punit_info->partno, '\0', 24);
                memcpy(punit_info->partno, partno->valuestring,
                       strlen(partno->valuestring));
            } else {
                log_error("Err: PARSER:: Failed to parse UnitInfo.partno.");
                goto cleanup;
            }

            skew = cJSON_GetObjectItemCaseSensitive(unit_info, "skew");
            if (cJSON_IsString(skew)) {
                memset(punit_info->skew, '\0', 24);
                memcpy(punit_info->skew, skew->valuestring,
                       strlen(skew->valuestring));
            } else {
                log_error("Err: PARSER:: Failed to parse UnitInfo.skew.");
                goto cleanup;
            }

            mac = cJSON_GetObjectItemCaseSensitive(unit_info, "mac");
            if (cJSON_IsString(mac)) {
                memset(punit_info->mac, '\0', 24);
                memcpy(punit_info->mac, mac->valuestring,
                       strlen(mac->valuestring));
            } else {
                log_error("Err: PARSER:: Failed to parse UnitInfo.mac.");
                goto cleanup;
            }

            sw_version =
                cJSON_GetObjectItemCaseSensitive(unit_info, "sw_version");
            Version *pversion = parse_schema_version(sw_version);
            if (pversion) {
                memcpy(&punit_info->swver, pversion, sizeof(Version));
                UBSP_FREE(pversion);
            } else {
                log_error("Err: PARSER:: Failed to parse UnitInfo.sw_version.");
                goto cleanup;
            }

            production_sw_version = cJSON_GetObjectItemCaseSensitive(
                unit_info, "production_sw_version");
            pversion = parse_schema_version(production_sw_version);
            if (pversion) {
                memcpy(&punit_info->pswver, pversion, sizeof(Version));
                UBSP_FREE(pversion);
            } else {
                log_error(
                    "Err: PARSER:: Failed to parse UnitInfo.production_sw_version.");
                goto cleanup;
            }

            assembly_date =
                cJSON_GetObjectItemCaseSensitive(unit_info, "assembly_date");
            if (cJSON_IsString(assembly_date)) {
                memset(punit_info->assm_date, '\0', 24);
                memcpy(punit_info->assm_date, assembly_date->valuestring,
                       strlen(assembly_date->valuestring));
            } else {
                log_error(
                    "Err: PARSER:: Failed to parse UnitInfo.assembly_date.");
                goto cleanup;
            }

            oem_name = cJSON_GetObjectItemCaseSensitive(unit_info, "OEM_name");
            if (cJSON_IsString(oem_name)) {
                memset(punit_info->oem_name, '\0', 24);
                memcpy(punit_info->oem_name, oem_name->valuestring,
                       strlen(oem_name->valuestring));
            } else {
                log_error("Err: PARSER:: Failed to parse UnitInfo.oem_name.");
                goto cleanup;
            }

            module_count =
                cJSON_GetObjectItemCaseSensitive(unit_info, "module_count");
            if (cJSON_IsNumber(module_count)) {
                punit_info->mod_count = module_count->valueint;
            } else {
                log_error("Err: PARSER:: Failed to parse UnitInfo.mod_count");
                goto cleanup;
            }

        } else {
            ret = ERR_UBSP_MEMORY_EXHAUSTED;
            log_error(
                "Err(%d):PARSER: Memory exhausted while parsing Unit Info.",
                ret);
            goto cleanup;
        }
    } else {
        ret = ERR_UBSP_UNEXPECTED_JSON_OBJECT;
        log_error(
            "Err(%d):PARSER: Unexpected JSON object found instead of Unit Info.",
            ret);
        goto cleanup;
    }
    return punit_info;
cleanup:
    UBSP_FREE(punit_info);
    return punit_info;
}

void *parse_schema_unit_config(const cJSON *schema, uint16_t count) {
    const cJSON *unit_cfgs = NULL;
    const cJSON *unit_cfg = NULL;
    const cJSON *uuid = NULL;
    const cJSON *name = NULL;
    const cJSON *sysfs = NULL;
    const cJSON *device = NULL;
    UnitCfg *punit_cfg = NULL;
    int iter = 0;
    if (!count) {
        goto cleanup;
    }
    unit_cfgs = cJSON_GetObjectItemCaseSensitive(schema, "unit_config");
    if (cJSON_IsArray(unit_cfgs)) {
        punit_cfg = malloc(sizeof(UnitCfg) * count);
        if (punit_cfg) {
            memset(punit_cfg, '\0', sizeof(UnitCfg) * count);
            cJSON_ArrayForEach(unit_cfg, unit_cfgs) {
                if (iter >= count) {
                    log_warn(
                        "PARSER: Expected modules are %d but seems like schema has more.",
                        count);
                    break;
                }
                uuid = cJSON_GetObjectItemCaseSensitive(unit_cfg, "UUID");
                if (cJSON_IsString(uuid)) {
                    memset(punit_cfg[iter].mod_uuid, '\0', 24);
                    memcpy(punit_cfg[iter].mod_uuid, uuid->valuestring,
                           strlen(uuid->valuestring));
                } else {
                    log_error(
                        "Err: PARSER:: Failed to parse UnitCfg[%d].mod_uuid",
                        iter);
                    goto cleanup;
                }

                name = cJSON_GetObjectItemCaseSensitive(unit_cfg, "name");
                if (cJSON_IsString(name)) {
                    memset(punit_cfg[iter].mod_name, '\0', 24);
                    memcpy(punit_cfg[iter].mod_name, name->valuestring,
                           strlen(name->valuestring));
                } else {
                    log_error(
                        "Err: PARSER:: Failed to parse UnitCfg[%d].mod_name");
                    goto cleanup;
                }

                sysfs = cJSON_GetObjectItemCaseSensitive(unit_cfg, "dbsysfs");
                if (cJSON_IsString(sysfs)) {
                    memset(punit_cfg[iter].sysfs, '\0', 64);
                    memcpy(punit_cfg[iter].sysfs, sysfs->valuestring,
                           strlen(sysfs->valuestring));
                } else {
                    log_error("Err: PARSER:: Failed to parse UnitCfg[%d].sysfs",
                              iter);
                    goto cleanup;
                }

                device = cJSON_GetObjectItemCaseSensitive(unit_cfg, "devicedb");
                DevI2cCfg *pdev = parse_schema_devices(device, DEV_CLASS_I2C);
                if (pdev) {
                    //memcpy(&punit_cfg[iter].eeprom_cfg, pdev, sizeof(DevI2cCfg));
                    punit_cfg[iter].eeprom_cfg = pdev;
                } else {
                    log_error(
                        "Err: PARSER:: Failed to parse UnitCfg[%d].eeprom_cfg",
                        iter);
                    goto cleanup;
                }
                iter++;
            }
            if (iter == count) {
                log_debug(
                    "PARSER:: All %d modules info read from the manufacturing data's unit config.",
                    count);
            } else {
                log_error(
                    "PARSER::Only %d modules info read from the manufacturing data's unit config expected were %d.",
                    iter, count);
            }
        }
    } else {
        log_error(
            "Err(%d):PARSER: Unexpected JSON object found instead of Unit Config.",
            ERR_UBSP_UNEXPECTED_JSON_OBJECT);
        goto cleanup;
    }
    return punit_cfg;
cleanup:
    parser_free_unit_cfg(&punit_cfg, count);
    return punit_cfg;
}

void *parse_schema_module_info(const cJSON *schema) {
    const cJSON *module_info = NULL;
    const cJSON *uuid = NULL;
    const cJSON *name = NULL;
    const cJSON *type = NULL;
    const cJSON *partno = NULL;
    const cJSON *hwver = NULL;
    const cJSON *mac = NULL;
    const cJSON *sw_version = NULL;
    const cJSON *production_sw_version = NULL;
    const cJSON *mfg_date = NULL;
    const cJSON *mfg_name = NULL;
    const cJSON *dev_count = NULL;

    ModuleInfo *pmodule_info = NULL;
    int ret = 0;
    /* Module Info */
    module_info = cJSON_GetObjectItemCaseSensitive(schema, "module_info");
    if (cJSON_IsObject(module_info)) {
        pmodule_info = malloc(sizeof(ModuleInfo));
        if (pmodule_info) {
            memset(pmodule_info, '\0', sizeof(ModuleInfo));
            uuid = cJSON_GetObjectItemCaseSensitive(module_info, "UUID");
            if (cJSON_IsString(uuid)) {
                memset(pmodule_info->uuid, '\0', 24);
                memcpy(pmodule_info->uuid, uuid->valuestring,
                       strlen(uuid->valuestring));
            } else {
                log_error("Err: PARSER:: Failed to parse ModuleInfo.uuid");
                goto cleanup;
            }

            name = cJSON_GetObjectItemCaseSensitive(module_info, "name");
            if (cJSON_IsString(name)) {
                memset(pmodule_info->name, '\0', 24);
                memcpy(pmodule_info->name, name->valuestring,
                       strlen(name->valuestring));
            } else {
                log_error("Err: PARSER:: Failed to parse ModuleInfo.name");
                goto cleanup;
            }

            type = cJSON_GetObjectItemCaseSensitive(module_info, "type");
            if (cJSON_IsNumber(type)) {
                pmodule_info->module = type->valueint;
            } else {
                log_error("Err: PARSER:: Failed to parse ModuleInfo.type");
                goto cleanup;
            }

            partno = cJSON_GetObjectItemCaseSensitive(module_info, "partno");
            if (cJSON_IsString(partno)) {
                memset(pmodule_info->partno, '\0', 24);
                memcpy(pmodule_info->partno, partno->valuestring,
                       strlen(partno->valuestring));
            } else {
                log_error("Err: PARSER:: Failed to parse ModuleInfo.partno");
                goto cleanup;
            }

            hwver = cJSON_GetObjectItemCaseSensitive(module_info, "hw_version");
            if (cJSON_IsString(hwver)) {
                memset(pmodule_info->hwver, '\0', 24);
                memcpy(pmodule_info->hwver, hwver->valuestring,
                       strlen(hwver->valuestring));
            } else {
                log_error("Err: PARSER:: Failed to parse ModuleInfo.hwver");
                goto cleanup;
            }

            mac = cJSON_GetObjectItemCaseSensitive(module_info, "mac");
            if (cJSON_IsString(mac)) {
                memset(pmodule_info->mac, '\0', 24);
                memcpy(pmodule_info->mac, mac->valuestring,
                       strlen(mac->valuestring));
            } else {
                log_error("Err: PARSER:: Failed to parse ModuleInfo.mac");
                goto cleanup;
            }

            sw_version =
                cJSON_GetObjectItemCaseSensitive(module_info, "sw_version");
            Version *pversion = parse_schema_version(sw_version);
            if (pversion) {
                memcpy(&pmodule_info->swver, pversion, sizeof(Version));
                UBSP_FREE(pversion);
            } else {
                log_error(
                    "Err: PARSER:: Failed to parse ModuleInfo.sw_version");
                goto cleanup;
            }

            production_sw_version = cJSON_GetObjectItemCaseSensitive(
                module_info, "production_sw_version");
            pversion = parse_schema_version(production_sw_version);
            if (pversion) {
                memcpy(&pmodule_info->pswver, pversion, sizeof(Version));
                UBSP_FREE(pversion);
            } else {
                log_error(
                    "Err: PARSER:: Failed to parse ModuleInfo.production_sw_version");
                goto cleanup;
            }

            mfg_date = cJSON_GetObjectItemCaseSensitive(module_info,
                                                        "manufacturing_date");
            if (cJSON_IsString(mfg_date)) {
                memset(pmodule_info->mfg_date, '\0', 24);
                memcpy(pmodule_info->mfg_date, mfg_date->valuestring,
                       strlen(mfg_date->valuestring));
            } else {
                log_error("Err: PARSER:: Failed to parse ModuleInfo.mfg_date");
                goto cleanup;
            }

            mfg_name = cJSON_GetObjectItemCaseSensitive(module_info,
                                                        "manufacturer_name");
            if (cJSON_IsString(mfg_name)) {
                memset(pmodule_info->mfg_name, '\0', 24);
                memcpy(pmodule_info->mfg_name, mfg_name->valuestring,
                       strlen(mfg_name->valuestring));
            } else {
                log_error("Err: PARSER:: Failed to parse ModuleInfo.mfg_name");
                goto cleanup;
            }

            dev_count =
                cJSON_GetObjectItemCaseSensitive(module_info, "device_count");
            if (cJSON_IsNumber(dev_count)) {
                pmodule_info->dev_count = dev_count->valueint;
            } else {
                log_error("Err: PARSER:: Failed to parse ModuleInfo.dev_count");
                goto cleanup;
            }

        } else {
            ret = ERR_UBSP_MEMORY_EXHAUSTED;
            log_error(
                "Err(%d):PARSER: Memory exhausted while parsing Module Info.",
                ret);
            goto cleanup;
        }
    } else {
        ret = ERR_UBSP_UNEXPECTED_JSON_OBJECT;
        log_error(
            "Err(%d):PARSER: Unexpected JSON object found instead of Module Info.",
            ret);
        goto cleanup;
    }
    return pmodule_info;
cleanup:
    UBSP_FREE(pmodule_info);
    return pmodule_info;
}

void *parse_schema_module_config(const cJSON *schema, uint16_t count) {
    const cJSON *module_cfgs = NULL;
    const cJSON *module_cfg = NULL;
    const cJSON *name = NULL;
    const cJSON *desc = NULL;
    const cJSON *type = NULL;
    const cJSON *class = NULL;
    const cJSON *sysfs = NULL;
    const cJSON *device = NULL;
    ModuleCfg *pmodule_cfg = NULL;
    int iter = 0;
    if (!count) {
        goto cleanup;
    }
    module_cfgs = cJSON_GetObjectItemCaseSensitive(schema, "module_config");
    if (cJSON_IsArray(module_cfgs)) {
        pmodule_cfg = malloc(sizeof(ModuleCfg) * count);
        if (pmodule_cfg) {
            memset(pmodule_cfg, '\0', sizeof(ModuleCfg) * count);
            cJSON_ArrayForEach(module_cfg, module_cfgs) {
                if (iter >= count) {
                    log_warn(
                        "PARSER: Expected devices are %d but seems like schema has more.",
                        count);
                    break;
                }
                name = cJSON_GetObjectItemCaseSensitive(module_cfg, "name");
                if (cJSON_IsString(name)) {
                    memset(pmodule_cfg[iter].dev_name, '\0', 24);
                    memcpy(pmodule_cfg[iter].dev_name, name->valuestring,
                           strlen(name->valuestring));
                } else {
                    log_error(
                        "Err: PARSER:: Failed to parse ModuleCfg[%d].name",
                        iter);
                    goto cleanup;
                }

                desc =
                    cJSON_GetObjectItemCaseSensitive(module_cfg, "description");
                if (cJSON_IsString(desc)) {
                    memset(pmodule_cfg[iter].dev_disc, '\0', 24);
                    memcpy(pmodule_cfg[iter].dev_disc, desc->valuestring,
                           strlen(desc->valuestring));
                } else {
                    log_error(
                        "Err: PARSER:: Failed to parse ModuleCfg[%d].dev_disc",
                        iter);
                    goto cleanup;
                }

                type = cJSON_GetObjectItemCaseSensitive(module_cfg, "type");
                if (cJSON_IsNumber(type)) {
                    pmodule_cfg[iter].dev_type = type->valueint;
                } else {
                    log_error(
                        "Err: PARSER:: Failed to parse ModuleCfg[%d].dev_type",
                        iter);
                    goto cleanup;
                }

                class = cJSON_GetObjectItemCaseSensitive(module_cfg, "class");
                if (cJSON_IsNumber(class)) {
                    pmodule_cfg[iter].dev_class = class->valueint;
                } else {
                    log_error(
                        "Err: PARSER:: Failed to parse ModuleCfg[%d].dev_class",
                        iter);
                    goto cleanup;
                }

                sysfs =
                    cJSON_GetObjectItemCaseSensitive(module_cfg, "devsysfs");
                if (cJSON_IsString(sysfs)) {
                    memset(pmodule_cfg[iter].sysfile, '\0', 64);
                    memcpy(pmodule_cfg[iter].sysfile, sysfs->valuestring,
                           strlen(sysfs->valuestring));
                } else {
                    log_error(
                        "Err: PARSER:: Failed to parse ModuleCfg[%d].sysfs",
                        iter);
                    goto cleanup;
                }

                //TODO:: Make this one to read a generic device type.
                device =
                    cJSON_GetObjectItemCaseSensitive(module_cfg, "dev_hwattrs");
                if (cJSON_IsObject(device)) {
                    /* Return pointer to the device config.*/
                    void *pdev = parse_schema_devices(
                        device, pmodule_cfg[iter].dev_class);
                    if (pdev) {
                        pmodule_cfg[iter].cfg = pdev;
                    } else {
                        log_error(
                            "Err: PARSER:: Failed to parse ModuleCfg[%d].cfg",
                            iter);
                        goto cleanup;
                    }
                } else {
                    pmodule_cfg[iter].cfg = NULL;
                    log_warn(
                        "PARSER:: HW attributes are unavailable for parsing in ModuleCfg[%d].cfg ()",
                        iter);
                }

                iter++;
            }
            if (iter == count) {
                log_debug(
                    "PARSER:: All %d device info read from the manufacturing data's Module config.",
                    count);
            } else {
                log_error(
                    "PARSER::Only %d device info read from the manufacturing data's Module config expected were %d.",
                    iter, count);
            }
        }
    } else {
        log_error(
            "Err(%d):PARSER: Unexpected JSON object found instead of Module Config.",
            ERR_UBSP_UNEXPECTED_JSON_OBJECT);
        goto cleanup;
    }
    return pmodule_cfg;
cleanup:
    UBSP_FREE(pmodule_cfg);
    return pmodule_cfg;
}

void *parse_schema_generic_file_data(const cJSON *schema, uint16_t *size,
                                     char *str) {
    const cJSON *gen_data = NULL;
    char *gen_filename = NULL;
    char *generic_data = NULL;
    int ret = 0;
    gen_filename = malloc(sizeof(char) * 64);
    if (gen_filename) {
        gen_data = cJSON_GetObjectItemCaseSensitive(schema, str);
        if (cJSON_IsString(gen_data)) {
            memset(gen_filename, '\0', 64);
            memcpy(gen_filename, gen_data->valuestring,
                   strlen(gen_data->valuestring));
            int rsize = read_mfg_data_size(gen_filename);
            if (rsize <= 0) {
                log_error("Err :PARSER: No data or file %s not found.",
                          gen_filename);
                goto cleanup;
            }
            generic_data = malloc(sizeof(char) * rsize);
            if (generic_data) {
                memset(generic_data, 0, sizeof(char) * rsize);
                ret = read_mfg_data(gen_filename, generic_data, rsize);
                if (ret == rsize) {
                    log_debug("PARSER:: File %s read data of %d bytes.",
                              gen_filename, rsize);
                } else {
                    log_error(
                        "Err: PARSER:: File %s read data of %d bytes expected %d bytes.",
                        gen_filename, ret, rsize);
                }
                *size = ret;
            } else {
                log_error(
                    "Err(%d):PARSER: Memory exhausted while reading data from file %s.",
                    ERR_UBSP_MEMORY_EXHAUSTED, gen_filename);
                goto cleanup;
            }
        } else {
            log_error("Err: PARSER:: Failed to parse generic data file name.");
            goto cleanup;
        }
    } else {
        log_error(
            "Err(%d):PARSER: Memory exhausted while reading generic data file name.",
            ERR_UBSP_MEMORY_EXHAUSTED);
    }
cleanup:
    UBSP_FREE(gen_filename);
    return generic_data;
}

int parse_schema_payload(const cJSON *schema, UKDB **pukdb, uint16_t id,
                         int iter) {
    int ret = 0;
    switch (id) {
    case FIELDID_UNIT_INFO: {
        UnitInfo *punit_info = parse_schema_unit_info(schema);
        if (punit_info) {
            memcpy(&(*pukdb)->unitinfo, punit_info, sizeof(UnitInfo));
            UBSP_FREE(punit_info);
        } else {
            ret = -1;
            goto cleanup;
        }
        break;
    }
    case FIELDID_UNIT_CONFIG: {
        uint16_t module_count = (*pukdb)->unitinfo.mod_count;
        UnitCfg *punit_cfg = parse_schema_unit_config(schema, module_count);
        if (punit_cfg) {
            (*pukdb)->unitcfg = punit_cfg;
        } else {
            ret = -1;
            goto cleanup;
        }
        break;
    }
    case FIELDID_MODULE_INFO: {
        ModuleInfo *pmodule_info = parse_schema_module_info(schema);
        if (pmodule_info) {
            memcpy(&(*pukdb)->modinfo, pmodule_info, sizeof(ModuleInfo));
            UBSP_FREE(pmodule_info);
        } else {
            ret = -1;
            goto cleanup;
        }
        break;
    }
    case FIELDID_MODULE_CONFIG: {
        uint16_t device_count = (*pukdb)->modinfo.dev_count;
        ModuleCfg *pmodule_cfg =
            parse_schema_module_config(schema, device_count);
        if (pmodule_cfg) {
            (*pukdb)->modcfg = pmodule_cfg;
        } else {
            ret = -1;
            goto cleanup;
        }
        break;
    }
    case FIELDID_FACT_CONFIG: {
        uint16_t size = 0;
        int exp_size = (*pukdb)->indextable[iter].payload_size;
        char *factcfg =
            parse_schema_generic_file_data(schema, &size, "factory_config");
        if (factcfg) {
            if (size != exp_size) {
                log_warn(
                    "PARSER:: Size read for Field id 0x%x is %d bytes "
                    "and size mentioned in index table [%d] is 0x%d bytes.",
                    size, iter, exp_size, size);
                log_debug(
                    "Parser:: Updating index table [%d] size to %d bytes.",
                    iter, size);
                (*pukdb)->indextable[iter].payload_size = size;
            }
            (*pukdb)->factcfg = factcfg;
        } else {
            ret = -1;
            goto cleanup;
        }
        break;
    }
    case FIELDID_USER_CONFIG: {
        uint16_t size = 0;
        int exp_size = (*pukdb)->indextable[iter].payload_size;
        char *usercfg =
            parse_schema_generic_file_data(schema, &size, "user_config");
        if (usercfg) {
            if (size != exp_size) {
                log_warn(
                    "PARSER:: Size read for Field id 0x%x is %d bytes "
                    "and size mentioned in index table [%d] is 0x%d bytes.",
                    size, iter, exp_size, size);
                log_debug(
                    "Parser:: Updating index table [%d] size to %d bytes.",
                    iter, size);
                (*pukdb)->indextable[iter].payload_size = size;
            }
            (*pukdb)->usercfg = usercfg;
        } else {
            ret = -1;
            goto cleanup;
        }
        break;
    }
    case FIELDID_FACT_CALIB: {
        uint16_t size = 0;
        int exp_size = (*pukdb)->indextable[iter].payload_size;
        char *factcalib =
            parse_schema_generic_file_data(schema, &size, "fact_calibaration");
        if (factcalib) {
            if (size != exp_size) {
                log_warn(
                    "PARSER:: Size read for Field id 0x%x is %d bytes "
                    "and size mentioned in index table [%d] is 0x%d bytes.",
                    size, iter, exp_size, size);
                log_debug(
                    "Parser:: Updating index table [%d] size to %d bytes.",
                    iter, size);
                (*pukdb)->indextable[iter].payload_size = size;
            }
            (*pukdb)->factcalib = factcalib;
        } else {
            ret = -1;
            goto cleanup;
        }
        break;
    }
    case FIELDID_USER_CALIB: {
        uint16_t size = 0;
        int exp_size = (*pukdb)->indextable[iter].payload_size;
        char *usercalib =
            parse_schema_generic_file_data(schema, &size, "user_calibration");
        if (usercalib) {
            if (size != exp_size) {
                log_warn(
                    "PARSER:: Size read for Field id 0x%x is %d bytes "
                    "and size mentioned in index table [%d] is 0x%d bytes.",
                    size, iter, exp_size, size);
                log_debug(
                    "Parser:: Updating index table [%d] size to %d bytes.",
                    iter, size);
                (*pukdb)->indextable[iter].payload_size = size;
            }
            (*pukdb)->usercalib = usercalib;
        } else {
            ret = -1;
            goto cleanup;
        }
        break;
    }
    case FIELDID_BS_CERTS: {
        uint16_t size = 0;
        int exp_size = (*pukdb)->indextable[iter].payload_size;
        char *bscerts =
            parse_schema_generic_file_data(schema, &size, "boot_strap_certs");
        if (bscerts) {
            if (size != exp_size) {
                log_warn(
                    "PARSER:: Size read for Field id 0x%x is %d bytes "
                    "and size mentioned in index table [%d] is 0x%d bytes.",
                    size, iter, exp_size, size);
                log_debug(
                    "Parser:: Updating index table [%d] size to %d bytes.",
                    iter, size);
                (*pukdb)->indextable[iter].payload_size = size;
            }
            (*pukdb)->bscerts = bscerts;
        } else {
            ret = -1;
            goto cleanup;
        }
        break;
    }
    case FIELDID_LWM2M_CERTS: {
        uint16_t size = 0;
        int exp_size = (*pukdb)->indextable[iter].payload_size;
        char *lwm2mcerts =
            parse_schema_generic_file_data(schema, &size, "lwm2m2_certs");
        if (lwm2mcerts) {
            if (size != exp_size) {
                log_warn(
                    "PARSER:: Size read for Field id 0x%x is %d bytes "
                    "and size mentioned in index table [%d] is 0x%d bytes.",
                    size, iter, exp_size, size);
                log_debug(
                    "Parser:: Updating index table [%d] size to %d bytes.",
                    iter, size);
                (*pukdb)->indextable[iter].payload_size = size;
            }
            (*pukdb)->lwm2mcerts = lwm2mcerts;
        } else {
            ret = -1;
            goto cleanup;
        }
        break;
    }
    default: {
        ret = ERR_UBSP_INVALID_FIELD;
        log_error("Err(%d): Invalid Field id supplied by Index entry.", ret);
    }
    }
cleanup:
    return ret;
}

int parse_mfg_schema(const char *mfgdata, uint8_t idx) {
    cJSON *schema = NULL;
    const cJSON *header = NULL;
    const cJSON *index_table = NULL;
    const cJSON *unit_info = NULL;
    const cJSON *unit_config = NULL;
    const cJSON *module_config = NULL;
    const cJSON *factory_config = NULL;
    const cJSON *user_config = NULL;
    const cJSON *fact_calibaration = NULL;
    const cJSON *user_calibration = NULL;
    const cJSON *boot_strap_certs = NULL;
    const cJSON *lwm2m2_certs = NULL;
    int ret = 0;
    UKDB *pukdb = NULL;
    g_pukdb[idx] = malloc(sizeof(UKDB));
    if (g_pukdb[idx]) {
        pukdb = g_pukdb[idx];
        memset(pukdb, '\0', sizeof(UKDB));
        schema = cJSON_Parse(mfgdata);
        if (schema == NULL) {
            const char *error_ptr = cJSON_GetErrorPtr();
            if (error_ptr != NULL) {
                fprintf(stderr, "Error before: %s\n", error_ptr);
            }
            ret = ERR_UBSP_JSON_PARSER;
            goto cleanup;
        }

        /* Header */
        header = cJSON_GetObjectItemCaseSensitive(schema, "header");
        if (cJSON_IsObject(header)) {
            UKDBHeader *pheader = parse_schema_header(header);
            if (pheader) {
                memcpy(&pukdb->header, pheader, sizeof(UKDBHeader));
                UBSP_FREE(pheader);
            } else {
                ret = ERR_UBSP_INVALID_JSON_OBJECT;
                goto cleanup;
            }
        } else {
            ret = -1;
        }

        /* Index Table */
        index_table = cJSON_GetObjectItemCaseSensitive(schema, "index_table");
        if (cJSON_IsArray(index_table)) {
            UKDBIdxTuple *pindex_table = parse_schema_index_table(
                index_table, pukdb->header.idx_cur_tpl);
            if (pindex_table) {
                /* Free me once done*/
                pukdb->indextable = pindex_table;
            } else {
                ret = ERR_UBSP_INVALID_JSON_OBJECT;
                goto cleanup;
            }
        } else {
            ret = -1;
        }

        for (int iter = 0; iter < pukdb->header.idx_cur_tpl; iter++) {
            uint16_t id = pukdb->indextable[iter].fieldid;
            ret = parse_schema_payload(schema, &pukdb, id, iter);
            if (ret) {
                log_error(
                    "Err(%d): PARSER:: Failed parsing for Field Id 0x%x from mfg data.",
                    ret, id);
                goto cleanup;
            } else {
                log_debug(
                    "PARSER:: Parsing for Field Id 0x%x from mfg data completed.",
                    id);
            }
        }
    }
cleanup:
    cJSON_Delete(schema);
    if (ret) {
        parser_free_ukdb_mfg_data(&pukdb);
    }
    return ret;
}

UKDB *parser_get_mfg_data_by_uuid(char *puuid) {
    int ret = 0;
    UKDB *db = NULL;
    /*Default section*/
    if ((!puuid) || !strcmp(puuid, "")) {
        db = g_pukdb[0];
        log_trace(
            "PARSER:: MFG Data set to the Module UUID %s with %d entries in Index Table.",
            db->modinfo.uuid, db->header.idx_cur_tpl);
    } else {
    	/* Searching for Module MFG data index.*/
    	for (uint8_t iter = 0; iter < MAX_JSON_SCHEMA; iter++) {
    		if (g_pukdb[iter]) {
    			if (!strcmp(puuid, g_pukdb[iter]->modinfo.uuid)) {
    				db = g_pukdb[iter];
    				log_trace(
    						"PARSER:: MFG Data set to the Module UUID %s with %d entries in Index Table.",
							db->modinfo.uuid, db->header.idx_cur_tpl);
    				break;
    			}
    		}
        }
    }
    return db;
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
            log_debug("PARSER:: Starting the parsing of %s.", fname);
            off_t size = read_mfg_data_size(fname);
            char *schemabuff = malloc((sizeof(char) * size) + 1);
            if (schemabuff) {
                memset(schemabuff, '\0', (sizeof(char) * size) + 1);
                ret = read_mfg_data(fname, schemabuff, size);
                if (ret == size) {
                    log_debug(
                        "PARSER:: File %s read manufacturing data of %d bytes.",
                        fname, size);
                    ret = parse_mfg_schema(schemabuff, iter);
                    if (ret) {
                        log_error("Err(%d): PARSER:: Parsing failed for %s.",
                                  ret, fname);
                    } else {
                        log_debug("PARSER: Parsing completed for %s.", fname);
                    }
                }
            }
            UBSP_FREE(schemabuff);
        }
    }
    return ret;
}

void parser_schema_exit() {
    for (int iter = 0; iter < MAX_JSON_SCHEMA; iter++) {
        if (g_pukdb[iter]) {
            parser_free_ukdb_mfg_data(&g_pukdb[iter]);
        }
    }
}

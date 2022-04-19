/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "prop_parser.h"

#include "errorcode.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"

PropertyMap *gPropMap[MAX_JSON_DEVICE] = { '\0' };

static int read_mfg_data(char *fname, char *buff, off_t size) {
    int readBytes = 0;
    /* Create input file descriptor */
    int fd = usys_open(fname, O_RDONLY, 0644);
    if (fd == -1) {
        usys_log_error("Failed to open file. Error: %s", usys_error(errno));
        return fd;
    }
    off_t off = usys_lseek(fd, 0, SEEK_SET);
    readBytes = usys_read(fd, buff, size);
    return readBytes;
}

static int read_mfg_data_size(char *fname) {
    int readBytes = 0;
    /* Create input file descriptor */
    int fd = usys_open(fname, O_RDONLY, 0644);
    if (fd == -1) {
        usys_log_error("Failed to open file. Error: %s", usys_error(errno));
        return fd;
    }

    off_t off = usys_lseek(fd, 0L, SEEK_END);
    return off;
}

void parser_free_prop(Property *prop, uint16_t count) {
    for (uint16_t itr = 0; itr < count; itr++) {
        usys_free(prop[itr].depProp);
        prop[itr].depProp = NULL;
    }
    usys_free(prop);
    prop = NULL;
}

void parser_free_pmap(PropertyMap **pMap) {
    if (*pMap) {
        if ((*pMap)->prop) {
            parser_free_prop((*pMap)->prop, (*pMap)->propCount);
        }
        usys_free(*pMap);
        *pMap = NULL;
    }
}

DepProperty *prop_parser_get_dependents(const JsonObj *jDepProp) {
    const JsonObj *jValProp = NULL;
    int ret = 0;
    DepProperty *dProp = usys_zmalloc(sizeof(DepProperty));
    if (dProp) {
        usys_memset(dProp, 0, sizeof(DepProperty));

        /* Current Value */
        int currIdx = 0;
        if (!parser_read_integer_object(jDepProp, JTAG_CURR_PROP_ID,
                                        &currIdx)) {
            goto cleanup;
        } else {
            dProp->currIdx = currIdx;
        }

        /* Limit value */
        int lmtIdx = 0;
        if (!parser_read_integer_object(jDepProp, JTAG_LIMIT_PROP_ID,
                                        &lmtIdx)) {
            goto cleanup;
        } else {
            dProp->lmtIdx = lmtIdx;
        }

        /* Alert Condition */
        char *alertCond;
        if (!parser_read_string_object(jDepProp, JTAG_ALERT_COND, &alertCond)) {
            goto cleanup;
        } else {
            dProp->cond = get_alert_cond(alertCond);
            usys_free(alertCond);
        }

    } else {
        ret = ERR_NODED_MEMORY_EXHAUSTED;
        usys_log_error(
            "Err(%d):PARSER: Memory exhausted while parsing DepProperty.", ret);
        goto cleanup;
    }

cleanup:
    if (ret) {
        usys_free(dProp);
        dProp = NULL;
    }
    return (dProp);
}

int prop_parse_table(const JsonObj *jPropTable, PropertyMap **pMap) {
    const JsonObj *jProp = NULL;
    const JsonObj *jDependent = NULL;
    int ret = 0;
    uint16_t iter = 0;
    Property *prop = NULL;

    uint16_t propCount = json_array_size(jPropTable);
    if (propCount <= 0) {
        goto cleanup;
    }
    prop = usys_zmalloc(sizeof(Property) * propCount);
    if (prop) {
        usys_memset(prop, '\0', sizeof(Property) * propCount);

        json_array_foreach(jPropTable, iter, jProp) {
            /* ID */
            int id = 0;
            if (!parser_read_integer_object(jProp, JTAG_ID, &id)) {
                goto cleanup;
            } else {
                prop[iter].id = (uint16_t)id;
            }

            /* Name */
            if (!parser_read_string_object_wrapper(jProp, JTAG_NAME,
                                                   prop[iter].name)) {
                goto cleanup;
            }

            /* Data Type */
            char *dataType = NULL;
            if (!parser_read_string_object(jProp, JTAG_DATA_TYPE, &dataType)) {
                goto cleanup;
            } else {
                prop[iter].dataType = get_prop_data_type(dataType);
                usys_free(dataType);
            }

            /* Permission */
            char *perm= NULL;
            if (!parser_read_string_object(jProp, JTAG_PERMISSION, &perm)) {
                goto cleanup;
            } else {
                prop[iter].perm = get_prop_perm(perm);
                usys_free(perm);
            }

            /* Availability */
            char *avail= NULL;
            if (!parser_read_string_object(jProp, JTAG_AVAILABILITY, &avail)) {
                goto cleanup;
            } else {
                prop[iter].available = get_prop_avail(avail);
                usys_free(avail);
            }

            /* Property Type */
            char *propType= NULL;
            if (!parser_read_string_object(jProp, JTAG_PROPERTY_TYPE,
                                           &propType)) {
                goto cleanup;
            } else {
                prop[iter].propType = get_prop_type(propType);
                usys_free(propType);
            }

            /* Units */
            if (!parser_read_string_object_wrapper(jProp, JTAG_UNITS,
                                                   prop[iter].units)) {
                goto cleanup;
            }

            /* Sysfile */
            /* SysFs */
            if (!parser_read_string_object_wrapper(jProp, JTAG_SYS_FS_FILE,
                                                   prop[iter].sysFname)) {
                goto cleanup;
            }

            /* Dependent Property */
            jDependent = json_object_get(jProp, JTAG_DEPENDENT);
            if (json_is_object(jDependent)) {
                DepProperty *dProp = prop_parser_get_dependents(jDependent);
                if (dProp) {
                    prop[iter].depProp = dProp;
                } else {
                    prop[iter].depProp = NULL;
                    usys_log_error(
                        "Err(%d): PARSER:: Failed to parse Property[%d].depProp for"
                        " %s device",
                        iter, (*pMap)->name);
                }
            } else {
                prop[iter].depProp = NULL;
            }
        }
        (*pMap)->propCount = iter;
        (*pMap)->prop = prop;
        usys_log_trace("PARSER:: %d property for device %s found in json.",
                       iter, (*pMap)->name);
    } else {
        ret = ERR_NODED_MEMORY_EXHAUSTED;
        usys_log_error(
            "Err(%d):PARSER:: Memory exhausted while parsing property table for %s device.",
            ret, (*pMap)->name);
        goto cleanup;
    }

cleanup:
    if (ret) {
        parser_free_prop(prop, iter);
    }
    return ret;
}

int prop_parse_dev(const JsonObj *jDevices) {
    const JsonObj *jDev = NULL;
    const JsonObj *jDevName = NULL;
    const JsonObj *jDevVersion = NULL;
    const JsonObj *jDevTable = NULL;
    PropertyMap *pMap = NULL;
    int iter = 0;
    int ret = 0;
    int count = 0;

    json_array_foreach(jDevices, iter, jDev) {
        pMap = usys_zmalloc(sizeof(PropertyMap));
        if (pMap) {
            gPropMap[count] = pMap;

            /* Name */
            if (!parser_read_string_object_wrapper(jDev, JTAG_NAME,
                                                   pMap->name)) {
                goto cleanup;
            }

            /* Version */
            jDevVersion = json_object_get(jDev, JTAG_VERSION);
            Version *pVersion = parse_version(jDevVersion);
            if (pVersion) {
                usys_log_info(
                    "Parser:: Device %s is using json property version %d.%d.",
                    pMap->name, pVersion->major, pVersion->minor);
                usys_memcpy(&pMap->ver, pVersion, sizeof(Version));
                usys_free(pVersion);
            } else {
                usys_log_error(
                    "Err(%d): PARSER:: Failed to parse Device  %s property version.",
                    ret, pMap->name);
                goto cleanup;
            }

            /* Property Table */
            jDevTable = json_object_get(jDev, JTAG_PROPERTY_TABLE);
            ret = prop_parse_table(jDevTable, &pMap);
            if (!ret) {
                usys_log_trace(
                    "Parser:: Device %s json property table parsing completed "
                    "with %d properties.",
                    pMap->name, pMap->propCount);
            } else {
                usys_log_error(
                    "Err(%d): Parser:: Device %s json property table parsing failed.",
                    pMap->name);
                goto cleanup;
            }

            count++;

        } else {
            ret = ERR_NODED_MEMORY_EXHAUSTED;
            usys_log_error(
                "Err(%d):PARSER:: Memory exhausted while parsing Device table.",
                ret);
            goto cleanup;
        }
        usys_log_trace("PARSER:: %d device read from json.", count);
    }

cleanup:
    if (ret) {
        parser_free_pmap(&pMap);
    }
    return ret;
}

int prop_parse_json(char *propBuff) {
    int ret = 0;
    JsonObj *jProp = NULL;
    JsonErrObj *jErr = NULL;
    const JsonObj *jName = NULL;
    const JsonObj *jDevice = NULL;
    Property *pProp = NULL;
    char name[NAME_LENGTH] = { '\0' };

    jProp = json_loads(propBuff, JSON_DECODE_ANY, jErr);
    if (!jProp) {
        parser_error(jErr, "Failed to parse property data");
        ret = ERR_NODED_JSON_PARSER;
        goto cleanup;
    }

    /* name */
    if (!parser_read_string_object_wrapper(jProp, JTAG_NAME, name)) {
        goto cleanup;
    }

    /* Devices */
    jDevice = json_object_get(jProp, JTAG_DEVICE);
    if (json_is_array(jDevice)) {
        ret = prop_parse_dev(jDevice);
        if (!ret) {
            usys_log_trace("PARSER: Property table parsed for Device %s.",
                           name);
        } else {
            usys_log_error(
                "Err(%d): PARSER: Property table not parsed for Device %s.",
                ret, name);
            goto cleanup;
        }
    } else {
        ret = ERR_NODED_INVALID_JSON_OBJECT;
        usys_log_error("Err(%d): PARSER: %s not parsed from JSON.", ret, name);
    }

cleanup:
    json_decref(jProp);
    return ret;
}

int prop_parser_get_count(char *name) {
    int count = -1;
    for (uint8_t iter = 0; iter < MAX_JSON_DEVICE; iter++) {
        if (gPropMap[iter]) {
            if (!usys_strcmp(name, gPropMap[iter]->name)) {
                count = gPropMap[iter]->propCount;
                break;
            }
        }
    }
    return count;
}

Version *prop_parser_get_table_version(char *name) {
    Version *ver = NULL;
    for (uint8_t iter = 0; iter < MAX_JSON_DEVICE; iter++) {
        if (gPropMap[iter]) {
            if (!usys_strcmp(name, gPropMap[iter]->name)) {
                ver = &gPropMap[iter]->ver;
                break;
            }
        }
    }
    return ver;
}

Property *prop_parser_get_table(char *name) {
    Property *table = NULL;
    for (uint8_t iter = 0; iter < MAX_JSON_DEVICE; iter++) {
        if (gPropMap[iter]) {
            if (!usys_strcmp(name, gPropMap[iter]->name)) {
                table = gPropMap[iter]->prop;
                break;
            }
        }
    }
    return table;
}

int prop_parser_init(char *ip) {
    int ret = 0;
    char *fname = NULL;
    if (ip) {
        fname = ip;
        usys_log_trace("PARSER:: Starting the parsing of %s.", fname);
        off_t size = read_mfg_data_size(fname);
        char *schemabuff = usys_zmalloc((sizeof(char) * size) + 1);
        if (schemabuff) {
            ret = read_mfg_data(fname, schemabuff, size);
            if (ret == size) {
                usys_log_trace(
                    "PARSER:: File %s read manufacturing data of %d bytes.",
                    fname, size);
                ret = prop_parse_json(schemabuff);
                if (ret) {
                    usys_log_error("Err(%d): PARSER:: Parsing failed for %s.",
                                   ret, fname);
                } else {
                    usys_log_trace("PARSER: Parsing completed for %s.", fname);
                }
            }
        }
        usys_free(schemabuff);
        schemabuff = NULL;
    }

    return ret;
}

void prop_parser_exit() {
    for (int iter = 0; iter < MAX_JSON_DEVICE; iter++) {
        if (gPropMap[iter]) {
            parser_free_pmap(&gPropMap[iter]);
        }
    }
}

/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "utils/pparser.h"

#include "headers/ubsp/devices.h"
#include "headers/errorcode.h"
#include "headers/ubsp/property.h"
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

PropertyMap *g_pmap[MAX_JSON_DEVICE] = { '\0' };

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

void parser_free_prop(Property *prop, uint16_t count) {
    for (uint16_t itr = 0; itr < count; itr++) {
        UBSP_FREE(prop[itr].dep_prop);
    }
    UBSP_FREE(prop);
}

void parser_free_pmap(PropertyMap **pmap) {
    if (*pmap) {
        if ((*pmap)->prop) {
            parser_free_prop((*pmap)->prop, (*pmap)->prop_count);
        }
        UBSP_FREE(*pmap);
    }
}

Version *parse_dev_prop_version(const cJSON *version) {
    const cJSON *major = NULL;
    const cJSON *minor = NULL;
    int ret = 0;
    Version *pversion = malloc(sizeof(Version));
    if (pversion) {
        /* Major */
        major = cJSON_GetObjectItemCaseSensitive(version, "major");
        if (cJSON_IsNumber(major)) {
            pversion->major = version->valueint;
        } else {
            ret = ERR_UBSP_INVALID_JSON_OBJECT;
            log_error("Err(%d):PARSER:: Parsing failed for Version.major.",
                      ret);
        }

        /* Minor */
        minor = cJSON_GetObjectItemCaseSensitive(version, "minor");
        if (cJSON_IsNumber(minor)) {
            pversion->minor = version->valueint;
        } else {
            ret = ERR_UBSP_INVALID_JSON_OBJECT;
            log_error("Err(%d):PARSER:: Parsing failed for Version.minor.",
                      ret);
        }

    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_error("Err(%d):PARSER: Memory exhausted while parsing Version.",
                  ret);
        goto cleanup;
    }

cleanup:
    if (ret) {
        UBSP_FREE(pversion);
    }
    return pversion;
}

DepProperty *parser_depproperty(const cJSON *jdprop) {
    const cJSON *jval_prop = NULL;
    const cJSON *jlimit_id = NULL;
    const cJSON *jalert_cond = NULL;
    int ret = 0;
    DepProperty *dprop = malloc(sizeof(DepProperty));
    if (dprop) {
        memset(dprop, 0, sizeof(DepProperty));
        /* Current Value*/
        jval_prop =
            cJSON_GetObjectItemCaseSensitive(jdprop, "current_val_property");
        if (cJSON_IsNumber(jval_prop)) {
            dprop->curr_idx = jval_prop->valueint;
        } else {
            ret = ERR_UBSP_INVALID_JSON_OBJECT;
            log_error(
                "Err(%d):PARSER:: Parsing failed for DepProperty.curr_idx.",
                ret);
            goto cleanup;
        }

        /* Limit value*/
        jlimit_id =
            cJSON_GetObjectItemCaseSensitive(jdprop, "limit_val_property");
        if (cJSON_IsNumber(jlimit_id)) {
            dprop->lmt_idx = jlimit_id->valueint;
        } else {
            ret = ERR_UBSP_INVALID_JSON_OBJECT;
            log_error(
                "Err(%d):PARSER:: Parsing failed for DepProperty.lmt_idx.",
                ret);
            goto cleanup;
        }

        /* Alert Condition */
        jalert_cond =
            cJSON_GetObjectItemCaseSensitive(jdprop, "alert_condition");
        if (cJSON_IsString(jalert_cond)) {
            dprop->cond = get_alert_cond(jalert_cond->valuestring);
        } else {
            ret = ERR_UBSP_INVALID_JSON_OBJECT;
            log_error(
                "Err(%d):PARSER:: Parsing failed for DepProperty.AlertCondition.",
                ret);
            goto cleanup;
        }

    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_error("Err(%d):PARSER: Memory exhausted while parsing DepProperty.",
                  ret);
        goto cleanup;
    }

cleanup:
    if (ret) {
        UBSP_FREE(dprop);
    }
    return (dprop);
}

int parse_property_table(const cJSON *jprop_table, PropertyMap **pmap) {
    const cJSON *jprop = NULL;
    const cJSON *jid = NULL;
    const cJSON *jname = NULL;
    const cJSON *jdata_type = NULL;
    const cJSON *jperm = NULL;
    const cJSON *javailable = NULL;
    const cJSON *jprop_type = NULL;
    const cJSON *junits = NULL;
    const cJSON *jsysfile = NULL;
    const cJSON *jdependent = NULL;
    int ret = 0;
    uint16_t iter = 0;
    Property *prop = NULL;
    uint16_t prop_count = cJSON_GetArraySize(jprop_table);
    if (prop_count <= 0) {
        goto cleanup;
    }
    prop = malloc(sizeof(Property) * prop_count);
    if (prop) {
        memset(prop, '\0', sizeof(Property) * prop_count);
        cJSON_ArrayForEach(jprop, jprop_table) {
            /* ID */
            jid = cJSON_GetObjectItemCaseSensitive(jprop, "id");
            if (cJSON_IsNumber(jid)) {
                prop[iter].id = jid->valueint;
            } else {
                ret = ERR_UBSP_INVALID_JSON_OBJECT;
                log_error(
                    "Err(%d): PARSER:: Failed to parse Property[%d].id for"
                    " %s device",
                    ret, iter, (*pmap)->name);
                goto cleanup;
            }

            /* Name */
            jname = cJSON_GetObjectItemCaseSensitive(jprop, "name");
            if (cJSON_IsString(jname)) {
                strcpy(prop[iter].name, jname->valuestring);
            } else {
                ret = ERR_UBSP_INVALID_JSON_OBJECT;
                log_error(
                    "Err(%d): PARSER:: Failed to parse Property[%d].name for"
                    " %s device",
                    iter, (*pmap)->name);
                goto cleanup;
            }

            /*Data Type*/
            jdata_type = cJSON_GetObjectItemCaseSensitive(jprop, "data_type");
            if (cJSON_IsString(jdata_type)) {
                prop[iter].data_type =
                    get_prop_datatype(jdata_type->valuestring);
            } else {
                ret = ERR_UBSP_INVALID_JSON_OBJECT;
                log_error(
                    "Err(%d): PARSER:: Failed to parse Property[%d].data_type for"
                    " %s device",
                    iter, (*pmap)->name);
                goto cleanup;
            }

            /* Permission */
            jperm = cJSON_GetObjectItemCaseSensitive(jprop, "perm");
            if (cJSON_IsString(jname)) {
                prop[iter].perm = get_prop_perm(jperm->valuestring);
            } else {
                ret = ERR_UBSP_INVALID_JSON_OBJECT;
                log_error(
                    "Err(%d): PARSER:: Failed to parse Property[%d].perm for"
                    " %s device",
                    iter, (*pmap)->name);
                goto cleanup;
            }

            /* Availability */
            javailable = cJSON_GetObjectItemCaseSensitive(jprop, "available");
            if (cJSON_IsString(javailable)) {
                prop[iter].available = get_prop_avail(javailable->valuestring);
            } else {
                ret = ERR_UBSP_INVALID_JSON_OBJECT;
                log_error(
                    "Err(%d): PARSER:: Failed to parse Property[%d].available for"
                    " %s device",
                    iter, (*pmap)->name);
                goto cleanup;
            }

            /* Property Type */
            jprop_type = cJSON_GetObjectItemCaseSensitive(jprop, "prop_type");
            if (cJSON_IsString(jprop_type)) {
                prop[iter].prop_type = get_prop_type(jprop_type->valuestring);
            } else {
                ret = ERR_UBSP_INVALID_JSON_OBJECT;
                log_error(
                    "Err(%d): PARSER:: Failed to parse Property[%d].prop_type for"
                    " %s device",
                    iter, (*pmap)->name);
                goto cleanup;
            }

            /* Units */
            junits = cJSON_GetObjectItemCaseSensitive(jprop, "units");
            if (cJSON_IsString(junits)) {
                strcpy(prop[iter].units, junits->valuestring);
            } else {
                ret = ERR_UBSP_INVALID_JSON_OBJECT;
                log_error(
                    "Err(%d): PARSER:: Failed to parse Property[%d].units for"
                    " %s device",
                    iter, (*pmap)->name);
                goto cleanup;
            }

            /* Sysfile */
            jsysfile = cJSON_GetObjectItemCaseSensitive(jprop, "sysfsfile");
            if (cJSON_IsString(jsysfile)) {
                strcpy(prop[iter].sysfname, jsysfile->valuestring);
            } else {
                ret = ERR_UBSP_INVALID_JSON_OBJECT;
                log_error(
                    "Err(%d): PARSER:: Failed to parse Property[%d].sysfsfile for"
                    " %s device",
                    iter, (*pmap)->name);
                goto cleanup;
            }

            /* Dependent Property */
            jdependent = cJSON_GetObjectItemCaseSensitive(jprop, "dependent");
            if (cJSON_IsObject(jdependent)) {
                DepProperty *dprop = parser_depproperty(jdependent);
                if (dprop) {
                    prop[iter].dep_prop = dprop;
                } else {
                    prop[iter].dep_prop = NULL;
                    log_error(
                        "Err(%d): PARSER:: Failed to parse Property[%d].dep_prop for"
                        " %s device",
                        iter, (*pmap)->name);
                }
            } else {
                prop[iter].dep_prop = NULL;
            }

            iter++;
        }
        (*pmap)->prop_count = iter;
        (*pmap)->prop = prop;
        log_trace("PARSER:: %d property for device %s found in json.", iter,
                  (*pmap)->name);
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_error(
            "Err(%d):PARSER:: Memory exhausted while parsing property table for %s device.",
            ret, (*pmap)->name);
        goto cleanup;
    }

cleanup:
    if (ret) {
        parser_free_prop(prop, iter);
    }
    return ret;
}

int parse_dev_prop(const cJSON *jdev_prop) {
    const cJSON *jdev = NULL;
    const cJSON *jdev_name = NULL;
    const cJSON *jdev_version = NULL;
    const cJSON *jdev_prop_table = NULL;
    PropertyMap *pmap = NULL;
    int ret = 0;
    int count = 0;

    cJSON_ArrayForEach(jdev, jdev_prop) {
        if (count >= MAX_JSON_DEVICE) {
            log_error(
                "Err(%d): PARSER:: More than expected devices(%d) found in property json.",
                MAX_JSON_DEVICE);
        }

        pmap = malloc(sizeof(PropertyMap));
        if (pmap) {
            g_pmap[count] = pmap;
            memset(pmap, 0, sizeof(PropertyMap));
            /* name */
            jdev_name = cJSON_GetObjectItemCaseSensitive(jdev, "name");
            if (cJSON_IsString(jdev_name)) {
                memcpy(pmap->name, jdev_name->valuestring,
                       strlen(jdev_name->valuestring));
            } else {
                ret = ERR_UBSP_INVALID_JSON_OBJECT;
                log_error(
                    "Err(%d): PARSER:: Failed to parse Property Device name.",
                    ret);
                goto cleanup;
            }

            /* Version */
            jdev_version = cJSON_GetObjectItemCaseSensitive(jdev, "version");
            Version *pversion = parse_dev_prop_version(jdev_version);
            if (pversion) {
                log_info(
                    "Parser:: Device %s is using json property version %d.%d.",
                    pmap->name, pversion->major, pversion->minor);
                memcpy(&pmap->ver, pversion, sizeof(Version));
                UBSP_FREE(pversion);
            } else {
                log_error(
                    "Err(%d): PARSER:: Failed to parse Device  %s property version.",
                    ret, pmap->name);
                goto cleanup;
            }

            /* Property Table */
            jdev_prop_table =
                cJSON_GetObjectItemCaseSensitive(jdev, "property_table");
            ret = parse_property_table(jdev_prop_table, &pmap);
            if (!ret) {
                log_trace(
                    "Parser:: Device %s json property table parsing completed.",
                    pmap->name);
            } else {
                log_error(
                    "Err(%d): Parser:: Device %s json property table parsing failed.",
                    pmap->name);
                goto cleanup;
            }

            count++;

        } else {
            ret = ERR_UBSP_MEMORY_EXHAUSTED;
            log_error(
                "Err(%d):PARSER:: Memory exhausted while parsing Device table.",
                ret);
            goto cleanup;
        }
        log_trace("PARSER:: %d device read from json.", count);
    }

cleanup:
    if (ret) {
        parser_free_pmap(&pmap);
    }
    return ret;
}

int parse_property_json(char *prop_buff) {
    int ret = 0;
    cJSON *jprop = NULL;
    const cJSON *jname = NULL;
    const cJSON *jdevice = NULL;
    Property *pprop = NULL;
    char name[PROP_NAME_LENGTH] = { '\0' };
    jprop = cJSON_Parse(prop_buff);
    if (jprop == NULL) {
        const char *error_ptr = cJSON_GetErrorPtr();
        if (error_ptr != NULL) {
            log_error("Err: PARSER:: Error before: %s\n", error_ptr);
        }
        ret = ERR_UBSP_JSON_PARSER;
        goto cleanup;
    }

    /* name */
    jname = cJSON_GetObjectItemCaseSensitive(jprop, "name");
    if (cJSON_IsString(jname)) {
        memcpy(name, jname->valuestring, strlen(jname->valuestring));
    } else {
        ret = ERR_UBSP_INVALID_JSON_OBJECT;
        log_error("Err(%d): PARSER:: Failed to parse Property Json name.", ret);
        goto cleanup;
    }

    /* Devices */
    jdevice = cJSON_GetObjectItemCaseSensitive(jprop, "device");
    if (cJSON_IsArray(jdevice)) {
        ret = parse_dev_prop(jdevice);
        if (!ret) {
            log_trace("PARSER: Property table parsed for Device %s.", name);
        } else {
            log_error(
                "Err(%d): PARSER: Property table not parsed for Device %s.",
                ret, name);
            goto cleanup;
        }
    } else {
        ret = ERR_UBSP_INVALID_JSON_OBJECT;
        log_error("Err(%d): PARSER: %s not parsed from JSON.", ret, name);
    }

cleanup:
    cJSON_Delete(jprop);
    return ret;
}

int parser_get_property_count(char *name) {
    int count = -1;
    for (uint8_t iter = 0; iter < MAX_JSON_DEVICE; iter++) {
        if (g_pmap[iter]) {
            if (!strcmp(name, g_pmap[iter]->name)) {
                count = g_pmap[iter]->prop_count;
                break;
            }
        }
    }
    return count;
}

Version *parser_get_property_table_version(char *name) {
    Version *ver = NULL;
    for (uint8_t iter = 0; iter < MAX_JSON_DEVICE; iter++) {
        if (g_pmap[iter]) {
            if (!strcmp(name, g_pmap[iter]->name)) {
                ver = &g_pmap[iter]->ver;
                break;
            }
        }
    }
    return ver;
}

Property *parser_get_property_table(char *name) {
    Property *table = NULL;
    for (uint8_t iter = 0; iter < MAX_JSON_DEVICE; iter++) {
        if (g_pmap[iter]) {
            if (!strcmp(name, g_pmap[iter]->name)) {
                table = g_pmap[iter]->prop;
                break;
            }
        }
    }
    return table;
}

int parser_property_init(char *ip) {
    int ret = 0;
    char *fname = NULL;
    if (ip) {
        fname = ip;
        log_trace("PARSER:: Starting the parsing of %s.", fname);
        off_t size = read_mfg_data_size(fname);
        char *schemabuff = calloc(1, ((sizeof(char) * size) + 1));
        if (schemabuff) {
            memset(schemabuff, '\0', (sizeof(char) * size) + 1);
            ret = read_mfg_data(fname, schemabuff, size);
            if (ret == size) {
                log_trace(
                    "PARSER:: File %s read manufacturing data of %d bytes.",
                    fname, size);
                ret = parse_property_json(schemabuff);
                if (ret) {
                    log_error("Err(%d): PARSER:: Parsing failed for %s.", ret,
                              fname);
                } else {
                    log_trace("PARSER: Parsing completed for %s.", fname);
                }
            }
        }
        UBSP_FREE(schemabuff);
    }

    return ret;
}

void parser_property_exit() {
    for (int iter = 0; iter < MAX_JSON_DEVICE; iter++) {
        if (g_pmap[iter]) {
            parser_free_pmap(&g_pmap[iter]);
        }
    }
}

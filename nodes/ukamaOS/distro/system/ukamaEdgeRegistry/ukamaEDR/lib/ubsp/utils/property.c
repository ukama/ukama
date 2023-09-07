/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "headers/ubsp/property.h"
#include "inc/driverfxn.h"
#include "devdb/sysfs/drvrsysfs.h"
#include "headers/utils/log.h"
#include "utils/pparser.h"

#include <stdbool.h>
#include <string.h>

uint16_t get_sizeof(DataType type) {
    uint16_t size = 0;
    switch (type) {
    case TYPE_NULL: {
        size = 0;
        break;
    }
    case TYPE_CHAR: {
        size = sizeof(char);
        break;
    }
    case TYPE_BOOL: {
        size = sizeof(bool);
        break;
    }
    case TYPE_UINT8: {
        size = sizeof(uint8_t);
        break;
    }
    case TYPE_INT8: {
        size = sizeof(int8_t);
        break;
    }
    case TYPE_UINT16:
        size = sizeof(uint16_t);
        break;
    case TYPE_INT16: {
        size = sizeof(int16_t);
        break;
        break;
    }
    case TYPE_UINT32: {
        size = sizeof(uint32_t);
        break;
    }
    case TYPE_INT32: {
        size = sizeof(int32_t);
        break;
    }
    case TYPE_INT: {
        size = sizeof(int);
        break;
    }
    case TYPE_FLOAT:
        size = sizeof(float);
        break;
    case TYPE_ENUM: {
        size = sizeof(int);
        break;
    }
    case TYPE_DOUBLE: {
        size = sizeof(double);
        break;
    }
    case TYPE_STRING: {
        size = sizeof(char) * PROP_MAX_STR_LENGTH;
        break;
    }
    default: {
        size = 0;
    }
    }
    return size;
};

int validate_irq_limits(double cur, double lmt, int cond) {
    int ret = 0;
    switch (cond) {
    case STRICTLYLESSTHEN: {
        if (cur < lmt) {
            ret = 1;
        }
        break;
    }
    case LESSTHENEQUALTO: {
        if (cur <= lmt) {
            ret = 1;
        }
        break;
    }
    case STRICTLYGREATERTHEN: {
        if (cur > lmt) {
            ret = 1;
        }
        break;
    }
    case GREATERTHENEQUALTO: {
        if (cur >= lmt) {
            ret = 1;
        }
        break;
    }
    default: {
        ret = -1;
    }
    }
    return ret;
}

int get_alert_cond(char *cond) {
    int ret = -1;
    if (cond) {
        if (!strcmp(cond, "STRICTLYLESSTHEN")) {
            ret = STRICTLYLESSTHEN;
        } else if (!strcmp(cond, "LESSTHENEQUALTO")) {
            ret = LESSTHENEQUALTO;
        } else if (!strcmp(cond, "STRICTLYGREATERTHEN")) {
            ret = STRICTLYGREATERTHEN;
        } else if (!strcmp(cond, "GREATERTHENEQUALTO")) {
            ret = GREATERTHENEQUALTO;
        }
    }
    return ret;
}

int get_prop_perm(char *perm) {
    int ret = -1;
    if (perm) {
        if (!strcmp(perm, "PERM_EX")) {
            ret = PERM_EX;
        } else if (!strcmp(perm, "PERM_RD")) {
            ret = PERM_RD;
        } else if (!strcmp(perm, "PERM_WR")) {
            ret = PERM_WR;
        } else if (!strcmp(perm, "PERM_RW")) {
            ret = PERM_RW;
        } else if (!strcmp(perm, "PERM_RWE")) {
            ret = PERM_RWE;
        }
    }
    return ret;
}

int get_prop_type(char *type) {
    int ret = -1;
    if (type) {
        if (!strcmp(type, "PROP_TYPE_CONFIG")) {
            ret = PROP_TYPE_CONFIG;
        } else if (!strcmp(type, "PROP_TYPE_STATUS")) {
            ret = PROP_TYPE_STATUS;
        } else if (!strcmp(type, "PROP_TYPE_ALERT")) {
            ret = PROP_TYPE_ALERT;
        } else if (!strcmp(type, "PROP_TYPE_EXEC")) {
            ret = PROP_TYPE_EXEC;
        }
    }
    return ret;
}

int get_prop_avail(char *avail) {
    int ret = -1;
    if (avail) {
        if (!strcmp(avail, "PROP_NOTAVAIL")) {
            ret = PROP_NOTAVAIL;
        } else if (!strcmp(avail, "PROP_AVAIL")) {
            ret = PROP_AVAIL;
        }
    }
    return ret;
}

int get_prop_datatype(char *type) {
    int ret = -1;
    if (type) {
        if (!strcmp(type, "TYPE_NULL")) {
            ret = TYPE_NULL;
        } else if (!strcmp(type, "TYPE_CHAR")) {
            ret = TYPE_CHAR;
        } else if (!strcmp(type, "TYPE_BOOL")) {
            ret = TYPE_BOOL;
        } else if (!strcmp(type, "TYPE_UINT8")) {
            ret = TYPE_UINT8;
        } else if (!strcmp(type, "TYPE_UINT16")) {
            ret = TYPE_UINT16;
        } else if (!strcmp(type, "TYPE_UINT32")) {
            ret = TYPE_UINT32;
        } else if (!strcmp(type, "TYPE_INT8")) {
            ret = TYPE_INT8;
        } else if (!strcmp(type, "TYPE_INT16")) {
            ret = TYPE_INT16;
        } else if (!strcmp(type, "TYPE_INT32")) {
            ret = TYPE_INT32;
        } else if (!strcmp(type, "TYPE_INT")) {
            ret = TYPE_INT;
        } else if (!strcmp(type, "TYPE_FLOAT")) {
            ret = TYPE_FLOAT;
        } else if (!strcmp(type, "TYPE_DOUBLE")) {
            ret = TYPE_DOUBLE;
        } else if (!strcmp(type, "TYPE_ENUM")) {
            ret = TYPE_ENUM;
        } else if (!strcmp(type, "TYPE_STRING")) {
            ret = TYPE_STRING;
        }
    }
    return ret;
}

int get_property_count(char *dev) {
    return parser_get_property_count(dev);
}

Property *get_property_table(char *dev) {
    return parser_get_property_table(dev);
}

/* Need to extract last part of file name from path /tmp/sys/class/hwmon/hwmon0/5/se98_1/temp1_min_alarm
 * i.e temp1_min_alarm.
 */
char *get_sysfs_name(char *fpath) {
    char *tok;
    char *prev;
    tok = strtok(fpath, "/");
    while (tok != NULL) {
        prev = tok;
        tok = strtok(NULL, "/");
    }
    return prev;
}

void print_properties(Property *prop, uint16_t count) {
    log_trace(
        "********************************************************************************");
    log_trace(
        "************************ Property Table ****************************************");
    for (uint16_t iter = 0; iter < count; iter++) {
        if (prop[iter].available == PROP_NOTAVAIL) {
            continue;
        }
        log_trace(
            "********************************************************************************");
        log_trace("* C-struct ID [%d] JIndex [%d]", iter, prop[iter].id);
        log_trace("* Name                      : %s", prop[iter].name);
        log_trace("* Data Type                 : 0x%x", prop[iter].data_type);
        log_trace("* Permission                : 0x%x", prop[iter].perm);
        log_trace("* Available                 : %d", prop[iter].available);
        log_trace("* Type                      : %d", prop[iter].prop_type);
        log_trace("* Units                     : %s", prop[iter].units);
        log_trace("* Sysfs                     : %s", prop[iter].sysfname);
        if (prop[iter].dep_prop) {
            int cur_idx = prop[iter].dep_prop->curr_idx;
            int lmt_idx = prop[iter].dep_prop->lmt_idx;
            log_trace("* Current value Index        : %d", cur_idx);
            log_trace("* Current value Name         : %s", prop[cur_idx].name);
            log_trace("* Limit value Index          : %d", lmt_idx);
            log_trace("* Limit value Name           : %s", prop[lmt_idx].name);
            log_trace("* Alert Condition           :  0x%x",
                      prop[iter].dep_prop->cond);
        }
        log_trace(
            "********************************************************************************");
    }
}

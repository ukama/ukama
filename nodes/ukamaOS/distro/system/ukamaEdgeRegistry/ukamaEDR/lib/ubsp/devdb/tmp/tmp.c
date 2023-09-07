/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devdb/tmp/tmp.h"

#include "headers/ubsp/devices.h"
#include "headers/errorcode.h"
#include "inc/devicefxn.h"
#include "headers/utils/log.h"
#include "devdb/tmp/adt7481.h"
#include "devdb/tmp/se98.h"
#include "devdb/tmp/tmp464.h"

#include <string.h>

static ListInfo tmpdb;
static int tmpdbflag = 0;

const DevFxnTable tmp464_fxn_table = { .init = tmp464_init,
                                       .registration = tmp464_registration,
                                       .read_prop_count =
                                           tmp464_read_prop_count,
                                       .read_prop = tmp464_read_properties,
                                       .configure = tmp464_configure,
                                       .read = tmp464_read,
                                       .write = tmp464_write,
                                       .enable = tmp464_enable,
                                       .disable = tmp464_disable,
                                       .register_cb = tmp464_reg_cb,
                                       .dregister_cb = tmp464_dreg_cb,
                                       .enable_irq = tmp464_enable_irq,
                                       .disable_irq = tmp464_disable_irq,
                                       .confirm_irq = tmp464_confirm_irq,
                                       .irq_type = tmp464_get_irq_type };

const DevFxnTable se98_fxn_table = { .init = se98_init,
                                     .registration = se98_registration,
                                     .read_prop_count = se98_read_prop_count,
                                     .read_prop = se98_read_properties,
                                     .configure = se98_configure,
                                     .read = se98_read,
                                     .write = se98_write,
                                     .enable = se98_enable,
                                     .disable = se98_disable,
                                     .register_cb = se98_reg_cb,
                                     .dregister_cb = se98_dreg_cb,
                                     .enable_irq = se98_enable_irq,
                                     .disable_irq = se98_disable_irq,
                                     .confirm_irq = se98_confirm_irq,
                                     .irq_type = se98_get_irq_type };

const DevFxnTable adt7481_fxn_table = { .init = adt7481_init,
                                        .registration = adt7481_registration,
                                        .read_prop_count =
                                            adt7481_read_prop_count,
                                        .read_prop = adt7481_read_properties,
                                        .configure = adt7481_configure,
                                        .read = adt7481_read,
                                        .write = adt7481_write,
                                        .enable = adt7481_enable,
                                        .disable = adt7481_disable,
                                        .register_cb = adt7481_reg_cb,
                                        .dregister_cb = adt7481_dreg_cb,
                                        .enable_irq = adt7481_enable_irq,
                                        .disable_irq = adt7481_disable_irq,
                                        .confirm_irq = adt7481_confirm_irq,
                                        .irq_type = adt7481_get_irq_type };

DevFxnMap tmp_dev_map[MAX_TEMP_SENSOR_TYPE] = {
    { .name = "TMP464", .fxn_table = &tmp464_fxn_table },
    { .name = "ADT7481", .fxn_table = &adt7481_fxn_table },
    { .name = "SE98", .fxn_table = &se98_fxn_table }
};

const DevFxnTable *get_dev_tmp_fxn_tbl(char *name) {
    const DevFxnTable *fxn_tbl = NULL;
    for (uint8_t iter = 0; iter < MAX_TEMP_SENSOR_TYPE; iter++) {
        if (!strcmp(name, tmp_dev_map[iter].name)) {
            fxn_tbl = tmp_dev_map[iter].fxn_table;
            break;
        }
    }
    return fxn_tbl;
}

static void free_tmp_dev(void *ip) {
    ListNode *node = (ListNode *)ip;
    if (node) {
        if (node->data) {
            free(node->data);
        }
        free(node);
    }
}

int compare_tmp_dev(void *ipt, void *sd) {
    Device *ip = (Device *)ipt;
    Device *op = (Device *)sd;
    int ret = 0;
    /* If module if  and device name, disc, type matches it means devices is same.*/
    if (!strcmp(ip->obj.mod_UUID, op->obj.mod_UUID) &&
        !strcmp(ip->obj.name, op->obj.name) &&
        !strcmp(ip->obj.disc, op->obj.disc) && (ip->obj.type == op->obj.type)) {
        ret = 1;
    }
    return ret;
}

ListInfo *get_dev_tmp_db() {
    if (tmpdbflag == 0) {
        list_new(&tmpdb, sizeof(Device), free_tmp_dev, compare_tmp_dev, NULL);
        tmpdbflag = 1;
        log_trace("TMP:: TEMP DB initialized.");
    }
    return &tmpdb;
}

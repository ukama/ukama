/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devdb/pwr/pwr.h"

#include "headers/errorcode.h"
#include "headers/ubsp/devices.h"
#include "inc/devicefxn.h"
#include "headers/utils/log.h"
#include "devdb/pwr/ina226.h"

#include <string.h>

static ListInfo pwrdb;
static int pwrdbflag = 0;

const DevFxnTable ina226_fxn_table = { .init = ina226_init,
                                       .registration = ina226_registration,
                                       .read_prop_count =
                                           ina226_read_prop_count,
                                       .read_prop = ina226_read_properties,
                                       .configure = ina226_configure,
                                       .read = ina226_read,
                                       .write = ina226_write,
                                       .enable = ina226_enable,
                                       .disable = ina226_disable,
                                       .register_cb = ina226_reg_cb,
                                       .dregister_cb = ina226_dreg_cb,
                                       .enable_irq = ina226_enable_irq,
                                       .disable_irq = ina226_disable_irq,
                                       .confirm_irq = ina226_confirm_irq,
                                       .irq_type = ina226_get_irq_type };

DevFxnMap pwr_dev_map[MAX_PWR_SENSOR_TYPE] = {
    { .name = "INA226", .fxn_table = &ina226_fxn_table }
};

const DevFxnTable *get_dev_pwr_fxn_tbl(char *name) {
    const DevFxnTable *fxn_tbl = NULL;
    for (uint8_t iter = 0; iter < MAX_PWR_SENSOR_TYPE; iter++) {
        if (!strcmp(name, pwr_dev_map[iter].name)) {
            fxn_tbl = pwr_dev_map[iter].fxn_table;
            break;
        }
    }
    return fxn_tbl;
}

static void free_pwr_dev(void *ip) {
    ListNode *node = (ListNode *)ip;
    if (node) {
        if (node->data) {
            free(node->data);
        }
        free(node);
    }
}

int compare_pwr_dev(void *ipt, void *sd) {
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

ListInfo *get_dev_pwr_db() {
    /* Initialize DB for the first time we try to access it.*/
    if (pwrdbflag == 0) {
        list_new(&pwrdb, sizeof(Device), free_pwr_dev, compare_pwr_dev, NULL);
        pwrdbflag = 1;
        log_trace("PWR:: PWR DB initialized.");
    }
    return &pwrdb;
}

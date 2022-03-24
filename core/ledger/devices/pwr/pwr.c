/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devices/pwr.h"

#include "device.h"
#include "device_ops.h"
#include "errorcode.h"
#include "devices/bsp_ina226.h"

#include "usys_list.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

static ListInfo pwrLdgr;
static int pwrLdgrflag = 0;

const DevOps ina226Ops = { .init = bsp_ina226_init,
                           .registration = bsp_ina226_registration,
                           .readPropCount = bsp_ina226_read_prop_count,
                           .readProp = bsp_ina226_read_properties,
                           .configure = bsp_ina226_configure,
                           .read = bsp_ina226_read,
                           .write = bsp_ina226_write,
                           .enable = bsp_ina226_enable,
                           .disable = bsp_ina226_disable,
                           .registerCb = bsp_ina226_reg_cb,
                           .dregisterCb = bsp_ina226_dreg_cb,
                           .enableIrq = bsp_ina226_enable_irq,
                           .disableIrq = bsp_ina226_disable_irq,
                           .confirmIrq = bsp_ina226_confirm_irq,
                           .irqType = bsp_ina226_get_irq_type };

DevOpsMap pwrDevMap[MAX_PWR_SENSOR_TYPE] = { { .name = "INA226",
                                                 .opsTable = &ina226Ops } };

const DevOps *get_pwr_dev_ops(char *name) {
    const DevOps *opsTbl = NULL;
    for (uint8_t iter = 0; iter < MAX_PWR_SENSOR_TYPE; iter++) {
        if (!usys_strcasecmp(name, pwrDevMap[iter].name)) {
            opsTbl = pwrDevMap[iter].opsTable;
            break;
        }
    }
    return opsTbl;
}

static void free_pwr_dev(void *ip) {
    ListNode *node = (ListNode *)ip;
    if (node) {
        if (node->data) {
            usys_free(node->data);
        }
        usys_free(node);
    }
}

int compare_pwr_dev(void *ipt, void *sd) {
    Device *ip = (Device *)ipt;
    Device *op = (Device *)sd;
    int ret = 0;
    /* If module if  and device name, disc, type matches it means devices is same.*/
    if (!usys_strcasecmp(ip->obj.modUuid, op->obj.modUuid) &&
        !usys_strcasecmp(ip->obj.name, op->obj.name) &&
        !usys_strcasecmp(ip->obj.desc, op->obj.desc) &&
        (ip->obj.type == op->obj.type)) {
        ret = 1;
    }
    return ret;
}

ListInfo *get_pwr_dev_ldgr() {
    /* Initialize ledger for the first time we try to access it.*/
    if (pwrLdgrflag == 0) {
        usys_list_new(&pwrLdgr, sizeof(Device), free_pwr_dev, compare_pwr_dev,
                      NULL);
        pwrLdgrflag = 1;
        usys_log_trace("PWR:: PWR ledger initialized.");
    }
    return &pwrLdgr;
}
